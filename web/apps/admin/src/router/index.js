import { createWebHistory, createRouter } from 'vue-router';
import Layout from '../layout/index.vue';
import { useUserStore } from '../store/userInfo';
import { usePermissionStore } from '../store/permission';

const routes = [
  {
    path: '/',
    name: '/',
    redirect: 'admin',
    component: Layout,
    children: [
      {
        path: '/profile',
        name: 'profile',
        component: () => import('../views/profile/index.vue'),
        meta: {
          title: '个人设置',
        },
      },
      {
        path: '/403',
        name: '403',
        component: () => import('../views/error-page/403.vue'),
        meta: {
          title: '无权访问',
        },
      },
      {
        hide: true,
        path: '/:catchAll(.*)',
        component: () => import('../views/error-page/404.vue'),
        meta: {
          title: '找不到页面',
        },
      },
      {
        path: '/500',
        name: '500',
        component: () => import('../views/error-page/500.vue'),
        meta: {
          title: '服务异常',
        },
      },
    ]
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('../views/login/index.vue'),
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach(async (to, from, next) => {
  const store = useUserStore();
  const permissionStore = usePermissionStore();

  if (!store.sysConfig) {
    await store.getSysConfig();
  }

  if (!store.token && to.name !== 'login') {
    next({ name: 'login' });
    return;
  }

  if (store.token && !store.roles) {
    try {
      await store.getUserInfo();
      const dynamicRoutes = await permissionStore.getMenuRole();
      dynamicRoutes.forEach((route) => {
        router.addRoute('/', route);
      });
      next({ ...to, replace: true });
    } catch (err) {
      console.error(err);
      permissionStore.ClearMenuList();
      next({ name: 'login', replace: true });
    }
    return;
  }

  next();
});

router.afterEach((to) => {
  const store = useUserStore();
  const appName = store.sysConfig?.sys_app_name || 'go-admin';
  if (to.name !== 'login') {
    document.title = `${to.meta.title || ''} - ${appName}`;
  } else {
    document.title = appName;
  }
});

export default router;
