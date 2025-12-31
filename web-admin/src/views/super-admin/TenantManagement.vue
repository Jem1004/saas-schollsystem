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
  message,
  Popconfirm,
  Card,
  Row,
  Col,
  Typography,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  PlusOutlined,
  SearchOutlined,
  EditOutlined,
  CheckCircleOutlined,
  StopOutlined,
  ReloadOutlined,
} from '@ant-design/icons-vue'
import { tenantService } from '@/services'
import type { School, CreateSchoolRequest, UpdateSchoolRequest } from '@/types/tenant'

const { Title } = Typography

// Table state
const loading = ref(false)
const schools = ref<School[]>([])
const total = ref(0)
const pagination = reactive({
  current: 1,
  pageSize: 10,
})
const searchText = ref('')

// Modal state
const modalVisible = ref(false)
const modalLoading = ref(false)
const isEditing = ref(false)
const editingSchool = ref<School | null>(null)

// Form state
const formRef = ref()
const formState = reactive<CreateSchoolRequest>({
  name: '',
  address: '',
  phone: '',
  email: '',
})

// Form rules
const formRules = {
  name: [{ required: true, message: 'Nama sekolah wajib diisi' }],
  email: [{ type: 'email' as const, message: 'Format email tidak valid' }],
}

// Mock data for development
const mockSchools: School[] = [
  { id: 1, name: 'SMP Negeri 1 Jakarta', address: 'Jl. Merdeka No. 1, Jakarta Pusat', phone: '021-1234567', email: 'smpn1@jakarta.sch.id', isActive: true, createdAt: '2024-01-15T10:00:00Z', updatedAt: '2024-01-15T10:00:00Z' },
  { id: 2, name: 'SMP Negeri 2 Bandung', address: 'Jl. Asia Afrika No. 10, Bandung', phone: '022-7654321', email: 'smpn2@bandung.sch.id', isActive: true, createdAt: '2024-01-14T09:00:00Z', updatedAt: '2024-01-14T09:00:00Z' },
  { id: 3, name: 'SMP Negeri 3 Surabaya', address: 'Jl. Pemuda No. 5, Surabaya', phone: '031-9876543', email: 'smpn3@surabaya.sch.id', isActive: false, createdAt: '2024-01-13T08:00:00Z', updatedAt: '2024-01-20T14:00:00Z' },
  { id: 4, name: 'SMP Negeri 4 Yogyakarta', address: 'Jl. Malioboro No. 20, Yogyakarta', phone: '0274-123456', email: 'smpn4@jogja.sch.id', isActive: true, createdAt: '2024-01-12T11:00:00Z', updatedAt: '2024-01-12T11:00:00Z' },
  { id: 5, name: 'SMP Negeri 5 Semarang', address: 'Jl. Pandanaran No. 15, Semarang', phone: '024-8765432', email: 'smpn5@semarang.sch.id', isActive: true, createdAt: '2024-01-11T07:00:00Z', updatedAt: '2024-01-11T07:00:00Z' },
]

// Table columns
const columns: TableProps['columns'] = [
  {
    title: 'Nama Sekolah',
    dataIndex: 'name',
    key: 'name',
    sorter: true,
  },
  {
    title: 'Email',
    dataIndex: 'email',
    key: 'email',
  },
  {
    title: 'Telepon',
    dataIndex: 'phone',
    key: 'phone',
  },
  {
    title: 'Status',
    dataIndex: 'isActive',
    key: 'isActive',
    width: 100,
    align: 'center',
  },
  {
    title: 'Tanggal Dibuat',
    dataIndex: 'createdAt',
    key: 'createdAt',
    width: 150,
  },
  {
    title: 'Aksi',
    key: 'action',
    width: 200,
    align: 'center',
  },
]

// Computed filtered data
const filteredSchools = computed(() => {
  if (!searchText.value) return schools.value
  const search = searchText.value.toLowerCase()
  return schools.value.filter(
    (school) =>
      school.name.toLowerCase().includes(search) ||
      school.email?.toLowerCase().includes(search) ||
      school.phone?.includes(search)
  )
})

// Load schools data
const loadSchools = async () => {
  loading.value = true
  try {
    const response = await tenantService.getSchools({
      page: pagination.current,
      pageSize: pagination.pageSize,
      search: searchText.value,
    })
    schools.value = response.data
    total.value = response.total
  } catch {
    // Use mock data on error
    schools.value = mockSchools
    total.value = mockSchools.length
  } finally {
    loading.value = false
  }
}

// Handle table change (pagination, sorting)
const handleTableChange: TableProps['onChange'] = (pag) => {
  pagination.current = pag.current || 1
  pagination.pageSize = pag.pageSize || 10
  loadSchools()
}

// Handle search
const handleSearch = () => {
  pagination.current = 1
  loadSchools()
}

// Open modal for create
const openCreateModal = () => {
  isEditing.value = false
  editingSchool.value = null
  resetForm()
  modalVisible.value = true
}

// Open modal for edit
const openEditModal = (school: School) => {
  isEditing.value = true
  editingSchool.value = school
  formState.name = school.name
  formState.address = school.address || ''
  formState.phone = school.phone || ''
  formState.email = school.email || ''
  modalVisible.value = true
}

