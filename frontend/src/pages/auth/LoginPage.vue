<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import Field from '@/components/ui/Field.vue'
import { controlClass } from '@/components/ui/control'
import { authProviders } from '@/features/auth/api'
import { googleErrorMessage } from '@/features/auth/google'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()
const email = ref('')
const password = ref('')
const error = ref(googleErrorMessage(route.query.error))
const loading = ref(false)
const google = ref(false)

onMounted(async () => {
  try {
    google.value = (await authProviders()).google
  } catch {
    google.value = false
  }
})

async function submit() {
  error.value = ''
  loading.value = true
  try {
    await auth.login(email.value, password.value)
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    const next = redirect.startsWith('/') && !redirect.startsWith('//') ? redirect : '/'
    await router.push(next)
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
      <h1 class="text-lg font-semibold leading-6 text-text">Sign in</h1>
      <p class="mt-1 text-sm text-muted">Use the email and password for this desk.</p>
    </div>
    <p v-if="error" class="rounded bg-bad-bg px-3 py-2 text-sm text-bad-fg">{{ error }}</p>
    <a
      v-if="google"
      class="flex h-8 w-full items-center justify-center rounded border border-border bg-surface text-sm font-medium text-text hover:bg-bg"
      href="/api/auth/google/start"
    >Continue with Google</a>
    <Field label="Email"><input v-model="email" :class="controlClass" type="email" autocomplete="email" required /></Field>
    <Field label="Password">
      <input v-model="password" :class="controlClass" type="password" autocomplete="current-password" required />
    </Field>
    <Button type="submit" :loading="loading" class="w-full">Sign in</Button>
    <p class="text-sm text-muted">
      <RouterLink class="font-medium text-text underline" to="/forgot">Forgot password</RouterLink>
    </p>
    <p class="text-sm text-muted">New here? <RouterLink class="font-medium text-text underline" to="/register">Create an account</RouterLink></p>
  </form>
</template>
