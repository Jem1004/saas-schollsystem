<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import {
  Table,
  Button,
  Input,
  Space,
  Modal,
  Form,
  FormItem,
  Select,
  SelectOption,
  message,
  Popconfirm,
  Card,
  Row,
  Col,
  Typography,
  Textarea,
  DatePicker,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  PlusOutlined,
  SearchOutlined,
  DeleteOutlined,
  ReloadOutlined,
  FilterOutlined,
  EyeOutlined,
  SettingOutlined,
  FilePdfOutlined,
} from '@ant-design/icons-vue'
import { useRouter } from 'vue-router'
import dayjs from 'dayjs'
import type { Dayjs } from 'dayjs'
import { bkService, schoolService } from '@/services'
import type { Violation, ViolationCategory, CreateViolationRequest } from '@/types/bk'
import type { Student } from '@/types/school'
import { VIOLATION_LEVELS } from '@/types/bk'
import { exportToPDF, formatViolationForExport } from '@/utils/pdfExport'

const { Title, Text } = Typography
const { RangePicker } = DatePicker

const router = useRouter()

// Table state
const loading = ref(false)
const violations = ref<Violation[]>([])
const total = ref(0)
const pagination = reactive({
  current: 1,
  pageSize: 10,
})
const searchText = ref('')
const filterLevel = ref<string | undefined>(undefined)
const filterCategory = ref<string | undefined>(undefined)
const dateRange = ref<[Dayjs, Dayjs] | undefined>(undefined)

// Students for dropdown
const students = ref<Student[]>([])
const loadingStudents = ref(false)

// Violation categories from API
const violationCategories = ref<ViolationCategory[]>([])
const loadingCategories = ref(false)

// Modal state
const modalVisible = ref(false)
const modalLoading = ref(false)

// Form state
const formRef = ref()
const formState = reactive<CreateViolationRequest>({
  studentId: 0,
  categoryId: undefined,
  category: '',
  level: 'ringan',
  point: -5,
  description: '',
})

// Form rules
const formRules = {
  studentId: [{ required: true, message: 'Siswa wajib dipilih' }],
  category: [{ required: true, message: 'Kategori wajib dipilih' }],
  description: [{ required: true, message: 'Deskripsi wajib diisi' }],
}

// Table columns
const columns: TableProps['columns'] = [
  { title: 'Tanggal', dataIndex: 'createdAt', key: 'createdAt', width: 110, sorter: true },
  { title: 'Siswa', dataIndex: 'studentName', key: 'studentName' },
  { title: 'Kelas', dataIndex: 'studentClass', key: 'studentClass', width: 90 },
  { title: 'Kategori', dataIndex: 'category', key: 'category', width: 130 },
  { title: 'Tingkat', dataIndex: 'level', key: 'level', width: 90, align: 'center' },
  { title: 'Poin', dataIndex: 'point', key: 'point', width: 80, align: 'center' },
  { title: 'Deskripsi', dataIndex: 'description', key: 'description', ellipsis: true },
  { title: 'Aksi', key: 'action', width: 100, align: 'center' },
]

// Computed filtered data
const filteredViolations = computed(() => {
  let result = violations.value
  if (filterLevel.value) result = result.filter(v => v.level === filterLevel.value)
  if (filterCategory.value) result = result.filter(v => v.category === filterCategory.value)
  if (dateRange.value) {
    const [start, end] = dateRange.value
    result = result.filter(v => {
      const date = dayjs(v.createdAt)
      return date.isAfter(start.startOf('day')) && date.isBefore(end.endOf('day'))
    })
  }
  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(v =>
      v.studentName?.toLowerCase().includes(search) ||
      v.description.toLowerCase().includes(search) ||
      v.category.toLowerCase().includes(search)
    )
  }
  return result
})

// Total points summary
const totalPoints = computed(() => {
  return filteredViolations.value.reduce((sum, v) => sum + (v.point || 0), 0)
})

// Category names for filter dropdown
const categoryNames = computed(() => {
  return violationCategories.value.map(c => c.name)
})

// Format date
const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

