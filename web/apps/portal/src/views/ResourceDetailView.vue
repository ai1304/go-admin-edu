<template>
  <PortalLayout>
    <a-spin :loading="loading" style="width: 100%">
      <template v-if="resource">
        <section class="resource-player-page">
          <main class="resource-main">
            <div class="resource-titlebar">
              <a-breadcrumb>
                <a-breadcrumb-item>
                  <router-link to="/resources">资源中心</router-link>
                </a-breadcrumb-item>
                <a-breadcrumb-item>{{ resource.title }}</a-breadcrumb-item>
              </a-breadcrumb>
              <h1>{{ resource.title }}</h1>
              <div class="meta-row">
                <a-tag color="blue">{{ resource.authorName || "平台资源" }}</a-tag>
                <a-tag v-for="tag in resource.tags || []" :key="tag.id" color="arcoblue">{{ tag.name }}</a-tag>
                <a-tag>{{ resource.viewCount || 0 }} 次浏览</a-tag>
                <a-tag>{{ resource.downloadCount || 0 }} 次下载</a-tag>
                <a-tag>{{ resource.favoriteCount || 0 }} 次收藏</a-tag>
              </div>
            </div>

            <section class="preview-shell">
              <div v-if="activeFile && activeKind === 'video'" class="video-stage">
                <video v-if="activeUrl" :src="activeUrl" controls playsinline preload="metadata"></video>
                <div v-else class="preview-loading">正在获取视频地址...</div>
              </div>

              <div v-else-if="activeFile && activeKind === 'pdf'" class="pdf-stage">
                <iframe v-if="activeUrl" :src="activeUrl" :title="activeFile.originalName"></iframe>
                <div v-else class="preview-loading">正在获取 PDF 地址...</div>
              </div>

              <div v-else class="empty-stage">
                <div class="empty-mark">{{ previewTypeLabel(activeFile) }}</div>
                <h2>{{ activeFile ? activeFile.originalName : "暂无可预览文件" }}</h2>
                <p>{{ previewHint }}</p>
                <a-space v-if="activeFile">
                  <a-button type="primary" @click="downloadFile(activeFile)">下载文件</a-button>
                  <a-button @click="openFile(activeFile)">新窗口打开</a-button>
                </a-space>
              </div>
            </section>

            <section class="player-actions">
              <a-space wrap>
                <a-button :type="favorited ? 'outline' : 'primary'" @click="toggleFavorite">
                  {{ favorited ? "取消收藏" : "收藏资源" }}
                </a-button>
                <a-button v-if="activeFile" @click="downloadFile(activeFile)">下载当前文件</a-button>
                <a-button v-if="activeFile" @click="openFile(activeFile)">新窗口预览</a-button>
              </a-space>
            </section>

            <section class="detail-panel resource-intro">
              <h2>资源介绍</h2>
              <p>{{ resource.summary || "暂未填写资源介绍。" }}</p>
              <div v-if="resource.keywords" class="keyword-row">
                <a-tag v-for="keyword in keywordList" :key="keyword">{{ keyword }}</a-tag>
              </div>
            </section>

            <section class="detail-panel attachment-panel">
              <div class="section-title">
                <h2>附件({{ downloadableFiles.length }})</h2>
                <span>{{ previewableFiles.length }} 个可在线预览</span>
              </div>
              <div v-if="downloadableFiles.length" class="attachment-grid">
                <button
                  v-for="file in downloadableFiles"
                  :key="file.id"
                  class="attachment-item"
                  :class="{ active: activeFile?.id === file.id }"
                  type="button"
                  @click="selectFile(file)"
                >
                  <span class="file-badge" :class="fileKind(file)">{{ fileTypeText(file) }}</span>
                  <span class="attachment-name">{{ file.originalName }}</span>
                  <small>{{ formatSize(file.size) }}</small>
                  <a-button size="mini" @click.stop="downloadFile(file)">下载</a-button>
                </button>
              </div>
              <a-empty v-else description="暂无附件" />
            </section>

            <section class="detail-panel comment-block">
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
          </main>

          <aside class="resource-sidebar">
            <section class="side-card">
              <h2>资源目录</h2>
              <div v-if="files.length" class="catalog-list">
                <button
                  v-for="(file, index) in files"
                  :key="file.id"
                  class="catalog-item"
                  :class="{ active: activeFile?.id === file.id }"
                  type="button"
                  @click="selectFile(file)"
                >
                  <span class="catalog-index">{{ index + 1 }}</span>
                  <span>
                    <strong>{{ file.originalName }}</strong>
                    <small>{{ fileTypeText(file) }} · {{ formatSize(file.size) }}</small>
                  </span>
                </button>
              </div>
              <a-empty v-else description="暂无目录" />
            </section>

            <section class="side-card">
              <h2>资源信息</h2>
              <dl class="info-list">
                <div>
                  <dt>作者</dt>
                  <dd>{{ resource.authorName || "平台资源" }}</dd>
                </div>
                <div>
                  <dt>文件数</dt>
                  <dd>{{ files.length }} 个</dd>
                </div>
                <div>
                  <dt>可预览</dt>
                  <dd>{{ previewableFiles.length }} 个</dd>
                </div>
              </dl>
            </section>
          </aside>
        </section>
      </template>
      <a-empty v-else description="暂无资源详情" />
    </a-spin>
  </PortalLayout>
