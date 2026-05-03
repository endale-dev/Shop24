<template>
  <section class="py-16 px-6 bg-white">
    <div class="container mx-auto">
      <div class="text-center mb-16 max-w-2xl mx-auto">
        <h2 class="text-3xl font-bold mb-4" style="color: var(--color-text-primary)">{{ title }}</h2>
      </div>

      <div class="grid gap-8" :class="gridClass">
        <div 
          v-for="(feature, index) in features" 
          :key="index"
          class="flex flex-col items-center text-center p-6 rounded-2xl bg-gray-50/50 hover:bg-white hover:shadow-lg transition-all duration-300 border border-transparent hover:border-gray-100"
        >
          <div class="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center mb-6 text-primary">
            <component 
              :is="feature.icon" 
              v-if="feature.icon" 
              class="w-8 h-8" 
              style="color: var(--color-primary, #10B981)" 
            />
            <span v-else class="text-2xl">✨</span>
          </div>
          <h3 class="text-xl font-bold mb-3" style="color: var(--color-text-primary)">
            {{ feature.title }}
          </h3>
          <p class="text-gray-600 leading-relaxed" style="color: var(--color-text-secondary)">
            {{ feature.description }}
          </p>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useResponsive } from '../../composables/useResponsive'
import { 
  CheckCircleIcon, 
  TruckIcon, 
  ShieldCheckIcon, 
  StarIcon, 
  HeartIcon, 
  SparklesIcon 
} from '@heroicons/vue/24/outline'

const props = defineProps({
  title: [String, Object],
  columns: [String, Object, Number],
  feature1Title: [String, Object],
  feature1Description: [String, Object],
  feature1Icon: [String, Object],
  feature2Title: [String, Object],
  feature2Description: [String, Object],
  feature3Title: [String, Object],
  feature3Description: [String, Object],
})

const { getProp } = useResponsive()

const title = computed(() => getProp(props.title, 'Why Choose Us'))
const columns = computed(() => String(getProp(props.columns, '3')))

const gridClass = computed(() => {
  if (columns.value === '2') return 'grid-cols-1 md:grid-cols-2'
  if (columns.value === '4') return 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-4'
  return 'grid-cols-1 md:grid-cols-3'
})

// Map string icon names to HeroIcons
const getIconComponent = (iconName: string) => {
  const icons: Record<string, any> = {
    'TruckIcon': TruckIcon,
    'ShieldCheckIcon': ShieldCheckIcon,
    'StarIcon': StarIcon,
    'HeartIcon': HeartIcon,
    'SparklesIcon': SparklesIcon,
    'CheckCircleIcon': CheckCircleIcon
  }
  return icons[iconName] || CheckCircleIcon
}

const features = computed(() => {
  const list = []
  
  if (getProp(props.feature1Title)) {
    list.push({
      title: getProp(props.feature1Title),
      description: getProp(props.feature1Description, ''),
      icon: getIconComponent(getProp(props.feature1Icon, 'CheckCircleIcon'))
    })
  }
  
  if (getProp(props.feature2Title)) {
    list.push({
      title: getProp(props.feature2Title),
      description: getProp(props.feature2Description, ''),
      icon: getIconComponent('TruckIcon') // simplified default
    })
  }
  
  if (getProp(props.feature3Title)) {
    list.push({
      title: getProp(props.feature3Title),
      description: getProp(props.feature3Description, ''),
      icon: getIconComponent('ShieldCheckIcon')
    })
  }

  // Fallback if empty
  if (list.length === 0) {
    list.push({ title: 'Quality Products', description: 'Curated selection of the best items.', icon: SparklesIcon })
    list.push({ title: 'Fast Shipping', description: 'Get your orders delivered quickly.', icon: TruckIcon })
    list.push({ title: 'Secure Checkout', description: '100% secure payment processing.', icon: ShieldCheckIcon })
  }
  
  return list
})
</script>
