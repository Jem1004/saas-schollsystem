package publicdisplay

import (
	"time"
)

// ==================== Public Display DTOs ====================
// Requirements: 5.4, 5.5, 5.6, 5.7, 5.8, 5.14 - Public display data without sensitive information

// PublicDisplayData represents the data shown on public display screens
// Requirements: 5.4, 5.5, 5.6, 5.7, 5.8 - Display live feed, stats, leaderboard, date/time, school name
// Requirements: 5.14 - SHALL NOT expose sensitive student information (only name, class, and attendance time)
type PublicDisplayData struct {
	SchoolName  string                   `json:"school_name"`
	CurrentTime time.Time                `json:"current_time"`
	Date        string                   `json:"date"` // Format: "Senin, 1 Januari 2024"
	Stats       PublicAttendanceStats    `json:"stats"`
	LiveFeed    []PublicLiveFeedEntry    `json:"live_feed"`
	Leaderboard []PublicLeaderboardEntry `json:"leaderboard"`
}

// PublicAttendanceStats represents attendance statistics for public display
// Requirements: 5.5 - Show real-time statistics (present, late, absent count)
type PublicAttendanceStats struct {
	TotalStudents int     `json:"total_students"`
	Present       int     `json:"present"`
	Late          int     `json:"late"`
	VeryLate      int     `json:"very_late"`
	Absent        int     `json:"absent"`
	Percentage    float64 `json:"percentage"`
}

// PublicLiveFeedEntry represents a single entry in the public live feed
// Requirements: 5.4 - Show live feed of recent attendance (last 10 records)
// Requirements: 5.14 - Only expose name, class, and attendance time (NO NIS, NISN)
type PublicLiveFeedEntry struct {
	StudentName string    `json:"student_name"`
	ClassName   string    `json:"class_name"`
	Time        time.Time `json:"time"`
	Status      string    `json:"status"` // "on_time", "late", "very_late"
	Type        string    `json:"type"`   // "check_in", "check_out"
}

// PublicLeaderboardEntry represents a single entry in the public leaderboard
// Requirements: 5.6 - Show leaderboard of top 10 earliest arrivals for the day
// Requirements: 5.14 - Only expose name, class, and arrival time (NO NIS, NISN)
type PublicLeaderboardEntry struct {
	Rank        int       `json:"rank"`
	StudentName string    `json:"student_name"`
	ClassName   string    `json:"class_name"`
	ArrivalTime time.Time `json:"arrival_time"`
}

// ==================== WebSocket DTOs ====================

// PublicWSMessage represents a WebSocket message for public display
type PublicWSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// PublicWSConnectionStatus represents connection status for public display
type PublicWSConnectionStatus struct {
	Connected bool   `json:"connected"`
	Message   string `json:"message,omitempty"`
}

// ==================== Error Response ====================

// PublicDisplayError represents an error response for public display
// Requirements: 5.13 - IF display token is invalid or revoked, show error message
type PublicDisplayError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
