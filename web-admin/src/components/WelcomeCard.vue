<script setup lang="ts">
import { computed } from 'vue'
import { Card, Typography, Button, Space, Tag, Steps, Step } from 'ant-design-vue'
import {
  SmileOutlined,
  RocketOutlined,
  CheckCircleOutlined,
  BookOutlined,
  TeamOutlined,
  SettingOutlined,
} from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { getRoleDisplayName } from '@/composables/usePermissions'

const { Title, Text } = Typography

interface Props {
  showSteps?: boolean
  dismissible?: boolean
}

withDefaults(defineProps<Props>(), {
  showSteps: true,
  dismissible: true,
})

const emit = defineEmits<{
  (e: 'dismiss'): void
  (e: 'action', action: string): void
}>()

const authStore = useAuthStore()

const userName = computed(() => authStore.user?.username || 'User')
const userRole = computed(() => authStore.userRole ? getRoleDisplayName(authStore.userRole) : 'User')

const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 12) return 'Selamat Pagi'
  if (hour < 15) return 'Selamat Siang'
  if (hour < 18) return 'Selamat Sore'
  return 'Selamat Malam'
})

// Role-specific quick actions
const quickActions = computed(() => {
  const role = authStore.userRole
  
  const actions: Record<string, { icon: typeof BookOutlined; label: string; action: string }[]> = {
    super_admin: [
      { icon: TeamOutlined, label: 'Kelola Sekolah', action: 'tenants' },
      { icon: SettingOutlined, label: 'Kelola Device', action: 'devices' },
    ],
    admin_sekolah: [
      { icon: TeamOutlined, label: 'Tambah Siswa', action: 'students' },
      { icon: BookOutlined, label: 'Kelola Kelas', action: 'classes' },
      { icon: SettingOutlined, label: 'Pengaturan', action: 'settings' },
    ],
    guru_bk: [
      { icon: TeamOutlined, label: 'Lihat Siswa', action: 'students' },
      { icon: BookOutlined, label: 'Catat Prestasi', action: 'achievements' },
    ],
    wali_kelas: [
      { icon: TeamOutlined, label: 'Absensi Kelas', action: 'attendance' },
      { icon: BookOutlined, label: 'Input Nilai', action: 'grades' },
    ],
  }
  
  return role ? actions[role] || [] : []
})

const handleAction = (action: string) => {
  emit('action', action)
}

const handleDismiss = () => {
  emit('dismiss')
}
</script>

<template>
  <Card class="welcome-card">
    <div class="welcome-content">
      <div class="welcome-header">
        <div class="welcome-icon">
          <SmileOutlined />
        </div>
        <div class="welcome-text">
          <Title :level="3" style="margin: 0">
            {{ greeting }}, {{ userName }}! 👋
          </Title>
          <div class="welcome-subtitle">
            <Text type="secondary">Selamat datang di Sistem Manajemen Sekolah</Text>
            <Tag color="blue">{{ userRole }}</Tag>
          </div>
        </div>
        <Button 
          v-if="dismissible" 
          type="text" 
          size="small" 
          class="dismiss-btn"
          @click="handleDismiss"
        >
          Tutup
        </Button>
      </div>

      <div v-if="quickActions.length > 0" class="quick-actions">
        <Text type="secondary" class="actions-label">Mulai dengan:</Text>
        <Space wrap>
          <Button 
            v-for="action in quickActions" 
            :key="action.action"
            @click="handleAction(action.action)"
          >
            <template #icon>
              <component :is="action.icon" />
            </template>
            {{ action.label }}
          </Button>
        </Space>
      </div>

      <div v-if="showSteps" class="onboarding-steps">
        <Text type="secondary" class="steps-label">Langkah Awal:</Text>
        <Steps size="small" :current="0" class="steps">
          <Step title="Login" description="Anda sudah login" status="finish">
            <template #icon><CheckCircleOutlined /></template>
          </Step>
          <Step title="Jelajahi" description="Lihat menu di sidebar" status="process">
            <template #icon><RocketOutlined /></template>
          </Step>
          <Step title="Mulai" description="Kelola data sekolah" status="wait" />
        </Steps>
      </div>
    </div>
  </Card>
</template>

<style scoped>
.welcome-card {
  background: linear-gradient(135deg, #fff7ed 0%, #ffedd5 100%);
  border: 1px solid #fed7aa;
  border-radius: 12px;
  margin-bottom: 24px;
}

.welcome-content {
  padding: 8px;
}

.welcome-header {
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

.welcome-icon {
  width: 48px;
  height: 48px;
  background: #f97316;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
  flex-shrink: 0;
}

.welcome-text {
  flex: 1;
}

.welcome-subtitle {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 4px;
  flex-wrap: wrap;
}

.dismiss-btn {
  color: #8c8c8c;
}

.quick-actions {
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid #fed7aa;
}

.actions-label,
.steps-label {
  display: block;
  margin-bottom: 12px;
  font-size: 13px;
}

.onboarding-steps {
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid #fed7aa;
}

.steps {
  margin-top: 8px;
}

:deep(.ant-steps-item-icon) {
  background: #f97316 !important;
  border-color: #f97316 !important;
}

:deep(.ant-steps-item-finish .ant-steps-item-icon) {
  background: #52c41a !important;
  border-color: #52c41a !important;
}
</style>
