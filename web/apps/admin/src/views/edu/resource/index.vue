<template>
  <div class="container">
    <a-card :bordered="false" class="cardStyle">
      <template #title>资源管理</template>
      <a-form :model="queryForm" layout="inline">
        <a-form-item label="关键词">
          <a-input v-model="queryForm.keyword" allow-clear placeholder="标题、简介、标签" />
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
            <a-button type="primary" status="success" @click="openCreate">新增资源</a-button>
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
            <a-button type="text" size="small" @click="openFiles(record)">附件</a-button>
            <a-button type="text" size="small" @click="openComments(record)">评论</a-button>
            <a-button v-if="record.status === 'draft' || record.status === 'rejected'" type="text" size="small" @click="handleSubmitReview(record)">提交审核</a-button>
            <a-button v-if="record.status === 'reviewing'" type="text" size="small" @click="handleReview(record, 'approve')">通过</a-button>
            <a-button v-if="record.status === 'reviewing'" type="text" status="warning" size="small" @click="handleReview(record, 'reject')">驳回</a-button>
            <a-button type="text" status="danger" size="small" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="formModel.id ? '编辑资源' : '新增资源'" width="760px" @before-ok="handleSave">
      <a-form :model="formModel" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item field="title" label="资源标题" required>
              <a-input v-model="formModel.title" placeholder="请输入资源标题" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="authorName" label="作者">
              <a-input v-model="formModel.authorName" placeholder="请输入作者/教师" />
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
            <a-form-item field="schoolId" label="学校 ID">
              <a-input-number v-model="formModel.schoolId" :min="0" placeholder="请输入学校 ID" style="width: 100%" />
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
          <a-col :span="12">
            <a-form-item field="resourceTypeId" label="资源类型">
              <a-select v-model="formModel.resourceTypeId" allow-clear placeholder="请选择资源类型">
                <a-option v-for="item in getCategoryOptions('resource_type')" :key="item.id" :value="item.id">{{ item.name }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="abilityDomainId" label="能力领域">
              <a-select v-model="formModel.abilityDomainId" allow-clear placeholder="请选择能力领域">
                <a-option v-for="item in getCategoryOptions('ability_domain')" :key="item.id" :value="item.id">{{ item.name }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="topicCategoryId" label="专题分类">
              <a-select v-model="formModel.topicCategoryId" allow-clear placeholder="请选择专题分类">
                <a-option v-for="item in getCategoryOptions('topic')" :key="item.id" :value="item.id">{{ item.name }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item field="keywords" label="关键词">
              <a-input v-model="formModel.keywords" placeholder="多个关键词用逗号分隔" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item field="summary" label="资源简介">
              <a-textarea v-model="formModel.summary" :auto-size="{ minRows: 3, maxRows: 6 }" placeholder="请输入资源简介" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>

    <a-modal v-model:visible="fileVisible" :title="`${currentResource?.title || ''} 附件`" width="760px" :footer="false">
      <a-space direction="vertical" fill>
        <a-space>
          <a-button type="primary" @click="triggerUpload('attachment')">上传附件</a-button>
          <a-button @click="triggerUpload('cover')">上传封面</a-button>
          <a-tag v-if="currentResource?.coverFileId" color="blue">当前封面文件 ID：{{ currentResource.coverFileId }}</a-tag>
          <input ref="fileInput" type="file" class="hidden-file-input" @change="handleFileChange" />
        </a-space>
        <a-table :columns="fileColumns" :data="fileList" :pagination="false" row-key="id">
          <template #usage="{ record }">
            <a-tag :color="record.usage === 'cover' ? 'blue' : 'gray'">{{ record.usage === 'cover' ? '封面' : '附件' }}</a-tag>
          </template>
          <template #size="{ record }">{{ formatSize(record.size) }}</template>
          <template #fileOperations="{ record }">
            <a-button type="text" status="danger" size="small" @click="handleDeleteFile(record)">删除</a-button>
          </template>
        </a-table>
      </a-space>
    </a-modal>

    <a-modal v-model:visible="commentVisible" :title="`${currentResource?.title || ''} 评论`" width="860px" :footer="false">
      <a-table :columns="commentColumns" :data="commentList" :pagination="false" row-key="id">
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '显示' : '隐藏' }}</a-tag>
        </template>
        <template #commentOperations="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="toggleCommentStatus(record)">{{ record.status === 1 ? '隐藏' : '显示' }}</a-button>
            <a-button type="text" status="danger" size="small" @click="handleDeleteComment(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-modal>
  </div>
</template>

<script setup>
import { Message, Modal } from '@arco-design/web-vue';
import { onMounted, reactive, ref } from 'vue';
import {
  addResource,
  getResourceCategories,
  getResourceComments,
  getResourceFiles,
  getResources,
  removeResourceComments,
  removeResourceFiles,
  removeResources,
  reviewResource,
  submitResourceReview,
  updateResourceComment,
  updateResource,
  uploadResourceFile
} from '@/api/edu/resource';

const statusOptions = [
  { label: '草稿', value: 'draft' },
  { label: '审核中', value: 'reviewing' },
  { label: '已发布', value: 'published' },
  { label: '已驳回', value: 'rejected' },
  { label: '已下架', value: 'offline' }
];
const statusText = Object.fromEntries(statusOptions.map((item) => [item.value, item.label]));
const statusColor = { draft: 'gray', reviewing: 'orange', published: 'green', rejected: 'red', offline: 'gray' };

const queryForm = reactive({ keyword: '', status: '', pageIndex: 1, pageSize: 10 });
const tableData = ref([]);
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const formVisible = ref(false);
const fileVisible = ref(false);
const commentVisible = ref(false);
const fileInput = ref(null);
const uploadUsage = ref('attachment');
const fileList = ref([]);
const commentList = ref([]);
const currentResource = ref(null);
const formModel = reactive(defaultForm());
const categoryOptions = reactive({});

const columns = [
  { title: '资源标题', dataIndex: 'title', ellipsis: true, tooltip: true },
  { title: '作者', dataIndex: 'authorName', width: 120 },
  { title: '学校 ID', dataIndex: 'schoolId', width: 100 },
  { title: '状态', slotName: 'status', width: 110 },
  { title: '浏览', dataIndex: 'viewCount', width: 90 },
  { title: '下载', dataIndex: 'downloadCount', width: 90 },
  { title: '收藏', dataIndex: 'favoriteCount', width: 90 },
  { title: '操作', slotName: 'operations', width: 360 }
];
const fileColumns = [
  { title: '文件名', dataIndex: 'originalName', ellipsis: true, tooltip: true },
  { title: '用途', slotName: 'usage', width: 90 },
  { title: '类型', dataIndex: 'contentType', width: 180 },
  { title: '大小', slotName: 'size', width: 120 },
  { title: '操作', slotName: 'fileOperations', width: 100 }
];
const commentColumns = [
  { title: '昵称', dataIndex: 'nickname', width: 140 },
  { title: '内容', dataIndex: 'content', ellipsis: true, tooltip: true },
  { title: '点赞', dataIndex: 'likeCount', width: 80 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'commentOperations', width: 140 }
];

function defaultForm() {
  return {
    id: undefined,
    title: '',
    summary: '',
    authorName: '',
    keywords: '',
    status: 'draft',
    schoolId: 0,
    stageCategoryId: undefined,
    disabilityTypeId: undefined,
    resourceTypeId: undefined,
    abilityDomainId: undefined,
    topicCategoryId: undefined
  };
}

function assignForm(data = {}) {
  Object.assign(formModel, defaultForm(), data);
}

function getPagePayload(res) {
  return res.data || {};
}

async function fetchData() {
  const res = await getResources(queryForm);
  const payload = getPagePayload(res);
  tableData.value = payload.list || payload || [];
  pagination.total = payload.count || res.total || 0;
}

async function fetchCategories() {
  const res = await getResourceCategories({ pageIndex: 1, pageSize: 1000, status: 1 });
  const payload = getPagePayload(res);
  const list = payload.list || payload || [];
  categoryOptions.stage = [];
  categoryOptions.disability = [];
  categoryOptions.resource_type = [];
  categoryOptions.ability_domain = [];
  categoryOptions.topic = [];
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
    Message.warning('请输入资源标题');
    return false;
  }
  const payload = { ...formModel };
  if (payload.id) {
    await updateResource(payload.id, payload);
  } else {
    await addResource(payload);
  }
  Message.success('保存成功');
  formVisible.value = false;
  fetchData();
}

