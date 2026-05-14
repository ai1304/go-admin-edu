import request from '../../utils/request';

const resourceUrl = '/api/v1/edu/resources';
const categoryUrl = '/api/v1/edu/resource-categories';
const tagUrl = '/api/v1/edu/resource-tags';
const fileUrl = '/api/v1/edu/resource-files';

export function getResources(params) {
  return request({ url: resourceUrl, method: 'get', params });
}

export function getResource(id) {
  return request({ url: `${resourceUrl}/${id}`, method: 'get' });
}

export function addResource(data) {
  return request({ url: resourceUrl, method: 'post', data });
}

export function updateResource(id, data) {
  return request({ url: `${resourceUrl}/${id}`, method: 'put', data });
}

export function removeResources(data) {
  return request({ url: resourceUrl, method: 'delete', data });
}

export function submitResourceReview(id) {
  return request({ url: `${resourceUrl}/${id}/submit-review`, method: 'put' });
}

export function reviewResource(id, data) {
  return request({ url: `${resourceUrl}/${id}/review`, method: 'put', data });
}

export function getResourceCategories(params) {
  return request({ url: categoryUrl, method: 'get', params });
}

export function addResourceCategory(data) {
  return request({ url: categoryUrl, method: 'post', data });
}

export function updateResourceCategory(id, data) {
  return request({ url: `${categoryUrl}/${id}`, method: 'put', data });
}

export function removeResourceCategories(data) {
  return request({ url: categoryUrl, method: 'delete', data });
}

export function getResourceTags(params) {
  return request({ url: tagUrl, method: 'get', params });
}

export function addResourceTag(data) {
  return request({ url: tagUrl, method: 'post', data });
}

export function updateResourceTag(id, data) {
  return request({ url: `${tagUrl}/${id}`, method: 'put', data });
}

export function removeResourceTags(data) {
  return request({ url: tagUrl, method: 'delete', data });
}

export function getResourceFiles(params) {
  return request({ url: fileUrl, method: 'get', params });
}

export function addResourceFile(data) {
  return request({ url: fileUrl, method: 'post', data });
}

export function uploadResourceFile(data) {
  return request({ url: `${fileUrl}/upload`, method: 'post', data });
}

export function removeResourceFiles(data) {
  return request({ url: fileUrl, method: 'delete', data });
}
