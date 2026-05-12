<template>
  <PortalLayout>
    <section class="page-heading">
      <h1>资源中心</h1>
      <a-input-search v-model="query.keyword" placeholder="搜索资源标题、标签、作者" search-button @search="fetchResources" />
    </section>
    <a-spin :loading="loading" style="width: 100%">
      <div v-if="resources.length" class="resource-grid">
        <router-link v-for="item in resources" :key="item.id" :to="`/resources/${item.id}`" class="resource-card">
          <strong>{{ item.title }}</strong>
          <span>{{ item.summary || "暂无简介" }}</span>
          <small>{{ item.authorName || "平台资源" }}</small>
        </router-link>
      </div>
      <a-empty v-else description="暂无资源" />
    </a-spin>
  </PortalLayout>
</template>

<script setup>
import { onMounted, reactive, ref } from "vue";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getPublishedResources } from "@/api/resources";

const loading = ref(false);
const resources = ref([]);
const query = reactive({ keyword: "", pageIndex: 1, pageSize: 12 });

async function fetchResources() {
  loading.value = true;
  try {
    const res = await getPublishedResources(query);
    resources.value = res.data?.list || res.data || [];
  } finally {
    loading.value = false;
  }
}

onMounted(fetchResources);
</script>

<style scoped>
.resource-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.resource-card {
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-height: 148px;
  padding: 18px;
  background: #fff;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
}

.resource-card strong {
  font-size: 17px;
}

.resource-card span {
  color: #4e5969;
  line-height: 1.7;
}

.resource-card small {
  margin-top: auto;
  color: #86909c;
}

@media (max-width: 860px) {
  .resource-grid {
    grid-template-columns: 1fr;
  }
}
</style>
