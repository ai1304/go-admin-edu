<template>
  <PortalLayout>
    <section class="page-heading case-heading">
      <div>
        <h1>特教案例</h1>
        <p>展示优秀特殊教育案例、IEP 摘要、评估记录与干预策略，为教学实践提供参考。</p>
      </div>
      <a-input-search v-model="query.keyword" placeholder="搜索案例名称、学校、障碍类型" search-button @search="searchCases" />
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
      <div v-if="cases.length" class="case-grid">
        <router-link v-for="item in cases" :key="item.id" :to="`/cases/${item.id}`" class="case-card">
          <div class="case-cover">
            <img v-if="item.coverUrl" :src="item.coverUrl" :alt="item.title" />
            <span v-else>{{ item.caseType || "特教案例" }}</span>
          </div>
          <div class="case-body">
            <strong>{{ item.title }}</strong>
            <p>{{ item.summary || "案例内容已脱敏展示。" }}</p>
            <div class="tag-row">
              <a-tag v-if="item.stage">{{ item.stage }}</a-tag>
              <a-tag v-if="item.disabilityType" color="green">{{ item.disabilityType }}</a-tag>
              <a-tag v-if="item.abilityDomain" color="blue">{{ item.abilityDomain }}</a-tag>
            </div>
            <small>{{ item.school || "平台案例库" }} · {{ item.viewCount || 0 }} 浏览</small>
          </div>
        </router-link>
      </div>
      <a-empty v-else description="暂无案例" />
    </a-spin>
  </PortalLayout>
</template>

<script setup>
import { onMounted, reactive, ref } from "vue";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getPublishedCases } from "@/api/cases";

const loading = ref(false);
const cases = ref([]);
const query = reactive({ keyword: "", stage: "", disabilityType: "", abilityDomain: "", caseType: "", pageIndex: 1, pageSize: 12 });
const filterGroups = [
  { label: "学段", field: "stage", options: ["学前", "小学", "初中", "高中", "高等教育"] },
  { label: "障碍类型", field: "disabilityType", options: ["视障", "听障", "智力障碍", "孤独症", "学习障碍", "肢体障碍"] },
  { label: "能力领域", field: "abilityDomain", options: ["认知", "语言沟通", "社会适应", "运动康复", "生活技能"] },
  { label: "案例类型", field: "caseType", options: ["IEP", "课堂教学", "康复训练", "融合支持", "行为干预"] }
];

async function fetchCases() {
  loading.value = true;
  try {
    const res = await getPublishedCases(query);
    const payload = res.data || {};
    cases.value = payload.list || payload || [];
  } finally {
    loading.value = false;
  }
}

function searchCases() {
  query.pageIndex = 1;
  fetchCases();
}

function setFilter(field, value) {
  query[field] = value;
  searchCases();
}

onMounted(fetchCases);
</script>

<style scoped>
.case-heading {
  grid-template-columns: minmax(0, 1fr) minmax(280px, 420px);
  align-items: end;
}
.filter-panel {
  display: grid;
  gap: 12px;
  margin-bottom: 18px;
  padding: 18px 20px;
  background: #fff;
  border: 1px solid #e6edf7;
  border-radius: 8px;
  box-shadow: 0 12px 28px rgba(27, 84, 150, 0.06);
}
.filter-row {
  display: grid;
  grid-template-columns: 84px minmax(0, 1fr);
  gap: 14px;
  align-items: start;
}
.filter-row > span {
  padding-top: 7px;
  color: #1d3557;
  font-weight: 600;
}
.filter-options {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  min-width: 0;
}
.filter-chip {
  min-width: 56px;
  height: 32px;
  padding: 0 14px;
  color: #344866;
  font-weight: 600;
  letter-spacing: 0;
  white-space: nowrap;
  background: #f4f7fb;
  border: 1px solid transparent;
  border-radius: 999px;
  cursor: pointer;
}
.filter-chip:hover {
  color: #0b6be8;
  background: #edf5ff;
}
.filter-chip.active {
  color: #fff;
  background: #176fd6;
  box-shadow: 0 8px 18px rgba(23, 111, 214, 0.22);
}
.case-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}
.case-card {
  overflow: hidden;
  background: #fff;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
}
.case-cover {
  display: grid;
  place-items: center;
  height: 150px;
  color: #0b5ed7;
  font-weight: 800;
  background: #eef6ff;
}
.case-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.case-body {
  display: grid;
  gap: 10px;
  padding: 14px;
}
.case-body p {
  min-height: 56px;
  margin: 0;
  color: #4e5969;
  line-height: 1.7;
}
.tag-row {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
.case-body small {
  color: #86909c;
}
@media (max-width: 960px) {
  .case-grid { grid-template-columns: 1fr; }
}
</style>
