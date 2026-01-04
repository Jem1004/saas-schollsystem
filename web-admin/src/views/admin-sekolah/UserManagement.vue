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
  Switch,
  Alert,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  PlusOutlined,
  SearchOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
  KeyOutlined,
  UserOutlined,
  CheckOutlined,
} from '@ant-design/icons-vue'
import { schoolService } from '@/services'
import type { SchoolUser, Class, UpdateUserRequest } from '@/types/school'

const { Title, Text } = Typography

// Table state
const loading = ref(false)
const users = ref<SchoolUser[]>([])
const total = ref(0)
const pagination = reactive({
  current: 1,
  pageSize: 10,
})
const searchText = ref('')
const filterRole = ref<string | undefined>(undefined)

// Classes for wali kelas assignment
const classes = ref<Class[]>([])
const loadingClasses = ref(false)

// Modal state
const modalVisible = ref(false)
const modalLoading = ref(false)
const isEditing = ref(false)
const editingUser = ref<SchoolUser | null>(null)

// Password reset modal
const passwordModalVisible = ref(false)
const newPassword = ref('')

// Form state
const formRef = ref()
const formState = reactive({
  role: 'guru' as 'guru' | 'wali_kelas' | 'guru_bk',
  username: '',
  email: '',
  name: '',
  password: '',
  assigned_class_id: undefined as number | undefined,  // For wali_kelas
  assigned_class_ids: [] as number[],  // For guru_bk
  is_active: true,
})

// Form rules
const formRules = {
  role: [{ required: true, message: 'Role wajib dipilih' }],
  username: [{ required: true, message: 'Username wajib diisi' }],
  name: [{ required: true, message: 'Nama wajib diisi' }],
  password: [
    { required: true, message: 'Password wajib diisi' },
    { min: 8, message: 'Password minimal 8 karakter' },
  ],
}

// Role options (removed colors as we now use CSS classes)
const roleOptions = [
  { value: 'guru', label: 'Guru' },
  { value: 'wali_kelas', label: 'Wali Kelas' },
  { value: 'guru_bk', label: 'Guru BK' },
]

// Table columns
const columns: TableProps['columns'] = [
  {
    title: 'Username',
    dataIndex: 'username',
    key: 'username',
  },
  {
    title: 'Nama',
    dataIndex: 'name',
    key: 'name',
    sorter: true,
  },
  {
    title: 'Role',
    dataIndex: 'role',
    key: 'role',
    width: 120,
  },
  {
    title: 'Kelas',
    dataIndex: 'assignedClassName',
    key: 'assignedClassName',
    width: 100,
  },
  {
    title: 'Status',
    dataIndex: 'isActive',
    key: 'isActive',
    width: 100,
    align: 'center',
  },
  {
    title: 'Login Terakhir',
    dataIndex: 'lastLoginAt',
    key: 'lastLoginAt',
    width: 150,
  },
  {
    title: 'Aksi',
    key: 'action',
    width: 200,
    align: 'center',
  },
]

// Computed filtered data
const filteredUsers = computed(() => {
  let result = users.value
  
  if (filterRole.value) {
    result = result.filter(u => u.role === filterRole.value)
  }
  
  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(
      (user) =>
        user.username.toLowerCase().includes(search) ||
        user.name?.toLowerCase().includes(search) ||
        user.email?.toLowerCase().includes(search)
    )
  }
  
  return result
})

