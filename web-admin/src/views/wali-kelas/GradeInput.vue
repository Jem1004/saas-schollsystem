<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import {
  Table,
  Button,
  Space,
  Tag,
  Card,
  Row,
  Col,
  Typography,
  Modal,
  Form,
  FormItem,
  Input,
  InputNumber,
  Select,
  SelectOption,
  Tabs,
  TabPane,
  message,
  Popconfirm,
  Drawer,
  List,
  ListItem,
  ListItemMeta,
  Avatar,
  Empty,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
  UserOutlined,
  HistoryOutlined,
} from '@ant-design/icons-vue'
import { homeroomService } from '@/services'
import type { Grade, ClassStudent, CreateGradeRequest, UpdateGradeRequest, BatchGradeRequest } from '@/types/homeroom'

const { Title, Text, Paragraph } = Typography
const { TextArea } = Input

// State
const loading = ref(false)
const grades = ref<Grade[]>([])
const students = ref<ClassStudent[]>([])
const className = ref('VII-A')
const activeTab = ref('individual')

// Modal state
const modalVisible = ref(false)
const modalLoading = ref(false)
const editingGrade = ref<Grade | null>(null)
const formRef = ref()

// Batch modal state
const batchModalVisible = ref(false)
const batchModalLoading = ref(false)
const batchFormRef = ref()

// History drawer state
const historyDrawerVisible = ref(false)
const selectedStudent = ref<ClassStudent | null>(null)
const studentGrades = ref<Grade[]>([])
const loadingHistory = ref(false)

// Form state
const formState = ref<{
  studentId: number | undefined
  title: string
  score: number | undefined
  description: string
}>({
  studentId: undefined,
  title: '',
  score: undefined,
  description: '',
})

// Batch form state
const batchFormState = ref<{
  title: string
  description: string
  grades: { studentId: number; studentName: string; score: number | undefined }[]
}>({
  title: '',
  description: '',
  grades: [],
})

// Mock data for development
const mockStudents: ClassStudent[] = [
  { id: 1, nis: '2024001', nisn: '0012345678', name: 'Ahmad Fauzi', isActive: true },
  { id: 2, nis: '2024002', nisn: '0012345679', name: 'Budi Santoso', isActive: true },
  { id: 3, nis: '2024003', nisn: '0012345680', name: 'Citra Dewi', isActive: true },
  { id: 4, nis: '2024004', nisn: '0012345681', name: 'Dian Pratama', isActive: true },
  { id: 5, nis: '2024005', nisn: '0012345682', name: 'Eka Putri', isActive: true },
]

const mockGrades: Grade[] = [
  { id: 1, studentId: 1, studentName: 'Ahmad Fauzi', studentNis: '2024001', title: 'Ulangan Matematika', score: 85, description: 'Ulangan Harian Bab 1', createdBy: 1, createdAt: new Date().toISOString(), updatedAt: new Date().toISOString() },
  { id: 2, studentId: 2, studentName: 'Budi Santoso', studentNis: '2024002', title: 'Ulangan Matematika', score: 78, description: 'Ulangan Harian Bab 1', createdBy: 1, createdAt: new Date().toISOString(), updatedAt: new Date().toISOString() },
  { id: 3, studentId: 3, studentName: 'Citra Dewi', studentNis: '2024003', title: 'Ulangan Matematika', score: 92, description: 'Ulangan Harian Bab 1', createdBy: 1, createdAt: new Date().toISOString(), updatedAt: new Date().toISOString() },
  { id: 4, studentId: 1, studentName: 'Ahmad Fauzi', studentNis: '2024001', title: 'Tugas IPA', score: 90, description: 'Tugas Praktikum', createdBy: 1, createdAt: new Date(Date.now() - 86400000).toISOString(), updatedAt: new Date(Date.now() - 86400000).toISOString() },
  { id: 5, studentId: 4, studentName: 'Dian Pratama', studentNis: '2024004', title: 'Ulangan Matematika', score: 75, description: 'Ulangan Harian Bab 1', createdBy: 1, createdAt: new Date().toISOString(), updatedAt: new Date().toISOString() },
]

