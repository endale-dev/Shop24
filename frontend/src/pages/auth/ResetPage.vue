<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import Field from '@/components/ui/Field.vue'
import { controlClass } from '@/components/ui/control'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()
const password = ref('')
const confirm = ref('')
const error = ref('')
const loading = ref(false)
const token = computed(() => (typeof route.query.token === 'string' ? route.query.token : ''))

async function submit() {
  error.value = ''
  if (password.value !== confirm.value) {
    error.value = 'Those passwords do not match.'
    return
  }
  loading.value = true
  try {
    await auth.resetPassword(token.value, password.value)
    await router.push('/')
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not reset the password'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <form class="space-y-4" @submit.prevent="submit">
    <div>
      <h1 class="text-lg font-semibold leading-6 text-text">Choose a new password</h1>
      <p class="mt-1 text-sm text-muted">This replaces the password on the account and signs you in.</p>
    </div>
    <p v-if="!token" class="rounded bg-bad-bg px-3 py-2 text-sm text-bad-fg">This reset link is missing its token.</p>
    <p v-if="error" class="rounded bg-bad-bg px-3 py-2 text-sm text-bad-fg">{{ error }}</p>
    <Field label="New password" hint="At least 8 characters.">
      <input v-model="password" :class="controlClass" type="password" autocomplete="new-password" required minlength="8" :disabled="!token" />
    </Field>
    <Field label="Repeat password">
      <input v-model="confirm" :class="controlClass" type="password" autocomplete="new-password" required minlength="8" :disabled="!token" />
    </Field>
    <Button type="submit" :loading="loading" :disabled="!token" class="w-full">Save password</Button>
    <p class="text-sm text-muted"><RouterLink class="font-medium text-text underline" to="/forgot">Request another link</RouterLink></p>
  </form>
</template>
