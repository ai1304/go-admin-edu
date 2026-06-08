import { defineStore } from 'pinia';
import { getUserMenuRole } from '@/api/admin/login';


const modules = import.meta.glob('../views/**/*.vue');

export const usePermissionStore = defineStore('permisson', {
  state: () => {
    return {
      addRouters: [],
      menuList: [],
    };
  },
  getters: {
    getRoutes: (state) => state.addRouters,
  },
  actions: {
    setMenuList(menus) {
      this.menuList = menus;
    },
    GenerateRoutes(routeList) {
      const routes = [];
      if (!Array.isArray(routeList)) {
        return routes;
      }

      routeList.forEach((item) => {
        if (!item) {
          return;
        }
        const route = {};
        // if (item.visible == 0) {
          if (item.menuType === 'M' || item.menuType === 'C') {
            route.path = item.path;
            route.name = item.menuName;
            if (item.menuType === 'M') {
              route.component = modules[`../views/index.vue`];
            } else if (item.menuType === 'C') {
              route.component = modules[`../views${item.component}.vue`] || modules['../views/error-page/888.vue'];
            }
            route.meta = {
              title: item.title,
              permission: item.permission,
            };
          }

          if (item.children) {
            route.children = this.GenerateRoutes(item.children);
          }
          routes.push(route);
        // }
      });

      return routes;
    },
    async getMenuRole() {
      const res = await getUserMenuRole();
      if (res.code !== 200 || !Array.isArray(res.data)) {
        this.setMenuList([]);
        this.addRouters = [];
        throw new Error(res.msg || 'Failed to load user menus.');
      }
      this.setMenuList(res.data);
      this.addRouters = await this.GenerateRoutes(res.data);
      return this.addRouters;
    },
    ClearMenuList() {
      this.menuList = [];
    }
  },
});
