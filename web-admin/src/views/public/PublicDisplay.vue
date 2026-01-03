<script setup lang="ts">
// Public Display Page for LCD/Monitor Kiosk
// Requirements: 5.4 - Show live feed of recent attendance (last 10 records)
// Requirements: 5.5 - Show real-time statistics (present, late, absent count)
// Requirements: 5.6 - Show leaderboard of top 10 earliest arrivals for the day
// Requirements: 5.7 - Show current date and time
// Requirements: 5.8 - Show school name
// Requirements: 5.9 - Update public display within 3 seconds when new attendance recorded
// Requirements: 5.12 - Use full-screen optimized layout with large fonts
// Requirements: 5.13 - Show error message if token is invalid or revoked

import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import {
  Row,
  Col,
  Card,
  Typography,
  Spin,
  Progress,
  List,
  ListItem,
  Avatar,
  Empty,
  Button,
} from 'ant-design-vue'
import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  ClockCircleOutlined,
  ExclamationCircleOutlined,
  UserOutlined,
  TrophyOutlined,
  WifiOutlined,
  DisconnectOutlined,
  SyncOutlined,
  ReloadOutlined,
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import 'dayjs/locale/id'

dayjs.locale('id')

const { Title, Text } = Typography

// Types for public display data
interface PublicAttendanceStats {
  totalStudents: number
  present: number
  late: number
  veryLate: number
  absent: number
  percentage: number
}

interface PublicLiveFeedEntry {
  studentName: string
  className: string
  time: string
  status: string
  type: string
}

interface PublicLeaderboardEntry {
  rank: number
  studentName: string
  className: string
  arrivalTime: string
}

interface PublicDisplayData {
  schoolName: string
  currentTime: string
  date: string
  stats: PublicAttendanceStats
  liveFeed: PublicLiveFeedEntry[]
  leaderboard: PublicLeaderboardEntry[]
}

interface WSConnectionStatus {
  connected: boolean
  message: string
}

// Route params
const route = useRoute()
const token = computed(() => route.params.token as string)

// State
const loading = ref(true)
const error = ref<string | null>(null)
const errorCode = ref<string | null>(null)
const displayData = ref<PublicDisplayData | null>(null)
const currentTime = ref(dayjs())
const connectionStatus = ref<WSConnectionStatus>({ connected: false, message: 'Menghubungkan...' })

// WebSocket connection
let ws: WebSocket | null = null
let reconnectTimeout: ReturnType<typeof setTimeout> | null = null
let reconnectAttempts = 0
const MAX_RECONNECT_ATTEMPTS = 10
const RECONNECT_DELAY = 3000

// Time update interval
let timeInterval: ReturnType<typeof setInterval> | null = null

// Computed values
const formattedTime = computed(() => currentTime.value.format('HH:mm:ss'))
const formattedDate = computed(() => displayData.value?.date || currentTime.value.format('dddd, D MMMM YYYY'))

// Get API base URL
const getApiBaseUrl = (): string => {
  // @ts-expect-error - Vite env types
  const apiUrl = import.meta.env?.VITE_API_URL as string | undefined
  if (apiUrl) {
    return apiUrl.replace(/\/api\/v1\/?$/, '')
  }
  return `${window.location.protocol}//${window.location.host}`
}

