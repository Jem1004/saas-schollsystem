package schedule

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/school-management/backend/internal/domain/models"
)

const (
	// MaxSchedulesPerSchool is the maximum number of schedules allowed per school
	// Requirements: 3.9 - Maximum 10 schedules per school
	MaxSchedulesPerSchool = 10
)

var (
	ErrNameRequired          = errors.New("nama jadwal wajib diisi")
	ErrStartTimeRequired     = errors.New("waktu mulai wajib diisi")
	ErrEndTimeRequired       = errors.New("waktu akhir wajib diisi")
	ErrLateThresholdRequired = errors.New("batas keterlambatan wajib diisi")
	ErrInvalidStartTime      = errors.New("format waktu mulai tidak valid (gunakan HH:MM)")
	ErrInvalidEndTime        = errors.New("format waktu akhir tidak valid (gunakan HH:MM)")
	ErrEndTimeBeforeStart    = errors.New("waktu akhir harus setelah waktu mulai")
	ErrVeryLateThreshold     = errors.New("batas sangat terlambat harus lebih besar dari batas terlambat")
)

// Service defines the interface for schedule business logic
type Service interface {
	// CRUD operations
	CreateSchedule(ctx context.Context, schoolID uint, req CreateScheduleRequest) (*ScheduleResponse, error)
	GetAllSchedules(ctx context.Context, schoolID uint) (*ScheduleListResponse, error)
	GetScheduleByID(ctx context.Context, schoolID, id uint) (*ScheduleResponse, error)
	UpdateSchedule(ctx context.Context, schoolID, id uint, req UpdateScheduleRequest) (*ScheduleResponse, error)
	DeleteSchedule(ctx context.Context, schoolID, id uint) error

	// Active schedule operations
	// Requirements: 3.4 - Determine which schedule is currently active
	GetActiveSchedule(ctx context.Context, schoolID uint, timestamp time.Time) (*ActiveScheduleResponse, error)

	// Default schedule operations
	SetDefaultSchedule(ctx context.Context, schoolID, id uint) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new schedule service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateSchedule creates a new attendance schedule
// Requirements: 3.1 - Schedule creation with validation and limit check
// Requirements: 3.9 - Maximum 10 schedules per school
func (s *service) CreateSchedule(ctx context.Context, schoolID uint, req CreateScheduleRequest) (*ScheduleResponse, error) {
	// Validate required fields
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Check schedule limit
	count, err := s.repo.CountBySchool(ctx, schoolID)
	if err != nil {
		return nil, err
	}
	if count >= MaxSchedulesPerSchool {
		return nil, ErrScheduleLimitExceeded
	}

	// Create schedule model
	schedule := &models.AttendanceSchedule{
		SchoolID:          schoolID,
		Name:              strings.TrimSpace(req.Name),
		StartTime:         req.StartTime,
		EndTime:           req.EndTime,
		LateThreshold:     req.LateThreshold,
		VeryLateThreshold: req.VeryLateThreshold,
		DaysOfWeek:        req.DaysOfWeek,
		IsActive:          true,
		IsDefault:         false,
	}

	// Set default days of week if not provided
	if schedule.DaysOfWeek == "" {
		schedule.DaysOfWeek = "1,2,3,4,5" // Monday to Friday
	}

	// Override IsActive if provided
	if req.IsActive != nil {
		schedule.IsActive = *req.IsActive
	}

	// Validate the model
	if err := schedule.Validate(); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, schedule); err != nil {
		return nil, err
	}

	return toScheduleResponse(schedule), nil
}