function handleDelete(record) {
  Modal.confirm({
    title: '确认删除资源',
    content: `确定删除「${record.title}」吗？`,
    async onOk() {
      await removeResources({ ids: [record.id] });
      Message.success('删除成功');
      fetchData();
    }
  });
}

async function handleSubmitReview(record) {
  await submitResourceReview(record.id);
  Message.success('已提交审核');
  fetchData();
}

async function handleReview(record, action) {
  await reviewResource(record.id, { action, comment: action === 'approve' ? '审核通过' : '审核驳回' });
  Message.success(action === 'approve' ? '已发布' : '已驳回');
  fetchData();
}

async function openFiles(record) {
  currentResource.value = record;
  fileVisible.value = true;
  await fetchFiles();
}

async function openComments(record) {
  currentResource.value = record;
  commentVisible.value = true;
  await fetchComments();
}

async function fetchFiles() {
  if (!currentResource.value?.id) return;
  const res = await getResourceFiles({ resourceId: currentResource.value.id, pageIndex: 1, pageSize: 100 });
  const payload = getPagePayload(res);
  fileList.value = payload.list || payload || [];
}

async function fetchComments() {
  if (!currentResource.value?.id) return;
  const res = await getResourceComments(currentResource.value.id);
  commentList.value = res.data || [];
}

