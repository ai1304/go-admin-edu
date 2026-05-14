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

export function addCaseIep(id, data) {
  return request({ url: `${url}/${id}/ieps`, method: 'post', data });
}

