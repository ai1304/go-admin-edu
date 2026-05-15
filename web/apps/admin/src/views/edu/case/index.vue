<template>
  <div class="container">
    <a-card :bordered="false" class="cardStyle">
      <template #title>特教案例管理</template>
      <a-form :model="queryForm" layout="inline">
        <a-form-item label="关键词">
          <a-input v-model="queryForm.keyword" allow-clear placeholder="案例、学生姓名或编号" />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model="queryForm.status" allow-clear placeholder="请选择状态" style="width: 160px">
            <a-option v-for="item in statusOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="脱敏">
          <a-switch v-model="queryForm.desensitize" />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="fetchData">查询</a-button>
            <a-button @click="resetQuery">重置</a-button>
            <a-button type="primary" status="success" @click="openCreate">新增案例</a-button>
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

    <a-modal v-model:visible="formVisible" :title="formModel.id ? '编辑案例' : '新增案例'" width="760px" @before-ok="handleSave">
      <a-form :model="formModel" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item field="title" label="案例名称" required>
              <a-input v-model="formModel.title" placeholder="请输入案例名称" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="studentName" label="学生姓名">
              <a-input v-model="formModel.studentName" placeholder="请输入学生姓名" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="studentCode" label="学生编号">
              <a-input v-model="formModel.studentCode" placeholder="请输入学生编号" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="gender" label="性别">
              <a-select v-model="formModel.gender" allow-clear placeholder="请选择性别">
                <a-option value="male">男</a-option>
                <a-option value="female">女</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="birthday" label="生日">
              <a-input v-model="formModel.birthday" placeholder="例如 2015-09-01" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="disabilityType" label="障碍类型">
              <a-input v-model="formModel.disabilityType" placeholder="请输入障碍类型" />
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
            <a-form-item field="summary" label="案例摘要">
              <a-textarea v-model="formModel.summary" :auto-size="{ minRows: 3, maxRows: 6 }" placeholder="请输入案例摘要" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>

    <a-modal v-model:visible="manageVisible" :title="`${currentCase?.title || ''} 业务管理`" width="1040px" :footer="false">
      <a-space class="manage-toolbar">
        <span>脱敏查看</span>
        <a-switch v-model="manageDesensitize" @change="fetchManageData" />
      </a-space>
      <a-tabs default-active-key="ieps">
        <a-tab-pane key="ieps" title="IEP">
          <a-space direction="vertical" fill>
            <a-button type="primary" status="success" @click="openIepCreate">新增 IEP</a-button>
            <a-table :columns="iepColumns" :data="iepList" :pagination="false" row-key="id">
              <template #status="{ record }">
                <a-tag :color="subStatusColor[record.status]">{{ subStatusText[record.status] || record.status }}</a-tag>
              </template>
              <template #operations="{ record }">
                <a-space>
                  <a-button type="text" size="small" @click="openIepEdit(record)">编辑</a-button>
                  <a-button type="text" status="danger" size="small" @click="handleIepDelete(record)">删除</a-button>
                </a-space>
              </template>
            </a-table>
          </a-space>
        </a-tab-pane>
        <a-tab-pane key="assessments" title="评估记录">
          <a-space direction="vertical" fill>
            <a-button type="primary" status="success" @click="openAssessmentCreate">新增评估</a-button>
            <a-table :columns="assessmentColumns" :data="assessmentList" :pagination="false" row-key="id">
              <template #operations="{ record }">
                <a-space>
                  <a-button type="text" size="small" @click="openAssessmentEdit(record)">编辑</a-button>
                  <a-button type="text" status="danger" size="small" @click="handleAssessmentDelete(record)">删除</a-button>
                </a-space>
              </template>
            </a-table>
          </a-space>
        </a-tab-pane>
        <a-tab-pane key="interventions" title="干预方案">
          <a-space direction="vertical" fill>
            <a-button type="primary" status="success" @click="openInterventionCreate">新增干预</a-button>
            <a-table :columns="interventionColumns" :data="interventionList" :pagination="false" row-key="id">
              <template #status="{ record }">
                <a-tag :color="interventionStatusColor[record.status]">{{ interventionStatusText[record.status] || record.status }}</a-tag>
              </template>
              <template #operations="{ record }">
                <a-space>
                  <a-button type="text" size="small" @click="openInterventionEdit(record)">编辑</a-button>
                  <a-button type="text" status="danger" size="small" @click="handleInterventionDelete(record)">删除</a-button>
                </a-space>
              </template>
            </a-table>
          </a-space>
        </a-tab-pane>
        <a-tab-pane key="authorizations" title="访问授权">
          <a-space direction="vertical" fill>
            <a-form :model="authorizationQuery" layout="inline">
              <a-form-item label="用户ID">
                <a-input-number v-model="authorizationQuery.userId" allow-clear placeholder="用户 ID" style="width: 140px" />
              </a-form-item>
              <a-form-item label="范围">
                <a-select v-model="authorizationQuery.scope" allow-clear placeholder="授权范围" style="width: 150px">
                  <a-option v-for="item in authorizationScopeOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="状态">
                <a-select v-model="authorizationQuery.status" allow-clear placeholder="状态" style="width: 140px">
                  <a-option v-for="item in authorizationStatusOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
                </a-select>
              </a-form-item>
              <a-form-item>
                <a-space>
                  <a-button type="primary" @click="searchAuthorizations">查询</a-button>
                  <a-button @click="resetAuthorizationQuery">重置</a-button>
                  <a-button type="primary" status="success" @click="openAuthorizationCreate">新增授权</a-button>
                </a-space>
              </a-form-item>
            </a-form>
            <a-table
              :columns="authorizationColumns"
              :data="authorizationList"
              :pagination="authorizationPagination"
              row-key="id"
              @page-change="handleAuthorizationPageChange"
            >
              <template #scope="{ record }">{{ authorizationScopeText[record.scope] || record.scope }}</template>
              <template #status="{ record }">
                <a-tag :color="authorizationStatusColor[record.status]">{{ authorizationStatusText[record.status] || record.status }}</a-tag>
              </template>
              <template #operations="{ record }">
                <a-space>
                  <a-button type="text" size="small" @click="openAuthorizationEdit(record)">编辑</a-button>
                  <a-button type="text" status="danger" size="small" @click="handleAuthorizationDelete(record)">删除</a-button>
                </a-space>
              </template>
            </a-table>
          </a-space>
        </a-tab-pane>
        <a-tab-pane key="accessLogs" title="访问日志">
          <a-space direction="vertical" fill>
            <a-form :model="accessLogQuery" layout="inline">
              <a-form-item label="动作">
                <a-select v-model="accessLogQuery.action" allow-clear placeholder="请选择动作" style="width: 160px">
                  <a-option v-for="item in accessActionOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="用户ID">
                <a-input-number v-model="accessLogQuery.userId" allow-clear placeholder="用户 ID" style="width: 140px" />
              </a-form-item>
              <a-form-item label="关键词">
                <a-input v-model="accessLogQuery.keyword" allow-clear placeholder="IP、路径或 User-Agent" />
              </a-form-item>
              <a-form-item>
                <a-space>
                  <a-button type="primary" @click="searchAccessLogs">查询</a-button>
                  <a-button @click="resetAccessLogQuery">重置</a-button>
                </a-space>
              </a-form-item>
            </a-form>
            <a-table
              :columns="accessLogColumns"
              :data="accessLogList"
              :pagination="accessLogPagination"
              row-key="id"
              @page-change="handleAccessLogPageChange"
            >
              <template #action="{ record }">{{ accessActionText[record.action] || record.action }}</template>
            </a-table>
          </a-space>
        </a-tab-pane>
      </a-tabs>
    </a-modal>

    <a-modal v-model:visible="iepVisible" :title="iepModel.id ? '编辑 IEP' : '新增 IEP'" width="760px" @before-ok="handleIepSave">
      <a-form :model="iepModel" layout="vertical">
        <a-form-item field="title" label="IEP 标题" required>
          <a-input v-model="iepModel.title" placeholder="请输入 IEP 标题" />
        </a-form-item>
        <a-form-item field="goal" label="目标">
          <a-textarea v-model="iepModel.goal" :auto-size="{ minRows: 3, maxRows: 6 }" placeholder="请输入阶段目标" />
        </a-form-item>
        <a-form-item field="plan" label="计划">
          <a-textarea v-model="iepModel.plan" :auto-size="{ minRows: 3, maxRows: 6 }" placeholder="请输入实施计划" />
        </a-form-item>
        <a-form-item field="evaluation" label="评价">
          <a-textarea v-model="iepModel.evaluation" :auto-size="{ minRows: 3, maxRows: 6 }" placeholder="请输入评价记录" />
        </a-form-item>
        <a-form-item field="status" label="状态">
          <a-select v-model="iepModel.status">
            <a-option v-for="item in subStatusOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:visible="assessmentVisible" :title="assessmentModel.id ? '编辑评估' : '新增评估'" width="680px" @before-ok="handleAssessmentSave">
      <a-form :model="assessmentModel" layout="vertical">
        <a-form-item field="toolName" label="评估工具" required>
          <a-input v-model="assessmentModel.toolName" placeholder="请输入评估工具或量表名称" />
        </a-form-item>
        <a-form-item field="assessedAt" label="评估时间">
          <a-input v-model="assessmentModel.assessedAt" placeholder="例如 2026-05-15" />
        </a-form-item>
        <a-form-item field="result" label="评估结果">
          <a-textarea v-model="assessmentModel.result" :auto-size="{ minRows: 4, maxRows: 8 }" placeholder="请输入评估结果" />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal
      v-model:visible="interventionVisible"
      :title="interventionModel.id ? '编辑干预方案' : '新增干预方案'"
      width="760px"
      @before-ok="handleInterventionSave"
    >
      <a-form :model="interventionModel" layout="vertical">
        <a-form-item field="title" label="干预标题" required>
          <a-input v-model="interventionModel.title" placeholder="请输入干预标题" />
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item field="startDate" label="开始日期">
              <a-input v-model="interventionModel.startDate" placeholder="例如 2026-05-15" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item field="endDate" label="结束日期">
              <a-input v-model="interventionModel.endDate" placeholder="例如 2026-06-15" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item field="status" label="状态">
              <a-select v-model="interventionModel.status">
                <a-option v-for="item in interventionStatusOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item field="content" label="干预内容">
          <a-textarea v-model="interventionModel.content" :auto-size="{ minRows: 4, maxRows: 8 }" placeholder="请输入干预内容" />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal
      v-model:visible="authorizationVisible"
      :title="authorizationModel.id ? '编辑访问授权' : '新增访问授权'"
      width="680px"
      @before-ok="handleAuthorizationSave"
    >
      <a-form :model="authorizationModel" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item field="userId" label="授权用户 ID" required>
              <a-input-number v-model="authorizationModel.userId" placeholder="请输入用户 ID" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="scope" label="授权范围">
              <a-select v-model="authorizationModel.scope">
                <a-option v-for="item in authorizationScopeOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="startAt" label="生效时间">
              <a-input v-model="authorizationModel.startAt" placeholder="例如 2026-05-15" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="endAt" label="失效时间">
              <a-input v-model="authorizationModel.endAt" placeholder="例如 2026-06-15" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="status" label="状态">
              <a-select v-model="authorizationModel.status">
                <a-option v-for="item in authorizationStatusOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item field="remark" label="备注">
              <a-textarea v-model="authorizationModel.remark" :auto-size="{ minRows: 3, maxRows: 5 }" placeholder="请输入授权说明" />
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
  addCase,
  addCaseAssessment,
  addCaseAuthorization,
  addCaseIep,
  addCaseIntervention,
  getCaseAccessLogs,
  getCaseAssessments,
  getCaseAuthorizations,
  getCaseIeps,
  getCaseInterventions,
  getCases,
  removeCaseAssessments,
  removeCaseAuthorizations,
  removeCaseIeps,
  removeCaseInterventions,
  removeCases,
  updateCase,
  updateCaseAssessment,
  updateCaseAuthorization,
  updateCaseIep,
  updateCaseIntervention
} from '@/api/edu/case';

