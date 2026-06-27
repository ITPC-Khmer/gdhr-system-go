<template>
  <div class="mx-auto max-w-2xl space-y-5">
    <RouterLink to="/users" class="inline-flex items-center gap-2 text-sm font-medium text-slate-500 hover:text-slate-800">
      <Icon name="arrowLeft" :size="18" /> Back to users
    </RouterLink>

    <div class="card p-6 sm:p-8">
      <h2 class="text-xl font-bold text-slate-900">{{ isEdit ? 'Edit user' : 'Add new user' }}</h2>
      <p class="mt-1 text-sm text-slate-500">{{ isEdit ? 'Update the account details below.' : 'Fill in the details to create a new account.' }}</p>

      <div v-if="error" class="mt-5 flex items-center gap-2 rounded-lg bg-red-50 px-4 py-3 text-sm text-red-700">
        <Icon name="x" :size="18" /> {{ error }}
      </div>

      <form class="mt-6 space-y-5" @submit.prevent="onSubmit">
        <div class="grid grid-cols-1 gap-5 sm:grid-cols-2">
          <div>
            <label class="label">Full name</label>
            <input v-model="form.name" type="text" required placeholder="Jane Doe" class="input" />
          </div>
          <div>
            <label class="label">Email address</label>
            <input v-model="form.email" type="email" required placeholder="jane@example.com" class="input" />
          </div>
        </div>

        <div>
          <label class="label">Password {{ isEdit ? '(leave blank to keep current)' : '' }}</label>
          <input v-model="form.password" type="password" :required="!isEdit" minlength="6" placeholder="••••••••" class="input" />
        </div>

        <div class="grid grid-cols-1 gap-5 sm:grid-cols-2">
          <div>
            <label class="label">Role</label>
            <select v-model="form.role" class="input">
              <option value="user">User</option>
              <option value="admin">Admin</option>
            </select>
          </div>
          <div>
            <label class="label">Status</label>
            <div class="flex items-center gap-3 pt-1.5">
              <button type="button" role="switch" :aria-checked="form.active"
                class="relative inline-flex h-6 w-11 items-center rounded-full transition"
                :class="form.active ? 'bg-primary-500' : 'bg-slate-300'"
                @click="form.active = !form.active">
                <span class="inline-block h-4 w-4 transform rounded-full bg-white transition" :class="form.active ? 'translate-x-6' : 'translate-x-1'" />
              </button>
              <span class="text-sm text-slate-600">{{ form.active ? 'Active' : 'Inactive' }}</span>
            </div>
          </div>
        </div>

        <div>
          <label class="label">Linked staff UID <span class="text-xs font-normal text-slate-400">(optional — authorizes leave approval actions)</span></label>
          <input v-model="form.staff_id" type="text" placeholder="staffs.uid this account approves as" class="input" />
        </div>

        <div class="flex justify-end gap-3 pt-2">
          <RouterLink to="/users" class="btn-outline">Cancel</RouterLink>
          <button type="submit" class="btn-primary" :disabled="loading">
            {{ loading ? 'Saving…' : isEdit ? 'Save changes' : 'Create user' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import Icon from '@/components/Icon.vue'
import api from '@/api/axios'

const route = useRoute()
const router = useRouter()

const isEdit = computed(() => !!route.params.id)
const loading = ref(false)
const error = ref('')

const form = ref({
  name: '',
  email: '',
  password: '',
  role: 'user',
  active: true,
  staff_id: '',
})

async function loadUser() {
  try {
    const { data } = await api.get(`/users/${route.params.id}`)
    form.value = { ...data.data, password: '' }
  } catch (e) {
    error.value = 'Failed to load user.'
  }
}

async function onSubmit() {
  loading.value = true
  error.value = ''
  try {
    const payload = { ...form.value }
    if (isEdit.value && !payload.password) delete payload.password
    if (isEdit.value) {
      await api.put(`/users/${route.params.id}`, payload)
    } else {
      await api.post('/users', payload)
    }
    router.push('/users')
  } catch (e) {
    error.value = e.response?.data?.message || 'Failed to save user.'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (isEdit.value) loadUser()
})
</script>
