package bk

import (
	"time"

	"github.com/school-management/backend/internal/domain/models"
)

// ==================== Pagination ====================

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// ==================== Violation DTOs ====================

// CreateViolationRequest represents the request to create a violation
// Requirements: 6.1 - Violation SHALL require category, level, description, and student identifier
type CreateViolationRequest struct {
	StudentID   uint                  `json:"student_id" validate:"required"`
	CategoryID  *uint                 `json:"category_id"`
	Category    string                `json:"category" validate:"required"`
	Level       models.ViolationLevel `json:"level" validate:"required"`
	Point       *int                  `json:"point"`
	Description string                `json:"description" validate:"required"`
}

// ViolationResponse represents a violation in responses
type ViolationResponse struct {
	ID           uint                  `json:"id"`
	StudentID    uint                  `json:"student_id"`
	StudentName  string                `json:"student_name,omitempty"`
	StudentNIS   string                `json:"student_nis,omitempty"`
	StudentNISN  string                `json:"student_nisn,omitempty"`
	ClassName    string                `json:"class_name,omitempty"`
	CategoryID   *uint                 `json:"category_id,omitempty"`
	Category     string                `json:"category"`
	Level        models.ViolationLevel `json:"level"`
	Point        int                   `json:"point"`
	Description  string                `json:"description"`
	CreatedBy    uint                  `json:"created_by"`
	CreatorName  string                `json:"creator_name,omitempty"`
	CreatedAt    time.Time             `json:"created_at"`
}

// ==================== Violation Category DTOs ====================

// CreateViolationCategoryRequest represents the request to create a violation category
type CreateViolationCategoryRequest struct {
	Name         string                `json:"name" validate:"required"`
	DefaultPoint int                   `json:"default_point" validate:"max=0"`
	DefaultLevel models.ViolationLevel `json:"default_level" validate:"required"`
	Description  string                `json:"description"`
}

// UpdateViolationCategoryRequest represents the request to update a violation category
type UpdateViolationCategoryRequest struct {
	Name         string                `json:"name"`
	DefaultPoint *int                  `json:"default_point"`
	DefaultLevel models.ViolationLevel `json:"default_level"`
	Description  string                `json:"description"`
	IsActive     *bool                 `json:"is_active"`
}

// ViolationCategoryResponse represents a violation category in responses
type ViolationCategoryResponse struct {
	ID           uint                  `json:"id"`
	SchoolID     uint                  `json:"school_id"`
	Name         string                `json:"name"`
	DefaultPoint int                   `json:"default_point"`
	DefaultLevel models.ViolationLevel `json:"default_level"`
	Description  string                `json:"description"`
	IsActive     bool                  `json:"is_active"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
}

// ViolationCategoryListResponse represents a list of violation categories
type ViolationCategoryListResponse struct {
	Categories []ViolationCategoryResponse `json:"categories"`
}

// ViolationListResponse represents a paginated list of violations
type ViolationListResponse struct {
	Violations []ViolationResponse `json:"violations"`
	Pagination PaginationMeta      `json:"pagination"`
}

// ViolationFilter represents filter options for listing violations
type ViolationFilter struct {
	StudentID *uint   `query:"student_id"`
	ClassID   *uint   `query:"class_id"`
	Level     *string `query:"level"`
	Category  *string `query:"category"`
	StartDate *string `query:"start_date"`
	EndDate   *string `query:"end_date"`
	Page      int     `query:"page"`
	PageSize  int     `query:"page_size"`
}

// ==================== Achievement DTOs ====================

// CreateAchievementRequest represents the request to create an achievement
// Requirements: 7.1 - Achievement SHALL require title, point value, and description
type CreateAchievementRequest struct {
	StudentID   uint   `json:"student_id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Point       int    `json:"point" validate:"required,min=1"`
	Description string `json:"description"`
}

