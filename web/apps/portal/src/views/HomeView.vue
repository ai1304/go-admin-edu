<template>
  <PortalLayout>
    <section class="home-hero">
      <div class="hero-copy">
        <p class="eyebrow">资源共建共享 · 教学科研协同 · 数智赋能创新 · 专业发展引领</p>
        <h1>汇聚优质资源 赋能高等特殊教育</h1>
        <p class="summary">
          面向高校特殊教育专业、教师发展与实践教学，提供资源检索、专题学习、教研协作、案例共享与名师引领的一站式服务。
        </p>
        <a-input-search
          v-model="keyword"
          class="hero-search"
          size="large"
          placeholder="搜索资源、课程、案例、名师资源..."
          search-button
          @search="handleSearch"
        />
        <div class="hot-searches">
          <span>热门搜索：</span>
          <router-link v-for="tag in hotSearches" :key="tag" :to="`/resources?keyword=${encodeURIComponent(tag)}`">{{ tag }}</router-link>
        </div>
      </div>
    </section>

    <section class="home-section">
      <div class="section-title">
        <h2>精品资源推荐</h2>
        <router-link :to="activeResourceTab.morePath">查看更多 ></router-link>
      </div>
      <a-tabs v-model:active-key="activeTab" class="resource-tabs">
        <a-tab-pane v-for="tab in resourceTabs" :key="tab.key" :title="tab.title" />
      </a-tabs>
      <div class="resource-grid">
        <router-link v-for="item in activeResourceTab.items" :key="item.title" :to="item.path" class="resource-card">
          <div class="resource-cover" :style="{ backgroundImage: item.gradient }">
            <span>{{ item.badge }}</span>
          </div>
          <strong>{{ item.title }}</strong>
          <small>{{ item.author }}</small>
          <div class="resource-meta">
            <span>浏览 {{ item.views }}</span>
            <span>收藏 {{ item.favorites }}</span>
          </div>
        </router-link>
      </div>
    </section>

    <section class="news-section" id="industry-news">
      <div class="news-list">
        <div class="section-title compact">
          <h2>行业资讯</h2>
          <router-link to="/news">查看更多资讯 ></router-link>
        </div>
        <router-link v-for="item in news" :key="item.title" :to="`/news?keyword=${encodeURIComponent(item.title)}`" class="news-row">
          <span>{{ item.title }}</span>
          <time>{{ item.date }}</time>
        </router-link>
      </div>
      <div class="news-cards">
        <router-link v-for="item in highlights" :key="item.title" :to="`/news?category=${encodeURIComponent(item.tag)}`" class="highlight-card">
          <div class="highlight-cover" :style="{ backgroundImage: item.gradient }">
            <span>{{ item.tag }}</span>
          </div>
          <strong>{{ item.title }}</strong>
          <p>{{ item.description }}</p>
          <time>{{ item.date }}</time>
        </router-link>
      </div>
    </section>
  </PortalLayout>
</template>

<script setup>
import { computed, ref } from "vue";
import { useRouter } from "vue-router";
import PortalLayout from "@/layouts/PortalLayout.vue";

const router = useRouter();
const keyword = ref("");
const activeTab = ref("featured");

const hotSearches = ["融合教育", "特殊教育学", "IEP与评估", "无障碍环境", "辅助技术", "实践教学"];

