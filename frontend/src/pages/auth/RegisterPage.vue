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
const name = ref('')
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
    await auth.register(name.value, email.value, password.value)
    await router.push('/')
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not create the account'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <form class="space-y-4" @submit.prevent="submit">
    <div>
      <h1 class="text-lg font-semibold leading-6 text-text">Create account</h1>
      <p class="mt-1 text-sm text-muted">This is the seller login. One account, one shop.</p>
    </div>
    <p v-if="error" class="rounded bg-bad-bg px-3 py-2 text-sm text-bad-fg">{{ error }}</p>
    <a
      v-if="google"
      class="flex h-8 w-full items-center justify-center rounded border border-border bg-surface text-sm font-medium text-text hover:bg-bg"
      href="/api/auth/google/start"
    >Continue with Google</a>
    <Field label="Name"><input v-model="name" :class="controlClass" autocomplete="name" required minlength="2" /></Field>
    <Field label="Email"><input v-model="email" :class="controlClass" type="email" autocomplete="email" required /></Field>
    <Field label="Password" hint="At least 8 characters.">
      <input v-model="password" :class="controlClass" type="password" autocomplete="new-password" required minlength="8" />
    </Field>
    <Button type="submit" :loading="loading" class="w-full">Create account</Button>
    <p class="text-sm text-muted">Already selling? <RouterLink class="font-medium text-text underline" to="/login">Sign in</RouterLink></p>
  </form>
</template>
