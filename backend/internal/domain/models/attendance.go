package models

import (
	"errors"
	"time"
)

// AttendanceMethod represents how attendance was recorded
type AttendanceMethod string

const (
	AttendanceMethodRFID   AttendanceMethod = "rfid"
	AttendanceMethodManual AttendanceMethod = "manual"
)

// IsValid checks if the attendance method is valid
func (m AttendanceMethod) IsValid() bool {
	switch m {
	case AttendanceMethodRFID, AttendanceMethodManual:
		return true
	}
	return false
}

// AttendanceStatus represents the attendance status
type AttendanceStatus string

const (
	AttendanceStatusOnTime   AttendanceStatus = "on_time"
	AttendanceStatusLate     AttendanceStatus = "late"
	AttendanceStatusVeryLate AttendanceStatus = "very_late"
	AttendanceStatusAbsent   AttendanceStatus = "absent"
)

// IsValid checks if the attendance status is valid
func (s AttendanceStatus) IsValid() bool {
	switch s {
	case AttendanceStatusOnTime, AttendanceStatusLate, AttendanceStatusVeryLate, AttendanceStatusAbsent:
		return true
	}
	return false
}

// Attendance represents daily attendance record
type Attendance struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	StudentID    uint             `gorm:"index;not null" json:"student_id"`
	ScheduleID   *uint            `gorm:"index" json:"schedule_id"`
	Date         time.Time        `gorm:"type:date;index;not null" json:"date"`
	CheckInTime  *time.Time       `json:"check_in_time"`
	CheckOutTime *time.Time       `json:"check_out_time"`
	Status       AttendanceStatus `gorm:"type:varchar(20)" json:"status"`
	Method       AttendanceMethod `gorm:"type:varchar(10);not null" json:"method"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`

	// Relations
	Student  Student             `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Schedule *AttendanceSchedule `gorm:"foreignKey:ScheduleID" json:"schedule,omitempty"`
}

// TableName specifies the table name for Attendance
func (Attendance) TableName() string {
	return "attendances"
}

// Validate validates the attendance data
// Requirements: 5.2 - Attendance SHALL record check-in or check-out based on existing records
func (a *Attendance) Validate() error {
	if a.StudentID == 0 {
		return errors.New("student_id is required")
	}
	if a.Date.IsZero() {
		return errors.New("date is required")
	}
	if !a.Method.IsValid() {
		return errors.New("method must be one of: rfid, manual")
	}
	return nil
}

// HasCheckedIn checks if the student has checked in
func (a *Attendance) HasCheckedIn() bool {
	return a.CheckInTime != nil
}

// HasCheckedOut checks if the student has checked out
func (a *Attendance) HasCheckedOut() bool {
	return a.CheckOutTime != nil
}

// SetCheckIn sets the check-in time
func (a *Attendance) SetCheckIn(t time.Time) {
	a.CheckInTime = &t
}

// SetCheckOut sets the check-out time
// Returns error if check-out time is before check-in time
func (a *Attendance) SetCheckOut(t time.Time) error {
	if a.CheckInTime != nil && t.Before(*a.CheckInTime) {
		return errors.New("check_out_time cannot be before check_in_time")
	}
	a.CheckOutTime = &t
	return nil
}