const resourceTabs = [
  {
    key: "featured",
    title: "精品课程",
    morePath: "/courses",
    items: [
      { title: "高等融合教育课程设计与实施", author: "张丽娜 教授", badge: "课程", views: "2.1万", favorites: 782, path: "/courses", gradient: "linear-gradient(135deg, #d7ebff, #4d8dff)" },
      { title: "高等教育 IEP 制定与评估实务", author: "陈静 副教授", badge: "案例", views: "1.6万", favorites: 542, path: "/courses?keyword=IEP", gradient: "linear-gradient(135deg, #f0e2ff, #8a63e6)" },
      { title: "高等特教教师教学方法与策略", author: "刘伟 副教授", badge: "方法", views: "1.2万", favorites: 431, path: "/courses", gradient: "linear-gradient(135deg, #dcf4ef, #20aa88)" },
      { title: "特殊教育学生实践能力培养", author: "周敏 副教授", badge: "实践", views: "1.3万", favorites: 487, path: "/courses", gradient: "linear-gradient(135deg, #ffe9d6, #f28b2e)" },
      { title: "无障碍课程资源建设专题", author: "平台教研组", badge: "专题", views: "9800", favorites: 358, path: "/courses?keyword=无障碍", gradient: "linear-gradient(135deg, #e1f8f3, #12a98f)" },
      { title: "融合教育课堂观察与评价", author: "王立方 教授", badge: "评价", views: "1.1万", favorites: 399, path: "/courses?keyword=融合教育", gradient: "linear-gradient(135deg, #e9edff, #5f75e8)" }
    ]
  },
  {
    key: "hot",
    title: "热门资源",
    morePath: "/resources?sort=view",
    items: [
      { title: "辅助技术在高等教育中的应用", author: "王立方 教授", badge: "资源", views: "3.3万", favorites: 965, path: "/resources?keyword=辅助技术", gradient: "linear-gradient(135deg, #d9e3ff, #4768d9)" },
      { title: "高等教育特殊支持政策解读与实践", author: "李晓慧 副教授", badge: "讲座", views: "1.8万", favorites: 643, path: "/resources?keyword=政策", gradient: "linear-gradient(135deg, #d8f6ff, #15a7c9)" },
      { title: "融合教育课程资源包", author: "平台资源中心", badge: "资源包", views: "1.7万", favorites: 612, path: "/resources?keyword=融合教育", gradient: "linear-gradient(135deg, #eef3ff, #3678dc)" },
      { title: "听障学生课堂支持工具清单", author: "资源共建小组", badge: "工具", views: "1.4万", favorites: 516, path: "/resources?keyword=听障", gradient: "linear-gradient(135deg, #fff1d8, #e79528)" },
      { title: "特教案例分析模板", author: "教研中心", badge: "模板", views: "1.2万", favorites: 488, path: "/resources?keyword=案例", gradient: "linear-gradient(135deg, #e9f8ea, #30a75b)" },
      { title: "学生发展支持记录表", author: "平台资源中心", badge: "表单", views: "1.1万", favorites: 453, path: "/resources?keyword=发展支持", gradient: "linear-gradient(135deg, #f3e7ff, #8f64d9)" }
    ]
  },
  {
    key: "latest",
    title: "最新上传",
    morePath: "/resources?sort=latest",
    items: [
      { title: "特殊教育专业实践教学案例集", author: "南京特教学院", badge: "新资源", views: "8200", favorites: 291, path: "/resources?sort=latest", gradient: "linear-gradient(135deg, #dff5ff, #23a1c8)" },
      { title: "课程思政与融合教育教学设计", author: "平台教研组", badge: "课程", views: "7600", favorites: 268, path: "/courses?keyword=课程思政", gradient: "linear-gradient(135deg, #ffe8e1, #e45b42)" },
      { title: "视障学生学习支持案例", author: "案例共建小组", badge: "案例", views: "6900", favorites: 245, path: "/resources?keyword=视障", gradient: "linear-gradient(135deg, #e6f4ff, #2976d8)" },
      { title: "AI辅助个别化学习工具指南", author: "智能应用团队", badge: "指南", views: "6500", favorites: 232, path: "/resources?keyword=AI", gradient: "linear-gradient(135deg, #eef1ff, #6366d9)" },
      { title: "高校无障碍环境建设观察表", author: "教研中心", badge: "表单", views: "5400", favorites: 197, path: "/resources?keyword=无障碍", gradient: "linear-gradient(135deg, #e3f8ef, #1ba676)" },
      { title: "实践基地共建方案样例", author: "平台运营组", badge: "方案", views: "5100", favorites: 188, path: "/resources?keyword=实践基地", gradient: "linear-gradient(135deg, #fff0dc, #ed8c22)" }
    ]
  }
];

const activeResourceTab = computed(() => resourceTabs.find((tab) => tab.key === activeTab.value) || resourceTabs[0]);

const news = [
  { title: "教育部印发关于加强高等学校特殊教育专业建设的指导意见", date: "2024-05-15" },
  { title: "第十届全国高等特殊教育发展论坛在京召开", date: "2024-05-10" },
  { title: "多所高校共建特殊教育实践教学基地 助力专业人才培养", date: "2024-05-08" },
  { title: "人工智能赋能特殊教育研究与实践应用趋势", date: "2024-05-06" },
  { title: "融合教育视野下高校资源建设和师资发展新探索", date: "2024-05-04" }
];

const highlights = [
  { tag: "政策动态", title: "推动高等特殊教育高质量发展", description: "聚焦资源共建、师资发展与实践育人，梳理政策重点与平台建设方向。", date: "2024-05-15", gradient: "linear-gradient(135deg, #dfeaff, #2c72d2)" },
  { tag: "学术前沿", title: "特殊教育领域最新研究热点", description: "梳理融合教育、辅助技术、评估支持等方向的研究进展。", date: "2024-05-12", gradient: "linear-gradient(135deg, #e6f8ff, #13a8c7)" },
  { tag: "行业发展", title: "多校合作共建产学研协同平台", description: "推进协同培养、案例共建和跨校资源共享机制。", date: "2024-05-09", gradient: "linear-gradient(135deg, #ecf8ee, #27ae60)" },
  { tag: "专家观点", title: "高等特教专业建设的实践思考", description: "从课程体系、实践基地、师资队伍等方面提出建议。", date: "2024-05-07", gradient: "linear-gradient(135deg, #fff0de, #f28b2e)" }
];

function handleSearch(value) {
  const query = (value || keyword.value || "").trim();
  router.push(query ? `/resources?keyword=${encodeURIComponent(query)}` : "/resources");
}
</script>