const statusOptions = [
  { label: '草稿', value: 'draft' },
  { label: '审核中', value: 'reviewing' },
  { label: '已归档', value: 'archived' }
];
const subStatusOptions = [
  { label: '草稿', value: 'draft' },
  { label: '执行中', value: 'active' },
  { label: '已完成', value: 'finished' }
];
const interventionStatusOptions = [
  { label: '执行中', value: 'active' },
  { label: '暂停', value: 'paused' },
  { label: '已完成', value: 'finished' }
];
const statusText = Object.fromEntries(statusOptions.map((item) => [item.value, item.label]));
const statusColor = { draft: 'gray', reviewing: 'orange', archived: 'blue' };
const subStatusText = Object.fromEntries(subStatusOptions.map((item) => [item.value, item.label]));
const subStatusColor = { draft: 'gray', active: 'green', finished: 'blue' };
const interventionStatusText = Object.fromEntries(interventionStatusOptions.map((item) => [item.value, item.label]));
const interventionStatusColor = { active: 'green', paused: 'orange', finished: 'blue' };
const queryForm = reactive({ keyword: '', status: '', desensitize: false, pageIndex: 1, pageSize: 10 });
const tableData = ref([]);
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const formVisible = ref(false);
const manageVisible = ref(false);
const iepVisible = ref(false);
const assessmentVisible = ref(false);
const interventionVisible = ref(false);
const authorizationVisible = ref(false);
const manageDesensitize = ref(false);
const currentCase = ref(null);
const formModel = reactive(defaultForm());
const iepModel = reactive(defaultIepForm());
const assessmentModel = reactive(defaultAssessmentForm());
const interventionModel = reactive(defaultInterventionForm());
const authorizationModel = reactive(defaultAuthorizationForm());
const iepList = ref([]);
const assessmentList = ref([]);
const interventionList = ref([]);
const authorizationList = ref([]);
const accessLogList = ref([]);
const accessActionText = {
  view_detail: '查看详情',
  view_ieps: '查看 IEP',
  view_assessments: '查看评估',
  view_interventions: '查看干预',
  view_detail_denied: '拒绝查看详情',
  edit_case_denied: '拒绝编辑案例',
  delete_case_denied: '拒绝删除案例',
  view_access_logs_denied: '拒绝查看日志',
  view_authorizations_denied: '拒绝查看授权',
  add_authorization_denied: '拒绝新增授权',
  update_authorization_denied: '拒绝编辑授权',
  delete_authorization_denied: '拒绝删除授权',
  view_ieps_denied: '拒绝查看 IEP',
  add_iep_denied: '拒绝新增 IEP',
  update_iep_denied: '拒绝编辑 IEP',
  delete_ieps_denied: '拒绝删除 IEP',
  view_assessments_denied: '拒绝查看评估',
  add_assessment_denied: '拒绝新增评估',
  update_assessment_denied: '拒绝编辑评估',
  delete_assessments_denied: '拒绝删除评估',
  view_interventions_denied: '拒绝查看干预',
  add_intervention_denied: '拒绝新增干预',
  update_intervention_denied: '拒绝编辑干预',
  delete_interventions_denied: '拒绝删除干预'
};
const accessActionOptions = Object.entries(accessActionText).map(([value, label]) => ({ value, label }));
const accessLogQuery = reactive({ action: '', userId: undefined, keyword: '', pageIndex: 1, pageSize: 10 });
const accessLogPagination = reactive({ current: 1, pageSize: 10, total: 0 });
const authorizationScopeOptions = [
  { label: '仅查看', value: 'view' },
  { label: '可编辑', value: 'edit' },
  { label: '可审核', value: 'review' }
];
const authorizationStatusOptions = [
  { label: '有效', value: 'active' },
  { label: '停用', value: 'disabled' },
  { label: '过期', value: 'expired' }
];
const authorizationScopeText = Object.fromEntries(authorizationScopeOptions.map((item) => [item.value, item.label]));
const authorizationStatusText = Object.fromEntries(authorizationStatusOptions.map((item) => [item.value, item.label]));
const authorizationStatusColor = { active: 'green', disabled: 'gray', expired: 'orange' };
const authorizationQuery = reactive({ userId: undefined, scope: '', status: '', pageIndex: 1, pageSize: 10 });
const authorizationPagination = reactive({ current: 1, pageSize: 10, total: 0 });

