import request from '../../utils/request';

const url = '/api/v1/edu/cases';

export function getCases(params) {
  return request({ url, method: 'get', params });
}

export function getCase(id, params) {
  return request({ url: `${url}/${id}`, method: 'get', params });
}

export function addCase(data) {
  return request({ url, method: 'post', data });
}

export function updateCase(id, data) {
  return request({ url: `${url}/${id}`, method: 'put', data });
}

export function removeCases(data) {
  return request({ url, method: 'delete', data });
}

export function submitCaseReview(id) {
  return request({ url: `${url}/${id}/submit-review`, method: 'put' });
}

export function reviewCase(id, data) {
  return request({ url: `${url}/${id}/review`, method: 'put', data });
}

export function getCaseReviews(id) {
  return request({ url: `${url}/${id}/reviews`, method: 'get' });
}

export function getCaseAttachments(id) {
  return request({ url: `${url}/${id}/attachments`, method: 'get' });
}

export function addCaseAttachment(id, data) {
  return request({ url: `${url}/${id}/attachments`, method: 'post', data });
}

export function removeCaseAttachments(id, data) {
  return request({ url: `${url}/${id}/attachments`, method: 'delete', data });
}

export function getCaseAttachmentFileUrl(id, attachmentId) {
  return request({ url: `${url}/${id}/attachments/${attachmentId}/file-url`, method: 'get' });
}

export function getCaseAccessLogs(id, params) {
  return request({ url: `${url}/${id}/access-logs`, method: 'get', params });
}

export function exportCaseAccessLogs(id, params) {
  return request({ url: `${url}/${id}/access-logs/export`, method: 'get', params, responseType: 'blob' });
}

export function getCaseAuthorizations(id, params) {
  return request({ url: `${url}/${id}/authorizations`, method: 'get', params });
}

export function addCaseAuthorization(id, data) {
  return request({ url: `${url}/${id}/authorizations`, method: 'post', data });
}

export function updateCaseAuthorization(id, authorizationId, data) {
  return request({ url: `${url}/${id}/authorizations/${authorizationId}`, method: 'put', data });
}

export function removeCaseAuthorizations(id, data) {
  return request({ url: `${url}/${id}/authorizations`, method: 'delete', data });
}

export function addCaseIep(id, data) {
  return request({ url: `${url}/${id}/ieps`, method: 'post', data });
}

export function getCaseIeps(id, params) {
  return request({ url: `${url}/${id}/ieps`, method: 'get', params });
}

export function updateCaseIep(id, iepId, data) {
  return request({ url: `${url}/${id}/ieps/${iepId}`, method: 'put', data });
}

export function removeCaseIeps(id, data) {
  return request({ url: `${url}/${id}/ieps`, method: 'delete', data });
}

export function getCaseAssessments(id, params) {
  return request({ url: `${url}/${id}/assessments`, method: 'get', params });
}

export function addCaseAssessment(id, data) {
  return request({ url: `${url}/${id}/assessments`, method: 'post', data });
}

export function updateCaseAssessment(id, assessmentId, data) {
  return request({ url: `${url}/${id}/assessments/${assessmentId}`, method: 'put', data });
}

export function removeCaseAssessments(id, data) {
  return request({ url: `${url}/${id}/assessments`, method: 'delete', data });
}

export function getCaseInterventions(id, params) {
  return request({ url: `${url}/${id}/interventions`, method: 'get', params });
}

export function addCaseIntervention(id, data) {
  return request({ url: `${url}/${id}/interventions`, method: 'post', data });
}

export function updateCaseIntervention(id, interventionId, data) {
  return request({ url: `${url}/${id}/interventions/${interventionId}`, method: 'put', data });
}

export function removeCaseInterventions(id, data) {
  return request({ url: `${url}/${id}/interventions`, method: 'delete', data });
}
