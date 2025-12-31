<script setup lang="ts">
import { computed } from 'vue'
import { Modal, Button, Typography, Tag } from 'ant-design-vue'
import { LockOutlined, ArrowLeftOutlined, HomeOutlined } from '@ant-design/icons-vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { getRoleDisplayName } from '@/composables/usePermissions'

const { Text, Paragraph } = Typography

interface Props {
  open: boolean
  resource?: string
  action?: string
  requiredRoles?: string[]
}

withDefaults(defineProps<Props>(), {
  resource: 'halaman ini',
  action: 'mengakses',
  requiredRoles: () => [],
})

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'back'): void
}>()

const router = useRouter()
const authStore = useAuthStore()

const userRole = computed(() => {
  return authStore.userRole ? getRoleDisplayName(authStore.userRole) : 'User'
})

const handleBack = () => {
  emit('update:open', false)
  router.back()
  emit('back')
}

const handleHome = () => {
  emit('update:open', false)
  router.push('/dashboard')
}

const handleClose = () => {
  emit('update:open', false)
}
</script>

<template>
  <Modal
    :open="open"
    :closable="true"
    centered
    width="450"
    :footer="null"
    @cancel="handleClose"
  >
    <div class="permission-denied">
      <div class="icon-wrapper">
        <LockOutlined class="lock-icon" />
      </div>
      <Text strong class="title">Akses Ditolak</Text>
      <Paragraph type="secondary" class="message">
        Anda tidak memiliki izin untuk {{ action }} {{ resource }}.
      </Paragraph>
      
      <div class="role-info">
        <Text type="secondary">Role Anda saat ini:</Text>
        <Tag color="blue" class="role-tag">{{ userRole }}</Tag>
      </div>

      <div v-if="requiredRoles.length > 0" class="required-roles">
        <Text type="secondary">Role yang diperlukan:</Text>
        <div class="roles-list">
          <Tag v-for="role in requiredRoles" :key="role" color="orange">
            {{ role }}
          </Tag>
        </div>
      </div>

      <div class="actions">
        <Button type="primary" @click="handleBack">
          <template #icon><ArrowLeftOutlined /></template>
          Kembali
        </Button>
        <Button @click="handleHome">
          <template #icon><HomeOutlined /></template>
          Ke Dashboard
        </Button>
      </div>

      <Paragraph type="secondary" class="hint">
        Jika Anda memerlukan akses, silakan hubungi administrator.
      </Paragraph>
    </div>
  </Modal>
</template>

<style scoped>
.permission-denied {
  text-align: center;
  padding: 24px 0;
}

.icon-wrapper {
  width: 80px;
  height: 80px;
  margin: 0 auto 24px;
  background: #fff1f0;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.lock-icon {
  font-size: 40px;
  color: #ff4d4f;
}

.title {
  display: block;
  font-size: 18px;
  margin-bottom: 8px;
}

.message {
  margin-bottom: 16px;
}

.role-info {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-bottom: 12px;
}

.role-tag {
  margin: 0;
}

.required-roles {
  margin-bottom: 24px;
}

.roles-list {
  margin-top: 8px;
  display: flex;
  justify-content: center;
  gap: 8px;
  flex-wrap: wrap;
}

.actions {
  display: flex;
  gap: 12px;
  justify-content: center;
  margin-bottom: 16px;
}

.hint {
  font-size: 12px;
  margin: 0;
}
</style>
