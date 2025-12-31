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
  Textarea,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  PlusOutlined,
  SearchOutlined,
  DeleteOutlined,
  ReloadOutlined,
  FilterOutlined,
  EyeOutlined,
} from '@ant-design/icons-vue'
import { useRouter } from 'vue-router'
import { bkService, schoolService } from '@/services'
import type { Violation, CreateViolationRequest } from '@/types/bk'
import type { Student } from '@/types/school'
import { VIOLATION_CATEGORIES, VIOLATION_LEVELS } from '@/types/bk'

const { Title } = Typography

const router = useRouter()

// Table state
const loading = ref(false)
const violations = ref<Violation[]>([])
const total = ref(0)
const pagination = reactive({
  current: 1,
  pageSize: 10,
})
const searchText = ref('')
const filterLevel = ref<string | undefined>(undefined)
const filterCategory = ref<string | undefined>(undefined)

// Students for dropdown
const students = ref<Student[]>([])
const loadingStudents = ref(false)

// Modal state
const modalVisible = ref(false)
const modalLoading = ref(false)

// Form state
const formRef = ref()
const formState = reactive<CreateViolationRequest>({
  studentId: 0,
  category: '',
  level: 'ringan',
  description: '',
})

// Form rules
const formRules = {
  studentId: [{ required: true, message: 'Siswa wajib dipilih' }],
  category: [{ required: true, message: 'Kategori wajib dipilih' }],
  level: [{ required: true, message: 'Tingkat wajib dipilih' }],
  description: [{ required: true, message: 'Deskripsi wajib diisi' }],
}

// Mock data for development
const mockViolations: Violation[] = [
  { id: 1, studentId: 1, studentName: 'Ahmad Fauzi', studentNis: '2024001', studentClass: 'VII-A', category: 'Keterlambatan', level: 'ringan', description: 'Terlambat 15 menit', createdBy: 1, createdByName: 'Guru BK', createdAt: new Date().toISOString() },
  { id: 2, studentId: 2, studentName: 'Budi Santoso', studentNis: '2024002', studentClass: 'VII-B', category: 'Seragam', level: 'ringan', description: 'Tidak memakai dasi', createdBy: 1, createdByName: 'Guru BK', createdAt: new Date(Date.now() - 86400000).toISOString() },
  { id: 3, studentId: 3, studentName: 'Citra Dewi', studentNis: '2024003', studentClass: 'VIII-A', category: 'Bolos', level: 'sedang', description: 'Tidak masuk tanpa keterangan', createdBy: 1, createdByName: 'Guru BK', createdAt: new Date(Date.now() - 172800000).toISOString() },
  { id: 4, studentId: 4, studentName: 'Dian Pratama', studentNis: '2024004', studentClass: 'IX-A', category: 'Perilaku', level: 'sedang', description: 'Mengganggu teman saat pelajaran', createdBy: 1, createdByName: 'Guru BK', createdAt: new Date(Date.now() - 259200000).toISOString() },
  { id: 5, studentId: 5, studentName: 'Eka Putri', studentNis: '2024005', studentClass: 'VIII-B', category: 'Merokok', level: 'berat', description: 'Kedapatan merokok di toilet', createdBy: 1, createdByName: 'Guru BK', createdAt: new Date(Date.now() - 345600000).toISOString() },
]

const mockStudents: Student[] = [
  { id: 1, schoolId: 1, classId: 1, className: 'VII-A', nis: '2024001', nisn: '0012345678', name: 'Ahmad Fauzi', isActive: true, createdAt: '', updatedAt: '' },
  { id: 2, schoolId: 1, classId: 2, className: 'VII-B', nis: '2024002', nisn: '0012345679', name: 'Budi Santoso', isActive: true, createdAt: '', updatedAt: '' },
  { id: 3, schoolId: 1, classId: 3, className: 'VIII-A', nis: '2024003', nisn: '0012345680', name: 'Citra Dewi', isActive: true, createdAt: '', updatedAt: '' },
  { id: 4, schoolId: 1, classId: 4, className: 'IX-A', nis: '2024004', nisn: '0012345681', name: 'Dian Pratama', isActive: true, createdAt: '', updatedAt: '' },
  { id: 5, schoolId: 1, classId: 5, className: 'VIII-B', nis: '2024005', nisn: '0012345682', name: 'Eka Putri', isActive: true, createdAt: '', updatedAt: '' },
]

