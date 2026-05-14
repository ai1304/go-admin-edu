import request from '../../utils/request';

export function getEduOverview() {
  return request({ url: '/api/v1/edu/stats/overview', method: 'get' });
}

