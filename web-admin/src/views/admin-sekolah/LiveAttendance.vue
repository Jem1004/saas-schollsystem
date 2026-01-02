<script setup lang="ts">
// Live Attendance Dashboard
// Requirements: 4.1 - Display current day's attendance statistics
// Requirements: 4.2 - Update dashboard within 3 seconds without page refresh
// Requirements: 4.3 - Show the 20 most recent attendance records
// Requirements: 4.4 - Add new attendance to top of live feed with visual highlight
// Requirements: 4.5 - Allow filtering by class
// Requirements: 4.9 - Display connection status and attempt to reconnect
// Requirements: 4.10 - Show percentage of attendance completion

import { ref, onMounted, onUnmounted, computed, watch, nextTick } from 'vue'
import {
  Card,
  Row,
  Col,
  Typography,
  Select,
  SelectOption,
  Statistic,
  Tag,
  Badge,
  List,
  ListItem,
  ListItemMeta,
  Avatar,
  Empty,
  Spin,
  Alert,
  Progress,
  Button,
  message,
} from 'ant-design-vue'
import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  ClockCircleOutlined,
  ExclamationCircleOutlined,
  WifiOutlined,
  DisconnectOutlined,
  ReloadOutlined,
  UserOutlined,
  SyncOutlined,
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import 'dayjs/locale/id'
import { realtimeService, schoolService } from '@/services'
import type { 
  AttendanceStats, 
  LiveFeedEntry, 
  AttendanceEvent,
  WSConnectionStatus 
} from '@/types/realtime'
import type { Class } from '@/types/school'

dayjs.locale('id')

const { Title, Text } = Typography

// State
const loading = ref(true)
const stats = ref<AttendanceStats>({
  totalStudents: 0,
  present: 0,
  late: 0,
  veryLate: 0,
  absent: 0,
  percentage: 0,
})
const liveFeed = ref<LiveFeedEntry[]>([])
const filterClassId = ref<number | undefined>(undefined)
const connectionStatus = ref<WSConnectionStatus>({ connected: false, message: 'Menghubungkan...' })
const highlightedId = ref<number | null>(null)

// Classes for filter
const classes = ref<Class[]>([])
const loadingClasses = ref(false)

// Feed container ref for auto-scroll
const feedContainer = ref<HTMLElement | null>(null)

// Computed values
const formattedDate = computed(() => {
  return dayjs().format('dddd, D MMMM YYYY')
})

const connectionStatusColor = computed(() => {
  return connectionStatus.value.connected ? 'success' : 'error'
})

const connectionIcon = computed(() => {
  return connectionStatus.value.connected ? WifiOutlined : DisconnectOutlined
})

// Get status color for attendance
const getStatusColor = (status: string): string => {
  switch (status) {
    case 'present':
      return '#22c55e'
    case 'late':
      return '#f97316'
    case 'very_late':
      return '#ef4444'
    case 'absent':
      return '#6b7280'
    default:
      return '#3b82f6'
  }
}

// Get status label
const getStatusLabel = (status: string): string => {
  switch (status) {
    case 'present':
      return 'Hadir'
    case 'late':
      return 'Terlambat'
    case 'very_late':
      return 'Sangat Terlambat'
    case 'absent':
      return 'Tidak Hadir'
    default:
      return status
  }
}

// Get status tag color
const getStatusTagColor = (status: string): string => {
  switch (status) {
    case 'present':
      return 'success'
    case 'late':
      return 'warning'
    case 'very_late':
      return 'error'
    case 'absent':
      return 'default'
    default:
      return 'processing'
  }
}

// Format time for display
const formatTime = (timeStr: string): string => {
  return dayjs(timeStr).format('HH:mm:ss')
}

// Format relative time
const formatRelativeTime = (timeStr: string): string => {
  const time = dayjs(timeStr)
  const now = dayjs()
  const diffSeconds = now.diff(time, 'second')
  
  if (diffSeconds < 60) {
    return 'Baru saja'
  } else if (diffSeconds < 3600) {
    const mins = Math.floor(diffSeconds / 60)
    return `${mins} menit lalu`
  } else {
    return time.format('HH:mm')
  }
}

// Load initial data via REST API
const loadInitialData = async () => {
  loading.value = true
  try {
    const [feedData, statsData] = await Promise.all([
      realtimeService.getLiveFeed(filterClassId.value),
      realtimeService.getStats(filterClassId.value),
    ])
    liveFeed.value = feedData
    stats.value = statsData
  } catch (error) {
    console.error('Failed to load initial data:', error)
    message.error('Gagal memuat data absensi')
  } finally {
    loading.value = false
  }
}

// Load classes for filter
const loadClasses = async () => {
  loadingClasses.value = true
  try {
    const response = await schoolService.getClasses({ page_size: 100 })
    classes.value = response.classes
  } catch (error) {
    console.error('Failed to load classes:', error)
    classes.value = []
  } finally {
    loadingClasses.value = false
  }
}

