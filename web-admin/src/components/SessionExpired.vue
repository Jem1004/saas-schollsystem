<script setup lang="ts">
import { Modal, Button, Typography, Space } from 'ant-design-vue'
import { ClockCircleOutlined, LoginOutlined } from '@ant-design/icons-vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const { Text, Paragraph } = Typography

interface Props {
  open: boolean
  title?: string
  message?: string
}

withDefaults(defineProps<Props>(), {
  title: 'Sesi Berakhir',
  message: 'Sesi Anda telah berakhir. Silakan login kembali untuk melanjutkan.',
})

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'login'): void
}>()

const router = useRouter()
const authStore = useAuthStore()

const handleLogin = () => {
  authStore.clearAuth()
  emit('update:open', false)
  router.push('/login')
  emit('login')
}

const handleStay = () => {
  emit('update:open', false)
}
</script>

<template>
  <Modal
    :open="open"
    :closable="false"
    :mask-closable="false"
    centered
    width="400"
    :footer="null"
  >
    <div class="session-expired">
      <div class="icon-wrapper">
        <ClockCircleOutlined class="expired-icon" />
      </div>
      <Text strong class="title">{{ title }}</Text>
      <Paragraph type="secondary" class="message">
        {{ message }}
      </Paragraph>
      <Space direction="vertical" class="actions">
        <Button type="primary" block @click="handleLogin">
          <template #icon><LoginOutlined /></template>
          Login Kembali
        </Button>
        <Button block @click="handleStay">
          Tetap di Halaman Ini
        </Button>
      </Space>
    </div>
  </Modal>
</template>

<style scoped>
.session-expired {
  text-align: center;
  padding: 24px 0;
}

.icon-wrapper {
  width: 80px;
  height: 80px;
  margin: 0 auto 24px;
  background: #fff7e6;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.expired-icon {
  font-size: 40px;
  color: #fa8c16;
}

.title {
  display: block;
  font-size: 18px;
  margin-bottom: 8px;
}

.message {
  margin-bottom: 24px;
}

.actions {
  width: 100%;
}
</style>
