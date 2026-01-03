<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import {
  Table, Button, Input, Space, Modal, Form, FormItem, Select, SelectOption,
  message, Popconfirm, Card, Row, Col, Typography, Textarea, InputNumber, Switch,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import { PlusOutlined, EditOutlined, DeleteOutlined, ReloadOutlined, ArrowLeftOutlined } from '@ant-design/icons-vue'
import { useRouter } from 'vue-router'
import { bkService } from '@/services'
import type { ViolationCategory, CreateViolationCategoryRequest, UpdateViolationCategoryRequest } from '@/types/bk'
import { VIOLATION_LEVELS } from '@/types/bk'

const { Title, Text } = Typography
const router = useRouter()

const loading = ref(false)
const categories = ref<ViolationCategory[]>([])
const showInactive = ref(false)

// Modal state
const modalVisible = ref(false)
const modalLoading = ref(false)
const isEditing = ref(false)
const editingId = ref<number | null>(null)

// Form state
const formRef = ref()
const formState = reactive<CreateViolationCategoryRequest & { isActive?: boolean }>({
  name: '',
  defaultPoint: -5,
  defaultLevel: 'ringan',
  description: '',
  isActive: true,
})

const formRules = {
  name: [{ required: true, message: 'Nama kategori wajib diisi' }],
  defaultPoint: [{ required: true, message: 'Poin default wajib diisi' }],
  defaultLevel: [{ required: true, message: 'Tingkat default wajib dipilih' }],
}

const columns: TableProps['columns'] = [
  { title: 'Nama Kategori', dataIndex: 'name', key: 'name' },
  { title: 'Poin Default', dataIndex: 'defaultPoint', key: 'defaultPoint', width: 120, align: 'center' },
  { title: 'Tingkat Default', dataIndex: 'defaultLevel', key: 'defaultLevel', width: 130, align: 'center' },
  { title: 'Deskripsi', dataIndex: 'description', key: 'description', ellipsis: true },
  { title: 'Status', dataIndex: 'isActive', key: 'isActive', width: 100, align: 'center' },
  { title: 'Aksi', key: 'action', width: 120, align: 'center' },
]

const getLevelClass = (level: string) => {
  switch (level) {
    case 'ringan': return 'success'
    case 'sedang': return 'warning'
    case 'berat': return 'error'
    case 'sangat_berat': return 'error'
    default: return 'default'
  }
}

const getLevelLabel = (level: string) => {
  const config = VIOLATION_LEVELS.find(l => l.value === level)
  return config?.label || level
}

const loadCategories = async () => {
  loading.value = true
  try {
    const response = await bkService.getViolationCategories(!showInactive.value)
    categories.value = response.categories || []
  } catch (err) {
    console.error('Failed to load categories:', err)
    categories.value = []
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  isEditing.value = false
  editingId.value = null
  resetForm()
  modalVisible.value = true
}

const openEditModal = (category: ViolationCategory) => {
  isEditing.value = true
  editingId.value = category.id
  formState.name = category.name
  formState.defaultPoint = category.defaultPoint
  formState.defaultLevel = category.defaultLevel
  formState.description = category.description || ''
  formState.isActive = category.isActive
  modalVisible.value = true
}

const resetForm = () => {
  formState.name = ''
  formState.defaultPoint = -5
  formState.defaultLevel = 'ringan'
  formState.description = ''
  formState.isActive = true
  formRef.value?.resetFields()
}

const handleModalCancel = () => {
  modalVisible.value = false
  resetForm()
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  modalLoading.value = true
  try {
    if (isEditing.value && editingId.value) {
      const updateData: UpdateViolationCategoryRequest = {
        name: formState.name,
        defaultPoint: formState.defaultPoint,
        defaultLevel: formState.defaultLevel,
        description: formState.description,
        isActive: formState.isActive,
      }
      await bkService.updateViolationCategory(editingId.value, updateData)
      
      // Show appropriate message based on active status change
      if (formState.isActive === false) {
        message.success('Kategori berhasil dinonaktifkan. Aktifkan "Tampilkan kategori tidak aktif" untuk melihatnya.')
        showInactive.value = true // Auto-show inactive categories
      } else {
        message.success('Kategori berhasil diperbarui')
      }
    } else {
      await bkService.createViolationCategory(formState)
      message.success('Kategori berhasil dibuat')
    }
    modalVisible.value = false
    resetForm()
    loadCategories()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Terjadi kesalahan')
  } finally {
    modalLoading.value = false
  }
}

const handleDelete = async (category: ViolationCategory) => {
  try {
    await bkService.deleteViolationCategory(category.id)
    message.success('Kategori berhasil dihapus')
    loadCategories()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string }; message?: string } } }
    message.error(err.response?.data?.error?.message || err.response?.data?.message || 'Gagal menghapus kategori')
  }
}

const goBack = () => {
  router.push('/bk/violations')
}

onMounted(() => {
  loadCategories()
})
</script>

