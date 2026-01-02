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
  DatePicker,
  Progress,
  Statistic,
  Modal,
  Form,
  FormItem,
  Select,
  SelectOption,
  TimePicker,
  message,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  ReloadOutlined,
  CalendarOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  ClockCircleOutlined,
  EditOutlined,
  PlusOutlined,
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import type { Dayjs } from 'dayjs'
import { homeroomService } from '@/services'
import type { StudentAttendance, ClassStudent, ManualAttendanceRequest } from '@/types/homeroom'

const { Title, Text } = Typography

// State
const loading = ref(false)
const attendanceData = ref<StudentAttendance[]>([])
const students = ref<ClassStudent[]>([])
const selectedDate = ref<Dayjs>(dayjs())
const className = ref('VII-A')

// Modal state
const modalVisible = ref(false)
const modalLoading = ref(false)
const editingRecord = ref<StudentAttendance | null>(null)
const formRef = ref()

// Form state
const formState = ref<{
  studentId: number | undefined
  status: 'present' | 'absent' | 'late' | 'excused'
  checkInTime: Dayjs | undefined
  checkOutTime: Dayjs | undefined
}>({
  studentId: undefined,
  status: 'present',
  checkInTime: undefined,
  checkOutTime: undefined,
})

// Mock data for development
const mockStudents: ClassStudent[] = [
  { id: 1, nis: '2024001', nisn: '0012345678', name: 'Ahmad Fauzi', isActive: true },
  { id: 2, nis: '2024002', nisn: '0012345679', name: 'Budi Santoso', isActive: true },
  { id: 3, nis: '2024003', nisn: '0012345680', name: 'Citra Dewi', isActive: true },
  { id: 4, nis: '2024004', nisn: '0012345681', name: 'Dian Pratama', isActive: true },
  { id: 5, nis: '2024005', nisn: '0012345682', name: 'Eka Putri', isActive: true },
  { id: 6, nis: '2024006', nisn: '0012345683', name: 'Fajar Nugroho', isActive: true },
  { id: 7, nis: '2024007', nisn: '0012345684', name: 'Galih Pratama', isActive: true },
  { id: 8, nis: '2024008', nisn: '0012345685', name: 'Hana Safitri', isActive: true },
]

const mockAttendance: StudentAttendance[] = [
  { id: 1, studentId: 1, studentName: 'Ahmad Fauzi', studentNis: '2024001', date: dayjs().format('YYYY-MM-DD'), checkInTime: '07:15', checkOutTime: '14:00', status: 'present', method: 'rfid', createdAt: '', updatedAt: '' },
  { id: 2, studentId: 2, studentName: 'Budi Santoso', studentNis: '2024002', date: dayjs().format('YYYY-MM-DD'), checkInTime: '07:35', checkOutTime: '14:00', status: 'late', method: 'rfid', createdAt: '', updatedAt: '' },
  { id: 3, studentId: 3, studentName: 'Citra Dewi', studentNis: '2024003', date: dayjs().format('YYYY-MM-DD'), status: 'absent', method: 'manual', createdAt: '', updatedAt: '' },
  { id: 4, studentId: 4, studentName: 'Dian Pratama', studentNis: '2024004', date: dayjs().format('YYYY-MM-DD'), checkInTime: '07:10', checkOutTime: '14:00', status: 'present', method: 'rfid', createdAt: '', updatedAt: '' },
  { id: 5, studentId: 5, studentName: 'Eka Putri', studentNis: '2024005', date: dayjs().format('YYYY-MM-DD'), checkInTime: '07:20', checkOutTime: '14:00', status: 'present', method: 'rfid', createdAt: '', updatedAt: '' },
  { id: 6, studentId: 6, studentName: 'Fajar Nugroho', studentNis: '2024006', date: dayjs().format('YYYY-MM-DD'), status: 'excused', method: 'manual', createdAt: '', updatedAt: '' },
  { id: 7, studentId: 7, studentName: 'Galih Pratama', studentNis: '2024007', date: dayjs().format('YYYY-MM-DD'), checkInTime: '07:05', checkOutTime: '14:00', status: 'present', method: 'rfid', createdAt: '', updatedAt: '' },
  { id: 8, studentId: 8, studentName: 'Hana Safitri', studentNis: '2024008', date: dayjs().format('YYYY-MM-DD'), checkInTime: '07:25', checkOutTime: '14:00', status: 'present', method: 'rfid', createdAt: '', updatedAt: '' },
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
    title: 'Status',
    dataIndex: 'status',
    key: 'status',
    width: 120,
    align: 'center',
  },
  {
    title: 'Jam Masuk',
    dataIndex: 'checkInTime',
    key: 'checkInTime',
    width: 120,
    align: 'center',
  },
  {
    title: 'Jam Pulang',
    dataIndex: 'checkOutTime',
    key: 'checkOutTime',
    width: 120,
    align: 'center',
  },
  {
    title: 'Metode',
    dataIndex: 'method',
    key: 'method',
    width: 100,
    align: 'center',
  },
  {
    title: 'Aksi',
    key: 'action',
    width: 100,
    align: 'center',
  },
]

