package models

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// AttendanceSchedule represents a configurable attendance time slot
// Requirements: 3.1, 3.2, 3.3 - Multi-schedule support for different activities
type AttendanceSchedule struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	SchoolID          uint       `gorm:"index;not null" json:"school_id"`
	Name              string     `gorm:"type:varchar(100);not null" json:"name"`
	StartTime         string     `gorm:"type:time without time zone;not null" json:"start_time"`
	EndTime           string     `gorm:"type:time without time zone;not null" json:"end_time"`
	LateThreshold     int        `gorm:"not null;default:15" json:"late_threshold"`
	VeryLateThreshold *int       `gorm:"" json:"very_late_threshold"`
	DaysOfWeek        string     `gorm:"type:varchar(20);default:'1,2,3,4,5'" json:"days_of_week"`
	IsActive          bool       `gorm:"default:true" json:"is_active"`
	IsDefault         bool       `gorm:"default:false" json:"is_default"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// Relations
	School School `gorm:"foreignKey:SchoolID;constraint:OnDelete:CASCADE" json:"school,omitempty"`
}

// TableName specifies the table name for AttendanceSchedule
func (AttendanceSchedule) TableName() string {
	return "attendance_schedules"
}

// Validate validates the attendance schedule data
// Requirements: 3.1 - Schedule creation SHALL require name, start_time, end_time, and late_threshold
func (s *AttendanceSchedule) Validate() error {
	if strings.TrimSpace(s.Name) == "" {
		return errors.New("name is required")
	}
	if s.Name != "" && len(s.Name) > 100 {
		return errors.New("name must be at most 100 characters")
	}
	if s.StartTime == "" {
		return errors.New("start_time is required")
	}
	if s.EndTime == "" {
		return errors.New("end_time is required")
	}
	if s.LateThreshold < 0 {
		return errors.New("late_threshold must be non-negative")
	}
	if s.VeryLateThreshold != nil && *s.VeryLateThreshold < s.LateThreshold {
		return errors.New("very_late_threshold must be greater than or equal to late_threshold")
	}
	if s.SchoolID == 0 {
		return errors.New("school_id is required")
	}

	// Validate time format (accept both HH:MM and HH:MM:SS)
	if !isValidTimeFormat(s.StartTime) {
		return errors.New("start_time must be in HH:MM or HH:MM:SS format")
	}
	if !isValidTimeFormat(s.EndTime) {
		return errors.New("end_time must be in HH:MM or HH:MM:SS format")
	}

	// Validate days_of_week format
	if s.DaysOfWeek != "" {
		if err := s.ValidateDaysOfWeek(); err != nil {
			return err
		}
	}

	return nil
}

// isValidTimeFormat checks if time string is in HH:MM or HH:MM:SS format
func isValidTimeFormat(timeStr string) bool {
	if _, err := time.Parse("15:04", timeStr); err == nil {
		return true
	}
	if _, err := time.Parse("15:04:05", timeStr); err == nil {
		return true
	}
	return false
}

// ValidateDaysOfWeek validates the days_of_week format
// Days are represented as 0-6 (Sunday-Saturday) or 1-7 (Monday-Sunday)
func (s *AttendanceSchedule) ValidateDaysOfWeek() error {
	if s.DaysOfWeek == "" {
		return nil
	}

	days := strings.Split(s.DaysOfWeek, ",")
	for _, day := range days {
		d, err := strconv.Atoi(strings.TrimSpace(day))
		if err != nil {
			return errors.New("days_of_week must contain comma-separated numbers")
		}
		if d < 0 || d > 7 {
			return errors.New("days_of_week values must be between 0 and 7")
		}
	}
	return nil
}

// IsActiveOnDay checks if the schedule is active on a given weekday
// weekday: 0=Sunday, 1=Monday, ..., 6=Saturday (Go's time.Weekday)
func (s *AttendanceSchedule) IsActiveOnDay(weekday time.Weekday) bool {
	if s.DaysOfWeek == "" {
		return true // Active on all days if not specified
	}

	days := strings.Split(s.DaysOfWeek, ",")
	for _, day := range days {
		d, err := strconv.Atoi(strings.TrimSpace(day))
		if err != nil {
			continue
		}
		// Support both 0-6 (Sunday-Saturday) and 1-7 (Monday-Sunday) formats
		if d == int(weekday) || (d >= 1 && d <= 7 && d == int(weekday)+1) {
			return true
		}
	}
	return false
}

// IsTimeInRange checks if a given time falls within the schedule's time range
func (s *AttendanceSchedule) IsTimeInRange(t time.Time) bool {
	timeStr := t.Format("15:04:05")
	
	// Normalize start and end times to HH:MM:SS format
	startTime := s.StartTime
	if len(startTime) == 5 {
		startTime += ":00"
	}
	endTime := s.EndTime
	if len(endTime) == 5 {
		endTime += ":00"
	}

	return timeStr >= startTime && timeStr <= endTime
}

// GetLateStatus determines the attendance status based on check-in time
// Returns: on_time, late, or very_late
func (s *AttendanceSchedule) GetLateStatus(checkInTime time.Time) AttendanceStatus {
	startTime, err := time.Parse("15:04", s.StartTime)
	if err != nil {
		startTime, _ = time.Parse("15:04:05", s.StartTime)
	}

	// Create a time on the same date as checkInTime for comparison
	scheduleStart := time.Date(
		checkInTime.Year(), checkInTime.Month(), checkInTime.Day(),
		startTime.Hour(), startTime.Minute(), startTime.Second(), 0,
		checkInTime.Location(),
	)

	minutesLate := int(checkInTime.Sub(scheduleStart).Minutes())

	if minutesLate <= 0 {
		return AttendanceStatusOnTime
	}

	if s.VeryLateThreshold != nil && minutesLate >= *s.VeryLateThreshold {
		return AttendanceStatusVeryLate
	}

	if minutesLate >= s.LateThreshold {
		return AttendanceStatusLate
	}

	return AttendanceStatusOnTime
}

// Deactivate deactivates the schedule
// Requirements: 3.8 - Deactivating a schedule SHALL stop using it for new attendance records
func (s *AttendanceSchedule) Deactivate() {
	s.IsActive = false
}

// Activate activates the schedule
func (s *AttendanceSchedule) Activate() {
	s.IsActive = true
}

// SetAsDefault sets this schedule as the default for the school
func (s *AttendanceSchedule) SetAsDefault() {
	s.IsDefault = true
}

// UnsetDefault removes the default flag from this schedule
func (s *AttendanceSchedule) UnsetDefault() {
	s.IsDefault = false
}
