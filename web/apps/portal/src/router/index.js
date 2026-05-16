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
    path: "/courses/:id",
    name: "course-detail",
    component: () => import("@/views/CourseDetailView.vue")
  },
  {
    path: "/activities",
    name: "activities",
    component: () => import("@/views/ActivityListView.vue")
  },
  {
    path: "/activities/:id",
    name: "activity-detail",
    component: () => import("@/views/ActivityDetailView.vue")
  },
  {
    path: "/experts",
    name: "experts",
    component: () => import("@/views/ExpertListView.vue")
  },
  {
    path: "/experts/:id",
    name: "expert-detail",
    component: () => import("@/views/ExpertDetailView.vue")
  },
  {
    path: "/teacher/workbench",
    name: "teacher-workbench",
    component: () => import("@/views/TeacherWorkbenchView.vue")
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
