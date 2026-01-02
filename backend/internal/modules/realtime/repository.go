package realtime

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

// Repository defines the interface for real-time data operations
type Repository interface {
	// GetLiveFeed retrieves the most recent attendance records
	// Requirements: 4.3 - Show the 20 most recent attendance records
	GetLiveFeed(ctx context.Context, schoolID uint, classID *uint, limit int) ([]LiveFeedEntry, error)

	// GetAttendanceStats retrieves current day's attendance statistics
	// Requirements: 4.1 - Display current day's attendance statistics
	GetAttendanceStats(ctx context.Context, schoolID uint, classID *uint, date time.Time) (*AttendanceStats, error)

	// GetLeaderboard retrieves top earliest arrivals for the day
	// Requirements: 5.6 - Show leaderboard of top 10 earliest arrivals
	GetLeaderboard(ctx context.Context, schoolID uint, limit int, date time.Time) ([]LeaderboardEntry, error)

	// GetTotalStudents retrieves total active students count
	GetTotalStudents(ctx context.Context, schoolID uint, classID *uint) (int, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new real-time repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// GetLiveFeed retrieves the most recent attendance records
// Requirements: 4.3 - Show the 20 most recent attendance records with student name, class, time, and status
func (r *repository) GetLiveFeed(ctx context.Context, schoolID uint, classID *uint, limit int) ([]LiveFeedEntry, error) {
	if limit <= 0 {
		limit = 20
	}

	// Get today's date
	today := time.Now()
	dateOnly := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	var attendances []models.Attendance

	query := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Joins("JOIN students ON students.id = attendances.student_id").
		Where("students.school_id = ?", schoolID).
		Where("attendances.date = ?", dateOnly).
		Where("attendances.check_in_time IS NOT NULL")

	// Apply class filter if provided
	if classID != nil {
		query = query.Where("students.class_id = ?", *classID)
	}

	// Order by check-in time descending (most recent first)
	err := query.Order("attendances.check_in_time DESC").
		Limit(limit).
		Find(&attendances).Error

	if err != nil {
		return nil, err
	}

	// Convert to LiveFeedEntry
	feed := make([]LiveFeedEntry, 0, len(attendances))
	for _, a := range attendances {
		entry := LiveFeedEntry{
			ID:          a.ID,
			StudentID:   a.StudentID,
			StudentName: a.Student.Name,
			Status:      a.Status,
			Type:        "check_in",
		}

		// Handle nullable Class pointer
		if a.Student.Class != nil && a.Student.Class.ID != 0 {
			entry.ClassName = a.Student.Class.Name
			entry.ClassID = a.Student.Class.ID
		}

		if a.CheckInTime != nil {
			entry.Time = *a.CheckInTime
		}

		feed = append(feed, entry)
	}

	return feed, nil
}

// GetAttendanceStats retrieves current day's attendance statistics
// Requirements: 4.1 - Display current day's attendance statistics (present, late, very late, absent count)
// Requirements: 4.10 - Show percentage of attendance completion
func (r *repository) GetAttendanceStats(ctx context.Context, schoolID uint, classID *uint, date time.Time) (*AttendanceStats, error) {
	// Normalize date to start of day
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	// Get total active students
	totalStudents, err := r.GetTotalStudents(ctx, schoolID, classID)
	if err != nil {
		return nil, err
	}

	// Get attendance counts by status
	type StatusCount struct {
		Status string
		Count  int
	}
	var statusCounts []StatusCount

	query := r.db.WithContext(ctx).
		Model(&models.Attendance{}).
		Select("attendances.status, COUNT(*) as count").
		Joins("JOIN students ON students.id = attendances.student_id").
		Where("students.school_id = ?", schoolID).
		Where("attendances.date = ?", dateOnly)

	// Apply class filter if provided
	if classID != nil {
		query = query.Where("students.class_id = ?", *classID)
	}

	err = query.Group("attendances.status").Scan(&statusCounts).Error
	if err != nil {
		return nil, err
	}

	// Build stats
	stats := &AttendanceStats{
		TotalStudents: totalStudents,
	}

	var presentCount int
	for _, sc := range statusCounts {
		switch models.AttendanceStatus(sc.Status) {
		case models.AttendanceStatusOnTime:
			stats.Present += sc.Count
			presentCount += sc.Count
		case models.AttendanceStatusLate:
			stats.Late = sc.Count
			presentCount += sc.Count
		case models.AttendanceStatusVeryLate:
			stats.VeryLate = sc.Count
			presentCount += sc.Count
		}
	}

	// Calculate absent (students without attendance record)
	stats.Absent = totalStudents - presentCount
	if stats.Absent < 0 {
		stats.Absent = 0
	}

	// Calculate percentage
	if totalStudents > 0 {
		stats.Percentage = float64(presentCount) / float64(totalStudents) * 100
	}

	return stats, nil
}

// GetLeaderboard retrieves top earliest arrivals for the day
// Requirements: 5.6 - Show leaderboard of top 10 earliest arrivals for the day
func (r *repository) GetLeaderboard(ctx context.Context, schoolID uint, limit int, date time.Time) ([]LeaderboardEntry, error) {
	if limit <= 0 {
		limit = 10
	}

	// Normalize date to start of day
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	var attendances []models.Attendance

	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Joins("JOIN students ON students.id = attendances.student_id").
		Where("students.school_id = ?", schoolID).
		Where("attendances.date = ?", dateOnly).
		Where("attendances.check_in_time IS NOT NULL").
		Where("attendances.status = ?", models.AttendanceStatusOnTime). // Only on-time arrivals
		Order("attendances.check_in_time ASC"). // Earliest first
		Limit(limit).
		Find(&attendances).Error

	if err != nil {
		return nil, err
	}

	// Convert to LeaderboardEntry
	leaderboard := make([]LeaderboardEntry, 0, len(attendances))
	for i, a := range attendances {
		entry := LeaderboardEntry{
			Rank:        i + 1,
			StudentID:   a.StudentID,
			StudentName: a.Student.Name,
		}

		if a.Student.Class.ID != 0 {
			entry.ClassName = a.Student.Class.Name
		}

		if a.CheckInTime != nil {
			entry.ArrivalTime = *a.CheckInTime
		}

		leaderboard = append(leaderboard, entry)
	}

	return leaderboard, nil
}

// GetTotalStudents retrieves total active students count
func (r *repository) GetTotalStudents(ctx context.Context, schoolID uint, classID *uint) (int, error) {
	var count int64

	query := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("school_id = ? AND is_active = ?", schoolID, true)

	// Apply class filter if provided
	if classID != nil {
		query = query.Where("class_id = ?", *classID)
	}

	err := query.Count(&count).Error
	return int(count), err
}
