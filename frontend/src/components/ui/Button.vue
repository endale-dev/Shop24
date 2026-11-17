<script setup lang="ts">
import { RouterLink } from 'vue-router'

withDefaults(
  defineProps<{
    variant?: 'primary' | 'secondary' | 'ghost' | 'danger'
    type?: 'button' | 'submit'
    disabled?: boolean
    loading?: boolean
    to?: string
  }>(),
  { variant: 'primary', type: 'button', disabled: false, loading: false },
)

const styles = {
  primary: 'bg-accent text-white hover:bg-accent-hover',
  secondary: 'border border-border bg-surface text-text hover:bg-bg',
  ghost: 'text-text hover:bg-bg',
  danger: 'bg-danger text-white hover:bg-danger-hover',
}
</script>

<template>
  <component
    :is="to ? RouterLink : 'button'"
    :to="to"
    :type="to ? undefined : type"
    :disabled="to ? undefined : disabled || loading"
    class="inline-flex h-8 items-center justify-center gap-2 rounded px-3 text-sm font-medium focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:cursor-not-allowed disabled:opacity-50"
    :class="[styles[variant], (disabled || loading) && to ? 'pointer-events-none opacity-50' : '']"
  >
    <span v-if="loading" class="h-3.5 w-3.5 animate-spin rounded-full border-2 border-current border-r-transparent" />
    <slot />
  </component>
</template>
