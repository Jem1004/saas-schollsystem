package schedule

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrScheduleNotFound     = errors.New("jadwal absensi tidak ditemukan")
	ErrScheduleLimitExceeded = errors.New("batas maksimum jadwal (10) telah tercapai")
	ErrScheduleTimeOverlap  = errors.New("waktu jadwal bertumpang tindih dengan jadwal lain")
	ErrScheduleInUse        = errors.New("jadwal tidak dapat dihapus karena masih digunakan")
	ErrInvalidTimeRange     = errors.New("waktu akhir harus setelah waktu mulai")
)

// Repository defines the interface for schedule data operations
type Repository interface {
	// CRUD operations
	Create(ctx context.Context, schedule *models.AttendanceSchedule) error
	FindByID(ctx context.Context, schoolID, id uint) (*models.AttendanceSchedule, error)
	FindAll(ctx context.Context, schoolID uint) ([]models.AttendanceSchedule, error)
	Update(ctx context.Context, schedule *models.AttendanceSchedule) error
	Delete(ctx context.Context, schoolID, id uint) error

	// Query operations
	// Requirements: 3.4 - Find active schedule based on current time and day
	FindActiveSchedule(ctx context.Context, schoolID uint, timestamp time.Time) (*models.AttendanceSchedule, error)
	
	// Requirements: 3.9 - Maximum 10 schedules per school
	CountBySchool(ctx context.Context, schoolID uint) (int64, error)
	
	// Default schedule operations
	FindDefaultSchedule(ctx context.Context, schoolID uint) (*models.AttendanceSchedule, error)
	ClearDefaultSchedule(ctx context.Context, schoolID uint) error
	SetDefaultSchedule(ctx context.Context, schoolID, scheduleID uint) error
	
	// Check if schedule has attendance records
	HasAttendanceRecords(ctx context.Context, scheduleID uint) (bool, error)
	
	// Check for time overlap with existing schedules
	CheckTimeOverlap(ctx context.Context, schoolID uint, startTime, endTime, daysOfWeek string, excludeID *uint) (bool, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new schedule repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create creates a new attendance schedule
// Requirements: 3.1 - Schedule creation
func (r *repository) Create(ctx context.Context, schedule *models.AttendanceSchedule) error {
	// Use raw SQL to properly handle TIME type
	query := `
		INSERT INTO attendance_schedules 
		(school_id, name, start_time, end_time, late_threshold, very_late_threshold, days_of_week, is_active, is_default, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`
	now := time.Now()
	schedule.CreatedAt = now
	schedule.UpdatedAt = now
	
	err := r.db.WithContext(ctx).Raw(query,
		schedule.SchoolID,
		schedule.Name,
		schedule.StartTime,
		schedule.EndTime,
		schedule.LateThreshold,
		schedule.VeryLateThreshold,
		schedule.DaysOfWeek,
		schedule.IsActive,
		schedule.IsDefault,
		schedule.CreatedAt,
		schedule.UpdatedAt,
	).Scan(&schedule.ID).Error
	
	return err
}

// FindByID retrieves a schedule by ID for a specific school
func (r *repository) FindByID(ctx context.Context, schoolID, id uint) (*models.AttendanceSchedule, error) {
	var schedule models.AttendanceSchedule
	err := r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", id, schoolID).
		First(&schedule).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScheduleNotFound
		}
		return nil, err
	}

	return &schedule, nil
}

// FindAll retrieves all schedules for a school
func (r *repository) FindAll(ctx context.Context, schoolID uint) ([]models.AttendanceSchedule, error) {
	var schedules []models.AttendanceSchedule
	err := r.db.WithContext(ctx).
		Where("school_id = ?", schoolID).
		Order("is_default DESC, name ASC").
		Find(&schedules).Error

	return schedules, err
}

// Update updates an existing schedule
// Requirements: 3.7 - Updates SHALL not affect existing attendance records
func (r *repository) Update(ctx context.Context, schedule *models.AttendanceSchedule) error {
	// Use raw SQL to properly handle TIME type
	query := `
		UPDATE attendance_schedules 
		SET name = $1, 
			start_time = $2, 
			end_time = $3, 
			late_threshold = $4, 
			very_late_threshold = $5, 
			days_of_week = $6, 
			is_active = $7, 
			is_default = $8,
			updated_at = $9
		WHERE id = $10 AND school_id = $11
	`
	schedule.UpdatedAt = time.Now()
	
	result := r.db.WithContext(ctx).Exec(query,
		schedule.Name,
		schedule.StartTime,
		schedule.EndTime,
		schedule.LateThreshold,
		schedule.VeryLateThreshold,
		schedule.DaysOfWeek,
		schedule.IsActive,
		schedule.IsDefault,
		schedule.UpdatedAt,
		schedule.ID,
		schedule.SchoolID,
	)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrScheduleNotFound
	}
	return nil
}

