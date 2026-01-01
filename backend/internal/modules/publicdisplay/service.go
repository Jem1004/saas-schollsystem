package publicdisplay

import (
	"context"
	"time"

	"github.com/school-management/backend/internal/modules/displaytoken"
)

// Service defines the interface for public display operations
type Service interface {
	// GetPublicDisplayData retrieves all data for public display using display token
	// Requirements: 5.3, 5.4, 5.5, 5.6 - Public display data access via token
	GetPublicDisplayData(ctx context.Context, token string) (*PublicDisplayData, error)

	// ValidateAndUpdateAccess validates token and updates last_accessed timestamp
	// Requirements: 5.10, 5.13, 6.7 - Token validation and access tracking
	ValidateAndUpdateAccess(ctx context.Context, token string) (*displaytoken.DisplayTokenValidation, error)
}

// service implements the Service interface
type service struct {
	repo                Repository
	displayTokenService displaytoken.Service
}

// NewService creates a new public display service
func NewService(repo Repository, displayTokenService displaytoken.Service) Service {
	return &service{
		repo:                repo,
		displayTokenService: displayTokenService,
	}
}

// GetPublicDisplayData retrieves all data for public display using display token
// Requirements: 5.3 - Accessing public display URL with valid token SHALL show attendance data without login
// Requirements: 5.4 - Show live feed of recent attendance (last 10 records)
// Requirements: 5.5 - Show real-time statistics (present, late, absent count)
// Requirements: 5.6 - Show leaderboard of top 10 earliest arrivals for the day
// Requirements: 5.7 - Show current date and time
// Requirements: 5.8 - Show school name
func (s *service) GetPublicDisplayData(ctx context.Context, token string) (*PublicDisplayData, error) {
	// Validate token and get school ID
	validation, err := s.ValidateAndUpdateAccess(ctx, token)
	if err != nil {
		return nil, err
	}

	if !validation.Valid {
		return nil, &PublicDisplayAccessError{
			Code:    "TOKEN_INVALID",
			Message: validation.Error,
		}
	}

	schoolID := validation.SchoolID

	// Get school information
	school, err := s.repo.GetSchoolByID(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	// Get live feed (last 10 records)
	// Requirements: 5.4 - Show live feed of recent attendance (last 10 records)
	liveFeed, err := s.repo.GetPublicLiveFeed(ctx, schoolID, 10)
	if err != nil {
		return nil, err
	}

	// Get attendance stats
	// Requirements: 5.5 - Show real-time statistics
	stats, err := s.repo.GetPublicStats(ctx, schoolID, now)
	if err != nil {
		return nil, err
	}

	// Get leaderboard (top 10 earliest arrivals)
	// Requirements: 5.6 - Show leaderboard of top 10 earliest arrivals
	leaderboard, err := s.repo.GetPublicLeaderboard(ctx, schoolID, 10, now)
	if err != nil {
		return nil, err
	}

	// Format date in Indonesian
	dateFormatted := formatIndonesianDate(now)

	return &PublicDisplayData{
		SchoolName:  school.Name,
		CurrentTime: now,
		Date:        dateFormatted,
		Stats: PublicAttendanceStats{
			TotalStudents: stats.TotalStudents,
			Present:       stats.Present,
			Late:          stats.Late,
			VeryLate:      stats.VeryLate,
			Absent:        stats.Absent,
			Percentage:    stats.Percentage,
		},
		LiveFeed:    liveFeed,
		Leaderboard: leaderboard,
	}, nil
}

// ValidateAndUpdateAccess validates token and updates last_accessed timestamp
// Requirements: 5.10, 5.13 - Token validation for access control
// Requirements: 6.7 - Track last accessed timestamp for each display token
func (s *service) ValidateAndUpdateAccess(ctx context.Context, token string) (*displaytoken.DisplayTokenValidation, error) {
	// Validate token
	validation, err := s.displayTokenService.ValidateToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// If valid, update last accessed timestamp
	if validation.Valid {
		// Update last accessed (ignore error as it's not critical)
		_ = s.displayTokenService.UpdateLastAccessed(ctx, token)
	}

	return validation, nil
}

// formatIndonesianDate formats a date in Indonesian format
// Example: "Senin, 1 Januari 2024"
func formatIndonesianDate(t time.Time) string {
	days := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	months := []string{
		"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}

	dayName := days[t.Weekday()]
	monthName := months[t.Month()]

	return dayName + ", " + t.Format("2") + " " + monthName + " " + t.Format("2006")
}

// PublicDisplayAccessError represents an access error for public display
type PublicDisplayAccessError struct {
	Code    string
	Message string
}

func (e *PublicDisplayAccessError) Error() string {
	return e.Message
}
