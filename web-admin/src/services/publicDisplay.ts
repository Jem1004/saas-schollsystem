// Public Display Service
// Requirements: 5.3 - Accessing public display URL with valid token SHALL show attendance data without login
// Requirements: 5.9 - Update public display within 3 seconds when new attendance recorded

// Types for public display data
export interface PublicAttendanceStats {
  totalStudents: number
  present: number
  late: number
  veryLate: number
  absent: number
  percentage: number
}

export interface PublicLiveFeedEntry {
  studentName: string
  className: string
  time: string
  status: string
  type: string
}

export interface PublicLeaderboardEntry {
  rank: number
  studentName: string
  className: string
  arrivalTime: string
}

export interface PublicDisplayData {
  schoolName: string
  currentTime: string
  date: string
  stats: PublicAttendanceStats
  liveFeed: PublicLiveFeedEntry[]
  leaderboard: PublicLeaderboardEntry[]
}

export interface PublicWSConnectionStatus {
  connected: boolean
  message: string
}

export interface PublicDisplayError {
  error: string
  message: string
}

// Event handler types
export type PublicDisplayEventHandler = (data: any) => void
export type PublicConnectionStatusHandler = (status: PublicWSConnectionStatus) => void

// WebSocket configuration
const WS_RECONNECT_DELAY = 3000 // 3 seconds
const WS_MAX_RECONNECT_ATTEMPTS = 10
const WS_PING_INTERVAL = 30000 // 30 seconds

class PublicDisplayService {
  private ws: WebSocket | null = null
  private token: string = ''
  private reconnectAttempts = 0
  private reconnectTimeout: ReturnType<typeof setTimeout> | null = null
  private pingInterval: ReturnType<typeof setInterval> | null = null
  private eventHandlers: Set<PublicDisplayEventHandler> = new Set()
  private statusHandlers: Set<PublicConnectionStatusHandler> = new Set()
  private connectionState: 'connecting' | 'connected' | 'disconnected' | 'reconnecting' = 'disconnected'

  // Get API base URL
  private getApiBaseUrl(): string {
    // @ts-expect-error - Vite env types
    const apiUrl = import.meta.env?.VITE_API_URL as string | undefined
    if (apiUrl) {
      return apiUrl.replace(/\/api\/v1\/?$/, '')
    }
    return `${window.location.protocol}//${window.location.host}`
  }

  // Get WebSocket URL for public display
  private getWsUrl(): string {
    const baseUrl = this.getApiBaseUrl()
    const protocol = baseUrl.startsWith('https') ? 'wss:' : 'ws:'
    const host = baseUrl.replace(/^https?:\/\//, '')
    return `${protocol}//${host}/api/v1/public/display/${this.token}/ws`
  }

  // Fetch display data via REST API
  // Requirements: 5.3 - Access public display with valid token
  async fetchDisplayData(token: string): Promise<{ success: boolean; data?: PublicDisplayData; error?: PublicDisplayError }> {
    try {
      const baseUrl = this.getApiBaseUrl()
      const response = await fetch(`${baseUrl}/api/v1/public/display/${token}`)
      const result = await response.json()

      if (!response.ok) {
        return {
          success: false,
          error: {
            error: result.error || 'UNKNOWN_ERROR',
            message: result.message || 'Gagal memuat data display',
          },
        }
      }

      if (result.success && result.data) {
        return {
          success: true,
          data: this.transformDisplayData(result.data),
        }
      }

      return {
        success: false,
        error: {
          error: 'INVALID_RESPONSE',
          message: 'Response tidak valid dari server',
        },
      }
    } catch (err) {
      console.error('Failed to fetch display data:', err)
      return {
        success: false,
        error: {
          error: 'NETWORK_ERROR',
          message: 'Gagal terhubung ke server. Periksa koneksi internet.',
        },
      }
    }
  }

  // Transform API response to frontend format
  private transformDisplayData(data: any): PublicDisplayData {
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
  connect(token: string): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      console.log('Public display WebSocket already connected')
      return
    }

    this.token = token
    this.connectionState = 'connecting'
    this.notifyStatusChange({ connected: false, message: 'Menghubungkan...' })

    try {
      this.ws = new WebSocket(this.getWsUrl())
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
      console.log('Public display WebSocket connected')
      this.connectionState = 'connected'
      this.reconnectAttempts = 0
      this.notifyStatusChange({ connected: true, message: 'Terhubung' })
      this.startPingInterval()
    }

    this.ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        this.notifyEventHandlers(data)
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error)
      }
    }

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error)
      this.notifyStatusChange({ connected: false, message: 'Kesalahan koneksi' })
    }

    this.ws.onclose = (event) => {
      console.log('WebSocket closed:', event.code, event.reason)
      this.connectionState = 'disconnected'
      this.stopPingInterval()

      if (event.code !== 1000) {
        this.handleReconnect()
      } else {
        this.notifyStatusChange({ connected: false, message: 'Terputus' })
      }
    }
  }

  // Handle reconnection logic
  private handleReconnect(): void {
    if (this.reconnectAttempts >= WS_MAX_RECONNECT_ATTEMPTS) {
      console.error('Max reconnection attempts reached')
      this.notifyStatusChange({
        connected: false,
        message: 'Gagal terhubung. Silakan refresh halaman.',
      })
      return
    }

    this.connectionState = 'reconnecting'
    this.reconnectAttempts++

    const delay = WS_RECONNECT_DELAY * Math.min(this.reconnectAttempts, 5)
    this.notifyStatusChange({
      connected: false,
      message: `Menghubungkan ulang (${this.reconnectAttempts}/${WS_MAX_RECONNECT_ATTEMPTS})...`,
    })

    this.reconnectTimeout = setTimeout(() => {
      console.log(`Reconnection attempt ${this.reconnectAttempts}`)
      this.connect(this.token)
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

  // Request data refresh via WebSocket
  requestRefresh(): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ type: 'refresh' }))
    }
  }

  // Subscribe to display events
  onDisplayEvent(handler: PublicDisplayEventHandler): () => void {
    this.eventHandlers.add(handler)
    return () => {
      this.eventHandlers.delete(handler)
    }
  }

  // Subscribe to connection status changes
  onConnectionStatus(handler: PublicConnectionStatusHandler): () => void {
    this.statusHandlers.add(handler)
    return () => {
      this.statusHandlers.delete(handler)
    }
  }

  // Notify all event handlers
  private notifyEventHandlers(event: any): void {
    this.eventHandlers.forEach((handler) => {
      try {
        handler(event)
      } catch (error) {
        console.error('Error in event handler:', error)
      }
    })
  }

  // Notify all status handlers
  private notifyStatusChange(status: PublicWSConnectionStatus): void {
    this.statusHandlers.forEach((handler) => {
      try {
        handler(status)
      } catch (error) {
        console.error('Error in status handler:', error)
      }
    })
  }

  // Get current connection state
  getConnectionState(): string {
    return this.connectionState
  }

  // Check if connected
  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN
  }
}

// Export singleton instance
export const publicDisplayService = new PublicDisplayService()

export default publicDisplayService
