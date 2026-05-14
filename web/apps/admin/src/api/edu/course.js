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

export function getCourseChapters(courseId) {
  return request({ url: `${url}/${courseId}/chapters`, method: 'get' });
}

export function addCourseChapter(courseId, data) {
  return request({ url: `${url}/${courseId}/chapters`, method: 'post', data });
}

export function updateCourseChapter(courseId, chapterId, data) {
  return request({ url: `${url}/${courseId}/chapters/${chapterId}`, method: 'put', data });
}

export function removeCourseChapters(courseId, data) {
  return request({ url: `${url}/${courseId}/chapters`, method: 'delete', data });
}

export function getCourseLessons(courseId, params) {
  return request({ url: `${url}/${courseId}/lessons`, method: 'get', params });
}

export function addCourseLesson(courseId, data) {
  return request({ url: `${url}/${courseId}/lessons`, method: 'post', data });
}

export function updateCourseLesson(courseId, lessonId, data) {
  return request({ url: `${url}/${courseId}/lessons/${lessonId}`, method: 'put', data });
}

export function removeCourseLessons(courseId, data) {
  return request({ url: `${url}/${courseId}/lessons`, method: 'delete', data });
}
