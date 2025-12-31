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
} from '@ant-design/icons-vue'
import { schoolService } from '@/services'
import type { Parent, Student, CreateParentRequest, UpdateParentRequest } from '@/types/school'

const { Title } = Typography

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

// Form state
const formRef = ref()
const formState = reactive<CreateParentRequest>({
  name: '',
  phone: '',
  email: '',
  password: '',
  studentIds: [],
})

// Form rules
const formRules = {
  name: [{ required: true, message: 'Nama orang tua wajib diisi' }],
  phone: [{ required: true, message: 'Nomor telepon wajib diisi' }],
  password: [
    { required: true, message: 'Password wajib diisi' },
    { min: 6, message: 'Password minimal 6 karakter' },
  ],
  studentIds: [{ required: true, message: 'Pilih minimal satu siswa', type: 'array' as const, min: 1 }],
}

// Mock data for development
const mockParents: Parent[] = [
  { id: 1, schoolId: 1, userId: 101, name: 'Pak Ahmad', phone: '081234567890', email: 'ahmad@email.com', studentIds: [1], studentNames: ['Ahmad Fauzi'], createdAt: '2024-01-15T10:00:00Z', updatedAt: '2024-01-15T10:00:00Z' },
  { id: 2, schoolId: 1, userId: 102, name: 'Bu Siti', phone: '081234567891', email: 'siti@email.com', studentIds: [2, 3], studentNames: ['Budi Santoso', 'Citra Dewi'], createdAt: '2024-01-15T10:00:00Z', updatedAt: '2024-01-15T10:00:00Z' },
  { id: 3, schoolId: 1, userId: 103, name: 'Pak Darmawan', phone: '081234567892', email: 'darmawan@email.com', studentIds: [4], studentNames: ['Dian Pratama'], createdAt: '2024-01-15T10:00:00Z', updatedAt: '2024-01-15T10:00:00Z' },
  { id: 4, schoolId: 1, userId: 104, name: 'Bu Eka', phone: '081234567893', email: 'eka@email.com', studentIds: [5], studentNames: ['Eka Putri'], createdAt: '2024-01-15T10:00:00Z', updatedAt: '2024-01-15T10:00:00Z' },
  { id: 5, schoolId: 1, userId: 105, name: 'Pak Fajar', phone: '081234567894', studentIds: [6], studentNames: ['Fajar Nugroho'], createdAt: '2024-01-15T10:00:00Z', updatedAt: '2024-01-15T10:00:00Z' },
]

const mockStudents: Student[] = [
  { id: 1, schoolId: 1, classId: 1, className: 'VII-A', nis: '2024001', nisn: '0012345678', name: 'Ahmad Fauzi', isActive: true, createdAt: '2024-01-15', updatedAt: '2024-01-15' },
  { id: 2, schoolId: 1, classId: 1, className: 'VII-A', nis: '2024002', nisn: '0012345679', name: 'Budi Santoso', isActive: true, createdAt: '2024-01-15', updatedAt: '2024-01-15' },
  { id: 3, schoolId: 1, classId: 2, className: 'VII-B', nis: '2024003', nisn: '0012345680', name: 'Citra Dewi', isActive: true, createdAt: '2024-01-15', updatedAt: '2024-01-15' },
  { id: 4, schoolId: 1, classId: 2, className: 'VII-B', nis: '2024004', nisn: '0012345681', name: 'Dian Pratama', isActive: true, createdAt: '2024-01-15', updatedAt: '2024-01-15' },
  { id: 5, schoolId: 1, classId: 3, className: 'VIII-A', nis: '2023001', nisn: '0012345682', name: 'Eka Putri', isActive: true, createdAt: '2024-01-15', updatedAt: '2024-01-15' },
  { id: 6, schoolId: 1, classId: 3, className: 'VIII-A', nis: '2023002', nisn: '0012345683', name: 'Fajar Nugroho', isActive: true, createdAt: '2024-01-15', updatedAt: '2024-01-15' },
]

// Table columns
const columns: TableProps['columns'] = [
  {
    title: 'Nama Orang Tua',
    dataIndex: 'name',
    key: 'name',
    sorter: true,
  },
  {
    title: 'Telepon',
    dataIndex: 'phone',
    key: 'phone',
    width: 150,
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
    width: 150,
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
      pageSize: pagination.pageSize,
      search: searchText.value,
    })
    parents.value = response.data
    total.value = response.total
  } catch {
    // Use mock data on error
    parents.value = mockParents
    total.value = mockParents.length
  } finally {
    loading.value = false
  }
}

// Load students for dropdown
const loadStudents = async () => {
  loadingStudents.value = true
  try {
    const response = await schoolService.getStudents({ pageSize: 500 })
    students.value = response.data.filter(s => s.isActive)
  } catch {
    students.value = mockStudents
  } finally {
    loadingStudents.value = false
  }
}

// Filter option for student select
const filterStudentOption = (input: string, option: { label?: string }) => {
  return option?.label?.toLowerCase().includes(input.toLowerCase()) ?? false
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
  formState.password = ''
  formState.studentIds = parent.studentIds || []
  modalVisible.value = true
}

// Reset form
const resetForm = () => {
  formState.name = ''
  formState.phone = ''
  formState.email = ''
  formState.password = ''
  formState.studentIds = []
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
      await formRef.value?.validate(['name', 'phone', 'studentIds'])
    } else {
      await formRef.value?.validate()
    }
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
        studentIds: formState.studentIds,
      }
      await schoolService.updateParent(editingParent.value.id, updateData)
      message.success('Data orang tua berhasil diperbarui')
    } else {
      await schoolService.createParent(formState)
      message.success('Orang tua berhasil ditambahkan')
    }
    modalVisible.value = false
    resetForm()
    loadParents()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } } }
    message.error(err.response?.data?.message || 'Terjadi kesalahan')
  } finally {
    modalLoading.value = false
  }
}

// Handle delete parent
const handleDelete = async (parent: Parent) => {
  try {
    await schoolService.deleteParent(parent.id)
    message.success(`Data ${parent.name} berhasil dihapus`)
    loadParents()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } } }
    message.error(err.response?.data?.message || 'Gagal menghapus data')
  }
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
          <template v-else-if="column.key === 'email'">
            {{ (record as Parent).email || '-' }}
          </template>
          <template v-else-if="column.key === 'studentNames'">
            <Space wrap>
              <Tag
                v-for="(name, index) in (record as Parent).studentNames"
                :key="index"
                color="blue"
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
              <Button size="small" @click="openEditModal(record as Parent)">
                <template #icon><EditOutlined /></template>
                Edit
              </Button>
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
            <FormItem label="Nomor Telepon" name="phone" required>
              <Input v-model:value="formState.phone" placeholder="Contoh: 081234567890" />
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
          label="Password"
          name="password"
          required
          extra="Password untuk login di aplikasi mobile"
        >
          <Input.Password v-model:value="formState.password" placeholder="Minimal 6 karakter" />
        </FormItem>
        <FormItem
          label="Anak (Siswa)"
          name="studentIds"
          required
          extra="Pilih siswa yang merupakan anak dari orang tua ini"
        >
          <Select
            v-model:value="formState.studentIds"
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

@media (max-width: 576px) {
  .toolbar-right {
    margin-top: 16px;
    justify-content: flex-start;
  }
}
</style>
