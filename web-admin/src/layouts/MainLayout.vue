<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  Layout,
  LayoutSider,
  LayoutHeader,
  LayoutContent,
  Menu,
  MenuItem,
  Dropdown,
  Avatar,
  Button,
  Badge,
  Tag,
  List,
  ListItem,
  ListItemMeta,
  Empty,
  Spin,
  Typography,
} from 'ant-design-vue'
import type { MenuInfo } from 'ant-design-vue/es/menu/src/interface'
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  DashboardOutlined,
  BankOutlined,
  TeamOutlined,
  UserOutlined,
  LogoutOutlined,
  SettingOutlined,
  BookOutlined,
  SafetyOutlined,
  TrophyOutlined,
  FileTextOutlined,
  ScheduleOutlined,
  BellOutlined,
  DesktopOutlined,
  SolutionOutlined,
  FormOutlined,
  EyeOutlined,
  CheckOutlined,
  ClockCircleOutlined,
  SyncOutlined,
} from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { usePermissions, getRoleDisplayName, getRoleColor } from '@/composables/usePermissions'
import { notificationService } from '@/services'
import type { Notification } from '@/types/notification'
import type { UserRole } from '@/types/user'

const { Text } = Typography

interface MenuItemConfig {
  key: string
  icon: typeof DashboardOutlined
  label: string
  path?: string
  roles: UserRole[]
  readOnly?: boolean
  children?: MenuItemConfig[]
}

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
usePermissions()

const collapsed = ref(false)
const selectedKeys = ref<string[]>(['dashboard'])
const openKeys = ref<string[]>([])

// Notification state
const notificationDropdownVisible = ref(false)
const notifications = ref<Notification[]>([])
const unreadCount = ref(0)
const notificationLoading = ref(false)
let notificationInterval: ReturnType<typeof setInterval> | null = null

// Check if user is super admin
const isSuperAdmin = computed(() => authStore.userRole === 'super_admin')

// Role-based menu configuration
const allMenuItems: MenuItemConfig[] = [
  {
    key: 'dashboard',
    icon: DashboardOutlined,
    label: 'Dashboard',
    path: '/dashboard',
    roles: ['super_admin', 'admin_sekolah', 'guru_bk', 'wali_kelas'],
  },
  // Super Admin menus
  {
    key: 'tenants',
    icon: BankOutlined,
    label: 'Manajemen Sekolah',
    path: '/tenants',
    roles: ['super_admin'],
  },
  {
    key: 'devices',
    icon: DesktopOutlined,
    label: 'Manajemen Device',
    path: '/devices',
    roles: ['super_admin'],
  },
  // Admin Sekolah menus
  {
    key: 'classes',
    icon: BookOutlined,
    label: 'Manajemen Kelas',
    path: '/classes',
    roles: ['admin_sekolah'],
  },
  {
    key: 'students',
    icon: TeamOutlined,
    label: 'Manajemen Siswa',
    path: '/students',
    roles: ['admin_sekolah'],
  },
  {
    key: 'parents',
    icon: SolutionOutlined,
    label: 'Manajemen Orang Tua',
    path: '/parents',
    roles: ['admin_sekolah'],
  },
  {
    key: 'users',
    icon: UserOutlined,
    label: 'Manajemen User',
    path: '/users',
    roles: ['admin_sekolah'],
  },
  {
    key: 'attendance-report',
    icon: ScheduleOutlined,
    label: 'Laporan Absensi',
    path: '/attendance',
    roles: ['admin_sekolah'],
  },
  {
    key: 'live-attendance',
    icon: SyncOutlined,
    label: 'Absensi Real-Time',
    path: '/attendance/live',
    roles: ['admin_sekolah'],
  },
  // Note: School settings removed from super_admin menu as they don't have school context
  // Super admin should manage school settings through tenant management
  {
    key: 'school-settings',
    icon: SettingOutlined,
    label: 'Pengaturan Sekolah',
    path: '/settings',
    roles: ['admin_sekolah'],
  },
  {
    key: 'schedules',
    icon: ClockCircleOutlined,
    label: 'Jadwal Absensi',
    path: '/schedules',
    roles: ['admin_sekolah'],
  },
  {
    key: 'display-tokens',
    icon: DesktopOutlined,
    label: 'Display Token',
    path: '/display-tokens',
    roles: ['admin_sekolah'],
  },
  // Guru BK menus
  {
    key: 'bk-students',
    icon: TeamOutlined,
    label: 'Profil Siswa BK',
    path: '/bk/students',
    roles: ['guru_bk'],
  },
  {
    key: 'violations',
    icon: SafetyOutlined,
    label: 'Pelanggaran',
    path: '/bk/violations',
    roles: ['guru_bk'],
  },
  {
    key: 'achievements',
    icon: TrophyOutlined,
    label: 'Prestasi',
    path: '/bk/achievements',
    roles: ['guru_bk'],
  },
  {
    key: 'permits',
    icon: FileTextOutlined,
    label: 'Izin Keluar',
    path: '/bk/permits',
    roles: ['guru_bk'],
  },
  {
    key: 'counseling',
    icon: FormOutlined,
    label: 'Catatan Konseling',
    path: '/bk/counseling',
    roles: ['guru_bk'],
  },
  // Wali Kelas menus
  {
    key: 'class-attendance',
    icon: ScheduleOutlined,
    label: 'Absensi Kelas',
    path: '/homeroom/attendance',
    roles: ['wali_kelas'],
  },
  {
    key: 'grades',
    icon: BookOutlined,
    label: 'Input Nilai',
    path: '/homeroom/grades',
    roles: ['wali_kelas'],
  },
  {
    key: 'homeroom-notes',
    icon: FormOutlined,
    label: 'Catatan Wali Kelas',
    path: '/homeroom/notes',
    roles: ['wali_kelas'],
  },
  {
    key: 'bk-monitoring',
    icon: EyeOutlined,
    label: 'Monitoring BK',
    path: '/homeroom/bk-monitoring',
    roles: ['wali_kelas'],
    readOnly: true,
  },
]

