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
  Descriptions,
  DescriptionsItem,
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
const filterSchoolOption = (input: string, option: { label?: string }) => {
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

const handleShowApiKey = (device: Device) => {
  selectedDevice.value = device
  showApiKey.value = false
  apiKeyModalVisible.value = true
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
    message.success('API Key device ' + device.deviceCode + ' berhasil dicabut')
    loadDevices()
  } catch (error: unknown) {
    const err = error as { message?: string; response?: { data?: { error?: { message?: string } } } }
    message.error(err.message || err.response?.data?.error?.message || 'Gagal mencabut API Key')
  }
}

const handleRegenerateApiKey = async (device: Device) => {
  try {
    const updatedDevice = await deviceService.regenerateApiKey(device.id)
    message.success('API Key device ' + device.deviceCode + ' berhasil di-regenerate')
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
    message.success('Device ' + deviceToDelete.value.deviceCode + ' berhasil dihapus')
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
      <Title :level="2" style="margin: 0">Manajemen Device</Title>
      <Text type="secondary">Kelola device RFID per sekolah</Text>
    </div>

    <Card>
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col :xs="24" :sm="12" :md="8">
          <Input v-model:value="searchText" placeholder="Cari device atau sekolah..." allow-clear>
            <template #prefix><SearchOutlined /></template>
          </Input>
        </Col>
        <Col :xs="24" :sm="12" :md="8" class="toolbar-right">
          <Space>
            <Text type="secondary">Total: {{ totalDevices }} device</Text>
            <Button @click="loadDevices" :loading="loading">
              <template #icon><ReloadOutlined /></template>
              Refresh
            </Button>
            <Button type="primary" @click="openRegisterModal">
              <template #icon><PlusOutlined /></template>
              Daftarkan Device
            </Button>
          </Space>
        </Col>
      </Row>

      <div v-if="loading" style="text-align: center; padding: 40px;">Memuat data...</div>
      <Empty v-else-if="filteredSchoolDevices.length === 0" description="Tidak ada device terdaftar" />
      
      <Collapse v-else v-model:activeKey="activeKeys" class="school-collapse">
        <CollapsePanel v-for="school in filteredSchoolDevices" :key="school.schoolId.toString()">
          <template #header>
            <Space>
              <BankOutlined />
              <span class="school-name">{{ school.schoolName }}</span>
              <Tag v-if="!school.isActive" color="default">Nonaktif</Tag>
            </Space>
          </template>
          <template #extra>
            <Space @click.stop>
              <Badge :count="getOnlineCount(school.devices)" :number-style="{ backgroundColor: '#52c41a' }" :show-zero="false">
                <Tag color="blue">{{ school.deviceCount }} device</Tag>
              </Badge>
            </Space>
          </template>
          
          <div class="device-list">
            <div v-for="device in school.devices" :key="device.id" class="device-item">
              <Row :gutter="16" align="middle">
                <Col :xs="24" :sm="8" :md="6">
                  <div class="device-code">
                    <Text strong>{{ device.deviceCode }}</Text>
                    <Tag :color="getDeviceStatus(device).color" size="small">
                      <template #icon>
                        <WifiOutlined v-if="isDeviceOnline(device)" />
                        <DisconnectOutlined v-else />
                      </template>
                      {{ getDeviceStatus(device).text }}
                    </Tag>
                  </div>
                </Col>
                <Col :xs="24" :sm="8" :md="6">
                  <Text type="secondary">{{ device.description || '-' }}</Text>
                </Col>
                <Col :xs="24" :sm="8" :md="6">
                  <Text type="secondary" style="font-size: 12px;">Terakhir: {{ formatDate(device.lastSeenAt) }}</Text>
                </Col>
                <Col :xs="24" :sm="24" :md="6" class="device-actions">
                  <Space>
                    <Tooltip title="Lihat API Key">
                      <Button size="small" @click="handleShowApiKey(device)">
                        <template #icon><KeyOutlined /></template>
                      </Button>
                    </Tooltip>
                    <Popconfirm title="Regenerate API Key?" ok-text="Ya" @confirm="handleRegenerateApiKey(device)">
                      <Tooltip title="Regenerate">
                        <Button size="small" type="primary" ghost>
                          <template #icon><ReloadOutlined /></template>
                        </Button>
                      </Tooltip>
                    </Popconfirm>
                    <Popconfirm v-if="device.isActive" title="Cabut API Key?" ok-text="Ya" @confirm="handleRevokeApiKey(device)">
                      <Tooltip title="Cabut">
                        <Button size="small" danger ghost>
                          <template #icon><StopOutlined /></template>
                        </Button>
                      </Tooltip>
                    </Popconfirm>
                    <Tooltip title="Hapus Device">
                      <Button size="small" danger type="primary" @click="openDeleteModal(device)">
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
    <Modal v-model:open="modalVisible" title="Daftarkan Device Baru" :confirm-loading="modalLoading" @ok="handleSubmit" @cancel="handleModalCancel">
      <Form ref="formRef" :model="formState" :rules="formRules" layout="vertical" style="margin-top: 16px">
        <FormItem label="Sekolah" name="schoolId" required>
          <Select v-model:value="formState.schoolId" placeholder="Pilih sekolah" :loading="schoolsLoading" show-search :filter-option="filterSchoolOption">
            <SelectOption v-for="school in schools" :key="school.id" :value="school.id" :label="school.name">{{ school.name }}</SelectOption>
          </Select>
        </FormItem>
        <FormItem label="Kode Device" name="deviceCode" required>
          <Input v-model:value="formState.deviceCode" placeholder="Contoh: ESP32-JKT-001" />
        </FormItem>
        <FormItem label="Deskripsi" name="description">
          <Input v-model:value="formState.description" placeholder="Contoh: Pintu Utama" />
        </FormItem>
      </Form>
    </Modal>

    <!-- API Key Modal -->
    <Modal v-model:open="apiKeyModalVisible" title="API Key Device" :footer="null">
      <div v-if="selectedDevice" class="api-key-modal">
        <Descriptions :column="1" bordered size="small">
          <DescriptionsItem label="Kode Device">{{ selectedDevice.deviceCode }}</DescriptionsItem>
          <DescriptionsItem label="Sekolah">{{ selectedDevice.schoolName }}</DescriptionsItem>
          <DescriptionsItem label="API Key">
            <div class="api-key-display">
              <Text code>{{ showApiKey ? selectedDevice.apiKey : maskApiKey(selectedDevice.apiKey) }}</Text>
              <Space>
                <Tooltip :title="showApiKey ? 'Sembunyikan' : 'Tampilkan'">
                  <Button size="small" type="text" @click="showApiKey = !showApiKey">
                    <template #icon><EyeInvisibleOutlined v-if="showApiKey" /><EyeOutlined v-else /></template>
                  </Button>
                </Tooltip>
                <Tooltip title="Salin">
                  <Button size="small" type="text" @click="copyApiKey">
                    <template #icon><CopyOutlined /></template>
                  </Button>
                </Tooltip>
              </Space>
            </div>
          </DescriptionsItem>
        </Descriptions>
        <div class="api-key-warning">
          <Text type="warning">⚠️ Simpan API Key ini dengan aman.</Text>
        </div>
      </div>
    </Modal>

    <!-- Delete Device Modal -->
    <Modal v-model:open="deleteModalVisible" title="Hapus Device" :confirm-loading="deleteLoading" ok-text="Hapus" ok-type="danger" @ok="handleDelete">
      <div v-if="deviceToDelete" class="delete-modal">
        <Alert type="warning" show-icon style="margin-bottom: 16px;">
          <template #icon><ExclamationCircleOutlined /></template>
          <template #message>
            <Text>Apakah Anda yakin ingin menghapus device ini?</Text>
          </template>
          <template #description>
            <div style="margin-top: 8px;">
              <Text strong>{{ deviceToDelete.deviceCode }}</Text><br />
              <Text type="secondary">{{ deviceToDelete.schoolName }}</Text>
            </div>
          </template>
        </Alert>
        <Text type="secondary">Device yang dihapus tidak dapat dikembalikan.</Text>
      </div>
    </Modal>
  </div>
</template>

<style scoped>
.device-management { padding: 0; }
.page-header { margin-bottom: 24px; }
.toolbar { margin-bottom: 16px; }
.toolbar-right { display: flex; justify-content: flex-end; }
@media (max-width: 768px) { .toolbar-right { margin-top: 16px; justify-content: flex-start; } }
.school-collapse { background: transparent; }
.school-collapse :deep(.ant-collapse-item) { margin-bottom: 8px; border: 1px solid #f0f0f0; border-radius: 8px; overflow: hidden; }
.school-collapse :deep(.ant-collapse-header) { background: #fafafa; }
.school-name { font-weight: 500; }
.device-list { display: flex; flex-direction: column; gap: 12px; }
.device-item { padding: 12px 16px; background: #fafafa; border-radius: 6px; border: 1px solid #f0f0f0; }
.device-code { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.device-actions { display: flex; justify-content: flex-end; }
@media (max-width: 768px) { .device-actions { margin-top: 12px; justify-content: flex-start; } }
.api-key-modal { padding: 8px 0; }
.api-key-display { display: flex; align-items: center; justify-content: space-between; gap: 8px; }
.api-key-warning { margin-top: 16px; padding: 12px; background: #fffbe6; border-radius: 6px; }
.delete-modal { padding: 8px 0; }
</style>
