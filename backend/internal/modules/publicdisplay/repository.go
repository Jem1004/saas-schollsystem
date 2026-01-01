package publicdisplay

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

// Repository defines the interface for public display data operations
type Repository interface {
	// GetSchoolByID retrieves school information by ID
	GetSchoolByID(ctx context.Context, schoolID uint) (*models.School, error)

	// GetPublicLiveFeed retrieves the most recent attendance records for public display
	// Requirements: 5.4 - Show live feed of recent attendance (last 10 records)
	GetPublicLiveFeed(ctx context.Context, schoolID uint, limit int) ([]PublicLiveFeedEntry, error)

	// GetPublicStats retrieves attendance statistics for public display
	// Requirements: 5.5 - Show real-time statistics
	GetPublicStats(ctx context.Context, schoolID uint, date time.Time) (*PublicAttendanceStats, error)

	// GetPublicLeaderboard retrieves top earliest arrivals for public display
	// Requirements: 5.6 - Show leaderboard of top 10 earliest arrivals
	GetPublicLeaderboard(ctx context.Context, schoolID uint, limit int, date time.Time) ([]PublicLeaderboardEntry, error)

	// GetTotalStudents retrieves total active students count
	GetTotalStudents(ctx context.Context, schoolID uint) (int, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new public display repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// GetSchoolByID retrieves school information by ID
func (r *repository) GetSchoolByID(ctx context.Context, schoolID uint) (*models.School, error) {
	var school models.School
	err := r.db.WithContext(ctx).
		Where("id = ?", schoolID).
		First(&school).Error

	if err != nil {
		return nil, err
	}

	return &school, nil
}

// GetPublicLiveFeed retrieves the most recent attendance records for public display
// Requirements: 5.4 - Show live feed of recent attendance (last 10 records)
// Requirements: 5.14 - Only expose name, class, and attendance time (NO NIS, NISN)
func (r *repository) GetPublicLiveFeed(ctx context.Context, schoolID uint, limit int) ([]PublicLiveFeedEntry, error) {
	if limit <= 0 {
		limit = 10
	}

	// Get today's date
	today := time.Now()
	dateOnly := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	var attendances []models.Attendance

	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Joins("JOIN students ON students.id = attendances.student_id").
		Where("students.school_id = ?", schoolID).
		Where("attendances.date = ?", dateOnly).
		Where("attendances.check_in_time IS NOT NULL").
		Order("attendances.check_in_time DESC").
		Limit(limit).
		Find(&attendances).Error

	if err != nil {
		return nil, err
	}

	// Convert to PublicLiveFeedEntry - only expose safe data
	feed := make([]PublicLiveFeedEntry, 0, len(attendances))
	for _, a := range attendances {
		entry := PublicLiveFeedEntry{
			StudentName: a.Student.Name,
			Status:      string(a.Status),
			Type:        "check_in",
		}

		if a.Student.Class.ID != 0 {
			entry.ClassName = a.Student.Class.Name
		}

		if a.CheckInTime != nil {
			entry.Time = *a.CheckInTime
		}

		feed = append(feed, entry)
	}

	return feed, nil
}

// GetPublicStats retrieves attendance statistics for public display
// Requirements: 5.5 - Show real-time statistics (present, late, absent count)
func (r *repository) GetPublicStats(ctx context.Context, schoolID uint, date time.Time) (*PublicAttendanceStats, error) {
	// Normalize date to start of day
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	// Get total active students
	totalStudents, err := r.GetTotalStudents(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	// Get attendance counts by status
	type StatusCount struct {
		Status string
		Count  int
	}
	var statusCounts []StatusCount

	err = r.db.WithContext(ctx).
		Model(&models.Attendance{}).
		Select("attendances.status, COUNT(*) as count").
		Joins("JOIN students ON students.id = attendances.student_id").
		Where("students.school_id = ?", schoolID).
		Where("attendances.date = ?", dateOnly).
		Group("attendances.status").
		Scan(&statusCounts).Error

	if err != nil {
		return nil, err
	}

	// Build stats
	stats := &PublicAttendanceStats{
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

// GetPublicLeaderboard retrieves top earliest arrivals for public display
// Requirements: 5.6 - Show leaderboard of top 10 earliest arrivals for the day
// Requirements: 5.14 - Only expose name, class, and arrival time (NO NIS, NISN)
func (r *repository) GetPublicLeaderboard(ctx context.Context, schoolID uint, limit int, date time.Time) ([]PublicLeaderboardEntry, error) {
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
		Order("attendances.check_in_time ASC").                         // Earliest first
		Limit(limit).
		Find(&attendances).Error

	if err != nil {
		return nil, err
	}

	// Convert to PublicLeaderboardEntry - only expose safe data
	leaderboard := make([]PublicLeaderboardEntry, 0, len(attendances))
	for i, a := range attendances {
		entry := PublicLeaderboardEntry{
			Rank:        i + 1,
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
func (r *repository) GetTotalStudents(ctx context.Context, schoolID uint) (int, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("school_id = ? AND is_active = ?", schoolID, true).
		Count(&count).Error

	return int(count), err
}
