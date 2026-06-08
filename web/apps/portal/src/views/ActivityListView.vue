<template>
  <PortalLayout>
    <section class="page-heading activity-heading">
      <div>
        <h1>教研活动</h1>
        <p>展示特殊教育教研活动、竞赛成果与经验交流</p>
      </div>
      <a-input-search v-model="query.keyword" placeholder="搜索活动名称、学校、教师" search-button @search="searchActivities" />
    </section>

    <section class="filter-panel">
      <div v-for="group in filterGroups" :key="group.field" class="filter-row">
        <span>{{ group.label }}</span>
        <div class="filter-options">
          <button class="filter-chip" :class="{ active: !query[group.field] }" type="button" @click="setFilter(group.field, '')">全部</button>
          <button
            v-for="option in group.options"
            :key="option"
            class="filter-chip"
            :class="{ active: query[group.field] === option }"
            type="button"
            @click="setFilter(group.field, option)"
          >
            {{ option }}
          </button>
        </div>
      </div>
    </section>

    <a-spin :loading="loading" style="width: 100%">
      <div v-if="activities.length" class="activity-grid">
        <router-link v-for="(item, index) in activities" :key="item.id" :to="`/activities/${item.id}`" class="activity-card">
          <div class="activity-cover">
            <img :src="cardCover(item, 'activity', index)" :alt="item.title" />
            <a-tag v-if="item.awardLevel" class="award" color="red">{{ item.awardLevel }}</a-tag>
          </div>
          <div class="activity-body">
            <strong>{{ item.title }}</strong>
            <p>{{ item.summary || "暂无简介" }}</p>
            <dl>
              <div><dt>学校</dt><dd>{{ item.school || item.organizer || "平台教研中心" }}</dd></div>
              <div><dt>教师</dt><dd>{{ item.teacher || "待公布" }}</dd></div>
              <div><dt>时间</dt><dd>{{ item.startTime || "时间待定" }}</dd></div>
            </dl>
          </div>
        </router-link>
      </div>
      <a-empty v-else description="暂无活动" />
    </a-spin>

    <div v-if="total > query.pageSize" class="pager">
      <a-pagination :current="query.pageIndex" :page-size="query.pageSize" :total="total" @change="handlePageChange" />
    </div>
  </PortalLayout>
</template>

<script setup>
import { onMounted, reactive, ref } from "vue";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getPublishedActivities } from "@/api/activities";
import { cardCover } from "@/utils/defaultCovers";

const loading = ref(false);
const activities = ref([]);
const total = ref(0);
const query = reactive({ keyword: "", track: "", schoolType: "", pageIndex: 1, pageSize: 12 });
const filterGroups = [
  { label: "赛道", field: "track", options: ["新工科", "新医科", "新农科", "新文科", "基础课程", "课程思政", "产教融合"] },
  { label: "学校类型", field: "schoolType", options: ["部属高校", "地方高校"] }
];

async function fetchActivities() {
  loading.value = true;
  try {
    const res = await getPublishedActivities(query);
    const payload = res.data || {};
    activities.value = payload.list || payload || [];
    total.value = payload.count || res.total || activities.value.length;
  } finally {
    loading.value = false;
  }
}

function searchActivities() {
  query.pageIndex = 1;
  fetchActivities();
}

function setFilter(field, value) {
  query[field] = value;
  searchActivities();
}

function handlePageChange(page) {
  query.pageIndex = page;
  fetchActivities();
}

onMounted(fetchActivities);
</script>

<style scoped>
.activity-heading {
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
.activity-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}
.activity-card {
  overflow: hidden;
  display: grid;
  grid-template-rows: 118px minmax(0, 1fr);
  min-height: 330px;
  background: #fff;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
}
.activity-cover {
  position: relative;
  min-height: 0;
  overflow: hidden;
  display: grid;
  place-items: center;
  height: 118px;
  color: #0b5ed7;
  font-weight: 800;
  background: #e8f3ff;
}
.activity-cover img {
  display: block;
  width: 100%;
  height: 100%;
  max-width: 100%;
  max-height: 100%;
  object-fit: cover;
}
.award {
  position: absolute;
  top: 10px;
  right: 10px;
}
.activity-body {
  min-height: 0;
  padding: 14px;
}
.activity-body strong {
  display: -webkit-box;
  overflow: hidden;
  font-size: 17px;
  line-height: 1.45;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}
.activity-body p {
  display: -webkit-box;
  overflow: hidden;
  min-height: 44px;
  color: #4e5969;
  line-height: 1.7;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}
.activity-body dl {
  display: grid;
  gap: 6px;
  margin: 0 0 12px;
}
.activity-body dl div {
  display: grid;
  grid-template-columns: 42px minmax(0, 1fr);
  gap: 8px;
}
dt {
  color: #86909c;
}
dd {
  margin: 0;
}
.tag-row {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
.pager {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}
@media (max-width: 960px) {
  .activity-grid { grid-template-columns: 1fr; }
}
</style>
