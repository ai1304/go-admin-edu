<template>
  <PortalLayout>
    <section class="page-heading">
      <h1>教研活动</h1>
      <a-input-search v-model="query.keyword" placeholder="搜索活动名称、主办方" search-button @search="fetchActivities" />
    </section>
    <a-spin :loading="loading" style="width: 100%">
      <div v-if="activities.length" class="content-grid">
        <article v-for="item in activities" :key="item.id" class="content-card">
          <strong>{{ item.title }}</strong>
          <span>{{ item.summary || "暂无简介" }}</span>
          <small>{{ item.startTime || "时间待定" }} · {{ item.location || "地点待定" }}</small>
        </article>
      </div>
      <a-empty v-else description="暂无活动" />
    </a-spin>
  </PortalLayout>
</template>

<script setup>
import { onMounted, reactive, ref } from "vue";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getPublishedActivities } from "@/api/activities";

const loading = ref(false);
const activities = ref([]);
const query = reactive({ keyword: "", pageIndex: 1, pageSize: 12 });

async function fetchActivities() {
  loading.value = true;
  try {
    const res = await getPublishedActivities(query);
    activities.value = res.data?.list || res.data || [];
  } finally {
    loading.value = false;
  }
}

onMounted(fetchActivities);
</script>
