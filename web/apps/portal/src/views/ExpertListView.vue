<template>
  <PortalLayout>
    <section class="page-heading expert-heading">
      <div>
        <h1>名师资源</h1>
        <p>汇聚特殊教育领域名师与专家，提供专业讲座、课程与实践指导。</p>
      </div>
      <a-input-search v-model="query.keyword" placeholder="搜索专家姓名、单位、领域" search-button @search="searchExperts" />
    </section>

    <section class="filter-panel">
      <a-form :model="query" layout="inline">
        <a-form-item label="职称">
          <a-select v-model="query.title" allow-clear placeholder="全部职称" style="width: 150px" @change="searchExperts">
            <a-option v-for="item in titleOptions" :key="item" :value="item">{{ item }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="擅长领域">
          <a-select v-model="query.specialty" allow-clear placeholder="全部方向" style="width: 180px" @change="searchExperts">
            <a-option v-for="item in specialtyOptions" :key="item" :value="item">{{ item }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-checkbox v-model="recommendedOnly" @change="searchExperts">只看推荐</a-checkbox>
        </a-form-item>
      </a-form>
    </section>

    <a-spin :loading="loading" style="width: 100%">
      <div v-if="experts.length" class="expert-grid">
        <router-link v-for="item in experts" :key="item.id" :to="`/experts/${item.id}`" class="expert-card">
          <div class="expert-avatar">
            <img v-if="item.avatarUrl" :src="item.avatarUrl" :alt="item.name" />
            <span v-else>{{ (item.name || "名").slice(0, 1) }}</span>
          </div>
          <div class="expert-body">
            <strong>{{ item.name }}</strong>
            <span>{{ item.title || "专家" }} · {{ item.organization || "平台专家库" }}</span>
            <p>{{ item.introduction || "暂无简介" }}</p>
          </div>
        </router-link>
      </div>
      <a-empty v-else description="暂无专家资源" />
    </a-spin>

    <div v-if="total > query.pageSize" class="pager">
      <a-pagination :current="query.pageIndex" :page-size="query.pageSize" :total="total" @change="handlePageChange" />
    </div>
  </PortalLayout>
</template>

<script setup>
import { onMounted, reactive, ref } from "vue";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getPublishedExperts } from "@/api/experts";

const loading = ref(false);
const experts = ref([]);
const total = ref(0);
const recommendedOnly = ref(false);
const titleOptions = ["教授", "副教授", "讲师", "研究员", "高级教师", "正高", "副高"];
const specialtyOptions = ["融合教育", "辅助技术", "康复训练", "个别化教育", "课程建设", "无障碍环境"];
const query = reactive({ keyword: "", title: undefined, specialty: undefined, isRecommended: 0, pageIndex: 1, pageSize: 12 });

async function fetchExperts() {
  loading.value = true;
  try {
    query.isRecommended = recommendedOnly.value ? 1 : 0;
    const res = await getPublishedExperts(query);
    const payload = res.data || {};
    experts.value = payload.list || payload || [];
    total.value = payload.count || res.total || experts.value.length;
  } finally {
    loading.value = false;
  }
}

function searchExperts() {
  query.pageIndex = 1;
  fetchExperts();
}

function handlePageChange(page) {
  query.pageIndex = page;
  fetchExperts();
}

onMounted(fetchExperts);
</script>

<style scoped>
.expert-heading {
  grid-template-columns: minmax(0, 1fr) minmax(280px, 420px);
  align-items: end;
}
.filter-panel {
  margin-bottom: 18px;
  padding: 16px;
  background: #fff;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
}
.expert-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
}
.expert-card {
  padding: 20px 16px;
  text-align: center;
  background: #fff;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
}
.expert-avatar {
  display: grid;
  place-items: center;
  width: 72px;
  height: 72px;
  margin: 0 auto 12px;
  color: #fff;
  font-size: 30px;
  font-weight: 800;
  background: linear-gradient(135deg, #176fd6, #12b886);
  border-radius: 50%;
  overflow: hidden;
}
.expert-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.expert-body {
  display: grid;
  gap: 8px;
}
.expert-body strong {
  display: flex;
  gap: 6px;
  align-items: center;
  justify-content: center;
  font-size: 17px;
}
.expert-body span {
  color: #86909c;
}
.expert-body p {
  min-height: 42px;
  margin: 0;
  color: #4e5969;
  line-height: 1.7;
  display: -webkit-box;
  overflow: hidden;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}
.pager {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}
@media (max-width: 1100px) {
  .expert-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}
@media (max-width: 640px) {
  .expert-grid { grid-template-columns: 1fr; }
}
</style>
