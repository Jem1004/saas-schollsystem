// Real-time attendance types
// Requirements: 4.1 - Display current day's attendance statistics

import type { AttendanceStatus } from './attendance'

// Event types for WebSocket messages
export type EventType = 'new_attendance' | 'stats_update' | 'leaderboard_update'

// Live feed entry representing a single attendance record
// Requirements: 4.3 - Show the 20 most recent attendance records with student name, class, time, and status
export interface LiveFeedEntry {
  id: number
  studentId: number
  studentName: string
  className: string
  classId: number
  time: string
  status: AttendanceStatus
  type: 'check_in' | 'check_out'
}

// Leaderboard entry for earliest arrivals
// Requirements: 5.6 - Show leaderboard of top 10 earliest arrivals
export interface LeaderboardEntry {
  rank: number
  studentId: number
  studentName: string
  className: string
  arrivalTime: string
}

// Real-time attendance statistics
// Requirements: 4.1 - Display current day's attendance statistics (present, late, very late, absent count)
// Requirements: 4.10 - Show percentage of attendance completion
export interface AttendanceStats {
  totalStudents: number
  present: number
  late: number
  veryLate: number
  absent: number
  percentage: number
}

// WebSocket attendance event
// Requirements: 4.2 - Update dashboard within 3 seconds without page refresh
export interface AttendanceEvent {
  type: EventType
  schoolId: number
  attendance?: LiveFeedEntry
  stats?: AttendanceStats
  leaderboard?: LeaderboardEntry[]
}

// WebSocket message wrapper
export interface WSMessage {
  type: string
  payload: unknown
}

// WebSocket subscription request
export interface WSSubscribeRequest {
  classId?: number
}

// WebSocket connection status
// Requirements: 4.9 - Display connection status and attempt to reconnect automatically
export interface WSConnectionStatus {
  connected: boolean
  message?: string
}

// Live feed response from REST API
export interface LiveFeedResponse {
  feed: LiveFeedEntry[]
}

// Stats response from REST API
export interface StatsResponse {
  stats: AttendanceStats
  date: string
}

// Leaderboard response from REST API
export interface LeaderboardResponse {
  leaderboard: LeaderboardEntry[]
  date: string
}
