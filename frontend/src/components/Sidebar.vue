<template>
  <!-- Mobile overlay -->
  <transition enter-active-class="transition-opacity duration-200" enter-from-class="opacity-0"
    leave-active-class="transition-opacity duration-200" leave-to-class="opacity-0">
    <div v-if="ui.sidebarMobileOpen" class="fixed inset-0 z-30 bg-slate-900/50 lg:hidden" @click="ui.closeMobile()" />
  </transition>

  <aside
    :class="[
      'fixed lg:static inset-y-0 left-0 z-40 flex flex-col bg-slate-900 text-slate-300 transition-all duration-300 ease-in-out',
      collapsed ? 'w-20' : 'w-64',
      ui.sidebarMobileOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0',
    ]"
  >
    <!-- Brand -->
    <div class="flex h-16 items-center gap-3 px-5 border-b border-white/5">
      <div class="grid h-9 w-9 shrink-0 place-items-center rounded-xl bg-primary-500 text-white font-extrabold shadow-lg shadow-primary-500/30">
        A
      </div>
      <span v-if="!collapsed" class="text-lg font-bold text-white tracking-tight whitespace-nowrap">Admin<span class="text-primary-500">System</span></span>
    </div>

    <!-- Nav -->
    <nav class="flex-1 overflow-y-auto px-3 py-4 space-y-1">
      <template v-for="item in menu" :key="item.label">
        <!-- Simple link -->
        <RouterLink
          v-if="!item.children"
          :to="item.to"
          class="group flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-colors"
          :class="isActive(item.to)
            ? 'bg-primary-500 text-white shadow-md shadow-primary-500/20'
            : 'text-slate-300 hover:bg-white/5 hover:text-white'"
          @click="ui.closeMobile()"
        >
          <Icon :name="item.icon" :size="20" class="shrink-0" />
          <span v-if="!collapsed" class="whitespace-nowrap">{{ item.label }}</span>
        </RouterLink>

        <!-- Group with submenu -->
        <div v-else>
          <button
            class="group flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-colors"
            :class="isGroupActive(item)
              ? 'text-white bg-white/5'
              : 'text-slate-300 hover:bg-white/5 hover:text-white'"
            @click="toggle(item.label)"
          >
            <Icon :name="item.icon" :size="20" class="shrink-0" />
            <span v-if="!collapsed" class="flex-1 text-left whitespace-nowrap">{{ item.label }}</span>
            <Icon
              v-if="!collapsed"
              name="chevron"
              :size="16"
              class="shrink-0 transition-transform duration-200"
              :class="isOpen(item.label) ? 'rotate-90' : ''"
            />
          </button>

          <!-- Submenu -->
          <transition
            enter-active-class="transition-all duration-200 ease-out overflow-hidden"
            enter-from-class="max-h-0 opacity-0"
            enter-to-class="max-h-60 opacity-100"
            leave-active-class="transition-all duration-200 ease-in overflow-hidden"
            leave-from-class="max-h-60 opacity-100"
            leave-to-class="max-h-0 opacity-0"
          >
            <div v-if="isOpen(item.label) && !collapsed" class="mt-1 space-y-1 pl-4">
              <RouterLink
                v-for="child in item.children"
                :key="child.to"
                :to="child.to"
                class="flex items-center gap-3 rounded-lg py-2 pl-4 pr-3 text-sm transition-colors border-l border-white/10"
                :class="isActive(child.to, child.exact)
                  ? 'text-primary-400 font-semibold border-primary-500'
                  : 'text-slate-400 hover:text-white'"
                @click="ui.closeMobile()"
              >
                <span class="h-1.5 w-1.5 rounded-full" :class="isActive(child.to, child.exact) ? 'bg-primary-500' : 'bg-slate-600'" />
                {{ child.label }}
              </RouterLink>
            </div>
          </transition>
        </div>
      </template>
    </nav>

    <!-- Footer -->
    <div class="border-t border-white/5 p-3">
      <div class="flex items-center gap-3 rounded-xl px-2 py-2" :class="collapsed ? 'justify-center' : ''">
        <div class="grid h-9 w-9 shrink-0 place-items-center rounded-full bg-primary-500/20 text-primary-400 text-sm font-bold">
          {{ auth.initials }}
        </div>
        <div v-if="!collapsed" class="min-w-0">
          <p class="truncate text-sm font-semibold text-white">{{ auth.user?.name }}</p>
          <p class="truncate text-xs text-slate-400">{{ auth.user?.role }}</p>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import Icon from '@/components/Icon.vue'
import { useUiStore } from '@/stores/ui'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const ui = useUiStore()
const auth = useAuthStore()

const collapsed = computed(() => ui.sidebarCollapsed)

// Full menu definition. `adminOnly` items (and children) are hidden for non-admins.
const rawMenu = [
  { label: 'Dashboard', icon: 'dashboard', to: '/dashboard' },
  {
    label: 'User Management',
    icon: 'users',
    adminOnly: true,
    children: [
      { label: 'All Users', to: '/users', exact: true },
      { label: 'Add User', to: '/users/new' },
    ],
  },
  {
    label: 'GDHR Data',
    icon: 'list',
    children: [
      { label: 'Institutes', to: '/institutes', exact: true },
      { label: 'Staffs', to: '/staffs', exact: true },
      { label: 'Ranks', to: '/ranks', exact: true },
      { label: 'Positions', to: '/positions', exact: true },
    ],
  },
  { label: 'Holidays', icon: 'calendar', to: '/holidays' },
  {
    label: 'Leave Management',
    icon: 'calendar',
    children: [
      { label: 'Leave Requests', to: '/leaves', exact: true },
      { label: 'Leave Approvals', to: '/leave_approvals', exact: true },
      { label: 'Leave Files', to: '/leave_files', exact: true },
      { label: 'Leave Balances', to: '/leave_years', exact: true },
      { label: 'Leave Types', to: '/leave_types', exact: true },
      { label: 'Leave Roles', to: '/leave_roles', exact: true },
      { label: 'Manage Roles', to: '/staff_institute_roles', exact: true },
    ],
  },
  { label: 'Settings', icon: 'settings', to: '/settings' },
]

// Role-filtered menu: drop admin-only entries and any group left with no children.
const menu = computed(() => {
  const isAdmin = auth.isAdmin
  return rawMenu
    .filter((item) => !item.adminOnly || isAdmin)
    .map((item) => {
      if (!item.children) return item
      return { ...item, children: item.children.filter((c) => !c.adminOnly || isAdmin) }
    })
    .filter((item) => !item.children || item.children.length > 0)
})

// Track which groups are expanded; open the active group by default.
const openGroups = ref(
  menu.value.filter((m) => m.children && isGroupActive(m)).map((m) => m.label)
)

function toggle(label) {
  const i = openGroups.value.indexOf(label)
  if (i === -1) openGroups.value.push(label)
  else openGroups.value.splice(i, 1)
}
function isOpen(label) {
  return openGroups.value.includes(label)
}
function isActive(to, exact = false) {
  return exact ? route.path === to : route.path === to || route.path.startsWith(to + '/')
}
function isGroupActive(item) {
  return item.children?.some((c) => isActive(c.to, c.exact))
}
</script>
