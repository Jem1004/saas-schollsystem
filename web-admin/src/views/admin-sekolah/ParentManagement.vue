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
  message,
  Popconfirm,
  Card,
  Row,
  Col,
  Typography,
  Tooltip,
  Alert,
  Upload,
  Divider,
  Progress,
} from 'ant-design-vue'
import type { TableProps, UploadProps } from 'ant-design-vue'
import {
  PlusOutlined,
  SearchOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
  UserOutlined,
  LinkOutlined,
  KeyOutlined,
  CopyOutlined,
  InfoCircleOutlined,
  DownloadOutlined,
  UploadOutlined,
  FileExcelOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  WarningOutlined,
} from '@ant-design/icons-vue'
import { schoolService, importService } from '@/services'
import type { StudentSearchResult } from '@/services/school'
import type { ImportResult, ImportError, ImportWarning } from '@/services/import'
import type { Parent, Student, UpdateParentRequest } from '@/types/school'

const { Title, Text } = Typography

// Table state
const loading = ref(false)
const parents = ref<Parent[]>([])
const total = ref(0)
const pagination = reactive({
  current: 1,
  pageSize: 10,
})
const searchText = ref('')

// Students for linking
const students = ref<Student[]>([])
const loadingStudents = ref(false)

// Student search state for enhanced linking
const studentSearchQuery = ref('')
const studentSearchResults = ref<StudentSearchResult[]>([])
const searchingStudents = ref(false)
const selectedStudentsForLinking = ref<StudentSearchResult[]>([])
let searchDebounceTimer: ReturnType<typeof setTimeout> | null = null

// Modal state
const modalVisible = ref(false)
const modalLoading = ref(false)
const isEditing = ref(false)
const editingParent = ref<Parent | null>(null)

// Credential modal state
const credentialModalVisible = ref(false)
const credentialData = ref<{ username: string; password: string; name: string } | null>(null)

// Import modal state
const importModalVisible = ref(false)
const importLoading = ref(false)
const importResultModalVisible = ref(false)
const importResult = ref<ImportResult | null>(null)

// Form state
const formRef = ref()
const formState = reactive({
  name: '',
  phone: '',
  email: '',
  student_ids: [] as number[],
})

// Form rules
const formRules = {
  name: [{ required: true, message: 'Nama orang tua wajib diisi' }],
  phone: [{ required: true, message: 'Nomor telepon wajib diisi (digunakan sebagai username)' }],
}

// Table columns
const columns: TableProps['columns'] = [
  {
    title: 'Nama Orang Tua',
    dataIndex: 'name',
    key: 'name',
    sorter: true,
  },
  {
    title: 'Username (No. HP)',
    dataIndex: 'phone',
    key: 'phone',
    width: 160,
  },
  {
    title: 'Email',
    dataIndex: 'email',
    key: 'email',
  },
  {
    title: 'Anak',
    dataIndex: 'studentNames',
    key: 'studentNames',
  },
  {
    title: 'Aksi',
    key: 'action',
    width: 200,
    align: 'center',
  },
]

// Computed filtered data
const filteredParents = computed(() => {
  if (!searchText.value) return parents.value
  const search = searchText.value.toLowerCase()
  return parents.value.filter(
    (parent) =>
      parent.name.toLowerCase().includes(search) ||
      parent.phone?.includes(search) ||
      parent.email?.toLowerCase().includes(search) ||
      parent.studentNames?.some(name => name.toLowerCase().includes(search))
  )
})

