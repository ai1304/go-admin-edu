<template>
  <PortalLayout>
    <section class="page-heading">
      <h1>教师工作台</h1>
      <p>查看授权范围内的特教个案、IEP、评估记录、干预方案和案例附件。</p>
      <a-input-search v-model="query.keyword" placeholder="搜索案例、学生姓名或编号" search-button @search="fetchCases" />
    </section>

    <a-spin :loading="loading" style="width: 100%">
      <div v-if="cases.length" class="workbench-layout">
        <aside class="workbench-list">
          <button v-for="item in cases" :key="item.id" class="case-list-item" :class="{ active: selectedCase?.id === item.id }" @click="selectCase(item)">
            <strong>{{ item.title }}</strong>
            <span>{{ item.studentName || "未填写学生姓名" }} · {{ statusText[item.status] || item.status }}</span>
          </button>
        </aside>

        <section class="workbench-panel" v-if="selectedCase">
          <div class="workbench-title">
            <div>
              <h2>{{ selectedCase.title }}</h2>
              <p>{{ selectedCase.summary || "暂无案例摘要。" }}</p>
            </div>
            <a-tag :color="statusColor[selectedCase.status]">{{ statusText[selectedCase.status] || selectedCase.status }}</a-tag>
          </div>
          <a-descriptions :column="3" bordered size="small">
            <a-descriptions-item label="学生姓名">{{ selectedCase.studentName || "-" }}</a-descriptions-item>
            <a-descriptions-item label="学生编号">{{ selectedCase.studentCode || "-" }}</a-descriptions-item>
            <a-descriptions-item label="障碍类型">{{ selectedCase.disabilityType || "-" }}</a-descriptions-item>
          </a-descriptions>

          <a-tabs class="workbench-tabs" default-active-key="ieps">
            <a-tab-pane key="ieps" title="IEP">
              <a-button class="tab-action" type="primary" status="success" @click="openIepCreate">新增 IEP</a-button>
              <div v-if="ieps.length" class="record-list">
                <article v-for="item in ieps" :key="item.id" class="record-item">
                  <strong>{{ item.title }}</strong>
                  <p>{{ item.goal || item.plan || "暂无 IEP 内容。" }}</p>
                  <div class="record-footer">
                    <small>{{ subStatusText[item.status] || item.status }}</small>
                    <a-space>
                      <a-button size="mini" @click="openIepEdit(item)">编辑</a-button>
                      <a-button size="mini" status="danger" @click="deleteIep(item)">删除</a-button>
                    </a-space>
                  </div>
                </article>
              </div>
              <a-empty v-else description="暂无 IEP" />
            </a-tab-pane>
            <a-tab-pane key="assessments" title="评估记录">
              <a-button class="tab-action" type="primary" status="success" @click="openAssessmentCreate">新增评估</a-button>
              <div v-if="assessments.length" class="record-list">
                <article v-for="item in assessments" :key="item.id" class="record-item">
                  <strong>{{ item.toolName }}</strong>
                  <p>{{ item.result || "暂无评估结果。" }}</p>
                  <div class="record-footer">
                    <small>{{ item.assessedAt || "未填写评估时间" }}</small>
                    <a-space>
                      <a-button size="mini" @click="openAssessmentEdit(item)">编辑</a-button>
                      <a-button size="mini" status="danger" @click="deleteAssessment(item)">删除</a-button>
                    </a-space>
                  </div>
                </article>
              </div>
              <a-empty v-else description="暂无评估记录" />
            </a-tab-pane>
            <a-tab-pane key="interventions" title="干预方案">
              <a-button class="tab-action" type="primary" status="success" @click="openInterventionCreate">新增干预</a-button>
              <div v-if="interventions.length" class="record-list">
                <article v-for="item in interventions" :key="item.id" class="record-item">
                  <strong>{{ item.title }}</strong>
                  <p>{{ item.content || "暂无干预内容。" }}</p>
                  <div class="record-footer">
                    <small>{{ item.startDate || "-" }} 至 {{ item.endDate || "-" }}</small>
                    <a-space>
                      <a-button size="mini" @click="openInterventionEdit(item)">编辑</a-button>
                      <a-button size="mini" status="danger" @click="deleteIntervention(item)">删除</a-button>
                    </a-space>
                  </div>
                </article>
              </div>
              <a-empty v-else description="暂无干预方案" />
            </a-tab-pane>
            <a-tab-pane key="attachments" title="案例附件">
              <div v-if="attachments.length" class="file-list">
                <div v-for="item in attachments" :key="item.id" class="file-item">
                  <div>
                    <strong>{{ item.title }}</strong>
                    <span>{{ item.remark || `文件 ID：${item.fileId}` }}</span>
                  </div>
                  <a-button size="small" @click="openAttachment(item)">打开</a-button>
                </div>
              </div>
              <a-empty v-else description="暂无案例附件" />
            </a-tab-pane>
          </a-tabs>
        </section>
      </div>
      <a-empty v-else description="暂无可访问案例，请确认已登录并具备个案访问授权" />
    </a-spin>

    <a-modal v-model:visible="iepVisible" :title="iepForm.id ? '编辑 IEP' : '新增 IEP'" width="680px" @before-ok="saveIep">
      <a-form :model="iepForm" layout="vertical">
        <a-form-item field="title" label="IEP 标题" required>
          <a-input v-model="iepForm.title" placeholder="请输入 IEP 标题" />
        </a-form-item>
        <a-form-item field="goal" label="目标">
          <a-textarea v-model="iepForm.goal" :auto-size="{ minRows: 3, maxRows: 6 }" placeholder="请输入阶段目标" />
        </a-form-item>
        <a-form-item field="plan" label="计划">
          <a-textarea v-model="iepForm.plan" :auto-size="{ minRows: 3, maxRows: 6 }" placeholder="请输入实施计划" />
        </a-form-item>
        <a-form-item field="evaluation" label="评价">
          <a-textarea v-model="iepForm.evaluation" :auto-size="{ minRows: 3, maxRows: 6 }" placeholder="请输入评价记录" />
        </a-form-item>
        <a-form-item field="status" label="状态">
          <a-select v-model="iepForm.status">
            <a-option value="draft">草稿</a-option>
            <a-option value="active">执行中</a-option>
            <a-option value="finished">已完成</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:visible="assessmentVisible" :title="assessmentForm.id ? '编辑评估' : '新增评估'" width="620px" @before-ok="saveAssessment">
      <a-form :model="assessmentForm" layout="vertical">
        <a-form-item field="toolName" label="评估工具" required>
          <a-input v-model="assessmentForm.toolName" placeholder="请输入评估工具或量表名称" />
        </a-form-item>
        <a-form-item field="assessedAt" label="评估时间">
          <a-input v-model="assessmentForm.assessedAt" placeholder="例如 2026-05-16" />
        </a-form-item>
        <a-form-item field="result" label="评估结果">
          <a-textarea v-model="assessmentForm.result" :auto-size="{ minRows: 4, maxRows: 8 }" placeholder="请输入评估结果" />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:visible="interventionVisible" :title="interventionForm.id ? '编辑干预方案' : '新增干预方案'" width="680px" @before-ok="saveIntervention">
      <a-form :model="interventionForm" layout="vertical">
        <a-form-item field="title" label="干预标题" required>
          <a-input v-model="interventionForm.title" placeholder="请输入干预标题" />
        </a-form-item>
        <a-form-item field="content" label="干预内容">
          <a-textarea v-model="interventionForm.content" :auto-size="{ minRows: 4, maxRows: 8 }" placeholder="请输入干预内容" />
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="8">
            <a-form-item field="startDate" label="开始日期">
              <a-input v-model="interventionForm.startDate" placeholder="2026-05-16" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item field="endDate" label="结束日期">
              <a-input v-model="interventionForm.endDate" placeholder="2026-06-16" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item field="status" label="状态">
              <a-select v-model="interventionForm.status">
                <a-option value="active">执行中</a-option>
                <a-option value="paused">暂停</a-option>
                <a-option value="finished">已完成</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>
  </PortalLayout>
