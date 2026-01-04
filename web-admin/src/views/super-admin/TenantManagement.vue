<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
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
  Divider,
  Tooltip,
  Collapse,
  CollapsePanel,
  Alert,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  PlusOutlined,
  SearchOutlined,
  EditOutlined,
  CheckCircleOutlined,
  StopOutlined,
  ReloadOutlined,
  CopyOutlined,
  EyeOutlined,
  EyeInvisibleOutlined,
  UserOutlined,
  DeleteOutlined,
  InfoCircleOutlined,
  ExclamationCircleOutlined,
} from '@ant-design/icons-vue'
import { tenantService } from '@/services'
import type { School, SchoolDetail, CreateSchoolRequest, UpdateSchoolRequest, AdminCredentials } from '@/types/tenant'

const { Title, Text } = Typography

// Table state
const loading = ref(false)
const schools = ref<School[]>([])
const total = ref(0)
const pagination = reactive({
  current: 1,
  pageSize: 10,
})
const searchText = ref('')

// Modal state
const modalVisible = ref(false)
const modalLoading = ref(false)
const isEditing = ref(false)
const editingSchool = ref<School | null>(null)

// Admin credentials modal state
const credentialsModalVisible = ref(false)
const adminCredentials = ref<AdminCredentials | null>(null)
const showPassword = ref(false)

// School detail modal state
const detailModalVisible = ref(false)
const detailLoading = ref(false)
const schoolDetail = ref<SchoolDetail | null>(null)

// Delete confirmation state
const deleteModalVisible = ref(false)
const deleteLoading = ref(false)
const schoolToDelete = ref<School | null>(null)
const deleteConfirmText = ref('')

// Form state
const formRef = ref()
const formState = reactive<CreateSchoolRequest>({
  name: '',
  address: '',
  phone: '',
  email: '',
  timezone: 'Asia/Makassar',
  adminUsername: '',
  adminPassword: '',
  adminName: '',
  adminEmail: '',
})

// Timezone options
const timezoneOptions = [
  { value: 'Asia/Jakarta', label: 'WIB - Waktu Indonesia Barat (UTC+7)' },
  { value: 'Asia/Makassar', label: 'WITA - Waktu Indonesia Tengah (UTC+8)' },
  { value: 'Asia/Jayapura', label: 'WIT - Waktu Indonesia Timur (UTC+9)' },
]

// Form rules
const formRules = {
  name: [{ required: true, message: 'Nama sekolah wajib diisi' }],
  email: [{ type: 'email' as const, message: 'Format email tidak valid' }],
  adminEmail: [{ type: 'email' as const, message: 'Format email tidak valid' }],
  adminUsername: [
    { pattern: /^[a-zA-Z0-9_]*$/, message: 'Username hanya boleh huruf, angka, dan underscore' },
  ],
}

// Table columns
const columns: TableProps['columns'] = [
  {
    title: 'Nama Sekolah',
    dataIndex: 'name',
    key: 'name',
    sorter: true,
  },
  {
    title: 'Email',
    dataIndex: 'email',
    key: 'email',
  },
  {
    title: 'Telepon',
    dataIndex: 'phone',
    key: 'phone',
  },
  {
    title: 'Status',
    dataIndex: 'isActive',
    key: 'isActive',
    width: 100,
    align: 'center',
  },
  {
    title: 'Tanggal Dibuat',
    dataIndex: 'createdAt',
    key: 'createdAt',
    width: 150,
  },
  {
    title: 'Aksi',
    key: 'action',
    width: 250, // Reduced slightly
    align: 'center',
  },
]

// Computed filtered data
const filteredSchools = computed(() => {
  return schools.value
})

// Check if delete is allowed
const canDelete = computed(() => {
  if (!schoolToDelete.value) return false
  return deleteConfirmText.value === schoolToDelete.value.name
})

// Load schools data
const loadSchools = async () => {
  loading.value = true
  try {
    const response = await tenantService.getSchools({
      page: pagination.current,
      pageSize: pagination.pageSize,
      search: searchText.value,
    })
    schools.value = response.data
    total.value = response.total
  } catch {
    message.error('Gagal memuat data sekolah')
  } finally {
    loading.value = false
  }
}