// Get WebSocket URL
const getWsUrl = (): string => {
  const baseUrl = getApiBaseUrl()
  const protocol = baseUrl.startsWith('https') ? 'wss:' : 'ws:'
  const host = baseUrl.replace(/^https?:\/\//, '')
  return `${protocol}//${host}/api/v1/public/display/${token.value}/ws`
}

// Fetch initial data via REST API
const fetchDisplayData = async () => {
  loading.value = true
  error.value = null
  errorCode.value = null

  try {
    const baseUrl = getApiBaseUrl()
    const response = await fetch(`${baseUrl}/api/v1/public/display/${token.value}`)
    const data = await response.json()

    if (!response.ok) {
      errorCode.value = data.error || 'UNKNOWN_ERROR'
      error.value = data.message || 'Gagal memuat data display'
      return
    }

    if (data.success && data.data) {
      displayData.value = transformDisplayData(data.data)
    }
  } catch (err) {
    console.error('Failed to fetch display data:', err)
    errorCode.value = 'NETWORK_ERROR'
    error.value = 'Gagal terhubung ke server. Periksa koneksi internet.'
  } finally {
    loading.value = false
  }
}

// Transform API response to frontend format
const transformDisplayData = (data: any): PublicDisplayData => {
  return {
    schoolName: data.school_name,
    currentTime: data.current_time,
    date: data.date,
    stats: {
      totalStudents: data.stats.total_students,
      present: data.stats.present,
      late: data.stats.late,
      veryLate: data.stats.very_late,
      absent: data.stats.absent,
      percentage: data.stats.percentage,
    },
    liveFeed: (data.live_feed || []).map((entry: any) => ({
      studentName: entry.student_name,
      className: entry.class_name,
      time: entry.time,
      status: entry.status,
      type: entry.type,
    })),
    leaderboard: (data.leaderboard || []).map((entry: any) => ({
      rank: entry.rank,
      studentName: entry.student_name,
      className: entry.class_name,
      arrivalTime: entry.arrival_time,
    })),
  }
}

// Connect to WebSocket for real-time updates
// Requirements: 5.9 - Update public display within 3 seconds
const connectWebSocket = () => {
  if (ws?.readyState === WebSocket.OPEN) {
    return
  }

  connectionStatus.value = { connected: false, message: 'Menghubungkan...' }

  try {
    ws = new WebSocket(getWsUrl())

    ws.onopen = () => {
      console.log('Public display WebSocket connected')
      connectionStatus.value = { connected: true, message: 'Terhubung' }
      reconnectAttempts = 0
    }

    ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data)
        handleWebSocketMessage(message)
      } catch (err) {
        console.error('Failed to parse WebSocket message:', err)
      }
    }

    ws.onerror = (err) => {
      console.error('WebSocket error:', err)
      connectionStatus.value = { connected: false, message: 'Kesalahan koneksi' }
    }

    ws.onclose = (event) => {
      console.log('WebSocket closed:', event.code, event.reason)
      connectionStatus.value = { connected: false, message: 'Terputus' }

      if (event.code !== 1000) {
        handleReconnect()
      }
    }
  } catch (err) {
    console.error('Failed to create WebSocket:', err)
    handleReconnect()
  }
}

// Handle WebSocket messages
const handleWebSocketMessage = (message: any) => {
  switch (message.type) {
    case 'connected':
      connectionStatus.value = { connected: true, message: 'Terhubung' }
      break

    case 'error':
      if (message.payload) {
        errorCode.value = message.payload.error
        error.value = message.payload.message
      }
      break

    case 'new_attendance':
    case 'stats_update':
    case 'leaderboard_update':
      handleAttendanceUpdate(message)
      break

    case 'refresh_data':
      if (message.payload) {
        displayData.value = transformDisplayData(message.payload)
      }
      break

    case 'pong':
      // Ping response received
      break
  }
}

// Handle attendance updates from WebSocket
const handleAttendanceUpdate = (message: any) => {
  if (!displayData.value) return

  // Update attendance entry
  if (message.attendance) {
    const newEntry: PublicLiveFeedEntry = {
      studentName: message.attendance.student_name || message.attendance.studentName,
      className: message.attendance.class_name || message.attendance.className,
      time: message.attendance.time,
      status: message.attendance.status,
      type: message.attendance.type,
    }
    // Add to top and keep only last 10
    displayData.value.liveFeed = [newEntry, ...displayData.value.liveFeed.slice(0, 9)]
  }

  // Update stats
  if (message.stats) {
    displayData.value.stats = {
      totalStudents: message.stats.total_students || message.stats.totalStudents,
      present: message.stats.present,
      late: message.stats.late,
      veryLate: message.stats.very_late || message.stats.veryLate,
      absent: message.stats.absent,
      percentage: message.stats.percentage,
    }
  }

  // Update leaderboard
  if (message.leaderboard) {
    displayData.value.leaderboard = message.leaderboard.map((entry: any) => ({
      rank: entry.rank,
      studentName: entry.student_name || entry.studentName,
      className: entry.class_name || entry.className,
      arrivalTime: entry.arrival_time || entry.arrivalTime,
    }))
  }
}

