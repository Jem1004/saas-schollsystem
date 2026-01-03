<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import {
  Button,
  Input,
  Space,
  Tag,
  Modal,
  Form,
  FormItem,
  Select,
  SelectOption,
  message,
  Popconfirm,
  Card,
  Row,
  Col,
  Typography,
  Tooltip,
  Collapse,
  CollapsePanel,
  Empty,
  Badge,
  Alert,
} from 'ant-design-vue'
import {
  PlusOutlined,
  SearchOutlined,
  KeyOutlined,
  StopOutlined,
  ReloadOutlined,
  WifiOutlined,
  DisconnectOutlined,
  CopyOutlined,
  EyeOutlined,
  EyeInvisibleOutlined,
  DeleteOutlined,
  BankOutlined,
  ExclamationCircleOutlined,
} from '@ant-design/icons-vue'
import { deviceService, tenantService } from '@/services'
import type { Device, SchoolDevices } from '@/types/device'
import type { School } from '@/types/tenant'

const { Title, Text } = Typography

// State
const loading = ref(false)
const schoolDevices = ref<SchoolDevices[]>([])
const totalDevices = ref(0)
const searchText = ref('')
const activeKeys = ref<string[]>([])

// Schools for dropdown
const schools = ref<School[]>([])
const schoolsLoading = ref(false)

// Modal state
const modalVisible = ref(false)
const modalLoading = ref(false)

// API Key modal state
const apiKeyModalVisible = ref(false)
const selectedDevice = ref<Device | null>(null)
const showApiKey = ref(false)

// Delete modal state
const deleteModalVisible = ref(false)
const deleteLoading = ref(false)
const deviceToDelete = ref<Device | null>(null)

// Form state
const formRef = ref()
const formState = reactive({
  schoolId: 0,
  deviceCode: '',
  description: '',
})

// Form rules
const formRules = {
  schoolId: [{ required: true, message: 'Sekolah wajib dipilih' }],
  deviceCode: [
    { required: true, message: 'Kode device wajib diisi' },
    { pattern: /^[A-Za-z0-9-_]+$/, message: 'Kode device hanya boleh huruf, angka, - dan _' },
  ],
}

// Filter option for school select
const filterSchoolOption = (input: string, option: any) => {
  return option?.label?.toLowerCase().includes(input.toLowerCase()) ?? false
}

// Check if device is online (last seen within 5 minutes)
const isDeviceOnline = (device: Device): boolean => {
  if (!device.lastSeenAt || !device.isActive) return false
  const lastSeen = new Date(device.lastSeenAt)
  const fiveMinutesAgo = new Date(Date.now() - 5 * 60 * 1000)
  return lastSeen > fiveMinutesAgo
}

// Get device status
const getDeviceStatus = (device: Device): { text: string; color: string } => {
  if (!device.isActive) return { text: 'Nonaktif', color: 'default' }
  if (isDeviceOnline(device)) return { text: 'Online', color: 'success' }
  return { text: 'Offline', color: 'warning' }
}

// Filtered school devices based on search
const filteredSchoolDevices = computed(() => {
  if (!searchText.value) return schoolDevices.value
  const search = searchText.value.toLowerCase()
  return schoolDevices.value
    .map(school => ({
      ...school,
      devices: school.devices.filter(d => 
        d.deviceCode.toLowerCase().includes(search) ||
        d.description?.toLowerCase().includes(search)
      )
    }))
    .filter(school => school.devices.length > 0 || school.schoolName.toLowerCase().includes(search))
})

// Count online devices for a school
const getOnlineCount = (devices: Device[]): number => {
  return devices.filter(d => isDeviceOnline(d)).length
}

// Load devices grouped by school
const loadDevices = async () => {
  loading.value = true
  try {
    const response = await deviceService.getDevicesGrouped()
    schoolDevices.value = response.schools
    totalDevices.value = response.total
    if (response.schools.length <= 5) {
      activeKeys.value = response.schools.map(s => s.schoolId.toString())
    }
  } catch {
    message.error('Gagal memuat data device')
  } finally {
    loading.value = false
  }
}

