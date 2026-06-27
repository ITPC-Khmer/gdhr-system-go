<template>
  <div class="space-y-5">
    <!-- Header -->
    <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h2 class="text-2xl font-bold text-slate-900">{{ cfg.title }}</h2>
        <p class="text-slate-500">Browse, search, filter and manage {{ cfg.title.toLowerCase() }}.</p>
      </div>
      <RouterLink v-if="isAdmin && !cfg.readOnly" :to="`/${resourceKey}/new`" class="btn-primary">
        <Icon name="plus" :size="18" /> Add {{ cfg.singular }}
      </RouterLink>
    </div>

    <!-- Summary count cards -->
    <div v-if="cfg.summary" class="grid grid-cols-2 gap-3 sm:grid-cols-4">
      <div v-for="c in cfg.summary.cards" :key="c.key" class="card p-4">
        <p class="text-xs font-medium uppercase tracking-wide text-slate-400">{{ c.label }}</p>
        <p class="mt-1 text-2xl font-bold" :class="c.class">{{ summary[c.key] ?? 0 }}</p>
      </div>
    </div>

    <div class="card overflow-hidden">
      <!-- Toolbar: search + filters -->
      <div class="space-y-3 border-b border-slate-100 p-4">
        <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
          <div class="relative w-full sm:max-w-sm">
            <Icon name="search" :size="18" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
            <input v-model="search" type="text" :placeholder="cfg.searchPlaceholder" class="input pl-10" @input="onSearch" />
          </div>
          <button v-if="cfg.filters?.length" class="btn-outline !py-2" @click="showFilters = !showFilters">
            <Icon name="list" :size="16" /> Filters
            <span v-if="activeFilterCount" class="ml-1 rounded-full bg-primary-500 px-1.5 text-xs text-white">{{ activeFilterCount }}</span>
          </button>
          <span class="text-sm text-slate-500 sm:ml-auto">{{ total }} record{{ total === 1 ? '' : 's' }}</span>
        </div>

        <!-- Filter row -->
        <div v-if="showFilters && cfg.filters?.length" class="grid grid-cols-2 gap-3 rounded-xl bg-slate-50 p-3 sm:grid-cols-3 lg:grid-cols-4">
          <div v-for="f in cfg.filters" :key="f.key">
            <label class="mb-1 block text-xs font-medium text-slate-500">{{ f.label }}</label>
            <select v-if="f.type === 'select'" v-model="filters[f.key]" class="input !py-2 text-sm" @change="applyFilters">
              <option v-for="o in f.options" :key="o.value" :value="o.value">{{ o.label }}</option>
            </select>
            <input v-else v-model="filters[f.key]" :type="f.type === 'number' ? 'number' : 'text'"
              class="input !py-2 text-sm" :placeholder="f.label" @input="onFilterInput" />
          </div>
          <div class="col-span-2 flex items-end sm:col-span-1">
            <button class="btn-ghost !py-2 text-sm" @click="clearFilters">Clear filters</button>
          </div>
        </div>
      </div>

      <!-- Table -->
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead class="bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
            <tr>
              <th v-for="col in cfg.columns" :key="col.key" class="whitespace-nowrap px-5 py-3 font-semibold">{{ col.label }}</th>
              <th class="px-5 py-3 text-right font-semibold">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-if="loading">
              <td :colspan="cfg.columns.length + 1" class="px-5 py-10 text-center text-slate-400">Loading…</td>
            </tr>
            <tr v-else-if="!rows.length">
              <td :colspan="cfg.columns.length + 1" class="px-5 py-12 text-center text-slate-400">
                <Icon :name="cfg.icon" :size="32" class="mx-auto mb-2 text-slate-300" />
                No {{ cfg.title.toLowerCase() }} found.
              </td>
            </tr>
            <tr v-for="row in rows" :key="row[cfg.pk]" class="hover:bg-slate-50/70 transition">
              <td v-for="col in cfg.columns" :key="col.key" class="px-5 py-3.5"
                :class="[col.mono ? 'font-mono text-xs text-slate-500' : 'text-slate-700', col.truncate ? 'max-w-[12rem] truncate' : '']">
                <!-- combined identity cell: avatar + names + badge + tags -->
                <template v-if="col.type === 'profile'">
                  <div class="flex items-center gap-3">
                    <img v-if="col.image && imgOk(row[col.image])" :src="proxied(row[col.image])" loading="lazy" referrerpolicy="no-referrer"
                      class="h-11 w-11 shrink-0 rounded-full object-cover ring-1 ring-slate-200" @error="brokenImg.add(row[col.image])" />
                    <div v-else-if="col.image || col.avatarIcon"
                      class="grid h-11 w-11 shrink-0 place-items-center rounded-full"
                      :class="col.avatarIcon ? 'bg-primary-100 text-primary-600' : 'bg-slate-100 text-slate-400'">
                      <Icon :name="col.avatarIcon || 'user'" :size="20" />
                    </div>
                    <div class="min-w-0">
                      <div class="flex items-center gap-2">
                        <span class="truncate font-semibold text-slate-800">{{ display(fieldText(row, col.title)) }}</span>
                        <span v-if="col.badge && row[col.badge]" class="shrink-0 rounded bg-slate-100 px-1.5 py-0.5 font-mono text-[11px] text-slate-500">{{ row[col.badge] }}</span>
                      </div>
                      <p v-if="col.subtitle && fieldText(row, col.subtitle)" class="truncate text-xs text-slate-500">{{ fieldText(row, col.subtitle) }}</p>
                      <p v-if="col.caption && captionText(row, col.caption)" class="flex items-center gap-1 text-[11px] text-slate-400">
                        <Icon v-if="captionIcon(col.caption)" :name="captionIcon(col.caption)" :size="12" class="shrink-0" />
                        <span class="truncate">{{ captionPrefix(col.caption) }}{{ captionText(row, col.caption) }}</span>
                      </p>
                      <div v-if="col.tags?.length" class="mt-1 flex flex-wrap gap-1">
                        <template v-for="t in col.tags" :key="tagKey(t)">
                          <span v-if="tagValue(row, t)" class="inline-flex items-center rounded-full px-2 py-0.5 text-[11px] font-medium" :class="tagClass(t)">
                            {{ tagValue(row, t) }}
                          </span>
                        </template>
                      </div>
                    </div>
                  </div>
                </template>
                <span v-else-if="col.type === 'chain'" class="block max-w-[26rem] line-clamp-2 text-slate-600" :title="chainText(row, col)">
                  {{ chainText(row, col) || '—' }}
                </span>
                <template v-else-if="col.type === 'image'">
                  <img v-if="imgOk(row[col.key])" :src="proxied(row[col.key])" loading="lazy" referrerpolicy="no-referrer"
                    class="h-10 w-10 rounded-full object-cover ring-1 ring-slate-200" @error="brokenImg.add(row[col.key])" />
                  <div v-else class="grid h-10 w-10 place-items-center rounded-full bg-slate-100 text-slate-400">
                    <Icon name="user" :size="18" />
                  </div>
                </template>
                <span v-else-if="col.type === 'bool'"
                  class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-semibold"
                  :class="row[col.key] ? 'bg-emerald-100 text-emerald-700' : 'bg-slate-100 text-slate-500'">
                  <span class="h-1.5 w-1.5 rounded-full" :class="row[col.key] ? 'bg-emerald-500' : 'bg-slate-400'" />
                  {{ row[col.key] ? 'Yes' : 'No' }}
                </span>
                <span v-else :title="String(row[col.key] ?? '')">{{ display(row[col.key]) }}</span>
              </td>
              <td class="px-5 py-3.5">
                <div class="flex items-center justify-end gap-1">
                  <RouterLink v-if="cfg.detail" :to="`/${resourceKey}/${encodeURIComponent(row[cfg.pk])}`" class="btn-ghost !p-2" title="View">
                    <Icon name="activity" :size="18" />
                  </RouterLink>
                  <button v-for="act in rowActions(row)" :key="act.key"
                    class="btn-ghost !p-2" :class="actionClass(act)" :title="act.label"
                    @click="openAction(act, row)">
                    <Icon :name="act.icon" :size="18" />
                  </button>
                  <RouterLink v-if="!cfg.readOnly" :to="`/${resourceKey}/${encodeURIComponent(row[cfg.pk])}/edit`" class="btn-ghost !p-2" title="Edit">
                    <Icon name="edit" :size="18" />
                  </RouterLink>
                  <button v-if="isAdmin && !cfg.readOnly" class="btn-ghost !p-2 text-red-500 hover:bg-red-50" title="Delete" @click="deleteTarget = row">
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
        <p class="text-sm text-slate-500">Page <span class="font-semibold text-slate-700">{{ page }}</span> of {{ totalPages }}</p>
        <div class="flex items-center gap-2">
          <select v-model.number="limit" class="input !w-auto !py-2 text-sm" @change="reset">
            <option :value="10">10</option>
            <option :value="20">20</option>
            <option :value="50">50</option>
            <option :value="100">100</option>
          </select>
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
          <h3 class="text-center text-lg font-bold text-slate-900">Delete {{ cfg.singular.toLowerCase() }}?</h3>
          <p class="mt-2 text-center text-sm text-slate-500">
            Delete <span class="font-semibold">{{ deleteLabel }}</span>? This cannot be undone.
          </p>
          <div class="mt-6 flex gap-3">
            <button class="btn-outline flex-1" @click="deleteTarget = null">Cancel</button>
            <button class="btn-danger flex-1" :disabled="deleting" @click="doDelete">{{ deleting ? 'Deleting…' : 'Delete' }}</button>
          </div>
        </div>
      </div>
    </transition>

    <!-- Row action confirm modal (approve / reject / …) -->
    <transition enter-active-class="transition duration-150" enter-from-class="opacity-0" leave-active-class="transition duration-100" leave-to-class="opacity-0">
      <div v-if="actionTarget" class="fixed inset-0 z-50 grid place-items-center bg-slate-900/50 p-4" @click.self="closeAction">
        <div class="w-full max-w-sm rounded-2xl bg-white p-6 shadow-xl">
          <div class="mx-auto mb-4 grid h-12 w-12 place-items-center rounded-full"
            :class="actionDef.variant === 'danger' ? 'bg-red-100 text-red-600' : 'bg-emerald-100 text-emerald-600'">
            <Icon :name="actionDef.icon" :size="24" />
          </div>
          <h3 class="text-center text-lg font-bold text-slate-900">{{ actionDef.label }}</h3>
          <p v-if="actionDef.confirm" class="mt-2 text-center text-sm text-slate-500">{{ actionDef.confirm }}</p>
          <div v-if="actionDef.prompt" class="mt-4">
            <label class="label">{{ actionDef.prompt.label }}</label>
            <textarea v-model="actionReason" rows="3" class="input" :placeholder="actionDef.prompt.label" />
          </div>
          <div v-if="actionError" class="mt-3 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{{ actionError }}</div>
          <div class="mt-6 flex gap-3">
            <button class="btn-outline flex-1" @click="closeAction">Cancel</button>
            <button class="flex-1" :class="actionDef.variant === 'danger' ? 'btn-danger' : 'btn-primary'"
              :disabled="actionBusy" @click="runAction">
              {{ actionBusy ? 'Working…' : actionDef.label }}
            </button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import Icon from '@/components/Icon.vue'
