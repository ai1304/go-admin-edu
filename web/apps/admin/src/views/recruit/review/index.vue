<template>
  <div class="container recruit-admin-page">
    <div class="page-title">信息审核</div>
    <section class="stat-grid">
      <article v-for="item in reviewStats" :key="item.label" class="stat-card">
        <span :class="['stat-icon', item.color]"><component :is="item.icon" /></span>
        <div>
          <p>{{ item.label }}</p>
          <strong>{{ item.value }}</strong>
        </div>
      </article>
    </section>

    <section class="panel-card">
      <a-tabs v-model:active-key="activeTab" @change="handleTabChange">
        <a-tab-pane key="company" title="企业入驻审核" />
        <a-tab-pane key="job" title="岗位发布审核" />
      </a-tabs>
      <a-form :model="queryForm" layout="inline" class="query-form">
        <a-form-item label="关键词">
          <a-input v-model="queryForm.keyword" allow-clear placeholder="请输入企业、岗位或联系人" @press-enter="fetchData" />
        </a-form-item>
        <a-form-item label="审核状态">
          <a-select v-model="queryForm.status" allow-clear placeholder="请选择" style="width: 150px" @change="fetchData">
            <a-option v-for="item in reviewOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="fetchData"><icon-search /> 查询</a-button>
        </a-form-item>
        <a-form-item>
          <a-button @click="resetQuery"><icon-loop /> 重置</a-button>
        </a-form-item>
      </a-form>
    </section>

    <section class="panel-card">
      <a-table :columns="columns" :data="tableData" :pagination="pagination" row-key="id" @page-change="handlePageChange">
        <template #reviewStatus="{ record }">
          <a-tag :color="reviewColor[record.reviewStatus]">{{ reviewText[record.reviewStatus] || record.reviewStatus }}</a-tag>
        </template>
        <template #operations="{ record }">
          <a-space>
            <a-button size="small" @click="openDetail(record)">查看详情</a-button>
            <a-button v-if="record.reviewStatus === 'pending'" size="small" status="success" @click="openApprove(record)">通过</a-button>
            <a-button v-if="record.reviewStatus === 'pending'" size="small" status="danger" @click="openReject(record)">驳回</a-button>
          </a-space>
        </template>
      </a-table>
    </section>

    <a-drawer v-model:visible="detailVisible" width="620px" title="审核详情">
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
import { IconBarChart, IconFile, IconHome, IconStorage } from '@arco-design/web-vue/es/icon';
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
const reviewColor = { pending: 'blue', approved: 'green', rejected: 'red', withdrawn: 'gray' };

const passRate = computed(() => {
  const total = stats.value.monthApplications || 0;
  const passed = stats.value.monthPassed || 0;
  return total ? `${Math.round((passed / total) * 100)}%` : '0%';
});
const reviewStats = computed(() => [
  { label: '待审核企业', value: stats.value.pendingCompanies || 0, icon: IconHome, color: 'blue' },
  { label: '待审核岗位', value: stats.value.pendingJobs || 0, icon: IconStorage, color: 'green' },
  { label: '今日新增申请', value: stats.value.todayApplications || 0, icon: IconFile, color: 'orange' },
  { label: '本月通过率', value: passRate.value, icon: IconBarChart, color: 'purple' }
]);
const columns = computed(() => activeTab.value === 'company'
  ? [
      { title: '企业名称', dataIndex: 'companyName', ellipsis: true, tooltip: true },
      { title: '统一社会信用代码', dataIndex: 'creditCode', width: 180 },
      { title: '企业性质', dataIndex: 'companyNature', width: 140 },
      { title: '所属行业', dataIndex: 'industry', width: 120 },
      { title: '联系人', dataIndex: 'contactName', width: 100 },
      { title: '联系电话', dataIndex: 'contactPhone', width: 140 },
      { title: '提交时间', dataIndex: 'createdAt', width: 180 },
      { title: '审核状态', slotName: 'reviewStatus', width: 110 },
      { title: '操作', slotName: 'operations', width: 220, fixed: 'right' }
    ]
  : [
      { title: '岗位名称', dataIndex: 'jobName', ellipsis: true, tooltip: true },
      { title: '所属企业', dataIndex: 'companyName', ellipsis: true, tooltip: true },
      { title: '岗位类型', dataIndex: 'jobType', width: 100 },
      { title: '招聘人数', dataIndex: 'headcount', width: 100 },
      { title: '工作地点', dataIndex: 'location', width: 130 },
      { title: '联系人', dataIndex: 'contactName', width: 100 },
      { title: '提交时间', dataIndex: 'createdAt', width: 180 },
      { title: '审核状态', slotName: 'reviewStatus', width: 110 },
      { title: '操作', slotName: 'operations', width: 220, fixed: 'right' }
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

function handleTabChange() {
  queryForm.pageIndex = 1;
  fetchData();
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
  if (reviewMode.value === 'approve') await approveRecruitReview(currentRecord.value.id, payload);
  else await rejectRecruitReview(currentRecord.value.id, payload);
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
.recruit-admin-page {
  padding: 20px;
  background: #f4f8ff;
}

.page-title {
  margin-bottom: 18px;
  color: #12264b;
  font-size: 26px;
  font-weight: 800;
}

.stat-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 18px;
  margin-bottom: 18px;
}

.stat-card,
.panel-card {
  background: #fff;
  border: 1px solid #dfe9f7;
  border-radius: 8px;
  box-shadow: 0 12px 32px rgba(20, 78, 148, 0.06);
}

.stat-card {
  display: flex;
  gap: 18px;
  align-items: center;
  padding: 22px;
}

.stat-icon {
  display: grid;
  place-items: center;
  width: 58px;
  height: 58px;
  font-size: 26px;
  border-radius: 50%;
}

.stat-icon.blue { color: #0969e8; background: #eaf3ff; }
.stat-icon.green { color: #19a86b; background: #e9f8f0; }
.stat-icon.orange { color: #ff8a1f; background: #fff2e5; }
.stat-icon.purple { color: #6d5dfc; background: #f0edff; }

.stat-card p {
  margin: 0 0 8px;
  color: #6b778c;
  font-weight: 700;
}

.stat-card strong {
  color: #12264b;
  font-size: 30px;
}

.panel-card {
  padding: 18px;
  margin-bottom: 18px;
}

.query-form {
  padding-top: 10px;
}

@media (max-width: 1000px) {
  .stat-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>
