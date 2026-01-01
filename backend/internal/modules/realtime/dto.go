package realtime

import (
	"time"

	"github.com/school-management/backend/internal/domain/models"
)

// ==================== Event Types ====================

// EventType represents the type of real-time event
type EventType string

const (
	EventTypeNewAttendance EventType = "new_attendance"
	EventTypeStatsUpdate   EventType = "stats_update"
	EventTypeLeaderboard   EventType = "leaderboard_update"
)

// ==================== Real-Time DTOs ====================

// AttendanceEvent represents a real-time attendance event
// Requirements: 4.2 - WHEN a student taps RFID card, THE System SHALL update the dashboard within 3 seconds
type AttendanceEvent struct {
	Type        EventType            `json:"type"`
	SchoolID    uint                 `json:"school_id"`
	Attendance  *LiveFeedEntry       `json:"attendance,omitempty"`
	Stats       *AttendanceStats     `json:"stats,omitempty"`
	Leaderboard []LeaderboardEntry   `json:"leaderboard,omitempty"`
}

// LiveFeedEntry represents a single entry in the live attendance feed
// Requirements: 4.3 - Show the 20 most recent attendance records with student name, class, time, and status
type LiveFeedEntry struct {
	ID          uint                    `json:"id"`
	StudentID   uint                    `json:"student_id"`
	StudentName string                  `json:"student_name"`
	ClassName   string                  `json:"class_name"`
	ClassID     uint                    `json:"class_id"`
	Time        time.Time               `json:"time"`
	Status      models.AttendanceStatus `json:"status"`
	Type        string                  `json:"type"` // "check_in" or "check_out"
}

// LeaderboardEntry represents a single entry in the leaderboard
// Requirements: 5.6 - Show leaderboard of top 10 earliest arrivals for the day
type LeaderboardEntry struct {
	Rank        int       `json:"rank"`
	StudentID   uint      `json:"student_id"`
	StudentName string    `json:"student_name"`
	ClassName   string    `json:"class_name"`
	ArrivalTime time.Time `json:"arrival_time"`
}

// AttendanceStats represents real-time attendance statistics
// Requirements: 4.1 - Display current day's attendance statistics (present, late, very late, absent count)
// Requirements: 4.10 - Show percentage of attendance completion (attended/total students)
type AttendanceStats struct {
	TotalStudents int     `json:"total_students"`
	Present       int     `json:"present"`
	Late          int     `json:"late"`
	VeryLate      int     `json:"very_late"`
	Absent        int     `json:"absent"`
	Percentage    float64 `json:"percentage"` // (present / total_students) * 100
}

// ==================== Request/Response DTOs ====================

// LiveFeedResponse represents the response for live feed request
type LiveFeedResponse struct {
	Feed []LiveFeedEntry `json:"feed"`
}

// StatsResponse represents the response for stats request
type StatsResponse struct {
	Stats AttendanceStats `json:"stats"`
	Date  string          `json:"date"` // Format: YYYY-MM-DD
}

// LeaderboardResponse represents the response for leaderboard request
type LeaderboardResponse struct {
	Leaderboard []LeaderboardEntry `json:"leaderboard"`
	Date        string             `json:"date"` // Format: YYYY-MM-DD
}

// ==================== WebSocket DTOs ====================

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// WSSubscribeRequest represents a subscription request
type WSSubscribeRequest struct {
	ClassID *uint `json:"class_id,omitempty"` // Optional filter by class
}

// WSConnectionStatus represents connection status
type WSConnectionStatus struct {
	Connected bool   `json:"connected"`
	Message   string `json:"message,omitempty"`
}
