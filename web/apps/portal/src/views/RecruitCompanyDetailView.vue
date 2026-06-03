<template>
  <PortalLayout>
    <a-spin :loading="loading" style="width: 100%">
      <section v-if="company.id" class="company-detail-page">
        <p class="breadcrumb">首页 / 招聘服务 / 企业主页</p>
        <section class="company-hero-card">
          <div class="company-logo">
            <img v-if="company.logoUrl" :src="company.logoUrl" :alt="company.companyName" />
            <span v-else>{{ companyInitial(company.companyName) }}</span>
          </div>
          <div>
            <a-tag color="blue">{{ company.certStatus === "approved" ? "已认证" : "认证中" }}</a-tag>
            <h1>{{ company.companyName }}</h1>
            <p>{{ company.companyNature || "用人单位" }} / {{ company.industry || "教育服务" }} / {{ company.companySize || "规模待完善" }} / {{ company.region || "全国" }}</p>
          </div>
          <a-button v-if="company.website" type="primary" @click="openExternal(company.website)">
            <template #icon><icon-link /></template>
            企业官网
          </a-button>
        </section>

        <section class="content-grid">
          <main class="main-stack">
            <section class="detail-card">
              <h2><icon-home /> 企业基础信息</h2>
              <div class="info-grid">
                <span>统一社会信用代码：{{ company.creditCode || "未公开" }}</span>
                <span>所在地区：{{ company.region || "待完善" }}</span>
                <span>详细地址：{{ company.address || "待完善" }}</span>
                <span>入驻时间：{{ formatDate(company.reviewedAt || company.createdAt) }}</span>
                <span>联系人：{{ company.contactName || "待完善" }} / {{ company.contactTitle || "联系人" }}</span>
                <span>联系电话：{{ company.contactPhone || "待完善" }}</span>
                <span>联系邮箱：{{ company.contactEmail || "待完善" }}</span>
              </div>
            </section>

            <section class="detail-card">
              <h2><icon-file /> 企业介绍</h2>
              <article v-for="item in introBlocks" :key="item.title" class="intro-block">
                <h3>{{ item.title }}</h3>
                <p>{{ item.value || item.fallback }}</p>
              </article>
            </section>

            <section class="detail-card">
              <h2><icon-safe /> 企业资质</h2>
              <div class="material-grid">
                <div v-for="item in materials" :key="item.name">
                  <strong>{{ item.name }}</strong>
                  <p>{{ item.desc }}</p>
                </div>
              </div>
            </section>
          </main>

          <aside class="detail-card jobs-card">
            <h2><icon-file /> 在招岗位</h2>
            <article v-for="item in jobs" :key="item.id" class="job-row">
              <router-link :to="`/recruit/jobs/${item.id}`">{{ item.jobName }}</router-link>
              <span>{{ item.jobType || "全职" }} / {{ item.location || "待定" }} / {{ item.headcount || 1 }}人 / {{ item.salaryRange || "面议" }}</span>
              <router-link :to="`/recruit/jobs/${item.id}`">
                <a-button size="small">查看详情</a-button>
              </router-link>
            </article>
            <a-empty v-if="!jobs.length" description="暂无在招岗位" />
          </aside>
        </section>
      </section>
      <a-empty v-else description="企业不存在或暂未开放展示" />
    </a-spin>
  </PortalLayout>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import { IconFile, IconHome, IconLink, IconSafe } from "@arco-design/web-vue/es/icon";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getRecruitCompany } from "@/api/recruit";

const route = useRoute();
const loading = ref(false);
const company = ref({});
const jobs = ref([]);

const introBlocks = computed(() => [
  { title: "企业简介", value: company.value.intro, fallback: "暂无企业简介。" },
  { title: "主营业务", value: company.value.mainBusiness, fallback: "暂无主营业务说明。" },
  { title: "人才需求方向", value: company.value.talentNeeds, fallback: "暂无人才需求说明。" },
  { title: "校企合作方向", value: company.value.cooperation, fallback: "暂无合作方向说明。" }
]);

const materials = computed(() => [
  { name: "营业执照或单位证明", desc: company.value.licenseMaterial || "后台可控制是否对前台展示下载" },
  { name: "企业资质材料", desc: company.value.qualification || "暂无公开材料" },
  { name: "联系人授权说明", desc: company.value.authorizationNote || "暂无公开材料" }
]);

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
    const res = await getRecruitCompany(route.params.id);
    company.value = res.data?.company || {};
    jobs.value = res.data?.jobs || [];
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
.company-detail-page {
  padding-top: 18px;
}

.breadcrumb {
  margin: 0 0 18px;
  color: #52657f;
  font-weight: 700;
}

.company-hero-card,
.detail-card {
  background: rgba(255, 255, 255, 0.98);
  border: 1px solid #dce8f7;
  border-radius: 8px;
  box-shadow: 0 14px 36px rgba(20, 78, 148, 0.08);
}

.company-hero-card {
  display: grid;
  grid-template-columns: 92px minmax(0, 1fr) auto;
  gap: 20px;
  align-items: center;
  padding: 26px;
  margin-bottom: 22px;
}

.company-logo {
  display: grid;
  place-items: center;
  width: 82px;
  height: 82px;
  overflow: hidden;
  color: #0969e8;
  font-size: 34px;
  font-weight: 800;
  background: #eef6ff;
  border: 2px solid #d4e6ff;
  border-radius: 50%;
}

.company-logo img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

h1,
h2,
h3 {
  color: #12325f;
}

h1 {
  margin: 10px 0 8px;
  font-size: 34px;
}

p,
span {
  color: #52657f;
  line-height: 1.8;
}

.content-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 360px;
  gap: 22px;
  align-items: stretch;
}

.main-stack {
  display: grid;
  gap: 18px;
}

.detail-card {
  padding: 24px;
}

.detail-card h2 {
  display: flex;
  gap: 8px;
  align-items: center;
  margin: 0 0 18px;
  font-size: 20px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px 18px;
}

.intro-block {
  padding: 14px 0;
  border-top: 1px solid #e4edf8;
}

.intro-block:first-of-type {
  border-top: 0;
}

.intro-block h3 {
  margin: 0 0 8px;
}

.intro-block p {
  margin: 0;
  white-space: pre-wrap;
}

.material-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 14px;
}

.material-grid div {
  padding: 14px;
  background: #f6fbff;
  border: 1px solid #dce8f7;
  border-radius: 8px;
}

.material-grid p {
  margin: 8px 0 0;
}

.jobs-card {
  align-self: stretch;
  min-width: 0;
}

.job-row {
  display: grid;
  gap: 8px;
  padding: 14px 0;
  border-top: 1px solid #e4edf8;
}

.job-row a:first-child {
  color: #12325f;
  font-weight: 800;
}

@media (max-width: 960px) {
  .company-hero-card,
  .content-grid,
  .info-grid,
  .material-grid {
    grid-template-columns: 1fr;
  }
}
</style>
