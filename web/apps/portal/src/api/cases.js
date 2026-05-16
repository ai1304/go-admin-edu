import { request } from "./request";

const caseUrl = "/edu/cases";

export function getTeacherCases(params) {
  return request({ url: caseUrl, method: "get", params });
}

export function getTeacherCaseIeps(id, params) {
  return request({ url: `${caseUrl}/${id}/ieps`, method: "get", params });
}

export function addTeacherCaseIep(id, data) {
  return request({ url: `${caseUrl}/${id}/ieps`, method: "post", data });
}

export function updateTeacherCaseIep(id, iepId, data) {
  return request({ url: `${caseUrl}/${id}/ieps/${iepId}`, method: "put", data });
}

export function removeTeacherCaseIeps(id, data) {
  return request({ url: `${caseUrl}/${id}/ieps`, method: "delete", data });
}

export function getTeacherCaseAssessments(id, params) {
  return request({ url: `${caseUrl}/${id}/assessments`, method: "get", params });
}

export function addTeacherCaseAssessment(id, data) {
  return request({ url: `${caseUrl}/${id}/assessments`, method: "post", data });
}

export function updateTeacherCaseAssessment(id, assessmentId, data) {
  return request({ url: `${caseUrl}/${id}/assessments/${assessmentId}`, method: "put", data });
}

export function removeTeacherCaseAssessments(id, data) {
  return request({ url: `${caseUrl}/${id}/assessments`, method: "delete", data });
}

export function getTeacherCaseInterventions(id, params) {
  return request({ url: `${caseUrl}/${id}/interventions`, method: "get", params });
}

export function addTeacherCaseIntervention(id, data) {
  return request({ url: `${caseUrl}/${id}/interventions`, method: "post", data });
}

export function updateTeacherCaseIntervention(id, interventionId, data) {
  return request({ url: `${caseUrl}/${id}/interventions/${interventionId}`, method: "put", data });
}

export function removeTeacherCaseInterventions(id, data) {
  return request({ url: `${caseUrl}/${id}/interventions`, method: "delete", data });
}

export function getTeacherCaseAttachments(id) {
  return request({ url: `${caseUrl}/${id}/attachments`, method: "get" });
}

export function getTeacherCaseAttachmentUrl(id, attachmentId) {
  return request({ url: `${caseUrl}/${id}/attachments/${attachmentId}/file-url`, method: "get" });
}
