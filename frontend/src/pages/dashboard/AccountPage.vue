<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import Field from '@/components/ui/Field.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import { controlClass } from '@/components/ui/control'
import { authProviders } from '@/features/auth/api'
import { googleErrorMessage } from '@/features/auth/google'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const route = useRoute()
const name = ref(auth.user?.name ?? '')
const current = ref('')
const password = ref('')
const confirm = ref('')
const nameError = ref('')
const nameMessage = ref('')
const passwordError = ref('')
const passwordMessage = ref('')
const googleNotice = ref(googleErrorMessage(route.query.error))
const nameLoading = ref(false)
const passwordLoading = ref(false)
const google = ref(false)

onMounted(async () => {
  try {
    google.value = (await authProviders()).google
  } catch {
    google.value = false
  }
})

async function saveName() {
  nameError.value = ''
  nameMessage.value = ''
  nameLoading.value = true
  try {
    await auth.saveName(name.value)
    name.value = auth.user?.name ?? name.value
    nameMessage.value = 'Name saved.'
  } catch (err) {
    nameError.value = err instanceof Error ? err.message : 'Could not save the name'
  } finally {
    nameLoading.value = false
  }
}

async function savePassword() {
  passwordError.value = ''
  passwordMessage.value = ''
  if (password.value !== confirm.value) {
    passwordError.value = 'Those passwords do not match.'
    return
  }
  passwordLoading.value = true
  try {
    await auth.savePassword(current.value, password.value)
    current.value = ''
    password.value = ''
    confirm.value = ''
    passwordMessage.value = 'Password saved.'
  } catch (err) {
    passwordError.value = err instanceof Error ? err.message : 'Could not change the password'
  } finally {
    passwordLoading.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-xl space-y-6">
    <PageHeader title="Account" subtitle="The name on this desk, the password, and Google." />

    <form class="space-y-4 rounded border border-border bg-surface p-4" @submit.prevent="saveName">
      <h2 class="text-sm font-medium text-text">Name</h2>
      <p v-if="nameError" class="rounded bg-bad-bg px-3 py-2 text-sm text-bad-fg">{{ nameError }}</p>
      <p v-if="nameMessage" class="rounded bg-good-bg px-3 py-2 text-sm text-good-fg">{{ nameMessage }}</p>
      <Field label="Name"><input v-model="name" :class="controlClass" autocomplete="name" required minlength="2" /></Field>
      <Field label="Email">
        <input :value="auth.user?.email" :class="controlClass" type="email" readonly />
      </Field>
      <p class="text-xs text-muted">Email is the sign-in address. It is not changed here.</p>
      <Button type="submit" :loading="nameLoading">Save name</Button>
    </form>

    <form class="space-y-4 rounded border border-border bg-surface p-4" @submit.prevent="savePassword">
      <h2 class="text-sm font-medium text-text">{{ auth.user?.hasPassword ? 'Change password' : 'Set a password' }}</h2>
      <p class="text-sm text-muted">
        {{ auth.user?.hasPassword ? 'Enter the current password, then the new one.' : 'This account has no password yet. Setting one lets you sign in without Google.' }}
      </p>
      <p v-if="passwordError" class="rounded bg-bad-bg px-3 py-2 text-sm text-bad-fg">{{ passwordError }}</p>
      <p v-if="passwordMessage" class="rounded bg-good-bg px-3 py-2 text-sm text-good-fg">{{ passwordMessage }}</p>
      <Field v-if="auth.user?.hasPassword" label="Current password">
        <input v-model="current" :class="controlClass" type="password" autocomplete="current-password" required />
      </Field>
      <Field label="New password" hint="At least 8 characters.">
        <input v-model="password" :class="controlClass" type="password" autocomplete="new-password" required minlength="8" />
      </Field>
      <Field label="Repeat password">
        <input v-model="confirm" :class="controlClass" type="password" autocomplete="new-password" required minlength="8" />
      </Field>
      <Button type="submit" :loading="passwordLoading">{{ auth.user?.hasPassword ? 'Change password' : 'Set password' }}</Button>
    </form>

    <section class="space-y-3 rounded border border-border bg-surface p-4">
      <h2 class="text-sm font-medium text-text">Google</h2>
      <p v-if="googleNotice" class="rounded bg-bad-bg px-3 py-2 text-sm text-bad-fg">{{ googleNotice }}</p>
      <p v-if="auth.user?.googleLinked" class="text-sm text-muted">Google is connected for {{ auth.user?.email }}.</p>
      <p v-else class="text-sm text-muted">Google is not connected. Connecting uses the same email as this account and does not replace it.</p>
      <a
        v-if="google && !auth.user?.googleLinked"
        class="inline-flex h-8 items-center justify-center rounded border border-border bg-surface px-3 text-sm font-medium text-text hover:bg-bg"
        href="/api/auth/google/start?mode=link"
      >Connect Google</a>
      <p v-if="!google" class="text-sm text-muted">Google sign-in is not set up on this server.</p>
    </section>
  </div>
</template>
