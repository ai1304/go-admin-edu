<template>
  <PortalLayout>
    <a-spin :loading="loading" style="width: 100%">
      <section v-if="company.id" class="company-detail">
        <div class="hero-card">
          <div class="logo-box">
            <img v-if="company.logoUrl" :src="company.logoUrl" :alt="company.companyName" />
            <span v-else>{{ company.companyName?.slice(0, 1) || "企" }}</span>
          </div>
          <div>
            <a-tag color="green">{{ company.certStatus === "approved" ? "已认证" : "认证中" }}</a-tag>
            <h1>{{ company.companyName }}</h1>
            <p>{{ company.companyNature }} · {{ company.industry }} · {{ company.companySize }} · {{ company.region }}</p>
          </div>
          <a-button v-if="company.website" type="primary" @click="openExternal(company.website)">企业官网链接</a-button>
        </div>

        <section class="content-grid">
          <main class="main-stack">
            <div class="detail-card">
              <h2>企业基础信息</h2>
              <div class="info-grid">
                <span>统一社会信用代码：{{ company.creditCode }}</span>
                <span>所在地区：{{ company.region }}</span>
                <span>详细地址：{{ company.address }}</span>
                <span>成立合作日期：{{ formatDate(company.reviewedAt || company.createdAt) }}</span>
                <span>联系人：{{ company.contactName }} · {{ company.contactTitle }}</span>
                <span>联系电话：{{ company.contactPhone }}</span>
                <span>联系邮箱：{{ company.contactEmail }}</span>
              </div>
            </div>

            <div class="detail-card">
              <h2>企业介绍</h2>
              <h3>企业简介</h3>
              <p>{{ company.intro || "暂无简介" }}</p>
              <h3>主营业务</h3>
              <p>{{ company.mainBusiness || "暂无说明" }}</p>
              <h3>人才需求方向</h3>
              <p>{{ company.talentNeeds || "暂无说明" }}</p>
              <h3>校企合作方向</h3>
              <p>{{ company.cooperation || "暂无说明" }}</p>
            </div>

            <div class="detail-card">
              <h2>企业资质</h2>
              <a-list :data="materials" :bordered="false">
                <template #item="{ item }">
                  <a-list-item>
                    <a-list-item-meta :title="item.name" :description="item.desc" />
                  </a-list-item>
                </template>
              </a-list>
            </div>
          </main>

          <aside class="detail-card">
            <h2>已发布岗位</h2>
            <div v-if="jobs.length" class="job-list">
              <article v-for="item in jobs" :key="item.id" class="job-row">
                <strong>{{ item.jobName }}</strong>
                <span>{{ item.jobType }} · {{ item.location }} · {{ item.headcount }} 人 · {{ item.salaryRange }}</span>
                <router-link :to="`/recruit/jobs/${item.id}`">
                  <a-button size="small" type="primary">查看详情</a-button>
                </router-link>
              </article>
            </div>
            <a-empty v-else description="暂无在招岗位" />
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
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getRecruitCompany } from "@/api/recruit";

const route = useRoute();
const loading = ref(false);
const company = ref({});
const jobs = ref([]);

const materials = computed(() => [
  { name: "营业执照或单位证明", desc: company.value.licenseMaterial || "后台可控制是否对前台展示下载" },
  { name: "企业资质材料", desc: company.value.qualification || "暂无公开材料" },
  { name: "联系人授权说明", desc: company.value.authorizationNote || "暂无公开材料" }
]);

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
.company-detail {
  display: grid;
  gap: 22px;
}

.hero-card,
.detail-card {
  background: #fff;
  border-radius: 18px;
  box-shadow: 0 12px 32px rgba(20, 78, 148, 0.08);
  padding: 24px;
}

.hero-card {
  display: flex;
  align-items: center;
  gap: 20px;
}

.logo-box {
  display: grid;
  place-items: center;
  flex: 0 0 82px;
  width: 82px;
  height: 82px;
  border-radius: 22px;
  background: #e8f3ff;
  color: #1677ff;
  font-size: 34px;
  font-weight: 700;
  overflow: hidden;
}

.logo-box img {
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
  font-size: 32px;
}

p,
span {
  color: #52657f;
  line-height: 1.8;
  white-space: pre-wrap;
}

.content-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 360px;
  gap: 22px;
}

.main-stack,
.job-list {
  display: grid;
  gap: 18px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.job-row {
  display: grid;
  gap: 8px;
  border-bottom: 1px solid #eef3f8;
  padding-bottom: 14px;
}

@media (max-width: 900px) {
  .hero-card {
    align-items: flex-start;
    flex-direction: column;
  }

  .content-grid,
  .info-grid {
    grid-template-columns: 1fr;
  }
}
</style>
