import { defineStore } from "pinia";

export const useSessionStore = defineStore("session", {
  state: () => ({
    token: localStorage.getItem("portalToken") || localStorage.getItem("token") || "",
    user: JSON.parse(localStorage.getItem("portalUser") || "null")
  }),
  getters: {
    isLoggedIn: state => Boolean(state.token)
  },
  actions: {
    setSession(payload) {
      this.token = payload?.token || "";
      this.user = payload?.user || null;
      if (this.token) {
        localStorage.setItem("portalToken", this.token);
      }
      if (this.user) {
        localStorage.setItem("portalUser", JSON.stringify(this.user));
      }
    },
    clearSession() {
      this.token = "";
      this.user = null;
      localStorage.removeItem("portalToken");
      localStorage.removeItem("portalUser");
    }
  }
});
