<template>
  <div class="container">
    <a-card :bordered="false" class="cardStyle">
      <template #title>资源标签管理</template>
      <a-form :model="queryForm" layout="inline">
        <a-form-item label="标签名称">
          <a-input v-model="queryForm.name" allow-clear placeholder="请输入标签名称" />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button v-has="'edu:tag:query'" type="primary" @click="fetchData">查询</a-button>
            <a-button @click="resetQuery">重置</a-button>
            <a-button v-has="'edu:tag:add'" type="primary" status="success" @click="openCreate">新增标签</a-button>
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
            <a-button v-has="'edu:tag:edit'" type="text" size="small" @click="openEdit(record)">编辑</a-button>
            <a-button v-has="'edu:tag:remove'" type="text" status="danger" size="small" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="formModel.id ? '编辑标签' : '新增标签'" width="520px" @before-ok="handleSave">
      <a-form :model="formModel" layout="vertical">
        <a-form-item field="name" label="标签名称" required>
          <a-input v-model="formModel.name" placeholder="请输入标签名称" />
        </a-form-item>
        <a-form-item field="status" label="状态">
          <a-select v-model="formModel.status">
            <a-option :value="1">启用</a-option>
            <a-option :value="0">停用</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { Message, Modal } from '@arco-design/web-vue';
import { onMounted, reactive, ref } from 'vue';
import { addResourceTag, getResourceTags, removeResourceTags, updateResourceTag } from '@/api/edu/resource';

const queryForm = reactive({ name: '', pageIndex: 1, pageSize: 10 });
const tableData = ref([]);
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const formVisible = ref(false);
const formModel = reactive(defaultForm());

const columns = [
  { title: '标签名称', dataIndex: 'name' },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'operations', width: 150 }
];

function defaultForm() {
  return { id: undefined, name: '', status: 1 };
}

function assignForm(data = {}) {
  Object.assign(formModel, defaultForm(), data);
}

async function fetchData() {
  const res = await getResourceTags(queryForm);
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
  queryForm.name = '';
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
    Message.warning('请输入标签名称');
    return false;
  }
  const payload = { ...formModel };
  if (payload.id) {
    await updateResourceTag(payload.id, payload);
  } else {
    await addResourceTag(payload);
  }
  Message.success('保存成功');
  formVisible.value = false;
  fetchData();
}

function handleDelete(record) {
  Modal.confirm({
    title: '确认删除标签',
    content: `确定删除「${record.name}」吗？`,
    async onOk() {
      await removeResourceTags({ ids: [record.id] });
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