</template>

<script setup>
import { Message, Modal } from "@arco-design/web-vue";
import { onMounted, reactive, ref } from "vue";
import {
  addTeacherCaseAssessment,
  addTeacherCaseIep,
  addTeacherCaseIntervention,
  getTeacherCaseAssessments,
  getTeacherCaseAttachmentUrl,
  getTeacherCaseAttachments,
  getTeacherCaseIeps,
  getTeacherCaseInterventions,
  getTeacherCases,
  removeTeacherCaseAssessments,
  removeTeacherCaseIeps,
  removeTeacherCaseInterventions,
  updateTeacherCaseAssessment,
  updateTeacherCaseIep,
  updateTeacherCaseIntervention
} from "@/api/cases";
import PortalLayout from "@/layouts/PortalLayout.vue";

const statusText = { draft: "草稿", reviewing: "审核中", rejected: "已驳回", archived: "已归档" };
const statusColor = { draft: "gray", reviewing: "orange", rejected: "red", archived: "blue" };
const subStatusText = { draft: "草稿", active: "执行中", finished: "已完成", paused: "暂停" };

const loading = ref(false);
const cases = ref([]);
const selectedCase = ref(null);
const ieps = ref([]);
const assessments = ref([]);
const interventions = ref([]);
const attachments = ref([]);
const query = reactive({ keyword: "", pageIndex: 1, pageSize: 20, desensitize: false });
const iepVisible = ref(false);
const assessmentVisible = ref(false);
const interventionVisible = ref(false);
const iepForm = reactive(defaultIepForm());
const assessmentForm = reactive(defaultAssessmentForm());
const interventionForm = reactive(defaultInterventionForm());

function defaultIepForm() {
  return { id: undefined, title: "", goal: "", plan: "", evaluation: "", status: "draft" };
}

function defaultAssessmentForm() {
  return { id: undefined, toolName: "", result: "", assessedAt: "" };
}