// Handle table change
const handleTableChange: TableProps['onChange'] = (pag) => {
  pagination.current = pag.current || 1
  pagination.pageSize = pag.pageSize || 10
  loadSchools()
}

// Handle search
const handleSearch = () => {
  pagination.current = 1
  loadSchools()
}

// Open modal for create
const openCreateModal = () => {
  isEditing.value = false
  editingSchool.value = null
  resetForm()
  modalVisible.value = true
}

// Open modal for edit
const openEditModal = (school: School) => {
  isEditing.value = true
  editingSchool.value = school
  formState.name = school.name
  formState.address = school.address || ''
  formState.phone = school.phone || ''
  formState.email = school.email || ''
  formState.timezone = school.timezone || 'Asia/Makassar'
  formState.adminUsername = ''
  formState.adminPassword = ''
  formState.adminName = ''
  formState.adminEmail = ''
  modalVisible.value = true
}

// Open detail modal
const openDetailModal = async (school: School) => {
  detailLoading.value = true
  detailModalVisible.value = true
  try {
    schoolDetail.value = await tenantService.getSchoolDetail(school.id)
  } catch {
    message.error('Gagal memuat detail sekolah')
    detailModalVisible.value = false
  } finally {
    detailLoading.value = false
  }
}

// Open delete modal
const openDeleteModal = (school: School) => {
  schoolToDelete.value = school
  deleteConfirmText.value = ''
  deleteModalVisible.value = true
}

// Reset form
const resetForm = () => {
  formState.name = ''
  formState.address = ''
  formState.phone = ''
  formState.email = ''
  formState.timezone = 'Asia/Makassar'
  formState.adminUsername = ''
  formState.adminPassword = ''
  formState.adminName = ''
  formState.adminEmail = ''
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
    if (isEditing.value && editingSchool.value) {
      const updateData: UpdateSchoolRequest = {
        name: formState.name,
        address: formState.address || undefined,
        phone: formState.phone || undefined,
        email: formState.email || undefined,
        timezone: formState.timezone || undefined,
      }
      await tenantService.updateSchool(editingSchool.value.id, updateData)
      message.success('Sekolah berhasil diperbarui')
    } else {
      const result = await tenantService.createSchool(formState)
      message.success('Sekolah berhasil ditambahkan')
      
      // Show admin credentials modal
      if (result.admin) {
        adminCredentials.value = result.admin
        showPassword.value = true
        credentialsModalVisible.value = true
      }
    }
    modalVisible.value = false
    resetForm()
    loadSchools()
  } catch (error: unknown) {
    const err = error as { message?: string; response?: { data?: { error?: { message?: string } } } }
    const errorMessage = err.message || err.response?.data?.error?.message || 'Terjadi kesalahan'
    message.error(errorMessage)
  } finally {
    modalLoading.value = false
  }
}

// Handle activate school
const handleActivate = async (school: School) => {
  try {
    await tenantService.activateSchool(school.id)
    message.success(`Sekolah ${school.name} berhasil diaktifkan`)
    loadSchools()
  } catch (error: unknown) {
    const err = error as { message?: string; response?: { data?: { error?: { message?: string } } } }
    const errorMessage = err.message || err.response?.data?.error?.message || 'Gagal mengaktifkan sekolah'
    message.error(errorMessage)
  }
}

// Handle deactivate school
const handleDeactivate = async (school: School) => {
  try {
    await tenantService.deactivateSchool(school.id)
    message.success(`Sekolah ${school.name} berhasil dinonaktifkan`)
    loadSchools()
  } catch (error: unknown) {
    const err = error as { message?: string; response?: { data?: { error?: { message?: string } } } }
    const errorMessage = err.message || err.response?.data?.error?.message || 'Gagal menonaktifkan sekolah'
    message.error(errorMessage)
  }
}

