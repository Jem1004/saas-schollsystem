<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { UserOutlined, LockOutlined } from '@ant-design/icons-vue'
import type { Rule } from 'ant-design-vue/es/form'
import type { ValidateErrorEntity } from 'ant-design-vue/es/form/interface'
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
  ],
  password: [
    { required: true, message: 'Password wajib diisi', trigger: 'blur' },
  ],
}

const handleLogin = async () => {
  // Prevent double submission
  if (loading.value) {
    console.log('Login already in progress, ignoring duplicate call')
    return
  }
  
  loading.value = true
  console.log('Starting login process...', formState)
  
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
    } else {
      message.error('Terjadi kesalahan')
    }
  } finally {
    loading.value = false
    console.log('Login process completed')
  }
}

const onFinishFailed = (errorInfo: ValidateErrorEntity) => {
  console.log('Form validation failed:', errorInfo)
}
</script>

<template>
  <div class="login-container">
    <div class="login-content">
      <div class="login-header">
        <div class="logo-wrapper">
          <img src="@/assets/styles/logo.png" alt="Logo" class="logo-img" />
        </div>
        <h1 class="brand-title">Sistem Manajemen Sekolah</h1>
        <p class="brand-subtitle">Silakan masuk dengan akun Anda untuk melanjutkan</p>
      </div>

      <div class="login-card">
        <a-form
          :model="formState"
          :rules="rules"
          layout="vertical"
          @finish="handleLogin"
          @finishFailed="onFinishFailed"
          class="auth-form"
          hideRequiredMark
        >
          <a-form-item name="username" label="Username / Email">
            <a-input
              v-model:value="formState.username"
              placeholder="Contoh: uks@sekolah.sch.id"
              size="large"
              :disabled="loading"
              class="modern-input"
            >
              <template #prefix>
                <UserOutlined class="field-icon" />
              </template>
            </a-input>
          </a-form-item>

          <a-form-item name="password" label="Password">
            <a-input-password
              v-model:value="formState.password"
              placeholder="Masukkan kata sandi"
              size="large"
              :disabled="loading"
              class="modern-input"
            >
              <template #prefix>
                <LockOutlined class="field-icon" />
              </template>
            </a-input-password>
          </a-form-item>

          <a-form-item class="submit-item">
            <a-button
              type="primary"
              html-type="submit"
              size="large"
              block
              :loading="loading"
              class="submit-btn"
            >
              {{ loading ? 'Memproses...' : 'Masuk ke Dashboard' }}
            </a-button>
          </a-form-item>
        </a-form>
      </div>

      <div class="login-footer">
        <p class="copyright">
          &copy; {{ new Date().getFullYear() }} Sistem Manajemen Sekolah. All rights reserved.
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  min-height: 100vh;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f8fafc;
  background-image: radial-gradient(#e2e8f0 1px, transparent 1px);
  background-size: 24px 24px;
  padding: 20px;
}

.login-content {
  width: 100%;
  max-width: 440px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo-wrapper {
  margin: 0 auto 20px;
  display: flex;
  justify-content: center;
}

.logo-img {
  width: 80px;
  height: 80px;
  object-fit: contain;
}

.brand-title {
  font-size: 24px;
  font-weight: 700;
  color: #0f172a;
  margin-bottom: 8px;
  letter-spacing: -0.5px;
}

.brand-subtitle {
  font-size: 15px;
  color: #64748b;
  line-height: 1.5;
  font-weight: 400;
}

.login-card {
  width: 100%;
  background: #ffffff;
  border-radius: 16px;
  padding: 32px;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.05), 0 4px 6px -2px rgba(0, 0, 0, 0.025);
  border: 1px solid #f1f5f9;
  margin-bottom: 24px;
}

/* Form Styles */
.auth-form :deep(.ant-form-item-label > label) {
  font-size: 14px;
  color: #334155;
  font-weight: 500;
}

.modern-input {
  border-radius: 8px;
  padding: 8px 11px;
  border-color: #e2e8f0;
  font-size: 15px;
  transition: all 0.2s ease;
}

.modern-input:hover,
.modern-input:focus,
.modern-input:focus-within {
  border-color: #f97316;
  box-shadow: 0 0 0 2px rgba(249, 115, 22, 0.1);
}

.field-icon {
  color: #94a3b8;
  margin-right: 6px;
}

.submit-item {
  margin-top: 24px;
  margin-bottom: 0;
}

.submit-btn {
  height: 46px;
  font-size: 15px;
  font-weight: 600;
  border-radius: 8px;
  background-color: #f97316;
  border-color: #f97316;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  transition: all 0.2s ease;
}

.submit-btn:hover {
  background-color: #ea580c;
  border-color: #ea580c;
  transform: translateY(-1px);
  box-shadow: 0 4px 6px -1px rgba(249, 115, 22, 0.2);
}

.login-footer {
  text-align: center;
}

.copyright {
  font-size: 13px;
  color: #94a3b8;
  margin: 0;
}

/* Responsiveness */
@media (max-width: 480px) {
  .login-card {
    padding: 24px;
    box-shadow: none;
    background: transparent;
    border: none;
  }
  
  .login-container {
    background: #fff;
    align-items: flex-start;
    padding-top: 60px;
  }

  .logo-wrapper {
    box-shadow: none;
    background: transparent;
    border: none;
  }
}
</style>