// Load users data
const loadUsers = async () => {
  loading.value = true
  try {
    const response = await schoolService.getUsers({
      page: pagination.current,
      page_size: pagination.pageSize,
      search: searchText.value,
    })
    users.value = response.users
    total.value = response.pagination.total
  } catch (err) {
    console.error('Failed to load users:', err)
    message.error('Gagal memuat data user')
    users.value = []
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
  loadUsers()
}

// Handle search
const handleSearch = () => {
  pagination.current = 1
  loadUsers()
}

// Handle role filter change
const handleRoleFilterChange = () => {
  pagination.current = 1
  loadUsers()
}

// Open modal for create
const openCreateModal = () => {
  isEditing.value = false
  editingUser.value = null
  resetForm()
  modalVisible.value = true
}

// Open modal for edit
const openEditModal = (user: SchoolUser) => {
  isEditing.value = true
  editingUser.value = user
  formState.role = user.role
  formState.username = user.username
  formState.email = user.email || ''
  formState.name = user.name || ''
  formState.password = ''
  formState.assigned_class_id = user.assignedClassId
  formState.assigned_class_ids = user.assignedClasses?.map(c => c.id) || []
  formState.is_active = user.isActive
  modalVisible.value = true
}

// Reset form
const resetForm = () => {
  formState.role = 'guru'
  formState.username = ''
  formState.email = ''
  formState.name = ''
  formState.password = ''
  formState.assigned_class_id = undefined
  formState.assigned_class_ids = []
  formState.is_active = true
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
    // Custom validation for editing (password not required)
    if (isEditing.value) {
      await formRef.value?.validate(['role', 'username', 'name'])
    } else {
      await formRef.value?.validate()
    }
  } catch {
    return
  }

  modalLoading.value = true
  try {
    if (isEditing.value && editingUser.value) {
      const updateData: UpdateUserRequest = {
        email: formState.email || undefined,
        name: formState.name || undefined,
        is_active: formState.is_active,
        assigned_class_id: formState.role === 'wali_kelas' ? formState.assigned_class_id : undefined,
        assigned_class_ids: formState.role === 'guru_bk' ? formState.assigned_class_ids : undefined,
      }
      await schoolService.updateUser(editingUser.value.id, updateData)
      message.success('User berhasil diperbarui')
    } else {
      await schoolService.createUser({
        role: formState.role,
        username: formState.username,
        email: formState.email || undefined,
        name: formState.name || undefined,
        password: formState.password,
        assigned_class_id: formState.role === 'wali_kelas' ? formState.assigned_class_id : undefined,
        assigned_class_ids: formState.role === 'guru_bk' ? formState.assigned_class_ids : undefined,
      })
      message.success('User berhasil ditambahkan')
    }
    modalVisible.value = false
    resetForm()
    loadUsers()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Terjadi kesalahan')
  } finally {
    modalLoading.value = false
  }
}

// Filter option for class select
const filterClassOption = (input: string, option: unknown) => {
  const opt = option as { children?: string }
  return opt.children?.toLowerCase().includes(input.toLowerCase()) ?? false
}

// Handle delete user
const handleDelete = async (user: SchoolUser) => {
  try {
    await schoolService.deleteUser(user.id)
    message.success(`User ${user.username} berhasil dihapus`)
    loadUsers()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Gagal menghapus user')
  }
}

// Handle reset password
const handleResetPassword = async (user: SchoolUser) => {
  try {
    const result = await schoolService.resetUserPassword(user.id)
    newPassword.value = result.temporaryPassword
    passwordModalVisible.value = true
    message.success('Password berhasil direset')
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Gagal mereset password')
  }
}

// Get role label and color
const getRoleInfo = (role: string) => {
  return roleOptions.find(r => r.value === role) || { label: role, color: 'default' }
}

// Format date
const formatDate = (dateStr?: string): string => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

onMounted(() => {
  loadUsers()
  loadClasses()
})
</script>

