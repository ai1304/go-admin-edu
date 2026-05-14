<template>
  <PortalLayout>
    <section class="page-heading">
      <h1>专题课程</h1>
      <a-input-search v-model="query.keyword" placeholder="搜索课程标题、教师" search-button @search="fetchCourses" />
    </section>
    <a-spin :loading="loading" style="width: 100%">
      <div v-if="courses.length" class="content-grid">
        <article v-for="item in courses" :key="item.id" class="content-card">
          <strong>{{ item.title }}</strong>
          <span>{{ item.summary || "暂无简介" }}</span>
          <small>{{ item.teacherName || "平台课程" }}</small>
        </article>
      </div>
      <a-empty v-else description="暂无课程" />
    </a-spin>
  </PortalLayout>
</template>

<script setup>
import { onMounted, reactive, ref } from "vue";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getPublishedCourses } from "@/api/courses";

const loading = ref(false);
const courses = ref([]);
const query = reactive({ keyword: "", pageIndex: 1, pageSize: 12 });

async function fetchCourses() {
  loading.value = true;
  try {
    const res = await getPublishedCourses(query);
    courses.value = res.data?.list || res.data || [];
  } finally {
    loading.value = false;
  }
}

onMounted(fetchCourses);
</script>
