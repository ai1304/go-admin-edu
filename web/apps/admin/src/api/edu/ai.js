import request from '../../utils/request';

const url = '/api/v1/edu/ai';

export function getAIStats() {
  return request({ url: `${url}/stats`, method: 'get' });
}

export function getAIConversations(params) {
  return request({ url: `${url}/conversations`, method: 'get', params });
}

export function getAIConversation(id) {
  return request({ url: `${url}/conversations/${id}`, method: 'get' });
}

export function removeAIConversation(id) {
  return request({ url: `${url}/conversations/${id}`, method: 'delete' });
}
