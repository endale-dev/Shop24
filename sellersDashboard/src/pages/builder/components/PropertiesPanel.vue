<template>
  <div class="flex flex-col h-full">
    <!-- Header -->
    <div class="p-4 border-b border-gray-200 bg-white">
      <h3 class="font-medium text-gray-900">Properties</h3>
      <p v-if="selectedSection" class="text-xs text-gray-500 mt-1">
        Editing: {{ selectedSection.label || selectedSection.type }}
      </p>
    </div>
    
    <!-- No Selection State -->
    <div v-if="!selectedSection" class="flex-1 p-6 bg-gray-50 flex flex-col items-center justify-center text-center">
      <div class="w-12 h-12 bg-gray-200 rounded-lg flex items-center justify-center mb-3">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122" />
        </svg>
      </div>
      <p class="text-gray-500 text-sm">Select a section to edit its properties</p>
    </div>
    
    <!-- Properties Form -->
    <div v-else class="flex-1 overflow-y-auto">
      <!-- Page Settings Tab -->
      <div class="border-b border-gray-200">
        <button 
          @click="activeTab = 'section'" 
          :class="['px-4 py-2 text-sm font-medium', activeTab === 'section' ? 'text-indigo-600 border-b-2 border-indigo-600' : 'text-gray-500']"
        >
          Section
        </button>
        <button 
          @click="activeTab = 'page'" 
          :class="['px-4 py-2 text-sm font-medium', activeTab === 'page' ? 'text-indigo-600 border-b-2 border-indigo-600' : 'text-gray-500']"
        >
          Page
        </button>
      </div>
      
      <!-- Section Properties -->
      <div v-if="activeTab === 'section'" class="p-4 space-y-4">
        <!-- Dynamic fields based on component type -->
        <div v-for="field in componentSchema" :key="field.key" class="space-y-1">
          <div class="flex items-center justify-between">
            <label class="block text-sm font-medium text-gray-700">{{ field.label }}</label>
            <span v-if="builderStore.viewMode === 'mobile'" class="text-xs text-indigo-600 font-medium">Mobile</span>
          </div>
          
          <!-- Text Input -->
          <input 
            v-if="field.type === 'text'"
            type="text"
            :value="getFieldValue(field.key)"
            @input="updateField(field.key, $event.target.value)"
            :placeholder="field.placeholder || ''"
            class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
          
          <!-- Textarea -->
          <textarea 
            v-else-if="field.type === 'textarea'"
            :value="getFieldValue(field.key)"
            @input="updateField(field.key, $event.target.value)"
            :placeholder="field.placeholder || ''"
            rows="3"
            class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
          
          <!-- Color Picker -->
          <div v-else-if="field.type === 'color'" class="flex items-center gap-2">
            <input 
              type="color"
              :value="getFieldValue(field.key) || '#000000'"
              @input="updateField(field.key, $event.target.value)"
              class="w-10 h-10 rounded border border-gray-300 cursor-pointer"
            />
            <input 
              type="text"
              :value="getFieldValue(field.key)"
              @input="updateField(field.key, $event.target.value)"
              class="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm"
              placeholder="#000000"
            />
          </div>
          
          <!-- Image URL -->
          <div v-else-if="field.type === 'image'" class="space-y-2">
            <input 
              type="text"
              :value="getFieldValue(field.key)"
              @input="updateField(field.key, $event.target.value)"
              placeholder="Enter image URL"
              class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
            <div v-if="getFieldValue(field.key)" class="relative aspect-video bg-gray-100 rounded-md overflow-hidden">
              <img :src="getFieldValue(field.key)" alt="Preview" class="w-full h-full object-cover" />
            </div>
          </div>
          
          <!-- Checkbox -->
          <label v-else-if="field.type === 'checkbox'" class="flex items-center gap-2 cursor-pointer">
            <input 
              type="checkbox"
              :checked="getFieldValue(field.key)"
              @change="updateField(field.key, $event.target.checked)"
              class="w-4 h-4 text-indigo-600 rounded"
            />
            <span class="text-sm text-gray-600">{{ field.checkboxLabel || 'Enable' }}</span>
          </label>
          
          <!-- Select -->
          <select 
            v-else-if="field.type === 'select'"
            :value="getFieldValue(field.key)"
            @change="updateField(field.key, $event.target.value)"
            class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
          >
            <option v-for="opt in field.options" :key="opt.value" :value="opt.value">
              {{ opt.label }}
            </option>
          </select>
          
          <!-- Number -->
          <input 
            v-else-if="field.type === 'number'"
            type="number"
            :value="getFieldValue(field.key)"
            @input="updateField(field.key, Number($event.target.value))"
            :min="field.min"
            :max="field.max"
            class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
          
          <p v-if="field.hint" class="text-xs text-gray-400">{{ field.hint }}</p>
        </div>
        
        <!-- Delete Section Button -->
        <div class="pt-4 border-t border-gray-200">
          <button 
            @click="deleteSection"
            class="w-full px-3 py-2 text-sm font-medium text-red-600 bg-red-50 rounded-md hover:bg-red-100 transition-colors"
          >
            Delete Section
          </button>
        </div>
      </div>
      
      <!-- Page Properties -->
      <div v-if="activeTab === 'page'" class="p-4 space-y-4">
        <div class="space-y-1">
          <label class="block text-sm font-medium text-gray-700">Page Title</label>
          <input 
            type="text"
            :value="builderStore.pageTitle"
            @input="builderStore.updatePageTitle($event.target.value)"
            class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="Home"
          />
        </div>
        
        <div class="space-y-1">
          <label class="block text-sm font-medium text-gray-700">Page Path</label>
          <input 
            type="text"
            :value="builderStore.pagePath"
            @input="builderStore.updatePagePath($event.target.value)"
            class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="/"
          />
          <p class="text-xs text-gray-400">URL path where this page will be accessible</p>
        </div>
        
        <div class="space-y-1">
          <label class="block text-sm font-medium text-gray-700">SEO Title</label>
          <input 
            type="text"
            :value="builderStore.pageSeo?.title"
            @input="builderStore.updatePageSeo({ title: $event.target.value })"
            class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="Page title for search engines"
          />
        </div>
        
        <div class="space-y-1">
          <label class="block text-sm font-medium text-gray-700">SEO Description</label>
          <textarea 
            :value="builderStore.pageSeo?.description"
            @input="builderStore.updatePageSeo({ description: $event.target.value })"
            rows="3"
            class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="Brief description for search engines"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useBuilderStore } from '@/store/builder';

