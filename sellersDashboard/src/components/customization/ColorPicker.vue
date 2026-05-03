<!-- sellersDashboard/src/components/customization/ColorPicker.vue -->
<template>
  <div class="color-picker-container">
    <label v-if="label" class="color-picker-label">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </label>
    
    <div class="color-picker-input-group">
      <!-- Color swatch -->
      <div class="color-swatch" :style="{ backgroundColor: modelValue }">
        <input
          type="color"
          :value="modelValue"
          @input="handleColorChange"
          class="color-input-hidden"
        />
      </div>
      
      <!-- Hex input -->
      <input
        type="text"
        :value="modelValue"
        @input="handleTextChange"
        @blur="validateColor"
        placeholder="#000000"
        class="color-hex-input"
        maxlength="7"
      />
      
      <!-- Copy button -->
      <button
        type="button"
        @click="copyColor"
        class="color-copy-btn"
        title="Copy color"
      >
        <svg
          v-if="!copied"
          xmlns="http://www.w3.org/2000/svg"
          class="h-4 w-4"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
          />
        </svg>
        <svg
          v-else
          xmlns="http://www.w3.org/2000/svg"
          class="h-4 w-4 text-green-600"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M5 13l4 4L19 7"
          />
        </svg>
      </button>
    </div>
    
    <!-- Description -->
    <p v-if="description" class="color-picker-description">
      {{ description }}
    </p>
    
    <!-- Error message -->
    <p v-if="error" class="color-picker-error">
      {{ error }}
    </p>
  </div>
</template>

<script>
export default {
  name: 'ColorPicker',
  props: {
    modelValue: {
      type: String,
      default: '#000000'
    },
    label: {
      type: String,
      default: ''
    },
    description: {
      type: String,
      default: ''
    },
    required: {
      type: Boolean,
      default: false
    }
  },
  emits: ['update:modelValue'],
  data() {
    return {
      copied: false,
      error: ''
    }
  },
  methods: {
    handleColorChange(event) {
      const color = event.target.value
      this.error = ''
      this.$emit('update:modelValue', color)
    },
    handleTextChange(event) {
      const color = event.target.value
      this.error = ''
      this.$emit('update:modelValue', color)
    },
    validateColor() {
      const hexRegex = /^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/
      if (!hexRegex.test(this.modelValue)) {
        this.error = 'Invalid hex color format. Use #RRGGBB'
        // Reset to a valid color
        this.$emit('update:modelValue', '#000000')
      } else {
        this.error = ''
      }
    },
    async copyColor() {
      try {
        await navigator.clipboard.writeText(this.modelValue)
        this.copied = true
        setTimeout(() => {
          this.copied = false
        }, 2000)
      } catch (err) {
        console.error('Failed to copy color:', err)
      }
    }
  }
}
</script>

<style scoped>
.color-picker-container {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.color-picker-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #374151;
}

.color-picker-input-group {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.color-swatch {
  position: relative;
  width: 3rem;
  height: 2.5rem;
  border-radius: 0.375rem;
  border: 2px solid #E5E7EB;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
}

.color-swatch:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.color-input-hidden {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  cursor: pointer;
}

.color-hex-input {
  flex: 1;
  padding: 0.625rem 0.75rem;
  border: 1px solid #D1D5DB;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  font-family: 'Courier New', monospace;
  text-transform: uppercase;
  transition: all 0.2s ease;
}

.color-hex-input:focus {
  outline: none;
  border-color: #3B82F6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.color-copy-btn {
  padding: 0.625rem;
  border: 1px solid #D1D5DB;
  border-radius: 0.375rem;
  background-color: white;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.color-copy-btn:hover {
  background-color: #F9FAFB;
  border-color: #9CA3AF;
}

.color-picker-description {
  font-size: 0.75rem;
  color: #6B7280;
  margin: 0;
}

.color-picker-error {
  font-size: 0.75rem;
  color: #EF4444;
  margin: 0;
}
</style>
