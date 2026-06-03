<template>
  <div class="container recruit-admin-page">
    <div class="page-title">岗位管理</div>
    <section class="stat-grid">
      <article v-for="item in statCards" :key="item.label" class="stat-card">
        <span :class="['stat-icon', item.color]"><component :is="item.icon" /></span>
        <div>
          <p>{{ item.label }}</p>
          <strong>{{ item.value }}</strong>
        </div>
      </article>
    </section>

    <section class="panel-card">
      <a-form :model="queryForm" layout="inline" class="query-form">
        <a-form-item label="关键词">
          <a-input v-model="queryForm.keyword" allow-clear placeholder="岗位名称/企业名称" @press-enter="fetchData">
            <template #suffix><icon-search /></template>
          </a-input>
        </a-form-item>
        <a-form-item label="岗位类型">
          <a-select v-model="queryForm.jobType" allow-clear placeholder="请选择" style="width: 150px">
            <a-option v-for="item in jobTypeOptions" :key="item" :value="item">{{ item }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="所属行业">
          <a-select v-model="queryForm.industry" allow-clear placeholder="请选择" style="width: 160px">
            <a-option v-for="item in industryOptions" :key="item" :value="item">{{ item }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="工作地点">
          <a-input v-model="queryForm.location" allow-clear placeholder="城市/区域" />
        </a-form-item>
        <a-form-item label="学历要求">
          <a-select v-model="queryForm.education" allow-clear placeholder="请选择" style="width: 150px">
            <a-option v-for="item in educationOptions" :key="item" :value="item">{{ item }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="岗位状态">
          <a-select v-model="queryForm.status" allow-clear placeholder="请选择" style="width: 150px">
            <a-option v-for="item in statusOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
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
      <div class="table-toolbar">
        <a-button type="primary" @click="openCreate">新增岗位</a-button>
      </div>
      <a-table :columns="columns" :data="tableData" :pagination="pagination" row-key="id" @page-change="handlePageChange">
        <template #jobName="{ record }">
          <a-button type="text" @click="openDetail(record)">{{ record.jobName }}</a-button>
        </template>
        <template #status="{ record }">
          <a-tag :color="statusColor[displayStatus(record)]">{{ statusText[displayStatus(record)] || displayStatus(record) }}</a-tag>
        </template>
        <template #operations="{ record }">
          <a-space>
            <a-button size="small" @click="openDetail(record)">查看</a-button>
            <a-button size="small" @click="openEdit(record)">编辑</a-button>
            <a-button v-if="record.status !== 'published'" size="small" status="success" @click="publishJob(record)">上架</a-button>
            <a-button v-else size="small" status="warning" @click="offlineJob(record)">下架</a-button>
            <a-button size="small" status="danger" @click="removeJob(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </section>

    <a-drawer v-model:visible="detailVisible" width="620px" title="岗位详情">
      <a-descriptions :data="detailRows" :column="1" bordered />
      <a-divider>审核记录</a-divider>
      <a-list :data="detailReviews" :bordered="false">
        <template #item="{ item }">
          <a-list-item>{{ item.status }} - {{ item.reason || item.opinion || '无意见' }}</a-list-item>
        </template>
      </a-list>
    </a-drawer>

    <a-modal v-model:visible="editVisible" :title="formMode === 'create' ? '新增岗位' : '编辑岗位'" width="900px" @before-ok="saveJob">
      <a-form :model="formModel" layout="vertical">
        <a-row :gutter="16">
          <a-col v-for="field in fields" :key="field.key" :span="field.span || 8">
            <a-form-item :label="field.label">
              <a-select v-if="field.select === 'company'" v-model="formModel.companyId" allow-clear placeholder="请选择企业" @change="handleCompanyPick">
                <a-option v-for="item in companyOptions" :key="item.id" :value="item.id">{{ item.companyName }}</a-option>
              </a-select>
              <component v-else :is="field.textarea ? 'a-textarea' : field.number ? 'a-input-number' : 'a-input'" v-model="formModel[field.key]" allow-clear style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { Message, Modal } from '@arco-design/web-vue';
import { IconClockCircle, IconFile, IconMinus, IconSend } from '@arco-design/web-vue/es/icon';
import { computed, onMounted, reactive, ref } from 'vue';
import { useRoute } from 'vue-router';
import { createRecruitJob, deleteRecruitJob, getRecruitCompanies, getRecruitJob, getRecruitJobs, getRecruitStats, offlineRecruitJob, publishRecruitJob, updateRecruitJob } from '@/api/edu/recruit';

const route = useRoute();
const tableData = ref([]);
const stats = ref({});
const detailVisible = ref(false);
const editVisible = ref(false);
const formMode = ref('edit');
const detailRows = ref([]);
const detailReviews = ref([]);
const companyOptions = ref([]);
const queryForm = reactive({ keyword: '', jobType: '', industry: '', location: '', education: '', salaryRange: '', status: '', companyId: 0, pageIndex: 1, pageSize: 10 });
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const formModel = reactive({});
const jobTypeOptions = ['全职', '实习', '兼职', '校招', '社招'];
const industryOptions = ['教育', '特殊教育', '融合教育', '康复机构', '教育科研', '医疗/康复'];
const educationOptions = ['不限', '大专及以上', '本科及以上', '硕士及以上'];
const statusOptions = [
  { label: '待审核', value: 'pending' },
  { label: '已发布', value: 'published' },
  { label: '已下架', value: 'offline' },
  { label: '已过期', value: 'expired' },
  { label: '审核驳回', value: 'rejected' }
];
const statusText = { pending: '待审核', published: '已发布', offline: '已下架', expired: '已过期', rejected: '审核驳回' };
const statusColor = { pending: 'blue', published: 'green', offline: 'gray', expired: 'orange', rejected: 'red' };
const statCards = computed(() => [
  { label: '岗位总数', value: stats.value.jobs || 0, icon: IconFile, color: 'blue' },
  { label: '已发布', value: stats.value.publishedJobs || 0, icon: IconSend, color: 'green' },
  { label: '已过期', value: stats.value.expiredJobs || 0, icon: IconClockCircle, color: 'orange' },
  { label: '已下架', value: stats.value.offlineJobs || 0, icon: IconMinus, color: 'purple' }
]);
const columns = [
  { title: '岗位名称', slotName: 'jobName', width: 150, fixed: 'left' },
  { title: '所属企业', dataIndex: 'companyName', ellipsis: true, tooltip: true },
  { title: '岗位类型', dataIndex: 'jobType', width: 90 },
  { title: '招聘人数', dataIndex: 'headcount', width: 90 },
  { title: '工作地点', dataIndex: 'location', width: 130 },
  { title: '薪资范围', dataIndex: 'salaryRange', width: 120 },
  { title: '学历要求', dataIndex: 'education', width: 110 },
  { title: '截止时间', dataIndex: 'deadline', width: 120 },
  { title: '岗位状态', slotName: 'status', width: 100 },
  { title: '联系人', dataIndex: 'contactName', width: 100 },
  { title: '发布时间', dataIndex: 'publishTime', width: 150 },
  { title: '操作', slotName: 'operations', width: 240, fixed: 'right' }
];
const fields = [
  { key: 'jobName', label: '岗位名称' },
  { key: 'companyId', label: '所属企业', select: 'company' },
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

function defaultJobForm() {
  return {
    companyId: undefined,
    companyName: '',
    jobName: '',
    jobType: '全职',
    headcount: 1,
    location: '',
    salaryRange: '面议',
    education: '本科及以上',
    majorRequirement: '特殊教育、教育学、心理学、康复治疗等相关专业',
    majorDirection: '特殊教育',
    responsibilities: '',
    requirements: '',
    workTime: '周一至周五 8:30-17:30',
    benefits: '五险一金、专业培训、带薪年假、节日福利',
    contactName: '',
    contactTitle: '招聘负责人',
    contactPhone: '',
    contactEmail: '',
    externalLink: '',
    status: 'published'
  };
}

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

async function fetchCompanyOptions() {
  const res = await getRecruitCompanies({ pageIndex: 1, pageSize: 1000, status: 'normal' });
  companyOptions.value = res.data?.list || [];
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
  Object.assign(formModel, defaultJobForm(), record);
  formMode.value = 'edit';
  editVisible.value = true;
}

function openCreate() {
  Object.keys(formModel).forEach((key) => delete formModel[key]);
  Object.assign(formModel, defaultJobForm());
  formMode.value = 'create';
  editVisible.value = true;
}

async function saveJob() {
  if (formMode.value === 'create') {
    await createRecruitJob(formModel);
  } else {
    await updateRecruitJob(formModel.id, formModel);
  }
  Message.success('保存成功');
  await Promise.all([fetchData(), fetchStats()]);
  return true;
}

function handleCompanyPick(companyId) {
  const company = companyOptions.value.find((item) => item.id === companyId);
  if (!company) return;
  formModel.companyName = company.companyName;
  formModel.industry = company.industry;
  if (!formModel.location) formModel.location = company.region;
  if (!formModel.contactName) formModel.contactName = company.contactName;
  if (!formModel.contactPhone) formModel.contactPhone = company.contactPhone;
  if (!formModel.contactEmail) formModel.contactEmail = company.contactEmail;
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
  if (route.query.companyId) queryForm.companyId = Number(route.query.companyId);
  fetchStats();
  fetchCompanyOptions();
  fetchData();
});
</script>

<style scoped>
@import '../shared.scss';
</style>
