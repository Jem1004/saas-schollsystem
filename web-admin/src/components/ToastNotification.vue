<script lang="ts">
import { notification } from 'ant-design-vue'
import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  InfoCircleOutlined,
  ExclamationCircleOutlined,
} from '@ant-design/icons-vue'
import { h, defineComponent } from 'vue'
import type { ToastType, ToastOptions } from '@/types/notification'

// Configure notification defaults
notification.config({
  placement: 'topRight',
  duration: 4.5,
  top: '24px',
})

// Icon mapping
const iconMap: Record<ToastType, typeof CheckCircleOutlined> = {
  success: CheckCircleOutlined,
  error: CloseCircleOutlined,
  warning: ExclamationCircleOutlined,
  info: InfoCircleOutlined,
}

// Color mapping
const colorMap: Record<ToastType, string> = {
  success: '#52c41a',
  error: '#ff4d4f',
  warning: '#faad14',
  info: '#1890ff',
}

// Get default title based on type
function getDefaultTitle(type: ToastType): string {
  const titles: Record<ToastType, string> = {
    success: 'Berhasil',
    error: 'Error',
    warning: 'Peringatan',
    info: 'Informasi',
  }
  return titles[type]
}

// Toast function
export function toast(options: ToastOptions) {
  const { type = 'info', title, message, duration = 4.5, placement = 'topRight' } = options
  
  const icon = h(iconMap[type], { style: { color: colorMap[type] } })
  
  notification[type]({
    message: title || getDefaultTitle(type),
    description: message,
    icon,
    duration,
    placement,
  })
}

// Shorthand methods
export function toastSuccess(message: string, title?: string) {
  toast({ type: 'success', message, title })
}

export function toastError(message: string, title?: string) {
  toast({ type: 'error', message, title: title || 'Error' })
}

export function toastWarning(message: string, title?: string) {
  toast({ type: 'warning', message, title: title || 'Peringatan' })
}

export function toastInfo(message: string, title?: string) {
  toast({ type: 'info', message, title: title || 'Informasi' })
}

// Close all notifications
export function closeAllToasts() {
  notification.destroy()
}

export default defineComponent({
  name: 'ToastNotification',
  setup() {
    return {}
  },
  render() {
    return null
  }
})
</script>
