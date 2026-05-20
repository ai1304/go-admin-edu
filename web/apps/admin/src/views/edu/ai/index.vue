<template>
  <div class="container">
    <a-row :gutter="16" class="summary-row">
      <a-col :span="8">
        <a-card :bordered="false" class="cardStyle stat-card">
          <span>会话数</span>
          <strong>{{ stats.conversationCount || 0 }}</strong>
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card :bordered="false" class="cardStyle stat-card">
          <span>消息数</span>
          <strong>{{ stats.messageCount || 0 }}</strong>
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card :bordered="false" class="cardStyle stat-card">
          <span>模型状态</span>
          <strong>{{ stats.enabled ? stats.model : '离线规则' }}</strong>
        </a-card>
      </a-col>
    </a-row>

    <a-card :bordered="false" class="cardStyle">
      <template #title>AI 会话监控</template>
      <a-form :model="queryForm" layout="inline">
        <a-form-item label="关键词">
          <a-input v-model="queryForm.keyword" allow-clear placeholder="会话标题或访客标识" />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button v-has="'edu:ai:query'" type="primary" @click="fetchData">查询</a-button>
            <a-button @click="resetQuery">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>

    <a-card :bordered="false" class="cardStyle table-card">
      <a-table :columns="columns" :data="tableData" :pagination="pagination" row-key="id" @page-change="handlePageChange">
        <template #mode="{ record }">
          <a-tag :color="record.mode === 'online' ? 'green' : 'gray'">{{ record.mode === 'online' ? '在线模型' : '离线规则' }}</a-tag>
        </template>
        <template #operations="{ record }">
          <a-space>
            <a-button v-has="'edu:ai:query'" type="text" size="small" @click="openDetail(record)">查看</a-button>
            <a-button v-has="'edu:ai:remove'" type="text" status="danger" size="small" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="detailVisible" :title="currentConversation?.title || '会话详情'" width="760px" :footer="false">
      <div class="message-list">
        <div v-for="message in messages" :key="message.id" class="message-item" :class="message.role">
          <div class="message-role">{{ message.role === 'user' ? '用户' : '助手' }}</div>
          <div class="message-content">{{ message.content }}</div>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { Message, Modal } from '@arco-design/web-vue';
import { onMounted, reactive, ref } from 'vue';
import { getAIConversation, getAIConversations, getAIStats, removeAIConversation } from '@/api/edu/ai';

const stats = reactive({ conversationCount: 0, messageCount: 0, enabled: false, model: '' });
const queryForm = reactive({ keyword: '', pageIndex: 1, pageSize: 10 });
const tableData = ref([]);
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const detailVisible = ref(false);
const currentConversation = ref(null);
const messages = ref([]);

const columns = [
  { title: '标题', dataIndex: 'title', ellipsis: true, tooltip: true },
  { title: '访客标识', dataIndex: 'clientKey', ellipsis: true, tooltip: true },
  { title: '用户ID', dataIndex: 'userId', width: 100 },
  { title: '模式', slotName: 'mode', width: 120 },
  { title: '更新时间', dataIndex: 'updatedAt', width: 180 },
  { title: '操作', slotName: 'operations', width: 150, fixed: 'right' }
];

async function fetchStats() {
  const res = await getAIStats();
  Object.assign(stats, res.data || {});
}

async function fetchData() {
  const res = await getAIConversations(queryForm);
  const payload = res.data || {};
  tableData.value = payload.list || [];
  pagination.total = payload.count || 0;
  pagination.current = queryForm.pageIndex;
}

function resetQuery() {
  Object.assign(queryForm, { keyword: '', pageIndex: 1 });
  fetchData();
}

function handlePageChange(page) {
  queryForm.pageIndex = page;
  fetchData();
}

async function openDetail(record) {
  const res = await getAIConversation(record.id);
  currentConversation.value = res.data || record;
  messages.value = currentConversation.value.messages || [];
  detailVisible.value = true;
}

function handleDelete(record) {
  Modal.warning({
    title: '确认删除',
    content: `确定删除会话“${record.title || record.id}”？`,
    hideCancel: false,
    onOk: async () => {
      await removeAIConversation(record.id);
      Message.success('删除成功');
      fetchStats();
      fetchData();
    }
  });
}

onMounted(() => {
  fetchStats();
  fetchData();
});
</script>

<style scoped>
.summary-row {
  margin-bottom: 16px;
}

.stat-card {
  min-height: 96px;
}

.stat-card span {
  display: block;
  color: #86909c;
  margin-bottom: 10px;
}

.stat-card strong {
  display: block;
  font-size: 26px;
  color: #1d2129;
}

.message-list {
  display: grid;
  gap: 12px;
  max-height: 560px;
  overflow: auto;
}

.message-item {
  padding: 12px;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
  background: #f7f8fa;
}

.message-item.assistant {
  background: #f2f7ff;
}

.message-role {
  margin-bottom: 6px;
  font-size: 12px;
  color: #86909c;
}

.message-content {
  white-space: pre-wrap;
  line-height: 1.8;
}
</style>
