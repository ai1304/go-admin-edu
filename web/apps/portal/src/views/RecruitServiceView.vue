<template>
  <PortalLayout>
    <section class="page-heading recruit-heading">
      <div class="hero-copy">
        <h1>{{ activeTab === "jobs" ? "招聘岗位" : "企业库" }}</h1>
        <p>{{ activeTab === "jobs" ? "汇聚实习与就业岗位信息，服务特殊教育人才发展" : "展示优质合作机构与用人单位，支持岗位发布与资源对接" }}</p>
      </div>
      <a-button type="primary" size="large" @click="openApply">
        <template #icon><icon-plus /></template>
        企业入驻申请
      </a-button>
    </section>

    <section class="switch-row">
      <button :class="{ active: activeTab === 'jobs' }" @click="setTab('jobs')">招聘岗位</button>
      <button :class="{ active: activeTab === 'companies' }" @click="setTab('companies')">企业库</button>
    </section>

    <section v-if="activeTab === 'jobs'" class="jobs-layout">
      <main class="job-main">
        <section class="filter-card">
          <a-input v-model="jobQuery.keyword" allow-clear class="search-input" placeholder="搜索岗位名称、企业名称或专业方向..." @press-enter="searchJobs">
            <template #suffix><icon-search /></template>
          </a-input>
          <a-button type="primary" @click="searchJobs">搜索</a-button>
          <div class="filter-line">
            <span>岗位类型</span>
            <button v-for="item in jobTypeTabs" :key="item.value" :class="{ active: jobQuery.jobType === item.value }" @click="pickJobType(item.value)">{{ item.label }}</button>
          </div>
          <div class="filter-line">
            <span>专业方向</span>
            <button v-for="item in majorTabs" :key="item.value" :class="{ active: jobQuery.industry === item.value }" @click="pickIndustry(item.value)">{{ item.label }}</button>
          </div>
          <div class="select-line">
            <label>
              地区
              <a-select v-model="jobQuery.location" allow-clear placeholder="全部地区" @change="searchJobs">
                <a-option v-for="item in cityOptions" :key="item" :value="item">{{ item }}</a-option>
              </a-select>
            </label>
            <label>
              学历
              <a-select v-model="jobQuery.education" allow-clear placeholder="全部学历" @change="searchJobs">
                <a-option v-for="item in educationOptions" :key="item" :value="item">{{ item }}</a-option>
              </a-select>
            </label>
            <label>
              薪资
              <a-select v-model="jobQuery.salaryRange" allow-clear placeholder="全部薪资" @change="searchJobs">
                <a-option v-for="item in salaryOptions" :key="item" :value="item">{{ item }}</a-option>
              </a-select>
            </label>
            <label>
              排序
              <a-select v-model="jobQuery.sort" @change="searchJobs">
                <a-option value="latest">最新发布</a-option>
                <a-option value="deadline">截止时间近</a-option>
                <a-option value="headcount">招聘人数多</a-option>
              </a-select>
            </label>
            <a-button @click="resetJobs">重置</a-button>
          </div>
        </section>

        <a-spin :loading="jobLoading" style="width: 100%">
          <section class="job-list-card">
            <div class="result-head">
              <span>共 {{ jobTotal }} 个岗位</span>
              <span>{{ jobSortText }}</span>
            </div>
            <article v-for="item in jobs" :key="item.id" class="job-row">
              <div class="logo-badge">{{ companyInitial(item.companyName) }}</div>
              <div class="job-title-block">
                <router-link :to="`/recruit/jobs/${item.id}`">{{ item.jobName }}</router-link>
                <span>{{ item.companyName }}</span>
              </div>
              <div class="job-tags">
                <a-tag :color="item.jobType === '实习' ? 'green' : 'arcoblue'">{{ item.jobType || "全职" }}</a-tag>
                <a-tag color="blue">{{ item.industry || item.majorDirection || "特殊教育" }}</a-tag>
              </div>
              <div class="job-location">
                <span><icon-location /> {{ item.location || "待定" }}</span>
                <span><icon-calendar /> {{ formatDate(item.publishTime || item.createdAt) }}</span>
              </div>
              <p>{{ item.responsibilities || item.requirements || "岗位说明待企业完善。" }}</p>
              <router-link :to="`/recruit/jobs/${item.id}`">
                <a-button>查看详情</a-button>
              </router-link>
            </article>
            <a-empty v-if="!jobs.length" description="暂无招聘岗位" />
            <a-pagination
              v-if="jobTotal > jobQuery.pageSize"
              v-model:current="jobQuery.pageIndex"
              :page-size="jobQuery.pageSize"
              :total="jobTotal"
              show-jumper
              @change="fetchJobs"
            />
          </section>
        </a-spin>
      </main>

      <aside class="recruit-sidebar">
        <SidePanel title="热门专业方向" icon="hot" :items="hotMajors" />
        <SidePanel title="热门地区" icon="pin" :items="hotCities" />
        <section class="side-panel">
          <h3><icon-home /> 推荐企业</h3>
          <div v-for="item in recommendedCompanies" :key="item.id" class="mini-company">
            <span>{{ companyInitial(item.companyName) }}</span>
            <router-link :to="`/recruit/companies/${item.id}`">{{ item.companyName }}</router-link>
            <router-link :to="`/recruit/companies/${item.id}`">查看</router-link>
          </div>
        </section>
      </aside>
    </section>

    <section v-else class="company-section">
      <section class="company-filter">
        <label class="company-keyword">
          关键词
          <a-input v-model="companyQuery.keyword" allow-clear placeholder="搜索企业名称、简介、关键词..." @press-enter="searchCompanies">
            <template #suffix><icon-search /></template>
          </a-input>
        </label>
        <label class="company-industry">
          行业类型
          <div class="pill-row">
            <button v-for="item in industryTabs" :key="item.value" :class="{ active: companyQuery.industry === item.value }" @click="pickCompanyIndustry(item.value)">{{ item.label }}</button>
          </div>
        </label>
        <label class="company-region">
          地区
          <a-select v-model="companyQuery.region" allow-clear placeholder="全部地区" @change="searchCompanies">
            <a-option v-for="item in cityOptions" :key="item" :value="item">{{ item }}</a-option>
          </a-select>
        </label>
        <label class="company-date">
          入驻时间
          <a-range-picker v-model="companyDateRange" @change="handleCompanyDateChange" />
        </label>
        <label class="company-sort">
          排序
          <a-select v-model="companyQuery.sort" @change="searchCompanies">
            <a-option value="latest">最新入驻</a-option>
            <a-option value="jobs">岗位数量多</a-option>
          </a-select>
        </label>
        <a-button type="primary" @click="searchCompanies">查询</a-button>
      </section>

      <a-spin :loading="companyLoading" style="width: 100%">
        <section class="company-grid">
          <article v-for="item in companies" :key="item.id" class="company-card">
            <div class="company-logo">
              <img v-if="item.logoUrl" :src="item.logoUrl" :alt="item.companyName" />
              <span v-else>{{ companyInitial(item.companyName) }}</span>
            </div>
            <div>
              <h3>{{ item.companyName }}</h3>
              <a-tag color="blue">{{ item.industry || "教育服务" }}</a-tag>
            </div>
            <p>{{ item.intro || item.mainBusiness || "专注特殊教育服务、课程资源与人才协同发展。" }}</p>
            <span><icon-location /> {{ item.region || "全国" }}</span>
            <span><icon-link /> {{ item.website || "暂未填写官网" }}</span>
            <span><icon-calendar /> 入驻时间：{{ formatDate(item.createdAt) }}</span>
            <router-link :to="`/recruit/companies/${item.id}`">
              <a-button>查看企业</a-button>
            </router-link>
          </article>
          <a-empty v-if="!companies.length" description="暂无企业" />
        </section>
        <a-pagination
          v-if="companyTotal > companyQuery.pageSize"
          v-model:current="companyQuery.pageIndex"
          class="company-pagination"
          :page-size="companyQuery.pageSize"
          :total="companyTotal"
          show-jumper
          @change="fetchCompanies"
        />
      </a-spin>
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
import { IconCalendar, IconHome, IconLink, IconLocation, IconPlus, IconSearch } from "@arco-design/web-vue/es/icon";
import { computed, defineComponent, h, onMounted, reactive, ref } from "vue";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getRecruitCompanies, getRecruitJobs, submitCompanyApplication, submitJobApplication } from "@/api/recruit";

