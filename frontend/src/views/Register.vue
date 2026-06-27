<template>
  <div class="flex min-h-screen items-center justify-center bg-slate-50 p-6">
    <div class="w-full max-w-md">
      <div class="mb-8 text-center">
        <div class="mx-auto mb-4 grid h-12 w-12 place-items-center rounded-2xl bg-primary-500 text-xl font-extrabold text-white shadow-lg shadow-primary-500/30">A</div>
        <h1 class="text-3xl font-bold text-slate-900">Create your account</h1>
        <p class="mt-2 text-slate-500">Get started in less than a minute.</p>
      </div>

      <div class="card p-7">
        <div v-if="error" class="mb-4 flex items-center gap-2 rounded-lg bg-red-50 px-4 py-3 text-sm text-red-700">
          <Icon name="x" :size="18" /> {{ error }}
        </div>

        <form class="space-y-5" @submit.prevent="onSubmit">
          <div>
            <label class="label">Full name</label>
            <input v-model="name" type="text" required placeholder="Jane Doe" class="input" />
          </div>
          <div>
            <label class="label">Email address</label>
            <input v-model="email" type="email" required placeholder="you@example.com" class="input" />
          </div>
          <div>
            <label class="label">Password</label>
            <input v-model="password" type="password" required minlength="6" placeholder="At least 6 characters" class="input" />
          </div>
          <button type="submit" class="btn-primary w-full !py-3" :disabled="loading">
            {{ loading ? 'Creating account…' : 'Create account' }}
          </button>
        </form>
      </div>

      <p class="mt-6 text-center text-sm text-slate-500">
        Already have an account?
        <RouterLink to="/login" class="font-semibold text-primary-600 hover:text-primary-700">Sign in</RouterLink>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import Icon from '@/components/Icon.vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()

const name = ref('')
const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

async function onSubmit() {
  loading.value = true
  error.value = ''
  try {
    await auth.register(name.value, email.value, password.value)
    router.push('/dashboard')
  } catch (e) {
    error.value = e.response?.data?.message || 'Registration failed.'
  } finally {
    loading.value = false
  }
}
</script>
