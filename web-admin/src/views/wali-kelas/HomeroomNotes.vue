<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import {
  Table, Button, Space, Card, Row, Col, Typography, Modal, Form,
  FormItem, Input, Select, SelectOption, message, Popconfirm, Drawer,
  List, ListItem, ListItemMeta, Avatar, Empty, Tag, Alert,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  PlusOutlined, EditOutlined, DeleteOutlined, ReloadOutlined,
  UserOutlined, HistoryOutlined, FormOutlined,
} from '@ant-design/icons-vue'
import { homeroomService } from '@/services'
import { useClassInfo, useClassStudents, useDateFormat, extractArrayFromResponse } from '@/composables/useWaliKelas'
import type { HomeroomNote, CreateHomeroomNoteRequest, UpdateHomeroomNoteRequest, ClassStudent } from '@/types/homeroom'

const { Title, Text, Paragraph } = Typography
const { TextArea } = Input

// Composables
const { className, loadClassInfo } = useClassInfo()
const { students, loadStudents, filterStudentOption } = useClassStudents()
const { formatDate, formatFullDate } = useDateFormat()

// Mounted state
const isMounted = ref(true)

// State
const loading = ref(false)
const error = ref<string | null>(null)
const notes = ref<HomeroomNote[]>([])

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
const formState = ref<{ studentId: number | undefined; content: string }>({
  studentId: undefined,
  content: '',
})

// Table columns
const columns: TableProps['columns'] = [
  { title: 'NIS', dataIndex: 'studentNis', key: 'studentNis', width: 100 },
  { title: 'Nama Siswa', dataIndex: 'studentName', key: 'studentName', width: 150 },
  { title: 'Catatan', dataIndex: 'content', key: 'content', ellipsis: true },
  { title: 'Tanggal', dataIndex: 'createdAt', key: 'createdAt', width: 150 },
  { title: 'Aksi', key: 'action', width: 150, align: 'center' },
]

// Load notes
const loadNotes = async () => {
  if (!isMounted.value) return
  loading.value = true
  error.value = null
  try {
    const response = await homeroomService.getNotes({ pageSize: 100 })
    if (isMounted.value) notes.value = extractArrayFromResponse<HomeroomNote>(response)
  } catch (err) {
    console.error('Failed to load notes:', err)
    if (isMounted.value) {
      notes.value = []
      error.value = 'Gagal memuat data catatan'
    }
  } finally {
    if (isMounted.value) loading.value = false
  }
}

// Open modal
const openNoteModal = (note?: HomeroomNote) => {
  editingNote.value = note || null
  formState.value = note
    ? { studentId: note.studentId, content: note.content }
    : { studentId: undefined, content: '' }
  modalVisible.value = true
}

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
      const data: UpdateHomeroomNoteRequest = { content: formState.value.content }
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
    if (err && typeof err === 'object' && 'errorFields' in err) return
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

// History drawer
const openHistoryDrawer = async (student: ClassStudent) => {
  selectedStudent.value = student
  historyDrawerVisible.value = true
  loadingHistory.value = true
  try {
    const response = await homeroomService.getStudentNotes(student.id, { pageSize: 50 })
    studentNotes.value = extractArrayFromResponse<HomeroomNote>(response)
  } catch (err) {
    console.error('Failed to load student notes:', err)
    studentNotes.value = []
  } finally {
    loadingHistory.value = false
  }
}

const closeHistoryDrawer = () => {
  historyDrawerVisible.value = false
  selectedStudent.value = null
  studentNotes.value = []
}

// Quick add note
const quickAddNote = (student: ClassStudent) => {
  formState.value = { studentId: student.id, content: '' }
  editingNote.value = null
  modalVisible.value = true
}

onMounted(() => {
  loadClassInfo()
  loadNotes()
  loadStudents()
})

onUnmounted(() => { isMounted.value = false })
</script>

