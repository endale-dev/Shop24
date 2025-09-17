<template>
  <div :style="containerStyle"><slot /></div>
</template>

<script setup>
const props = defineProps({ maxWidth: { type: String, default: '1200px' }, padding: { type: String, default: '16px' } })
const containerStyle = {
  maxWidth: props.maxWidth,
  margin: '0 auto',
  padding: props.padding
}
</script>
