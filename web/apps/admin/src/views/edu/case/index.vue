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
            <a-option value="draft">草稿</a-option>
            <a-option value="reviewing">审核中</a-option>
            <a-option value="archived">已归档</a-option>
          </a-select>
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
import { getCases } from '@/api/edu/case';

const queryForm = reactive({ keyword: '', status: '', pageIndex: 1, pageSize: 10 });
const tableData = ref([]);
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const columns = [
  { title: '案例名称', dataIndex: 'title' },
  { title: '学生姓名', dataIndex: 'studentName' },
  { title: '学生编号', dataIndex: 'studentCode' },
  { title: '障碍类型', dataIndex: 'disabilityType' },
  { title: '状态', dataIndex: 'status' }
];

async function fetchData() {
  const res = await getCases(queryForm);
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
  queryForm.status = '';
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