// Load parents data
const loadParents = async () => {
  loading.value = true
  try {
    const response = await schoolService.getParents({
      page: pagination.current,
      page_size: pagination.pageSize,
      search: searchText.value,
    })
    parents.value = response.parents
    total.value = response.pagination.total
  } catch (err) {
    console.error('Failed to load parents:', err)
    message.error('Gagal memuat data orang tua')
    parents.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// Load students for dropdown
const loadStudents = async () => {
  loadingStudents.value = true
  try {
    const response = await schoolService.getStudents({ page_size: 500 })
    students.value = response.students.filter(s => s.isActive)
  } catch (err) {
    console.error('Failed to load students:', err)
    students.value = []
  } finally {
    loadingStudents.value = false
  }
}

// ==================== Student Search Functions ====================

// Search students by NISN or name
const handleStudentSearch = (query: string) => {
  studentSearchQuery.value = query
  
  // Clear previous timer
  if (searchDebounceTimer) {
    clearTimeout(searchDebounceTimer)
  }
  
  // Debounce search
  searchDebounceTimer = setTimeout(async () => {
    if (!query || query.length < 2) {
      studentSearchResults.value = []
      return
    }
    
    searchingStudents.value = true
    try {
      const results = await schoolService.searchStudents(query)
      // Filter out already selected students
      studentSearchResults.value = results.filter(
        r => !formState.student_ids.includes(r.id)
      )
    } catch (err) {
      console.error('Failed to search students:', err)
      studentSearchResults.value = []
    } finally {
      searchingStudents.value = false
    }
  }, 300)
}

// Add student to selection
const addStudentToSelection = (student: StudentSearchResult) => {
  if (!formState.student_ids.includes(student.id)) {
    formState.student_ids.push(student.id)
    selectedStudentsForLinking.value.push(student)
  }
  // Clear search
  studentSearchQuery.value = ''
  studentSearchResults.value = []
}

// Remove student from selection
const removeStudentFromSelection = (studentId: number) => {
  formState.student_ids = formState.student_ids.filter(id => id !== studentId)
  selectedStudentsForLinking.value = selectedStudentsForLinking.value.filter(s => s.id !== studentId)
}

// Get student name by ID (for display)
const getStudentDisplayName = (studentId: number): string => {
  // First check in selected students for linking
  const selected = selectedStudentsForLinking.value.find(s => s.id === studentId)
  if (selected) {
    return `${selected.name}${selected.className ? ` - ${selected.className}` : ''}`
  }
  // Then check in loaded students
  const student = students.value.find(s => s.id === studentId)
  if (student) {
    return `${student.name}${student.className ? ` - ${student.className}` : ''}`
  }
  return `Siswa #${studentId}`
}

// Handle table change (pagination, sorting)
const handleTableChange: TableProps['onChange'] = (pag) => {
  pagination.current = pag.current || 1
  pagination.pageSize = pag.pageSize || 10
  loadParents()
}

// Handle search
const handleSearch = () => {
  pagination.current = 1
  loadParents()
}

// Open modal for create
const openCreateModal = () => {
  isEditing.value = false
  editingParent.value = null
  resetForm()
  modalVisible.value = true
}

// Open modal for edit
const openEditModal = async (parent: Parent) => {
  isEditing.value = true
  editingParent.value = parent
  formState.name = parent.name
  formState.phone = parent.phone || ''
  formState.email = parent.email || ''
  formState.student_ids = parent.studentIds || []
  
  // Populate selectedStudentsForLinking with existing linked students
  selectedStudentsForLinking.value = []
  if (parent.studentIds && parent.studentIds.length > 0 && parent.studentNames) {
    for (let i = 0; i < parent.studentIds.length; i++) {
      const student = students.value.find(s => s.id === parent.studentIds![i])
      selectedStudentsForLinking.value.push({
        id: parent.studentIds[i],
        nis: student?.nis || '',
        nisn: student?.nisn || '',
        name: parent.studentNames[i] || `Siswa #${parent.studentIds[i]}`,
        className: student?.className,
        classId: student?.classId,
      })
    }
  }
  
  modalVisible.value = true
}

// Reset form
const resetForm = () => {
  formState.name = ''
  formState.phone = ''
  formState.email = ''
  formState.student_ids = []
  selectedStudentsForLinking.value = []
  studentSearchQuery.value = ''
  studentSearchResults.value = []
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
    if (isEditing.value && editingParent.value) {
      const updateData: UpdateParentRequest = {
        name: formState.name,
        phone: formState.phone || undefined,
        email: formState.email || undefined,
        student_ids: formState.student_ids,
      }
      await schoolService.updateParent(editingParent.value.id, updateData)
      message.success('Data orang tua berhasil diperbarui')
      modalVisible.value = false
    } else {
      const result = await schoolService.createParent({
        name: formState.name,
        phone: formState.phone,
        email: formState.email || undefined,
        student_ids: formState.student_ids,
      })
      modalVisible.value = false
      
      // Show credential modal
      if (result.temporaryPassword) {
        credentialData.value = {
          username: result.username || formState.phone,
          password: result.temporaryPassword,
          name: formState.name,
        }
        credentialModalVisible.value = true
      } else {
        message.success('Orang tua berhasil ditambahkan')
      }
    }
    resetForm()
    loadParents()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Terjadi kesalahan')
  } finally {
    modalLoading.value = false
  }
}

// Handle reset password
const handleResetPassword = async (parent: Parent) => {
  try {
    const result = await schoolService.resetParentPassword(parent.id)
    credentialData.value = {
      username: result.username,
      password: result.temporaryPassword,
      name: parent.name,
    }
    credentialModalVisible.value = true
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Gagal reset password')
  }
}

// Handle delete parent
const handleDelete = async (parent: Parent) => {
  try {
    await schoolService.deleteParent(parent.id)
    message.success(`Data ${parent.name} berhasil dihapus`)
    loadParents()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Gagal menghapus data')
  }
}

// Copy to clipboard
const copyToClipboard = (text: string) => {
  navigator.clipboard.writeText(text)
  message.success('Berhasil disalin!')
}

// ==================== Import Functions ====================

// Download parent template
const handleDownloadTemplate = async () => {
  try {
    const blob = await importService.downloadParentTemplate()
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = 'template_import_orang_tua.xlsx'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    message.success('Template berhasil diunduh')
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string } } } }
    message.error(err.response?.data?.error?.message || 'Gagal mengunduh template')
  }
}

