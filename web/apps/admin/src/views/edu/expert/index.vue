<template>
  <div class="container">
    <a-card :bordered="false" class="cardStyle">
      <template #title>名师/专家管理</template>
      <a-form :model="queryForm" layout="inline">
        <a-form-item label="关键词">
          <a-input v-model="queryForm.keyword" allow-clear placeholder="专家姓名、职称、领域" />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button v-has="'edu:expert:query'" type="primary" @click="fetchData">查询</a-button>
            <a-button @click="resetQuery">重置</a-button>
            <a-button v-has="'edu:expert:add'" type="primary" status="success" @click="openCreate">新增专家</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>

    <a-card :bordered="false" class="cardStyle table-card">
      <a-table :columns="columns" :data="tableData" :pagination="pagination" row-key="id" @page-change="handlePageChange">
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '停用' }}</a-tag>
        </template>
        <template #operations="{ record }">
          <a-space>
            <a-button v-has="'edu:expert:edit'" type="text" size="small" @click="openEdit(record)">编辑</a-button>
            <a-button v-has="'edu:expert:resources'" type="text" size="small" @click="openResources(record)">资源</a-button>
            <a-button v-has="'edu:expert:remove'" type="text" status="danger" size="small" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="formModel.id ? '编辑专家' : '新增专家'" width="760px" @before-ok="handleSave">
      <a-form :model="formModel" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item field="name" label="姓名" required>
              <a-input v-model="formModel.name" placeholder="请输入专家姓名" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="title" label="职称">
              <a-input v-model="formModel.title" placeholder="请输入职称" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="organization" label="机构">
              <a-input v-model="formModel.organization" placeholder="请输入机构" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="status" label="状态">
              <a-select v-model="formModel.status">
                <a-option :value="1">启用</a-option>
                <a-option :value="0">停用</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="isRecommended" label="门户推荐">
              <a-switch v-model="formModel.isRecommended" :checked-value="1" :unchecked-value="0" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="sort" label="门户排序">
              <a-input-number v-model="formModel.sort" :min="0" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item field="specialties" label="擅长领域">
              <a-input v-model="formModel.specialties" placeholder="多个领域用逗号分隔" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item field="introduction" label="简介">
              <a-textarea v-model="formModel.introduction" :auto-size="{ minRows: 4, maxRows: 8 }" placeholder="请输入专家简介" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>

    <a-modal v-model:visible="resourceVisible" :title="`${currentExpert?.name || ''} 资源管理`" width="960px" :footer="false">
      <a-space direction="vertical" fill>
        <a-button v-has="'edu:expert:resources'" type="primary" status="success" @click="openResourceCreate">新增关联资源</a-button>
        <a-table :columns="resourceColumns" :data="resourceList" :pagination="false" row-key="id">
          <template #type="{ record }">
            <a-tag>{{ resourceTypeText[record.type] || record.type }}</a-tag>
          </template>
          <template #status="{ record }">
            <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '停用' }}</a-tag>
          </template>
          <template #operations="{ record }">
            <a-space>
              <a-button v-has="'edu:expert:resources'" type="text" size="small" @click="openResourceEdit(record)">编辑</a-button>
              <a-button v-has="'edu:expert:resources'" type="text" status="danger" size="small" @click="handleResourceDelete(record)">删除</a-button>
            </a-space>
          </template>
        </a-table>
      </a-space>
    </a-modal>

    <a-modal v-model:visible="resourceFormVisible" :title="resourceModel.id ? '编辑关联资源' : '新增关联资源'" width="680px" @before-ok="handleResourceSave">
      <a-form :model="resourceModel" layout="vertical">
        <a-form-item field="title" label="标题" required>
          <a-input v-model="resourceModel.title" placeholder="请输入门户展示标题" />
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item field="type" label="类型">
              <a-select v-model="resourceModel.type">
                <a-option v-for="item in resourceTypeOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="status" label="状态">
              <a-select v-model="resourceModel.status">
                <a-option :value="1">启用</a-option>
                <a-option :value="0">停用</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item field="resourceId" label="资源 ID">
              <a-input-number v-model="resourceModel.resourceId" :min="0" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item field="courseId" label="课程 ID">
              <a-input-number v-model="resourceModel.courseId" :min="0" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item field="fileId" label="文件 ID">
              <a-input-number v-model="resourceModel.fileId" :min="0" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { Message, Modal } from '@arco-design/web-vue';
import { onMounted, reactive, ref } from 'vue';
import {
  addExpert,
  addExpertResource,
  getExpertResources,
  getExperts,
  removeExpertResources,
  removeExperts,
  updateExpert,
  updateExpertResource
} from '@/api/edu/expert';