const columns = [
  { title: '案例名称', dataIndex: 'title', ellipsis: true, tooltip: true },
  { title: '学生姓名', dataIndex: 'studentName', width: 120 },
  { title: '学生编号', dataIndex: 'studentCode', width: 130 },
  { title: '障碍类型', dataIndex: 'disabilityType', width: 130 },
  { title: '状态', slotName: 'status', width: 110 },
  { title: '操作', slotName: 'operations', width: 220 }
];
const iepColumns = [
  { title: 'IEP 标题', dataIndex: 'title', ellipsis: true, tooltip: true },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'operations', width: 150 }
];
const assessmentColumns = [
  { title: '评估工具', dataIndex: 'toolName', width: 180 },
  { title: '评估时间', dataIndex: 'assessedAt', width: 130 },
  { title: '评估结果', dataIndex: 'result', ellipsis: true, tooltip: true },
  { title: '操作', slotName: 'operations', width: 150 }
];
const interventionColumns = [
  { title: '干预标题', dataIndex: 'title', ellipsis: true, tooltip: true },
  { title: '开始日期', dataIndex: 'startDate', width: 120 },
  { title: '结束日期', dataIndex: 'endDate', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'operations', width: 150 }
];
const authorizationColumns = [
  { title: '用户 ID', dataIndex: 'userId', width: 100 },
  { title: '范围', slotName: 'scope', width: 100 },
  { title: '生效时间', dataIndex: 'startAt', width: 120 },
  { title: '失效时间', dataIndex: 'endAt', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '备注', dataIndex: 'remark', ellipsis: true, tooltip: true },
  { title: '操作', slotName: 'operations', width: 150 }
];
const accessLogColumns = [
  { title: '动作', slotName: 'action', width: 120 },
  { title: '用户 ID', dataIndex: 'userId', width: 100 },
  { title: 'IP', dataIndex: 'ip', width: 140 },
  { title: '路径', dataIndex: 'path', ellipsis: true, tooltip: true },
  { title: 'User-Agent', dataIndex: 'userAgent', ellipsis: true, tooltip: true },
  { title: '时间', dataIndex: 'createdAt', width: 170 }
];

