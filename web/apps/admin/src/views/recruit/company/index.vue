<template>
  <div class="container recruit-page">
    <div class="stat-grid">
      <a-card v-for="item in statCards" :key="item.label" :bordered="false" class="stat-card">
        <span>{{ item.label }}</span>
        <strong>{{ item.value }}</strong>
      </a-card>
    </div>

    <a-card :bordered="false" class="cardStyle">
      <template #title>企业管理</template>
      <a-form :model="queryForm" layout="inline">
        <a-form-item label="关键词">
          <a-input v-model="queryForm.keyword" allow-clear placeholder="企业名称 / 联系人" @press-enter="fetchData" />
        </a-form-item>
        <a-form-item label="企业性质">
          <a-input v-model="queryForm.companyNature" allow-clear placeholder="企业性质" />
        </a-form-item>
        <a-form-item label="行业">
          <a-input v-model="queryForm.industry" allow-clear placeholder="所属行业" />
        </a-form-item>
        <a-form-item label="规模">
          <a-input v-model="queryForm.companySize" allow-clear placeholder="企业规模" />
        </a-form-item>
        <a-form-item label="地区">
          <a-input v-model="queryForm.region" allow-clear placeholder="所在地区" />
        </a-form-item>
        <a-form-item label="入驻状态">
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
          <a-tag :color="statusColor[record.status]">{{ statusText[record.status] || record.status }}</a-tag>
        </template>
        <template #operations="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openDetail(record)">查看</a-button>
            <a-button type="text" size="small" @click="openEdit(record)">编辑</a-button>
            <a-button type="text" size="small" @click="viewJobs(record)">查看岗位</a-button>
            <a-button v-if="record.status === 'disabled'" type="text" status="success" size="small" @click="toggleCompany(record, 'enable')">启用</a-button>
            <a-button v-else type="text" status="danger" size="small" @click="toggleCompany(record, 'disable')">禁用</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-drawer v-model:visible="detailVisible" width="560px" title="企业详情">
      <a-descriptions :data="detailRows" :column="1" bordered />
      <a-divider>已发布岗位</a-divider>
      <a-list :data="detailJobs" :bordered="false">
        <template #item="{ item }">
          <a-list-item>{{ item.jobName }} - {{ item.status }}</a-list-item>
        </template>
      </a-list>
      <a-divider>审核记录</a-divider>
      <a-list :data="detailReviews" :bordered="false">
        <template #item="{ item }">
          <a-list-item>{{ item.status }} - {{ item.reason || item.opinion || '无意见' }}</a-list-item>
        </template>
      </a-list>
    </a-drawer>

    <a-modal v-model:visible="editVisible" title="编辑企业" width="860px" @before-ok="saveCompany">
      <a-form :model="formModel" layout="vertical">
        <a-row :gutter="16">
          <a-col v-for="field in fields" :key="field.key" :span="field.span || 8">
            <a-form-item :label="field.label">
              <component :is="field.textarea ? 'a-textarea' : 'a-input'" v-model="formModel[field.key]" allow-clear />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { Message } from '@arco-design/web-vue';
