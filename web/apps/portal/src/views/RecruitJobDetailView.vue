<template>
  <PortalLayout>
    <a-spin :loading="loading" style="width: 100%">
      <section v-if="job.id" class="detail-layout">
        <main class="detail-main">
          <div class="detail-card">
            <a-space wrap>
              <a-tag color="arcoblue">{{ job.jobType }}</a-tag>
              <a-tag>{{ job.majorDirection || "专业方向不限" }}</a-tag>
              <a-tag color="green">{{ job.education }}</a-tag>
            </a-space>
            <h1>{{ job.jobName }}</h1>
            <p class="company-line">{{ job.companyName }} · {{ job.location }} · {{ job.salaryRange }}</p>
            <div class="info-grid">
              <span>招聘人数：{{ job.headcount }} 人</span>
              <span>专业要求：{{ job.majorRequirement }}</span>
              <span>经验要求：{{ job.experience || "不限" }}</span>
              <span>招聘对象：{{ job.recruitTarget || "不限" }}</span>
              <span>发布时间：{{ job.publishTime || "待发布" }}</span>
              <span>截止时间：{{ job.deadline || "长期有效" }}</span>
            </div>
          </div>

          <div class="detail-card">
            <h2>岗位职责</h2>
            <p>{{ job.responsibilities || "暂无说明" }}</p>
            <h2>任职要求</h2>
            <p>{{ job.requirements || "暂无说明" }}</p>
            <h2>工作时间</h2>
            <p>{{ job.workTime || "以企业说明为准" }}</p>
            <h2>福利待遇</h2>
            <p>{{ job.benefits || "以企业说明为准" }}</p>
            <h2>培养机制</h2>
            <p>{{ job.training || "暂无说明" }}</p>
            <h2>其他说明</h2>
            <p>{{ job.otherNotes || "暂无说明" }}</p>
            <div class="tag-row">
              <a-tag v-for="tag in splitTags(job.tags)" :key="tag" color="blue">{{ tag }}</a-tag>
            </div>
          </div>
        </main>

        <aside class="side-stack">
          <div class="detail-card">
            <h3>企业信息</h3>
            <strong>{{ company.companyName || job.companyName }}</strong>
            <p>{{ company.companyNature }} · {{ company.industry }}</p>
            <p>{{ company.region }} · {{ company.companySize }}</p>
            <a-tag color="green">{{ company.certStatus === "approved" ? "已认证" : "认证中" }}</a-tag>
            <router-link v-if="company.id" :to="`/recruit/companies/${company.id}`">
              <a-button long type="primary">查看企业主页</a-button>
            </router-link>
          </div>
          <div class="detail-card">
            <h3>岗位联系人</h3>
            <p>{{ job.contactName }} · {{ job.contactTitle }}</p>
            <p>{{ job.contactPhone }}</p>
            <p>{{ job.contactEmail }}</p>
            <p>{{ job.contactAddress || "联系地址以企业确认为准" }}</p>
            <a-space direction="vertical" fill>
              <a-button long @click="copyContact">复制联系方式</a-button>
              <a-button v-if="job.externalLink" long type="primary" @click="openExternal(job.externalLink)">访问岗位链接</a-button>
              <a-button v-if="company.website" long @click="openExternal(company.website)">企业官网链接</a-button>
            </a-space>
          </div>
        </aside>
      </section>
      <a-empty v-else description="岗位不存在或已过期" />
    </a-spin>
  </PortalLayout>
</template>

<script setup>
import { Message } from "@arco-design/web-vue";
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getRecruitJob } from "@/api/recruit";

const route = useRoute();
const loading = ref(false);
const job = ref({});
const company = ref({});

function splitTags(tags) {
  return String(tags || "").split(/[,，\s]+/).filter(Boolean);
}

function openExternal(url) {
  window.open(url, "_blank", "noopener,noreferrer");
}

async function copyContact() {
  const text = `${job.value.contactName} ${job.value.contactTitle}\n${job.value.contactPhone}\n${job.value.contactEmail}\n${job.value.contactAddress || ""}`;
  await navigator.clipboard.writeText(text);
  Message.success("联系方式已复制");
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
.detail-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 320px;
  gap: 22px;
}

.detail-main,
.side-stack {
  display: grid;
  gap: 18px;
  align-content: start;
}

.detail-card {
  background: #fff;
  border-radius: 18px;
  box-shadow: 0 12px 32px rgba(20, 78, 148, 0.08);
  padding: 24px;
}

h1,
h2,
h3 {
  color: #12325f;
}

h1 {
  margin: 12px 0;
  font-size: 32px;
}

h2 {
  margin: 22px 0 8px;
  font-size: 18px;
}

p {
  color: #52657f;
  line-height: 1.8;
  white-space: pre-wrap;
}

.company-line {
  font-size: 16px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  margin-top: 16px;
  color: #44546f;
}

.tag-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 18px;
}

@media (max-width: 900px) {
  .detail-layout,
  .info-grid {
    grid-template-columns: 1fr;
  }
}
</style>
