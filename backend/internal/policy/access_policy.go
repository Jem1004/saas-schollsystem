package policy

import (
	"context"
	"errors"

	"github.com/school-management/backend/internal/domain/models"
)

// Access policy errors
var (
	ErrAccessDenied           = errors.New("access denied")
	ErrStudentNotInClass      = errors.New("student is not in your assigned class")
	ErrNotHomeroomTeacher     = errors.New("you are not the homeroom teacher for this class")
	ErrCannotAccessResource   = errors.New("you cannot access this resource")
	ErrInternalNotesRestricted = errors.New("internal notes are restricted to Guru BK only")
)

// AccessLevel represents the level of access a user has to a resource
type AccessLevel string

const (
	AccessLevelNone     AccessLevel = "none"
	AccessLevelReadOnly AccessLevel = "readonly"
	AccessLevelLimited  AccessLevel = "limited"  // Can see summary only
	AccessLevelFull     AccessLevel = "full"
)

// ResourceType represents the type of resource being accessed
type ResourceType string

const (
	ResourceStudent        ResourceType = "student"
	ResourceClass          ResourceType = "class"
	ResourceGrade          ResourceType = "grade"
	ResourceViolation      ResourceType = "violation"
	ResourceAchievement    ResourceType = "achievement"
	ResourcePermit         ResourceType = "permit"
	ResourceCounselingNote ResourceType = "counseling_note"
	ResourceHomeroomNote   ResourceType = "homeroom_note"
	ResourceAttendance     ResourceType = "attendance"
)

// UserContext holds the current user's context for access control
type UserContext struct {
	UserID   uint
	SchoolID *uint
	Role     models.UserRole
	ClassID  *uint // For wali kelas - their assigned class
}

// AccessPolicy defines the interface for access control decisions
// Requirements: 4.5 - THE System SHALL enforce role-based access control for all protected resources
// Requirements: 6.5 - WHEN a Wali_Kelas views BK data, THE System SHALL provide read-only access
// Requirements: 9.3 - THE System SHALL keep internal_note private and accessible only to Guru_BK
// Requirements: 9.4 - WHEN a Wali_Kelas views counseling data, THE System SHALL show only parent_summary
// Requirements: 10.5 - THE System SHALL validate that Wali_Kelas can only input grades for students in their assigned class
// Requirements: 11.4 - THE System SHALL validate that Wali_Kelas can only create notes for students in their assigned class
type AccessPolicy interface {
	// CanAccessStudent checks if user can access a specific student's data
	CanAccessStudent(ctx context.Context, user *UserContext, studentID uint) (bool, error)

	// CanModifyStudent checks if user can modify a specific student's data
	CanModifyStudent(ctx context.Context, user *UserContext, studentID uint) (bool, error)

	// CanAccessBKData checks if user can access BK data and returns the access level
	CanAccessBKData(user *UserContext) AccessLevel

	// CanViewInternalNotes checks if user can view internal counseling notes
	CanViewInternalNotes(user *UserContext) bool

	// CanModifyBKData checks if user can create/modify BK data (violations, achievements, permits, counseling)
	CanModifyBKData(user *UserContext) bool

	// CanAccessGrade checks if user can access grade data and returns the access level
	CanAccessGrade(ctx context.Context, user *UserContext, studentID uint) (AccessLevel, error)

	// CanModifyGrade checks if user can modify a specific grade
	CanModifyGrade(ctx context.Context, user *UserContext, studentID uint) (bool, error)

	// CanAccessHomeroomNote checks if user can access homeroom notes
	CanAccessHomeroomNote(ctx context.Context, user *UserContext, studentID uint) (AccessLevel, error)

	// CanModifyHomeroomNote checks if user can modify homeroom notes for a student
	CanModifyHomeroomNote(ctx context.Context, user *UserContext, studentID uint) (bool, error)

	// GetVisibleFields returns the fields visible to a user for a resource type
	GetVisibleFields(user *UserContext, resourceType ResourceType) []string

	// IsStudentInUserClass checks if a student is in the user's assigned class (for wali kelas)
	IsStudentInUserClass(ctx context.Context, user *UserContext, studentID uint) (bool, error)

	// GetUserAssignedClassID returns the class ID assigned to a wali kelas user
	GetUserAssignedClassID(ctx context.Context, userID uint) (*uint, error)
}