// Open import modal
const openImportModal = () => {
  importModalVisible.value = true
}

// Handle file upload for import
const handleImportUpload: UploadProps['customRequest'] = async (options) => {
  const { file, onSuccess, onError } = options
  
  if (!(file instanceof File)) {
    message.error('File tidak valid')
    onError?.(new Error('Invalid file'))
    return
  }

  // Validate file type
  const isExcel = file.type === 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' ||
                  file.name.endsWith('.xlsx')
  if (!isExcel) {
    message.error('Hanya file Excel (.xlsx) yang diperbolehkan')
    onError?.(new Error('Invalid file type'))
    return
  }

  // Validate file size (max 5MB)
  const maxSize = 5 * 1024 * 1024
  if (file.size > maxSize) {
    message.error('Ukuran file maksimal 5MB')
    onError?.(new Error('File too large'))
    return
  }

  importLoading.value = true
  try {
    const result = await importService.importParents(file)
    importResult.value = result
    importModalVisible.value = false
    importResultModalVisible.value = true
    
    // Reload parents list
    loadParents()
    onSuccess?.(result)
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string } } } }
    message.error(err.response?.data?.error?.message || 'Gagal mengimpor data orang tua')
    onError?.(error as Error)
  } finally {
    importLoading.value = false
  }
}

// Close import result modal
const closeImportResultModal = () => {
  importResultModalVisible.value = false
  importResult.value = null
}

onMounted(() => {
  loadParents()
  loadStudents()
})
</script>

