<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import {
  Table,
  Button,
  Space,
  Tag,
  Modal,
  Form,
  FormItem,
  Input,
  InputNumber,
  Switch,
  Checkbox,
  CheckboxGroup,
  TimePicker,
  message,
  Popconfirm,
  Card,
  Row,
  Col,
  Typography,
  Tooltip,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
  StarOutlined,
  StarFilled,
  ClockCircleOutlined,
} from '@ant-design/icons-vue'
import { scheduleService } from '@/services'
import type { AttendanceSchedule, CreateScheduleRequest, UpdateScheduleRequest } from '@/types/schedule'
import { DAYS_OF_WEEK, parseDaysOfWeek, formatDaysOfWeek, getDayLabels } from '@/types/schedule'

const { Title, Text } = Typography

// Table state
const loading = ref(false)
const schedules = ref<AttendanceSchedule[]>([])

// Modal state
const modalVisible = ref(false)
const modalLoading = ref(false)
const isEditing = ref(false)
const editingSchedule = ref<AttendanceSchedule | null>(null)

// Form state
const formRef = ref()
const formState = reactive({
  name: '',
  startTime: '' as string,
  endTime: '' as string,
  lateThreshold: 15,
  veryLateThreshold: undefined as number | undefined,
  daysOfWeek: ['1', '2', '3', '4', '5'] as string[],
  isActive: true,
})

// Form rules
const formRules = {
  name: [{ required: true, message: 'Nama jadwal wajib diisi' }],
  startTime: [{ 
    required: true, 
    message: 'Waktu mulai wajib diisi',
    validator: (_rule: unknown, value: string) => {
      if (!value || value.trim() === '') {
        return Promise.reject('Waktu mulai wajib diisi')
      }
      return Promise.resolve()
    }
  }],
  endTime: [{ 
    required: true, 
    message: 'Waktu selesai wajib diisi',
    validator: (_rule: unknown, value: string) => {
      if (!value || value.trim() === '') {
        return Promise.reject('Waktu selesai wajib diisi')
      }
      return Promise.resolve()
    }
  }],
  lateThreshold: [{ required: true, message: 'Batas terlambat wajib diisi' }],
  daysOfWeek: [{ required: true, message: 'Pilih minimal satu hari', type: 'array' as const, min: 1 }],
}

// Table columns
const columns: TableProps['columns'] = [
  {
    title: 'Nama Jadwal',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: 'Waktu',
    key: 'time',
    width: 150,
  },
  {
    title: 'Batas Terlambat',
    key: 'threshold',
    width: 180,
  },
  {
    title: 'Hari Aktif',
    key: 'days',
    width: 200,
  },
  {
    title: 'Status',
    key: 'status',
    width: 120,
    align: 'center',
  },
  {
    title: 'Aksi',
    key: 'action',
    width: 200,
    align: 'center',
  },
]

// Computed: check if max schedules reached (10 per school)
const maxSchedulesReached = computed(() => schedules.value.length >= 10)

// Load schedules data
const loadSchedules = async () => {
  loading.value = true
  try {
    const response = await scheduleService.getSchedules()
    schedules.value = response.schedules
  } catch (err) {
    console.error('Failed to load schedules:', err)
    message.error('Gagal memuat data jadwal absensi')
    schedules.value = []
  } finally {
    loading.value = false
  }
}

// Open modal for create
const openCreateModal = () => {
  if (maxSchedulesReached.value) {
    message.warning('Maksimal 10 jadwal per sekolah')
    return
  }
  isEditing.value = false
  editingSchedule.value = null
  resetForm()
  modalVisible.value = true
}

// Open modal for edit
const openEditModal = (schedule: AttendanceSchedule) => {
  isEditing.value = true
  editingSchedule.value = schedule
  formState.name = schedule.name
  formState.startTime = schedule.startTime
  formState.endTime = schedule.endTime
  formState.lateThreshold = schedule.lateThreshold
  formState.veryLateThreshold = schedule.veryLateThreshold || undefined
  formState.daysOfWeek = parseDaysOfWeek(schedule.daysOfWeek)
  formState.isActive = schedule.isActive
  modalVisible.value = true
}

