<template>
  <PortalLayout>
    <section class="page-heading recruit-heading">
      <div>
        <h1>招聘服务</h1>
        <p>连接人才与机构企业，助力特殊教育事业高质量发展</p>
      </div>
      <a-button type="primary" size="large" @click="openApply">企业入驻申请</a-button>
    </section>

    <section class="recruit-tabs">
      <a-tabs v-model:active-key="activeTab" type="rounded" @change="handleTabChange">
        <a-tab-pane key="jobs" title="招聘岗位">
          <section class="filter-panel">
            <a-form :model="jobQuery" layout="inline">
              <a-form-item label="关键词">
                <a-input v-model="jobQuery.keyword" allow-clear placeholder="岗位名称、企业名称、专业方向" @press-enter="searchJobs" />
              </a-form-item>
              <a-form-item label="岗位类型">
                <a-select v-model="jobQuery.jobType" allow-clear placeholder="全部类型" style="width: 130px" @change="searchJobs">
                  <a-option v-for="item in jobTypeOptions" :key="item" :value="item">{{ item }}</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="工作地点">
                <a-input v-model="jobQuery.location" allow-clear placeholder="城市/地区" style="width: 140px" @press-enter="searchJobs" />
              </a-form-item>
              <a-form-item label="学历">
                <a-select v-model="jobQuery.education" allow-clear placeholder="全部学历" style="width: 130px" @change="searchJobs">
                  <a-option v-for="item in educationOptions" :key="item" :value="item">{{ item }}</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="薪资">
                <a-select v-model="jobQuery.salaryRange" allow-clear placeholder="全部薪资" style="width: 130px" @change="searchJobs">
                  <a-option v-for="item in salaryOptions" :key="item" :value="item">{{ item }}</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="排序">
                <a-select v-model="jobQuery.sort" style="width: 130px" @change="searchJobs">
                  <a-option value="latest">最新发布</a-option>
                  <a-option value="deadline">截止时间近</a-option>
                  <a-option value="headcount">招聘人数多</a-option>
                </a-select>
              </a-form-item>
              <a-form-item>
                <a-space>
                  <a-button type="primary" @click="searchJobs">查询</a-button>
                  <a-button @click="resetJobs">重置</a-button>
                </a-space>
              </a-form-item>
            </a-form>
          </section>

          <a-spin :loading="jobLoading" style="width: 100%">
            <div v-if="jobs.length" class="job-grid">
              <article v-for="item in jobs" :key="item.id" class="job-card">
                <div class="card-head">
                  <div>
                    <h3>{{ item.jobName }}</h3>
                    <p>{{ item.companyName }}</p>
                  </div>
                  <a-tag color="arcoblue">{{ item.jobType }}</a-tag>
                </div>
                <div class="job-meta">
                  <span>{{ item.majorDirection || "专业方向不限" }}</span>
                  <span>{{ item.location }}</span>
                  <span>{{ item.headcount }} 人</span>
                  <span>{{ item.salaryRange }}</span>
                  <span>{{ item.education }}</span>
                </div>
                <div class="tag-row">
                  <a-tag v-for="tag in splitTags(item.tags)" :key="tag" color="blue">{{ tag }}</a-tag>
                </div>
                <div class="card-foot">
                  <small>{{ item.publishTime || "待发布" }}</small>
                  <router-link :to="`/recruit/jobs/${item.id}`">
                    <a-button type="primary">查看详情</a-button>
                  </router-link>
                </div>
              </article>
            </div>
            <a-empty v-else description="暂无招聘岗位" />
          </a-spin>
        </a-tab-pane>

        <a-tab-pane key="companies" title="企业库">
          <section class="filter-panel">
            <a-form :model="companyQuery" layout="inline">
              <a-form-item label="关键词">
                <a-input v-model="companyQuery.keyword" allow-clear placeholder="企业名称、所属行业、联系人" @press-enter="searchCompanies" />
              </a-form-item>
              <a-form-item label="企业性质">
                <a-select v-model="companyQuery.companyNature" allow-clear placeholder="全部性质" style="width: 140px" @change="searchCompanies">
                  <a-option v-for="item in natureOptions" :key="item" :value="item">{{ item }}</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="行业">
                <a-select v-model="companyQuery.industry" allow-clear placeholder="全部行业" style="width: 150px" @change="searchCompanies">
                  <a-option v-for="item in industryOptions" :key="item" :value="item">{{ item }}</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="规模">
                <a-select v-model="companyQuery.companySize" allow-clear placeholder="全部规模" style="width: 130px" @change="searchCompanies">
                  <a-option v-for="item in sizeOptions" :key="item" :value="item">{{ item }}</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="地区">
                <a-input v-model="companyQuery.region" allow-clear placeholder="所在地" style="width: 130px" @press-enter="searchCompanies" />
              </a-form-item>
              <a-form-item>
                <a-checkbox v-model="companyHasJobs" @change="searchCompanies">有在招岗位</a-checkbox>
              </a-form-item>
              <a-form-item>
                <a-space>
                  <a-button type="primary" @click="searchCompanies">查询</a-button>
                  <a-button @click="resetCompanies">重置</a-button>
                </a-space>
              </a-form-item>
            </a-form>
          </section>

          <a-spin :loading="companyLoading" style="width: 100%">
            <div v-if="companies.length" class="company-grid">
              <article v-for="item in companies" :key="item.id" class="company-card">
                <div class="company-logo">
                  <img v-if="item.logoUrl" :src="item.logoUrl" :alt="item.companyName" />
                  <span v-else>{{ item.companyName?.slice(0, 1) || "企" }}</span>
                </div>
                <div class="company-main">
                  <h3>{{ item.companyName }}</h3>
                  <p>{{ item.companyNature }} · {{ item.industry }} · {{ item.companySize }}</p>
                  <small>{{ item.region }} · 已发布 {{ item.jobCount || 0 }} 个岗位</small>
                  <div class="tag-row">
                    <a-tag v-for="tag in splitTags(item.tags || item.talentNeeds)" :key="tag">{{ tag }}</a-tag>
                  </div>
                </div>
                <router-link :to="`/recruit/companies/${item.id}`">
                  <a-button type="primary">查看企业</a-button>
                </router-link>
              </article>
            </div>
            <a-empty v-else description="暂无企业" />
          </a-spin>
        </a-tab-pane>
      </a-tabs>
    </section>

    <section class="value-grid">
      <div v-for="item in values" :key="item.title" class="value-card">
        <strong>{{ item.title }}</strong>
        <span>{{ item.desc }}</span>
      </div>
    </section>

    <a-modal v-model:visible="applyVisible" title="企业入驻 / 岗位发布" width="900px" :footer="false">
      <a-tabs v-model:active-key="applyTab">
        <a-tab-pane key="company" title="企业入驻申请">
          <a-form :model="companyForm" layout="vertical" class="apply-form">
            <a-row :gutter="16">
              <a-col v-for="field in companyFields" :key="field.key" :span="field.span || 8">
                <a-form-item :label="field.label" :required="field.required">
                  <component :is="field.textarea ? 'a-textarea' : 'a-input'" v-model="companyForm[field.key]" :max-length="field.max" allow-clear />
                </a-form-item>
              </a-col>
            </a-row>
            <a-button type="primary" :loading="submitting" @click="submitCompany">提交申请</a-button>
          </a-form>
        </a-tab-pane>
        <a-tab-pane key="job" title="岗位发布申请">
          <a-alert v-if="!approvedCompanies.length" type="info" class="apply-tip">企业入驻审核通过后，可提交岗位发布信息。</a-alert>
          <a-form v-else :model="jobForm" layout="vertical" class="apply-form">
            <a-row :gutter="16">
              <a-col :span="8">
                <a-form-item label="所属企业" required>
                  <a-select v-model="jobForm.companyId">
                    <a-option v-for="item in approvedCompanies" :key="item.id" :value="item.id">{{ item.companyName }}</a-option>
                  </a-select>
                </a-form-item>
              </a-col>
              <a-col v-for="field in jobFields" :key="field.key" :span="field.span || 8">
                <a-form-item :label="field.label" :required="field.required">
                  <component :is="field.textarea ? 'a-textarea' : field.number ? 'a-input-number' : 'a-input'" v-model="jobForm[field.key]" allow-clear style="width: 100%" />
                </a-form-item>
              </a-col>
            </a-row>
            <a-button type="primary" :loading="submitting" @click="submitJob">提交岗位信息</a-button>
          </a-form>
        </a-tab-pane>
      </a-tabs>
    </a-modal>
  </PortalLayout>
