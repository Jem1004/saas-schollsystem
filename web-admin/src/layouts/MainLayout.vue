<script setup lang="ts">
import { ref, computed, watch } from 'vue'
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
} from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { usePermissions, getRoleDisplayName, getRoleColor } from '@/composables/usePermissions'
import type { UserRole } from '@/types/user'

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

// Update selected keys based on current route
watch(
  () => route.path,
  (path) => {
    const key = path.split('/')[1] || 'dashboard'
    selectedKeys.value = [key]
  },
  { immediate: true }
)

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
    key: 'attendance-admin',
    icon: ScheduleOutlined,
    label: 'Monitoring Absensi',
    path: '/attendance',
    roles: ['admin_sekolah'],
  },
  {
    key: 'school-settings',
    icon: SettingOutlined,
    label: 'Pengaturan Sekolah',
    path: '/settings',
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

// Filter menu items based on user role
const menuItems = computed(() => {
  const role = authStore.userRole as UserRole | null
  if (!role) {
    // Return empty menu if no role (should redirect to login)
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

// Get role display name
const roleDisplayName = computed(() => {
  return authStore.userRole ? getRoleDisplayName(authStore.userRole) : 'User'
})

// Get role color
const roleTagColor = computed(() => {
  return authStore.userRole ? getRoleColor(authStore.userRole) : 'default'
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
          <Badge :count="0" class="notification-badge">
            <Button type="text" class="icon-btn">
              <BellOutlined />
            </Button>
          </Badge>

          <Dropdown placement="bottomRight">
            <template #overlay>
              <Menu>
                <MenuItem key="profile">
                  <UserOutlined />
                  <span style="margin-left: 8px">Profile</span>
                </MenuItem>
                <MenuItem key="settings">
                  <SettingOutlined />
                  <span style="margin-left: 8px">Settings</span>
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
        <router-view />
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
}

/* Handle collapsed state */
:deep(.ant-layout-sider-collapsed) + .ant-layout .header {
  margin-left: 80px;
}

:deep(.ant-layout-sider-collapsed) + .ant-layout .content {
  margin-left: calc(80px + 24px);
}
</style>
