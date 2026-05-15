<template>
  <PortalLayout>
    <a-spin :loading="loading" style="width: 100%">
      <template v-if="resource">
        <section class="detail-hero">
          <img v-if="resource.coverUrl" class="detail-cover" :src="resource.coverUrl" :alt="resource.title" />
          <a-breadcrumb>
            <a-breadcrumb-item>
              <router-link to="/resources">资源中心</router-link>
            </a-breadcrumb-item>
            <a-breadcrumb-item>{{ resource.title }}</a-breadcrumb-item>
          </a-breadcrumb>
          <h1>{{ resource.title }}</h1>
          <p>{{ resource.summary || "暂无简介" }}</p>
          <div class="meta-row">
            <a-tag color="blue">{{ resource.authorName || "平台资源" }}</a-tag>
            <a-tag>{{ resource.keywords || "暂无关键词" }}</a-tag>
            <a-tag v-for="tag in resource.tags || []" :key="tag.id" color="arcoblue">{{ tag.name }}</a-tag>
            <a-tag>{{ resource.viewCount || 0 }} 次浏览</a-tag>
            <a-tag>{{ resource.downloadCount || 0 }} 次下载</a-tag>
            <a-tag>{{ resource.favoriteCount || 0 }} 次收藏</a-tag>
          </div>
        </section>

        <section class="detail-layout">
          <article class="detail-panel">
            <h2>资源介绍</h2>
            <p>{{ resource.summary || "暂未填写资源介绍。" }}</p>
            <div class="comment-block">
              <h2>评论</h2>
              <a-form :model="commentForm" layout="vertical" class="comment-form">
                <a-form-item field="nickname" label="昵称">
                  <a-input v-model="commentForm.nickname" placeholder="请输入昵称" />
                </a-form-item>
                <a-form-item field="content" label="评论内容" required>
                  <a-textarea v-model="commentForm.content" :auto-size="{ minRows: 3, maxRows: 5 }" placeholder="说说你的想法" />
                </a-form-item>
                <a-button type="primary" @click="submitComment">发布评论</a-button>
              </a-form>
              <div v-if="comments.length" class="comment-list">
                <section v-for="item in comments" :key="item.id" class="comment-item">
                  <strong>{{ item.nickname || "访客" }}</strong>
                  <p>{{ item.content }}</p>
                </section>
              </div>
              <a-empty v-else description="暂无评论" />
            </div>
          </article>
          <aside class="side-panel">
            <a-button :type="favorited ? 'outline' : 'primary'" long class="side-action" @click="toggleFavorite">
              {{ favorited ? "取消收藏" : "收藏资源" }}
            </a-button>
            <h2>附件</h2>
            <div v-if="files.length" class="file-list">
              <div v-for="file in files" :key="file.id" class="file-item">
                <div>
                  <strong>{{ file.originalName }}</strong>
                  <span>{{ file.contentType || file.ext || "未知类型" }} · {{ formatSize(file.size) }}</span>
                </div>
                <a-button type="primary" size="small" @click="openFile(file)">下载/预览</a-button>
              </div>
            </div>
            <a-empty v-else description="暂无附件" />
          </aside>
        </section>
      </template>
      <a-empty v-else description="暂无资源详情" />
    </a-spin>
  </PortalLayout>
</template>

<script setup>
import { Message } from "@arco-design/web-vue";
import { onMounted, reactive, ref } from "vue";
import { useRoute } from "vue-router";
import PortalLayout from "@/layouts/PortalLayout.vue";
import {
  createResourceComment,
  favoriteResource,
  getPublishedResource,
  getResourceComments,
  getResourceFavoriteState,
  getResourceFileAccessUrl,
  unfavoriteResource
} from "@/api/resources";

const route = useRoute();
const loading = ref(false);
const resource = ref(null);
const files = ref([]);
const comments = ref([]);
const favorited = ref(false);
const commentForm = reactive({ nickname: "", content: "" });

function clientKey() {
  const storageKey = "edu_portal_client_key";
  let value = window.localStorage.getItem(storageKey);
  if (!value) {
    value = `guest-${Date.now()}-${Math.random().toString(16).slice(2)}`;
    window.localStorage.setItem(storageKey, value);
  }
  return value;
}

async function fetchResource() {
  loading.value = true;
  try {
    const res = await getPublishedResource(route.params.id);
    resource.value = res.data?.resource || res.data || null;
    files.value = res.data?.files || [];
    await Promise.all([fetchFavoriteState(), fetchComments()]);
  } finally {
    loading.value = false;
  }
}

async function fetchFavoriteState() {
  const res = await getResourceFavoriteState(route.params.id, { clientKey: clientKey() });
  favorited.value = !!res.data?.favorited;
}

async function fetchComments() {
  const res = await getResourceComments(route.params.id);
  comments.value = res.data || [];
}

function formatSize(size = 0) {
  if (size < 1024) return `${size} B`;
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
  return `${(size / 1024 / 1024).toFixed(1)} MB`;
}

async function openFile(file) {
  try {
    const res = await getResourceFileAccessUrl(route.params.id, file.id);
    const url = res.data?.url || res.url;
    if (!url) {
      Message.warning("暂未获取到文件访问地址");
      return;
    }
    window.open(url, "_blank", "noopener,noreferrer");
  } catch (error) {
    Message.error("获取文件访问地址失败");
  }
}

async function toggleFavorite() {
  const data = { clientKey: clientKey() };
  if (favorited.value) {
    await unfavoriteResource(route.params.id, data);
    favorited.value = false;
    if (resource.value?.favoriteCount > 0) {
      resource.value.favoriteCount -= 1;
    }
    Message.success("已取消收藏");
  } else {
    await favoriteResource(route.params.id, data);
    favorited.value = true;
    if (resource.value) {
      resource.value.favoriteCount = (resource.value.favoriteCount || 0) + 1;
    }
    Message.success("收藏成功");
  }
}

async function submitComment() {
  if (!commentForm.content) {
    Message.warning("请输入评论内容");
    return;
  }
  await createResourceComment(route.params.id, { ...commentForm });
  Message.success("评论成功");
  commentForm.content = "";
  await fetchComments();
}

onMounted(fetchResource);
</script>

<style scoped>
.comment-block {
  margin-top: 28px;
}

.detail-cover {
  width: 100%;
  max-height: 320px;
  object-fit: cover;
  border-radius: 8px;
}

.comment-form {
  margin-bottom: 20px;
}

.comment-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.comment-item {
  border: 1px solid #e5e6eb;
  border-radius: 8px;
  padding: 14px 16px;
}

.comment-item p {
  margin: 8px 0 0;
  color: #4e5969;
  line-height: 1.7;
}

.side-action {
  margin-bottom: 18px;
}
</style>