// Table columns
const columns: TableProps['columns'] = [
  {
    title: 'NIS',
    dataIndex: 'studentNis',
    key: 'studentNis',
    width: 100,
  },
  {
    title: 'Nama Siswa',
    dataIndex: 'studentName',
    key: 'studentName',
  },
  {
    title: 'Judul',
    dataIndex: 'title',
    key: 'title',
  },
  {
    title: 'Nilai',
    dataIndex: 'score',
    key: 'score',
    width: 100,
    align: 'center',
  },
  {
    title: 'Tanggal',
    dataIndex: 'createdAt',
    key: 'createdAt',
    width: 150,
  },
  {
    title: 'Aksi',
    key: 'action',
    width: 150,
    align: 'center',
  },
]

// Get score color
const getScoreColor = (score: number): string => {
  if (score >= 85) return '#22c55e'
  if (score >= 70) return '#f97316'
  return '#ef4444'
}

// Get score tag color
const getScoreTagColor = (score: number): string => {
  if (score >= 85) return 'success'
  if (score >= 70) return 'warning'
  return 'error'
}

// Format date
const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

// Load grades
const loadGrades = async () => {
  loading.value = true
  try {
    const response = await homeroomService.getGrades({ pageSize: 100 })
    grades.value = response.data
  } catch {
    grades.value = mockGrades
  } finally {
    loading.value = false
  }
}

// Load students
const loadStudents = async () => {
  try {
    const response = await homeroomService.getClassStudents({ pageSize: 100 })
    students.value = response.data
  } catch {
    students.value = mockStudents
  }
}

// Open modal for individual grade
const openGradeModal = (grade?: Grade) => {
  editingGrade.value = grade || null
  
  if (grade) {
    formState.value = {
      studentId: grade.studentId,
      title: grade.title,
      score: grade.score,
      description: grade.description || '',
    }
  } else {
    formState.value = {
      studentId: undefined,
      title: '',
      score: undefined,
      description: '',
    }
  }
  
  modalVisible.value = true
}

// Close modal
const closeModal = () => {
  modalVisible.value = false
  editingGrade.value = null
  formRef.value?.resetFields()
}

// Filter option for student select
const filterStudentOption = (input: string, option: { label?: string }) => {
  return option.label?.toLowerCase().includes(input.toLowerCase()) ?? false
}

// Submit individual grade
const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    modalLoading.value = true

    if (editingGrade.value) {
      const data: UpdateGradeRequest = {
        title: formState.value.title,
        score: formState.value.score,
        description: formState.value.description || undefined,
      }
      await homeroomService.updateGrade(editingGrade.value.id, data)
      message.success('Nilai berhasil diperbarui')
    } else {
      const data: CreateGradeRequest = {
        studentId: formState.value.studentId!,
        title: formState.value.title,
        score: formState.value.score!,
        description: formState.value.description || undefined,
      }
      await homeroomService.createGrade(data)
      message.success('Nilai berhasil ditambahkan')
    }

    closeModal()
    loadGrades()
  } catch (err: unknown) {
    if (err && typeof err === 'object' && 'errorFields' in err) {
      return
    }
    message.error('Gagal menyimpan nilai')
  } finally {
    modalLoading.value = false
  }
}

// Open batch modal
const openBatchModal = () => {
  batchFormState.value = {
    title: '',
    description: '',
    grades: students.value.map(s => ({
      studentId: s.id,
      studentName: s.name,
      score: undefined,
    })),
  }
  batchModalVisible.value = true
}

// Close batch modal
const closeBatchModal = () => {
  batchModalVisible.value = false
  batchFormRef.value?.resetFields()
}