<template>
  <div class="parent-management">
    <div class="page-header">
      <Title :level="2" style="margin: 0">Manajemen Orang Tua</Title>
    </div>

    <Alert
      type="info"
      show-icon
      style="margin-bottom: 16px"
      closable
    >
      <template #icon><InfoCircleOutlined /></template>
      <template #message>Informasi Akun Orang Tua</template>
      <template #description>
        <ul style="margin: 8px 0 0 0; padding-left: 20px;">
          <li><strong>Username:</strong> Nomor HP orang tua</li>
          <li><strong>Password Default:</strong> <code>password123</code></li>
          <li>Orang tua wajib mengganti password saat login pertama kali.</li>
        </ul>
      </template>
    </Alert>

    <Card>
      <!-- Toolbar -->
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col :xs="24" :sm="12" :md="8">
          <Input
            v-model:value="searchText"
            placeholder="Cari orang tua..."
            allow-clear
            @press-enter="handleSearch"
          >
            <template #prefix>
              <SearchOutlined />
            </template>
          </Input>
        </Col>
        <Col :xs="24" :sm="12" :md="16" class="toolbar-right">
          <Space wrap>
            <Button @click="loadParents">
              <template #icon><ReloadOutlined /></template>
            </Button>
            <Button @click="handleDownloadTemplate">
              <template #icon><DownloadOutlined /></template>
              Download Template
            </Button>
            <Button type="default" @click="openImportModal">
              <template #icon><UploadOutlined /></template>
              Import Excel
            </Button>
            <Button type="primary" @click="openCreateModal">
              <template #icon><PlusOutlined /></template>
              Tambah Orang Tua
            </Button>
          </Space>
        </Col>
      </Row>

      <!-- Table -->
      <Table
        :columns="columns"
        :data-source="filteredParents"
        :loading="loading"
        :pagination="{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: total,
          showSizeChanger: true,
          showTotal: (total: number) => `Total ${total} orang tua`,
        }"
        row-key="id"
        @change="handleTableChange"
        class="custom-table"
        :scroll="{ x: 800 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <div class="user-info">
              <div class="user-avatar">
                <UserOutlined />
              </div>
              <span class="user-name">{{ (record as Parent).name }}</span>
            </div>
          </template>
          <template v-else-if="column.key === 'phone'">
            <span class="phone-badge">{{ (record as Parent).phone || '-' }}</span>
          </template>
          <template v-else-if="column.key === 'email'">
            <span class="text-secondary">{{ (record as Parent).email || '-' }}</span>
          </template>
          <template v-else-if="column.key === 'studentNames'">
            <div class="student-list">
              <span
                v-for="(name, index) in (record as Parent).studentNames"
                :key="index"
                class="student-badge"
              >
                {{ name }}
              </span>
              <span v-if="!(record as Parent).studentNames?.length" class="text-secondary italic">
                Belum ada anak
              </span>
            </div>
          </template>
          <template v-else-if="column.key === 'action'">
            <Space>
               <Tooltip title="Reset Password">
                <Popconfirm
                  title="Reset password orang tua ini?"
                  description="Password baru akan di-generate otomatis."
                  ok-text="Ya, Reset"
                  cancel-text="Batal"
                  @confirm="handleResetPassword(record as Parent)"
                >
                  <Button size="small" type="text">
                    <template #icon><KeyOutlined style="color: #f59e0b;" /></template>
                  </Button>
                </Popconfirm>
              </Tooltip>
              <Tooltip title="Edit">
                <Button size="small" type="text" @click="openEditModal(record as Parent)">
                  <template #icon><EditOutlined style="color: #64748b;" /></template>
                </Button>
              </Tooltip>
              <Tooltip title="Hapus">
                <Popconfirm
                  title="Hapus data orang tua ini?"
                  description="Akun orang tua akan dihapus."
                  ok-text="Ya, Hapus"
                  cancel-text="Batal"
                  ok-type="danger"
                  @confirm="handleDelete(record as Parent)"
                >
                  <Button size="small" type="text" danger>
                    <template #icon><DeleteOutlined /></template>
                  </Button>
                </Popconfirm>
              </Tooltip>
            </Space>
          </template>
        </template>
      </Table>
    </Card>

    <!-- Create/Edit Modal -->
    <Modal
      v-model:open="modalVisible"
      :title="isEditing ? 'Edit Orang Tua' : 'Tambah Orang Tua Baru'"
      :confirm-loading="modalLoading"
      width="600px"
      wrap-class-name="modern-modal"
      @ok="handleSubmit"
      @cancel="handleModalCancel"
    >
      <Form
        ref="formRef"
        :model="formState"
        :rules="formRules"
        layout="vertical"
        class="modern-form"
      >
        <FormItem label="Nama Lengkap" name="name" required>
          <Input v-model:value="formState.name" placeholder="Nama lengkap orang tua" size="large" />
        </FormItem>
        <Row :gutter="16">
          <Col :span="12">
            <FormItem 
              label="Nomor Telepon (Username)" 
              name="phone" 
              required
              :extra="!isEditing ? 'Digunakan sebagai username untuk login' : ''"
            >
              <Input 
                v-model:value="formState.phone" 
                placeholder="Contoh: 081234567890"
                :disabled="isEditing"
                size="large"
              />
            </FormItem>
          </Col>
          <Col :span="12">
            <FormItem label="Email" name="email">
              <Input v-model:value="formState.email" placeholder="email@example.com" size="large" />
            </FormItem>
          </Col>
        </Row>
        <FormItem
          v-if="!isEditing"
        >
          <template #label>
            <div class="form-label-with-tag">
              <span>Password</span>
              <span class="default-badge">Default</span>
            </div>
          </template>
          <div class="password-info">
            <Text type="secondary">
              Password default: <Text code>password123</Text> — Orang tua wajib mengganti saat login pertama.
            </Text>
          </div>
        </FormItem>
        <FormItem
          label="Anak (Siswa)"
          name="student_ids"
        >
          <template #extra>
            <span v-if="!isEditing">Pilih siswa yang merupakan anak dari orang tua ini</span>
            <span v-else>Cari siswa berdasarkan NISN atau nama untuk menghubungkan</span>
          </template>
          
          <!-- Student Search Input -->
          <Input
            v-model:value="studentSearchQuery"
            placeholder="Cari siswa berdasarkan NISN atau nama..."
            allow-clear
            size="large"
            @input="(e: Event) => handleStudentSearch((e.target as HTMLInputElement).value)"
            style="margin-bottom: 8px"
          >
            <template #prefix>
              <SearchOutlined />
            </template>
            <template #suffix>
              <span v-if="searchingStudents" class="ant-spin-dot ant-spin-dot-spin" style="font-size: 12px">
                <i class="ant-spin-dot-item"></i>
              </span>
            </template>
          </Input>
          
          <!-- Search Results -->
          <div v-if="studentSearchResults.length > 0" class="student-search-results">
            <div
              v-for="student in studentSearchResults"
              :key="student.id"
              class="student-search-item"
              @click="addStudentToSelection(student)"
            >
              <div class="student-search-info">
                <span class="student-name">{{ student.name }}</span>
                <span class="student-details">
                  NISN: {{ student.nisn }} | NIS: {{ student.nis }}
                </span>
              </div>
              <span v-if="student.className" class="class-badge-small blue">{{ student.className }}</span>
              <span v-else class="class-badge-small orange">Belum ada kelas</span>
            </div>
          </div>
          
          <!-- No Results Message -->
          <div v-else-if="studentSearchQuery.length >= 2 && !searchingStudents && studentSearchResults.length === 0" class="student-search-empty">
            <Text type="secondary">Tidak ada siswa ditemukan</Text>
          </div>
          
          <!-- Selected Students List -->
          <div v-if="formState.student_ids.length > 0" class="selected-students">
            <Divider style="margin: 12px 0 8px 0">
              <Text type="secondary" style="font-size: 12px">
                <LinkOutlined /> Siswa Terhubung ({{ formState.student_ids.length }})
              </Text>
            </Divider>
            <div class="selected-students-list">
              <div
                v-for="studentId in formState.student_ids"
                :key="studentId"
                class="student-chip"
              >
                {{ getStudentDisplayName(studentId) }}
                <CloseCircleOutlined @click="removeStudentFromSelection(studentId)" class="remove-icon" />
              </div>
            </div>
          </div>
          
          <!-- Empty State for Edit Mode -->
          <Alert
            v-else-if="isEditing"
            type="info"
            show-icon
            style="margin-top: 8px"
          >
            <template #message>Belum ada siswa terhubung</template>
            <template #description>
              Gunakan pencarian di atas untuk menghubungkan orang tua dengan siswa.
            </template>
          </Alert>
        </FormItem>
      </Form>
    </Modal>

    <!-- Credential Modal -->
    <Modal
      v-model:open="credentialModalVisible"
      title="Kredensial Akun Orang Tua"
      :footer="null"
      width="500px"
    >
      <div v-if="credentialData" class="credential-info">
        <div class="credential-header">
          <UserOutlined style="font-size: 48px; color: #52c41a" />
          <Title :level="4" style="margin: 16px 0 8px">{{ credentialData.name }}</Title>
          <Text type="secondary">Akun berhasil dibuat. Simpan kredensial berikut:</Text>
        </div>
        
        <Card class="credential-card">
          <div class="credential-item">
            <Text type="secondary">Username:</Text>
            <div class="credential-value">
              <Text strong copyable>{{ credentialData.username }}</Text>
              <Button size="small" @click="copyToClipboard(credentialData.username)">
                <template #icon><CopyOutlined /></template>
              </Button>
            </div>
          </div>
          <div class="credential-item">
            <Text type="secondary">Password:</Text>
            <div class="credential-value">
              <Text strong code>{{ credentialData.password }}</Text>
              <Button size="small" @click="copyToClipboard(credentialData.password)">
                <template #icon><CopyOutlined /></template>
              </Button>
            </div>
          </div>
        </Card>

        <div class="credential-note">
          <Text type="warning">
            ⚠️ Password ini hanya ditampilkan sekali. Pastikan untuk menyimpan atau memberikan ke orang tua.
          </Text>
        </div>

        <Button type="primary" block @click="credentialModalVisible = false" style="margin-top: 16px">
          Tutup
        </Button>
      </div>
    </Modal>

    <!-- Import Modal -->
    <Modal
      v-model:open="importModalVisible"
      title="Import Data Orang Tua dari Excel"
      :footer="null"
      width="500px"
    >
      <div class="import-info">
        <div class="import-header">
          <FileExcelOutlined style="font-size: 48px; color: #52c41a" />
          <Title :level="4" style="margin: 16px 0 8px">Import Orang Tua</Title>
          <Text type="secondary">
            Upload file Excel (.xlsx) untuk mengimpor data orang tua secara massal
          </Text>
        </div>

        <Alert
          type="info"
          show-icon
          style="margin: 16px 0"
        >
          <template #message>Format Template</template>
          <template #description>
            <ul style="margin: 8px 0 0 0; padding-left: 20px;">
              <li><strong>Nama</strong> - Nama lengkap orang tua (wajib)</li>
              <li><strong>No_HP</strong> - Nomor telepon/username (wajib)</li>
              <li><strong>Email</strong> - Alamat email (wajib)</li>
            </ul>
          </template>
        </Alert>

        <Alert
          type="warning"
          show-icon
          style="margin-bottom: 16px"
        >
          <template #message>Catatan Penting</template>
          <template #description>
            <ul style="margin: 8px 0 0 0; padding-left: 20px;">
              <li>Orang tua yang diimpor belum terhubung dengan siswa</li>
              <li>Gunakan fitur edit untuk menghubungkan orang tua dengan siswa</li>
              <li>Password default: <code>password123</code></li>
            </ul>
          </template>
        </Alert>

        <Upload.Dragger
          :custom-request="handleImportUpload"
          :show-upload-list="false"
          accept=".xlsx"
          :disabled="importLoading"
        >
          <p class="ant-upload-drag-icon">
            <UploadOutlined v-if="!importLoading" />
            <Progress v-else type="circle" :percent="0" :size="48" status="active" />
          </p>
          <p class="ant-upload-text">
            {{ importLoading ? 'Mengimpor data...' : 'Klik atau drag file ke area ini' }}
          </p>
          <p class="ant-upload-hint">
            Hanya file Excel (.xlsx), maksimal 5MB
          </p>
        </Upload.Dragger>

        <Divider />

        <Button block @click="handleDownloadTemplate">
          <template #icon><DownloadOutlined /></template>
          Download Template Excel
        </Button>
      </div>
    </Modal>

    <!-- Import Result Modal -->
    <Modal
      v-model:open="importResultModalVisible"
      title="Hasil Import Orang Tua"
      :footer="null"
      width="600px"
      @cancel="closeImportResultModal"
    >
      <div v-if="importResult" class="import-result">
        <!-- Summary -->
        <Row :gutter="16" style="margin-bottom: 24px">
          <Col :span="8">
            <Card size="small" class="result-card">
              <div class="result-number">{{ importResult.total_rows }}</div>
              <div class="result-label">Total Baris</div>
            </Card>
          </Col>
          <Col :span="8">
            <Card size="small" class="result-card success">
              <div class="result-number">{{ importResult.success_count }}</div>
              <div class="result-label">Berhasil</div>
            </Card>
          </Col>
          <Col :span="8">
            <Card size="small" class="result-card error">
              <div class="result-number">{{ importResult.failed_count + importResult.warning_count }}</div>
              <div class="result-label">Gagal/Skip</div>
            </Card>
          </Col>
        </Row>

        <!-- Success Message -->
        <Alert
          v-if="importResult.success_count > 0 && importResult.failed_count === 0"
          type="success"
          show-icon
          style="margin-bottom: 16px"
        >
          <template #icon><CheckCircleOutlined /></template>
          <template #message>
            Semua data berhasil diimpor!
          </template>
          <template #description>
            Orang tua yang diimpor belum terhubung dengan siswa. Gunakan fitur edit untuk menghubungkan.
          </template>
        </Alert>

        <!-- Partial Success Message -->
        <Alert
          v-else-if="importResult.success_count > 0"
          type="info"
          show-icon
          style="margin-bottom: 16px"
        >
          <template #icon><InfoCircleOutlined /></template>
          <template #message>
            {{ importResult.success_count }} data berhasil diimpor
          </template>
          <template #description>
            Orang tua yang diimpor belum terhubung dengan siswa. Gunakan fitur edit untuk menghubungkan.
          </template>
        </Alert>

        <!-- Errors List -->
        <div v-if="importResult.errors && importResult.errors.length > 0" class="result-section">
          <Title :level="5">
            <CloseCircleOutlined style="color: #ff4d4f" /> Error ({{ importResult.errors.length }})
          </Title>
          <List
            size="small"
            :data-source="importResult.errors"
            :bordered="true"
            style="max-height: 200px; overflow-y: auto"
          >
            <template #renderItem="{ item }">
              <ListItem>
                <ListItemMeta>
                  <template #title>
                    <span class="row-badge error">Baris {{ (item as ImportError).row }}</span>
                    {{ (item as ImportError).field }}
                  </template>
                  <template #description>
                    {{ (item as ImportError).message }}
                  </template>
                </ListItemMeta>
              </ListItem>
            </template>
          </List>
        </div>

        <!-- Warnings List -->
        <div v-if="importResult.warnings && importResult.warnings.length > 0" class="result-section">
          <Title :level="5">
            <WarningOutlined style="color: #faad14" /> Warning ({{ importResult.warnings.length }})
          </Title>
          <List
            size="small"
            :data-source="importResult.warnings"
            :bordered="true"
            style="max-height: 200px; overflow-y: auto"
          >
            <template #renderItem="{ item }">
              <ListItem>
                <ListItemMeta>
                  <template #title>
                    <span class="row-badge warning">Baris {{ (item as ImportWarning).row }}</span>
                    {{ (item as ImportWarning).field }}
                  </template>
                  <template #description>
                    {{ (item as ImportWarning).message }}
                  </template>
                </ListItemMeta>
              </ListItem>
            </template>
          </List>
        </div>

        <Divider />

        <Button type="primary" block @click="closeImportResultModal">
          Tutup
        </Button>
      </div>
    </Modal>
  </div>
