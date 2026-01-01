package schedule

import (
	"time"
)

// ==================== Request DTOs ====================

// CreateScheduleRequest represents the request to create a new attendance schedule
// Requirements: 3.1 - Schedule creation SHALL require name, start_time, end_time, and late_threshold
type CreateScheduleRequest struct {
	Name              string `json:"name" validate:"required,max=100"`
	StartTime         string `json:"start_time" validate:"required"`         // Format: HH:MM
	EndTime           string `json:"end_time" validate:"required"`           // Format: HH:MM
	LateThreshold     int    `json:"late_threshold" validate:"required,min=0"` // minutes after start_time
	VeryLateThreshold *int   `json:"very_late_threshold,omitempty"`          // optional, minutes after start_time
	DaysOfWeek        string `json:"days_of_week,omitempty"`                 // e.g., "1,2,3,4,5" (Mon-Fri)
	IsActive          *bool  `json:"is_active,omitempty"`                    // defaults to true
}

// UpdateScheduleRequest represents the request to update an attendance schedule
// Requirements: 3.7 - Schedule updates SHALL not affect existing attendance records
type UpdateScheduleRequest struct {
	Name              *string `json:"name,omitempty" validate:"omitempty,max=100"`
	StartTime         *string `json:"start_time,omitempty"`
	EndTime           *string `json:"end_time,omitempty"`
	LateThreshold     *int    `json:"late_threshold,omitempty" validate:"omitempty,min=0"`
	VeryLateThreshold *int    `json:"very_late_threshold,omitempty"`
	DaysOfWeek        *string `json:"days_of_week,omitempty"`
	IsActive          *bool   `json:"is_active,omitempty"`
}

// ==================== Response DTOs ====================

// ScheduleResponse represents an attendance schedule in responses
type ScheduleResponse struct {
	ID                uint      `json:"id"`
	SchoolID          uint      `json:"school_id"`
	Name              string    `json:"name"`
	StartTime         string    `json:"start_time"`
	EndTime           string    `json:"end_time"`
	LateThreshold     int       `json:"late_threshold"`
	VeryLateThreshold *int      `json:"very_late_threshold,omitempty"`
	DaysOfWeek        string    `json:"days_of_week"`
	IsActive          bool      `json:"is_active"`
	IsDefault         bool      `json:"is_default"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// ScheduleListResponse represents a list of attendance schedules
type ScheduleListResponse struct {
	Schedules []ScheduleResponse `json:"schedules"`
	Total     int                `json:"total"`
}

// ActiveScheduleResponse represents the currently active schedule
type ActiveScheduleResponse struct {
	Schedule *ScheduleResponse `json:"schedule,omitempty"`
	Message  string            `json:"message,omitempty"`
}
