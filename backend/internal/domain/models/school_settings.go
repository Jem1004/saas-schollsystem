package models

import (
	"errors"
	"regexp"
	"time"
)

// SchoolSettings represents configurable settings per school
type SchoolSettings struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	SchoolID uint `gorm:"uniqueIndex;not null" json:"school_id"`

	// Attendance Settings
	AttendanceStartTime         string `gorm:"type:varchar(5);default:'07:00'" json:"attendance_start_time"` // HH:MM format
	AttendanceEndTime           string `gorm:"type:varchar(5);default:'07:30'" json:"attendance_end_time"`   // Late after this
	AttendanceLateThreshold     int    `gorm:"default:30" json:"attendance_late_threshold"`                  // Minutes after start to be considered late
	AttendanceVeryLateThreshold int    `gorm:"default:60" json:"attendance_very_late_threshold"`             // Minutes after start to be considered very late

	// Notification Settings
	EnableAttendanceNotification bool `gorm:"default:true" json:"enable_attendance_notification"`
	EnableGradeNotification      bool `gorm:"default:true" json:"enable_grade_notification"`
	EnableBKNotification         bool `gorm:"default:true" json:"enable_bk_notification"`
	EnableHomeroomNotification   bool `gorm:"default:true" json:"enable_homeroom_notification"`

	// General Settings
	AcademicYear string `gorm:"type:varchar(10)" json:"academic_year"` // e.g., "2024/2025"
	Semester     int    `gorm:"default:1" json:"semester"`             // 1 or 2

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	School School `gorm:"foreignKey:SchoolID" json:"school,omitempty"`
}

// TableName specifies the table name for SchoolSettings
func (SchoolSettings) TableName() string {
	return "school_settings"
}

// Validate validates the school settings data
func (s *SchoolSettings) Validate() error {
	if s.SchoolID == 0 {
		return errors.New("school_id is required")
	}

	// Validate time format (HH:MM)
	timeRegex := regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]$`)
	if !timeRegex.MatchString(s.AttendanceStartTime) {
		return errors.New("attendance_start_time must be in HH:MM format")
	}
	if !timeRegex.MatchString(s.AttendanceEndTime) {
		return errors.New("attendance_end_time must be in HH:MM format")
	}

	// Validate thresholds
	if s.AttendanceLateThreshold < 0 {
		return errors.New("attendance_late_threshold must be non-negative")
	}
	if s.AttendanceVeryLateThreshold < s.AttendanceLateThreshold {
		return errors.New("attendance_very_late_threshold must be greater than or equal to attendance_late_threshold")
	}

	// Validate semester
	if s.Semester != 1 && s.Semester != 2 {
		return errors.New("semester must be 1 or 2")
	}

	return nil
}

// AttendanceTimeWindow represents the valid time window for attendance
type AttendanceTimeWindow struct {
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	LateTime     time.Time `json:"late_time"`
	VeryLateTime time.Time `json:"very_late_time"`
}

// DefaultSchoolSettings returns default settings for a new school
func DefaultSchoolSettings(schoolID uint) *SchoolSettings {
	return &SchoolSettings{
		SchoolID:                     schoolID,
		AttendanceStartTime:          "07:00",
		AttendanceEndTime:            "16:00",  // Extended to 4 PM for check-out
		AttendanceLateThreshold:      30,
		AttendanceVeryLateThreshold:  60,
		EnableAttendanceNotification: true,
		EnableGradeNotification:      true,
		EnableBKNotification:         true,
		EnableHomeroomNotification:   true,
		Semester:                     1,
	}
}

// ShouldSendNotification checks if a notification should be sent based on settings
func (s *SchoolSettings) ShouldSendNotification(notificationType NotificationType) bool {
	switch notificationType {
	case NotificationTypeAttendanceIn, NotificationTypeAttendanceOut:
		return s.EnableAttendanceNotification
	case NotificationTypeGrade:
		return s.EnableGradeNotification
	case NotificationTypeViolation, NotificationTypeAchievement, NotificationTypePermit, NotificationTypeCounseling:
		return s.EnableBKNotification
	case NotificationTypeHomeroomNote:
		return s.EnableHomeroomNotification
	default:
		return true
	}
}