function defaultForm() {
  return {
    id: undefined,
    title: '',
    studentName: '',
    studentCode: '',
    gender: '',
    birthday: '',
    disabilityType: '',
    summary: '',
    status: 'draft'
  };
}

function defaultIepForm() {
  return { id: undefined, title: '', goal: '', plan: '', evaluation: '', status: 'draft' };
}

function defaultAssessmentForm() {
  return { id: undefined, toolName: '', result: '', assessedAt: '' };
}

function defaultInterventionForm() {
  return { id: undefined, title: '', content: '', startDate: '', endDate: '', status: 'active' };
}

function defaultAuthorizationForm() {
  return { id: undefined, userId: undefined, scope: 'view', startAt: '', endAt: '', status: 'active', remark: '' };
}

function assignForm(data = {}) {
  Object.assign(formModel, defaultForm(), data);
}

async function fetchData() {
  const res = await getCases(queryForm);
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
  queryForm.desensitize = false;
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
    Message.warning('请输入案例名称');
    return false;
  }
  const payload = { ...formModel };
  if (payload.id) {
    await updateCase(payload.id, payload);
  } else {
    await addCase(payload);
  }
  Message.success('保存成功');
  formVisible.value = false;
  fetchData();
}

function handleDelete(record) {
  Modal.confirm({
    title: '确认删除案例',
    content: `确定删除「${record.title}」吗？`,
    async onOk() {
      await removeCases({ ids: [record.id] });
      Message.success('删除成功');
      fetchData();
    }
  });
}

