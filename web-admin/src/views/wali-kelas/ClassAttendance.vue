<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import {
  Table,
  Button,
  Space,
  Card,
  Row,
  Col,
  Typography,
  DatePicker,
  Progress, Statistic, Modal, Form, FormItem, Select, SelectOption,
  TimePicker, message, Alert,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  ReloadOutlined, CalendarOutlined, CheckCircleOutlined,
  CloseCircleOutlined, ClockCircleOutlined, EditOutlined, PlusOutlined,
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import type { Dayjs } from 'dayjs'
import { homeroomService } from '@/services'
import { useClassInfo, useClassStudents } from '@/composables/useWaliKelas'
import type { StudentAttendance, ManualAttendanceRequest } from '@/types/homeroom'
import type { AttendanceSchedule } from '@/types/schedule'

const { Title, Text } = Typography

// Composables
const { className, loadClassInfo } = useClassInfo()
const { students, loadingStudents, loadStudents, filterStudentOption } = useClassStudents()

// Mounted state
const isMounted = ref(true)

// State
const loading = ref(false)
const loadingSchedules = ref(false)
const error = ref<string | null>(null)
const attendanceData = ref<StudentAttendance[]>([])
const schedules = ref<AttendanceSchedule[]>([])
const selectedDate = ref<Dayjs>(dayjs())

// Modal state
const modalVisible = ref(false)
const modalLoading = ref(false)
const editingRecord = ref<StudentAttendance | null>(null)
const formRef = ref()

// Form state
const formState = ref<{
  studentId: number | undefined
  scheduleId: number | undefined
  status: 'on_time' | 'late' | 'sick' | 'excused'
  checkInTime: Dayjs | undefined
}>({
  studentId: undefined,
  scheduleId: undefined,
  status: 'on_time',
  checkInTime: undefined,
})

// Table columns
const columns: TableProps['columns'] = [
  { title: 'NIS', dataIndex: 'studentNis', key: 'studentNis', width: 100 },
  { title: 'Nama Siswa', dataIndex: 'studentName', key: 'studentName' },
  { title: 'Status', dataIndex: 'status', key: 'status', width: 120, align: 'center' },
  { title: 'Jam Masuk', dataIndex: 'checkInTime', key: 'checkInTime', width: 120, align: 'center' },
  { title: 'Jam Pulang', dataIndex: 'checkOutTime', key: 'checkOutTime', width: 120, align: 'center' },
  { title: 'Metode', dataIndex: 'method', key: 'method', width: 100, align: 'center' },
  { title: 'Aksi', key: 'action', width: 100, align: 'center' },
]

// Computed summary statistics
const summaryStats = computed(() => {
  const data = attendanceData.value
  const total = data.length
  const present = data.filter(item => item.status === 'present' || item.status === 'on_time').length
  const late = data.filter(item => item.status === 'late' || item.status === 'very_late').length
  const absent = data.filter(item => item.status === 'absent').length
  const excused = data.filter(item => item.status === 'excused').length
  const sick = data.filter(item => item.status === 'sick').length
  
  return {
    total, present, late, absent, excused, sick,
    percentage: total > 0 ? Math.round(((present + late) / total) * 100) : 0,
  }
})

// Status helpers
const getStatusColor = (status: string): string => {
  const colors: Record<string, string> = {
    present: 'success', on_time: 'success', late: 'warning',
    very_late: 'orange', absent: 'error', excused: 'blue', sick: 'purple',
  }
  return colors[status] || 'default'
}

const getStatusLabel = (status: string): string => {
  const labels: Record<string, string> = {
    present: 'Hadir', on_time: 'Hadir', late: 'Terlambat',
    very_late: 'Sangat Terlambat', absent: 'Tidak Hadir', excused: 'Izin', sick: 'Sakit',
  }
  return labels[status] || status
}

const getMethodLabel = (method: string): string => method === 'rfid' ? 'RFID' : 'Manual'

// Load attendance data
const loadAttendance = async () => {
  if (!isMounted.value) return
  loading.value = true
  error.value = null
  try {
    const response = await homeroomService.getClassAttendance(selectedDate.value.format('YYYY-MM-DD'))
    if (isMounted.value) attendanceData.value = response.data || []
  } catch (err) {
    console.error('Failed to load attendance:', err)
    if (isMounted.value) {
      attendanceData.value = []
      error.value = 'Gagal memuat data absensi'
    }
  } finally {
    if (isMounted.value) loading.value = false
  }
}

// Load schedules
const loadSchedules = async () => {
  loadingSchedules.value = true
  try {
    const response = await homeroomService.getActiveSchedules(selectedDate.value.format('YYYY-MM-DD'))
    schedules.value = response || []
    const defaultSchedule = schedules.value.find(s => s.isDefault)
    if (defaultSchedule && !formState.value.scheduleId) {
      formState.value.scheduleId = defaultSchedule.id
    }
  } catch (err) {
    console.error('Failed to load schedules:', err)
    schedules.value = []
  } finally {
    loadingSchedules.value = false
  }
}

// Handle date change
const handleDateChange = (date: string | Dayjs) => {
  if (date && typeof date !== 'string') {
    selectedDate.value = date
    loadAttendance()
    schedules.value = []
  }
}

// Open modal
const openManualAttendanceModal = async (record?: StudentAttendance) => {
  editingRecord.value = record || null
  if (students.value.length === 0) await loadStudents()
  await loadSchedules()
  
  if (record) {
    let status: 'on_time' | 'late' | 'sick' | 'excused' = 'on_time'
    if (record.status === 'sick') status = 'sick'
    else if (record.status === 'excused') status = 'excused'
    else if (record.status === 'late' || record.status === 'very_late') status = 'late'
    
    formState.value = {
      studentId: record.studentId,
      scheduleId: schedules.value.length > 0 ? (schedules.value.find(s => s.isDefault)?.id || schedules.value[0].id) : undefined,
      status,
      checkInTime: record.checkInTime ? dayjs(record.checkInTime, 'HH:mm') : undefined,
    }
  } else {
    const defaultSchedule = schedules.value.find(s => s.isDefault) || schedules.value[0]
    // Set waktu default berdasarkan startTime jadwal yang dipilih
    const defaultTime = defaultSchedule?.startTime || '07:00'
    formState.value = {
      studentId: undefined,
      scheduleId: defaultSchedule?.id,
      status: 'on_time',
      checkInTime: dayjs(defaultTime, 'HH:mm'),
    }
  }
  modalVisible.value = true
}

// Update waktu default ketika jadwal berubah
watch(() => formState.value.scheduleId, (newScheduleId) => {
  if (newScheduleId && !editingRecord.value) {
    const selectedSchedule = schedules.value.find(s => s.id === newScheduleId)
    if (selectedSchedule) {
      formState.value.checkInTime = dayjs(selectedSchedule.startTime, 'HH:mm')
    }
  }
})

const closeModal = () => {
  modalVisible.value = false
  editingRecord.value = null
  formRef.value?.resetFields()
}

// Submit
const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    modalLoading.value = true

    const data: ManualAttendanceRequest = {
      studentId: formState.value.studentId!,
      scheduleId: formState.value.scheduleId!,
      date: selectedDate.value.format('YYYY-MM-DD'),
      status: formState.value.status,
      checkInTime: requiresTimeInput.value ? formState.value.checkInTime?.format('HH:mm') : undefined,
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
    if (err && typeof err === 'object' && 'errorFields' in err) return
    message.error('Gagal menyimpan absensi')
  } finally {
    modalLoading.value = false
  }
}