// Reset form
const resetForm = () => {
  formState.name = ''
  formState.address = ''
  formState.phone = ''
  formState.email = ''
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
    if (isEditing.value && editingSchool.value) {
      const updateData: UpdateSchoolRequest = {
        name: formState.name,
        address: formState.address || undefined,
        phone: formState.phone || undefined,
        email: formState.email || undefined,
      }
      await tenantService.updateSchool(editingSchool.value.id, updateData)
      message.success('Sekolah berhasil diperbarui')
    } else {
      await tenantService.createSchool(formState)
      message.success('Sekolah berhasil ditambahkan')
    }
    modalVisible.value = false
    resetForm()
    loadSchools()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } } }
    message.error(err.response?.data?.message || 'Terjadi kesalahan')
  } finally {
    modalLoading.value = false
  }
}

// Handle activate school
const handleActivate = async (school: School) => {
  try {
    await tenantService.activateSchool(school.id)
    message.success(`Sekolah ${school.name} berhasil diaktifkan`)
    loadSchools()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } } }
    message.error(err.response?.data?.message || 'Gagal mengaktifkan sekolah')
  }
}

// Handle deactivate school
const handleDeactivate = async (school: School) => {
  try {
    await tenantService.deactivateSchool(school.id)
    message.success(`Sekolah ${school.name} berhasil dinonaktifkan`)
    loadSchools()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } } }
    message.error(err.response?.data?.message || 'Gagal menonaktifkan sekolah')
  }
}

// Format date
const formatDate = (dateStr: string): string => {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

onMounted(() => {
  loadSchools()
})
</script>

<template>
  <div class="tenant-management">
    <div class="page-header">
      <Title :level="2" style="margin: 0">Manajemen Sekolah</Title>
    </div>

    <Card>
      <!-- Toolbar -->
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col :xs="24" :sm="12" :md="8">
          <Input
            v-model:value="searchText"
            placeholder="Cari sekolah..."
            allow-clear
            @press-enter="handleSearch"
          >
            <template #prefix>
              <SearchOutlined />
            </template>
          </Input>
        </Col>
        <Col :xs="24" :sm="12" :md="8" class="toolbar-right">
          <Space>
            <Button @click="loadSchools">
              <template #icon><ReloadOutlined /></template>
              Refresh
            </Button>
            <Button type="primary" @click="openCreateModal">
              <template #icon><PlusOutlined /></template>
              Tambah Sekolah
            </Button>
          </Space>
        </Col>
      </Row>

      <!-- Table -->
      <Table
        :columns="columns"
        :data-source="filteredSchools"
        :loading="loading"
        :pagination="{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: total,
          showSizeChanger: true,
          showTotal: (total: number) => `Total ${total} sekolah`,
        }"
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'isActive'">
            <Tag :color="(record as School).isActive ? 'success' : 'default'">
              {{ (record as School).isActive ? 'Aktif' : 'Nonaktif' }}
            </Tag>
          </template>
          <template v-else-if="column.key === 'createdAt'">
            {{ formatDate((record as School).createdAt) }}
          </template>
          <template v-else-if="column.key === 'action'">
            <Space>
              <Button size="small" @click="openEditModal(record as School)">
                <template #icon><EditOutlined /></template>
                Edit
              </Button>
              <Popconfirm
                v-if="(record as School).isActive"
                title="Nonaktifkan sekolah ini?"
                description="Semua user di sekolah ini tidak akan bisa mengakses sistem."
                ok-text="Ya, Nonaktifkan"
                cancel-text="Batal"
                @confirm="handleDeactivate(record as School)"
              >
                <Button size="small" danger>
                  <template #icon><StopOutlined /></template>
                  Nonaktifkan
                </Button>
              </Popconfirm>
              <Popconfirm
                v-else
                title="Aktifkan sekolah ini?"
                ok-text="Ya, Aktifkan"
                cancel-text="Batal"
                @confirm="handleActivate(record as School)"
              >
                <Button size="small" type="primary" ghost>
                  <template #icon><CheckCircleOutlined /></template>
                  Aktifkan
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
      :title="isEditing ? 'Edit Sekolah' : 'Tambah Sekolah Baru'"
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
        <FormItem label="Nama Sekolah" name="name" required>
          <Input v-model:value="formState.name" placeholder="Masukkan nama sekolah" />
        </FormItem>
        <FormItem label="Alamat" name="address">
          <Input.TextArea
            v-model:value="formState.address"
            placeholder="Masukkan alamat sekolah"
            :rows="3"
          />
        </FormItem>
        <Row :gutter="16">
          <Col :span="12">
            <FormItem label="Telepon" name="phone">
              <Input v-model:value="formState.phone" placeholder="Contoh: 021-1234567" />
            </FormItem>
          </Col>
          <Col :span="12">
            <FormItem label="Email" name="email">
              <Input v-model:value="formState.email" placeholder="email@sekolah.sch.id" />
            </FormItem>
          </Col>
        </Row>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.tenant-management {
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
