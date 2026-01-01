<script setup lang="ts">
import { ref, reactive, onMounted, computed, onUnmounted } from 'vue'
import {
  Table,
  Button,
  Input,
  Space,
  Tag,
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
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
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
} from '@ant-design/icons-vue'
import { schoolService } from '@/services'
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
  try {
    const response = await schoolService.getStudents({
      page: pagination.current,
      page_size: pagination.pageSize,
      search: searchText.value,
      class_id: filterClassId.value,
    })
    students.value = response.students
    total.value = response.pagination.total
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
        
        // Check if pairing was successful (student now has RFID)
        if (status.message.includes('berhasil')) {
          message.success('Kartu RFID berhasil dipasangkan!')
          pairingModalVisible.value = false
          loadStudents()
        } else {
          pairingSession.value = status
        }
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

onMounted(() => {
  loadStudents()
  loadClasses()
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
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col :xs="24" :sm="24" :md="16">
          <Space wrap>
            <Input
              v-model:value="searchText"
              placeholder="Cari siswa (nama/NIS/NISN)..."
              allow-clear
              style="width: 250px"
              @press-enter="handleSearch"
            >
              <template #prefix>
                <SearchOutlined />
              </template>
            </Input>
            <Select
              v-model:value="filterClassId"
              placeholder="Filter Kelas"
              allow-clear
              style="width: 150px"
              :loading="loadingClasses"
              @change="handleClassFilterChange"
            >
              <template #suffixIcon>
                <FilterOutlined />
              </template>
              <SelectOption v-for="cls in classes" :key="cls.id" :value="cls.id">
                {{ cls.name }}
              </SelectOption>
            </Select>
          </Space>
        </Col>
        <Col :xs="24" :sm="24" :md="8" class="toolbar-right">
          <Space>
            <Button @click="loadStudents">
              <template #icon><ReloadOutlined /></template>
              Refresh
            </Button>
            <Button type="primary" @click="openCreateModal">
              <template #icon><PlusOutlined /></template>
              Tambah Siswa
            </Button>
          </Space>
        </Col>
      </Row>

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
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'className'">
            <Tag color="blue">{{ (record as Student).className }}</Tag>
          </template>
          <template v-else-if="column.key === 'rfidCode'">
            <template v-if="(record as Student).rfidCode">
              <Tooltip :title="(record as Student).rfidCode">
                <Tag color="green">
                  <WifiOutlined /> Terpasang
                </Tag>
              </Tooltip>
            </template>
            <template v-else>
              <Tag color="default">Belum</Tag>
            </template>
          </template>
          <template v-else-if="column.key === 'hasAccount'">
            <Tooltip :title="(record as Student).hasAccount ? 'Sudah punya akun' : 'Belum punya akun'">
              <Tag :color="(record as Student).hasAccount ? 'green' : 'default'">
                <MobileOutlined />
              </Tag>
            </Tooltip>
          </template>
          <template v-else-if="column.key === 'isActive'">
            <Tag :color="(record as Student).isActive ? 'success' : 'default'">
              {{ (record as Student).isActive ? 'Aktif' : 'Nonaktif' }}
            </Tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <Space>
              <Tooltip title="Edit">
                <Button size="small" @click="openEditModal(record as Student)">
                  <template #icon><EditOutlined /></template>
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
                    @confirm="handleClearRFID(record as Student)"
                  >
                    <Button size="small" danger>
                      <template #icon><CloseCircleOutlined /></template>
                    </Button>
                  </Popconfirm>
                </Tooltip>
              </template>
              <template v-else>
                <Tooltip title="Pasangkan Kartu RFID">
                  <Button size="small" type="primary" ghost @click="openPairingModal(record as Student)">
                    <template #icon><WifiOutlined /></template>
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
                    <Button size="small">
                      <template #icon><KeyOutlined /></template>
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
                    <Button size="small" type="primary" ghost>
                      <template #icon><UserOutlined /></template>
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
                  @confirm="handleDelete(record as Student)"
                >
                  <Button size="small" danger>
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
      @ok="handleSubmit"
      @cancel="handleModalCancel"
    >
      <Form
        ref="formRef"
        :model="formState"
        :rules="formRules"
        layout="vertical"
        style="margin-top: 16px"
      >
        <FormItem label="Kelas" name="class_id" required>
          <Select
            v-model:value="formState.class_id"
            placeholder="Pilih kelas"
            :loading="loadingClasses"
          >
            <SelectOption v-for="cls in classes" :key="cls.id" :value="cls.id">
              {{ cls.name }}
            </SelectOption>
          </Select>
        </FormItem>
        <Row :gutter="16">
          <Col :span="12">
            <FormItem label="NIS" name="nis" required>
              <Input v-model:value="formState.nis" placeholder="Nomor Induk Siswa" />
            </FormItem>
          </Col>
          <Col :span="12">
            <FormItem label="NISN" name="nisn" required>
              <Input
                v-model:value="formState.nisn"
                placeholder="Nomor Induk Siswa Nasional"
                :disabled="isEditing"
              />
            </FormItem>
          </Col>
        </Row>
        <FormItem label="Nama Lengkap" name="name" required>
          <Input v-model:value="formState.name" placeholder="Nama lengkap siswa" />
        </FormItem>
        <FormItem label="Kode RFID" name="rfid_code">
          <Input v-model:value="formState.rfid_code" placeholder="Kode kartu RFID (opsional)" />
        </FormItem>
        <FormItem v-if="isEditing" label="Status" name="is_active">
          <Switch v-model:checked="formState.is_active" checked-children="Aktif" un-checked-children="Nonaktif" />
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
            <Text type="secondary">Username (NIS):</Text>
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
            ⚠️ Password ini hanya ditampilkan sekali. Pastikan untuk menyimpan atau memberikan ke siswa.
          </Text>
        </div>

        <Button type="primary" block @click="credentialModalVisible = false" style="margin-top: 16px">
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
      @cancel="closePairingModal"
    >
      <div v-if="pairingStudent" class="pairing-info">
        <div class="pairing-header">
          <WifiOutlined style="font-size: 48px; color: #1890ff" />
          <Title :level="4" style="margin: 16px 0 8px">{{ pairingStudent.name }}</Title>
          <Text type="secondary">NIS: {{ pairingStudent.nis }}</Text>
        </div>

        <!-- Device Selection -->
        <div v-if="!pairingSession?.active" class="pairing-device-select">
          <FormItem label="Pilih Perangkat RFID" required>
            <Select
              v-model:value="selectedDeviceId"
              placeholder="Pilih perangkat"
              :loading="loadingDevices"
              style="width: 100%"
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
            description="Hubungi Super Admin untuk mendaftarkan perangkat RFID ke sekolah Anda."
            show-icon
            style="margin-bottom: 16px"
          />

          <Button
            type="primary"
            block
            :loading="pairingLoading"
            :disabled="!selectedDeviceId"
            @click="startPairing"
          >
            <template #icon><WifiOutlined /></template>
            Mulai Pairing
          </Button>
        </div>

        <!-- Pairing In Progress -->
        <div v-else class="pairing-progress">
          <Alert
            type="info"
            message="Menunggu Tap Kartu..."
            :description="pairingSession.message"
            show-icon
            style="margin-bottom: 16px"
          />

          <div class="countdown-container">
            <Progress
              type="circle"
              :percent="Math.round((pairingCountdown / 60) * 100)"
              :format="() => `${pairingCountdown}s`"
              :status="pairingCountdown > 10 ? 'active' : 'exception'"
            />
            <Text type="secondary" style="margin-top: 8px">
              Silakan tap kartu RFID pada perangkat
            </Text>
          </div>

          <Button
            danger
            block
            @click="cancelPairing"
            style="margin-top: 16px"
          >
            <template #icon><CloseCircleOutlined /></template>
            Batalkan Pairing
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
          <Button type="primary" block @click="pairingSession = null">
            Coba Lagi
          </Button>
        </div>
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
  margin-bottom: 16px;
}

.toolbar-right {
  display: flex;
  justify-content: flex-end;
}

.account-info {
  margin-top: 8px;
  padding: 8px 12px;
  background: #f6ffed;
  border-radius: 4px;
  border: 1px solid #b7eb8f;
}

.credential-info {
  text-align: center;
}

.credential-header {
  margin-bottom: 24px;
}

.credential-card {
  text-align: left;
  margin-bottom: 16px;
}

.credential-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
}

.credential-item:not(:last-child) {
  border-bottom: 1px solid #f0f0f0;
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
}

.pairing-info {
  text-align: center;
}

.pairing-header {
  margin-bottom: 24px;
}

.pairing-device-select {
  text-align: left;
}

.pairing-progress {
  text-align: center;
}

.countdown-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24px 0;
}

.pairing-result {
  margin-top: 16px;
}

@media (max-width: 768px) {
  .toolbar-right {
    margin-top: 16px;
    justify-content: flex-start;
  }
}
</style>
