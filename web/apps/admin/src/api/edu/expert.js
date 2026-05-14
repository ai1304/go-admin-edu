import request from '../../utils/request';

const url = '/api/v1/edu/experts';

export function getExperts(params) {
  return request({ url, method: 'get', params });
}

export function getExpert(id) {
  return request({ url: `${url}/${id}`, method: 'get' });
}

export function addExpert(data) {
  return request({ url, method: 'post', data });
}

export function updateExpert(id, data) {
  return request({ url: `${url}/${id}`, method: 'put', data });
}

export function removeExperts(data) {
  return request({ url, method: 'delete', data });
}

