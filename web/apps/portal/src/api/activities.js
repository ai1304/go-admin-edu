import { request } from "./request";

export function getPublishedActivities(params) {
  return request({ url: "/portal/activities", method: "get", params });
}

export function getPublishedActivity(id) {
  return request({ url: `/portal/activities/${id}`, method: "get" });
}