// Delete deletes a schedule
func (r *repository) Delete(ctx context.Context, schoolID, id uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", id, schoolID).
		Delete(&models.AttendanceSchedule{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrScheduleNotFound
	}
	return nil
}


// FindActiveSchedule finds the active schedule for a given time and day
// Requirements: 3.4 - Determine which schedule is currently active based on current time
// Property 8: Active Schedule Selection
func (r *repository) FindActiveSchedule(ctx context.Context, schoolID uint, timestamp time.Time) (*models.AttendanceSchedule, error) {
	var schedules []models.AttendanceSchedule
	
	// Get current time in HH:MM:SS format
	currentTime := timestamp.Format("15:04:05")
	
	// Find all active schedules for this school
	err := r.db.WithContext(ctx).
		Where("school_id = ? AND is_active = ?", schoolID, true).
		Find(&schedules).Error

	if err != nil {
		return nil, err
	}

	// Find the schedule that matches the current time and day
	for _, schedule := range schedules {
		// Check if schedule is active on this day
		if !schedule.IsActiveOnDay(timestamp.Weekday()) {
			continue
		}

		// Normalize times for comparison
		startTime := schedule.StartTime
		if len(startTime) == 5 {
			startTime += ":00"
		}
		endTime := schedule.EndTime
		if len(endTime) == 5 {
			endTime += ":00"
		}

		// Check if current time is within schedule range
		if currentTime >= startTime && currentTime <= endTime {
			return &schedule, nil
		}
	}

	// No active schedule found for current time, try to find default
	defaultSchedule, err := r.FindDefaultSchedule(ctx, schoolID)
	if err == nil && defaultSchedule != nil {
		// Check if default schedule is active on this day
		if defaultSchedule.IsActive && defaultSchedule.IsActiveOnDay(timestamp.Weekday()) {
			return defaultSchedule, nil
		}
	}

	// Return nil if no schedule matches (will use default behavior)
	// Requirements: 3.6 - IF no schedule is active, use default schedule or reject
	return nil, nil
}

// CountBySchool counts the number of schedules for a school
// Requirements: 3.9 - Maximum 10 schedules per school
func (r *repository) CountBySchool(ctx context.Context, schoolID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.AttendanceSchedule{}).
		Where("school_id = ?", schoolID).
		Count(&count).Error

	return count, err
}

// FindDefaultSchedule finds the default schedule for a school
func (r *repository) FindDefaultSchedule(ctx context.Context, schoolID uint) (*models.AttendanceSchedule, error) {
	var schedule models.AttendanceSchedule
	err := r.db.WithContext(ctx).
		Where("school_id = ? AND is_default = ?", schoolID, true).
		First(&schedule).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No default schedule is not an error
		}
		return nil, err
	}

	return &schedule, nil
}

// ClearDefaultSchedule clears the default flag from all schedules for a school
func (r *repository) ClearDefaultSchedule(ctx context.Context, schoolID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.AttendanceSchedule{}).
		Where("school_id = ? AND is_default = ?", schoolID, true).
		Update("is_default", false).Error
}

// SetDefaultSchedule sets a schedule as the default for a school
func (r *repository) SetDefaultSchedule(ctx context.Context, schoolID, scheduleID uint) error {
	// First, clear any existing default
	if err := r.ClearDefaultSchedule(ctx, schoolID); err != nil {
		return err
	}

	// Then set the new default
	result := r.db.WithContext(ctx).
		Model(&models.AttendanceSchedule{}).
		Where("id = ? AND school_id = ?", scheduleID, schoolID).
		Update("is_default", true)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrScheduleNotFound
	}
	return nil
}

// HasAttendanceRecords checks if a schedule has any attendance records
func (r *repository) HasAttendanceRecords(ctx context.Context, scheduleID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Attendance{}).
		Where("schedule_id = ?", scheduleID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}


// CheckTimeOverlap checks if a schedule's time range overlaps with existing schedules
// This helps prevent ambiguous active schedule selection
func (r *repository) CheckTimeOverlap(ctx context.Context, schoolID uint, startTime, endTime, daysOfWeek string, excludeID *uint) (bool, error) {
	// Get all active schedules for this school
	var schedules []models.AttendanceSchedule
	query := r.db.WithContext(ctx).
		Where("school_id = ? AND is_active = ?", schoolID, true)
	
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	
	if err := query.Find(&schedules).Error; err != nil {
		return false, err
	}

	// Normalize input times
	if len(startTime) == 5 {
		startTime += ":00"
	}
	if len(endTime) == 5 {
		endTime += ":00"
	}

	// Parse days of week for the new schedule
	newDays := parseDaysOfWeek(daysOfWeek)

	// Check each existing schedule for overlap
	for _, schedule := range schedules {
		// Check if days overlap
		existingDays := parseDaysOfWeek(schedule.DaysOfWeek)
		if !daysOverlap(newDays, existingDays) {
			continue
		}

		// Normalize existing schedule times
		existingStart := schedule.StartTime
		if len(existingStart) == 5 {
			existingStart += ":00"
		}
		existingEnd := schedule.EndTime
		if len(existingEnd) == 5 {
			existingEnd += ":00"
		}

		// Check time overlap: (start1 < end2) AND (end1 > start2)
		if startTime < existingEnd && endTime > existingStart {
			return true, nil
		}
	}

	return false, nil
}

// parseDaysOfWeek parses days of week string (e.g., "1,2,3,4,5") into a map
func parseDaysOfWeek(daysStr string) map[int]bool {
	days := make(map[int]bool)
	if daysStr == "" {
		// Default to weekdays
		for i := 1; i <= 5; i++ {
			days[i] = true
		}
		return days
	}
	
	for _, d := range daysStr {
		if d >= '0' && d <= '6' {
			days[int(d-'0')] = true
		}
	}
	return days
}

// daysOverlap checks if two sets of days have any overlap
func daysOverlap(days1, days2 map[int]bool) bool {
	for day := range days1 {
		if days2[day] {
			return true
		}
	}
	return false
}