// Load schools for dropdown
const loadSchools = async () => {
  schoolsLoading.value = true
  try {
    const response = await tenantService.getSchools({ pageSize: 100 })
    schools.value = response.data.filter((s) => s.isActive)
  } catch {
    message.error('Gagal memuat data sekolah')
  } finally {
    schoolsLoading.value = false
  }
}

const openRegisterModal = () => {
  resetForm()
  modalVisible.value = true
}

const resetForm = () => {
  formState.schoolId = 0
  formState.deviceCode = ''
  formState.description = ''
  formRef.value?.resetFields()
}

const handleModalCancel = () => {
  modalVisible.value = false
  resetForm()
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
  } catch { return }

  modalLoading.value = true
  try {
    const newDevice = await deviceService.registerDevice(formState)
    message.success('Device berhasil didaftarkan')
    modalVisible.value = false
    resetForm()
    selectedDevice.value = newDevice
    showApiKey.value = true
    apiKeyModalVisible.value = true
    loadDevices()
  } catch (error: unknown) {
    const err = error as { message?: string; response?: { data?: { error?: { message?: string } } } }
    message.error(err.message || err.response?.data?.error?.message || 'Gagal mendaftarkan device')
  } finally {
    modalLoading.value = false
  }
}

const handleShowApiKey = async (device: Device) => {
  try {
    const apiKeyData = await deviceService.getDeviceApiKey(device.id)
    selectedDevice.value = {
      ...device,
      apiKey: apiKeyData.apiKey,
    }
    showApiKey.value = false
    apiKeyModalVisible.value = true
  } catch (error: unknown) {
    const err = error as { message?: string }
    message.error(err.message || 'Gagal mengambil API Key')
  }
}

const copyApiKey = async () => {
  if (!selectedDevice.value) return
  try {
    await navigator.clipboard.writeText(selectedDevice.value.apiKey)
    message.success('API Key berhasil disalin')
  } catch {
    message.error('Gagal menyalin API Key')
  }
}

const handleRevokeApiKey = async (device: Device) => {
  try {
    await deviceService.revokeApiKey(device.id)
    message.success(`API Key device ${device.deviceCode} berhasil dicabut`)
    loadDevices()
  } catch (error: unknown) {
    const err = error as { message?: string; response?: { data?: { error?: { message?: string } } } }
    message.error(err.message || err.response?.data?.error?.message || 'Gagal mencabut API Key')
  }
}

const handleRegenerateApiKey = async (device: Device) => {
  try {
    const updatedDevice = await deviceService.regenerateApiKey(device.id)
    message.success(`API Key device ${device.deviceCode} berhasil di-regenerate`)
    selectedDevice.value = updatedDevice
    showApiKey.value = true
    apiKeyModalVisible.value = true
    loadDevices()
  } catch (error: unknown) {
    const err = error as { message?: string; response?: { data?: { error?: { message?: string } } } }
    message.error(err.message || err.response?.data?.error?.message || 'Gagal regenerate API Key')
  }
}

const openDeleteModal = (device: Device) => {
  deviceToDelete.value = device
  deleteModalVisible.value = true
}

const handleDelete = async () => {
  if (!deviceToDelete.value) return
  deleteLoading.value = true
  try {
    await deviceService.deleteDevice(deviceToDelete.value.id)
    message.success(`Device ${deviceToDelete.value.deviceCode} berhasil dihapus`)
    deleteModalVisible.value = false
    deviceToDelete.value = null
    loadDevices()
  } catch (error: unknown) {
    const err = error as { message?: string; response?: { data?: { error?: { message?: string } } } }
    message.error(err.message || err.response?.data?.error?.message || 'Gagal menghapus device')
  } finally {
    deleteLoading.value = false
  }
}

const formatDate = (dateStr?: string): string => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('id-ID', {
    day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit',
  })
}

