<script setup lang="ts">
import { useRouter } from 'vue-router'
import { FileSearchOutlined } from '@ant-design/icons-vue'

const router = useRouter()

const goBack = () => {
  router.back()
}

const goHome = () => {
  router.push('/dashboard')
}
</script>

<template>
  <div class="error-container">
    <div class="error-content">
      <FileSearchOutlined class="error-icon" />
      <h1 class="error-title">404</h1>
      <h2 class="error-subtitle">Halaman Tidak Ditemukan</h2>
      <p class="error-message">
        Maaf, halaman yang Anda cari tidak ditemukan.
        Halaman mungkin telah dipindahkan atau dihapus.
      </p>
      <div class="error-actions">
        <a-button type="primary" size="large" @click="goHome">
          Ke Dashboard
        </a-button>
        <a-button size="large" @click="goBack">
          Kembali
        </a-button>
      </div>
    </div>
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

.error-content {
  text-align: center;
  max-width: 500px;
}

.error-icon {
  font-size: 80px;
  color: #f97316;
  margin-bottom: 24px;
}

.error-title {
  font-size: 72px;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0 0 8px 0;
  line-height: 1;
}

.error-subtitle {
  font-size: 24px;
  font-weight: 600;
  color: #333333;
  margin: 0 0 16px 0;
}

.error-message {
  font-size: 16px;
  color: #666666;
  margin: 0 0 32px 0;
  line-height: 1.6;
}

.error-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
}
</style>