// Submit batch grades
const handleBatchSubmit = async () => {
  try {
    await batchFormRef.value?.validate()
    
    const validGrades = batchFormState.value.grades.filter(g => g.score !== undefined && g.score !== null)
    
    if (validGrades.length === 0) {
      message.warning('Masukkan minimal satu nilai')
      return
    }

    batchModalLoading.value = true

    const data: BatchGradeRequest = {
      title: batchFormState.value.title,
      description: batchFormState.value.description || undefined,
      grades: validGrades.map(g => ({
        studentId: g.studentId,
        score: g.score!,
      })),
    }

    await homeroomService.createBatchGrades(data)
    message.success(`${validGrades.length} nilai berhasil ditambahkan`)

    closeBatchModal()
    loadGrades()
  } catch (err: unknown) {
    if (err && typeof err === 'object' && 'errorFields' in err) {
      return
    }
    message.error('Gagal menyimpan nilai')
  } finally {
    batchModalLoading.value = false
  }
}

// Delete grade
const handleDelete = async (id: number) => {
  try {
    await homeroomService.deleteGrade(id)
    message.success('Nilai berhasil dihapus')
    loadGrades()
  } catch {
    message.error('Gagal menghapus nilai')
  }
}

// Open history drawer
const openHistoryDrawer = async (student: ClassStudent) => {
  selectedStudent.value = student
  historyDrawerVisible.value = true
  loadingHistory.value = true
  
  try {
    const response = await homeroomService.getStudentGrades(student.id, { pageSize: 50 })
    studentGrades.value = response.data
  } catch {
    studentGrades.value = mockGrades.filter(g => g.studentId === student.id)
  } finally {
    loadingHistory.value = false
  }
}

// Close history drawer
const closeHistoryDrawer = () => {
  historyDrawerVisible.value = false
  selectedStudent.value = null
  studentGrades.value = []
}

// Calculate average score for student
const calculateAverage = computed(() => {
  if (studentGrades.value.length === 0) return 0
  const sum = studentGrades.value.reduce((acc, g) => acc + g.score, 0)
  return Math.round(sum / studentGrades.value.length)
})

onMounted(() => {
  loadGrades()
  loadStudents()
})
</script>

