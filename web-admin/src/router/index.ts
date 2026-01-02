import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw, RouteLocationNormalized } from 'vue-router'
import type { UserRole } from '@/types/user'

// Define route meta types
declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth?: boolean
    roles?: UserRole[]
    title?: string
  }
}

const routes: RouteRecordRaw[] = [
  // Public Display Route (no auth required)
  // Requirements: 5.3 - Accessing public display URL with valid token SHALL show attendance data without login
  {
    path: '/display/:token',
    name: 'PublicDisplay',
    component: () => import('@/views/public/PublicDisplay.vue'),
    meta: { requiresAuth: false, title: 'Public Display' },
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { requiresAuth: false, title: 'Login' },
  },
  {
    path: '/change-password',
    name: 'ChangePassword',
    component: () => import('@/views/auth/ChangePasswordView.vue'),
    meta: { requiresAuth: true, title: 'Ubah Password' },
  },
  {
    path: '/access-denied',
    name: 'AccessDenied',
    component: () => import('@/views/errors/AccessDeniedView.vue'),
    meta: { requiresAuth: false, title: 'Akses Ditolak' },
  },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/dashboard',
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/DashboardView.vue'),
        meta: { title: 'Dashboard' },
      },
      // Super Admin Routes
      {
        path: 'tenants',
        name: 'Tenants',
        component: () => import('@/views/super-admin/TenantManagement.vue'),
        meta: { title: 'Manajemen Sekolah', roles: ['super_admin'] },
      },
      {
        path: 'devices',
        name: 'Devices',
        component: () => import('@/views/super-admin/DeviceManagement.vue'),
        meta: { title: 'Manajemen Device', roles: ['super_admin'] },
      },
      // Admin Sekolah Routes
      {
        path: 'classes',
        name: 'Classes',
        component: () => import('@/views/admin-sekolah/ClassManagement.vue'),
        meta: { title: 'Manajemen Kelas', roles: ['admin_sekolah'] },
      },
      {
        path: 'students',
        name: 'Students',
        component: () => import('@/views/admin-sekolah/StudentManagement.vue'),
        meta: { title: 'Manajemen Siswa', roles: ['admin_sekolah'] },
      },
      {
        path: 'parents',
        name: 'Parents',
        component: () => import('@/views/admin-sekolah/ParentManagement.vue'),
        meta: { title: 'Manajemen Orang Tua', roles: ['admin_sekolah'] },
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/admin-sekolah/UserManagement.vue'),
        meta: { title: 'Manajemen User', roles: ['admin_sekolah'] },
      },
      {
        path: 'attendance',
        name: 'AttendanceReport',
        component: () => import('@/views/admin-sekolah/AttendanceReport.vue'),
        meta: { title: 'Laporan Absensi', roles: ['admin_sekolah', 'wali_kelas'] },
      },
      {
        path: 'attendance/live',
        name: 'LiveAttendance',
        component: () => import('@/views/admin-sekolah/LiveAttendance.vue'),
        meta: { title: 'Absensi Real-Time', roles: ['admin_sekolah', 'wali_kelas', 'guru_bk'] },
      },
      {
        path: 'settings',
        name: 'SchoolSettings',
        component: () => import('@/views/admin-sekolah/SchoolSettings.vue'),
        meta: { title: 'Pengaturan Sekolah', roles: ['admin_sekolah'], requiresAuth: true },
      },
      {
        path: 'schedules',
        name: 'ScheduleManagement',
        component: () => import('@/views/admin-sekolah/ScheduleManagement.vue'),
        meta: { title: 'Jadwal Absensi', roles: ['admin_sekolah'] },
      },
      {
        path: 'display-tokens',
        name: 'DisplayTokenManagement',
        component: () => import('@/views/admin-sekolah/DisplayTokenManagement.vue'),
        meta: { title: 'Display Token', roles: ['admin_sekolah'] },
      },
      // Guru BK Routes
      {
        path: 'bk',
        name: 'BKDashboard',
        component: () => import('@/views/guru-bk/GuruBKDashboard.vue'),
        meta: { title: 'Dashboard BK', roles: ['guru_bk'] },
      },
      {
        path: 'bk/students/:id',
        name: 'StudentBKProfile',
        component: () => import('@/views/guru-bk/StudentBKProfile.vue'),
        meta: { title: 'Profil Siswa BK', roles: ['guru_bk', 'wali_kelas'] },
      },
      {
        path: 'bk/violations',
        name: 'ViolationManagement',
        component: () => import('@/views/guru-bk/ViolationManagement.vue'),
        meta: { title: 'Manajemen Pelanggaran', roles: ['guru_bk'] },
      },
      {
        path: 'bk/achievements',
        name: 'AchievementManagement',
        component: () => import('@/views/guru-bk/AchievementManagement.vue'),
        meta: { title: 'Manajemen Prestasi', roles: ['guru_bk'] },
      },
      {
        path: 'bk/permits',
        name: 'PermitManagement',
        component: () => import('@/views/guru-bk/PermitManagement.vue'),
        meta: { title: 'Izin Keluar', roles: ['guru_bk'] },
      },
      {
        path: 'bk/counseling',
        name: 'CounselingManagement',
        component: () => import('@/views/guru-bk/CounselingManagement.vue'),
        meta: { title: 'Catatan Konseling', roles: ['guru_bk'] },
      },
      // Wali Kelas Routes
      {
        path: 'homeroom',
        name: 'HomeroomDashboard',
        component: () => import('@/views/wali-kelas/WaliKelasDashboard.vue'),
        meta: { title: 'Dashboard Wali Kelas', roles: ['wali_kelas'] },
      },
      {
        path: 'homeroom/attendance',
        name: 'ClassAttendance',
        component: () => import('@/views/wali-kelas/ClassAttendance.vue'),
        meta: { title: 'Absensi Kelas', roles: ['wali_kelas'] },
      },
      {
        path: 'homeroom/grades',
        name: 'GradeInput',
        component: () => import('@/views/wali-kelas/GradeInput.vue'),
        meta: { title: 'Input Nilai', roles: ['wali_kelas'] },
      },
      {
        path: 'homeroom/notes',
        name: 'HomeroomNotes',
        component: () => import('@/views/wali-kelas/HomeroomNotes.vue'),
        meta: { title: 'Catatan Wali Kelas', roles: ['wali_kelas'] },
      },
      {
        path: 'homeroom/bk-monitoring',
        name: 'BKMonitoring',
        component: () => import('@/views/wali-kelas/BKMonitoring.vue'),
        meta: { title: 'Monitoring BK', roles: ['wali_kelas'] },
      },
      // Account Routes (all authenticated users)
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/account/ProfileView.vue'),
        meta: { title: 'Profil Saya' },
      },
      {
        path: 'account-settings',
        name: 'AccountSettings',
        component: () => import('@/views/account/SettingsView.vue'),
        meta: { title: 'Pengaturan' },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/errors/NotFoundView.vue'),
    meta: { requiresAuth: false, title: 'Halaman Tidak Ditemukan' },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Helper function to check authentication
function isAuthenticated(): boolean {
  return !!localStorage.getItem('accessToken')
}

// Helper function to get user from storage
function getUser(): { role: UserRole; mustResetPwd: boolean } | null {
  const userStr = localStorage.getItem('user')
  if (!userStr) return null
  try {
    return JSON.parse(userStr)
  } catch {
    return null
  }
}

// Helper function to check if user has required role
function hasRequiredRole(userRole: UserRole, requiredRoles?: UserRole[]): boolean {
  if (!requiredRoles || requiredRoles.length === 0) return true
  return requiredRoles.includes(userRole)
}

// Route guard
router.beforeEach((to: RouteLocationNormalized, from: RouteLocationNormalized, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth !== false)
  const requiredRoles = to.meta.roles as UserRole[] | undefined
  
  // Update document title
  const title = to.meta.title as string | undefined
  document.title = title ? `${title} - Sistem Manajemen Sekolah` : 'Sistem Manajemen Sekolah'

  // Prevent redirect loops
  if (to.path === from.path) {
    return next()
  }

  // Check if route requires authentication
  if (requiresAuth) {
    if (!isAuthenticated()) {
      // Not authenticated - redirect to login
      return next({
        path: '/login',
        query: { redirect: to.fullPath },
      })
    }

    const user = getUser()
    
    // Check if user must change password
    if (user?.mustResetPwd && to.name !== 'ChangePassword') {
      return next({ path: '/change-password' })
    }

    // Check role-based access
    if (requiredRoles && user) {
      if (!hasRequiredRole(user.role, requiredRoles)) {
        return next({ path: '/access-denied' })
      }
    }
  } else {
    // Route doesn't require auth
    // If user is authenticated and trying to access login, redirect to dashboard
    if (to.name === 'Login' && isAuthenticated()) {
      const user = getUser()
      if (user?.mustResetPwd) {
        return next({ path: '/change-password' })
      }
      return next({ path: '/dashboard' })
    }
  }

  next()
})

export default router
