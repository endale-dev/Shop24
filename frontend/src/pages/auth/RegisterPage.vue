<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import Field from '@/components/ui/Field.vue'
import { controlClass } from '@/components/ui/control'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const name = ref('')
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

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
      <h1 class="font-serif text-4xl">Create your account</h1>
      <p class="mt-1 text-sm text-stone-600">One login, one store, a catalog you can publish.</p>
    </div>
    <p v-if="error" class="rounded-xl bg-red-50 px-3 py-2 text-sm text-red-800">{{ error }}</p>
    <Field label="Name"><input v-model="name" :class="controlClass" required minlength="2" /></Field>
    <Field label="Email"><input v-model="email" :class="controlClass" type="email" autocomplete="email" required /></Field>
    <Field label="Password" hint="At least 8 characters."><input v-model="password" :class="controlClass" type="password" autocomplete="new-password" required minlength="8" /></Field>
    <Button type="submit" :loading="loading" class="w-full">Create account</Button>
    <p class="text-sm text-stone-600">Already selling? <RouterLink class="font-semibold text-rust" to="/login">Sign in</RouterLink></p>
  </form>
</template>
