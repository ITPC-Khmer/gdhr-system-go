<template>
  <div class="space-y-5">
    <!-- Header -->
    <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h2 class="text-2xl font-bold text-slate-900">Users</h2>
        <p class="text-slate-500">Manage all user accounts in your system.</p>
      </div>
      <RouterLink to="/users/new" class="btn-primary">
        <Icon name="plus" :size="18" /> Add User
      </RouterLink>
    </div>

    <!-- Toolbar + table -->
    <div class="card overflow-hidden">
      <div class="flex flex-col gap-3 border-b border-slate-100 p-4 sm:flex-row sm:items-center">
        <div class="relative w-full sm:max-w-xs">
          <Icon name="search" :size="18" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
          <input v-model="search" type="text" placeholder="Search by name or email…" class="input pl-10" @input="onSearch" />
        </div>
        <span class="text-sm text-slate-500 sm:ml-auto">{{ total }} user{{ total === 1 ? '' : 's' }} total</span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead class="bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
            <tr>
              <th class="px-5 py-3 font-semibold">User</th>
              <th class="px-5 py-3 font-semibold">Role</th>
              <th class="px-5 py-3 font-semibold">Status</th>
              <th class="px-5 py-3 font-semibold">Joined</th>
              <th class="px-5 py-3 text-right font-semibold">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-if="loading">
              <td colspan="5" class="px-5 py-10 text-center text-slate-400">Loading…</td>
            </tr>
            <tr v-else-if="!users.length">
              <td colspan="5" class="px-5 py-12 text-center text-slate-400">
                <Icon name="users" :size="32" class="mx-auto mb-2 text-slate-300" />
                No users found.
              </td>
            </tr>
            <tr v-for="u in users" :key="u.id" class="hover:bg-slate-50/70 transition">
              <td class="px-5 py-3.5">
                <div class="flex items-center gap-3">
                  <div class="grid h-9 w-9 shrink-0 place-items-center rounded-full bg-gradient-to-br from-primary-400 to-primary-600 text-xs font-bold text-white">
                    {{ getInitials(u.name) }}
                  </div>
                  <div class="min-w-0">
                    <p class="font-semibold text-slate-800">{{ u.name }}</p>
                    <p class="truncate text-xs text-slate-500">{{ u.email }}</p>
                  </div>
                </div>
              </td>
              <td class="px-5 py-3.5">
                <span class="inline-flex rounded-full px-2.5 py-1 text-xs font-semibold"
                  :class="u.role === 'admin' ? 'bg-blue-100 text-blue-700' : 'bg-slate-100 text-slate-600'">
                  {{ u.role }}
                </span>
              </td>
              <td class="px-5 py-3.5">
                <span class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-semibold"
                  :class="u.active ? 'bg-emerald-100 text-emerald-700' : 'bg-rose-100 text-rose-700'">
                  <span class="h-1.5 w-1.5 rounded-full" :class="u.active ? 'bg-emerald-500' : 'bg-rose-500'" />
                  {{ u.active ? 'Active' : 'Inactive' }}
                </span>
              </td>
              <td class="px-5 py-3.5 text-slate-500">{{ formatDate(u.created_at) }}</td>
              <td class="px-5 py-3.5">
                <div class="flex items-center justify-end gap-1">
                  <RouterLink :to="`/users/${u.id}/edit`" class="btn-ghost !p-2" title="Edit">
                    <Icon name="edit" :size="18" />
                  </RouterLink>
                  <button class="btn-ghost !p-2 text-red-500 hover:bg-red-50" title="Delete" @click="confirmDelete(u)">
                    <Icon name="trash" :size="18" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="flex flex-col items-center gap-3 border-t border-slate-100 p-4 sm:flex-row sm:justify-between">
        <p class="text-sm text-slate-500">
          Page <span class="font-semibold text-slate-700">{{ page }}</span> of {{ totalPages }}
        </p>
        <div class="flex gap-2">
          <button class="btn-outline !py-2" :disabled="page <= 1" @click="goTo(page - 1)">Previous</button>
          <button class="btn-outline !py-2" :disabled="page >= totalPages" @click="goTo(page + 1)">Next</button>
        </div>
      </div>
    </div>

    <!-- Delete confirm modal -->
    <transition enter-active-class="transition duration-150" enter-from-class="opacity-0" leave-active-class="transition duration-100" leave-to-class="opacity-0">
      <div v-if="deleteTarget" class="fixed inset-0 z-50 grid place-items-center bg-slate-900/50 p-4" @click.self="deleteTarget = null">
        <div class="w-full max-w-sm rounded-2xl bg-white p-6 shadow-xl">
          <div class="mx-auto mb-4 grid h-12 w-12 place-items-center rounded-full bg-red-100 text-red-600"><Icon name="trash" :size="24" /></div>
          <h3 class="text-center text-lg font-bold text-slate-900">Delete user?</h3>
          <p class="mt-2 text-center text-sm text-slate-500">
            Are you sure you want to delete <span class="font-semibold">{{ deleteTarget.name }}</span>? This cannot be undone.
          </p>
          <div class="mt-6 flex gap-3">
            <button class="btn-outline flex-1" @click="deleteTarget = null">Cancel</button>
            <button class="btn-danger flex-1" :disabled="deleting" @click="doDelete">{{ deleting ? 'Deleting…' : 'Delete' }}</button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import Icon from '@/components/Icon.vue'
import api from '@/api/axios'

const users = ref([])
const total = ref(0)
const page = ref(1)
const limit = ref(10)
const search = ref('')
const loading = ref(false)

const deleteTarget = ref(null)
const deleting = ref(false)

let searchTimer = null

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / limit.value)))

async function load() {
  loading.value = true
  try {
    const { data } = await api.get('/users', {
      params: { page: page.value, limit: limit.value, search: search.value },
    })
    users.value = data.data
    total.value = data.total
  } catch (e) {
    users.value = []
  } finally {
    loading.value = false
  }
}

function onSearch() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    page.value = 1
    load()
  }, 350)
}

function goTo(p) {
  if (p < 1 || p > totalPages.value) return
  page.value = p
  load()
}

function confirmDelete(u) {
  deleteTarget.value = u
}

async function doDelete() {
  deleting.value = true
  try {
    await api.delete(`/users/${deleteTarget.value.id}`)
    deleteTarget.value = null
    if (users.value.length === 1 && page.value > 1) page.value--
    await load()
  } finally {
    deleting.value = false
  }
}

function getInitials(name) {
  return name.split(' ').map((p) => p[0]).slice(0, 2).join('').toUpperCase()
}
function formatDate(d) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })
}

onMounted(load)
</script>
