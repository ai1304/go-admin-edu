<template>
  <PortalLayout>
    <a-spin :loading="loading" style="width: 100%">
      <template v-if="course">
        <section class="detail-hero">
          <a-breadcrumb>
            <a-breadcrumb-item>
              <router-link to="/courses">专题课程</router-link>
            </a-breadcrumb-item>
            <a-breadcrumb-item>{{ course.title }}</a-breadcrumb-item>
          </a-breadcrumb>
          <h1>{{ course.title }}</h1>
          <p>{{ course.summary || "暂无课程简介" }}</p>
          <div class="meta-row">
            <a-tag color="blue">{{ course.teacherName || "平台课程" }}</a-tag>
            <a-tag>{{ difficultyText[course.difficulty] || course.difficulty || "未设置难度" }}</a-tag>
            <a-tag>{{ course.learnerCount || 0 }} 人学习</a-tag>
          </div>
        </section>

        <section class="detail-layout">
          <article class="detail-panel">
            <h2>教学目标</h2>
            <p>{{ course.objectives || "暂未填写教学目标。" }}</p>
            <div class="outline-list">
              <h2>课程大纲</h2>
              <div v-if="chapters.length" class="outline-chapters">
                <section v-for="chapter in chapters" :key="chapter.id" class="outline-chapter">
                  <strong>{{ chapter.title }}</strong>
                  <div v-if="chapterLessons(chapter.id).length" class="outline-lessons">
                    <span v-for="lesson in chapterLessons(chapter.id)" :key="lesson.id">{{ lesson.title }} · {{ formatDuration(lesson.durationSeconds) }}</span>
                  </div>
                  <small v-else>暂无课时</small>
                </section>
              </div>
              <a-empty v-else description="暂无课程大纲" />
            </div>
          </article>
          <aside class="side-panel">
            <h2>课程信息</h2>
            <dl class="info-list">
              <div>
                <dt>授课教师</dt>
                <dd>{{ course.teacherName || "平台课程" }}</dd>
              </div>
              <div>
                <dt>课程难度</dt>
                <dd>{{ difficultyText[course.difficulty] || course.difficulty || "未设置" }}</dd>
              </div>
              <div>
                <dt>浏览量</dt>
                <dd>{{ course.viewCount || 0 }}</dd>
              </div>
            </dl>
          </aside>
        </section>
      </template>
      <a-empty v-else description="暂无课程详情" />
    </a-spin>
  </PortalLayout>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getPublishedCourse } from "@/api/courses";

const route = useRoute();
const loading = ref(false);
const course = ref(null);
const chapters = ref([]);
const lessons = ref([]);
const difficultyText = {
  basic: "基础",
  advanced: "进阶",
  expert: "专家"
};

async function fetchCourse() {
  loading.value = true;
  try {
    const res = await getPublishedCourse(route.params.id);
    course.value = res.data?.course || res.data || null;
    chapters.value = res.data?.chapters || [];
    lessons.value = res.data?.lessons || [];
  } finally {
    loading.value = false;
  }
}

function chapterLessons(chapterId) {
  return lessons.value.filter((item) => item.chapterId === chapterId);
}

function formatDuration(seconds = 0) {
  if (!seconds) return "未设置时长";
  const minutes = Math.floor(seconds / 60);
  const rest = seconds % 60;
  return `${minutes}分${rest}秒`;
}

onMounted(fetchCourse);
</script>
