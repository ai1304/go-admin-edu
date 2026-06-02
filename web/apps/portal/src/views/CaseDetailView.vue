<template>
  <PortalLayout>
    <a-spin :loading="loading" style="width: 100%">
      <article v-if="data" class="case-detail">
        <router-link class="back-link" to="/cases">返回案例列表</router-link>
        <h1>{{ data.title }}</h1>
        <div class="meta">
          <a-tag v-if="data.stage">{{ data.stage }}</a-tag>
          <a-tag v-if="data.disabilityType" color="green">{{ data.disabilityType }}</a-tag>
          <a-tag v-if="data.abilityDomain" color="blue">{{ data.abilityDomain }}</a-tag>
          <span>{{ data.school || "平台案例库" }}</span>
          <span>{{ data.viewCount || 0 }} 浏览</span>
        </div>
        <img v-if="data.coverUrl" class="case-cover" :src="data.coverUrl" :alt="data.title" />
        <p class="summary">{{ data.summary || "案例信息已脱敏展示。" }}</p>
        <section v-if="ieps.length">
          <h2>IEP 摘要</h2>
          <div v-for="item in ieps" :key="item.id" class="info-block">
            <strong>{{ item.title }}</strong>
            <p>{{ item.goal }}</p>
            <p>{{ item.plan }}</p>
          </div>
        </section>
        <section v-if="assessments.length">
          <h2>评估记录</h2>
          <div v-for="item in assessments" :key="item.id" class="info-block">
            <strong>{{ item.toolName }}</strong>
            <p>{{ item.result }}</p>
          </div>
        </section>
        <section v-if="interventions.length">
          <h2>干预策略</h2>
          <div v-for="item in interventions" :key="item.id" class="info-block">
            <strong>{{ item.title }}</strong>
            <p>{{ item.content }}</p>
          </div>
        </section>
      </article>
      <a-empty v-else description="案例不存在或未发布" />
    </a-spin>
  </PortalLayout>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { getPublishedCase } from "@/api/cases";

const route = useRoute();
const loading = ref(false);
const data = ref(null);
const ieps = ref([]);
const assessments = ref([]);
const interventions = ref([]);

async function fetchDetail() {
  loading.value = true;
  try {
    const res = await getPublishedCase(route.params.id);
    const payload = res.data || {};
    data.value = payload.case || null;
    ieps.value = payload.ieps || [];
    assessments.value = payload.assessments || [];
    interventions.value = payload.interventions || [];
  } finally {
    loading.value = false;
  }
}

onMounted(fetchDetail);
</script>

<style scoped>
.case-detail {
  max-width: 900px;
  margin: 32px auto 0;
  padding: 28px;
  background: #fff;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
}
.back-link {
  color: #176fd6;
}
h1 {
  margin: 18px 0 12px;
  font-size: 34px;
}
.meta {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  color: #86909c;
}
.summary {
  margin: 22px 0;
  color: #4e5969;
  font-size: 16px;
  line-height: 1.9;
}
.case-cover {
  width: 100%;
  max-height: 420px;
  margin-top: 22px;
  object-fit: cover;
  border-radius: 8px;
}
h2 {
  margin-top: 28px;
}
.info-block {
  margin-top: 12px;
  padding: 14px;
  background: #f7f9fc;
  border-radius: 8px;
}
.info-block p {
  color: #4e5969;
  line-height: 1.8;
  white-space: pre-wrap;
}
</style>