import { computed, onMounted, reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import { disableRecruitCompany, enableRecruitCompany, getRecruitCompany, getRecruitCompanies, getRecruitStats, updateRecruitCompany } from '@/api/edu/recruit';

const router = useRouter();
const tableData = ref([]);
const stats = ref({});
const detailVisible = ref(false);
const editVisible = ref(false);
const detailRows = ref([]);
const detailJobs = ref([]);
const detailReviews = ref([]);
const queryForm = reactive({ keyword: '', companyNature: '', industry: '', companySize: '', region: '', status: '', pageIndex: 1, pageSize: 10 });
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const formModel = reactive({});
const statusOptions = [
  { label: '正常', value: 'normal' },
  { label: '待完善', value: 'incomplete' },
  { label: '已禁用', value: 'disabled' },
  { label: '已驳回', value: 'rejected' }
];
const statusText = { normal: '正常', incomplete: '待完善', disabled: '已禁用', rejected: '已驳回', pending: '待审核' };
const statusColor = { normal: 'green', incomplete: 'orange', disabled: 'gray', rejected: 'red', pending: 'orange' };
const statCards = computed(() => [
  { label: '入驻企业总数', value: stats.value.companies || 0 },
  { label: '正常企业', value: stats.value.normalCompanies || 0 },
  { label: '待完善', value: stats.value.incompleteCompanies || 0 },
  { label: '已禁用', value: stats.value.disabledCompanies || 0 }
]);
const columns = [
  { title: '企业名称', dataIndex: 'companyName', ellipsis: true, tooltip: true },
  { title: '企业性质', dataIndex: 'companyNature', width: 120 },
  { title: '所属行业', dataIndex: 'industry', width: 120 },
  { title: '企业规模', dataIndex: 'companySize', width: 120 },
  { title: '所在地区', dataIndex: 'region', width: 120 },
  { title: '联系人', dataIndex: 'contactName', width: 100 },
  { title: '联系电话', dataIndex: 'contactPhone', width: 130 },
  { title: '入驻状态', slotName: 'status', width: 100 },
  { title: '岗位数量', dataIndex: 'jobCount', width: 100 },
  { title: '入驻时间', dataIndex: 'createdAt', width: 180 },
  { title: '操作', slotName: 'operations', width: 230, fixed: 'right' }
];
const fields = [
  { key: 'companyName', label: '企业名称' },
  { key: 'logoUrl', label: '企业Logo' },
  { key: 'companyNature', label: '企业性质' },
  { key: 'industry', label: '所属行业' },
  { key: 'companySize', label: '企业规模' },
  { key: 'region', label: '所在地区' },
  { key: 'address', label: '详细地址', span: 16 },
  { key: 'website', label: '企业官网链接' },
  { key: 'contactName', label: '联系人' },
  { key: 'contactTitle', label: '联系人职务' },
  { key: 'contactPhone', label: '联系电话' },
  { key: 'contactEmail', label: '联系邮箱' },
  { key: 'intro', label: '企业简介', textarea: true, span: 12 },
  { key: 'mainBusiness', label: '主营业务', textarea: true, span: 12 },
  { key: 'talentNeeds', label: '人才需求方向', span: 12 },
  { key: 'cooperation', label: '校企合作方向', textarea: true, span: 12 },
  { key: 'status', label: '企业状态' }
];

async function fetchStats() {
  const res = await getRecruitStats();
  stats.value = res.data || {};
}

async function fetchData() {
  const res = await getRecruitCompanies(queryForm);
  const payload = res.data || {};
  tableData.value = payload.list || [];
  pagination.total = payload.count || 0;
  pagination.current = queryForm.pageIndex;
}

function resetQuery() {
  Object.assign(queryForm, { keyword: '', companyNature: '', industry: '', companySize: '', region: '', status: '', pageIndex: 1 });
  fetchData();
}

function handlePageChange(page) {
  queryForm.pageIndex = page;
  fetchData();
}

async function openDetail(record) {
  const res = await getRecruitCompany(record.id);
  const data = res.data || {};
  detailJobs.value = data.jobs || [];
  detailReviews.value = data.reviews || [];
  detailRows.value = Object.entries(data.company || {}).filter(([key]) => !['deletedAt'].includes(key)).map(([label, value]) => ({ label, value: String(value ?? '') }));
  detailVisible.value = true;
}

function openEdit(record) {
  Object.assign(formModel, record);
  editVisible.value = true;
}

async function saveCompany() {
  await updateRecruitCompany(formModel.id, formModel);
  Message.success('保存成功');
  await Promise.all([fetchData(), fetchStats()]);
  return true;
}

async function toggleCompany(record, action) {
  if (action === 'enable') {
    await enableRecruitCompany(record.id);
  } else {
    await disableRecruitCompany(record.id);
  }
  Message.success('操作成功');
  await Promise.all([fetchData(), fetchStats()]);
}

function viewJobs(record) {
  router.push(`/recruit/job?companyId=${record.id}`);
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
