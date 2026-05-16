import { request } from "./request";

export function getPublishedActivities(params) {
  return request({ url: "/portal/activities", method: "get", params });
}

export function getPublishedActivity(id) {
  return request({ url: `/portal/activities/${id}`, method: "get" });
}

export function signupActivity(id, data) {
  return request({ url: `/portal/activities/${id}/signup`, method: "post", data });
}

export function getActivitySignupState(id, params) {
  return request({ url: `/portal/activities/${id}/signup-state`, method: "get", params });
}

export function cancelActivitySignup(id, data) {
  return request({ url: `/portal/activities/${id}/signup`, method: "delete", data });
}

export function checkinActivity(id, data) {
  return request({ url: `/portal/activities/${id}/checkin`, method: "post", data });
}

export function uploadActivityOutcomeFile(id, data) {
  return request({
    url: `/portal/activities/${id}/outcomes/files/upload`,
    method: "post",
    data,
    headers: { "Content-Type": "multipart/form-data" }
  });
}

export function submitActivityOutcome(id, data) {
  return request({ url: `/portal/activities/${id}/outcomes`, method: "post", data });
}