</template>

<style scoped>
.parent-management {
  padding: 0;
}

.page-header {
  margin-bottom: 24px;
}

/* Toolbar */
.toolbar {
  margin-bottom: 24px;
}

.toolbar-right {
  display: flex;
  justify-content: flex-end;
}

/* Custom Table */
.custom-table :deep(.ant-table-thead > tr > th) {
  background: #f8fafc;
  color: #475569;
  font-weight: 600;
  border-bottom: 1px solid #f1f5f9;
}

.custom-table :deep(.ant-table-tbody > tr > td) {
  padding: 16px;
  border-bottom: 1px solid #f1f5f9;
}

/* User Info */
.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  width: 32px;
  height: 32px;
  background: #eff6ff;
  color: #3b82f6;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
}

.user-name {
  font-weight: 500;
  color: #1e293b;
}

/* Badges */
.phone-badge {
  background: #f8fafc;
  color: #64748b;
  padding: 2px 8px;
  border-radius: 6px;
  font-family: monospace;
  font-size: 13px;
  border: 1px solid #e2e8f0;
}

.text-secondary {
  color: #94a3b8;
}

.italic {
  font-style: italic;
}

.student-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.student-badge {
  background: #f0fdf4;
  color: #166534;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  border: 1px solid #dcfce7;
}

