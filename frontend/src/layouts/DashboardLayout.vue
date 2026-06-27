<template>
  <div class="flex h-screen overflow-hidden bg-slate-50">
    <Sidebar />
    <div class="flex flex-1 flex-col overflow-hidden">
      <Topbar />
      <main class="flex-1 overflow-y-auto p-4 lg:p-6">
        <RouterView v-slot="{ Component }">
          <transition enter-active-class="transition duration-200 ease-out" enter-from-class="opacity-0 translate-y-2" mode="out-in">
            <component :is="Component" />
          </transition>
        </RouterView>
      </main>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import Topbar from '@/components/Topbar.vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()

// Refresh the current user on mount (validates the token too).
onMounted(() => {
  if (auth.isAuthenticated) auth.fetchMe().catch(() => {})
})
</script>
