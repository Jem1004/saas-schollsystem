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
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
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
} from '@ant-design/icons-vue'
import { deviceService, tenantService } from '@/services'
import type { Device, RegisterDeviceRequest } from '@/types/device'
import type { School } from '@/types/tenant'

const { Title, Text } = Typography

// Table state
const loading = ref(false)
const devices = ref<Device[]>([])
const total = ref(0)
const pagination = reactive({
  current: 1,
  pageSize: 10,
})
const searchText = ref('')
const filterSchoolId = ref<number | undefined>(undefined)

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

// Form state
const formRef = ref()
const formState = reactive<RegisterDeviceRequest>({
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

// Mock data for development
const mockDevices: Device[] = [
  { id: 1, schoolId: 1, schoolName: 'SMP Negeri 1 Jakarta', deviceCode: 'ESP32-JKT-001', apiKey: 'sk_live_abc123xyz789', description: 'Pintu Utama', isActive: true, lastSeenAt: new Date().toISOString(), createdAt: '2024-01-15T10:00:00Z', updatedAt: '2024-01-15T10:00:00Z' },
  { id: 2, schoolId: 1, schoolName: 'SMP Negeri 1 Jakarta', deviceCode: 'ESP32-JKT-002', apiKey: 'sk_live_def456uvw012', description: 'Pintu Belakang', isActive: true, lastSeenAt: new Date(Date.now() - 3600000).toISOString(), createdAt: '2024-01-14T09:00:00Z', updatedAt: '2024-01-14T09:00:00Z' },
  { id: 3, schoolId: 2, schoolName: 'SMP Negeri 2 Bandung', deviceCode: 'ESP32-BDG-001', apiKey: 'sk_live_ghi789rst345', description: 'Gerbang Utama', isActive: true, lastSeenAt: new Date(Date.now() - 7200000).toISOString(), createdAt: '2024-01-13T08:00:00Z', updatedAt: '2024-01-13T08:00:00Z' },
  { id: 4, schoolId: 3, schoolName: 'SMP Negeri 3 Surabaya', deviceCode: 'ESP32-SBY-001', apiKey: 'sk_live_jkl012mno678', description: 'Pintu Masuk', isActive: false, createdAt: '2024-01-12T11:00:00Z', updatedAt: '2024-01-20T14:00:00Z' },
  { id: 5, schoolId: 4, schoolName: 'SMP Negeri 4 Yogyakarta', deviceCode: 'ESP32-JOG-001', apiKey: 'sk_live_pqr345stu901', description: 'Lobby', isActive: true, lastSeenAt: new Date(Date.now() - 300000).toISOString(), createdAt: '2024-01-11T07:00:00Z', updatedAt: '2024-01-11T07:00:00Z' },
]

const mockSchools: School[] = [
  { id: 1, name: 'SMP Negeri 1 Jakarta', isActive: true, createdAt: '2024-01-15', updatedAt: '2024-01-15' },
  { id: 2, name: 'SMP Negeri 2 Bandung', isActive: true, createdAt: '2024-01-14', updatedAt: '2024-01-14' },
  { id: 3, name: 'SMP Negeri 3 Surabaya', isActive: false, createdAt: '2024-01-13', updatedAt: '2024-01-13' },
  { id: 4, name: 'SMP Negeri 4 Yogyakarta', isActive: true, createdAt: '2024-01-12', updatedAt: '2024-01-12' },
]

// Filter option for school select
const filterSchoolOption = (input: string, option: { label?: string }) => {
  return option?.label?.toLowerCase().includes(input.toLowerCase()) ?? false
}

// Table columns
const columns: TableProps['columns'] = [
  {
    title: 'Kode Device',
    dataIndex: 'deviceCode',
    key: 'deviceCode',
    sorter: true,
  },
  {
    title: 'Sekolah',
    dataIndex: 'schoolName',
    key: 'schoolName',
  },
  {
    title: 'Deskripsi',
    dataIndex: 'description',
    key: 'description',
  },
  {
    title: 'Status',
    key: 'status',
    width: 120,
    align: 'center',
  },
  {
    title: 'Terakhir Online',
    dataIndex: 'lastSeenAt',
    key: 'lastSeenAt',
    width: 150,
  },
  {
    title: 'Aksi',
    key: 'action',
    width: 250,
    align: 'center',
  },
]

// Check if device is online (last seen within 5 minutes)
const isDeviceOnline = (device: Device): boolean => {
  if (!device.lastSeenAt || !device.isActive) return false
  const lastSeen = new Date(device.lastSeenAt)
  const fiveMinutesAgo = new Date(Date.now() - 5 * 60 * 1000)
  return lastSeen > fiveMinutesAgo
}

// Get device status
const getDeviceStatus = (device: Device): { text: string; color: string } => {
  if (!device.isActive) {
    return { text: 'Nonaktif', color: 'default' }
  }
  if (isDeviceOnline(device)) {
    return { text: 'Online', color: 'success' }
  }
  return { text: 'Offline', color: 'warning' }
}

// Computed filtered data
const filteredDevices = computed(() => {
  let result = devices.value
  
  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(
      (device) =>
        device.deviceCode.toLowerCase().includes(search) ||
        device.schoolName?.toLowerCase().includes(search) ||
        device.description?.toLowerCase().includes(search)
    )
  }
  
  if (filterSchoolId.value) {
    result = result.filter((device) => device.schoolId === filterSchoolId.value)
  }
  
  return result
})

// Load devices data
const loadDevices = async () => {
  loading.value = true
  try {
    const response = await deviceService.getDevices({
      page: pagination.current,
      pageSize: pagination.pageSize,
      search: searchText.value,
      schoolId: filterSchoolId.value,
    })
    devices.value = response.data
    total.value = response.total
  } catch {
    // Use mock data on error
    devices.value = mockDevices
    total.value = mockDevices.length
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
    schools.value = mockSchools.filter((s) => s.isActive)
  } finally {
    schoolsLoading.value = false
  }
}

// Handle table change
const handleTableChange: TableProps['onChange'] = (pag) => {
  pagination.current = pag.current || 1
  pagination.pageSize = pag.pageSize || 10
  loadDevices()
}

// Handle search
const handleSearch = () => {
  pagination.current = 1
  loadDevices()
}

// Open modal for register
const openRegisterModal = () => {
  resetForm()
  modalVisible.value = true
}

// Reset form
const resetForm = () => {
  formState.schoolId = 0
  formState.deviceCode = ''
  formState.description = ''
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
    const newDevice = await deviceService.registerDevice(formState)
    message.success('Device berhasil didaftarkan')
    modalVisible.value = false
    resetForm()
    
    // Show API key modal
    selectedDevice.value = newDevice
    showApiKey.value = true
    apiKeyModalVisible.value = true
    
    loadDevices()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } } }
    message.error(err.response?.data?.message || 'Gagal mendaftarkan device')
  } finally {
    modalLoading.value = false
  }
}

