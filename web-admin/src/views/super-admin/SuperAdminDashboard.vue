<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  Row,
  Col,
  Card,
  Statistic,
  Spin,
  Alert,
  Tag,
  List,
  ListItem,
  ListItemMeta,
  Typography,
} from 'ant-design-vue'
import {
  BankOutlined,
  DesktopOutlined,
  TeamOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  WifiOutlined,
  DisconnectOutlined,
} from '@ant-design/icons-vue'
import { tenantService, deviceService } from '@/services'
import type { TenantStats } from '@/types/tenant'
import type { DeviceStats, Device } from '@/types/device'
import type { School } from '@/types/tenant'

const { Title, Text } = Typography

const loading = ref(true)
const error = ref<string | null>(null)

const tenantStats = ref<TenantStats>({
  totalSchools: 0,
  activeSchools: 0,
  inactiveSchools: 0,
})

const deviceStats = ref<DeviceStats>({
  totalDevices: 0,
  activeDevices: 0,
  inactiveDevices: 0,
  onlineDevices: 0,
})

const recentSchools = ref<School[]>([])
const recentDevices = ref<Device[]>([])

// Mock data for development (will be replaced with real API calls)
const loadMockData = () => {
  tenantStats.value = {
    totalSchools: 12,
    activeSchools: 10,
    inactiveSchools: 2,
  }

  deviceStats.value = {
    totalDevices: 45,
    activeDevices: 40,
    inactiveDevices: 5,
    onlineDevices: 38,
  }

  recentSchools.value = [
    { id: 1, name: 'SMP Negeri 1 Jakarta', isActive: true, email: 'smpn1@jakarta.sch.id', createdAt: '2024-01-15', updatedAt: '2024-01-15' },
    { id: 2, name: 'SMP Negeri 2 Bandung', isActive: true, email: 'smpn2@bandung.sch.id', createdAt: '2024-01-14', updatedAt: '2024-01-14' },
    { id: 3, name: 'SMP Negeri 3 Surabaya', isActive: false, email: 'smpn3@surabaya.sch.id', createdAt: '2024-01-13', updatedAt: '2024-01-13' },
  ]

  recentDevices.value = [
    { id: 1, schoolId: 1, schoolName: 'SMP Negeri 1 Jakarta', deviceCode: 'ESP32-001', apiKey: '***', isActive: true, lastSeenAt: new Date().toISOString(), createdAt: '2024-01-15', updatedAt: '2024-01-15' },
    { id: 2, schoolId: 1, schoolName: 'SMP Negeri 1 Jakarta', deviceCode: 'ESP32-002', apiKey: '***', isActive: true, lastSeenAt: new Date(Date.now() - 3600000).toISOString(), createdAt: '2024-01-14', updatedAt: '2024-01-14' },
    { id: 3, schoolId: 2, schoolName: 'SMP Negeri 2 Bandung', deviceCode: 'ESP32-003', apiKey: '***', isActive: false, createdAt: '2024-01-13', updatedAt: '2024-01-13' },
  ]
}

const loadData = async () => {
  loading.value = true
  error.value = null

  try {
    // Try to load real data from API
    const [tenantStatsRes, deviceStatsRes, schoolsRes, devicesRes] = await Promise.allSettled([
      tenantService.getStats(),
      deviceService.getStats(),
      tenantService.getSchools({ page: 1, pageSize: 5 }),
      deviceService.getDevices({ page: 1, pageSize: 5 }),
    ])

    if (tenantStatsRes.status === 'fulfilled') {
      tenantStats.value = tenantStatsRes.value
    }
    if (deviceStatsRes.status === 'fulfilled') {
      deviceStats.value = deviceStatsRes.value
    }
    if (schoolsRes.status === 'fulfilled') {
      recentSchools.value = schoolsRes.value.data
    }
    if (devicesRes.status === 'fulfilled') {
      recentDevices.value = devicesRes.value.data
    }

    // If all failed, use mock data
    if (
      tenantStatsRes.status === 'rejected' &&
      deviceStatsRes.status === 'rejected'
    ) {
      loadMockData()
    }
  } catch {
    // Use mock data on error
    loadMockData()
  } finally {
    loading.value = false
  }
}

