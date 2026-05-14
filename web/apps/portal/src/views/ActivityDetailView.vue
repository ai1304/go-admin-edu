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
            <div class="outline-list">
              <h2>活动成果</h2>
              <div v-if="outcomes.length" class="outline-chapters">
                <section v-for="item in outcomes" :key="item.id" class="outline-chapter">
                  <strong>{{ item.title }}</strong>
                  <span>{{ item.content || "暂无成果说明" }}</span>
                </section>
              </div>
              <a-empty v-else description="暂无活动成果" />
            </div>
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
            <a-button type="primary" long class="side-action" @click="signupVisible = true">我要报名</a-button>
          </aside>
        </section>
      </template>
      <a-empty v-else description="暂无活动详情" />
    </a-spin>
    <a-modal v-model:visible="signupVisible" title="活动报名" width="460px" @before-ok="handleSignup">
      <a-form :model="signupForm" layout="vertical">
        <a-form-item field="name" label="姓名" required>
          <a-input v-model="signupForm.name" placeholder="请输入姓名" />
        </a-form-item>
        <a-form-item field="phone" label="电话">
          <a-input v-model="signupForm.phone" placeholder="请输入联系电话" />
        </a-form-item>
      </a-form>
    </a-modal>
  </PortalLayout>
</template>

<script setup>
import { Message } from "@arco-design/web-vue";
import { onMounted, reactive, ref } from "vue";
import { useRoute } from "vue-router";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getPublishedActivity, signupActivity } from "@/api/activities";

const route = useRoute();
const loading = ref(false);
const activity = ref(null);
const outcomes = ref([]);
const signupVisible = ref(false);
const signupForm = reactive({ name: "", phone: "" });

async function fetchActivity() {
  loading.value = true;
  try {
    const res = await getPublishedActivity(route.params.id);
    activity.value = res.data?.activity || res.data || null;
    outcomes.value = res.data?.outcomes || [];
  } finally {
    loading.value = false;
  }
}

async function handleSignup() {
  if (!signupForm.name) {
    Message.warning("请输入姓名");
    return false;
  }
  await signupActivity(route.params.id, { ...signupForm });
  Message.success("报名成功");
  signupVisible.value = false;
  signupForm.name = "";
  signupForm.phone = "";
  fetchActivity();
}

onMounted(fetchActivity);
</script>
