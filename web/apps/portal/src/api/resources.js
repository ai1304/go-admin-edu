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