// Handle WebSocket reconnection
// Validate token before reconnecting to avoid infinite loop
const handleReconnect = async () => {
  if (reconnectAttempts >= MAX_RECONNECT_ATTEMPTS) {
    connectionStatus.value = { connected: false, message: 'Gagal terhubung. Refresh halaman.' }
    return
  }

  // Validate token before reconnecting
  if (reconnectAttempts > 0 && reconnectAttempts % 3 === 0) {
    // Every 3 attempts, validate token via REST API
    try {
      const baseUrl = getApiBaseUrl()
      const response = await fetch(`${baseUrl}/api/v1/public/display/${token.value}`)
      const data = await response.json()
      
      if (!response.ok) {
        // Token is invalid/revoked/expired - stop reconnecting
        errorCode.value = data.error?.code || data.error || 'TOKEN_INVALID'
        error.value = getErrorMessage(errorCode.value)
        connectionStatus.value = { connected: false, message: 'Token tidak valid' }
        return
      }
    } catch {
      // Network error - continue trying to reconnect
    }
  }

  reconnectAttempts++
  const delay = RECONNECT_DELAY * Math.min(reconnectAttempts, 5)
  connectionStatus.value = { 
    connected: false, 
    message: `Menghubungkan ulang (${reconnectAttempts}/${MAX_RECONNECT_ATTEMPTS})...` 
  }

  reconnectTimeout = setTimeout(() => {
    connectWebSocket()
  }, delay)
}

// Get user-friendly error message
const getErrorMessage = (code: string | null): string => {
  switch (code) {
    case 'TOKEN_INVALID':
    case 'TOKEN_NOT_FOUND':
      return 'Token display tidak valid atau tidak ditemukan.'
    case 'TOKEN_REVOKED':
      return 'Token display telah dicabut oleh administrator.'
    case 'TOKEN_EXPIRED':
      return 'Token display telah kedaluwarsa.'
    case 'NETWORK_ERROR':
      return 'Gagal terhubung ke server. Periksa koneksi internet.'
    default:
      return 'Terjadi kesalahan. Silakan hubungi administrator.'
  }
}

// Manual refresh function
const handleManualRefresh = async () => {
  loading.value = true
  error.value = null
  errorCode.value = null
  reconnectAttempts = 0
  
  disconnectWebSocket()
  await fetchDisplayData()
  
  if (!error.value) {
    connectWebSocket()
  }
}

// Disconnect WebSocket
const disconnectWebSocket = () => {
  if (reconnectTimeout) {
    clearTimeout(reconnectTimeout)
    reconnectTimeout = null
  }

  if (ws) {
    ws.close(1000, 'Client disconnect')
    ws = null
  }
}

