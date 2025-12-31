<script setup lang="ts">
import { ref, onErrorCaptured } from 'vue'
import { Button, Result, Typography, Card } from 'ant-design-vue'
import { ReloadOutlined, HomeOutlined, BugOutlined } from '@ant-design/icons-vue'
import { useRouter } from 'vue-router'

const { Text } = Typography

interface Props {
  fallbackTitle?: string
  fallbackDescription?: string
  showDetails?: boolean
}

withDefaults(defineProps<Props>(), {
  fallbackTitle: 'Terjadi Kesalahan',
  fallbackDescription: 'Maaf, terjadi kesalahan yang tidak terduga. Silakan coba lagi.',
  showDetails: false,
})

const router = useRouter()

const hasError = ref(false)
const errorMessage = ref('')
const errorStack = ref('')

onErrorCaptured((err: Error) => {
  hasError.value = true
  errorMessage.value = err.message || 'Unknown error'
  errorStack.value = err.stack || ''
  
  // Log error for debugging
  console.error('ErrorBoundary caught:', err)
  
  // Return false to prevent error from propagating
  return false
})

const handleRetry = () => {
  hasError.value = false
  errorMessage.value = ''
  errorStack.value = ''
  // Force re-render by reloading the page
  window.location.reload()
}

const handleGoHome = () => {
  hasError.value = false
  router.push('/dashboard')
}

const handleReportError = () => {
  // In production, this would send error to monitoring service
  console.log('Error reported:', { message: errorMessage.value, stack: errorStack.value })
  alert('Laporan error telah dikirim. Terima kasih!')
}
</script>

<template>
  <div v-if="hasError" class="error-boundary">
    <Result
      status="error"
      :title="fallbackTitle"
      :sub-title="fallbackDescription"
    >
      <template #extra>
        <div class="error-actions">
          <Button type="primary" @click="handleRetry">
            <template #icon><ReloadOutlined /></template>
            Coba Lagi
          </Button>
          <Button @click="handleGoHome">
            <template #icon><HomeOutlined /></template>
            Ke Dashboard
          </Button>
        </div>
      </template>
      
      <div v-if="showDetails && errorMessage" class="error-details">
        <Card size="small" class="error-card">
          <template #title>
            <BugOutlined /> Detail Error
          </template>
          <template #extra>
            <Button type="link" size="small" @click="handleReportError">
              Laporkan
            </Button>
          </template>
          <Text code class="error-message">{{ errorMessage }}</Text>
          <details v-if="errorStack" class="error-stack">
            <summary>Stack Trace</summary>
            <pre>{{ errorStack }}</pre>
          </details>
        </Card>
      </div>
    </Result>
  </div>
  <slot v-else />
</template>

<style scoped>
.error-boundary {
  min-height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.error-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.error-details {
  margin-top: 24px;
  max-width: 600px;
  text-align: left;
}

.error-card {
  background: #fff2f0;
  border-color: #ffccc7;
}

.error-message {
  display: block;
  word-break: break-all;
}

.error-stack {
  margin-top: 12px;
}

.error-stack summary {
  cursor: pointer;
  color: #8c8c8c;
  font-size: 12px;
}

.error-stack pre {
  margin-top: 8px;
  padding: 12px;
  background: #f5f5f5;
  border-radius: 4px;
  font-size: 11px;
  overflow-x: auto;
  max-height: 200px;
}
</style>
