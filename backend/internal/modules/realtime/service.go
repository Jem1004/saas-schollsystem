package realtime

import (
	"context"
	"time"

	"github.com/school-management/backend/internal/domain/models"
)

// Service defines the interface for real-time operations
type Service interface {
	// GetLiveFeed retrieves the most recent attendance records
	// Requirements: 4.3 - Show the 20 most recent attendance records
	GetLiveFeed(ctx context.Context, schoolID uint, classID *uint) (*LiveFeedResponse, error)

	// GetAttendanceStats retrieves current day's attendance statistics
	// Requirements: 4.1 - Display current day's attendance statistics
	GetAttendanceStats(ctx context.Context, schoolID uint, classID *uint) (*StatsResponse, error)

	// GetLeaderboard retrieves top 10 earliest arrivals for the day
	// Requirements: 5.6 - Show leaderboard of top 10 earliest arrivals
	GetLeaderboard(ctx context.Context, schoolID uint) (*LeaderboardResponse, error)

	// BroadcastAttendance broadcasts a new attendance event to all connected clients
	// Requirements: 4.2 - Update dashboard within 3 seconds without page refresh
	BroadcastAttendance(ctx context.Context, schoolID uint, attendance *models.Attendance, student *models.Student, attendanceType string)

	// GetHub returns the WebSocket hub
	GetHub() *Hub
}

// service implements the Service interface
type service struct {
	repo Repository
	hub  *Hub
}

// NewService creates a new real-time service
func NewService(repo Repository, hub *Hub) Service {
	return &service{
		repo: repo,
		hub:  hub,
	}
}

// GetLiveFeed retrieves the most recent attendance records
// Requirements: 4.3 - Show the 20 most recent attendance records with student name, class, time, and status
func (s *service) GetLiveFeed(ctx context.Context, schoolID uint, classID *uint) (*LiveFeedResponse, error) {
	feed, err := s.repo.GetLiveFeed(ctx, schoolID, classID, 20)
	if err != nil {
		return nil, err
	}

	return &LiveFeedResponse{
		Feed: feed,
	}, nil
}

// GetAttendanceStats retrieves current day's attendance statistics
// Requirements: 4.1 - Display current day's attendance statistics (present, late, very late, absent count)
// Requirements: 4.10 - Show percentage of attendance completion
func (s *service) GetAttendanceStats(ctx context.Context, schoolID uint, classID *uint) (*StatsResponse, error) {
	today := time.Now()
	stats, err := s.repo.GetAttendanceStats(ctx, schoolID, classID, today)
	if err != nil {
		return nil, err
	}

	return &StatsResponse{
		Stats: *stats,
		Date:  today.Format("2006-01-02"),
	}, nil
}

// GetLeaderboard retrieves top 10 earliest arrivals for the day
// Requirements: 5.6 - Show leaderboard of top 10 earliest arrivals for the day
func (s *service) GetLeaderboard(ctx context.Context, schoolID uint) (*LeaderboardResponse, error) {
	today := time.Now()
	leaderboard, err := s.repo.GetLeaderboard(ctx, schoolID, 10, today)
	if err != nil {
		return nil, err
	}

	return &LeaderboardResponse{
		Leaderboard: leaderboard,
		Date:        today.Format("2006-01-02"),
	}, nil
}

// BroadcastAttendance broadcasts a new attendance event to all connected clients
// Requirements: 4.2 - WHEN a student taps RFID card, THE System SHALL update the dashboard within 3 seconds
func (s *service) BroadcastAttendance(ctx context.Context, schoolID uint, attendance *models.Attendance, student *models.Student, attendanceType string) {
	if s.hub == nil {
		return
	}

	// Create live feed entry
	liveFeedEntry := &LiveFeedEntry{
		ID:          attendance.ID,
		StudentID:   attendance.StudentID,
		StudentName: student.Name,
		Status:      attendance.Status,
		Type:        attendanceType,
	}

	// Handle nullable Class pointer
	if student.Class != nil && student.Class.ID != 0 {
		liveFeedEntry.ClassName = student.Class.Name
		liveFeedEntry.ClassID = student.Class.ID
	}

	if attendanceType == "check_in" && attendance.CheckInTime != nil {
		liveFeedEntry.Time = *attendance.CheckInTime
	} else if attendanceType == "check_out" && attendance.CheckOutTime != nil {
		liveFeedEntry.Time = *attendance.CheckOutTime
	}

	// Get updated stats
	stats, err := s.repo.GetAttendanceStats(ctx, schoolID, nil, time.Now())
	if err != nil {
		stats = nil
	}

	// Get updated leaderboard
	leaderboard, err := s.repo.GetLeaderboard(ctx, schoolID, 10, time.Now())
	if err != nil {
		leaderboard = nil
	}

	// Create and broadcast event
	event := &AttendanceEvent{
		Type:        EventTypeNewAttendance,
		SchoolID:    schoolID,
		Attendance:  liveFeedEntry,
		Stats:       stats,
		Leaderboard: leaderboard,
	}

	s.hub.Broadcast(event)
}

// GetHub returns the WebSocket hub
func (s *service) GetHub() *Hub {
	return s.hub
}
