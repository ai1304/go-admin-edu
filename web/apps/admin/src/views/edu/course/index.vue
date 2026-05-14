<template>
  <div class="container">
    <a-card :bordered="false" class="cardStyle">
      <template #title>课程管理</template>
      <a-form :model="queryForm" layout="inline">
        <a-form-item label="关键词">
          <a-input v-model="queryForm.keyword" allow-clear placeholder="课程标题、教师" />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model="queryForm.status" allow-clear placeholder="请选择状态" style="width: 160px">
            <a-option v-for="item in statusOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="fetchData">查询</a-button>
            <a-button @click="resetQuery">重置</a-button>
            <a-button type="primary" status="success" @click="openCreate">新增课程</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>

    <a-card :bordered="false" class="cardStyle table-card">
      <a-table :columns="columns" :data="tableData" :pagination="pagination" row-key="id" @page-change="handlePageChange">
        <template #status="{ record }">
          <a-tag :color="statusColor[record.status]">{{ statusText[record.status] || record.status }}</a-tag>
        </template>
        <template #operations="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openEdit(record)">编辑</a-button>
            <a-button type="text" size="small" @click="openStructure(record)">章节/课时</a-button>
            <a-button type="text" status="danger" size="small" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="formModel.id ? '编辑课程' : '新增课程'" width="760px" @before-ok="handleSave">
      <a-form :model="formModel" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item field="title" label="课程标题" required>
              <a-input v-model="formModel.title" placeholder="请输入课程标题" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="teacherName" label="授课教师">
              <a-input v-model="formModel.teacherName" placeholder="请输入授课教师" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="difficulty" label="难度">
              <a-select v-model="formModel.difficulty" allow-clear placeholder="请选择难度">
                <a-option value="basic">基础</a-option>
                <a-option value="advanced">进阶</a-option>
                <a-option value="expert">专家</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="status" label="状态">
              <a-select v-model="formModel.status">
                <a-option v-for="item in statusOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="stageCategoryId" label="学段">
              <a-select v-model="formModel.stageCategoryId" allow-clear placeholder="请选择学段">
                <a-option v-for="item in getCategoryOptions('stage')" :key="item.id" :value="item.id">{{ item.name }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="disabilityTypeId" label="障碍类型">
              <a-select v-model="formModel.disabilityTypeId" allow-clear placeholder="请选择障碍类型">
                <a-option v-for="item in getCategoryOptions('disability')" :key="item.id" :value="item.id">{{ item.name }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item field="summary" label="课程简介">
              <a-textarea v-model="formModel.summary" :auto-size="{ minRows: 3, maxRows: 5 }" placeholder="请输入课程简介" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item field="objectives" label="教学目标">
              <a-textarea v-model="formModel.objectives" :auto-size="{ minRows: 3, maxRows: 5 }" placeholder="请输入教学目标" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>

    <a-modal v-model:visible="structureVisible" :title="`${currentCourse?.title || ''} 章节/课时`" width="980px" :footer="false">
      <a-tabs default-active-key="chapters">
        <a-tab-pane key="chapters" title="章节">
          <a-space direction="vertical" fill>
            <a-button type="primary" status="success" @click="openChapterCreate">新增章节</a-button>
            <a-table :columns="chapterColumns" :data="chapterList" :pagination="false" row-key="id">
              <template #status="{ record }">
                <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '停用' }}</a-tag>
              </template>
              <template #operations="{ record }">
                <a-space>
                  <a-button type="text" size="small" @click="openChapterEdit(record)">编辑</a-button>
                  <a-button type="text" status="danger" size="small" @click="handleChapterDelete(record)">删除</a-button>
                </a-space>
              </template>
            </a-table>
          </a-space>
        </a-tab-pane>
        <a-tab-pane key="lessons" title="课时">
          <a-space direction="vertical" fill>
            <a-button type="primary" status="success" @click="openLessonCreate">新增课时</a-button>
            <a-table :columns="lessonColumns" :data="lessonList" :pagination="false" row-key="id">
              <template #chapter="{ record }">{{ chapterName(record.chapterId) }}</template>
              <template #duration="{ record }">{{ formatDuration(record.durationSeconds) }}</template>
              <template #status="{ record }">
                <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '停用' }}</a-tag>
              </template>
              <template #operations="{ record }">
                <a-space>
                  <a-button type="text" size="small" @click="openLessonEdit(record)">编辑</a-button>
                  <a-button type="text" status="danger" size="small" @click="handleLessonDelete(record)">删除</a-button>
                </a-space>
              </template>
            </a-table>
          </a-space>
        </a-tab-pane>
      </a-tabs>
    </a-modal>

    <a-modal v-model:visible="chapterVisible" :title="chapterModel.id ? '编辑章节' : '新增章节'" width="520px" @before-ok="handleChapterSave">
      <a-form :model="chapterModel" layout="vertical">
        <a-form-item field="title" label="章节标题" required>
          <a-input v-model="chapterModel.title" placeholder="请输入章节标题" />
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item field="sort" label="排序">
              <a-input-number v-model="chapterModel.sort" :min="0" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="status" label="状态">
              <a-select v-model="chapterModel.status">
                <a-option :value="1">启用</a-option>
                <a-option :value="0">停用</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>

    <a-modal v-model:visible="lessonVisible" :title="lessonModel.id ? '编辑课时' : '新增课时'" width="640px" @before-ok="handleLessonSave">
      <a-form :model="lessonModel" layout="vertical">
        <a-form-item field="title" label="课时标题" required>
          <a-input v-model="lessonModel.title" placeholder="请输入课时标题" />
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item field="chapterId" label="所属章节" required>
              <a-select v-model="lessonModel.chapterId" placeholder="请选择章节">
                <a-option v-for="item in chapterList" :key="item.id" :value="item.id">{{ item.title }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="durationSeconds" label="时长（秒）">
              <a-input-number v-model="lessonModel.durationSeconds" :min="0" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item field="videoFileId" label="视频文件 ID">
              <a-input-number v-model="lessonModel.videoFileId" :min="0" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item field="sort" label="排序">
              <a-input-number v-model="lessonModel.sort" :min="0" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item field="status" label="状态">
              <a-select v-model="lessonModel.status">
                <a-option :value="1">启用</a-option>
                <a-option :value="0">停用</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { Message, Modal } from '@arco-design/web-vue';
import { onMounted, reactive, ref } from 'vue';
import {
  addCourse,
  addCourseChapter,
  addCourseLesson,
  getCourseChapters,
  getCourseLessons,
  getCourses,
  removeCourseChapters,
  removeCourseLessons,
  removeCourses,
  updateCourse,
  updateCourseChapter,
  updateCourseLesson
} from '@/api/edu/course';
import { getResourceCategories } from '@/api/edu/resource';

const statusOptions = [
  { label: '草稿', value: 'draft' },
  { label: '已发布', value: 'published' },
  { label: '已下架', value: 'offline' }
];
const statusText = Object.fromEntries(statusOptions.map((item) => [item.value, item.label]));
const statusColor = { draft: 'gray', published: 'green', offline: 'gray' };
const queryForm = reactive({ keyword: '', status: '', pageIndex: 1, pageSize: 10 });
const tableData = ref([]);
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const formVisible = ref(false);
const formModel = reactive(defaultForm());
const categoryOptions = reactive({});
const structureVisible = ref(false);
const chapterVisible = ref(false);
const lessonVisible = ref(false);
const currentCourse = ref(null);
const chapterList = ref([]);
const lessonList = ref([]);
const chapterModel = reactive(defaultChapterForm());
const lessonModel = reactive(defaultLessonForm());

const columns = [
  { title: '课程标题', dataIndex: 'title', ellipsis: true, tooltip: true },
  { title: '教师', dataIndex: 'teacherName', width: 120 },
  { title: '难度', dataIndex: 'difficulty', width: 100 },
  { title: '状态', slotName: 'status', width: 110 },
  { title: '学习人数', dataIndex: 'learnerCount', width: 110 },
  { title: '操作', slotName: 'operations', width: 220 }
];
const chapterColumns = [
  { title: '章节标题', dataIndex: 'title', ellipsis: true, tooltip: true },
  { title: '排序', dataIndex: 'sort', width: 90 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'operations', width: 150 }
];
const lessonColumns = [
  { title: '课时标题', dataIndex: 'title', ellipsis: true, tooltip: true },
  { title: '所属章节', slotName: 'chapter', width: 180 },
  { title: '时长', slotName: 'duration', width: 110 },
  { title: '视频文件 ID', dataIndex: 'videoFileId', width: 120 },
  { title: '排序', dataIndex: 'sort', width: 90 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'operations', width: 150 }
];

function defaultForm() {
  return {
    id: undefined,
    title: '',
    summary: '',
    teacherName: '',
    difficulty: '',
    objectives: '',
    status: 'draft',
    stageCategoryId: undefined,
    disabilityTypeId: undefined
  };
}

function assignForm(data = {}) {
  Object.assign(formModel, defaultForm(), data);
}

function defaultChapterForm() {
  return { id: undefined, title: '', sort: 0, status: 1 };
}

function defaultLessonForm() {
  return { id: undefined, chapterId: undefined, title: '', videoFileId: 0, durationSeconds: 0, sort: 0, status: 1 };
}

async function fetchData() {
  const res = await getCourses(queryForm);
  const payload = res.data || {};
  tableData.value = payload.list || payload || [];
  pagination.total = payload.count || res.total || 0;
}

async function fetchCategories() {
  const res = await getResourceCategories({ pageIndex: 1, pageSize: 1000, status: 1 });
  const payload = res.data || {};
  const list = payload.list || payload || [];
  categoryOptions.stage = [];
  categoryOptions.disability = [];
  list.forEach((item) => {
    if (!categoryOptions[item.type]) {
      categoryOptions[item.type] = [];
    }
    categoryOptions[item.type].push(item);
  });
}

function getCategoryOptions(type) {
  return categoryOptions[type] || [];
}

function handlePageChange(page) {
  queryForm.pageIndex = page;
  pagination.current = page;
  fetchData();
}

function resetQuery() {
  queryForm.keyword = '';
  queryForm.status = '';
  queryForm.pageIndex = 1;
  pagination.current = 1;
  fetchData();
}

function openCreate() {
  assignForm();
  formVisible.value = true;
}

function openEdit(record) {
  assignForm(record);
  formVisible.value = true;
}

async function handleSave() {
  if (!formModel.title) {
    Message.warning('请输入课程标题');
    return false;
  }
  const payload = { ...formModel };
  if (payload.id) {
    await updateCourse(payload.id, payload);
  } else {
    await addCourse(payload);
  }
  Message.success('保存成功');
  formVisible.value = false;
  fetchData();
}

function handleDelete(record) {
  Modal.confirm({
    title: '确认删除课程',
    content: `确定删除「${record.title}」吗？`,
    async onOk() {
      await removeCourses({ ids: [record.id] });
      Message.success('删除成功');
      fetchData();
    }
  });
}

async function openStructure(record) {
  currentCourse.value = record;
  structureVisible.value = true;
  await fetchStructure();
}

async function fetchStructure() {
  if (!currentCourse.value?.id) return;
  const [chaptersRes, lessonsRes] = await Promise.all([
    getCourseChapters(currentCourse.value.id),
    getCourseLessons(currentCourse.value.id)
  ]);
  chapterList.value = chaptersRes.data || [];
  lessonList.value = lessonsRes.data || [];
}

function openChapterCreate() {
  Object.assign(chapterModel, defaultChapterForm());
  chapterVisible.value = true;
}

function openChapterEdit(record) {
  Object.assign(chapterModel, defaultChapterForm(), record);
  chapterVisible.value = true;
}

async function handleChapterSave() {
  if (!chapterModel.title) {
    Message.warning('请输入章节标题');
    return false;
  }
  const payload = { ...chapterModel };
  if (payload.id) {
    await updateCourseChapter(currentCourse.value.id, payload.id, payload);
  } else {
    await addCourseChapter(currentCourse.value.id, payload);
  }
  Message.success('保存成功');
  chapterVisible.value = false;
  fetchStructure();
}

function handleChapterDelete(record) {
  Modal.confirm({
    title: '确认删除章节',
    content: `确定删除「${record.title}」吗？`,
    async onOk() {
      await removeCourseChapters(currentCourse.value.id, { ids: [record.id] });
      Message.success('删除成功');
      fetchStructure();
    }
  });
}

function openLessonCreate() {
  Object.assign(lessonModel, defaultLessonForm(), {
    chapterId: chapterList.value[0]?.id
  });
  lessonVisible.value = true;
}

function openLessonEdit(record) {
  Object.assign(lessonModel, defaultLessonForm(), record);
  lessonVisible.value = true;
}

async function handleLessonSave() {
  if (!lessonModel.title) {
    Message.warning('请输入课时标题');
    return false;
  }
  if (!lessonModel.chapterId) {
    Message.warning('请选择所属章节');
    return false;
  }
  const payload = { ...lessonModel };
  if (payload.id) {
    await updateCourseLesson(currentCourse.value.id, payload.id, payload);
  } else {
    await addCourseLesson(currentCourse.value.id, payload);
  }
  Message.success('保存成功');
  lessonVisible.value = false;
  fetchStructure();
}

function handleLessonDelete(record) {
  Modal.confirm({
    title: '确认删除课时',
    content: `确定删除「${record.title}」吗？`,
    async onOk() {
      await removeCourseLessons(currentCourse.value.id, { ids: [record.id] });
      Message.success('删除成功');
      fetchStructure();
    }
  });
}

function chapterName(chapterId) {
  return chapterList.value.find((item) => item.id === chapterId)?.title || '未分配章节';
}

function formatDuration(seconds = 0) {
  if (!seconds) return '未设置';
  const minutes = Math.floor(seconds / 60);
  const rest = seconds % 60;
  return `${minutes}分${rest}秒`;
}

onMounted(() => {
  fetchData();
  fetchCategories();
});
</script>

<style scoped>
.table-card {
  margin-top: 16px;
}
</style>
