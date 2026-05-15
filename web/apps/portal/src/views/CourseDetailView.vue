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
                    <div v-for="lesson in chapterLessons(chapter.id)" :key="lesson.id" class="lesson-row">
                      <span>{{ lesson.title }} · {{ formatDuration(lesson.durationSeconds) }}</span>
                      <a-space>
                        <a-tag :color="lessonProgress(lesson.id) >= 100 ? 'green' : 'orange'">{{ lessonProgress(lesson.id) }}%</a-tag>
                        <a-button size="mini" type="primary" @click="markLessonFinished(lesson)">标记完成</a-button>
                      </a-space>
                    </div>
                  </div>
                  <small v-else>暂无课时</small>
                </section>
              </div>
              <a-empty v-else description="暂无课程大纲" />
            </div>
            <div class="outline-list">
              <h2>课程作业</h2>
              <div v-if="assignments.length" class="outline-chapters">
                <section v-for="item in assignments" :key="item.id" class="outline-chapter">
                  <strong>{{ item.title }}</strong>
                  <span>{{ item.content || "暂无作业说明" }}</span>
                  <a-button type="primary" size="small" @click="openAssignment(item)">提交作业</a-button>
                </section>
              </div>
              <a-empty v-else description="暂无课程作业" />
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
            <a-progress :percent="courseProgress / 100" :show-text="true" />
          </aside>
        </section>
      </template>
      <a-empty v-else description="暂无课程详情" />
    </a-spin>
    <a-modal v-model:visible="assignmentVisible" :title="currentAssignment?.title || '提交作业'" width="560px" @before-ok="submitAssignment">
      <a-form :model="assignmentForm" layout="vertical">
        <a-form-item field="nickname" label="昵称">
          <a-input v-model="assignmentForm.nickname" placeholder="请输入昵称" />
        </a-form-item>
        <a-form-item field="content" label="提交内容" required>
          <a-textarea v-model="assignmentForm.content" :auto-size="{ minRows: 4, maxRows: 8 }" placeholder="请输入作业内容" />
        </a-form-item>
      </a-form>
    </a-modal>
  </PortalLayout>
</template>

<script setup>
import { Message } from "@arco-design/web-vue";
import { computed, onMounted, reactive, ref } from "vue";
import { useRoute } from "vue-router";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getCourseLearningRecords, getPublishedCourse, submitCourseAssignment, trackCourseLesson } from "@/api/courses";

const route = useRoute();
const loading = ref(false);
const course = ref(null);
const chapters = ref([]);
const lessons = ref([]);
const assignments = ref([]);
const learningRecords = ref([]);
const assignmentVisible = ref(false);
const currentAssignment = ref(null);
const assignmentForm = reactive({ nickname: "", content: "" });
const difficultyText = {
  basic: "基础",
  advanced: "进阶",
  expert: "专家"
};
const courseProgress = computed(() => {
  if (!lessons.value.length) return 0;
  const total = lessons.value.reduce((sum, item) => sum + lessonProgress(item.id), 0);
  return Math.round(total / lessons.value.length);
});

function clientKey() {
  const storageKey = "edu_portal_client_key";
  let value = window.localStorage.getItem(storageKey);
  if (!value) {
    value = `guest-${Date.now()}-${Math.random().toString(16).slice(2)}`;
    window.localStorage.setItem(storageKey, value);
  }
  return value;
}

async function fetchCourse() {
  loading.value = true;
  try {
    const res = await getPublishedCourse(route.params.id);
    course.value = res.data?.course || res.data || null;
    chapters.value = res.data?.chapters || [];
    lessons.value = res.data?.lessons || [];
    assignments.value = res.data?.assignments || [];
    await fetchLearningRecords();
  } finally {
    loading.value = false;
  }
}

async function fetchLearningRecords() {
  const res = await getCourseLearningRecords(route.params.id, { clientKey: clientKey() });
  learningRecords.value = res.data || [];
}

function chapterLessons(chapterId) {
  return lessons.value.filter((item) => item.chapterId === chapterId);
}

function lessonProgress(lessonId) {
  return learningRecords.value.find((item) => item.lessonId === lessonId)?.progress || 0;
}

async function markLessonFinished(lesson) {
  const res = await trackCourseLesson(route.params.id, lesson.id, {
    clientKey: clientKey(),
    progress: 100,
    status: "finished"
  });
  const record = res.data;
  const index = learningRecords.value.findIndex((item) => item.lessonId === lesson.id);
  if (index >= 0) {
    learningRecords.value[index] = record;
  } else {
    learningRecords.value.push(record);
  }
  Message.success("学习进度已更新");
}

function openAssignment(item) {
  currentAssignment.value = item;
  assignmentForm.content = "";
  assignmentVisible.value = true;
}

async function submitAssignment() {
  if (!assignmentForm.content) {
    Message.warning("请输入作业内容");
    return false;
  }
  await submitCourseAssignment(route.params.id, currentAssignment.value.id, {
    clientKey: clientKey(),
    ...assignmentForm
  });
  Message.success("作业提交成功");
  assignmentVisible.value = false;
}

function formatDuration(seconds = 0) {
  if (!seconds) return "未设置时长";
  const minutes = Math.floor(seconds / 60);
  const rest = seconds % 60;
  return `${minutes}分${rest}秒`;
}

onMounted(fetchCourse);
</script>

<style scoped>
.lesson-row {
  align-items: center;
  display: flex;
  justify-content: space-between;
  gap: 12px;
}
</style>
