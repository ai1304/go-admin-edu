<template>
  <div class="container">
    <a-card :bordered="false" class="cardStyle">
      <template #title>教研活动管理</template>
      <a-form :model="queryForm" layout="inline">
        <a-form-item label="关键词">
          <a-input v-model="queryForm.keyword" allow-clear placeholder="活动名称、主办方" />
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
            <a-button type="primary" status="success" @click="openCreate">新增活动</a-button>
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
            <a-button type="text" size="small" @click="openManage(record)">业务管理</a-button>
            <a-button type="text" status="danger" size="small" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="formModel.id ? '编辑活动' : '新增活动'" width="760px" @before-ok="handleSave">
      <a-form :model="formModel" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item field="title" label="活动名称" required>
              <a-input v-model="formModel.title" placeholder="请输入活动名称" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="organizer" label="主办方">
              <a-input v-model="formModel.organizer" placeholder="请输入主办方" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="startTime" label="开始时间">
              <a-input v-model="formModel.startTime" placeholder="例如 2026-05-20 09:00" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="endTime" label="结束时间">
              <a-input v-model="formModel.endTime" placeholder="例如 2026-05-20 17:00" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="location" label="地点">
              <a-input v-model="formModel.location" placeholder="请输入活动地点" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="status" label="状态">
              <a-select v-model="formModel.status">
                <a-option v-for="item in statusOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item field="summary" label="活动简介">
              <a-textarea v-model="formModel.summary" :auto-size="{ minRows: 3, maxRows: 6 }" placeholder="请输入活动简介" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>

    <a-modal v-model:visible="manageVisible" :title="`${currentActivity?.title || ''} 业务管理`" width="980px" :footer="false">
      <a-tabs default-active-key="signups">
        <a-tab-pane key="signups" title="报名">
          <a-space direction="vertical" fill>
            <a-button type="primary" status="success" @click="openSignupCreate">新增报名</a-button>
            <a-table :columns="signupColumns" :data="signupList" :pagination="false" row-key="id">
              <template #operations="{ record }">
                <a-space>
                  <a-button type="text" size="small" @click="openSignupEdit(record)">编辑</a-button>
                  <a-button type="text" status="danger" size="small" @click="handleSignupDelete(record)">删除</a-button>
                </a-space>
              </template>
            </a-table>
          </a-space>
        </a-tab-pane>
        <a-tab-pane key="checkins" title="签到">
          <a-space direction="vertical" fill>
            <a-button type="primary" status="success" @click="openCheckinCreate">新增签到</a-button>
            <a-table :columns="checkinColumns" :data="checkinList" :pagination="false" row-key="id">
              <template #operations="{ record }">
                <a-button type="text" status="danger" size="small" @click="handleCheckinDelete(record)">删除</a-button>
              </template>
            </a-table>
          </a-space>
        </a-tab-pane>
        <a-tab-pane key="outcomes" title="成果">
          <a-space direction="vertical" fill>
            <a-button type="primary" status="success" @click="openOutcomeCreate">新增成果</a-button>
            <a-table :columns="outcomeColumns" :data="outcomeList" :pagination="false" row-key="id">
              <template #status="{ record }">
                <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '停用' }}</a-tag>
              </template>
              <template #operations="{ record }">
                <a-space>
                  <a-button type="text" size="small" @click="openOutcomeEdit(record)">编辑</a-button>
                  <a-button type="text" status="danger" size="small" @click="handleOutcomeDelete(record)">删除</a-button>
                </a-space>
              </template>
            </a-table>
          </a-space>
        </a-tab-pane>
      </a-tabs>
    </a-modal>

    <a-modal v-model:visible="signupVisible" :title="signupModel.id ? '编辑报名' : '新增报名'" width="520px" @before-ok="handleSignupSave">
      <a-form :model="signupModel" layout="vertical">
        <a-form-item field="name" label="姓名" required>
          <a-input v-model="signupModel.name" placeholder="请输入姓名" />
        </a-form-item>
        <a-form-item field="phone" label="电话">
          <a-input v-model="signupModel.phone" placeholder="请输入电话" />
        </a-form-item>
        <a-form-item field="status" label="状态">
          <a-select v-model="signupModel.status">
            <a-option value="signed">已报名</a-option>
            <a-option value="cancelled">已取消</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:visible="checkinVisible" title="新增签到" width="520px" @before-ok="handleCheckinSave">
      <a-form :model="checkinModel" layout="vertical">
        <a-form-item field="userId" label="用户 ID">
          <a-input-number v-model="checkinModel.userId" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item field="checkinAt" label="签到时间">
          <a-input v-model="checkinModel.checkinAt" placeholder="留空则使用当前时间" />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:visible="outcomeVisible" :title="outcomeModel.id ? '编辑成果' : '新增成果'" width="680px" @before-ok="handleOutcomeSave">
      <a-form :model="outcomeModel" layout="vertical">
        <a-form-item field="title" label="成果标题" required>
          <a-input v-model="outcomeModel.title" placeholder="请输入成果标题" />
        </a-form-item>
        <a-form-item field="content" label="成果内容">
          <a-textarea v-model="outcomeModel.content" :auto-size="{ minRows: 3, maxRows: 6 }" placeholder="请输入成果内容" />
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item field="fileId" label="附件文件 ID">
              <a-input-number v-model="outcomeModel.fileId" :min="0" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="status" label="状态">
              <a-select v-model="outcomeModel.status">
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
  addActivity,
  addActivityCheckin,
  addActivityOutcome,
  addActivitySignup,
  getActivities,
  getActivityCheckins,
  getActivityOutcomes,
  getActivitySignups,
  removeActivities,
  removeActivityCheckins,
  removeActivityOutcomes,
  removeActivitySignups,
  updateActivity,
  updateActivityOutcome,
  updateActivitySignup
} from '@/api/edu/activity';

