<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import {
  Table,
  Button,
  Space,
  Modal,
  Form,
  FormItem,
  Input,
  DatePicker,
  message,
  Popconfirm,
  Card,
  Row,
  Col,
  Typography,
  Tooltip,
  Alert,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  PlusOutlined,
  DeleteOutlined,
  ReloadOutlined,
  CopyOutlined,
  StopOutlined,
  SyncOutlined,
  KeyOutlined,
  LinkOutlined,
  CheckOutlined,
  ExclamationCircleOutlined,
} from '@ant-design/icons-vue'
import { displayTokenService } from '@/services'
import type { DisplayToken, DisplayTokenWithSecret, CreateDisplayTokenRequest } from '@/types/displayToken'
import dayjs from 'dayjs'

const { Title, Text } = Typography

// Table state
const loading = ref(false)
const tokens = ref<DisplayToken[]>([])

// Create modal state
const createModalVisible = ref(false)
const createModalLoading = ref(false)

// Token secret modal state (shown after create/regenerate)
const secretModalVisible = ref(false)
const newTokenData = ref<DisplayTokenWithSecret | null>(null)
const tokenCopied = ref(false)
const urlCopied = ref(false)

// Form state
const formRef = ref()
const formState = reactive({
  name: '',
  expiresAt: undefined as dayjs.Dayjs | undefined,
})

// Form rules
const formRules = {
  name: [
    { required: true, message: 'Nama display wajib diisi' },
    { max: 100, message: 'Nama maksimal 100 karakter' },
  ],
}

// Table columns
const columns: TableProps['columns'] = [
  {
    title: 'Nama Display',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: 'Status',
    key: 'status',
    width: 120,
    align: 'center',
  },
  {
    title: 'Terakhir Diakses',
    key: 'lastAccessed',
    width: 180,
  },
  {
    title: 'Kedaluwarsa',
    key: 'expires',
    width: 150,
  },
  {
    title: 'Dibuat',
    key: 'created',
    width: 150,
  },
  {
    title: 'Aksi',
    key: 'action',
    width: 280,
    align: 'center',
  },
]

// Load tokens data
const loadTokens = async () => {
  loading.value = true
  try {
    const response = await displayTokenService.getTokens()
    tokens.value = response.tokens
  } catch (err) {
    console.error('Failed to load display tokens:', err)
    message.error('Gagal memuat data display token')
    tokens.value = []
  } finally {
    loading.value = false
  }
}

// Open create modal
const openCreateModal = () => {
  resetForm()
  createModalVisible.value = true
}

// Reset form
const resetForm = () => {
  formState.name = ''
  formState.expiresAt = undefined
  formRef.value?.resetFields()
}

// Handle modal cancel
const handleCreateModalCancel = () => {
  createModalVisible.value = false
  resetForm()
}

// Handle form submit
const handleCreate = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  createModalLoading.value = true
  try {
    const createData: CreateDisplayTokenRequest = {
      name: formState.name,
      expires_at: formState.expiresAt?.toISOString(),
    }
    const result = await displayTokenService.createToken(createData)
    
    // Close create modal and show secret modal
    createModalVisible.value = false
    resetForm()
    
    // Show the token secret modal
    newTokenData.value = result
    tokenCopied.value = false
    urlCopied.value = false
    secretModalVisible.value = true
    
    // Reload tokens list
    loadTokens()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Gagal membuat display token')
  } finally {
    createModalLoading.value = false
  }
}

// Handle revoke token
const handleRevoke = async (token: DisplayToken) => {
  try {
    await displayTokenService.revokeToken(token.id)
    message.success(`Token "${token.name}" berhasil dicabut`)
    loadTokens()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Gagal mencabut token')
  }
}

// Handle regenerate token
const handleRegenerate = async (token: DisplayToken) => {
  try {
    const result = await displayTokenService.regenerateToken(token.id)
    
    // Show the new token secret modal
    newTokenData.value = result
    tokenCopied.value = false
    urlCopied.value = false
    secretModalVisible.value = true
    
    message.success(`Token "${token.name}" berhasil di-regenerate`)
    loadTokens()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Gagal regenerate token')
  }
}

