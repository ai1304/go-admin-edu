import { request } from "./request";

export function getPublishedCourses(params) {
  return request({ url: "/portal/courses", method: "get", params });
}

export function getPublishedCourse(id) {
  return request({ url: `/portal/courses/${id}`, method: "get" });
}

export function getCourseLearningRecords(id, params) {
  return request({ url: `/portal/courses/${id}/learning-records`, method: "get", params });
}

export function getCourseLessonVideoUrl(id, lessonId) {
  return request({ url: `/portal/courses/${id}/lessons/${lessonId}/video-url`, method: "get" });
}

export function trackCourseLesson(id, lessonId, data) {
  return request({ url: `/portal/courses/${id}/lessons/${lessonId}/learning-records`, method: "post", data });
}

export function submitCourseAssignment(id, assignmentId, data) {
  return request({ url: `/portal/courses/${id}/assignments/${assignmentId}/submissions`, method: "post", data });
}

export function uploadCourseAssignmentFile(id, assignmentId, data) {
  return request({
    url: `/portal/courses/${id}/assignments/${assignmentId}/files/upload`,
    method: "post",
    data,
    headers: { "Content-Type": "multipart/form-data" }
  });
}
