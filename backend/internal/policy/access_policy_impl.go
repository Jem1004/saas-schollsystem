package policy

import (
	"context"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

// accessPolicy implements the AccessPolicy interface
type accessPolicy struct {
	db *gorm.DB
}

// NewAccessPolicy creates a new access policy instance
func NewAccessPolicy(db *gorm.DB) AccessPolicy {
	return &accessPolicy{db: db}
}

// CanAccessStudent checks if user can access a specific student's data
// Requirements: 4.5 - Role-based access control
func (p *accessPolicy) CanAccessStudent(ctx context.Context, user *UserContext, studentID uint) (bool, error) {
	switch user.Role {
	case models.RoleSuperAdmin:
		// Super admin can access all students
		return true, nil

	case models.RoleAdminSekolah:
		// Admin can access students in their school
		return p.isStudentInSchool(ctx, studentID, user.SchoolID)

	case models.RoleGuruBK:
		// Guru BK can access all students in their school
		return p.isStudentInSchool(ctx, studentID, user.SchoolID)

	case models.RoleWaliKelas:
		// Wali kelas can access students in their assigned class
		return p.IsStudentInUserClass(ctx, user, studentID)

	case models.RoleParent:
		// Parent can only access their linked children
		return p.isParentOfStudent(ctx, user.UserID, studentID)

	case models.RoleStudent:
		// Student can only access their own data
		return p.isOwnStudentRecord(ctx, user.UserID, studentID)

	default:
		return false, nil
	}
}

// CanModifyStudent checks if user can modify a specific student's data
func (p *accessPolicy) CanModifyStudent(ctx context.Context, user *UserContext, studentID uint) (bool, error) {
	switch user.Role {
	case models.RoleSuperAdmin:
		return true, nil

	case models.RoleAdminSekolah:
		return p.isStudentInSchool(ctx, studentID, user.SchoolID)

	default:
		// Other roles cannot modify student data
		return false, nil
	}
}

// CanAccessBKData checks if user can access BK data and returns the access level
// Requirements: 6.5 - WHEN a Wali_Kelas views BK data, THE System SHALL provide read-only access
func (p *accessPolicy) CanAccessBKData(user *UserContext) AccessLevel {
	switch user.Role {
	case models.RoleSuperAdmin, models.RoleAdminSekolah:
		return AccessLevelFull

	case models.RoleGuruBK:
		return AccessLevelFull

	case models.RoleWaliKelas:
		// Read-only access for wali kelas
		return AccessLevelReadOnly

	case models.RoleParent, models.RoleStudent:
		// Limited access - summary only
		return AccessLevelLimited

	default:
		return AccessLevelNone
	}
}

// CanViewInternalNotes checks if user can view internal counseling notes
// Requirements: 9.3 - THE System SHALL keep internal_note private and accessible only to Guru_BK
// Requirements: 9.4 - WHEN a Wali_Kelas views counseling data, THE System SHALL show only parent_summary
func (p *accessPolicy) CanViewInternalNotes(user *UserContext) bool {
	// Only Guru BK can view internal notes
	return user.Role == models.RoleGuruBK
}

// CanModifyBKData checks if user can create/modify BK data
func (p *accessPolicy) CanModifyBKData(user *UserContext) bool {
	switch user.Role {
	case models.RoleSuperAdmin, models.RoleAdminSekolah:
		return true

	case models.RoleGuruBK:
		return true

	default:
		// Wali kelas and others cannot modify BK data
		return false
	}
}

// CanAccessGrade checks if user can access grade data and returns the access level
func (p *accessPolicy) CanAccessGrade(ctx context.Context, user *UserContext, studentID uint) (AccessLevel, error) {
	switch user.Role {
	case models.RoleSuperAdmin, models.RoleAdminSekolah:
		return AccessLevelFull, nil

	case models.RoleWaliKelas:
		// Wali kelas can access grades for students in their class
		inClass, err := p.IsStudentInUserClass(ctx, user, studentID)
		if err != nil {
			return AccessLevelNone, err
		}
		if inClass {
			return AccessLevelFull, nil
		}
		return AccessLevelNone, nil

	case models.RoleGuruBK:
		// Can view grades in their school
		inSchool, err := p.isStudentInSchool(ctx, studentID, user.SchoolID)
		if err != nil {
			return AccessLevelNone, err
		}
		if inSchool {
			return AccessLevelReadOnly, nil
		}
		return AccessLevelNone, nil

	case models.RoleParent:
		// Parent can view their children's grades
		isParent, err := p.isParentOfStudent(ctx, user.UserID, studentID)
		if err != nil {
			return AccessLevelNone, err
		}
		if isParent {
			return AccessLevelReadOnly, nil
		}
		return AccessLevelNone, nil

	case models.RoleStudent:
		// Student can view their own grades
		isOwn, err := p.isOwnStudentRecord(ctx, user.UserID, studentID)
		if err != nil {
			return AccessLevelNone, err
		}
		if isOwn {
			return AccessLevelReadOnly, nil
		}
		return AccessLevelNone, nil

	default:
		return AccessLevelNone, nil
	}
}

// CanModifyGrade checks if user can modify a specific grade
// Requirements: 10.5 - THE System SHALL validate that Wali_Kelas can only input grades for students in their assigned class
func (p *accessPolicy) CanModifyGrade(ctx context.Context, user *UserContext, studentID uint) (bool, error) {
	switch user.Role {
	case models.RoleSuperAdmin, models.RoleAdminSekolah:
		return true, nil

	case models.RoleWaliKelas:
		// Wali kelas can only modify grades for students in their class
		return p.IsStudentInUserClass(ctx, user, studentID)

	default:
		return false, nil
	}
}

// CanAccessHomeroomNote checks if user can access homeroom notes
func (p *accessPolicy) CanAccessHomeroomNote(ctx context.Context, user *UserContext, studentID uint) (AccessLevel, error) {
	switch user.Role {
	case models.RoleSuperAdmin, models.RoleAdminSekolah:
		return AccessLevelFull, nil

	case models.RoleWaliKelas:
		// Wali kelas can access notes for students in their class
		inClass, err := p.IsStudentInUserClass(ctx, user, studentID)
		if err != nil {
			return AccessLevelNone, err
		}
		if inClass {
			return AccessLevelFull, nil
		}
		return AccessLevelNone, nil

	case models.RoleGuruBK:
		// Guru BK can view homeroom notes
		inSchool, err := p.isStudentInSchool(ctx, studentID, user.SchoolID)
		if err != nil {
			return AccessLevelNone, err
		}
		if inSchool {
			return AccessLevelReadOnly, nil
		}
		return AccessLevelNone, nil

	case models.RoleParent:
		// Parent can view their children's notes
		isParent, err := p.isParentOfStudent(ctx, user.UserID, studentID)
		if err != nil {
			return AccessLevelNone, err
		}
		if isParent {
			return AccessLevelReadOnly, nil
		}
		return AccessLevelNone, nil

	case models.RoleStudent:
		// Student can view their own notes
		isOwn, err := p.isOwnStudentRecord(ctx, user.UserID, studentID)
		if err != nil {
			return AccessLevelNone, err
		}
		if isOwn {
			return AccessLevelReadOnly, nil
		}
		return AccessLevelNone, nil

	default:
		return AccessLevelNone, nil
	}
}

// CanModifyHomeroomNote checks if user can modify homeroom notes for a student
// Requirements: 11.4 - THE System SHALL validate that Wali_Kelas can only create notes for students in their assigned class
func (p *accessPolicy) CanModifyHomeroomNote(ctx context.Context, user *UserContext, studentID uint) (bool, error) {
	switch user.Role {
	case models.RoleSuperAdmin, models.RoleAdminSekolah:
		return true, nil

	case models.RoleWaliKelas:
		// Wali kelas can only modify notes for students in their class
		return p.IsStudentInUserClass(ctx, user, studentID)

	default:
		return false, nil
	}
}

// GetVisibleFields returns the fields visible to a user for a resource type
// Requirements: 9.3, 9.4 - Control visibility of internal notes
func (p *accessPolicy) GetVisibleFields(user *UserContext, resourceType ResourceType) []string {
	switch resourceType {
	case ResourceCounselingNote:
		if p.CanViewInternalNotes(user) {
			// Guru BK can see all fields
			return []string{"id", "student_id", "internal_note", "parent_summary", "created_by", "created_at"}
		}
		// Others can only see parent_summary
		return []string{"id", "student_id", "parent_summary", "created_by", "created_at"}

	case ResourceViolation:
		if user.Role == models.RoleParent || user.Role == models.RoleStudent {
			// Parents and students see summary only
			return []string{"id", "student_id", "category", "level", "created_at"}
		}
		return []string{"id", "student_id", "category", "level", "description", "created_by", "created_at"}

	case ResourceAchievement:
		return []string{"id", "student_id", "title", "point", "description", "created_by", "created_at"}

	case ResourcePermit:
		return []string{"id", "student_id", "reason", "exit_time", "return_time", "responsible_teacher", "document_url", "created_by", "created_at"}

	case ResourceGrade:
		return []string{"id", "student_id", "title", "score", "description", "created_by", "created_at", "updated_at"}

	case ResourceHomeroomNote:
		return []string{"id", "student_id", "teacher_id", "content", "created_at", "updated_at"}

	default:
		return []string{}
	}
}

// IsStudentInUserClass checks if a student is in the user's assigned class
func (p *accessPolicy) IsStudentInUserClass(ctx context.Context, user *UserContext, studentID uint) (bool, error) {
	if user.Role != models.RoleWaliKelas {
		return false, nil
	}

	// Get the user's assigned class
	classID, err := p.GetUserAssignedClassID(ctx, user.UserID)
	if err != nil {
		return false, err
	}
	if classID == nil {
		return false, nil
	}

	// Check if student is in that class
	var count int64
	err = p.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("id = ? AND class_id = ?", studentID, *classID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetUserAssignedClassID returns the class ID assigned to a wali kelas user
func (p *accessPolicy) GetUserAssignedClassID(ctx context.Context, userID uint) (*uint, error) {
	var class models.Class
	err := p.db.WithContext(ctx).
		Where("homeroom_teacher_id = ?", userID).
		First(&class).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &class.ID, nil
}

// Helper methods

// isStudentInSchool checks if a student belongs to a specific school
func (p *accessPolicy) isStudentInSchool(ctx context.Context, studentID uint, schoolID *uint) (bool, error) {
	if schoolID == nil {
		return false, nil
	}

	var count int64
	err := p.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("id = ? AND school_id = ?", studentID, *schoolID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// isParentOfStudent checks if a user (parent) is linked to a student
func (p *accessPolicy) isParentOfStudent(ctx context.Context, userID uint, studentID uint) (bool, error) {
	// First get the parent record for this user
	var parent models.Parent
	err := p.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&parent).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	// Check if parent is linked to the student
	var count int64
	err = p.db.WithContext(ctx).
		Table("student_parents").
		Where("parent_id = ? AND student_id = ?", parent.ID, studentID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// isOwnStudentRecord checks if a user's student record matches the given student ID
func (p *accessPolicy) isOwnStudentRecord(ctx context.Context, userID uint, studentID uint) (bool, error) {
	var count int64
	err := p.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("id = ? AND user_id = ?", studentID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
