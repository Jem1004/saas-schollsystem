<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import {
  Card,
  Form,
  FormItem,
  InputNumber,
  Switch,
  Button,
  Row,
  Col,
  Typography,
  Divider,
  TimePicker,
  Select,
  SelectOption,
  message,
  Spin,
  Popconfirm,
} from 'ant-design-vue'
import {
  SaveOutlined,
  ReloadOutlined,
  ClockCircleOutlined,
  BellOutlined,
  CalendarOutlined,
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import type { Dayjs } from 'dayjs'
import { schoolService } from '@/services'
import type { SchoolSettings, UpdateSchoolSettingsRequest } from '@/types/school'

const { Title, Text } = Typography

// State
const loading = ref(false)
const saving = ref(false)
const resetting = ref(false)

// Form state
const formRef = ref()
const formState = reactive<{
  attendanceStartTime: Dayjs | undefined
  attendanceEndTime: Dayjs | undefined
  attendanceLateThreshold: number
  attendanceVeryLateThreshold: number
  enableAttendanceNotification: boolean
  enableGradeNotification: boolean
  enableBKNotification: boolean
  enableHomeroomNotification: boolean
  academicYear: string
  semester: number
}>({
  attendanceStartTime: undefined,
  attendanceEndTime: undefined,
  attendanceLateThreshold: 30,
  attendanceVeryLateThreshold: 60,
  enableAttendanceNotification: true,
  enableGradeNotification: true,
  enableBKNotification: true,
  enableHomeroomNotification: true,
  academicYear: '',
  semester: 1,
})

// Form rules
const formRules = {
  attendanceStartTime: [{ required: true, message: 'Waktu mulai absensi wajib diisi' }],
  attendanceEndTime: [{ required: true, message: 'Waktu akhir absensi wajib diisi' }],
  academicYear: [{ required: true, message: 'Tahun ajaran wajib diisi' }],
}

// Mock data for development
const mockSettings: SchoolSettings = {
  id: 1,
  schoolId: 1,
  attendanceStartTime: '07:00',
  attendanceEndTime: '07:30',
  attendanceLateThreshold: 30,
  attendanceVeryLateThreshold: 60,
  enableAttendanceNotification: true,
  enableGradeNotification: true,
  enableBKNotification: true,
  enableHomeroomNotification: true,
  academicYear: '2024/2025',
  semester: 1,
  createdAt: '2024-01-15T10:00:00Z',
  updatedAt: '2024-01-15T10:00:00Z',
}

// Load settings
const loadSettings = async () => {
  loading.value = true
  try {
    const settings = await schoolService.getSettings()
    applySettings(settings)
  } catch {
    // Use mock data on error
    applySettings(mockSettings)
  } finally {
    loading.value = false
  }
}

// Apply settings to form
const applySettings = (settings: SchoolSettings) => {
  formState.attendanceStartTime = settings.attendanceStartTime 
    ? dayjs(settings.attendanceStartTime, 'HH:mm') 
    : undefined
  formState.attendanceEndTime = settings.attendanceEndTime 
    ? dayjs(settings.attendanceEndTime, 'HH:mm') 
    : undefined
  formState.attendanceLateThreshold = settings.attendanceLateThreshold
  formState.attendanceVeryLateThreshold = settings.attendanceVeryLateThreshold
  formState.enableAttendanceNotification = settings.enableAttendanceNotification
  formState.enableGradeNotification = settings.enableGradeNotification
  formState.enableBKNotification = settings.enableBKNotification
  formState.enableHomeroomNotification = settings.enableHomeroomNotification
  formState.academicYear = settings.academicYear
  formState.semester = settings.semester
}

// Handle save
const handleSave = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  saving.value = true
  try {
    const updateData: UpdateSchoolSettingsRequest = {
      attendanceStartTime: formState.attendanceStartTime?.format('HH:mm'),
      attendanceEndTime: formState.attendanceEndTime?.format('HH:mm'),
      attendanceLateThreshold: formState.attendanceLateThreshold,
      attendanceVeryLateThreshold: formState.attendanceVeryLateThreshold,
      enableAttendanceNotification: formState.enableAttendanceNotification,
      enableGradeNotification: formState.enableGradeNotification,
      enableBKNotification: formState.enableBKNotification,
      enableHomeroomNotification: formState.enableHomeroomNotification,
      academicYear: formState.academicYear,
      semester: formState.semester,
    }
    
    await schoolService.updateSettings(updateData)
    message.success('Pengaturan berhasil disimpan')
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } } }
    message.error(err.response?.data?.message || 'Gagal menyimpan pengaturan')
  } finally {
    saving.value = false
  }
}

