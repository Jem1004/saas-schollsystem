<script setup lang="ts">
import { ref, reactive, onMounted, computed, onUnmounted } from 'vue'
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
  Switch,
  Checkbox,
  Tooltip,
  Alert,
  Progress,
  Upload,
  Divider,
} from 'ant-design-vue'
import type { TableProps, UploadProps } from 'ant-design-vue'
import {
  PlusOutlined,
  SearchOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
  FilterOutlined,
  UserOutlined,
  KeyOutlined,
  CopyOutlined,
  MobileOutlined,
  InfoCircleOutlined,
  WifiOutlined,
  CloseCircleOutlined,
  DownloadOutlined,
  UploadOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
} from '@ant-design/icons-vue'
import { schoolService, importService } from '@/services'
import type { ImportResult, ImportError } from '@/services/import'
import type { Student, Class, UpdateStudentRequest, Device, PairingSessionResponse } from '@/types/school'

const { Title, Text } = Typography

// Table state
const loading = ref(false)
const students = ref<Student[]>([])
const total = ref(0)
const pagination = reactive({
  current: 1,
  pageSize: 10,
})
const searchText = ref('')
const filterClassId = ref<number | undefined>(undefined)

// Classes for filter and form
const classes = ref<Class[]>([])
const loadingClasses = ref(false)

// Modal state
const modalVisible = ref(false)
const modalLoading = ref(false)
const isEditing = ref(false)
const editingStudent = ref<Student | null>(null)

// Credential modal state
const credentialModalVisible = ref(false)
const credentialData = ref<{ username: string; password: string; name: string } | null>(null)

// Pairing modal state
const pairingModalVisible = ref(false)
const pairingLoading = ref(false)
const pairingStudent = ref<Student | null>(null)
const pairingSession = ref<PairingSessionResponse | null>(null)
const pairingCountdown = ref(0)
const pairingTimer = ref<ReturnType<typeof setInterval> | null>(null)
const selectedDeviceId = ref<number | undefined>(undefined)
const devices = ref<Device[]>([])
const loadingDevices = ref(false)

// Import modal state
const importModalVisible = ref(false)
const importLoading = ref(false)
const importResultModalVisible = ref(false)
const importResult = ref<ImportResult | null>(null)
const studentsWithoutClass = ref(0)

// Bulk class assignment state
const selectedStudentIds = ref<number[]>([])
const bulkAssignModalVisible = ref(false)
const bulkAssignLoading = ref(false)
const bulkAssignClassId = ref<number | undefined>(undefined)
const filterWithoutClass = ref(false)

// Form state
const formRef = ref()
const formState = reactive({
  class_id: undefined as number | undefined,
  nis: '',
  nisn: '',
  name: '',
  rfid_code: '',
  is_active: true,
  create_account: false,
})

// Form rules
const formRules = {
  class_id: [{ required: true, message: 'Kelas wajib dipilih' }],
  nis: [{ required: true, message: 'NIS wajib diisi' }],
  nisn: [{ required: true, message: 'NISN wajib diisi' }],
  name: [{ required: true, message: 'Nama siswa wajib diisi' }],
}

// Table columns
const columns: TableProps['columns'] = [
  {
    title: 'NIS',
    dataIndex: 'nis',
    key: 'nis',
    width: 100,
  },
  {
    title: 'NISN',
    dataIndex: 'nisn',
    key: 'nisn',
    width: 120,
  },
  {
    title: 'Nama Siswa',
    dataIndex: 'name',
    key: 'name',
    sorter: true,
  },
  {
    title: 'Kelas',
    dataIndex: 'className',
    key: 'className',
    width: 100,
  },
  {
    title: 'RFID',
    dataIndex: 'rfidCode',
    key: 'rfidCode',
    width: 120,
    align: 'center',
  },
  {
    title: 'Akun',
    dataIndex: 'hasAccount',
    key: 'hasAccount',
    width: 80,
    align: 'center',
  },
  {
    title: 'Status',
    dataIndex: 'isActive',
    key: 'isActive',
    width: 100,
    align: 'center',
  },
  {
    title: 'Aksi',
    key: 'action',
    width: 250,
    align: 'center',
  },
]

// Computed filtered data
const filteredStudents = computed(() => {
  let result = students.value
  
  if (filterClassId.value) {
    result = result.filter(s => s.classId === filterClassId.value)
  }
  
  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(
      (student) =>
        student.name.toLowerCase().includes(search) ||
        student.nis.includes(search) ||
        student.nisn.includes(search)
    )
  }
  
  return result
})