// Get status color for attendance
const getStatusColor = (status: string): string => {
  switch (status) {
    case 'present':
    case 'on_time':
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
    case 'on_time':
      return 'Tepat Waktu'
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
    case 'on_time':
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



// Start time update interval
const startTimeUpdate = () => {
  timeInterval = setInterval(() => {
    currentTime.value = dayjs()
  }, 1000)
}

// Stop time update interval
const stopTimeUpdate = () => {
  if (timeInterval) {
    clearInterval(timeInterval)
    timeInterval = null
  }
}

// Lifecycle hooks
onMounted(async () => {
  startTimeUpdate()
  await fetchDisplayData()
  
  if (!error.value) {
    connectWebSocket()
  }
})

onUnmounted(() => {
  stopTimeUpdate()
  disconnectWebSocket()
})
</script>

<template>
  <div class="public-display">
    <!-- Loading State -->
    <div v-if="loading" class="loading-container">
      <Spin size="large" tip="Memuat data..." />
    </div>

    <!-- Error State -->
    <!-- Requirements: 5.13 - Show error message if token is invalid or revoked -->
    <div v-else-if="error" class="error-container">
      <div class="error-content">
        <ExclamationCircleOutlined class="error-icon" />
        <Title :level="2" class="error-title">{{ getErrorMessage(errorCode) }}</Title>
        <Text type="secondary" class="error-code" v-if="errorCode">Kode Error: {{ errorCode }}</Text>
        <Text type="secondary" class="error-hint">
          Silakan hubungi administrator sekolah untuk mendapatkan token display yang valid.
        </Text>
        <Button type="primary" size="large" class="refresh-button" @click="handleManualRefresh">
          <template #icon><ReloadOutlined /></template>
          Coba Lagi
        </Button>
      </div>
    </div>

    <!-- Main Display Content -->
    <div v-else-if="displayData" class="display-content">
      <!-- Header Section -->
      <header class="display-header">
        <div class="header-left">
          <Title :level="1" class="school-name">{{ displayData.schoolName }}</Title>
          <Text class="date-text">{{ formattedDate }}</Text>
        </div>
        <div class="header-right">
          <div class="time-display">{{ formattedTime }}</div>
          <div 
            class="connection-badge"
            :title="connectionStatus.connected ? 'Terhubung' : 'Terputus'"
          >
            <component 
              :is="connectionStatus.connected ? WifiOutlined : DisconnectOutlined" 
              class="connection-icon"
              :style="{ color: connectionStatus.connected ? '#22c55e' : '#ef4444' }"
            />
          </div>
        </div>
      </header>

      <!-- Stats Section -->
      <section class="stats-section">
        <Row :gutter="[24, 24]">
          <Col :xs="12" :sm="8" :md="4">
            <Card class="stat-card stat-total" :bordered="false">
              <div class="stat-content">
                <UserOutlined class="stat-icon" />
                <div class="stat-value">{{ displayData.stats.totalStudents }}</div>
                <div class="stat-label">Total Siswa</div>
              </div>
            </Card>
          </Col>
          <Col :xs="12" :sm="8" :md="4">
            <Card class="stat-card stat-present" :bordered="false">
              <div class="stat-content">
                <CheckCircleOutlined class="stat-icon" />
                <div class="stat-value">{{ displayData.stats.present }}</div>
                <div class="stat-label">Hadir</div>
              </div>
            </Card>
          </Col>
          <Col :xs="12" :sm="8" :md="4">
            <Card class="stat-card stat-late" :bordered="false">
              <div class="stat-content">
                <ClockCircleOutlined class="stat-icon" />
                <div class="stat-value">{{ displayData.stats.late }}</div>
                <div class="stat-label">Terlambat</div>
              </div>
            </Card>
          </Col>
          <Col :xs="12" :sm="8" :md="4">
            <Card class="stat-card stat-very-late" :bordered="false">
              <div class="stat-content">
                <ExclamationCircleOutlined class="stat-icon" />
                <div class="stat-value">{{ displayData.stats.veryLate }}</div>
                <div class="stat-label">Sangat Terlambat</div>
              </div>
            </Card>
          </Col>
          <Col :xs="12" :sm="8" :md="4">
            <Card class="stat-card stat-absent" :bordered="false">
              <div class="stat-content">
                <CloseCircleOutlined class="stat-icon" />
                <div class="stat-value">{{ displayData.stats.absent }}</div>
                <div class="stat-label">Tidak Hadir</div>
              </div>
            </Card>
          </Col>
          <Col :xs="12" :sm="8" :md="4">
            <Card class="stat-card stat-percentage" :bordered="false">
              <div class="stat-content percentage-content">
                <Progress
                  type="circle"
                  :percent="Math.round(displayData.stats.percentage)"
                  :size="50"
                  :stroke-color="displayData.stats.percentage >= 90 ? '#22c55e' : displayData.stats.percentage >= 75 ? '#f97316' : '#ef4444'"
                  :stroke-width="8"
                  :trail-color="'#f1f5f9'"
                />
                <div class="stat-label">Kehadiran</div>
              </div>
            </Card>
          </Col>
        </Row>
      </section>

      <!-- Main Content: Live Feed and Leaderboard -->
      <section class="main-content">
        <Row :gutter="[24, 24]" style="height: 100%">
          <!-- Live Feed Section -->
          <Col :xs="24" :lg="14" style="height: 100%">
            <Card class="content-card live-feed-card" :bordered="false">
              <template #title>
                <div class="card-title">
                  <SyncOutlined v-if="connectionStatus.connected" spin class="live-icon" />
                  <span>Absensi Terbaru</span>
                  <span v-if="connectionStatus.connected" class="status-badge success" style="font-size: 11px; padding: 2px 8px;">LIVE</span>
                </div>
              </template>
              
              <div class="feed-container">
                <Empty v-if="displayData.liveFeed.length === 0" description="Belum ada data absensi hari ini" />
                <List v-else :data-source="displayData.liveFeed" :split="false">
                  <template #renderItem="{ item }">
                    <ListItem class="feed-item">
                      <div class="feed-item-content">
                        <Avatar 
                          :style="{ backgroundColor: getStatusColor(item.status) }"
                          class="feed-avatar"
                        >
                          {{ item.studentName.charAt(0).toUpperCase() }}
                        </Avatar>
                        <div class="feed-info">
                          <div class="feed-name">{{ item.studentName }}</div>
                          <div class="feed-meta">
                            <span class="class-badge">{{ item.className }}</span>
                            <span class="feed-time">{{ formatTime(item.time) }}</span>
                          </div>
                        </div>
                        <span :class="['status-badge', getStatusTagColor(item.status)]" class="feed-status">
                          {{ getStatusLabel(item.status) }}
                        </span>
                      </div>
                    </ListItem>
                  </template>
                </List>
              </div>
            </Card>
          </Col>

          <!-- Leaderboard Section -->
          <Col :xs="24" :lg="10" style="height: 100%">
            <Card class="content-card leaderboard-card" :bordered="false">
              <template #title>
                <div class="card-title">
                  <TrophyOutlined class="trophy-icon" />
                  <span>Siswa Terawal Hari Ini</span>
                </div>
              </template>
              
              <div class="leaderboard-container">
                <Empty v-if="displayData.leaderboard.length === 0" description="Belum ada data kehadiran" />
                <List v-else :data-source="displayData.leaderboard" :split="false">
                  <template #renderItem="{ item }">
                    <ListItem class="leaderboard-item">
                      <div class="leaderboard-item-content">
                        <div 
                          class="rank-badge"
                          :class="item.rank <= 3 ? `rank-badge-${item.rank}` : 'rank-badge-other'"
                        >
                          {{ item.rank }}
                        </div>
                        <div class="leaderboard-info">
                          <div class="leaderboard-name">{{ item.studentName }}</div>
                          <div class="leaderboard-meta">
                            <span class="class-badge">{{ item.className }}</span>
                          </div>
                        </div>
                        <div class="leaderboard-time">
                          {{ formatTime(item.arrivalTime) }}
                        </div>
                      </div>
                    </ListItem>
                  </template>
                </List>
              </div>
            </Card>
          </Col>
        </Row>
      </section>
    </div>
  </div>
</template>


<style scoped>
/* Requirements: 5.12 - Use full-screen optimized layout with large fonts */
.public-display {
  min-height: 100vh;
  background: #f8fafc;
  color: #1e293b;
  padding: 16px;
  overflow: hidden;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
}

/* Loading State */
.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
}

.loading-container :deep(.ant-spin-text) {
  color: #64748b;
  font-size: 16px;
  margin-top: 16px;
}

/* Error State */
.error-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
}