<template>
  <div class="grade-input">
    <div class="page-header">
      <div>
        <Title :level="2" style="margin: 0">Input Nilai</Title>
        <Text type="secondary">Kelas {{ className }}</Text>
      </div>
    </div>

    <Card>
      <Tabs v-model:activeKey="activeTab">
        <TabPane key="individual" tab="Input Individual">
          <!-- Toolbar -->
          <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
            <Col :xs="24" :sm="12">
              <Text type="secondary">Total {{ grades.length }} nilai tercatat</Text>
            </Col>
            <Col :xs="24" :sm="12" class="toolbar-right">
              <Space>
                <Button @click="loadGrades">
                  <template #icon><ReloadOutlined /></template>
                  Refresh
                </Button>
                <Button type="primary" @click="openGradeModal()">
                  <template #icon><PlusOutlined /></template>
                  Tambah Nilai
                </Button>
              </Space>
            </Col>
          </Row>

          <!-- Table -->
          <Table
            :columns="columns"
            :data-source="grades"
            :loading="loading"
            :pagination="{ pageSize: 10 }"
            row-key="id"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'score'">
                <Tag :color="getScoreTagColor((record as Grade).score)">
                  {{ (record as Grade).score }}
                </Tag>
              </template>
              <template v-else-if="column.key === 'createdAt'">
                {{ formatDate((record as Grade).createdAt) }}
              </template>
              <template v-else-if="column.key === 'action'">
                <Space>
                  <Button type="link" size="small" @click="openGradeModal(record as Grade)">
                    <template #icon><EditOutlined /></template>
                  </Button>
                  <Popconfirm
                    title="Hapus nilai ini?"
                    ok-text="Ya"
                    cancel-text="Tidak"
                    @confirm="handleDelete((record as Grade).id)"
                  >
                    <Button type="link" size="small" danger>
                      <template #icon><DeleteOutlined /></template>
                    </Button>
                  </Popconfirm>
                </Space>
              </template>
            </template>
          </Table>
        </TabPane>

        <TabPane key="batch" tab="Input Batch">
          <!-- Batch Input -->
          <Row :gutter="16" class="toolbar" justify="end">
            <Col>
              <Button type="primary" @click="openBatchModal">
                <template #icon><PlusOutlined /></template>
                Input Nilai Batch
              </Button>
            </Col>
          </Row>

          <div class="batch-info">
            <Paragraph type="secondary">
              Input nilai batch memungkinkan Anda memasukkan nilai untuk semua siswa sekaligus dengan judul yang sama.
              Cocok untuk ulangan harian, tugas, atau penilaian lainnya.
            </Paragraph>
          </div>
        </TabPane>

        <TabPane key="history" tab="Riwayat per Siswa">
          <!-- Student List for History -->
          <List
            :data-source="students"
            :loading="loading"
          >
            <template #renderItem="{ item }">
              <ListItem class="student-item" @click="openHistoryDrawer(item as ClassStudent)">
                <ListItemMeta>
                  <template #avatar>
                    <Avatar :style="{ backgroundColor: '#f97316' }">
                      <template #icon><UserOutlined /></template>
                    </Avatar>
                  </template>
                  <template #title>
                    {{ (item as ClassStudent).name }}
                  </template>
                  <template #description>
                    NIS: {{ (item as ClassStudent).nis }}
                  </template>
                </ListItemMeta>
                <template #actions>
                  <Button type="link">
                    <HistoryOutlined /> Lihat Riwayat
                  </Button>
                </template>
              </ListItem>
            </template>
          </List>
        </TabPane>
      </Tabs>
    </Card>

    <!-- Individual Grade Modal -->
    <Modal
      v-model:open="modalVisible"
      :title="editingGrade ? 'Edit Nilai' : 'Tambah Nilai'"
      :confirm-loading="modalLoading"
      @ok="handleSubmit"
      @cancel="closeModal"
    >
      <Form
        ref="formRef"
        :model="formState"
        layout="vertical"
        style="margin-top: 16px"
      >
        <FormItem
          label="Siswa"
          name="studentId"
          :rules="[{ required: true, message: 'Pilih siswa' }]"
        >
          <Select
            v-model:value="formState.studentId"
            placeholder="Pilih siswa"
            :disabled="!!editingGrade"
            show-search
            :filter-option="filterStudentOption"
          >
            <SelectOption
              v-for="student in students"
              :key="student.id"
              :value="student.id"
              :label="student.name"
            >
              {{ student.nis }} - {{ student.name }}
            </SelectOption>
          </Select>
        </FormItem>

        <FormItem
          label="Judul"
          name="title"
          :rules="[{ required: true, message: 'Masukkan judul nilai' }]"
        >
          <Input v-model:value="formState.title" placeholder="Contoh: Ulangan Matematika Bab 1" />
        </FormItem>

        <FormItem
          label="Nilai"
          name="score"
          :rules="[{ required: true, message: 'Masukkan nilai' }]"
        >
          <InputNumber
            v-model:value="formState.score"
            :min="0"
            :max="100"
            placeholder="0-100"
            style="width: 100%"
          />
        </FormItem>

        <FormItem label="Keterangan" name="description">
          <TextArea
            v-model:value="formState.description"
            placeholder="Keterangan tambahan (opsional)"
            :rows="3"
          />
        </FormItem>
      </Form>
    </Modal>

    <!-- Batch Grade Modal -->
    <Modal
      v-model:open="batchModalVisible"
      title="Input Nilai Batch"
      :confirm-loading="batchModalLoading"
      width="700px"
      @ok="handleBatchSubmit"
      @cancel="closeBatchModal"
    >
      <Form
        ref="batchFormRef"
        :model="batchFormState"
        layout="vertical"
        style="margin-top: 16px"
      >
        <FormItem
          label="Judul"
          name="title"
          :rules="[{ required: true, message: 'Masukkan judul nilai' }]"
        >
          <Input v-model:value="batchFormState.title" placeholder="Contoh: Ulangan Matematika Bab 1" />
        </FormItem>

        <FormItem label="Keterangan" name="description">
          <TextArea
            v-model:value="batchFormState.description"
            placeholder="Keterangan tambahan (opsional)"
            :rows="2"
          />
        </FormItem>

        <div class="batch-grades-header">
          <Text strong>Nilai Siswa</Text>
          <Text type="secondary">(Kosongkan jika tidak ingin memasukkan nilai)</Text>
        </div>

        <div class="batch-grades-list">
          <Row
            v-for="(grade, index) in batchFormState.grades"
            :key="grade.studentId"
            :gutter="16"
            class="batch-grade-row"
          >
            <Col :span="16">
              <Text>{{ grade.studentName }}</Text>
            </Col>
            <Col :span="8">
              <InputNumber
                v-model:value="batchFormState.grades[index].score"
                :min="0"
                :max="100"
                placeholder="Nilai"
                style="width: 100%"
              />
            </Col>
          </Row>
        </div>
      </Form>
    </Modal>

    <!-- History Drawer -->
    <Drawer
      v-model:open="historyDrawerVisible"
      :title="`Riwayat Nilai - ${selectedStudent?.name || ''}`"
      width="500"
      @close="closeHistoryDrawer"
    >
      <div v-if="selectedStudent" class="history-content">
        <Card class="average-card" size="small">
          <Row :gutter="16" align="middle">
            <Col :span="12">
              <Text type="secondary">Rata-rata Nilai</Text>
            </Col>
            <Col :span="12" style="text-align: right">
              <Text :style="{ fontSize: '24px', fontWeight: 'bold', color: getScoreColor(calculateAverage) }">
                {{ calculateAverage }}
              </Text>
            </Col>
          </Row>
        </Card>

        <List
          v-if="studentGrades.length > 0"
          :data-source="studentGrades"
          :loading="loadingHistory"
          style="margin-top: 16px"
        >
          <template #renderItem="{ item }">
            <ListItem>
              <ListItemMeta>
                <template #avatar>
                  <Avatar :style="{ backgroundColor: getScoreColor((item as Grade).score) }">
                    {{ (item as Grade).score }}
                  </Avatar>
                </template>
                <template #title>
                  {{ (item as Grade).title }}
                </template>
                <template #description>
                  <div v-if="(item as Grade).description">{{ (item as Grade).description }}</div>
                  <Text type="secondary" style="font-size: 12px">
                    {{ formatDate((item as Grade).createdAt) }}
                  </Text>
                </template>
              </ListItemMeta>
            </ListItem>
          </template>
        </List>
        <Empty v-else description="Belum ada nilai tercatat" />
      </div>
    </Drawer>
  </div>
</template>

<style scoped>
.grade-input {
  padding: 0;
}

.page-header {
  margin-bottom: 24px;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  flex-wrap: wrap;
  gap: 8px;
}

.toolbar {
  margin-bottom: 16px;
}

.toolbar-right {
  display: flex;
  justify-content: flex-end;
}

.batch-info {
  padding: 24px;
  background: #fafafa;
  border-radius: 8px;
  text-align: center;
}

.student-item {
  cursor: pointer;
  transition: background-color 0.2s;
}

.student-item:hover {
  background-color: #fafafa;
}

.batch-grades-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f0f0f0;
}

.batch-grades-list {
  max-height: 400px;
  overflow-y: auto;
}

.batch-grade-row {
  padding: 8px 0;
  align-items: center;
}

.batch-grade-row:not(:last-child) {
  border-bottom: 1px solid #f0f0f0;
}

.history-content {
  padding: 0;
}

.average-card {
  background: #f6ffed;
  border-color: #b7eb8f;
}

@media (max-width: 768px) {
  .toolbar-right {
    margin-top: 16px;
    justify-content: flex-start;
  }
}
</style>
