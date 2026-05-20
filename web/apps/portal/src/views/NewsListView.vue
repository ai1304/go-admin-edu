<template>
  <PortalLayout>
    <section class="page-heading news-heading">
      <div>
        <h1>行业资讯</h1>
        <p>聚合政策法规、学术前沿、行业动态和优秀实践，帮助教师与研究者及时了解特殊教育发展方向。</p>
      </div>
      <a-input-search v-model="query.keyword" placeholder="搜索标题、摘要、来源" search-button @search="searchNews" />
    </section>

    <section class="filter-panel">
      <a-radio-group v-model="query.moduleType" type="button" @change="searchNews">
        <a-radio value="">全部</a-radio>
        <a-radio v-for="item in moduleOptions" :key="item.value" :value="item.value">{{ item.label }}</a-radio>
      </a-radio-group>
    </section>

    <a-spin :loading="loading" style="width: 100%">
      <div v-if="list.length" class="news-page-list">
        <router-link v-for="item in list" :key="item.id" :to="`/news/${item.id}`" class="news-page-item">
          <div v-if="item.coverUrl" class="news-cover">
            <img :src="item.coverUrl" :alt="item.title" />
          </div>
          <div class="news-body">
            <a-tag color="arcoblue">{{ moduleText[item.moduleType] || item.moduleType }}</a-tag>
            <h2>{{ item.title }}</h2>
            <p>{{ item.summary || "暂无摘要" }}</p>
            <footer>
              <span>{{ item.source || "平台发布" }}</span>
              <time>{{ item.publishTime || "未设置时间" }}</time>
              <span>{{ item.viewCount || 0 }} 浏览</span>
              <span>{{ item.likeCount || 0 }} 点赞</span>
            </footer>
          </div>
        </router-link>
      </div>
      <a-empty v-else description="暂无资讯" />
    </a-spin>

    <div v-if="total > query.pageSize" class="pager">
      <a-pagination :current="query.pageIndex" :page-size="query.pageSize" :total="total" @change="handlePageChange" />
    </div>
  </PortalLayout>
</template>

<script setup>
import { onMounted, reactive, ref } from "vue";
import { useRoute } from "vue-router";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getPublishedNews } from "@/api/news";

const route = useRoute();
const loading = ref(false);
const list = ref([]);
const total = ref(0);
const moduleOptions = [
  { label: "政策法规", value: "POLICY" },
  { label: "学术前沿", value: "ACADEMIC" },
  { label: "行业动态", value: "INDUSTRY" },
  { label: "优秀实践", value: "PRACTICE" }
];
const moduleText = Object.fromEntries(moduleOptions.map((item) => [item.value, item.label]));
const query = reactive({
  keyword: String(route.query.keyword || ""),
  moduleType: String(route.query.category || route.query.moduleType || ""),
  pageIndex: 1,
  pageSize: 10
});

function pagePayload(res) {
  return res.data || {};
}

async function fetchNews() {
  loading.value = true;
  try {
    const res = await getPublishedNews(query);
    const payload = pagePayload(res);
    list.value = payload.list || [];
    total.value = payload.count || 0;
  } finally {
    loading.value = false;
  }
}

function searchNews() {
  query.pageIndex = 1;
  fetchNews();
}

function handlePageChange(page) {
  query.pageIndex = page;
  fetchNews();
}

onMounted(fetchNews);
</script>

<style scoped>
.news-heading {
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

.news-page-list {
  display: grid;
  gap: 14px;
}

.news-page-item {
  display: grid;
  grid-template-columns: 210px minmax(0, 1fr);
  gap: 18px;
  padding: 16px;
  background: #fff;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
}

.news-cover {
  height: 132px;
  overflow: hidden;
  border-radius: 8px;
  background: #edf4ff;
}

.news-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.news-body h2 {
  margin: 10px 0 8px;
  color: #1d2129;
  font-size: 20px;
}

.news-body p {
  margin: 0 0 14px;
  color: #4e5969;
  line-height: 1.8;
}

.news-body footer {
  display: flex;
  flex-wrap: wrap;
  gap: 14px;
  color: #86909c;
  font-size: 13px;
}

.pager {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

@media (max-width: 720px) {
  .news-page-item {
    grid-template-columns: 1fr;
  }
}
</style>
