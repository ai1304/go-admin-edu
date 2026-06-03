<template>
  <div class="container">
    <a-card :bordered="false" class="cardStyle">
      <template #title>区域管理</template>
      <a-form :model="queryForm" layout="inline">
        <a-form-item label="区域名称">
          <a-input v-model="queryForm.name" allow-clear placeholder="请输入区域名称" />
        </a-form-item>
        <a-form-item label="区域编码">
          <a-input v-model="queryForm.code" allow-clear placeholder="请输入区域编码" />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button v-has="'edu:region:query'" type="primary" @click="fetchData">查询</a-button>
            <a-button @click="resetQuery">重置</a-button>
            <a-button v-has="'edu:region:add'" type="primary" status="success" @click="openCreate">新增区域</a-button>
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
            <a-button v-has="'edu:region:edit'" type="text" size="small" @click="openEdit(record)">编辑</a-button>
            <a-button v-has="'edu:region:remove'" type="text" status="danger" size="small" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="formModel.id ? '编辑区域' : '新增区域'" width="640px" @before-ok="handleSave">
      <a-form :model="formModel" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item field="name" label="区域名称" required>
              <a-input v-model="formModel.name" placeholder="请输入区域名称" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="code" label="区域编码">
              <a-input v-model="formModel.code" placeholder="请输入区域编码" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="parentId" label="上级区域 ID">
              <a-input-number v-model="formModel.parentId" :min="0" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="sort" label="排序">
              <a-input-number v-model="formModel.sort" :min="0" style="width: 100%" />
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
          <a-col :span="24">
            <a-form-item field="remark" label="备注">
              <a-textarea v-model="formModel.remark" :auto-size="{ minRows: 3, maxRows: 5 }" placeholder="请输入备注" />
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
import { addRegion, getRegions, removeRegions, updateRegion } from '@/api/edu/region';

const queryForm = reactive({ name: '', code: '', pageIndex: 1, pageSize: 10 });
const tableData = ref([]);
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const formVisible = ref(false);
const formModel = reactive(defaultForm());

const columns = [
  { title: '区域名称', dataIndex: 'name' },
  { title: '区域编码', dataIndex: 'code', width: 140 },
  { title: '上级区域', dataIndex: 'parentId', width: 110 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '排序', dataIndex: 'sort', width: 90 },
  { title: '操作', slotName: 'operations', width: 150 }
];

function defaultForm() {
  return { id: undefined, name: '', code: '', parentId: 0, sort: 0, status: 1, remark: '' };
}

function assignForm(data = {}) {
  Object.assign(formModel, defaultForm(), data);
}

async function fetchData() {
  const res = await getRegions(queryForm);
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
  queryForm.code = '';
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
    Message.warning('请输入区域名称');
    return false;
  }
  const payload = { ...formModel };
  if (payload.id) {
    await updateRegion(payload.id, payload);
  } else {
    await addRegion(payload);
  }
  Message.success('保存成功');
  formVisible.value = false;
  fetchData();
}

function handleDelete(record) {
  Modal.confirm({
    title: '确认删除区域',
    content: `确定删除「${record.name}」吗？`,
    async onOk() {
      await removeRegions({ ids: [record.id] });
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
