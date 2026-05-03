<template>
  <div class="h-screen flex flex-col overflow-hidden bg-gray-50">
      <!-- Top Bar -->
      <div class="h-16 bg-white border-b border-gray-200 flex items-center justify-between px-4 z-20 relative">
        <div class="flex items-center gap-4">
          <router-link :to="`/shop/${route.params.shopId}/dashboard`" class="text-gray-500 hover:text-gray-900">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
          </router-link>
          
          <div class="h-6 w-px bg-gray-300 mx-2"></div>
          
          <h1 class="font-bold text-lg text-gray-800">
            Store Designer
            <span v-if="builderStore.isDirty" class="ml-2 text-xs bg-yellow-100 text-yellow-800 px-2 py-0.5 rounded-full font-medium">Unsaved</span>
          </h1>

        <div class="flex items-center gap-4 ml-8">
          <!-- Mode Switcher -->
          <div class="flex bg-gray-100 p-1 rounded-lg">
            <button 
              @click="activeMode = 'page'"
              class="px-4 py-1.5 rounded-md text-sm font-medium transition-all"
              :class="activeMode === 'page' ? 'bg-white shadow text-gray-900' : 'text-gray-500 hover:text-gray-900'"
            >
              Page Editor
            </button>
            <button 
              @click="activeMode = 'theme'"
              class="px-4 py-1.5 rounded-md text-sm font-medium transition-all"
              :class="activeMode === 'theme' ? 'bg-white shadow text-gray-900' : 'text-gray-500 hover:text-gray-900'"
            >
              Theme Settings
            </button>
          </div>

          <!-- View Mode Toggles -->
          <div class="flex bg-gray-100 p-1 rounded-lg">
             <button 
              @click="builderStore.setViewMode('desktop')"
              class="p-1.5 rounded-md transition-all"
              :class="builderStore.viewMode === 'desktop' ? 'bg-white shadow text-indigo-600' : 'text-gray-500 hover:text-gray-900'"
              title="Desktop View"
            >
               <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
              </svg>
            </button>
            <button 
              @click="builderStore.setViewMode('mobile')"
              class="p-1.5 rounded-md transition-all"
              :class="builderStore.viewMode === 'mobile' ? 'bg-white shadow text-indigo-600' : 'text-gray-500 hover:text-gray-900'"
              title="Mobile View"
            >
               <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
              </svg>
            </button>
          </div>
        </div>
        </div>

        <div v-if="activeMode === 'page'" class="flex items-center gap-2">
            <!-- Page Selector -->
            <select 
              v-model="selectedPageId" 
              @change="onPageChange"
              class="border border-gray-300 rounded-md px-3 py-1.5 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option disabled value="">Select Page</option>
              <option v-for="page in builderStore.pages" :key="page.id" :value="page.id">
                {{ page.title || page.path }} {{ page.status === 'published' ? '✓' : '(draft)' }}
              </option>
              <option value="new">+ Create New Page</option>
            </select>
        </div>

        <div class="flex items-center gap-3">
          <button 
            @click="saveChanges" 
            :disabled="builderStore.isSaving || !builderStore.isDirty"
            class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50"
          >
            {{ builderStore.isSaving ? 'Saving...' : 'Save Changes' }}
          </button>
          <button 
            @click="publishChanges" 
            :disabled="builderStore.isPublishing || (!builderStore.currentPageId && activeMode === 'page')"
            class="px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 disabled:opacity-50 shadow-sm"
          >
            {{ builderStore.isPublishing ? 'Publishing...' : 'Publish' }}
          </button>
        </div>
      </div>

    <!-- Error Alert -->
    <div v-if="builderStore.error" class="bg-red-50 border-b border-red-200 px-4 py-2 flex items-center justify-between">
      <span class="text-sm text-red-700">{{ builderStore.error }}</span>
      <button @click="builderStore.error = null" class="text-red-500 hover:text-red-700">×</button>
    </div>

    <!-- Loading State -->
    <div v-if="builderStore.isLoading" class="flex-1 flex items-center justify-center">
      <div class="flex flex-col items-center gap-3">
        <div class="w-8 h-8 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
        <span class="text-gray-500">Loading page...</span>
      </div>
    </div>

    <!-- Main Workspace -->
    <div v-else class="flex-1 flex overflow-hidden">
      <!-- Left Sidebar: Components -->
      <aside class="w-72 bg-white border-r border-gray-200 flex flex-col z-10">
        <ComponentsPanel />
      </aside>

      <!-- Center: Canvas -->
      <main class="flex-1 relative bg-gray-100 overflow-hidden flex flex-col">
        <PreviewCanvas />
      </main>

      <!-- Right Sidebar: Properties -->
      <aside class="w-80 bg-white border-l border-gray-200 flex flex-col z-10">
        <PropertiesPanel />
        <PropertiesPanel v-if="activeMode === 'page'" />
        <ThemeSettingsPanel v-else-if="activeMode === 'theme'" />
      </aside>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useBuilderStore } from '@/store/builder';
