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
  message,
  Popconfirm,
  Card,
  Row,
  Col,
  Typography,
  Textarea,
  Alert,
  Divider,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  PlusOutlined,
  SearchOutlined,
  DeleteOutlined,
  ReloadOutlined,
  EyeOutlined,
  LockOutlined,
  UnlockOutlined,
  SafetyOutlined,
} from '@ant-design/icons-vue'
import { useRouter } from 'vue-router'
import { bkService, schoolService } from '@/services'
import type { CounselingNote, CreateCounselingNoteRequest } from '@/types/bk'
import type { Student } from '@/types/school'
import { SensitiveDataField, ConfidentialBadge, ConfirmationDialog } from '@/components'

const { Title, Text, Paragraph } = Typography

const router = useRouter()

// Table state
const loading = ref(false)
const counselingNotes = ref<CounselingNote[]>([])
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
const viewModalVisible = ref(false)
const selectedNote = ref<CounselingNote | null>(null)

// Form state
const formRef = ref()
const formState = reactive<CreateCounselingNoteRequest>({
  studentId: 0,
  internalNote: '',
  parentSummary: '',
})

// Form rules
const formRules = {
  studentId: [{ required: true, message: 'Siswa wajib dipilih' }],
  internalNote: [{ required: true, message: 'Catatan internal wajib diisi' }],
}

// Mock data for development
const mockCounselingNotes: CounselingNote[] = [
  { id: 1, studentId: 1, studentName: 'Ahmad Fauzi', studentNis: '2024001', studentClass: 'VII-A', internalNote: 'Siswa menunjukkan tanda-tanda stres karena tekanan akademik. Perlu pendampingan lebih lanjut dalam manajemen waktu dan teknik belajar efektif. Orang tua perlu dilibatkan dalam proses ini.', parentSummary: 'Siswa memerlukan dukungan dalam mengelola waktu belajar. Mohon bantuan orang tua untuk memantau jadwal belajar di rumah.', createdBy: 1, createdByName: 'Guru BK', createdAt: new Date().toISOString() },
  { id: 2, studentId: 2, studentName: 'Budi Santoso', studentNis: '2024002', studentClass: 'VII-B', internalNote: 'Sesi konseling terkait konflik dengan teman sekelas. Siswa merasa dikucilkan. Perlu mediasi dengan pihak terkait.', parentSummary: 'Siswa sedang dalam proses adaptasi sosial. Mohon dukungan orang tua untuk memberikan semangat.', createdBy: 1, createdByName: 'Guru BK', createdAt: new Date(Date.now() - 86400000).toISOString() },
  { id: 3, studentId: 3, studentName: 'Citra Dewi', studentNis: '2024003', studentClass: 'VIII-A', internalNote: 'Konseling rutin bulanan. Siswa menunjukkan perkembangan positif dalam manajemen emosi. Tetap perlu monitoring.', parentSummary: 'Perkembangan positif dalam pengelolaan emosi. Terima kasih atas dukungan orang tua.', createdBy: 1, createdByName: 'Guru BK', createdAt: new Date(Date.now() - 172800000).toISOString() },
  { id: 4, studentId: 4, studentName: 'Dian Pratama', studentNis: '2024004', studentClass: 'IX-A', internalNote: 'Siswa mengalami kecemasan menghadapi ujian nasional. Diberikan teknik relaksasi dan strategi menghadapi ujian.', createdBy: 1, createdByName: 'Guru BK', createdAt: new Date(Date.now() - 259200000).toISOString() },
]