const SidePanel = defineComponent({
  props: { title: String, icon: String, items: Array },
  setup(props) {
    return () => h("section", { class: "side-panel" }, [
      h("h3", [props.icon === "hot" ? h("i", { class: "panel-dot hot" }) : h("i", { class: "panel-dot pin" }), props.title]),
      h("div", { class: "side-list" }, (props.items || []).map((item) => h("button", { onClick: item.onClick }, [
        h("span", item.name),
        h("b", item.count)
      ])))
    ]);
  }
});

const activeTab = ref("jobs");
const applyVisible = ref(false);
const applyTab = ref("company");
const jobLoading = ref(false);
const companyLoading = ref(false);
const submitting = ref(false);
const jobs = ref([]);
const companies = ref([]);
const approvedCompanies = ref([]);
const recommendedCompanies = ref([]);
const companyDateRange = ref([]);
const jobTotal = ref(0);
const companyTotal = ref(0);

const jobTypeTabs = [{ label: "全部", value: "" }, { label: "实习", value: "实习" }, { label: "就业", value: "全职" }];
const majorTabs = ["", "特殊教育", "融合教育", "言语听觉", "康复", "心理支持", "辅助技术"].map((value) => ({ label: value || "全部", value }));
const industryTabs = ["", "特殊教育学校", "康复机构", "融合教育机构", "公益组织", "教育科技", "社会服务"].map((value) => ({ label: value || "全部", value }));
const educationOptions = ["不限", "大专及以上", "本科及以上", "硕士及以上"];
const salaryOptions = ["面议", "3K-4K/月", "6K-8K/月", "8K-12K/月", "10K-15K/月", "12K-18K/月"];
const cityOptions = ["北京", "上海", "广州", "深圳", "杭州", "南京", "成都", "武汉"];
const natureOptions = ["事业单位", "民办非企业单位", "民办非企业", "国有企业", "民营企业", "社会组织"];
const sizeOptions = ["50人以下", "50-99人", "100-499人", "500人以上"];