// GetAllSchedules retrieves all schedules for a school
func (s *service) GetAllSchedules(ctx context.Context, schoolID uint) (*ScheduleListResponse, error) {
	schedules, err := s.repo.FindAll(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	responses := make([]ScheduleResponse, len(schedules))
	for i, schedule := range schedules {
		responses[i] = *toScheduleResponse(&schedule)
	}

	return &ScheduleListResponse{
		Schedules: responses,
		Total:     len(responses),
	}, nil
}

// GetScheduleByID retrieves a schedule by ID
func (s *service) GetScheduleByID(ctx context.Context, schoolID, id uint) (*ScheduleResponse, error) {
	schedule, err := s.repo.FindByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	return toScheduleResponse(schedule), nil
}


// UpdateSchedule updates an existing schedule
// Requirements: 3.7 - Updates SHALL not affect existing attendance records
func (s *service) UpdateSchedule(ctx context.Context, schoolID, id uint, req UpdateScheduleRequest) (*ScheduleResponse, error) {
	// Get existing schedule
	schedule, err := s.repo.FindByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if req.Name != nil {
		schedule.Name = strings.TrimSpace(*req.Name)
	}
	if req.StartTime != nil {
		schedule.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		schedule.EndTime = *req.EndTime
	}
	if req.LateThreshold != nil {
		schedule.LateThreshold = *req.LateThreshold
	}
	if req.VeryLateThreshold != nil {
		schedule.VeryLateThreshold = req.VeryLateThreshold
	}
	if req.DaysOfWeek != nil {
		schedule.DaysOfWeek = *req.DaysOfWeek
	}
	if req.IsActive != nil {
		schedule.IsActive = *req.IsActive
	}

	// Validate the updated model
	if err := schedule.Validate(); err != nil {
		return nil, err
	}

	// Validate time range
	if err := s.validateTimeRange(schedule.StartTime, schedule.EndTime); err != nil {
		return nil, err
	}

	// Validate very late threshold
	if schedule.VeryLateThreshold != nil && *schedule.VeryLateThreshold < schedule.LateThreshold {
		return nil, ErrVeryLateThreshold
	}

	// Update in database
	if err := s.repo.Update(ctx, schedule); err != nil {
		return nil, err
	}

	return toScheduleResponse(schedule), nil
}

// DeleteSchedule deletes a schedule
// Requirements: 3.8 - Deactivating a schedule SHALL stop using it for new attendance records
func (s *service) DeleteSchedule(ctx context.Context, schoolID, id uint) error {
	// Check if schedule exists
	_, err := s.repo.FindByID(ctx, schoolID, id)
	if err != nil {
		return err
	}

	// Check if schedule has attendance records
	hasRecords, err := s.repo.HasAttendanceRecords(ctx, id)
	if err != nil {
		return err
	}
	if hasRecords {
		return ErrScheduleInUse
	}

	return s.repo.Delete(ctx, schoolID, id)
}

// GetActiveSchedule finds the active schedule for a given time
// Requirements: 3.4 - Determine which schedule is currently active based on current time
// Requirements: 3.6 - IF no schedule is active, use default schedule or reject
func (s *service) GetActiveSchedule(ctx context.Context, schoolID uint, timestamp time.Time) (*ActiveScheduleResponse, error) {
	schedule, err := s.repo.FindActiveSchedule(ctx, schoolID, timestamp)
	if err != nil {
		return nil, err
	}

	if schedule == nil {
		// Try to get default schedule
		defaultSchedule, err := s.repo.FindDefaultSchedule(ctx, schoolID)
		if err != nil {
			return nil, err
		}

		if defaultSchedule != nil && defaultSchedule.IsActive {
			return &ActiveScheduleResponse{
				Schedule: toScheduleResponse(defaultSchedule),
				Message:  "Menggunakan jadwal default",
			}, nil
		}

		return &ActiveScheduleResponse{
			Schedule: nil,
			Message:  "Tidak ada jadwal aktif untuk waktu ini",
		}, nil
	}

	return &ActiveScheduleResponse{
		Schedule: toScheduleResponse(schedule),
	}, nil
}

// SetDefaultSchedule sets a schedule as the default for a school
func (s *service) SetDefaultSchedule(ctx context.Context, schoolID, id uint) error {
	// Verify schedule exists
	_, err := s.repo.FindByID(ctx, schoolID, id)
	if err != nil {
		return err
	}

	return s.repo.SetDefaultSchedule(ctx, schoolID, id)
}

// validateCreateRequest validates the create schedule request
func (s *service) validateCreateRequest(req CreateScheduleRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return ErrNameRequired
	}
	if req.StartTime == "" {
		return ErrStartTimeRequired
	}
	if req.EndTime == "" {
		return ErrEndTimeRequired
	}
	if req.LateThreshold < 0 {
		return ErrLateThresholdRequired
	}

	// Validate time formats (accept both HH:MM and HH:MM:SS)
	if !isValidTimeFormat(req.StartTime) {
		return ErrInvalidStartTime
	}
	if !isValidTimeFormat(req.EndTime) {
		return ErrInvalidEndTime
	}

	// Validate time range
	if err := s.validateTimeRange(req.StartTime, req.EndTime); err != nil {
		return err
	}

	// Validate very late threshold
	if req.VeryLateThreshold != nil && *req.VeryLateThreshold < req.LateThreshold {
		return ErrVeryLateThreshold
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

// parseTimeString parses time string in HH:MM or HH:MM:SS format
func parseTimeString(timeStr string) (time.Time, error) {
	if t, err := time.Parse("15:04", timeStr); err == nil {
		return t, nil
	}
	if t, err := time.Parse("15:04:05", timeStr); err == nil {
		return t, nil
	}
	return time.Time{}, errors.New("invalid time format")
}

// validateTimeRange validates that end time is after start time
func (s *service) validateTimeRange(startTime, endTime string) error {
	start, err := parseTimeString(startTime)
	if err != nil {
		return ErrInvalidStartTime
	}
	end, err := parseTimeString(endTime)
	if err != nil {
		return ErrInvalidEndTime
	}

	if !end.After(start) {
		return ErrEndTimeBeforeStart
	}

	return nil
}

// toScheduleResponse converts a model to response DTO
func toScheduleResponse(schedule *models.AttendanceSchedule) *ScheduleResponse {
	return &ScheduleResponse{
		ID:                schedule.ID,
		SchoolID:          schedule.SchoolID,
		Name:              schedule.Name,
		StartTime:         schedule.StartTime,
		EndTime:           schedule.EndTime,
		LateThreshold:     schedule.LateThreshold,
		VeryLateThreshold: schedule.VeryLateThreshold,
		DaysOfWeek:        schedule.DaysOfWeek,
		IsActive:          schedule.IsActive,
		IsDefault:         schedule.IsDefault,
		CreatedAt:         schedule.CreatedAt,
		UpdatedAt:         schedule.UpdatedAt,
	}
}
