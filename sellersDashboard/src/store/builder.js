import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  fetchPages,
  fetchPage,
  createPage,
  updatePage,
  deletePage,
  publishPage,
  fetchActiveTheme,
  updateThemeColors,
  updateThemeTypography,
  updateThemeLayout
} from '@/services/pageBuilderService'

// ... (existing code)

  // ============ Theme Actions ============
  
  async function loadTheme() {
    if (!currentShopId.value) return
    try {
      const themeData = await fetchActiveTheme(currentShopId.value)
      if (themeData) {
        theme.value = themeData
      }
    } catch (e) {
      console.error("Failed to load theme", e)
      // Fallback to default if needed, or handle error
    }
  }

  // Update local state (for preview)
  function updateTheme(newTheme) {
    theme.value = newTheme
    isDirty.value = true
  }
  
  // Persist theme changes to backend
  async function saveTheme() {
    if (!currentShopId.value || !theme.value || !theme.value.id) return
    
    isSaving.value = true
    try {
      // Save specific sections
      // ideally we track what changed, but for now we can save all 3
      // or we can just call the specific ones based on what changed if we tracked it
      
      const p1 = updateThemeColors(currentShopId.value, theme.value.id, theme.value.colors)
      const p2 = updateThemeTypography(currentShopId.value, theme.value.id, theme.value.typography)
      const p3 = updateThemeLayout(currentShopId.value, theme.value.id, theme.value.layout)
      
      await Promise.all([p1, p2, p3])
      
      // reload to get fresh state? or just mark clean
      isDirty.value = false
    } catch (err) {
      error.value = err.message || 'Failed to save theme'
      console.error('Save theme error:', err)
    } finally {
      isSaving.value = false
    }
  }

