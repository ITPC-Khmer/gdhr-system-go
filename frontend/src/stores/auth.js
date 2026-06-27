import { defineStore } from 'pinia'
import api from '@/api/axios'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('user') || 'null'),
    token: localStorage.getItem('token') || null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    isAdmin: (state) => state.user?.role === 'admin',
    initials: (state) => {
      if (!state.user?.name) return '?'
      return state.user.name
        .split(' ')
        .map((p) => p[0])
        .slice(0, 2)
        .join('')
        .toUpperCase()
    },
  },

  actions: {
    persist(token, user) {
      this.token = token
      this.user = user
      localStorage.setItem('token', token)
      localStorage.setItem('user', JSON.stringify(user))
    },

    async login(email, password) {
      const { data } = await api.post('/auth/login', { email, password })
      this.persist(data.token, data.user)
      return data
    },

    async register(name, email, password) {
      const { data } = await api.post('/auth/register', { name, email, password })
      this.persist(data.token, data.user)
      return data
    },

    async fetchMe() {
      const { data } = await api.get('/auth/me')
      this.user = data.user
      localStorage.setItem('user', JSON.stringify(data.user))
      return data.user
    },

    logout() {
      this.token = null
      this.user = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    },
  },
})
