<template>
  <PortalLayout>
    <section class="page-heading course-heading">
      <div>
        <h1>专题课程</h1>
        <p>按学段、障碍类型、课程分类和难度筛选在线课程，查看课时、作业和学习进度。</p>
      </div>
      <a-input-search v-model="query.keyword" placeholder="搜索课程标题、教师、教学目标" search-button @search="searchCourses" />
    </section>

    <section class="filter-panel">
      <div class="filter-row" v-for="group in filterGroups" :key="group.field">
        <span>{{ group.label }}</span>
        <div class="filter-options">
          <button class="filter-chip" :class="{ active: !query[group.field] }" type="button" @click="setFilter(group.field, undefined)">全部</button>
          <button
            v-for="item in group.options"
            :key="item.value"
            class="filter-chip"
            :class="{ active: query[group.field] === item.value }"
            type="button"
            @click="setFilter(group.field, item.value)"
          >
            {{ item.label }}
          </button>
        </div>
      </div>
    </section>

    <a-spin :loading="loading" style="width: 100%">
      <div v-if="courses.length" class="course-grid">
        <router-link v-for="(item, index) in courses" :key="item.id" :to="`/courses/${item.id}`" class="course-card">
          <div class="course-cover">
            <img :src="cardCover(item, 'course', index)" :alt="item.title" />
          </div>
          <div class="course-body">
            <strong>{{ item.title }}</strong>
            <p>{{ item.summary || "暂无简介" }}</p>
            <div class="meta-row">
              <a-tag v-if="categoryName(item.stageCategoryId, 'stage')">{{ categoryName(item.stageCategoryId, 'stage') }}</a-tag>
              <a-tag v-if="categoryName(item.disabilityTypeId, 'disability')">{{ categoryName(item.disabilityTypeId, 'disability') }}</a-tag>
              <a-tag v-if="item.category" color="arcoblue">{{ item.category }}</a-tag>
            </div>
            <small>{{ item.teacherName || "平台课程" }} · {{ item.learnerCount || 0 }} 人学习 · {{ item.viewCount || 0 }} 浏览</small>
          </div>
        </router-link>
      </div>
      <a-empty v-else description="暂无课程" />
    </a-spin>

    <div v-if="total > query.pageSize" class="pager">
      <a-pagination :current="query.pageIndex" :page-size="query.pageSize" :total="total" @change="handlePageChange" />
    </div>
  </PortalLayout>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import { useRoute } from "vue-router";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getPublishedCourses } from "@/api/courses";
import { getResourceCategories } from "@/api/resources";
import { cardCover } from "@/utils/defaultCovers";

const route = useRoute();
const loading = ref(false);
const courses = ref([]);
const total = ref(0);
const query = reactive({ keyword: "", stageCategoryId: undefined, disabilityTypeId: undefined, category: undefined, difficulty: undefined, pageIndex: 1, pageSize: 12 });
const categoryOptions = reactive({ stage: [], disability: [] });
const difficultyOptions = [
  { label: "入门", value: "basic" },
  { label: "进阶", value: "advanced" },
  { label: "专家", value: "expert" },
  { label: "高阶", value: "hard" }
];
const courseCategories = ["融合教育", "康复训练", "辅助技术", "IEP 编制", "课程建设"].map((value) => ({ label: value, value }));
const difficultyText = Object.fromEntries(difficultyOptions.map((item) => [item.value, item.label]));

const filterGroups = computed(() => [
  { label: "学段", field: "stageCategoryId", options: categoryOptions.stage.map((item) => ({ label: item.name, value: item.id })) },
  { label: "障碍类型", field: "disabilityTypeId", options: categoryOptions.disability.map((item) => ({ label: item.name, value: item.id })) },
  { label: "课程分类", field: "category", options: courseCategories },
  { label: "难度", field: "difficulty", options: difficultyOptions }
]);

function pagePayload(res) {
  return res.data || {};
}

async function fetchCategories() {
  const res = await getResourceCategories({ pageIndex: 1, pageSize: 1000, status: 1 });
  const list = pagePayload(res).list || pagePayload(res) || [];
  categoryOptions.stage = list.filter((item) => item.type === "stage");
  categoryOptions.disability = list.filter((item) => item.type === "disability");
}

async function fetchCourses() {
  loading.value = true;
  try {
    const res = await getPublishedCourses(query);
    const payload = pagePayload(res);
    courses.value = payload.list || payload || [];
    total.value = payload.count || res.total || courses.value.length;
  } finally {
    loading.value = false;
  }
}

function searchCourses() {
  query.pageIndex = 1;
  fetchCourses();
}

function setFilter(field, value) {
  query[field] = value;
  searchCourses();
}

function handlePageChange(page) {
  query.pageIndex = page;
  fetchCourses();
}

function categoryName(id, type) {
  return (categoryOptions[type] || []).find((item) => item.id === id)?.name || "";
}

onMounted(async () => {
  query.keyword = String(route.query.keyword || "");
  query.difficulty = route.query.difficulty ? String(route.query.difficulty) : undefined;
  await fetchCategories();
  await fetchCourses();
});
</script>

<style scoped>
.course-heading {
  grid-template-columns: minmax(0, 1fr) minmax(280px, 420px);
  align-items: end;
}
.filter-panel {
  display: grid;
  gap: 0;
  margin-bottom: 18px;
  padding: 16px;
  background: #fff;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
  box-shadow: none;
}
.filter-row + .filter-row {
  margin-top: 10px;
}
.filter-row {
  display: grid;
  grid-template-columns: 84px minmax(0, 1fr);
  gap: 14px;
  align-items: start;
}
.filter-row > span {
  padding-top: 7px;
  color: #1f2f4a;
  font-weight: 600;
}
.filter-options {
  display: flex;
  flex-wrap: wrap;
  gap: 0;
  min-width: 0;
}
.filter-chip {
  min-width: 54px;
  height: 30px;
  padding: 0 13px;
  color: #4e5969;
  font-weight: 500;
  letter-spacing: 0;
  white-space: nowrap;
  background: #f4f6f8;
  border: 1px solid transparent;
  border-radius: 0;
  cursor: pointer;
}
.filter-chip + .filter-chip {
  border-left-color: #e5e6eb;
}
.filter-chip:hover {
  color: #0b6be8;
  background: #eef6ff;
}
.filter-chip.active {
  color: #0b6be8;
  background: #eef6ff;
  box-shadow: none;
}
.course-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
}
.course-card {
  overflow: hidden;
  min-height: 330px;
  background: #fff;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
}
.course-cover {
  overflow: hidden;
  display: grid;
  place-items: center;
  height: 150px;
  color: #0b5ed7;
  font-weight: 800;
  background: linear-gradient(135deg, #e8f3ff, #eefbf7);
}
.course-cover img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.course-body {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 14px;
}
.course-body p {
  min-height: 44px;
  margin: 0;
  color: #4e5969;
  line-height: 1.7;
}
.meta-row {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
.course-body small {
  color: #86909c;
}
.pager {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}
@media (max-width: 1100px) {
  .course-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}
@media (max-width: 640px) {
  .course-grid { grid-template-columns: 1fr; }
}
</style>