<template>
  <div class="user-management">
    <div class="page-header">
      <Title :level="2" style="margin: 0">Manajemen User</Title>
    </div>

    <Card>
      <!-- Toolbar -->
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col :xs="24" :sm="24" :md="16">
          <Space wrap>
            <Input
              v-model:value="searchText"
              placeholder="Cari user..."
              allow-clear
              style="width: 200px"
              @press-enter="handleSearch"
            >
              <template #prefix>
                <SearchOutlined />
              </template>
            </Input>
            <Select
              v-model:value="filterRole"
              placeholder="Filter Role"
              allow-clear
              style="width: 150px"
              @change="handleRoleFilterChange"
            >
              <SelectOption v-for="role in roleOptions" :key="role.value" :value="role.value">
                {{ role.label }}
              </SelectOption>
            </Select>
          </Space>
        </Col>
        <Col :xs="24" :sm="24" :md="8" class="toolbar-right">
          <Space>
            <Button @click="loadUsers">
              <template #icon><ReloadOutlined /></template>
            </Button>
            <Button type="primary" @click="openCreateModal">
              <template #icon><PlusOutlined /></template>
              Tambah User
            </Button>
          </Space>
        </Col>
      </Row>

      <!-- Table -->
      <Table
        :columns="columns"
        :data-source="filteredUsers"
        :loading="loading"
        :pagination="{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: total,
          showSizeChanger: true,
          showTotal: (total: number) => `Total ${total} user`,
        }"
        :scroll="{ x: 1000 }"
        row-key="id"
        @change="handleTableChange"
        class="custom-table"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'username'">
            <div class="user-cell">
              <div class="user-avatar">
                <UserOutlined />
              </div>
              <div class="user-info">
                 <Text strong>{{ (record as SchoolUser).username }}</Text>
                 <Text type="secondary" style="font-size: 12px">{{ (record as SchoolUser).email || '-' }}</Text>
              </div>
            </div>
          </template>
          <template v-else-if="column.key === 'role'">
            <span :class="['role-badge', `role-${(record as SchoolUser).role}`]">
              {{ getRoleInfo((record as SchoolUser).role).label }}
            </span>
          </template>
          <template v-else-if="column.key === 'assignedClassName'">
            <!-- For wali_kelas: single class -->
            <span v-if="(record as SchoolUser).assignedClassName" class="class-badge">
              {{ (record as SchoolUser).assignedClassName }}
            </span>
            <!-- For guru_bk: multiple classes -->
            <template v-else-if="(record as SchoolUser).assignedClasses?.length">
              <div class="class-tags-wrapper">
                <span v-for="cls in (record as SchoolUser).assignedClasses" :key="cls.id" class="class-badge">
                  {{ cls.name }}
                </span>
              </div>
            </template>
            <span v-else class="text-secondary">-</span>
          </template>
          <template v-else-if="column.key === 'isActive'">
            <div class="status-indicator">
              <span :class="['status-dot', (record as SchoolUser).isActive ? 'active' : 'inactive']"></span>
              <span>{{ (record as SchoolUser).isActive ? 'Aktif' : 'Nonaktif' }}</span>
            </div>
          </template>
          <template v-else-if="column.key === 'lastLoginAt'">
            <span class="text-secondary">{{ formatDate((record as SchoolUser).lastLoginAt) }}</span>
          </template>
          <template v-else-if="column.key === 'action'">
            <Space>
              <Button type="text" @click="openEditModal(record as SchoolUser)">
                <template #icon><EditOutlined style="color: #64748b;" /></template>
              </Button>
              <Popconfirm
                title="Reset password user ini?"
                description="Password baru akan digenerate."
                ok-text="Ya, Reset"
                cancel-text="Batal"
                @confirm="handleResetPassword(record as SchoolUser)"
              >
                <Button type="text">
                  <template #icon><KeyOutlined style="color: #f59e0b;" /></template>
                </Button>
              </Popconfirm>
              <Popconfirm
                title="Hapus user ini?"
                description="User akan dihapus permanen."
                ok-text="Ya, Hapus"
                cancel-text="Batal"
                ok-type="danger"
                @confirm="handleDelete(record as SchoolUser)"
              >
                <Button type="text" danger>
                  <template #icon><DeleteOutlined /></template>
                </Button>
              </Popconfirm>
            </Space>
          </template>
        </template>
      </Table>
    </Card>

    <!-- Create/Edit Modal -->
    <Modal
      v-model:open="modalVisible"
      :title="isEditing ? 'Edit User' : 'Tambah User Baru'"
      :confirm-loading="modalLoading"
      :ok-text="isEditing ? 'Simpan' : 'Buat User'"
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
        <FormItem label="Role" name="role" required>
          <Select 
            v-model:value="formState.role" 
            placeholder="Pilih role" 
            :disabled="isEditing"
          >
            <SelectOption v-for="role in roleOptions" :key="role.value" :value="role.value">
              {{ role.label }}
            </SelectOption>
          </Select>
        </FormItem>
        <Row :gutter="16">
          <Col :span="12">
            <FormItem label="Username" name="username" required>
              <Input
                v-model:value="formState.username"
                placeholder="Username"
                :disabled="isEditing"
              />
            </FormItem>
          </Col>
          <Col :span="12">
            <FormItem label="Email" name="email">
              <Input v-model:value="formState.email" placeholder="email@sekolah.id" />
            </FormItem>
          </Col>
        </Row>
        <FormItem label="Nama Lengkap" name="name" required>
          <Input v-model:value="formState.name" placeholder="Nama lengkap" />
        </FormItem>
        <FormItem
          v-if="!isEditing"
          label="Password"
          name="password"
          required
        >
          <Input.Password v-model:value="formState.password" placeholder="Minimal 8 karakter" />
          <Text type="secondary" style="font-size: 12px; margin-top: 4px; display: block;">User akan diminta mengganti password saat login pertama</Text>
        </FormItem>
        <FormItem
          v-if="formState.role === 'wali_kelas'"
          label="Kelas yang Diampu"
          name="assigned_class_id"
        >
          <Select
            v-model:value="formState.assigned_class_id"
            placeholder="Pilih kelas"
            allow-clear
            :loading="loadingClasses"
          >
            <SelectOption v-for="cls in classes" :key="cls.id" :value="cls.id">
              {{ cls.name }}
            </SelectOption>
          </Select>
        </FormItem>
        <FormItem
          v-if="formState.role === 'guru_bk'"
          label="Kelas yang Ditangani"
          name="assigned_class_ids"
        >
          <Select
            v-model:value="formState.assigned_class_ids"
            mode="multiple"
            placeholder="Pilih kelas"
            allow-clear
            :loading="loadingClasses"
            :filter-option="filterClassOption"
          >
            <SelectOption v-for="cls in classes" :key="cls.id" :value="cls.id">
              {{ cls.name }}
            </SelectOption>
          </Select>
        </FormItem>
        <FormItem v-if="isEditing" label="Status" name="is_active">
           <div class="status-switch-wrapper">
             <Switch v-model:checked="formState.is_active" />
             <span :class="['status-label', formState.is_active ? 'active' : 'inactive']">
               {{ formState.is_active ? 'User Aktif' : 'User Nonaktif' }}
             </span>
           </div>
        </FormItem>
      </Form>
    </Modal>

    <!-- Password Reset Result Modal -->
    <Modal
      v-model:open="passwordModalVisible"
      title="Password Berhasil Direset"
      :footer="null"
      width="400px"
      wrap-class-name="modern-modal"
    >
      <div class="password-result">
        <div class="success-icon">
          <CheckOutlined />
        </div>
        <Text strong style="font-size: 16px; margin-bottom: 8px; display: block;">Password Baru</Text>
        <div class="password-display">
          <Text strong copyable class="password-text">{{ newPassword }}</Text>
        </div>
        <Alert
          message="Penting"
          description="Salin password ini dan berikan kepada user. Password harus diganti saat login pertama."
          type="warning"
          show-icon
          style="text-align: left; border-radius: 8px;"
        />
      </div>
    </Modal>
  </div>
