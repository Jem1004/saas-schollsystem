package student

import (
	"context"
	"time"

	"github.com/school-management/backend/internal/domain/models"
)

// Service defines the interface for student business logic
type Service interface {
	// Profile operations
	GetProfile(ctx context.Context, userID uint) (*StudentProfileResponse, error)
	GetDashboard(ctx context.Context, userID uint) (*DashboardResponse, error)
	GetSummary(ctx context.Context, userID uint) (*StudentSummaryResponse, error)

	// Attendance operations
	GetAttendance(ctx context.Context, userID uint, filter AttendanceFilter) (*AttendanceListResponse, error)
	GetAttendanceSummary(ctx context.Context, userID uint, startDate, endDate string) (*AttendanceSummaryResponse, error)

	// Grade operations
	GetGrades(ctx context.Context, userID uint, filter GradeFilter) (*GradeListResponse, error)
	GetGradeSummary(ctx context.Context, userID uint) (*GradeSummaryResponse, error)

	// BK operations
	GetBKInfo(ctx context.Context, userID uint) (*BKInfoResponse, error)
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new student service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// GetProfile retrieves the student's profile
// Requirements: 16.2 - Student views personal information
func (s *service) GetProfile(ctx context.Context, userID uint) (*StudentProfileResponse, error) {
	student, err := s.repo.FindStudentByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return toStudentProfileResponse(student), nil
}

// GetDashboard retrieves the student's dashboard
// Requirements: 16.2, 16.3, 16.4, 16.5 - Student views profile, attendance, grades, BK info
func (s *service) GetDashboard(ctx context.Context, userID uint) (*DashboardResponse, error) {
	student, err := s.repo.FindStudentByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get attendance summary (last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	attendanceSummary, _ := s.repo.GetAttendanceSummary(ctx, student.ID, thirtyDaysAgo, time.Now())
	if attendanceSummary == nil {
		attendanceSummary = &AttendanceSummaryResponse{}
	}

	// Get grade summary
	gradeSummary, _ := s.repo.GetGradeSummary(ctx, student.ID)
	if gradeSummary == nil {
		gradeSummary = &GradeSummaryResponse{}
	}

	// Get BK summary
	totalPoints, _ := s.repo.GetStudentAchievementPoints(ctx, student.ID)
	violations, achievements, _ := s.repo.GetStudentBKCounts(ctx, student.ID)

	// Get recent attendance (5)
	recentAttendances, _, _ := s.repo.GetStudentAttendance(ctx, student.ID, time.Time{}, time.Time{}, 1, 5)
	recentAttendanceResponses := make([]AttendanceResponse, len(recentAttendances))
	for i, a := range recentAttendances {
		recentAttendanceResponses[i] = toAttendanceResponse(&a)
	}

	// Get recent grades (5)
	recentGrades, _, _ := s.repo.GetStudentGrades(ctx, student.ID, 1, 5)
	recentGradeResponses := make([]GradeResponse, len(recentGrades))
	for i, g := range recentGrades {
		recentGradeResponses[i] = toGradeResponse(&g)
	}

	// Get recent achievements (5)
	recentAchievements, _ := s.repo.GetStudentAchievements(ctx, student.ID, 5)
	recentAchievementResponses := make([]AchievementResponse, len(recentAchievements))
	for i, a := range recentAchievements {
		recentAchievementResponses[i] = toAchievementResponse(&a)
	}

	return &DashboardResponse{
		Profile: *toStudentProfileResponse(student),
		Summary: StudentSummaryResponse{
			AttendanceSummary: *attendanceSummary,
			GradeSummary:      *gradeSummary,
			BKSummary: BKSummaryResponse{
				TotalPoints:      totalPoints,
				ViolationCount:   violations,
				AchievementCount: achievements,
			},
		},
		RecentAttendance:   recentAttendanceResponses,
		RecentGrades:       recentGradeResponses,
		RecentAchievements: recentAchievementResponses,
	}, nil
}

// GetSummary retrieves summary statistics for the student
func (s *service) GetSummary(ctx context.Context, userID uint) (*StudentSummaryResponse, error) {
	student, err := s.repo.FindStudentByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get attendance summary (last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	attendanceSummary, _ := s.repo.GetAttendanceSummary(ctx, student.ID, thirtyDaysAgo, time.Now())
	if attendanceSummary == nil {
		attendanceSummary = &AttendanceSummaryResponse{}
	}

	// Get grade summary
	gradeSummary, _ := s.repo.GetGradeSummary(ctx, student.ID)
	if gradeSummary == nil {
		gradeSummary = &GradeSummaryResponse{}
	}

	// Get BK summary
	totalPoints, _ := s.repo.GetStudentAchievementPoints(ctx, student.ID)
	violations, achievements, _ := s.repo.GetStudentBKCounts(ctx, student.ID)

	return &StudentSummaryResponse{
		AttendanceSummary: *attendanceSummary,
		GradeSummary:      *gradeSummary,
		BKSummary: BKSummaryResponse{
			TotalPoints:      totalPoints,
			ViolationCount:   violations,
			AchievementCount: achievements,
		},
	}, nil
}

// GetAttendance retrieves attendance records for the student
// Requirements: 16.3 - Student views attendance history
func (s *service) GetAttendance(ctx context.Context, userID uint, filter AttendanceFilter) (*AttendanceListResponse, error) {
	student, err := s.repo.FindStudentByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	var startDate, endDate time.Time
	if filter.StartDate != "" {
		startDate, _ = time.Parse("2006-01-02", filter.StartDate)
	}
	if filter.EndDate != "" {
		endDate, _ = time.Parse("2006-01-02", filter.EndDate)
	}

	attendances, total, err := s.repo.GetStudentAttendance(ctx, student.ID, startDate, endDate, filter.Page, filter.PageSize)
	if err != nil {
		return nil, err
	}

	responses := make([]AttendanceResponse, len(attendances))
	for i, a := range attendances {
		responses[i] = toAttendanceResponse(&a)
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &AttendanceListResponse{
		Attendances: responses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetAttendanceSummary retrieves attendance summary for the student
func (s *service) GetAttendanceSummary(ctx context.Context, userID uint, startDateStr, endDateStr string) (*AttendanceSummaryResponse, error) {
	student, err := s.repo.FindStudentByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, _ = time.Parse("2006-01-02", startDateStr)
	}
	if endDateStr != "" {
		endDate, _ = time.Parse("2006-01-02", endDateStr)
	}

	return s.repo.GetAttendanceSummary(ctx, student.ID, startDate, endDate)
}

// GetGrades retrieves grades for the student
// Requirements: 16.4 - Student views all grades
func (s *service) GetGrades(ctx context.Context, userID uint, filter GradeFilter) (*GradeListResponse, error) {
	student, err := s.repo.FindStudentByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	grades, total, err := s.repo.GetStudentGrades(ctx, student.ID, filter.Page, filter.PageSize)
	if err != nil {
		return nil, err
	}

	responses := make([]GradeResponse, len(grades))
	for i, g := range grades {
		responses[i] = toGradeResponse(&g)
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &GradeListResponse{
		Grades: responses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetGradeSummary retrieves grade summary for the student
func (s *service) GetGradeSummary(ctx context.Context, userID uint) (*GradeSummaryResponse, error) {
	student, err := s.repo.FindStudentByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetGradeSummary(ctx, student.ID)
}

// GetBKInfo retrieves BK information for the student
// Requirements: 16.5 - Student views achievements and violations (summary only)
func (s *service) GetBKInfo(ctx context.Context, userID uint) (*BKInfoResponse, error) {
	student, err := s.repo.FindStudentByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get counts
	totalPoints, _ := s.repo.GetStudentAchievementPoints(ctx, student.ID)
	violationCount, achievementCount, _ := s.repo.GetStudentBKCounts(ctx, student.ID)

	// Get recent items (5 each)
	violations, _ := s.repo.GetStudentViolations(ctx, student.ID, 5)
	achievements, _ := s.repo.GetStudentAchievements(ctx, student.ID, 5)

	recentViolations := make([]ViolationSummaryResponse, len(violations))
	for i, v := range violations {
		recentViolations[i] = toViolationSummaryResponse(&v)
	}

	recentAchievements := make([]AchievementResponse, len(achievements))
	for i, a := range achievements {
		recentAchievements[i] = toAchievementResponse(&a)
	}

	return &BKInfoResponse{
		TotalPoints:        totalPoints,
		ViolationCount:     violationCount,
		AchievementCount:   achievementCount,
		RecentViolations:   recentViolations,
		RecentAchievements: recentAchievements,
	}, nil
}

// ==================== Response Converters ====================

func toStudentProfileResponse(s *models.Student) *StudentProfileResponse {
	className := ""
	grade := 0
	year := ""
	if s.Class.ID != 0 {
		className = s.Class.Name
		grade = s.Class.Grade
		year = s.Class.Year
	}

	return &StudentProfileResponse{
		ID:        s.ID,
		NIS:       s.NIS,
		NISN:      s.NISN,
		Name:      s.Name,
		ClassName: className,
		Grade:     grade,
		Year:      year,
		IsActive:  s.IsActive,
	}
}

func toAttendanceResponse(a *models.Attendance) AttendanceResponse {
	return AttendanceResponse{
		ID:           a.ID,
		Date:         a.Date.Format("2006-01-02"),
		CheckInTime:  a.CheckInTime,
		CheckOutTime: a.CheckOutTime,
		Status:       string(a.Status),
		Method:       string(a.Method),
	}
}

func toGradeResponse(g *models.Grade) GradeResponse {
	teacherName := ""
	if g.Creator.ID != 0 {
		teacherName = g.Creator.Name
		if teacherName == "" {
			teacherName = g.Creator.Username
		}
	}

	return GradeResponse{
		ID:          g.ID,
		Title:       g.Title,
		Score:       g.Score,
		Description: g.Description,
		TeacherName: teacherName,
		CreatedAt:   g.CreatedAt,
	}
}

func toViolationSummaryResponse(v *models.Violation) ViolationSummaryResponse {
	return ViolationSummaryResponse{
		ID:          v.ID,
		Category:    v.Category,
		Level:       string(v.Level),
		Description: v.Description,
		CreatedAt:   v.CreatedAt,
	}
}

func toAchievementResponse(a *models.Achievement) AchievementResponse {
	return AchievementResponse{
		ID:          a.ID,
		Title:       a.Title,
		Point:       a.Point,
		Description: a.Description,
		CreatedAt:   a.CreatedAt,
	}
}