async function openManage(record) {
  currentCase.value = record;
  manageDesensitize.value = false;
  resetAccessLogQuery(false);
  resetAuthorizationQuery(false);
  manageVisible.value = true;
  await fetchManageData();
}

async function fetchManageData() {
  if (!currentCase.value?.id) return;
  const sensitiveParams = { desensitize: manageDesensitize.value };
  const [iepsRes, assessmentsRes, interventionsRes, authorizationsRes, accessLogsRes] = await Promise.all([
    getCaseIeps(currentCase.value.id, sensitiveParams),
    getCaseAssessments(currentCase.value.id, sensitiveParams),
    getCaseInterventions(currentCase.value.id, sensitiveParams),
    getCaseAuthorizations(currentCase.value.id, authorizationQuery),
    getCaseAccessLogs(currentCase.value.id, accessLogQuery)
  ]);
  iepList.value = iepsRes.data || [];
  assessmentList.value = assessmentsRes.data || [];
  interventionList.value = interventionsRes.data || [];
  setAuthorizations(authorizationsRes);
  setAccessLogs(accessLogsRes);
}

async function fetchAccessLogs() {
  if (!currentCase.value?.id) return;
  const res = await getCaseAccessLogs(currentCase.value.id, accessLogQuery);
  setAccessLogs(res);
}

