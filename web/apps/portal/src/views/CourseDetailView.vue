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
            <section v-if="currentLesson" class="video-panel">
              <div class="section-title">
                <h2>正在学习：{{ currentLesson.title }}</h2>
                <a-tag :color="lessonProgress(currentLesson.id) >= 100 ? 'green' : 'orange'">
                  {{ lessonProgress(currentLesson.id) }}%
                </a-tag>
              </div>
              <video
                v-if="videoUrl"
                class="lesson-video"
                :src="videoUrl"
                controls
                @timeupdate="handleVideoTimeUpdate"
                @ended="handleVideoEnded"
              ></video>
              <a-empty v-else description="该课时暂未配置视频" />
            </section>
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
                        <a-button size="mini" @click="openLesson(lesson)">学习</a-button>
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
        <a-form-item field="content" label="提交内容">
          <a-textarea v-model="assignmentForm.content" :auto-size="{ minRows: 4, maxRows: 8 }" placeholder="请输入作业内容" />
        </a-form-item>
        <a-form-item field="fileId" label="附件文件 ID">
          <a-input-number v-model="assignmentForm.fileId" :min="0" placeholder="已上传附件的文件 ID，可选" style="width: 100%" />
        </a-form-item>
        <a-form-item label="上传附件">
          <input type="file" :disabled="uploadingAssignmentFile" @change="handleAssignmentFileChange" />
          <p v-if="uploadedAssignmentFile" class="upload-tip">
            已上传：{{ uploadedAssignmentFile.originalName }}，文件 ID：{{ uploadedAssignmentFile.id }}
          </p>
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
import {
  getCourseLearningRecords,
  getCourseLessonVideoUrl,
  getPublishedCourse,
  submitCourseAssignment,
  trackCourseLesson,
  uploadCourseAssignmentFile
} from "@/api/courses";

const route = useRoute();
const loading = ref(false);
const course = ref(null);
const chapters = ref([]);
const lessons = ref([]);
const assignments = ref([]);
const learningRecords = ref([]);
const currentLesson = ref(null);
const videoUrl = ref("");
const lastTrackedSecond = ref(0);
const assignmentVisible = ref(false);
const currentAssignment = ref(null);
const assignmentForm = reactive({ nickname: "", content: "", fileId: undefined });
const uploadedAssignmentFile = ref(null);
const uploadingAssignmentFile = ref(false);
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

function upsertLearningRecord(record) {
  const index = learningRecords.value.findIndex((item) => item.lessonId === record.lessonId);
  if (index >= 0) {
    learningRecords.value[index] = record;
  } else {
    learningRecords.value.push(record);
  }
}

async function openLesson(lesson) {
  currentLesson.value = lesson;
  videoUrl.value = "";
  lastTrackedSecond.value = 0;
  if (!lesson.videoFileId) {
    Message.info("该课时暂未配置视频");
    return;
  }
  const res = await getCourseLessonVideoUrl(route.params.id, lesson.id);
  videoUrl.value = res.data?.url || res.url || "";
}

async function saveLessonProgress(lesson, progress, watchedSeconds, status = "learning") {
  const res = await trackCourseLesson(route.params.id, lesson.id, {
    clientKey: clientKey(),
    progress,
    watchedSeconds,
    status
  });
  upsertLearningRecord(res.data);
}

async function handleVideoTimeUpdate(event) {
  if (!currentLesson.value) return;
  const video = event.target;
  const watchedSeconds = Math.floor(video.currentTime || 0);
  if (watchedSeconds - lastTrackedSecond.value < 10) {
    return;
  }
  lastTrackedSecond.value = watchedSeconds;
  const duration = Math.floor(video.duration || currentLesson.value.durationSeconds || 0);
  const progress = duration > 0 ? Math.min(99, Math.round((watchedSeconds / duration) * 100)) : lessonProgress(currentLesson.value.id);
  await saveLessonProgress(currentLesson.value, progress, watchedSeconds, "learning");
}

async function handleVideoEnded(event) {
  if (!currentLesson.value) return;
  const watchedSeconds = Math.floor(event.target.currentTime || currentLesson.value.durationSeconds || 0);
  await saveLessonProgress(currentLesson.value, 100, watchedSeconds, "finished");
  Message.success("课时学习已完成");
}

async function markLessonFinished(lesson) {
  const res = await trackCourseLesson(route.params.id, lesson.id, {
    clientKey: clientKey(),
    progress: 100,
    watchedSeconds: lesson.durationSeconds || lessonProgress(lesson.id),
    status: "finished"
  });
  upsertLearningRecord(res.data);
  Message.success("学习进度已更新");
}

function openAssignment(item) {
  currentAssignment.value = item;
  assignmentForm.content = "";
  assignmentForm.fileId = undefined;
  uploadedAssignmentFile.value = null;
  assignmentVisible.value = true;
}

async function handleAssignmentFileChange(event) {
  const file = event.target.files?.[0];
  if (!file || !currentAssignment.value) return;
  const formData = new FormData();
  formData.append("file", file);
  uploadingAssignmentFile.value = true;
  try {
    const res = await uploadCourseAssignmentFile(route.params.id, currentAssignment.value.id, formData);
    uploadedAssignmentFile.value = res.data || res;
    assignmentForm.fileId = uploadedAssignmentFile.value.id;
    Message.success("附件上传成功");
  } finally {
    uploadingAssignmentFile.value = false;
    event.target.value = "";
  }
}

async function submitAssignment() {
  if (!assignmentForm.content && !assignmentForm.fileId) {
    Message.warning("请输入作业内容或附件文件 ID");
    return false;
  }
  await submitCourseAssignment(route.params.id, currentAssignment.value.id, {
    clientKey: clientKey(),
    ...assignmentForm,
    fileId: assignmentForm.fileId || 0
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

.video-panel {
  margin: 24px 0;
  padding: 16px;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
  background: #fff;
}

.section-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.lesson-video {
  width: 100%;
  max-height: 420px;
  margin-top: 12px;
  border-radius: 8px;
  background: #000;
}

.upload-tip {
  margin: 8px 0 0;
  color: #165dff;
  font-size: 13px;
}
</style>
