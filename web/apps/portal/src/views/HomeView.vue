<template>
  <PortalLayout>
    <section class="home-hero">
      <div class="hero-copy">
        <p class="eyebrow">资源共建共享 · 教学科研协同 · 数智赋能创新 · 专业发展引领</p>
        <h1>汇聚优质资源 赋能高等特殊教育</h1>
        <p class="summary">面向高校特殊教育专业、教师发展与实践教学，提供资源检索、专题学习、教研协作、案例共享与名师引领的一站式服务。</p>
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

    <section class="shortcut-bar" aria-label="快捷入口">
      <router-link v-for="item in shortcuts" :key="item.path" :to="item.path" class="shortcut-item">
        <span class="shortcut-icon">{{ item.icon }}</span>
        <strong>{{ item.title }}</strong>
      </router-link>
    </section>

    <section class="feature-grid">
      <router-link v-for="item in featureLinks" :key="item.path" :to="item.path" class="feature-card">
        <span class="feature-icon">{{ item.icon }}</span>
        <span>
          <strong>{{ item.title }}</strong>
          <small>{{ item.description }}</small>
        </span>
        <i>→</i>
      </router-link>
    </section>

    <section class="home-section">
      <div class="section-title">
        <h2>精品资源推荐</h2>
        <router-link to="/resources">查看更多 ></router-link>
      </div>
      <a-tabs default-active-key="featured" class="resource-tabs">
        <a-tab-pane key="featured" title="精品课程" />
        <a-tab-pane key="hot" title="热门资源" />
        <a-tab-pane key="latest" title="最新上传" />
      </a-tabs>
      <div class="resource-grid">
        <router-link v-for="item in resources" :key="item.title" :to="item.path" class="resource-card">
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

    <section class="news-section">
      <div class="news-list">
        <div class="section-title compact">
          <h2>行业资讯</h2>
          <a href="#">查看更多资讯 ></a>
        </div>
        <a v-for="item in news" :key="item.title" href="#" class="news-row">
          <span>{{ item.title }}</span>
          <time>{{ item.date }}</time>
        </a>
      </div>
      <div class="news-cards">
        <a v-for="item in highlights" :key="item.title" href="#" class="highlight-card">
          <div class="highlight-cover" :style="{ backgroundImage: item.gradient }">
            <span>{{ item.tag }}</span>
          </div>
          <strong>{{ item.title }}</strong>
          <p>{{ item.description }}</p>
          <time>{{ item.date }}</time>
        </a>
      </div>
    </section>
  </PortalLayout>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import PortalLayout from "@/layouts/PortalLayout.vue";

const router = useRouter();
const keyword = ref("");

const hotSearches = ["融合教育", "特殊教育学", "IEP与评估", "无障碍环境", "辅助技术", "实践教学"];

const shortcuts = [
  { title: "资源检索", icon: "检", path: "/resources" },
  { title: "课程学习", icon: "课", path: "/courses" },
  { title: "活动报名", icon: "报", path: "/activities" },
  { title: "案例库", icon: "案", path: "/resources?type=case" },
  { title: "名师名家", icon: "师", path: "/experts" },
  { title: "机构共建", icon: "建", path: "/teacher/workbench" }
];

const featureLinks = [
  { title: "资源中心", description: "海量资源一站检索", icon: "资源", path: "/resources" },
  { title: "专题课程", description: "系统学习 提升专业", icon: "课程", path: "/courses" },
  { title: "教研活动", description: "协同教研 交流共进", icon: "活动", path: "/activities" },
  { title: "特教案例", description: "实践案例 经验共享", icon: "案例", path: "/resources?type=case" },
  { title: "名师资源", description: "专家引领 专业成长", icon: "名师", path: "/experts" }
];

const resources = [
  { title: "高等融合教育课程设计与实施", author: "张丽娜 教授", badge: "课程", views: "2.1万", favorites: 782, path: "/courses", gradient: "linear-gradient(135deg, #d7ebff, #4d8dff)" },
  { title: "高等教育特殊支持政策解读与实践", author: "李晓慧 副教授", badge: "讲座", views: "1.8万", favorites: 643, path: "/resources", gradient: "linear-gradient(135deg, #d8f6ff, #15a7c9)" },
  { title: "辅助技术在高等教育中的应用", author: "王立新 教授", badge: "资源", views: "3.3万", favorites: 965, path: "/resources", gradient: "linear-gradient(135deg, #d9e3ff, #4768d9)" },
  { title: "高等教育 IEP 制定与评估实务", author: "陈静 副教授", badge: "案例", views: "1.6万", favorites: 542, path: "/resources?keyword=IEP", gradient: "linear-gradient(135deg, #f0e2ff, #8a63e6)" },
  { title: "高等特教教师教学方法与策略", author: "刘伟 副教授", badge: "方法", views: "1.2万", favorites: 431, path: "/courses", gradient: "linear-gradient(135deg, #dcf4ef, #20aa88)" },
  { title: "特殊教育学生实践能力培养", author: "周敏 副教授", badge: "实践", views: "1.3万", favorites: 487, path: "/activities", gradient: "linear-gradient(135deg, #ffe9d6, #f28b2e)" }
];

const news = [
  { title: "教育部印发《关于加强高等学校特殊教育专业建设的指导意见》", date: "2024-05-15" },
  { title: "第十届全国高等特殊教育发展论坛在京召开", date: "2024-05-10" },
  { title: "多所高校共建特殊教育实践教学基地 助力专业人才培养", date: "2024-05-08" },
  { title: "人工智能赋能特殊教育研究与实践应用趋势", date: "2024-05-06" },
  { title: "融合教育视域下高校资源建设和师资发展新探索", date: "2024-05-04" }
];

const highlights = [
  { tag: "政策动态", title: "最新政策：推动高等教育特殊教育高质量发展", description: "深入解读政策重点，聚焦资源共建、师资发展与实践育人。", date: "2024-05-15", gradient: "linear-gradient(135deg, #dfeaff, #2c72d2)" },
  { tag: "学术前沿", title: "研究趋势：特殊教育领域最新研究热点", description: "梳理融合教育、辅助技术、评估支持等方向的研究进展。", date: "2024-05-12", gradient: "linear-gradient(135deg, #e6f8ff, #13a8c7)" },
  { tag: "行业发展", title: "多校合作：共建特殊教育产学研协同创新平台", description: "推进协同培养、案例共建和跨校资源共享机制。", date: "2024-05-09", gradient: "linear-gradient(135deg, #ecf8ee, #27ae60)" },
  { tag: "专家观点", title: "专家观点：高等特殊教育专业建设的思考", description: "专家从课程体系、实践基地、师资队伍等方面提出建议。", date: "2024-05-07", gradient: "linear-gradient(135deg, #fff0de, #f28b2e)" }
];

function handleSearch(value) {
  const query = (value || keyword.value || "").trim();
  router.push(query ? `/resources?keyword=${encodeURIComponent(query)}` : "/resources");
}
</script>