// Handle delete token
const handleDelete = async (token: DisplayToken) => {
  try {
    await displayTokenService.deleteToken(token.id)
    message.success(`Token "${token.name}" berhasil dihapus`)
    loadTokens()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Gagal menghapus token')
  }
}

// Copy token to clipboard
const copyToken = async () => {
  if (!newTokenData.value?.token) return
  try {
    await navigator.clipboard.writeText(newTokenData.value.token)
    tokenCopied.value = true
    message.success('Token berhasil disalin')
  } catch {
    message.error('Gagal menyalin token')
  }
}

// Copy display URL to clipboard
const copyDisplayUrl = async (url?: string) => {
  const urlToCopy = url || newTokenData.value?.displayUrl
  if (!urlToCopy) return
  try {
    await navigator.clipboard.writeText(urlToCopy)
    if (!url) urlCopied.value = true
    message.success('URL display berhasil disalin')
  } catch {
    message.error('Gagal menyalin URL')
  }
}

// Close secret modal
const handleSecretModalClose = () => {
  secretModalVisible.value = false
  newTokenData.value = null
  tokenCopied.value = false
  urlCopied.value = false
}

// Format date for display
const formatDate = (dateStr?: string) => {
  if (!dateStr) return '-'
  return dayjs(dateStr).format('DD MMM YYYY HH:mm')
}

// Format relative time
const formatRelativeTime = (dateStr?: string) => {
  if (!dateStr) return 'Belum pernah'
  const date = dayjs(dateStr)
  const now = dayjs()
  const diffMins = now.diff(date, 'minute')
  const diffHours = now.diff(date, 'hour')
  const diffDays = now.diff(date, 'day')
  
  if (diffMins < 1) return 'Baru saja'
  if (diffMins < 60) return `${diffMins} menit lalu`
  if (diffHours < 24) return `${diffHours} jam lalu`
  if (diffDays < 7) return `${diffDays} hari lalu`
  return date.format('DD MMM YYYY')
}

// Check if token is expired
const isExpired = (token: DisplayToken) => {
  if (!token.expiresAt) return false
  return dayjs(token.expiresAt).isBefore(dayjs())
}

// Get status tag color
const getStatusColor = (token: DisplayToken) => {
  if (!token.isActive) return 'error'
  if (isExpired(token)) return 'warning'
  return 'success'
}

// Get status text
const getStatusText = (token: DisplayToken) => {
  if (!token.isActive) return 'Dicabut'
  if (isExpired(token)) return 'Kedaluwarsa'
  return 'Aktif'
}

onMounted(() => {
  loadTokens()
})
</script>

