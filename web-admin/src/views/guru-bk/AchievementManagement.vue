<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import {
  Table,
  Button,
  Input,
  InputNumber,
  Space,
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
  DatePicker,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  PlusOutlined,
  SearchOutlined,
  DeleteOutlined,
  ReloadOutlined,
  EyeOutlined,
  TrophyOutlined,
  FilePdfOutlined,
} from '@ant-design/icons-vue'
import { useRouter } from 'vue-router'
import dayjs from 'dayjs'
import type { Dayjs } from 'dayjs'
import { bkService, schoolService } from '@/services'
import type { Achievement, CreateAchievementRequest } from '@/types/bk'
import type { Student } from '@/types/school'
import { exportToPDF, formatAchievementForExport } from '@/utils/pdfExport'

const { Title } = Typography
const { RangePicker } = DatePicker

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
const dateRange = ref<[Dayjs, Dayjs] | undefined>(undefined)

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
  let result = achievements.value
  
  if (dateRange.value) {
    const [start, end] = dateRange.value
    result = result.filter(a => {
      const date = dayjs(a.createdAt)
      return date.isAfter(start.startOf('day')) && date.isBefore(end.endOf('day'))
    })
  }

  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(
      (a) =>
        a.studentName?.toLowerCase().includes(search) ||
        a.title.toLowerCase().includes(search) ||
        a.description?.toLowerCase().includes(search)
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

// Load achievements data
const loadAchievements = async () => {
  loading.value = true
  try {
    const response = await bkService.getAchievements({
      page: pagination.current,
      pageSize: pagination.pageSize,
      search: searchText.value,
    })
    achievements.value = response.data || []
    total.value = response.total || 0
  } catch (err) {
    console.error('Failed to load achievements:', err)
    achievements.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// Load students for dropdown
const loadStudents = async () => {
  loadingStudents.value = true
  try {
    const response = await schoolService.getStudents({ page_size: 1000 })
    students.value = response.students || []
  } catch (err) {
    console.error('Failed to load students:', err)
    students.value = []
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

// Handle filter change
const handleFilterChange = () => {
  pagination.current = 1
}

// Export to PDF
const handleExportPDF = () => {
  if (filteredAchievements.value.length === 0) {
    message.warning('Tidak ada data untuk diekspor')
    return
  }

  const dateRangeStr = dateRange.value
    ? { start: dateRange.value[0].format('DD/MM/YYYY'), end: dateRange.value[1].format('DD/MM/YYYY') }
    : undefined

  exportToPDF({
    title: 'Laporan Data Prestasi Siswa',
    filename: `laporan-prestasi-${dayjs().format('YYYY-MM-DD')}`,
    columns: [
      { header: 'Tanggal', dataKey: 'createdAt' },
      { header: 'Siswa', dataKey: 'studentName' },
      { header: 'Kelas', dataKey: 'studentClass' },
      { header: 'Prestasi', dataKey: 'title' },
      { header: 'Poin', dataKey: 'point' },
      { header: 'Deskripsi', dataKey: 'description' },
    ],
    data: filteredAchievements.value.map(a => formatAchievementForExport(a as unknown as Record<string, unknown>)),
    dateRange: dateRangeStr,
  })
  message.success('PDF berhasil diunduh')
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
const filterStudentOption = (input: string, option: any) => {
  return option?.label?.toLowerCase().includes(input.toLowerCase())
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
              size="large"
              style="width: 250px"
              @press-enter="handleSearch"
            >
              <template #prefix>
                <SearchOutlined />
              </template>
            </Input>
            <RangePicker v-model:value="dateRange" format="DD/MM/YYYY" size="large" :placeholder="['Dari Tanggal', 'Sampai Tanggal']" style="width: 250px" @change="handleFilterChange" />
          </Space>
        </Col>
        <Col :xs="24" :sm="24" :md="8" class="toolbar-right">
          <Space>
            <Button size="large" @click="handleExportPDF">
              <template #icon><FilePdfOutlined /></template>
              Export PDF
            </Button>
            <Button size="large" @click="loadAchievements">
              <template #icon><ReloadOutlined /></template>
            </Button>
            <Button type="primary" size="large" @click="openCreateModal">
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
        class="custom-table"
        :scroll="{ x: 800 }"
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
            <span class="class-badge">{{ (record as Achievement).studentClass }}</span>
          </template>
          <template v-else-if="column.key === 'point'">
            <span class="status-badge success">+{{ (record as Achievement).point }}</span>
          </template>
          <template v-else-if="column.key === 'action'">
            <Space>
              <Button type="text" style="color: #3b82f6" @click="viewStudentProfile((record as Achievement).studentId)">
                <template #icon><EyeOutlined /></template>
              </Button>
              <Popconfirm
                title="Hapus prestasi ini?"
                description="Data prestasi akan dihapus permanen."
                ok-text="Ya, Hapus"
                cancel-text="Batal"
                @confirm="handleDelete(record as Achievement)"
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

    <!-- Create Modal -->
    <Modal
      v-model:open="modalVisible"
      title="Catat Prestasi Baru"
      :confirm-loading="modalLoading"
      @ok="handleSubmit"
      @cancel="handleModalCancel"
      width="600px"
      wrap-class-name="modern-modal"
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
            size="large"
            :filter-option="filterStudentOption"
            :options="students.map(s => ({ value: s.id, label: `${s.name} (${s.className})` }))"
          />
        </FormItem>
        <FormItem label="Judul Prestasi" name="title" required>
          <Input v-model:value="formState.title" size="large" placeholder="Contoh: Juara 1 Olimpiade Matematika" />
        </FormItem>
        <FormItem label="Poin" name="point" required>
          <InputNumber
            v-model:value="formState.point"
            :min="1"
            :max="1000"
            size="large"
            style="width: 100%"
            placeholder="Masukkan poin prestasi"
          />
        </FormItem>
        <FormItem label="Deskripsi (Opsional)" name="description">
          <Textarea
            v-model:value="formState.description"
            placeholder="Jelaskan detail prestasi..."
            :rows="3"
            class="custom-textarea"
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

/* Custom Table Styles */
.custom-table :deep(.ant-table-thead > tr > th) {
  background: #fafafa;
  font-weight: 600;
  color: #475569;
}

.custom-table :deep(.ant-table-tbody > tr > td) {
  padding: 16px;
}

.custom-table :deep(.ant-table-tbody > tr:hover > td) {
  background: #f8fafc;
}

/* Badge Styles */
.status-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap;
}

.status-badge.success {
  background-color: #f6ffed;
  color: #52c41a;
  border: 1px solid #b7eb8f;
}

.class-badge {
  background-color: #e6f7ff;
  color: #1890ff;
  border: 1px solid #91d5ff;
  padding: 0 8px;
  border-radius: 4px;
  font-size: 12px;
}

@media (max-width: 768px) {
  .toolbar-right {
    margin-top: 16px;
    justify-content: flex-start;
  }
}
</style>
