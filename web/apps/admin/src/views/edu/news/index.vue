<template>
  <div class="container">
    <a-card :bordered="false" class="cardStyle">
      <template #title>资讯管理</template>
      <a-form :model="queryForm" layout="inline">
        <a-form-item label="关键词">
          <a-input v-model="queryForm.keyword" allow-clear placeholder="标题、摘要、来源" />
        </a-form-item>
        <a-form-item label="类型">
          <a-select v-model="queryForm.moduleType" allow-clear placeholder="全部类型" style="width: 160px">
            <a-option v-for="item in moduleOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model="queryForm.status" allow-clear placeholder="全部状态" style="width: 140px">
            <a-option v-for="item in statusOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button v-has="'edu:news:query'" type="primary" @click="fetchData">查询</a-button>
            <a-button @click="resetQuery">重置</a-button>
            <a-button v-has="'edu:news:add'" type="primary" status="success" @click="openCreate">新增资讯</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>

    <a-card :bordered="false" class="cardStyle table-card">
      <a-table :columns="columns" :data="tableData" :pagination="pagination" row-key="id" @page-change="handlePageChange">
        <template #moduleType="{ record }">
          <a-tag color="arcoblue">{{ moduleText[record.moduleType] || record.moduleType }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="statusColor[record.status]">{{ statusText[record.status] || record.status }}</a-tag>
        </template>
        <template #isTop="{ record }">
          <a-tag :color="record.isTop ? 'orange' : 'gray'">{{ record.isTop ? '置顶' : '普通' }}</a-tag>
        </template>
        <template #operations="{ record }">
          <a-space>
            <a-button v-has="'edu:news:edit'" type="text" size="small" @click="openEdit(record)">编辑</a-button>
            <a-button v-has="'edu:news:remove'" type="text" status="danger" size="small" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="formModel.id ? '编辑资讯' : '新增资讯'" width="860px" @before-ok="handleSave">
      <a-form :model="formModel" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="16">
            <a-form-item field="title" label="标题" required>
              <a-input v-model="formModel.title" placeholder="请输入资讯标题" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item field="moduleType" label="类型">
              <a-select v-model="formModel.moduleType">
                <a-option v-for="item in moduleOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item field="source" label="来源">
              <a-input v-model="formModel.source" placeholder="发布单位或来源" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item field="publishTime" label="发布时间">
              <a-input v-model="formModel.publishTime" placeholder="2026-05-20 09:00" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item field="status" label="状态">
              <a-select v-model="formModel.status">
                <a-option v-for="item in statusOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item field="coverUrl" label="封面地址">
              <a-input v-model="formModel.coverUrl" placeholder="可填写外部图片 URL" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item field="keywords" label="关键词">
              <a-input v-model="formModel.keywords" placeholder="逗号分隔" />
            </a-form-item>
          </a-col>
          <a-col :span="4">
            <a-form-item field="isTop" label="置顶">
              <a-switch v-model="formModel.isTop" :checked-value="1" :unchecked-value="0" />
            </a-form-item>
          </a-col>
          <a-col :span="4">
            <a-form-item field="sort" label="排序">
              <a-input-number v-model="formModel.sort" :min="0" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item field="summary" label="摘要">
              <a-textarea v-model="formModel.summary" :auto-size="{ minRows: 2, maxRows: 4 }" placeholder="请输入摘要" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item field="content" label="正文">
              <a-textarea v-model="formModel.content" :auto-size="{ minRows: 8, maxRows: 16 }" placeholder="请输入正文内容" />
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
import { addNews, getNews, removeNews, updateNews } from '@/api/edu/news';

const moduleOptions = [
  { label: '政策法规', value: 'POLICY' },
  { label: '学术前沿', value: 'ACADEMIC' },
  { label: '行业动态', value: 'INDUSTRY' },
  { label: '优秀实践', value: 'PRACTICE' }
];
const statusOptions = [
  { label: '草稿', value: 'draft' },
  { label: '已发布', value: 'published' },
  { label: '已下线', value: 'offline' }
];
const moduleText = Object.fromEntries(moduleOptions.map((item) => [item.value, item.label]));
const statusText = Object.fromEntries(statusOptions.map((item) => [item.value, item.label]));
const statusColor = { draft: 'gray', published: 'green', offline: 'red' };
const queryForm = reactive({ keyword: '', moduleType: '', status: '', pageIndex: 1, pageSize: 10 });
const tableData = ref([]);
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const formVisible = ref(false);
const formModel = reactive(defaultForm());

const columns = [
  { title: '标题', dataIndex: 'title', ellipsis: true, tooltip: true },
  { title: '类型', slotName: 'moduleType', width: 120 },
  { title: '来源', dataIndex: 'source', width: 160 },
  { title: '发布时间', dataIndex: 'publishTime', width: 160 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '置顶', slotName: 'isTop', width: 90 },
  { title: '浏览', dataIndex: 'viewCount', width: 90 },
  { title: '操作', slotName: 'operations', width: 150, fixed: 'right' }
];

function defaultForm() {
  return {
    id: 0,
    title: '',
    moduleType: 'POLICY',
    source: '',
    coverUrl: '',
    summary: '',
    content: '',
    keywords: '',
    publishTime: '',
    status: 'draft',
    isTop: 0,
    sort: 0
  };
}

function assignForm(data) {
  Object.assign(formModel, defaultForm(), data || {});
}

async function fetchData() {
  const res = await getNews(queryForm);
  const payload = res.data || {};
  tableData.value = payload.list || [];
  pagination.total = payload.count || 0;
  pagination.current = queryForm.pageIndex;
}

function resetQuery() {
  Object.assign(queryForm, { keyword: '', moduleType: '', status: '', pageIndex: 1 });
  fetchData();
}

function handlePageChange(page) {
  queryForm.pageIndex = page;
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
    Message.warning('请输入资讯标题');
    return false;
  }
  if (formModel.id) {
    await updateNews(formModel.id, formModel);
  } else {
    await addNews(formModel);
  }
  Message.success('保存成功');
  await fetchData();
  return true;
}

function handleDelete(record) {
  Modal.warning({
    title: '确认删除',
    content: `确定删除“${record.title}”？`,
    hideCancel: false,
    onOk: async () => {
      await removeNews({ ids: [record.id] });
      Message.success('删除成功');
      fetchData();
    }
  });
}

onMounted(fetchData);
</script>