// Update selected keys based on current route
watch(
  () => route.path,
  (path) => {
    // Find matching menu item by path
    const matchingItem = allMenuItems.find(item => item.path === path)
    if (matchingItem) {
      selectedKeys.value = [matchingItem.key]
    } else {
      // Fallback to first segment
      const key = path.split('/')[1] || 'dashboard'
      selectedKeys.value = [key]
    }
  },
  { immediate: true }
)

// Filter menu items based on user role
const menuItems = computed(() => {
  const role = authStore.userRole as UserRole | null
  if (!role) {
    return []
  }
  return allMenuItems.filter(item => item.roles.includes(role))
})

function handleMenuClick(info: MenuInfo) {
  const item = allMenuItems.find(m => m.key === String(info.key))
  if (item?.path) {
    router.push(item.path)
  }
}

function handleLogout() {
  authStore.clearAuth()
  router.push('/login')
}

function goToProfile() {
  router.push('/profile')
}

function goToSettings() {
  router.push('/account-settings')
}

// Notification functions
async function loadNotifications() {
  if (isSuperAdmin.value) return
  
  notificationLoading.value = true
  try {
    const response = await notificationService.getNotifications({ pageSize: 10, unreadOnly: false })
    notifications.value = response.data || []
    unreadCount.value = response.unreadCount || 0
  } catch {
    // Silently fail - notifications are not critical
    notifications.value = []
    unreadCount.value = 0
  } finally {
    notificationLoading.value = false
  }
}

async function loadUnreadCount() {
  if (isSuperAdmin.value) return
  
  try {
    unreadCount.value = await notificationService.getUnreadCount()
  } catch {
    // Silently fail
  }
}

async function markAsRead(notification: Notification) {
  if (notification.isRead) return
  
  try {
    await notificationService.markAsRead(notification.id)
    notification.isRead = true
    unreadCount.value = Math.max(0, unreadCount.value - 1)
  } catch {
    // Silently fail
  }
}

async function markAllAsRead() {
  try {
    await notificationService.markAllAsRead()
    notifications.value.forEach(n => n.isRead = true)
    unreadCount.value = 0
  } catch {
    // Silently fail
  }
}

function formatNotificationTime(dateStr: string): string {
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)
  
  if (diffMins < 1) return 'Baru saja'
  if (diffMins < 60) return diffMins + ' menit lalu'
  if (diffHours < 24) return diffHours + ' jam lalu'
  if (diffDays < 7) return diffDays + ' hari lalu'
  
  return date.toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })
}

