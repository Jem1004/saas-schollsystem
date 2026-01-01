// Real-time WebSocket service for live attendance updates
// Requirements: 4.8 - Use WebSocket for real-time updates
// Requirements: 4.9 - Display connection status and attempt to reconnect automatically

import api from './api'
import type {
  AttendanceEvent,
  AttendanceStats,
  LiveFeedEntry,
  LeaderboardEntry,
  LiveFeedResponse,
  StatsResponse,
  LeaderboardResponse,
  WSConnectionStatus,
} from '@/types/realtime'

// WebSocket connection states
export type ConnectionState = 'connecting' | 'connected' | 'disconnected' | 'reconnecting'

// Event handler types
export type AttendanceEventHandler = (event: AttendanceEvent) => void
export type ConnectionStatusHandler = (status: WSConnectionStatus) => void

// WebSocket configuration
const WS_RECONNECT_DELAY = 3000 // 3 seconds
const WS_MAX_RECONNECT_ATTEMPTS = 10
const WS_PING_INTERVAL = 30000 // 30 seconds

class RealtimeService {
  private ws: WebSocket | null = null
  private reconnectAttempts = 0
  private reconnectTimeout: ReturnType<typeof setTimeout> | null = null
  private pingInterval: ReturnType<typeof setInterval> | null = null
  private eventHandlers: Set<AttendanceEventHandler> = new Set()
  private statusHandlers: Set<ConnectionStatusHandler> = new Set()
  private classIdFilter: number | undefined = undefined
  private connectionState: ConnectionState = 'disconnected'

  // Get WebSocket URL from environment or construct from API base URL
  private getWsUrl(): string {
    // @ts-expect-error - Vite env types
    const wsUrl = import.meta.env?.VITE_WS_URL as string | undefined
    if (wsUrl) {
      return wsUrl
    }

    // Construct WebSocket URL from current location
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    return `${protocol}//${host}/api/v1/ws/attendance`
  }