// Handle delete school
const handleDelete = async () => {
  if (!schoolToDelete.value || !canDelete.value) return
  
  deleteLoading.value = true
  try {
    const result = await tenantService.deleteSchool(schoolToDelete.value.id)
    message.success(result.message)
    deleteModalVisible.value = false
    schoolToDelete.value = null
    deleteConfirmText.value = ''
    loadSchools()
  } catch (error: unknown) {
    const err = error as { message?: string; response?: { data?: { error?: { message?: string } } } }
    const errorMessage = err.message || err.response?.data?.error?.message || 'Gagal menghapus sekolah'
    message.error(errorMessage)
  } finally {
    deleteLoading.value = false
  }
}

// Copy to clipboard
const copyToClipboard = async (text: string, label: string) => {
  try {
    await navigator.clipboard.writeText(text)
    message.success(`${label} berhasil disalin`)
  } catch {
    message.error('Gagal menyalin')
  }
}

// Format date
const formatDate = (dateStr: string): string => {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

onMounted(() => {
  loadSchools()
})
</script>

<template>
  <div class="tenant-management">
    <div class="page-header">
      <div class="header-content">
        <Title :level="2" class="page-title">Manajemen Sekolah</Title>
        <Text class="page-subtitle">Kelola data sekolah dan admin</Text>
      </div>
    </div>

    <Card class="content-card" :bordered="false">
      <!-- Toolbar -->
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col :xs="24" :sm="12" :md="10">
          <Input
            v-model:value="searchText"
            placeholder="Cari nama sekolah..."
            allow-clear
            @press-enter="handleSearch"
            class="search-input"
          >
            <template #prefix>
              <SearchOutlined />
            </template>
          </Input>
        </Col>
        <Col :xs="24" :sm="12" :md="8" class="toolbar-right">
          <Space>
            <Button @click="loadSchools">
              <template #icon><ReloadOutlined /></template>
            </Button>
            <Button type="primary" @click="openCreateModal">
              <template #icon><PlusOutlined /></template>
              Tambah Sekolah
            </Button>
          </Space>
        </Col>
      </Row>

      <!-- Table -->
      <Table
        :columns="columns"
        :data-source="filteredSchools"
        :loading="loading"
        :pagination="{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: total,
          showSizeChanger: true,
          showTotal: (total: number) => `Total ${total} sekolah`,
        }"
        :scroll="{ x: 1000 }"
        row-key="id"
        @change="handleTableChange"
        class="custom-table"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'isActive'">
            <Tag :color="(record as School).isActive ? 'success' : 'default'" :bordered="false">
              {{ (record as School).isActive ? 'Aktif' : 'Nonaktif' }}
            </Tag>
          </template>
          <template v-else-if="column.key === 'createdAt'">
            <span class="date-text">{{ formatDate((record as School).createdAt) }}</span>
          </template>
          <template v-else-if="column.key === 'action'">
            <Space :size="4">
              <Tooltip title="Lihat Detail">
                <Button type="text" @click="openDetailModal(record as School)">
                  <template #icon><InfoCircleOutlined style="color: #64748b" /></template>
                </Button>
              </Tooltip>
              <Tooltip title="Edit">
                <Button type="text" @click="openEditModal(record as School)">
                  <template #icon><EditOutlined style="color: #f97316" /></template>
                </Button>
              </Tooltip>
              <Popconfirm
                v-if="(record as School).isActive"
                title="Nonaktifkan sekolah ini?"
                description="Semua user di sekolah ini tidak akan bisa mengakses sistem."
                ok-text="Ya, Nonaktifkan"
                cancel-text="Batal"
                @confirm="handleDeactivate(record as School)"
              >
                <Tooltip title="Nonaktifkan">
                  <Button type="text" danger>
                    <template #icon><StopOutlined /></template>
                  </Button>
                </Tooltip>
              </Popconfirm>
              <Popconfirm
                v-else
                title="Aktifkan sekolah ini?"
                ok-text="Ya, Aktifkan"
                cancel-text="Batal"
                @confirm="handleActivate(record as School)"
              >
                <Tooltip title="Aktifkan">
                  <Button type="text" style="color: #22c55e">
                    <template #icon><CheckCircleOutlined /></template>
                  </Button>
                </Tooltip>
              </Popconfirm>
              <Popconfirm
                title="Hapus Sekolah Permanen?"
                ok-text="Ya, Hapus"
                cancel-text="Batal"
                ok-type="danger"
                @confirm="openDeleteModal(record as School)"
              >
                <Tooltip title="Hapus Permanen">
                  <Button type="text" danger>
                    <template #icon><DeleteOutlined /></template>
                  </Button>
                </Tooltip>
              </Popconfirm>
            </Space>
          </template>
        </template>
      </Table>
    </Card>


    <!-- Create/Edit Modal -->
    <Modal
      v-model:open="modalVisible"
      :title="isEditing ? 'Edit Sekolah' : 'Tambah Sekolah Baru'"
      :confirm-loading="modalLoading"
      width="600px"
      wrap-class-name="modern-modal"
      :ok-text="isEditing ? 'Simpan Perubahan' : 'Buat Sekolah'"
      cancel-text="Batal"
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
        <FormItem label="Nama Sekolah" name="name" required>
          <Input v-model:value="formState.name" placeholder="Masukkan nama sekolah" />
        </FormItem>
        <FormItem label="Alamat" name="address">
          <Input.TextArea
            v-model:value="formState.address"
            placeholder="Masukkan alamat sekolah"
            :rows="3"
            class="modern-textarea"
          />
        </FormItem>
        <Row :gutter="16">
          <Col :span="12">
            <FormItem label="Telepon" name="phone">
              <Input v-model:value="formState.phone" placeholder="Contoh: 021-1234567" />
            </FormItem>
          </Col>
          <Col :span="12">
            <FormItem label="Email Sekolah" name="email">
              <Input v-model:value="formState.email" placeholder="email@sekolah.sch.id" />
            </FormItem>
          </Col>
        </Row>
        <FormItem label="Zona Waktu" name="timezone">
          <Select v-model:value="formState.timezone" placeholder="Pilih zona waktu">
            <SelectOption 
              v-for="tz in timezoneOptions" 
              :key="tz.value" 
              :value="tz.value"
            >
              {{ tz.label }}
            </SelectOption>
          </Select>
        </FormItem>

        <!-- Admin Section (only for create) -->
        <template v-if="!isEditing">
          <Divider orientation="left" style="margin: 24px 0 16px 0; font-size: 14px; color: #94a3b8;">
            <Space>
              <UserOutlined />
              <span>Akun Admin Sekolah</span>
            </Space>
          </Divider>
          
          <div class="admin-section-bg">
            <Collapse ghost :bordered="false" class="minimal-collapse">
              <CollapsePanel key="admin" header="Kustomisasi Akun Admin (Opsional)">
                <Text type="secondary" style="display: block; margin-bottom: 16px; font-size: 13px;">
                  Biarkan kosong untuk generate otomatis.
                </Text>
                <Row :gutter="16">
                  <Col :span="12">
                    <FormItem label="Username Admin" name="adminUsername">
                      <Input 
                        v-model:value="formState.adminUsername" 
                        placeholder="Otomatis"
                      />
                    </FormItem>
                  </Col>
                  <Col :span="12">
                    <FormItem label="Password Admin" name="adminPassword">
                      <Input.Password 
                        v-model:value="formState.adminPassword" 
                        placeholder="Otomatis"
                      />
                    </FormItem>
                  </Col>
                </Row>
                <Row :gutter="16">
                  <Col :span="12">
                    <FormItem label="Nama Admin" name="adminName">
                      <Input 
                        v-model:value="formState.adminName" 
                        placeholder="Otomatis"
                      />
                    </FormItem>
                  </Col>
                  <Col :span="12">
                    <FormItem label="Email Admin" name="adminEmail">
                      <Input 
                        v-model:value="formState.adminEmail" 
                        placeholder="Otomatis"
                      />
                    </FormItem>
                  </Col>
                </Row>
              </CollapsePanel>
            </Collapse>
          </div>
        </template>
      </Form>
    </Modal>

    <!-- Admin Credentials Modal -->
    <Modal
      v-model:open="credentialsModalVisible"
      title="Kredensial Admin Sekolah"
      :footer="null"
      :closable="true"
      :mask-closable="false"
      width="480px"
      wrap-class-name="modern-modal"
      @cancel="credentialsModalVisible = false"
    >
      <div v-if="adminCredentials" class="credentials-modal-content">
        <div class="credentials-warning">
          <ExclamationCircleOutlined style="color: #f59e0b; font-size: 20px; margin-right: 12px;" />
          <div>
            <span style="font-weight: 600; color: #d97706; display: block; margin-bottom: 4px;">Penting!</span>
            <span style="font-size: 13px; color: #b45309;">Simpan kredensial ini. Password tidak akan ditampilkan lagi.</span>
          </div>
        </div>

        <div class="credentials-card">
          <div class="credential-row">
            <span class="label">Username</span>
            <div class="value-group">
              <span class="value code">{{ adminCredentials.username }}</span>
              <Button type="text" @click="copyToClipboard(adminCredentials.username, 'Username')">
                <template #icon><CopyOutlined /></template>
              </Button>
            </div>
          </div>
          <div class="divider"></div>
          <div class="credential-row">
            <span class="label">Password</span>
            <div class="value-group">
              <span class="value code">{{ showPassword ? adminCredentials.password : '••••••••••••' }}</span>
              <Space :size="4">
                <Button type="text" @click="showPassword = !showPassword">
                  <template #icon>
                    <EyeInvisibleOutlined v-if="showPassword" />
                    <EyeOutlined v-else />
                  </template>
                </Button>
                <Button type="text" @click="copyToClipboard(adminCredentials.password, 'Password')">
                  <template #icon><CopyOutlined /></template>
                </Button>
              </Space>
            </div>
          </div>
        </div>

        <div class="credential-details">
          <div class="detail-item">
            <span class="label">Nama:</span>
            <span class="value">{{ adminCredentials.name }}</span>
          </div>
          <div class="detail-item" v-if="adminCredentials.email">
            <span class="label">Email:</span>
            <span class="value">{{ adminCredentials.email }}</span>
          </div>
        </div>

        <div class="credentials-actions">
          <Button 
            type="primary" 
            block 
            class="copy-all-btn"
            @click="copyToClipboard(`Username: ${adminCredentials.username}\nPassword: ${adminCredentials.password}`, 'Kredensial')"
          >
            <template #icon><CopyOutlined /></template>
            Salin Semua Kredensial
          </Button>
        </div>
      </div>
    </Modal>

    <!-- School Detail Modal -->
    <Modal
      v-model:open="detailModalVisible"
      title="Detail Sekolah"
      :footer="null"
      width="700px"
      wrap-class-name="modern-modal"
      @cancel="detailModalVisible = false"
    >
      <div v-if="detailLoading" class="loading-state">
         <ReloadOutlined spin /> Memuat...
      </div>
      <div v-else-if="schoolDetail" class="school-detail-content">
        <div class="detail-header">
           <div class="school-avatar">
             <BankOutlined />
           </div>
           <div class="school-title-info">
             <h3>{{ schoolDetail.name }}</h3>
             <Text type="secondary">{{ schoolDetail.address || 'Alamat tidak tersedia' }}</Text>
             <div style="margin-top: 8px;">
               <Tag :color="schoolDetail.isActive ? 'success' : 'default'" :bordered="false" class="status-tag">
                 {{ schoolDetail.isActive ? 'Aktif' : 'Nonaktif' }}
               </Tag>
             </div>
           </div>
        </div>

        <div class="detail-grid">
           <div class="detail-box">
              <span class="label">Email</span>
              <span class="value">{{ schoolDetail.email || '-' }}</span>
           </div>
           <div class="detail-box">
              <span class="label">Telepon</span>
              <span class="value">{{ schoolDetail.phone || '-' }}</span>
           </div>
           <div class="detail-box">
              <span class="label">Zona Waktu</span>
              <span class="value">{{ schoolDetail.timezone || 'WITA' }}</span>
           </div>
           <div class="detail-box">
              <span class="label">Bergabung</span>
              <span class="value">{{ formatDate(schoolDetail.createdAt) }}</span>
           </div>
        </div>

        <!-- Stats -->
        <div class="stats-section" v-if="schoolDetail.stats">
          <Row :gutter="16">
            <Col :span="6">
              <div class="stat-mini-card">
                <span class="stat-value">{{ schoolDetail.stats.totalUsers || 0 }}</span>
                <span class="stat-label">User</span>
              </div>
            </Col>
            <Col :span="6">
              <div class="stat-mini-card">
                <span class="stat-value">{{ schoolDetail.stats.totalStudents || 0 }}</span>
                <span class="stat-label">Siswa</span>
              </div>
            </Col>
            <Col :span="6">
              <div class="stat-mini-card">
                <span class="stat-value">{{ schoolDetail.stats.totalClasses || 0 }}</span>
                <span class="stat-label">Kelas</span>
              </div>
            </Col>
            <Col :span="6">
              <div class="stat-mini-card">
                <span class="stat-value">{{ schoolDetail.stats.totalDevices || 0 }}</span>
                <span class="stat-label">Device</span>
              </div>
            </Col>
          </Row>
        </div>

        <!-- Admin Users -->
        <div v-if="schoolDetail.admins && schoolDetail.admins.length > 0" class="admin-list-section">
          <Text strong style="margin-bottom: 12px; display: block;">Admin Sekolah</Text>
          <div v-for="admin in schoolDetail.admins" :key="admin.id" class="admin-list-item">
            <div class="admin-info">
              <div class="admin-avatar">
                <UserOutlined />
              </div>
              <div>
                <span class="admin-name">{{ admin.name }}</span>
                <span class="admin-username">@{{ admin.username }}</span>
              </div>
            </div>
            <Tag :color="admin.isActive ? 'success' : 'default'" :bordered="false">{{ admin.isActive ? 'Aktif' : 'Nonaktif' }}</Tag>
          </div>
        </div>
      </div>
    </Modal>

    <!-- Delete Confirmation Modal -->
    <Modal
      v-model:open="deleteModalVisible"
      title="Hapus Sekolah Permanen"
      :confirm-loading="deleteLoading"
      ok-text="Ya, Hapus Permanen"
      cancel-text="Batal"
      ok-type="danger"
      wrap-class-name="modern-modal"
      :ok-button-props="{ disabled: !canDelete }"
      @ok="handleDelete"
      @cancel="deleteModalVisible = false"
    >
      <div v-if="schoolToDelete" class="delete-confirmation-content">
        <Alert
          message="Tindakan ini tidak dapat dibatalkan"
          description="Menghapus sekolah akan menghapus seluruh data yang terkait secara permanen."
          type="error"
          show-icon
          class="delete-alert"
        />

        <div class="confirm-input-wrapper">
          <Text class="confirm-label">
            Ketik <Text strong code>{{ schoolToDelete.name }}</Text> untuk konfirmasi:
          </Text>
          <Input 
            v-model:value="deleteConfirmText" 
            :placeholder="schoolToDelete.name"
            class="confirm-input"
          />
        </div>
      </div>
    </Modal>
  </div>
