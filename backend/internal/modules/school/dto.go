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
	ID                uint               `json:"id"`
	SchoolID          uint               `json:"school_id"`
	Name              string             `json:"name"`
	Grade             int                `json:"grade"`
	Year              string             `json:"year"`
	HomeroomTeacherID *uint              `json:"homeroom_teacher_id"`
	HomeroomTeacher   *TeacherResponse   `json:"homeroom_teacher,omitempty"`
	Counselors        []CounselorResponse `json:"counselors,omitempty"`
	StudentCount      int64              `json:"student_count"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
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
	NIS           string `json:"nis" validate:"required"`
	NISN          string `json:"nisn" validate:"required"`
	Name          string `json:"name" validate:"required"`
	ClassID       uint   `json:"class_id" validate:"required"`
	RFIDCode      string `json:"rfid_code"`
	CreateAccount bool   `json:"create_account"` // If true, create user account for mobile login
}

// UpdateStudentRequest represents the request to update a student
type UpdateStudentRequest struct {
	NIS      *string `json:"nis"`
	Name     *string `json:"name"`
	ClassID  *int    `json:"class_id"`
	RFIDCode *string `json:"rfid_code"`
	IsActive *bool   `json:"is_active"`
}

// CreateStudentAccountRequest represents the request to create account for existing student
type CreateStudentAccountRequest struct {
	Password string `json:"password"` // Optional, will auto-generate if empty
}

// StudentResponse represents the student data in responses
type StudentResponse struct {
	ID                uint           `json:"id"`
	SchoolID          uint           `json:"school_id"`
	ClassID           uint           `json:"class_id"`
	ClassName         string         `json:"class_name,omitempty"`
	NIS               string         `json:"nis"`
	NISN              string         `json:"nisn"`
	Name              string         `json:"name"`
	RFIDCode          string         `json:"rfid_code"`
	IsActive          bool           `json:"is_active"`
	HasAccount        bool           `json:"has_account"`
	Username          string         `json:"username,omitempty"`
	TemporaryPassword string         `json:"temporary_password,omitempty"`
	Class             *ClassResponse `json:"class,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
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
// Username will be phone number (primary) or email (secondary)
type CreateParentRequest struct {
	Name       string `json:"name" validate:"required"`
	Phone      string `json:"phone" validate:"required"` // Required, used as primary username
	Email      string `json:"email"`                     // Optional, used as secondary login
	Password   string `json:"password"`                  // Optional, will auto-generate if empty
	StudentIDs []uint `json:"student_ids"`
}

// UpdateParentRequest represents the request to update a parent
type UpdateParentRequest struct {
	Name       *string `json:"name"`
	Phone      *string `json:"phone"`
	Email      *string `json:"email"`
	StudentIDs []uint  `json:"student_ids"`
}

// LinkParentStudentRequest represents the request to link a parent to students
type LinkParentStudentRequest struct {
	StudentIDs []uint `json:"student_ids" validate:"required,min=1"`
}

// ParentResponse represents the parent data in responses
type ParentResponse struct {
	ID                uint              `json:"id"`
	SchoolID          uint              `json:"school_id"`
	UserID            uint              `json:"user_id"`
	Name              string            `json:"name"`
	Phone             string            `json:"phone"`
	Email             string            `json:"email,omitempty"`
	Username          string            `json:"username,omitempty"`
	TemporaryPassword string            `json:"temporary_password,omitempty"`
	Students          []StudentResponse `json:"students,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
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

// ResetPasswordResponse represents the response after resetting a password
type ResetPasswordResponse struct {
	Username          string `json:"username"`
	TemporaryPassword string `json:"temporary_password"`
	Message           string `json:"message"`
}

// TeacherResponse represents a simplified teacher response
type TeacherResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

// CounselorResponse represents a simplified counselor (guru BK) response
type CounselorResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

// AssignHomeroomTeacherRequest represents the request to assign a homeroom teacher
type AssignHomeroomTeacherRequest struct {
	TeacherID uint `json:"teacher_id" validate:"required"`
}

// AssignCounselorsRequest represents the request to assign counselors to a class
type AssignCounselorsRequest struct {
	CounselorIDs []uint `json:"counselor_ids" validate:"required"`
}

// ClassCounselorResponse represents the class counselor assignment response
type ClassCounselorResponse struct {
	ID          uint              `json:"id"`
	ClassID     uint              `json:"class_id"`
	ClassName   string            `json:"class_name"`
	CounselorID uint              `json:"counselor_id"`
	Counselor   *CounselorResponse `json:"counselor,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
}

