<template>
  <div class="container">
    <a-card :bordered="false" class="cardStyle">
      <template #title>课程管理</template>
      <a-form :model="queryForm" layout="inline">
        <a-form-item label="关键词">
          <a-input v-model="queryForm.keyword" allow-clear placeholder="课程标题、教师" />
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
            <a-button type="primary" status="success" @click="openCreate">新增课程</a-button>
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

    <a-modal v-model:visible="formVisible" :title="formModel.id ? '编辑课程' : '新增课程'" width="760px" @before-ok="handleSave">
      <a-form :model="formModel" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item field="title" label="课程标题" required>
              <a-input v-model="formModel.title" placeholder="请输入课程标题" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="teacherName" label="授课教师">
              <a-input v-model="formModel.teacherName" placeholder="请输入授课教师" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="difficulty" label="难度">
              <a-select v-model="formModel.difficulty" allow-clear placeholder="请选择难度">
                <a-option value="basic">基础</a-option>
                <a-option value="advanced">进阶</a-option>
                <a-option value="expert">专家</a-option>
              </a-select>
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
            <a-form-item field="summary" label="课程简介">
              <a-textarea v-model="formModel.summary" :auto-size="{ minRows: 3, maxRows: 5 }" placeholder="请输入课程简介" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item field="objectives" label="教学目标">
              <a-textarea v-model="formModel.objectives" :auto-size="{ minRows: 3, maxRows: 5 }" placeholder="请输入教学目标" />
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
import { addCourse, getCourses, removeCourses, updateCourse } from '@/api/edu/course';

const statusOptions = [
  { label: '草稿', value: 'draft' },
  { label: '已发布', value: 'published' },
  { label: '已下架', value: 'offline' }
];
const statusText = Object.fromEntries(statusOptions.map((item) => [item.value, item.label]));
const statusColor = { draft: 'gray', published: 'green', offline: 'gray' };
const queryForm = reactive({ keyword: '', status: '', pageIndex: 1, pageSize: 10 });
const tableData = ref([]);
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const formVisible = ref(false);
const formModel = reactive(defaultForm());

const columns = [
  { title: '课程标题', dataIndex: 'title', ellipsis: true, tooltip: true },
  { title: '教师', dataIndex: 'teacherName', width: 120 },
  { title: '难度', dataIndex: 'difficulty', width: 100 },
  { title: '状态', slotName: 'status', width: 110 },
  { title: '学习人数', dataIndex: 'learnerCount', width: 110 },
  { title: '操作', slotName: 'operations', width: 150 }
];

function defaultForm() {
  return { id: undefined, title: '', summary: '', teacherName: '', difficulty: '', objectives: '', status: 'draft' };
}

function assignForm(data = {}) {
  Object.assign(formModel, defaultForm(), data);
}

async function fetchData() {
  const res = await getCourses(queryForm);
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
    Message.warning('请输入课程标题');
    return false;
  }
  const payload = { ...formModel };
  if (payload.id) {
    await updateCourse(payload.id, payload);
  } else {
    await addCourse(payload);
  }
  Message.success('保存成功');
  formVisible.value = false;
  fetchData();
}

function handleDelete(record) {
  Modal.confirm({
    title: '确认删除课程',
    content: `确定删除「${record.title}」吗？`,
    async onOk() {
      await removeCourses({ ids: [record.id] });
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
