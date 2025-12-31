package settings

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrInvalidTimeFormat           = errors.New("time must be in HH:MM format")
	ErrInvalidLateThreshold        = errors.New("late threshold must be non-negative")
	ErrInvalidVeryLateThreshold    = errors.New("very late threshold must be greater than or equal to late threshold")
	ErrInvalidSemester             = errors.New("semester must be 1 or 2")
	ErrSchoolIDRequired            = errors.New("school_id is required")
)

// Service defines the interface for SchoolSettings business logic
type Service interface {
	// Settings operations
	GetSchoolSettings(ctx context.Context, schoolID uint) (*SettingsResponse, error)
	UpdateSchoolSettings(ctx context.Context, schoolID uint, req UpdateSettingsRequest) (*SettingsResponse, error)
	ResetToDefaults(ctx context.Context, schoolID uint) (*SettingsResponse, error)

	// Partial update operations
	UpdateAttendanceSettings(ctx context.Context, schoolID uint, req UpdateAttendanceSettingsRequest) (*SettingsResponse, error)
	UpdateNotificationSettings(ctx context.Context, schoolID uint, req UpdateNotificationSettingsRequest) (*SettingsResponse, error)
	UpdateAcademicSettings(ctx context.Context, schoolID uint, req UpdateAcademicSettingsRequest) (*SettingsResponse, error)

	// Utility operations
	GetAttendanceTimeWindow(ctx context.Context, schoolID uint, date time.Time) (*AttendanceTimeWindowResponse, error)
	ShouldSendNotification(ctx context.Context, schoolID uint, notificationType models.NotificationType) (bool, error)

	// Ensure settings exist (creates defaults if not)
	EnsureSettingsExist(ctx context.Context, schoolID uint) (*models.SchoolSettings, error)
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new Settings service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// GetSchoolSettings retrieves settings for a school, creating defaults if not exist
func (s *service) GetSchoolSettings(ctx context.Context, schoolID uint) (*SettingsResponse, error) {
	if schoolID == 0 {
		return nil, ErrSchoolIDRequired
	}

	settings, err := s.EnsureSettingsExist(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	return toSettingsResponse(settings), nil
}

// UpdateSchoolSettings updates all school settings
func (s *service) UpdateSchoolSettings(ctx context.Context, schoolID uint, req UpdateSettingsRequest) (*SettingsResponse, error) {
	if schoolID == 0 {
		return nil, ErrSchoolIDRequired
	}

	settings, err := s.EnsureSettingsExist(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if req.AttendanceStartTime != nil {
		if err := validateTimeFormat(*req.AttendanceStartTime); err != nil {
			return nil, err
		}
		settings.AttendanceStartTime = *req.AttendanceStartTime
	}
	if req.AttendanceEndTime != nil {
		if err := validateTimeFormat(*req.AttendanceEndTime); err != nil {
			return nil, err
		}
		settings.AttendanceEndTime = *req.AttendanceEndTime
	}
	if req.AttendanceLateThreshold != nil {
		if *req.AttendanceLateThreshold < 0 {
			return nil, ErrInvalidLateThreshold
		}
		settings.AttendanceLateThreshold = *req.AttendanceLateThreshold
	}
	if req.AttendanceVeryLateThreshold != nil {
		if *req.AttendanceVeryLateThreshold < settings.AttendanceLateThreshold {
			return nil, ErrInvalidVeryLateThreshold
		}
		settings.AttendanceVeryLateThreshold = *req.AttendanceVeryLateThreshold
	}
	if req.EnableAttendanceNotification != nil {
		settings.EnableAttendanceNotification = *req.EnableAttendanceNotification
	}
	if req.EnableGradeNotification != nil {
		settings.EnableGradeNotification = *req.EnableGradeNotification
	}
	if req.EnableBKNotification != nil {
		settings.EnableBKNotification = *req.EnableBKNotification
	}
	if req.EnableHomeroomNotification != nil {
		settings.EnableHomeroomNotification = *req.EnableHomeroomNotification
	}
	if req.AcademicYear != nil {
		settings.AcademicYear = *req.AcademicYear
	}
	if req.Semester != nil {
		if *req.Semester != 1 && *req.Semester != 2 {
			return nil, ErrInvalidSemester
		}
		settings.Semester = *req.Semester
	}

	// Validate the complete settings
	if err := settings.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, settings); err != nil {
		return nil, err
	}

	return toSettingsResponse(settings), nil
}

// ResetToDefaults resets school settings to default values
func (s *service) ResetToDefaults(ctx context.Context, schoolID uint) (*SettingsResponse, error) {
	if schoolID == 0 {
		return nil, ErrSchoolIDRequired
	}

	// Check if school exists
	_, err := s.repo.FindSchoolByID(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	// Get existing settings or create new
	settings, err := s.repo.FindBySchoolID(ctx, schoolID)
	if err != nil && !errors.Is(err, ErrSettingsNotFound) {
		return nil, err
	}

	// Create default settings
	defaultSettings := models.DefaultSchoolSettings(schoolID)

	if settings != nil {
		// Update existing settings with defaults
		defaultSettings.ID = settings.ID
		defaultSettings.CreatedAt = settings.CreatedAt
		if err := s.repo.Update(ctx, defaultSettings); err != nil {
			return nil, err
		}
	} else {
		// Create new settings with defaults
		if err := s.repo.Create(ctx, defaultSettings); err != nil {
			return nil, err
		}
	}

	return toSettingsResponse(defaultSettings), nil
}

// UpdateAttendanceSettings updates only attendance-related settings
func (s *service) UpdateAttendanceSettings(ctx context.Context, schoolID uint, req UpdateAttendanceSettingsRequest) (*SettingsResponse, error) {
	return s.UpdateSchoolSettings(ctx, schoolID, UpdateSettingsRequest{
		AttendanceStartTime:         req.AttendanceStartTime,
		AttendanceEndTime:           req.AttendanceEndTime,
		AttendanceLateThreshold:     req.AttendanceLateThreshold,
		AttendanceVeryLateThreshold: req.AttendanceVeryLateThreshold,
	})
}

// UpdateNotificationSettings updates only notification-related settings
func (s *service) UpdateNotificationSettings(ctx context.Context, schoolID uint, req UpdateNotificationSettingsRequest) (*SettingsResponse, error) {
	return s.UpdateSchoolSettings(ctx, schoolID, UpdateSettingsRequest{
		EnableAttendanceNotification: req.EnableAttendanceNotification,
		EnableGradeNotification:      req.EnableGradeNotification,
		EnableBKNotification:         req.EnableBKNotification,
		EnableHomeroomNotification:   req.EnableHomeroomNotification,
	})
}

// UpdateAcademicSettings updates only academic-related settings
func (s *service) UpdateAcademicSettings(ctx context.Context, schoolID uint, req UpdateAcademicSettingsRequest) (*SettingsResponse, error) {
	return s.UpdateSchoolSettings(ctx, schoolID, UpdateSettingsRequest{
		AcademicYear: req.AcademicYear,
		Semester:     req.Semester,
	})
}


// GetAttendanceTimeWindow calculates the attendance time window for a specific date
// Property 17: School Settings Policy Enforcement - Attendance status SHALL be determined based on school's configured time thresholds
func (s *service) GetAttendanceTimeWindow(ctx context.Context, schoolID uint, date time.Time) (*AttendanceTimeWindowResponse, error) {
	if schoolID == 0 {
		return nil, ErrSchoolIDRequired
	}

	settings, err := s.EnsureSettingsExist(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	// Parse start time
	startTime, err := parseTimeOnDate(settings.AttendanceStartTime, date)
	if err != nil {
		return nil, fmt.Errorf("invalid start time: %w", err)
	}

	// Parse end time
	endTime, err := parseTimeOnDate(settings.AttendanceEndTime, date)
	if err != nil {
		return nil, fmt.Errorf("invalid end time: %w", err)
	}

	// Calculate late and very late times
	lateTime := startTime.Add(time.Duration(settings.AttendanceLateThreshold) * time.Minute)
	veryLateTime := startTime.Add(time.Duration(settings.AttendanceVeryLateThreshold) * time.Minute)

	return &AttendanceTimeWindowResponse{
		StartTime:    startTime.Format("15:04"),
		EndTime:      endTime.Format("15:04"),
		LateTime:     lateTime.Format("15:04"),
		VeryLateTime: veryLateTime.Format("15:04"),
	}, nil
}

// ShouldSendNotification checks if a notification should be sent based on settings
// Property 17: Notifications SHALL only be sent if the corresponding notification setting is enabled
func (s *service) ShouldSendNotification(ctx context.Context, schoolID uint, notificationType models.NotificationType) (bool, error) {
	if schoolID == 0 {
		return false, ErrSchoolIDRequired
	}

	settings, err := s.EnsureSettingsExist(ctx, schoolID)
	if err != nil {
		return false, err
	}

	return settings.ShouldSendNotification(notificationType), nil
}

// EnsureSettingsExist ensures settings exist for a school, creating defaults if not
// Property 17: Default settings SHALL be applied when a new school is created
func (s *service) EnsureSettingsExist(ctx context.Context, schoolID uint) (*models.SchoolSettings, error) {
	settings, err := s.repo.FindBySchoolID(ctx, schoolID)
	if err == nil {
		return settings, nil
	}

	if !errors.Is(err, ErrSettingsNotFound) {
		return nil, err
	}

	// Settings not found, check if school exists
	_, err = s.repo.FindSchoolByID(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	// Create default settings
	defaultSettings := models.DefaultSchoolSettings(schoolID)
	if err := s.repo.Create(ctx, defaultSettings); err != nil {
		return nil, err
	}

	return defaultSettings, nil
}

// ==================== Helper Functions ====================

// validateTimeFormat validates that a string is in HH:MM format
func validateTimeFormat(timeStr string) error {
	timeRegex := regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]$`)
	if !timeRegex.MatchString(timeStr) {
		return ErrInvalidTimeFormat
	}
	return nil
}

// parseTimeOnDate parses a time string (HH:MM) and combines it with a date
func parseTimeOnDate(timeStr string, date time.Time) (time.Time, error) {
	var hour, minute int
	_, err := fmt.Sscanf(timeStr, "%d:%d", &hour, &minute)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(
		date.Year(), date.Month(), date.Day(),
		hour, minute, 0, 0,
		date.Location(),
	), nil
}

// toSettingsResponse converts a SchoolSettings model to a response DTO
func toSettingsResponse(s *models.SchoolSettings) *SettingsResponse {
	return &SettingsResponse{
		ID:                           s.ID,
		SchoolID:                     s.SchoolID,
		AttendanceStartTime:          s.AttendanceStartTime,
		AttendanceEndTime:            s.AttendanceEndTime,
		AttendanceLateThreshold:      s.AttendanceLateThreshold,
		AttendanceVeryLateThreshold:  s.AttendanceVeryLateThreshold,
		EnableAttendanceNotification: s.EnableAttendanceNotification,
		EnableGradeNotification:      s.EnableGradeNotification,
		EnableBKNotification:         s.EnableBKNotification,
		EnableHomeroomNotification:   s.EnableHomeroomNotification,
		AcademicYear:                 s.AcademicYear,
		Semester:                     s.Semester,
		CreatedAt:                    s.CreatedAt,
		UpdatedAt:                    s.UpdatedAt,
	}
}