// Handle WebSocket attendance event
// Requirements: 4.2 - Update dashboard within 3 seconds
// Requirements: 4.4 - Add new attendance to top with visual highlight
const handleAttendanceEvent = (event: AttendanceEvent) => {
  if (event.type === 'new_attendance' && event.attendance) {
    // Add to top of feed with highlight
    // Requirements: 4.3 - Show the 20 most recent records
    liveFeed.value = [event.attendance, ...liveFeed.value.slice(0, 19)]
    
    // Highlight the new entry
    highlightedId.value = event.attendance.id
    setTimeout(() => {
      highlightedId.value = null
    }, 3000)

    // Auto-scroll to top
    nextTick(() => {
      if (feedContainer.value) {
        feedContainer.value.scrollTop = 0
      }
    })
  }

  if (event.type === 'stats_update' && event.stats) {
    stats.value = event.stats
  }
}

// Handle connection status change
// Requirements: 4.9 - Display connection status
const handleConnectionStatus = (status: WSConnectionStatus) => {
  connectionStatus.value = status
}

// Handle class filter change
// Requirements: 4.5 - Allow filtering by class
const handleClassFilterChange = () => {
  loadInitialData()
  realtimeService.setClassFilter(filterClassId.value)
}

// Manual reconnect
const handleReconnect = () => {
  realtimeService.disconnect()
  realtimeService.connect(filterClassId.value)
}

// Manual refresh
const handleRefresh = () => {
  loadInitialData()
}

// Setup WebSocket connection and event handlers
onMounted(async () => {
  // Load classes first
  await loadClasses()
  
  // Load initial data
  await loadInitialData()
  
  // Subscribe to events
  realtimeService.onAttendanceEvent(handleAttendanceEvent)
  realtimeService.onConnectionStatus(handleConnectionStatus)
  
  // Connect to WebSocket
  realtimeService.connect(filterClassId.value)
})

// Cleanup on unmount
onUnmounted(() => {
  realtimeService.disconnect()
})

// Watch for class filter changes
watch(filterClassId, () => {
  handleClassFilterChange()
})
</script>

