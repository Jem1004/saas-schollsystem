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
  TimePicker,
  Descriptions,
  DescriptionsItem,
  DatePicker,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  PlusOutlined,
  SearchOutlined,
  DeleteOutlined,
  ReloadOutlined,
  EyeOutlined,
  FileProtectOutlined,
  CheckOutlined,
  PrinterOutlined,
  FilePdfOutlined,
} from '@ant-design/icons-vue'
import { useRouter } from 'vue-router'
import dayjs from 'dayjs'
import type { Dayjs } from 'dayjs'
import { bkService, schoolService } from '@/services'
import type { Permit, CreatePermitRequest } from '@/types/bk'
import type { Student, SchoolUser } from '@/types/school'
import { exportToPDF, formatPermitForExport, generatePermitPDF } from '@/utils/pdfExport'

const { Title, Text } = Typography
const { RangePicker } = DatePicker

const router = useRouter()

// Table state
const loading = ref(false)
const permits = ref<Permit[]>([])
const total = ref(0)
const pagination = reactive({
  current: 1,
  pageSize: 10,
})
const searchText = ref('')
const filterStatus = ref<string | undefined>(undefined)
const dateRange = ref<[Dayjs, Dayjs] | undefined>(undefined)

// Students and teachers for dropdown
const students = ref<Student[]>([])
const teachers = ref<SchoolUser[]>([])
const loadingStudents = ref(false)
const loadingTeachers = ref(false)

// Modal state
const modalVisible = ref(false)
const modalLoading = ref(false)
const returnModalVisible = ref(false)
const returnModalLoading = ref(false)
const previewModalVisible = ref(false)
const selectedPermit = ref<Permit | null>(null)

// Form state
const formRef = ref()
const formState = reactive<CreatePermitRequest & { exitTimeValue?: Dayjs }>({
  studentId: 0,
  reason: '',
  exitTime: '',
  responsibleTeacherId: 0,
  exitTimeValue: undefined,
})

// Return form state
const returnFormRef = ref()
const returnFormState = reactive<{ returnTimeValue?: Dayjs }>({
  returnTimeValue: undefined,
})

// Form rules
const formRules = {
  studentId: [{ required: true, message: 'Siswa wajib dipilih' }],
  reason: [{ required: true, message: 'Alasan wajib diisi' }],
  exitTimeValue: [{ required: true, message: 'Waktu keluar wajib diisi' }],
  responsibleTeacherId: [{ required: true, message: 'Guru penanggung jawab wajib dipilih' }],
}

const returnFormRules = {
  returnTimeValue: [{ required: true, message: 'Waktu kembali wajib diisi' }],
}

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
    title: 'Alasan',
    dataIndex: 'reason',
    key: 'reason',
    ellipsis: true,
  },
  {
    title: 'Waktu Keluar',
    dataIndex: 'exitTime',
    key: 'exitTime',
    width: 100,
  },
  {
    title: 'Status',
    key: 'status',
    width: 130,
    align: 'center',
  },
  {
    title: 'Aksi',
    key: 'action',
    width: 180,
    align: 'center',
  },
]

