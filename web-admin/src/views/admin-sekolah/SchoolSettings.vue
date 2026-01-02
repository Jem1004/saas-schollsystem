<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  Card,
  Form,
  FormItem,
  Switch,
  Button,
  Row,
  Col,
  Typography,
  Divider,
  Select,
  SelectOption,
  message,
  Spin,
  Popconfirm,
  Alert,
} from 'ant-design-vue'
import {
  SaveOutlined,
  ReloadOutlined,
  BellOutlined,
  CalendarOutlined,
} from '@ant-design/icons-vue'
import { schoolService } from '@/services'
import type { SchoolSettings, UpdateSchoolSettingsRequest } from '@/types/school'

const { Title, Text } = Typography
const router = useRouter()

// State
const loading = ref(false)
const saving = ref(false)
const resetting = ref(false)

// Form state - removed attendance settings
const formRef = ref()
const formState = reactive<{
  enableAttendanceNotification: boolean
  enableGradeNotification: boolean
  enableBKNotification: boolean
  enableHomeroomNotification: boolean
  academicYear: string
  semester: number
}>({
  enableAttendanceNotification: true,
  enableGradeNotification: true,
  enableBKNotification: true,
  enableHomeroomNotification: true,
  academicYear: '',
  semester: 1,
})

// Form rules
const formRules = {
  academicYear: [{ required: true, message: 'Tahun ajaran wajib diisi' }],
}

// Handle authorization error
const handleAuthError = (error: unknown) => {
  const err = error as { response?: { data?: { error?: { code?: string; message?: string; debug?: Record<string, unknown> }; message?: string }; status?: number } }
  const errorCode = err.response?.data?.error?.code
  const errorMessage = err.response?.data?.error?.message || err.response?.data?.message
  const debugInfo = err.response?.data?.error?.debug
  
  // Log debug info for troubleshooting
  if (debugInfo) {
    console.error('Authorization Error Debug Info:', debugInfo)
  }
  
  if (errorCode === 'AUTHZ_ROLE_DENIED' || err.response?.status === 403) {
    const yourRole = debugInfo?.your_role as string | undefined
    const allowedRoles = debugInfo?.allowed_roles as string[] | undefined
    
    let errorMsg = 'Anda tidak memiliki akses ke halaman ini.'
    if (yourRole && allowedRoles) {
      errorMsg = `Role Anda (${yourRole}) tidak diizinkan. Halaman ini hanya untuk: ${allowedRoles.join(', ')}`
    }
    
    message.error(errorMsg)
    router.push('/access-denied')
    return true
  }
  
  if (errorCode === 'AUTHZ_TENANT_REQUIRED') {
    message.error('Konteks sekolah tidak ditemukan. Silakan login ulang.')
    router.push('/login')
    return true
  }
  
  message.error(errorMessage || 'Terjadi kesalahan')
  return false
}

// Load settings
const loadSettings = async () => {
  loading.value = true
  try {
    const settings = await schoolService.getSettings()
    applySettings(settings)
  } catch (err) {
    console.error('Failed to load settings:', err)
    if (!handleAuthError(err)) {
      message.error('Gagal memuat pengaturan sekolah')
    }
  } finally {
    loading.value = false
  }
}

// Apply settings to form - removed attendance settings
const applySettings = (settings: SchoolSettings) => {
  formState.enableAttendanceNotification = settings.enableAttendanceNotification
  formState.enableGradeNotification = settings.enableGradeNotification
  formState.enableBKNotification = settings.enableBKNotification
  formState.enableHomeroomNotification = settings.enableHomeroomNotification
  formState.academicYear = settings.academicYear
  formState.semester = settings.semester
}

// Handle save - removed attendance settings
const handleSave = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  saving.value = true
  try {
    const updateData: UpdateSchoolSettingsRequest = {
      enable_attendance_notification: formState.enableAttendanceNotification,
      enable_grade_notification: formState.enableGradeNotification,
      enable_bk_notification: formState.enableBKNotification,
      enable_homeroom_notification: formState.enableHomeroomNotification,
      academic_year: formState.academicYear,
      semester: formState.semester,
    }
    
    await schoolService.updateSettings(updateData)
    message.success('Pengaturan berhasil disimpan')
  } catch (error: unknown) {
    if (!handleAuthError(error)) {
      const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
      message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Gagal menyimpan pengaturan')
    }
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
    if (!handleAuthError(error)) {
      const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
      message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Gagal mereset pengaturan')
    }
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
      <Text type="secondary">Konfigurasi notifikasi dan tahun ajaran</Text>
    </div>

    <Spin :spinning="loading">
      <Form
        ref="formRef"
        :model="formState"
        :rules="formRules"
        layout="vertical"
      >
        <!-- Notification Settings -->
        <Card class="settings-card">
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
