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
            <a-button type="primary" @click="fetchData">查询</a-button>
            <a-button @click="resetQuery">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>
    <a-card :bordered="false" class="cardStyle table-card">
      <a-table :columns="columns" :data="tableData" :pagination="pagination" row-key="id" @page-change="handlePageChange" />
    </a-card>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue';
import { getRegions } from '@/api/edu/region';

const queryForm = reactive({ name: '', code: '', pageIndex: 1, pageSize: 10 });
const tableData = ref([]);
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });

const columns = [
  { title: '区域名称', dataIndex: 'name' },
  { title: '区域编码', dataIndex: 'code' },
  { title: '上级区域', dataIndex: 'parentId' },
  { title: '状态', dataIndex: 'status' },
  { title: '排序', dataIndex: 'sort' }
];

async function fetchData() {
  const res = await getRegions(queryForm);
  tableData.value = res.data?.list || res.data || [];
  pagination.total = res.data?.count || res.total || 0;
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

onMounted(fetchData);
</script>

<style scoped>
.table-card {
  margin-top: 16px;
}
</style>

