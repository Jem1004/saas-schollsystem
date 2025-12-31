<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  Table,
  Button,
  Space,
  Card,
  Row,
  Col,
  Typography,
  Modal,
  Form,
  FormItem,
  Input,
  Select,
  SelectOption,
  message,
  Popconfirm,
  Drawer,
  List,
  ListItem,
  ListItemMeta,
  Avatar,
  Empty,
  Tag,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
  UserOutlined,
  HistoryOutlined,
  FormOutlined,
} from '@ant-design/icons-vue'
import { homeroomService } from '@/services'
import type { HomeroomNote, ClassStudent, CreateHomeroomNoteRequest, UpdateHomeroomNoteRequest } from '@/types/homeroom'

const { Title, Text, Paragraph } = Typography
const { TextArea } = Input

// State
const loading = ref(false)
const notes = ref<HomeroomNote[]>([])
const students = ref<ClassStudent[]>([])
const className = ref('VII-A')

// Modal state
const modalVisible = ref(false)
const modalLoading = ref(false)
const editingNote = ref<HomeroomNote | null>(null)
const formRef = ref()

// History drawer state
const historyDrawerVisible = ref(false)
const selectedStudent = ref<ClassStudent | null>(null)
const studentNotes = ref<HomeroomNote[]>([])
const loadingHistory = ref(false)

// Form state
const formState = ref<{
  studentId: number | undefined
  content: string
}>({
  studentId: undefined,
  content: '',
})

// Mock data for development
const mockStudents: ClassStudent[] = [
  { id: 1, nis: '2024001', nisn: '0012345678', name: 'Ahmad Fauzi', isActive: true },
  { id: 2, nis: '2024002', nisn: '0012345679', name: 'Budi Santoso', isActive: true },
  { id: 3, nis: '2024003', nisn: '0012345680', name: 'Citra Dewi', isActive: true },
  { id: 4, nis: '2024004', nisn: '0012345681', name: 'Dian Pratama', isActive: true },
  { id: 5, nis: '2024005', nisn: '0012345682', name: 'Eka Putri', isActive: true },
]

const mockNotes: HomeroomNote[] = [
  { id: 1, studentId: 1, studentName: 'Ahmad Fauzi', studentNis: '2024001', teacherId: 1, teacherName: 'Ibu Sari', content: 'Siswa menunjukkan peningkatan dalam partisipasi kelas. Aktif bertanya dan menjawab pertanyaan.', createdAt: new Date().toISOString(), updatedAt: new Date().toISOString() },
  { id: 2, studentId: 4, studentName: 'Dian Pratama', studentNis: '2024004', teacherId: 1, teacherName: 'Ibu Sari', content: 'Perlu perhatian lebih dalam mata pelajaran Bahasa Inggris. Disarankan untuk mengikuti les tambahan.', createdAt: new Date(Date.now() - 86400000).toISOString(), updatedAt: new Date(Date.now() - 86400000).toISOString() },
  { id: 3, studentId: 2, studentName: 'Budi Santoso', studentNis: '2024002', teacherId: 1, teacherName: 'Ibu Sari', content: 'Siswa sangat rajin dan disiplin. Selalu mengumpulkan tugas tepat waktu.', createdAt: new Date(Date.now() - 172800000).toISOString(), updatedAt: new Date(Date.now() - 172800000).toISOString() },
  { id: 4, studentId: 3, studentName: 'Citra Dewi', studentNis: '2024003', teacherId: 1, teacherName: 'Ibu Sari', content: 'Siswa memiliki bakat dalam bidang seni. Disarankan untuk mengikuti ekstrakurikuler seni.', createdAt: new Date(Date.now() - 259200000).toISOString(), updatedAt: new Date(Date.now() - 259200000).toISOString() },
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
    width: 150,
  },
  {
    title: 'Catatan',
    dataIndex: 'content',
    key: 'content',
    ellipsis: true,
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

// Format date
const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

// Format full date
const formatFullDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
}

