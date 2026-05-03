<!-- sellersDashboard/src/pages/customization/ThemeColorsPage.vue -->
<template>
  <div class="theme-colors-page">
    <PageHeader title="Theme Colors" subtitle="Customize your storefront color palette" />

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center items-center min-h-96">
      <LoadingSpinner size="lg" text="Loading theme..." />
    </div>

    <!-- Error State -->
    <BaseAlert v-else-if="error" variant="error" :title="error" dismissible @dismiss="error = null" />

    <!-- Content -->
    <div v-else class="colors-content">
      <!-- Color Grid -->
      <BaseCard title="Color Palette" description="Define your brand colors">
        <div class="colors-grid">
          <ColorPicker
            v-model="colors.primary"
            label="Primary Color"
            description="Main brand color used for key actions"
            required
          />
          <ColorPicker
            v-model="colors.secondary"
            label="Secondary Color"
            description="Supporting brand color for accents"
            required
          />
          <ColorPicker
            v-model="colors.accent"
            label="Accent Color"
            description="Call-to-action and highlight color"
            required
          />
          <ColorPicker
            v-model="colors.background"
            label="Background"
            description="Main page background color"
            required
          />
          <ColorPicker
            v-model="colors.surface"
            label="Surface"
            description="Card and container background color"
            required
          />
          <ColorPicker
            v-model="colors.textPrimary"
            label="Primary Text"
            description="Main text color for headings and body"
            required
          />
          <ColorPicker
            v-model="colors.textSecondary"
            label="Secondary Text"
            description="Supporting text and descriptions"
            required
          />
          <ColorPicker
            v-model="colors.border"
            label="Border"
            description="Borders and dividers"
            required
          />
          <ColorPicker
            v-model="colors.success"
            label="Success"
            description="Success messages and states"
            required
          />
          <ColorPicker
            v-model="colors.warning"
            label="Warning"
            description="Warning messages and alerts"
            required
          />
          <ColorPicker
            v-model="colors.error"
            label="Error"
            description="Error messages and states"
            required
          />
          <ColorPicker
            v-model="colors.info"
            label="Info"
            description="Informational messages"
            required
          />
        </div>

        <!-- Actions -->
        <div class="card-actions">
          <BaseButton variant="secondary" @click="resetToDefaults">
            Reset to Defaults
          </BaseButton>
          <BaseButton variant="primary" :loading="saving" @click="saveColors">
            Save Changes
          </BaseButton>
        </div>
      </BaseCard>

      <!-- Preview Section -->
      <BaseCard title="Preview" description="See how your colors look">
        <div class="preview-section">
          <div class="preview-banner" :style="{ backgroundColor: colors.primary }">
            <h3 :style="{ color: colors.background }">Your Store Name</h3>
            <p :style="{ color: colors.background }">Welcome to your storefront</p>
          </div>
          
          <div class="preview-card" :style="{ backgroundColor: colors.surface, borderColor: colors.border }">
            <h4 :style="{ color: colors.textPrimary }">Product Card</h4>
            <p :style="{ color: colors.textSecondary }">This is how your product cards will look</p>
            <button class="preview-button" :style="{ backgroundColor: colors.accent, color: colors.background }">
              Add to Cart
            </button>
          </div>

          <div class="preview-states">
            <div class="state-badge" :style="{ backgroundColor: colors.success }">Success</div>
            <div class="state-badge" :style="{ backgroundColor: colors.warning }">Warning</div>
            <div class="state-badge" :style="{ backgroundColor: colors.error }">Error</div>
            <div class="state-badge" :style="{ backgroundColor: colors.info }">Info</div>
          </div>
        </div>
      </BaseCard>
    </div>

    <!-- Success Toast -->
    <BaseAlert
      v-if="showSuccess"
      variant="success"
      title="Colors saved successfully!"
      class="success-toast"
      dismissible
      @dismiss="showSuccess = false"
    />
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import PageHeader from '@/components/common/PageHeader.vue'
import BaseCard from '@/components/base/BaseCard.vue'
import BaseButton from '@/components/base/BaseButton.vue'
import BaseAlert from '@/components/base/BaseAlert.vue'
import LoadingSpinner from '@/components/base/LoadingSpinner.vue'
import ColorPicker from '@/components/customization/ColorPicker.vue'
import { themesAPI, defaultColors } from '@/services/api/endpoints/themes'

export default {
  name: 'ThemeColorsPage',
  components: {
    PageHeader,
    BaseCard,
    BaseButton,
    BaseAlert,
    LoadingSpinner,
    ColorPicker
  },
  setup() {
    const route = useRoute()
    const shopId = route.params.shopId

    const loading = ref(true)
    const saving = ref(false)
    const error = ref(null)
    const showSuccess = ref(false)
    const themeId = ref(null)

    const colors = reactive({ ...defaultColors })

    const loadTheme = async () => {
      try {
        loading.value = true
        error.value = null
        const response = await themesAPI.getActiveTheme(shopId)
        
        if (response.theme && response.theme.colors) {
          Object.assign(colors, response.theme.colors)
          themeId.value = response.theme.id
        }
      } catch (err) {
        error.value = err.message || 'Failed to load theme'
        console.error('Error loading theme:', err)
      } finally {
        loading.value = false
      }
    }

    const saveColors = async () => {
      try {
        saving.value = true
        error.value = null
        
        if (!themeId.value) {
          throw new Error('No active theme found. Please create a theme first.')
        }

        await themesAPI.updateColors(shopId, themeId.value, colors)
        showSuccess.value = true
        
        setTimeout(() => {
          showSuccess.value = false
        }, 3000)
      } catch (err) {
        error.value = err.message || 'Failed to save colors'
        console.error('Error saving colors:', err)
      } finally {
        saving.value = false
      }
    }

    const resetToDefaults = () => {
      Object.assign(colors, defaultColors)
    }

    onMounted(() => {
      loadTheme()
    })

    return {
      loading,
      saving,
      error,
      showSuccess,
      colors,
      saveColors,
      resetToDefaults
    }
  }
}
</script>

<style scoped>
.theme-colors-page {
  padding: 1.5rem;
  max-width: 1400px;
  margin: 0 auto;
}

.colors-content {
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 1.5rem;
  margin-top: 1.5rem;
}

.colors-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}

.card-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding-top: 1.5rem;
  border-top: 1px solid #E5E7EB;
}

.preview-section {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.preview-banner {
  padding: 2rem;
  border-radius: 0.5rem;
  text-align: center;
}

.preview-banner h3 {
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0 0 0.5rem 0;
}

.preview-banner p {
  margin: 0;
}

.preview-card {
  padding: 1.5rem;
  border-radius: 0.5rem;
  border: 2px solid;
}

.preview-card h4 {
  font-size: 1.125rem;
  font-weight: 600;
  margin: 0 0 0.5rem 0;
}

.preview-card p {
  margin: 0 0 1rem 0;
}

.preview-button {
  padding: 0.5rem 1rem;
  border-radius: 0.375rem;
  border: none;
  font-weight: 500;
  cursor: pointer;
}

.preview-states {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.state-badge {
  padding: 0.5rem 1rem;
  border-radius: 0.375rem;
  color: white;
  font-size: 0.875rem;
  font-weight: 500;
}

.success-toast {
  position: fixed;
  bottom: 2rem;
  right: 2rem;
  z-index: 1000;
  max-width: 400px;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
}

@media (max-width: 1024px) {
  .colors-content {
    grid-template-columns: 1fr;
  }
  
  .colors-grid {
    grid-template-columns: 1fr;
  }
}
</style>
