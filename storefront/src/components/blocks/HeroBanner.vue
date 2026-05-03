<template>
  <section 
    class="relative flex items-center justify-center overflow-hidden transition-all duration-300"
    :class="[heightClass, alignmentClass]"
    :style="sectionStyle"
  >
    <!-- Background overlay if image exists -->
    <div v-if="bgImage" class="absolute inset-0 bg-black/40 z-0"></div>

    <div class="relative z-10 container mx-auto px-6 py-12 flex flex-col gap-6 animate-fade-in-up" :class="alignmentClass">
      <h1 class="text-4xl md:text-6xl font-extrabold tracking-tight drop-shadow-md" :style="{ color: textColor }">
        {{ title }}
      </h1>
      
      <p v-if="subtitle" class="text-lg md:text-xl font-medium max-w-2xl opacity-90 drop-shadow" :style="{ color: textColor }">
        {{ subtitle }}
      </p>
      
      <div v-if="buttonText" class="mt-4">
        <a 
          :href="buttonLink || '#'"
          class="inline-flex items-center justify-center px-8 py-3 text-base font-semibold rounded-full shadow-lg hover:shadow-xl transition-all hover:-translate-y-1"
          :style="{ backgroundColor: 'var(--color-primary, #10B981)', color: '#fff' }"
        >
          {{ buttonText }}
        </a>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useResponsive } from '../../composables/useResponsive'

const props = defineProps({
  title: [String, Object],
  subtitle: [String, Object],
  buttonText: [String, Object],
  buttonLink: [String, Object],
  backgroundImage: [String, Object],
  backgroundColor: [String, Object],
  textColor: [String, Object],
  alignment: [String, Object],
  height: [String, Object],
})

const { getProp } = useResponsive()

const title = computed(() => getProp(props.title, 'Welcome'))
const subtitle = computed(() => getProp(props.subtitle, ''))
const buttonText = computed(() => getProp(props.buttonText, ''))
const buttonLink = computed(() => getProp(props.buttonLink, '#'))
const bgImage = computed(() => getProp(props.backgroundImage, ''))
const bgColor = computed(() => getProp(props.backgroundColor, 'var(--color-surface, #f9fafb)'))
const textColor = computed(() => getProp(props.textColor, bgImage.value ? '#ffffff' : 'var(--color-text-primary, #111827)'))

const alignmentClass = computed(() => {
  const align = getProp(props.alignment, 'center')
  if (align === 'left') return 'items-start text-left'
  if (align === 'right') return 'items-end text-right'
  return 'items-center text-center'
})

const heightClass = computed(() => {
  const h = getProp(props.height, 'medium')
  if (h === 'small') return 'min-h-[300px]'
  if (h === 'large') return 'min-h-[700px]'
  if (h === 'full') return 'min-h-screen'
  return 'min-h-[500px]' // medium
})

const sectionStyle = computed(() => {
  return {
    backgroundColor: bgColor.value,
    backgroundImage: bgImage.value ? `url(${bgImage.value})` : 'none',
    backgroundSize: 'cover',
    backgroundPosition: 'center',
  }
})
</script>

<style scoped>
.animate-fade-in-up {
  animation: fadeInUp 0.8s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
