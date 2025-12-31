import { ref, computed } from 'vue'
import { Modal } from 'ant-design-vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { h } from 'vue'

export interface SensitiveDataOptions {
  title?: string
  description?: string
  requireConfirmation?: boolean
  blurByDefault?: boolean
}

export function useSensitiveData(options: SensitiveDataOptions = {}) {
  const {
    title = 'Data Sensitif',
    description = 'Data ini bersifat rahasia dan hanya dapat diakses oleh pihak yang berwenang.',
    requireConfirmation = true,
    blurByDefault = true,
  } = options

  const isRevealed = ref(!blurByDefault)
  const hasConfirmed = ref(false)

  const showConfirmation = (): Promise<boolean> => {
    return new Promise((resolve) => {
      Modal.confirm({
        title: title,
        icon: h(ExclamationCircleOutlined),
        content: h('div', [
          h('p', description),
          h('p', { style: { marginTop: '8px', color: '#8c8c8c', fontSize: '12px' } }, 
            'Dengan melanjutkan, Anda menyetujui bahwa akses ini akan dicatat dalam sistem.')
        ]),
        okText: 'Ya, Tampilkan',
        cancelText: 'Batal',
        okType: 'primary',
        centered: true,
        onOk() {
          hasConfirmed.value = true
          isRevealed.value = true
          resolve(true)
        },
        onCancel() {
          resolve(false)
        },
      })
    })
  }

  const reveal = async (): Promise<boolean> => {
    if (isRevealed.value) return true
    
    if (requireConfirmation && !hasConfirmed.value) {
      return await showConfirmation()
    }
    
    isRevealed.value = true
    return true
  }

  const hide = () => {
    isRevealed.value = false
  }

  const toggle = async (): Promise<boolean> => {
    if (isRevealed.value) {
      hide()
      return true
    }
    return await reveal()
  }

  const displayValue = computed(() => (value: string, maskChar = '•') => {
    if (isRevealed.value) return value
    return maskChar.repeat(Math.min(value.length, 20))
  })

  return {
    isRevealed,
    hasConfirmed,
    reveal,
    hide,
    toggle,
    displayValue,
    showConfirmation,
  }
}

// Utility function to mask sensitive text
export function maskSensitiveText(text: string, visibleChars = 0, maskChar = '•'): string {
  if (!text) return ''
  if (visibleChars >= text.length) return text
  
  const visible = text.slice(0, visibleChars)
  const masked = maskChar.repeat(Math.min(text.length - visibleChars, 15))
  return visible + masked
}

// Utility function to check if user has permission to view sensitive data
export function canViewSensitiveData(userRole: string, dataType: string): boolean {
  const permissions: Record<string, string[]> = {
    internal_counseling_note: ['guru_bk'],
    parent_summary: ['guru_bk', 'wali_kelas'],
    student_violation_detail: ['guru_bk', 'wali_kelas'],
    student_achievement_detail: ['guru_bk', 'wali_kelas'],
  }

  const allowedRoles = permissions[dataType] || []
  return allowedRoles.includes(userRole)
}
