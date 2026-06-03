import request from '../../utils/request';

const url = '/api/v1/edu/regions';

export function getRegions(params) {
  return request({ url, method: 'get', params });
}

export function addRegion(data) {
  return request({ url, method: 'post', data });
}

export function updateRegion(id, data) {
  return request({ url: `${url}/${id}`, method: 'put', data });
}

export function removeRegions(data) {
  return request({ url, method: 'delete', data });
}

