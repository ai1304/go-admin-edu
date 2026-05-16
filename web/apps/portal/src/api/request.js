import axios from "axios";
import { API_PREFIX } from "@special-edu/shared";

export const request = axios.create({
  baseURL: API_PREFIX,
  timeout: 15000
});

request.interceptors.request.use(config => {
  const token = localStorage.getItem("portalToken") || localStorage.getItem("token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

request.interceptors.response.use(
  response => response.data,
  error => Promise.reject(error)
);