/* Form Styles */
.form-label-with-tag {
  display: flex;
  align-items: center;
  gap: 8px;
}

.default-badge {
  background: #eff6ff;
  color: #3b82f6;
  font-size: 10px;
  padding: 1px 6px;
  border-radius: 4px;
  border: 1px solid #dbeafe;
}

/* Student Search */
.student-search-results {
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  max-height: 200px;
  overflow-y: auto;
  margin-bottom: 8px;
  background: white;
}

.student-search-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  cursor: pointer;
  transition: all 0.2s;
  border-bottom: 1px solid #f1f5f9;
}

.student-search-item:last-child {
  border-bottom: none;
}

.student-search-item:hover {
  background-color: #f8fafc;
}

.student-search-info {
  display: flex;
  flex-direction: column;
}

.student-name {
  font-weight: 500;
  color: #334155;
}

.student-details {
  font-size: 11px;
  color: #94a3b8;
}

.class-badge-small {
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 4px;
}

.class-badge-small.blue {
  background: #eff6ff;
  color: #3b82f6;
}

.class-badge-small.orange {
  background: #fff7ed;
  color: #f97316;
}

.selected-students-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.student-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: #f0fdf4;
  color: #15803d;
  padding: 4px 10px;
  border-radius: 16px;
  font-size: 13px;
  border: 1px solid #dcfce7;
}