// Load notes
const loadNotes = async () => {
  loading.value = true
  try {
    const response = await homeroomService.getNotes({ pageSize: 100 })
    notes.value = response.data
  } catch {
    notes.value = mockNotes
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

// Filter option for student select
const filterStudentOption = (input: string, option: { label?: string }) => {
  return option.label?.toLowerCase().includes(input.toLowerCase()) ?? false
}

// Open modal for note
const openNoteModal = (note?: HomeroomNote) => {
  editingNote.value = note || null
  
  if (note) {
    formState.value = {
      studentId: note.studentId,
      content: note.content,
    }
  } else {
    formState.value = {
      studentId: undefined,
      content: '',
    }
  }
  
  modalVisible.value = true
}

// Close modal
const closeModal = () => {
  modalVisible.value = false
  editingNote.value = null
  formRef.value?.resetFields()
}

// Submit note
const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    modalLoading.value = true

    if (editingNote.value) {
      const data: UpdateHomeroomNoteRequest = {
        content: formState.value.content,
      }
      await homeroomService.updateNote(editingNote.value.id, data)
      message.success('Catatan berhasil diperbarui')
    } else {
      const data: CreateHomeroomNoteRequest = {
        studentId: formState.value.studentId!,
        content: formState.value.content,
      }
      await homeroomService.createNote(data)
      message.success('Catatan berhasil ditambahkan')
    }

    closeModal()
    loadNotes()
  } catch (err: unknown) {
    if (err && typeof err === 'object' && 'errorFields' in err) {
      return
    }
    message.error('Gagal menyimpan catatan')
  } finally {
    modalLoading.value = false
  }
}

// Delete note
const handleDelete = async (id: number) => {
  try {
    await homeroomService.deleteNote(id)
    message.success('Catatan berhasil dihapus')
    loadNotes()
  } catch {
    message.error('Gagal menghapus catatan')
  }
}

// Open history drawer
const openHistoryDrawer = async (student: ClassStudent) => {
  selectedStudent.value = student
  historyDrawerVisible.value = true
  loadingHistory.value = true
  
  try {
    const response = await homeroomService.getStudentNotes(student.id, { pageSize: 50 })
    studentNotes.value = response.data
  } catch {
    studentNotes.value = mockNotes.filter(n => n.studentId === student.id)
  } finally {
    loadingHistory.value = false
  }
}

// Close history drawer
const closeHistoryDrawer = () => {
  historyDrawerVisible.value = false
  selectedStudent.value = null
  studentNotes.value = []
}

// Quick add note for student
const quickAddNote = (student: ClassStudent) => {
  formState.value = {
    studentId: student.id,
    content: '',
  }
  editingNote.value = null
  modalVisible.value = true
}

onMounted(() => {
  loadNotes()
  loadStudents()
})
</script>

