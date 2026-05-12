<template>
  <PortalLayout>
    <a-spin :loading="loading" style="width: 100%">
      <section v-if="resource" class="page-heading">
        <h1>{{ resource.title }}</h1>
        <p>{{ resource.summary || "暂无简介" }}</p>
        <a-space>
          <a-tag>{{ resource.authorName || "平台资源" }}</a-tag>
          <a-tag>{{ resource.status }}</a-tag>
        </a-space>
      </section>
      <a-empty v-else description="暂无资源详情" />
    </a-spin>
  </PortalLayout>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getPublishedResource } from "@/api/resources";

const route = useRoute();
const loading = ref(false);
const resource = ref(null);

async function fetchResource() {
  loading.value = true;
  try {
    const res = await getPublishedResource(route.params.id);
    resource.value = res.data?.resource || res.data || null;
  } finally {
    loading.value = false;
  }
}

onMounted(fetchResource);
</script>