const mockStudents: Student[] = [
  { id: 1, schoolId: 1, classId: 1, className: 'VII-A', nis: '2024001', nisn: '0012345678', name: 'Ahmad Fauzi', isActive: true, createdAt: '', updatedAt: '' },
  { id: 2, schoolId: 1, classId: 2, className: 'VII-B', nis: '2024002', nisn: '0012345679', name: 'Budi Santoso', isActive: true, createdAt: '', updatedAt: '' },
  { id: 3, schoolId: 1, classId: 3, className: 'VIII-A', nis: '2024003', nisn: '0012345680', name: 'Citra Dewi', isActive: true, createdAt: '', updatedAt: '' },
  { id: 4, schoolId: 1, classId: 4, className: 'IX-A', nis: '2024004', nisn: '0012345681', name: 'Dian Pratama', isActive: true, createdAt: '', updatedAt: '' },
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
    title: 'Ringkasan untuk Orang Tua',
    key: 'parentSummary',
    ellipsis: true,
  },
  {
    title: 'Status',
    key: 'status',
    width: 150,
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
const filteredNotes = computed(() => {
  if (!searchText.value) return counselingNotes.value

  const search = searchText.value.toLowerCase()
  return counselingNotes.value.filter(
    (n) =>
      n.studentName?.toLowerCase().includes(search) ||
      n.parentSummary?.toLowerCase().includes(search)
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

// Load counseling notes data
const loadCounselingNotes = async () => {
  loading.value = true
  try {
    const response = await bkService.getCounselingNotes({
      page: pagination.current,
      pageSize: pagination.pageSize,
      search: searchText.value,
    })
    counselingNotes.value = response.data
    total.value = response.total
  } catch {
    counselingNotes.value = mockCounselingNotes
    total.value = mockCounselingNotes.length
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
  loadCounselingNotes()
}

// Handle search
const handleSearch = () => {
  pagination.current = 1
  loadCounselingNotes()
}

// Open create modal
const openCreateModal = () => {
  resetForm()
  modalVisible.value = true
}

// Reset form
const resetForm = () => {
  formState.studentId = 0
  formState.internalNote = ''
  formState.parentSummary = ''
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
    await bkService.createCounselingNote(formState)
    message.success('Catatan konseling berhasil disimpan')
    modalVisible.value = false
    resetForm()
    loadCounselingNotes()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } } }
    message.error(err.response?.data?.message || 'Terjadi kesalahan')
  } finally {
    modalLoading.value = false
  }
}

// Handle delete
const handleDelete = async (note: CounselingNote) => {
  try {
    await bkService.deleteCounselingNote(note.id)
    message.success('Catatan konseling berhasil dihapus')
    loadCounselingNotes()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } } }
    message.error(err.response?.data?.message || 'Gagal menghapus catatan')
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

// Sensitive data access confirmation
const sensitiveAccessConfirmVisible = ref(false)
const pendingViewNote = ref<CounselingNote | null>(null)

const requestViewInternalNote = (note: CounselingNote) => {
  pendingViewNote.value = note
  sensitiveAccessConfirmVisible.value = true
}

const confirmViewInternalNote = () => {
  if (pendingViewNote.value) {
    selectedNote.value = pendingViewNote.value
    viewModalVisible.value = true
  }
  sensitiveAccessConfirmVisible.value = false
  pendingViewNote.value = null
}

const cancelViewInternalNote = () => {
  sensitiveAccessConfirmVisible.value = false
  pendingViewNote.value = null
}

onMounted(() => {
  loadCounselingNotes()
  loadStudents()
})
</script>