</template>

<script setup>
import { Message } from "@arco-design/web-vue";
import { onMounted, reactive, ref } from "vue";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getRecruitCompanies, getRecruitJobs, submitCompanyApplication, submitJobApplication } from "@/api/recruit";

const activeTab = ref("jobs");
const applyVisible = ref(false);
const applyTab = ref("company");
const jobLoading = ref(false);
const companyLoading = ref(false);
const submitting = ref(false);
const jobs = ref([]);
const companies = ref([]);
const approvedCompanies = ref([]);
const companyHasJobs = ref(false);

const jobTypeOptions = ["全职", "实习", "兼职", "校招", "社招"];
const educationOptions = ["不限", "大专", "本科", "硕士", "博士"];
const salaryOptions = ["面议", "3k-5k", "5k-8k", "8k-12k", "12k以上"];
const natureOptions = ["事业单位", "民办非企业", "国有企业", "民营企业", "社会组织"];
const industryOptions = ["特殊教育", "儿童康复", "融合教育", "心理服务", "教育科技"];
const sizeOptions = ["50人以下", "50-100人", "100-500人", "500人以上"];
const values = [
  { title: "真实可靠", desc: "企业与岗位均经平台审核后展示" },
  { title: "精准匹配", desc: "聚焦特殊教育、康复和融合教育方向" },
  { title: "高效对接", desc: "清晰呈现联系人与外部岗位链接" },
  { title: "专业专注", desc: "服务特殊教育事业人才发展" }
];

