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
  FilterOutlined,
} from '@ant-design/icons-vue'
import { schoolService } from '@/services'
import type { Student, Class, CreateStudentRequest, UpdateStudentRequest } from '@/types/school'

const { Title } = Typography

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

// Form state
const formRef = ref()
const formState = reactive<CreateStudentRequest & { isActive?: boolean }>({
  classId: 0,
  nis: '',
  nisn: '',
  name: '',
  rfidCode: '',
  isActive: true,
})

// Form rules
const formRules = {
  classId: [{ required: true, message: 'Kelas wajib dipilih' }],
  nis: [{ required: true, message: 'NIS wajib diisi' }],
  nisn: [{ required: true, message: 'NISN wajib diisi' }],
  name: [{ required: true, message: 'Nama siswa wajib diisi' }],
}

// Mock data for development
const mockStudents: Student[] = [
  { id: 1, schoolId: 1, classId: 1, className: 'VII-A', nis: '2024001', nisn: '0012345678', name: 'Ahmad Fauzi', rfidCode: 'RF001', isActive: true, createdAt: '2024-01-15T10:00:00Z', updatedAt: '2024-01-15T10:00:00Z' },
  { id: 2, schoolId: 1, classId: 1, className: 'VII-A', nis: '2024002', nisn: '0012345679', name: 'Budi Santoso', rfidCode: 'RF002', isActive: true, createdAt: '2024-01-15T10:00:00Z', updatedAt: '2024-01-15T10:00:00Z' },
  { id: 3, schoolId: 1, classId: 2, className: 'VII-B', nis: '2024003', nisn: '0012345680', name: 'Citra Dewi', rfidCode: 'RF003', isActive: true, createdAt: '2024-01-15T10:00:00Z', updatedAt: '2024-01-15T10:00:00Z' },
  { id: 4, schoolId: 1, classId: 2, className: 'VII-B', nis: '2024004', nisn: '0012345681', name: 'Dian Pratama', rfidCode: 'RF004', isActive: true, createdAt: '2024-01-15T10:00:00Z', updatedAt: '2024-01-15T10:00:00Z' },
  { id: 5, schoolId: 1, classId: 3, className: 'VIII-A', nis: '2023001', nisn: '0012345682', name: 'Eka Putri', rfidCode: 'RF005', isActive: true, createdAt: '2024-01-15T10:00:00Z', updatedAt: '2024-01-15T10:00:00Z' },
  { id: 6, schoolId: 1, classId: 3, className: 'VIII-A', nis: '2023002', nisn: '0012345683', name: 'Fajar Nugroho', isActive: false, createdAt: '2024-01-15T10:00:00Z', updatedAt: '2024-01-15T10:00:00Z' },
]

const mockClasses: Class[] = [
  { id: 1, schoolId: 1, name: 'VII-A', grade: 7, year: '2024/2025', studentCount: 30, createdAt: '2024-01-15', updatedAt: '2024-01-15' },
  { id: 2, schoolId: 1, name: 'VII-B', grade: 7, year: '2024/2025', studentCount: 30, createdAt: '2024-01-15', updatedAt: '2024-01-15' },
  { id: 3, schoolId: 1, name: 'VIII-A', grade: 8, year: '2024/2025', studentCount: 32, createdAt: '2024-01-15', updatedAt: '2024-01-15' },
  { id: 4, schoolId: 1, name: 'VIII-B', grade: 8, year: '2024/2025', studentCount: 28, createdAt: '2024-01-15', updatedAt: '2024-01-15' },
  { id: 5, schoolId: 1, name: 'IX-A', grade: 9, year: '2024/2025', studentCount: 30, createdAt: '2024-01-15', updatedAt: '2024-01-15' },
]

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
    title: 'Aksi',
    key: 'action',
    width: 150,
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
      pageSize: pagination.pageSize,
      search: searchText.value,
      classId: filterClassId.value,
    })
    students.value = response.data
    total.value = response.total
  } catch {
    // Use mock data on error
    students.value = mockStudents
    total.value = mockStudents.length
  } finally {
    loading.value = false
  }
}

// Load classes for dropdown
const loadClasses = async () => {
  loadingClasses.value = true
  try {
    const response = await schoolService.getClasses({ pageSize: 100 })
    classes.value = response.data
  } catch {
    classes.value = mockClasses
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
  formState.classId = student.classId
  formState.nis = student.nis
  formState.nisn = student.nisn
  formState.name = student.name
  formState.rfidCode = student.rfidCode || ''
  formState.isActive = student.isActive
  modalVisible.value = true
}

// Reset form
const resetForm = () => {
  formState.classId = 0
  formState.nis = ''
  formState.nisn = ''
  formState.name = ''
  formState.rfidCode = ''
  formState.isActive = true
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
    if (isEditing.value && editingStudent.value) {
      const updateData: UpdateStudentRequest = {
        classId: formState.classId,
        nis: formState.nis,
        name: formState.name,
        rfidCode: formState.rfidCode || undefined,
        isActive: formState.isActive,
      }
      await schoolService.updateStudent(editingStudent.value.id, updateData)
      message.success('Siswa berhasil diperbarui')
    } else {
      await schoolService.createStudent({
        classId: formState.classId,
        nis: formState.nis,
        nisn: formState.nisn,
        name: formState.name,
        rfidCode: formState.rfidCode || undefined,
      })
      message.success('Siswa berhasil ditambahkan')
    }
    modalVisible.value = false
    resetForm()
    loadStudents()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } } }
    message.error(err.response?.data?.message || 'Terjadi kesalahan')
  } finally {
    modalLoading.value = false
  }
}

// Handle delete student
const handleDelete = async (student: Student) => {
  try {
    await schoolService.deleteStudent(student.id)
    message.success(`Siswa ${student.name} berhasil dihapus`)
    loadStudents()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } } }
    message.error(err.response?.data?.message || 'Gagal menghapus siswa')
  }
}

onMounted(() => {
  loadStudents()
  loadClasses()
})
</script>

<template>
  <div class="student-management">
    <div class="page-header">
      <Title :level="2" style="margin: 0">Manajemen Siswa</Title>
    </div>

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
            <Tag v-if="(record as Student).rfidCode" color="green">
              {{ (record as Student).rfidCode }}
            </Tag>
            <Tag v-else color="default">Belum ada</Tag>
          </template>
          <template v-else-if="column.key === 'isActive'">
            <Tag :color="(record as Student).isActive ? 'success' : 'default'">
              {{ (record as Student).isActive ? 'Aktif' : 'Nonaktif' }}
            </Tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <Space>
              <Button size="small" @click="openEditModal(record as Student)">
                <template #icon><EditOutlined /></template>
                Edit
              </Button>
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
        <FormItem label="Kelas" name="classId" required>
          <Select
            v-model:value="formState.classId"
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
        <FormItem label="Kode RFID" name="rfidCode">
          <Input v-model:value="formState.rfidCode" placeholder="Kode kartu RFID (opsional)" />
        </FormItem>
        <FormItem v-if="isEditing" label="Status" name="isActive">
          <Switch v-model:checked="formState.isActive" checked-children="Aktif" un-checked-children="Nonaktif" />
        </FormItem>
      </Form>
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

@media (max-width: 768px) {
  .toolbar-right {
    margin-top: 16px;
    justify-content: flex-start;
  }
}
</style>
