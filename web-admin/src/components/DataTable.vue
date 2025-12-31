<script setup lang="ts" generic="T extends Record<string, unknown>">
import { ref, computed, watch } from 'vue'
import {
  Table,
  Button,
  Input,
  Space,
  Dropdown,
  Menu,
  MenuItem,
  Tooltip,
} from 'ant-design-vue'
import type { TableProps, TablePaginationConfig } from 'ant-design-vue'
import type { SorterResult, FilterValue } from 'ant-design-vue/es/table/interface'
import {
  SearchOutlined,
  ReloadOutlined,
  DownloadOutlined,
  ColumnHeightOutlined,
  FilterOutlined,
} from '@ant-design/icons-vue'

export interface DataTableColumn<R = unknown> {
  title: string
  dataIndex?: string
  key: string
  width?: number | string
  align?: 'left' | 'center' | 'right'
  sorter?: boolean | ((a: R, b: R) => number)
  sortDirections?: ('ascend' | 'descend')[]
  filters?: { text: string; value: string | number | boolean }[]
  filterMultiple?: boolean
  fixed?: 'left' | 'right' | boolean
  ellipsis?: boolean
  customRender?: (params: { text: unknown; record: R; index: number }) => unknown
}

export interface DataTableProps<R> {
  columns: DataTableColumn<R>[]
  dataSource: R[]
  loading?: boolean
  rowKey?: string | ((record: R) => string | number)
  pagination?: TablePaginationConfig | false
  total?: number
  showSearch?: boolean
  searchPlaceholder?: string
  showRefresh?: boolean
  showExport?: boolean
  showDensity?: boolean
  exportFileName?: string
  bordered?: boolean
  size?: 'large' | 'middle' | 'small'
  scroll?: { x?: number | string; y?: number | string }
  rowSelection?: TableProps['rowSelection']
}

const props = withDefaults(defineProps<DataTableProps<T>>(), {
  loading: false,
  rowKey: 'id',
  showSearch: true,
  searchPlaceholder: 'Cari...',
  showRefresh: true,
  showExport: true,
  showDensity: true,
  exportFileName: 'data-export',
  bordered: false,
  size: 'middle',
})

const emit = defineEmits<{
  (e: 'search', value: string): void
  (e: 'refresh'): void
  (e: 'change', pagination: TablePaginationConfig, filters: Record<string, FilterValue | null>, sorter: SorterResult<T> | SorterResult<T>[]): void
  (e: 'export', format: 'csv' | 'excel'): void
}>()

// Local state
const searchText = ref('')
const tableSize = ref<'large' | 'middle' | 'small'>(props.size)

// Computed pagination config
const paginationConfig = computed<TablePaginationConfig | false>(() => {
  if (props.pagination === false) return false
  
  return {
    showSizeChanger: true,
    showQuickJumper: true,
    showTotal: (total: number, range: [number, number]) => 
      `${range[0]}-${range[1]} dari ${total} data`,
    pageSizeOptions: ['10', '20', '50', '100'],
    ...props.pagination,
    total: props.total ?? props.pagination?.total,
  }
})

// Handle search
const handleSearch = () => {
  emit('search', searchText.value)
}

// Handle search on enter
const handleSearchEnter = () => {
  handleSearch()
}

// Handle refresh
const handleRefresh = () => {
  emit('refresh')
}

// Handle table change
const handleTableChange = (
  pagination: TablePaginationConfig,
  filters: Record<string, FilterValue | null>,
  sorter: SorterResult<T> | SorterResult<T>[]
) => {
  emit('change', pagination, filters, sorter)
}

// Handle export
const handleExport = (format: 'csv' | 'excel') => {
  emit('export', format)
}

