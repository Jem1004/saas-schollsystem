<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import {
  Dropdown,
  Badge,
  Button,
  List,
  ListItem,
  ListItemMeta,
  Empty,
  Spin,
  Typography,
  Divider,
  Tag,
} from 'ant-design-vue'
import {
  BellOutlined,
  CheckOutlined,
  ClockCircleOutlined,
  SafetyOutlined,
  TrophyOutlined,
  FileTextOutlined,
  BookOutlined,
  FormOutlined,
  SettingOutlined,
} from '@ant-design/icons-vue'
import { notificationService } from '@/services/notification'
import type { Notification, NotificationType } from '@/types/notification'

const { Text } = Typography

// State
const loading = ref(false)
const notifications = ref<Notification[]>([])
const unreadCount = ref(0)
const dropdownVisible = ref(false)

// Mock data for development
const mockNotifications: Notification[] = [
  {
    id: 1,
    userId: 1,
    type: 'attendance_checkin',
    title: 'Absensi Masuk',
    message: 'Ahmad Fauzi telah check-in pada pukul 07:15',
    isRead: false,
    createdAt: new Date(Date.now() - 5 * 60 * 1000).toISOString(),
  },
  {
    id: 2,
    userId: 1,
    type: 'violation',
    title: 'Pelanggaran Baru',
    message: 'Budi Santoso melakukan pelanggaran: Terlambat masuk kelas',
    isRead: false,
    createdAt: new Date(Date.now() - 30 * 60 * 1000).toISOString(),
  },
  {
    id: 3,
    userId: 1,
    type: 'achievement',
    title: 'Prestasi Baru',
    message: 'Citra Dewi mendapat prestasi: Juara 1 Olimpiade Matematika',
    isRead: true,
    createdAt: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
  },
  {
    id: 4,
    userId: 1,
    type: 'permit',
    title: 'Izin Keluar',
    message: 'Dian Pratama mengajukan izin keluar: Sakit',
    isRead: true,
    createdAt: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(),
  },
]

// Icon mapping
const iconMap: Record<NotificationType, typeof BellOutlined> = {
  attendance_checkin: ClockCircleOutlined,
  attendance_checkout: ClockCircleOutlined,
  violation: SafetyOutlined,
  achievement: TrophyOutlined,
  permit: FileTextOutlined,
  counseling: FormOutlined,
  grade: BookOutlined,
  homeroom_note: FormOutlined,
  system: SettingOutlined,
}

// Color mapping
const colorMap: Record<NotificationType, string> = {
  attendance_checkin: '#52c41a',
  attendance_checkout: '#1890ff',
  violation: '#ff4d4f',
  achievement: '#faad14',
  permit: '#722ed1',
  counseling: '#13c2c2',
  grade: '#2f54eb',
  homeroom_note: '#eb2f96',
  system: '#8c8c8c',
}

// Load notifications
const loadNotifications = async () => {
  loading.value = true
  try {
    const response = await notificationService.getNotifications({ pageSize: 10 })
    notifications.value = response.data
    unreadCount.value = response.unreadCount
  } catch {
    // Use mock data on error
    notifications.value = mockNotifications
    unreadCount.value = mockNotifications.filter(n => !n.isRead).length
  } finally {
    loading.value = false
  }
}

// Mark notification as read
const handleMarkAsRead = async (notification: Notification) => {
  if (notification.isRead) return
  
  try {
    await notificationService.markAsRead(notification.id)
    notification.isRead = true
    unreadCount.value = Math.max(0, unreadCount.value - 1)
  } catch {
    // Update locally on error
    notification.isRead = true
    unreadCount.value = Math.max(0, unreadCount.value - 1)
  }
}

// Mark all as read
const handleMarkAllAsRead = async () => {
  try {
    await notificationService.markAllAsRead()
    notifications.value.forEach(n => n.isRead = true)
    unreadCount.value = 0
  } catch {
    // Update locally on error
    notifications.value.forEach(n => n.isRead = true)
    unreadCount.value = 0
  }
}

// Format time ago
const formatTimeAgo = (dateString: string): string => {
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)
  
  if (diffMins < 1) return 'Baru saja'
  if (diffMins < 60) return `${diffMins} menit lalu`
  if (diffHours < 24) return `${diffHours} jam lalu`
  if (diffDays < 7) return `${diffDays} hari lalu`
  
  return date.toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })
}

