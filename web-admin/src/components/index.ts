// Shared components exports

// Table components
export { default as DataTable } from './DataTable.vue'

// Form components
export { default as StudentSelect } from './StudentSelect.vue'
export { default as ClassSelect } from './ClassSelect.vue'
export { default as TeacherSelect } from './TeacherSelect.vue'
export { default as DatePickerID } from './DatePickerID.vue'

// Notification components
export { default as NotificationDropdown } from './NotificationDropdown.vue'
export { default as ToastNotification } from './ToastNotification.vue'

// Document components
export { default as PermitDocumentPreview } from './PermitDocumentPreview.vue'

// UX Guards & Security components
export { default as SensitiveDataField } from './SensitiveDataField.vue'
export { default as ConfidentialBadge } from './ConfidentialBadge.vue'
export { default as ConfirmationDialog } from './ConfirmationDialog.vue'
export { default as ReadOnlyBanner } from './ReadOnlyBanner.vue'
export { default as ReadOnlyWrapper } from './ReadOnlyWrapper.vue'

// Empty & First-Time State components
export { default as EmptyState } from './EmptyState.vue'
export { default as WelcomeCard } from './WelcomeCard.vue'
export { default as FirstTimeHint } from './FirstTimeHint.vue'

// Error & Permission components
export { default as ErrorBoundary } from './ErrorBoundary.vue'
export { default as NetworkError } from './NetworkError.vue'
export { default as SessionExpired } from './SessionExpired.vue'
export { default as PermissionDenied } from './PermissionDenied.vue'

// Form Safety components
export { default as DestructiveActionDialog } from './DestructiveActionDialog.vue'
export { default as RequiredFieldIndicator } from './RequiredFieldIndicator.vue'

// Composables
export { useToast } from '@/composables/useToast'
export { useSensitiveData, maskSensitiveText, canViewSensitiveData } from '@/composables/useSensitiveData'
export { usePermissions, hasRoleAccess, getRoleDisplayName, getRoleColor } from '@/composables/usePermissions'
export { 
  useFormSafety, 
  sanitizeInput, 
  validateRequired, 
  validateEmail, 
  validatePhone, 
  validateNISN, 
  validatePassword, 
  validateLength, 
  validateNumberRange 
} from '@/composables/useFormSafety'
