import { request } from "./request";

export function getPublishedNews(params) {
  return request({
    url: "/portal/news",
    method: "get",
    params
  });
}

export function getPublishedNewsItem(id) {
  return request({
    url: `/portal/news/${id}`,
    method: "get"
  });
}

export function likeNews(id) {
  return request({
    url: `/portal/news/${id}/like`,
    method: "post"
  });
}
