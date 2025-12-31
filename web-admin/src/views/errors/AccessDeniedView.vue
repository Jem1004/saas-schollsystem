<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { Button, Typography, Card } from 'ant-design-vue'
import { StopOutlined, HomeOutlined, ArrowLeftOutlined } from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { getRoleDisplayName } from '@/composables/usePermissions'

const { Title, Text, Paragraph } = Typography

const router = useRouter()
const authStore = useAuthStore()

const userRole = computed(() => {
  return authStore.userRole ? getRoleDisplayName(authStore.userRole) : 'User'
})

const goBack = () => {
  router.back()
}

const goHome = () => {
  router.push('/dashboard')
}

const goLogin = () => {
  authStore.clearAuth()
  router.push('/login')
}
</script>

<template>
  <div class="error-container">
    <Card class="error-card">
      <div class="error-content">
        <div class="error-icon-wrapper">
          <StopOutlined class="error-icon" />
        </div>
        <Title :level="1" class="error-title">403</Title>
        <Title :level="3" class="error-subtitle">Akses Ditolak</Title>
        <Paragraph class="error-message">
          Maaf, Anda tidak memiliki izin untuk mengakses halaman ini.
        </Paragraph>
        <div v-if="authStore.isAuthenticated" class="role-info">
          <Text type="secondary">
            Anda login sebagai: <Text strong>{{ userRole }}</Text>
          </Text>
        </div>
        <Paragraph type="secondary" class="error-hint">
          Jika Anda merasa ini adalah kesalahan, silakan hubungi administrator sistem.
        </Paragraph>
        <div class="error-actions">
          <Button type="primary" size="large" @click="goHome">
            <template #icon><HomeOutlined /></template>
            Ke Dashboard
          </Button>
          <Button size="large" @click="goBack">
            <template #icon><ArrowLeftOutlined /></template>
            Kembali
          </Button>
        </div>
        <div v-if="!authStore.isAuthenticated" class="login-link">
          <Button type="link" @click="goLogin">
            Login dengan akun lain
          </Button>
        </div>
      </div>
    </Card>
  </div>
</template>

<style scoped>
.error-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f5f5f5 0%, #e8e8e8 100%);
  padding: 20px;
}

.error-card {
  max-width: 500px;
  width: 100%;
  text-align: center;
  border-radius: 12px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.1);
}

.error-content {
  padding: 24px;
}

.error-icon-wrapper {
  width: 80px;
  height: 80px;
  margin: 0 auto 24px;
  background: #fff2f0;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.error-icon {
  font-size: 40px;
  color: #ff4d4f;
}

.error-title {
  font-size: 64px !important;
  font-weight: 700 !important;
  color: #1a1a1a !important;
  margin: 0 0 8px 0 !important;
  line-height: 1 !important;
}

.error-subtitle {
  font-weight: 600 !important;
  color: #333333 !important;
  margin: 0 0 16px 0 !important;
}

.error-message {
  font-size: 16px;
  color: #666666;
  margin: 0 0 8px 0;
}

.role-info {
  margin-bottom: 16px;
  padding: 12px;
  background: #fafafa;
  border-radius: 8px;
}

.error-hint {
  font-size: 14px;
  margin: 0 0 24px 0;
}

.error-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
  flex-wrap: wrap;
}

.login-link {
  margin-top: 16px;
}
</style>
