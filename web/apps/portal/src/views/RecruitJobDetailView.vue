<template>
  <PortalLayout>
    <a-spin :loading="loading" style="width: 100%">
      <section v-if="job.id" class="job-detail-page">
        <p class="breadcrumb">首页 / 招聘服务 / 岗位详情</p>
        <div class="detail-layout">
          <main class="detail-card main-card">
            <h1>{{ job.jobName }}</h1>
            <div class="top-meta">
              <span><icon-home /> {{ job.companyName }}</span>
              <span><icon-file /> {{ job.jobType || "全职" }}</span>
              <span><icon-location /> {{ job.location || "工作地点待定" }}</span>
              <span><icon-safe /> {{ job.salaryRange || "面议" }}</span>
              <span><icon-book /> {{ job.education || "学历不限" }}</span>
              <span><icon-user /> {{ job.headcount || 1 }}人</span>
              <span><icon-calendar /> 发布日期：{{ formatDate(job.publishTime || job.createdAt) }}</span>
              <span><icon-calendar-clock /> 截止日期：{{ job.deadline || "长期有效" }}</span>
            </div>

            <section v-for="block in contentBlocks" :key="block.title" class="content-block">
              <div class="section-icon">
                <component :is="block.icon" />
              </div>
              <div>
                <h2>{{ block.title }}</h2>
                <ul v-if="block.items.length">
                  <li v-for="item in block.items" :key="item">{{ item }}</li>
                </ul>
                <p v-else>{{ block.fallback }}</p>
              </div>
            </section>

            <section class="tag-band">
              <span><icon-tags /> 岗位标签</span>
              <a-tag v-for="tag in splitTags(job.tags || job.industry || job.majorDirection)" :key="tag" color="blue">{{ tag }}</a-tag>
            </section>
          </main>

          <aside class="side-stack">
            <section class="detail-card side-card">
              <h3><icon-home /> 企业信息</h3>
              <div class="company-head">
                <div class="company-logo">{{ companyInitial(company.companyName || job.companyName) }}</div>
                <div>
                  <strong>{{ company.companyName || job.companyName }}</strong>
                  <a-tag color="blue">{{ company.certStatus === "approved" ? "已认证" : "认证中" }}</a-tag>
                </div>
              </div>
              <dl>
                <dt>企业性质</dt><dd>{{ company.companyNature || "待完善" }}</dd>
                <dt>所属行业</dt><dd>{{ company.industry || job.industry || "教育" }}</dd>
                <dt>所在地</dt><dd>{{ company.region || job.location || "待完善" }}</dd>
              </dl>
              <router-link v-if="company.id" :to="`/recruit/companies/${company.id}`">
                <a-button long>查看企业主页</a-button>
              </router-link>
            </section>

            <section class="detail-card side-card">
              <h3><icon-user /> 岗位联系人</h3>
              <div class="contact-head">
                <div class="avatar"><icon-user /></div>
                <div>
                  <strong>{{ job.contactName || "联系人待完善" }}</strong>
                  <a-tag>{{ job.contactTitle || "岗位负责人" }}</a-tag>
                </div>
              </div>
              <p><icon-phone /> {{ job.contactPhone || "待完善" }}</p>
              <p><icon-email /> {{ job.contactEmail || "待完善" }}</p>
              <p><icon-location /> {{ job.contactAddress || company.address || "联系地址以企业确认为准" }}</p>
              <a-space direction="vertical" fill>
                <a-button v-if="job.externalLink" long type="primary" @click="openExternal(job.externalLink)">访问岗位链接</a-button>
              </a-space>
            </section>
          </aside>
        </div>
      </section>
      <a-empty v-else description="岗位不存在或已过期" />
    </a-spin>
  </PortalLayout>
</template>

<script setup>
import {
  IconBook,
  IconCalendar,
  IconCalendarClock,
  IconClockCircle,
  IconEmail,
  IconFile,
  IconGift,
  IconHome,
  IconInfoCircle,
  IconLocation,
  IconPhone,
  IconSafe,
  IconTags,
  IconUser,
  IconUserGroup
} from "@arco-design/web-vue/es/icon";
import { computed, onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getRecruitJob } from "@/api/recruit";

const route = useRoute();
const loading = ref(false);
const job = ref({});
const company = ref({});

