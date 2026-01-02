<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import {
  Table,
  Button,
  Input,
  InputNumber,
  Space,
  Tag,
  Modal,
  Form,
  FormItem,
  Select,
  message,
  Popconfirm,
  Card,
  Row,
  Col,
  Typography,
  Textarea,
  Statistic,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  PlusOutlined,
  SearchOutlined,
  DeleteOutlined,
  ReloadOutlined,
  EyeOutlined,
  TrophyOutlined,
} from '@ant-design/icons-vue'
import { useRouter } from 'vue-router'
import { bkService, schoolService } from '@/services'
import type { Achievement, CreateAchievementRequest } from '@/types/bk'
import type { Student } from '@/types/school'

const { Title } = Typography

const router = useRouter()

// Table state
const loading = ref(false)
const achievements = ref<Achievement[]>([])
const total = ref(0)
const pagination = reactive({
  current: 1,
  pageSize: 10,
})
const searchText = ref('')

// Students for dropdown
const students = ref<Student[]>([])
const loadingStudents = ref(false)

// Modal state
const modalVisible = ref(false)
const modalLoading = ref(false)

// Form state
const formRef = ref()
const formState = reactive<CreateAchievementRequest>({
  studentId: 0,
  title: '',
  point: 10,
  description: '',
})

// Form rules
const formRules = {
  studentId: [{ required: true, message: 'Siswa wajib dipilih' }],
  title: [{ required: true, message: 'Judul prestasi wajib diisi' }],
  point: [{ required: true, message: 'Poin wajib diisi' }],
}

// Total points stat
const totalPoints = computed(() => {
  return achievements.value.reduce((sum, a) => sum + a.point, 0)
})

// Mock data for development
const mockAchievements: Achievement[] = [
  { id: 1, studentId: 1, studentName: 'Ahmad Fauzi', studentNis: '2024001', studentClass: 'VII-A', title: 'Juara 1 Olimpiade Matematika Tingkat Kota', point: 100, description: 'Meraih juara 1 dalam olimpiade matematika tingkat kota', createdBy: 1, createdByName: 'Guru BK', createdAt: new Date().toISOString() },
  { id: 2, studentId: 2, studentName: 'Budi Santoso', studentNis: '2024002', studentClass: 'VII-B', title: 'Juara 2 Lomba Pidato', point: 75, description: 'Meraih juara 2 dalam lomba pidato tingkat kabupaten', createdBy: 1, createdByName: 'Guru BK', createdAt: new Date(Date.now() - 86400000).toISOString() },
  { id: 3, studentId: 3, studentName: 'Citra Dewi', studentNis: '2024003', studentClass: 'VIII-A', title: 'Siswa Teladan Bulan Ini', point: 50, description: 'Terpilih sebagai siswa teladan bulan ini', createdBy: 1, createdByName: 'Guru BK', createdAt: new Date(Date.now() - 172800000).toISOString() },
  { id: 4, studentId: 4, studentName: 'Dian Pratama', studentNis: '2024004', studentClass: 'IX-A', title: 'Juara 3 Lomba Cerdas Cermat', point: 50, description: 'Meraih juara 3 dalam lomba cerdas cermat', createdBy: 1, createdByName: 'Guru BK', createdAt: new Date(Date.now() - 259200000).toISOString() },
  { id: 5, studentId: 5, studentName: 'Eka Putri', studentNis: '2024005', studentClass: 'VIII-B', title: 'Penghargaan Kehadiran Sempurna', point: 25, description: 'Tidak pernah absen selama 1 semester', createdBy: 1, createdByName: 'Guru BK', createdAt: new Date(Date.now() - 345600000).toISOString() },
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
    title: 'Prestasi',
    dataIndex: 'title',
    key: 'title',
  },
  {
    title: 'Poin',
    dataIndex: 'point',
    key: 'point',
    width: 100,
    align: 'center',
  },
  {
    title: 'Aksi',
    key: 'action',
    width: 120,
    align: 'center',
  },
]

// Computed filtered data
const filteredAchievements = computed(() => {
  if (!searchText.value) return achievements.value

  const search = searchText.value.toLowerCase()
  return achievements.value.filter(
    (a) =>
      a.studentName?.toLowerCase().includes(search) ||
      a.title.toLowerCase().includes(search) ||
      a.description?.toLowerCase().includes(search)
  )
})

// Format date
const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

// Load achievements data
const loadAchievements = async () => {
  loading.value = true
  try {
    const response = await bkService.getAchievements({
      page: pagination.current,
      pageSize: pagination.pageSize,
      search: searchText.value,
    })
    achievements.value = response.data
    total.value = response.total
  } catch {
    achievements.value = mockAchievements
    total.value = mockAchievements.length
  } finally {
    loading.value = false
  }
}

