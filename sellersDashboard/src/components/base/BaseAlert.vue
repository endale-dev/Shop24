<template>
  <div v-if="show" class="rounded-md p-4" :class="alertClasses">
    <div class="flex">
      <div class="flex-shrink-0">
        <component :is="iconComponent" class="h-5 w-5" :class="iconColorClass" />
      </div>
      <div class="ml-3 flex-1">
        <h3 v-if="title" class="text-sm font-medium" :class="titleColorClass">
          {{ title }}
        </h3>
        <div class="text-sm" :class="[title ? 'mt-2' : '', messageColorClass]">
          <slot>
            <p>{{ message }}</p>
          </slot>
        </div>
      </div>
      <div v-if="dismissible" class="ml-auto pl-3">
        <button @click="dismiss" class="-mx-1.5 -my-1.5 inline-flex rounded-md p-1.5 focus:outline-none focus:ring-2 focus:ring-offset-2" :class="dismissButtonClass">
          <span class="sr-only">Dismiss</span>
          <XMarkIcon class="h-5 w-5" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { 
  CheckCircleIcon, 
  ExclamationTriangleIcon, 
  InformationCircleIcon, 
  XCircleIcon,
  XMarkIcon 
} from '@heroicons/vue/24/outline'

const props = defineProps({
  variant: {
    type: String,
    default: 'info',
    validator: (value) => ['success', 'warning', 'error', 'info'].includes(value)
  },
  title: {
    type: String,
    default: ''
  },
  message: {
    type: String,
    default: ''
  },
  dismissible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['dismiss'])

const show = ref(true)

const dismiss = () => {
  show.value = false
  emit('dismiss')
}

const variantConfig = {
  success: {
    icon: CheckCircleIcon,
    alertClass: 'bg-green-50',
    iconColor: 'text-green-400',
    titleColor: 'text-green-800',
    messageColor: 'text-green-700',
    dismissButton: 'text-green-500 hover:bg-green-100 focus:ring-green-600'
  },
  warning: {
    icon: ExclamationTriangleIcon,
    alertClass: 'bg-amber-50',
    iconColor: 'text-amber-400',
    titleColor: 'text-amber-800',
    messageColor: 'text-amber-700',
    dismissButton: 'text-amber-500 hover:bg-amber-100 focus:ring-amber-600'
  },
  error: {
    icon: XCircleIcon,
    alertClass: 'bg-red-50',
    iconColor: 'text-red-400',
    titleColor: 'text-red-800',
    messageColor: 'text-red-700',
    dismissButton: 'text-red-500 hover:bg-red-100 focus:ring-red-600'
  },
  info: {
    icon: InformationCircleIcon,
    alertClass: 'bg-blue-50',
    iconColor: 'text-blue-400',
    titleColor: 'text-blue-800',
    messageColor: 'text-blue-700',
    dismissButton: 'text-blue-500 hover:bg-blue-100 focus:ring-blue-600'
  }
}

const config = computed(() => variantConfig[props.variant])
const iconComponent = computed(() => config.value.icon)
const alertClasses = computed(() => config.value.alertClass)
const iconColorClass = computed(() => config.value.iconColor)
const titleColorClass = computed(() => config.value.titleColor)
const messageColorClass = computed(() => config.value.messageColor)
const dismissButtonClass = computed(() => config.value.dismissButton)
</script>
