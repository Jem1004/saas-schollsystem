package school

import "time"

// ==================== Class DTOs ====================

// CreateClassRequest represents the request to create a new class
// Requirements: 3.1 - WHEN an Admin_Sekolah creates a class, THE System SHALL associate it with the school tenant
type CreateClassRequest struct {
	Name              string `json:"name" validate:"required"`
	Grade             int    `json:"grade" validate:"required,min=1"`
	Year              string `json:"year" validate:"required"`
	HomeroomTeacherID *uint  `json:"homeroom_teacher_id"`
}

// UpdateClassRequest represents the request to update a class
type UpdateClassRequest struct {
	Name              *string `json:"name"`
	Grade             *int    `json:"grade"`
	Year              *string `json:"year"`
	HomeroomTeacherID *uint   `json:"homeroom_teacher_id"`
}

// ClassResponse represents the class data in responses
type ClassResponse struct {
	ID                uint             `json:"id"`
	SchoolID          uint             `json:"school_id"`
	Name              string           `json:"name"`
	Grade             int              `json:"grade"`
	Year              string           `json:"year"`
	HomeroomTeacherID *uint            `json:"homeroom_teacher_id"`
	HomeroomTeacher   *TeacherResponse `json:"homeroom_teacher,omitempty"`
	StudentCount      int64            `json:"student_count"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}

// ClassListResponse represents a paginated list of classes
type ClassListResponse struct {
	Classes    []ClassResponse `json:"classes"`
	Pagination PaginationMeta  `json:"pagination"`
}

// ClassFilter represents filter options for listing classes
type ClassFilter struct {
	Name     string `query:"name"`
	Grade    *int   `query:"grade"`
	Year     string `query:"year"`
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
}

// DefaultClassFilter returns default filter values
func DefaultClassFilter() ClassFilter {
	return ClassFilter{
		Page:     1,
		PageSize: 20,
	}
}

// ==================== Student DTOs ====================

// CreateStudentRequest represents the request to create a new student
// Requirements: 3.2 - WHEN an Admin_Sekolah registers a student, THE System SHALL require NIS, NISN, name, and class assignment
type CreateStudentRequest struct {
	NIS      string `json:"nis" validate:"required"`
	NISN     string `json:"nisn" validate:"required"`
	Name     string `json:"name" validate:"required"`
	ClassID  uint   `json:"class_id" validate:"required"`
	RFIDCode string `json:"rfid_code"`
}

// UpdateStudentRequest represents the request to update a student
type UpdateStudentRequest struct {
	NIS      *string `json:"nis"`
	Name     *string `json:"name"`
	ClassID  *uint   `json:"class_id"`
	RFIDCode *string `json:"rfid_code"`
	IsActive *bool   `json:"is_active"`
}

// StudentResponse represents the student data in responses
type StudentResponse struct {
	ID        uint           `json:"id"`
	SchoolID  uint           `json:"school_id"`
	ClassID   uint           `json:"class_id"`
	NIS       string         `json:"nis"`
	NISN      string         `json:"nisn"`
	Name      string         `json:"name"`
	RFIDCode  string         `json:"rfid_code"`
	IsActive  bool           `json:"is_active"`
	Class     *ClassResponse `json:"class,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// StudentListResponse represents a paginated list of students
type StudentListResponse struct {
	Students   []StudentResponse `json:"students"`
	Pagination PaginationMeta    `json:"pagination"`
}

// StudentFilter represents filter options for listing students
type StudentFilter struct {
	Name     string `query:"name"`
	NIS      string `query:"nis"`
	NISN     string `query:"nisn"`
	ClassID  *uint  `query:"class_id"`
	IsActive *bool  `query:"is_active"`
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
}

// DefaultStudentFilter returns default filter values
func DefaultStudentFilter() StudentFilter {
	return StudentFilter{
		Page:     1,
		PageSize: 20,
	}
}


// ==================== Parent DTOs ====================

// CreateParentRequest represents the request to create a new parent
// Requirements: 3.3 - WHEN an Admin_Sekolah registers a parent, THE System SHALL link the parent to one or more students
type CreateParentRequest struct {
	Name       string `json:"name" validate:"required"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Password   string `json:"password" validate:"required,min=8"`
	StudentIDs []uint `json:"student_ids"`
}

// UpdateParentRequest represents the request to update a parent
type UpdateParentRequest struct {
	Name  *string `json:"name"`
	Phone *string `json:"phone"`
}

// LinkParentStudentRequest represents the request to link a parent to students
type LinkParentStudentRequest struct {
	StudentIDs []uint `json:"student_ids" validate:"required,min=1"`
}

// ParentResponse represents the parent data in responses
type ParentResponse struct {
	ID        uint              `json:"id"`
	SchoolID  uint              `json:"school_id"`
	UserID    uint              `json:"user_id"`
	Name      string            `json:"name"`
	Phone     string            `json:"phone"`
	Email     string            `json:"email,omitempty"`
	Students  []StudentResponse `json:"students,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// ParentListResponse represents a paginated list of parents
type ParentListResponse struct {
	Parents    []ParentResponse `json:"parents"`
	Pagination PaginationMeta   `json:"pagination"`
}

// ParentFilter represents filter options for listing parents
type ParentFilter struct {
	Name     string `query:"name"`
	Phone    string `query:"phone"`
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
}

// DefaultParentFilter returns default filter values
func DefaultParentFilter() ParentFilter {
	return ParentFilter{
		Page:     1,
		PageSize: 20,
	}
}

// ==================== Common DTOs ====================

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// TeacherResponse represents a simplified teacher response
type TeacherResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

// AssignHomeroomTeacherRequest represents the request to assign a homeroom teacher
type AssignHomeroomTeacherRequest struct {
	TeacherID uint `json:"teacher_id" validate:"required"`
}
