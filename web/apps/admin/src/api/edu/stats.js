import request from '../../utils/request';

export function getEduOverview() {
  return request({ url: '/api/v1/edu/stats/overview', method: 'get' });
}

export function getEduResourceStats() {
  return request({ url: '/api/v1/edu/stats/resources', method: 'get' });
}

export function getEduCourseStats() {
  return request({ url: '/api/v1/edu/stats/courses', method: 'get' });
}

export function getEduActivityStats() {
  return request({ url: '/api/v1/edu/stats/activities', method: 'get' });
}

export function getEduSchoolStats() {
  return request({ url: '/api/v1/edu/stats/schools', method: 'get' });
}

export function getEduTeacherStats() {
  return request({ url: '/api/v1/edu/stats/teachers', method: 'get' });
}

export function getEduStudentStats() {
  return request({ url: '/api/v1/edu/stats/students', method: 'get' });
}

export function getEduCaseStats() {
  return request({ url: '/api/v1/edu/stats/cases', method: 'get' });
}

export function exportEduStats() {
  return request({ url: '/api/v1/edu/stats/export', method: 'get', responseType: 'blob' });
}
