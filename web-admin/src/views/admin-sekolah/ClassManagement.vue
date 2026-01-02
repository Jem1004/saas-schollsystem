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
} from '@ant-design/icons-vue'
import { schoolService } from '@/services'
import type { Class, UpdateClassRequest, SchoolUser } from '@/types/school'

const { Title } = Typography

// Table state
const loading = ref(false)
const classes = ref<Class[]>([])
const total = ref(0)
const pagination = reactive({
  current: 1,
  pageSize: 10,
})
const searchText = ref('')

// Teachers for homeroom assignment
const teachers = ref<SchoolUser[]>([])
const loadingTeachers = ref(false)

// Modal state
const modalVisible = ref(false)
const modalLoading = ref(false)
const isEditing = ref(false)
const editingClass = ref<Class | null>(null)

// Form state
const formRef = ref()
const formState = reactive({
  name: '',
  grade: 7,
  year: '',
  homeroom_teacher_id: undefined as number | undefined,
})

// Form rules
const formRules = {
  name: [{ required: true, message: 'Nama kelas wajib diisi' }],
  grade: [{ required: true, message: 'Tingkat kelas wajib diisi' }],
  year: [{ required: true, message: 'Tahun ajaran wajib diisi' }],
}

// Table columns
const columns: TableProps['columns'] = [
  {
    title: 'Nama Kelas',
    dataIndex: 'name',
    key: 'name',
    sorter: true,
  },
  {
    title: 'Tingkat',
    dataIndex: 'grade',
    key: 'grade',
    width: 100,
    align: 'center',
  },
  {
    title: 'Tahun Ajaran',
    dataIndex: 'year',
    key: 'year',
    width: 120,
  },
  {
    title: 'Wali Kelas',
    dataIndex: 'homeroomTeacherName',
    key: 'homeroomTeacherName',
  },
  {
    title: 'Jumlah Siswa',
    dataIndex: 'studentCount',
    key: 'studentCount',
    width: 120,
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
const filteredClasses = computed(() => {
  if (!searchText.value) return classes.value
  const search = searchText.value.toLowerCase()
  return classes.value.filter(
    (cls) =>
      cls.name.toLowerCase().includes(search) ||
      cls.homeroomTeacherName?.toLowerCase().includes(search)
  )
})

// Get current academic year
const getCurrentAcademicYear = (): string => {
  const now = new Date()
  const year = now.getFullYear()
  const month = now.getMonth() + 1
  // Academic year starts in July
  if (month >= 7) {
    return `${year}/${year + 1}`
  }
  return `${year - 1}/${year}`
}

// Load classes data
const loadClasses = async () => {
  loading.value = true
  try {
    const response = await schoolService.getClasses({
      page: pagination.current,
      page_size: pagination.pageSize,
      search: searchText.value,
    })
    classes.value = response.classes
    total.value = response.pagination.total
  } catch (err) {
    console.error('Failed to load classes:', err)
    message.error('Gagal memuat data kelas')
    classes.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// Load teachers for dropdown
const loadTeachers = async () => {
  loadingTeachers.value = true
  try {
    const response = await schoolService.getUsers({ page_size: 100 })
    teachers.value = response.users.filter(u => u.role === 'wali_kelas' || u.role === 'guru')
  } catch (err) {
    console.error('Failed to load teachers:', err)
    teachers.value = []
  } finally {
    loadingTeachers.value = false
  }
}

// Filter option for teacher select
const filterTeacherOption = (input: string, option: { label?: string }) => {
  return option?.label?.toLowerCase().includes(input.toLowerCase()) ?? false
}

// Handle table change (pagination, sorting)
const handleTableChange: TableProps['onChange'] = (pag) => {
  pagination.current = pag.current || 1
  pagination.pageSize = pag.pageSize || 10
  loadClasses()
}

// Handle search
const handleSearch = () => {
  pagination.current = 1
  loadClasses()
}

// Open modal for create
const openCreateModal = () => {
  isEditing.value = false
  editingClass.value = null
  resetForm()
  formState.year = getCurrentAcademicYear()
  modalVisible.value = true
}

// Open modal for edit
const openEditModal = (cls: Class) => {
  isEditing.value = true
  editingClass.value = cls
  formState.name = cls.name
  formState.grade = cls.grade
  formState.year = cls.year
  formState.homeroom_teacher_id = cls.homeroomTeacherId
  modalVisible.value = true
}

// Reset form
const resetForm = () => {
  formState.name = ''
  formState.grade = 7
  formState.year = ''
  formState.homeroom_teacher_id = undefined
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
    if (isEditing.value && editingClass.value) {
      const updateData: UpdateClassRequest = {
        name: formState.name,
        grade: formState.grade,
        year: formState.year,
        homeroom_teacher_id: formState.homeroom_teacher_id,
      }
      await schoolService.updateClass(editingClass.value.id, updateData)
      message.success('Kelas berhasil diperbarui')
    } else {
      await schoolService.createClass({
        name: formState.name,
        grade: formState.grade,
        year: formState.year,
        homeroom_teacher_id: formState.homeroom_teacher_id,
      })
      message.success('Kelas berhasil ditambahkan')
    }
    modalVisible.value = false
    resetForm()
    loadClasses()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Terjadi kesalahan')
  } finally {
    modalLoading.value = false
  }
}

// Handle delete class
const handleDelete = async (cls: Class) => {
  try {
    await schoolService.deleteClass(cls.id)
    message.success(`Kelas ${cls.name} berhasil dihapus`)
    loadClasses()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Gagal menghapus kelas')
  }
}

// Get grade label
const getGradeLabel = (grade: number): string => {
  const gradeMap: Record<number, string> = {
    7: 'VII',
    8: 'VIII',
    9: 'IX',
    10: 'X',
    11: 'XI',
    12: 'XII',
  }
  return gradeMap[grade] || grade.toString()
}

onMounted(() => {
  loadClasses()
  loadTeachers()
})
</script>

<template>
  <div class="class-management">
    <div class="page-header">
      <Title :level="2" style="margin: 0">Manajemen Kelas</Title>
    </div>

    <Card>
      <!-- Toolbar -->
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col :xs="24" :sm="12" :md="8">
          <Input
            v-model:value="searchText"
            placeholder="Cari kelas..."
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
            <Button @click="loadClasses">
              <template #icon><ReloadOutlined /></template>
            </Button>
            <Button type="primary" @click="openCreateModal">
              <template #icon><PlusOutlined /></template>
              Tambah Kelas
            </Button>
          </Space>
        </Col>
      </Row>

      <!-- Table -->
      <Table
        :columns="columns"
        :data-source="filteredClasses"
        :loading="loading"
        :pagination="{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: total,
          showSizeChanger: true,
          showTotal: (total: number) => `Total ${total} kelas`,
        }"
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'grade'">
            <Tag color="blue">{{ getGradeLabel((record as Class).grade) }}</Tag>
          </template>
          <template v-else-if="column.key === 'homeroomTeacherName'">
            <span v-if="(record as Class).homeroomTeacherName">
              <UserOutlined /> {{ (record as Class).homeroomTeacherName }}
            </span>
            <Tag v-else color="default">Belum ditentukan</Tag>
          </template>
          <template v-else-if="column.key === 'studentCount'">
            <Tag>{{ (record as Class).studentCount || 0 }} siswa</Tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <Space>
              <Button size="small" @click="openEditModal(record as Class)">
                <template #icon><EditOutlined /></template>
                Edit
              </Button>
              <Popconfirm
                title="Hapus kelas ini?"
                description="Semua data terkait kelas ini akan dihapus."
                ok-text="Ya, Hapus"
                cancel-text="Batal"
                @confirm="handleDelete(record as Class)"
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
      :title="isEditing ? 'Edit Kelas' : 'Tambah Kelas Baru'"
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
        <FormItem label="Nama Kelas" name="name" required>
          <Input v-model:value="formState.name" placeholder="Contoh: VII-A, VIII-B" />
        </FormItem>
        <Row :gutter="16">
          <Col :span="12">
            <FormItem label="Tingkat" name="grade" required>
              <Select v-model:value="formState.grade" placeholder="Pilih tingkat">
                <SelectOption :value="7">VII (Kelas 7)</SelectOption>
                <SelectOption :value="8">VIII (Kelas 8)</SelectOption>
                <SelectOption :value="9">IX (Kelas 9)</SelectOption>
                <SelectOption :value="10">X (Kelas 10)</SelectOption>
                <SelectOption :value="11">XI (Kelas 11)</SelectOption>
                <SelectOption :value="12">XII (Kelas 12)</SelectOption>
              </Select>
            </FormItem>
          </Col>
          <Col :span="12">
            <FormItem label="Tahun Ajaran" name="year" required>
              <Input v-model:value="formState.year" placeholder="Contoh: 2024/2025" />
            </FormItem>
          </Col>
        </Row>
        <FormItem label="Wali Kelas" name="homeroom_teacher_id">
          <Select
            v-model:value="formState.homeroom_teacher_id"
            placeholder="Pilih wali kelas (opsional)"
            allow-clear
            show-search
            :filter-option="filterTeacherOption"
            :loading="loadingTeachers"
          >
            <SelectOption
              v-for="teacher in teachers"
              :key="teacher.id"
              :value="teacher.id"
              :label="teacher.name || teacher.username"
            >
              {{ teacher.name || teacher.username }}
              <Tag v-if="teacher.role === 'wali_kelas'" size="small" color="blue" style="margin-left: 8px">
                Wali Kelas
              </Tag>
            </SelectOption>
          </Select>
        </FormItem>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.class-management {
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
</style>