</template>

<script setup>
import { Message } from "@arco-design/web-vue";
import { computed, nextTick, onMounted, reactive, ref } from "vue";
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
const activeFile = ref(null);
const activeUrl = ref("");
const urlCache = reactive({});
const commentForm = reactive({ nickname: "", content: "", parentId: 0 });

const videoExts = new Set(["mp4", "webm", "ogg", "mov", "m4v"]);
const pdfExts = new Set(["pdf"]);
const docExts = new Set(["ppt", "pptx", "doc", "docx", "xls", "xlsx"]);

const activeKind = computed(() => fileKind(activeFile.value));
const previewableFiles = computed(() => files.value.filter((file) => ["video", "pdf"].includes(fileKind(file))));
const downloadableFiles = computed(() => files.value.filter((file) => file.usage !== "cover"));
const keywordList = computed(() => (resource.value?.keywords || "").split(/[,，\s]+/).filter(Boolean));
const previewHint = computed(() => {
  if (!activeFile.value) return "后台上传视频或 PDF 后，可在这里直接在线播放或在线预览。";
  if (activeKind.value === "document") return "PPT、Word、Excel 类型暂不做浏览器内嵌预览，请下载后查看。";
  return "当前文件暂不支持在线预览，请下载后查看。";
});

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
    pickInitialFile();
    await Promise.all([fetchFavoriteState(), fetchComments()]);
  } finally {
    loading.value = false;
  }
}