const statusOptions = [
  { label: '草稿', value: 'draft' },
  { label: '已发布', value: 'published' },
  { label: '已结束', value: 'finished' }
];
const statusText = Object.fromEntries(statusOptions.map((item) => [item.value, item.label]));
const statusColor = { draft: 'gray', published: 'green', finished: 'blue' };
const queryForm = reactive({ keyword: '', status: '', pageIndex: 1, pageSize: 10 });
const tableData = ref([]);
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const formVisible = ref(false);
const formModel = reactive(defaultForm());
const manageVisible = ref(false);
const signupVisible = ref(false);
const checkinVisible = ref(false);
const outcomeVisible = ref(false);
const currentActivity = ref(null);
const signupList = ref([]);
const checkinList = ref([]);
const outcomeList = ref([]);
const signupModel = reactive(defaultSignupForm());
const checkinModel = reactive(defaultCheckinForm());
const outcomeModel = reactive(defaultOutcomeForm());

const columns = [
  { title: '活动名称', dataIndex: 'title', ellipsis: true, tooltip: true },
  { title: '主办方', dataIndex: 'organizer', width: 130 },
  { title: '开始时间', dataIndex: 'startTime', width: 160 },
  { title: '地点', dataIndex: 'location', width: 150 },
  { title: '报名人数', dataIndex: 'signupCount', width: 100 },
  { title: '状态', slotName: 'status', width: 110 },
  { title: '操作', slotName: 'operations', width: 220 }
];
const signupColumns = [
  { title: '姓名', dataIndex: 'name' },
  { title: '电话', dataIndex: 'phone', width: 150 },
  { title: '状态', dataIndex: 'status', width: 110 },
  { title: '操作', slotName: 'operations', width: 150 }
];
const checkinColumns = [
  { title: '用户 ID', dataIndex: 'userId', width: 120 },
  { title: '签到时间', dataIndex: 'checkinAt' },
  { title: '状态', dataIndex: 'status', width: 110 },
  { title: '操作', slotName: 'operations', width: 100 }
];
const outcomeColumns = [
  { title: '成果标题', dataIndex: 'title', ellipsis: true, tooltip: true },
  { title: '附件文件 ID', dataIndex: 'fileId', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'operations', width: 150 }
];

function defaultForm() {
  return { id: undefined, title: '', summary: '', startTime: '', endTime: '', location: '', organizer: '', status: 'draft' };
}

function assignForm(data = {}) {
  Object.assign(formModel, defaultForm(), data);
}

function defaultSignupForm() {
  return { id: undefined, name: '', phone: '', status: 'signed' };
}

function defaultCheckinForm() {
  return { userId: 0, checkinAt: '', status: 'checked' };
}

