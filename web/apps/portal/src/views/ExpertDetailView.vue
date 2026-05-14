<template>
  <PortalLayout>
    <a-spin :loading="loading" style="width: 100%">
      <template v-if="expert">
        <section class="detail-hero">
          <a-breadcrumb>
            <a-breadcrumb-item>
              <router-link to="/experts">名师资源</router-link>
            </a-breadcrumb-item>
            <a-breadcrumb-item>{{ expert.name }}</a-breadcrumb-item>
          </a-breadcrumb>
          <h1>{{ expert.name }}</h1>
          <p>{{ expert.introduction || expert.specialties || "暂无专家简介" }}</p>
          <div class="meta-row">
            <a-tag color="purple">{{ expert.title || "专家" }}</a-tag>
            <a-tag>{{ expert.organization || "平台专家库" }}</a-tag>
          </div>
        </section>

        <section class="detail-layout">
          <article class="detail-panel">
            <h2>专家介绍</h2>
            <p>{{ expert.introduction || "暂未填写专家介绍。" }}</p>
            <div class="outline-list">
              <h2>相关资源</h2>
              <div v-if="resources.length" class="outline-chapters">
                <component
                  :is="resourceTarget(item) ? 'router-link' : 'section'"
                  v-for="item in resources"
                  :key="item.id"
                  :to="resourceTarget(item)"
                  class="outline-chapter"
                >
                  <strong>{{ item.title }}</strong>
                  <span>{{ resourceTypeText[item.type] || "资源" }}</span>
                </component>
              </div>
              <a-empty v-else description="暂无关联资源" />
            </div>
          </article>
          <aside class="side-panel">
            <h2>专业方向</h2>
            <p>{{ expert.specialties || "暂未设置专业方向。" }}</p>
          </aside>
        </section>
      </template>
      <a-empty v-else description="暂无专家详情" />
    </a-spin>
  </PortalLayout>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getPublishedExpert } from "@/api/experts";

const resourceTypeText = {
  resource: "资源",
  course: "课程",
  lecture: "讲座",
  file: "文件"
};
const route = useRoute();
const loading = ref(false);
const expert = ref(null);
const resources = ref([]);

function resourceTarget(item) {
  if (item.type === "resource" && item.resourceId) {
    return `/resources/${item.resourceId}`;
  }
  if ((item.type === "course" || item.type === "lecture") && item.courseId) {
    return `/courses/${item.courseId}`;
  }
  return "";
}

async function fetchExpert() {
  loading.value = true;
  try {
    const res = await getPublishedExpert(route.params.id);
    expert.value = res.data?.expert || res.data || null;
    resources.value = res.data?.resources || [];
  } finally {
    loading.value = false;
  }
}

onMounted(fetchExpert);
</script>
