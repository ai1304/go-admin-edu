import { createRouter, createWebHistory } from "vue-router";

const routes = [
  {
    path: "/",
    name: "home",
    component: () => import("@/views/HomeView.vue")
  },
  {
    path: "/resources",
    name: "resources",
    component: () => import("@/views/ResourceListView.vue")
  },
  {
    path: "/resources/:id",
    name: "resource-detail",
    component: () => import("@/views/ResourceDetailView.vue")
  },
  {
    path: "/courses",
    name: "courses",
    component: () => import("@/views/CourseListView.vue")
  },
  {
    path: "/activities",
    name: "activities",
    component: () => import("@/views/ActivityListView.vue")
  },
  {
    path: "/experts",
    name: "experts",
    component: () => import("@/views/ExpertListView.vue")
  },
  {
    path: "/login",
    name: "login",
    component: () => import("@/views/LoginView.vue")
  }
];

export default createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 })
});
