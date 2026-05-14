<template>
  <div class="container">
    <a-card :bordered="false" class="cardStyle">
      <template #title>学校管理</template>
      <a-form :model="queryForm" layout="inline">
        <a-form-item label="学校名称">
          <a-input v-model="queryForm.name" allow-clear placeholder="请输入学校名称" />
        </a-form-item>
        <a-form-item label="学校编码">
          <a-input v-model="queryForm.code" allow-clear placeholder="请输入学校编码" />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="fetchData">查询</a-button>
            <a-button @click="resetQuery">重置</a-button>
            <a-button type="primary" status="success" @click="openCreate">新增学校</a-button>
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
            <a-button type="text" size="small" @click="openEdit(record)">编辑</a-button>
            <a-button type="text" status="danger" size="small" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="formModel.id ? '编辑学校' : '新增学校'" width="720px" @before-ok="handleSave">
      <a-form :model="formModel" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item field="name" label="学校名称" required>
              <a-input v-model="formModel.name" placeholder="请输入学校名称" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="code" label="学校编码">
              <a-input v-model="formModel.code" placeholder="请输入学校编码" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="contact" label="联系人">
              <a-input v-model="formModel.contact" placeholder="请输入联系人" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="phone" label="电话">
              <a-input v-model="formModel.phone" placeholder="请输入电话" />
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
            <a-form-item field="address" label="地址">
              <a-input v-model="formModel.address" placeholder="请输入学校地址" />
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
import { addSchool, getSchools, removeSchools, updateSchool } from '@/api/edu/school';

const queryForm = reactive({ name: '', code: '', pageIndex: 1, pageSize: 10 });
const tableData = ref([]);
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const formVisible = ref(false);
const formModel = reactive(defaultForm());

const columns = [
  { title: '学校名称', dataIndex: 'name', ellipsis: true, tooltip: true },
  { title: '学校编码', dataIndex: 'code', width: 140 },
  { title: '联系人', dataIndex: 'contact', width: 120 },
  { title: '电话', dataIndex: 'phone', width: 140 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'operations', width: 150 }
];

function defaultForm() {
  return { id: undefined, name: '', code: '', address: '', contact: '', phone: '', status: 1, remark: '' };
}

function assignForm(data = {}) {
  Object.assign(formModel, defaultForm(), data);
}

async function fetchData() {
  const res = await getSchools(queryForm);
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
    Message.warning('请输入学校名称');
    return false;
  }
  const payload = { ...formModel };
  if (payload.id) {
    await updateSchool(payload.id, payload);
  } else {
    await addSchool(payload);
  }
  Message.success('保存成功');
  formVisible.value = false;
  fetchData();
}

function handleDelete(record) {
  Modal.confirm({
    title: '确认删除学校',
    content: `确定删除「${record.name}」吗？`,
    async onOk() {
      await removeSchools({ ids: [record.id] });
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
