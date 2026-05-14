<template>
  <PortalLayout>
    <section class="page-heading">
      <h1>名师资源</h1>
      <a-input-search v-model="query.keyword" placeholder="搜索专家姓名、领域" search-button @search="fetchExperts" />
    </section>
    <a-spin :loading="loading" style="width: 100%">
      <div v-if="experts.length" class="content-grid">
        <router-link v-for="item in experts" :key="item.id" :to="`/experts/${item.id}`" class="content-card">
          <strong>{{ item.name }}</strong>
          <span>{{ item.specialties || item.introduction || "暂无简介" }}</span>
          <small>{{ item.title || "专家" }} · {{ item.organization || "平台专家库" }}</small>
        </router-link>
      </div>
      <a-empty v-else description="暂无专家资源" />
    </a-spin>
  </PortalLayout>
</template>

<script setup>
import { onMounted, reactive, ref } from "vue";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getPublishedExperts } from "@/api/experts";

const loading = ref(false);
const experts = ref([]);
const query = reactive({ keyword: "", pageIndex: 1, pageSize: 12 });

async function fetchExperts() {
  loading.value = true;
  try {
    const res = await getPublishedExperts(query);
    experts.value = res.data?.list || res.data || [];
  } finally {
    loading.value = false;
  }
}

onMounted(fetchExperts);
</script>
