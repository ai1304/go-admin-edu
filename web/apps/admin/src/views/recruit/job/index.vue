<template>
  <div class="container recruit-page">
    <div class="stat-grid">
      <a-card v-for="item in statCards" :key="item.label" :bordered="false" class="stat-card">
        <span>{{ item.label }}</span>
        <strong>{{ item.value }}</strong>
      </a-card>
    </div>

    <a-card :bordered="false" class="cardStyle">
      <template #title>岗位管理</template>
      <a-form :model="queryForm" layout="inline">
        <a-form-item label="关键词">
          <a-input v-model="queryForm.keyword" allow-clear placeholder="岗位名称 / 企业名称" @press-enter="fetchData" />
        </a-form-item>
        <a-form-item label="岗位类型">
          <a-input v-model="queryForm.jobType" allow-clear placeholder="岗位类型" />
        </a-form-item>
        <a-form-item label="所属行业">
          <a-input v-model="queryForm.industry" allow-clear placeholder="所属行业" />
        </a-form-item>
        <a-form-item label="工作地点">
          <a-input v-model="queryForm.location" allow-clear placeholder="工作地点" />
        </a-form-item>
        <a-form-item label="学历">
          <a-input v-model="queryForm.education" allow-clear placeholder="学历要求" />
        </a-form-item>
        <a-form-item label="岗位状态">
          <a-select v-model="queryForm.status" allow-clear placeholder="全部状态" style="width: 140px">
            <a-option v-for="item in statusOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
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
        <template #status="{ record }">
          <a-tag :color="statusColor[displayStatus(record)]">{{ statusText[displayStatus(record)] || displayStatus(record) }}</a-tag>
        </template>
        <template #operations="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openDetail(record)">查看</a-button>
            <a-button type="text" size="small" @click="openEdit(record)">编辑</a-button>
            <a-button v-if="record.status !== 'published'" type="text" status="success" size="small" @click="publishJob(record)">上架</a-button>
            <a-button v-else type="text" status="warning" size="small" @click="offlineJob(record)">下架</a-button>
            <a-button type="text" status="danger" size="small" @click="removeJob(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-drawer v-model:visible="detailVisible" width="560px" title="岗位详情">
      <a-descriptions :data="detailRows" :column="1" bordered />
      <a-divider>审核记录</a-divider>
      <a-list :data="detailReviews" :bordered="false">
        <template #item="{ item }">
          <a-list-item>{{ item.status }} - {{ item.reason || item.opinion || '无意见' }}</a-list-item>
        </template>
      </a-list>
    </a-drawer>

    <a-modal v-model:visible="editVisible" title="编辑岗位" width="900px" @before-ok="saveJob">
      <a-form :model="formModel" layout="vertical">
        <a-row :gutter="16">
          <a-col v-for="field in fields" :key="field.key" :span="field.span || 8">
            <a-form-item :label="field.label">
              <component :is="field.textarea ? 'a-textarea' : field.number ? 'a-input-number' : 'a-input'" v-model="formModel[field.key]" allow-clear style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { Message, Modal } from '@arco-design/web-vue';
import { computed, onMounted, reactive, ref } from 'vue';
import { useRoute } from 'vue-router';
import { deleteRecruitJob, getRecruitJob, getRecruitJobs, getRecruitStats, offlineRecruitJob, publishRecruitJob, updateRecruitJob } from '@/api/edu/recruit';

