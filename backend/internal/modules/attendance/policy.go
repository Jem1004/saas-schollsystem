package attendance

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

// AttendancePolicy determines attendance rules based on school settings
// Requirements: 5.2 - Status determination based on school settings
type AttendancePolicy interface {
	// DetermineAttendanceStatus determines the attendance status based on check-in time
	// Requirements: Property 17 - School Settings Policy Enforcement
	DetermineAttendanceStatus(checkInTime time.Time, schoolID uint) models.AttendanceStatus
	
	// IsWithinAttendanceWindow checks if the given time is within the attendance window
	IsWithinAttendanceWindow(schoolID uint, timestamp time.Time) bool
	
	// ShouldSendNotification checks if a notification should be sent for the event type
	ShouldSendNotification(schoolID uint, eventType models.NotificationType) bool
	
	// GetAttendanceTimeWindow returns the attendance time window for a school
	GetAttendanceTimeWindow(schoolID uint, date time.Time) (*models.AttendanceTimeWindow, error)
}

// policy implements the AttendancePolicy interface
type policy struct {
	db *gorm.DB
}

// NewAttendancePolicy creates a new attendance policy
func NewAttendancePolicy(db *gorm.DB) AttendancePolicy {
	return &policy{db: db}
}

// DetermineAttendanceStatus determines the attendance status based on check-in time and school settings
// Requirements: 5.2 - THE System SHALL record check-in or check-out based on existing records
// Requirements: Property 17 - Attendance status SHALL be determined based on school's configured time thresholds
func (p *policy) DetermineAttendanceStatus(checkInTime time.Time, schoolID uint) models.AttendanceStatus {
	settings := p.getSchoolSettings(schoolID)
	
	// Parse attendance start time
	startTime, err := parseTimeString(settings.AttendanceStartTime)
	if err != nil {
		// Default to 07:00 if parsing fails
		startTime = time.Date(0, 1, 1, 7, 0, 0, 0, time.UTC)
	}

	// Create time window for the check-in date
	checkInDate := time.Date(checkInTime.Year(), checkInTime.Month(), checkInTime.Day(), 0, 0, 0, 0, checkInTime.Location())
	
	// Calculate threshold times for the check-in date
	attendanceStart := time.Date(
		checkInDate.Year(), checkInDate.Month(), checkInDate.Day(),
		startTime.Hour(), startTime.Minute(), 0, 0, checkInTime.Location(),
	)
	
	lateThreshold := attendanceStart.Add(time.Duration(settings.AttendanceLateThreshold) * time.Minute)
	veryLateThreshold := attendanceStart.Add(time.Duration(settings.AttendanceVeryLateThreshold) * time.Minute)

	// Determine status based on check-in time
	// Requirements: Property 7 - Attendance Check-In/Check-Out Logic
	switch {
	case checkInTime.Before(lateThreshold) || checkInTime.Equal(lateThreshold):
		return models.AttendanceStatusOnTime
	case checkInTime.Before(veryLateThreshold) || checkInTime.Equal(veryLateThreshold):
		return models.AttendanceStatusLate
	default:
		return models.AttendanceStatusVeryLate
	}
}

// IsWithinAttendanceWindow checks if the given time is within the attendance window
func (p *policy) IsWithinAttendanceWindow(schoolID uint, timestamp time.Time) bool {
	settings := p.getSchoolSettings(schoolID)
	
	// Parse start and end times
	startTime, err := parseTimeString(settings.AttendanceStartTime)
	if err != nil {
		startTime = time.Date(0, 1, 1, 7, 0, 0, 0, time.UTC)
	}
	
	endTime, err := parseTimeString(settings.AttendanceEndTime)
	if err != nil {
		endTime = time.Date(0, 1, 1, 17, 0, 0, 0, time.UTC) // Default end at 5 PM
	}

	// Create time window for the timestamp date
	date := time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, timestamp.Location())
	
	windowStart := time.Date(
		date.Year(), date.Month(), date.Day(),
		startTime.Hour(), startTime.Minute(), 0, 0, timestamp.Location(),
	)
	windowEnd := time.Date(
		date.Year(), date.Month(), date.Day(),
		endTime.Hour(), endTime.Minute(), 0, 0, timestamp.Location(),
	)

	// Check if timestamp is within window
	return (timestamp.After(windowStart) || timestamp.Equal(windowStart)) &&
		(timestamp.Before(windowEnd) || timestamp.Equal(windowEnd))
}

// ShouldSendNotification checks if a notification should be sent based on school settings
// Requirements: Property 17 - Notifications SHALL only be sent if the corresponding notification setting is enabled
func (p *policy) ShouldSendNotification(schoolID uint, eventType models.NotificationType) bool {
	settings := p.getSchoolSettings(schoolID)
	return settings.ShouldSendNotification(eventType)
}

// GetAttendanceTimeWindow returns the attendance time window for a school on a specific date
func (p *policy) GetAttendanceTimeWindow(schoolID uint, date time.Time) (*models.AttendanceTimeWindow, error) {
	settings := p.getSchoolSettings(schoolID)
	
	// Parse start time
	startTime, err := parseTimeString(settings.AttendanceStartTime)
	if err != nil {
		startTime = time.Date(0, 1, 1, 7, 0, 0, 0, time.UTC)
	}
	
	// Parse end time
	endTime, err := parseTimeString(settings.AttendanceEndTime)
	if err != nil {
		endTime = time.Date(0, 1, 1, 7, 30, 0, 0, time.UTC)
	}

	// Create time window for the given date
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	
	windowStart := time.Date(
		dateOnly.Year(), dateOnly.Month(), dateOnly.Day(),
		startTime.Hour(), startTime.Minute(), 0, 0, date.Location(),
	)
	windowEnd := time.Date(
		dateOnly.Year(), dateOnly.Month(), dateOnly.Day(),
		endTime.Hour(), endTime.Minute(), 0, 0, date.Location(),
	)
	lateTime := windowStart.Add(time.Duration(settings.AttendanceLateThreshold) * time.Minute)
	veryLateTime := windowStart.Add(time.Duration(settings.AttendanceVeryLateThreshold) * time.Minute)

	return &models.AttendanceTimeWindow{
		StartTime:    windowStart,
		EndTime:      windowEnd,
		LateTime:     lateTime,
		VeryLateTime: veryLateTime,
	}, nil
}

// getSchoolSettings retrieves school settings from database or returns defaults
func (p *policy) getSchoolSettings(schoolID uint) *models.SchoolSettings {
	var settings models.SchoolSettings
	
	err := p.db.WithContext(context.Background()).
		Where("school_id = ?", schoolID).
		First(&settings).Error
	
	if err != nil {
		// Return default settings if not found
		return models.DefaultSchoolSettings(schoolID)
	}
	
	return &settings
}

// parseTimeString parses a time string in HH:MM format
func parseTimeString(timeStr string) (time.Time, error) {
	return time.Parse("15:04", timeStr)
}
