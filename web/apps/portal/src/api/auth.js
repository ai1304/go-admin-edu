import { request } from "./request";

export function login(data) {
  return request({ url: "/login", method: "post", data });
}

export function getCaptcha() {
  return request({ url: "/captcha", method: "get" });
}

export function getInfo() {
  return request({ url: "/getinfo", method: "get" });
}

export function logout() {
  return request({ url: "/logout", method: "post" });
}
