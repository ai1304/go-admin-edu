<template>
  <a-layout class="portal-layout">
    <a-layout-header class="portal-header">
      <router-link class="brand" to="/" aria-label="特殊教育资源库首页">
        <span class="brand-mark">教</span>
        <span>
          <strong>特殊教育资源库</strong>
          <small>专业 · 共享 · 赋能 · 发展</small>
        </span>
      </router-link>
      <nav class="nav">
        <router-link to="/">首页</router-link>
        <router-link to="/resources">资源中心</router-link>
        <router-link to="/courses">专题课程</router-link>
        <router-link to="/activities">教研活动</router-link>
        <router-link to="/resources?type=case">特教案例</router-link>
        <router-link to="/experts">名师资源</router-link>
        <router-link to="/teacher/workbench">AI应用</router-link>
      </nav>
      <a-dropdown v-if="session.isLoggedIn">
        <a-button class="login-button">{{ session.user?.name || session.user?.userName || "已登录" }}</a-button>
        <template #content>
          <a-doption @click="goWorkbench">教师工作台</a-doption>
          <a-doption @click="handleLogout">退出登录</a-doption>
        </template>
      </a-dropdown>
      <router-link v-else to="/login">
        <a-button class="login-button" type="primary">登录 / 注册</a-button>
      </router-link>
    </a-layout-header>
    <a-layout-content class="portal-main">
      <slot />
    </a-layout-content>
    <footer class="portal-footer">
      <div class="footer-inner">
        <div class="footer-brand">
          <div class="brand footer-logo">
            <span class="brand-mark">教</span>
            <span>
              <strong>特殊教育资源库</strong>
              <small>专业 · 共享 · 赋能 · 发展</small>
            </span>
          </div>
          <p>致力于汇聚优质特殊教育资源，推动教学科研协同与创新，服务高校、教师与研究者。</p>
        </div>
        <div>
          <h3>快速导航</h3>
          <router-link to="/resources">资源中心</router-link>
          <router-link to="/courses">专题课程</router-link>
          <router-link to="/activities">教研活动</router-link>
          <router-link to="/experts">名师资源</router-link>
        </div>
        <div>
          <h3>联系我们</h3>
          <p>010-1234 5678</p>
          <p>service@tejiaoyuku.com</p>
          <p>北京市海淀区中关村东路1号院</p>
        </div>
        <div>
          <h3>关注我们</h3>
          <div class="qr-row">
            <span>微信公众号</span>
            <span>微信视频号</span>
          </div>
        </div>
      </div>
      <p class="copyright">© 2024 特殊教育资源库 版权所有 | 京ICP备2024001234号-1</p>
    </footer>
  </a-layout>
</template>

<script setup>
import { useRouter } from "vue-router";
import { logout } from "@/api/auth";
import { useSessionStore } from "@/stores/session";

const router = useRouter();
const session = useSessionStore();

function goWorkbench() {
  router.push("/teacher/workbench");
}

async function handleLogout() {
  try {
    await logout();
  } finally {
    session.clearSession();
    router.push("/");
  }
}
</script>
