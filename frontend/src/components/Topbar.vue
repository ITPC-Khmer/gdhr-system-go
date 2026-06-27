<template>
  <header class="sticky top-0 z-20 flex h-16 items-center gap-3 border-b border-slate-200 bg-white/80 px-4 backdrop-blur-md lg:px-6">
    <!-- Mobile menu -->
    <button class="btn-ghost !p-2 lg:hidden" @click="ui.toggleMobile()">
      <Icon name="menu" :size="22" />
    </button>

    <!-- Desktop collapse -->
    <button class="btn-ghost !p-2 hidden lg:inline-flex" @click="ui.toggleSidebar()" title="Toggle sidebar">
      <Icon name="menu" :size="22" />
    </button>

    <!-- Page title -->
    <div class="min-w-0">
      <h1 class="truncate text-lg font-bold text-slate-800">{{ title }}</h1>
    </div>

    <!-- Search (desktop) -->
    <div class="ml-auto hidden md:block">
      <div class="relative">
        <Icon name="search" :size="18" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
        <input v-model="search" type="text" placeholder="Search…" class="input !py-2 w-56 pl-10" />
      </div>
    </div>

    <!-- Notifications -->
    <button class="btn-ghost relative !p-2 ml-auto md:ml-0">
      <Icon name="bell" :size="22" />
      <span class="absolute right-1.5 top-1.5 h-2 w-2 rounded-full bg-primary-500 ring-2 ring-white" />
    </button>

    <!-- User dropdown -->
    <div class="relative" ref="dropdownRef">
      <button class="flex items-center gap-2 rounded-xl p-1 pr-2 hover:bg-slate-100 transition" @click="open = !open">
        <div class="grid h-9 w-9 place-items-center rounded-full bg-gradient-to-br from-primary-400 to-primary-600 text-sm font-bold text-white">
          {{ auth.initials }}
        </div>
        <div class="hidden text-left sm:block">
          <p class="text-sm font-semibold leading-tight text-slate-800">{{ auth.user?.name }}</p>
          <p class="text-xs leading-tight text-slate-500">{{ auth.user?.email }}</p>
        </div>
        <Icon name="chevron" :size="16" class="hidden text-slate-400 sm:block rotate-90" />
      </button>

      <transition enter-active-class="transition duration-150 ease-out" enter-from-class="opacity-0 scale-95"
        leave-active-class="transition duration-100 ease-in" leave-to-class="opacity-0 scale-95">
        <div v-if="open" class="absolute right-0 mt-2 w-56 origin-top-right rounded-xl border border-slate-200 bg-white p-1.5 shadow-lg">
          <div class="px-3 py-2 border-b border-slate-100 mb-1">
            <p class="text-sm font-semibold text-slate-800">{{ auth.user?.name }}</p>
            <p class="text-xs text-slate-500 truncate">{{ auth.user?.email }}</p>
          </div>
          <RouterLink to="/settings" class="flex items-center gap-2.5 rounded-lg px-3 py-2 text-sm text-slate-700 hover:bg-slate-100" @click="open = false">
            <Icon name="settings" :size="18" /> Settings
          </RouterLink>
          <button class="flex w-full items-center gap-2.5 rounded-lg px-3 py-2 text-sm text-red-600 hover:bg-red-50" @click="logout">
            <Icon name="logout" :size="18" /> Sign out
          </button>
        </div>
      </transition>
    </div>
  </header>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import Icon from '@/components/Icon.vue'
import { useUiStore } from '@/stores/ui'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const ui = useUiStore()
const auth = useAuthStore()

const search = ref('')
const open = ref(false)
const dropdownRef = ref(null)

const title = computed(() => route.meta.title || 'Dashboard')

function logout() {
  auth.logout()
  router.push({ name: 'login' })
}

function onClickOutside(e) {
  if (dropdownRef.value && !dropdownRef.value.contains(e.target)) open.value = false
}
onMounted(() => document.addEventListener('click', onClickOutside))
onUnmounted(() => document.removeEventListener('click', onClickOutside))
</script>