// Show API key
const handleShowApiKey = (device: Device) => {
  selectedDevice.value = device
  showApiKey.value = false
  apiKeyModalVisible.value = true
}

// Copy API key to clipboard
const copyApiKey = async () => {
  if (!selectedDevice.value) return
  try {
    await navigator.clipboard.writeText(selectedDevice.value.apiKey)
    message.success('API Key berhasil disalin')
  } catch {
    message.error('Gagal menyalin API Key')
  }
}

// Handle revoke API key
const handleRevokeApiKey = async (device: Device) => {
  try {
    await deviceService.revokeApiKey(device.id)
    message.success(`API Key device ${device.deviceCode} berhasil dicabut`)
    loadDevices()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } } }
    message.error(err.response?.data?.message || 'Gagal mencabut API Key')
  }
}

// Handle regenerate API key
const handleRegenerateApiKey = async (device: Device) => {
  try {
    const updatedDevice = await deviceService.regenerateApiKey(device.id)
    message.success(`API Key device ${device.deviceCode} berhasil di-regenerate`)
    
    // Show new API key
    selectedDevice.value = updatedDevice
    showApiKey.value = true
    apiKeyModalVisible.value = true
    
    loadDevices()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } } }
    message.error(err.response?.data?.message || 'Gagal regenerate API Key')
  }
}

// Format date
const formatDate = (dateStr?: string): string => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

