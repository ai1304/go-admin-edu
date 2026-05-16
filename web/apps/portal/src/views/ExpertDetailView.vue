<template>
  <PortalLayout>
    <a-spin :loading="loading" style="width: 100%">
      <template v-if="expert">
        <section class="detail-hero">
          <a-breadcrumb>
            <a-breadcrumb-item>
              <router-link to="/experts">名师资源</router-link>
            </a-breadcrumb-item>
            <a-breadcrumb-item>{{ expert.name }}</a-breadcrumb-item>
          </a-breadcrumb>
          <h1>{{ expert.name }}</h1>
          <p>{{ expert.introduction || expert.specialties || "暂无专家简介" }}</p>
          <div class="meta-row">
            <a-tag color="purple">{{ expert.title || "专家" }}</a-tag>
            <a-tag>{{ expert.organization || "平台专家库" }}</a-tag>
            <a-tag>{{ expert.favoriteCount || 0 }} 次收藏</a-tag>
            <a-tag>{{ expert.shareCount || 0 }} 次分享</a-tag>
          </div>
        </section>

        <section class="detail-layout">
          <article class="detail-panel">
            <h2>专家介绍</h2>
            <p>{{ expert.introduction || "暂未填写专家介绍。" }}</p>
            <div class="outline-list">
              <h2>相关资源</h2>
              <div v-if="resources.length" class="outline-chapters">
                <component
                  :is="resourceTarget(item) ? 'router-link' : 'button'"
                  v-for="item in resources"
                  :key="item.id"
                  :to="resourceTarget(item)"
                  class="outline-chapter"
                  @click="handleResourceClick(item)"
                >
                  <strong>{{ item.title }}</strong>
                  <span>{{ resourceTypeText[item.type] || "资源" }}</span>
                </component>
              </div>
              <a-empty v-else description="暂无关联资源" />
            </div>
          </article>
          <aside class="side-panel">
            <a-button :type="favorited ? 'outline' : 'primary'" long class="side-action" @click="toggleFavorite">
              {{ favorited ? "取消收藏" : "收藏专家" }}
            </a-button>
            <a-button long class="side-action" @click="handleShare">分享专家主页</a-button>
            <h2>专业方向</h2>
            <p>{{ expert.specialties || "暂未设置专业方向。" }}</p>
          </aside>
        </section>
      </template>
      <a-empty v-else description="暂无专家详情" />
    </a-spin>
  </PortalLayout>
</template>

<script setup>
import { Message } from "@arco-design/web-vue";
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import PortalLayout from "@/layouts/PortalLayout.vue";
import {
  favoriteExpert,
  getExpertFavoriteState,
  getExpertResourceAccessUrl,
  getPublishedExpert,
  shareExpert,
  unfavoriteExpert
} from "@/api/experts";

const resourceTypeText = {
  resource: "资源",
  course: "课程",
  lecture: "讲座",
  file: "文件"
};
const route = useRoute();
const loading = ref(false);
const expert = ref(null);
const resources = ref([]);
const favorited = ref(false);

function clientKey() {
  const storageKey = "edu_portal_client_key";
  let value = window.localStorage.getItem(storageKey);
  if (!value) {
    value = `guest-${Date.now()}-${Math.random().toString(16).slice(2)}`;
    window.localStorage.setItem(storageKey, value);
  }
  return value;
}

function resourceTarget(item) {
  if (item.type === "resource" && item.resourceId) {
    return `/resources/${item.resourceId}`;
  }
  if ((item.type === "course" || item.type === "lecture") && item.courseId) {
    return `/courses/${item.courseId}`;
  }
  return "";
}

async function handleResourceClick(item) {
  if (resourceTarget(item)) return;
  if (item.type !== "file" || !item.fileId) {
    Message.info("该资源暂未配置可访问内容");
    return;
  }
  const res = await getExpertResourceAccessUrl(route.params.id, item.id);
  const url = res.data?.url || res.url;
  if (!url) {
    Message.warning("暂未获取到文件访问地址");
    return;
  }
  window.open(url, "_blank", "noopener,noreferrer");
}

async function fetchFavoriteState() {
  const res = await getExpertFavoriteState(route.params.id, { clientKey: clientKey() });
  favorited.value = !!res.data?.favorited;
}

async function toggleFavorite() {
  const data = { clientKey: clientKey() };
  if (favorited.value) {
    await unfavoriteExpert(route.params.id, data);
    favorited.value = false;
    if (expert.value?.favoriteCount > 0) {
      expert.value.favoriteCount -= 1;
    }
    Message.success("已取消收藏");
    return;
  }
  await favoriteExpert(route.params.id, data);
  favorited.value = true;
  if (expert.value) {
    expert.value.favoriteCount = (expert.value.favoriteCount || 0) + 1;
  }
  Message.success("收藏成功");
}

async function handleShare() {
  const shareUrl = window.location.href;
  if (navigator.share) {
    await navigator.share({ title: expert.value?.name || "专家主页", text: expert.value?.specialties || "", url: shareUrl });
  } else if (navigator.clipboard) {
    await navigator.clipboard.writeText(shareUrl);
    Message.success("分享链接已复制");
  } else {
    Message.info(shareUrl);
  }
  await shareExpert(route.params.id);
  if (expert.value) {
    expert.value.shareCount = (expert.value.shareCount || 0) + 1;
  }
}

async function fetchExpert() {
  loading.value = true;
  try {
    const res = await getPublishedExpert(route.params.id);
    expert.value = res.data?.expert || res.data || null;
    resources.value = res.data?.resources || [];
    await fetchFavoriteState();
  } finally {
    loading.value = false;
  }
}

onMounted(fetchExpert);
</script>

<style scoped>
.outline-chapter {
  width: 100%;
  border: 0;
  text-align: left;
  cursor: pointer;
}

.side-action {
  margin-bottom: 12px;
}
</style>
