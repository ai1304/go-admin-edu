<template>
  <a-layout class="portal-layout">
    <a-layout-header class="portal-header">
      <router-link class="brand" to="/">特殊教育资源库</router-link>
      <nav class="nav">
        <router-link to="/resources">资源中心</router-link>
        <router-link to="/courses">专题课程</router-link>
        <router-link to="/activities">教研活动</router-link>
        <router-link to="/experts">名师资源</router-link>
      </nav>
      <router-link to="/login">
        <a-button type="primary">登录</a-button>
      </router-link>
    </a-layout-header>
    <a-layout-content class="portal-main">
      <slot />
    </a-layout-content>
  </a-layout>
</template>