<template>
  <div class="display-token-management">
    <div class="page-header">
      <Title :level="2" style="margin: 0">Manajemen Display Token</Title>
      <Text type="secondary">Kelola token untuk akses public display absensi di LCD/monitor sekolah</Text>
    </div>

    <Card>
      <!-- Toolbar -->
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col>
          <Text type="secondary">
            <KeyOutlined /> {{ tokens.length }} token
          </Text>
        </Col>
        <Col class="toolbar-right">
          <Space>
            <Button @click="loadTokens">
              <template #icon><ReloadOutlined /></template>
            </Button>
            <Button type="primary" @click="openCreateModal">
              <template #icon><PlusOutlined /></template>
              Buat Token Baru
            </Button>
          </Space>
        </Col>
      </Row>

      <!-- Table -->
      <Table
        :columns="columns"
        :data-source="tokens"
        :loading="loading"
        :pagination="false"
        row-key="id"
        class="custom-table"
        :scroll="{ x: 800 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <div style="display: flex; flex-direction: column;">
              <Text strong>{{ (record as DisplayToken).name }}</Text>
              <Text v-if="(record as DisplayToken).displayUrl" type="secondary" style="font-size: 12px; color: #94a3b8;">
                <LinkOutlined /> {{ (record as DisplayToken).displayUrl }}
              </Text>
            </div>
          </template>
          <template v-else-if="column.key === 'status'">
            <span class="status-badge" :class="getStatusColor(record as DisplayToken)">
              {{ getStatusText(record as DisplayToken) }}
            </span>
          </template>
          <template v-else-if="column.key === 'lastAccessed'">
            <span class="text-secondary">{{ formatRelativeTime((record as DisplayToken).lastAccessedAt) }}</span>
          </template>
          <template v-else-if="column.key === 'expires'">
            <span v-if="(record as DisplayToken).expiresAt" :class="isExpired(record as DisplayToken) ? 'text-red-500' : 'text-secondary'">
              {{ formatDate((record as DisplayToken).expiresAt) }}
            </span>
            <span v-else class="text-secondary">Tidak ada</span>
          </template>
          <template v-else-if="column.key === 'created'">
            <span class="text-secondary">{{ formatDate((record as DisplayToken).createdAt) }}</span>
          </template>
          <template v-else-if="column.key === 'action'">
            <Space>
              <Tooltip title="Salin URL Display">
                <Button 
                  type="text"
                  :disabled="!(record as DisplayToken).displayUrl"
                  @click="copyDisplayUrl((record as DisplayToken).displayUrl)"
                >
                  <template #icon><CopyOutlined style="color: #3b82f6;" /></template>
                </Button>
              </Tooltip>
              <Tooltip title="Regenerate Token">
                <Popconfirm
                  title="Regenerate token ini?"
                  description="Token lama akan tidak valid dan diganti dengan token baru."
                  ok-text="Ya, Regenerate"
                  cancel-text="Batal"
                  @confirm="handleRegenerate(record as DisplayToken)"
                >
                  <Button type="text">
                    <template #icon><SyncOutlined style="color: #f59e0b;" /></template>
                  </Button>
                </Popconfirm>
              </Tooltip>
              <Tooltip v-if="(record as DisplayToken).isActive" title="Cabut Token">
                <Popconfirm
                  title="Cabut token ini?"
                  description="Display yang menggunakan token ini tidak akan bisa mengakses data."
                  ok-text="Ya, Cabut"
                  cancel-text="Batal"
                  @confirm="handleRevoke(record as DisplayToken)"
                >
                  <Button type="text" danger>
                    <template #icon><StopOutlined /></template>
                  </Button>
                </Popconfirm>
              </Tooltip>
              <Popconfirm
                title="Hapus token ini?"
                description="Token akan dihapus permanen."
                ok-text="Ya, Hapus"
                cancel-text="Batal"
                ok-type="danger"
                @confirm="handleDelete(record as DisplayToken)"
              >
                <Button type="text" danger>
                  <template #icon><DeleteOutlined /></template>
                </Button>
              </Popconfirm>
            </Space>
          </template>
        </template>
      </Table>
    </Card>

    <!-- Create Token Modal -->
    <Modal
      v-model:open="createModalVisible"
      title="Buat Display Token Baru"
      :confirm-loading="createModalLoading"
      width="500px"
      wrap-class-name="modern-modal"
      @ok="handleCreate"
      @cancel="handleCreateModalCancel"
    >
      <Form
        ref="formRef"
        :model="formState"
        :rules="formRules"
        layout="vertical"
        class="modern-form"
      >
        <FormItem label="Nama Display" name="name" required>
          <Input 
            v-model:value="formState.name" 
            placeholder="Contoh: Display Pintu Utama, LCD Lobby" 
          />
          <Text type="secondary" style="font-size: 12px">
            Nama untuk mengidentifikasi lokasi display
          </Text>
        </FormItem>
        
        <FormItem label="Tanggal Kedaluwarsa (Opsional)" name="expiresAt">
          <DatePicker
            v-model:value="formState.expiresAt"
            format="DD MMMM YYYY"
            placeholder="Pilih tanggal"
            style="width: 100%"
            :disabled-date="(current: dayjs.Dayjs) => current && current < dayjs().startOf('day')"
          />
          <Text type="secondary" style="font-size: 12px">
            Kosongkan jika token tidak memiliki batas waktu
          </Text>
        </FormItem>
      </Form>
    </Modal>

    <!-- Token Secret Modal (shown after create/regenerate) -->
    <Modal
      v-model:open="secretModalVisible"
      title="Token Berhasil Dibuat"
      :footer="null"
      :closable="false"
      :maskClosable="false"
      width="600px"
      wrap-class-name="modern-modal"
    >
      <Alert
        type="warning"
        show-icon
        style="margin-bottom: 24px"
      >
        <template #icon><ExclamationCircleOutlined /></template>
        <template #message>Simpan token ini sekarang!</template>
        <template #description>
          Token hanya ditampilkan sekali dan tidak dapat dilihat lagi setelah modal ini ditutup.
        </template>
      </Alert>

      <div v-if="newTokenData" class="token-info">
        <div class="token-field">
          <Text strong>Nama Display:</Text>
          <Text>{{ newTokenData.name }}</Text>
        </div>

        <div class="token-field">
          <Text strong>Token:</Text>
          <div class="token-value">
            <Input.Password
              :value="newTokenData.token"
              readonly
              style="flex: 1"
            />
            <Button 
              :type="tokenCopied ? 'primary' : 'default'"
              @click="copyToken"
            >
              <template #icon>
                <CheckOutlined v-if="tokenCopied" />
                <CopyOutlined v-else />
              </template>
              {{ tokenCopied ? 'Tersalin' : 'Salin' }}
            </Button>
          </div>
        </div>

        <div class="token-field">
          <Text strong>URL Display:</Text>
          <div class="token-value">
            <Input
              :value="newTokenData.displayUrl"
              readonly
              style="flex: 1"
            />
            <Button 
              :type="urlCopied ? 'primary' : 'default'"
              @click="copyDisplayUrl()"
            >
              <template #icon>
                <CheckOutlined v-if="urlCopied" />
                <CopyOutlined v-else />
              </template>
              {{ urlCopied ? 'Tersalin' : 'Salin' }}
            </Button>
          </div>
          <Text type="secondary" style="font-size: 13px; margin-top: 8px; display: block">
            Gunakan URL ini untuk mengakses public display di browser
          </Text>
        </div>
      </div>

      <div style="margin-top: 32px; text-align: right">
        <Button type="primary" @click="handleSecretModalClose">
          Saya Sudah Menyimpan Token
        </Button>
      </div>
    </Modal>
  </div>
