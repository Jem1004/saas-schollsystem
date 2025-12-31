<script setup lang="ts">
import { reactive, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { LockOutlined, CheckCircleOutlined, CloseCircleOutlined } from '@ant-design/icons-vue'
import type { Rule } from 'ant-design-vue/es/form'
import { useAuthStore } from '@/stores/auth'
import { authService } from '@/services'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)

interface ChangePasswordForm {
  oldPassword: string
  newPassword: string
  confirmPassword: string
}

const formState = reactive<ChangePasswordForm>({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

// Password strength validation
const passwordStrength = computed(() => {
  const password = formState.newPassword
  const checks = {
    minLength: password.length >= 8,
    hasUppercase: /[A-Z]/.test(password),
    hasLowercase: /[a-z]/.test(password),
    hasNumber: /[0-9]/.test(password),
    hasSpecial: /[!@#$%^&*(),.?":{}|<>]/.test(password),
  }
  
  const passedChecks = Object.values(checks).filter(Boolean).length
  
  return {
    checks,
    score: passedChecks,
    level: passedChecks <= 2 ? 'weak' : passedChecks <= 3 ? 'medium' : 'strong',
  }
})

const validatePassword = async (_rule: Rule, value: string) => {
  if (!value) {
    return Promise.reject('Password baru wajib diisi')
  }
  if (value.length < 8) {
    return Promise.reject('Password minimal 8 karakter')
  }
  if (passwordStrength.value.score < 3) {
    return Promise.reject('Password terlalu lemah')
  }
  return Promise.resolve()
}

const validateConfirmPassword = async (_rule: Rule, value: string) => {
  if (!value) {
    return Promise.reject('Konfirmasi password wajib diisi')
  }
  if (value !== formState.newPassword) {
    return Promise.reject('Password tidak cocok')
  }
  return Promise.resolve()
}

const rules: Record<string, Rule[]> = {
  oldPassword: [
    { required: true, message: 'Password lama wajib diisi', trigger: 'blur' },
  ],
  newPassword: [
    { validator: validatePassword, trigger: 'change' },
  ],
  confirmPassword: [
    { validator: validateConfirmPassword, trigger: 'change' },
  ],
}

const handleChangePassword = async () => {
  loading.value = true
  try {
    await authService.changePassword(formState.oldPassword, formState.newPassword)
    
    // Update user state to reflect password has been changed
    if (authStore.user) {
      authStore.setUser({ ...authStore.user, mustResetPwd: false })
    }
    
    message.success('Password berhasil diubah!')
    router.push('/dashboard')
  } catch (error: any) {
    const errorMessage = error.response?.data?.error?.message || 'Gagal mengubah password'
    message.error(errorMessage)
  } finally {
    loading.value = false
  }
}

const getStrengthColor = (level: string) => {
  switch (level) {
    case 'weak': return '#ff4d4f'
    case 'medium': return '#faad14'
    case 'strong': return '#52c41a'
    default: return '#d9d9d9'
  }
}

const getStrengthText = (level: string) => {
  switch (level) {
    case 'weak': return 'Lemah'
    case 'medium': return 'Sedang'
    case 'strong': return 'Kuat'
    default: return ''
  }
}
</script>

<template>
  <div class="change-password-container">
    <div class="change-password-card">
      <div class="card-header">
        <h1 class="title">Ubah Password</h1>
        <p class="subtitle">
          {{ authStore.user?.mustResetPwd 
            ? 'Anda harus mengubah password untuk melanjutkan' 
            : 'Masukkan password lama dan password baru Anda' 
          }}
        </p>
      </div>

      <a-form
        :model="formState"
        :rules="rules"
        layout="vertical"
        @finish="handleChangePassword"
        class="change-password-form"
      >
        <a-form-item name="oldPassword" label="Password Lama">
          <a-input-password
            v-model:value="formState.oldPassword"
            placeholder="Masukkan password lama"
            size="large"
            :disabled="loading"
          >
            <template #prefix>
              <LockOutlined class="input-icon" />
            </template>
          </a-input-password>
        </a-form-item>

        <a-form-item name="newPassword" label="Password Baru">
          <a-input-password
            v-model:value="formState.newPassword"
            placeholder="Masukkan password baru"
            size="large"
            :disabled="loading"
          >
            <template #prefix>
              <LockOutlined class="input-icon" />
            </template>
          </a-input-password>
        </a-form-item>

        <!-- Password Strength Indicator -->
        <div v-if="formState.newPassword" class="password-strength">
          <div class="strength-bar">
            <div 
              class="strength-fill" 
              :style="{ 
                width: `${(passwordStrength.score / 5) * 100}%`,
                backgroundColor: getStrengthColor(passwordStrength.level)
              }"
            />
          </div>
          <span 
            class="strength-text"
            :style="{ color: getStrengthColor(passwordStrength.level) }"
          >
            {{ getStrengthText(passwordStrength.level) }}
          </span>
        </div>

        <!-- Password Requirements Checklist -->
        <div v-if="formState.newPassword" class="password-requirements">
          <div 
            v-for="(passed, key) in passwordStrength.checks" 
            :key="key"
            class="requirement-item"
            :class="{ passed }"
          >
            <CheckCircleOutlined v-if="passed" class="check-icon passed" />
            <CloseCircleOutlined v-else class="check-icon" />
            <span>{{ getRequirementText(key) }}</span>
          </div>
        </div>

        <a-form-item name="confirmPassword" label="Konfirmasi Password Baru">
          <a-input-password
            v-model:value="formState.confirmPassword"
            placeholder="Masukkan ulang password baru"
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
            {{ loading ? 'Memproses...' : 'Ubah Password' }}
          </a-button>
        </a-form-item>

        <a-form-item v-if="!authStore.user?.mustResetPwd">
          <a-button
            size="large"
            block
            @click="router.push('/dashboard')"
            :disabled="loading"
          >
            Batal
          </a-button>
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script lang="ts">
function getRequirementText(key: string): string {
  const texts: Record<string, string> = {
    minLength: 'Minimal 8 karakter',
    hasUppercase: 'Mengandung huruf besar',
    hasLowercase: 'Mengandung huruf kecil',
    hasNumber: 'Mengandung angka',
    hasSpecial: 'Mengandung karakter khusus',
  }
  return texts[key] || key
}
</script>

<style scoped>
.change-password-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f5f5f5 0%, #e8e8e8 100%);
  padding: 20px;
}

.change-password-card {
  width: 100%;
  max-width: 420px;
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
  padding: 40px;
}

.card-header {
  text-align: center;
  margin-bottom: 32px;
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

.change-password-form {
  margin-bottom: 16px;
}

.input-icon {
  color: #999999;
}

.password-strength {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.strength-bar {
  flex: 1;
  height: 4px;
  background: #f0f0f0;
  border-radius: 2px;
  overflow: hidden;
}

.strength-fill {
  height: 100%;
  transition: width 0.3s, background-color 0.3s;
}

.strength-text {
  font-size: 12px;
  font-weight: 500;
  min-width: 50px;
}

.password-requirements {
  background: #fafafa;
  border-radius: 8px;
  padding: 12px 16px;
  margin-bottom: 16px;
}

.requirement-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #999999;
  margin-bottom: 4px;
}

.requirement-item:last-child {
  margin-bottom: 0;
}

.requirement-item.passed {
  color: #52c41a;
}

.check-icon {
  font-size: 14px;
  color: #d9d9d9;
}

.check-icon.passed {
  color: #52c41a;
}
</style>
