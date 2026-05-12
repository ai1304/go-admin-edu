import { defineStore } from "pinia";

export const useSessionStore = defineStore("session", {
  state: () => ({
    token: "",
    user: null
  }),
  getters: {
    isLoggedIn: state => Boolean(state.token)
  },
  actions: {
    setSession(payload) {
      this.token = payload?.token || "";
      this.user = payload?.user || null;
    },
    clearSession() {
      this.token = "";
      this.user = null;
    }
  }
});