// Load students for dropdown
const loadStudents = async () => {
  loadingStudents.value = true
  try {
    const response = await schoolService.getStudents({ page_size: 1000 })
    students.value = response.students
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
  loadAchievements()
}

// Handle search
const handleSearch = () => {
  pagination.current = 1
  loadAchievements()
}

// Open create modal
const openCreateModal = () => {
  resetForm()
  modalVisible.value = true
}

// Reset form
const resetForm = () => {
  formState.studentId = 0
  formState.title = ''
  formState.point = 10
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
    await bkService.createAchievement(formState)
    message.success('Prestasi berhasil dicatat')
    modalVisible.value = false
    resetForm()
    loadAchievements()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Terjadi kesalahan')
  } finally {
    modalLoading.value = false
  }
}

// Handle delete
const handleDelete = async (achievement: Achievement) => {
  try {
    await bkService.deleteAchievement(achievement.id)
    message.success('Prestasi berhasil dihapus')
    loadAchievements()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Gagal menghapus prestasi')
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
  loadAchievements()
  loadStudents()
})
</script>

<template>
  <div class="achievement-management">
    <div class="page-header">
      <Title :level="2" style="margin: 0">Manajemen Prestasi</Title>
    </div>

    <!-- Stats Card -->
    <Row :gutter="24" style="margin-bottom: 24px">
      <Col :xs="24" :sm="12" :md="8">
        <Card>
          <Statistic
            title="Total Prestasi Tercatat"
            :value="total"
            :value-style="{ color: '#22c55e' }"
          >
            <template #prefix>
              <TrophyOutlined />
            </template>
          </Statistic>
        </Card>
      </Col>
      <Col :xs="24" :sm="12" :md="8">
        <Card>
          <Statistic
            title="Total Poin Diberikan"
            :value="totalPoints"
            :value-style="{ color: '#f97316' }"
            suffix="poin"
          />
        </Card>
      </Col>
    </Row>

    <Card>
      <!-- Toolbar -->
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col :xs="24" :sm="24" :md="16">
          <Space wrap>
            <Input
              v-model:value="searchText"
              placeholder="Cari siswa atau prestasi..."
              allow-clear
              style="width: 300px"
              @press-enter="handleSearch"
            >
              <template #prefix>
                <SearchOutlined />
              </template>
            </Input>
          </Space>
        </Col>
        <Col :xs="24" :sm="24" :md="8" class="toolbar-right">
          <Space>
            <Button @click="loadAchievements">
              <template #icon><ReloadOutlined /></template>
            </Button>
            <Button type="primary" @click="openCreateModal">
              <template #icon><PlusOutlined /></template>
              Catat Prestasi
            </Button>
          </Space>
        </Col>
      </Row>

      <!-- Table -->
      <Table
        :columns="columns"
        :data-source="filteredAchievements"
        :loading="loading"
        :pagination="{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: total,
          showSizeChanger: true,
          showTotal: (total: number) => `Total ${total} prestasi`,
        }"
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'createdAt'">
            {{ formatDate((record as Achievement).createdAt) }}
          </template>
          <template v-else-if="column.key === 'studentName'">
            <a @click="viewStudentProfile((record as Achievement).studentId)">
              {{ (record as Achievement).studentName }}
            </a>
          </template>
          <template v-else-if="column.key === 'studentClass'">
            <Tag color="blue">{{ (record as Achievement).studentClass }}</Tag>
          </template>
          <template v-else-if="column.key === 'point'">
            <Tag color="success">+{{ (record as Achievement).point }}</Tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <Space>
              <Button size="small" @click="viewStudentProfile((record as Achievement).studentId)">
                <template #icon><EyeOutlined /></template>
              </Button>
              <Popconfirm
                title="Hapus prestasi ini?"
                description="Data prestasi akan dihapus permanen."
                ok-text="Ya, Hapus"
                cancel-text="Batal"
                @confirm="handleDelete(record as Achievement)"
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
      title="Catat Prestasi Baru"
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
        <FormItem label="Judul Prestasi" name="title" required>
          <Input v-model:value="formState.title" placeholder="Contoh: Juara 1 Olimpiade Matematika" />
        </FormItem>
        <FormItem label="Poin" name="point" required>
          <InputNumber
            v-model:value="formState.point"
            :min="1"
            :max="1000"
            style="width: 100%"
            placeholder="Masukkan poin prestasi"
          />
        </FormItem>
        <FormItem label="Deskripsi (Opsional)" name="description">
          <Textarea
            v-model:value="formState.description"
            placeholder="Jelaskan detail prestasi..."
            :rows="3"
          />
        </FormItem>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.achievement-management {
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
