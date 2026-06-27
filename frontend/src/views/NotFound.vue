<template>
  <div class="grid min-h-screen place-items-center bg-slate-50 p-6 text-center">
    <div>
      <p class="text-7xl font-extrabold text-primary-500">404</p>
      <h1 class="mt-4 text-2xl font-bold text-slate-900">Page not found</h1>
      <p class="mt-2 text-slate-500">The page you're looking for doesn't exist or was moved.</p>
      <RouterLink to="/dashboard" class="btn-primary mt-6">Back to dashboard</RouterLink>
    </div>
  </div>
</template>

<script setup>
import { RouterLink } from 'vue-router'
</script>
