import { request } from "./request";

export function getPublishedResources(params) {
  return request({
    url: "/portal/resources",
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

export function getResourceFileAccessUrl(resourceId, fileId) {
  return request({
    url: `/portal/resources/${resourceId}/files/${fileId}/access-url`,
    method: "get"
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