// Available students for input
const availableStudentsForInput = computed(() => {
  const recordedIds = attendanceData.value.filter(a => a.id > 0).map(a => a.studentId)
  const unrecorded = students.value.filter(s => !recordedIds.includes(s.id))
  return unrecorded.length > 0 ? unrecorded : students.value
})

const requiresTimeInput = computed(() => ['on_time', 'late'].includes(formState.value.status))

watch(() => formState.value.status, (newStatus) => {
  if (['sick', 'excused'].includes(newStatus)) formState.value.checkInTime = undefined
})

const formattedDate = computed(() => selectedDate.value.format('dddd, D MMMM YYYY'))

onMounted(() => {
  loadClassInfo()
  loadAttendance()
  loadStudents()
})

onUnmounted(() => { isMounted.value = false })
</script>

<template>
  <div class="wali-kelas-page">
    <div class="page-header">
      <div>
        <Title :level="2" style="margin: 0">Absensi Kelas</Title>
        <Text type="secondary">Kelas {{ className }}</Text>
      </div>
      <Text type="secondary"><CalendarOutlined /> {{ formattedDate }}</Text>
    </div>

    <Alert v-if="error" type="error" :message="error" show-icon closable style="margin-bottom: 16px" @close="error = null" />

    <Row :gutter="[16, 16]" class="summary-row">
      <Col :xs="12" :sm="8" :lg="4">
        <Card class="stat-card"><Statistic title="Total Siswa" :value="summaryStats.total" :value-style="{ color: '#3b82f6' }" /></Card>
      </Col>
      <Col :xs="12" :sm="8" :lg="4">
        <Card class="stat-card"><Statistic title="Hadir" :value="summaryStats.present" :value-style="{ color: '#22c55e' }"><template #prefix><CheckCircleOutlined /></template></Statistic></Card>
      </Col>
      <Col :xs="12" :sm="8" :lg="4">
        <Card class="stat-card"><Statistic title="Terlambat" :value="summaryStats.late" :value-style="{ color: '#f97316' }"><template #prefix><ClockCircleOutlined /></template></Statistic></Card>
      </Col>
      <Col :xs="12" :sm="8" :lg="4">
        <Card class="stat-card"><Statistic title="Sakit" :value="summaryStats.sick" :value-style="{ color: '#a855f7' }" /></Card>
      </Col>
      <Col :xs="12" :sm="8" :lg="4">
        <Card class="stat-card"><Statistic title="Izin" :value="summaryStats.excused" :value-style="{ color: '#3b82f6' }" /></Card>
      </Col>
      <Col :xs="12" :sm="8" :lg="4">
        <Card class="stat-card"><Statistic title="Tidak Hadir" :value="summaryStats.absent" :value-style="{ color: '#ef4444' }"><template #prefix><CloseCircleOutlined /></template></Statistic></Card>
      </Col>
    </Row>

    <Card style="margin-top: 24px">
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col :xs="24" :sm="24" :md="16">
          <DatePicker :value="selectedDate" format="DD MMMM YYYY" :allow-clear="false" style="width: 200px" @change="handleDateChange" />
        </Col>
        <Col :xs="24" :sm="24" :md="8" class="toolbar-right">
          <Space>
            <Button @click="loadAttendance"><template #icon><ReloadOutlined /></template></Button>
            <Button type="primary" @click="openManualAttendanceModal()"><template #icon><PlusOutlined /></template>Input Manual</Button>
          </Space>
        </Col>
      </Row>

      <Table :columns="columns" :data-source="attendanceData" :loading="loading" :pagination="false" row-key="id">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <span :class="['status-text', getStatusColor((record as StudentAttendance).status)]">{{ getStatusLabel((record as StudentAttendance).status) }}</span>
          </template>
          <template v-else-if="column.key === 'checkInTime'">
            <span v-if="(record as StudentAttendance).checkInTime">{{ (record as StudentAttendance).checkInTime }}</span>
            <Text v-else type="secondary">-</Text>
          </template>
          <template v-else-if="column.key === 'checkOutTime'">
            <span v-if="(record as StudentAttendance).checkOutTime">{{ (record as StudentAttendance).checkOutTime }}</span>
            <Text v-else type="secondary">-</Text>
          </template>
          <template v-else-if="column.key === 'method'">
            <span :class="['method-text', (record as StudentAttendance).method === 'rfid' ? 'blue' : 'orange']">{{ getMethodLabel((record as StudentAttendance).method) }}</span>
          </template>
          <template v-else-if="column.key === 'action'">
            <Button type="link" @click="openManualAttendanceModal(record as StudentAttendance)"><template #icon><EditOutlined /></template>Edit</Button>
          </template>
        </template>
        <template #summary>
          <Table.Summary fixed>
            <Table.Summary.Row class="summary-row-table">
              <Table.Summary.Cell :index="0" :col-span="2"><Text strong>Persentase Kehadiran</Text></Table.Summary.Cell>
              <Table.Summary.Cell :index="2" :col-span="5">
                <Progress :percent="summaryStats.percentage" :stroke-color="summaryStats.percentage >= 90 ? '#22c55e' : summaryStats.percentage >= 75 ? '#f97316' : '#ef4444'" :show-info="true" style="max-width: 300px" />
              </Table.Summary.Cell>
            </Table.Summary.Row>
          </Table.Summary>
        </template>
      </Table>
    </Card>

    <Modal v-model:open="modalVisible" :title="editingRecord ? 'Edit Absensi' : 'Input Absensi Manual'" :confirm-loading="modalLoading" @ok="handleSubmit" @cancel="closeModal">
      <Alert v-if="schedules.length === 0 && !loadingSchedules" type="warning" message="Tidak ada jadwal aktif" description="Tidak ada jadwal absensi yang aktif untuk tanggal ini." show-icon style="margin-bottom: 16px" />
      <Form ref="formRef" :model="formState" layout="vertical" style="margin-top: 16px">
        <FormItem label="Jadwal Absensi" name="scheduleId" :rules="[{ required: true, message: 'Pilih jadwal absensi' }]">
          <Select v-model:value="formState.scheduleId" placeholder="Pilih jadwal" :loading="loadingSchedules" :disabled="schedules.length === 0">
            <SelectOption v-for="schedule in schedules" :key="schedule.id" :value="schedule.id">
              {{ schedule.name }} ({{ schedule.startTime }} - {{ schedule.endTime }})
              <Tag v-if="schedule.isDefault" color="blue" style="margin-left: 8px">Default</Tag>
            </SelectOption>
          </Select>
        </FormItem>
        <FormItem label="Siswa" name="studentId" :rules="[{ required: true, message: 'Pilih siswa' }]">
          <Select v-model:value="formState.studentId" placeholder="Pilih siswa" :disabled="!!editingRecord" :loading="loadingStudents" show-search :filter-option="filterStudentOption" :not-found-content="loadingStudents ? 'Memuat...' : 'Tidak ada siswa'">
            <SelectOption v-for="student in (editingRecord ? students : availableStudentsForInput)" :key="student.id" :value="student.id" :label="student.name">{{ student.nis }} - {{ student.name }}</SelectOption>
          </Select>
        </FormItem>
        <FormItem label="Status" name="status" :rules="[{ required: true, message: 'Pilih status' }]">
          <Select v-model:value="formState.status" placeholder="Pilih status">
            <SelectOption value="on_time">Hadir (Tepat Waktu)</SelectOption>
            <SelectOption value="late">Terlambat</SelectOption>
            <SelectOption value="sick">Sakit</SelectOption>
            <SelectOption value="excused">Izin</SelectOption>
          </Select>
        </FormItem>
        <FormItem v-if="requiresTimeInput" label="Jam Masuk" name="checkInTime">
          <TimePicker v-model:value="formState.checkInTime" format="HH:mm" placeholder="Jam masuk" style="width: 100%" />
        </FormItem>
        <div v-else class="info-box">
          <Text type="secondary">
            <template v-if="formState.status === 'sick'">Siswa tidak masuk karena sakit. Jam masuk tidak diperlukan.</template>
            <template v-else-if="formState.status === 'excused'">Siswa tidak masuk karena izin. Jam masuk tidak diperlukan.</template>
          </Text>
        </div>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.wali-kelas-page { padding: 0; }
