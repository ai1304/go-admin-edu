<template>
  <PortalLayout>
    <a-spin :loading="loading" style="width: 100%">
      <template v-if="resource">
        <section class="detail-hero">
          <a-breadcrumb>
            <a-breadcrumb-item>
              <router-link to="/resources">资源中心</router-link>
            </a-breadcrumb-item>
            <a-breadcrumb-item>{{ resource.title }}</a-breadcrumb-item>
          </a-breadcrumb>
          <h1>{{ resource.title }}</h1>
          <p>{{ resource.summary || "暂无简介" }}</p>
          <div class="meta-row">
            <a-tag color="blue">{{ resource.authorName || "平台资源" }}</a-tag>
            <a-tag>{{ resource.keywords || "暂无关键词" }}</a-tag>
            <a-tag>{{ resource.viewCount || 0 }} 次浏览</a-tag>
          </div>
        </section>

        <section class="detail-layout">
          <article class="detail-panel">
            <h2>资源介绍</h2>
            <p>{{ resource.summary || "暂未填写资源介绍。" }}</p>
          </article>
          <aside class="side-panel">
            <h2>附件</h2>
            <div v-if="files.length" class="file-list">
              <div v-for="file in files" :key="file.id" class="file-item">
                <strong>{{ file.originalName }}</strong>
                <span>{{ file.contentType || file.ext || "未知类型" }} · {{ formatSize(file.size) }}</span>
              </div>
            </div>
            <a-empty v-else description="暂无附件" />
          </aside>
        </section>
      </template>
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
const files = ref([]);

async function fetchResource() {
  loading.value = true;
  try {
    const res = await getPublishedResource(route.params.id);
    resource.value = res.data?.resource || res.data || null;
    files.value = res.data?.files || [];
  } finally {
    loading.value = false;
  }
}

function formatSize(size = 0) {
  if (size < 1024) return `${size} B`;
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
  return `${(size / 1024 / 1024).toFixed(1)} MB`;
}

onMounted(fetchResource);
</script>
