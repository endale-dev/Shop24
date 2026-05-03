<template>
  <div class="min-h-full bg-gray-50 p-4 sm:p-6 lg:p-8">
    <!-- Header Section -->
    <div class="mb-8">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">Dashboard</h1>
          <p class="mt-1 text-sm text-gray-600">
            Welcome back! Here's an overview of your shop's performance.
          </p>
        </div>
        <BaseButton @click="refreshData" :loading="loading" icon-left="RefreshIcon" variant="ghost">
          Refresh
        </BaseButton>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center py-20">
      <LoadingSpinner size="lg" text="Loading dashboard data..." />
    </div>

    <!-- Error State -->
   <BaseAlert 
      v-else-if="error" 
      variant="error" 
      title="Failed to load dashboard"
      :message="error.message"
      dismissible
      @dismiss="error = null"
    />

    <!-- Dashboard Content -->
    <div v-else class="space-y-6">
      
      <!-- Stats Grid -->
      <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard
          title="Total Revenue"
          :value="`$${formatCurrency(stats.totalRevenue)}`"
          :change="stats.revenueGrowth"
          subtitle="+${{formatCurrency(stats.revenueToday)}} today"
          :icon="CurrencyDollarIcon"
          icon-color="green"
        />
        
        <StatCard
          title="Total Orders"
          :value="stats.orders"
          :change="stats.ordersGrowth"
          subtitle="+{{stats.newOrdersToday}} new today"
          :icon="ShoppingBagIcon"
          icon-color="blue"
        />
        
        <StatCard
          title="Total Customers"
          :value="stats.customers"
          :change="stats.customersGrowth"
          subtitle="+{{stats.newCustomersToday}} new today"
          :icon="UsersIcon"
          icon-color="purple"
        />
        
        <StatCard
          title="Total Products"
          :value="stats.products"
          :subtitle="`${stats.lowStockProducts} low stock`"
          :icon="CubeIcon"
          icon-color="amber"
        />
      </div>

      <!-- Main Grid -->
      <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
        
        <!-- Revenue Chart (2 columns) -->
        <BaseCard title="Revenue Overview" subtitle="Last 7 days" class="lg:col-span-2">
          <div class="h-64 flex items-end justify-center">
            <EmptyState 
              v-if="!weeklySales || weeklySales.length === 0"
              title="No data available"
              description="Sales data will appear here once you have orders"
            />
            <div v-else class="w-full">
              <div class="grid grid-cols-7 gap-2 h-56">
                <div 
                  v-for="(day, index) in weeklySales" 
                  :key="index"
                  class="text-center flex flex-col items-center justify-end relative group"
                >
                  <div class="absolute -top-8 bg-gray-900 text-white text-xs px-2 py-1 rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap">
                    ${{ formatCurrency(day.amount) }}
                  </div>
                  <div 
                    class="bg-gradient-to-t from-blue-500 to-blue-600 w-full rounded-t hover:from-blue-600 hover:to-blue-700 cursor-pointer transition-colors"
                    :style="{ height: `${Math.max(day.amount / maxWeeklySales * 200, 4)}px` }"
                  ></div>
                  <div class="text-xs text-gray-600 font-medium mt-2">{{ day.day }}</div>
                </div>
              </div>
            </div>
          </div>
        </BaseCard>

        <!-- Quick Stats Panel -->
        <BaseCard title="Quick Stats">
          <template #default>
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-sm text-gray-600">Avg Order Value</p>
                  <p class="text-2xl font-semibold text-gray-900">${{ formatCurrency(stats.averageOrderValue) }}</p>
                </div>
                <div class="p-3 bg-blue-100 rounded-lg">
                  <ShoppingCartIcon class="h-6 w-6 text-blue-600" />
                </div>
              </div>
              
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-sm text-gray-600">Pending Orders</p>
                  <p class="text-2xl font-semibold text-gray-900">{{ stats.pendingOrders }}</p>
                </div>
                <BaseBadge variant="warning" rounded>Action Needed</BaseBadge>
              </div>
              
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-sm text-gray-600">Top Product</p>
                  <p class="text-sm font-medium text-gray-900 truncate">{{ stats.topProduct?.name || 'N/A' }}</p>
                </div>
                <StarIcon class="h-5 w-5 text-amber-500" />
              </div>
            </div>
          </template>
        </BaseCard>
      </div>

      <!-- Recent Activity -->
      <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
        
        <!-- Recent Orders -->
        <BaseCard title="Recent Orders">
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-gray-900">Recent Orders</h3>
              <router-link to="/dashboard/orders" class="text-sm text-blue-600 hover:text-blue-700 font-medium">
                View all →
              </router-link>
            </div>
          </template>
          
          <EmptyState 
            v-if="!recentOrders || recentOrders.length === 0"
            :icon="ShoppingBagIcon"
            title="No orders yet"
            description="New orders will appear here once placed"
          />
          
          <div v-else class="divide-y divide-gray-100">
            <div 
              v-for="order in recentOrders.slice(0, 5)" 
              :key="order.id"
              class="flex items-center justify-between py-3 hover:bg-gray-50 cursor-pointer rounded px-2 -mx-2 transition-colors"
              @click="goToOrder(order.id)"
            >
              <div class="flex items-center space-x-3">
                <div class="p-2 bg-blue-50 rounded">
                  <ShoppingBagIcon class="h-5 w-5 text-blue-600" />
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-900">#{{ order.orderNumber }}</p>
                  <p class="text-xs text-gray-500">{{ formatDateTime(order.createdAt) }}</p>
                </div>
              </div>
              <div class="text-right">
                <p class="text-sm font-medium text-gray-900">${{ formatCurrency(order.total) }}</p>
                <BaseBadge :variant="getStatusVariant(order.status)" size="sm">
                  {{ formatStatus(order.status) }}
                </BaseBadge>
              </div>
            </div>
          </div>
        </BaseCard>

        <!-- Low Stock Alerts -->
        <BaseCard title="Inventory Alerts">
          <EmptyState 
            v-if="lowStockItems.length === 0"
            :icon="CheckCircleIcon"
            title="All stocked up!"
            description="No inventory alerts at this time"
            variant="success"
          />
          
          <div v-else class="space-y-3">
            <BaseAlert 
              v-for="item in lowStockItems.slice(0, 5)" 
              :key="item.id"
              variant="warning"
              :title="item.name"
              :message="`Only ${item.stock} items remaining in stock`"
            />
          </div>
        </BaseCard>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useShopStore } from '@/store/shops'