.page-header { margin-bottom: 24px; display: flex; justify-content: space-between; align-items: flex-start; flex-wrap: wrap; gap: 8px; }
.summary-row { margin-bottom: 0; }
.stat-card { height: 100%; }
.stat-card :deep(.ant-statistic-title) { font-size: 14px; color: #8c8c8c; }
.stat-card :deep(.ant-statistic-content-prefix) { margin-right: 8px; }
.toolbar { margin-bottom: 16px; }
.toolbar-right { display: flex; justify-content: flex-end; }
:deep(.summary-row-table) { background: #fafafa; }
.info-box { padding: 12px 16px; background: #f5f5f5; border-radius: 6px; margin-top: 8px; }
.status-text { font-size: 13px; font-weight: 500; }
.status-text.success { color: #22c55e; }
.status-text.warning { color: #f59e0b; }
.status-text.orange { color: #f97316; }
.status-text.error { color: #ef4444; }
.status-text.purple { color: #a855f7; }
.status-text.blue { color: #3b82f6; }
.method-text { font-size: 12px; font-weight: 500; }
.method-text.blue { color: #3b82f6; }
.method-text.orange { color: #f97316; }
@media (max-width: 768px) { .toolbar-right { margin-top: 16px; justify-content: flex-start; } }
</style>
