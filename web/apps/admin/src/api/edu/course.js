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

export function getCourseAssignments(courseId) {
  return request({ url: `${url}/${courseId}/assignments`, method: 'get' });
}

export function addCourseAssignment(courseId, data) {
  return request({ url: `${url}/${courseId}/assignments`, method: 'post', data });
}

export function updateCourseAssignment(courseId, assignmentId, data) {
  return request({ url: `${url}/${courseId}/assignments/${assignmentId}`, method: 'put', data });
}

export function removeCourseAssignments(courseId, data) {
  return request({ url: `${url}/${courseId}/assignments`, method: 'delete', data });
}

export function getCourseAssignmentSubmissions(courseId, assignmentId) {
  return request({ url: `${url}/${courseId}/assignments/${assignmentId}/submissions`, method: 'get' });
}

export function addCourseAssignmentSubmission(courseId, assignmentId, data) {
  return request({ url: `${url}/${courseId}/assignments/${assignmentId}/submissions`, method: 'post', data });
}

export function updateCourseAssignmentSubmission(courseId, assignmentId, submissionId, data) {
  return request({ url: `${url}/${courseId}/assignments/${assignmentId}/submissions/${submissionId}`, method: 'put', data });
}

export function getCourseAssignmentSubmissionFileUrl(courseId, assignmentId, submissionId) {
  return request({ url: `${url}/${courseId}/assignments/${assignmentId}/submissions/${submissionId}/file-url`, method: 'get' });
}

export function removeCourseAssignmentSubmissions(courseId, assignmentId, data) {
  return request({ url: `${url}/${courseId}/assignments/${assignmentId}/submissions`, method: 'delete', data });
}

export function getCourseLearningRecords(courseId, params) {
  return request({ url: `${url}/${courseId}/learning-records`, method: 'get', params });
}

export function addCourseLearningRecord(courseId, data) {
  return request({ url: `${url}/${courseId}/learning-records`, method: 'post', data });
}

export function updateCourseLearningRecord(courseId, recordId, data) {
  return request({ url: `${url}/${courseId}/learning-records/${recordId}`, method: 'put', data });
}

export function removeCourseLearningRecords(courseId, data) {
  return request({ url: `${url}/${courseId}/learning-records`, method: 'delete', data });
}