<template>
  <div class="counseling-management">
    <div class="page-header">
      <Title :level="2" style="margin: 0">Catatan Konseling</Title>
    </div>

    <!-- Privacy Notice -->
    <Alert
      type="warning"
      show-icon
      style="margin-bottom: 24px"
    >
      <template #icon>
        <SafetyOutlined />
      </template>
      <template #message>
        <Text strong>Perhatian Kerahasiaan Data</Text>
      </template>
      <template #description>
        <div>
          Catatan internal bersifat rahasia dan hanya dapat diakses oleh Guru BK. 
          Ringkasan untuk orang tua akan ditampilkan di aplikasi mobile orang tua.
          <ConfidentialBadge type="internal" size="small" style="margin-left: 8px" />
        </div>
      </template>
    </Alert>

    <Card>
      <!-- Toolbar -->
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col :xs="24" :sm="24" :md="16">
          <Space wrap>
            <Input
              v-model:value="searchText"
              placeholder="Cari siswa..."
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
            <Button @click="loadCounselingNotes">
              <template #icon><ReloadOutlined /></template>
              Refresh
            </Button>
            <Button type="primary" @click="openCreateModal">
              <template #icon><PlusOutlined /></template>
              Buat Catatan
            </Button>
          </Space>
        </Col>
      </Row>

      <!-- Table -->
      <Table
        :columns="columns"
        :data-source="filteredNotes"
        :loading="loading"
        :pagination="{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: total,
          showSizeChanger: true,
          showTotal: (total: number) => `Total ${total} catatan`,
        }"
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'createdAt'">
            {{ formatDate((record as CounselingNote).createdAt) }}
          </template>
          <template v-else-if="column.key === 'studentName'">
            <a @click="viewStudentProfile((record as CounselingNote).studentId)">
              {{ (record as CounselingNote).studentName }}
            </a>
          </template>
          <template v-else-if="column.key === 'studentClass'">
            <Tag color="blue">{{ (record as CounselingNote).studentClass }}</Tag>
          </template>
          <template v-else-if="column.key === 'parentSummary'">
            <span v-if="(record as CounselingNote).parentSummary">
              {{ (record as CounselingNote).parentSummary }}
            </span>
            <Text v-else type="secondary" italic>Tidak ada ringkasan untuk orang tua</Text>
          </template>
          <template v-else-if="column.key === 'status'">
            <Tag v-if="(record as CounselingNote).parentSummary" color="success">
              <UnlockOutlined /> Dibagikan
            </Tag>
            <Tag v-else color="default">
              <LockOutlined /> Internal
            </Tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <Space>
              <Button size="small" @click="requestViewInternalNote(record as CounselingNote)">
                <template #icon><EyeOutlined /></template>
              </Button>
              <Popconfirm
                title="Hapus catatan konseling ini?"
                description="Data catatan akan dihapus permanen."
                ok-text="Ya, Hapus"
                cancel-text="Batal"
                @confirm="handleDelete(record as CounselingNote)"
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
      title="Buat Catatan Konseling Baru"
      :confirm-loading="modalLoading"
      @ok="handleSubmit"
      @cancel="handleModalCancel"
      width="700px"
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

        <Divider orientation="left">
          <LockOutlined /> Catatan Internal (Rahasia)
        </Divider>
        
        <Alert
          type="info"
          message="Catatan ini hanya dapat dilihat oleh Guru BK"
          style="margin-bottom: 16px"
        />
        
        <FormItem label="Catatan Internal" name="internalNote" required>
          <Textarea
            v-model:value="formState.internalNote"
            placeholder="Tuliskan catatan detail hasil konseling. Catatan ini bersifat rahasia dan tidak akan dibagikan kepada orang tua atau wali kelas..."
            :rows="6"
          />
        </FormItem>

        <Divider orientation="left">
          <UnlockOutlined /> Ringkasan untuk Orang Tua (Opsional)
        </Divider>
        
        <Alert
          type="success"
          message="Ringkasan ini akan ditampilkan di aplikasi mobile orang tua"
          style="margin-bottom: 16px"
        />
        
        <FormItem label="Ringkasan untuk Orang Tua" name="parentSummary">
          <Textarea
            v-model:value="formState.parentSummary"
            placeholder="Tuliskan ringkasan yang aman untuk dibagikan kepada orang tua. Hindari informasi sensitif..."
            :rows="4"
          />
        </FormItem>
      </Form>
    </Modal>

    <!-- View Modal -->
    <Modal
      v-model:open="viewModalVisible"
      title="Detail Catatan Konseling"
      :footer="null"
      width="700px"
    >
      <div v-if="selectedNote" class="counseling-detail">
        <div class="student-info">
          <Text strong>{{ selectedNote.studentName }}</Text>
          <Tag color="blue" style="margin-left: 8px">{{ selectedNote.studentClass }}</Tag>
          <Text type="secondary" style="margin-left: 16px">
            {{ formatDate(selectedNote.createdAt) }}
          </Text>
        </div>

        <Divider orientation="left">
          <LockOutlined /> Catatan Internal
          <ConfidentialBadge type="internal" size="small" style="margin-left: 8px" />
        </Divider>
        
        <Card class="internal-note-card" size="small">
          <SensitiveDataField
            :value="selectedNote.internalNote"
            :blur-by-default="true"
            :require-confirmation="false"
            :show-indicator="true"
            confirm-title="Catatan Internal Rahasia"
            confirm-description="Catatan ini bersifat rahasia dan hanya untuk keperluan konseling internal."
          />
        </Card>

        <template v-if="selectedNote.parentSummary">
          <Divider orientation="left">
            <UnlockOutlined /> Ringkasan untuk Orang Tua
          </Divider>
          
          <Card class="parent-summary-card" size="small">
            <Paragraph>{{ selectedNote.parentSummary }}</Paragraph>
          </Card>
        </template>

        <div class="note-meta">
          <Text type="secondary">Dibuat oleh: {{ selectedNote.createdByName }}</Text>
        </div>
      </div>
    </Modal>

    <!-- Sensitive Data Access Confirmation -->
    <ConfirmationDialog
      v-model:open="sensitiveAccessConfirmVisible"
      title="Akses Data Sensitif"
      message="Anda akan mengakses catatan konseling yang bersifat rahasia."
      description="Akses ke data ini akan dicatat dalam sistem audit. Pastikan Anda memiliki keperluan yang sah untuk mengakses data ini."
      type="sensitive"
      confirm-text="Ya, Tampilkan"
      cancel-text="Batal"
      @confirm="confirmViewInternalNote"
      @cancel="cancelViewInternalNote"
    />
  </div>
</template>

<style scoped>
.counseling-management {
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

.counseling-detail {
  padding: 16px 0;
}

.student-info {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.internal-note-card {
  background-color: #fff7e6;
  border-color: #ffd591;
}

.parent-summary-card {
  background-color: #f6ffed;
  border-color: #b7eb8f;
}

.note-meta {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

@media (max-width: 768px) {
  .toolbar-right {
    margin-top: 16px;
    justify-content: flex-start;
  }
}
</style>
