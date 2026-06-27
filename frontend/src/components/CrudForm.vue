<template>
  <div class="mx-auto max-w-3xl space-y-5">
    <RouterLink :to="`/${resourceKey}`" class="inline-flex items-center gap-2 text-sm font-medium text-slate-500 hover:text-slate-800">
      <Icon name="arrowLeft" :size="18" /> Back to {{ cfg.title.toLowerCase() }}
    </RouterLink>

    <div class="card p-6 sm:p-8">
      <h2 class="text-xl font-bold text-slate-900">{{ isEdit ? `Edit ${cfg.singular.toLowerCase()}` : `Add ${cfg.singular.toLowerCase()}` }}</h2>
      <p class="mt-1 text-sm text-slate-500">
        {{ isEdit ? 'Update the fields below.' : `Fill in the details to create a new ${cfg.singular.toLowerCase()}.` }}
      </p>

      <div v-if="error" class="mt-5 flex items-center gap-2 rounded-lg bg-red-50 px-4 py-3 text-sm text-red-700">
        <Icon name="x" :size="18" /> {{ error }}
      </div>

      <div v-if="loadingRecord" class="mt-8 text-center text-slate-400">Loading…</div>

      <form v-else class="mt-6 grid grid-cols-1 gap-5 sm:grid-cols-2" @submit.prevent="onSubmit">
        <div v-for="f in cfg.fields" :key="f.key" :class="f.span === 2 ? 'sm:col-span-2' : ''">
          <label class="label">
            {{ f.label }}
            <span v-if="f.required" class="text-red-500">*</span>
            <span v-if="isPk(f) && isEdit" class="text-xs font-normal text-slate-400">(read-only)</span>
          </label>

          <!-- toggle -->
          <div v-if="f.type === 'bool'" class="flex items-center gap-3 pt-1.5">
            <button type="button" role="switch" :aria-checked="!!form[f.key]"
              class="relative inline-flex h-6 w-11 items-center rounded-full transition"
              :class="form[f.key] ? 'bg-primary-500' : 'bg-slate-300'"
              @click="form[f.key] = !form[f.key]">
              <span class="inline-block h-4 w-4 transform rounded-full bg-white transition" :class="form[f.key] ? 'translate-x-6' : 'translate-x-1'" />
            </button>
            <span class="text-sm text-slate-600">{{ form[f.key] ? 'Yes' : 'No' }}</span>
          </div>

          <!-- image: url input + live preview -->
          <div v-else-if="f.type === 'image'" class="flex items-center gap-4">
            <img v-if="imgOk(form[f.key])" :src="proxied(form[f.key])" referrerpolicy="no-referrer"
              class="h-20 w-20 shrink-0 rounded-lg object-cover ring-1 ring-slate-200" @error="brokenImg.add(form[f.key])" />
            <div v-else class="grid h-20 w-20 shrink-0 place-items-center rounded-lg bg-slate-100 text-slate-400">
              <Icon name="user" :size="28" />
            </div>
            <input v-model="form[f.key]" type="text" :placeholder="`${f.label} URL`" class="input" @input="brokenImg.delete(form[f.key])" />
          </div>

          <!-- select -->
          <select v-else-if="f.type === 'select'" v-model="form[f.key]" class="input">
            <option v-for="o in f.options" :key="o.value" :value="o.value">{{ o.label }}</option>
          </select>

          <!-- textarea -->
          <textarea v-else-if="f.type === 'textarea'" v-model="form[f.key]" rows="3" class="input" :placeholder="f.label" />

          <!-- number / date / text -->
          <input v-else v-model="form[f.key]"
            :type="f.type === 'number' ? 'number' : f.type === 'date' ? 'date' : 'text'"
            :required="f.required"
            :disabled="isPk(f) && isEdit"
            :placeholder="f.label"
            class="input"
            :class="(isPk(f) && isEdit) ? 'cursor-not-allowed bg-slate-100' : ''" />
        </div>

        <div class="flex justify-end gap-3 pt-2 sm:col-span-2">
          <RouterLink :to="`/${resourceKey}`" class="btn-outline">Cancel</RouterLink>
          <button type="submit" class="btn-primary" :disabled="saving">
            {{ saving ? 'Saving…' : isEdit ? 'Save changes' : `Create ${cfg.singular.toLowerCase()}` }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import Icon from '@/components/Icon.vue'
import api from '@/api/axios'
import { proxied } from '@/api/image'
import { getResource } from '@/config/resources'

const props = defineProps({ resourceKey: { type: String, required: true } })

const route = useRoute()
const router = useRouter()

const cfg = computed(() => getResource(props.resourceKey))
const isEdit = computed(() => route.params.id != null)

const form = reactive({})
const saving = ref(false)
const loadingRecord = ref(false)
const error = ref('')

const brokenImg = reactive(new Set())
function imgOk(src) {
  return !!src && !brokenImg.has(src)
}

function isPk(f) {
  return f.key === cfg.value.pk
}

function blankForm() {
  for (const k of Object.keys(form)) delete form[k]
  for (const f of cfg.value.fields) {
    if (f.type === 'bool') form[f.key] = false
    else if (f.type === 'number') form[f.key] = null
    else if (f.type === 'select') form[f.key] = f.options?.[0]?.value ?? ''
    else form[f.key] = ''
  }
}

async function loadRecord() {
  loadingRecord.value = true
  try {
    const { data } = await api.get(`${cfg.value.endpoint}/${encodeURIComponent(route.params.id)}`)
    for (const f of cfg.value.fields) {
      const v = data.data[f.key]
      form[f.key] = v === null || v === undefined ? (f.type === 'bool' ? false : '') : v
    }
  } catch (e) {
    error.value = 'Failed to load record.'
  } finally {
    loadingRecord.value = false
  }
}

function buildPayload() {
  const payload = {}
  for (const f of cfg.value.fields) {
    let v = form[f.key]
    if (f.type === 'number') v = v === '' || v === null ? 0 : Number(v)
    if (f.type === 'bool') v = !!v
    payload[f.key] = v
  }
  return payload
}

async function onSubmit() {
  saving.value = true
  error.value = ''
  try {
    const payload = buildPayload()
    if (isEdit.value) {
      await api.put(`${cfg.value.endpoint}/${encodeURIComponent(route.params.id)}`, payload)
    } else {
      await api.post(cfg.value.endpoint, payload)
    }
    router.push(`/${props.resourceKey}`)
  } catch (e) {
    error.value = e.response?.data?.message || 'Failed to save record.'
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  blankForm()
  if (isEdit.value) loadRecord()
})
</script>