  // Connect to WebSocket server
  // Requirements: 4.8 - Use WebSocket for real-time updates
  connect(classId?: number): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      console.log('WebSocket already connected')
      return
    }

    this.classIdFilter = classId
    this.connectionState = 'connecting'
    this.notifyStatusChange({ connected: false, message: 'Menghubungkan...' })

    const token = localStorage.getItem('accessToken')
    if (!token) {
      console.error('No access token available for WebSocket connection')
      this.notifyStatusChange({ connected: false, message: 'Token tidak tersedia' })
      return
    }

    // Build WebSocket URL with token and optional class filter
    let wsUrl = this.getWsUrl()
    const params = new URLSearchParams()
    params.append('token', token)
    if (classId) {
      params.append('class_id', classId.toString())
    }
    wsUrl += `?${params.toString()}`

    try {
      this.ws = new WebSocket(wsUrl)
      this.setupEventListeners()
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error)
      this.handleReconnect()
    }
  }

  // Setup WebSocket event listeners
  private setupEventListeners(): void {
    if (!this.ws) return

    this.ws.onopen = () => {
      console.log('WebSocket connected')
      this.connectionState = 'connected'
      this.reconnectAttempts = 0
      this.notifyStatusChange({ connected: true, message: 'Terhubung' })
      this.startPingInterval()
    }

    this.ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data) as AttendanceEvent
        this.notifyEventHandlers(data)
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error)
      }
    }

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error)
      this.notifyStatusChange({ connected: false, message: 'Kesalahan koneksi' })
    }

    // Requirements: 4.9 - Attempt to reconnect automatically
    this.ws.onclose = (event) => {
      console.log('WebSocket closed:', event.code, event.reason)
      this.connectionState = 'disconnected'
      this.stopPingInterval()
      
      if (event.code !== 1000) {
        // Abnormal closure, attempt reconnect
        this.handleReconnect()
      } else {
        this.notifyStatusChange({ connected: false, message: 'Terputus' })
      }
    }
  }

  // Handle reconnection logic
  // Requirements: 4.9 - Attempt to reconnect automatically
  private handleReconnect(): void {
    if (this.reconnectAttempts >= WS_MAX_RECONNECT_ATTEMPTS) {
      console.error('Max reconnection attempts reached')
      this.notifyStatusChange({ 
        connected: false, 
        message: 'Gagal terhubung. Silakan refresh halaman.' 
      })
      return
    }

    this.connectionState = 'reconnecting'
    this.reconnectAttempts++
    
    const delay = WS_RECONNECT_DELAY * Math.min(this.reconnectAttempts, 5)
    this.notifyStatusChange({ 
      connected: false, 
      message: `Menghubungkan ulang (${this.reconnectAttempts}/${WS_MAX_RECONNECT_ATTEMPTS})...` 
    })

    this.reconnectTimeout = setTimeout(() => {
      console.log(`Reconnection attempt ${this.reconnectAttempts}`)
      this.connect(this.classIdFilter)
    }, delay)
  }

  // Start ping interval to keep connection alive
  private startPingInterval(): void {
    this.stopPingInterval()
    this.pingInterval = setInterval(() => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        this.ws.send(JSON.stringify({ type: 'ping' }))
      }
    }, WS_PING_INTERVAL)
  }

  // Stop ping interval
  private stopPingInterval(): void {
    if (this.pingInterval) {
      clearInterval(this.pingInterval)
      this.pingInterval = null
    }
  }

  // Disconnect from WebSocket server
  disconnect(): void {
    this.stopPingInterval()
    
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout)
      this.reconnectTimeout = null
    }

    if (this.ws) {
      this.ws.close(1000, 'Client disconnect')
      this.ws = null
    }

    this.connectionState = 'disconnected'
    this.reconnectAttempts = 0
    this.notifyStatusChange({ connected: false, message: 'Terputus' })
  }

  // Subscribe to attendance events
  onAttendanceEvent(handler: AttendanceEventHandler): () => void {
    this.eventHandlers.add(handler)
    return () => {
      this.eventHandlers.delete(handler)
    }
  }

  // Subscribe to connection status changes
  onConnectionStatus(handler: ConnectionStatusHandler): () => void {
    this.statusHandlers.add(handler)
    return () => {
      this.statusHandlers.delete(handler)
    }
  }

  // Notify all event handlers
  private notifyEventHandlers(event: AttendanceEvent): void {
    this.eventHandlers.forEach(handler => {
      try {
        handler(event)
      } catch (error) {
        console.error('Error in event handler:', error)
      }
    })
  }

  // Notify all status handlers
  private notifyStatusChange(status: WSConnectionStatus): void {
    this.statusHandlers.forEach(handler => {
      try {
        handler(status)
      } catch (error) {
        console.error('Error in status handler:', error)
      }
    })
  }

  // Get current connection state
  getConnectionState(): ConnectionState {
    return this.connectionState
  }

  // Check if connected
  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN
  }

  // Update class filter (reconnects if needed)
  setClassFilter(classId?: number): void {
    if (this.classIdFilter !== classId) {
      this.classIdFilter = classId
      if (this.isConnected()) {
        this.disconnect()
        this.connect(classId)
      }
    }
  }

  // REST API fallback methods for initial data load

  // Get live feed via REST API
  async getLiveFeed(classId?: number): Promise<LiveFeedEntry[]> {
    const params: Record<string, string> = {}
    if (classId) {
      params.class_id = classId.toString()
    }
    const response = await api.get<LiveFeedResponse>('/realtime/feed', { params })
    return response.data.feed || []
  }

  // Get attendance stats via REST API
  async getStats(classId?: number): Promise<AttendanceStats> {
    const params: Record<string, string> = {}
    if (classId) {
      params.class_id = classId.toString()
    }
    const response = await api.get<StatsResponse>('/realtime/stats', { params })
    return response.data.stats
  }

  // Get leaderboard via REST API
  async getLeaderboard(classId?: number): Promise<LeaderboardEntry[]> {
    const params: Record<string, string> = {}
    if (classId) {
      params.class_id = classId.toString()
    }
    const response = await api.get<LeaderboardResponse>('/realtime/leaderboard', { params })
    return response.data.leaderboard || []
  }
}

// Export singleton instance
export const realtimeService = new RealtimeService()

export default realtimeService