// AchievementResponse represents an achievement in responses
type AchievementResponse struct {
	ID          uint      `json:"id"`
	StudentID   uint      `json:"student_id"`
	StudentName string    `json:"student_name,omitempty"`
	StudentNIS  string    `json:"student_nis,omitempty"`
	StudentNISN string    `json:"student_nisn,omitempty"`
	ClassName   string    `json:"class_name,omitempty"`
	Title       string    `json:"title"`
	Point       int       `json:"point"`
	Description string    `json:"description"`
	CreatedBy   uint      `json:"created_by"`
	CreatorName string    `json:"creator_name,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// AchievementListResponse represents a paginated list of achievements
type AchievementListResponse struct {
	Achievements []AchievementResponse `json:"achievements"`
	Pagination   PaginationMeta        `json:"pagination"`
}

// AchievementPointsResponse represents total achievement points for a student
// Requirements: 7.3 - THE System SHALL display total achievement points
type AchievementPointsResponse struct {
	StudentID   uint   `json:"student_id"`
	StudentName string `json:"student_name"`
	TotalPoints int    `json:"total_points"`
}

// AchievementFilter represents filter options for listing achievements
type AchievementFilter struct {
	StudentID *uint   `query:"student_id"`
	ClassID   *uint   `query:"class_id"`
	StartDate *string `query:"start_date"`
	EndDate   *string `query:"end_date"`
	Page      int     `query:"page"`
	PageSize  int     `query:"page_size"`
}


// ==================== Permit DTOs ====================

// CreatePermitRequest represents the request to create an exit permit
// Requirements: 8.1 - Permit SHALL require reason, exit time, and responsible teacher
type CreatePermitRequest struct {
	StudentID          uint      `json:"student_id" validate:"required"`
	Reason             string    `json:"reason" validate:"required"`
	ExitTime           time.Time `json:"exit_time" validate:"required"`
	ResponsibleTeacher uint      `json:"responsible_teacher" validate:"required"`
}

// RecordReturnRequest represents the request to record student return
// Requirements: 8.4 - THE System SHALL allow recording of return time
type RecordReturnRequest struct {
	ReturnTime time.Time `json:"return_time" validate:"required"`
}

// PermitResponse represents a permit in responses
// Requirements: 8.5 - Permit document SHALL contain student info, reason, exit time, teacher, timestamp
type PermitResponse struct {
	ID                 uint       `json:"id"`
	StudentID          uint       `json:"student_id"`
	StudentName        string     `json:"student_name,omitempty"`
	StudentNIS         string     `json:"student_nis,omitempty"`
	StudentNISN        string     `json:"student_nisn,omitempty"`
	ClassName          string     `json:"class_name,omitempty"`
	Reason             string     `json:"reason"`
	ExitTime           time.Time  `json:"exit_time"`
	ReturnTime         *time.Time `json:"return_time,omitempty"`
	ResponsibleTeacher uint       `json:"responsible_teacher"`
	TeacherName        string     `json:"teacher_name,omitempty"`
	DocumentURL        string     `json:"document_url,omitempty"`
	CreatedBy          uint       `json:"created_by"`
	CreatorName        string     `json:"creator_name,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	HasReturned        bool       `json:"has_returned"`
}

// PermitListResponse represents a paginated list of permits
type PermitListResponse struct {
	Permits    []PermitResponse `json:"permits"`
	Pagination PaginationMeta   `json:"pagination"`
}

// PermitDocumentData represents data for generating permit document
// Requirements: 8.2, 8.5 - Permit document content
type PermitDocumentData struct {
	PermitID           uint      `json:"permit_id"`
	StudentName        string    `json:"student_name"`
	StudentNIS         string    `json:"student_nis"`
	StudentNISN        string    `json:"student_nisn"`
	ClassName          string    `json:"class_name"`
	SchoolName         string    `json:"school_name"`
	Reason             string    `json:"reason"`
	ExitTime           time.Time `json:"exit_time"`
	ResponsibleTeacher string    `json:"responsible_teacher"`
	GeneratedAt        time.Time `json:"generated_at"`
}

// PermitFilter represents filter options for listing permits
type PermitFilter struct {
	StudentID  *uint   `query:"student_id"`
	ClassID    *uint   `query:"class_id"`
	TeacherID  *uint   `query:"teacher_id"`
	HasReturned *bool  `query:"has_returned"`
	StartDate  *string `query:"start_date"`
	EndDate    *string `query:"end_date"`
	Page       int     `query:"page"`
	PageSize   int     `query:"page_size"`
}

// ==================== Counseling Note DTOs ====================

// CreateCounselingNoteRequest represents the request to create a counseling note
// Requirements: 9.1 - Counseling note SHALL require internal_note and optional parent_summary
type CreateCounselingNoteRequest struct {
	StudentID     uint   `json:"student_id" validate:"required"`
	InternalNote  string `json:"internal_note" validate:"required"`
	ParentSummary string `json:"parent_summary"`
}

// UpdateCounselingNoteRequest represents the request to update a counseling note
type UpdateCounselingNoteRequest struct {
	InternalNote  string `json:"internal_note"`
	ParentSummary string `json:"parent_summary"`
}