import api from '@/api/axios'
import { proxied } from '@/api/image'
import { getResource } from '@/config/resources'
import { useAuthStore } from '@/stores/auth'

const props = defineProps({ resourceKey: { type: String, required: true } })

const auth = useAuthStore()
const isAdmin = computed(() => auth.isAdmin)

const cfg = computed(() => getResource(props.resourceKey))

const rows = ref([])
const total = ref(0)
const summary = ref({})
const page = ref(1)
const limit = ref(20)
const search = ref('')
const loading = ref(false)
const showFilters = ref(false)
const filters = reactive({})

const deleteTarget = ref(null)
const deleting = ref(false)

// Config-driven per-row actions (e.g. approve/reject on leave_approvals).
const actionTarget = ref(null)
const actionDef = ref({})
const actionReason = ref('')
const actionBusy = ref(false)
const actionError = ref('')

// rowActions returns the actions visible for a given row. showWhen gates by row
// fields; authApprover restricts to the assigned approver (admins are view-only,
// so they never see approver actions).
function rowActions(row) {
  return (cfg.value.actions || []).filter((a) => {
    if (a.authApprover && !isAssignedApprover(row)) return false
    if (!a.showWhen) return true
    return Object.entries(a.showWhen).every(([k, v]) => row[k] === v)
  })
}
function isAssignedApprover(row) {
  if (isAdmin.value) return false
  const sid = auth.user?.staff_id
  return !!sid && sid === row.staff_id
}
function actionClass(act) {
  if (act.variant === 'danger') return 'text-red-500 hover:bg-red-50'
  if (act.variant === 'success') return 'text-emerald-600 hover:bg-emerald-50'
  return 'text-slate-500'
}
function openAction(act, row) {
  actionDef.value = act
  actionTarget.value = row
  actionReason.value = ''
  actionError.value = ''
}
function closeAction() {
  actionTarget.value = null
}
async function runAction() {
  const act = actionDef.value
  const row = actionTarget.value
  if (!act || !row) return
  actionBusy.value = true
  actionError.value = ''
  try {
    const path = act.path.replace('{id}', encodeURIComponent(row[cfg.value.pk]))
    const body = act.prompt ? { [act.prompt.field]: actionReason.value } : {}
    await api[act.method || 'post'](path, body)
    actionTarget.value = null
    await load()
  } catch (e) {
    actionError.value = e.response?.data?.message || 'Action failed.'
  } finally {
    actionBusy.value = false
  }
}

