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
import { getExperts } from '@/api/edu/expert';

const queryForm = reactive({ keyword: '', pageIndex: 1, pageSize: 10 });
const tableData = ref([]);
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const columns = [
  { title: '姓名', dataIndex: 'name' },
  { title: '职称', dataIndex: 'title' },
  { title: '机构', dataIndex: 'organization' },
  { title: '擅长领域', dataIndex: 'specialties' },
  { title: '状态', dataIndex: 'status' }
];

async function fetchData() {
  const res = await getExperts(queryForm);
  tableData.value = res.data?.list || res.data || [];
  pagination.total = res.data?.count || res.total || 0;
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

onMounted(fetchData);
</script>

<style scoped>
.table-card {
  margin-top: 16px;
}
</style>

