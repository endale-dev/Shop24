<script setup lang="ts">
import LandingPage from '@/pages/LandingPage.vue'
import { useAuthStore } from '@/stores/auth'
import DashboardLayout from './DashboardLayout.vue'

const auth = useAuthStore()
</script>

<template>
  <DashboardLayout v-if="auth.user" />
  <LandingPage v-else />
</template>
