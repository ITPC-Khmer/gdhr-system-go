<template>
  <div class="space-y-6">
    <!-- Greeting -->
    <div class="flex flex-col gap-1">
      <h2 class="text-2xl font-bold text-slate-900">Welcome back, {{ auth.user?.name?.split(' ')[0] }} 👋</h2>
      <p class="text-slate-500">Here's what's happening with your system today.</p>
    </div>

    <!-- Stat cards -->
    <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 xl:grid-cols-4">
      <div v-for="card in cards" :key="card.label" class="card p-5">
        <div class="flex items-center justify-between">
          <div class="grid h-11 w-11 place-items-center rounded-xl" :class="card.bg">
            <Icon :name="card.icon" :size="22" :class="card.fg" />
          </div>
          <span class="text-xs font-semibold" :class="card.trendColor">{{ card.trend }}</span>
        </div>
        <p class="mt-4 text-3xl font-extrabold text-slate-900">{{ card.value }}</p>
        <p class="mt-1 text-sm text-slate-500">{{ card.label }}</p>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
      <!-- Activity / chart placeholder -->
      <div class="card p-6 lg:col-span-2">
        <div class="mb-5 flex items-center justify-between">
          <h3 class="font-bold text-slate-800">Activity overview</h3>
          <span class="text-xs text-slate-400">Last 7 days</span>
        </div>
        <div class="flex h-56 items-end gap-3">
          <div v-for="(bar, i) in bars" :key="i" class="flex flex-1 flex-col items-center gap-2">
            <div class="w-full rounded-t-lg bg-gradient-to-t from-primary-500 to-primary-400 transition-all hover:from-primary-600"
              :style="{ height: bar + '%' }" />
            <span class="text-xs text-slate-400">{{ days[i] }}</span>
          </div>
        </div>
      </div>

      <!-- Quick actions -->
      <div class="card p-6">
        <h3 class="mb-5 font-bold text-slate-800">Quick actions</h3>
        <div class="space-y-3">
          <RouterLink v-if="auth.isAdmin" to="/users/new" class="flex items-center gap-3 rounded-xl border border-slate-200 p-3 hover:border-primary-300 hover:bg-primary-50/50 transition">
            <div class="grid h-9 w-9 place-items-center rounded-lg bg-primary-100 text-primary-600"><Icon name="plus" :size="18" /></div>
            <div><p class="text-sm font-semibold text-slate-800">Add a user</p><p class="text-xs text-slate-500">Create a new account</p></div>
          </RouterLink>
          <RouterLink v-if="auth.isAdmin" to="/users" class="flex items-center gap-3 rounded-xl border border-slate-200 p-3 hover:border-primary-300 hover:bg-primary-50/50 transition">
            <div class="grid h-9 w-9 place-items-center rounded-lg bg-blue-100 text-blue-600"><Icon name="list" :size="18" /></div>
            <div><p class="text-sm font-semibold text-slate-800">Manage users</p><p class="text-xs text-slate-500">View & edit all users</p></div>
          </RouterLink>
          <RouterLink to="/settings" class="flex items-center gap-3 rounded-xl border border-slate-200 p-3 hover:border-primary-300 hover:bg-primary-50/50 transition">
            <div class="grid h-9 w-9 place-items-center rounded-lg bg-emerald-100 text-emerald-600"><Icon name="settings" :size="18" /></div>
            <div><p class="text-sm font-semibold text-slate-800">Settings</p><p class="text-xs text-slate-500">Configure your system</p></div>
          </RouterLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import Icon from '@/components/Icon.vue'
import api from '@/api/axios'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const stats = ref({ total_users: 0, active_users: 0, admins: 0, inactive: 0 })

const days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
const bars = [55, 75, 45, 90, 65, 80, 40]

const cards = computed(() => [
  { label: 'Total Users', value: stats.value.total_users, icon: 'users', bg: 'bg-primary-100', fg: 'text-primary-600', trend: '+12%', trendColor: 'text-emerald-600' },
  { label: 'Active Users', value: stats.value.active_users, icon: 'activity', bg: 'bg-emerald-100', fg: 'text-emerald-600', trend: '+8%', trendColor: 'text-emerald-600' },
  { label: 'Administrators', value: stats.value.admins, icon: 'shield', bg: 'bg-blue-100', fg: 'text-blue-600', trend: 'stable', trendColor: 'text-slate-400' },
  { label: 'Inactive', value: stats.value.inactive, icon: 'user', bg: 'bg-rose-100', fg: 'text-rose-600', trend: '-3%', trendColor: 'text-rose-600' },
])

onMounted(async () => {
  try {
    const { data } = await api.get('/stats')
    stats.value = data
  } catch (e) {
    /* keep zeros on failure */
  }
})
</script>
