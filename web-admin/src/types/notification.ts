// Notification types

export interface Notification {
  id: number
  userId: number
  type: NotificationType
  title: string
  message: string
  data?: Record<string, unknown>
  isRead: boolean
  createdAt: string
}

export type NotificationType = 
  | 'attendance_checkin'
  | 'attendance_checkout'
  | 'violation'
  | 'achievement'
  | 'permit'
  | 'counseling'
  | 'grade'
  | 'homeroom_note'
  | 'system'

export interface NotificationListResponse {
  data: Notification[]
  total: number
  unreadCount: number
  page: number
  pageSize: number
}

// Toast notification types
export type ToastType = 'success' | 'error' | 'warning' | 'info'

export interface ToastOptions {
  type?: ToastType
  title?: string
  message: string
  duration?: number
  closable?: boolean
  placement?: 'topLeft' | 'topRight' | 'bottomLeft' | 'bottomRight'
}
