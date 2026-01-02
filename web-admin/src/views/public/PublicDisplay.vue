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
  Tag,
  Badge,
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

// Get rank badge color
const getRankColor = (rank: number): string => {
  switch (rank) {
    case 1:
      return '#ffd700' // Gold
    case 2:
      return '#c0c0c0' // Silver
    case 3:
      return '#cd7f32' // Bronze
    default:
      return '#3b82f6'
  }
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
      <!-- Requirements: 5.7, 5.8 - Show school name, current date and time -->
      <header class="display-header">
        <div class="header-left">
          <Title :level="1" class="school-name">{{ displayData.schoolName }}</Title>
          <Text class="date-text">{{ formattedDate }}</Text>
        </div>
        <div class="header-right">
          <div class="time-display">{{ formattedTime }}</div>
          <Badge 
            :status="connectionStatus.connected ? 'success' : 'error'" 
            class="connection-badge"
          >
            <template #text>
              <component 
                :is="connectionStatus.connected ? WifiOutlined : DisconnectOutlined" 
                class="connection-icon"
              />
            </template>
          </Badge>
        </div>
      </header>

      <!-- Stats Section -->
      <!-- Requirements: 5.5 - Show real-time statistics -->
      <section class="stats-section">
        <Row :gutter="[24, 24]">
          <Col :xs="12" :sm="8" :md="4">
            <Card class="stat-card stat-total">
              <div class="stat-content">
                <UserOutlined class="stat-icon" />
                <div class="stat-value">{{ displayData.stats.totalStudents }}</div>
                <div class="stat-label">Total Siswa</div>
              </div>
            </Card>
          </Col>
          <Col :xs="12" :sm="8" :md="4">
            <Card class="stat-card stat-present">
              <div class="stat-content">
                <CheckCircleOutlined class="stat-icon" />
                <div class="stat-value">{{ displayData.stats.present }}</div>
                <div class="stat-label">Hadir</div>
              </div>
            </Card>
          </Col>
          <Col :xs="12" :sm="8" :md="4">
            <Card class="stat-card stat-late">
              <div class="stat-content">
                <ClockCircleOutlined class="stat-icon" />
                <div class="stat-value">{{ displayData.stats.late }}</div>
                <div class="stat-label">Terlambat</div>
              </div>
            </Card>
          </Col>
          <Col :xs="12" :sm="8" :md="4">
            <Card class="stat-card stat-very-late">
              <div class="stat-content">
                <ExclamationCircleOutlined class="stat-icon" />
                <div class="stat-value">{{ displayData.stats.veryLate }}</div>
                <div class="stat-label">Sangat Terlambat</div>
              </div>
            </Card>
          </Col>
          <Col :xs="12" :sm="8" :md="4">
            <Card class="stat-card stat-absent">
              <div class="stat-content">
                <CloseCircleOutlined class="stat-icon" />
                <div class="stat-value">{{ displayData.stats.absent }}</div>
                <div class="stat-label">Tidak Hadir</div>
              </div>
            </Card>
          </Col>
          <Col :xs="12" :sm="8" :md="4">
            <Card class="stat-card stat-percentage">
              <div class="stat-content percentage-content">
                <Progress
                  type="circle"
                  :percent="Math.round(displayData.stats.percentage)"
                  :size="80"
                  :stroke-color="displayData.stats.percentage >= 90 ? '#22c55e' : displayData.stats.percentage >= 75 ? '#f97316' : '#ef4444'"
                  :stroke-width="8"
                />
                <div class="stat-label">Kehadiran</div>
              </div>
            </Card>
          </Col>
        </Row>
      </section>

      <!-- Main Content: Live Feed and Leaderboard -->
      <section class="main-content">
        <Row :gutter="[24, 24]">
          <!-- Live Feed Section -->
          <!-- Requirements: 5.4 - Show live feed of recent attendance (last 10 records) -->
          <Col :xs="24" :lg="14">
            <Card class="content-card live-feed-card">
              <template #title>
                <div class="card-title">
                  <SyncOutlined v-if="connectionStatus.connected" spin class="live-icon" />
                  <span>Absensi Terbaru</span>
                  <Tag v-if="connectionStatus.connected" color="success" class="live-tag">LIVE</Tag>
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
                          size="large"
                          class="feed-avatar"
                        >
                          {{ item.studentName.charAt(0).toUpperCase() }}
                        </Avatar>
                        <div class="feed-info">
                          <div class="feed-name">{{ item.studentName }}</div>
                          <div class="feed-meta">
                            <Tag color="blue">{{ item.className }}</Tag>
                            <span class="feed-time">{{ formatTime(item.time) }}</span>
                          </div>
                        </div>
                        <Tag :color="getStatusTagColor(item.status)" class="feed-status">
                          {{ getStatusLabel(item.status) }}
                        </Tag>
                      </div>
                    </ListItem>
                  </template>
                </List>
              </div>
            </Card>
          </Col>

          <!-- Leaderboard Section -->
          <!-- Requirements: 5.6 - Show leaderboard of top 10 earliest arrivals -->
          <Col :xs="24" :lg="10">
            <Card class="content-card leaderboard-card">
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
                          :style="{ backgroundColor: getRankColor(item.rank) }"
                        >
                          {{ item.rank }}
                        </div>
                        <div class="leaderboard-info">
                          <div class="leaderboard-name">{{ item.studentName }}</div>
                          <div class="leaderboard-meta">
                            <Tag color="blue">{{ item.className }}</Tag>
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
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  color: #ffffff;
  padding: 24px;
  overflow: hidden;
}

