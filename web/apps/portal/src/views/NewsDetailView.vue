<template>
  <PortalLayout>
    <a-spin :loading="loading" style="width: 100%">
      <article v-if="news" class="news-detail">
        <router-link class="back-link" to="/news">返回资讯列表</router-link>
        <a-tag color="arcoblue">{{ moduleText[news.moduleType] || news.moduleType }}</a-tag>
        <h1>{{ news.title }}</h1>
        <div class="meta">
          <span>{{ news.source || "平台发布" }}</span>
          <span>{{ news.publishTime || "未设置时间" }}</span>
          <span>{{ news.viewCount || 0 }} 浏览</span>
          <span>{{ news.likeCount || 0 }} 点赞</span>
        </div>
        <img v-if="news.coverUrl" class="cover" :src="news.coverUrl" :alt="news.title" />
        <p class="summary">{{ news.summary }}</p>
        <div v-if="news.content" class="content" v-html="news.content"></div>
        <div v-else class="content">暂无正文内容。</div>
        <a-button type="primary" @click="handleLike">点赞</a-button>
      </article>
      <a-empty v-else description="资讯不存在或已下线" />
    </a-spin>
  </PortalLayout>
</template>

<script setup>
import { Message } from "@arco-design/web-vue";
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getPublishedNewsItem, likeNews } from "@/api/news";

const route = useRoute();
const loading = ref(false);
const news = ref(null);
const moduleText = {
  POLICY: "政策法规",
  ACADEMIC: "学术前沿",
  INDUSTRY: "行业动态",
  PRACTICE: "优秀实践"
};

async function fetchDetail() {
  loading.value = true;
  try {
    const res = await getPublishedNewsItem(route.params.id);
    news.value = res.data || null;
  } finally {
    loading.value = false;
  }
}

async function handleLike() {
  await likeNews(route.params.id);
  if (news.value) {
    news.value.likeCount = (news.value.likeCount || 0) + 1;
  }
  Message.success("已点赞");
}

onMounted(fetchDetail);
</script>

<style scoped>
.news-detail {
  max-width: 880px;
  margin: 32px auto 0;
  padding: 28px;
  background: #fff;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
}

.back-link {
  display: inline-block;
  margin-bottom: 18px;
  color: #176fd6;
}

.news-detail h1 {
  margin: 14px 0;
  font-size: 34px;
  line-height: 1.3;
}

.meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 22px;
  color: #86909c;
}

.cover {
  width: 100%;
  max-height: 420px;
  object-fit: cover;
  border-radius: 8px;
}

.summary {
  color: #4e5969;
  font-size: 17px;
  line-height: 1.8;
}

.content {
  margin: 24px 0;
  color: #1d2129;
  font-size: 16px;
  line-height: 2;
  white-space: normal;
}

.content :deep(p),
.content :deep(ul),
.content :deep(ol) {
  margin: 0 0 14px;
}
</style>
