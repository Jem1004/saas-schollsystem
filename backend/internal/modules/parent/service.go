package parent

import (
	"context"
	"time"

	"github.com/school-management/backend/internal/domain/models"
)

// Service defines the interface for parent business logic
type Service interface {
	// Child operations
	GetLinkedChildren(ctx context.Context, userID uint) (*ChildListResponse, error)
	GetChildDashboard(ctx context.Context, userID, studentID uint) (*ChildDashboardResponse, error)

	// Attendance operations
	GetChildAttendance(ctx context.Context, userID, studentID uint, filter AttendanceFilter) (*ChildAttendanceListResponse, error)
	GetChildAttendanceSummary(ctx context.Context, userID, studentID uint, startDate, endDate string) (*AttendanceSummaryResponse, error)

	// Grade operations
	GetChildGrades(ctx context.Context, userID, studentID uint, filter GradeFilter) (*ChildGradeListResponse, error)
	GetChildGradeSummary(ctx context.Context, userID, studentID uint) (*GradeSummaryResponse, error)

	// Homeroom note operations
	GetChildNotes(ctx context.Context, userID, studentID uint, filter NoteFilter) (*ChildNoteListResponse, error)

	// BK operations
	GetChildBKInfo(ctx context.Context, userID, studentID uint) (*ChildBKInfoResponse, error)
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new parent service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// GetLinkedChildren retrieves all children linked to a parent
// Requirements: 12.2 - Authentication succeeds SHALL return access to all linked children's data
func (s *service) GetLinkedChildren(ctx context.Context, userID uint) (*ChildListResponse, error) {
	parent, err := s.repo.FindParentByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	students, err := s.repo.GetLinkedStudents(ctx, parent.ID)
	if err != nil {
		return nil, err
	}

	children := make([]ChildResponse, len(students))
	for i, student := range students {
		children[i] = toChildResponse(&student)
	}

	return &ChildListResponse{Children: children}, nil
}

// GetChildDashboard retrieves dashboard data for a specific child
// Requirements: 15.1, 15.2 - Parent can view grades and notes for their child
func (s *service) GetChildDashboard(ctx context.Context, userID, studentID uint) (*ChildDashboardResponse, error) {
	// Validate parent has access to this student
	if err := s.validateAccess(ctx, userID, studentID); err != nil {
		return nil, err
	}

	student, err := s.repo.FindStudentByID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	// Get attendance summary (last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	attendanceSummary, _ := s.repo.GetAttendanceSummary(ctx, studentID, thirtyDaysAgo, time.Now())
	if attendanceSummary == nil {
		attendanceSummary = &AttendanceSummaryResponse{StudentID: studentID, StudentName: student.Name}
	}

	// Get grade summary
	gradeSummary, _ := s.repo.GetGradeSummary(ctx, studentID)
	if gradeSummary == nil {
		gradeSummary = &GradeSummaryResponse{StudentID: studentID, StudentName: student.Name}
	}

	// Get BK summary
	totalPoints, _ := s.repo.GetStudentAchievementPoints(ctx, studentID)
	violations, achievements, _, _, _ := s.repo.GetStudentBKCounts(ctx, studentID)

	// Get recent grades (5)
	grades, _, _ := s.repo.GetStudentGrades(ctx, studentID, 1, 5)
	recentGrades := make([]ChildGradeResponse, len(grades))
	for i, g := range grades {
		recentGrades[i] = toChildGradeResponse(&g)
	}

	// Get recent notes (5)
	notes, _, _ := s.repo.GetStudentNotes(ctx, studentID, 1, 5)
	recentNotes := make([]ChildNoteResponse, len(notes))
	for i, n := range notes {
		recentNotes[i] = toChildNoteResponse(&n)
	}

	return &ChildDashboardResponse{
		Student:           toChildResponse(student),
		AttendanceSummary: *attendanceSummary,
		GradeSummary:      *gradeSummary,
		BKSummary: ChildBKSummaryResponse{
			TotalPoints:      totalPoints,
			ViolationCount:   violations,
			AchievementCount: achievements,
		},
		RecentGrades: recentGrades,
		RecentNotes:  recentNotes,
	}, nil
}

// GetChildAttendance retrieves attendance records for a child
func (s *service) GetChildAttendance(ctx context.Context, userID, studentID uint, filter AttendanceFilter) (*ChildAttendanceListResponse, error) {
	if err := s.validateAccess(ctx, userID, studentID); err != nil {
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

	attendances, total, err := s.repo.GetStudentAttendance(ctx, studentID, startDate, endDate, filter.Page, filter.PageSize)
	if err != nil {
		return nil, err
	}

	responses := make([]ChildAttendanceResponse, len(attendances))
	for i, a := range attendances {
		responses[i] = toChildAttendanceResponse(&a)
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &ChildAttendanceListResponse{
		Attendances: responses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetChildAttendanceSummary retrieves attendance summary for a child
func (s *service) GetChildAttendanceSummary(ctx context.Context, userID, studentID uint, startDateStr, endDateStr string) (*AttendanceSummaryResponse, error) {
	if err := s.validateAccess(ctx, userID, studentID); err != nil {
		return nil, err
	}

	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, _ = time.Parse("2006-01-02", startDateStr)
	}
	if endDateStr != "" {
		endDate, _ = time.Parse("2006-01-02", endDateStr)
	}

	return s.repo.GetAttendanceSummary(ctx, studentID, startDate, endDate)
}

// GetChildGrades retrieves grades for a child
// Requirements: 15.1 - Parent opens grades section SHALL display all grades sorted by date
func (s *service) GetChildGrades(ctx context.Context, userID, studentID uint, filter GradeFilter) (*ChildGradeListResponse, error) {
	if err := s.validateAccess(ctx, userID, studentID); err != nil {
		return nil, err
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	grades, total, err := s.repo.GetStudentGrades(ctx, studentID, filter.Page, filter.PageSize)
	if err != nil {
		return nil, err
	}

	responses := make([]ChildGradeResponse, len(grades))
	for i, g := range grades {
		responses[i] = toChildGradeResponse(&g)
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &ChildGradeListResponse{
		Grades: responses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetChildGradeSummary retrieves grade summary for a child
func (s *service) GetChildGradeSummary(ctx context.Context, userID, studentID uint) (*GradeSummaryResponse, error) {
	if err := s.validateAccess(ctx, userID, studentID); err != nil {
		return nil, err
	}

	return s.repo.GetGradeSummary(ctx, studentID)
}

// GetChildNotes retrieves homeroom notes for a child
// Requirements: 15.2 - Parent opens notes section SHALL display all homeroom notes
func (s *service) GetChildNotes(ctx context.Context, userID, studentID uint, filter NoteFilter) (*ChildNoteListResponse, error) {
	if err := s.validateAccess(ctx, userID, studentID); err != nil {
		return nil, err
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	notes, total, err := s.repo.GetStudentNotes(ctx, studentID, filter.Page, filter.PageSize)
	if err != nil {
		return nil, err
	}

	responses := make([]ChildNoteResponse, len(notes))
	for i, n := range notes {
		responses[i] = toChildNoteResponse(&n)
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &ChildNoteListResponse{
		Notes: responses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetChildBKInfo retrieves BK information for a child
// Requirements: 14.4 - Parent views violations, achievements, and permits
// Requirements: 14.5 - Only parent_summary for counseling notes
func (s *service) GetChildBKInfo(ctx context.Context, userID, studentID uint) (*ChildBKInfoResponse, error) {
	if err := s.validateAccess(ctx, userID, studentID); err != nil {
		return nil, err
	}

	// Get counts
	totalPoints, _ := s.repo.GetStudentAchievementPoints(ctx, studentID)
	violationCount, achievementCount, permitCount, counselingCount, _ := s.repo.GetStudentBKCounts(ctx, studentID)

	// Get recent items (5 each)
	violations, _ := s.repo.GetStudentViolations(ctx, studentID, 5)
	achievements, _ := s.repo.GetStudentAchievements(ctx, studentID, 5)
	permits, _ := s.repo.GetStudentPermits(ctx, studentID, 5)
	counselingNotes, _ := s.repo.GetStudentCounselingNotes(ctx, studentID, 5)

	recentViolations := make([]ChildViolationResponse, len(violations))
	for i, v := range violations {
		recentViolations[i] = toChildViolationResponse(&v)
	}

	recentAchievements := make([]ChildAchievementResponse, len(achievements))
	for i, a := range achievements {
		recentAchievements[i] = toChildAchievementResponse(&a)
	}

	recentPermits := make([]ChildPermitResponse, len(permits))
	for i, p := range permits {
		recentPermits[i] = toChildPermitResponse(&p)
	}

	recentCounseling := make([]ChildCounselingResponse, len(counselingNotes))
	for i, c := range counselingNotes {
		recentCounseling[i] = toChildCounselingResponse(&c)
	}

	return &ChildBKInfoResponse{
		TotalPoints:        totalPoints,
		ViolationCount:     violationCount,
		AchievementCount:   achievementCount,
		PermitCount:        permitCount,
		CounselingCount:    counselingCount,
		RecentViolations:   recentViolations,
		RecentAchievements: recentAchievements,
		RecentPermits:      recentPermits,
		RecentCounseling:   recentCounseling,
	}, nil
}

// validateAccess validates that a parent has access to a student
// Requirements: 12.2 - Parent SHALL only access linked children's data
func (s *service) validateAccess(ctx context.Context, userID, studentID uint) error {
	parent, err := s.repo.FindParentByUserID(ctx, userID)
	if err != nil {
		return err
	}

	linked, err := s.repo.IsStudentLinked(ctx, parent.ID, studentID)
	if err != nil {
		return err
	}

	if !linked {
		return ErrNotLinked
	}

	return nil
}

// ==================== Response Converters ====================

func toChildResponse(s *models.Student) ChildResponse {
	className := ""
	grade := 0
	if s.Class.ID != 0 {
		className = s.Class.Name
		grade = s.Class.Grade
	}

	return ChildResponse{
		ID:        s.ID,
		NIS:       s.NIS,
		NISN:      s.NISN,
		Name:      s.Name,
		ClassName: className,
		Grade:     grade,
		IsActive:  s.IsActive,
	}
}

func toChildAttendanceResponse(a *models.Attendance) ChildAttendanceResponse {
	return ChildAttendanceResponse{
		ID:           a.ID,
		Date:         a.Date.Format("2006-01-02"),
		CheckInTime:  a.CheckInTime,
		CheckOutTime: a.CheckOutTime,
		Status:       string(a.Status),
		Method:       string(a.Method),
	}
}

func toChildGradeResponse(g *models.Grade) ChildGradeResponse {
	teacherName := ""
	if g.Creator.ID != 0 {
		teacherName = g.Creator.Name
		if teacherName == "" {
			teacherName = g.Creator.Username
		}
	}

	return ChildGradeResponse{
		ID:          g.ID,
		Title:       g.Title,
		Score:       g.Score,
		Description: g.Description,
		TeacherName: teacherName,
		CreatedAt:   g.CreatedAt,
	}
}

func toChildNoteResponse(n *models.HomeroomNote) ChildNoteResponse {
	teacherName := ""
	if n.Teacher.ID != 0 {
		teacherName = n.Teacher.Name
		if teacherName == "" {
			teacherName = n.Teacher.Username
		}
	}

	return ChildNoteResponse{
		ID:          n.ID,
		Content:     n.Content,
		TeacherName: teacherName,
		CreatedAt:   n.CreatedAt,
	}
}

func toChildViolationResponse(v *models.Violation) ChildViolationResponse {
	return ChildViolationResponse{
		ID:          v.ID,
		Category:    v.Category,
		Level:       string(v.Level),
		Description: v.Description,
		CreatedAt:   v.CreatedAt,
	}
}

func toChildAchievementResponse(a *models.Achievement) ChildAchievementResponse {
	return ChildAchievementResponse{
		ID:          a.ID,
		Title:       a.Title,
		Point:       a.Point,
		Description: a.Description,
		CreatedAt:   a.CreatedAt,
	}
}

func toChildPermitResponse(p *models.Permit) ChildPermitResponse {
	teacherName := ""
	if p.Teacher.ID != 0 {
		teacherName = p.Teacher.Name
		if teacherName == "" {
			teacherName = p.Teacher.Username
		}
	}

	return ChildPermitResponse{
		ID:          p.ID,
		Reason:      p.Reason,
		ExitTime:    p.ExitTime,
		ReturnTime:  p.ReturnTime,
		TeacherName: teacherName,
		HasReturned: p.HasReturned(),
		CreatedAt:   p.CreatedAt,
	}
}

func toChildCounselingResponse(c *models.CounselingNote) ChildCounselingResponse {
	return ChildCounselingResponse{
		ID:            c.ID,
		ParentSummary: c.ParentSummary,
		CreatedAt:     c.CreatedAt,
	}
}