<template>
  <div class="homeroom-notes">
    <div class="page-header">
      <div>
        <Title :level="2" style="margin: 0">Catatan Wali Kelas</Title>
        <Text type="secondary">Kelas {{ className }}</Text>
      </div>
    </div>

    <Row :gutter="24">
      <!-- Notes List -->
      <Col :xs="24" :lg="16">
        <Card title="Daftar Catatan">
          <!-- Toolbar -->
          <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
            <Col :xs="24" :sm="12">
              <Text type="secondary">Total {{ notes.length }} catatan</Text>
            </Col>
            <Col :xs="24" :sm="12" class="toolbar-right">
              <Space>
                <Button @click="loadNotes">
                  <template #icon><ReloadOutlined /></template>
                  Refresh
                </Button>
                <Button type="primary" @click="openNoteModal()">
                  <template #icon><PlusOutlined /></template>
                  Tambah Catatan
                </Button>
              </Space>
            </Col>
          </Row>

          <!-- Table -->
          <Table
            :columns="columns"
            :data-source="notes"
            :loading="loading"
            :pagination="{ pageSize: 10 }"
            row-key="id"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'content'">
                <Paragraph :ellipsis="{ rows: 2 }" style="margin: 0">
                  {{ (record as HomeroomNote).content }}
                </Paragraph>
              </template>
              <template v-else-if="column.key === 'createdAt'">
                {{ formatDate((record as HomeroomNote).createdAt) }}
              </template>
              <template v-else-if="column.key === 'action'">
                <Space>
                  <Button type="link" size="small" @click="openNoteModal(record as HomeroomNote)">
                    <template #icon><EditOutlined /></template>
                  </Button>
                  <Popconfirm
                    title="Hapus catatan ini?"
                    ok-text="Ya"
                    cancel-text="Tidak"
                    @confirm="handleDelete((record as HomeroomNote).id)"
                  >
                    <Button type="link" size="small" danger>
                      <template #icon><DeleteOutlined /></template>
                    </Button>
                  </Popconfirm>
                </Space>
              </template>
            </template>
          </Table>
        </Card>
      </Col>

      <!-- Students Quick Access -->
      <Col :xs="24" :lg="8">
        <Card title="Siswa" class="students-card">
          <List
            :data-source="students"
            :loading="loading"
            size="small"
          >
            <template #renderItem="{ item }">
              <ListItem class="student-item">
                <ListItemMeta>
                  <template #avatar>
                    <Avatar :style="{ backgroundColor: '#f97316' }" size="small">
                      <template #icon><UserOutlined /></template>
                    </Avatar>
                  </template>
                  <template #title>
                    <span class="student-name">{{ (item as ClassStudent).name }}</span>
                  </template>
                  <template #description>
                    NIS: {{ (item as ClassStudent).nis }}
                  </template>
                </ListItemMeta>
                <template #actions>
                  <Space size="small">
                    <Button type="link" size="small" @click="quickAddNote(item as ClassStudent)">
                      <FormOutlined />
                    </Button>
                    <Button type="link" size="small" @click="openHistoryDrawer(item as ClassStudent)">
                      <HistoryOutlined />
                    </Button>
                  </Space>
                </template>
              </ListItem>
            </template>
          </List>
        </Card>
      </Col>
    </Row>

    <!-- Note Modal -->
    <Modal
      v-model:open="modalVisible"
      :title="editingNote ? 'Edit Catatan' : 'Tambah Catatan'"
      :confirm-loading="modalLoading"
      width="600px"
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
            :disabled="!!editingNote"
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
          label="Catatan"
          name="content"
          :rules="[{ required: true, message: 'Masukkan catatan' }]"
        >
          <TextArea
            v-model:value="formState.content"
            placeholder="Tulis catatan tentang perkembangan siswa..."
            :rows="6"
            show-count
            :maxlength="1000"
          />
        </FormItem>

        <div class="note-tips">
          <Text type="secondary">
            <strong>Tips:</strong> Catatan ini akan dapat dilihat oleh orang tua siswa melalui aplikasi mobile.
            Pastikan catatan bersifat konstruktif dan informatif.
          </Text>
        </div>
      </Form>
    </Modal>

    <!-- History Drawer -->
    <Drawer
      v-model:open="historyDrawerVisible"
      :title="`Riwayat Catatan - ${selectedStudent?.name || ''}`"
      width="500"
      @close="closeHistoryDrawer"
    >
      <div v-if="selectedStudent" class="history-content">
        <div class="student-info-header">
          <Avatar :style="{ backgroundColor: '#f97316' }" size="large">
            <template #icon><UserOutlined /></template>
          </Avatar>
          <div class="student-info-text">
            <Text strong>{{ selectedStudent.name }}</Text>
            <br />
            <Text type="secondary">NIS: {{ selectedStudent.nis }}</Text>
          </div>
        </div>

        <div class="add-note-btn">
          <Button type="primary" block @click="quickAddNote(selectedStudent)">
            <template #icon><PlusOutlined /></template>
            Tambah Catatan Baru
          </Button>
        </div>

        <div class="notes-timeline">
          <List
            v-if="studentNotes.length > 0"
            :data-source="studentNotes"
            :loading="loadingHistory"
          >
            <template #renderItem="{ item }">
              <Card size="small" class="note-card">
                <div class="note-date">
                  <Tag color="blue">{{ formatFullDate((item as HomeroomNote).createdAt) }}</Tag>
                </div>
                <Paragraph style="margin: 8px 0 0 0">
                  {{ (item as HomeroomNote).content }}
                </Paragraph>
                <div class="note-actions">
                  <Space>
                    <Button type="link" size="small" @click="openNoteModal(item as HomeroomNote)">
                      <EditOutlined /> Edit
                    </Button>
                    <Popconfirm
                      title="Hapus catatan ini?"
                      ok-text="Ya"
                      cancel-text="Tidak"
                      @confirm="handleDelete((item as HomeroomNote).id)"
                    >
                      <Button type="link" size="small" danger>
                        <DeleteOutlined /> Hapus
                      </Button>
                    </Popconfirm>
                  </Space>
                </div>
              </Card>
            </template>
          </List>
          <Empty v-else description="Belum ada catatan untuk siswa ini" />
        </div>
      </div>
    </Drawer>
  </div>
</template>

<style scoped>
.homeroom-notes {
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

.students-card {
  height: fit-content;
}

.student-item {
  cursor: pointer;
  transition: background-color 0.2s;
}

.student-item:hover {
  background-color: #fafafa;
}

.student-name {
  font-weight: 500;
}

.note-tips {
  padding: 12px;
  background: #fffbe6;
  border: 1px solid #ffe58f;
  border-radius: 6px;
}

.history-content {
  padding: 0;
}

.student-info-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.student-info-text {
  flex: 1;
}

.add-note-btn {
  margin: 16px 0;
}

.notes-timeline {
  margin-top: 16px;
}

.note-card {
  margin-bottom: 12px;
}

.note-date {
  margin-bottom: 8px;
}

.note-actions {
  margin-top: 12px;
  padding-top: 8px;
  border-top: 1px solid #f0f0f0;
}

@media (max-width: 768px) {
  .toolbar-right {
    margin-top: 16px;
    justify-content: flex-start;
  }
}
</style>