// Computed filtered data
const filteredPermits = computed(() => {
  let result = permits.value

  if (filterStatus.value === 'returned') {
    result = result.filter(p => p.returnTime)
  } else if (filterStatus.value === 'pending') {
    result = result.filter(p => !p.returnTime)
  }

  if (dateRange.value) {
    const [start, end] = dateRange.value
    result = result.filter(p => {
      const date = dayjs(p.createdAt)
      return date.isAfter(start.startOf('day')) && date.isBefore(end.endOf('day'))
    })
  }

  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(
      (p) =>
        p.studentName?.toLowerCase().includes(search) ||
        p.reason.toLowerCase().includes(search)
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

// Format time
const formatTime = (dateStr: string) => {
  return new Date(dateStr).toLocaleTimeString('id-ID', {
    hour: '2-digit',
    minute: '2-digit',
  })
}

// Load permits data
const loadPermits = async () => {
  loading.value = true
  try {
    const response = await bkService.getPermits({
      page: pagination.current,
      pageSize: pagination.pageSize,
      search: searchText.value,
    })
    permits.value = response.data || []
    total.value = response.total || 0
  } catch (err) {
    console.error('Failed to load permits:', err)
    permits.value = []
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

// Load teachers for dropdown
const loadTeachers = async () => {
  loadingTeachers.value = true
  try {
    const response = await schoolService.getUsers({ page_size: 1000 })
    teachers.value = response.users || []
  } catch (err) {
    console.error('Failed to load teachers:', err)
    teachers.value = []
  } finally {
    loadingTeachers.value = false
  }
}

// Handle table change
const handleTableChange: TableProps['onChange'] = (pag) => {
  pagination.current = pag.current || 1
  pagination.pageSize = pag.pageSize || 10
  loadPermits()
}

// Handle search
const handleSearch = () => {
  pagination.current = 1
  loadPermits()
}

// Handle filter change
const handleFilterChange = () => {
  pagination.current = 1
  loadPermits()
}

// Export to PDF
const handleExportPDF = () => {
  if (filteredPermits.value.length === 0) {
    message.warning('Tidak ada data untuk diekspor')
    return
  }

  const dateRangeStr = dateRange.value
    ? { start: dateRange.value[0].format('DD/MM/YYYY'), end: dateRange.value[1].format('DD/MM/YYYY') }
    : undefined

  exportToPDF({
    title: 'Laporan Izin Keluar Siswa',
    filename: `laporan-izin-keluar-${dayjs().format('YYYY-MM-DD')}`,
    columns: [
      { header: 'Tanggal', dataKey: 'createdAt' },
      { header: 'Siswa', dataKey: 'studentName' },
      { header: 'Kelas', dataKey: 'studentClass' },
      { header: 'Alasan', dataKey: 'reason' },
      { header: 'Waktu Keluar', dataKey: 'exitTime' },
      { header: 'Waktu Kembali', dataKey: 'returnTime' },
      { header: 'Status', dataKey: 'status' },
    ],
    data: filteredPermits.value.map(p => formatPermitForExport(p as unknown as Record<string, unknown>)),
    dateRange: dateRangeStr,
  })
  message.success('PDF berhasil diunduh')
}

// Open create modal
const openCreateModal = () => {
  resetForm()
  formState.exitTimeValue = dayjs()
  modalVisible.value = true
}

// Reset form
const resetForm = () => {
  formState.studentId = 0
  formState.reason = ''
  formState.exitTime = ''
  formState.responsibleTeacherId = 0
  formState.exitTimeValue = undefined
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

  // Convert time to ISO string
  if (formState.exitTimeValue) {
    const now = new Date()
    const exitTime = formState.exitTimeValue.toDate()
    exitTime.setFullYear(now.getFullYear(), now.getMonth(), now.getDate())
    formState.exitTime = exitTime.toISOString()
  }

  modalLoading.value = true
  try {
    await bkService.createPermit({
      studentId: formState.studentId,
      reason: formState.reason,
      exitTime: formState.exitTime,
      responsibleTeacherId: formState.responsibleTeacherId,
    })
    message.success('Izin keluar berhasil dibuat')
    modalVisible.value = false
    resetForm()
    loadPermits()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Terjadi kesalahan')
  } finally {
    modalLoading.value = false
  }
}

// Open return modal
const openReturnModal = (permit: Permit) => {
  selectedPermit.value = permit
  returnFormState.returnTimeValue = dayjs()
  returnModalVisible.value = true
}

// Handle return modal cancel
const handleReturnModalCancel = () => {
  returnModalVisible.value = false
  selectedPermit.value = null
  returnFormState.returnTimeValue = undefined
}

// Handle record return
const handleRecordReturn = async () => {
  try {
    await returnFormRef.value?.validate()
  } catch {
    return
  }

  if (!selectedPermit.value || !returnFormState.returnTimeValue) return

  const now = new Date()
  const returnTime = returnFormState.returnTimeValue.toDate()
  returnTime.setFullYear(now.getFullYear(), now.getMonth(), now.getDate())

  returnModalLoading.value = true
  try {
    await bkService.recordReturn(selectedPermit.value.id, {
      returnTime: returnTime.toISOString(),
    })
    message.success('Waktu kembali berhasil dicatat')
    returnModalVisible.value = false
    selectedPermit.value = null
    returnFormState.returnTimeValue = undefined
    loadPermits()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Terjadi kesalahan')
  } finally {
    returnModalLoading.value = false
  }
}

// Open preview modal
const openPreviewModal = (permit: Permit) => {
  selectedPermit.value = permit
  previewModalVisible.value = true
}

// Handle print permit
const handlePrintPermit = async () => {
  if (!selectedPermit.value) return

  try {
    const docData = await bkService.getPermitDocument(selectedPermit.value.id)
    generatePermitPDF(docData)
  } catch (error) {
    console.error('Failed to generate permit document:', error)
    message.error('Gagal membuat dokumen izin keluar')
  }
}

// Handle delete
const handleDelete = async (permit: Permit) => {
  try {
    await bkService.deletePermit(permit.id)
    message.success('Izin keluar berhasil dihapus')
    loadPermits()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Gagal menghapus izin keluar')
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

// Filter teacher options
const filterTeacherOption = (input: string, option: { label: string }) => {
  return option.label.toLowerCase().includes(input.toLowerCase())
}

onMounted(() => {
  loadPermits()
  loadStudents()
  loadTeachers()
})
</script>

<template>
  <div class="permit-management">
    <div class="page-header">
      <Title :level="2" style="margin: 0">Izin Keluar Sekolah</Title>
    </div>

    <Card>
      <!-- Toolbar -->
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col :xs="24" :sm="24" :md="16">
          <Space wrap>
            <Input
              v-model:value="searchText"
              placeholder="Cari siswa atau alasan..."
              allow-clear
              style="width: 220px"
              @press-enter="handleSearch"
            >
              <template #prefix>
                <SearchOutlined />
              </template>
            </Input>
            <RangePicker v-model:value="dateRange" format="DD/MM/YYYY" :placeholder="['Dari Tanggal', 'Sampai Tanggal']" style="width: 250px" @change="handleFilterChange" />
            <Select
              v-model:value="filterStatus"
              placeholder="Filter Status"
              allow-clear
              style="width: 150px"
              @change="handleFilterChange"
            >
              <Select.Option value="pending">Belum Kembali</Select.Option>
              <Select.Option value="returned">Sudah Kembali</Select.Option>
            </Select>
          </Space>
        </Col>
        <Col :xs="24" :sm="24" :md="8" class="toolbar-right">
          <Space>
            <Button @click="handleExportPDF">
              <template #icon><FilePdfOutlined /></template>
              Export PDF
            </Button>
            <Button @click="loadPermits">
              <template #icon><ReloadOutlined /></template>
            </Button>
            <Button type="primary" @click="openCreateModal">
              <template #icon><PlusOutlined /></template>
              Buat Izin
            </Button>
          </Space>
        </Col>
      </Row>

      <!-- Table -->
      <Table
        :columns="columns"
        :data-source="filteredPermits"
        :loading="loading"
        :pagination="{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: total,
          showSizeChanger: true,
          showTotal: (total: number) => `Total ${total} izin`,
        }"
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'createdAt'">
            {{ formatDate((record as Permit).createdAt) }}
          </template>
          <template v-else-if="column.key === 'studentName'">
            <a @click="viewStudentProfile((record as Permit).studentId)">
              {{ (record as Permit).studentName }}
            </a>
          </template>
          <template v-else-if="column.key === 'studentClass'">
            <Tag color="blue">{{ (record as Permit).studentClass }}</Tag>
          </template>
          <template v-else-if="column.key === 'exitTime'">
            {{ formatTime((record as Permit).exitTime) }}
          </template>
          <template v-else-if="column.key === 'status'">
            <Tag :color="(record as Permit).returnTime ? 'success' : 'processing'">
              {{ (record as Permit).returnTime ? 'Sudah Kembali' : 'Belum Kembali' }}
            </Tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <Space>
              <Button size="small" @click="openPreviewModal(record as Permit)">
                <template #icon><EyeOutlined /></template>
              </Button>
              <Button
                v-if="!(record as Permit).returnTime"
                size="small"
                type="primary"
                @click="openReturnModal(record as Permit)"
              >
                <template #icon><CheckOutlined /></template>
                Kembali
              </Button>
              <Popconfirm
                title="Hapus izin keluar ini?"
                description="Data izin akan dihapus permanen."
                ok-text="Ya, Hapus"
                cancel-text="Batal"
                @confirm="handleDelete(record as Permit)"
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
      title="Buat Izin Keluar Baru"
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
        <FormItem label="Alasan" name="reason" required>
          <Textarea
            v-model:value="formState.reason"
            placeholder="Jelaskan alasan izin keluar..."
            :rows="3"
          />
        </FormItem>
        <Row :gutter="16">
          <Col :span="12">
            <FormItem label="Waktu Keluar" name="exitTimeValue" required>
              <TimePicker
                v-model:value="formState.exitTimeValue"
                format="HH:mm"
                style="width: 100%"
                placeholder="Pilih waktu"
              />
            </FormItem>
          </Col>
          <Col :span="12">
            <FormItem label="Guru Penanggung Jawab" name="responsibleTeacherId" required>
              <Select
                v-model:value="formState.responsibleTeacherId"
                placeholder="Pilih guru"
                :loading="loadingTeachers"
                show-search
                :filter-option="filterTeacherOption"
                :options="teachers.map(t => ({ value: t.id, label: t.name || t.username }))"
              />
            </FormItem>
          </Col>
        </Row>
      </Form>
    </Modal>

    <!-- Return Modal -->
    <Modal
      v-model:open="returnModalVisible"
      title="Catat Waktu Kembali"
      :confirm-loading="returnModalLoading"
      @ok="handleRecordReturn"
      @cancel="handleReturnModalCancel"
      width="400px"
    >
      <div v-if="selectedPermit" style="margin-bottom: 16px">
        <Text strong>{{ selectedPermit.studentName }}</Text>
        <Text type="secondary"> - {{ selectedPermit.studentClass }}</Text>
      </div>
      <Form
        ref="returnFormRef"
        :model="returnFormState"
        :rules="returnFormRules"
        layout="vertical"
      >
        <FormItem label="Waktu Kembali" name="returnTimeValue" required>
          <TimePicker
            v-model:value="returnFormState.returnTimeValue"
            format="HH:mm"
            style="width: 100%"
            placeholder="Pilih waktu kembali"
          />
        </FormItem>
      </Form>
    </Modal>

    <!-- Preview Modal -->
    <Modal
      v-model:open="previewModalVisible"
      title="Detail Izin Keluar"
      :footer="null"
      width="600px"
    >
      <div v-if="selectedPermit" class="permit-preview">
        <div class="permit-header">
          <FileProtectOutlined class="permit-icon" />
          <Title :level="4" style="margin: 0">Surat Izin Keluar Sekolah</Title>
        </div>

        <Descriptions :column="1" bordered size="small" style="margin-top: 16px">
          <DescriptionsItem label="Nama Siswa">{{ selectedPermit.studentName }}</DescriptionsItem>
          <DescriptionsItem label="NIS">{{ selectedPermit.studentNis }}</DescriptionsItem>
          <DescriptionsItem label="NISN">{{ selectedPermit.studentNisn }}</DescriptionsItem>
          <DescriptionsItem label="Kelas">{{ selectedPermit.studentClass }}</DescriptionsItem>
          <DescriptionsItem label="Alasan">{{ selectedPermit.reason }}</DescriptionsItem>
          <DescriptionsItem label="Waktu Keluar">
            {{ formatDate(selectedPermit.exitTime) }} {{ formatTime(selectedPermit.exitTime) }}
          </DescriptionsItem>
          <DescriptionsItem label="Waktu Kembali">
            <span v-if="selectedPermit.returnTime">
              {{ formatDate(selectedPermit.returnTime) }} {{ formatTime(selectedPermit.returnTime) }}
            </span>
            <Tag v-else color="processing">Belum Kembali</Tag>
          </DescriptionsItem>
          <DescriptionsItem label="Guru Penanggung Jawab">
            {{ selectedPermit.responsibleTeacherName }}
          </DescriptionsItem>
          <DescriptionsItem label="Dibuat Oleh">{{ selectedPermit.createdByName }}</DescriptionsItem>
          <DescriptionsItem label="Tanggal Dibuat">
            {{ formatDate(selectedPermit.createdAt) }}
          </DescriptionsItem>
        </Descriptions>

        <div class="permit-actions">
          <Button type="primary" @click="handlePrintPermit">
            <template #icon><PrinterOutlined /></template>
            Cetak Dokumen
          </Button>
        </div>
      </div>
    </Modal>
  </div>
</template>

<style scoped>
.permit-management {
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

.permit-preview {
  padding: 16px 0;
}

.permit-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding-bottom: 16px;
  border-bottom: 2px solid #f97316;
}

.permit-icon {
  font-size: 32px;
  color: #f97316;
}

.permit-actions {
  margin-top: 24px;
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
