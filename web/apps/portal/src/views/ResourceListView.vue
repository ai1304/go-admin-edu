<template>
  <PortalLayout>
    <section class="page-heading resource-heading">
      <div>
        <h1>资源中心</h1>
        <p>按学段、障碍类型、资源类型和能力领域筛选优质特教资源。</p>
      </div>
      <a-input-search v-model="query.keyword" placeholder="搜索标题、简介、关键词、作者" search-button @search="searchResources" />
    </section>

    <section class="filter-panel">
      <a-form :model="query" layout="inline">
        <a-form-item label="学段">
          <a-select v-model="query.stageCategoryId" allow-clear placeholder="全部学段" style="width: 150px" @change="searchResources">
            <a-option v-for="item in categoryOptions.stage" :key="item.id" :value="item.id">{{ item.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="障碍类型">
          <a-select v-model="query.disabilityTypeId" allow-clear placeholder="全部类型" style="width: 160px" @change="searchResources">
            <a-option v-for="item in categoryOptions.disability" :key="item.id" :value="item.id">{{ item.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="资源类型">
          <a-select v-model="query.resourceTypeId" allow-clear placeholder="全部资源" style="width: 160px" @change="searchResources">
            <a-option v-for="item in categoryOptions.resource_type" :key="item.id" :value="item.id">{{ item.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="能力领域">
          <a-select v-model="query.abilityDomainId" allow-clear placeholder="全部领域" style="width: 160px" @change="searchResources">
            <a-option v-for="item in categoryOptions.ability_domain" :key="item.id" :value="item.id">{{ item.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="标签">
          <a-select v-model="query.tagId" allow-clear placeholder="全部标签" style="width: 150px" @change="searchResources">
            <a-option v-for="item in tagOptions" :key="item.id" :value="item.id">{{ item.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="排序">
          <a-select v-model="query.sort" style="width: 140px" @change="searchResources">
            <a-option value="latest">最新发布</a-option>
            <a-option value="view">浏览最多</a-option>
            <a-option value="download">下载最多</a-option>
            <a-option value="favorite">收藏最多</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="searchResources">查询</a-button>
            <a-button @click="resetFilters">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </section>

    <a-spin :loading="loading" style="width: 100%">
      <div v-if="resources.length" class="resource-grid">
        <router-link v-for="item in resources" :key="item.id" :to="`/resources/${item.id}`" class="resource-card">
          <div class="cover" :class="{ empty: !item.coverUrl }">
            <img v-if="item.coverUrl" :src="item.coverUrl" :alt="item.title" />
            <span v-else>{{ categoryName(item.resourceTypeId, 'resource_type') || '资源' }}</span>
          </div>
          <div class="resource-body">
            <strong>{{ item.title }}</strong>
            <span>{{ item.summary || "暂无简介" }}</span>
            <div class="resource-tags">
              <a-tag v-if="categoryName(item.stageCategoryId, 'stage')">{{ categoryName(item.stageCategoryId, 'stage') }}</a-tag>
              <a-tag v-if="categoryName(item.disabilityTypeId, 'disability')">{{ categoryName(item.disabilityTypeId, 'disability') }}</a-tag>
              <a-tag v-if="categoryName(item.resourceTypeId, 'resource_type')" color="blue">{{ categoryName(item.resourceTypeId, 'resource_type') }}</a-tag>
              <a-tag v-for="tag in item.tags || []" :key="tag.id" color="arcoblue">{{ tag.name }}</a-tag>
            </div>
            <small>{{ item.authorName || "平台资源" }} · {{ item.viewCount || 0 }} 浏览 · {{ item.downloadCount || 0 }} 下载</small>
          </div>
        </router-link>
      </div>
      <a-empty v-else description="暂无资源" />
    </a-spin>

    <div v-if="total > query.pageSize" class="pager">
      <a-pagination :current="query.pageIndex" :page-size="query.pageSize" :total="total" @change="handlePageChange" />
    </div>
  </PortalLayout>
</template>

<script setup>
import { onMounted, reactive, ref } from "vue";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getResourceCategories, getResourceTags, searchPublishedResources } from "@/api/resources";

const loading = ref(false);
const resources = ref([]);
const total = ref(0);
const tagOptions = ref([]);
const query = reactive({
  keyword: "",
  tagId: undefined,
  stageCategoryId: undefined,
  disabilityTypeId: undefined,
  resourceTypeId: undefined,
  abilityDomainId: undefined,
  sort: "latest",
  pageIndex: 1,
  pageSize: 12
});
const categoryOptions = reactive({
  stage: [],
  disability: [],
  resource_type: [],
  ability_domain: [],
  topic: []
});

function pagePayload(res) {
  return res.data || {};
}

async function fetchCategories() {
  const res = await getResourceCategories({ pageIndex: 1, pageSize: 1000, status: 1 });
  const list = pagePayload(res).list || pagePayload(res) || [];
  Object.keys(categoryOptions).forEach((key) => {
    categoryOptions[key] = [];
  });
  list.forEach((item) => {
    if (!categoryOptions[item.type]) {
      categoryOptions[item.type] = [];
    }
    categoryOptions[item.type].push(item);
  });
}

async function fetchTags() {
  const res = await getResourceTags({ pageIndex: 1, pageSize: 1000, status: 1 });
  const list = pagePayload(res).list || pagePayload(res) || [];
  tagOptions.value = list;
}

async function fetchResources() {
  loading.value = true;
  try {
    const res = await searchPublishedResources(query);
    const payload = pagePayload(res);
    resources.value = payload.list || payload || [];
    total.value = payload.count || res.total || 0;
  } finally {
    loading.value = false;
  }
}

function searchResources() {
  query.pageIndex = 1;
  fetchResources();
}

function resetFilters() {
  query.keyword = "";
  query.tagId = undefined;
  query.stageCategoryId = undefined;
  query.disabilityTypeId = undefined;
  query.resourceTypeId = undefined;
  query.abilityDomainId = undefined;
  query.sort = "latest";
  query.pageIndex = 1;
  fetchResources();
}

function handlePageChange(page) {
  query.pageIndex = page;
  fetchResources();
}

function categoryName(id, type) {
  return (categoryOptions[type] || []).find((item) => item.id === id)?.name || "";
}

onMounted(async () => {
  await Promise.all([fetchCategories(), fetchTags()]);
  await fetchResources();
});
</script>

<style scoped>
.resource-heading {
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

.resource-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.resource-card {
  overflow: hidden;
  display: grid;
  grid-template-rows: 150px minmax(0, 1fr);
  min-height: 340px;
  background: #fff;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
}

.cover {
  display: grid;
  place-items: center;
  background: #e8f3ff;
  color: #165dff;
  font-weight: 700;
}

.cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.cover.empty span {
  padding: 0 18px;
}

.resource-body {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 16px;
}

.resource-body strong {
  font-size: 17px;
}

.resource-body span {
  color: #4e5969;
  line-height: 1.7;
}

.resource-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.resource-body small {
  margin-top: auto;
  color: #86909c;
}

.pager {
  display: flex;
  justify-content: center;
  margin-top: 22px;
}

@media (max-width: 960px) {
  .resource-heading,
  .resource-grid {
    grid-template-columns: 1fr;
  }
}
</style>
