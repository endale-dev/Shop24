<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import Field from '@/components/ui/Field.vue'
import { controlClass } from '@/components/ui/control'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function submit() {
  error.value = ''
  loading.value = true
  try {
    await auth.login(email.value, password.value)
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    await router.push(redirect)
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not sign in'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <form class="space-y-4" @submit.prevent="submit">
    <div>
      <h1 class="font-serif text-4xl">Sign in</h1>
      <p class="mt-1 text-sm text-stone-600">Back to the shop you are building.</p>
    </div>
    <p v-if="error" class="rounded-xl bg-red-50 px-3 py-2 text-sm text-red-800">{{ error }}</p>
    <Field label="Email"><input v-model="email" :class="controlClass" type="email" autocomplete="email" required /></Field>
    <Field label="Password"><input v-model="password" :class="controlClass" type="password" autocomplete="current-password" required /></Field>
    <Button type="submit" :loading="loading" class="w-full">Sign in</Button>
    <p class="text-sm text-stone-600">New here? <RouterLink class="font-semibold text-rust" to="/register">Create an account</RouterLink></p>
  </form>
</template>
