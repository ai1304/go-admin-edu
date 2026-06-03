import request from '../../utils/request';

const url = '/api/v1/edu/news';

export function getNews(params) {
  return request({ url, method: 'get', params });
}

export function getNewsItem(id) {
  return request({ url: `${url}/${id}`, method: 'get' });
}

export function addNews(data) {
  return request({ url, method: 'post', data });
}

export function updateNews(id, data) {
  return request({ url: `${url}/${id}`, method: 'put', data });
}

export function removeNews(data) {
  return request({ url, method: 'delete', data });
}