const brokenImg = reactive(new Set())
function imgOk(src) {
  return !!src && !brokenImg.has(src)
}

// Join an array field (e.g. institute_hierarchy) into a path string, optionally
// dropping the last N entries (e.g. the root ministry).
function chainText(row, col) {
  const arr = row[col.key]
  if (!Array.isArray(arr) || arr.length === 0) return ''
  const drop = col.dropLast || 0
  let items = drop > 0 ? arr.slice(0, Math.max(0, arr.length - drop)) : arr.slice()
  if (col.reverse) items = items.reverse()
  return items.join(col.separator || ' → ')
}

// Resolve a profile field spec: a single key, or an array of keys joined by a space.
function fieldText(row, spec) {
  if (Array.isArray(spec)) {
    return spec
      .map((k) => row[k])
      .filter((v) => v !== null && v !== undefined && v !== '')
      .join(' ')
  }
  return row[spec]
}

// Caption is a small muted line (string key, or { key, prefix }).
function captionText(row, cap) {
  return row[typeof cap === 'string' ? cap : cap.key]
}
function captionPrefix(cap) {
  return typeof cap === 'object' && cap.prefix ? cap.prefix : ''
}
function captionIcon(cap) {
  return typeof cap === 'object' && cap.icon ? cap.icon : ''
}

