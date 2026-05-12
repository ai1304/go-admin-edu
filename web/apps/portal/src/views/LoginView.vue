<template>
  <PortalLayout>
    <section class="login-shell">
      <a-card :bordered="false" class="login-panel">
        <h1>登录</h1>
        <a-form layout="vertical">
          <a-form-item label="账号">
            <a-input placeholder="请输入账号" />
          </a-form-item>
          <a-form-item label="密码">
            <a-input-password placeholder="请输入密码" />
          </a-form-item>
          <a-button type="primary" long>登录</a-button>
        </a-form>
      </a-card>
    </section>
  </PortalLayout>
</template>

<script setup>
import PortalLayout from "@/layouts/PortalLayout.vue";
</script>