function searchAccessLogs() {
  accessLogQuery.pageIndex = 1;
  accessLogPagination.current = 1;
  fetchAccessLogs();
}

function setAccessLogs(res) {
  const payload = res.data || {};
  accessLogList.value = payload.list || payload || [];
  accessLogPagination.total = payload.count || res.total || 0;
  accessLogPagination.current = accessLogQuery.pageIndex;
}

function handleAccessLogPageChange(page) {
  accessLogQuery.pageIndex = page;
  accessLogPagination.current = page;
  fetchAccessLogs();
}

function resetAccessLogQuery(shouldFetch = true) {
  accessLogQuery.action = '';
  accessLogQuery.userId = undefined;
  accessLogQuery.keyword = '';
  accessLogQuery.pageIndex = 1;
  accessLogPagination.current = 1;
  if (shouldFetch) {
    fetchAccessLogs();
  }
}

async function fetchAuthorizations() {
  if (!currentCase.value?.id) return;
  const res = await getCaseAuthorizations(currentCase.value.id, authorizationQuery);
  setAuthorizations(res);
}

function searchAuthorizations() {
  authorizationQuery.pageIndex = 1;
  authorizationPagination.current = 1;
  fetchAuthorizations();
}

function setAuthorizations(res) {
  const payload = res.data || {};
  authorizationList.value = payload.list || payload || [];
  authorizationPagination.total = payload.count || res.total || 0;
  authorizationPagination.current = authorizationQuery.pageIndex;
}

function handleAuthorizationPageChange(page) {
  authorizationQuery.pageIndex = page;
  authorizationPagination.current = page;
  fetchAuthorizations();
}

function resetAuthorizationQuery(shouldFetch = true) {
  authorizationQuery.userId = undefined;
  authorizationQuery.scope = '';
  authorizationQuery.status = '';
  authorizationQuery.pageIndex = 1;
  authorizationPagination.current = 1;
  if (shouldFetch) {
    fetchAuthorizations();
  }
}

function openAuthorizationCreate() {
  Object.assign(authorizationModel, defaultAuthorizationForm());
  authorizationVisible.value = true;
}

function openAuthorizationEdit(record) {
  Object.assign(authorizationModel, defaultAuthorizationForm(), record);
  authorizationVisible.value = true;
}

