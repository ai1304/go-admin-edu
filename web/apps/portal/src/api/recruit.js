import { request } from "./request";

export function getRecruitJobs(params) {
  return request({ url: "/recruit/jobs", method: "get", params });
}

export function getRecruitJob(id) {
  return request({ url: `/recruit/jobs/${id}`, method: "get" });
}

export function getRecruitCompanies(params) {
  return request({ url: "/recruit/companies", method: "get", params });
}

export function getRecruitCompany(id) {
  return request({ url: `/recruit/companies/${id}`, method: "get" });
}

export function submitCompanyApplication(data) {
  return request({ url: "/recruit/company-applications", method: "post", data });
}

export function submitJobApplication(data) {
  return request({ url: "/recruit/job-applications", method: "post", data });
}