// Load students data
const loadStudents = async () => {
  loading.value = true
  selectedStudentIds.value = [] // Clear selection when loading
  try {
    if (filterWithoutClass.value) {
      // Load students without class from import service
      const studentsData = await importService.getStudentsWithoutClass()
      students.value = studentsData.map(s => ({
        id: s.id,
        schoolId: 0,
        classId: 0,
        className: undefined,
        nis: s.nis,
        nisn: s.nisn,
        name: s.name,
        rfidCode: undefined,
        isActive: s.is_active,
        hasAccount: false,
        createdAt: '',
        updatedAt: '',
      })) as Student[]
      total.value = studentsData.length
    } else {
      const response = await schoolService.getStudents({
        page: pagination.current,
        page_size: pagination.pageSize,
        search: searchText.value,
        class_id: filterClassId.value,
      })
      students.value = response.students
      total.value = response.pagination.total
    }
  } catch (err) {
    console.error('Failed to load students:', err)
    message.error('Gagal memuat data siswa')
    students.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// Load classes for dropdown
const loadClasses = async () => {
  loadingClasses.value = true
  try {
    const response = await schoolService.getClasses({ page_size: 100 })
    classes.value = response.classes
  } catch (err) {
    console.error('Failed to load classes:', err)
    classes.value = []
  } finally {
    loadingClasses.value = false
  }
}

// Handle table change (pagination, sorting)
const handleTableChange: TableProps['onChange'] = (pag) => {
  pagination.current = pag.current || 1
  pagination.pageSize = pag.pageSize || 10
  loadStudents()
}

// Handle search
const handleSearch = () => {
  pagination.current = 1
  loadStudents()
}

// Handle class filter change
const handleClassFilterChange = () => {
  pagination.current = 1
  loadStudents()
}

// Open modal for create
const openCreateModal = () => {
  isEditing.value = false
  editingStudent.value = null
  resetForm()
  modalVisible.value = true
}

// Open modal for edit
const openEditModal = (student: Student) => {
  isEditing.value = true
  editingStudent.value = student
  formState.class_id = student.classId
  formState.nis = student.nis
  formState.nisn = student.nisn
  formState.name = student.name
  formState.rfid_code = student.rfidCode || ''
  formState.is_active = student.isActive
  formState.create_account = false
  modalVisible.value = true
}

// Reset form
const resetForm = () => {
  formState.class_id = undefined
  formState.nis = ''
  formState.nisn = ''
  formState.name = ''
  formState.rfid_code = ''
  formState.is_active = true
  formState.create_account = false
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

  // Pastikan class_id valid
  if (!formState.class_id || formState.class_id <= 0) {
    message.error('Kelas wajib dipilih')
    return
  }

  modalLoading.value = true
  try {
    if (isEditing.value && editingStudent.value) {
      const updateData: UpdateStudentRequest = {
        class_id: formState.class_id,
        nis: formState.nis,
        name: formState.name,
        rfid_code: formState.rfid_code || undefined,
        is_active: formState.is_active,
      }
      await schoolService.updateStudent(editingStudent.value.id, updateData)
      message.success('Siswa berhasil diperbarui')
      modalVisible.value = false
    } else {
      const result = await schoolService.createStudent({
        class_id: formState.class_id,
        nis: formState.nis,
        nisn: formState.nisn,
        name: formState.name,
        rfid_code: formState.rfid_code || undefined,
        create_account: formState.create_account,
      })
      modalVisible.value = false
      
      // Show credential modal if account was created
      if (formState.create_account && result.temporaryPassword) {
        credentialData.value = {
          username: result.username || formState.nis,
          password: result.temporaryPassword,
          name: formState.name,
        }
        credentialModalVisible.value = true
      } else {
        message.success('Siswa berhasil ditambahkan')
      }
    }
    resetForm()
    loadStudents()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Terjadi kesalahan')
  } finally {
    modalLoading.value = false
  }
}

// Handle create account for existing student
const handleCreateAccount = async (student: Student) => {
  try {
    const result = await schoolService.createStudentAccount(student.id)
    if (result.temporaryPassword) {
      credentialData.value = {
        username: result.username || student.nis,
        password: result.temporaryPassword,
        name: student.name,
      }
      credentialModalVisible.value = true
    }
    loadStudents()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Gagal membuat akun')
  }
}

// Handle reset password
const handleResetPassword = async (student: Student) => {
  try {
    const result = await schoolService.resetStudentPassword(student.id)
    credentialData.value = {
      username: result.username,
      password: result.temporaryPassword,
      name: student.name,
    }
    credentialModalVisible.value = true
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Gagal reset password')
  }
}

// Handle delete student
const handleDelete = async (student: Student) => {
  try {
    await schoolService.deleteStudent(student.id)
    message.success(`Siswa ${student.name} berhasil dihapus`)
    loadStudents()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Gagal menghapus siswa')
  }
}

// Copy to clipboard
const copyToClipboard = (text: string) => {
  navigator.clipboard.writeText(text)
  message.success('Berhasil disalin!')
}

// Load devices for pairing
const loadDevices = async () => {
  loadingDevices.value = true
  try {
    devices.value = await schoolService.getSchoolDevices()
  } catch (err) {
    console.error('Failed to load devices:', err)
    devices.value = []
  } finally {
    loadingDevices.value = false
  }
}

// Open pairing modal
const openPairingModal = async (student: Student) => {
  pairingStudent.value = student
  pairingSession.value = null
  pairingCountdown.value = 0
  selectedDeviceId.value = undefined
  pairingModalVisible.value = true
  await loadDevices()
}

// Start pairing session
const startPairing = async () => {
  if (!selectedDeviceId.value || !pairingStudent.value) {
    message.error('Pilih perangkat terlebih dahulu')
    return
  }

  pairingLoading.value = true
  try {
    const response = await schoolService.startPairing(selectedDeviceId.value, pairingStudent.value.id)
    pairingSession.value = response
    
    if (response.active && response.expiresAt) {
      // Start countdown
      const expiresAt = new Date(response.expiresAt).getTime()
      const updateCountdown = () => {
        const now = Date.now()
        const remaining = Math.max(0, Math.floor((expiresAt - now) / 1000))
        pairingCountdown.value = remaining
        
        if (remaining <= 0) {
          if (pairingTimer.value) {
            clearInterval(pairingTimer.value)
            pairingTimer.value = null
          }
          pairingSession.value = { ...pairingSession.value!, active: false, message: 'Sesi pairing kadaluarsa' }
        }
      }
      
      updateCountdown()
      pairingTimer.value = setInterval(updateCountdown, 1000)
      
      // Poll for pairing status
      pollPairingStatus()
    }
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string } } } }
    message.error(err.response?.data?.error?.message || 'Gagal memulai pairing')
  } finally {
    pairingLoading.value = false
  }
}