import ComponentsPanel from './components/ComponentsPanel.vue';
import PreviewCanvas from './components/PreviewCanvas.vue';
import PropertiesPanel from './components/PropertiesPanel.vue';
import ThemeSettingsPanel from './components/ThemeSettingsPanel.vue';

const route = useRoute();
const router = useRouter();
const builderStore = useBuilderStore();
const shopId = route.params.shopId;
const pageId = route.params.pageId;

const previewMode = ref(false);
const iframeRef = ref(null);
const selectedPageId = ref('');
const activeMode = ref('page'); // 'page' | 'theme'

// Initialize shop context from route
onMounted(async () => {
  if (shopId) {
    builderStore.setShopId(shopId);
    await builderStore.loadPages();
    
    // Load page from route or first available
    if (pageId) {
      selectedPageId.value = pageId;
      await builderStore.loadPage(pageId);
    } else if (builderStore.pages.length > 0) {
      // Auto-select home page or first page
      const homePage = builderStore.pages.find(p => p.path === '/');
      const firstPage = homePage || builderStore.pages[0];
      if (firstPage) {
        selectedPageId.value = firstPage.id;
        await builderStore.loadPage(firstPage.id);
      }
    }
  }
});

// Watch for page selection changes
function onPageChange() {
  if (selectedPageId.value) {
    builderStore.loadPage(selectedPageId.value);
    router.replace({ query: { ...route.query, pageId: selectedPageId.value } });
  }
}

// Create new page
async function createPage() {
  const title = prompt('Enter page title:', 'New Page');
  if (!title) return;
  
  const path = prompt('Enter page path (e.g., /about):', '/new-page');
  if (!path) return;
  
  const newPage = await builderStore.createNewPage(title, path);
  if (newPage) {
    selectedPageId.value = newPage.id;
  }
}

// Save changes based on active mode
async function saveChanges() {
  if (activeMode.value === 'page') {
     await builderStore.savePage();
  } else if (activeMode.value === 'theme') {
     await builderStore.saveTheme();
  }
}

// Publish current page
async function publishChanges() {
  if (activeMode.value === 'page') {
      if (confirm('Are you sure you want to publish this page?')) {
        await builderStore.publishCurrentPage();
      }
  } else {
      // For themes, save is essentially publish for now as there's no draft theme state yet
      if (confirm('Are you sure you want to publish these theme settings?')) {
        await builderStore.saveTheme();
      }
  }
}

// Open preview in new tab
function openPreview() {
  const shopSlug = localStorage.getItem('selectedShopSlug') || 'demo';
  const previewUrl = `http://${shopSlug}.localhost:5175${builderStore.pagePath}`;
  window.open(previewUrl, '_blank');
}

// Go back
function goBack() {
  if (builderStore.isDirty) {
    if (!confirm('You have unsaved changes. Are you sure you want to leave?')) {
      return;
    }
  }
  router.push('/dashboard/customization');
}

// Warn on page unload
window.addEventListener('beforeunload', (e) => {
  if (builderStore.isDirty) {
    e.preventDefault();
    e.returnValue = '';
  }
});

// Sync Theme to Preview
watch(() => builderStore.theme, (newTheme) => {
  if (newTheme) {
    // We need to access the PreviewCanvas iframe
    // Since PreviewCanvas is a child, we can rely on it exposing a method or we can use a global event bus
    // Or we can let PreviewCanvas handle the store subscription itself.
    // However, PreviewCanvas likely has the iframe ref.
    // Let's rely on PreviewCanvas watching the store, which is cleaner.
    // BUT, we need to ensure the store is actually updated.
    // builderStore.updateTheme() already updates the state.
  }
}, { deep: true });
</script>