const jobQuery = reactive({ keyword: "", jobType: "", industry: "", location: "", education: "", salaryRange: "", sort: "latest", pageIndex: 1, pageSize: 8 });
const companyQuery = reactive({ keyword: "", companyNature: "", industry: "", companySize: "", region: "", hasJobs: "", sort: "latest", startDate: "", endDate: "", pageIndex: 1, pageSize: 8 });
const companyForm = reactive(defaultCompanyForm());
const jobForm = reactive(defaultJobForm());

const hotMajors = computed(() => buildHotItems(majorTabs.filter((item) => item.value).map((item) => item.value), "industry"));
const hotCities = computed(() => buildHotItems(cityOptions, "location"));
const jobSortText = computed(() => ({ latest: "最新发布", deadline: "截止时间近", headcount: "招聘人数多" }[jobQuery.sort] || "最新发布"));

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

function defaultCompanyForm() {
  return { companyName: "", creditCode: "", companyNature: "", industry: "", companySize: "", region: "", address: "", website: "", logoUrl: "", contactName: "", contactTitle: "", contactPhone: "", contactEmail: "", intro: "", mainBusiness: "", talentNeeds: "", cooperation: "", licenseMaterial: "", qualification: "", authorizationNote: "" };
}

function defaultJobForm() {
  return { companyId: undefined, jobName: "", jobType: "", headcount: 1, location: "", salaryRange: "", education: "", majorRequirement: "", majorDirection: "", experience: "", recruitTarget: "", deadline: "", externalLink: "", responsibilities: "", requirements: "", workTime: "", benefits: "", training: "", otherNotes: "", tags: "", contactName: "", contactTitle: "", contactPhone: "", contactEmail: "", contactAddress: "" };
}