function defaultInterventionForm() {
  return { id: undefined, title: "", content: "", startDate: "", endDate: "", status: "active" };
}

async function fetchCases() {
  loading.value = true;
  try {
    const res = await getTeacherCases(query);
    cases.value = res.data?.list || res.data || [];
    if (cases.value.length) {
      await selectCase(cases.value[0]);
    } else {
      selectedCase.value = null;
      resetCaseDetail();
    }
  } catch (error) {
    Message.warning("暂时无法加载个案，请确认登录状态和访问授权");
  } finally {
    loading.value = false;
  }
}

async function selectCase(item) {
  selectedCase.value = item;
  resetCaseDetail();
  const params = { desensitize: false };
  try {
    const [iepRes, assessmentRes, interventionRes, attachmentRes] = await Promise.all([
      getTeacherCaseIeps(item.id, params),
      getTeacherCaseAssessments(item.id, params),
      getTeacherCaseInterventions(item.id, params),
      getTeacherCaseAttachments(item.id)
    ]);
    ieps.value = iepRes.data || [];
    assessments.value = assessmentRes.data || [];
    interventions.value = interventionRes.data || [];
    attachments.value = attachmentRes.data || [];
  } catch (error) {
    Message.warning("暂时无法加载该案例详情，请确认访问授权");
  }
}

function resetCaseDetail() {
  ieps.value = [];
  assessments.value = [];
  interventions.value = [];
  attachments.value = [];
}

async function openAttachment(item) {
  const res = await getTeacherCaseAttachmentUrl(selectedCase.value.id, item.id);
  const url = res.data?.url;
  if (url) {
    window.open(url, "_blank");
  }
}

function openIepCreate() {
  Object.assign(iepForm, defaultIepForm());
  iepVisible.value = true;
}

function openIepEdit(item) {
  Object.assign(iepForm, defaultIepForm(), item);
  iepVisible.value = true;
}

async function saveIep() {
  if (!iepForm.title) {
    Message.warning("请输入 IEP 标题");
    return false;
  }
  const payload = { ...iepForm };
  if (payload.id) {
    await updateTeacherCaseIep(selectedCase.value.id, payload.id, payload);
  } else {
    await addTeacherCaseIep(selectedCase.value.id, payload);
  }
  Message.success("IEP 已保存");
  iepVisible.value = false;
  selectCase(selectedCase.value);
}

function deleteIep(item) {
  Modal.confirm({
    title: "确认删除 IEP",
    content: `确定删除「${item.title}」吗？`,
    async onOk() {
      await removeTeacherCaseIeps(selectedCase.value.id, { ids: [item.id] });
      Message.success("删除成功");
      selectCase(selectedCase.value);
    }
  });
}

function openAssessmentCreate() {
  Object.assign(assessmentForm, defaultAssessmentForm());
  assessmentVisible.value = true;
}

function openAssessmentEdit(item) {
  Object.assign(assessmentForm, defaultAssessmentForm(), item);
  assessmentVisible.value = true;
}

async function saveAssessment() {
  if (!assessmentForm.toolName) {
    Message.warning("请输入评估工具");
    return false;
  }
  const payload = { ...assessmentForm };
  if (payload.id) {
    await updateTeacherCaseAssessment(selectedCase.value.id, payload.id, payload);
  } else {
    await addTeacherCaseAssessment(selectedCase.value.id, payload);
  }
  Message.success("评估记录已保存");
  assessmentVisible.value = false;
  selectCase(selectedCase.value);
}

function deleteAssessment(item) {
  Modal.confirm({
    title: "确认删除评估记录",
    content: `确定删除「${item.toolName}」吗？`,
    async onOk() {
      await removeTeacherCaseAssessments(selectedCase.value.id, { ids: [item.id] });
      Message.success("删除成功");
      selectCase(selectedCase.value);
    }
  });
}

function openInterventionCreate() {
  Object.assign(interventionForm, defaultInterventionForm());
  interventionVisible.value = true;
}

function openInterventionEdit(item) {
  Object.assign(interventionForm, defaultInterventionForm(), item);
  interventionVisible.value = true;
}

async function saveIntervention() {
  if (!interventionForm.title) {
    Message.warning("请输入干预标题");
    return false;
  }
  const payload = { ...interventionForm };
  if (payload.id) {
    await updateTeacherCaseIntervention(selectedCase.value.id, payload.id, payload);
  } else {
    await addTeacherCaseIntervention(selectedCase.value.id, payload);
  }
  Message.success("干预方案已保存");
  interventionVisible.value = false;
  selectCase(selectedCase.value);
}

function deleteIntervention(item) {
  Modal.confirm({
    title: "确认删除干预方案",
    content: `确定删除「${item.title}」吗？`,
    async onOk() {
      await removeTeacherCaseInterventions(selectedCase.value.id, { ids: [item.id] });
      Message.success("删除成功");
      selectCase(selectedCase.value);
    }
  });
}

onMounted(fetchCases);
</script>