export const useBuilderStore = defineStore('builder', () => {
  // ============ State ============
  
  // Current shop context
  const currentShopId = ref(null)
  
  // Page list and selection
  const pages = ref([])
  const currentPageId = ref(null)
  const currentPage = computed(() => 
    pages.value.find(p => p.id === currentPageId.value) || null
  )
  
  // Page sections (blocks)
  const sections = ref([])
  const selectedSectionId = ref(null)
  
  // Theme State
  const theme = ref(null)
  
  // Page metadata
  const pageTitle = ref('')
  const pagePath = ref('/')
  const pageSeo = ref({})
  const pageVisibility = ref('public')
  const pageStatus = ref('draft')
  
  // Theme settings (for live preview)
  const themeSettings = ref({
    colors: {
      primary: '#4F46E5',
      background: '#FFFFFF',
      text: '#1F2937'
    },
    fonts: {
      heading: 'Inter',
      body: 'Inter'
    }
  })
  
  // UI states
  const activeMode = ref('page') // 'page' | 'theme'
  const viewMode = ref('desktop') // 'desktop' | 'mobile'
  const isLoading = ref(false)
  const isSaving = ref(false)
  const isPublishing = ref(false)
  const isDirty = ref(false)
  const error = ref(null)
  
  // ============ Getters ============
  
  const selectedSection = computed(() => 
    sections.value.find(s => s.id === selectedSectionId.value) || null
  )
  
  const hasUnsavedChanges = computed(() => isDirty.value)
  
  // ============ Actions ============
  
  // Set current shop context
  function setShopId(shopId) {
    currentShopId.value = shopId
  }
  
  // Load all pages for current shop
  async function loadPages() {
    if (!currentShopId.value) return
    isLoading.value = true
    error.value = null
    try {
      pages.value = await fetchPages(currentShopId.value)
    } catch (err) {
      error.value = err.message || 'Failed to load pages'
      console.error('Load pages error:', err)
    } finally {
      isLoading.value = false
    }
  }
  
  // Load a specific page into the editor
  async function loadPage(pageId) {
    if (!currentShopId.value || !pageId) return
    isLoading.value = true
    error.value = null
    try {
      const page = await fetchPage(currentShopId.value, pageId)
      currentPageId.value = page.id
      pageTitle.value = page.title || ''
      pagePath.value = page.path || '/'
      pageSeo.value = page.seo || {}
      pageVisibility.value = page.visibility || 'public'
      pageStatus.value = page.status || 'draft'
      sections.value = page.blocks || []
      isDirty.value = false
    } catch (err) {
      error.value = err.message || 'Failed to load page'
      console.error('Load page error:', err)
    } finally {
      isLoading.value = false
    }
  }
  
  // Create a new page
  async function createNewPage(title = 'New Page', path = '/new-page') {
    if (!currentShopId.value) return null
    isSaving.value = true
    error.value = null
    try {
      const newPage = await createPage(currentShopId.value, {
        title,
        path,
        blocks: [],
        seo: {},
        visibility: 'public',
        status: 'draft'
      })
      await loadPages()
      await loadPage(newPage.id)
      return newPage
    } catch (err) {
      error.value = err.message || 'Failed to create page'
      console.error('Create page error:', err)
      return null
    } finally {
      isSaving.value = false
    }
  }
  
  // Save current page
  async function savePage() {
    if (!currentShopId.value || !currentPageId.value) {
      // If no page exists yet, create a new one
      if (!currentPageId.value && sections.value.length > 0) {
        return createNewPage(pageTitle.value || 'Home', pagePath.value || '/')
      }
      return
    }
    
    isSaving.value = true
    error.value = null
    try {
      await updatePage(currentShopId.value, currentPageId.value, {
        title: pageTitle.value,
        path: pagePath.value,
        blocks: sections.value,
        seo: pageSeo.value,
        visibility: pageVisibility.value,
        status: pageStatus.value
      })
      isDirty.value = false
      await loadPages()
    } catch (err) {
      error.value = err.message || 'Failed to save page'
      console.error('Save page error:', err)
    } finally {
      isSaving.value = false
    }
  }
  
  // Publish current page
  async function publishCurrentPage() {
    if (!currentShopId.value || !currentPageId.value) return
    
    isPublishing.value = true
    error.value = null
    try {
      // First save any pending changes
      await savePage()
      
      // Then publish
      await publishPage(currentShopId.value, currentPageId.value, {
        title: pageTitle.value,
        path: pagePath.value,
        blocks: sections.value,
        seo: pageSeo.value,
        visibility: pageVisibility.value,
        status: 'published'
      })
      
      pageStatus.value = 'published'
      await loadPages()
    } catch (err) {
      error.value = err.message || 'Failed to publish page'
      console.error('Publish page error:', err)
    } finally {
      isPublishing.value = false
    }
  }
  
  // Delete a page
  async function removePageFromList(pageId) {
    if (!currentShopId.value || !pageId) return
    try {
      await deletePage(currentShopId.value, pageId)
      if (currentPageId.value === pageId) {
        resetEditor()
      }
      await loadPages()
    } catch (err) {
      error.value = err.message || 'Failed to delete page'
      console.error('Delete page error:', err)
    }
  }
  
  // ============ Section Actions ============
  
  function addSection(section) {
    sections.value.push(section)
    isDirty.value = true
  }

  function removeSection(sectionId) {
    const index = sections.value.findIndex(s => s.id === sectionId)
    if (index !== -1) {
      sections.value.splice(index, 1)
      isDirty.value = true
    }
    if (selectedSectionId.value === sectionId) {
      selectedSectionId.value = null
    }
  }

  function updateSection(sectionId, newProps) {
    const section = sections.value.find(s => s.id === sectionId)
    if (section) {
      section.props = { ...section.props, ...newProps }
      isDirty.value = true
    }
  }

  function reorderSections(newOrder) {
    sections.value = newOrder
    isDirty.value = true
  }

  function selectSection(id) {
    selectedSectionId.value = id
  }
  
  // ============ Theme Actions ============
  
  async function loadTheme() {
    if (!currentShopId.value) return
    try {
      const themeData = await fetchActiveTheme(currentShopId.value)
      // Ensure local theme state matches backend structure
      if (themeData) {
        theme.value = themeData
      }
    } catch (e) {
      console.error("Failed to load theme", e)
    }
  }

  function updateTheme(newTheme) {
    theme.value = newTheme
    isDirty.value = true
  }

  async function saveTheme() {
    if (!currentShopId.value || !theme.value || !theme.value.id) return
    
    isSaving.value = true
    try {
       // Save all sections
      await Promise.all([
        updateThemeColors(currentShopId.value, theme.value.id, theme.value.colors),
        updateThemeTypography(currentShopId.value, theme.value.id, theme.value.typography),
        updateThemeLayout(currentShopId.value, theme.value.id, theme.value.layout)
      ])
      
      isDirty.value = false
    } catch (err) {
      error.value = "Failed to save theme"
      console.error(err)
    } finally {
      isSaving.value = false
    }
  }
  
  // ============ Page Metadata Actions ============
  
  function updatePageTitle(title) {
    pageTitle.value = title
    isDirty.value = true
  }
  
  function updatePagePath(path) {
    pagePath.value = path
    isDirty.value = true
  }
  
  function updatePageSeo(seo) {
    pageSeo.value = { ...pageSeo.value, ...seo }
    isDirty.value = true
  }
  
  // ============ Reset ============
  
  // ============ Reset ============
  
  function resetEditor() {
    currentPageId.value = null
    sections.value = []
    selectedSectionId.value = null
    pageTitle.value = ''
    pagePath.value = '/'
    pageSeo.value = {}
    pageVisibility.value = 'public'
    pageStatus.value = 'draft'
    isDirty.value = false
    error.value = null
  }

  function setViewMode(mode) {
    viewMode.value = mode
  }

  return {
    // State
    currentShopId,
    pages,
    currentPageId,
    currentPage,
    sections,
    selectedSectionId,
    selectedSection,
    pageTitle,
    pagePath,
    pageSeo,
    pageVisibility,
    pageStatus,
    themeSettings,
    activeMode,
    viewMode,
    isLoading,
    isSaving,
    isPublishing,
    isDirty,
    hasUnsavedChanges,
    error,
    
    // Actions
    setShopId,
    loadPages,
    loadPage,
    createNewPage,
    savePage,
    publishCurrentPage,
    removePageFromList,
    addSection,
    removeSection,
    updateSection,
    reorderSections,
    selectSection,
    theme,
    loadTheme,
    updateTheme,
    saveTheme,
    updatePageTitle,
    updatePagePath,
    updatePageSeo,
    resetEditor,
    setViewMode
  }
}, {
  persist: {
    paths: ['currentShopId', 'currentPageId', 'sections', 'pageTitle', 'pagePath']
  }
})