const route = useRoute();
const tableData = ref([]);
const stats = ref({});
const detailVisible = ref(false);
const editVisible = ref(false);
const detailRows = ref([]);
const detailReviews = ref([]);
const queryForm = reactive({ keyword: '', jobType: '', industry: '', location: '', education: '', salaryRange: '', status: '', companyId: 0, pageIndex: 1, pageSize: 10 });
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const formModel = reactive({});
const statusOptions = [
  { label: '待审核', value: 'pending' },
  { label: '已发布', value: 'published' },
  { label: '已下架', value: 'offline' },
  { label: '已过期', value: 'expired' },
  { label: '审核驳回', value: 'rejected' }
];
const statusText = { pending: '待审核', published: '已发布', offline: '已下架', expired: '已过期', rejected: '审核驳回' };
const statusColor = { pending: 'orange', published: 'green', offline: 'gray', expired: 'red', rejected: 'red' };
const statCards = computed(() => [
  { label: '岗位总数', value: stats.value.jobs || 0 },
  { label: '已发布', value: stats.value.publishedJobs || 0 },
  { label: '已过期', value: stats.value.expiredJobs || 0 },
  { label: '已下架', value: stats.value.offlineJobs || 0 }
]);
const columns = [
  { title: '岗位名称', dataIndex: 'jobName', ellipsis: true, tooltip: true },
  { title: '所属企业', dataIndex: 'companyName', ellipsis: true, tooltip: true },
  { title: '岗位类型', dataIndex: 'jobType', width: 100 },
  { title: '招聘人数', dataIndex: 'headcount', width: 100 },
  { title: '工作地点', dataIndex: 'location', width: 120 },
  { title: '薪资范围', dataIndex: 'salaryRange', width: 120 },
  { title: '学历要求', dataIndex: 'education', width: 100 },
  { title: '截止时间', dataIndex: 'deadline', width: 130 },
  { title: '岗位状态', slotName: 'status', width: 100 },
  { title: '联系人', dataIndex: 'contactName', width: 100 },
  { title: '发布时间', dataIndex: 'publishTime', width: 180 },
  { title: '操作', slotName: 'operations', width: 220, fixed: 'right' }
];
const fields = [
  { key: 'jobName', label: '岗位名称' },
  { key: 'companyName', label: '所属企业' },
  { key: 'jobType', label: '岗位类型' },
  { key: 'headcount', label: '招聘人数', number: true },
  { key: 'location', label: '工作地点' },
  { key: 'salaryRange', label: '薪资范围' },
  { key: 'education', label: '学历要求' },
  { key: 'majorRequirement', label: '专业要求' },
  { key: 'majorDirection', label: '专业方向' },
  { key: 'responsibilities', label: '岗位职责', textarea: true, span: 12 },
  { key: 'requirements', label: '任职要求', textarea: true, span: 12 },
  { key: 'workTime', label: '工作时间', textarea: true, span: 12 },
  { key: 'benefits', label: '福利待遇', textarea: true, span: 12 },
  { key: 'contactName', label: '联系人' },
  { key: 'contactTitle', label: '联系人职务' },
  { key: 'contactPhone', label: '联系电话' },
  { key: 'contactEmail', label: '联系邮箱' },
  { key: 'externalLink', label: '岗位外部链接', span: 16 },
  { key: 'status', label: '岗位状态' }
];

function displayStatus(record) {
  if (record.deadline && record.deadline < new Date().toISOString().slice(0, 10)) return 'expired';
  return record.status;
}

async function fetchStats() {
  const res = await getRecruitStats();
  stats.value = res.data || {};
}

async function fetchData() {
  const res = await getRecruitJobs(queryForm);
  const payload = res.data || {};
  tableData.value = payload.list || [];
  pagination.total = payload.count || 0;
  pagination.current = queryForm.pageIndex;
}

function resetQuery() {
  Object.assign(queryForm, { keyword: '', jobType: '', industry: '', location: '', education: '', salaryRange: '', status: '', companyId: 0, pageIndex: 1 });
  fetchData();
}

function handlePageChange(page) {
  queryForm.pageIndex = page;
  fetchData();
}

async function openDetail(record) {
  const res = await getRecruitJob(record.id);
  const data = res.data || {};
  detailReviews.value = data.reviews || [];
  detailRows.value = Object.entries(data.job || {}).filter(([key]) => !['deletedAt'].includes(key)).map(([label, value]) => ({ label, value: String(value ?? '') }));
  detailVisible.value = true;
}

function openEdit(record) {
  Object.assign(formModel, record);
  editVisible.value = true;
}

async function saveJob() {
  await updateRecruitJob(formModel.id, formModel);
  Message.success('保存成功');
  await Promise.all([fetchData(), fetchStats()]);
  return true;
}

async function publishJob(record) {
  await publishRecruitJob(record.id);
  Message.success('上架成功');
  await Promise.all([fetchData(), fetchStats()]);
}

async function offlineJob(record) {
  await offlineRecruitJob(record.id);
  Message.success('下架成功');
  await Promise.all([fetchData(), fetchStats()]);
}

function removeJob(record) {
  Modal.warning({
    title: '确认删除岗位',
    content: `确定删除“${record.jobName}”吗？`,
    hideCancel: false,
    onOk: async () => {
      await deleteRecruitJob(record.id);
      Message.success('删除成功');
      await Promise.all([fetchData(), fetchStats()]);
    }
  });
}

onMounted(() => {
  if (route.query.companyId) {
    queryForm.companyId = Number(route.query.companyId);
  }
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
