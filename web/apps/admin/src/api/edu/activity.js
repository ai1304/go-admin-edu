import request from '../../utils/request';

const url = '/api/v1/edu/activities';

export function getActivities(params) {
  return request({ url, method: 'get', params });
}

export function getActivity(id) {
  return request({ url: `${url}/${id}`, method: 'get' });
}

export function addActivity(data) {
  return request({ url, method: 'post', data });
}

export function updateActivity(id, data) {
  return request({ url: `${url}/${id}`, method: 'put', data });
}

export function removeActivities(data) {
  return request({ url, method: 'delete', data });
}