const jobQuery = reactive({ keyword: "", jobType: "", location: "", education: "", salaryRange: "", sort: "latest", pageIndex: 1, pageSize: 12 });
const companyQuery = reactive({ keyword: "", companyNature: "", industry: "", companySize: "", region: "", hasJobs: "", pageIndex: 1, pageSize: 12 });
const companyForm = reactive(defaultCompanyForm());
const jobForm = reactive(defaultJobForm());

const companyFields = [
  { key: "companyName", label: "企业名称", required: true },
  { key: "creditCode", label: "统一社会信用代码", required: true },
  { key: "companyNature", label: "企业性质", required: true },
  { key: "industry", label: "所属行业", required: true },
  { key: "companySize", label: "企业规模", required: true },
  { key: "region", label: "所在地区", required: true },
  { key: "address", label: "详细地址", required: true, span: 16 },
  { key: "website", label: "企业官网链接" },
  { key: "logoUrl", label: "企业Logo链接" },
  { key: "contactName", label: "联系人姓名", required: true },
  { key: "contactTitle", label: "联系人职务", required: true },
  { key: "contactPhone", label: "联系电话", required: true },
  { key: "contactEmail", label: "联系邮箱", required: true },
  { key: "intro", label: "企业简介", required: true, textarea: true, span: 24, max: 500 },
  { key: "mainBusiness", label: "主营业务", required: true, textarea: true, span: 12 },
  { key: "talentNeeds", label: "人才需求方向", required: true, span: 12 },
  { key: "cooperation", label: "校企合作意向", textarea: true, span: 12 },
  { key: "licenseMaterial", label: "营业执照或单位证明", required: true, span: 12 },
  { key: "qualification", label: "企业资质材料", span: 12 },
  { key: "authorizationNote", label: "联系人授权说明", span: 12 }
];
const jobFields = [
  { key: "jobName", label: "岗位名称", required: true },
  { key: "jobType", label: "岗位类型", required: true },
  { key: "headcount", label: "招聘人数", required: true, number: true },
  { key: "location", label: "工作地点", required: true },
  { key: "salaryRange", label: "薪资范围", required: true },
  { key: "education", label: "学历要求", required: true },
  { key: "majorRequirement", label: "专业要求", required: true },
  { key: "majorDirection", label: "专业方向" },
  { key: "experience", label: "工作经验要求" },
  { key: "recruitTarget", label: "招聘对象", required: true },
  { key: "deadline", label: "截止时间", required: true },
  { key: "externalLink", label: "岗位外部链接" },
  { key: "responsibilities", label: "岗位职责", required: true, textarea: true, span: 12 },
  { key: "requirements", label: "任职要求", required: true, textarea: true, span: 12 },
  { key: "workTime", label: "工作时间", textarea: true, span: 12 },
  { key: "benefits", label: "福利待遇", textarea: true, span: 12 },
  { key: "training", label: "培养机制", textarea: true, span: 12 },
  { key: "otherNotes", label: "其他说明", textarea: true, span: 12 },
  { key: "tags", label: "岗位标签", span: 12 },
  { key: "contactName", label: "岗位联系人", required: true },
  { key: "contactTitle", label: "联系人职务", required: true },
  { key: "contactPhone", label: "联系电话", required: true },
  { key: "contactEmail", label: "联系邮箱", required: true },
  { key: "contactAddress", label: "联系地址", span: 16 }
];