// CounselingNoteResponse represents a counseling note in responses (for parents/wali kelas)
// Requirements: 9.3, 9.4 - Only parent_summary visible to parents and wali kelas
type CounselingNoteResponse struct {
	ID            uint      `json:"id"`
	StudentID     uint      `json:"student_id"`
	StudentName   string    `json:"student_name,omitempty"`
	StudentNIS    string    `json:"student_nis,omitempty"`
	StudentNISN   string    `json:"student_nisn,omitempty"`
	ClassName     string    `json:"class_name,omitempty"`
	ParentSummary string    `json:"parent_summary"`
	CreatedBy     uint      `json:"created_by"`
	CreatorName   string    `json:"creator_name,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// CounselingNoteFullResponse represents a counseling note with internal note (for Guru BK only)
// Requirements: 9.3 - Internal note accessible only to Guru BK
type CounselingNoteFullResponse struct {
	ID            uint      `json:"id"`
	StudentID     uint      `json:"student_id"`
	StudentName   string    `json:"student_name,omitempty"`
	StudentNIS    string    `json:"student_nis,omitempty"`
	StudentNISN   string    `json:"student_nisn,omitempty"`
	ClassName     string    `json:"class_name,omitempty"`
	InternalNote  string    `json:"internal_note"`
	ParentSummary string    `json:"parent_summary"`
	CreatedBy     uint      `json:"created_by"`
	CreatorName   string    `json:"creator_name,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// CounselingNoteListResponse represents a paginated list of counseling notes
type CounselingNoteListResponse struct {
	Notes      []CounselingNoteResponse `json:"notes"`
	Pagination PaginationMeta           `json:"pagination"`
}

// CounselingNoteFullListResponse represents a paginated list of counseling notes with internal notes
type CounselingNoteFullListResponse struct {
	Notes      []CounselingNoteFullResponse `json:"notes"`
	Pagination PaginationMeta               `json:"pagination"`
}

// CounselingNoteFilter represents filter options for listing counseling notes
type CounselingNoteFilter struct {
	StudentID *uint   `query:"student_id"`
	ClassID   *uint   `query:"class_id"`
	StartDate *string `query:"start_date"`
	EndDate   *string `query:"end_date"`
	Page      int     `query:"page"`
	PageSize  int     `query:"page_size"`
}

// ==================== Student BK Profile DTOs ====================

// StudentBKProfileResponse represents a student's complete BK profile
// Requirements: 6.3, 7.5, 8.4, 9.5 - Timeline view per student
type StudentBKProfileResponse struct {
	StudentID         uint                         `json:"student_id"`
	StudentName       string                       `json:"student_name"`
	StudentNIS        string                       `json:"student_nis"`
	StudentNISN       string                       `json:"student_nisn"`
	ClassName         string                       `json:"class_name"`
	TotalPoints       int                          `json:"total_points"`
	ViolationCount    int                          `json:"violation_count"`
	AchievementCount  int                          `json:"achievement_count"`
	PermitCount       int                          `json:"permit_count"`
	CounselingCount   int                          `json:"counseling_count"`
	RecentViolations  []ViolationResponse          `json:"recent_violations,omitempty"`
	RecentAchievements []AchievementResponse       `json:"recent_achievements,omitempty"`
	RecentPermits     []PermitResponse             `json:"recent_permits,omitempty"`
	RecentCounseling  []CounselingNoteResponse     `json:"recent_counseling,omitempty"`
}

// StudentBKProfileFullResponse includes internal counseling notes (for Guru BK)
type StudentBKProfileFullResponse struct {
	StudentID         uint                         `json:"student_id"`
	StudentName       string                       `json:"student_name"`
	StudentNIS        string                       `json:"student_nis"`
	StudentNISN       string                       `json:"student_nisn"`
	ClassName         string                       `json:"class_name"`
	TotalPoints       int                          `json:"total_points"`
	ViolationCount    int                          `json:"violation_count"`
	AchievementCount  int                          `json:"achievement_count"`
	PermitCount       int                          `json:"permit_count"`
	CounselingCount   int                          `json:"counseling_count"`
	RecentViolations  []ViolationResponse          `json:"recent_violations,omitempty"`
	RecentAchievements []AchievementResponse       `json:"recent_achievements,omitempty"`
	RecentPermits     []PermitResponse             `json:"recent_permits,omitempty"`
	RecentCounseling  []CounselingNoteFullResponse `json:"recent_counseling,omitempty"`
}

// ==================== Dashboard DTOs ====================

// BKDashboardResponse represents the BK dashboard data
// Requirements: 6.1, 7.1 - Overview: recent violations, achievements
type BKDashboardResponse struct {
	TotalViolations    int                   `json:"total_violations"`
	TotalAchievements  int                   `json:"total_achievements"`
	TotalPermits       int                   `json:"total_permits"`
	ActivePermits      int                   `json:"active_permits"`
	TotalCounseling    int                   `json:"total_counseling"`
	RecentViolations   []ViolationResponse   `json:"recent_violations"`
	RecentAchievements []AchievementResponse `json:"recent_achievements"`
	StudentsNeedingAttention []StudentAttentionItem `json:"students_needing_attention"`
}

// StudentAttentionItem represents a student that needs attention
type StudentAttentionItem struct {
	StudentID      uint   `json:"student_id"`
	StudentName    string `json:"student_name"`
	ClassName      string `json:"class_name"`
	ViolationCount int    `json:"violation_count"`
	Reason         string `json:"reason"`
}
