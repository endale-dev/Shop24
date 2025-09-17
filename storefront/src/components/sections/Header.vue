<template>
  <header :style="headerStyle">
    <img v-if="logoUrl" :src="logoUrl" alt="logo" style="height:40px" />
    <slot />
  </header>
</template>

<script setup>
const props = defineProps({ logoUrl: { type: String, default: '' }, showSearch: { type: Boolean, default: true }, showCart: { type: Boolean, default: true } })

const headerStyle = {
  display: 'flex',
  alignItems: 'center',
  gap: '12px',
  padding: '8px 16px'
}
</script>
