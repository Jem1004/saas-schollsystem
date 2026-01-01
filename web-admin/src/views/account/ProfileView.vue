<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  Card,
  Row,
  Col,
  Typography,
  Descriptions,
  DescriptionsItem,
  Avatar,
  Tag,
  Button,
  Modal,
  Form,
  FormItem,
  Input,
  message,
  Divider,
} from 'ant-design-vue'
import {
  UserOutlined,
  MailOutlined,
  CalendarOutlined,
  SafetyOutlined,
  LockOutlined,
} from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { getRoleDisplayName, getRoleColor } from '@/composables/usePermissions'

const { Title, Text } = Typography

const authStore = useAuthStore()

// Change password modal
const passwordModalVisible = ref(false)
const passwordLoading = ref(false)
const passwordFormRef = ref()
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const passwordRules = {
  oldPassword: [{ required: true, message: 'Password lama wajib diisi' }],
  newPassword: [
    { required: true, message: 'Password baru wajib diisi' },
    { min: 6, message: 'Password minimal 6 karakter' },
  ],
  confirmPassword: [
    { required: true, message: 'Konfirmasi password wajib diisi' },
    {
      validator: (_rule: unknown, value: string) => {
        if (value && value !== passwordForm.value.newPassword) {
          return Promise.reject('Password tidak cocok')
        }
        return Promise.resolve()
      },
    },
  ],
}

const user = computed(() => authStore.user)
const roleDisplayName = computed(() => user.value ? getRoleDisplayName(user.value.role) : '-')
const roleColor = computed(() => user.value ? getRoleColor(user.value.role) : 'default')

const formatDate = (dateStr?: string): string => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('id-ID', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const openPasswordModal = () => {
  passwordForm.value = { oldPassword: '', newPassword: '', confirmPassword: '' }
  passwordModalVisible.value = true
}

const handleChangePassword = async () => {
  try {
    await passwordFormRef.value?.validate()
  } catch {
    return
  }

  passwordLoading.value = true
  try {
    await authStore.changePassword(passwordForm.value.oldPassword, passwordForm.value.newPassword)
    message.success('Password berhasil diubah')
    passwordModalVisible.value = false
  } catch (error: unknown) {
    const err = error as { response?: { data?: { error?: { message?: string } } } }
    message.error(err.response?.data?.error?.message || 'Gagal mengubah password')
  } finally {
    passwordLoading.value = false
  }
}
</script>

<template>
  <div class="profile-view">
    <div class="page-header">
      <Title :level="2" style="margin: 0">Profil Saya</Title>
      <Text type="secondary">Informasi akun dan pengaturan profil</Text>
    </div>

    <Row :gutter="24">
      <Col :xs="24" :md="8">
        <Card class="profile-card">
          <div class="profile-avatar-section">
            <Avatar :size="100" class="profile-avatar">
              <template #icon><UserOutlined /></template>
            </Avatar>
            <Title :level="4" style="margin: 16px 0 4px">{{ user?.username || '-' }}</Title>
            <Tag :color="roleColor">{{ roleDisplayName }}</Tag>
          </div>
          
          <Divider />
          
          <div class="profile-actions">
            <Button block @click="openPasswordModal">
              <template #icon><LockOutlined /></template>
              Ubah Password
            </Button>
          </div>
        </Card>
      </Col>

      <Col :xs="24" :md="16">
        <Card title="Informasi Akun">
          <Descriptions :column="1" bordered>
            <DescriptionsItem label="Username">
              <UserOutlined style="margin-right: 8px; color: #8c8c8c" />
              {{ user?.username || '-' }}
            </DescriptionsItem>
            <DescriptionsItem label="Email">
              <MailOutlined style="margin-right: 8px; color: #8c8c8c" />
              {{ user?.email || '-' }}
            </DescriptionsItem>
            <DescriptionsItem label="Role">
              <SafetyOutlined style="margin-right: 8px; color: #8c8c8c" />
              <Tag :color="roleColor">{{ roleDisplayName }}</Tag>
            </DescriptionsItem>
            <DescriptionsItem label="Status">
              <Tag :color="user?.isActive ? 'success' : 'default'">
                {{ user?.isActive ? 'Aktif' : 'Nonaktif' }}
              </Tag>
            </DescriptionsItem>
            <DescriptionsItem label="Login Terakhir">
              <CalendarOutlined style="margin-right: 8px; color: #8c8c8c" />
              {{ formatDate(user?.lastLoginAt) }}
            </DescriptionsItem>
            <DescriptionsItem label="Terdaftar Sejak">
              <CalendarOutlined style="margin-right: 8px; color: #8c8c8c" />
              {{ formatDate(user?.createdAt) }}
            </DescriptionsItem>
          </Descriptions>
        </Card>
      </Col>
    </Row>

    <!-- Change Password Modal -->
    <Modal
      v-model:open="passwordModalVisible"
      title="Ubah Password"
      :confirm-loading="passwordLoading"
      ok-text="Simpan"
      cancel-text="Batal"
      @ok="handleChangePassword"
    >
      <Form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        layout="vertical"
        style="margin-top: 16px"
      >
        <FormItem label="Password Lama" name="oldPassword">
          <Input.Password v-model:value="passwordForm.oldPassword" placeholder="Masukkan password lama" />
        </FormItem>
        <FormItem label="Password Baru" name="newPassword">
          <Input.Password v-model:value="passwordForm.newPassword" placeholder="Masukkan password baru" />
        </FormItem>
        <FormItem label="Konfirmasi Password" name="confirmPassword">
          <Input.Password v-model:value="passwordForm.confirmPassword" placeholder="Konfirmasi password baru" />
        </FormItem>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.profile-view {
  padding: 0;
}

.page-header {
  margin-bottom: 24px;
}

.profile-card {
  text-align: center;
}

.profile-avatar-section {
  padding: 24px 0;
}

.profile-avatar {
  background-color: #f97316;
}

.profile-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
