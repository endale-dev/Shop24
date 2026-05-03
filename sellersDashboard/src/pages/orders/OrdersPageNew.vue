<template>
  <div class="min-h-full bg-gray-50 p-4 sm:p-6 lg:p-8">
    <!-- Header -->
    <div class="mb-6">
      <div class="sm:flex sm:items-center sm:justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">Orders</h1>
          <p class="mt-1 text-sm text-gray-600">
            Track and manage customer orders
          </p>
        </div>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 gap-5 sm:grid-cols-4 mb-6">
      <StatCard
        title="Total Orders"
        :value="stats.total"
        :icon="ShoppingBagIcon"
        icon-color="blue"
      />
      <StatCard
        title="Pending"
        :value="stats.pending"
        subtitle="Awaiting action"
        :icon="ClockIcon"
        icon-color="amber"
      />
      <StatCard
        title="Completed"
        :value="stats.completed"
        subtitle="Delivered"
        :icon="CheckCircleIcon"
        icon-color="green"
      />
      <StatCard
        title="Revenue"
        :value="`$${formatCurrency(stats.revenue)}`"
        :icon="CurrencyDollarIcon"
        icon-color="green"
      />
    </div>

    <!-- Filters -->
    <BaseCard class="mb-6">
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <BaseInput
          v-model="searchQuery"
          placeholder="Search by order number..."
          :icon-left="MagnifyingGlassIcon"
          @input="debouncedSearch"
        />
        
        <div>
          <select 
            v-model="statusFilter"
            class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            @change="loadOrders"
          >
            <option value="">All Status</option>
            <option value="pending">Pending</option>
            <option value="paid">Paid</option>
            <option value="shipped">Shipped</option>
            <option value="delivered">Delivered</option>
            <option value="cancelled">Cancelled</option>
          </select>
        </div>
      </div>
    </BaseCard>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center py-12">
      <LoadingSpinner size="lg" text="Loading orders..." />
    </div>

    <!-- Error State -->
    <BaseAlert 
      v-else-if="error"
      variant="error"
      title="Failed to load orders"
      :message="error.message"
      dismissible
      @dismiss="error = null"
    />

    <!-- Orders Table -->
    <BaseCard v-else no-padding>
      <DataTable
        :columns="columns"
        :data="orders"
        row-key="id"
        :empty-icon="ShoppingBagIcon"
        empty-title="No orders found"
        empty-description="Orders will appear here once customers place them"
      >
        <!-- Order Number column -->
        <template #cell-orderNumber="{ item }">
          <div class="font-medium text-blue-600 hover:text-blue-800 cursor-pointer" @click="goToOrderDetail(item.id)">
            #{{ item.orderNumber }}
          </div>
        </template>

        <!-- Customer column -->
        <template #cell-customer="{ item }">
          <div class="text-sm text-gray-900">{{ item.customerName || 'N/A' }}</div>
        </template>

        <!-- Total column -->
        <template #cell-total="{ item }">
          <span class="font-medium">${{ formatCurrency(item.total) }}</span>
        </template>

        <!-- Status column -->
        <template #cell-status="{ item }">
          <BaseBadge :variant="getStatusVariant(item.status)" size="sm">
            {{ formatStatus(item.status) }}
          </BaseBadge>
        </template>

        <!-- Date column -->
        <template #cell-createdAt="{ item }">
          <div class="text-sm text-gray-600">{{ formatDateTime(item.createdAt) }}</div>
        </template>

        <!-- Actions column -->
        <template #actions="{ item }">
          <div class="flex items-center justify-end space-x-3">
            <button
              @click="goToOrderDetail(item.id)"
              class="text-blue-600 hover:text-blue-800 font-medium text-sm"
            >
              View
            </button>
            <button
              v-if="item.status === 'pending'"
              @click="updateOrderStatus(item, 'paid')"
              class="text-green-600 hover:text-green-800 font-medium text-sm"
            >
              Mark Paid
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
            {{ pagination.total }} orders
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
import { ordersAPI } from '@/services/api/endpoints'
import { getUserMessage } from '@/services/api/errors'
import { format } from 'date-fns'
import { 
  BaseButton, 
  BaseCard, 
  BaseBadge, 
  BaseInput, 
  LoadingSpinner, 
  BaseAlert,
  DataTable
} from '@/components/base'
import StatCard from '@/components/common/StatCard.vue'
import { 
  ShoppingBagIcon, 
  MagnifyingGlassIcon,
  ClockIcon,
  CheckCircleIcon,
  CurrencyDollarIcon
} from '@heroicons/vue/24/outline'