const maskApiKey = (apiKey: string): string => {
  if (!apiKey || apiKey.length <= 8) return '********'
  return apiKey.substring(0, 8) + '...' + apiKey.substring(apiKey.length - 4)
}

onMounted(() => {
  loadDevices()
  loadSchools()
})
</script>

<template>
  <div class="device-management">
    <div class="page-header">
      <div class="header-content">
        <Title :level="2" class="page-title">Manajemen Device</Title>
        <Text class="page-subtitle">Kelola device RFID per sekolah</Text>
      </div>
    </div>

    <Card class="content-card" :bordered="false">
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col :xs="24" :sm="12" :md="10">
          <Input 
            v-model:value="searchText" 
            placeholder="Cari device atau sekolah..." 
            allow-clear
            size="large"
            class="search-input"
          >
            <template #prefix><SearchOutlined /></template>
          </Input>
        </Col>
        <Col :xs="24" :sm="12" :md="10" class="toolbar-right">
          <Space size="middle">
            <div class="device-count">
              <Text type="secondary">Total:</Text>
              <Text strong>{{ totalDevices }}</Text>
            </div>
            <Button @click="loadDevices" :loading="loading" size="large">
              <template #icon><ReloadOutlined /></template>
            </Button>
            <Button type="primary" @click="openRegisterModal" size="large">
              <template #icon><PlusOutlined /></template>
              Daftarkan Device
            </Button>
          </Space>
        </Col>
      </Row>

      <div v-if="loading" style="text-align: center; padding: 40px;">
        <Space direction="vertical">
             <ReloadOutlined spin style="font-size: 24px; color: #f97316" />
             <Text type="secondary">Memuat data device...</Text>
        </Space>
      </div>
      <Empty v-else-if="filteredSchoolDevices.length === 0" description="Tidak ada device terdaftar" />
      
      <Collapse 
        v-else 
        v-model:activeKey="activeKeys" 
        class="custom-collapse" 
        ghost 
        accordion
        expand-icon-position="end"
      >
        <CollapsePanel v-for="school in filteredSchoolDevices" :key="school.schoolId.toString()">
          <template #header>
            <div class="panel-header">
              <div class="school-info">
                <div class="school-icon">
                  <BankOutlined />
                </div>
                <div class="school-text">
                  <Text strong class="school-title">{{ school.schoolName }}</Text>
                  <Tag v-if="!school.isActive" color="default" :bordered="false" size="small">Nonaktif</Tag>
                </div>
              </div>
              <div class="school-stats" @click.stop>
                <Badge 
                  :count="getOnlineCount(school.devices)" 
                  :number-style="{ backgroundColor: '#22c55e' }" 
                  :show-zero="true" 
                  title="Device Online"
                >
                  <span class="device-badge-label">Online</span>
                </Badge>
                <div class="divider-vertical"></div>
                <Text type="secondary">{{ school.deviceCount }} Total</Text>
              </div>
            </div>
          </template>
          
          <div class="device-list">
            <div v-for="device in school.devices" :key="device.id" class="device-item">
              <Row :gutter="16" align="middle">
                <Col :xs="24" :sm="8" :md="6" class="device-col">
                  <div class="device-info">
                    <div class="device-icon" :class="isDeviceOnline(device) ? 'online' : 'offline'">
                       <WifiOutlined v-if="isDeviceOnline(device)" />
                       <DisconnectOutlined v-else />
                    </div>
                    <div class="device-text">
                      <Text strong class="device-code-text">{{ device.deviceCode }}</Text>
                      <Text type="secondary" class="device-status-text">
                        {{ getDeviceStatus(device).text }}
                      </Text>
                    </div>
                  </div>
                </Col>
                <Col :xs="24" :sm="8" :md="8" class="device-col">
                  <div class="device-meta">
                    <Text type="secondary" style="font-size: 13px">Lokasi / Deskripsi</Text>
                    <Text>{{ device.description || '-' }}</Text>
                  </div>
                </Col>
                <Col :xs="24" :sm="8" :md="6" class="device-col">
                  <div class="device-meta">
                    <Text type="secondary" style="font-size: 13px">Terakhir Dilihat</Text>
                    <Text>{{ formatDate(device.lastSeenAt) }}</Text>
                  </div>
                </Col>
                <Col :xs="24" :sm="24" :md="4" class="device-actions">
                  <Space size="small">
                    <Tooltip title="Lihat API Key">
                      <Button type="text" size="small" @click="handleShowApiKey(device)">
                        <template #icon><KeyOutlined style="color: #64748b" /></template>
                      </Button>
                    </Tooltip>
                    <Popconfirm title="Regenerate API Key?" ok-text="Ya" @confirm="handleRegenerateApiKey(device)">
                      <Tooltip title="Regenerate">
                        <Button type="text" size="small">
                          <template #icon><ReloadOutlined style="color: #f97316" /></template>
                        </Button>
                      </Tooltip>
                    </Popconfirm>
                    <Popconfirm v-if="device.isActive" title="Cabut API Key?" ok-text="Ya" @confirm="handleRevokeApiKey(device)">
                      <Tooltip title="Cabut Akses">
                        <Button type="text" size="small" danger>
                          <template #icon><StopOutlined /></template>
                        </Button>
                      </Tooltip>
                    </Popconfirm>
                    <Tooltip title="Hapus Device">
                      <Button type="text" size="small" danger @click="openDeleteModal(device)">
                        <template #icon><DeleteOutlined /></template>
                      </Button>
                    </Tooltip>
                  </Space>
                </Col>
              </Row>
            </div>
          </div>
        </CollapsePanel>
      </Collapse>
    </Card>


    <!-- Register Device Modal -->
    <Modal
      v-model:open="modalVisible"
      title="Daftarkan Device Baru"
      :confirm-loading="modalLoading"
      ok-text="Daftarkan Device"
      cancel-text="Batal"
      width="500px"
      wrap-class-name="modern-modal"
      @ok="handleSubmit"
      @cancel="handleModalCancel"
    >
      <Form
        ref="formRef"
        :model="formState"
        :rules="formRules"
        layout="vertical"
        class="modern-form"
      >
        <FormItem label="Sekolah" name="schoolId" required>
          <Select 
            v-model:value="formState.schoolId" 
            placeholder="Pilih sekolah" 
            :loading="schoolsLoading" 
            show-search 
            :filter-option="filterSchoolOption"
            size="large"
          >
            <SelectOption v-for="school in schools" :key="school.id" :value="school.id" :label="school.name">{{ school.name }}</SelectOption>
          </Select>
        </FormItem>
        <FormItem label="Kode Device" name="deviceCode" required>
          <Input v-model:value="formState.deviceCode" placeholder="Contoh: ESP32-JKT-001" size="large" />
        </FormItem>
        <FormItem label="Deskripsi" name="description">
          <Input v-model:value="formState.description" placeholder="Contoh: Pintu Utama" size="large" />
        </FormItem>
      </Form>
    </Modal>

    <!-- API Key Modal -->
    <Modal
      v-model:open="apiKeyModalVisible"
      title="API Key Device"
      :footer="null"
      width="480px"
      wrap-class-name="modern-modal"
    >
      <div v-if="selectedDevice" class="api-key-modal-content">
        <div class="device-summary">
          <div class="summary-item">
            <span class="label">Device</span>
            <span class="value">{{ selectedDevice.deviceCode }}</span>
          </div>
          <div class="summary-item">
            <span class="label">Sekolah</span>
            <span class="value">{{ selectedDevice.schoolName }}</span>
          </div>
        </div>

        <div class="api-key-box">
          <Text type="secondary" class="box-label">API Key</Text>
          <div class="key-display">
            <Text code class="key-text">{{ showApiKey ? selectedDevice.apiKey : maskApiKey(selectedDevice.apiKey) }}</Text>
            <Space :size="4">
              <Tooltip :title="showApiKey ? 'Sembunyikan' : 'Tampilkan'">
                <Button type="text" size="small" @click="showApiKey = !showApiKey">
                  <template #icon><EyeInvisibleOutlined v-if="showApiKey" /><EyeOutlined v-else /></template>
                </Button>
              </Tooltip>
              <Tooltip title="Salin">
                <Button type="text" size="small" @click="copyApiKey">
                  <template #icon><CopyOutlined /></template>
                </Button>
              </Tooltip>
            </Space>
          </div>
        </div>

        <Alert
          message="Simpan API Key ini dengan aman"
          description="Jangan bagikan key ini kepada pihak yang tidak berwenang."
          type="warning"
          show-icon
          style="margin-top: 24px; border-radius: 8px;"
        />

        <div class="modal-actions" style="margin-top: 24px;">
           <Button type="primary" block size="large" @click="apiKeyModalVisible = false">Selesai</Button>
        </div>
      </div>
    </Modal>

    <!-- Delete Device Modal -->
    <Modal
      v-model:open="deleteModalVisible"
      title="Hapus Device"
      :confirm-loading="deleteLoading"
      ok-text="Ya, Hapus"
      cancel-text="Batal"
      ok-type="danger"
      wrap-class-name="modern-modal"
      :ok-button-props="{ size: 'large' }"
      :cancel-button-props="{ size: 'large' }"
      @ok="handleDelete"
    >
      <div v-if="deviceToDelete" class="delete-confirmation-content">
        <Alert
          type="warning"
          show-icon
          class="delete-alert"
        >
          <template #icon><ExclamationCircleOutlined /></template>
          <template #message>
            <Text strong>Konfirmasi Penghapusan</Text>
          </template>
          <template #description>
             Apakah Anda yakin ingin menghapus device <Text strong>{{ deviceToDelete.deviceCode }}</Text>? 
             Tindakan ini tidak dapat dibatalkan.
          </template>
        </Alert>
        
        <div class="device-preview">
           <div class="preview-row">
             <span class="label">Sekolah:</span>
             <span class="value">{{ deviceToDelete.schoolName }}</span>
           </div>
           <div class="preview-row">
             <span class="label">Status:</span>
             <Tag :color="deviceToDelete.isActive ? 'green' : 'default'" :bordered="false" size="small">
                {{ deviceToDelete.isActive ? 'Aktif' : 'Nonaktif' }}
             </Tag>
           </div>
        </div>
      </div>
    </Modal>
  </div>