// Computed summary statistics
const summaryStats = computed(() => {
  const data = attendanceData.value
  const total = data.length
  const present = data.filter(item => item.status === 'present').length
  const late = data.filter(item => item.status === 'late').length
  const absent = data.filter(item => item.status === 'absent').length
  const excused = data.filter(item => item.status === 'excused').length
  
  return {
    total,
    present,
    late,
    absent,
    excused,
    percentage: total > 0 ? Math.round(((present + late) / total) * 100) : 0,
  }
})

// Get status tag color
const getStatusColor = (status: string): string => {
  switch (status) {
    case 'present': return 'success'
    case 'late': return 'warning'
    case 'absent': return 'error'
    case 'excused': return 'blue'
    default: return 'default'
  }
}

// Get status label
const getStatusLabel = (status: string): string => {
  switch (status) {
    case 'present': return 'Hadir'
    case 'late': return 'Terlambat'
    case 'absent': return 'Tidak Hadir'
    case 'excused': return 'Izin'
    default: return status
  }
}

// Get method label
const getMethodLabel = (method: string): string => {
  return method === 'rfid' ? 'RFID' : 'Manual'
}

// Load attendance data
const loadAttendance = async () => {
  loading.value = true
  try {
    const response = await homeroomService.getClassAttendance(selectedDate.value.format('YYYY-MM-DD'))
    attendanceData.value = response.data
  } catch {
    // Use mock data on error
    attendanceData.value = mockAttendance
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

// Handle date change
const handleDateChange = (date: string | Dayjs) => {
  if (date && typeof date !== 'string') {
    selectedDate.value = date
    loadAttendance()
  }
}

// Open modal for manual attendance
const openManualAttendanceModal = (record?: StudentAttendance) => {
  editingRecord.value = record || null
  
  if (record) {
    formState.value = {
      studentId: record.studentId,
      status: record.status,
      checkInTime: record.checkInTime ? dayjs(record.checkInTime, 'HH:mm') : undefined,
      checkOutTime: record.checkOutTime ? dayjs(record.checkOutTime, 'HH:mm') : undefined,
    }
  } else {
    formState.value = {
      studentId: undefined,
      status: 'present',
      checkInTime: dayjs('07:00', 'HH:mm'),
      checkOutTime: undefined,
    }
  }
  
  modalVisible.value = true
}

// Close modal
const closeModal = () => {
  modalVisible.value = false
  editingRecord.value = null
  formRef.value?.resetFields()
}

// Filter option for student select
const filterStudentOption = (input: string, option: { label?: string }) => {
  return option.label?.toLowerCase().includes(input.toLowerCase()) ?? false
}

// Submit manual attendance
const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    modalLoading.value = true

    const data: ManualAttendanceRequest = {
      studentId: formState.value.studentId!,
      date: selectedDate.value.format('YYYY-MM-DD'),
      status: formState.value.status,
      checkInTime: formState.value.checkInTime?.format('HH:mm'),
      checkOutTime: formState.value.checkOutTime?.format('HH:mm'),
    }

    if (editingRecord.value) {
      await homeroomService.updateAttendance(editingRecord.value.id, data)
      message.success('Absensi berhasil diperbarui')
    } else {
      await homeroomService.recordManualAttendance(data)
      message.success('Absensi manual berhasil dicatat')
    }

    closeModal()
    loadAttendance()
  } catch (err: unknown) {
    if (err && typeof err === 'object' && 'errorFields' in err) {
      // Form validation error
      return
    }
    message.error('Gagal menyimpan absensi')
  } finally {
    modalLoading.value = false
  }
}

