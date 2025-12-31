<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { UserOutlined, LockOutlined } from '@ant-design/icons-vue'
import type { Rule } from 'ant-design-vue/es/form'
import { useAuthStore } from '@/stores/auth'
import { authService } from '@/services'
import type { LoginRequest } from '@/types/user'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)

const formState = reactive<LoginRequest>({
  username: '',
  password: '',
})

const rules: Record<string, Rule[]> = {
  username: [
    { required: true, message: 'Username atau email wajib diisi', trigger: 'blur' },
    { min: 3, message: 'Minimal 3 karakter', trigger: 'blur' },
  ],
  password: [
    { required: true, message: 'Password wajib diisi', trigger: 'blur' },
    { min: 6, message: 'Minimal 6 karakter', trigger: 'blur' },
  ],
}

const handleLogin = async () => {
  // Prevent double submission
  if (loading.value) {
    console.log('Login already in progress, ignoring duplicate call')
    return
  }
  
  loading.value = true
  console.log('Starting login process...')
  
  try {
    const response = await authService.login(formState)
    console.log('Login API response received:', { user: response.user.username, role: response.user.role })
    
    // Store tokens and user data
    authStore.setTokens(response.accessToken, response.refreshToken)
    authStore.setUser(response.user)
    console.log('Tokens and user data stored')
    
    message.success('Login berhasil!')
    
    // Small delay to ensure localStorage is updated
    await new Promise(resolve => setTimeout(resolve, 50))
    
    // Check if user needs to change password
    const targetPath = response.user.mustResetPwd ? '/change-password' : '/dashboard'
    console.log('Navigating to:', targetPath)
    
    // Navigate to target path
    router.push(targetPath).catch((err) => {
      // Ignore navigation duplicated errors
      if (err.name !== 'NavigationDuplicated') {
        console.error('Navigation error:', err)
      }
    })
  } catch (error: any) {
    console.error('Login error:', error)
    
    // Only show error if it's actually an error response from API
    if (error.response?.data?.error?.message) {
      message.error(error.response.data.error.message)
    } else if (error.response?.status === 401) {
      message.error('Username atau password salah')
    } else if (error.response?.status) {
      message.error('Login gagal. Silakan coba lagi.')
    } else if (error.code === 'ERR_NETWORK') {
      message.error('Tidak dapat terhubung ke server')
    }
    // Don't show error for other cases (like navigation errors)
  } finally {
    loading.value = false
    console.log('Login process completed')
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <div class="logo">
          <img src="/vite.svg" alt="Logo" class="logo-img" />
        </div>
        <h1 class="title">Sistem Manajemen Sekolah</h1>
        <p class="subtitle">Silakan masuk untuk melanjutkan</p>
      </div>

      <a-form
        :model="formState"
        :rules="rules"
        layout="vertical"
        @finish="handleLogin"
        @submit.prevent
        class="login-form"
      >
        <a-form-item name="username" label="Username / Email">
          <a-input
            v-model:value="formState.username"
            placeholder="Masukkan username atau email"
            size="large"
            :disabled="loading"
          >
            <template #prefix>
              <UserOutlined class="input-icon" />
            </template>
          </a-input>
        </a-form-item>

        <a-form-item name="password" label="Password">
          <a-input-password
            v-model:value="formState.password"
            placeholder="Masukkan password"
            size="large"
            :disabled="loading"
          >
            <template #prefix>
              <LockOutlined class="input-icon" />
            </template>
          </a-input-password>
        </a-form-item>

        <a-form-item>
          <a-button
            type="primary"
            html-type="submit"
            size="large"
            block
            :loading="loading"
          >
            {{ loading ? 'Memproses...' : 'Masuk' }}
          </a-button>
        </a-form-item>
      </a-form>

      <div class="login-footer">
        <p class="footer-text">
          &copy; {{ new Date().getFullYear() }} Sistem Manajemen Sekolah
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f5f5f5 0%, #e8e8e8 100%);
  padding: 20px;
}

.login-card {
  width: 100%;
  max-width: 400px;
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
  padding: 40px;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo {
  margin-bottom: 16px;
}

.logo-img {
  width: 64px;
  height: 64px;
}

.title {
  font-size: 24px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0 0 8px 0;
}

.subtitle {
  font-size: 14px;
  color: #666666;
  margin: 0;
}

.login-form {
  margin-bottom: 24px;
}

.input-icon {
  color: #999999;
}

.login-footer {
  text-align: center;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

.footer-text {
  font-size: 12px;
  color: #999999;
  margin: 0;
}
</style>
