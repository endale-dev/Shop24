<template>
  <div class="h-full overflow-y-auto">
    <div class="p-4 border-b border-gray-200">
      <h3 class="font-bold text-gray-800">Theme Settings</h3>
      <p class="text-xs text-gray-500">Global styles for your store</p>
    </div>

    <!-- Colors -->
    <div class="p-4 border-b border-gray-200">
      <div 
        @click="toggleGroup('colors')"
        class="flex items-center justify-between cursor-pointer mb-2"
      >
        <h4 class="font-medium text-gray-700">Colors</h4>
        <span class="text-xs text-indigo-600">{{ openGroups.colors ? 'Hide' : 'Show' }}</span>
      </div>
      
      <div v-if="openGroups.colors" class="space-y-4">
        <!-- Primary -->
        <div>
          <label class="block text-xs font-medium text-gray-700 mb-1">Primary Color</label>
          <div class="flex items-center gap-2">
            <input 
              type="color" 
              v-model="theme.colors.primary"
              class="w-8 h-8 rounded cursor-pointer border-0 p-0"
              @change="updateTheme"
            />
            <input 
              type="text" 
              v-model="theme.colors.primary"
              class="flex-1 text-sm border border-gray-300 rounded px-2 py-1"
              @change="updateTheme"
            />
          </div>
        </div>

        <!-- Secondary -->
        <div>
          <label class="block text-xs font-medium text-gray-700 mb-1">Secondary Color</label>
          <div class="flex items-center gap-2">
            <input 
              type="color" 
              v-model="theme.colors.secondary"
              class="w-8 h-8 rounded cursor-pointer border-0 p-0"
              @change="updateTheme"
            />
            <input 
              type="text" 
              v-model="theme.colors.secondary"
              class="flex-1 text-sm border border-gray-300 rounded px-2 py-1"
              @change="updateTheme"
            />
          </div>
        </div>

        <!-- Background -->
        <div>
          <label class="block text-xs font-medium text-gray-700 mb-1">Background</label>
          <div class="flex items-center gap-2">
            <input 
              type="color" 
              v-model="theme.colors.background"
              class="w-8 h-8 rounded cursor-pointer border-0 p-0"
              @change="updateTheme"
            />
             <input 
              type="text" 
              v-model="theme.colors.background"
              class="flex-1 text-sm border border-gray-300 rounded px-2 py-1"
              @change="updateTheme"
            />
          </div>
        </div>

        <!-- Text -->
        <div>
          <label class="block text-xs font-medium text-gray-700 mb-1">Text Color</label>
           <div class="flex items-center gap-2">
            <input 
              type="color" 
              v-model="theme.colors.textPrimary"
              class="w-8 h-8 rounded cursor-pointer border-0 p-0"
              @change="updateTheme"
            />
             <input 
              type="text" 
              v-model="theme.colors.textPrimary"
              class="flex-1 text-sm border border-gray-300 rounded px-2 py-1"
              @change="updateTheme"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Typography -->
    <div class="p-4 border-b border-gray-200">
      <div 
        @click="toggleGroup('typography')"
        class="flex items-center justify-between cursor-pointer mb-2"
      >
        <h4 class="font-medium text-gray-700">Typography</h4>
        <span class="text-xs text-indigo-600">{{ openGroups.typography ? 'Hide' : 'Show' }}</span>
      </div>

       <div v-if="openGroups.typography" class="space-y-4">
        <!-- Heading Font -->
        <div>
          <label class="block text-xs font-medium text-gray-700 mb-1">Heading Font</label>
          <select 
            v-model="theme.typography.headingFont"
            @change="updateTheme"
            class="w-full text-sm border border-gray-300 rounded px-2 py-1.5"
          >
            <option value="Inter">Inter</option>
            <option value="Roboto">Roboto</option>
            <option value="Open Sans">Open Sans</option>
            <option value="Poppins">Poppins</option>
            <option value="Playfair Display">Playfair Display (Serif)</option>
            <option value="Montserrat">Montserrat</option>
          </select>
        </div>

        <!-- Body Font -->
        <div>
          <label class="block text-xs font-medium text-gray-700 mb-1">Body Font</label>
          <select 
            v-model="theme.typography.bodyFont"
            @change="updateTheme"
            class="w-full text-sm border border-gray-300 rounded px-2 py-1.5"
          >
            <option value="Inter">Inter</option>
            <option value="Roboto">Roboto</option>
            <option value="Open Sans">Open Sans</option>
            <option value="Lato">Lato</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Layout -->
    <div class="p-4 border-b border-gray-200">
       <div 
        @click="toggleGroup('layout')"
        class="flex items-center justify-between cursor-pointer mb-2"
      >
        <h4 class="font-medium text-gray-700">Layout</h4>
         <span class="text-xs text-indigo-600">{{ openGroups.layout ? 'Hide' : 'Show' }}</span>
      </div>

      <div v-if="openGroups.layout" class="space-y-4">
        <!-- Container Width -->
        <div>
          <label class="block text-xs font-medium text-gray-700 mb-1">Container Width</label>
          <select 
            v-model="theme.layout.containerWidth" 
            @change="updateTheme"
            class="w-full text-sm border border-gray-300 rounded px-2 py-1.5"
          >
            <option value="boxed">Boxed (1200px)</option>
            <option value="wide">Wide (1400px)</option>
            <option value="full">Full Width</option>
          </select>
        </div>

        <!-- Border Radius -->
        <div>
          <label class="block text-xs font-medium text-gray-700 mb-1">Border Radius</label>
          <select 
            v-model="theme.layout.borderRadius" 
            @change="updateTheme"
            class="w-full text-sm border border-gray-300 rounded px-2 py-1.5"
          >
            <option value="0px">None (Sharp)</option>
            <option value="0.25rem">Small (4px)</option>
            <option value="0.5rem">Medium (8px)</option>
            <option value="1rem">Large (16px)</option>
            <option value="9999px">Pill (Rounded)</option>
          </select>
        </div>

         <!-- Button Style -->
        <div>
          <label class="block text-xs font-medium text-gray-700 mb-1">Button Style</label>
          <select 
            v-model="theme.layout.buttonStyle" 
            @change="updateTheme"
            class="w-full text-sm border border-gray-300 rounded px-2 py-1.5"
          >
            <option value="filled">Filled</option>
            <option value="outline">Outline</option>
            <option value="flat">Flat</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Reset -->
    <div class="p-4">
       <button 
        @click="resetToDefaults"
        class="w-full py-2 text-sm text-red-600 border border-red-200 rounded hover:bg-red-50"
      >
        Reset to Defaults
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import { useBuilderStore } from '@/store/builder'

const builderStore = useBuilderStore()

// Local state for UI toggles
const openGroups = reactive({
  colors: true,
  typography: false,
  layout: false
})

// Bind directly to store's theme object
// We assume builderStore.theme is populated from backend
const theme = ref(builderStore.theme || {
  colors: {
    primary: '#4F46E5',
    secondary: '#10B981',
    background: '#FFFFFF',
    textPrimary: '#1F2937'
  },
  typography: {
    headingFont: 'Inter',
    bodyFont: 'Inter'
  },
  layout: {
    containerWidth: 'boxed',
    borderRadius: '0.5rem',
    buttonStyle: 'filled'
  }
})

// Sync local ref if store updates (e.g. initial load)
watch(() => builderStore.theme, (newTheme) => {
  if (newTheme) {
    theme.value = newTheme
  }
}, { deep: true, immediate: true })

function toggleGroup(group) {
  openGroups[group] = !openGroups[group]
}

function updateTheme() {
  // Trigger store update to mark as dirty and update preview
  builderStore.updateTheme(theme.value)
}

function resetToDefaults() {
  if (confirm('Are you sure? This will reset all theme settings.')) {
     // TODO: Implement reset logic in store
  }
}
</script>