</template>

<style scoped>
.device-management {
  padding: 0;
  max-width: 1600px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 32px;
}

.page-title {
  font-weight: 600 !important;
  color: #1e293b !important;
  margin-bottom: 8px !important;
  letter-spacing: -0.5px;
}

.page-subtitle {
  font-size: 14px;
  color: #64748b;
}

.content-card {
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
  border: 1px solid #f1f5f9;
}

.toolbar {
  margin-bottom: 24px;
}

.toolbar-right {
  display: flex;
  justify-content: flex-end;
}

.device-count {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-right: 16px;
}

/* Custom Collapse Styling */
.custom-collapse {
  background: transparent;
  border: none;
}

.custom-collapse :deep(.ant-collapse-item) {
  margin-bottom: 12px;
  border: 1px solid #f1f5f9;
  border-radius: 8px !important;
  background: #ffffff;
  overflow: hidden;
}

.custom-collapse :deep(.ant-collapse-header) {
  background: #ffffff;
  padding: 16px 24px !important;
  align-items: center;
  border-bottom: 1px solid transparent;
  transition: all 0.2s ease;
}

.custom-collapse :deep(.ant-collapse-item-active .ant-collapse-header) {
  border-bottom-color: #f1f5f9;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.school-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.school-icon {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background: #fff7ed;
  color: #f97316;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}

.school-text {
  display: flex;
  align-items: center;
  gap: 12px;
}

.school-title {
  font-size: 15px;
  color: #1e293b;
}

.school-stats {
  display: flex;
  align-items: center;
  gap: 12px;
}

.device-badge-label {
  margin-right: 8px;
  font-size: 12px;
  color: #64748b;
}

.divider-vertical {
  width: 1px;
  height: 16px;
  background: #e2e8f0;
}

.device-list {
  display: flex;
  flex-direction: column;
}

.device-item {
  padding: 16px 24px;
  border-bottom: 1px solid #f8fafc;
  transition: background 0.2s;
}

.device-item:last-child {
  border-bottom: none;
}

.device-item:hover {
  background: #f8fafc;
}

.device-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.device-icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}

