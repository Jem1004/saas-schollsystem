import api from './api'
import type { Notification, NotificationListResponse } from '@/types/notification'

export interface GetNotificationsParams {
  page?: number
  pageSize?: number
  unreadOnly?: boolean
}

export const notificationService = {
  // Get notifications list
  async getNotifications(params?: GetNotificationsParams): Promise<NotificationListResponse> {
    const response = await api.get<NotificationListResponse>('/notifications', { params })
    return response.data
  },

  // Get unread count
  async getUnreadCount(): Promise<number> {
    const response = await api.get<{ count: number }>('/notifications/unread-count')
    return response.data.count
  },

  // Mark notification as read
  async markAsRead(notificationId: number): Promise<Notification> {
    const response = await api.put<Notification>(`/notifications/${notificationId}/read`)
    return response.data
  },

  // Mark all notifications as read
  async markAllAsRead(): Promise<void> {
    await api.put('/notifications/read-all')
  },

  // Delete notification
  async deleteNotification(notificationId: number): Promise<void> {
    await api.delete(`/notifications/${notificationId}`)
  },
}

export default notificationService
