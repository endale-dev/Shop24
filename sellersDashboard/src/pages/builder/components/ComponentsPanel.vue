<template>
  <div class="flex flex-col h-full">
    <!-- Active Sections (The "Page Structure") -->
    <div class="flex-1 flex flex-col min-h-0 border-b border-gray-200">
      <div class="p-4 border-b border-gray-100 bg-gray-50 flex justify-between items-center">
        <h3 class="font-medium text-gray-900">Page Sections</h3>
        <span class="text-xs text-gray-400">Drag to reorder</span>
      </div>
      <div class="flex-1 overflow-y-auto p-4 bg-white">
        <draggable
          class="space-y-2 min-h-[100px]"
          v-model="builderStore.sections"
          group="sections"
          item-key="id"
          handle=".drag-handle"
          @change="handleSectionChange"
        >
          <template #item="{ element }">
            <div 
              class="group p-3 border rounded-lg bg-white shadow-sm hover:border-indigo-500 transition-colors cursor-pointer flex items-center justify-between"
              :class="{'border-indigo-600 ring-1 ring-indigo-600': builderStore.selectedSectionId === element.id, 'border-gray-200': builderStore.selectedSectionId !== element.id}"
              @click="builderStore.selectSection(element.id)"
            >
              <div class="flex items-center gap-3">
                <div class="drag-handle cursor-move text-gray-400 hover:text-gray-600">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M3 5a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 15a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
                  </svg>
                </div>
                <div class="flex flex-col">
                    <span class="text-sm font-medium text-gray-900">{{ element.label || element.type }}</span>
                    <span class="text-xs text-gray-400 capitalize">{{ element.type }}</span>
                </div>
              </div>
              <button @click.stop="builderStore.removeSection(element.id)" class="text-gray-400 hover:text-red-500 opacity-0 group-hover:opacity-100 transition-opacity">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 000-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
                </svg>
              </button>
            </div>
          </template>
        </draggable>
        
        <div v-if="builderStore.sections.length === 0" class="text-center py-8 text-gray-400 text-sm border-2 border-dashed border-gray-200 rounded-lg">
            Drop sections here
        </div>
      </div>
    </div>

    <!-- Available Sections (Source) -->
    <div class="h-1/3 flex flex-col border-t border-gray-200 shadow-[0_-4px_6px_-1px_rgba(0,0,0,0.05)]">
      <div class="p-3 bg-gray-50 border-b border-gray-200">
         <h3 class="font-medium text-gray-900 text-sm">Add Section</h3>
      </div>
      <div class="flex-1 overflow-y-auto p-4 bg-gray-50/50">
        <draggable 
            class="grid grid-cols-2 gap-2" 
            :list="availableSections" 
            :group="{ name: 'sections', pull: 'clone', put: false }" 
            :clone="cloneSection" 
            item-key="type"
            :sort="false"
        >
            <template #item="{ element }">
            <div class="p-2 bg-white border border-gray-200 rounded-md shadow-sm cursor-move hover:border-indigo-500 hover:shadow-md transition-all flex flex-col items-center text-center gap-2">
                <component :is="getIcon(element.icon)" class="w-6 h-6 text-gray-500" />
                <span class="text-xs font-medium text-gray-700">{{ element.label }}</span>
            </div>
            </template>
        </draggable>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import draggable from 'vuedraggable';
import { useBuilderStore } from '@/store/builder';
import { 
  PhotoIcon, 
  ShoppingBagIcon, 
  NewspaperIcon, 
  ChatBubbleBottomCenterTextIcon 
} from '@heroicons/vue/24/outline';

const builderStore = useBuilderStore();

const availableSections = ref([
  { type: 'HeroBanner', label: 'Hero Banner', description: 'Large promotional banner', icon: 'PhotoIcon' },
  { type: 'ProductGrid', label: 'Product Grid', description: 'Display a collection of products', icon: 'ShoppingBagIcon' },
  { type: 'CollectionGrid', label: 'Collections', description: 'Shop by collection', icon: 'ShoppingBagIcon' },
  { type: 'Features', label: 'Features', description: 'Store specific features', icon: 'NewspaperIcon' },
  { type: 'Newsletter', label: 'Newsletter', description: 'Email signup form', icon: 'ChatBubbleBottomCenterTextIcon' },
  { type: 'RichText', label: 'Rich Text', description: 'Text block with heading', icon: 'NewspaperIcon' },
]);

const getIcon = (name) => {
  const icons = { PhotoIcon, ShoppingBagIcon, NewspaperIcon, ChatBubbleBottomCenterTextIcon };
  return icons[name] || PhotoIcon;
};

const cloneSection = (original) => {
  return {
    ...original,
    id: crypto.randomUUID(),
    props: {} // Initialize with default props if needed
  }
};

const handleSectionChange = () => {
    // Trigger any necessary updates or history saving
};
</script>