// Table columns
const columns: TableProps['columns'] = [
  {
    title: 'Tanggal',
    dataIndex: 'createdAt',
    key: 'createdAt',
    width: 120,
    sorter: true,
  },
  {
    title: 'Siswa',
    dataIndex: 'studentName',
    key: 'studentName',
  },
  {
    title: 'Kelas',
    dataIndex: 'studentClass',
    key: 'studentClass',
    width: 100,
  },
  {
    title: 'Kategori',
    dataIndex: 'category',
    key: 'category',
    width: 150,
  },
  {
    title: 'Tingkat',
    dataIndex: 'level',
    key: 'level',
    width: 100,
    align: 'center',
  },
  {
    title: 'Deskripsi',
    dataIndex: 'description',
    key: 'description',
    ellipsis: true,
  },
  {
    title: 'Aksi',
    key: 'action',
    width: 120,
    align: 'center',
  },
]

// Computed filtered data
const filteredViolations = computed(() => {
  let result = violations.value

  if (filterLevel.value) {
    result = result.filter(v => v.level === filterLevel.value)
  }

  if (filterCategory.value) {
    result = result.filter(v => v.category === filterCategory.value)
  }

  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(
      (v) =>
        v.studentName?.toLowerCase().includes(search) ||
        v.description.toLowerCase().includes(search) ||
        v.category.toLowerCase().includes(search)
    )
  }

  return result
})

// Format date
const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

// Get level color
const getLevelColor = (level: string) => {
  const levelConfig = VIOLATION_LEVELS.find(l => l.value === level)
  return levelConfig?.color || 'default'
}

// Get level label
const getLevelLabel = (level: string) => {
  const levelConfig = VIOLATION_LEVELS.find(l => l.value === level)
  return levelConfig?.label || level
}

// Load violations data
const loadViolations = async () => {
  loading.value = true
  try {
    const response = await bkService.getViolations({
      page: pagination.current,
      pageSize: pagination.pageSize,
      search: searchText.value,
      level: filterLevel.value,
      category: filterCategory.value,
    })
    violations.value = response.data
    total.value = response.total
  } catch {
    violations.value = mockViolations
    total.value = mockViolations.length
  } finally {
    loading.value = false
  }
}

// Load students for dropdown
const loadStudents = async () => {
  loadingStudents.value = true
  try {
    const response = await schoolService.getStudents({ pageSize: 1000 })
    students.value = response.data
  } catch {
    students.value = mockStudents
  } finally {
    loadingStudents.value = false
  }
}

// Handle table change
const handleTableChange: TableProps['onChange'] = (pag) => {
  pagination.current = pag.current || 1
  pagination.pageSize = pag.pageSize || 10
  loadViolations()
}

// Handle search
const handleSearch = () => {
  pagination.current = 1
  loadViolations()
}

// Handle filter change
const handleFilterChange = () => {
  pagination.current = 1
  loadViolations()
}

// Open create modal
const openCreateModal = () => {
  resetForm()
  modalVisible.value = true
}

// Reset form
const resetForm = () => {
  formState.studentId = 0
  formState.category = ''
  formState.level = 'ringan'
  formState.description = ''
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
    await bkService.createViolation(formState)
    message.success('Pelanggaran berhasil dicatat')
    modalVisible.value = false
    resetForm()
    loadViolations()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } } }
    message.error(err.response?.data?.message || 'Terjadi kesalahan')
  } finally {
    modalLoading.value = false
  }
}

// Handle delete
const handleDelete = async (violation: Violation) => {
  try {
    await bkService.deleteViolation(violation.id)
    message.success('Pelanggaran berhasil dihapus')
    loadViolations()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } } }
    message.error(err.response?.data?.message || 'Gagal menghapus pelanggaran')
  }
}

// View student profile
const viewStudentProfile = (studentId: number) => {
  router.push(`/bk/students/${studentId}`)
}

// Filter student options
const filterStudentOption = (input: string, option: { label: string }) => {
  return option.label.toLowerCase().includes(input.toLowerCase())
}

onMounted(() => {
  loadViolations()
  loadStudents()
})
</script>

