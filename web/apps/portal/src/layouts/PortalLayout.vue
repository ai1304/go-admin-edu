<template>
  <a-layout class="portal-layout">
    <a-layout-header class="portal-header">
      <router-link class="brand" to="/">特殊教育资源库</router-link>
      <nav class="nav">
        <router-link to="/resources">资源中心</router-link>
        <router-link to="/courses">专题课程</router-link>
        <router-link to="/activities">教研活动</router-link>
        <router-link to="/experts">名师资源</router-link>
        <router-link to="/teacher/workbench">教师工作台</router-link>
      </nav>
      <a-dropdown v-if="session.isLoggedIn">
        <a-button>{{ session.user?.name || session.user?.userName || "已登录" }}</a-button>
        <template #content>
          <a-doption @click="goWorkbench">教师工作台</a-doption>
          <a-doption @click="handleLogout">退出登录</a-doption>
        </template>
      </a-dropdown>
      <router-link v-else to="/login">
        <a-button type="primary">登录</a-button>
      </router-link>
    </a-layout-header>
    <a-layout-content class="portal-main">
      <slot />
    </a-layout-content>
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
