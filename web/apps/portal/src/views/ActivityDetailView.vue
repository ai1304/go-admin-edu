<template>
  <PortalLayout>
    <a-spin :loading="loading" style="width: 100%">
      <template v-if="activity">
        <section class="detail-hero">
          <div class="breadcrumb-line">
            <a-button size="small" @click="router.push('/activities')">返回</a-button>
            <a-breadcrumb>
              <a-breadcrumb-item>
                <router-link to="/activities">教研活动</router-link>
              </a-breadcrumb-item>
              <a-breadcrumb-item>{{ activity.title }}</a-breadcrumb-item>
            </a-breadcrumb>
          </div>
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
              <h2>活动详情</h2>
              <div v-if="outcomes.length" class="outline-chapters">
                <section v-for="item in outcomes" :key="item.id" class="outline-chapter">
                  <strong>{{ item.title }}</strong>
                  <span>{{ item.content || "暂无成果说明" }}</span>
                </section>
              </div>
              <a-empty v-else description="暂无活动详情" />
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
            <a-button v-if="!signed" type="primary" long class="side-action" @click="signupVisible = true">我要报名</a-button>
            <a-button v-else status="warning" long class="side-action" @click="handleCancelSignup">取消报名</a-button>
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
        <a-form-item field="phone" label="电话" required>
          <a-input v-model="signupForm.phone" placeholder="请输入联系电话" />
        </a-form-item>
      </a-form>
    </a-modal>
    <a-modal v-model:visible="outcomeVisible" title="上传活动详情" width="560px" @before-ok="handleSubmitOutcome">
      <a-form :model="outcomeForm" layout="vertical">
        <a-form-item field="title" label="成果标题" required>
          <a-input v-model="outcomeForm.title" placeholder="请输入成果标题" />
        </a-form-item>
        <a-form-item field="content" label="成果说明">
          <a-textarea v-model="outcomeForm.content" :auto-size="{ minRows: 3, maxRows: 6 }" placeholder="请输入成果说明" />
        </a-form-item>
        <a-form-item label="成果附件">
          <input type="file" :disabled="outcomeUploading" @change="handleOutcomeFileChange" />
          <p v-if="uploadedOutcomeFile" class="upload-tip">
            已上传：{{ uploadedOutcomeFile.originalName }}，文件 ID：{{ uploadedOutcomeFile.id }}
          </p>
        </a-form-item>
      </a-form>
    </a-modal>
  </PortalLayout>
</template>

<script setup>
import { Message } from "@arco-design/web-vue";
import { onMounted, reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import PortalLayout from "@/layouts/PortalLayout.vue";
import {
  cancelActivitySignup,
  checkinActivity,
  getActivitySignupState,
  getPublishedActivity,
  signupActivity,
  submitActivityOutcome,
  uploadActivityOutcomeFile
} from "@/api/activities";

const route = useRoute();
const router = useRouter();
const loading = ref(false);
const activity = ref(null);
const outcomes = ref([]);
const signupVisible = ref(false);
const outcomeVisible = ref(false);
const signed = ref(false);
const checked = ref(false);
const uploadedOutcomeFile = ref(null);
const outcomeUploading = ref(false);
const signupForm = reactive({ name: "", phone: "" });
const outcomeForm = reactive({ title: "", content: "", fileId: 0 });

function clientKey() {
  const storageKey = "edu_portal_client_key";
  let value = window.localStorage.getItem(storageKey);
  if (!value) {
    value = `guest-${Date.now()}-${Math.random().toString(16).slice(2)}`;
    window.localStorage.setItem(storageKey, value);
  }
  return value;
}

async function fetchActivity() {
  loading.value = true;
  try {
    const res = await getPublishedActivity(route.params.id);
    activity.value = res.data?.activity || res.data || null;
    outcomes.value = res.data?.outcomes || [];
    await fetchSignupState();
  } finally {
    loading.value = false;
  }
}

async function fetchSignupState() {
  const res = await getActivitySignupState(route.params.id, { clientKey: clientKey() });
  signed.value = !!res.data?.signed;
  checked.value = !!res.data?.checked;
}

async function handleSignup() {
  if (!signupForm.name) {
    Message.warning("请输入姓名");
    return false;
  }
  if (!/^1[3-9]\d{9}$/.test(signupForm.phone.trim())) {
    Message.warning("请输入正确的手机号");
    return false;
  }
  await signupActivity(route.params.id, { ...signupForm, clientKey: clientKey() });
  Message.success("报名成功");
  signupVisible.value = false;
  signupForm.name = "";
  signupForm.phone = "";
  fetchActivity();
}

async function handleCancelSignup() {
  await cancelActivitySignup(route.params.id, { clientKey: clientKey() });
  Message.success("已取消报名");
  signed.value = false;
  checked.value = false;
  fetchActivity();
}

async function handleCheckin() {
  await checkinActivity(route.params.id, { clientKey: clientKey() });
  Message.success("签到成功");
  checked.value = true;
}

async function handleOutcomeFileChange(event) {
  const file = event.target.files?.[0];
  if (!file) return;
  const formData = new FormData();
  formData.append("file", file);
  outcomeUploading.value = true;
  try {
    const res = await uploadActivityOutcomeFile(route.params.id, formData);
    uploadedOutcomeFile.value = res.data || res;
    outcomeForm.fileId = uploadedOutcomeFile.value.id;
    Message.success("成果附件上传成功");
  } finally {
    outcomeUploading.value = false;
    event.target.value = "";
  }
}

async function handleSubmitOutcome() {
  if (!outcomeForm.title) {
    Message.warning("请输入成果标题");
    return false;
  }
  await submitActivityOutcome(route.params.id, { ...outcomeForm, clientKey: clientKey() });
  Message.success("成果提交成功");
  outcomeVisible.value = false;
  outcomeForm.title = "";
  outcomeForm.content = "";
  outcomeForm.fileId = 0;
  uploadedOutcomeFile.value = null;
  fetchActivity();
}

onMounted(fetchActivity);
</script>

<style scoped>
.upload-tip {
  margin: 8px 0 0;
  color: #165dff;
  font-size: 13px;
}

.breadcrumb-line {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
}
</style>
