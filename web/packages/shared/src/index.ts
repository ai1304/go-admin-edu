export const API_PREFIX = "/api/v1";

export const APP_SURFACES = {
  admin: "admin",
  portal: "portal"
} as const;

export const USER_TYPES = {
  superAdmin: "super_admin",
  regionAdmin: "region_admin",
  schoolAdmin: "school_admin",
  teacher: "teacher",
  student: "student"
} as const;

export type AppSurface = (typeof APP_SURFACES)[keyof typeof APP_SURFACES];
export type UserType = (typeof USER_TYPES)[keyof typeof USER_TYPES];
