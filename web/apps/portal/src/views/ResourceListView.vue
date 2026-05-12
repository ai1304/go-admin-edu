<template>
  <PortalLayout>
    <section class="page-heading">
      <h1>资源中心</h1>
      <a-input-search placeholder="搜索资源标题、标签、作者" search-button />
    </section>
    <a-empty description="暂无资源" />
  </PortalLayout>
</template>

<script setup>
import PortalLayout from "@/layouts/PortalLayout.vue";
</script>
