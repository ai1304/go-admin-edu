<template>
  <div class="container">
    <a-card :bordered="false" class="cardStyle">
      <template #title>特教案例管理</template>
      <a-form :model="queryForm" layout="inline">
        <a-form-item label="关键词">
          <a-input v-model="queryForm.keyword" allow-clear placeholder="案例、学生姓名或编号" />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model="queryForm.status" allow-clear placeholder="请选择状态" style="width: 160px">
            <a-option v-for="item in statusOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="fetchData">查询</a-button>
            <a-button @click="resetQuery">重置</a-button>
            <a-button type="primary" status="success" @click="openCreate">新增案例</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>

    <a-card :bordered="false" class="cardStyle table-card">
      <a-table :columns="columns" :data="tableData" :pagination="pagination" row-key="id" @page-change="handlePageChange">
        <template #status="{ record }">
          <a-tag :color="statusColor[record.status]">{{ statusText[record.status] || record.status }}</a-tag>
        </template>
        <template #operations="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openEdit(record)">编辑</a-button>
            <a-button type="text" status="danger" size="small" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="formModel.id ? '编辑案例' : '新增案例'" width="760px" @before-ok="handleSave">
      <a-form :model="formModel" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item field="title" label="案例名称" required>
              <a-input v-model="formModel.title" placeholder="请输入案例名称" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="studentName" label="学生姓名">
              <a-input v-model="formModel.studentName" placeholder="请输入学生姓名" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="studentCode" label="学生编号">
              <a-input v-model="formModel.studentCode" placeholder="请输入学生编号" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="gender" label="性别">
              <a-select v-model="formModel.gender" allow-clear placeholder="请选择性别">
                <a-option value="male">男</a-option>
                <a-option value="female">女</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="birthday" label="生日">
              <a-input v-model="formModel.birthday" placeholder="例如 2015-09-01" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="disabilityType" label="障碍类型">
              <a-input v-model="formModel.disabilityType" placeholder="请输入障碍类型" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="status" label="状态">
              <a-select v-model="formModel.status">
                <a-option v-for="item in statusOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item field="summary" label="案例摘要">
              <a-textarea v-model="formModel.summary" :auto-size="{ minRows: 3, maxRows: 6 }" placeholder="请输入案例摘要" />
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
import { addCase, getCases, removeCases, updateCase } from '@/api/edu/case';

const statusOptions = [
  { label: '草稿', value: 'draft' },
  { label: '审核中', value: 'reviewing' },
  { label: '已归档', value: 'archived' }
];
const statusText = Object.fromEntries(statusOptions.map((item) => [item.value, item.label]));
const statusColor = { draft: 'gray', reviewing: 'orange', archived: 'blue' };
const queryForm = reactive({ keyword: '', status: '', pageIndex: 1, pageSize: 10 });
const tableData = ref([]);
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const formVisible = ref(false);
const formModel = reactive(defaultForm());

const columns = [
  { title: '案例名称', dataIndex: 'title', ellipsis: true, tooltip: true },
  { title: '学生姓名', dataIndex: 'studentName', width: 120 },
  { title: '学生编号', dataIndex: 'studentCode', width: 130 },
  { title: '障碍类型', dataIndex: 'disabilityType', width: 130 },
  { title: '状态', slotName: 'status', width: 110 },
  { title: '操作', slotName: 'operations', width: 150 }
];

function defaultForm() {
  return {
    id: undefined,
    title: '',
    studentName: '',
    studentCode: '',
    gender: '',
    birthday: '',
    disabilityType: '',
    summary: '',
    status: 'draft'
  };
}

function assignForm(data = {}) {
  Object.assign(formModel, defaultForm(), data);
}

async function fetchData() {
  const res = await getCases(queryForm);
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
  queryForm.status = '';
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
  if (!formModel.title) {
    Message.warning('请输入案例名称');
    return false;
  }
  const payload = { ...formModel };
  if (payload.id) {
    await updateCase(payload.id, payload);
  } else {
    await addCase(payload);
  }
  Message.success('保存成功');
  formVisible.value = false;
  fetchData();
}

function handleDelete(record) {
  Modal.confirm({
    title: '确认删除案例',
    content: `确定删除「${record.title}」吗？`,
    async onOk() {
      await removeCases({ ids: [record.id] });
      Message.success('删除成功');
      fetchData();
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