.error-content {
  text-align: center;
  max-width: 480px;
  padding: 32px;
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  border: 1px solid #e2e8f0;
}

.error-icon {
  font-size: 48px;
  color: #ef4444;
  margin-bottom: 16px;
}

.error-title {
  color: #1e293b !important;
  font-size: 24px !important;
  margin-bottom: 12px !important;
}

.error-code {
  display: block;
  font-size: 13px;
  color: #64748b !important;
  margin-bottom: 12px;
  font-family: monospace;
}

.error-hint {
  display: block;
  font-size: 14px;
  color: #475569 !important;
  margin-bottom: 20px;
}

.refresh-button {
  margin-top: 12px;
}

/* Display Content */
.display-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
  height: calc(100vh - 32px);
}

/* Header Section */
.display-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
  border: 1px solid #f1f5f9;
}

.header-left {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.school-name {
  color: #1e293b !important;
  margin: 0 !important;
  font-size: 24px !important;
  font-weight: 700 !important;
  letter-spacing: -0.5px;
}

.date-text {
  font-size: 14px;
  color: #64748b;
  font-weight: 500;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 24px;
}

.time-display {
  font-size: 32px;
  font-weight: 700;
  font-variant-numeric: tabular-nums;
  color: #f97316;
  letter-spacing: -0.5px;
  background: #fff7ed;
  padding: 4px 16px;
  border-radius: 8px;
  border: 1px solid #ffedd5;
}

.connection-badge {
  font-size: 20px;
  display: flex;
  align-items: center;
}

.connection-icon {
  font-size: 20px;
  color: #94a3b8;
}

/* Stats Section */
.stats-section {
  flex-shrink: 0;
}

.stat-card {
  background: #ffffff !important;
  border: 1px solid #f1f5f9 !important;
  border-radius: 12px !important;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  height: 100%;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.stat-card :deep(.ant-card-body) {
  padding: 16px !important;
}

.stat-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.stat-icon {
  font-size: 24px;
  padding: 8px;
  border-radius: 8px;
  margin-bottom: 2px;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1e293b;
  line-height: 1;
}

.stat-label {
  font-size: 12px;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  font-weight: 600;
}

/* Stat Variants */
.stat-total .stat-icon { color: #3b82f6; background: #eff6ff; }
.stat-total .stat-value { color: #1e293b; }

.stat-present .stat-icon { color: #22c55e; background: #f0fdf4; }
.stat-present .stat-value { color: #1e293b; }

.stat-late .stat-icon { color: #f97316; background: #fff7ed; }
.stat-late .stat-value { color: #1e293b; }

.stat-very-late .stat-icon { color: #ef4444; background: #fef2f2; }
.stat-very-late .stat-value { color: #1e293b; }

.stat-absent .stat-icon { color: #64748b; background: #f8fafc; }
.stat-absent .stat-value { color: #1e293b; }

/* Main Content */
.main-content {
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

.content-card {
  background: #ffffff !important;
  border: 1px solid #f1f5f9 !important;
  border-radius: 12px !important;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1) !important;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.content-card :deep(.ant-card-head) {
  background: #f8fafc !important;
  border-bottom: 1px solid #e2e8f0 !important;
  padding: 12px 16px !important;
  min-height: 56px;
  border-top-left-radius: 12px !important;
  border-top-right-radius: 12px !important;
}

.content-card :deep(.ant-card-head-title) {
  padding: 0;
}

.content-card :deep(.ant-card-body) {
  padding: 0 !important;
  flex: 1;
  overflow: hidden;
  background: #ffffff;
  border-bottom-left-radius: 12px;
  border-bottom-right-radius: 12px;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #1e293b;
  font-size: 16px;
  font-weight: 600;
}

.live-icon {
  color: #22c55e;
  font-size: 16px;
}

.trophy-icon {
  color: #eab308;
  font-size: 20px;
}

/* Feed & Leaderboard Containers */
.feed-container,
.leaderboard-container {
  height: 100%;
  overflow-y: auto;
  padding: 0 16px;
}

/* List Items */
.feed-item,
.leaderboard-item {
  padding: 10px 0 !important;
  border-bottom: 1px solid #f1f5f9 !important;
  transition: background-color 0.2s;
}

.feed-item:last-child,
.leaderboard-item:last-child {
  border-bottom: none !important;
}

.feed-item-content,
.leaderboard-item-content {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.feed-avatar {
  flex-shrink: 0;
  font-weight: 600;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  font-size: 14px !important;
  width: 32px !important;
  height: 32px !important;
  line-height: 32px !important;
}

.feed-info,
.leaderboard-info {
  flex: 1;
  min-width: 0;
}

.feed-name,
.leaderboard-name {
  font-size: 14px;
  font-weight: 600;
  color: #1e293b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 2px;
}

.feed-meta,
.leaderboard-meta {
  display: flex;
  align-items: center;
  gap: 8px;
}

.feed-time,
.leaderboard-time {
  font-size: 12px;
  color: #64748b;
  font-variant-numeric: tabular-nums;
  font-weight: 500;
}

.leaderboard-time {
  font-size: 14px;
  color: #15803d;
  font-weight: 700;
}

/* Badges */
.status-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  border-width: 1px;
  border-style: solid;
}

.status-badge.success { background: #f0fdf4; color: #166534; border-color: #dcfce7; }
.status-badge.warning { background: #fff7ed; color: #9a3412; border-color: #ffedd5; }
.status-badge.error { background: #fef2f2; color: #991b1b; border-color: #fee2e2; }
.status-badge.default { background: #f1f5f9; color: #475569; border-color: #e2e8f0; }

.class-badge {
  background: #f1f5f9;
  color: #475569;
  padding: 1px 6px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
  border: 1px solid #e2e8f0;
}

.rank-badge {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 800;
  color: #ffffff;
  flex-shrink: 0;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  text-shadow: 0 1px 2px rgba(0,0,0,0.1);
}

/* Rank Colors */
.rank-badge-1 { background: linear-gradient(135deg, #fbbf24 0%, #d97706 100%); }
.rank-badge-2 { background: linear-gradient(135deg, #94a3b8 0%, #475569 100%); }
.rank-badge-3 { background: linear-gradient(135deg, #d97706 0%, #b45309 100%); }
.rank-badge-other { background: #f1f5f9; color: #64748b; border: 1px solid #e2e8f0; }

/* Empty State */
.content-card :deep(.ant-empty-description) {
  color: #94a3b8;
}

/* Scrollbar Styling */
.feed-container::-webkit-scrollbar,
.leaderboard-container::-webkit-scrollbar {
  width: 4px;
}

.feed-container::-webkit-scrollbar-track,
.leaderboard-container::-webkit-scrollbar-track {
  background: transparent;
}

.feed-container::-webkit-scrollbar-thumb,
.leaderboard-container::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 2px;
}

.feed-container::-webkit-scrollbar-thumb:hover,
.leaderboard-container::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

/* Responsive Design */
@media (max-width: 1200px) {
  .school-name { font-size: 20px !important; }
  .time-display { font-size: 28px; }
  .stat-value { font-size: 24px; }
}

@media (max-width: 992px) {
  .public-display { padding: 12px; }
  .display-header { flex-direction: column; gap: 12px; text-align: center; }
  .header-right { flex-direction: column; gap: 8px; }
  .school-name { font-size: 20px !important; }
  .time-display { font-size: 28px; }
}

@media (max-width: 576px) {
  .school-name { font-size: 18px !important; }
  .date-text { font-size: 13px; }
  .time-display { font-size: 24px; }
  .stat-value { font-size: 20px; }
  .feed-name, .leaderboard-name { font-size: 13px; }
  .feed-time, .leaderboard-time { font-size: 11px; }
}
</style>
```
