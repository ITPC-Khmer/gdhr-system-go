<template>
  <div class="mx-auto max-w-2xl space-y-5">
    <div>
      <h2 class="text-2xl font-bold text-slate-900">Settings</h2>
      <p class="text-slate-500">Manage your profile and preferences.</p>
    </div>

    <div class="card p-6 sm:p-8">
      <div class="flex items-center gap-4">
        <div class="grid h-16 w-16 place-items-center rounded-2xl bg-gradient-to-br from-primary-400 to-primary-600 text-xl font-bold text-white">
          {{ auth.initials }}
        </div>
        <div>
          <p class="text-lg font-bold text-slate-900">{{ auth.user?.name }}</p>
          <p class="text-sm text-slate-500">{{ auth.user?.email }}</p>
          <span class="mt-1 inline-flex rounded-full bg-blue-100 px-2.5 py-0.5 text-xs font-semibold text-blue-700">{{ auth.user?.role }}</span>
        </div>
      </div>

      <div class="mt-8 grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div class="rounded-xl border border-slate-200 p-4">
          <p class="text-xs uppercase tracking-wide text-slate-400">Primary color</p>
          <div class="mt-2 flex items-center gap-2">
            <span class="h-6 w-6 rounded-lg bg-primary-500 ring-2 ring-primary-200" />
            <span class="font-mono text-sm text-slate-700">DarkOrange · #ff8c00</span>
          </div>
        </div>
        <div class="rounded-xl border border-slate-200 p-4">
          <p class="text-xs uppercase tracking-wide text-slate-400">Account ID</p>
          <p class="mt-2 font-mono text-sm text-slate-700">#{{ auth.user?.id }}</p>
        </div>
      </div>

      <div class="mt-8 border-t border-slate-100 pt-6">
        <button class="btn-danger" @click="logout">
          <Icon name="logout" :size="18" /> Sign out
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import Icon from '@/components/Icon.vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()

function logout() {
  auth.logout()
  router.push({ name: 'login' })
}
</script>