// Handle dropdown visibility change
const handleVisibleChange = (visible: boolean) => {
  dropdownVisible.value = visible
  if (visible) {
    loadNotifications()
  }
}

// Polling interval
let pollInterval: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  loadNotifications()
  // Poll for new notifications every 30 seconds
  pollInterval = setInterval(loadNotifications, 30000)
})

onUnmounted(() => {
  if (pollInterval) {
    clearInterval(pollInterval)
  }
})
</script>

<template>
  <Dropdown
    :open="dropdownVisible"
    trigger="click"
    placement="bottomRight"
    @openChange="handleVisibleChange"
  >
    <template #overlay>
      <div class="notification-dropdown">
        <div class="notification-header">
          <Text strong>Notifikasi</Text>
          <Button
            v-if="unreadCount > 0"
            type="link"
            size="small"
            @click="handleMarkAllAsRead"
          >
            <template #icon><CheckOutlined /></template>
            Tandai semua dibaca
          </Button>
        </div>
        
        <Divider style="margin: 0" />
        
        <div class="notification-content">
          <Spin v-if="loading" class="loading-spin" />
          
          <Empty
            v-else-if="notifications.length === 0"
            description="Tidak ada notifikasi"
            :image="Empty.PRESENTED_IMAGE_SIMPLE"
          />
          
          <List
            v-else
            :data-source="notifications"
            item-layout="horizontal"
            size="small"
          >
            <template #renderItem="{ item }: { item: Notification }">
              <ListItem
                :class="['notification-item', { unread: !item.isRead }]"
                @click="handleMarkAsRead(item)"
              >
                <ListItemMeta>
                  <template #avatar>
                    <div
                      class="notification-icon"
                      :style="{ backgroundColor: colorMap[item.type] + '20', color: colorMap[item.type] }"
                    >
                      <component :is="iconMap[item.type]" />
                    </div>
                  </template>
                  <template #title>
                    <div class="notification-title">
                      <Text :strong="!item.isRead">{{ item.title }}</Text>
                      <Tag v-if="!item.isRead" color="blue" size="small">Baru</Tag>
                    </div>
                  </template>
                  <template #description>
                    <div class="notification-desc">
                      <Text type="secondary" class="notification-message">
                        {{ item.message }}
                      </Text>
                      <Text type="secondary" class="notification-time">
                        {{ formatTimeAgo(item.createdAt) }}
                      </Text>
                    </div>
                  </template>
                </ListItemMeta>
              </ListItem>
            </template>
          </List>
        </div>
        
        <Divider style="margin: 0" />
        
        <div class="notification-footer">
          <Button type="link" block>
            Lihat semua notifikasi
          </Button>
        </div>
      </div>
    </template>
    
    <Badge :count="unreadCount" :overflow-count="99" class="notification-badge">
      <Button type="text" class="notification-btn">
        <BellOutlined />
      </Button>
    </Badge>
  </Dropdown>
</template>

<style scoped>
.notification-dropdown {
  width: 380px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 6px 16px 0 rgba(0, 0, 0, 0.08),
              0 3px 6px -4px rgba(0, 0, 0, 0.12),
              0 9px 28px 8px rgba(0, 0, 0, 0.05);
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
}

.notification-content {
  max-height: 400px;
  overflow-y: auto;
}

.loading-spin {
  display: flex;
  justify-content: center;
  padding: 24px;
}

.notification-item {
  padding: 12px 16px !important;
  cursor: pointer;
  transition: background-color 0.2s;
}

.notification-item:hover {
  background-color: #f5f5f5;
}

.notification-item.unread {
  background-color: #e6f7ff;
}

.notification-item.unread:hover {
  background-color: #bae7ff;
}

.notification-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

.notification-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.notification-desc {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.notification-message {
  font-size: 13px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.notification-time {
  font-size: 12px;
}

.notification-footer {
  padding: 8px 16px;
}

.notification-badge {
  margin-right: 8px;
}

.notification-btn {
  font-size: 18px;
  width: 40px;
  height: 40px;
}
</style>