.remove-icon {
  font-size: 12px;
  cursor: pointer;
  color: #15803d;
  opacity: 0.7;
}

.remove-icon:hover {
  opacity: 1;
}

/* Credential */
.credential-info {
  text-align: center;
}

.credential-header {
  margin-bottom: 24px;
}

.credential-card {
  text-align: left;
  margin-bottom: 16px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
}

.credential-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
}

.credential-item:not(:last-child) {
  border-bottom: 1px solid #e2e8f0;
}

.credential-value {
  display: flex;
  align-items: center;
  gap: 8px;
}

.credential-note {
  background: #fffbe6;
  padding: 12px;
  border-radius: 6px;
  text-align: left;
  border: 1px solid #ffe58f;
}

/* Import Styles */
.import-info {
  text-align: center;
}

.import-header {
  margin-bottom: 24px;
}

.import-result {
  text-align: left;
}

.result-card {
  text-align: center;
  height: 100%;
}

.result-card.success .result-number { color: #22c55e; }
.result-card.error .result-number { color: #ef4444; }

.result-number {
  font-size: 24px;
  font-weight: 700;
  color: #334155;
  line-height: 1.2;
}

.result-label {
  font-size: 12px;
  color: #64748b;
  margin-top: 4px;
}

.result-section {
  margin-bottom: 24px;
}

.row-badge {
  display: inline-block;
  padding: 1px 6px;
  border-radius: 4px;
  font-size: 11px;
  margin-right: 8px;
}

.row-badge.error {
  background: #fef2f2;
  color: #ef4444;
  border: 1px solid #fee2e2;
}

.row-badge.warning {
  background: #fffbeb;
  color: #f59e0b;
  border: 1px solid #fef3c7;
}

@media (max-width: 576px) {
  .toolbar-right {
    margin-top: 16px;
    justify-content: flex-start;
  }
}
</style>