</template>

<style scoped>
.user-management {
  padding: 0;
}

.page-header {
  margin-bottom: 24px;
}

.toolbar {
  margin-bottom: 24px;
}

.toolbar-right {
  display: flex;
  justify-content: flex-end;
}

/* Custom Component Styles */
.user-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  width: 32px;
  height: 32px;
  background: #f1f5f9;
  color: #64748b;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #e2e8f0;
}

.user-info {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
}

.text-secondary {
  color: #94a3b8;
  font-size: 13px;
}

/* Role Badges */
.role-badge {
  font-size: 13px;
  font-weight: 500;
}

.role-guru {
  color: #475569;
}

.role-wali_kelas {
  color: #c2410c;
}

.role-guru_bk {
  color: #0369a1;
}

/* Class Tags */
.class-badge {
  color: #475569;
  font-size: 13px;
}

.class-tags-wrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

/* Status Indicator */
.status-indicator {
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

.status-dot.active {
  background-color: #22c55e;
  box-shadow: 0 0 0 2px rgba(34, 197, 94, 0.2);
}

.status-dot.inactive {
  background-color: #94a3b8;
}

/* Password Result */
.password-result {
  text-align: center;
  padding: 16px 0;
}

.success-icon {
  width: 48px;
  height: 48px;
  background: #ecfdf5;
  color: #10b981;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  margin: 0 auto 16px auto;
}

.password-display {
  background: #f8fafc;
  padding: 16px;
  border-radius: 8px;
  margin-bottom: 24px;
  border: 1px solid #e2e8f0;
}

.password-text {
  font-family: monospace;
  font-size: 20px;
  color: #0f172a;
}

.status-switch-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
}

.status-label {
  font-size: 14px;
}

.status-label.active { color: #22c55e; font-weight: 500; }
.status-label.inactive { color: #64748b; }

@media (max-width: 768px) {
  .toolbar-right {
    margin-top: 16px;
    justify-content: flex-start;
  }
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
</style>

<!-- Global Styles for Modals -->
<style>
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