const builderStore = useBuilderStore();
const activeTab = ref('section');

// Component property schemas
const componentSchemas = {
  HeroBanner: [
    { key: 'title', label: 'Title', type: 'text', placeholder: 'Welcome to our store' },
    { key: 'subtitle', label: 'Subtitle', type: 'textarea', placeholder: 'Discover amazing products' },
    { key: 'buttonText', label: 'Button Text', type: 'text', placeholder: 'Shop Now' },
    { key: 'buttonLink', label: 'Button Link', type: 'text', placeholder: '/products' },
    { key: 'backgroundImage', label: 'Background Image', type: 'image' },
    { key: 'backgroundColor', label: 'Background Color', type: 'color' },
    { key: 'textColor', label: 'Text Color', type: 'color' },
    { key: 'alignment', label: 'Alignment', type: 'select', options: [
      { value: 'left', label: 'Left' },
      { value: 'center', label: 'Center' },
      { value: 'right', label: 'Right' }
    ]},
    { key: 'height', label: 'Height', type: 'select', options: [
      { value: 'small', label: 'Small (300px)' },
      { value: 'medium', label: 'Medium (500px)' },
      { value: 'large', label: 'Large (700px)' },
      { value: 'full', label: 'Full Screen' }
    ]}
  ],
  ProductGrid: [
    { key: 'title', label: 'Section Title', type: 'text', placeholder: 'Featured Products' },
    { key: 'subtitle', label: 'Subtitle', type: 'text', placeholder: 'Check out our latest items' },
    { key: 'collectionId', label: 'Collection ID', type: 'text', hint: 'Leave empty to show all products' },
    { key: 'columns', label: 'Columns', type: 'select', options: [
      { value: '2', label: '2 Columns' },
      { value: '3', label: '3 Columns' },
      { value: '4', label: '4 Columns' }
    ]},
    { key: 'limit', label: 'Max Products', type: 'number', min: 1, max: 20 },
    { key: 'showPrice', label: 'Show Price', type: 'checkbox', checkboxLabel: 'Display product prices' },
    { key: 'showRating', label: 'Show Rating', type: 'checkbox', checkboxLabel: 'Display product ratings' }
  ],
  CollectionGrid: [
    { key: 'title', label: 'Section Title', type: 'text', placeholder: 'Shop by Category' },
    { key: 'columns', label: 'Columns', type: 'select', options: [
      { value: '2', label: '2 Columns' },
      { value: '3', label: '3 Columns' },
      { value: '4', label: '4 Columns' }
    ]},
    { key: 'showDescription', label: 'Show Description', type: 'checkbox', checkboxLabel: 'Display collection descriptions' }
  ],
  Features: [
    { key: 'title', label: 'Section Title', type: 'text', placeholder: 'Why Choose Us' },
    { key: 'columns', label: 'Columns', type: 'select', options: [
      { value: '2', label: '2 Columns' },
      { value: '3', label: '3 Columns' },
      { value: '4', label: '4 Columns' }
    ]},
    { key: 'feature1Title', label: 'Feature 1 Title', type: 'text' },
    { key: 'feature1Description', label: 'Feature 1 Description', type: 'textarea' },
    { key: 'feature1Icon', label: 'Feature 1 Icon', type: 'text', hint: 'Icon name or emoji' },
    { key: 'feature2Title', label: 'Feature 2 Title', type: 'text' },
    { key: 'feature2Description', label: 'Feature 2 Description', type: 'textarea' },
    { key: 'feature3Title', label: 'Feature 3 Title', type: 'text' },
    { key: 'feature3Description', label: 'Feature 3 Description', type: 'textarea' }
  ],
  Newsletter: [
    { key: 'title', label: 'Title', type: 'text', placeholder: 'Subscribe to our newsletter' },
    { key: 'subtitle', label: 'Subtitle', type: 'textarea', placeholder: 'Get the latest updates and offers' },
    { key: 'buttonText', label: 'Button Text', type: 'text', placeholder: 'Subscribe' },
    { key: 'backgroundColor', label: 'Background Color', type: 'color' },
    { key: 'textColor', label: 'Text Color', type: 'color' }
  ],
  RichText: [
    { key: 'heading', label: 'Heading', type: 'text', placeholder: 'About Us' },
    { key: 'content', label: 'Content', type: 'textarea', placeholder: 'Enter your content here...' },
    { key: 'alignment', label: 'Alignment', type: 'select', options: [
      { value: 'left', label: 'Left' },
      { value: 'center', label: 'Center' },
      { value: 'right', label: 'Right' }
    ]},
    { key: 'backgroundColor', label: 'Background Color', type: 'color' }
  ]
};

