import { defineStore } from 'pinia'

export const useUiStore = defineStore('ui', {
  state: () => ({
    sidebarCollapsed: false, // desktop collapse
    sidebarMobileOpen: false, // mobile drawer
  }),
  actions: {
    toggleSidebar() {
      this.sidebarCollapsed = !this.sidebarCollapsed
    },
    toggleMobile() {
      this.sidebarMobileOpen = !this.sidebarMobileOpen
    },
    closeMobile() {
      this.sidebarMobileOpen = false
    },
  },
})