<template>
  <div class="wali-kelas-page">
    <div class="page-header">
      <div>
        <Title :level="2" style="margin: 0">Catatan Wali Kelas</Title>
        <Text type="secondary">Kelas {{ className }}</Text>
      </div>
    </div>

    <Alert v-if="error" type="error" :message="error" show-icon closable style="margin-bottom: 16px" @close="error = null" />

    <Row :gutter="24">
      <Col :xs="24" :lg="16">
        <Card title="Daftar Catatan">
          <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
            <Col :xs="24" :sm="12"><Text type="secondary">Total {{ notes.length }} catatan</Text></Col>
            <Col :xs="24" :sm="12" class="toolbar-right">
              <Space>
                <Button @click="loadNotes"><template #icon><ReloadOutlined /></template></Button>
                <Button type="primary" @click="openNoteModal()"><template #icon><PlusOutlined /></template>Tambah Catatan</Button>
              </Space>
            </Col>
          </Row>

          <Table :columns="columns" :data-source="notes" :loading="loading" :pagination="{ pageSize: 10 }" row-key="id">
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'content'">
                <Paragraph :ellipsis="{ rows: 2 }" :content="(record as HomeroomNote).content" style="margin: 0" />
              </template>
              <template v-else-if="column.key === 'createdAt'">{{ formatDate((record as HomeroomNote).createdAt) }}</template>
              <template v-else-if="column.key === 'action'">
                <Space>
                  <Button type="link" size="small" @click="openNoteModal(record as HomeroomNote)"><template #icon><EditOutlined /></template></Button>
                  <Popconfirm title="Hapus catatan ini?" ok-text="Ya" cancel-text="Tidak" @confirm="handleDelete((record as HomeroomNote).id)">
                    <Button type="link" size="small" danger><template #icon><DeleteOutlined /></template></Button>
                  </Popconfirm>
                </Space>
              </template>
            </template>
          </Table>
        </Card>
      </Col>

      <Col :xs="24" :lg="8">
        <Card title="Siswa" class="students-card">
          <List :data-source="students" :loading="loading" size="small">
            <template #renderItem="{ item }">
              <ListItem class="clickable-item">
                <ListItemMeta>
                  <template #avatar><Avatar :style="{ backgroundColor: '#f97316' }" size="small"><template #icon><UserOutlined /></template></Avatar></template>
                  <template #title><span class="student-name">{{ (item as ClassStudent).name }}</span></template>
                  <template #description>NIS: {{ (item as ClassStudent).nis }}</template>
                </ListItemMeta>
                <template #actions>
                  <Space size="small">
                    <Button type="link" size="small" @click="quickAddNote(item as ClassStudent)"><FormOutlined /></Button>
                    <Button type="link" size="small" @click="openHistoryDrawer(item as ClassStudent)"><HistoryOutlined /></Button>
                  </Space>
                </template>
              </ListItem>
            </template>
          </List>
        </Card>
      </Col>
    </Row>

    <!-- Note Modal -->
    <Modal v-model:open="modalVisible" :title="editingNote ? 'Edit Catatan' : 'Tambah Catatan'" :confirm-loading="modalLoading" width="600px" @ok="handleSubmit" @cancel="closeModal">
      <Form ref="formRef" :model="formState" layout="vertical" style="margin-top: 16px">
        <FormItem label="Siswa" name="studentId" :rules="[{ required: true, message: 'Pilih siswa' }]">
          <Select v-model:value="formState.studentId" placeholder="Pilih siswa" :disabled="!!editingNote" show-search :filter-option="filterStudentOption">
            <SelectOption v-for="student in students" :key="student.id" :value="student.id" :label="student.name">{{ student.nis }} - {{ student.name }}</SelectOption>
          </Select>
        </FormItem>
        <FormItem label="Catatan" name="content" :rules="[{ required: true, message: 'Masukkan catatan' }]">
          <TextArea v-model:value="formState.content" placeholder="Tulis catatan tentang perkembangan siswa..." :rows="6" show-count :maxlength="1000" />
        </FormItem>
        <div class="note-tips">
          <Text type="secondary"><strong>Tips:</strong> Catatan ini akan dapat dilihat oleh orang tua siswa melalui aplikasi mobile. Pastikan catatan bersifat konstruktif dan informatif.</Text>
        </div>
      </Form>
    </Modal>

    <!-- History Drawer -->
    <Drawer v-model:open="historyDrawerVisible" :title="`Riwayat Catatan - ${selectedStudent?.name || ''}`" width="500" @close="closeHistoryDrawer">
      <div v-if="selectedStudent" class="history-content">
        <div class="student-info-header">
          <Avatar :style="{ backgroundColor: '#f97316' }" size="large"><template #icon><UserOutlined /></template></Avatar>
          <div class="student-info-text">
            <Text strong>{{ selectedStudent.name }}</Text><br />
            <Text type="secondary">NIS: {{ selectedStudent.nis }}</Text>
          </div>
        </div>
        <div class="add-note-btn">
          <Button type="primary" block @click="quickAddNote(selectedStudent)"><template #icon><PlusOutlined /></template>Tambah Catatan Baru</Button>
        </div>
        <div class="notes-timeline">
          <List v-if="studentNotes.length > 0" :data-source="studentNotes" :loading="loadingHistory">
            <template #renderItem="{ item }">
              <Card size="small" class="note-card">
                <div class="note-date"><Tag color="blue">{{ formatFullDate((item as HomeroomNote).createdAt) }}</Tag></div>
                <Paragraph style="margin: 8px 0 0 0">{{ (item as HomeroomNote).content }}</Paragraph>
                <div class="note-actions">
                  <Space>
                    <Button type="link" size="small" @click="openNoteModal(item as HomeroomNote)"><EditOutlined /> Edit</Button>
                    <Popconfirm title="Hapus catatan ini?" ok-text="Ya" cancel-text="Tidak" @confirm="handleDelete((item as HomeroomNote).id)">
                      <Button type="link" size="small" danger><DeleteOutlined /> Hapus</Button>
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
.wali-kelas-page { padding: 0; }
.page-header { margin-bottom: 24px; display: flex; justify-content: space-between; align-items: flex-start; flex-wrap: wrap; gap: 8px; }
.toolbar { margin-bottom: 16px; }
.toolbar-right { display: flex; justify-content: flex-end; }
.students-card { height: fit-content; }
.clickable-item { cursor: pointer; transition: background-color 0.2s; }
.clickable-item:hover { background-color: #fafafa; }
.student-name { font-weight: 500; }
.note-tips { padding: 12px; background: #fffbe6; border: 1px solid #ffe58f; border-radius: 6px; }
.history-content { padding: 0; }
.student-info-header { display: flex; align-items: center; gap: 16px; padding-bottom: 16px; border-bottom: 1px solid #f0f0f0; }
.student-info-text { flex: 1; }
.add-note-btn { margin: 16px 0; }
.notes-timeline { margin-top: 16px; }
.note-card { margin-bottom: 12px; }
.note-date { margin-bottom: 8px; }
.note-actions { margin-top: 12px; padding-top: 8px; border-top: 1px solid #f0f0f0; }
@media (max-width: 768px) { .toolbar-right { margin-top: 16px; justify-content: flex-start; } }
</style>
