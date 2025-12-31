package settings

import (
	"time"
)

// ==================== Settings DTOs ====================

// UpdateSettingsRequest represents the request to update school settings
type UpdateSettingsRequest struct {
	// Attendance Settings
	AttendanceStartTime         *string `json:"attendance_start_time"`
	AttendanceEndTime           *string `json:"attendance_end_time"`
	AttendanceLateThreshold     *int    `json:"attendance_late_threshold"`
	AttendanceVeryLateThreshold *int    `json:"attendance_very_late_threshold"`

	// Notification Settings
	EnableAttendanceNotification *bool `json:"enable_attendance_notification"`
	EnableGradeNotification      *bool `json:"enable_grade_notification"`
	EnableBKNotification         *bool `json:"enable_bk_notification"`
	EnableHomeroomNotification   *bool `json:"enable_homeroom_notification"`

	// General Settings
	AcademicYear *string `json:"academic_year"`
	Semester     *int    `json:"semester"`
}

// UpdateAttendanceSettingsRequest represents the request to update attendance settings only
type UpdateAttendanceSettingsRequest struct {
	AttendanceStartTime         *string `json:"attendance_start_time"`
	AttendanceEndTime           *string `json:"attendance_end_time"`
	AttendanceLateThreshold     *int    `json:"attendance_late_threshold"`
	AttendanceVeryLateThreshold *int    `json:"attendance_very_late_threshold"`
}

// UpdateNotificationSettingsRequest represents the request to update notification settings only
type UpdateNotificationSettingsRequest struct {
	EnableAttendanceNotification *bool `json:"enable_attendance_notification"`
	EnableGradeNotification      *bool `json:"enable_grade_notification"`
	EnableBKNotification         *bool `json:"enable_bk_notification"`
	EnableHomeroomNotification   *bool `json:"enable_homeroom_notification"`
}

// UpdateAcademicSettingsRequest represents the request to update academic settings only
type UpdateAcademicSettingsRequest struct {
	AcademicYear *string `json:"academic_year"`
	Semester     *int    `json:"semester"`
}

// SettingsResponse represents school settings in responses
type SettingsResponse struct {
	ID       uint `json:"id"`
	SchoolID uint `json:"school_id"`

	// Attendance Settings
	AttendanceStartTime         string `json:"attendance_start_time"`
	AttendanceEndTime           string `json:"attendance_end_time"`
	AttendanceLateThreshold     int    `json:"attendance_late_threshold"`
	AttendanceVeryLateThreshold int    `json:"attendance_very_late_threshold"`

	// Notification Settings
	EnableAttendanceNotification bool `json:"enable_attendance_notification"`
	EnableGradeNotification      bool `json:"enable_grade_notification"`
	EnableBKNotification         bool `json:"enable_bk_notification"`
	EnableHomeroomNotification   bool `json:"enable_homeroom_notification"`

	// General Settings
	AcademicYear string `json:"academic_year"`
	Semester     int    `json:"semester"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AttendanceTimeWindowResponse represents the attendance time window response
type AttendanceTimeWindowResponse struct {
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	LateTime     string `json:"late_time"`
	VeryLateTime string `json:"very_late_time"`
}
