<template>
  <div class="mx-auto max-w-3xl space-y-5">
    <RouterLink to="/leaves" class="inline-flex items-center gap-2 text-sm font-medium text-slate-500 hover:text-slate-800">
      <Icon name="arrowLeft" :size="18" /> Back to leave requests
    </RouterLink>

    <div v-if="loading" class="card p-10 text-center text-slate-400">Loading…</div>
    <div v-else-if="error" class="card flex items-center gap-2 p-6 text-sm text-red-700">
      <Icon name="x" :size="18" /> {{ error }}
    </div>

    <template v-else>
      <!-- Summary -->
      <div class="card p-6 sm:p-8">
        <div class="flex items-start justify-between gap-4">
          <div>
            <h2 class="text-xl font-bold text-slate-900">Leave request #{{ leave.id }}</h2>
            <p class="text-sm text-slate-500">{{ leave.leave_type_name || ('Type ' + leave.leave_type_id) }}</p>
          </div>
          <span class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-semibold" :class="statusClass(leave.status)">
            <span class="h-1.5 w-1.5 rounded-full" :class="dotClass(leave.status)" />
            {{ leave.status }}
          </span>
        </div>

        <dl class="mt-5 grid grid-cols-2 gap-4 text-sm sm:grid-cols-3">
          <div><dt class="text-xs text-slate-400">Staff</dt><dd class="font-mono text-slate-700">{{ leave.staff_id }}</dd></div>
          <div><dt class="text-xs text-slate-400">Start</dt><dd class="text-slate-700">{{ leave.start_date || '—' }}</dd></div>
          <div><dt class="text-xs text-slate-400">End</dt><dd class="text-slate-700">{{ leave.end_date || '—' }}</dd></div>
          <div><dt class="text-xs text-slate-400">Total days</dt><dd class="text-slate-700">{{ leave.total_day }}</dd></div>
          <div v-if="leave.ref_number"><dt class="text-xs text-slate-400">Reference</dt><dd class="text-slate-700">{{ leave.ref_number }}</dd></div>
          <div v-if="leave.phone"><dt class="text-xs text-slate-400">Phone</dt><dd class="text-slate-700">{{ leave.phone }}</dd></div>
        </dl>

        <div v-if="leave.reason" class="mt-4 rounded-lg bg-slate-50 p-3 text-sm text-slate-600">
          <span class="font-medium text-slate-500">Reason: </span>{{ leave.reason }}
        </div>
        <div v-if="leave.reject_reason" class="mt-3 rounded-lg bg-red-50 p-3 text-sm text-red-700">
          <span class="font-medium">Reject reason: </span>{{ leave.reject_reason }}
        </div>
        <p v-if="leave.approved_by" class="mt-3 text-xs text-slate-400">
          Finalized by <span class="font-mono">{{ leave.approved_by }}</span><span v-if="leave.approved_at"> · {{ fmt(leave.approved_at) }}</span>
        </p>
      </div>

      <!-- Approval timeline -->
      <div class="card p-6 sm:p-8">
        <h3 class="mb-5 text-lg font-bold text-slate-900">Approval timeline</h3>

        <p v-if="!approvals.length" class="text-sm text-slate-400">No approval steps for this request.</p>

        <ol v-else class="relative ml-3 space-y-7 border-l border-slate-200">
          <li v-for="ap in approvals" :key="ap.id" class="ml-6">
            <span class="absolute -left-[9px] grid h-4 w-4 place-items-center rounded-full ring-4 ring-white" :class="dotClass(ap.status)" />

            <div class="flex items-center justify-between gap-3">
              <div class="min-w-0">
                <p class="font-semibold text-slate-800">
                  Level {{ ap.approve_level ?? ap.l_level ?? '—' }}
                  <span class="ml-1 text-xs font-normal text-slate-400">{{ ap.role_name }}</span>
                </p>
                <p class="truncate font-mono text-sm text-slate-500">{{ ap.staff_id || '—' }}</p>
              </div>
              <span class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-semibold" :class="statusClass(ap.status)">
                {{ ap.status }}
              </span>
            </div>

            <p class="mt-1 text-xs text-slate-400">
              <span v-if="ap.institute_id">Institute <span class="font-mono">{{ ap.institute_id }}</span> · </span>
              created {{ fmt(ap.created_at) }}
              <span v-if="ap.approved_at"> · acted {{ fmt(ap.approved_at) }}</span>
            </p>

            <!-- Inline actions for a pending step (assigned approver only;
                 admins are view-only) -->
            <div v-if="ap.status === 'pending' && canAct(ap)" class="mt-3">
              <div v-if="rejectingId === ap.id" class="space-y-2">
                <textarea v-model="rejectReason" rows="2" class="input" placeholder="Reason for rejection (optional)" />
                <div class="flex gap-2">
                  <button class="btn-danger !py-1.5 text-sm" :disabled="busy" @click="act(ap, 'reject')">{{ busy ? 'Working…' : 'Confirm reject' }}</button>
                  <button class="btn-outline !py-1.5 text-sm" @click="cancelReject">Cancel</button>
                </div>
              </div>
              <div v-else class="flex gap-2">
                <button class="btn-primary !py-1.5 text-sm" :disabled="busy" @click="act(ap, 'approve')">
                  <Icon name="check" :size="16" /> Approve
                </button>
                <button class="btn-outline !py-1.5 text-sm text-red-600" :disabled="busy" @click="startReject(ap)">
                  <Icon name="x" :size="16" /> Reject
                </button>
              </div>
              <p v-if="actionError && actingId === ap.id" class="mt-2 text-sm text-red-600">{{ actionError }}</p>
            </div>
          </li>
        </ol>
      </div>

      <!-- Admin break-glass override: only for admins, only while pending -->
      <div v-if="auth.isAdmin && leave.status === 'pending'" class="card border border-amber-200 p-6 sm:p-8">
        <div class="flex items-center gap-2">
          <Icon name="shield" :size="18" class="text-amber-600" />
          <h3 class="text-lg font-bold text-slate-900">Admin override (break-glass)</h3>
        </div>
        <p class="mt-1 text-sm text-slate-500">
          Use only if this request is stuck (no approver could be routed, or the assigned approver can't act).
          This bypasses the normal approval chain and is recorded for audit.
        </p>
        <textarea v-model="overrideNote" rows="2" class="input mt-4" placeholder="Reason / note (recommended)" />
        <div v-if="overrideError" class="mt-3 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{{ overrideError }}</div>
        <div class="mt-4 flex gap-2">
          <button class="btn-primary !py-1.5 text-sm" :disabled="overrideBusy" @click="doOverride('approve')">
            <Icon name="check" :size="16" /> Force approve
          </button>
          <button class="btn-danger !py-1.5 text-sm" :disabled="overrideBusy" @click="doOverride('reject')">
            <Icon name="x" :size="16" /> Force reject
          </button>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import Icon from '@/components/Icon.vue'
import api from '@/api/axios'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const auth = useAuthStore()

// Only the assigned approver may act; admins have view-only access to leave.
function canAct(ap) {
  if (auth.isAdmin) return false
  const sid = auth.user?.staff_id
  return !!sid && sid === ap.staff_id
}

const leave = ref({})
const approvals = ref([])
const loading = ref(true)
const error = ref('')

const busy = ref(false)
const actionError = ref('')
const actingId = ref(null)
const rejectingId = ref(null)
const rejectReason = ref('')

// Admin break-glass override (stuck requests only).
const overrideNote = ref('')
const overrideBusy = ref(false)
const overrideError = ref('')

async function doOverride(kind) {
  overrideBusy.value = true
  overrideError.value = ''
  try {
    const id = route.params.id
    const body = kind === 'approve' ? { note: overrideNote.value } : { reject_reason: overrideNote.value }
    await api.post(`/leaves/${encodeURIComponent(id)}/override-${kind}`, body)
    overrideNote.value = ''
    await load()
  } catch (e) {
    overrideError.value = e.response?.data?.message || `Failed to override-${kind}.`
  } finally {
    overrideBusy.value = false
  }
}

function statusClass(s) {
  if (s === 'approved') return 'bg-emerald-100 text-emerald-700'
  if (s === 'rejected') return 'bg-red-100 text-red-700'
  return 'bg-amber-100 text-amber-700'
}
function dotClass(s) {
  if (s === 'approved') return 'bg-emerald-500'
  if (s === 'rejected') return 'bg-red-500'
  return 'bg-amber-400'
}
function fmt(ts) {
  if (!ts) return '—'
  const d = new Date(ts)
  return isNaN(d) ? ts : d.toLocaleString()
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    const id = route.params.id
    const [lv, ap] = await Promise.all([
      api.get(`/leaves/${encodeURIComponent(id)}`),
      api.get('/leave-approvals', { params: { leave_id: id, limit: 200 } }),
    ])
    leave.value = lv.data.data || {}
    // Chronological chain order (creation order).
    approvals.value = (ap.data.data || []).slice().sort((a, b) => a.id - b.id)
  } catch (e) {
    error.value = e.response?.data?.message || 'Failed to load leave request.'
  } finally {
    loading.value = false
  }
}

function startReject(ap) {
  rejectingId.value = ap.id
  rejectReason.value = ''
  actionError.value = ''
}
function cancelReject() {
  rejectingId.value = null
}

async function act(ap, kind) {
  busy.value = true
  actingId.value = ap.id
  actionError.value = ''
  try {
    const body = kind === 'reject' ? { reject_reason: rejectReason.value } : {}
    await api.post(`/leave-approvals/${ap.id}/${kind}`, body)
    rejectingId.value = null
    await load()
  } catch (e) {
    actionError.value = e.response?.data?.message || `Failed to ${kind} task.`
  } finally {
    busy.value = false
  }
}

onMounted(load)
</script>