// Handle reset to defaults
const handleReset = async () => {
  resetting.value = true
  try {
    const settings = await schoolService.resetSettings()
    applySettings(settings)
    message.success('Pengaturan berhasil direset ke default')
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } } }
    message.error(err.response?.data?.message || 'Gagal mereset pengaturan')
  } finally {
    resetting.value = false
  }
}

// Get current academic year options
const academicYearOptions = () => {
  const currentYear = new Date().getFullYear()
  const options = []
  for (let i = -1; i <= 2; i++) {
    const year = currentYear + i
    options.push(`${year}/${year + 1}`)
  }
  return options
}

onMounted(() => {
  loadSettings()
})
</script>

<template>
  <div class="school-settings">
    <div class="page-header">
      <Title :level="2" style="margin: 0">Pengaturan Sekolah</Title>
      <Text type="secondary">Konfigurasi pengaturan absensi, notifikasi, dan tahun ajaran</Text>
    </div>

    <Spin :spinning="loading">
      <Form
        ref="formRef"
        :model="formState"
        :rules="formRules"
        layout="vertical"
      >
        <!-- Attendance Settings -->
        <Card class="settings-card">
          <template #title>
            <Space>
              <ClockCircleOutlined />
              <span>Pengaturan Absensi</span>
            </Space>
          </template>
          
          <Row :gutter="24">
            <Col :xs="24" :sm="12" :md="6">
              <FormItem label="Waktu Mulai Absensi" name="attendanceStartTime" required>
                <TimePicker
                  v-model:value="formState.attendanceStartTime"
                  format="HH:mm"
                  placeholder="07:00"
                  style="width: 100%"
                />
              </FormItem>
            </Col>
            <Col :xs="24" :sm="12" :md="6">
              <FormItem label="Waktu Akhir Absensi" name="attendanceEndTime" required>
                <TimePicker
                  v-model:value="formState.attendanceEndTime"
                  format="HH:mm"
                  placeholder="07:30"
                  style="width: 100%"
                />
              </FormItem>
            </Col>
            <Col :xs="24" :sm="12" :md="6">
              <FormItem
                label="Batas Terlambat (menit)"
                name="attendanceLateThreshold"
                extra="Menit setelah waktu mulai untuk dianggap terlambat"
              >
                <InputNumber
                  v-model:value="formState.attendanceLateThreshold"
                  :min="1"
                  :max="120"
                  style="width: 100%"
                />
              </FormItem>
            </Col>
            <Col :xs="24" :sm="12" :md="6">
              <FormItem
                label="Batas Sangat Terlambat (menit)"
                name="attendanceVeryLateThreshold"
                extra="Menit setelah waktu mulai untuk dianggap sangat terlambat"
              >
                <InputNumber
                  v-model:value="formState.attendanceVeryLateThreshold"
                  :min="1"
                  :max="180"
                  style="width: 100%"
                />
              </FormItem>
            </Col>
          </Row>
        </Card>

        <!-- Notification Settings -->
        <Card class="settings-card" style="margin-top: 24px">
          <template #title>
            <Space>
              <BellOutlined />
              <span>Pengaturan Notifikasi</span>
            </Space>
          </template>
          
          <Row :gutter="[24, 16]">
            <Col :xs="24" :sm="12" :md="6">
              <div class="notification-item">
                <div class="notification-label">
                  <Text strong>Notifikasi Absensi</Text>
                  <Text type="secondary" class="notification-desc">
                    Kirim notifikasi saat siswa check-in/check-out
                  </Text>
                </div>
                <Switch v-model:checked="formState.enableAttendanceNotification" />
              </div>
            </Col>
            <Col :xs="24" :sm="12" :md="6">
              <div class="notification-item">
                <div class="notification-label">
                  <Text strong>Notifikasi Nilai</Text>
                  <Text type="secondary" class="notification-desc">
                    Kirim notifikasi saat nilai baru ditambahkan
                  </Text>
                </div>
                <Switch v-model:checked="formState.enableGradeNotification" />
              </div>
            </Col>
            <Col :xs="24" :sm="12" :md="6">
              <div class="notification-item">
                <div class="notification-label">
                  <Text strong>Notifikasi BK</Text>
                  <Text type="secondary" class="notification-desc">
                    Kirim notifikasi untuk pelanggaran, prestasi, izin
                  </Text>
                </div>
                <Switch v-model:checked="formState.enableBKNotification" />
              </div>
            </Col>
            <Col :xs="24" :sm="12" :md="6">
              <div class="notification-item">
                <div class="notification-label">
                  <Text strong>Notifikasi Wali Kelas</Text>
                  <Text type="secondary" class="notification-desc">
                    Kirim notifikasi untuk catatan wali kelas
                  </Text>
                </div>
                <Switch v-model:checked="formState.enableHomeroomNotification" />
              </div>
            </Col>
          </Row>
        </Card>

        <!-- Academic Year Settings -->
        <Card class="settings-card" style="margin-top: 24px">
          <template #title>
            <Space>
              <CalendarOutlined />
              <span>Pengaturan Tahun Ajaran</span>
            </Space>
          </template>
          
          <Row :gutter="24">
            <Col :xs="24" :sm="12" :md="8">
              <FormItem label="Tahun Ajaran" name="academicYear" required>
                <Select v-model:value="formState.academicYear" placeholder="Pilih tahun ajaran">
                  <SelectOption v-for="year in academicYearOptions()" :key="year" :value="year">
                    {{ year }}
                  </SelectOption>
                </Select>
              </FormItem>
            </Col>
            <Col :xs="24" :sm="12" :md="8">
              <FormItem label="Semester" name="semester">
                <Select v-model:value="formState.semester" placeholder="Pilih semester">
                  <SelectOption :value="1">Semester 1 (Ganjil)</SelectOption>
                  <SelectOption :value="2">Semester 2 (Genap)</SelectOption>
                </Select>
              </FormItem>
            </Col>
          </Row>
        </Card>

        <!-- Action Buttons -->
        <Divider />
        <div class="form-actions">
          <Space>
            <Popconfirm
              title="Reset pengaturan ke default?"
              description="Semua pengaturan akan dikembalikan ke nilai awal."
              ok-text="Ya, Reset"
              cancel-text="Batal"
              @confirm="handleReset"
            >
              <Button :loading="resetting">
                <template #icon><ReloadOutlined /></template>
                Reset ke Default
              </Button>
            </Popconfirm>
            <Button type="primary" :loading="saving" @click="handleSave">
              <template #icon><SaveOutlined /></template>
              Simpan Pengaturan
            </Button>
          </Space>
        </div>
      </Form>
    </Spin>
  </div>
</template>

<style scoped>
.school-settings {
  padding: 0;
}

.page-header {
  margin-bottom: 24px;
}

.page-header h2 {
  margin-bottom: 4px;
}

.settings-card {
  background: #fff;
}

.settings-card :deep(.ant-card-head-title) {
  font-size: 16px;
}

.notification-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 12px 16px;
  background: #fafafa;
  border-radius: 8px;
  height: 100%;
}

.notification-label {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
  margin-right: 16px;
}

.notification-desc {
  font-size: 12px;
  line-height: 1.4;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  padding-top: 16px;
}

@media (max-width: 576px) {
  .form-actions {
    justify-content: center;
  }
}
</style>
