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
  message,
  Popconfirm,
  Card,
  Row,
  Col,
  Typography,
  Divider,
  Descriptions,
  DescriptionsItem,
  Tooltip,
  Collapse,
  CollapsePanel,
  Statistic,
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
  adminUsername: '',
  adminPassword: '',
  adminName: '',
  adminEmail: '',
})

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
    width: 280,
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
      <Title :level="2" style="margin: 0">Manajemen Sekolah</Title>
    </div>

    <Card>
      <!-- Toolbar -->
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col :xs="24" :sm="12" :md="8">
          <Input
            v-model:value="searchText"
            placeholder="Cari sekolah..."
            allow-clear
            @press-enter="handleSearch"
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
              Refresh
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
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'isActive'">
            <Tag :color="(record as School).isActive ? 'success' : 'default'">
              {{ (record as School).isActive ? 'Aktif' : 'Nonaktif' }}
            </Tag>
          </template>
          <template v-else-if="column.key === 'createdAt'">
            {{ formatDate((record as School).createdAt) }}
          </template>
          <template v-else-if="column.key === 'action'">
            <Space>
              <Tooltip title="Lihat Detail">
                <Button size="small" @click="openDetailModal(record as School)">
                  <template #icon><InfoCircleOutlined /></template>
                </Button>
              </Tooltip>
              <Button size="small" @click="openEditModal(record as School)">
                <template #icon><EditOutlined /></template>
                Edit
              </Button>
              <Popconfirm
                v-if="(record as School).isActive"
                title="Nonaktifkan sekolah ini?"
                description="Semua user di sekolah ini tidak akan bisa mengakses sistem."
                ok-text="Ya, Nonaktifkan"
                cancel-text="Batal"
                @confirm="handleDeactivate(record as School)"
              >
                <Button size="small" danger>
                  <template #icon><StopOutlined /></template>
                </Button>
              </Popconfirm>
              <Popconfirm
                v-else
                title="Aktifkan sekolah ini?"
                ok-text="Ya, Aktifkan"
                cancel-text="Batal"
                @confirm="handleActivate(record as School)"
              >
                <Button size="small" type="primary" ghost>
                  <template #icon><CheckCircleOutlined /></template>
                </Button>
              </Popconfirm>
              <Tooltip title="Hapus Permanen">
                <Button size="small" danger type="primary" @click="openDeleteModal(record as School)">
                  <template #icon><DeleteOutlined /></template>
                </Button>
              </Tooltip>
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
        <FormItem label="Nama Sekolah" name="name" required>
          <Input v-model:value="formState.name" placeholder="Masukkan nama sekolah" />
        </FormItem>
        <FormItem label="Alamat" name="address">
          <Input.TextArea
            v-model:value="formState.address"
            placeholder="Masukkan alamat sekolah"
            :rows="3"
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

        <!-- Admin Section (only for create) -->
        <template v-if="!isEditing">
          <Divider orientation="left">
            <Space>
              <UserOutlined />
              <span>Akun Admin Sekolah</span>
            </Space>
          </Divider>
          
          <Collapse ghost>
            <CollapsePanel key="admin" header="Kustomisasi Akun Admin (Opsional)">
              <Text type="secondary" style="display: block; margin-bottom: 16px;">
                Jika dikosongkan, username dan password akan di-generate otomatis.
              </Text>
              <Row :gutter="16">
                <Col :span="12">
                  <FormItem label="Username Admin" name="adminUsername">
                    <Input 
                      v-model:value="formState.adminUsername" 
                      placeholder="Otomatis jika kosong"
                    />
                  </FormItem>
                </Col>
                <Col :span="12">
                  <FormItem label="Password Admin" name="adminPassword">
                    <Input.Password 
                      v-model:value="formState.adminPassword" 
                      placeholder="Otomatis jika kosong"
                    />
                  </FormItem>
                </Col>
              </Row>
              <Row :gutter="16">
                <Col :span="12">
                  <FormItem label="Nama Admin" name="adminName">
                    <Input 
                      v-model:value="formState.adminName" 
                      placeholder="Otomatis jika kosong"
                    />
                  </FormItem>
                </Col>
                <Col :span="12">
                  <FormItem label="Email Admin" name="adminEmail">
                    <Input 
                      v-model:value="formState.adminEmail" 
                      placeholder="Sama dengan email sekolah jika kosong"
                    />
                  </FormItem>
                </Col>
              </Row>
            </CollapsePanel>
          </Collapse>
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
      width="500px"
      @cancel="credentialsModalVisible = false"
    >
      <div v-if="adminCredentials" class="credentials-modal">
        <div class="credentials-warning">
          <Text type="warning">
            ⚠️ {{ adminCredentials.message }}
          </Text>
        </div>

        <Descriptions :column="1" bordered size="small" style="margin-top: 16px;">
          <DescriptionsItem label="Username">
            <div class="credential-item">
              <Text code>{{ adminCredentials.username }}</Text>
              <Tooltip title="Salin Username">
                <Button 
                  size="small" 
                  type="text" 
                  @click="copyToClipboard(adminCredentials.username, 'Username')"
                >
                  <template #icon><CopyOutlined /></template>
                </Button>
              </Tooltip>
            </div>
          </DescriptionsItem>
          <DescriptionsItem label="Password">
            <div class="credential-item">
              <Text code>
                {{ showPassword ? adminCredentials.password : '••••••••••••' }}
              </Text>
              <Space>
                <Tooltip :title="showPassword ? 'Sembunyikan' : 'Tampilkan'">
                  <Button size="small" type="text" @click="showPassword = !showPassword">
                    <template #icon>
                      <EyeInvisibleOutlined v-if="showPassword" />
                      <EyeOutlined v-else />
                    </template>
                  </Button>
                </Tooltip>
                <Tooltip title="Salin Password">
                  <Button 
                    size="small" 
                    type="text" 
                    @click="copyToClipboard(adminCredentials.password, 'Password')"
                  >
                    <template #icon><CopyOutlined /></template>
                  </Button>
                </Tooltip>
              </Space>
            </div>
          </DescriptionsItem>
          <DescriptionsItem label="Nama">
            {{ adminCredentials.name }}
          </DescriptionsItem>
          <DescriptionsItem v-if="adminCredentials.email" label="Email">
            {{ adminCredentials.email }}
          </DescriptionsItem>
        </Descriptions>

        <div class="credentials-actions">
          <Button 
            type="primary" 
            block 
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
      width="600px"
      @cancel="detailModalVisible = false"
    >
      <div v-if="detailLoading" style="text-align: center; padding: 40px;">
        Memuat...
      </div>
      <div v-else-if="schoolDetail" class="detail-modal">
        <Descriptions :column="2" bordered size="small">
          <DescriptionsItem label="Nama Sekolah" :span="2">
            {{ schoolDetail.name }}
          </DescriptionsItem>
          <DescriptionsItem label="Email">
            {{ schoolDetail.email || '-' }}
          </DescriptionsItem>
          <DescriptionsItem label="Telepon">
            {{ schoolDetail.phone || '-' }}
          </DescriptionsItem>
          <DescriptionsItem label="Alamat" :span="2">
            {{ schoolDetail.address || '-' }}
          </DescriptionsItem>
          <DescriptionsItem label="Status">
            <Tag :color="schoolDetail.isActive ? 'success' : 'default'">
              {{ schoolDetail.isActive ? 'Aktif' : 'Nonaktif' }}
            </Tag>
          </DescriptionsItem>
          <DescriptionsItem label="Tanggal Dibuat">
            {{ formatDate(schoolDetail.createdAt) }}
          </DescriptionsItem>
        </Descriptions>

        <!-- Stats -->
        <div v-if="schoolDetail.stats" style="margin-top: 24px;">
          <Divider orientation="left">Statistik</Divider>
          <Row :gutter="16">
            <Col :span="6">
              <Statistic title="User" :value="schoolDetail.stats.totalUsers" />
            </Col>
            <Col :span="6">
              <Statistic title="Siswa" :value="schoolDetail.stats.totalStudents" />
            </Col>
            <Col :span="6">
              <Statistic title="Kelas" :value="schoolDetail.stats.totalClasses" />
            </Col>
            <Col :span="6">
              <Statistic title="Device" :value="schoolDetail.stats.totalDevices" />
            </Col>
          </Row>
        </div>

        <!-- Admin Users -->
        <div v-if="schoolDetail.admins && schoolDetail.admins.length > 0" style="margin-top: 24px;">
          <Divider orientation="left">
            <Space>
              <UserOutlined />
              <span>Admin Sekolah</span>
            </Space>
          </Divider>
          <div v-for="admin in schoolDetail.admins" :key="admin.id" class="admin-item">
            <Descriptions :column="2" size="small" bordered>
              <DescriptionsItem label="Username">
                <div class="credential-item">
                  <Text code>{{ admin.username }}</Text>
                  <Tooltip title="Salin Username">
                    <Button 
                      size="small" 
                      type="text" 
                      @click="copyToClipboard(admin.username, 'Username')"
                    >
                      <template #icon><CopyOutlined /></template>
                    </Button>
                  </Tooltip>
                </div>
              </DescriptionsItem>
              <DescriptionsItem label="Status">
                <Tag :color="admin.isActive ? 'success' : 'default'">
                  {{ admin.isActive ? 'Aktif' : 'Nonaktif' }}
                </Tag>
              </DescriptionsItem>
              <DescriptionsItem label="Nama">
                {{ admin.name }}
              </DescriptionsItem>
              <DescriptionsItem label="Email">
                {{ admin.email || '-' }}
              </DescriptionsItem>
            </Descriptions>
          </div>
        </div>
      </div>
    </Modal>

    <!-- Delete Confirmation Modal -->
    <Modal
      v-model:open="deleteModalVisible"
      title="Hapus Sekolah Permanen"
      :confirm-loading="deleteLoading"
      ok-text="Hapus Permanen"
      ok-type="danger"
      :ok-button-props="{ disabled: !canDelete }"
      @ok="handleDelete"
      @cancel="deleteModalVisible = false"
    >
      <div v-if="schoolToDelete" class="delete-modal">
        <Alert
          type="error"
          show-icon
          style="margin-bottom: 16px;"
        >
          <template #icon><ExclamationCircleOutlined /></template>
          <template #message>
            <Text strong>Peringatan: Tindakan ini tidak dapat dibatalkan!</Text>
          </template>
          <template #description>
            <div>
              Menghapus sekolah <Text strong>{{ schoolToDelete.name }}</Text> akan menghapus:
              <ul style="margin: 8px 0; padding-left: 20px;">
                <li>Semua data user (admin, guru, siswa, orang tua)</li>
                <li>Semua data siswa dan kelas</li>
                <li>Semua data absensi</li>
                <li>Semua data BK (pelanggaran, prestasi, izin, konseling)</li>
                <li>Semua data nilai dan catatan wali kelas</li>
                <li>Semua device yang terdaftar</li>
              </ul>
            </div>
          </template>
        </Alert>

        <div style="margin-top: 16px;">
          <Text>
            Untuk mengkonfirmasi, ketik nama sekolah: <Text strong code>{{ schoolToDelete.name }}</Text>
          </Text>
          <Input
            v-model:value="deleteConfirmText"
            placeholder="Ketik nama sekolah untuk konfirmasi"
            style="margin-top: 8px;"
          />
        </div>
      </div>
    </Modal>
  </div>
</template>

<style scoped>
.tenant-management {
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

@media (max-width: 576px) {
  .toolbar-right {
    margin-top: 16px;
    justify-content: flex-start;
  }
}

.credentials-modal {
  padding: 8px 0;
}

.credentials-warning {
  padding: 12px;
  background: #fffbe6;
  border-radius: 6px;
  border: 1px solid #ffe58f;
}

.credential-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.credentials-actions {
  margin-top: 24px;
}

.detail-modal {
  padding: 8px 0;
}

.admin-item {
  margin-bottom: 12px;
}

.admin-item:last-child {
  margin-bottom: 0;
}

.delete-modal {
  padding: 8px 0;
}
</style>