.device-icon.online {
  background: #f0fdf4;
  color: #22c55e;
}

.device-icon.offline {
  background: #f1f5f9;
  color: #94a3b8;
}

.device-code-text {
  font-size: 14px;
  color: #334155;
  display: block;
}

.device-status-text {
  font-size: 12px;
  display: block;
}

.device-meta {
  display: flex;
  flex-direction: column;
}

.device-actions {
  display: flex;
  justify-content: flex-end;
}

/* API Key Modal Styling */
.api-key-modal-content {
  padding-top: 8px;
}

.device-summary {
  display: flex;
  background: #f8fafc;
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 20px;
  border: 1px solid #f1f5f9;
}

.summary-item {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.summary-item .label {
  font-size: 12px;
  color: #64748b;
  margin-bottom: 2px;
}

.summary-item .value {
  font-weight: 500;
  color: #1e293b;
}

.api-key-box {
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 16px;
  background: white;
}

.box-label {
  display: block;
  font-size: 12px;
  margin-bottom: 8px;
}

.key-display {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.key-text {
  font-size: 16px;
  color: #0f172a;
}

/* Delete Modal Styling */
.delete-confirmation-content {
  padding-top: 8px;
}

.delete-alert {
  margin-bottom: 20px;
  border-radius: 8px;
}

.device-preview {
  background: #f8fafc;
  padding: 16px;
  border-radius: 8px;
  border: 1px solid #f1f5f9;
}

.preview-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 14px;
}

.preview-row:last-child {
  margin-bottom: 0;
}

.preview-row .label {
  color: #64748b;
}

.preview-row .value {
  font-weight: 500;
  color: #334155;
}

/* Responsiveness */
@media (max-width: 768px) {
  .toolbar-right { margin-top: 16px; justify-content: flex-start; }
  .device-actions { margin-top: 12px; justify-content: flex-start; }
  .panel-header { flex-direction: column; align-items: flex-start; gap: 8px; }
  .school-stats { width: 100%; justify-content: space-between; }
  
  .device-col {
    margin-bottom: 12px;
  }
}
</style>

<!-- Global Styles for Modals (Teleported) -->
<style>
/* Modern Modal Global Overrides */
.modern-modal .ant-modal-content {
  border-radius: 16px !important;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1) !important;
  padding: 0 !important;
  overflow: hidden;
}

.modern-modal .ant-modal-header {
  border-bottom: 1px solid #f1f5f9 !important;
  padding: 20px 24px !important;
  background: #fff !important;
}

.modern-modal .ant-modal-title {
  font-size: 18px !important;
  font-weight: 600 !important;
  color: #0f172a !important;
}

.modern-modal .ant-modal-body {
  padding: 24px !important;
}

.modern-modal .ant-modal-footer {
  border-top: 1px solid #f1f5f9 !important;
  padding: 16px 24px !important;
  background: #ffffff !important;
}

.modern-modal .ant-btn {
  border-radius: 8px !important;
  height: 40px !important;
  font-weight: 500 !important;
}

.modern-modal .ant-btn-primary {
  box-shadow: 0 4px 6px -1px rgba(249, 115, 22, 0.2) !important;
}
</style>