/* Loading State */
.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
}

.loading-container :deep(.ant-spin-text) {
  color: #ffffff;
  font-size: 18px;
}

/* Error State */
/* Requirements: 5.13 - Show error message if token is invalid */
.error-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
}

.error-content {
  text-align: center;
  max-width: 500px;
  padding: 48px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  backdrop-filter: blur(10px);
}

.error-icon {
  font-size: 80px;
  color: #ef4444;
  margin-bottom: 24px;
}

.error-title {
  color: #ffffff !important;
  margin-bottom: 16px !important;
}

.error-code {
  display: block;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6) !important;
  margin-bottom: 16px;
}

.error-hint {
  display: block;
  font-size: 16px;
  color: rgba(255, 255, 255, 0.8) !important;
  margin-bottom: 24px;
}

.refresh-button {
  margin-top: 16px;
}

/* Display Content */
.display-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
  height: calc(100vh - 48px);
}

/* Header Section */
/* Requirements: 5.7, 5.8 - Show school name, date, time */
.display-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  backdrop-filter: blur(10px);
}

.header-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.school-name {
  color: #ffffff !important;
  margin: 0 !important;
  font-size: 32px !important;
  font-weight: 700 !important;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
}

.date-text {
  font-size: 20px;
  color: rgba(255, 255, 255, 0.8);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 24px;
}

.time-display {
  font-size: 48px;
  font-weight: 700;
  font-family: 'Courier New', monospace;
  color: #22c55e;
  text-shadow: 0 0 20px rgba(34, 197, 94, 0.5);
}

.connection-badge {
  font-size: 24px;
}

.connection-icon {
  font-size: 24px;
  color: rgba(255, 255, 255, 0.8);
}

/* Stats Section */
/* Requirements: 5.5 - Show real-time statistics */
.stats-section {
  flex-shrink: 0;
}

.stat-card {
  background: rgba(255, 255, 255, 0.1) !important;
  border: none !important;
  border-radius: 16px !important;
  backdrop-filter: blur(10px);
  height: 100%;
}

.stat-card :deep(.ant-card-body) {
  padding: 20px !important;
}

.stat-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.stat-icon {
  font-size: 32px;
  color: rgba(255, 255, 255, 0.8);
}

.stat-value {
  font-size: 36px;
  font-weight: 700;
  color: #ffffff;
}

.stat-label {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.7);
  text-transform: uppercase;
  letter-spacing: 1px;
}

