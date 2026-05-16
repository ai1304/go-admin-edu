<template>
  <div class="container">
    <a-card :bordered="false" class="cardStyle">
      <template #title>数据中心</template>
      <a-space>
        <a-button type="primary" @click="fetchData">刷新</a-button>
        <a-button @click="handleExport">导出概览 CSV</a-button>
      </a-space>
    </a-card>

    <div class="stats-grid">
      <a-card v-for="item in items" :key="item.key" :bordered="false" class="stat-card">
        <div class="label">{{ item.label }}</div>
        <div class="value">{{ overview[item.key] ?? 0 }}</div>
      </a-card>
    </div>

    <div class="panel-grid">
      <a-card :bordered="false" title="资源状态分布">
        <a-table :columns="pairColumns" :data="resourceStats.byStatus || []" :pagination="false" size="small" />
      </a-card>
      <a-card :bordered="false" title="资源类型分布">
        <a-table :columns="pairColumns" :data="resourceStats.byResourceType || []" :pagination="false" size="small" />
      </a-card>
      <a-card :bordered="false" title="课程学习概览">
        <a-descriptions :column="1" bordered>
          <a-descriptions-item label="学习人数">{{ courseStats.totalLearners || 0 }}</a-descriptions-item>
          <a-descriptions-item label="完成课时">{{ courseStats.finishedLessons || 0 }}</a-descriptions-item>
          <a-descriptions-item label="作业提交">{{ courseStats.submissions || 0 }}</a-descriptions-item>
        </a-descriptions>
      </a-card>
      <a-card :bordered="false" title="活动参与概览">
        <a-descriptions :column="1" bordered>
          <a-descriptions-item label="报名数">{{ activityStats.signups || 0 }}</a-descriptions-item>
          <a-descriptions-item label="签到数">{{ activityStats.checkins || 0 }}</a-descriptions-item>
          <a-descriptions-item label="成果数">{{ activityStats.outcomes || 0 }}</a-descriptions-item>
        </a-descriptions>
      </a-card>
    </div>

    <div class="wide-grid">
      <a-card :bordered="false" title="学校贡献排行">
        <a-table :columns="schoolColumns" :data="schoolStats" :pagination="false" size="small" />
      </a-card>
      <a-card :bordered="false" title="教师活跃排行">
        <a-table :columns="teacherColumns" :data="teacherStats" :pagination="false" size="small" />
      </a-card>
      <a-card :bordered="false" title="学生学习排行">
        <a-table :columns="studentColumns" :data="studentStats" :pagination="false" size="small" />
      </a-card>
      <a-card :bordered="false" title="敏感案例概览">
        <a-descriptions :column="2" bordered>
          <a-descriptions-item label="IEP">{{ caseStats.ieps || 0 }}</a-descriptions-item>
          <a-descriptions-item label="评估记录">{{ caseStats.assessments || 0 }}</a-descriptions-item>
          <a-descriptions-item label="干预方案">{{ caseStats.interventions || 0 }}</a-descriptions-item>
          <a-descriptions-item label="访问日志">{{ caseStats.accessLogs || 0 }}</a-descriptions-item>
          <a-descriptions-item label="授权记录">{{ caseStats.authorizations || 0 }}</a-descriptions-item>
        </a-descriptions>
      </a-card>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import {
  exportEduStats,
  getEduActivityStats,
  getEduCaseStats,
  getEduCourseStats,
  getEduOverview,
  getEduResourceStats,
  getEduSchoolStats,
  getEduStudentStats,
  getEduTeacherStats
} from '@/api/edu/stats';

const overview = ref({});
const resourceStats = ref({});
const courseStats = ref({});
const activityStats = ref({});
const caseStats = ref({});
const schoolStats = ref([]);
const teacherStats = ref([]);
const studentStats = ref([]);

const items = [
  { key: 'regions', label: '区域数' },
  { key: 'schools', label: '学校数' },
  { key: 'resources', label: '资源数' },
  { key: 'publishedResources', label: '已发布资源' },
  { key: 'courses', label: '课程数' },
  { key: 'activities', label: '活动数' },
  { key: 'cases', label: '案例数' },
  { key: 'experts', label: '专家数' },
  { key: 'views', label: '资源浏览' },
  { key: 'downloads', label: '资源下载' },
  { key: 'learners', label: '学习人数' }
];

const pairColumns = [
  { title: '名称', dataIndex: 'name' },
  { title: '数量', dataIndex: 'value', width: 120 }
];
const schoolColumns = [
  { title: '学校', dataIndex: 'schoolName' },
  { title: '资源', dataIndex: 'resourceCount', width: 90 },
  { title: '课程', dataIndex: 'courseCount', width: 90 },
  { title: '活动', dataIndex: 'activityCount', width: 90 },
  { title: '总计', dataIndex: 'total', width: 90 }
];
const teacherColumns = [
  { title: '用户 ID', dataIndex: 'userId', width: 100 },
  { title: '资源', dataIndex: 'resourceCount', width: 90 },
  { title: '课程', dataIndex: 'courseCount', width: 90 },
  { title: '案例', dataIndex: 'caseCount', width: 90 },
  { title: '总计', dataIndex: 'total', width: 90 }
];
const studentColumns = [
  { title: '学习标识', dataIndex: 'identity', ellipsis: true, tooltip: true },
  { title: '课时', dataIndex: 'lessonCount', width: 90 },
  { title: '完成', dataIndex: 'finishedCount', width: 90 },
  { title: '作业', dataIndex: 'submissionCount', width: 90 }
];

async function fetchData() {
  const [overviewRes, resourceRes, courseRes, activityRes, schoolRes, teacherRes, studentRes, caseRes] = await Promise.all([
    getEduOverview(),
    getEduResourceStats(),
    getEduCourseStats(),
    getEduActivityStats(),
    getEduSchoolStats(),
    getEduTeacherStats(),
    getEduStudentStats(),
    getEduCaseStats()
  ]);
  overview.value = overviewRes.data || {};
  resourceStats.value = resourceRes.data || {};
  courseStats.value = courseRes.data || {};
  activityStats.value = activityRes.data || {};
  schoolStats.value = schoolRes.data || [];
  teacherStats.value = teacherRes.data || [];
  studentStats.value = studentRes.data || [];
  caseStats.value = caseRes.data || {};
}

async function handleExport() {
  const blob = await exportEduStats();
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = `edu-stats-${Date.now()}.csv`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  window.URL.revokeObjectURL(url);
}

onMounted(fetchData);
</script>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
  margin-top: 16px;
}

.panel-grid,
.wide-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
  margin-top: 16px;
}

.label {
  color: #86909c;
}

.value {
  margin-top: 8px;
  font-size: 30px;
  font-weight: 700;
}
</style>
