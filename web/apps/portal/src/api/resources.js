import { request } from "./request";

export function getPublishedResources(params) {
  return request({
    url: "/portal/resources",
    method: "get",
    params
  });
}

export function searchPublishedResources(params) {
  return request({
    url: "/portal/search/resources",
    method: "get",
    params
  });
}

export function getResourceCategories(params) {
  return request({
    url: "/portal/resource-categories",
    method: "get",
    params
  });
}

export function getResourceTags(params) {
  return request({
    url: "/portal/resource-tags",
    method: "get",
    params
  });
}

export function getPublishedResource(id) {
  return request({
    url: `/portal/resources/${id}`,
    method: "get"
  });
}

export function getResourceFileAccessUrl(resourceId, fileId, params) {
  return request({
    url: `/portal/resources/${resourceId}/files/${fileId}/access-url`,
    method: "get",
    params
  });
}

export function getResourceFavoriteState(resourceId, params) {
  return request({
    url: `/portal/resources/${resourceId}/favorite-state`,
    method: "get",
    params
  });
}

export function favoriteResource(resourceId, data) {
  return request({
    url: `/portal/resources/${resourceId}/favorite`,
    method: "post",
    data
  });
}

export function unfavoriteResource(resourceId, data) {
  return request({
    url: `/portal/resources/${resourceId}/favorite`,
    method: "delete",
    data
  });
}

export function getResourceComments(resourceId) {
  return request({
    url: `/portal/resources/${resourceId}/comments`,
    method: "get"
  });
}

export function createResourceComment(resourceId, data) {
  return request({
    url: `/portal/resources/${resourceId}/comments`,
    method: "post",
    data
  });
}

export function likeResourceComment(resourceId, commentId) {
  return request({
    url: `/portal/resources/${resourceId}/comments/${commentId}/like`,
    method: "put"
  });
}