const contentBlocks = computed(() => [
  { title: "岗位职责", icon: IconFile, items: toLines(job.value.responsibilities), fallback: "岗位职责待企业完善。" },
  { title: "任职要求", icon: IconUserGroup, items: toLines(job.value.requirements || job.value.majorRequirement), fallback: "任职要求待企业完善。" },
  { title: "工作时间", icon: IconClockCircle, items: toLines(job.value.workTime), fallback: "以企业说明为准。" },
  { title: "福利待遇", icon: IconGift, items: toLines(job.value.benefits), fallback: "以企业说明为准。" },
  { title: "其他说明", icon: IconInfoCircle, items: toLines([job.value.training, job.value.otherNotes].filter(Boolean).join("\n")), fallback: "暂无补充说明。" }
]);

function toLines(value) {
  return String(value || "").split(/\n|；|;/).map((item) => item.replace(/^[-·\s]+/, "").trim()).filter(Boolean);
}

function splitTags(tags) {
  return String(tags || "融合教育,特殊教育").split(/[,，\s]+/).filter(Boolean).slice(0, 6);
}

function companyInitial(name) {
  return String(name || "企").slice(0, 1);
}

function formatDate(value) {
  return value ? String(value).slice(0, 10) : "待确认";
}

function openExternal(url) {
  window.open(url, "_blank", "noopener,noreferrer");
}

onMounted(async () => {
  loading.value = true;
  try {
    const res = await getRecruitJob(route.params.id);
    job.value = res.data?.job || {};
    company.value = res.data?.company || {};
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
.job-detail-page {
  padding-top: 18px;
}

.breadcrumb {
  margin: 0 0 18px;
  color: #52657f;
  font-weight: 700;
}

.detail-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 360px;
  gap: 22px;
  align-items: stretch;
}

.detail-card {
  background: rgba(255, 255, 255, 0.98);
  border: 1px solid #dce8f7;
  border-radius: 8px;
  box-shadow: 0 14px 36px rgba(20, 78, 148, 0.08);
}

.main-card {
  padding: 28px;
}

h1,
h2,
h3 {
  color: #12325f;
}

h1 {
  margin: 0 0 18px;
  font-size: 34px;
}

.top-meta {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px 22px;
  padding-bottom: 22px;
  border-bottom: 1px solid #dce8f7;
}

.top-meta span,
.side-card p {
  display: inline-flex;
  gap: 8px;
  align-items: center;
  color: #223452;
}

.content-block {
  display: grid;
  grid-template-columns: 54px minmax(0, 1fr);
  gap: 16px;
  padding: 22px 0;
  border-bottom: 1px solid #e4edf8;
}

.section-icon {
  display: grid;
  place-items: center;
  width: 46px;
  height: 46px;
  color: #0969e8;
  font-size: 22px;
  background: #eef6ff;
  border-radius: 50%;
}

.content-block h2 {
  margin: 0 0 10px;
  font-size: 18px;
}

.content-block ul {
  padding-left: 18px;
  margin: 0;
  color: #223452;
  line-height: 1.9;
}

.content-block p {
  margin: 0;
  color: #52657f;
}

.tag-band {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
  padding: 16px;
  background: #eef6ff;
  border-radius: 8px;
}

.tag-band span {
  display: inline-flex;
  gap: 8px;
  align-items: center;
  color: #12325f;
  font-weight: 800;
}

.side-stack {
  display: flex;
  flex-direction: column;
  gap: 18px;
  height: 100%;
  min-width: 0;
}

.side-stack .side-card:last-child {
  flex: 1;
}

.side-card {
  display: grid;
  gap: 18px;
  padding: 24px;
}

.side-card h3 {
  display: flex;
  gap: 8px;
  align-items: center;
  margin: 0;
}

.company-head,
.contact-head {
  display: flex;
  gap: 14px;
  align-items: center;
}

.company-logo,
.avatar {
  display: grid;
  place-items: center;
  width: 72px;
  height: 72px;
  color: #0969e8;
  font-size: 30px;
  font-weight: 800;
  background: #eef6ff;
  border: 2px solid #d4e6ff;
  border-radius: 50%;
}

.avatar {
  font-size: 32px;
}

dl {
  display: grid;
  grid-template-columns: 88px minmax(0, 1fr);
  gap: 14px 10px;
  margin: 0;
}

dt {
  color: #52657f;
}

dd {
  margin: 0;
  color: #223452;
}

@media (max-width: 960px) {
  .detail-layout,
  .top-meta {
    grid-template-columns: 1fr;
  }
}
</style>