// Poll pairing status
const pollPairingStatus = async () => {
  if (!selectedDeviceId.value || !pairingSession.value?.active) return
  
  const pollInterval = setInterval(async () => {
    try {
      const status = await schoolService.getPairingStatus(selectedDeviceId.value!)
      
      if (!status.active) {
        clearInterval(pollInterval)
        if (pairingTimer.value) {
          clearInterval(pairingTimer.value)
          pairingTimer.value = null
        }
        
        // Reload student data to check if RFID was paired successfully
        if (pairingStudent.value) {
          try {
            const updatedStudent = await schoolService.getStudent(pairingStudent.value.id)
            if (updatedStudent.rfidCode) {
              // Pairing was successful - student now has RFID
              message.success(`Kartu RFID berhasil dipasangkan ke ${updatedStudent.name}!`)
              pairingModalVisible.value = false
              loadStudents()
              return
            }
          } catch {
            // Ignore error, fall through to show session status
          }
        }
        
        // Session ended without successful pairing (expired or cancelled)
        pairingSession.value = status
      }
    } catch {
      // Ignore polling errors
    }
  }, 2000)
  
  // Stop polling after 65 seconds
  setTimeout(() => clearInterval(pollInterval), 65000)
}

// Cancel pairing
const cancelPairing = async () => {
  if (!selectedDeviceId.value) return
  
  try {
    await schoolService.cancelPairing(selectedDeviceId.value)
    if (pairingTimer.value) {
      clearInterval(pairingTimer.value)
      pairingTimer.value = null
    }
    pairingSession.value = null
    pairingCountdown.value = 0
    message.info('Sesi pairing dibatalkan')
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string } } } }
    message.error(err.response?.data?.error?.message || 'Gagal membatalkan pairing')
  }
}

// Close pairing modal
const closePairingModal = () => {
  if (pairingSession.value?.active && selectedDeviceId.value) {
    cancelPairing()
  }
  pairingModalVisible.value = false
  pairingStudent.value = null
  pairingSession.value = null
}

// Clear student RFID
const handleClearRFID = async (student: Student) => {
  try {
    await schoolService.clearStudentRFID(student.id)
    message.success('Kartu RFID berhasil dihapus dari siswa')
    loadStudents()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string } } } }
    message.error(err.response?.data?.error?.message || 'Gagal menghapus kartu RFID')
  }
}

// ==================== Import Functions ====================

// Download student template
const handleDownloadTemplate = async () => {
  try {
    const blob = await importService.downloadStudentTemplate()
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = 'template_import_siswa.xlsx'
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
    const result = await importService.importStudents(file)
    importResult.value = result
    importModalVisible.value = false
    importResultModalVisible.value = true
    
    // Update students without class count
    if (result.students_without_class) {
      studentsWithoutClass.value = result.students_without_class
    }
    
    // Reload students list
    loadStudents()
    onSuccess?.(result)
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string } } } }
    message.error(err.response?.data?.error?.message || 'Gagal mengimpor data siswa')
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

// Load students without class count
const loadStudentsWithoutClassCount = async () => {
  try {
    const students = await importService.getStudentsWithoutClass()
    studentsWithoutClass.value = students.length
  } catch {
    // Ignore error, just don't show the count
    studentsWithoutClass.value = 0
  }
}

// ==================== Bulk Class Assignment Functions ====================

// Handle filter without class change
const handleFilterWithoutClassChange = () => {
  pagination.current = 1
  if (filterWithoutClass.value) {
    filterClassId.value = undefined // Clear class filter when showing without class
  }
  loadStudents()
}

// Row selection config
const rowSelection = computed(() => ({
  selectedRowKeys: selectedStudentIds.value,
  onChange: (selectedKeys: (string | number)[]) => {
    selectedStudentIds.value = selectedKeys.map(k => Number(k))
  },
  getCheckboxProps: () => ({
    disabled: !filterWithoutClass.value, // Only allow selection when filtering without class
  }),
}))

// Open bulk assign modal
const openBulkAssignModal = () => {
  if (selectedStudentIds.value.length === 0) {
    message.warning('Pilih siswa terlebih dahulu')
    return
  }
  bulkAssignClassId.value = undefined
  bulkAssignModalVisible.value = true
}

