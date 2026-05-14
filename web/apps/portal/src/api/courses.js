import { request } from "./request";

export function getPublishedCourses(params) {
  return request({ url: "/portal/courses", method: "get", params });
}

export function getPublishedCourse(id) {
  return request({ url: `/portal/courses/${id}`, method: "get" });
}

