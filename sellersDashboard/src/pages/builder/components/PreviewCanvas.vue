<template>
  <div class="flex-1 flex flex-col h-full bg-gray-100">
    <!-- Canvas Area -->
    <div class="flex-1 overflow-hidden flex items-start justify-center p-8 active-area relative">
      <div 
        :class="['bg-white shadow-2xl transition-all duration-300 relative overflow-hidden', builderStore.viewMode === 'mobile' ? 'w-[375px] h-[812px] rounded-3xl border-8 border-gray-900' : 'w-full h-full max-w-[1200px] border border-gray-200']"
      >
        <!-- The Iframe Storefront -->
        <iframe 
          ref="iframeRef"
          :src="storefrontUrl"
          class="w-full h-full border-0"
          title="Store Preview"
          @load="onIframeLoad"
        ></iframe>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { useBuilderStore } from '@/store/builder';

const iframeRef = ref(null);
const builderStore = useBuilderStore();

// Point to the running storefront dev server
const storefrontUrl = computed(() => {
    // In production this would be the actual shop URL. In dev we point to localhost.
    return 'http://localhost:5174/?mode=editor'; 
});

// Watch for state changes and send to iframe
watch(
  () => builderStore.sections,
  (newSections) => {
    if(iframeRef.value?.contentWindow) {
        // We send a clone to avoid proxy issues across postMessage
        const sectionsPayload = JSON.parse(JSON.stringify(newSections));
        iframeRef.value.contentWindow.postMessage({ type: 'UPDATE_SECTIONS', payload: sectionsPayload }, '*');
    }
  },
  { deep: true }
);

// Also send initial state when iframe loads
const onIframeLoad = () => {
    if(builderStore.sections.length > 0 && iframeRef.value?.contentWindow) {
         const sectionsPayload = JSON.parse(JSON.stringify(builderStore.sections));
         iframeRef.value.contentWindow.postMessage({ type: 'UPDATE_SECTIONS', payload: sectionsPayload }, '*');
    }
}
// Watch for theme changes
watch(() => builderStore.theme, (newTheme) => {
  if (!newTheme || !iframeRef.value?.contentWindow) return
  
  // Direct postMessage since sendMessage helper isn't defined
  iframeRef.value.contentWindow.postMessage({
    type: 'preview:theme',
    payload: {
      theme: JSON.parse(JSON.stringify(newTheme))
    }
  }, '*')
}, { deep: true })
</script>

<style scoped>
.active-area {
  perspective: 1000px;
}
</style>
