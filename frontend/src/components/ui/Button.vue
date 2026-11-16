<script setup lang="ts">
withDefaults(
  defineProps<{
    variant?: 'primary' | 'secondary' | 'ghost' | 'danger'
    type?: 'button' | 'submit'
    disabled?: boolean
    loading?: boolean
  }>(),
  { variant: 'primary', type: 'button', disabled: false, loading: false },
)

const styles = {
  primary: 'bg-rust text-white hover:bg-orange-800',
  secondary: 'border border-line bg-white text-ink hover:bg-stone-50',
  ghost: 'text-stone-600 hover:bg-stone-200/60',
  danger: 'bg-red-800 text-white hover:bg-red-900',
}
</script>

<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    class="inline-flex items-center justify-center gap-2 rounded-full px-4 py-2 text-sm font-semibold transition disabled:cursor-not-allowed disabled:opacity-50"
    :class="styles[variant]"
  >
    <span v-if="loading" class="h-3.5 w-3.5 animate-spin rounded-full border-2 border-current border-r-transparent" />
    <slot />
  </button>
</template>