</template>

<style scoped>
.tenant-management {
  padding: 0;
  max-width: 1600px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 32px;
}

.page-title {
  font-weight: 600 !important;
  color: #1e293b !important;
  margin-bottom: 8px !important;
  letter-spacing: -0.5px;
}

.page-subtitle {
  font-size: 14px;
  color: #64748b;
}

.content-card {
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
  border: 1px solid #f1f5f9;
}

.toolbar {
  margin-bottom: 24px;
}

.toolbar-right {
  display: flex;
  justify-content: flex-end;
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

.date-text {
  color: #64748b;
  font-size: 13px;
}

/* Form Styles inside Modal */
.admin-section-bg {
  background: #f8fafc;
  border-radius: 8px;
  padding: 16px;
  border: 1px solid #f1f5f9;
}

.minimal-collapse :deep(.ant-collapse-header) {
  padding: 0 !important;
}

.minimal-collapse :deep(.ant-collapse-content-box) {
  padding: 16px 0 0 0 !important;
}

/* Credentials Modal Styles */
.credentials-modal-content {
  padding-top: 8px;
}

.credentials-warning {
  display: flex;
  align-items: flex-start;
  padding: 16px;
  background: #fffbeb;
  border-radius: 8px;
  margin-bottom: 20px;
  border: 1px solid #fef3c7;
}

.credentials-card {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  margin-bottom: 20px;
}

.credential-row {
  display: flex;
  flex-direction: column;
  padding: 12px 16px;
}

.credential-row .label {
  font-size: 12px;
  color: #64748b;
  margin-bottom: 4px;
  text-transform: uppercase;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.credential-row .value-group {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.credential-row .value {
  font-family: 'SF Mono', 'Roboto Mono', monospace;
  font-size: 15px;
  color: #0f172a;
  font-weight: 500;
}

.credentials-card .divider {
  height: 1px;
  background: #e2e8f0;
  margin: 0 16px;
}

.credential-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 24px;
  padding: 0 4px;
}

.detail-item {
  display: flex;
  font-size: 14px;
}

.detail-item .label {
  color: #64748b;
  width: 60px;
}

.detail-item .value {
  color: #334155;
  font-weight: 500;
}

.copy-all-btn {
  height: 48px; 
  border-radius: 8px; 
  font-weight: 500;
  box-shadow: 0 4px 6px -1px rgba(249, 115, 22, 0.2);
}

/* School Detail Modal Styles */
.loading-state {
  text-align: center;
  padding: 40px;
  color: #64748b;
}

.detail-header {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 24px;
}

.school-avatar {
  width: 64px;
  height: 64px;
  background: #fff7ed;
  color: #f97316;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
}

.school-title-info h3 {
  margin: 0;
  font-size: 20px;
  color: #1e293b;
  font-weight: 600;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  padding: 20px;
  background: #f8fafc;
  border-radius: 12px;
  margin-bottom: 24px;
}

.detail-box {
  display: flex;
  flex-direction: column;
}

.detail-box .label {
  font-size: 12px;
  color: #94a3b8;
  margin-bottom: 4px;
}

.detail-box .value {
  font-size: 14px;
  color: #334155;
  font-weight: 500;
}

.stat-mini-card {
  background: white;
  border: 1px solid #f1f5f9;
  border-radius: 8px;
  padding: 12px;
  text-align: center;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.stat-value {
  font-size: 18px;
  font-weight: 700;
  color: #f97316;
}

.stat-label {
  font-size: 12px;
  color: #64748b;
}

.admin-list-section {
  margin-top: 24px;
}

.admin-list-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  border: 1px solid #f1f5f9;
  border-radius: 8px;
  margin-bottom: 8px;
}

.admin-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.admin-avatar {
  width: 36px;
  height: 36px;
  background: #f1f5f9;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #64748b;
}

.admin-name {
  display: block;
  font-weight: 500;
  color: #334155;
  font-size: 14px;
}

.admin-username {
  display: block;
  font-size: 12px;
  color: #94a3b8;
}

/* Delete Modal */
.delete-confirmation-content {
  padding-top: 8px;
}

.delete-alert {
  margin-bottom: 20px;
  border-radius: 8px;
}

.confirm-input-wrapper {
  background: #fdf2f2;
  padding: 16px;
  border-radius: 8px;
  border: 1px solid #fee2e2;
}

.confirm-label {
  display: block;
  margin-bottom: 12px;
  color: #7f1d1d;
}

.confirm-input {
  border-color: #fca5a5;
}

.confirm-input:hover, .confirm-input:focus {
  border-color: #ef4444;
  box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.1);
}

/* Responsiveness */
@media (max-width: 768px) {
  .toolbar-right {
    margin-top: 16px;
    justify-content: flex-start;
  }
}
</style>

<!-- Global Styles for Modals (Teleported) -->
<style>
/* Modern Modal Global Overrides */
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

.modern-modal .ant-modal-footer {
  border-top: 1px solid #f1f5f9 !important;
  padding: 16px 24px !important;
  background: #ffffff !important;
}

.modern-modal .ant-btn {
  border-radius: 8px !important;
  height: 40px !important;
  font-weight: 500 !important;
}

.modern-modal .ant-btn-primary {
  box-shadow: 0 4px 6px -1px rgba(249, 115, 22, 0.2) !important;
}
</style>