// Handle bulk class assignment
const handleBulkAssign = async () => {
  if (!bulkAssignClassId.value) {
    message.error('Pilih kelas terlebih dahulu')
    return
  }

  bulkAssignLoading.value = true
  try {
    const result = await importService.bulkAssignClass(selectedStudentIds.value, bulkAssignClassId.value)
    message.success(`${result.updated_count} siswa berhasil diassign ke kelas`)
    bulkAssignModalVisible.value = false
    selectedStudentIds.value = []
    loadStudents()
    loadStudentsWithoutClassCount()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string } } } }
    message.error(err.response?.data?.error?.message || 'Gagal mengassign kelas')
  } finally {
    bulkAssignLoading.value = false
  }
}

// Close bulk assign modal
const closeBulkAssignModal = () => {
  bulkAssignModalVisible.value = false
  bulkAssignClassId.value = undefined
}

onMounted(() => {
  loadStudents()
  loadClasses()
  loadStudentsWithoutClassCount()
})

onUnmounted(() => {
  if (pairingTimer.value) {
    clearInterval(pairingTimer.value)
  }
})
</script>

<template>
  <div class="student-management">
    <div class="page-header">
      <Title :level="2" style="margin: 0">Manajemen Siswa</Title>
    </div>

    <Alert
      type="info"
      show-icon
      style="margin-bottom: 16px"
      closable
    >
      <template #icon><InfoCircleOutlined /></template>
      <template #message>Informasi Akun Siswa untuk Mobile Apps</template>
      <template #description>
        <ul style="margin: 8px 0 0 0; padding-left: 20px;">
          <li><strong>Username:</strong> NIS (Nomor Induk Siswa)</li>
          <li><strong>Password Default:</strong> <code>password123</code></li>
          <li>Centang "Buat akun untuk login mobile apps" saat menambah siswa, atau gunakan tombol <UserOutlined /> untuk membuat akun nanti.</li>
          <li>Siswa wajib mengganti password saat login pertama kali.</li>
        </ul>
      </template>
    </Alert>

    <Card>
      <!-- Toolbar -->
      <div class="toolbar">
        <Space wrap class="toolbar-filters">
          <Input
            v-model:value="searchText"
            placeholder="Cari siswa (nama/NIS/NISN)..."
            allow-clear
            style="width: 220px"
            @press-enter="handleSearch"
            :disabled="filterWithoutClass"
          >
            <template #prefix>
              <SearchOutlined />
            </template>
          </Input>
          <Select
            v-model:value="filterClassId"
            placeholder="Filter Kelas"
            allow-clear
            style="width: 140px"
            :loading="loadingClasses"
            :disabled="filterWithoutClass"
            @change="handleClassFilterChange"
          >
            <template #suffixIcon>
              <FilterOutlined />
            </template>
            <SelectOption v-for="cls in classes" :key="cls.id" :value="cls.id">
              {{ cls.name }}
            </SelectOption>
          </Select>
          <Checkbox v-model:checked="filterWithoutClass" @change="handleFilterWithoutClassChange">
            <Space>
              <ExclamationCircleOutlined style="color: #faad14" />
              <span>Tanpa Kelas</span>
            </Space>
          </Checkbox>
        </Space>
        <Space wrap class="toolbar-actions">
          <Button @click="loadStudents">
            <template #icon><ReloadOutlined /></template>
          </Button>
          <Button @click="handleDownloadTemplate">
            <template #icon><DownloadOutlined /></template>
            Template
          </Button>
          <Button @click="openImportModal">
            <template #icon><UploadOutlined /></template>
            Import
          </Button>
          <Button
            v-if="filterWithoutClass && selectedStudentIds.length > 0"
            type="primary"
            @click="openBulkAssignModal"
          >
            <template #icon><CheckCircleOutlined /></template>
            Assign ({{ selectedStudentIds.length }})
          </Button>
          <Button type="primary" @click="openCreateModal">
            <template #icon><PlusOutlined /></template>
            Tambah Siswa
          </Button>
        </Space>
      </div>

      <!-- Students Without Class Alert -->
      <Alert
        v-if="studentsWithoutClass > 0 && !filterWithoutClass"
        type="warning"
        show-icon
        style="margin-bottom: 16px"
      >
        <template #icon><ExclamationCircleOutlined /></template>
        <template #message>
          <Space>
            <span>Ada <strong>{{ studentsWithoutClass }}</strong> siswa yang belum memiliki kelas.</span>
            <Button type="link" size="small" @click="filterWithoutClass = true; handleFilterWithoutClassChange()">
              Lihat &amp; Assign Kelas
            </Button>
          </Space>
        </template>
      </Alert>

      <!-- Table -->
      <Table
        :columns="columns"
        :data-source="filteredStudents"
        :loading="loading"
        :pagination="{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: total,
          showSizeChanger: true,
          showTotal: (total: number) => `Total ${total} siswa`,
        }"
        :scroll="{ x: 1000 }"
        :row-selection="filterWithoutClass ? rowSelection : undefined"
        row-key="id"
        @change="handleTableChange"
        class="custom-table"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'nis'">
            <Text strong>{{ (record as Student).nis }}</Text>
          </template>
          <template v-else-if="column.key === 'nisn'">
             <Text type="secondary">{{ (record as Student).nisn }}</Text>
          </template>
          <template v-else-if="column.key === 'name'">
             <div class="student-name">
                <Text strong>{{ (record as Student).name }}</Text>
             </div>
          </template>
          <template v-else-if="column.key === 'className'">
            <span v-if="(record as Student).className" class="class-badge">
              {{ (record as Student).className }}
            </span>
            <span v-else class="no-class-badge">
              Belum ada kelas
            </span>
          </template>
          <template v-else-if="column.key === 'rfidCode'">
            <template v-if="(record as Student).rfidCode">
              <Tooltip :title="(record as Student).rfidCode">
                <div class="status-indicator success">
                   <WifiOutlined /> <span>Terpasang</span>
                </div>
              </Tooltip>
            </template>
            <template v-else>
               <div class="status-indicator">
                   <span class="text-secondary">Belum</span>
               </div>
            </template>
          </template>
          <template v-else-if="column.key === 'hasAccount'">
            <Tooltip :title="(record as Student).hasAccount ? 'Sudah punya akun' : 'Belum punya akun'">
              <div v-if="(record as Student).hasAccount" class="account-status active">
                <MobileOutlined />
              </div>
              <div v-else class="account-status inactive">
                <MobileOutlined />
              </div>
            </Tooltip>
          </template>
          <template v-else-if="column.key === 'isActive'">
            <div class="status-dot-wrapper">
              <span :class="['status-dot', (record as Student).isActive ? 'active' : 'inactive']"></span>
              <span>{{ (record as Student).isActive ? 'Aktif' : 'Nonaktif' }}</span>
            </div>
          </template>
          <template v-else-if="column.key === 'action'">
            <Space>
              <Tooltip title="Edit">
                <Button size="small" type="text" @click="openEditModal(record as Student)">
                  <template #icon><EditOutlined style="color: #64748b;" /></template>
                </Button>
              </Tooltip>
              <!-- RFID Pairing -->
              <template v-if="(record as Student).rfidCode">
                <Tooltip title="Hapus Kartu RFID">
                  <Popconfirm
                    title="Hapus kartu RFID dari siswa ini?"
                    description="Siswa tidak akan bisa absen dengan kartu ini lagi."
                    ok-text="Ya, Hapus"
                    cancel-text="Batal"
                    ok-type="danger"
                    @confirm="handleClearRFID(record as Student)"
                  >
                    <Button size="small" type="text" danger>
                      <template #icon><CloseCircleOutlined /></template>
                    </Button>
                  </Popconfirm>
                </Tooltip>
              </template>
              <template v-else>
                <Tooltip title="Pasangkan Kartu RFID">
                  <Button size="small" type="text" @click="openPairingModal(record as Student)">
                    <template #icon><WifiOutlined style="color: #3b82f6;" /></template>
                  </Button>
                </Tooltip>
              </template>
              <!-- Account Management -->
              <template v-if="(record as Student).hasAccount">
                <Tooltip title="Reset Password">
                  <Popconfirm
                    title="Reset password siswa ini?"
                    description="Password baru akan di-generate otomatis."
                    ok-text="Ya, Reset"
                    cancel-text="Batal"
                    @confirm="handleResetPassword(record as Student)"
                  >
                    <Button size="small" type="text">
                      <template #icon><KeyOutlined style="color: #f59e0b;" /></template>
                    </Button>
                  </Popconfirm>
                </Tooltip>
              </template>
              <template v-else>
                <Tooltip title="Buat Akun">
                  <Popconfirm
                    title="Buat akun untuk siswa ini?"
                    description="Akun akan dibuat dengan NIS sebagai username."
                    ok-text="Ya, Buat"
                    cancel-text="Batal"
                    @confirm="handleCreateAccount(record as Student)"
                  >
                    <Button size="small" type="text">
                      <template #icon><UserOutlined style="color: #8b5cf6;" /></template>
                    </Button>
                  </Popconfirm>
                </Tooltip>
              </template>
              <Tooltip title="Hapus">
                <Popconfirm
                  title="Hapus siswa ini?"
                  description="Data siswa akan dihapus permanen."
                  ok-text="Ya, Hapus"
                  cancel-text="Batal"
                  ok-type="danger"
                  @confirm="handleDelete(record as Student)"
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
      :title="isEditing ? 'Edit Siswa' : 'Tambah Siswa Baru'"
      :confirm-loading="modalLoading"
      :ok-text="isEditing ? 'Simpan' : 'Buat Siswa'"
      cancel-text="Batal"
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
        <FormItem label="Kelas" name="class_id" required>
          <Select
            v-model:value="formState.class_id"
            placeholder="Pilih kelas"
            :loading="loadingClasses"
            size="large"
          >
            <SelectOption v-for="cls in classes" :key="cls.id" :value="cls.id">
              {{ cls.name }}
            </SelectOption>
          </Select>
        </FormItem>
        <Row :gutter="16">
          <Col :span="12">
            <FormItem label="NIS" name="nis" required>
              <Input v-model:value="formState.nis" placeholder="NIS" size="large" />
            </FormItem>
          </Col>
          <Col :span="12">
            <FormItem label="NISN" name="nisn" required>
              <Input
                v-model:value="formState.nisn"
                placeholder="NISN"
                :disabled="isEditing"
                size="large"
              />
            </FormItem>
          </Col>
        </Row>
        <FormItem label="Nama Lengkap" name="name" required>
          <Input v-model:value="formState.name" placeholder="Nama lengkap siswa" size="large" />
        </FormItem>
        <FormItem label="Kode RFID" name="rfid_code">
          <Input v-model:value="formState.rfid_code" placeholder="Kode kartu RFID (opsional)" size="large" />
        </FormItem>
        <FormItem v-if="isEditing" label="Status" name="is_active">
           <div class="status-switch-wrapper">
             <Switch v-model:checked="formState.is_active" />
             <span :class="['status-label', formState.is_active ? 'active' : 'inactive']">
               {{ formState.is_active ? 'Siswa Aktif' : 'Siswa Nonaktif' }}
             </span>
           </div>
        </FormItem>
        <FormItem v-if="!isEditing" name="create_account">
          <Checkbox v-model:checked="formState.create_account">
            <Space>
              <MobileOutlined />
              <span>Buat akun untuk login mobile apps</span>
            </Space>
          </Checkbox>
          <div v-if="formState.create_account" class="account-info">
            <Text type="secondary">
              Username: NIS ({{ formState.nis || '...' }}), Password: <Text code>password123</Text>
            </Text>
          </div>
        </FormItem>
      </Form>
    </Modal>

    <!-- Credential Modal -->
    <Modal
      v-model:open="credentialModalVisible"
      title="Kredensial Akun Siswa"
      :footer="null"
      width="450px"
      wrap-class-name="modern-modal"
    >
      <div v-if="credentialData" class="credential-info">
        <div class="credential-header">
           <div class="success-icon-large">
             <UserOutlined />
           </div>
          <Title :level="4" style="margin: 16px 0 8px">{{ credentialData.name }}</Title>
          <Text type="secondary">Akun berhasil dibuat. Simpan kredensial berikut:</Text>
        </div>
        
        <div class="credential-box">
          <div class="credential-row">
            <span class="label">Username (NIS)</span>
            <div class="value-row">
               <span class="value">{{ credentialData.username }}</span>
               <Button size="small" type="text" @click="copyToClipboard(credentialData.username)">
                 <template #icon><CopyOutlined /></template>
               </Button>
            </div>
          </div>
          <div class="credential-row">
            <span class="label">Password</span>
             <div class="value-row">
               <span class="value code">{{ credentialData.password }}</span>
               <Button size="small" type="text" @click="copyToClipboard(credentialData.password)">
                 <template #icon><CopyOutlined /></template>
               </Button>
            </div>
          </div>
        </div>

        <Alert
          type="warning"
          message="Penting"
          description="Password ini hanya ditampilkan sekali. Simpan sekarang."
          show-icon
          style="text-align: left; margin-top: 16px; border-radius: 8px;"
        />

        <Button type="primary" block size="large" @click="credentialModalVisible = false" style="margin-top: 24px">
          Tutup
        </Button>
      </div>
    </Modal>

    <!-- Pairing Modal -->
    <Modal
      v-model:open="pairingModalVisible"
      title="Pasangkan Kartu RFID"
      :footer="null"
      width="500px"
      :maskClosable="false"
      wrap-class-name="modern-modal"
      @cancel="closePairingModal"
    >
      <div v-if="pairingStudent" class="pairing-container">
        <div class="pairing-header-content">
          <div class="pairing-icon-wrapper">
             <WifiOutlined />
          </div>
          <div class="pairing-details">
             <Title :level="4" style="margin: 0">{{ pairingStudent.name }}</Title>
             <Text type="secondary">NIS: {{ pairingStudent.nis }}</Text>
          </div>
        </div>

        <Divider style="margin: 16px 0;" />

        <!-- Device Selection -->
        <div v-if="!pairingSession?.active" class="pairing-device-select">
          <FormItem label="Pilih Perangkat Pairing" required style="margin-bottom: 8px;">
            <Select
              v-model:value="selectedDeviceId"
              placeholder="Pilih perangkat"
              :loading="loadingDevices"
              style="width: 100%"
              size="large"
            >
              <SelectOption v-for="device in devices" :key="device.id" :value="device.id">
                {{ device.deviceCode }} - {{ device.description || 'Tanpa deskripsi' }}
              </SelectOption>
            </Select>
          </FormItem>
          
          <Alert
            v-if="devices.length === 0 && !loadingDevices"
            type="warning"
            message="Tidak ada perangkat RFID"
            description="Hubungi Super Admin untuk mendaftarkan perangkat."
            show-icon
            style="margin: 16px 0"
          />

          <Button
            type="primary"
            block
            size="large"
            :loading="pairingLoading"
            :disabled="!selectedDeviceId"
            @click="startPairing"
            style="margin-top: 16px"
          >
            <template #icon><WifiOutlined /></template>
            Mulai Pairing
          </Button>
        </div>

        <!-- Pairing In Progress -->
        <div v-else class="pairing-progress">
          <div class="pairing-status-card">
              <Text strong>Menunggu Tap Kartu...</Text>
              <Text type="secondary">{{ pairingSession.message }}</Text>
          </div>

          <div class="countdown-circle">
            <Progress
              type="circle"
              :percent="Math.round((pairingCountdown / 60) * 100)"
              :format="() => `${pairingCountdown}s`"
              :status="pairingCountdown > 10 ? 'active' : 'exception'"
              :width="120"
            />
          </div>
          
           <Text type="secondary" class="instruction-text">
              Silakan tap kartu RFID pada perangkat sekarang.
           </Text>

          <Button
            danger
            block
            size="large"
            @click="cancelPairing"
            style="margin-top: 24px"
          >
            Batalkan
          </Button>
        </div>

        <!-- Pairing Failed/Expired -->
        <div v-if="pairingSession && !pairingSession.active && pairingCountdown === 0" class="pairing-result">
          <Alert
            type="warning"
            :message="pairingSession.message"
            show-icon
            style="margin-bottom: 16px"
          />
          <Button type="primary" block size="large" @click="pairingSession = null">
            Coba Lagi
          </Button>
        </div>
      </div>
    </Modal>

    <!-- Import Modal -->
    <Modal
      v-model:open="importModalVisible"
      title="Import Data Siswa"
      :footer="null"
      width="500px"
      wrap-class-name="modern-modal"
    >
      <div class="import-container">
        <Typography.Text type="secondary" style="display: block; margin-bottom: 16px;">
          Upload file Excel untuk import data siswa.
        </Typography.Text>

        <Upload.Dragger
          :custom-request="handleImportUpload"
          :show-upload-list="false"
          accept=".xlsx"
          :disabled="importLoading"
          class="modern-upload"
        >
          <p class="ant-upload-drag-icon">
            <UploadOutlined v-if="!importLoading" style="color: #f97316;" />
            <Progress v-else type="circle" :percent="0" :size="48" status="active" />
          </p>
          <p class="ant-upload-text">
            {{ importLoading ? 'Mengimpor data...' : 'Klik atau drag file Excel di sini' }}
          </p>
          <p class="ant-upload-hint">
            Maksimal 5MB. Format .xlsx
          </p>
        </Upload.Dragger>

        <div class="template-section">
            <Button type="link" @click="handleDownloadTemplate">
              <DownloadOutlined /> Download Template Excel
            </Button>
        </div>
      </div>
    </Modal>

    <!-- Import Result Modal -->
    <Modal
      v-model:open="importResultModalVisible"
      title="Hasil Import"
      :footer="null"
      width="600px"
      wrap-class-name="modern-modal"
      @cancel="closeImportResultModal"
    >
      <div v-if="importResult" class="import-result">
        <!-- Summary Stats -->
        <div class="result-stats-grid">
           <div class="stat-item total">
              <span class="stat-val">{{ importResult.total_rows }}</span>
              <span class="stat-lbl">Total</span>
           </div>
           <div class="stat-item success">
              <span class="stat-val">{{ importResult.success_count }}</span>
              <span class="stat-lbl">Berhasil</span>
           </div>
           <div class="stat-item warning">
              <span class="stat-val">{{ importResult.warning_count }}</span>
              <span class="stat-lbl">Warning</span>
           </div>
           <div class="stat-item error">
               <span class="stat-val">{{ importResult.failed_count }}</span>
               <span class="stat-lbl">Gagal</span>
           </div>
        </div>

        <Divider style="margin: 16px 0;" />

        <!-- Students Without Class Alert -->
        <Alert
          v-if="importResult.students_without_class && importResult.students_without_class > 0"
          type="warning"
          show-icon
          style="margin-bottom: 16px"
        >
          <template #message>
            <strong>{{ importResult.students_without_class }}</strong> siswa tanpa kelas.
            Anda dapat assign kelas nanti.
          </template>
        </Alert>

        <div v-if="importResult.errors && importResult.errors.length > 0" class="log-section">
           <Text strong type="danger">Error Log</Text>
           <div class="log-box">
              <div v-for="(err, idx) in importResult.errors" :key="idx" class="log-item error">
                 <span class="log-row">Baris {{ (err as ImportError).row }}</span>
                 <span class="log-msg">{{ (err as ImportError).message }}</span>
              </div>
           </div>
        </div>

        <Button type="primary" block size="large" @click="closeImportResultModal" style="margin-top: 16px;">
          Selesai
        </Button>
      </div>
    </Modal>

    <!-- Bulk Assign Class Modal -->
    <Modal
      v-model:open="bulkAssignModalVisible"
      title="Assign Kelas"
      :confirm-loading="bulkAssignLoading"
      @ok="handleBulkAssign"
      @cancel="closeBulkAssignModal"
      ok-text="Assign"
      cancel-text="Batal"
      width="400px"
      wrap-class-name="modern-modal"
    >
      <div class="bulk-assign-container">
        <Typography.Text style="display: block; margin-bottom: 16px;">
          Pilih kelas untuk <strong>{{ selectedStudentIds.length }}</strong> siswa terpilih.
        </Typography.Text>

        <Form layout="vertical">
          <FormItem label="Kelas Tujuan" required>
            <Select
              v-model:value="bulkAssignClassId"
              placeholder="Pilih kelas"
              :loading="loadingClasses"
              style="width: 100%"
              size="large"
            >
              <SelectOption v-for="cls in classes" :key="cls.id" :value="cls.id">
                {{ cls.name }} ({{ cls.studentCount || 0 }} siswa)
              </SelectOption>
            </Select>
          </FormItem>
        </Form>
      </div>
    </Modal>
  </div>
