// Package models contains all domain models for the School Management SaaS system.
// These models represent the core entities and their relationships.
package models

import "errors"

// This file serves as the package documentation.
// All models are defined in their respective files:
//
// Core Models:
//   - school.go: School (tenant) model
//   - user.go: User model with roles
//   - class.go: Class model
//   - student.go: Student model
//   - parent.go: Parent model
//
// Attendance:
//   - attendance.go: Attendance record model
//   - attendance_schedule.go: Attendance schedule model for multi-schedule support
//
// BK (Counseling) Models:
//   - violation.go: Violation record model
//   - achievement.go: Achievement record model
//   - permit.go: Exit permit model
//   - counseling_note.go: Counseling note model
//
// Academic Models:
//   - grade.go: Grade entry model
//   - homeroom_note.go: Homeroom teacher note model
//
// Device & Notification:
//   - device.go: RFID device (ESP32) model
//   - notification.go: Notification and FCM token models
//
// Display:
//   - display_token.go: Display token model for public display access
//
// Settings:
//   - school_settings.go: School-specific settings model
//
// Event Outbox:
//   - outbox.go: Outbox event model for reliable event publishing

// Common validation errors
var (
	ErrRequiredFieldMissing = errors.New("required field is missing")
	ErrInvalidFieldValue    = errors.New("invalid field value")
	ErrDuplicateEntry       = errors.New("duplicate entry")
)

// AllModels returns all models for GORM auto-migration
// This ensures all models are registered in a single place
func AllModels() []interface{} {
	return []interface{}{
		// Core models
		&School{},
		&User{},
		&Class{},
		&Student{},
		&Parent{},
		&StudentParent{},
		&ClassCounselor{},

		// Attendance
		&Attendance{},
		&AttendanceSchedule{},

		// BK models
		&Violation{},
		&Achievement{},
		&Permit{},
		&CounselingNote{},

		// Academic models
		&Grade{},
		&HomeroomNote{},

		// Device & Notification
		&Device{},
		&Notification{},
		&FCMToken{},

		// Settings
		&SchoolSettings{},

		// Display
		&DisplayToken{},

		// Outbox
		&OutboxEvent{},
	}
}

// Pagination represents pagination parameters
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int64 `json:"total"`
}

// DefaultPagination returns default pagination settings
func DefaultPagination() Pagination {
	return Pagination{
		Page:     1,
		PageSize: 20,
	}
}

// Offset calculates the offset for database queries
func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit returns the page size for database queries
func (p Pagination) Limit() int {
	if p.PageSize <= 0 {
		return 20
	}
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}