function pagePayload(res) {
  return res.data || {};
}

function splitTags(tags) {
  return String(tags || "").split(/[,，\s]+/).filter(Boolean).slice(0, 5);
}

function defaultCompanyForm() {
  return { companyName: "", creditCode: "", companyNature: "", industry: "", companySize: "", region: "", address: "", website: "", logoUrl: "", contactName: "", contactTitle: "", contactPhone: "", contactEmail: "", intro: "", mainBusiness: "", talentNeeds: "", cooperation: "", licenseMaterial: "", qualification: "", authorizationNote: "" };
}

function defaultJobForm() {
  return { companyId: undefined, jobName: "", jobType: "", headcount: 1, location: "", salaryRange: "", education: "", majorRequirement: "", majorDirection: "", experience: "", recruitTarget: "", deadline: "", externalLink: "", responsibilities: "", requirements: "", workTime: "", benefits: "", training: "", otherNotes: "", tags: "", contactName: "", contactTitle: "", contactPhone: "", contactEmail: "", contactAddress: "" };
}

async function fetchJobs() {
  jobLoading.value = true;
  try {
    const payload = pagePayload(await getRecruitJobs(jobQuery));
    jobs.value = payload.list || [];
  } finally {
    jobLoading.value = false;
  }
}

async function fetchCompanies() {
  companyLoading.value = true;
  try {
    companyQuery.hasJobs = companyHasJobs.value ? "1" : "";
    const payload = pagePayload(await getRecruitCompanies(companyQuery));
    companies.value = payload.list || [];
  } finally {
    companyLoading.value = false;
  }
}

async function fetchApprovedCompanies() {
  const payload = pagePayload(await getRecruitCompanies({ pageIndex: 1, pageSize: 1000 }));
  approvedCompanies.value = payload.list || [];
}

function searchJobs() {
  jobQuery.pageIndex = 1;
  fetchJobs();
}