// Computed
const selectedSection = computed(() => builderStore.selectedSection);

const componentSchema = computed(() => {
  if (!selectedSection.value) return [];
  return componentSchemas[selectedSection.value.type] || [];
});

// Methods
function getFieldValue(key) {
  const mode = builderStore.viewMode || 'desktop';
  const val = selectedSection.value?.props?.[key];
  
  // If val is an object (responsive), return correct mode's value
  if (val && typeof val === 'object' && ('desktop' in val || 'mobile' in val)) {
    return val[mode] ?? val['desktop'] ?? '';
  }
  
  // If val is primitive but we are in mobile mode, fallback to primitive (desktop)
  return val ?? '';
}

function updateField(key, value) {
  if (!selectedSection.value) return;
  
  const mode = builderStore.viewMode || 'desktop';
  const currentVal = selectedSection.value.props?.[key];
  
  let newVal;
  
  // Check if we need to upgrade to responsive object
  if (mode === 'mobile') {
    if (typeof currentVal === 'object' && currentVal !== null) {
       newVal = { ...currentVal, mobile: value };
    } else {
       // Upgrade primitive to object, keeping primitive as desktop default
       newVal = { desktop: currentVal, mobile: value };
    }
  } else {
    // Desktop mode
    if (typeof currentVal === 'object' && currentVal !== null) {
      newVal = { ...currentVal, desktop: value };
    } else {
      newVal = value; // Keep as primitive if never switched to mobile
    }
  }
  
  builderStore.updateSection(selectedSection.value.id, { [key]: newVal });
}

function deleteSection() {
  if (!selectedSection.value) return;
  if (confirm('Are you sure you want to delete this section?')) {
    builderStore.removeSection(selectedSection.value.id);
  }
}
</script>

