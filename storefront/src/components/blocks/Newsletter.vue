<template>
  <section class="py-20 px-6" :style="{ backgroundColor: bgColor }">
    <div class="container mx-auto max-w-4xl text-center">
      <h2 class="text-3xl md:text-4xl font-bold mb-4" :style="{ color: textColor }">
        {{ title }}
      </h2>
      <p class="text-lg md:text-xl mb-8 opacity-90 mx-auto max-w-2xl" :style="{ color: textColor }">
        {{ subtitle }}
      </p>

      <form @submit.prevent="submitNewsletter" class="flex flex-col sm:flex-row max-w-md mx-auto gap-3">
        <input 
          type="email" 
          v-model="email" 
          required 
          placeholder="Enter your email" 
          class="flex-1 px-5 py-3 rounded-full border-none focus:ring-4 focus:outline-none shadow-sm text-gray-900 focus:ring-primary"

        />
        <button 
          type="submit" 
          class="px-8 py-3 rounded-full font-semibold shadow-md hover:shadow-lg transition-all transform hover:-translate-y-0.5 whitespace-nowrap"
          :style="{ backgroundColor: 'var(--color-primary, #10B981)', color: '#fff' }"
        >
          {{ buttonText }}
        </button>
      </form>
      <p v-if="success" class="mt-4 text-sm font-medium" :style="{ color: textColor }">
        Thanks for subscribing!
      </p>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useResponsive } from '../../composables/useResponsive'

const props = defineProps({
  title: [String, Object],
  subtitle: [String, Object],
  buttonText: [String, Object],
  backgroundColor: [String, Object],
  textColor: [String, Object],
})

const { getProp } = useResponsive()

const title = computed(() => getProp(props.title, 'Subscribe to our newsletter'))
const subtitle = computed(() => getProp(props.subtitle, 'Get the latest updates and offers directly in your inbox.'))
const buttonText = computed(() => getProp(props.buttonText, 'Subscribe'))
const bgColor = computed(() => getProp(props.backgroundColor, '#111827')) // Default to dark grey
const textColor = computed(() => getProp(props.textColor, '#ffffff')) // Default to white

const email = ref('')
const success = ref(false)

const submitNewsletter = () => {
  if (email.value) {
    success.value = true
    email.value = ''
    setTimeout(() => { success.value = false }, 3000)
  }
}
</script>
