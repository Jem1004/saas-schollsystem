import api from './api'
import type { NotificationListResponse } from '@/types/notification'

export interface GetNotificationsParams {
  page?: number
  pageSize?: number
  unreadOnly?: boolean
}

interface ApiResponse<T> {
  success: boolean
  data: T
}

export const notificationService = {
  // Get notifications list
  async getNotifications(params?: GetNotificationsParams): Promise<NotificationListResponse> {
    const response = await api.get<ApiResponse<NotificationListResponse>>('/notifications', { params })
    return response.data.data
  },

  // Get unread count
  async getUnreadCount(): Promise<number> {
    const response = await api.get<ApiResponse<{ count: number }>>('/notifications/unread-count')
    return response.data.data.count
  },

  // Mark notification as read
  async markAsRead(notificationId: number): Promise<void> {
    await api.post(`/notifications/${notificationId}/read`)
  },

  // Mark all notifications as read
  async markAllAsRead(): Promise<void> {
    await api.post('/notifications/read-all')
  },

  // Delete notification
  async deleteNotification(notificationId: number): Promise<void> {
    await api.delete(`/notifications/${notificationId}`)
  },
}

export default notificationService
