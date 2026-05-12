<template>
  <PortalLayout>
    <section class="hero">
      <div>
        <p class="eyebrow">Special Education Resource Platform</p>
        <h1>特殊教育资源库平台</h1>
        <p class="summary">汇聚教学资源、专题课程、教研活动、特教案例与专家资源，服务教师教学和学生个性化学习。</p>
      </div>
      <a-input-search class="hero-search" size="large" placeholder="搜索课程、资源、案例或专家" search-button />
    </section>

    <section class="quick-grid">
      <router-link v-for="item in quickLinks" :key="item.path" :to="item.path" class="quick-item">
        <strong>{{ item.title }}</strong>
        <span>{{ item.description }}</span>
      </router-link>
    </section>
  </PortalLayout>
</template>

<script setup>
import PortalLayout from "@/layouts/PortalLayout.vue";

const quickLinks = [
  { title: "资源中心", description: "检索、预览和收藏教学资源", path: "/resources" },
  { title: "专题课程", description: "进入课程学习和进度记录", path: "/courses" },
  { title: "教研活动", description: "查看活动、报名和成果", path: "/activities" },
  { title: "名师资源", description: "浏览专家讲座和优质资料", path: "/experts" }
];
</script>