// ==================== Stats DTOs ====================

// TodayAttendanceStats represents today's attendance statistics
type TodayAttendanceStats struct {
	Present int64 `json:"present"`
	Absent  int64 `json:"absent"`
	Late    int64 `json:"late"`
	Total   int64 `json:"total"`
}

// SchoolStatsResponse represents the school statistics for admin sekolah dashboard
type SchoolStatsResponse struct {
	TotalStudents   int64                `json:"totalStudents"`
	TotalClasses    int64                `json:"totalClasses"`
	TotalTeachers   int64                `json:"totalTeachers"`
	TotalParents    int64                `json:"totalParents"`
	TodayAttendance TodayAttendanceStats `json:"todayAttendance"`
}


// ==================== User DTOs ====================

// UserFilter represents filter options for listing users
type UserFilter struct {
	Name     string `query:"name"`
	Role     string `query:"role"`
	IsActive *bool  `query:"is_active"`
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
}

// DefaultUserFilter returns default filter values
func DefaultUserFilter() UserFilter {
	return UserFilter{
		Page:     1,
		PageSize: 20,
	}
}

// UserResponse represents the user data in responses
type UserResponse struct {
	ID                 uint                `json:"id"`
	SchoolID           uint                `json:"school_id"`
	Role               string              `json:"role"`
	Username           string              `json:"username"`
	Email              string              `json:"email,omitempty"`
	Name               string              `json:"name,omitempty"`
	IsActive           bool                `json:"is_active"`
	MustResetPwd       bool                `json:"must_reset_pwd"`
	AssignedClassID    *uint               `json:"assigned_class_id,omitempty"`    // For wali_kelas
	AssignedClassName  string              `json:"assigned_class_name,omitempty"`  // For wali_kelas
	AssignedClasses    []AssignedClassInfo `json:"assigned_classes,omitempty"`     // For guru_bk
	LastLoginAt        *string             `json:"last_login_at,omitempty"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
}

// AssignedClassInfo represents simplified class info for user assignments
type AssignedClassInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// UserListResponse represents a paginated list of users
type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	Pagination PaginationMeta `json:"pagination"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Role             string `json:"role" validate:"required,oneof=guru wali_kelas guru_bk admin_sekolah"`
	Username         string `json:"username" validate:"required"`
	Email            string `json:"email"`
	Name             string `json:"name"`
	Password         string `json:"password" validate:"required,min=8"`
	AssignedClassID  *uint  `json:"assigned_class_id"`   // For wali_kelas
	AssignedClassIDs []uint `json:"assigned_class_ids"`  // For guru_bk
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Email            *string `json:"email"`
	Name             *string `json:"name"`
	IsActive         *bool   `json:"is_active"`
	AssignedClassID  *uint   `json:"assigned_class_id"`   // For wali_kelas
	AssignedClassIDs []uint  `json:"assigned_class_ids"`  // For guru_bk
}

// ==================== Device DTOs ====================

// DeviceResponse represents the device data in responses
type DeviceResponse struct {
	ID          uint       `json:"id"`
	SchoolID    uint       `json:"school_id"`
	DeviceCode  string     `json:"device_code"`
	Description string     `json:"description"`
	IsActive    bool       `json:"is_active"`
	LastSeenAt  *time.Time `json:"last_seen_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
