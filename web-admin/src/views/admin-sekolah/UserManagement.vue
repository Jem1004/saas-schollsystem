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
  Switch,
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

// Role options
const roleOptions = [
  { value: 'guru', label: 'Guru', color: 'blue' },
  { value: 'wali_kelas', label: 'Wali Kelas', color: 'green' },
  { value: 'guru_bk', label: 'Guru BK', color: 'purple' },
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
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'username'">
            <Space>
              <UserOutlined />
              {{ (record as SchoolUser).username }}
            </Space>
          </template>
          <template v-else-if="column.key === 'role'">
            <Tag :color="getRoleInfo((record as SchoolUser).role).color">
              {{ getRoleInfo((record as SchoolUser).role).label }}
            </Tag>
          </template>
          <template v-else-if="column.key === 'assignedClassName'">
            <!-- For wali_kelas: single class -->
            <Tag v-if="(record as SchoolUser).assignedClassName" color="blue">
              {{ (record as SchoolUser).assignedClassName }}
            </Tag>
            <!-- For guru_bk: multiple classes -->
            <template v-else-if="(record as SchoolUser).assignedClasses?.length">
              <Tag v-for="cls in (record as SchoolUser).assignedClasses" :key="cls.id" color="purple" style="margin-bottom: 2px;">
                {{ cls.name }}
              </Tag>
            </template>
            <span v-else>-</span>
          </template>
          <template v-else-if="column.key === 'isActive'">
            <Tag :color="(record as SchoolUser).isActive ? 'success' : 'default'">
              {{ (record as SchoolUser).isActive ? 'Aktif' : 'Nonaktif' }}
            </Tag>
          </template>
          <template v-else-if="column.key === 'lastLoginAt'">
            {{ formatDate((record as SchoolUser).lastLoginAt) }}
          </template>
          <template v-else-if="column.key === 'action'">
            <Space>
              <Button size="small" @click="openEditModal(record as SchoolUser)">
                <template #icon><EditOutlined /></template>
                Edit
              </Button>
              <Popconfirm
                title="Reset password user ini?"
                description="Password baru akan digenerate."
                ok-text="Ya, Reset"
                cancel-text="Batal"
                @confirm="handleResetPassword(record as SchoolUser)"
              >
                <Button size="small">
                  <template #icon><KeyOutlined /></template>
                </Button>
              </Popconfirm>
              <Popconfirm
                title="Hapus user ini?"
                description="User akan dihapus permanen."
                ok-text="Ya, Hapus"
                cancel-text="Batal"
                @confirm="handleDelete(record as SchoolUser)"
              >
                <Button size="small" danger>
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
        <FormItem label="Role" name="role" required>
          <Select v-model:value="formState.role" placeholder="Pilih role" :disabled="isEditing">
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
                placeholder="Username untuk login"
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
          extra="User akan diminta mengganti password saat login pertama"
        >
          <Input.Password v-model:value="formState.password" placeholder="Minimal 8 karakter" />
        </FormItem>
        <FormItem
          v-if="formState.role === 'wali_kelas'"
          label="Kelas yang Diampu"
          name="assigned_class_id"
          extra="Pilih kelas yang akan diampu sebagai wali kelas"
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
          extra="Pilih satu atau lebih kelas yang akan ditangani oleh guru BK"
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
          <Switch v-model:checked="formState.is_active" checked-children="Aktif" un-checked-children="Nonaktif" />
        </FormItem>
      </Form>
    </Modal>

    <!-- Password Reset Result Modal -->
    <Modal
      v-model:open="passwordModalVisible"
      title="Password Berhasil Direset"
      :footer="null"
    >
      <div class="password-result">
        <Text>Password baru:</Text>
        <div class="password-display">
          <Text strong copyable>{{ newPassword }}</Text>
        </div>
        <Text type="secondary">
          Salin password ini dan berikan kepada user. Password harus diganti saat login pertama.
        </Text>
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
  margin-bottom: 16px;
}

.toolbar-right {
  display: flex;
  justify-content: flex-end;
}

.password-result {
  text-align: center;
  padding: 16px 0;
}

.password-display {
  background: #f5f5f5;
  padding: 16px;
  border-radius: 8px;
  margin: 16px 0;
  font-size: 18px;
}

@media (max-width: 768px) {
  .toolbar-right {
    margin-top: 16px;
    justify-content: flex-start;
  }
}
</style>
