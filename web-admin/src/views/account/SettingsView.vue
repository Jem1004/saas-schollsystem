<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  Card,
  Row,
  Col,
  Typography,
  Switch,
  List,
  ListItem,
  ListItemMeta,
  Button,
  message,
  Divider,
  Alert,
} from 'ant-design-vue'
import {
  BellOutlined,
  SecurityScanOutlined,
  GlobalOutlined,
  MobileOutlined,
  MailOutlined,
} from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/auth'

const { Title, Text } = Typography

const authStore = useAuthStore()

// Settings state (stored in localStorage for now)
const notificationEnabled = ref(localStorage.getItem('notificationEnabled') !== 'false')
const emailNotification = ref(localStorage.getItem('emailNotification') !== 'false')
const soundEnabled = ref(localStorage.getItem('soundEnabled') !== 'false')

const isSuperAdmin = computed(() => authStore.userRole === 'super_admin')

const handleNotificationToggle = (checked: boolean) => {
  notificationEnabled.value = checked
  localStorage.setItem('notificationEnabled', String(checked))
  message.success(checked ? 'Notifikasi diaktifkan' : 'Notifikasi dinonaktifkan')
}

const handleEmailToggle = (checked: boolean) => {
  emailNotification.value = checked
  localStorage.setItem('emailNotification', String(checked))
  message.success(checked ? 'Notifikasi email diaktifkan' : 'Notifikasi email dinonaktifkan')
}

const handleSoundToggle = (checked: boolean) => {
  soundEnabled.value = checked
  localStorage.setItem('soundEnabled', String(checked))
  message.success(checked ? 'Suara notifikasi diaktifkan' : 'Suara notifikasi dinonaktifkan')
}

const handleClearCache = () => {
  // Clear only non-essential cache
  const keysToKeep = ['accessToken', 'refreshToken', 'user']
  const allKeys = Object.keys(localStorage)
  
  allKeys.forEach(key => {
    if (!keysToKeep.includes(key)) {
      localStorage.removeItem(key)
    }
  })
  
  message.success('Cache berhasil dibersihkan')
}
</script>

<template>
  <div class="settings-view">
    <div class="page-header">
      <Title :level="2" style="margin: 0">Pengaturan</Title>
      <Text type="secondary">Kelola preferensi dan pengaturan akun Anda</Text>
    </div>

    <Row :gutter="24">
      <Col :xs="24" :lg="16">
        <!-- Notification Settings -->
        <Card title="Notifikasi" class="settings-card" v-if="!isSuperAdmin">
          <template #extra>
            <BellOutlined />
          </template>
          
          <List item-layout="horizontal">
            <ListItem>
              <ListItemMeta
                title="Notifikasi Push"
                description="Terima notifikasi langsung di browser"
              />
              <template #actions>
                <Switch :checked="notificationEnabled" @change="handleNotificationToggle" />
              </template>
            </ListItem>
            <ListItem>
              <ListItemMeta
                title="Notifikasi Email"
                description="Terima notifikasi penting via email"
              />
              <template #actions>
                <Switch :checked="emailNotification" @change="handleEmailToggle" />
              </template>
            </ListItem>
            <ListItem>
              <ListItemMeta
                title="Suara Notifikasi"
                description="Mainkan suara saat ada notifikasi baru"
              />
              <template #actions>
                <Switch :checked="soundEnabled" @change="handleSoundToggle" />
              </template>
            </ListItem>
          </List>
        </Card>

        <Alert 
          v-if="isSuperAdmin" 
          type="info" 
          show-icon 
          message="Notifikasi tidak tersedia untuk Super Admin"
          description="Super Admin tidak menerima notifikasi sistem karena tidak terkait dengan sekolah tertentu."
          style="margin-bottom: 24px"
        />

        <!-- Security Settings -->
        <Card title="Keamanan" class="settings-card">
          <template #extra>
            <SecurityScanOutlined />
          </template>
          
          <List item-layout="horizontal">
            <ListItem>
              <ListItemMeta
                title="Ubah Password"
                description="Perbarui password akun Anda secara berkala"
              />
              <template #actions>
                <Button type="link" @click="$router.push('/profile')">Ubah</Button>
              </template>
            </ListItem>
            <ListItem>
              <ListItemMeta
                title="Sesi Aktif"
                description="Anda sedang login di perangkat ini"
              />
              <template #actions>
                <Text type="secondary">1 perangkat</Text>
              </template>
            </ListItem>
          </List>
        </Card>

        <!-- System Settings -->
        <Card title="Sistem" class="settings-card">
          <template #extra>
            <GlobalOutlined />
          </template>
          
          <List item-layout="horizontal">
            <ListItem>
              <ListItemMeta
                title="Bahasa"
                description="Bahasa tampilan aplikasi"
              />
              <template #actions>
                <Text>Bahasa Indonesia</Text>
              </template>
            </ListItem>
            <ListItem>
              <ListItemMeta
                title="Bersihkan Cache"
                description="Hapus data cache lokal untuk memperbaiki masalah"
              />
              <template #actions>
                <Button type="link" danger @click="handleClearCache">Bersihkan</Button>
              </template>
            </ListItem>
          </List>
        </Card>
      </Col>

      <Col :xs="24" :lg="8">
        <Card title="Informasi Aplikasi">
          <div class="app-info">
            <div class="info-item">
              <Text type="secondary">Versi</Text>
              <Text strong>1.0.0</Text>
            </div>
            <Divider style="margin: 12px 0" />
            <div class="info-item">
              <Text type="secondary">Terakhir Diperbarui</Text>
              <Text strong>1 Januari 2026</Text>
            </div>
            <Divider style="margin: 12px 0" />
            <div class="info-item">
              <Text type="secondary">Dukungan</Text>
              <a href="mailto:support@schooladmin.com">support@schooladmin.com</a>
            </div>
          </div>
        </Card>
      </Col>
    </Row>
  </div>
</template>

<style scoped>
.settings-view {
  padding: 0;
}

.page-header {
  margin-bottom: 24px;
}

.settings-card {
  margin-bottom: 24px;
}

.app-info {
  display: flex;
  flex-direction: column;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