function resetJobs() {
  Object.assign(jobQuery, { keyword: "", jobType: "", location: "", education: "", salaryRange: "", sort: "latest", pageIndex: 1 });
  fetchJobs();
}

function searchCompanies() {
  companyQuery.pageIndex = 1;
  fetchCompanies();
}

function resetCompanies() {
  Object.assign(companyQuery, { keyword: "", companyNature: "", industry: "", companySize: "", region: "", hasJobs: "", pageIndex: 1 });
  companyHasJobs.value = false;
  fetchCompanies();
}

function handleTabChange() {
  if (activeTab.value === "companies" && !companies.value.length) fetchCompanies();
}

async function openApply() {
  applyVisible.value = true;
  await fetchApprovedCompanies();
}

async function submitCompany() {
  if (!companyForm.companyName || !companyForm.creditCode || !companyForm.contactName) {
    Message.warning("请完善企业必填信息");
    return;
  }
  submitting.value = true;
  try {
    await submitCompanyApplication(companyForm);
    Message.success("入驻申请已提交，请等待平台审核。");
    Object.assign(companyForm, defaultCompanyForm());
    applyVisible.value = false;
  } finally {
    submitting.value = false;
  }
}

async function submitJob() {
  if (!jobForm.companyId || !jobForm.jobName || !jobForm.responsibilities || !jobForm.requirements) {
    Message.warning("请完善岗位必填信息");
    return;
  }
  submitting.value = true;
  try {
    await submitJobApplication(jobForm);
    Message.success("岗位信息已提交，请等待平台审核。");
    Object.assign(jobForm, defaultJobForm());
    applyVisible.value = false;
  } finally {
    submitting.value = false;
  }
}

onMounted(() => {
  fetchJobs();
  fetchCompanies();
});
</script>

<style scoped>
.recruit-heading {
  align-items: center;
}

.recruit-tabs,
.filter-panel,
.value-card,
.job-card,
.company-card {
  background: #fff;
  border-radius: 18px;
  box-shadow: 0 12px 32px rgba(20, 78, 148, 0.08);
}

.recruit-tabs {
  padding: 20px;
}

.filter-panel {
  margin-bottom: 18px;
  padding: 18px;
}

.job-grid,
.company-grid,
.value-grid {
  display: grid;
  gap: 18px;
}

.job-grid {
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
}

.company-grid {
  grid-template-columns: repeat(auto-fit, minmax(360px, 1fr));
}

.job-card,
.company-card,
.value-card {
  padding: 20px;
}

.card-head,
.card-foot,
.company-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.card-head h3,
.company-main h3 {
  margin: 0 0 6px;
  color: #12325f;
}

.card-head p,
.company-main p,
.company-main small {
  margin: 0;
  color: #5f6f89;
}

.job-meta {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
  margin: 16px 0;
  color: #44546f;
}

.tag-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin: 12px 0;
}

.company-logo {
  display: grid;
  place-items: center;
  flex: 0 0 64px;
  width: 64px;
  height: 64px;
  border-radius: 16px;
  background: #e8f3ff;
  color: #1677ff;
  font-size: 28px;
  font-weight: 700;
  overflow: hidden;
}

.company-logo img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.company-main {
  flex: 1;
}

.value-grid {
  grid-template-columns: repeat(4, minmax(0, 1fr));
  margin-top: 22px;
}

.value-card {
  display: grid;
  gap: 8px;
}

.value-card strong {
  color: #12325f;
  font-size: 18px;
}

.value-card span {
  color: #637792;
}

.apply-form {
  max-height: 62vh;
  overflow: auto;
  padding-right: 6px;
}

.apply-tip {
  margin: 12px 0;
}

@media (max-width: 900px) {
  .value-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .company-card,
  .card-head,
  .card-foot {
    align-items: flex-start;
    flex-direction: column;
  }

  .company-grid,
  .job-grid,
  .value-grid {
    grid-template-columns: 1fr;
  }
}
</style>
