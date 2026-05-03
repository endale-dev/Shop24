<template>
  <div class="min-h-full bg-gray-50 p-4 sm:p-6 lg:p-8">
    <!-- Header -->
    <div class="mb-6">
      <div class="sm:flex sm:items-center sm:justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">Products</h1>
          <p class="mt-1 text-sm text-gray-600">
            Manage your product catalog
          </p>
        </div>
        <div class="mt-4 sm:mt-0">
          <BaseButton @click="goToCreateProduct" icon-left="PlusIcon">
            Add Product
          </BaseButton>
        </div>
      </div>
    </div>

    <!-- Filters & Search -->
    <BaseCard class="mb-6">
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
        <BaseInput
          v-model="searchQuery"
          placeholder="Search products..."
          :icon-left="MagnifyingGlassIcon"
          @input="debouncedSearch"
        />
        
        <div>
          <select 
            v-model="stockFilter"
            class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
          >
            <option value="">All Stock Status</option>
            <option value="in_stock">In Stock</option>
            <option value="out_of_stock">Out of Stock</option>
          </select>
        </div>
        
        <div>
          <select 
            v-model="collectionFilter"
            class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
          >
            <option value="">All Collections</option>
            <option v-for="collection in collections" :key="collection.id" :value="collection.id">
              {{ collection.name }}
            </option>
          </select>
        </div>
      </div>
    </BaseCard>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center py-12">
      <LoadingSpinner size="lg" text="Loading products..." />
    </div>

    <!-- Error State -->
    <BaseAlert 
      v-else-if="error"
      variant="error"
      title="Failed to load products"
      :message="error.message"
      dismissible
      @dismiss="error = null"
    />

    <!-- Products Table -->
    <BaseCard v-else no-padding>
      <DataTable
        :columns="columns"
        :data="products"
        row-key="id"
        :empty-icon="CubeIcon"
        empty-title="No products found"
        empty-description="Get started by creating your first product"
      >
        <!-- Product column with image -->
        <template #cell-name="{ item }">
          <div class="flex items-center">
            <img 
              :src="getProductImage(item)" 
              :alt="item.name"
              class="h-10 w-10 rounded object-cover mr-3"
              @error="handleImageError"
            />
            <div>
              <div class="font-medium text-gray-900">{{ item.name }}</div>
              <div class="text-xs text-gray-500">{{ item.slug }}</div>
            </div>
          </div>
        </template>

        <!-- Price column -->
        <template #cell-price="{ item }">
          <span class="font-medium">${{ formatCurrency(item.price) }}</span>
        </template>

        <!-- Stock column -->
        <template #cell-stock="{ item }">
          <BaseBadge 
            :variant="item.stock > 0 ? 'success' : 'danger'"
            size="sm"
          >
            {{ item.stock }} in stock
          </BaseBadge>
        </template>

        <!-- Actions column -->
        <template #actions="{ item }">
          <div class="flex items-center justify-end space-x-2">
            <button
              @click="goToEditProduct(item.id)"
              class="text-blue-600 hover:text-blue-800 font-medium"
            >
              Edit
            </button>
            <button
              @click="confirmDelete(item)"
              class="text-red-600 hover:text-red-800 font-medium"
            >
              Delete
            </button>
          </div>
        </template>
      </DataTable>

      <!-- Pagination -->
      <div v-if="pagination.total > 0" class="border-t border-gray-200 px-6 py-4 bg-gray-50">
        <div class="flex items-center justify-between">
          <div class="text-sm text-gray-700">
            Showing {{ (pagination.page - 1) * pagination.limit + 1 }} to 
            {{ Math.min(pagination.page * pagination.limit, pagination.total) }} of 
            {{ pagination.total }} products
          </div>
          <div class="flex space-x-2">
            <BaseButton
              @click="prevPage"
              :disabled="pagination.page === 1"
              variant="secondary"
              size="sm"
            >
              Previous
            </BaseButton>
            <BaseButton
              @click="nextPage"
              :disabled="pagination.page >= Math.ceil(pagination.total / pagination.limit)"
              variant="secondary"
              size="sm"
            >
              Next
            </BaseButton>
          </div>
        </div>
      </div>
    </BaseCard>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useShopStore } from '@/store/shops'
