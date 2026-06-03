import request from '../../utils/request';

const url = '/api/v1/edu/activities';

export function getActivities(params) {
  return request({ url, method: 'get', params });
}

export function getActivity(id) {
  return request({ url: `${url}/${id}`, method: 'get' });
}

export function addActivity(data) {
  return request({ url, method: 'post', data });
}

export function updateActivity(id, data) {
  return request({ url: `${url}/${id}`, method: 'put', data });
}

export function removeActivities(data) {
  return request({ url, method: 'delete', data });
}

export function getActivitySignups(activityId) {
  return request({ url: `${url}/${activityId}/signups`, method: 'get' });
}

export function addActivitySignup(activityId, data) {
  return request({ url: `${url}/${activityId}/signups`, method: 'post', data });
}

export function updateActivitySignup(activityId, signupId, data) {
  return request({ url: `${url}/${activityId}/signups/${signupId}`, method: 'put', data });
}

export function removeActivitySignups(activityId, data) {
  return request({ url: `${url}/${activityId}/signups`, method: 'delete', data });
}

export function getActivityCheckins(activityId) {
  return request({ url: `${url}/${activityId}/checkins`, method: 'get' });
}

export function addActivityCheckin(activityId, data) {
  return request({ url: `${url}/${activityId}/checkins`, method: 'post', data });
}

export function removeActivityCheckins(activityId, data) {
  return request({ url: `${url}/${activityId}/checkins`, method: 'delete', data });
}

export function getActivityOutcomes(activityId) {
  return request({ url: `${url}/${activityId}/outcomes`, method: 'get' });
}

export function addActivityOutcome(activityId, data) {
  return request({ url: `${url}/${activityId}/outcomes`, method: 'post', data });
}

export function updateActivityOutcome(activityId, outcomeId, data) {
  return request({ url: `${url}/${activityId}/outcomes/${outcomeId}`, method: 'put', data });
}

export function removeActivityOutcomes(activityId, data) {
  return request({ url: `${url}/${activityId}/outcomes`, method: 'delete', data });
}