// Get level class
const getLevelClass = (level: string) => {
  switch (level) {
    case 'ringan': return 'success'
    case 'sedang': return 'warning'
    case 'berat': return 'error'
    case 'sangat_berat': return 'error'
    default: return 'default'
  }
}

// Get level label
const getLevelLabel = (level: string) => {
  const levelConfig = VIOLATION_LEVELS.find(l => l.value === level)
  return levelConfig?.label || level
}

// Load violation categories
const loadCategories = async () => {
  loadingCategories.value = true
  try {
    const response = await bkService.getViolationCategories(true)
    violationCategories.value = response.categories || []
    // If no categories, initialize defaults
    if (violationCategories.value.length === 0) {
      await bkService.initializeDefaultCategories()
      const retryResponse = await bkService.getViolationCategories(true)
      violationCategories.value = retryResponse.categories || []
    }
  } catch (err) {
    console.error('Failed to load categories:', err)
    violationCategories.value = []
  } finally {
    loadingCategories.value = false
  }
}

// Load violations data
const loadViolations = async () => {
  loading.value = true
  try {
    const response = await bkService.getViolations({
      page: pagination.current,
      pageSize: pagination.pageSize,
      search: searchText.value,
      level: filterLevel.value,
      category: filterCategory.value,
    })
    violations.value = response.data || []
    total.value = response.total || 0
  } catch (err) {
    console.error('Failed to load violations:', err)
    violations.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// Load students for dropdown
const loadStudents = async () => {
  loadingStudents.value = true
  try {
    const response = await schoolService.getStudents({ page_size: 1000 })
    students.value = response.students || []
  } catch (err) {
    console.error('Failed to load students:', err)
    students.value = []
  } finally {
    loadingStudents.value = false
  }
}

// Handle category selection - auto-fill level and point
const handleCategoryChange = (categoryName: unknown) => {
  const category = violationCategories.value.find(c => c.name === categoryName)
  if (category) {
    formState.categoryId = category.id
    formState.level = category.defaultLevel
    formState.point = category.defaultPoint
  }
}

// Handle table change
const handleTableChange: TableProps['onChange'] = (pag) => {
  pagination.current = pag.current || 1
  pagination.pageSize = pag.pageSize || 10
  loadViolations()
}

// Handle search
const handleSearch = () => {
  pagination.current = 1
  loadViolations()
}

// Handle filter change
const handleFilterChange = () => {
  pagination.current = 1
  loadViolations()
}

// Open create modal
const openCreateModal = () => {
  resetForm()
  modalVisible.value = true
}

// Reset form
const resetForm = () => {
  formState.studentId = 0
  formState.categoryId = undefined
  formState.category = ''
  formState.level = 'ringan'
  formState.point = -5
  formState.description = ''
  formRef.value?.resetFields()
}

// Handle modal cancel
const handleModalCancel = () => {
  modalVisible.value = false
  resetForm()
}

// Handle form submit
const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  modalLoading.value = true
  try {
    await bkService.createViolation(formState)
    message.success('Pelanggaran berhasil dicatat')
    modalVisible.value = false
    resetForm()
    loadViolations()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Terjadi kesalahan')
  } finally {
    modalLoading.value = false
  }
}

// Handle delete
const handleDelete = async (violation: Violation) => {
  try {
    await bkService.deleteViolation(violation.id)
    message.success('Pelanggaran berhasil dihapus')
    loadViolations()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Gagal menghapus pelanggaran')
  }
}

// View student profile
const viewStudentProfile = (studentId: number) => {
  router.push(`/bk/students/${studentId}`)
}

// Go to category management
const goToCategoryManagement = () => {
  router.push('/bk/violation-categories')
}

// Export to PDF
const handleExportPDF = () => {
  if (filteredViolations.value.length === 0) {
    message.warning('Tidak ada data untuk diekspor')
    return
  }

  const dateRangeStr = dateRange.value
    ? { start: dateRange.value[0].format('DD/MM/YYYY'), end: dateRange.value[1].format('DD/MM/YYYY') }
    : undefined

  exportToPDF({
    title: 'Laporan Data Pelanggaran Siswa',
    filename: `laporan-pelanggaran-${dayjs().format('YYYY-MM-DD')}`,
    columns: [
      { header: 'Tanggal', dataKey: 'createdAt' },
      { header: 'Siswa', dataKey: 'studentName' },
      { header: 'Kelas', dataKey: 'studentClass' },
      { header: 'Kategori', dataKey: 'category' },
      { header: 'Tingkat', dataKey: 'level' },
      { header: 'Poin', dataKey: 'point' },
      { header: 'Deskripsi', dataKey: 'description' },
    ],
    data: filteredViolations.value.map(v => formatViolationForExport(v as unknown as Record<string, unknown>)),
    dateRange: dateRangeStr,
  })
  message.success('PDF berhasil diunduh')
}

// Filter student options
const filterStudentOption = (input: string, option: any) => {
  return option?.label?.toLowerCase().includes(input.toLowerCase())
}

onMounted(() => {
  loadCategories()
  loadViolations()
  loadStudents()
})
</script>

<template>
  <div class="violation-management">
    <div class="page-header">
      <Row justify="space-between" align="middle">
        <Col><Title :level="2" style="margin: 0">Manajemen Pelanggaran</Title></Col>
        <Col>
          <Button @click="goToCategoryManagement">
            <template #icon><SettingOutlined /></template>
            Kelola Kategori
          </Button>
        </Col>
      </Row>
    </div>

    <Card>
      <!-- Toolbar -->
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col :xs="24" :sm="24" :md="18">
          <Space wrap>
            <Input v-model:value="searchText" placeholder="Cari siswa atau deskripsi..." allow-clear style="width: 220px" @press-enter="handleSearch">
              <template #prefix><SearchOutlined /></template>
            </Input>
            <RangePicker v-model:value="dateRange" format="DD/MM/YYYY" :placeholder="['Dari Tanggal', 'Sampai Tanggal']" style="width: 250px" @change="handleFilterChange" />
            <Select v-model:value="filterLevel" placeholder="Filter Tingkat" allow-clear style="width: 140px" @change="handleFilterChange">
              <template #suffixIcon><FilterOutlined /></template>
              <SelectOption v-for="level in VIOLATION_LEVELS" :key="level.value" :value="level.value">{{ level.label }}</SelectOption>
            </Select>
            <Select v-model:value="filterCategory" placeholder="Filter Kategori" allow-clear style="width: 160px" @change="handleFilterChange">
              <SelectOption v-for="cat in categoryNames" :key="cat" :value="cat">{{ cat }}</SelectOption>
            </Select>
          </Space>
        </Col>
        <Col :xs="24" :sm="24" :md="6" class="toolbar-right">
          <Space>
            <Text type="secondary">Total Poin: <Text strong :style="{ color: '#ef4444' }">{{ totalPoints }}</Text></Text>
            <Button @click="handleExportPDF"><template #icon><FilePdfOutlined /></template>Export PDF</Button>
            <Button @click="loadViolations"><template #icon><ReloadOutlined /></template></Button>
            <Button type="primary" @click="openCreateModal"><template #icon><PlusOutlined /></template>Catat</Button>
          </Space>
        </Col>
      </Row>

      <!-- Table -->
      <Table 
        :columns="columns" 
        :data-source="filteredViolations" 
        :loading="loading"
        :pagination="{ current: pagination.current, pageSize: pagination.pageSize, total, showSizeChanger: true, showTotal: (t: number) => `Total ${t} pelanggaran` }"
        row-key="id" 
        @change="handleTableChange"
        class="custom-table"
        :scroll="{ x: 800 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'createdAt'">{{ formatDate((record as Violation).createdAt) }}</template>
          <template v-else-if="column.key === 'studentName'">
            <a @click="viewStudentProfile((record as Violation).studentId)">{{ (record as Violation).studentName }}</a>
          </template>
          <template v-else-if="column.key === 'studentClass'">
            <span class="class-badge">{{ (record as Violation).studentClass }}</span>
          </template>
          <template v-else-if="column.key === 'level'">
            <span :class="['status-badge', getLevelClass((record as Violation).level)]">{{ getLevelLabel((record as Violation).level) }}</span>
          </template>
          <template v-else-if="column.key === 'point'">
            <Text strong :style="{ color: '#ef4444' }">{{ (record as Violation).point }}</Text>
          </template>
          <template v-else-if="column.key === 'action'">
            <Space>
              <Button type="text" style="color: #3b82f6" @click="viewStudentProfile((record as Violation).studentId)"><template #icon><EyeOutlined /></template></Button>
              <Popconfirm title="Hapus pelanggaran ini?" description="Data pelanggaran akan dihapus permanen." ok-text="Ya, Hapus" cancel-text="Batal" @confirm="handleDelete(record as Violation)">
                <Button type="text" danger><template #icon><DeleteOutlined /></template></Button>
              </Popconfirm>
            </Space>
          </template>
        </template>
      </Table>
    </Card>

    <!-- Create Modal -->
    <Modal v-model:open="modalVisible" title="Catat Pelanggaran Baru" :confirm-loading="modalLoading" @ok="handleSubmit" @cancel="handleModalCancel" width="600px" wrap-class-name="modern-modal">
      <Form ref="formRef" :model="formState" :rules="formRules" layout="vertical" style="margin-top: 16px">
        <FormItem label="Siswa" name="studentId" required>
          <Select v-model:value="formState.studentId" placeholder="Pilih siswa" :loading="loadingStudents" show-search :filter-option="filterStudentOption"
            :options="students.map(s => ({ value: s.id, label: `${s.name} (${s.className})` }))" />
        </FormItem>
        <FormItem label="Kategori Pelanggaran" name="category" required>
          <Select v-model:value="formState.category" placeholder="Pilih kategori pelanggaran" :loading="loadingCategories" @change="handleCategoryChange">
            <SelectOption v-for="cat in violationCategories" :key="cat.id" :value="cat.name">
              {{ cat.name }}
            </SelectOption>
          </Select>
        </FormItem>
        <Row :gutter="16" v-if="formState.categoryId">
          <Col :span="12">
            <FormItem label="Tingkat">
              <span :class="['status-badge', getLevelClass(formState.level)]">{{ getLevelLabel(formState.level) }}</span>
            </FormItem>
          </Col>
          <Col :span="12">
            <FormItem label="Poin">
              <Text strong :style="{ color: '#ef4444', fontSize: '16px' }">{{ formState.point }} poin</Text>
            </FormItem>
          </Col>
        </Row>
        <FormItem label="Deskripsi" name="description" required>
          <Textarea v-model:value="formState.description" placeholder="Jelaskan detail pelanggaran..." :rows="4" class="custom-textarea" />
        </FormItem>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.violation-management { padding: 0; }
.page-header { margin-bottom: 24px; }
.toolbar { margin-bottom: 16px; }
.toolbar-right { display: flex; justify-content: flex-end; align-items: center; gap: 16px; }

/* Custom Table Styles */
.custom-table :deep(.ant-table-thead > tr > th) {
  background: #fafafa;
  font-weight: 600;
  color: #475569;
}

.custom-table :deep(.ant-table-tbody > tr > td) {
  padding: 16px;
}

.custom-table :deep(.ant-table-tbody > tr:hover > td) {
  background: #f8fafc;
}

/* Badge Styles */
.status-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap;
}

.status-badge.success {
  background-color: #f6ffed;
  color: #52c41a;
  border: 1px solid #b7eb8f;
}

.status-badge.warning {
  background-color: #fffbe6;
  color: #faad14;
  border: 1px solid #ffe58f;
}

.status-badge.error {
  background-color: #fff2f0;
  color: #ff4d4f;
  border: 1px solid #ffccc7;
}

.status-badge.default {
  background-color: #f5f5f5;
  color: #000000d9;
  border: 1px solid #d9d9d9;
}

.class-badge {
  color: #1890ff;
  font-size: 12px;
}

@media (max-width: 768px) { .toolbar-right { margin-top: 16px; justify-content: flex-start; flex-wrap: wrap; } }
</style>
