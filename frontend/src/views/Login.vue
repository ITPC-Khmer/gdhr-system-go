<template>
  <div class="grid min-h-screen lg:grid-cols-2">
    <!-- Brand panel -->
    <div class="relative hidden flex-col justify-between overflow-hidden bg-slate-900 p-12 text-white lg:flex">
      <div class="absolute -right-24 -top-24 h-96 w-96 rounded-full bg-primary-500/30 blur-3xl" />
      <div class="absolute -bottom-32 -left-20 h-96 w-96 rounded-full bg-primary-600/20 blur-3xl" />
      <div class="relative z-10 flex items-center gap-3">
        <div class="grid h-10 w-10 place-items-center rounded-xl bg-primary-500 font-extrabold shadow-lg shadow-primary-500/40">A</div>
        <span class="text-xl font-bold">Admin<span class="text-primary-500">System</span></span>
      </div>
      <div class="relative z-10 max-w-md">
        <h2 class="text-4xl font-extrabold leading-tight">Manage everything from one beautiful dashboard.</h2>
        <p class="mt-4 text-slate-400">Go + MySQL backend. Vue 3, Tailwind & Pinia frontend. Fast, secure, and gorgeous.</p>
      </div>
      <p class="relative z-10 text-sm text-slate-500">© {{ new Date().getFullYear() }} AdminSystem. All rights reserved.</p>
    </div>

    <!-- Form panel -->
    <div class="flex items-center justify-center bg-slate-50 p-6">
      <div class="w-full max-w-md">
        <div class="mb-8 text-center lg:text-left">
          <h1 class="text-3xl font-bold text-slate-900">Welcome back 👋</h1>
          <p class="mt-2 text-slate-500">Sign in to your account to continue.</p>
        </div>

        <div v-if="error" class="mb-4 flex items-center gap-2 rounded-lg bg-red-50 px-4 py-3 text-sm text-red-700">
          <Icon name="x" :size="18" /> {{ error }}
        </div>

        <form class="space-y-5" @submit.prevent="onSubmit">
          <div>
            <label class="label">Email address</label>
            <input v-model="email" type="email" required placeholder="admin@example.com" class="input" />
          </div>
          <div>
            <label class="label">Password</label>
            <input v-model="password" type="password" required placeholder="••••••••" class="input" />
          </div>
          <button type="submit" class="btn-primary w-full !py-3" :disabled="loading">
            <span v-if="loading">Signing in…</span>
            <span v-else>Sign in</span>
          </button>
        </form>

        <p class="mt-6 text-center text-sm text-slate-500">
          Don't have an account?
          <RouterLink to="/register" class="font-semibold text-primary-600 hover:text-primary-700">Create one</RouterLink>
        </p>

        <div class="mt-6 rounded-lg border border-dashed border-slate-300 bg-white px-4 py-3 text-center text-xs text-slate-500">
          Demo: <span class="font-semibold text-slate-700">admin@example.com</span> / <span class="font-semibold text-slate-700">admin123</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import Icon from '@/components/Icon.vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const email = ref('admin@example.com')
const password = ref('admin123')
const loading = ref(false)
const error = ref('')

async function onSubmit() {
  loading.value = true
  error.value = ''
  try {
    await auth.login(email.value, password.value)
    router.push(route.query.redirect || '/dashboard')
  } catch (e) {
    error.value = e.response?.data?.message || 'Login failed. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>
