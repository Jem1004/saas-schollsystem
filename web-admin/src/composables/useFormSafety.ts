import { ref, onMounted, onUnmounted } from 'vue'
import { Modal } from 'ant-design-vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { h } from 'vue'
import { onBeforeRouteLeave } from 'vue-router'

export interface FormSafetyOptions {
  enabled?: boolean
  message?: string
  title?: string
}

export function useFormSafety(
  isDirty: () => boolean,
  options: FormSafetyOptions = {}
) {
  const {
    enabled = true,
    message = 'Anda memiliki perubahan yang belum disimpan. Yakin ingin meninggalkan halaman ini?',
    title = 'Perubahan Belum Disimpan',
  } = options

  const isNavigating = ref(false)

  // Handle browser beforeunload event
  const handleBeforeUnload = (e: BeforeUnloadEvent) => {
    if (enabled && isDirty() && !isNavigating.value) {
      e.preventDefault()
      e.returnValue = message
      return message
    }
  }

  // Setup beforeunload listener
  onMounted(() => {
    window.addEventListener('beforeunload', handleBeforeUnload)
  })

  onUnmounted(() => {
    window.removeEventListener('beforeunload', handleBeforeUnload)
  })

  // Handle Vue Router navigation
  onBeforeRouteLeave((_to, _from, next) => {
    if (enabled && isDirty() && !isNavigating.value) {
      Modal.confirm({
        title: title,
        icon: h(ExclamationCircleOutlined),
        content: message,
        okText: 'Ya, Tinggalkan',
        cancelText: 'Tetap di Sini',
        okType: 'danger',
        centered: true,
        onOk() {
          isNavigating.value = true
          next()
        },
        onCancel() {
          next(false)
        },
      })
    } else {
      next()
    }
  })

  // Allow programmatic navigation bypass
  const allowNavigation = () => {
    isNavigating.value = true
  }

  // Reset navigation flag
  const resetNavigation = () => {
    isNavigating.value = false
  }

  return {
    isNavigating,
    allowNavigation,
    resetNavigation,
  }
}

// Utility function to sanitize input for XSS prevention
export function sanitizeInput(input: string): string {
  if (!input) return ''
  
  // Basic HTML entity encoding
  const map: Record<string, string> = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#x27;',
    '/': '&#x2F;',
  }
  
  return input.replace(/[&<>"'/]/g, (char) => map[char] || char)
}

// Utility function to validate required fields
export function validateRequired(value: unknown, fieldName: string): string | null {
  if (value === null || value === undefined || value === '') {
    return `${fieldName} wajib diisi`
  }
  if (typeof value === 'string' && value.trim() === '') {
    return `${fieldName} wajib diisi`
  }
  if (Array.isArray(value) && value.length === 0) {
    return `${fieldName} wajib dipilih`
  }
  return null
}

// Utility function to validate email format
export function validateEmail(email: string): string | null {
  if (!email) return null
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(email)) {
    return 'Format email tidak valid'
  }
  return null
}

// Utility function to validate phone number (Indonesian format)
export function validatePhone(phone: string): string | null {
  if (!phone) return null
  const phoneRegex = /^(\+62|62|0)[0-9]{9,12}$/
  if (!phoneRegex.test(phone.replace(/[\s-]/g, ''))) {
    return 'Format nomor telepon tidak valid'
  }
  return null
}

// Utility function to validate NISN (10 digits)
export function validateNISN(nisn: string): string | null {
  if (!nisn) return null
  const nisnRegex = /^[0-9]{10}$/
  if (!nisnRegex.test(nisn)) {
    return 'NISN harus 10 digit angka'
  }
  return null
}

// Utility function to validate password strength
export function validatePassword(password: string): string | null {
  if (!password) return null
  if (password.length < 8) {
    return 'Password minimal 8 karakter'
  }
  if (!/[A-Z]/.test(password)) {
    return 'Password harus mengandung huruf besar'
  }
  if (!/[a-z]/.test(password)) {
    return 'Password harus mengandung huruf kecil'
  }
  if (!/[0-9]/.test(password)) {
    return 'Password harus mengandung angka'
  }
  return null
}

// Utility function to validate min/max length
export function validateLength(
  value: string,
  fieldName: string,
  min?: number,
  max?: number
): string | null {
  if (!value) return null
  if (min !== undefined && value.length < min) {
    return `${fieldName} minimal ${min} karakter`
  }
  if (max !== undefined && value.length > max) {
    return `${fieldName} maksimal ${max} karakter`
  }
  return null
}

// Utility function to validate number range
export function validateNumberRange(
  value: number,
  fieldName: string,
  min?: number,
  max?: number
): string | null {
  if (value === null || value === undefined) return null
  if (min !== undefined && value < min) {
    return `${fieldName} minimal ${min}`
  }
  if (max !== undefined && value > max) {
    return `${fieldName} maksimal ${max}`
  }
  return null
}