const resourceTypeOptions = [
  { label: '资源', value: 'resource' },
  { label: '课程', value: 'course' },
  { label: '讲座', value: 'lecture' },
  { label: '文件', value: 'file' }
];
const resourceTypeText = Object.fromEntries(resourceTypeOptions.map((item) => [item.value, item.label]));
const queryForm = reactive({ keyword: '', pageIndex: 1, pageSize: 10 });
const tableData = ref([]);
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const formVisible = ref(false);
const resourceVisible = ref(false);
const resourceFormVisible = ref(false);
const currentExpert = ref(null);
const resourceList = ref([]);
const formModel = reactive(defaultForm());
const resourceModel = reactive(defaultResourceForm());

const columns = [
  { title: '姓名', dataIndex: 'name', width: 120 },
  { title: '职称', dataIndex: 'title', width: 130 },
  { title: '机构', dataIndex: 'organization', ellipsis: true, tooltip: true },
  { title: '擅长领域', dataIndex: 'specialties', ellipsis: true, tooltip: true },
  { title: '推荐', dataIndex: 'isRecommended', width: 80 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'operations', width: 190 }
];
const resourceColumns = [
  { title: '标题', dataIndex: 'title', ellipsis: true, tooltip: true },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '资源 ID', dataIndex: 'resourceId', width: 90 },
  { title: '课程 ID', dataIndex: 'courseId', width: 90 },
  { title: '文件 ID', dataIndex: 'fileId', width: 90 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'operations', width: 150 }
];

function defaultForm() {
  return { id: undefined, name: '', title: '', organization: '', specialties: '', introduction: '', isRecommended: 0, sort: 0, status: 1 };
}

function defaultResourceForm() {
  return { id: undefined, title: '', type: 'resource', resourceId: 0, courseId: 0, fileId: 0, status: 1 };
}

function assignForm(data = {}) {
  Object.assign(formModel, defaultForm(), data);
}

async function fetchData() {
  const res = await getExperts(queryForm);
  const payload = res.data || {};
  tableData.value = payload.list || payload || [];
  pagination.total = payload.count || res.total || 0;
}

function handlePageChange(page) {
  queryForm.pageIndex = page;
  pagination.current = page;
  fetchData();
}

function resetQuery() {
  queryForm.keyword = '';
  queryForm.pageIndex = 1;
  pagination.current = 1;
  fetchData();
}

function openCreate() {
  assignForm();
  formVisible.value = true;
}

function openEdit(record) {
  assignForm(record);
  formVisible.value = true;
}

async function handleSave() {
  if (!formModel.name) {
    Message.warning('请输入专家姓名');
    return false;
  }
  const payload = { ...formModel };
  if (payload.id) {
    await updateExpert(payload.id, payload);
  } else {
    await addExpert(payload);
  }
  Message.success('保存成功');
  formVisible.value = false;
  fetchData();
}

function handleDelete(record) {
  Modal.confirm({
    title: '确认删除专家',
    content: `确定删除「${record.name}」吗？`,
    async onOk() {
      await removeExperts({ ids: [record.id] });
      Message.success('删除成功');
      fetchData();
    }
  });
}

async function openResources(record) {
  currentExpert.value = record;
  resourceVisible.value = true;
  await fetchResources();
}

async function fetchResources() {
  if (!currentExpert.value?.id) return;
  const res = await getExpertResources(currentExpert.value.id);
  resourceList.value = res.data || [];
}

function openResourceCreate() {
  Object.assign(resourceModel, defaultResourceForm());
  resourceFormVisible.value = true;
}

function openResourceEdit(record) {
  Object.assign(resourceModel, defaultResourceForm(), record);
  resourceFormVisible.value = true;
}

async function handleResourceSave() {
  if (!resourceModel.title) {
    Message.warning('请输入资源标题');
    return false;
  }
  const payload = { ...resourceModel };
  if (payload.id) {
    await updateExpertResource(currentExpert.value.id, payload.id, payload);
  } else {
    await addExpertResource(currentExpert.value.id, payload);
  }
  Message.success('保存成功');
  resourceFormVisible.value = false;
  fetchResources();
}

function handleResourceDelete(record) {
  Modal.confirm({
    title: '确认删除关联资源',
    content: `确定删除「${record.title}」吗？`,
    async onOk() {
      await removeExpertResources(currentExpert.value.id, { ids: [record.id] });
      Message.success('删除成功');
      fetchResources();
    }
  });
}

onMounted(fetchData);
</script>

<style scoped>
.table-card {
  margin-top: 16px;
}
</style>
