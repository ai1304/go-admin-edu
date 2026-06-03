<template>
  <PortalLayout>
    <section class="login-shell">
      <a-card :bordered="false" class="login-panel">
        <h1>教师登录</h1>
        <a-form ref="formRef" :model="form" :rules="rules" layout="vertical">
          <a-form-item field="username" label="账号">
            <a-input v-model="form.username" placeholder="请输入账号" />
          </a-form-item>
          <a-form-item field="password" label="密码">
            <a-input-password v-model="form.password" placeholder="请输入密码" />
          </a-form-item>
          <a-form-item field="code" label="验证码">
            <div class="captcha-row">
              <a-input v-model="form.code" placeholder="请输入验证码" />
              <img v-if="captchaUrl" class="captcha-img" :src="captchaUrl" alt="验证码" @click="loadCaptcha" />
              <a-button v-else @click="loadCaptcha">获取验证码</a-button>
            </div>
          </a-form-item>
          <a-button type="primary" long :loading="loading" @click="handleLogin">登录</a-button>
        </a-form>
      </a-card>
    </section>
  </PortalLayout>
</template>

<script setup>
import { Message } from "@arco-design/web-vue";
import { onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { getCaptcha, getInfo, login } from "@/api/auth";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { useSessionStore } from "@/stores/session";

const router = useRouter();
const session = useSessionStore();
const formRef = ref(null);
const loading = ref(false);
const captchaUrl = ref("");
const form = reactive({ username: "", password: "", code: "", uuid: "" });
const rules = {
  username: [{ required: true, message: "请输入账号" }],
  password: [{ required: true, message: "请输入密码" }],
  code: [{ required: true, message: "请输入验证码" }]
};

async function loadCaptcha() {
  const res = await getCaptcha();
  captchaUrl.value = res.data || "";
  form.uuid = res.id || "";
  if (!form.uuid) {
    form.uuid = "0";
  }
}

async function handleLogin() {
  const errors = await formRef.value?.validate();
  if (errors) return;
  loading.value = true;
  try {
    const res = await login({ ...form, UserName: form.username, Password: form.password, Code: form.code, UUID: form.uuid });
    if (res.code !== 200 || !res.token) {
      Message.error(res.msg || "登录失败");
      await loadCaptcha();
      return;
    }
    localStorage.setItem("portalToken", res.token);
    const info = await getInfo();
    session.setSession({ token: res.token, user: info.data || null });
    Message.success("登录成功");
    router.push("/teacher/workbench");
  } catch (error) {
    Message.error("登录失败，请检查账号、密码和验证码");
    await loadCaptcha();
  } finally {
    loading.value = false;
  }
}

onMounted(loadCaptcha);
</script>