.stat-total .stat-icon { color: #3b82f6; }
.stat-total .stat-value { color: #3b82f6; }

.stat-present .stat-icon { color: #22c55e; }
.stat-present .stat-value { color: #22c55e; }

.stat-late .stat-icon { color: #f97316; }
.stat-late .stat-value { color: #f97316; }

.stat-very-late .stat-icon { color: #ef4444; }
.stat-very-late .stat-value { color: #ef4444; }

.stat-absent .stat-icon { color: #6b7280; }
.stat-absent .stat-value { color: #6b7280; }

.percentage-content {
  padding: 8px 0;
}

/* Main Content */
.main-content {
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

.content-card {
  background: rgba(255, 255, 255, 0.1) !important;
  border: none !important;
  border-radius: 16px !important;
  backdrop-filter: blur(10px);
  height: 100%;
}

.content-card :deep(.ant-card-head) {
  background: transparent !important;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1) !important;
  padding: 16px 24px !important;
}

.content-card :deep(.ant-card-head-title) {
  color: #ffffff !important;
  font-size: 20px !important;
  font-weight: 600 !important;
}

.content-card :deep(.ant-card-body) {
  padding: 16px 24px !important;
  height: calc(100% - 65px);
  overflow-y: auto;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 12px;
  color: #ffffff;
}

.live-icon {
  color: #22c55e;
  font-size: 20px;
}

.live-tag {
  font-size: 12px;
}

.trophy-icon {
  color: #ffd700;
  font-size: 24px;
}

/* Live Feed */
/* Requirements: 5.4 - Show live feed of recent attendance */
.feed-container {
  height: 100%;
  overflow-y: auto;
}

.feed-item {
  padding: 12px 0 !important;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1) !important;
}

.feed-item:last-child {
  border-bottom: none !important;
}

.feed-item-content {
  display: flex;
  align-items: center;
  gap: 16px;
  width: 100%;
}

.feed-avatar {
  flex-shrink: 0;
  font-size: 20px !important;
  font-weight: 600;
}

.feed-info {
  flex: 1;
  min-width: 0;
}

.feed-name {
  font-size: 18px;
  font-weight: 600;
  color: #ffffff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.feed-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 4px;
}

.feed-time {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.7);
  font-family: 'Courier New', monospace;
}

.feed-status {
  flex-shrink: 0;
  font-size: 14px;
}

/* Leaderboard */
/* Requirements: 5.6 - Show leaderboard of top 10 earliest arrivals */
.leaderboard-container {
  height: 100%;
  overflow-y: auto;
}

.leaderboard-item {
  padding: 12px 0 !important;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1) !important;
}

.leaderboard-item:last-child {
  border-bottom: none !important;
}

.leaderboard-item-content {
  display: flex;
  align-items: center;
  gap: 16px;
  width: 100%;
}

.rank-badge {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 700;
  color: #1a1a2e;
  flex-shrink: 0;
}

.leaderboard-info {
  flex: 1;
  min-width: 0;
}

.leaderboard-name {
  font-size: 18px;
  font-weight: 600;
  color: #ffffff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.leaderboard-meta {
  margin-top: 4px;
}

.leaderboard-time {
  font-size: 18px;
  font-weight: 600;
  color: #22c55e;
  font-family: 'Courier New', monospace;
  flex-shrink: 0;
}

/* Empty State */
.content-card :deep(.ant-empty-description) {
  color: rgba(255, 255, 255, 0.6);
}

/* Scrollbar Styling */
.feed-container::-webkit-scrollbar,
.leaderboard-container::-webkit-scrollbar,
.content-card :deep(.ant-card-body)::-webkit-scrollbar {
  width: 6px;
}

.feed-container::-webkit-scrollbar-track,
.leaderboard-container::-webkit-scrollbar-track,
.content-card :deep(.ant-card-body)::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 3px;
}

.feed-container::-webkit-scrollbar-thumb,
.leaderboard-container::-webkit-scrollbar-thumb,
.content-card :deep(.ant-card-body)::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.3);
  border-radius: 3px;
}

.feed-container::-webkit-scrollbar-thumb:hover,
.leaderboard-container::-webkit-scrollbar-thumb:hover,
.content-card :deep(.ant-card-body)::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.5);
}

/* Responsive Design */
@media (max-width: 1200px) {
  .school-name {
    font-size: 28px !important;
  }

  .time-display {
    font-size: 40px;
  }

  .stat-value {
    font-size: 28px;
  }
}

@media (max-width: 992px) {
  .public-display {
    padding: 16px;
  }

  .display-header {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }

  .header-right {
    flex-direction: column;
    gap: 12px;
  }

  .school-name {
    font-size: 24px !important;
  }

  .time-display {
    font-size: 36px;
  }

  .stat-value {
    font-size: 24px;
  }

  .stat-icon {
    font-size: 24px;
  }
}

@media (max-width: 576px) {
  .school-name {
    font-size: 20px !important;
  }

  .date-text {
    font-size: 16px;
  }

  .time-display {
    font-size: 28px;
  }

  .stat-value {
    font-size: 20px;
  }

  .feed-name,
  .leaderboard-name {
    font-size: 16px;
  }

  .feed-time,
  .leaderboard-time {
    font-size: 14px;
  }
}
</style>
