<template>
  <PortalLayout>
    <a-spin :loading="loading" style="width: 100%">
      <template v-if="activity">
        <section class="detail-hero">
          <a-breadcrumb>
            <a-breadcrumb-item>
              <router-link to="/activities">教研活动</router-link>
            </a-breadcrumb-item>
            <a-breadcrumb-item>{{ activity.title }}</a-breadcrumb-item>
          </a-breadcrumb>
          <h1>{{ activity.title }}</h1>
          <p>{{ activity.summary || "暂无活动简介" }}</p>
          <div class="meta-row">
            <a-tag color="green">{{ activity.organizer || "平台活动" }}</a-tag>
            <a-tag>{{ activity.startTime || "时间待定" }}</a-tag>
            <a-tag>{{ activity.location || "地点待定" }}</a-tag>
          </div>
        </section>

        <section class="detail-layout">
          <article class="detail-panel">
            <h2>活动介绍</h2>
            <p>{{ activity.summary || "暂未填写活动介绍。" }}</p>
          </article>
          <aside class="side-panel">
            <h2>活动信息</h2>
            <dl class="info-list">
              <div>
                <dt>开始时间</dt>
                <dd>{{ activity.startTime || "待定" }}</dd>
              </div>
              <div>
                <dt>结束时间</dt>
                <dd>{{ activity.endTime || "待定" }}</dd>
              </div>
              <div>
                <dt>报名人数</dt>
                <dd>{{ activity.signupCount || 0 }}</dd>
              </div>
            </dl>
          </aside>
        </section>
      </template>
      <a-empty v-else description="暂无活动详情" />
    </a-spin>
  </PortalLayout>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getPublishedActivity } from "@/api/activities";

const route = useRoute();
const loading = ref(false);
const activity = ref(null);

async function fetchActivity() {
  loading.value = true;
  try {
    const res = await getPublishedActivity(route.params.id);
    activity.value = res.data?.activity || res.data || null;
  } finally {
    loading.value = false;
  }
}

onMounted(fetchActivity);
</script>