// Profile-cell tag helpers (a tag is either a field-key string or { key, class }).
function tagKey(t) {
  return typeof t === 'string' ? t : t.key
}
function tagValue(row, t) {
  return row[tagKey(t)]
}
function tagClass(t) {
  return typeof t === 'object' && t.class ? t.class : 'bg-slate-100 text-slate-600'
}

let searchTimer = null

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / limit.value)))
const activeFilterCount = computed(() => Object.values(filters).filter((v) => v !== '' && v != null).length)
const deleteLabel = computed(() => {
  const r = deleteTarget.value
  if (!r) return ''
  return r.name || r.name_kh || r.rank_name || r.position_name || r[cfg.value.pk]
})

function initFilters() {
  for (const k of Object.keys(filters)) delete filters[k]
  for (const f of cfg.value.filters || []) filters[f.key] = ''
}

async function load() {
  loading.value = true
  try {
    const params = { page: page.value, limit: limit.value }
    if (search.value) params.search = search.value
    for (const [k, v] of Object.entries(filters)) {
      if (v !== '' && v != null) params[k] = v
    }
    // Always-applied fixed filters (e.g. staff status_id=2); win over user filters.
    Object.assign(params, cfg.value.baseFilter || {})
    const { data } = await api.get(cfg.value.endpoint, { params })
    rows.value = data.data || []
    total.value = data.total || 0
  } catch (e) {
    rows.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
  loadSummary()
}

// loadSummary fetches the resource's summary counts (if configured). Counts are
// global (not affected by paging/filters), so failures are non-fatal.
async function loadSummary() {
  if (!cfg.value.summary) return
  try {
    const { data } = await api.get(cfg.value.summary.endpoint)
    summary.value = data.data || {}
  } catch (e) {
    summary.value = {}
  }
}

function reset() {
  page.value = 1
  load()
}
function onSearch() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(reset, 350)
}
function onFilterInput() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(reset, 350)
}
function applyFilters() {
  reset()
}
function clearFilters() {
  initFilters()
  reset()
}
function goTo(p) {
  if (p < 1 || p > totalPages.value) return
  page.value = p
  load()
}

async function doDelete() {
  deleting.value = true
  try {
    await api.delete(`${cfg.value.endpoint}/${encodeURIComponent(deleteTarget.value[cfg.value.pk])}`)
    deleteTarget.value = null
    if (rows.value.length === 1 && page.value > 1) page.value--
    await load()
  } finally {
    deleting.value = false
  }
}

function display(v) {
  if (v === null || v === undefined || v === '') return '—'
  return v
}

// Re-init when the route switches between resources (component is reused).
watch(
  () => props.resourceKey,
  () => {
    search.value = ''
    page.value = 1
    initFilters()
    load()
  }
)

onMounted(() => {
  initFilters()
  load()
})
</script>