function triggerUpload(usage = 'attachment') {
  uploadUsage.value = usage;
  fileInput.value?.click();
}

async function handleFileChange(event) {
  const file = event.target.files?.[0];
  if (!file || !currentResource.value?.id) return;
  const formData = new FormData();
  formData.append('file', file);
  formData.append('resourceId', currentResource.value.id);
  formData.append('usage', uploadUsage.value);
  const res = await uploadResourceFile(formData);
  const uploadedFile = res.data || {};
  if (uploadUsage.value === 'cover' && uploadedFile.id) {
    await updateResource(currentResource.value.id, { ...currentResource.value, coverFileId: uploadedFile.id });
    currentResource.value.coverFileId = uploadedFile.id;
  }
  event.target.value = '';
  Message.success(uploadUsage.value === 'cover' ? '封面上传成功' : '上传成功');
  fetchFiles();
  fetchData();
}

function handleDeleteFile(record) {
  Modal.confirm({
    title: '确认删除附件',
    content: `确定删除「${record.originalName}」吗？`,
    async onOk() {
      await removeResourceFiles({ ids: [record.id] });
      Message.success('删除成功');
      fetchFiles();
    }
  });
}

async function toggleCommentStatus(record) {
  const status = record.status === 1 ? 0 : 1;
  await updateResourceComment(currentResource.value.id, record.id, { ...record, status });
  Message.success(status === 1 ? '评论已显示' : '评论已隐藏');
  fetchComments();
}

function handleDeleteComment(record) {
  Modal.confirm({
    title: '确认删除评论',
    content: `确定删除「${record.nickname || '访客'}」的评论吗？`,
    async onOk() {
      await removeResourceComments(currentResource.value.id, { ids: [record.id] });
      Message.success('删除成功');
      fetchComments();
    }
  });
}

function formatSize(size = 0) {
  if (size < 1024) return `${size} B`;
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
  return `${(size / 1024 / 1024).toFixed(1)} MB`;
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

.hidden-file-input {
  display: none;
}
</style>