const router = useRouter()
const shopStore = useShopStore()

// State
const loading = ref(true)
const error = ref(null)
const orders = ref([])
const searchQuery = ref('')
const statusFilter = ref('')

const stats = ref({
  total: 0,
  pending: 0,
  completed: 0,
  revenue: 0
})

const pagination = ref({
  page: 1,
  limit: 10,
  total: 0
})

// Computed
const shopId = computed(() => shopStore.active?.id)

// Table columns
const columns = [
  { key: 'orderNumber', label: 'Order #' },
  { key: 'customer', label: 'Customer' },
  { key: 'total', label: 'Total' },
  { key: 'status', label: 'Status' },
  { key: 'createdAt', label: 'Date' }
]

// Methods
function formatCurrency(amount) {
  return (amount || 0).toFixed(2)
}

function formatDateTime(dateString) {
  if (!dateString) return 'N/A'
  try {
    return format(new Date(dateString), 'MMM d, yyyy HH:mm')
  } catch (e) {
    return 'Invalid Date'
  }
}

function formatStatus(status) {
  const statusMap = {
    'pending': 'Pending',
    'paid': 'Paid',
    'shipped': 'Shipped',
    'delivered': 'Delivered',
    'cancelled': 'Cancelled'
  }
  return statusMap[status] || status
}

function getStatusVariant(status) {
  const variantMap = {
    'pending': 'warning',
    'paid': 'info',
    'shipped': 'primary',
    'delivered': 'success',
    'cancelled': 'danger'
  }
  return variantMap[status] || 'neutral'
}

function goToOrderDetail(orderId) {
  router.push({ name: 'OrderDetail', params: { orderId } })
}

async function updateOrderStatus(order, newStatus) {
  try {
    await ordersAPI.update(shopId.value, order.id, { status: newStatus })
    // Refresh orders
    await loadOrders()
  } catch (err) {
    console.error('Update order error:', err)
    error.value = { message: getUserMessage(err) }
  }
}

async function loadOrders() {
  if (!shopId.value) return

  try {
    loading.value = true
    error.value = null

    const params = {
      page: pagination.value.page,
      limit: pagination.value.limit,
      status: statusFilter.value || undefined
    }

    const data = await ordersAPI.getAll(shopId.value, params)

    orders.value = (data.orders || []).map(o => ({
      id: o.id || o._id,
      orderNumber: o.order_number || o.orderNumber,
      customerName: o.customer_name || o.customerName,
      total: o.total || 0,
      status: o.status,
      createdAt: o.created_at || o.createdAt
    }))

    pagination.value = {
      page: data.pagination?.page || 1,
      limit: data.pagination?.limit || 10,
      total: data.pagination?.total || orders.value.length
    }

    // Load stats
    await loadStats()

  } catch (err) {
    console.error('Load orders error:', err)
    error.value = { message: getUserMessage(err) }
  } finally {
    loading.value = false
  }
}

async function loadStats() {
  try {
    const data = await ordersAPI.getStats(shopId.value)
    
    stats.value = {
      total: data.total || 0,
      pending: data.pending || 0,
      completed: data.delivered || 0,
      revenue: data.total_revenue || 0
    }
  } catch (err) {
    console.error('Load stats error:', err)
  }
}

function prevPage() {
  if (pagination.value.page > 1) {
    pagination.value.page--
    loadOrders()
  }
}

function nextPage() {
  const maxPage = Math.ceil(pagination.value.total / pagination.value.limit)
  if (pagination.value.page < maxPage) {
    pagination.value.page++
    loadOrders()
  }
}

// Debounced search
let searchTimeout
function debouncedSearch() {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    pagination.value.page = 1
    loadOrders()
  }, 300)
}

// Watch shop changes
watch(() => shopStore.active, (newShop) => {
  if (newShop?.id) {
    loadOrders()
  }
})

// Load on mount
onMounted(() => {
  if (shopId.value) {
    loadOrders()
  }
})
</script>
