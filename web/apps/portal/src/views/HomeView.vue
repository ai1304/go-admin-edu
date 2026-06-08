<template>
  <PortalLayout>
    <section class="home-hero">
      <div class="hero-copy">
        <p class="eyebrow">资源共建共享 · 教学科研协同 · 数智赋能创新 · 专业发展引领</p>
        <h1>汇聚优质资源 赋能高等特殊教育</h1>
        <p class="summary">
          面向高校特殊教育专业、教师发展与实践教学，提供资源检索、专题学习、教研协作、案例共享与名师引领的一站式服务。
        </p>
      </div>
    </section>

    <section class="home-section">
      <div class="section-title">
        <h2>精品课程</h2>
        <router-link to="/courses">查看更多 ></router-link>
      </div>
      <div class="resource-grid">
        <router-link v-for="(item, index) in featuredCourses" :key="item.id || item.title" :to="item.id ? `/courses/${item.id}` : item.path" class="resource-card">
          <div class="resource-cover image-cover">
            <img v-if="item.coverUrl" :src="item.coverUrl" :alt="item.title" />
            <img v-else :src="cardCover(item, 'course', index)" :alt="item.title" />
            <span>{{ item.category || "课程" }}</span>
          </div>
          <strong>{{ item.title }}</strong>
          <small>{{ item.teacherName || item.author || "平台课程" }}</small>
          <div class="resource-meta">
            <span>浏览 {{ formatCount(item.viewCount, item.views) }}</span>
          </div>
        </router-link>
      </div>
    </section>

    <section class="home-section">
      <div class="section-title">
        <h2>名师在线</h2>
        <router-link to="/experts">更多专家 ></router-link>
      </div>
      <div class="home-expert-grid">
        <router-link v-for="item in onlineExperts" :key="item.id || item.name" :to="item.id ? `/experts/${item.id}` : '/experts'" class="home-expert-card">
          <div class="home-expert-avatar">
            <img v-if="item.avatarUrl" :src="item.avatarUrl" :alt="item.name" />
            <span v-else>{{ (item.name || "名").slice(0, 1) }}</span>
          </div>
          <strong>{{ item.name }}</strong>
          <small>{{ item.title || "专家" }} · {{ item.organization || "平台名师库" }}</small>
          <p>{{ item.summary || item.profile || "长期关注特殊教育课程建设、教师发展与实践指导。" }}</p>
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
import { onMounted, ref } from "vue";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getPublishedCourses } from "@/api/courses";
import { getPublishedExperts } from "@/api/experts";
import { cardCover } from "@/utils/defaultCovers";

const featuredCourses = ref([]);
const onlineExperts = ref([]);

const fallbackCourses = [
  { title: "高等融合教育课程设计与实施", author: "张丽娜 教授", views: "2.1万", path: "/courses" },
  { title: "高等教育 IEP 制定与评估实务", author: "陈静 副教授", views: "1.6万", path: "/courses?keyword=IEP" },
  { title: "高等特教教师教学方法与策略", author: "刘伟 副教授", views: "1.2万", path: "/courses" },
  { title: "特殊教育学生实践能力培养", author: "周敏 副教授", views: "1.3万", path: "/courses" },
  { title: "无障碍课程资源建设专题", author: "平台教研组", views: "9800", path: "/courses?keyword=无障碍" },
  { title: "融合教育课堂观察与评价", author: "王立方 教授", views: "1.1万", path: "/courses?keyword=融合教育" }
];

const fallbackExperts = [
  { name: "张丽娜", title: "教授", organization: "特殊教育学院", summary: "融合教育课程建设与教师发展专家。" },
  { name: "王立方", title: "教授", organization: "教育科学学院", summary: "长期研究辅助技术与课堂支持。" },
  { name: "陈静", title: "副教授", organization: "实践教学中心", summary: "聚焦 IEP 制定、评估与案例督导。" },
  { name: "刘伟", title: "副教授", organization: "教师发展中心", summary: "关注特教教师专业能力培养。" },
  { name: "周敏", title: "副教授", organization: "康复教育研究所", summary: "深耕学生实践能力与康复支持。" }
];

const news = [
  { title: "教育部印发关于加强高等学校特殊教育专业建设的指导意见", date: "2024-05-15" },
  { title: "第十届全国高等特殊教育发展论坛在京召开", date: "2024-05-10" },
  { title: "多所高校共建特殊教育实践教学基地 助力专业人才培养", date: "2024-05-08" },
  { title: "人工智能赋能特殊教育研究与实践应用趋势", date: "2024-05-06" },
  { title: "融合教育视野下高校资源建设和师资发展新探索", date: "2024-05-04" }
];

const highlights = [
  { tag: "政策法规", title: "推动高等特殊教育高质量发展", description: "聚焦资源共建、师资发展与实践育人，梳理政策重点与平台建设方向。", date: "2024-05-15", gradient: "linear-gradient(135deg, #dfeaff, #2c72d2)" },
  { tag: "学术前沿", title: "特殊教育领域最新研究热点", description: "梳理融合教育、辅助技术、评估支持等方向的研究进展。", date: "2024-05-12", gradient: "linear-gradient(135deg, #e6f8ff, #13a8c7)" },
  { tag: "行业动态", title: "多校合作共建产学研协同平台", description: "推进协同培养、案例共建和跨校资源共享机制。", date: "2024-05-09", gradient: "linear-gradient(135deg, #ecf8ee, #27ae60)" },
  { tag: "优秀实践", title: "高等特教专业建设的实践思考", description: "从课程体系、实践基地、师资队伍等方面提出建议。", date: "2024-05-07", gradient: "linear-gradient(135deg, #fff0de, #f28b2e)" }
];

function pagePayload(res) {
  return res.data || {};
}

function formatCount(count, fallback = "0") {
  if (fallback && count === undefined) return fallback;
  const value = Number(count || 0);
  if (value >= 10000) return `${(value / 10000).toFixed(1)}万`;
  return String(value);
}

onMounted(async () => {
  const [coursesRes, expertsRes] = await Promise.allSettled([
    getPublishedCourses({ pageIndex: 1, pageSize: 6, sort: "view" }),
    getPublishedExperts({ pageIndex: 1, pageSize: 5, sort: "view" })
  ]);
  if (coursesRes.status === "fulfilled") {
    const payload = pagePayload(coursesRes.value);
    featuredCourses.value = payload.list || payload || [];
  }
  if (!featuredCourses.value.length) featuredCourses.value = fallbackCourses;
  if (expertsRes.status === "fulfilled") {
    const payload = pagePayload(expertsRes.value);
    onlineExperts.value = payload.list || payload || [];
  }
  if (!onlineExperts.value.length) onlineExperts.value = fallbackExperts;
});
</script>