// Export to CSV
const exportToCSV = () => {
  const headers = props.columns
    .filter(col => col.dataIndex)
    .map(col => col.title)
  
  const rows = props.dataSource.map(record => 
    props.columns
      .filter(col => col.dataIndex)
      .map(col => {
        const value = col.dataIndex ? record[col.dataIndex as keyof T] : ''
        return typeof value === 'string' ? `"${value.replace(/"/g, '""')}"` : value
      })
  )
  
  const csvContent = [
    headers.join(','),
    ...rows.map(row => row.join(','))
  ].join('\n')
  
  const blob = new Blob(['\ufeff' + csvContent], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = `${props.exportFileName}-${new Date().toISOString().split('T')[0]}.csv`
  link.click()
  URL.revokeObjectURL(link.href)
}

// Change table density
const handleDensityChange = (size: 'large' | 'middle' | 'small') => {
  tableSize.value = size
}

// Watch for external size changes
watch(() => props.size, (newSize) => {
  tableSize.value = newSize
})

// Expose methods
defineExpose({
  exportToCSV,
})
</script>

<template>
  <div class="data-table">
    <!-- Toolbar -->
    <div class="data-table-toolbar">
      <div class="toolbar-left">
        <slot name="toolbar-left">
          <Input
            v-if="showSearch"
            v-model:value="searchText"
            :placeholder="searchPlaceholder"
            allow-clear
            class="search-input"
            @press-enter="handleSearchEnter"
            @change="handleSearch"
          >
            <template #prefix>
              <SearchOutlined />
            </template>
          </Input>
        </slot>
        <slot name="toolbar-filters" />
      </div>
      
      <div class="toolbar-right">
        <slot name="toolbar-right" />
        <Space>
          <Tooltip v-if="showRefresh" title="Refresh">
            <Button @click="handleRefresh">
              <template #icon><ReloadOutlined /></template>
            </Button>
          </Tooltip>
          
          <Dropdown v-if="showExport">
            <template #overlay>
              <Menu>
                <MenuItem key="csv" @click="exportToCSV">
                  Export CSV
                </MenuItem>
                <MenuItem key="excel" @click="handleExport('excel')">
                  Export Excel
                </MenuItem>
              </Menu>
            </template>
            <Tooltip title="Export">
              <Button>
                <template #icon><DownloadOutlined /></template>
              </Button>
            </Tooltip>
          </Dropdown>
          
          <Dropdown v-if="showDensity">
            <template #overlay>
              <Menu :selected-keys="[tableSize]">
                <MenuItem key="large" @click="handleDensityChange('large')">
                  Besar
                </MenuItem>
                <MenuItem key="middle" @click="handleDensityChange('middle')">
                  Sedang
                </MenuItem>
                <MenuItem key="small" @click="handleDensityChange('small')">
                  Kecil
                </MenuItem>
              </Menu>
            </template>
            <Tooltip title="Ukuran Tabel">
              <Button>
                <template #icon><ColumnHeightOutlined /></template>
              </Button>
            </Tooltip>
          </Dropdown>
        </Space>
      </div>
    </div>
    
    <!-- Table -->
    <Table
      :columns="columns as TableProps['columns']"
      :data-source="dataSource"
      :loading="loading"
      :row-key="rowKey"
      :pagination="paginationConfig"
      :bordered="bordered"
      :size="tableSize"
      :scroll="scroll"
      :row-selection="rowSelection"
      @change="handleTableChange"
    >
      <template #bodyCell="slotProps">
        <slot name="bodyCell" v-bind="slotProps" />
      </template>
      
      <template #headerCell="slotProps">
        <slot name="headerCell" v-bind="slotProps" />
      </template>
      
      <template #emptyText>
        <slot name="empty">
          <div class="empty-state">
            <FilterOutlined class="empty-icon" />
            <p>Tidak ada data</p>
          </div>
        </slot>
      </template>
    </Table>
  </div>
</template>

<style scoped>
.data-table {
  width: 100%;
}

.data-table-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 12px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.search-input {
  width: 250px;
}

.empty-state {
  padding: 32px;
  text-align: center;
  color: #8c8c8c;
}

.empty-icon {
  font-size: 48px;
  color: #d9d9d9;
  margin-bottom: 16px;
}

@media (max-width: 768px) {
  .data-table-toolbar {
    flex-direction: column;
    align-items: stretch;
  }
  
  .toolbar-left,
  .toolbar-right {
    width: 100%;
    justify-content: flex-start;
  }
  
  .search-input {
    width: 100%;
  }
}
</style>
