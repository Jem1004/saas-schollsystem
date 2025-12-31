<script setup lang="ts">
import { Result, Button, Typography, Space } from 'ant-design-vue'
import { WifiOutlined, ReloadOutlined, CloudOutlined } from '@ant-design/icons-vue'

const { Text, Paragraph } = Typography

interface Props {
  title?: string
  description?: string
  showRetry?: boolean
  retryText?: string
  loading?: boolean
}

withDefaults(defineProps<Props>(), {
  title: 'Koneksi Terputus',
  description: 'Tidak dapat terhubung ke server. Periksa koneksi internet Anda dan coba lagi.',
  showRetry: true,
  retryText: 'Coba Lagi',
  loading: false,
})

const emit = defineEmits<{
  (e: 'retry'): void
}>()

const handleRetry = () => {
  emit('retry')
}
</script>

<template>
  <div class="network-error">
    <Result status="warning">
      <template #icon>
        <div class="error-icon">
          <WifiOutlined class="wifi-icon" />
          <div class="error-badge">!</div>
        </div>
      </template>
      <template #title>
        <Text strong>{{ title }}</Text>
      </template>
      <template #subTitle>
        <Paragraph type="secondary">{{ description }}</Paragraph>
      </template>
      <template #extra>
        <Space direction="vertical" align="center">
          <Button 
            v-if="showRetry"
            type="primary" 
            :loading="loading"
            @click="handleRetry"
          >
            <template #icon><ReloadOutlined /></template>
            {{ retryText }}
          </Button>
          <Text type="secondary" class="hint-text">
            <CloudOutlined /> Pastikan Anda terhubung ke internet
          </Text>
        </Space>
      </template>
    </Result>
  </div>
</template>

<style scoped>
.network-error {
  padding: 48px 24px;
  text-align: center;
}

.error-icon {
  position: relative;
  display: inline-block;
}

.wifi-icon {
  font-size: 64px;
  color: #faad14;
}

.error-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  width: 24px;
  height: 24px;
  background: #ff4d4f;
  border-radius: 50%;
  color: white;
  font-weight: bold;
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.hint-text {
  font-size: 12px;
}
</style>
