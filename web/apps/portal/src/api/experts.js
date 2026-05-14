import { request } from "./request";

export function getPublishedExperts(params) {
  return request({ url: "/portal/experts", method: "get", params });
}

export function getPublishedExpert(id) {
  return request({ url: `/portal/experts/${id}`, method: "get" });
}

