import { request } from "./request";

export function getPublishedExperts(params) {
  return request({ url: "/portal/experts", method: "get", params });
}

export function getPublishedExpert(id) {
  return request({ url: `/portal/experts/${id}`, method: "get" });
}

export function getExpertResourceAccessUrl(expertId, resourceId) {
  return request({
    url: `/portal/experts/${expertId}/resources/${resourceId}/access-url`,
    method: "get"
  });
}

export function getExpertFavoriteState(expertId, params) {
  return request({
    url: `/portal/experts/${expertId}/favorite-state`,
    method: "get",
    params
  });
}

export function favoriteExpert(expertId, data) {
  return request({
    url: `/portal/experts/${expertId}/favorite`,
    method: "post",
    data
  });
}

export function unfavoriteExpert(expertId, data) {
  return request({
    url: `/portal/experts/${expertId}/favorite`,
    method: "delete",
    data
  });
}

export function shareExpert(expertId) {
  return request({
    url: `/portal/experts/${expertId}/share`,
    method: "put"
  });
}
