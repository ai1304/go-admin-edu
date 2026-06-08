import { defineStore } from 'pinia';
import { storage } from '@/utils/storage';
import { getInfo } from '@/api/admin/sys-user';
import { getAppConfig } from '@/api/admin/login';

export const useUserStore = defineStore('user', {
  state: () => {
    return {
      token: storage.getItem('token') || null,
      uid: storage.getItem('uid') || null,
      sysConfig: null,
      userInfo: null,
    }
  },
  getters: {
    roles: (state) => state?.userInfo?.roles,
  },
  actions: {
    setToken(token) {
      this.token = token;
      storage.setItem('token', token);
    },
    async getUserInfo() {
      try {
        const { code, data, msg } = await getInfo();
        if (code !== 200 || !data || !data.userId) {
          this.userLogout();
          throw new Error(msg || 'Failed to load current user info.');
        }
        storage.setItem('uid', data.userId);
        this.userInfo = data;
        return data;
      } catch (err) {
        console.error(err);
        this.userLogout();
        throw err;
      }
    },
    async getSysConfig() {
      const sysConfig = storage.getItem('sysConfig');
      if (sysConfig) {
        this.sysConfig = sysConfig;
      } else {
        try {
          const { data, code } = await getAppConfig();
          if (code === 200) {
            storage.setItem('sysConfig', data);
            this.sysConfig = data;
          }
        } catch (err) {
          console.error(err);
        }
      }
    },
    userLogout() {
      this.token = null;
      this.uid = null;
      this.userInfo = null;
      storage.removeItem('token');
      storage.removeItem('uid');
    }
  }
})