<template>
  <div class="violation-management">
    <div class="page-header">
      <Title :level="2" style="margin: 0">Manajemen Pelanggaran</Title>
    </div>

    <Card>
      <!-- Toolbar -->
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col :xs="24" :sm="24" :md="18">
          <Space wrap>
            <Input
              v-model:value="searchText"
              placeholder="Cari siswa atau deskripsi..."
              allow-clear
              style="width: 250px"
              @press-enter="handleSearch"
            >
              <template #prefix>
                <SearchOutlined />
              </template>
            </Input>
            <Select
              v-model:value="filterLevel"
              placeholder="Filter Tingkat"
              allow-clear
              style="width: 150px"
              @change="handleFilterChange"
            >
              <template #suffixIcon>
                <FilterOutlined />
              </template>
              <SelectOption v-for="level in VIOLATION_LEVELS" :key="level.value" :value="level.value">
                {{ level.label }}
              </SelectOption>
            </Select>
            <Select
              v-model:value="filterCategory"
              placeholder="Filter Kategori"
              allow-clear
              style="width: 180px"
              @change="handleFilterChange"
            >
              <SelectOption v-for="cat in VIOLATION_CATEGORIES" :key="cat" :value="cat">
                {{ cat }}
              </SelectOption>
            </Select>
          </Space>
        </Col>
        <Col :xs="24" :sm="24" :md="6" class="toolbar-right">
          <Space>
            <Button @click="loadViolations">
              <template #icon><ReloadOutlined /></template>
              Refresh
            </Button>
            <Button type="primary" @click="openCreateModal">
              <template #icon><PlusOutlined /></template>
              Catat Pelanggaran
            </Button>
          </Space>
        </Col>
      </Row>

      <!-- Table -->
      <Table
        :columns="columns"
        :data-source="filteredViolations"
        :loading="loading"
        :pagination="{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: total,
          showSizeChanger: true,
          showTotal: (total: number) => `Total ${total} pelanggaran`,
        }"
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'createdAt'">
            {{ formatDate((record as Violation).createdAt) }}
          </template>
          <template v-else-if="column.key === 'studentName'">
            <a @click="viewStudentProfile((record as Violation).studentId)">
              {{ (record as Violation).studentName }}
            </a>
          </template>
          <template v-else-if="column.key === 'studentClass'">
            <Tag color="blue">{{ (record as Violation).studentClass }}</Tag>
          </template>
          <template v-else-if="column.key === 'level'">
            <Tag :color="getLevelColor((record as Violation).level)">
              {{ getLevelLabel((record as Violation).level) }}
            </Tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <Space>
              <Button size="small" @click="viewStudentProfile((record as Violation).studentId)">
                <template #icon><EyeOutlined /></template>
              </Button>
              <Popconfirm
                title="Hapus pelanggaran ini?"
                description="Data pelanggaran akan dihapus permanen."
                ok-text="Ya, Hapus"
                cancel-text="Batal"
                @confirm="handleDelete(record as Violation)"
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

    <!-- Create Modal -->
    <Modal
      v-model:open="modalVisible"
      title="Catat Pelanggaran Baru"
      :confirm-loading="modalLoading"
      @ok="handleSubmit"
      @cancel="handleModalCancel"
      width="600px"
    >
      <Form
        ref="formRef"
        :model="formState"
        :rules="formRules"
        layout="vertical"
        style="margin-top: 16px"
      >
        <FormItem label="Siswa" name="studentId" required>
          <Select
            v-model:value="formState.studentId"
            placeholder="Pilih siswa"
            :loading="loadingStudents"
            show-search
            :filter-option="filterStudentOption"
            :options="students.map(s => ({ value: s.id, label: `${s.name} (${s.className})` }))"
          />
        </FormItem>
        <Row :gutter="16">
          <Col :span="12">
            <FormItem label="Kategori" name="category" required>
              <Select v-model:value="formState.category" placeholder="Pilih kategori">
                <SelectOption v-for="cat in VIOLATION_CATEGORIES" :key="cat" :value="cat">
                  {{ cat }}
                </SelectOption>
              </Select>
            </FormItem>
          </Col>
          <Col :span="12">
            <FormItem label="Tingkat" name="level" required>
              <Select v-model:value="formState.level" placeholder="Pilih tingkat">
                <SelectOption v-for="level in VIOLATION_LEVELS" :key="level.value" :value="level.value">
                  <Tag :color="level.color" style="margin-right: 8px">{{ level.label }}</Tag>
                </SelectOption>
              </Select>
            </FormItem>
          </Col>
        </Row>
        <FormItem label="Deskripsi" name="description" required>
          <Textarea
            v-model:value="formState.description"
            placeholder="Jelaskan detail pelanggaran..."
            :rows="4"
          />
        </FormItem>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.violation-management {
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