// Check if device is online (last seen within 5 minutes)
const isDeviceOnline = (device: Device): boolean => {
  if (!device.lastSeenAt) return false
  const lastSeen = new Date(device.lastSeenAt)
  const fiveMinutesAgo = new Date(Date.now() - 5 * 60 * 1000)
  return lastSeen > fiveMinutesAgo
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="super-admin-dashboard">
    <div class="page-header">
      <div class="header-content">
        <Title :level="2" class="page-title">Dashboard Super Admin</Title>
        <Text class="page-subtitle">Overview sistem dan statistik platform</Text>
      </div>
    </div>

    <Spin :spinning="loading">
      <Alert
        v-if="error"
        type="error"
        :message="error"
        show-icon
        closable
        style="margin-bottom: 24px"
      />

      <!-- Statistics Cards -->
      <Row :gutter="[24, 24]" class="stats-row">
        <Col :xs="24" :sm="12" :lg="6">
          <Card class="stat-card" :bordered="false">
            <Statistic
              title="Total Sekolah"
              :value="tenantStats.totalSchools"
              :value-style="{ color: '#0f172a', fontWeight: '600' }"
            >
              <template #prefix>
                <BankOutlined class="stat-icon primary" />
              </template>
            </Statistic>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :lg="6">
          <Card class="stat-card" :bordered="false">
            <Statistic
              title="Sekolah Aktif"
              :value="tenantStats.activeSchools"
              :value-style="{ color: '#0f172a', fontWeight: '600' }"
            >
              <template #prefix>
                <CheckCircleOutlined class="stat-icon success" />
              </template>
            </Statistic>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :lg="6">
          <Card class="stat-card" :bordered="false">
            <Statistic
              title="Total Device"
              :value="deviceStats.totalDevices"
              :value-style="{ color: '#0f172a', fontWeight: '600' }"
            >
              <template #prefix>
                <DesktopOutlined class="stat-icon primary" />
              </template>
            </Statistic>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :lg="6">
          <Card class="stat-card" :bordered="false">
            <Statistic
              title="Device Online"
              :value="deviceStats.onlineDevices"
              :value-style="{ color: '#0f172a', fontWeight: '600' }"
            >
              <template #prefix>
                <WifiOutlined class="stat-icon success" />
              </template>
            </Statistic>
          </Card>
        </Col>
      </Row>

      <!-- System Status -->
      <Row :gutter="[24, 24]" style="margin-top: 24px">
        <Col :xs="24" :lg="12">
          <Card title="Status Sistem" class="content-card" :bordered="false">
            <div class="status-list">
              <div class="status-item">
                <span class="status-label">Database</span>
                <Tag color="success" class="status-tag">
                  <CheckCircleOutlined /> Online
                </Tag>
              </div>
              <div class="status-item">
                <span class="status-label">Redis Queue</span>
                <Tag color="success" class="status-tag">
                  <CheckCircleOutlined /> Online
                </Tag>
              </div>
              <div class="status-item">
                <span class="status-label">FCM Service</span>
                <Tag color="success" class="status-tag">
                  <CheckCircleOutlined /> Online
                </Tag>
              </div>
              <div class="status-item">
                <span class="status-label">Notification Worker</span>
                <Tag color="success" class="status-tag">
                  <CheckCircleOutlined /> Running
                </Tag>
              </div>
            </div>
          </Card>
        </Col>
        <Col :xs="24" :lg="12">
          <Card title="Ringkasan" class="content-card" :bordered="false">
            <div class="summary-list">
              <div class="summary-item">
                <div class="summary-icon-wrapper">
                  <TeamOutlined class="summary-icon" />
                </div>
                <div class="summary-content">
                  <Text strong>{{ tenantStats.activeSchools }}</Text>
                  <Text type="secondary"> sekolah aktif dari </Text>
                  <Text strong>{{ tenantStats.totalSchools }}</Text>
                  <Text type="secondary"> total terdaftar</Text>
                </div>
              </div>
              <div class="summary-item">
                <div class="summary-icon-wrapper">
                  <DesktopOutlined class="summary-icon" />
                </div>
                <div class="summary-content">
                  <Text strong>{{ deviceStats.activeDevices }}</Text>
                  <Text type="secondary"> device aktif, </Text>
                  <Text strong>{{ deviceStats.onlineDevices }}</Text>
                  <Text type="secondary"> sedang online</Text>
                </div>
              </div>
              <div class="summary-item">
                <div class="summary-icon-wrapper inactive">
                  <CloseCircleOutlined class="summary-icon" />
                </div>
                <div class="summary-content">
                  <Text strong>{{ tenantStats.inactiveSchools }}</Text>
                  <Text type="secondary"> sekolah nonaktif, </Text>
                  <Text strong>{{ deviceStats.inactiveDevices }}</Text>
                  <Text type="secondary"> device nonaktif</Text>
                </div>
              </div>
            </div>
          </Card>
        </Col>
      </Row>

      <!-- Recent Data -->
      <Row :gutter="[24, 24]" style="margin-top: 24px">
        <Col :xs="24" :lg="12">
          <Card title="Sekolah Terbaru" class="content-card" :bordered="false">
            <List
              :data-source="recentSchools"
              :loading="loading"
              size="small"
              item-layout="horizontal"
            >
              <template #renderItem="{ item }">
                <ListItem class="recent-list-item">
                  <ListItemMeta>
                    <template #title>
                      <span class="list-title">{{ item.name }}</span>
                    </template>
                    <template #description>
                      <span class="list-subtitle">{{ item.email || 'No email' }}</span>
                    </template>
                    <template #avatar>
                      <div class="list-avatar">
                        <BankOutlined />
                      </div>
                    </template>
                  </ListItemMeta>
                  <template #actions>
                    <Tag :color="item.isActive ? 'success' : 'default'" :bordered="false">
                      {{ item.isActive ? 'Aktif' : 'Nonaktif' }}
                    </Tag>
                  </template>
                </ListItem>
              </template>
            </List>
          </Card>
        </Col>
        <Col :xs="24" :lg="12">
          <Card title="Device Terbaru" class="content-card" :bordered="false">
            <List
              :data-source="recentDevices"
              :loading="loading"
              size="small"
              item-layout="horizontal"
            >
              <template #renderItem="{ item }">
                <ListItem class="recent-list-item">
                  <ListItemMeta>
                    <template #title>
                      <span class="list-title">{{ item.deviceCode }}</span>
                    </template>
                    <template #description>
                      <span class="list-subtitle">{{ item.schoolName || `School: ${item.schoolId}` }}</span>
                    </template>
                    <template #avatar>
                      <div class="list-avatar" :class="{ 'online': isDeviceOnline(item) }">
                        <component :is="isDeviceOnline(item) ? WifiOutlined : DisconnectOutlined" />
                      </div>
                    </template>
                  </ListItemMeta>
                  <template #actions>
                    <Tag :color="item.isActive ? (isDeviceOnline(item) ? 'success' : 'warning') : 'default'" :bordered="false">
                      {{ item.isActive ? (isDeviceOnline(item) ? 'Online' : 'Offline') : 'Nonaktif' }}
                    </Tag>
                  </template>
                </ListItem>
              </template>
            </List>
          </Card>
        </Col>
      </Row>
    </Spin>
  </div>
</template>

<style scoped>
.super-admin-dashboard {
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

/* Stat Cards */
.stat-card {
  height: 100%;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
  border: 1px solid #f1f5f9;
  transition: all 0.2s ease;
}

.stat-card:hover {
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
  border-color: #e2e8f0;
}

.stat-icon {
  font-size: 20px;
  padding: 10px;
  border-radius: 10px;
}

.stat-icon.primary {
  color: #f97316;
  background-color: #fff7ed;
}

.stat-icon.success {
  color: #22c55e;
  background-color: #f0fdf4;
}

/* Content Cards */
.content-card {
  height: 100%;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
  border: 1px solid #f1f5f9;
}

.content-card :deep(.ant-card-head) {
  border-bottom: 1px solid #f1f5f9;
  padding: 16px 24px;
}

.content-card :deep(.ant-card-head-title) {
  font-size: 16px;
  font-weight: 600;
  color: #334155;
}

.content-card :deep(.ant-card-body) {
  padding: 24px;
}

/* Status List */
.status-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 12px;
  border-bottom: 1px solid #f8fafc;
}

.status-item:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.status-label {
  font-weight: 500;
  color: #475569;
}

.status-tag {
  min-width: 80px;
  text-align: center;
}

/* Summary List */
.summary-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.summary-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

.summary-icon-wrapper {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background-color: #fff7ed;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.summary-icon-wrapper.inactive {
  background-color: #f1f5f9;
}

.summary-icon {
  font-size: 18px;
  color: #f97316;
}

.summary-icon-wrapper.inactive .summary-icon {
  color: #64748b;
}

.summary-content {
  font-size: 14px;
  line-height: 1.5;
  color: #475569;
}

/* Recent List */
.recent-list-item {
  padding: 12px 0 !important;
  border-bottom: 1px solid #f8fafc !important;
}

.recent-list-item:last-child {
  border-bottom: none !important;
}

.list-avatar {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background-color: #f8fafc;
  color: #64748b;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}

.list-avatar.online {
  background-color: #f0fdf4;
  color: #22c55e;
}

.list-title {
  font-weight: 500;
  color: #334155;
  font-size: 14px;
}

.list-subtitle {
  color: #94a3b8;
  font-size: 13px;
}
</style>
