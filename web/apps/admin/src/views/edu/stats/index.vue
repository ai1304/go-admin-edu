<template>
  <div class="container">
    <a-card :bordered="false" class="cardStyle">
      <template #title>数据中心</template>
      <a-button type="primary" @click="fetchData">刷新</a-button>
    </a-card>
    <div class="stats-grid">
      <a-card v-for="item in items" :key="item.key" :bordered="false" class="stat-card">
        <div class="label">{{ item.label }}</div>
        <div class="value">{{ overview[item.key] ?? 0 }}</div>
      </a-card>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import { getEduOverview } from '@/api/edu/stats';

const overview = ref({});
const items = computed(() => [
  { key: 'regions', label: '区域数' },
  { key: 'schools', label: '学校数' },
  { key: 'resources', label: '资源数' },
  { key: 'publishedResources', label: '已发布资源' },
  { key: 'courses', label: '课程数' },
  { key: 'activities', label: '活动数' },
  { key: 'cases', label: '案例数' },
  { key: 'experts', label: '专家数' }
]);

async function fetchData() {
  const res = await getEduOverview();
  overview.value = res.data || {};
}

onMounted(fetchData);
</script>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
  margin-top: 16px;
}

.label {
  color: #86909c;
}

.value {
  margin-top: 8px;
  font-size: 30px;
  font-weight: 700;
}
</style>