</template>

<style scoped>
.display-token-management {
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

.token-info {
  background: #fafafa;
  border-radius: 8px;
  padding: 16px;
}

.token-field {
  margin-bottom: 16px;
}

.token-field:last-child {
  margin-bottom: 0;
}

.token-field > span:first-child {
  display: block;
  margin-bottom: 8px;
}

.token-value {
  display: flex;
  gap: 8px;
}

@media (max-width: 576px) {
  .toolbar-right {
    margin-top: 16px;
    justify-content: flex-start;
  }
  
  .token-value {
    flex-direction: column;
  }
}

/* Custom Table */
.custom-table :deep(.ant-table-thead > tr > th) {
  background: #f8fafc;
  color: #475569;
  font-weight: 600;
  border-bottom: 1px solid #f1f5f9;
}

.custom-table :deep(.ant-table-tbody > tr > td) {
  padding: 16px;
  border-bottom: 1px solid #f1f5f9;
}

/* Status Badges */
.status-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  min-width: 80px;
}

.status-badge.success { background: #f0fdf4; color: #166534; border: 1px solid #dcfce7; }
.status-badge.warning { background: #fffbeb; color: #b45309; border: 1px solid #fef3c7; }
.status-badge.error { background: #fef2f2; color: #991b1b; border: 1px solid #fee2e2; }

/* Text Styles */
.text-secondary {
  color: #64748b;
  font-size: 13px;
}

/* Token Info in Modal */
.token-info {
  background: #f8fafc;
  border-radius: 8px;
  padding: 20px;
  border: 1px solid #e2e8f0;
}
</style>
