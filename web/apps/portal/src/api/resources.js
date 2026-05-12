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