function getNotificationIcon(type: string) {
  switch (type) {
    case 'attendance_checkin':
    case 'attendance_checkout':
      return ClockCircleOutlined
    case 'violation':
      return SafetyOutlined
    case 'achievement':
      return TrophyOutlined
    default:
      return BellOutlined
  }
}

// Get role display name
const roleDisplayName = computed(() => {
  return authStore.userRole ? getRoleDisplayName(authStore.userRole) : 'User'
})

// Get role color
const roleTagColor = computed(() => {
  return authStore.userRole ? getRoleColor(authStore.userRole) : 'default'
})

onMounted(() => {
  if (!isSuperAdmin.value) {
    loadNotifications()
    // Poll for new notifications every 60 seconds
    notificationInterval = setInterval(loadUnreadCount, 60000)
  }
})

onUnmounted(() => {
  if (notificationInterval) {
    clearInterval(notificationInterval)
  }
})
</script>

<template>
  <Layout class="main-layout">
    <LayoutSider
      v-model:collapsed="collapsed"
      :trigger="null"
      collapsible
      theme="light"
      class="sider"
      :width="240"
    >
      <div class="logo">
        <img src="/vite.svg" alt="Logo" class="logo-img" />
        <span v-if="!collapsed" class="logo-text">School Admin</span>
      </div>
      <Menu
        v-model:selectedKeys="selectedKeys"
        v-model:openKeys="openKeys"
        mode="inline"
        @click="handleMenuClick"
      >
        <MenuItem v-for="item in menuItems" :key="item.key">
          <component :is="item.icon" />
          <span>{{ item.label }}</span>
          <Tag v-if="item.readOnly" color="default" class="readonly-tag">
            <EyeOutlined /> Lihat
          </Tag>
        </MenuItem>
      </Menu>
    </LayoutSider>

    <Layout>
      <LayoutHeader class="header">
        <div class="header-left">
          <Button
            type="text"
            class="trigger-btn"
            @click="collapsed = !collapsed"
          >
            <MenuUnfoldOutlined v-if="collapsed" />
            <MenuFoldOutlined v-else />
          </Button>
        </div>

        <div class="header-right">
          <!-- Notification Bell (hidden for super_admin) -->
          <Dropdown 
            v-if="!isSuperAdmin"
            v-model:open="notificationDropdownVisible"
            placement="bottomRight"
            :trigger="['click']"
            @openChange="(visible: boolean) => visible && loadNotifications()"
          >
            <template #overlay>
              <div class="notification-dropdown">
                <div class="notification-header">
                  <Text strong>Notifikasi</Text>
                  <Button v-if="unreadCount > 0" type="link" size="small" @click="markAllAsRead">
                    <CheckOutlined /> Tandai semua dibaca
                  </Button>
                </div>
                <div class="notification-content">
                  <Spin v-if="notificationLoading" />
                  <Empty v-else-if="notifications.length === 0" description="Tidak ada notifikasi" :image="Empty.PRESENTED_IMAGE_SIMPLE" />
                  <List v-else item-layout="horizontal" :data-source="notifications" size="small">
                    <template #renderItem="{ item }">
                      <ListItem 
                        :class="['notification-item', { unread: !item.isRead }]"
                        @click="markAsRead(item)"
                      >
                        <ListItemMeta>
                          <template #avatar>
                            <Avatar :style="{ backgroundColor: item.isRead ? '#d9d9d9' : '#f97316' }" size="small">
                              <template #icon><component :is="getNotificationIcon(item.type)" /></template>
                            </Avatar>
                          </template>
                          <template #title>
                            <Text :strong="!item.isRead" style="font-size: 13px">{{ item.title }}</Text>
                          </template>
                          <template #description>
                            <div>
                              <Text type="secondary" style="font-size: 12px">{{ item.message }}</Text>
                              <br />
                              <Text type="secondary" style="font-size: 11px">{{ formatNotificationTime(item.createdAt) }}</Text>
                            </div>
                          </template>
                        </ListItemMeta>
                      </ListItem>
                    </template>
                  </List>
                </div>
              </div>
            </template>
            <Badge :count="unreadCount" :overflow-count="99" class="notification-badge">
              <Button type="text" class="icon-btn">
                <BellOutlined />
              </Button>
            </Badge>
          </Dropdown>

          <!-- User Dropdown -->
          <Dropdown placement="bottomRight">
            <template #overlay>
              <Menu>
                <MenuItem key="profile" @click="goToProfile">
                  <UserOutlined />
                  <span style="margin-left: 8px">Profil</span>
                </MenuItem>
                <MenuItem key="settings" @click="goToSettings">
                  <SettingOutlined />
                  <span style="margin-left: 8px">Pengaturan</span>
                </MenuItem>
                <Menu.Divider />
                <MenuItem key="logout" @click="handleLogout">
                  <LogoutOutlined />
                  <span style="margin-left: 8px">Logout</span>
                </MenuItem>
              </Menu>
            </template>
            <div class="user-info">
              <Avatar :size="36" class="user-avatar">
                <template #icon><UserOutlined /></template>
              </Avatar>
              <div v-if="!collapsed" class="user-details">
                <span class="username">{{ authStore.user?.username || 'User' }}</span>
                <Tag :color="roleTagColor" class="role-tag">{{ roleDisplayName }}</Tag>
              </div>
            </div>
          </Dropdown>
        </div>
      </LayoutHeader>

      <LayoutContent class="content">
        <router-view v-slot="{ Component, route }">
          <component :is="Component" :key="route.fullPath" />
        </router-view>
      </LayoutContent>
    </Layout>
  </Layout>
