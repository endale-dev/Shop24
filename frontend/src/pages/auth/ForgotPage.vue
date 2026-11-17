<script setup lang="ts">
import { ref } from 'vue'
import Button from '@/components/ui/Button.vue'
import Field from '@/components/ui/Field.vue'
import { controlClass } from '@/components/ui/control'
import { requestReset } from '@/features/auth/api'

const email = ref('')
const error = ref('')
const sent = ref(false)
const loading = ref(false)

async function submit() {
  error.value = ''
  loading.value = true
  try {
    await requestReset(email.value)
    sent.value = true
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not send the reset'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <form class="space-y-4" @submit.prevent="submit">
    <div>
      <h1 class="text-lg font-semibold leading-6 text-text">Reset password</h1>
      <p class="mt-1 text-sm text-muted">Enter the email on the account. If it matches, a link is sent. The link works for 30 minutes.</p>
    </div>
    <p v-if="error" class="rounded bg-bad-bg px-3 py-2 text-sm text-bad-fg">{{ error }}</p>
    <p v-if="sent" class="rounded bg-good-bg px-3 py-2 text-sm text-good-fg">
      If that email has an account, the reset link is on its way. When this server is not sending mail, the link is written to the API log and is not shown here.
    </p>
    <Field v-if="!sent" label="Email">
      <input v-model="email" :class="controlClass" type="email" autocomplete="email" required />
    </Field>
    <Button v-if="!sent" type="submit" :loading="loading" class="w-full">Send reset link</Button>
    <p class="text-sm text-muted"><RouterLink class="font-medium text-text underline" to="/login">Back to sign in</RouterLink></p>
  </form>
</template>