import { analyticsAPI } from '@/services/api/endpoints'
import { getUserMessage } from '@/services/api/errors'
import { format } from 'date-fns'

// Components
import { BaseButton, BaseCard, BaseBadge, LoadingSpinner, EmptyState, BaseAlert } from '@/components/base'
import StatCard from '@/components/common/StatCard.vue'

// Icons
import {
  CurrencyDollarIcon,
  ShoppingBagIcon,
  UsersIcon,
  CubeIcon,
  ShoppingCartIcon,
  StarIcon,
  CheckCircleIcon
} from '@heroicons/vue/24/outline'

const router = useRouter()
const shopStore = useShopStore()

// Computed
const activeShop = computed(() => shopStore.active)
const shopId = computed(() => activeShop.value?.id)

// State
const loading = ref(true)
const error = ref(null)

const stats = ref({
  totalRevenue: 0,
  revenueToday: 0,
  revenueGrowth: 0,
  orders: 0,
  newOrdersToday: 0,
  ordersGrowth: 0,
  customers: 0,
  newCustomersToday: 0,
  customersGrowth: 0,
  products: 0,
  lowStockProducts: 0,
  pendingOrders: 0,
  averageOrderValue: 0,
  topProduct: null
})

const weeklySales = ref([])
const recentOrders = ref([])
const lowStockItems = ref([])

// Computed
const maxWeeklySales = computed(() => {
  if (!weeklySales.value || weeklySales.value.length === 0) return 1
  return Math.max(...weeklySales.value.map(day => day.amount))
})

// Methods
function formatCurrency(amount) {
  return (amount || 0).toFixed(2)
}

function formatDateTime(dateString) {
  if (!dateString) return 'N/A'
  try {
    return format(new Date(dateString), 'MMM d, HH:mm')
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

function goToOrder(orderId) {
  router.push({ name: 'OrderDetail', params: { orderId } })
}

async function loadDashboardData() {
  if (!shopId.value) {
    loading.value = false
    return
  }

  try {
    loading.value = true
    error.value = null

    // Fetch dashboard analytics
    const data = await analyticsAPI.getDashboard(shopId.value)

    // Map data to state
    stats.value = {
      totalRevenue: data.total_revenue || 0,
      revenueToday: data.revenue_today || 0,
      revenueGrowth: data.revenue_growth || 0,
      orders: data.total_orders || 0,
      newOrdersToday: data.new_orders_today || 0,
      ordersGrowth: data.orders_growth || 0,
      customers: data.total_customers || 0,
      newCustomersToday: data.new_customers_today || 0,
      customersGrowth: data.customers_growth || 0,
      products: data.total_products || 0,
      lowStockProducts: data.low_stock_count || 0,
      pendingOrders: data.pending_orders || 0,
      averageOrderValue: data.average_order_value || 0,
      topProduct: data.top_product || null
    }

    weeklySales.value = (data.weekly_sales || []).map(day => ({
      day: day.day,
      amount: day.amount || 0
    }))

    recentOrders.value = (data.recent_orders) || []).map(order => ({
      id: order.id || order._id,
      orderNumber: order.order_number || order.orderNumber,
      total: order.total || 0,
      status: order.status,
      createdAt: order.created_at || order.createdAt
    }))

    lowStockItems.value = (data.low_stock_products || []).map(product => ({
      id: product.id || product._id,
      name: product.name,
      stock: product.stock || 0
    }))

  } catch (err) {
    console.error('Dashboard load error:', err)
    error.value = {
      message: getUserMessage(err)
    }
  } finally {
    loading.value = false
  }
}

function refreshData() {
  loadDashboardData()
}

// Watch for shop changes
watch(() => shopStore.active, (newShop) => {
  if (newShop?.id) {
    loadDashboardData()
  }
}, { immediate: false })

// Load data on mount
onMounted(() => {
  if (shopId.value) {
    loadDashboardData()
  }
})
</script>