<template>
  <div class="category-management">
    <div class="page-header">
      <Row justify="space-between" align="middle">
        <Col>
          <Space>
            <Button @click="goBack"><template #icon><ArrowLeftOutlined /></template></Button>
            <Title :level="2" style="margin: 0">Kelola Kategori Pelanggaran</Title>
          </Space>
        </Col>
      </Row>
      <Text type="secondary">Atur kategori pelanggaran dan poin default untuk sekolah Anda</Text>
    </div>

    <Card>
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col>
          <Space>
            <Switch v-model:checked="showInactive" @change="loadCategories" />
            <Text>Tampilkan kategori tidak aktif</Text>
          </Space>
        </Col>
        <Col>
          <Space>
            <Button size="large" @click="loadCategories"><template #icon><ReloadOutlined /></template></Button>
            <Button type="primary" size="large" @click="openCreateModal"><template #icon><PlusOutlined /></template>Tambah Kategori</Button>
          </Space>
        </Col>
      </Row>

      <Table 
        :columns="columns" 
        :data-source="categories" 
        :loading="loading" 
        :pagination="false" 
        row-key="id"
        class="custom-table"
        :scroll="{ x: 800 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'defaultPoint'">
            <Text strong :style="{ color: '#ef4444' }">{{ (record as ViolationCategory).defaultPoint }}</Text>
          </template>
          <template v-else-if="column.key === 'defaultLevel'">
            <span :class="['status-badge', getLevelClass((record as ViolationCategory).defaultLevel)]">
              {{ getLevelLabel((record as ViolationCategory).defaultLevel) }}
            </span>
          </template>
          <template v-else-if="column.key === 'isActive'">
            <span :class="['status-badge', (record as ViolationCategory).isActive ? 'success' : 'default']">
              {{ (record as ViolationCategory).isActive ? 'Aktif' : 'Tidak Aktif' }}
            </span>
          </template>
          <template v-else-if="column.key === 'action'">
            <Space>
              <Button type="text" style="color: #3b82f6" @click="openEditModal(record as ViolationCategory)">
                <template #icon><EditOutlined /></template>
              </Button>
              <Popconfirm title="Hapus kategori ini?" description="Kategori akan dihapus permanen." ok-text="Ya, Hapus" cancel-text="Batal" @confirm="handleDelete(record as ViolationCategory)">
                <Button type="text" danger><template #icon><DeleteOutlined /></template></Button>
              </Popconfirm>
            </Space>
          </template>
        </template>
      </Table>
    </Card>

    <!-- Create/Edit Modal -->
    <Modal v-model:open="modalVisible" :title="isEditing ? 'Edit Kategori' : 'Tambah Kategori Baru'" :confirm-loading="modalLoading" @ok="handleSubmit" @cancel="handleModalCancel" width="500px" wrap-class-name="modern-modal">
      <Form ref="formRef" :model="formState" :rules="formRules" layout="vertical" style="margin-top: 16px">
        <FormItem label="Nama Kategori" name="name" required>
          <Input v-model:value="formState.name" placeholder="Contoh: Keterlambatan, Bolos, dll" size="large" />
        </FormItem>
        <Row :gutter="16">
          <Col :span="12">
            <FormItem label="Poin Default" name="defaultPoint" required>
              <InputNumber v-model:value="formState.defaultPoint" :max="0" style="width: 100%" size="large" />
              <Text type="secondary" style="font-size: 12px">Harus 0 atau negatif</Text>
            </FormItem>
          </Col>
          <Col :span="12">
            <FormItem label="Tingkat Default" name="defaultLevel" required>
              <Select v-model:value="formState.defaultLevel" size="large">
                <SelectOption v-for="level in VIOLATION_LEVELS" :key="level.value" :value="level.value">
                   {{ level.label }}
                </SelectOption>
              </Select>
            </FormItem>
          </Col>
        </Row>
        <FormItem label="Deskripsi" name="description">
          <Textarea v-model:value="formState.description" placeholder="Deskripsi kategori (opsional)" :rows="3" class="custom-textarea" />
        </FormItem>
        <FormItem v-if="isEditing" label="Status">
          <Switch v-model:checked="formState.isActive" />
          <Text style="margin-left: 8px">{{ formState.isActive ? 'Aktif' : 'Tidak Aktif' }}</Text>
        </FormItem>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.category-management { padding: 0; }
.page-header { margin-bottom: 24px; }
.toolbar { margin-bottom: 16px; }

/* Custom Table Styles */
.custom-table :deep(.ant-table-thead > tr > th) {
  background: #fafafa;
  font-weight: 600;
  color: #475569;
}

.custom-table :deep(.ant-table-tbody > tr > td) {
  padding: 16px;
}

.custom-table :deep(.ant-table-tbody > tr:hover > td) {
  background: #f8fafc;
}

/* Badge Styles */
.status-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap;
}

.status-badge.success {
  background-color: #f6ffed;
  color: #52c41a;
  border: 1px solid #b7eb8f;
}

.status-badge.warning {
  background-color: #fffbe6;
  color: #faad14;
  border: 1px solid #ffe58f;
}

.status-badge.error {
  background-color: #fff2f0;
  color: #ff4d4f;
  border: 1px solid #ffccc7;
}

.status-badge.default {
  background-color: #f5f5f5;
  color: #000000d9;
  border: 1px solid #d9d9d9;
}
</style>
