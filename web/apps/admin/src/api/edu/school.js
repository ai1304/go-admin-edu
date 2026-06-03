import request from '../../utils/request';

const url = '/api/v1/edu/schools';

export function getSchools(params) {
  return request({ url, method: 'get', params });
}

export function addSchool(data) {
  return request({ url, method: 'post', data });
}

export function updateSchool(id, data) {
  return request({ url: `${url}/${id}`, method: 'put', data });
}

export function removeSchools(data) {
  return request({ url, method: 'delete', data });
}