async function handleAuthorizationSave() {
  if (!authorizationModel.userId) {
    Message.warning('请输入授权用户 ID');
    return false;
  }
  const payload = { ...authorizationModel };
  if (payload.id) {
    await updateCaseAuthorization(currentCase.value.id, payload.id, payload);
  } else {
    await addCaseAuthorization(currentCase.value.id, payload);
  }
  Message.success('保存成功');
  authorizationVisible.value = false;
  fetchAuthorizations();
}

function handleAuthorizationDelete(record) {
  Modal.confirm({
    title: '确认删除访问授权',
    content: `确定删除用户 ${record.userId} 的访问授权吗？`,
    async onOk() {
      await removeCaseAuthorizations(currentCase.value.id, { ids: [record.id] });
      Message.success('删除成功');
      fetchAuthorizations();
    }
  });
}

function openIepCreate() {
  Object.assign(iepModel, defaultIepForm());
  iepVisible.value = true;
}

function openIepEdit(record) {
  Object.assign(iepModel, defaultIepForm(), record);
  iepVisible.value = true;
}

async function handleIepSave() {
  if (!iepModel.title) {
    Message.warning('请输入 IEP 标题');
    return false;
  }
  const payload = { ...iepModel };
  if (payload.id) {
    await updateCaseIep(currentCase.value.id, payload.id, payload);
  } else {
    await addCaseIep(currentCase.value.id, payload);
  }
  Message.success('保存成功');
  iepVisible.value = false;
  fetchManageData();
}

function handleIepDelete(record) {
  Modal.confirm({
    title: '确认删除 IEP',
    content: `确定删除「${record.title}」吗？`,
    async onOk() {
      await removeCaseIeps(currentCase.value.id, { ids: [record.id] });
      Message.success('删除成功');
      fetchManageData();
    }
  });
}

function openAssessmentCreate() {
  Object.assign(assessmentModel, defaultAssessmentForm());
  assessmentVisible.value = true;
}

function openAssessmentEdit(record) {
  Object.assign(assessmentModel, defaultAssessmentForm(), record);
  assessmentVisible.value = true;
}

async function handleAssessmentSave() {
  if (!assessmentModel.toolName) {
    Message.warning('请输入评估工具');
    return false;
  }
  const payload = { ...assessmentModel };
  if (payload.id) {
    await updateCaseAssessment(currentCase.value.id, payload.id, payload);
  } else {
    await addCaseAssessment(currentCase.value.id, payload);
  }
  Message.success('保存成功');
  assessmentVisible.value = false;
  fetchManageData();
}

function handleAssessmentDelete(record) {
  Modal.confirm({
    title: '确认删除评估记录',
    content: `确定删除「${record.toolName}」吗？`,
    async onOk() {
      await removeCaseAssessments(currentCase.value.id, { ids: [record.id] });
      Message.success('删除成功');
      fetchManageData();
    }
  });
}

function openInterventionCreate() {
  Object.assign(interventionModel, defaultInterventionForm());
  interventionVisible.value = true;
}

function openInterventionEdit(record) {
  Object.assign(interventionModel, defaultInterventionForm(), record);
  interventionVisible.value = true;
}

async function handleInterventionSave() {
  if (!interventionModel.title) {
    Message.warning('请输入干预标题');
    return false;
  }
  const payload = { ...interventionModel };
  if (payload.id) {
    await updateCaseIntervention(currentCase.value.id, payload.id, payload);
  } else {
    await addCaseIntervention(currentCase.value.id, payload);
  }
  Message.success('保存成功');
  interventionVisible.value = false;
  fetchManageData();
}

function handleInterventionDelete(record) {
  Modal.confirm({
    title: '确认删除干预方案',
    content: `确定删除「${record.title}」吗？`,
    async onOk() {
      await removeCaseInterventions(currentCase.value.id, { ids: [record.id] });
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

.manage-toolbar {
  margin-bottom: 12px;
}
</style>
