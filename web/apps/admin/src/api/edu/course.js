import request from '../../utils/request';

const url = '/api/v1/edu/courses';

export function getCourses(params) {
  return request({ url, method: 'get', params });
}

export function getCourse(id) {
  return request({ url: `${url}/${id}`, method: 'get' });
}

export function addCourse(data) {
  return request({ url, method: 'post', data });
}

export function updateCourse(id, data) {
  return request({ url: `${url}/${id}`, method: 'put', data });
}

export function removeCourses(data) {
  return request({ url, method: 'delete', data });
}

