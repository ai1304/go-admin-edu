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

            <section class="comment-block">
              <div class="section-title">
                <h2>评论</h2>
                <span>{{ comments.length }} 条</span>
              </div>
              <a-form :model="commentForm" layout="vertical" class="comment-form">
                <div v-if="replyTarget" class="reply-target">
                  正在回复 {{ replyTarget.nickname || "访客" }}
                  <a-button type="text" size="small" @click="cancelReply">取消</a-button>
                </div>
                <a-form-item field="nickname" label="昵称">
                  <a-input v-model="commentForm.nickname" placeholder="请输入昵称" />
                </a-form-item>
                <a-form-item field="content" label="评论内容" required>
                  <a-textarea v-model="commentForm.content" :auto-size="{ minRows: 3, maxRows: 5 }" placeholder="说说你的想法" />
                </a-form-item>
                <a-button type="primary" @click="submitComment">{{ replyTarget ? "发布回复" : "发布评论" }}</a-button>
              </a-form>

              <div v-if="commentTree.length" class="comment-list">
                <section v-for="item in commentTree" :key="item.id" class="comment-item">
                  <div class="comment-main">
                    <strong>{{ item.nickname || "访客" }}</strong>
                    <p>{{ item.content }}</p>
                    <div class="comment-actions">
                      <a-button type="text" size="mini" @click="likeComment(item)">点赞 {{ item.likeCount || 0 }}</a-button>
                      <a-button type="text" size="mini" @click="replyTo(item)">回复</a-button>
                    </div>
                  </div>
                  <div v-if="item.replies?.length" class="reply-list">
                    <section v-for="reply in item.replies" :key="reply.id" class="reply-item">
                      <strong>{{ reply.nickname || "访客" }}</strong>
                      <p>{{ reply.content }}</p>
                      <div class="comment-actions">
                        <a-button type="text" size="mini" @click="likeComment(reply)">点赞 {{ reply.likeCount || 0 }}</a-button>
                        <a-button type="text" size="mini" @click="replyTo(item)">回复</a-button>
                      </div>
                    </section>
                  </div>
                </section>
              </div>
              <a-empty v-else description="暂无评论" />
            </section>
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
import { computed, onMounted, reactive, ref } from "vue";
import { useRoute } from "vue-router";
import PortalLayout from "@/layouts/PortalLayout.vue";
import {
  createResourceComment,
  favoriteResource,
  getPublishedResource,
  getResourceComments,
  getResourceFavoriteState,
  getResourceFileAccessUrl,
  likeResourceComment,
  unfavoriteResource
} from "@/api/resources";

const route = useRoute();
const loading = ref(false);
const resource = ref(null);
const files = ref([]);
const comments = ref([]);
const favorited = ref(false);
const replyTarget = ref(null);
const commentForm = reactive({ nickname: "", content: "", parentId: 0 });

const commentTree = computed(() => {
  const nodes = new Map();
  const roots = [];
  [...comments.value]
    .sort((a, b) => (a.id || 0) - (b.id || 0))
    .forEach((item) => {
      nodes.set(item.id, { ...item, replies: [] });
    });
  nodes.forEach((item) => {
    if (item.parentId && nodes.has(item.parentId)) {
      nodes.get(item.parentId).replies.push(item);
      return;
    }
    roots.push(item);
  });
  return roots.sort((a, b) => (b.id || 0) - (a.id || 0));
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
    return;
  }
  await favoriteResource(route.params.id, data);
  favorited.value = true;
  if (resource.value) {
    resource.value.favoriteCount = (resource.value.favoriteCount || 0) + 1;
  }
  Message.success("收藏成功");
}

function replyTo(comment) {
  replyTarget.value = comment;
  commentForm.parentId = comment.id;
}

function cancelReply() {
  replyTarget.value = null;
  commentForm.parentId = 0;
}

async function submitComment() {
  if (!commentForm.content.trim()) {
    Message.warning("请输入评论内容");
    return;
  }
  await createResourceComment(route.params.id, { ...commentForm, content: commentForm.content.trim() });
  Message.success(replyTarget.value ? "回复成功" : "评论成功");
  commentForm.content = "";
  cancelReply();
  await fetchComments();
}

async function likeComment(comment) {
  await likeResourceComment(route.params.id, comment.id);
  comment.likeCount = (comment.likeCount || 0) + 1;
}

onMounted(fetchResource);
</script>

<style scoped>
.detail-cover {
  width: 100%;
  max-height: 320px;
  object-fit: cover;
  border-radius: 8px;
}

.comment-block {
  margin-top: 28px;
}

.section-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.section-title span {
  color: #86909c;
  font-size: 13px;
}

.comment-form {
  margin-bottom: 20px;
}

.reply-target {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
  padding: 8px 12px;
  border: 1px solid #bedaff;
  border-radius: 8px;
  color: #165dff;
  background: #f2f7ff;
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
  background: #fff;
}

.comment-main p,
.reply-item p {
  margin: 8px 0 0;
  color: #4e5969;
  line-height: 1.7;
  word-break: break-word;
}

.comment-actions {
  display: flex;
  gap: 4px;
  margin-top: 6px;
}

.reply-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 12px;
  padding-left: 14px;
  border-left: 2px solid #e5e6eb;
}

.reply-item {
  padding: 10px 12px;
  border-radius: 8px;
  background: #f7f8fa;
}

.side-action {
  margin-bottom: 18px;
}
</style>
