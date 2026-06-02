import request from '../../utils/request';

const url = '/api/v1/admin/recruit';

export function getRecruitStats() {
  return request({ url: `${url}/stats`, method: 'get' });
}

export function getReviewCompanies(params) {
  return request({ url: `${url}/reviews/companies`, method: 'get', params });
}

export function getReviewJobs(params) {
  return request({ url: `${url}/reviews/jobs`, method: 'get', params });
}

export function approveRecruitReview(id, data) {
  return request({ url: `${url}/reviews/${id}/approve`, method: 'post', data });
}

export function rejectRecruitReview(id, data) {
  return request({ url: `${url}/reviews/${id}/reject`, method: 'post', data });
}

export function getRecruitCompanies(params) {
  return request({ url: `${url}/companies`, method: 'get', params });
}

export function getRecruitCompany(id) {
  return request({ url: `${url}/companies/${id}`, method: 'get' });
}

export function updateRecruitCompany(id, data) {
  return request({ url: `${url}/companies/${id}`, method: 'put', data });
}

export function enableRecruitCompany(id) {
  return request({ url: `${url}/companies/${id}/enable`, method: 'post' });
}

export function disableRecruitCompany(id) {
  return request({ url: `${url}/companies/${id}/disable`, method: 'post' });
}

export function getRecruitJobs(params) {
  return request({ url: `${url}/jobs`, method: 'get', params });
}

export function getRecruitJob(id) {
  return request({ url: `${url}/jobs/${id}`, method: 'get' });
}

export function updateRecruitJob(id, data) {
  return request({ url: `${url}/jobs/${id}`, method: 'put', data });
}

export function publishRecruitJob(id) {
  return request({ url: `${url}/jobs/${id}/publish`, method: 'post' });
}

export function offlineRecruitJob(id) {
  return request({ url: `${url}/jobs/${id}/offline`, method: 'post' });
}

export function deleteRecruitJob(id) {
  return request({ url: `${url}/jobs/${id}`, method: 'delete' });
}
