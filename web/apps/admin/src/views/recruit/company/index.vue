<template>
  <div class="container recruit-admin-page">
    <div class="page-title">企业管理</div>
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
          <a-input v-model="queryForm.keyword" allow-clear placeholder="企业名称/联系人" @press-enter="fetchData">
            <template #suffix><icon-search /></template>
          </a-input>
        </a-form-item>
        <a-form-item label="企业性质">
          <a-select v-model="queryForm.companyNature" allow-clear placeholder="请选择" style="width: 170px">
            <a-option v-for="item in natureOptions" :key="item" :value="item">{{ item }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="所属行业">
          <a-select v-model="queryForm.industry" allow-clear placeholder="请选择" style="width: 160px">
            <a-option v-for="item in industryOptions" :key="item" :value="item">{{ item }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="企业规模">
          <a-select v-model="queryForm.companySize" allow-clear placeholder="请选择" style="width: 150px">
            <a-option v-for="item in sizeOptions" :key="item" :value="item">{{ item }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="所在地区">
          <a-input v-model="queryForm.region" allow-clear placeholder="省市区" />
        </a-form-item>
        <a-form-item label="入驻状态">
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
        <a-button type="primary" @click="openCreate">新增企业</a-button>
      </div>
      <a-table :columns="columns" :data="tableData" :pagination="pagination" row-key="id" @page-change="handlePageChange">
        <template #status="{ record }">
          <a-tag :color="statusColor[record.status]">{{ statusText[record.status] || record.status }}</a-tag>
        </template>
        <template #operations="{ record }">
          <a-space>
            <a-button size="small" @click="openDetail(record)">查看</a-button>
            <a-button size="small" @click="openEdit(record)">编辑</a-button>
            <a-button size="small" @click="viewJobs(record)">查看岗位</a-button>
            <a-button v-if="record.status === 'disabled'" size="small" status="success" @click="toggleCompany(record, 'enable')">启用</a-button>
            <a-button v-else size="small" status="danger" @click="toggleCompany(record, 'disable')">禁用</a-button>
          </a-space>
        </template>
      </a-table>
    </section>

    <a-drawer v-model:visible="detailVisible" width="620px" title="企业详情">
      <a-descriptions :data="detailRows" :column="1" bordered />
      <a-divider>已发布岗位</a-divider>
      <a-list :data="detailJobs" :bordered="false">
        <template #item="{ item }">
          <a-list-item>{{ item.jobName }} - {{ statusText[item.status] || item.status }}</a-list-item>
        </template>
      </a-list>
      <a-divider>审核记录</a-divider>
      <a-list :data="detailReviews" :bordered="false">
        <template #item="{ item }">
          <a-list-item>{{ item.status }} - {{ item.reason || item.opinion || '无意见' }}</a-list-item>
        </template>
      </a-list>
    </a-drawer>

    <a-modal v-model:visible="editVisible" :title="formMode === 'create' ? '新增企业' : '编辑企业'" width="860px" @before-ok="saveCompany">
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
import { IconCheckCircle, IconFile, IconHome, IconStop } from '@arco-design/web-vue/es/icon';
import { computed, onMounted, reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import { createRecruitCompany, disableRecruitCompany, enableRecruitCompany, getRecruitCompany, getRecruitCompanies, getRecruitStats, updateRecruitCompany } from '@/api/edu/recruit';

const router = useRouter();
const tableData = ref([]);
const stats = ref({});
const detailVisible = ref(false);
const editVisible = ref(false);
const formMode = ref('edit');
const detailRows = ref([]);
const detailJobs = ref([]);
const detailReviews = ref([]);
const queryForm = reactive({ keyword: '', companyNature: '', industry: '', companySize: '', region: '', status: '', pageIndex: 1, pageSize: 10 });
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const formModel = reactive({});
const natureOptions = ['事业单位', '民办非企业单位', '民办非企业', '社会组织', '教育科技'];
const industryOptions = ['教育', '特殊教育', '康复机构', '融合教育', '教育科研', '医疗/康复'];
const sizeOptions = ['50人以下', '50-99人', '100-499人', '500人以上'];
const statusOptions = [
  { label: '正常', value: 'normal' },
  { label: '待完善', value: 'incomplete' },
  { label: '已禁用', value: 'disabled' },
  { label: '已驳回', value: 'rejected' },
  { label: '待审核', value: 'pending' }
];
const statusText = { normal: '正常', incomplete: '待完善', disabled: '已禁用', rejected: '已驳回', pending: '待审核', published: '已发布', offline: '已下架' };
const statusColor = { normal: 'green', incomplete: 'orange', disabled: 'red', rejected: 'red', pending: 'blue', published: 'green', offline: 'gray' };
const statCards = computed(() => [
  { label: '入驻企业总数', value: stats.value.companies || 0, icon: IconHome, color: 'blue' },
  { label: '正常企业', value: stats.value.normalCompanies || 0, icon: IconCheckCircle, color: 'green' },
  { label: '待完善', value: stats.value.incompleteCompanies || 0, icon: IconFile, color: 'orange' },
  { label: '已禁用', value: stats.value.disabledCompanies || 0, icon: IconStop, color: 'red' }
]);
const columns = [
  { title: '企业名称', dataIndex: 'companyName', ellipsis: true, tooltip: true },
  { title: '企业性质', dataIndex: 'companyNature', width: 140 },
  { title: '所属行业', dataIndex: 'industry', width: 120 },
  { title: '企业规模', dataIndex: 'companySize', width: 120 },
  { title: '所在地区', dataIndex: 'region', width: 150 },
  { title: '联系人', dataIndex: 'contactName', width: 100 },
  { title: '联系电话', dataIndex: 'contactPhone', width: 140 },
  { title: '入驻状态', slotName: 'status', width: 100 },
  { title: '岗位数量', dataIndex: 'jobCount', width: 100 },
  { title: '入驻时间', dataIndex: 'createdAt', width: 180 },
  { title: '操作', slotName: 'operations', width: 250, fixed: 'right' }
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

function defaultCompanyForm() {
  return {
    companyName: '',
    creditCode: '',
    companyNature: '民办非企业单位',
    industry: '特殊教育',
    companySize: '50-99人',
    region: '',
    address: '',
    website: '',
    logoUrl: '',
    contactName: '',
    contactTitle: '人事负责人',
    contactPhone: '',
    contactEmail: '',
    intro: '',
    mainBusiness: '',
    talentNeeds: '',
    cooperation: '',
    status: 'normal'
  };
}

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
  Object.assign(formModel, defaultCompanyForm(), record);
  formMode.value = 'edit';
  editVisible.value = true;
}

function openCreate() {
  Object.keys(formModel).forEach((key) => delete formModel[key]);
  Object.assign(formModel, defaultCompanyForm());
  formMode.value = 'create';
  editVisible.value = true;
}

async function saveCompany() {
  if (formMode.value === 'create') {
    await createRecruitCompany(formModel);
  } else {
    await updateRecruitCompany(formModel.id, formModel);
  }
  Message.success('保存成功');
  await Promise.all([fetchData(), fetchStats()]);
  return true;
}

async function toggleCompany(record, action) {
  if (action === 'enable') await enableRecruitCompany(record.id);
  else await disableRecruitCompany(record.id);
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
@import '../shared.scss';
</style>
