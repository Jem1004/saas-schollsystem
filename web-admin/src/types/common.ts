// Common/shared types

// Pagination
export interface PaginationParams {
  page?: number
  pageSize?: number
  search?: string
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  pageSize: number
}

// API Response wrapper
export interface ApiResponse<T> {
  success: boolean
  data: T
  message?: string
}

export interface ApiErrorResponse {
  success: false
  message: string
  code?: string
  details?: Record<string, unknown>
}

// Date range filter
export interface DateRangeParams {
  startDate?: string
  endDate?: string
}

// Sort params
export interface SortParams {
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}

// Combined filter params
export interface FilterParams extends PaginationParams, DateRangeParams, SortParams {}

// Select option type
export interface SelectOption<T = string | number> {
  value: T
  label: string
  disabled?: boolean
}

// Table column type
export interface TableColumn {
  key: string
  title: string
  dataIndex?: string
  width?: number | string
  align?: 'left' | 'center' | 'right'
  sortable?: boolean
  fixed?: 'left' | 'right'
}

// Confirmation dialog options
export interface ConfirmOptions {
  title: string
  content: string
  okText?: string
  cancelText?: string
  okType?: 'primary' | 'danger'
}

// Export options
export interface ExportOptions {
  format: 'xlsx' | 'csv' | 'pdf'
  filename?: string
  columns?: string[]
}
