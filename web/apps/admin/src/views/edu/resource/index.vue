<template>
  <div class="container">
    <a-card :bordered="false" class="cardStyle">
      <template #title>资源管理</template>
      <a-form :model="queryForm" layout="inline">
        <a-form-item label="关键词">
          <a-input v-model="queryForm.keyword" allow-clear placeholder="标题、简介、标签" />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model="queryForm.status" allow-clear placeholder="请选择状态" style="width: 160px">
            <a-option value="draft">草稿</a-option>
            <a-option value="reviewing">审核中</a-option>
            <a-option value="published">已发布</a-option>
            <a-option value="rejected">已驳回</a-option>
            <a-option value="offline">已下架</a-option>
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
      <a-table :columns="columns" :data="tableData" :pagination="pagination" row-key="id" @page-change="handlePageChange">
        <template #status="{ record }">
          <a-tag>{{ statusText[record.status] || record.status }}</a-tag>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue';
import { getResources } from '@/api/edu/resource';

const queryForm = reactive({ keyword: '', status: '', pageIndex: 1, pageSize: 10 });
const tableData = ref([]);
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const statusText = {
  draft: '草稿',
  reviewing: '审核中',
  published: '已发布',
  rejected: '已驳回',
  offline: '已下架'
};

const columns = [
  { title: '资源标题', dataIndex: 'title' },
  { title: '作者', dataIndex: 'authorName' },
  { title: '学校', dataIndex: 'schoolId' },
  { title: '状态', slotName: 'status' },
  { title: '浏览', dataIndex: 'viewCount' },
  { title: '下载', dataIndex: 'downloadCount' }
];

async function fetchData() {
  const res = await getResources(queryForm);
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

