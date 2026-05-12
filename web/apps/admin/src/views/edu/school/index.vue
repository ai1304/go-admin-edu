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
import { getSchools } from '@/api/edu/school';

const queryForm = reactive({ name: '', code: '', pageIndex: 1, pageSize: 10 });
const tableData = ref([]);
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });

const columns = [
  { title: '学校名称', dataIndex: 'name' },
  { title: '学校编码', dataIndex: 'code' },
  { title: '区域', dataIndex: 'regionId' },
  { title: '联系人', dataIndex: 'contact' },
  { title: '电话', dataIndex: 'phone' },
  { title: '状态', dataIndex: 'status' }
];

async function fetchData() {
  const res = await getSchools(queryForm);
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