</template>

<style scoped>
.main-layout {
  min-height: 100vh;
}

.sider {
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.05);
  overflow: auto;
  height: 100vh;
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  z-index: 10;
}

.sider :deep(.ant-layout-sider-children) {
  display: flex;
  flex-direction: column;
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 0 16px;
  border-bottom: 1px solid #f0f0f0;
  background: #fff;
}

.logo-img {
  width: 32px;
  height: 32px;
}

.logo-text {
  font-size: 18px;
  font-weight: 600;
  color: #f97316;
  white-space: nowrap;
}

.header {
  background: #fff;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  position: sticky;
  top: 0;
  z-index: 9;
  margin-left: 240px;
  transition: margin-left 0.2s;
}

.sider.ant-layout-sider-collapsed + .ant-layout .header {
  margin-left: 80px;
}

.header-left {
  display: flex;
  align-items: center;
}

.trigger-btn {
  font-size: 18px;
  width: 48px;
  height: 48px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.notification-badge {
  margin-right: 8px;
}

.icon-btn {
  font-size: 18px;
  width: 40px;
  height: 40px;
}

.notification-dropdown {
  width: 360px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #f0f0f0;
}

.notification-content {
  max-height: 400px;
  overflow-y: auto;
  padding: 8px 0;
}

.notification-item {
  padding: 8px 16px !important;
  cursor: pointer;
  transition: background-color 0.2s;
}

.notification-item:hover {
  background-color: #f5f5f5;
}

.notification-item.unread {
  background-color: #fff7e6;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 8px;
  transition: background-color 0.2s;
}

.user-info:hover {
  background-color: #f5f5f5;
}

.user-avatar {
  background-color: #f97316;
}

.user-details {
  display: flex;
  flex-direction: column;
  line-height: 1.3;
}

.username {
  font-size: 14px;
  font-weight: 500;
  color: #262626;
}

.user-role {
  font-size: 12px;
  color: #8c8c8c;
}

.role-tag {
  font-size: 11px;
  padding: 0 6px;
  line-height: 18px;
  margin: 0;
}

.readonly-tag {
  font-size: 10px;
  padding: 0 4px;
  line-height: 16px;
  margin-left: 8px;
}

.content {
  margin: 24px;
  margin-left: calc(240px + 24px);
  padding: 24px;
  background: #fff;
  border-radius: 8px;
  min-height: calc(100vh - 64px - 48px);
  transition: margin-left 0.2s;
  position: relative;
  z-index: 1;
  overflow: auto;
}

/* Handle collapsed state */
:deep(.ant-layout-sider-collapsed) + .ant-layout .header {
  margin-left: 80px;
}

:deep(.ant-layout-sider-collapsed) + .ant-layout .content {
  margin-left: calc(80px + 24px);
}
</style>
