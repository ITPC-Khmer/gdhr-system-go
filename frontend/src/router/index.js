import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

import DashboardLayout from '@/layouts/DashboardLayout.vue'
import { resources } from '@/config/resources'

// Build list (+ new/edit unless read-only) routes for every resource from config.
const resourceRoutes = Object.entries(resources).flatMap(([key, cfg]) => {
  const list = {
    path: key,
    name: key,
    component: () => import('@/components/CrudList.vue'),
    props: { resourceKey: key },
    meta: { title: cfg.title },
  }
  if (cfg.readOnly) return [list]
  return [
    list,
    {
      path: `${key}/new`,
      name: `${key}-create`,
      component: () => import('@/components/CrudForm.vue'),
      props: { resourceKey: key },
      meta: { title: `Add ${cfg.singular}`, roles: ['admin'] },
    },
    {
      path: `${key}/:id/edit`,
      name: `${key}-edit`,
      component: () => import('@/components/CrudForm.vue'),
      props: { resourceKey: key },
      meta: { title: `Edit ${cfg.singular}`, roles: ['admin'] },
    },
  ]
})

const routes = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/Login.vue'),
    meta: { guestOnly: true },
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('@/views/Register.vue'),
    meta: { guestOnly: true },
  },
  {
    path: '/',
    component: DashboardLayout,
    meta: { requiresAuth: true },
    children: [
      { path: '', redirect: '/dashboard' },
      {
        path: 'dashboard',
        name: 'dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: 'Dashboard' },
      },
      {
        path: 'users',
        name: 'users',
        component: () => import('@/views/users/UserList.vue'),
        meta: { title: 'All Users', roles: ['admin'] },
      },
      {
        path: 'users/new',
        name: 'user-create',
        component: () => import('@/views/users/UserForm.vue'),
        meta: { title: 'Add User', roles: ['admin'] },
      },
      {
        path: 'users/:id/edit',
        name: 'user-edit',
        component: () => import('@/views/users/UserForm.vue'),
        meta: { title: 'Edit User', roles: ['admin'] },
      },
      ...resourceRoutes,
      {
        path: 'leaves/:id',
        name: 'leave-detail',
        component: () => import('@/views/leaves/LeaveDetail.vue'),
        meta: { title: 'Leave Request' },
      },
      {
        path: 'settings',
        name: 'settings',
        component: () => import('@/views/Settings.vue'),
        meta: { title: 'Settings' },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('@/views/NotFound.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  },
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.meta.guestOnly && auth.isAuthenticated) {
    return { name: 'dashboard' }
  }
  // Role guard: block routes that require a role the user doesn't have.
  if (to.meta.roles && !to.meta.roles.includes(auth.user?.role)) {
    return { name: 'dashboard' }
  }
})

export default router