function pickInitialFile() {
  const initial = previewableFiles.value[0] || downloadableFiles.value[0] || null;
  if (initial) {
    nextTick(() => selectFile(initial));
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

function normalizedExt(file = {}) {
  const fromExt = (file.ext || "").toLowerCase().replace(/^\./, "");
  if (fromExt) return fromExt;
  return (file.originalName || "").split(".").pop()?.toLowerCase() || "";
}

function fileKind(file) {
  if (!file) return "empty";
  const ext = normalizedExt(file);
  const contentType = (file.contentType || "").toLowerCase();
  if (contentType.startsWith("video/") || videoExts.has(ext)) return "video";
  if (contentType.includes("pdf") || pdfExts.has(ext)) return "pdf";
  if (docExts.has(ext)) return "document";
  return "other";
}

function fileTypeText(file) {
  const kind = fileKind(file);
  if (kind === "video") return "视频";
  if (kind === "pdf") return "PDF";
  if (kind === "document") return (normalizedExt(file) || "文档").toUpperCase();
  return (normalizedExt(file) || "附件").toUpperCase();
}

function previewTypeLabel(file) {
  if (!file) return "资源";
  return fileTypeText(file);
}

function formatSize(size = 0) {
  if (size < 1024) return `${size} B`;
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
  return `${(size / 1024 / 1024).toFixed(1)} MB`;
}

async function getFileUrl(file, purpose = "preview") {
  const cacheKey = `${file.id}:${purpose}`;
  if (purpose === "preview" && urlCache[cacheKey]) return urlCache[cacheKey];
  const res = await getResourceFileAccessUrl(route.params.id, file.id, { purpose });
  const url = res.data?.url || res.url;
  if (!url) throw new Error("missing file url");
  if (purpose === "preview") {
    urlCache[cacheKey] = url;
  }
  return url;
}

async function selectFile(file) {
  activeFile.value = file;
  activeUrl.value = "";
  if (!["video", "pdf"].includes(fileKind(file))) return;
  try {
    activeUrl.value = await getFileUrl(file, "preview");
  } catch (error) {
    Message.error("获取预览地址失败");
  }
}

async function openFile(file) {
  try {
    const url = await getFileUrl(file, "preview");
    window.open(url, "_blank", "noopener,noreferrer");
  } catch (error) {
    Message.error("获取文件访问地址失败");
  }
}

async function downloadFile(file) {
  try {
    const url = await getFileUrl(file, "download");
    window.open(url, "_blank", "noopener,noreferrer");
    if (resource.value) {
      resource.value.downloadCount = (resource.value.downloadCount || 0) + 1;
    }
  } catch (error) {
    Message.error("获取下载地址失败");
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
.resource-player-page {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 360px;
  gap: 22px;
  align-items: start;
  padding-top: 24px;
}

.resource-main {
  display: grid;
  gap: 16px;
  min-width: 0;
}

.resource-titlebar {
  display: grid;
  gap: 12px;
}

.resource-titlebar h1 {
  margin: 0;
  color: #17233d;
  font-size: 30px;
  line-height: 1.28;
}

.preview-shell {
  min-height: 560px;
  overflow: hidden;
  background: #080b16;
  border: 1px solid rgba(20, 32, 56, 0.16);
  border-radius: 8px;
  box-shadow: 0 18px 45px rgba(24, 48, 92, 0.12);
}

.video-stage,
.pdf-stage,
.empty-stage {
  width: 100%;
  min-height: 560px;
}

.video-stage {
  display: grid;
  place-items: center;
  background:
    linear-gradient(180deg, rgba(6, 9, 22, 0.65), rgba(6, 9, 22, 0.96)),
    radial-gradient(circle at 50% 30%, rgba(49, 120, 255, 0.22), transparent 34%);
}

.video-stage video {
  width: 100%;
  max-height: 560px;
  background: #000;
}

.pdf-stage iframe {
  width: 100%;
  height: 720px;
  border: 0;
  background: #fff;
}

.empty-stage,
.preview-loading {
  display: grid;
  place-items: center;
  align-content: center;
  gap: 16px;
  padding: 48px;
  color: #dbe7ff;
  text-align: center;
}

.empty-stage h2 {
  max-width: 720px;
  margin: 0;
  color: #fff;
  font-size: 24px;
  line-height: 1.45;
  word-break: break-word;
}

.empty-stage p {
  max-width: 560px;
  margin: 0;
  color: #aebbd2;
  line-height: 1.8;
}

.empty-mark {
  display: grid;
  place-items: center;
  width: 72px;
  height: 72px;
  color: #fff;
  font-weight: 800;
  background: rgba(255, 255, 255, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.22);
  border-radius: 999px;
}

.player-actions {
  padding: 14px 0 2px;
}

.resource-intro,
.attachment-panel,
.comment-block {
  box-shadow: 0 12px 28px rgba(35, 84, 150, 0.06);
}

.keyword-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 14px;
}

.attachment-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.attachment-item {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto auto;
  gap: 10px;
  align-items: center;
  min-height: 58px;
  padding: 10px 12px;
  color: #1d2129;
  text-align: left;
  background: #f7f8fa;
  border: 1px solid transparent;
  border-radius: 8px;
  cursor: pointer;
}

.attachment-item:hover,
.attachment-item.active {
  background: #fff7f0;
  border-color: #ff7d00;
}

.attachment-name {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.attachment-item small {
  color: #86909c;
}

.file-badge {
  min-width: 46px;
  padding: 5px 8px;
  color: #fff;
  font-size: 12px;
  font-weight: 700;
  text-align: center;
  background: #4e5969;
  border-radius: 6px;
}

.file-badge.video {
  background: #f53f3f;
}

.file-badge.pdf {
  background: #ff7d00;
}

.file-badge.document {
  background: #165dff;
}

.resource-sidebar {
  position: sticky;
  top: 92px;
  display: grid;
  gap: 16px;
}

.side-card {
  padding: 20px;
  background: #fff;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
  box-shadow: 0 12px 28px rgba(35, 84, 150, 0.06);
}

.side-card h2 {
  margin: 0 0 16px;
  font-size: 20px;
}

.catalog-list {
  display: grid;
  gap: 10px;
  max-height: 640px;
  overflow: auto;
  padding-right: 2px;
}

.catalog-item {
  display: grid;
  grid-template-columns: 28px minmax(0, 1fr);
  gap: 10px;
  align-items: start;
  width: 100%;
  padding: 12px;
  color: #1d2129;
  text-align: left;
  background: #fff;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
  cursor: pointer;
}

.catalog-item:hover,
.catalog-item.active {
  color: #ff5a00;
  background: #fff7f0;
  border-color: #ff7d00;
}

.catalog-index {
  display: grid;
  place-items: center;
  width: 28px;
  height: 28px;
  color: #4e5969;
  font-size: 12px;
  font-weight: 700;
  background: #f2f3f5;
  border-radius: 999px;
}

.catalog-item.active .catalog-index {
  color: #fff;
  background: #ff7d00;
}

.catalog-item strong,
.catalog-item small {
  display: block;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
}

.catalog-item strong {
  white-space: nowrap;
}

.catalog-item small {
  margin-top: 5px;
  color: #86909c;
  font-size: 12px;
  white-space: nowrap;
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
  color: #165dff;
  background: #f2f7ff;
  border: 1px solid #bedaff;
  border-radius: 8px;
}

.comment-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.comment-item {
  padding: 14px 16px;
  background: #fff;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
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
  background: #f7f8fa;
  border-radius: 8px;
}

@media (max-width: 980px) {
  .resource-player-page {
    grid-template-columns: 1fr;
  }

  .resource-sidebar {
    position: static;
  }

  .preview-shell,
  .video-stage,
  .pdf-stage,
  .empty-stage {
    min-height: 360px;
  }

  .pdf-stage iframe {
    height: 520px;
  }

  .attachment-grid {
    grid-template-columns: 1fr;
  }
}
</style>