<template>
  <div class="live-attendance">
    <div class="page-header">
      <div class="header-left">
        <Title :level="2" style="margin: 0">Absensi Real-Time</Title>
        <Text type="secondary">{{ formattedDate }}</Text>
      </div>
      <div class="header-right">
        <!-- Connection Status Indicator -->
        <!-- Requirements: 4.9 - Display connection status -->
        <Badge :status="connectionStatusColor" :text="connectionStatus.message" />
        <Button 
          v-if="!connectionStatus.connected" 
          type="link" 
          size="small"
          @click="handleReconnect"
        >
          <template #icon><ReloadOutlined /></template>
          Hubungkan Ulang
        </Button>
      </div>
    </div>

    <!-- Connection Alert -->
    <Alert
      v-if="!connectionStatus.connected"
      type="warning"
      show-icon
      style="margin-bottom: 16px"
    >
      <template #message>
        <span>
          <component :is="connectionIcon" style="margin-right: 8px" />
          {{ connectionStatus.message }}
        </span>
      </template>
      <template #description>
        Data mungkin tidak diperbarui secara real-time. Klik "Hubungkan Ulang" atau refresh halaman.
      </template>
    </Alert>

    <!-- Stats Cards -->
    <!-- Requirements: 4.1 - Display current day's attendance statistics -->
    <Row :gutter="[16, 16]" class="stats-row">
      <Col :xs="12" :sm="12" :md="6" :lg="4">
        <Card class="stat-card">
          <Statistic
            title="Total Siswa"
            :value="stats.totalStudents"
            :value-style="{ color: '#3b82f6' }"
          >
            <template #prefix>
              <UserOutlined />
            </template>
          </Statistic>
        </Card>
      </Col>
      <Col :xs="12" :sm="12" :md="6" :lg="4">
        <Card class="stat-card stat-present">
          <Statistic
            title="Hadir"
            :value="stats.present"
            :value-style="{ color: '#22c55e' }"
          >
            <template #prefix>
              <CheckCircleOutlined />
            </template>
          </Statistic>
        </Card>
      </Col>
      <Col :xs="12" :sm="12" :md="6" :lg="4">
        <Card class="stat-card stat-late">
          <Statistic
            title="Terlambat"
            :value="stats.late"
            :value-style="{ color: '#f97316' }"
          >
            <template #prefix>
              <ClockCircleOutlined />
            </template>
          </Statistic>
        </Card>
      </Col>
      <Col :xs="12" :sm="12" :md="6" :lg="4">
        <Card class="stat-card stat-very-late">
          <Statistic
            title="Sangat Terlambat"
            :value="stats.veryLate"
            :value-style="{ color: '#ef4444' }"
          >
            <template #prefix>
              <ExclamationCircleOutlined />
            </template>
          </Statistic>
        </Card>
      </Col>
      <Col :xs="12" :sm="12" :md="6" :lg="4">
        <Card class="stat-card stat-absent">
          <Statistic
            title="Tidak Hadir"
            :value="stats.absent"
            :value-style="{ color: '#6b7280' }"
          >
            <template #prefix>
              <CloseCircleOutlined />
            </template>
          </Statistic>
        </Card>
      </Col>
      <!-- Requirements: 4.10 - Show percentage of attendance completion -->
      <Col :xs="12" :sm="12" :md="6" :lg="4">
        <Card class="stat-card stat-percentage">
          <div class="percentage-stat">
            <Text type="secondary" class="stat-title">Persentase Kehadiran</Text>
            <Progress
              type="circle"
              :percent="Math.round(stats.percentage)"
              :size="80"
              :stroke-color="stats.percentage >= 90 ? '#22c55e' : stats.percentage >= 75 ? '#f97316' : '#ef4444'"
            />
          </div>
        </Card>
      </Col>
    </Row>

    <!-- Live Feed Section -->
    <Card style="margin-top: 16px">
      <template #title>
        <div class="feed-header">
          <div class="feed-title">
            <SyncOutlined v-if="connectionStatus.connected" spin style="color: #22c55e; margin-right: 8px" />
            <span>Live Feed</span>
            <Tag v-if="connectionStatus.connected" color="success" style="margin-left: 8px">LIVE</Tag>
          </div>
          <!-- Requirements: 4.5 - Allow filtering by class -->
          <div class="feed-actions">
            <Select
              v-model:value="filterClassId"
              placeholder="Semua Kelas"
              allow-clear
              style="width: 150px; margin-right: 8px"
              :loading="loadingClasses"
            >
              <SelectOption v-for="cls in classes" :key="cls.id" :value="cls.id">
                {{ cls.name }}
              </SelectOption>
            </Select>
            <Button @click="handleRefresh" :loading="loading">
              <template #icon><ReloadOutlined /></template>
            </Button>
          </div>
        </div>
      </template>

      <!-- Requirements: 4.3 - Show the 20 most recent attendance records -->
      <div ref="feedContainer" class="feed-container">
        <Spin v-if="loading" tip="Memuat data..." />
        <Empty v-else-if="liveFeed.length === 0" description="Belum ada data absensi hari ini" />
        <List v-else item-layout="horizontal" :data-source="liveFeed">
          <template #renderItem="{ item }">
            <!-- Requirements: 4.4 - Add new attendance with visual highlight -->
            <ListItem 
              :class="['feed-item', { highlighted: highlightedId === item.id }]"
            >
              <ListItemMeta>
                <template #avatar>
                  <Avatar 
                    :style="{ backgroundColor: getStatusColor(item.status) }"
                    size="large"
                  >
                    {{ item.studentName.charAt(0).toUpperCase() }}
                  </Avatar>
                </template>
                <template #title>
                  <div class="feed-item-title">
                    <Text strong>{{ item.studentName }}</Text>
                    <Tag :color="getStatusTagColor(item.status)">
                      {{ getStatusLabel(item.status) }}
                    </Tag>
                  </div>
                </template>
                <template #description>
                  <div class="feed-item-desc">
                    <Tag color="blue">{{ item.className }}</Tag>
                    <Text type="secondary">
                      {{ item.type === 'check_in' ? 'Masuk' : 'Keluar' }} • 
                      {{ formatTime(item.time) }}
                    </Text>
                    <Text type="secondary" class="relative-time">
                      {{ formatRelativeTime(item.time) }}
                    </Text>
                  </div>
                </template>
              </ListItemMeta>
            </ListItem>
          </template>
        </List>
      </div>
    </Card>
  </div>
</template>

<style scoped>
.live-attendance {
  padding: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.header-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stats-row {
  margin-bottom: 0;
}

.stat-card {
  height: 100%;
}

.stat-card :deep(.ant-statistic-title) {
  font-size: 13px;
  color: #8c8c8c;
}

.stat-card :deep(.ant-statistic-content-prefix) {
  margin-right: 8px;
}

.percentage-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.percentage-stat .stat-title {
  font-size: 13px;
}

.feed-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.feed-title {
  display: flex;
  align-items: center;
  font-weight: 600;
}

.feed-actions {
  display: flex;
  align-items: center;
}

.feed-container {
  max-height: 500px;
  overflow-y: auto;
}

.feed-item {
  transition: background-color 0.3s ease;
  padding: 12px 0 !important;
  border-bottom: 1px solid #f0f0f0;
}

.feed-item:last-child {
  border-bottom: none;
}

/* Requirements: 4.4 - Visual highlight for new entries */
.feed-item.highlighted {
  background-color: #fff7e6;
  animation: highlight-fade 3s ease-out;
}

@keyframes highlight-fade {
  0% {
    background-color: #ffe7ba;
  }
  100% {
    background-color: transparent;
  }
}

.feed-item-title {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.feed-item-desc {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  margin-top: 4px;
}

.relative-time {
  margin-left: auto;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .header-right {
    width: 100%;
    justify-content: flex-start;
  }

  .feed-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .feed-actions {
    width: 100%;
    flex-wrap: wrap;
  }

  .relative-time {
    margin-left: 0;
    width: 100%;
  }
}
</style>