function defaultOutcomeForm() {
  return { id: undefined, title: '', content: '', fileId: 0, status: 1 };
}

async function fetchData() {
  const res = await getActivities(queryForm);
  const payload = res.data || {};
  tableData.value = payload.list || payload || [];
  pagination.total = payload.count || res.total || 0;
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
    Message.warning('请输入活动名称');
    return false;
  }
  const payload = { ...formModel };
  if (payload.id) {
    await updateActivity(payload.id, payload);
  } else {
    await addActivity(payload);
  }
  Message.success('保存成功');
  formVisible.value = false;
  fetchData();
}

function handleDelete(record) {
  Modal.confirm({
    title: '确认删除活动',
    content: `确定删除「${record.title}」吗？`,
    async onOk() {
      await removeActivities({ ids: [record.id] });
      Message.success('删除成功');
      fetchData();
    }
  });
}

async function openManage(record) {
  currentActivity.value = record;
  manageVisible.value = true;
  await fetchManageData();
}

async function fetchManageData() {
  if (!currentActivity.value?.id) return;
  const [signupsRes, checkinsRes, outcomesRes] = await Promise.all([
    getActivitySignups(currentActivity.value.id),
    getActivityCheckins(currentActivity.value.id),
    getActivityOutcomes(currentActivity.value.id)
  ]);
  signupList.value = signupsRes.data || [];
  checkinList.value = checkinsRes.data || [];
  outcomeList.value = outcomesRes.data || [];
}

function openSignupCreate() {
  Object.assign(signupModel, defaultSignupForm());
  signupVisible.value = true;
}

function openSignupEdit(record) {
  Object.assign(signupModel, defaultSignupForm(), record);
  signupVisible.value = true;
}

async function handleSignupSave() {
  if (!signupModel.name) {
    Message.warning('请输入姓名');
    return false;
  }
  const payload = { ...signupModel };
  if (payload.id) {
    await updateActivitySignup(currentActivity.value.id, payload.id, payload);
  } else {
    await addActivitySignup(currentActivity.value.id, payload);
  }
  Message.success('保存成功');
  signupVisible.value = false;
  fetchManageData();
  fetchData();
}

function handleSignupDelete(record) {
  Modal.confirm({
    title: '确认删除报名',
    content: `确定删除「${record.name}」的报名记录吗？`,
    async onOk() {
      await removeActivitySignups(currentActivity.value.id, { ids: [record.id] });
      Message.success('删除成功');
      fetchManageData();
      fetchData();
    }
  });
}

function openCheckinCreate() {
  Object.assign(checkinModel, defaultCheckinForm());
  checkinVisible.value = true;
}

async function handleCheckinSave() {
  await addActivityCheckin(currentActivity.value.id, { ...checkinModel });
  Message.success('签到成功');
  checkinVisible.value = false;
  fetchManageData();
}

function handleCheckinDelete(record) {
  Modal.confirm({
    title: '确认删除签到',
    content: '确定删除该签到记录吗？',
    async onOk() {
      await removeActivityCheckins(currentActivity.value.id, { ids: [record.id] });
      Message.success('删除成功');
      fetchManageData();
    }
  });
}

function openOutcomeCreate() {
  Object.assign(outcomeModel, defaultOutcomeForm());
  outcomeVisible.value = true;
}

function openOutcomeEdit(record) {
  Object.assign(outcomeModel, defaultOutcomeForm(), record);
  outcomeVisible.value = true;
}

async function handleOutcomeSave() {
  if (!outcomeModel.title) {
    Message.warning('请输入成果标题');
    return false;
  }
  const payload = { ...outcomeModel };
  if (payload.id) {
    await updateActivityOutcome(currentActivity.value.id, payload.id, payload);
  } else {
    await addActivityOutcome(currentActivity.value.id, payload);
  }
  Message.success('保存成功');
  outcomeVisible.value = false;
  fetchManageData();
}

function handleOutcomeDelete(record) {
  Modal.confirm({
    title: '确认删除成果',
    content: `确定删除「${record.title}」吗？`,
    async onOk() {
      await removeActivityOutcomes(currentActivity.value.id, { ids: [record.id] });
      Message.success('删除成功');
      fetchManageData();
    }
  });
}

onMounted(fetchData);
</script>

<style scoped>
.table-card {
  margin-top: 16px;
}
</style>
