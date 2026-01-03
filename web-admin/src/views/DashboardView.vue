<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import SuperAdminDashboard from '@/views/super-admin/SuperAdminDashboard.vue'
import AdminSekolahDashboard from '@/views/admin-sekolah/AdminSekolahDashboard.vue'
import GuruBKDashboard from '@/views/guru-bk/GuruBKDashboard.vue'
import WaliKelasDashboard from '@/views/wali-kelas/WaliKelasDashboard.vue'

const authStore = useAuthStore()

const userRole = computed(() => authStore.userRole)

// Redirect to role-specific dashboard or show appropriate component
onMounted(() => {
  // For now, we show the Super Admin dashboard for super_admin role
  // Other roles will get their dashboards in later tasks
})
</script>

<template>
  <div class="dashboard">
    <!-- Super Admin Dashboard -->
    <SuperAdminDashboard v-if="userRole === 'super_admin'" />
    
    <!-- Admin Sekolah Dashboard -->
    <AdminSekolahDashboard v-else-if="userRole === 'admin_sekolah'" />
    
    <!-- Guru BK Dashboard -->
    <GuruBKDashboard v-else-if="userRole === 'guru_bk'" />
    
    <!-- Wali Kelas Dashboard -->
    <WaliKelasDashboard v-else-if="userRole === 'wali_kelas'" />
    
    <!-- Placeholder for other roles (will be implemented in later tasks) -->
    <div v-else class="placeholder-dashboard">
      <h1>Dashboard</h1>
      <p>Dashboard untuk role {{ userRole }} akan diimplementasikan pada task selanjutnya.</p>
    </div>
  </div>
</template>

<style scoped>
.dashboard {
  padding: 0;
}

.placeholder-dashboard {
  padding: 24px;
  text-align: center;
}

.placeholder-dashboard h1 {
  color: #262626;
  margin-bottom: 16px;
}

.placeholder-dashboard p {
  color: #8c8c8c;
}
</style>