import { productsAPI } from '@/services/api/endpoints'
import { getUserMessage } from '@/services/api/errors'
import { 
  BaseButton, 
  BaseCard, 
  BaseBadge, 
  BaseInput, 
  LoadingSpinner, 
  BaseAlert,
  DataTable
} from '@/components/base'
import { 
  PlusIcon, 
  MagnifyingGlassIcon, 
  CubeIcon 
} from '@heroicons/vue/24/outline'

const router = useRouter()
const shopStore = useShopStore()

// State
const loading = ref(true)
const error = ref(null)
const products = ref([])
const collections = ref([])
const searchQuery = ref('')
const stockFilter = ref('')
const collectionFilter = ref('')

const pagination = ref({
  page: 1,
  limit: 10,
  total: 0
})

// Computed
const shopId = computed(() => shopStore.active?.id)

// Table columns
const columns = [
  { key: 'name', label: 'Product', class: 'font-medium' },
  { key: 'price', label: 'Price' },
  { key: 'stock', label: 'Stock' },
  { key: 'createdAt', label: 'Created' }
]

// Methods
function formatCurrency(amount) {
  return (amount || 0).toFixed(2)
}

function getProductImage(product) {
  return product.main_image || product.mainImage || product.images?.[0] || '/placeholder-product.jpg'
}

function handleImageError(event) {
  event.target.src = '/placeholder-product.jpg'
}

function goToCreateProduct() {
  router.push({ name: 'CreateProduct' })
}

function goToEditProduct(productId) {
  router.push({ name: 'EditProduct', params: { productId } })
}

function confirmDelete(product) {
  if (confirm(`Are you sure you want to delete "${product.name}"?`)) {
    deleteProduct(product.id)
  }
}

async function deleteProduct(productId) {
  try {
    await productsAPI.delete(shopId.value, productId)
    // Refresh products list
    await loadProducts()
  } catch (err) {
    console.error('Delete error:', err)
    error.value = { message: getUserMessage(err) }
  }
}

async function loadProducts() {
  if (!shopId.value) return

  try {
    loading.value = true
    error.value = null

    const params = {
      page: pagination.value.page,
      limit: pagination.value.limit,
      search: searchQuery.value || undefined,
      stock_status: stockFilter.value || undefined,
      collection_id: collectionFilter.value || undefined
    }

    const data = await productsAPI.getAll(shopId.value, params)

    products.value = (data.products || []).map(p => ({
      id: p.id || p._id,
      name: p.name,
      slug: p.slug,
      price: p.price || p.starting_price || 0,
      stock: p.stock || p.total_stock || 0,
      main_image: p.main_image || p.mainImage,
      images: p.images || [],
      createdAt: p.createdAt || p.created_at
    }))

    pagination.value = {
      page: data.pagination?.page || 1,
      limit: data.pagination?.limit || 10,
      total: data.pagination?.total || products.value.length
    }

  } catch (err) {
    console.error('Load products error:', err)
    error.value = { message: getUserMessage(err) }
  } finally {
    loading.value = false
  }
}

function prevPage() {
  if (pagination.value.page > 1) {
    pagination.value.page--
    loadProducts()
  }
}

function nextPage() {
  const maxPage = Math.ceil(pagination.value.total / pagination.value.limit)
  if (pagination.value.page < maxPage) {
    pagination.value.page++
    loadProducts()
  }
}

// Debounced search
let searchTimeout
function debouncedSearch() {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    pagination.value.page = 1 // Reset to first page
    loadProducts()
  }, 300)
}

// Watch filters
watch([stockFilter, collectionFilter], () => {
  pagination.value.page = 1
  loadProducts()
})

// Watch shop changes
watch(() => shopStore.active, (newShop) => {
  if (newShop?.id) {
    loadProducts()
  }
})

// Load on mount
onMounted(() => {
  if (shopId.value) {
    loadProducts()
  }
})
</script>
