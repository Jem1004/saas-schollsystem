import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import type { UserRole } from '@/types/user'

export interface Permission {
  action: 'view' | 'create' | 'edit' | 'delete' | 'export'
  resource: string
}

// Define role-based permissions
const rolePermissions: Record<UserRole, Permission[]> = {
  super_admin: [
    { action: 'view', resource: 'tenants' },
    { action: 'create', resource: 'tenants' },
    { action: 'edit', resource: 'tenants' },
    { action: 'delete', resource: 'tenants' },
    { action: 'view', resource: 'devices' },
    { action: 'create', resource: 'devices' },
    { action: 'edit', resource: 'devices' },
    { action: 'delete', resource: 'devices' },
    { action: 'view', resource: 'all_schools' },
    { action: 'export', resource: 'system_reports' },
  ],
  admin_sekolah: [
    { action: 'view', resource: 'classes' },
    { action: 'create', resource: 'classes' },
    { action: 'edit', resource: 'classes' },
    { action: 'delete', resource: 'classes' },
    { action: 'view', resource: 'students' },
    { action: 'create', resource: 'students' },
    { action: 'edit', resource: 'students' },
    { action: 'delete', resource: 'students' },
    { action: 'view', resource: 'parents' },
    { action: 'create', resource: 'parents' },
    { action: 'edit', resource: 'parents' },
    { action: 'view', resource: 'users' },
    { action: 'create', resource: 'users' },
    { action: 'edit', resource: 'users' },
    { action: 'delete', resource: 'users' },
    { action: 'view', resource: 'attendance' },
    { action: 'export', resource: 'attendance' },
    { action: 'view', resource: 'school_settings' },
    { action: 'edit', resource: 'school_settings' },
  ],
  guru_bk: [
    { action: 'view', resource: 'students' },
    { action: 'view', resource: 'violations' },
    { action: 'create', resource: 'violations' },
    { action: 'edit', resource: 'violations' },
    { action: 'delete', resource: 'violations' },
    { action: 'view', resource: 'achievements' },
    { action: 'create', resource: 'achievements' },
    { action: 'edit', resource: 'achievements' },
    { action: 'delete', resource: 'achievements' },
    { action: 'view', resource: 'permits' },
    { action: 'create', resource: 'permits' },
    { action: 'edit', resource: 'permits' },
    { action: 'view', resource: 'counseling_notes' },
    { action: 'create', resource: 'counseling_notes' },
    { action: 'edit', resource: 'counseling_notes' },
    { action: 'delete', resource: 'counseling_notes' },
    { action: 'view', resource: 'internal_notes' },
    { action: 'export', resource: 'bk_reports' },
  ],
  wali_kelas: [
    { action: 'view', resource: 'class_students' },
    { action: 'view', resource: 'class_attendance' },
    { action: 'create', resource: 'manual_attendance' },
    { action: 'view', resource: 'grades' },
    { action: 'create', resource: 'grades' },
    { action: 'edit', resource: 'grades' },
    { action: 'delete', resource: 'grades' },
    { action: 'view', resource: 'homeroom_notes' },
    { action: 'create', resource: 'homeroom_notes' },
    { action: 'edit', resource: 'homeroom_notes' },
    { action: 'delete', resource: 'homeroom_notes' },
    { action: 'view', resource: 'bk_monitoring' }, // Read-only BK access
    { action: 'view', resource: 'violations' }, // Read-only
    { action: 'view', resource: 'achievements' }, // Read-only
    { action: 'view', resource: 'permits' }, // Read-only
    { action: 'view', resource: 'parent_summary' }, // Only parent summary, not internal notes
  ],
}

// Resources that are read-only for certain roles
const readOnlyResources: Record<UserRole, string[]> = {
  super_admin: [],
  admin_sekolah: [],
  guru_bk: [],
  wali_kelas: ['violations', 'achievements', 'permits', 'bk_monitoring', 'parent_summary'],
}

export function usePermissions() {
  const authStore = useAuthStore()

  const userRole = computed(() => authStore.userRole as UserRole | null)

  const hasPermission = (action: Permission['action'], resource: string): boolean => {
    if (!userRole.value) return false
    
    const permissions = rolePermissions[userRole.value] || []
    return permissions.some(p => p.action === action && p.resource === resource)
  }

  const canView = (resource: string): boolean => hasPermission('view', resource)
  const canCreate = (resource: string): boolean => hasPermission('create', resource)
  const canEdit = (resource: string): boolean => hasPermission('edit', resource)
  const canDelete = (resource: string): boolean => hasPermission('delete', resource)
  const canExport = (resource: string): boolean => hasPermission('export', resource)

  const isReadOnly = (resource: string): boolean => {
    if (!userRole.value) return true
    return readOnlyResources[userRole.value]?.includes(resource) || false
  }

  const getPermissionsForResource = (resource: string): Permission['action'][] => {
    if (!userRole.value) return []
    
    const permissions = rolePermissions[userRole.value] || []
    return permissions
      .filter(p => p.resource === resource)
      .map(p => p.action)
  }

  const isSuperAdmin = computed(() => userRole.value === 'super_admin')
  const isAdminSekolah = computed(() => userRole.value === 'admin_sekolah')
  const isGuruBK = computed(() => userRole.value === 'guru_bk')
  const isWaliKelas = computed(() => userRole.value === 'wali_kelas')

  // Check if user can access internal counseling notes
  const canViewInternalNotes = computed(() => hasPermission('view', 'internal_notes'))

  // Check if user can only see parent summary (not internal notes)
  const canOnlyViewParentSummary = computed(() => {
    return hasPermission('view', 'parent_summary') && !hasPermission('view', 'internal_notes')
  })

  return {
    userRole,
    hasPermission,
    canView,
    canCreate,
    canEdit,
    canDelete,
    canExport,
    isReadOnly,
    getPermissionsForResource,
    isSuperAdmin,
    isAdminSekolah,
    isGuruBK,
    isWaliKelas,
    canViewInternalNotes,
    canOnlyViewParentSummary,
  }
}

// Utility function to check role access
export function hasRoleAccess(userRole: UserRole | null, allowedRoles: UserRole[]): boolean {
  if (!userRole) return false
  return allowedRoles.includes(userRole)
}

// Get role display name in Indonesian
export function getRoleDisplayName(role: UserRole): string {
  const roleNames: Record<UserRole, string> = {
    super_admin: 'Super Admin',
    admin_sekolah: 'Admin Sekolah',
    guru_bk: 'Guru BK',
    wali_kelas: 'Wali Kelas',
  }
  return roleNames[role] || role
}

// Get role color for UI
export function getRoleColor(role: UserRole): string {
  const roleColors: Record<UserRole, string> = {
    super_admin: 'purple',
    admin_sekolah: 'blue',
    guru_bk: 'orange',
    wali_kelas: 'green',
  }
  return roleColors[role] || 'default'
}