// Mask API key
const maskApiKey = (apiKey: string): string => {
  if (apiKey.length <= 8) return '********'
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
    </div>

    <Card>
      <!-- Toolbar -->
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col :xs="24" :sm="24" :md="16">
          <Space wrap>
            <Input
              v-model:value="searchText"
              placeholder="Cari device..."
              allow-clear
              style="width: 200px"
              @press-enter="handleSearch"
            >
              <template #prefix>
                <SearchOutlined />
              </template>
            </Input>
            <Select
              v-model:value="filterSchoolId"
              placeholder="Filter Sekolah"
              allow-clear
              style="width: 200px"
              :loading="schoolsLoading"
              @change="handleSearch"
            >
              <SelectOption v-for="school in schools" :key="school.id" :value="school.id">
                {{ school.name }}
              </SelectOption>
            </Select>
          </Space>
        </Col>
        <Col :xs="24" :sm="24" :md="8" class="toolbar-right">
          <Space>
            <Button @click="loadDevices">
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

      <!-- Table -->
      <Table
        :columns="columns"
        :data-source="filteredDevices"
        :loading="loading"
        :pagination="{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: total,
          showSizeChanger: true,
          showTotal: (total: number) => `Total ${total} device`,
        }"
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <Tag :color="getDeviceStatus(record as Device).color">
              <template #icon>
                <WifiOutlined v-if="isDeviceOnline(record as Device)" />
                <DisconnectOutlined v-else />
              </template>
              {{ getDeviceStatus(record as Device).text }}
            </Tag>
          </template>
          <template v-else-if="column.key === 'lastSeenAt'">
            {{ formatDate((record as Device).lastSeenAt) }}
          </template>
          <template v-else-if="column.key === 'action'">
            <Space>
              <Tooltip title="Lihat API Key">
                <Button size="small" @click="handleShowApiKey(record as Device)">
                  <template #icon><KeyOutlined /></template>
                </Button>
              </Tooltip>
              <Popconfirm
                title="Regenerate API Key?"
                description="API Key lama akan tidak berlaku lagi."
                ok-text="Ya, Regenerate"
                cancel-text="Batal"
                @confirm="handleRegenerateApiKey(record as Device)"
              >
                <Tooltip title="Regenerate API Key">
                  <Button size="small" type="primary" ghost>
                    <template #icon><ReloadOutlined /></template>
                  </Button>
                </Tooltip>
              </Popconfirm>
              <Popconfirm
                v-if="(record as Device).isActive"
                title="Cabut API Key?"
                description="Device tidak akan bisa mengirim data absensi."
                ok-text="Ya, Cabut"
                cancel-text="Batal"
                @confirm="handleRevokeApiKey(record as Device)"
              >
                <Tooltip title="Cabut API Key">
                  <Button size="small" danger>
                    <template #icon><StopOutlined /></template>
                  </Button>
                </Tooltip>
              </Popconfirm>
            </Space>
          </template>
        </template>
      </Table>
    </Card>

    <!-- Register Device Modal -->
    <Modal
      v-model:open="modalVisible"
      title="Daftarkan Device Baru"
      :confirm-loading="modalLoading"
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
        <FormItem label="Sekolah" name="schoolId" required>
          <Select
            v-model:value="formState.schoolId"
            placeholder="Pilih sekolah"
            :loading="schoolsLoading"
            show-search
            :filter-option="filterSchoolOption"
          >
            <SelectOption
              v-for="school in schools"
              :key="school.id"
              :value="school.id"
              :label="school.name"
            >
              {{ school.name }}
            </SelectOption>
          </Select>
        </FormItem>
        <FormItem label="Kode Device" name="deviceCode" required>
          <Input
            v-model:value="formState.deviceCode"
            placeholder="Contoh: ESP32-JKT-001"
          />
        </FormItem>
        <FormItem label="Deskripsi" name="description">
          <Input
            v-model:value="formState.description"
            placeholder="Contoh: Pintu Utama"
          />
        </FormItem>
      </Form>
    </Modal>

    <!-- API Key Modal -->
    <Modal
      v-model:open="apiKeyModalVisible"
      title="API Key Device"
      :footer="null"
      @cancel="apiKeyModalVisible = false"
    >
      <div v-if="selectedDevice" class="api-key-modal">
        <Descriptions :column="1" bordered size="small">
          <DescriptionsItem label="Kode Device">
            {{ selectedDevice.deviceCode }}
          </DescriptionsItem>
          <DescriptionsItem label="Sekolah">
            {{ selectedDevice.schoolName }}
          </DescriptionsItem>
          <DescriptionsItem label="API Key">
            <div class="api-key-display">
              <Text code>
                {{ showApiKey ? selectedDevice.apiKey : maskApiKey(selectedDevice.apiKey) }}
              </Text>
              <Space>
                <Tooltip :title="showApiKey ? 'Sembunyikan' : 'Tampilkan'">
                  <Button size="small" type="text" @click="showApiKey = !showApiKey">
                    <template #icon>
                      <EyeInvisibleOutlined v-if="showApiKey" />
                      <EyeOutlined v-else />
                    </template>
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
          <Text type="warning">
            ⚠️ Simpan API Key ini dengan aman. API Key hanya ditampilkan sekali saat pendaftaran.
          </Text>
        </div>
      </div>
    </Modal>
  </div>
</template>

<style scoped>
.device-management {
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

@media (max-width: 768px) {
  .toolbar-right {
    margin-top: 16px;
    justify-content: flex-start;
  }
}

.api-key-modal {
  padding: 8px 0;
}

.api-key-display {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.api-key-warning {
  margin-top: 16px;
  padding: 12px;
  background: #fffbe6;
  border-radius: 6px;
}
</style>
