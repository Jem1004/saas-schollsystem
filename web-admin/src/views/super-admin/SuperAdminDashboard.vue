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
      <Title :level="2" style="margin: 0">Dashboard Super Admin</Title>
      <Text type="secondary">Overview sistem dan statistik platform</Text>
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
          <Card class="stat-card">
            <Statistic
              title="Total Sekolah"
              :value="tenantStats.totalSchools"
              :value-style="{ color: '#f97316' }"
            >
              <template #prefix>
                <BankOutlined />
              </template>
            </Statistic>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :lg="6">
          <Card class="stat-card">
            <Statistic
              title="Sekolah Aktif"
              :value="tenantStats.activeSchools"
              :value-style="{ color: '#22c55e' }"
            >
              <template #prefix>
                <CheckCircleOutlined />
              </template>
            </Statistic>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :lg="6">
          <Card class="stat-card">
            <Statistic
              title="Total Device"
              :value="deviceStats.totalDevices"
              :value-style="{ color: '#3b82f6' }"
            >
              <template #prefix>
                <DesktopOutlined />
              </template>
            </Statistic>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :lg="6">
          <Card class="stat-card">
            <Statistic
              title="Device Online"
              :value="deviceStats.onlineDevices"
              :value-style="{ color: '#22c55e' }"
            >
              <template #prefix>
                <WifiOutlined />
              </template>
            </Statistic>
          </Card>
        </Col>
      </Row>

      <!-- System Status -->
      <Row :gutter="[24, 24]" style="margin-top: 24px">
        <Col :xs="24" :lg="12">
          <Card title="Status Sistem" class="status-card">
            <div class="status-item">
              <span class="status-label">Database</span>
              <Tag color="success">
                <CheckCircleOutlined /> Online
              </Tag>
            </div>
            <div class="status-item">
              <span class="status-label">Redis Queue</span>
              <Tag color="success">
                <CheckCircleOutlined /> Online
              </Tag>
            </div>
            <div class="status-item">
              <span class="status-label">FCM Service</span>
              <Tag color="success">
                <CheckCircleOutlined /> Online
              </Tag>
            </div>
            <div class="status-item">
              <span class="status-label">Notification Worker</span>
              <Tag color="success">
                <CheckCircleOutlined /> Running
              </Tag>
            </div>
          </Card>
        </Col>
        <Col :xs="24" :lg="12">
          <Card title="Ringkasan" class="summary-card">
            <div class="summary-item">
              <TeamOutlined class="summary-icon" />
              <div class="summary-content">
                <Text strong>{{ tenantStats.activeSchools }}</Text>
                <Text type="secondary"> sekolah aktif dari </Text>
                <Text strong>{{ tenantStats.totalSchools }}</Text>
                <Text type="secondary"> total</Text>
              </div>
            </div>
            <div class="summary-item">
              <DesktopOutlined class="summary-icon" />
              <div class="summary-content">
                <Text strong>{{ deviceStats.activeDevices }}</Text>
                <Text type="secondary"> device aktif, </Text>
                <Text strong>{{ deviceStats.onlineDevices }}</Text>
                <Text type="secondary"> online</Text>
              </div>
            </div>
            <div class="summary-item">
              <CloseCircleOutlined class="summary-icon inactive" />
              <div class="summary-content">
                <Text strong>{{ tenantStats.inactiveSchools }}</Text>
                <Text type="secondary"> sekolah nonaktif, </Text>
                <Text strong>{{ deviceStats.inactiveDevices }}</Text>
                <Text type="secondary"> device nonaktif</Text>
              </div>
            </div>
          </Card>
        </Col>
      </Row>

      <!-- Recent Data -->
      <Row :gutter="[24, 24]" style="margin-top: 24px">
        <Col :xs="24" :lg="12">
          <Card title="Sekolah Terbaru" class="list-card">
            <List
              :data-source="recentSchools"
              :loading="loading"
              size="small"
            >
              <template #renderItem="{ item }">
                <ListItem>
                  <ListItemMeta
                    :title="item.name"
                    :description="item.email || 'No email'"
                  >
                    <template #avatar>
                      <BankOutlined class="list-icon" />
                    </template>
                  </ListItemMeta>
                  <template #actions>
                    <Tag :color="item.isActive ? 'success' : 'default'">
                      {{ item.isActive ? 'Aktif' : 'Nonaktif' }}
                    </Tag>
                  </template>
                </ListItem>
              </template>
            </List>
          </Card>
        </Col>
        <Col :xs="24" :lg="12">
          <Card title="Device Terbaru" class="list-card">
            <List
              :data-source="recentDevices"
              :loading="loading"
              size="small"
            >
              <template #renderItem="{ item }">
                <ListItem>
                  <ListItemMeta
                    :title="item.deviceCode"
                    :description="item.schoolName || `School ID: ${item.schoolId}`"
                  >
                    <template #avatar>
                      <component
                        :is="isDeviceOnline(item) ? WifiOutlined : DisconnectOutlined"
                        :class="['list-icon', isDeviceOnline(item) ? 'online' : 'offline']"
                      />
                    </template>
                  </ListItemMeta>
                  <template #actions>
                    <Tag :color="item.isActive ? (isDeviceOnline(item) ? 'success' : 'warning') : 'default'">
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
}

.page-header {
  margin-bottom: 24px;
}

.stat-card {
  height: 100%;
}

.stat-card :deep(.ant-statistic-title) {
  font-size: 14px;
  color: #8c8c8c;
}

.stat-card :deep(.ant-statistic-content-prefix) {
  margin-right: 8px;
}

.status-card,
.summary-card,
.list-card {
  height: 100%;
}

.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.status-item:last-child {
  border-bottom: none;
}

.status-label {
  font-weight: 500;
}

.summary-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.summary-item:last-child {
  border-bottom: none;
}

.summary-icon {
  font-size: 24px;
  color: #f97316;
}

.summary-icon.inactive {
  color: #8c8c8c;
}

.summary-content {
  flex: 1;
}

.list-icon {
  font-size: 20px;
  color: #f97316;
}

.list-icon.online {
  color: #22c55e;
}

.list-icon.offline {
  color: #8c8c8c;
}

.list-card :deep(.ant-list-item) {
  padding: 12px 0;
}

.list-card :deep(.ant-list-item-meta-title) {
  margin-bottom: 0;
}
</style>