// Reset form
const resetForm = () => {
  formState.name = ''
  formState.startTime = ''
  formState.endTime = ''
  formState.lateThreshold = 15
  formState.veryLateThreshold = undefined
  formState.daysOfWeek = ['1', '2', '3', '4', '5']
  formState.isActive = true
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
    if (isEditing.value && editingSchedule.value) {
      const updateData: UpdateScheduleRequest = {
        name: formState.name,
        start_time: formState.startTime,
        end_time: formState.endTime,
        late_threshold: formState.lateThreshold,
        very_late_threshold: formState.veryLateThreshold || undefined,
        days_of_week: formatDaysOfWeek(formState.daysOfWeek),
        is_active: formState.isActive,
      }
      await scheduleService.updateSchedule(editingSchedule.value.id, updateData)
      message.success('Jadwal berhasil diperbarui')
    } else {
      const createData: CreateScheduleRequest = {
        name: formState.name,
        start_time: formState.startTime,
        end_time: formState.endTime,
        late_threshold: formState.lateThreshold,
        very_late_threshold: formState.veryLateThreshold || undefined,
        days_of_week: formatDaysOfWeek(formState.daysOfWeek),
        is_active: formState.isActive,
      }
      await scheduleService.createSchedule(createData)
      message.success('Jadwal berhasil ditambahkan')
    }
    modalVisible.value = false
    resetForm()
    loadSchedules()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Terjadi kesalahan')
  } finally {
    modalLoading.value = false
  }
}

// Handle delete schedule
const handleDelete = async (schedule: AttendanceSchedule) => {
  try {
    await scheduleService.deleteSchedule(schedule.id)
    message.success(`Jadwal "${schedule.name}" berhasil dihapus`)
    loadSchedules()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Gagal menghapus jadwal')
  }
}

// Handle set default schedule
const handleSetDefault = async (schedule: AttendanceSchedule) => {
  try {
    await scheduleService.setDefaultSchedule(schedule.id)
    message.success(`Jadwal "${schedule.name}" ditetapkan sebagai default`)
    loadSchedules()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Gagal menetapkan jadwal default')
  }
}

onMounted(() => {
  loadSchedules()
})
</script>