function companyInitial(name) {
  return String(name || "企").slice(0, 1);
}

function formatDate(value) {
  return value ? String(value).slice(0, 10) : "待确认";
}

function buildHotItems(names, field) {
  return names.map((name, index) => ({
    name,
    count: Math.max(6, (field === "industry" ? 68 : 32) - index * (field === "industry" ? 8 : 4)),
    onClick: () => {
      if (field === "industry") pickIndustry(name);
      else {
        jobQuery.location = name;
        searchJobs();
      }
    }
  }));
}

function setTab(tab) {
  activeTab.value = tab;
  if (tab === "companies" && !companies.value.length) fetchCompanies();
}

function pickJobType(value) {
  jobQuery.jobType = value;
  searchJobs();
}

function pickIndustry(value) {
  jobQuery.industry = value;
  searchJobs();
}

function pickCompanyIndustry(value) {
  companyQuery.industry = value;
  searchCompanies();
}

function handleCompanyDateChange(value) {
  companyQuery.startDate = value?.[0] || "";
  companyQuery.endDate = value?.[1] || "";
  searchCompanies();
}

async function fetchJobs() {
  jobLoading.value = true;
  try {
    const payload = pagePayload(await getRecruitJobs(jobQuery));
    jobs.value = payload.list || [];
    jobTotal.value = payload.count || 0;
  } finally {
    jobLoading.value = false;
  }
}

async function fetchCompanies() {
  companyLoading.value = true;
  try {
    const payload = pagePayload(await getRecruitCompanies(companyQuery));
    companies.value = payload.list || [];
    companyTotal.value = payload.count || 0;
    if (!recommendedCompanies.value.length) recommendedCompanies.value = (payload.list || []).slice(0, 5);
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
  Object.assign(jobQuery, { keyword: "", jobType: "", industry: "", location: "", education: "", salaryRange: "", sort: "latest", pageIndex: 1 });
  fetchJobs();
}

function searchCompanies() {
  companyQuery.pageIndex = 1;
  fetchCompanies();
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
  grid-template-columns: minmax(0, 1fr) auto;
  align-items: end;
}

.recruit-heading :deep(.arco-btn) {
  flex: 0 0 auto;
  min-width: 168px;
  height: 42px;
  border-radius: 4px;
}

.switch-row {
  display: flex;
  gap: 10px;
  margin-bottom: 14px;
}

.switch-row button,
.filter-line button,
.pill-row button {
  height: 32px;
  padding: 0 14px;
  color: #4e5969;
  background: #f4f6f8;
  border: 1px solid transparent;
  border-radius: 0;
  cursor: pointer;
}

.switch-row button.active,
.filter-line button.active,
.pill-row button.active {
  color: #0b6be8;
  background: #eef6ff;
  border-color: transparent;
}

.jobs-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 300px;
  gap: 18px;
}

.filter-card,
.job-list-card,
.side-panel,
.company-filter,
.company-card {
  background: rgba(255, 255, 255, 0.96);
  border: 1px solid #e5e6eb;
  border-radius: 8px;
  box-shadow: none;
}

.filter-card {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 96px;
  gap: 14px;
  padding: 18px;
  margin-bottom: 14px;
}

.filter-line,
.select-line {
  grid-column: 1 / -1;
}

.filter-line {
  display: flex;
  flex-wrap: wrap;
  gap: 0;
  align-items: center;
}

.filter-line span {
  width: 72px;
  color: #21314f;
  font-weight: 700;
}

.select-line {
  display: grid;
  grid-template-columns: repeat(4, minmax(130px, 1fr)) auto;
  gap: 14px;
  align-items: end;
}

