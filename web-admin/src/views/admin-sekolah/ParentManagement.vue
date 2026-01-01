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
  Tooltip,
  Alert,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
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
} from '@ant-design/icons-vue'
import { schoolService } from '@/services'
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

// Modal state
const modalVisible = ref(false)
const modalLoading = ref(false)
const isEditing = ref(false)
const editingParent = ref<Parent | null>(null)

// Credential modal state
const credentialModalVisible = ref(false)
const credentialData = ref<{ username: string; password: string; name: string } | null>(null)

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
  student_ids: [{ required: true, message: 'Pilih minimal satu siswa', type: 'array' as const, min: 1 }],
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

// Filter option for student select
const filterStudentOption = (input: string, option: unknown) => {
  const opt = option as { label?: string }
  return opt?.label?.toLowerCase().includes(input.toLowerCase()) ?? false
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
const openEditModal = (parent: Parent) => {
  isEditing.value = true
  editingParent.value = parent
  formState.name = parent.name
  formState.phone = parent.phone || ''
  formState.email = parent.email || ''
  formState.student_ids = parent.studentIds || []
  modalVisible.value = true
}

// Reset form
const resetForm = () => {
  formState.name = ''
  formState.phone = ''
  formState.email = ''
  formState.student_ids = []
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
        <Col :xs="24" :sm="12" :md="8" class="toolbar-right">
          <Space>
            <Button @click="loadParents">
              <template #icon><ReloadOutlined /></template>
              Refresh
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
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <Space>
              <UserOutlined />
              {{ (record as Parent).name }}
            </Space>
          </template>
          <template v-else-if="column.key === 'phone'">
            <Tag color="blue">{{ (record as Parent).phone || '-' }}</Tag>
          </template>
          <template v-else-if="column.key === 'email'">
            {{ (record as Parent).email || '-' }}
          </template>
          <template v-else-if="column.key === 'studentNames'">
            <Space wrap>
              <Tag
                v-for="(name, index) in (record as Parent).studentNames"
                :key="index"
                color="green"
              >
                <LinkOutlined /> {{ name }}
              </Tag>
              <Tag v-if="!(record as Parent).studentNames?.length" color="default">
                Belum ada anak
              </Tag>
            </Space>
          </template>
          <template v-else-if="column.key === 'action'">
            <Space>
              <Tooltip title="Edit">
                <Button size="small" @click="openEditModal(record as Parent)">
                  <template #icon><EditOutlined /></template>
                </Button>
              </Tooltip>
              <Tooltip title="Reset Password">
                <Popconfirm
                  title="Reset password orang tua ini?"
                  description="Password baru akan di-generate otomatis."
                  ok-text="Ya, Reset"
                  cancel-text="Batal"
                  @confirm="handleResetPassword(record as Parent)"
                >
                  <Button size="small">
                    <template #icon><KeyOutlined /></template>
                  </Button>
                </Popconfirm>
              </Tooltip>
              <Tooltip title="Hapus">
                <Popconfirm
                  title="Hapus data orang tua ini?"
                  description="Akun orang tua akan dihapus."
                  ok-text="Ya, Hapus"
                  cancel-text="Batal"
                  @confirm="handleDelete(record as Parent)"
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
      :title="isEditing ? 'Edit Orang Tua' : 'Tambah Orang Tua Baru'"
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
        <FormItem label="Nama Lengkap" name="name" required>
          <Input v-model:value="formState.name" placeholder="Nama lengkap orang tua" />
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
              />
            </FormItem>
          </Col>
          <Col :span="12">
            <FormItem label="Email" name="email">
              <Input v-model:value="formState.email" placeholder="email@example.com" />
            </FormItem>
          </Col>
        </Row>
        <FormItem
          v-if="!isEditing"
        >
          <template #label>
            <Space>
              <span>Password</span>
              <Tag color="blue">Default</Tag>
            </Space>
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
          required
          extra="Pilih siswa yang merupakan anak dari orang tua ini"
        >
          <Select
            v-model:value="formState.student_ids"
            mode="multiple"
            placeholder="Pilih siswa"
            :loading="loadingStudents"
            show-search
            :filter-option="filterStudentOption"
          >
            <SelectOption
              v-for="student in students"
              :key="student.id"
              :value="student.id"
              :label="`${student.name} - ${student.className}`"
            >
              <div class="student-option">
                <span>{{ student.name }}</span>
                <Tag size="small" color="blue">{{ student.className }}</Tag>
              </div>
            </SelectOption>
          </Select>
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
  </div>
</template>

<style scoped>
.parent-management {
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

.student-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
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

.password-info {
  padding: 8px 12px;
  background: #f6ffed;
  border-radius: 4px;
  border: 1px solid #b7eb8f;
}

@media (max-width: 576px) {
  .toolbar-right {
    margin-top: 16px;
    justify-content: flex-start;
  }
}
</style>