<template>
  <div class="schedule-management">
    <div class="page-header">
      <Title :level="2" style="margin: 0">Pengaturan Jadwal Absensi</Title>
      <Text type="secondary">Kelola jadwal waktu absensi untuk berbagai kegiatan (masuk, pulang, sholat, dll)</Text>
    </div>

    <Card>
      <!-- Toolbar -->
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col>
          <Text type="secondary">
            <ClockCircleOutlined /> {{ schedules.length }}/10 jadwal
          </Text>
        </Col>
        <Col class="toolbar-right">
          <Space>
            <Button @click="loadSchedules">
              <template #icon><ReloadOutlined /></template>
            </Button>
            <Tooltip v-if="maxSchedulesReached" title="Maksimal 10 jadwal per sekolah">
              <Button type="primary" disabled>
                <template #icon><PlusOutlined /></template>
                Tambah Jadwal
              </Button>
            </Tooltip>
            <Button v-else type="primary" @click="openCreateModal">
              <template #icon><PlusOutlined /></template>
              Tambah Jadwal
            </Button>
          </Space>
        </Col>
      </Row>

      <!-- Table -->
      <Table
        :columns="columns"
        :data-source="schedules"
        :loading="loading"
        :pagination="false"
        row-key="id"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <Space>
              <span>{{ (record as AttendanceSchedule).name }}</span>
              <Tag v-if="(record as AttendanceSchedule).isDefault" color="gold">
                <StarFilled /> Default
              </Tag>
            </Space>
          </template>
          <template v-else-if="column.key === 'time'">
            <Tag color="blue">
              {{ (record as AttendanceSchedule).startTime }} - {{ (record as AttendanceSchedule).endTime }}
            </Tag>
          </template>
          <template v-else-if="column.key === 'threshold'">
            <Space direction="vertical" size="small">
              <Text>Terlambat: {{ (record as AttendanceSchedule).lateThreshold }} menit</Text>
              <Text v-if="(record as AttendanceSchedule).veryLateThreshold" type="secondary">
                Sangat terlambat: {{ (record as AttendanceSchedule).veryLateThreshold }} menit
              </Text>
            </Space>
          </template>
          <template v-else-if="column.key === 'days'">
            <Text>{{ getDayLabels((record as AttendanceSchedule).daysOfWeek) }}</Text>
          </template>
          <template v-else-if="column.key === 'status'">
            <Tag :color="(record as AttendanceSchedule).isActive ? 'success' : 'default'">
              {{ (record as AttendanceSchedule).isActive ? 'Aktif' : 'Nonaktif' }}
            </Tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <Space>
              <Tooltip v-if="!(record as AttendanceSchedule).isDefault" title="Jadikan Default">
                <Button size="small" @click="handleSetDefault(record as AttendanceSchedule)">
                  <template #icon><StarOutlined /></template>
                </Button>
              </Tooltip>
              <Button size="small" @click="openEditModal(record as AttendanceSchedule)">
                <template #icon><EditOutlined /></template>
                Edit
              </Button>
              <Popconfirm
                title="Hapus jadwal ini?"
                description="Jadwal yang sudah digunakan tidak dapat dihapus."
                ok-text="Ya, Hapus"
                cancel-text="Batal"
                @confirm="handleDelete(record as AttendanceSchedule)"
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

    <!-- Create/Edit Modal -->
    <Modal
      v-model:open="modalVisible"
      :title="isEditing ? 'Edit Jadwal Absensi' : 'Tambah Jadwal Absensi'"
      :confirm-loading="modalLoading"
      width="600px"
      @ok="handleSubmit"
      @cancel="handleModalCancel"
    >
      <Form
        ref="formRef"
        :model="formState"
        :rules="formRules"
        layout="vertical"
        style="margin-top: 16px"
      >
        <FormItem label="Nama Jadwal" name="name" required>
          <Input v-model:value="formState.name" placeholder="Contoh: Masuk Pagi, Pulang, Sholat Dzuhur" />
        </FormItem>
        
        <Row :gutter="16">
          <Col :span="12">
            <FormItem label="Waktu Mulai" name="startTime" required>
              <TimePicker
                v-model:value="formState.startTime"
                format="HH:mm"
                value-format="HH:mm"
                placeholder="Pilih waktu"
                style="width: 100%"
              />
            </FormItem>
          </Col>
          <Col :span="12">
            <FormItem label="Waktu Selesai" name="endTime" required>
              <TimePicker
                v-model:value="formState.endTime"
                format="HH:mm"
                value-format="HH:mm"
                placeholder="Pilih waktu"
                style="width: 100%"
              />
            </FormItem>
          </Col>
        </Row>

        <Row :gutter="16">
          <Col :span="12">
            <FormItem label="Batas Terlambat (menit)" name="lateThreshold" required>
              <InputNumber
                v-model:value="formState.lateThreshold"
                :min="0"
                :max="120"
                placeholder="15"
                style="width: 100%"
              />
              <Text type="secondary" style="font-size: 12px">
                Menit setelah waktu mulai untuk dianggap terlambat
              </Text>
            </FormItem>
          </Col>
          <Col :span="12">
            <FormItem label="Batas Sangat Terlambat (menit)" name="veryLateThreshold">
              <InputNumber
                v-model:value="formState.veryLateThreshold"
                :min="0"
                :max="180"
                placeholder="Opsional"
                style="width: 100%"
              />
              <Text type="secondary" style="font-size: 12px">
                Opsional. Menit setelah waktu mulai untuk sangat terlambat
              </Text>
            </FormItem>
          </Col>
        </Row>

        <FormItem label="Hari Aktif" name="daysOfWeek" required>
          <CheckboxGroup v-model:value="formState.daysOfWeek">
            <Row>
              <Col v-for="day in DAYS_OF_WEEK" :key="day.value" :span="6">
                <Checkbox :value="day.value">{{ day.label }}</Checkbox>
              </Col>
            </Row>
          </CheckboxGroup>
        </FormItem>

        <FormItem label="Status" name="isActive">
          <Switch v-model:checked="formState.isActive" checked-children="Aktif" un-checked-children="Nonaktif" />
          <Text type="secondary" style="margin-left: 12px; font-size: 12px">
            Jadwal nonaktif tidak akan digunakan untuk absensi baru
          </Text>
        </FormItem>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.schedule-management {
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

@media (max-width: 576px) {
  .toolbar-right {
    margin-top: 16px;
    justify-content: flex-start;
  }
}
</style>