</template>

<style scoped>
.student-management {
  padding: 0;
}

.page-header {
  margin-bottom: 24px;
}

.toolbar {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.toolbar-filters {
  flex: 1;
  min-width: 0;
}

.toolbar-actions {
  flex-shrink: 0;
}

/* Custom Table Styles */
.custom-table :deep(.ant-table-thead > tr > th) {
  background: #f8fafc;
  color: #475569;
  font-weight: 600;
  border-bottom: 1px solid #f1f5f9;
}

.custom-table :deep(.ant-table-tbody > tr > td) {
  padding: 16px 16px;
  border-bottom: 1px solid #f1f5f9;
}

/* Badges & Indicators */
.class-badge {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  color: #475569;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  display: inline-block;
}

.no-class-badge {
  color: #f59e0b;
  font-size: 11px;
  background: #fffbeb;
  padding: 2px 8px;
  border-radius: 4px;
  border: 1px solid #fcd34d;
}

.status-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  color: #94a3b8;
}

.status-indicator.success {
  color: #22c55e;
}

.account-status {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 4px;
  margin: 0 auto;
}

.account-status.active {
  color: #3b82f6;
  background: #eff6ff;
}

.account-status.inactive {
  color: #94a3b8;
  background: #f1f5f9;
}

.status-dot-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 13px;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.status-dot.active { background-color: #22c55e; box-shadow: 0 0 0 2px rgba(34, 197, 94, 0.2); }
.status-dot.inactive { background-color: #94a3b8; }

.text-secondary { color: #94a3b8; font-size: 13px; }

/* Credential Modal */
.credential-header { text-align: center; margin-bottom: 24px; }
.success-icon-large { font-size: 48px; color: #10b981; background: #ecfdf5; width: 80px; height: 80px; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; }
.credential-box { background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 8px; padding: 16px; }
.credential-row { display: flex; justify-content: space-between; align-items: center; padding: 8px 0; border-bottom: 1px solid #f1f5f9; }
.credential-row:last-child { border-bottom: none; }
.label { color: #64748b; font-size: 13px; }
.value-row { display: flex; align-items: center; gap: 8px; }
.value { font-weight: 600; color: #0f172a; }
.value.code { font-family: monospace; font-size: 16px; background: #e2e8f0; padding: 2px 6px; border-radius: 4px; }

/* Pairing */
.pairing-container { text-align: center; }
.pairing-header-content { display: flex; align-items: center; gap: 16px; margin-bottom: 24px; text-align: left; background: #f8fafc; padding: 16px; border-radius: 12px; }
.pairing-icon-wrapper { width: 48px; height: 48px; background: #dbeafe; color: #3b82f6; border-radius: 12px; display: flex; align-items: center; justify-content: center; font-size: 24px; }
.pairing-progress { margin-top: 24px; }
.pairing-status-card { background: #eff6ff; border: 1px solid #dbeafe; border-radius: 8px; padding: 16px; margin-bottom: 24px; }
.countdown-circle { margin: 24px 0; display: flex; justify-content: center; }
.instruction-text { display: block; margin-top: 16px; font-size: 14px; }

/* Import Results */
.result-stats-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; margin-bottom: 24px; }
.stat-item { padding: 12px; border-radius: 8px; text-align: center; display: flex; flex-direction: column; }
.stat-item.total { background: #f1f5f9; color: #475569; }
.stat-item.success { background: #ecfdf5; color: #059669; }
.stat-item.warning { background: #fffbeb; color: #d97706; }
.stat-item.error { background: #fef2f2; color: #dc2626; }
.stat-val { font-size: 20px; font-weight: 700; display: block; }
.stat-lbl { font-size: 11px; text-transform: uppercase; letter-spacing: 0.5px; margin-top: 4px; }

.log-section { margin-top: 16px; }
.log-box { max-height: 150px; overflow-y: auto; background: #000; border-radius: 8px; padding: 12px; margin-top: 8px; }
.log-item { display: flex; justify-content: space-between; font-family: monospace; font-size: 12px; margin-bottom: 4px; }
.log-item.error { color: #fca5a5; }

/* Status Switch */
.status-switch-wrapper { display: flex; align-items: center; gap: 12px; }
.status-label { font-size: 14px; }
.status-label.active { color: #22c55e; font-weight: 500; }
.status-label.inactive { color: #64748b; }

@media (max-width: 768px) {
  .toolbar {
    flex-direction: column;
    align-items: stretch;
  }
  
  .toolbar-filters,
  .toolbar-actions {
    width: 100%;
  }
  
  .toolbar-actions {
    justify-content: flex-start;
  }
}
</style>

<!-- Global Styles for Modals (Duplicated here for safety if scoped logic fails, but relies on existing global styles) -->
<style>
/* Modern Modal Styles Overrides */
.modern-modal .ant-modal-content {
  border-radius: 16px !important;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1) !important;
  padding: 0 !important;
  overflow: hidden;
}
.modern-modal .ant-modal-header {
  border-bottom: 1px solid #f1f5f9 !important;
  padding: 20px 24px !important;
  background: #fff !important;
}
.modern-modal .ant-modal-title {
  font-size: 18px !important;
  font-weight: 600 !important;
  color: #0f172a !important;
}
.modern-modal .ant-modal-body {
  padding: 24px !important;
}
.modern-modal .ant-btn {
  border-radius: 8px !important;
  height: 40px !important;
  font-weight: 500 !important;
}
</style>