// Get students not yet recorded
const unrecordedStudents = computed(() => {
  const recordedIds = attendanceData.value.map(a => a.studentId)
  return students.value.filter(s => !recordedIds.includes(s.id))
})

// Format date for display
const formattedDate = computed(() => {
  return selectedDate.value.format('dddd, D MMMM YYYY')
})

onMounted(() => {
  loadAttendance()
  loadStudents()
})
</script>

<template>
  <div class="class-attendance">
    <div class="page-header">
      <div>
        <Title :level="2" style="margin: 0">Absensi Kelas</Title>
        <Text type="secondary">Kelas {{ className }}</Text>
      </div>
      <Text type="secondary">
        <CalendarOutlined /> {{ formattedDate }}
      </Text>
    </div>

    <!-- Summary Cards -->
    <Row :gutter="[24, 24]" class="summary-row">
      <Col :xs="24" :sm="12" :lg="6">
        <Card class="stat-card">
          <Statistic
            title="Total Siswa"
            :value="summaryStats.total"
            :value-style="{ color: '#3b82f6' }"
          />
        </Card>
      </Col>
      <Col :xs="24" :sm="12" :lg="6">
        <Card class="stat-card">
          <Statistic
            title="Hadir"
            :value="summaryStats.present"
            :value-style="{ color: '#22c55e' }"
          >
            <template #prefix>
              <CheckCircleOutlined />
            </template>
          </Statistic>
        </Card>
      </Col>
      <Col :xs="24" :sm="12" :lg="6">
        <Card class="stat-card">
          <Statistic
            title="Terlambat"
            :value="summaryStats.late"
            :value-style="{ color: '#f97316' }"
          >
            <template #prefix>
              <ClockCircleOutlined />
            </template>
          </Statistic>
        </Card>
      </Col>
      <Col :xs="24" :sm="12" :lg="6">
        <Card class="stat-card">
          <Statistic
            title="Tidak Hadir"
            :value="summaryStats.absent"
            :value-style="{ color: '#ef4444' }"
          >
            <template #prefix>
              <CloseCircleOutlined />
            </template>
          </Statistic>
        </Card>
      </Col>
    </Row>

    <Card style="margin-top: 24px">
      <!-- Toolbar -->
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col :xs="24" :sm="24" :md="16">
          <Space wrap>
            <DatePicker
              :value="selectedDate"
              format="DD MMMM YYYY"
              :allow-clear="false"
              style="width: 200px"
              @change="handleDateChange"
            />
          </Space>
        </Col>
        <Col :xs="24" :sm="24" :md="8" class="toolbar-right">
          <Space>
            <Button @click="loadAttendance">
              <template #icon><ReloadOutlined /></template>
            </Button>
            <Button type="primary" @click="openManualAttendanceModal()">
              <template #icon><PlusOutlined /></template>
              Input Manual
            </Button>
          </Space>
        </Col>
      </Row>

      <!-- Table -->
      <Table
        :columns="columns"
        :data-source="attendanceData"
        :loading="loading"
        :pagination="false"
        row-key="id"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <Tag :color="getStatusColor((record as StudentAttendance).status)">
              {{ getStatusLabel((record as StudentAttendance).status) }}
            </Tag>
          </template>
          <template v-else-if="column.key === 'checkInTime'">
            <span v-if="(record as StudentAttendance).checkInTime">
              {{ (record as StudentAttendance).checkInTime }}
            </span>
            <Text v-else type="secondary">-</Text>
          </template>
          <template v-else-if="column.key === 'checkOutTime'">
            <span v-if="(record as StudentAttendance).checkOutTime">
              {{ (record as StudentAttendance).checkOutTime }}
            </span>
            <Text v-else type="secondary">-</Text>
          </template>
          <template v-else-if="column.key === 'method'">
            <Tag :color="(record as StudentAttendance).method === 'rfid' ? 'blue' : 'orange'">
              {{ getMethodLabel((record as StudentAttendance).method) }}
            </Tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <Button type="link" size="small" @click="openManualAttendanceModal(record as StudentAttendance)">
              <template #icon><EditOutlined /></template>
              Edit
            </Button>
          </template>
        </template>
        
        <!-- Summary Footer -->
        <template #summary>
          <Table.Summary fixed>
            <Table.Summary.Row class="summary-row-table">
              <Table.Summary.Cell :index="0" :col-span="2">
                <Text strong>Persentase Kehadiran</Text>
              </Table.Summary.Cell>
              <Table.Summary.Cell :index="2" :col-span="5">
                <Progress
                  :percent="summaryStats.percentage"
                  :stroke-color="summaryStats.percentage >= 90 ? '#22c55e' : summaryStats.percentage >= 75 ? '#f97316' : '#ef4444'"
                  :show-info="true"
                  size="small"
                  style="max-width: 300px"
                />
              </Table.Summary.Cell>
            </Table.Summary.Row>
          </Table.Summary>
        </template>
      </Table>
    </Card>

    <!-- Manual Attendance Modal -->
    <Modal
      v-model:open="modalVisible"
      :title="editingRecord ? 'Edit Absensi' : 'Input Absensi Manual'"
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
            :disabled="!!editingRecord"
            show-search
            :filter-option="filterStudentOption"
          >
            <SelectOption
              v-for="student in (editingRecord ? students : unrecordedStudents)"
              :key="student.id"
              :value="student.id"
              :label="student.name"
            >
              {{ student.nis }} - {{ student.name }}
            </SelectOption>
          </Select>
        </FormItem>

        <FormItem
          label="Status"
          name="status"
          :rules="[{ required: true, message: 'Pilih status' }]"
        >
          <Select v-model:value="formState.status" placeholder="Pilih status">
            <SelectOption value="present">Hadir</SelectOption>
            <SelectOption value="late">Terlambat</SelectOption>
            <SelectOption value="absent">Tidak Hadir</SelectOption>
            <SelectOption value="excused">Izin</SelectOption>
          </Select>
        </FormItem>

        <Row :gutter="16">
          <Col :span="12">
            <FormItem label="Jam Masuk" name="checkInTime">
              <TimePicker
                v-model:value="formState.checkInTime"
                format="HH:mm"
                placeholder="Jam masuk"
                style="width: 100%"
              />
            </FormItem>
          </Col>
          <Col :span="12">
            <FormItem label="Jam Pulang" name="checkOutTime">
              <TimePicker
                v-model:value="formState.checkOutTime"
                format="HH:mm"
                placeholder="Jam pulang"
                style="width: 100%"
              />
            </FormItem>
          </Col>
        </Row>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.class-attendance {
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

.summary-row {
  margin-bottom: 0;
}

.stat-card {
  height: 100%;
}

.stat-card :deep(.ant-statistic-title) {
  font-size: 14px;
  color: #8c8c8c;
}

.stat-card :deep(.ant-statistic-content-prefix) {
  margin-right: 8px;
}

.toolbar {
  margin-bottom: 16px;
}

.toolbar-right {
  display: flex;
  justify-content: flex-end;
}

:deep(.summary-row-table) {
  background: #fafafa;
}

@media (max-width: 768px) {
  .toolbar-right {
    margin-top: 16px;
    justify-content: flex-start;
  }
}
</style>