.select-line label,
.company-filter label {
  display: grid;
  gap: 8px;
  color: #223452;
  font-weight: 700;
}

.job-list-card {
  padding: 0 16px 16px;
}

.result-head {
  display: flex;
  justify-content: space-between;
  padding: 14px 2px;
  color: #52657f;
}

.job-row {
  display: grid;
  grid-template-columns: 58px 210px 180px 160px minmax(0, 1fr) 104px;
  gap: 14px;
  align-items: center;
  min-height: 86px;
  padding: 12px 0;
  border-top: 1px solid #e8eef7;
}

.logo-badge,
.company-logo {
  display: grid;
  place-items: center;
  color: #0969e8;
  font-weight: 800;
  background: #eef6ff;
  border: 1px solid #cfe2fb;
}

.logo-badge {
  width: 46px;
  height: 46px;
  border-radius: 50%;
}

.job-title-block {
  display: grid;
  gap: 5px;
}

.job-title-block a,
.company-card h3 {
  color: #12325f;
  font-weight: 800;
}

.job-title-block span,
.job-location span,
.job-row p {
  color: #52657f;
}

.job-tags,
.job-location {
  display: grid;
  gap: 8px;
}

.job-row p {
  margin: 0;
  line-height: 1.7;
}

.recruit-sidebar {
  display: grid;
  gap: 14px;
  align-content: start;
}

.side-panel {
  padding: 18px;
}

.side-panel h3 {
  display: flex;
  gap: 8px;
  align-items: center;
  margin: 0 0 12px;
  color: #12325f;
}

.panel-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.panel-dot.hot {
  background: #ff6b3d;
}

.panel-dot.pin {
  background: #0969e8;
}

.side-list {
  display: grid;
  gap: 8px;
}

.side-list button,
.mini-company {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  align-items: center;
  color: #52657f;
  background: transparent;
  border: 0;
}

.side-list button {
  cursor: pointer;
}

.side-list b,
.mini-company a:last-child {
  color: #0969e8;
}

.mini-company span {
  display: grid;
  place-items: center;
  width: 24px;
  height: 24px;
  color: #0969e8;
  background: #eef6ff;
  border-radius: 50%;
}

.company-section {
  padding: 20px;
  background: rgba(255, 255, 255, 0.92);
  border: 1px solid #e5e6eb;
  border-radius: 8px;
}

.company-filter {
  display: grid;
  grid-template-columns: minmax(220px, 1.25fr) minmax(160px, 0.8fr) minmax(260px, 1.25fr) minmax(150px, 0.7fr) auto;
  gap: 14px 18px;
  align-items: end;
  padding: 18px;
  margin-bottom: 18px;
}

.company-industry {
  grid-column: 1 / -1;
}

.company-keyword,
.company-region,
.company-date,
.company-sort {
  min-width: 0;
}

.pill-row {
  display: flex;
  flex-wrap: wrap;
  gap: 0;
}

.company-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
}

.company-card {
  display: grid;
  gap: 12px;
  padding: 20px;
}

.company-logo {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  overflow: hidden;
}

.company-logo img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.company-card h3 {
  margin: 0 0 6px;
}

.company-card p,
.company-card span {
  margin: 0;
  color: #52657f;
  line-height: 1.7;
}

.company-pagination {
  justify-content: center;
  margin-top: 18px;
}

.apply-form {
  max-height: 62vh;
  overflow: auto;
  padding-right: 6px;
}

.apply-tip {
  margin: 12px 0;
}

@media (max-width: 1080px) {
  .jobs-layout,
  .company-filter {
    grid-template-columns: 1fr;
  }

  .company-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .job-row {
    grid-template-columns: 58px minmax(0, 1fr);
  }

  .job-row p,
  .job-tags,
  .job-location,
  .job-row > a {
    grid-column: 2;
  }
}

@media (max-width: 640px) {
  .recruit-heading {
    grid-template-columns: 1fr;
  }

  .filter-card,
  .select-line,
  .company-grid {
    grid-template-columns: 1fr;
  }
}
</style>
