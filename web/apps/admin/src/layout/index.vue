<template>
  <a-layout :style="{ height: '100vh' }">
    <a-layout-sider breakpoint="lg" :width="220" collapsible :collapsed="collapsed" @collapse="onCollapse">
      <Menu :collapsed="collapsed" />
    </a-layout-sider>
    <a-layout>
      <a-layout-header>
        <Navbar :collapsed="collapsed" @on-collapse="onCollapse" />
        <tag-view />
      </a-layout-header>
      <a-layout-content class="layout-content">
        <AppMain />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup>
import { ref } from 'vue';
import { AppMain, Navbar } from './components';
import Menu from './components/Menu/Menu.vue';
import TagView from './components/TagView/TagView.vue';

const collapsed = ref(false);
const onCollapse = () => {
  collapsed.value = !collapsed.value;
};
</script>

<style lang="scss">
@use '../style/index.scss';
@use '../style/transition.scss';
@use '../style/dark-theme.scss';

.layout-content {
  padding: 16px;
  background-color: #ffffff;
}
</style>
