<template>
  <div class="container recruit-page">
    <div class="stat-grid">
      <a-card v-for="item in reviewStats" :key="item.label" :bordered="false" class="stat-card">
        <span>{{ item.label }}</span>
        <strong>{{ item.value }}</strong>
      </a-card>
    </div>

    <a-card :bordered="false" class="cardStyle">
      <template #title>信息审核</template>
      <a-tabs v-model:active-key="activeTab" @change="fetchData">
        <a-tab-pane key="company" title="企业入驻审核" />
        <a-tab-pane key="job" title="岗位发布审核" />
      </a-tabs>
      <a-form :model="queryForm" layout="inline">
        <a-form-item label="关键词">
          <a-input v-model="queryForm.keyword" allow-clear placeholder="企业/岗位/联系人" @press-enter="fetchData" />
        </a-form-item>
        <a-form-item label="审核状态">
          <a-select v-model="queryForm.status" allow-clear placeholder="全部状态" style="width: 140px" @change="fetchData">
            <a-option v-for="item in reviewOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="fetchData">查询</a-button>
            <a-button @click="resetQuery">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>

    <a-card :bordered="false" class="cardStyle">
      <a-table :columns="columns" :data="tableData" :pagination="pagination" row-key="id" @page-change="handlePageChange">
        <template #reviewStatus="{ record }">
          <a-tag :color="reviewColor[record.reviewStatus]">{{ reviewText[record.reviewStatus] || record.reviewStatus }}</a-tag>
        </template>
        <template #operations="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openDetail(record)">查看详情</a-button>
            <a-button v-if="record.reviewStatus === 'pending'" type="text" status="success" size="small" @click="openApprove(record)">通过</a-button>
            <a-button v-if="record.reviewStatus === 'pending'" type="text" status="danger" size="small" @click="openReject(record)">驳回</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-drawer v-model:visible="detailVisible" width="520px" title="审核详情">
      <a-descriptions :data="detailRows" :column="1" bordered />
      <a-divider>审核记录</a-divider>
      <a-list :data="detailReviews" :bordered="false">
        <template #item="{ item }">
          <a-list-item>{{ reviewText[item.status] || item.status }} - {{ item.reason || item.opinion || '无意见' }}</a-list-item>
        </template>
      </a-list>
    </a-drawer>

    <a-modal v-model:visible="reviewVisible" :title="reviewMode === 'approve' ? '审核通过' : '审核驳回'" @before-ok="submitReview">
      <a-form :model="reviewForm" layout="vertical">
        <a-form-item v-if="reviewMode === 'reject'" label="驳回原因" required>
          <a-select v-model="reviewForm.reason" placeholder="请选择驳回原因">
            <a-option v-for="item in rejectReasons" :key="item" :value="item">{{ item }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="审核意见">
          <a-textarea v-model="reviewForm.opinion" :auto-size="{ minRows: 3, maxRows: 5 }" placeholder="可填写补充说明" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { Message } from '@arco-design/web-vue';
import { computed, onMounted, reactive, ref } from 'vue';
import { approveRecruitReview, getRecruitCompany, getRecruitJob, getRecruitStats, getReviewCompanies, getReviewJobs, rejectRecruitReview } from '@/api/edu/recruit';

const activeTab = ref('company');
const tableData = ref([]);
const stats = ref({});
const detailVisible = ref(false);
const reviewVisible = ref(false);
const detailRows = ref([]);
const detailReviews = ref([]);
const currentRecord = ref(null);
const reviewMode = ref('approve');
const queryForm = reactive({ keyword: '', status: '', pageIndex: 1, pageSize: 10 });
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const reviewForm = reactive({ reason: '', opinion: '' });
const reviewOptions = [
  { label: '待审核', value: 'pending' },
  { label: '已通过', value: 'approved' },
  { label: '已驳回', value: 'rejected' },
  { label: '已撤回', value: 'withdrawn' }
];
const rejectReasons = ['信息不完整', '联系人信息异常', '材料不清晰', '岗位描述不清晰', '薪资范围缺失', '企业状态异常', '其他原因'];
const reviewText = { pending: '待审核', approved: '已通过', rejected: '已驳回', withdrawn: '已撤回' };
const reviewColor = { pending: 'orange', approved: 'green', rejected: 'red', withdrawn: 'gray' };

const reviewStats = computed(() => [
  { label: '待审核企业', value: stats.value.pendingCompanies || 0 },
  { label: '待审核岗位', value: stats.value.pendingJobs || 0 },
  { label: '今日新增申请', value: stats.value.todayApplications || 0 },
  { label: '本月通过率', value: passRate.value }
]);
const passRate = computed(() => {
  const total = (stats.value.pendingCompanies || 0) + (stats.value.pendingJobs || 0) + (stats.value.normalCompanies || 0) + (stats.value.publishedJobs || 0);
  const passed = (stats.value.normalCompanies || 0) + (stats.value.publishedJobs || 0);
  return total ? `${Math.round((passed / total) * 100)}%` : '0%';
});
const columns = computed(() => activeTab.value === 'company'
  ? [
      { title: '企业名称', dataIndex: 'companyName', ellipsis: true, tooltip: true },
      { title: '统一社会信用代码', dataIndex: 'creditCode', width: 170 },
      { title: '企业性质', dataIndex: 'companyNature', width: 120 },
      { title: '所属行业', dataIndex: 'industry', width: 120 },
      { title: '联系人', dataIndex: 'contactName', width: 100 },
      { title: '联系电话', dataIndex: 'contactPhone', width: 130 },
      { title: '提交时间', dataIndex: 'createdAt', width: 180 },
      { title: '审核状态', slotName: 'reviewStatus', width: 100 },
      { title: '操作', slotName: 'operations', width: 180, fixed: 'right' }
    ]
  : [
      { title: '岗位名称', dataIndex: 'jobName', ellipsis: true, tooltip: true },
      { title: '所属企业', dataIndex: 'companyName', ellipsis: true, tooltip: true },
      { title: '岗位类型', dataIndex: 'jobType', width: 100 },
      { title: '招聘人数', dataIndex: 'headcount', width: 100 },
      { title: '工作地点', dataIndex: 'location', width: 120 },
      { title: '联系人', dataIndex: 'contactName', width: 100 },
      { title: '提交时间', dataIndex: 'createdAt', width: 180 },
      { title: '审核状态', slotName: 'reviewStatus', width: 100 },
      { title: '操作', slotName: 'operations', width: 180, fixed: 'right' }
    ]);

async function fetchStats() {
  const res = await getRecruitStats();
  stats.value = res.data || {};
}

async function fetchData() {
  const api = activeTab.value === 'company' ? getReviewCompanies : getReviewJobs;
  const res = await api(queryForm);
  const payload = res.data || {};
  tableData.value = payload.list || [];
  pagination.total = payload.count || 0;
  pagination.current = queryForm.pageIndex;
}

function resetQuery() {
  Object.assign(queryForm, { keyword: '', status: '', pageIndex: 1 });
  fetchData();
}

function handlePageChange(page) {
  queryForm.pageIndex = page;
  fetchData();
}

async function openDetail(record) {
  currentRecord.value = record;
  const res = activeTab.value === 'company' ? await getRecruitCompany(record.id) : await getRecruitJob(record.id);
  const data = res.data || {};
  const entity = activeTab.value === 'company' ? data.company : data.job;
  detailReviews.value = data.reviews || [];
  detailRows.value = Object.entries(entity || {}).filter(([key]) => !['deletedAt'].includes(key)).map(([label, value]) => ({ label, value: String(value ?? '') }));
  detailVisible.value = true;
}

function openApprove(record) {
  currentRecord.value = record;
  reviewMode.value = 'approve';
  Object.assign(reviewForm, { reason: '', opinion: '' });
  reviewVisible.value = true;
}

function openReject(record) {
  currentRecord.value = record;
  reviewMode.value = 'reject';
  Object.assign(reviewForm, { reason: '', opinion: '' });
  reviewVisible.value = true;
}

async function submitReview() {
  if (reviewMode.value === 'reject' && !reviewForm.reason) {
    Message.warning('请选择驳回原因');
    return false;
  }
  const payload = { targetType: activeTab.value, reason: reviewForm.reason, opinion: reviewForm.opinion };
  if (reviewMode.value === 'approve') {
    await approveRecruitReview(currentRecord.value.id, payload);
  } else {
    await rejectRecruitReview(currentRecord.value.id, payload);
  }
  Message.success('审核成功');
  await Promise.all([fetchData(), fetchStats()]);
  return true;
}

onMounted(() => {
  fetchStats();
  fetchData();
});
</script>

<style scoped>
.recruit-page {
  padding: 16px;
}

.stat-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
  margin-bottom: 16px;
}

.stat-card span {
  color: #6b778c;
}

.stat-card strong {
  display: block;
  margin-top: 8px;
  color: #172b4d;
  font-size: 28px;
}
</style>
