import request from '../../utils/request';

const url = '/api/v1/edu/cases';

export function getCases(params) {
  return request({ url, method: 'get', params });
}

export function getCase(id) {
  return request({ url: `${url}/${id}`, method: 'get' });
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

export function getCaseAccessLogs(id) {
  return request({ url: `${url}/${id}/access-logs`, method: 'get' });
}

export function addCaseIep(id, data) {
  return request({ url: `${url}/${id}/ieps`, method: 'post', data });
}

export function getCaseIeps(id) {
  return request({ url: `${url}/${id}/ieps`, method: 'get' });
}

export function updateCaseIep(id, iepId, data) {
  return request({ url: `${url}/${id}/ieps/${iepId}`, method: 'put', data });
}

export function removeCaseIeps(id, data) {
  return request({ url: `${url}/${id}/ieps`, method: 'delete', data });
}

export function getCaseAssessments(id) {
  return request({ url: `${url}/${id}/assessments`, method: 'get' });
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

export function getCaseInterventions(id) {
  return request({ url: `${url}/${id}/interventions`, method: 'get' });
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
