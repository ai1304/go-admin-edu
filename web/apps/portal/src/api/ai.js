import { request } from "./request";

export function chatWithAI(data) {
  return request({
    url: "/portal/ai/chat",
    method: "post",
    data
  });
}

export function getAIConversations(params) {
  return request({
    url: "/portal/ai/conversations",
    method: "get",
    params
  });
}

export function getAIConversation(id) {
  return request({
    url: `/portal/ai/conversations/${id}`,
    method: "get"
  });
}

export function deleteAIConversation(id, params) {
  return request({
    url: `/portal/ai/conversations/${id}`,
    method: "delete",
    params
  });
}
