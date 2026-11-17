<script setup lang="ts">
import Logo from '@/components/ui/Logo.vue'
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-bg px-4 py-12 text-text">
    <div class="w-full max-w-[360px]">
      <RouterLink to="/" class="mb-4 inline-flex items-center gap-2 text-text">
        <Logo />
        <span class="text-sm font-medium leading-none text-text">Sable</span>
      </RouterLink>
      <RouterView />
    </div>
  </div>
</template>
