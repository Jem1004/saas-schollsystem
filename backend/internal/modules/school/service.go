package school

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrNameRequired        = errors.New("nama wajib diisi")
	ErrNISRequired         = errors.New("NIS wajib diisi")
	ErrNISNRequired        = errors.New("NISN wajib diisi")
	ErrClassIDRequired     = errors.New("kelas wajib dipilih")
	ErrGradeRequired       = errors.New("tingkat wajib diisi")
	ErrYearRequired        = errors.New("tahun ajaran wajib diisi")
	ErrPasswordRequired    = errors.New("password wajib diisi")
	ErrPasswordTooShort    = errors.New("password minimal 8 karakter")
	ErrStudentIDsRequired  = errors.New("minimal satu siswa wajib dipilih")
	ErrClassHasStudents    = errors.New("tidak dapat menghapus kelas yang masih memiliki siswa")
	ErrPhoneRequired       = errors.New("nomor HP wajib diisi")
	ErrStudentHasAccount   = errors.New("siswa sudah memiliki akun")
	ErrStudentNoAccount    = errors.New("siswa belum memiliki akun")
)

// Default password for parent and student accounts
const DefaultPassword = "password123"

// Service defines the interface for school business logic
type Service interface {
	// Class operations
	CreateClass(ctx context.Context, schoolID uint, req CreateClassRequest) (*ClassResponse, error)
	GetAllClasses(ctx context.Context, schoolID uint, filter ClassFilter) (*ClassListResponse, error)
	GetClassByID(ctx context.Context, schoolID uint, id uint) (*ClassResponse, error)
	UpdateClass(ctx context.Context, schoolID uint, id uint, req UpdateClassRequest) (*ClassResponse, error)
	DeleteClass(ctx context.Context, schoolID uint, id uint) error
	AssignHomeroomTeacher(ctx context.Context, schoolID uint, classID uint, teacherID uint) (*ClassResponse, error)

	// Student operations
	CreateStudent(ctx context.Context, schoolID uint, req CreateStudentRequest) (*StudentResponse, error)
	GetAllStudents(ctx context.Context, schoolID uint, filter StudentFilter) (*StudentListResponse, error)
	GetStudentByID(ctx context.Context, schoolID uint, id uint) (*StudentResponse, error)
	GetStudentByNISN(ctx context.Context, nisn string) (*StudentResponse, error)
	GetStudentsByClass(ctx context.Context, schoolID uint, classID uint) ([]StudentResponse, error)
	UpdateStudent(ctx context.Context, schoolID uint, id uint, req UpdateStudentRequest) (*StudentResponse, error)
	DeleteStudent(ctx context.Context, schoolID uint, id uint) error
	ClearStudentRFID(ctx context.Context, schoolID uint, studentID uint) error

	// Parent operations
	CreateParent(ctx context.Context, schoolID uint, req CreateParentRequest) (*ParentResponse, error)
	GetAllParents(ctx context.Context, schoolID uint, filter ParentFilter) (*ParentListResponse, error)
	GetParentByID(ctx context.Context, schoolID uint, id uint) (*ParentResponse, error)
	UpdateParent(ctx context.Context, schoolID uint, id uint, req UpdateParentRequest) (*ParentResponse, error)
	DeleteParent(ctx context.Context, schoolID uint, id uint) error
	LinkParentToStudents(ctx context.Context, schoolID uint, parentID uint, studentIDs []uint) (*ParentResponse, error)
	ResetParentPassword(ctx context.Context, schoolID uint, parentID uint) (*ResetPasswordResponse, error)

	// Student operations with account
	CreateStudentAccount(ctx context.Context, schoolID uint, studentID uint) (*StudentResponse, error)
	ResetStudentPassword(ctx context.Context, schoolID uint, studentID uint) (*ResetPasswordResponse, error)

	// Stats operations
	GetStats(ctx context.Context, schoolID uint) (*SchoolStatsResponse, error)

	// User operations
	GetAllUsers(ctx context.Context, schoolID uint, filter UserFilter) (*UserListResponse, error)
	GetUserByID(ctx context.Context, schoolID uint, id uint) (*UserResponse, error)
	CreateUser(ctx context.Context, schoolID uint, req CreateUserRequest) (*UserResponse, error)
	UpdateUser(ctx context.Context, schoolID uint, id uint, req UpdateUserRequest) (*UserResponse, error)
	DeleteUser(ctx context.Context, schoolID uint, id uint) error
	ResetUserPassword(ctx context.Context, schoolID uint, id uint) (*ResetPasswordResponse, error)

	// Device operations
	GetSchoolDevices(ctx context.Context, schoolID uint) ([]DeviceResponse, error)

	// Bulk operations for import
	GetStudentsWithoutClass(ctx context.Context, schoolID uint) ([]StudentResponse, error)
	BulkAssignClass(ctx context.Context, schoolID uint, req BulkAssignClassRequest) (*BulkAssignClassResponse, error)

	// Search operations for parent linking
	SearchStudents(ctx context.Context, schoolID uint, query string) ([]StudentSearchResponse, error)
}

// service implements the Service interface
type service struct {
	repo     Repository
	userRepo UserRepository
}

// UserRepository defines the interface for user operations needed by school service
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByUsername(ctx context.Context, username string) (*models.User, error)
}

// NewService creates a new school service
func NewService(repo Repository, userRepo UserRepository) Service {
	return &service{
		repo:     repo,
		userRepo: userRepo,
	}
}


// ==================== Class Service Methods ====================

// CreateClass creates a new class
// Requirements: 3.1 - WHEN an Admin_Sekolah creates a class, THE System SHALL associate it with the school tenant and allow student assignment
func (s *service) CreateClass(ctx context.Context, schoolID uint, req CreateClassRequest) (*ClassResponse, error) {
	// Validate required fields
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrNameRequired
	}
	if req.Grade <= 0 {
		return nil, ErrGradeRequired
	}
	year := strings.TrimSpace(req.Year)
	if year == "" {
		return nil, ErrYearRequired
	}

	// Check for duplicate class
	existing, err := s.repo.FindClassByNameGradeYear(ctx, schoolID, name, req.Grade, year)
	if err != nil && !errors.Is(err, ErrClassNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrDuplicateClass
	}

	// Validate homeroom teacher if provided
	if req.HomeroomTeacherID != nil {
		if err := s.repo.ValidateHomeroomTeacher(ctx, schoolID, *req.HomeroomTeacherID); err != nil {
			return nil, err
		}
	}

	// Create class
	class := &models.Class{
		SchoolID:          schoolID,
		Name:              name,
		Grade:             req.Grade,
		Year:              year,
		HomeroomTeacherID: req.HomeroomTeacherID,
	}

	if err := s.repo.CreateClass(ctx, class); err != nil {
		return nil, err
	}

	// Reload with relations
	class, err = s.repo.FindClassByID(ctx, schoolID, class.ID)
	if err != nil {
		return nil, err
	}

	return s.toClassResponse(ctx, class), nil
}

// GetAllClasses retrieves all classes for a school
func (s *service) GetAllClasses(ctx context.Context, schoolID uint, filter ClassFilter) (*ClassListResponse, error) {
	// Set defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	classes, total, err := s.repo.FindAllClasses(ctx, schoolID, filter)
	if err != nil {
		return nil, err
	}

	// Convert to response
	classResponses := make([]ClassResponse, len(classes))
	for i, class := range classes {
		classResponses[i] = *s.toClassResponse(ctx, &class)
	}

	// Calculate total pages
	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &ClassListResponse{
		Classes: classResponses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetClassByID retrieves a class by ID
func (s *service) GetClassByID(ctx context.Context, schoolID uint, id uint) (*ClassResponse, error) {
	class, err := s.repo.FindClassByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}
	return s.toClassResponse(ctx, class), nil
}

// UpdateClass updates a class
func (s *service) UpdateClass(ctx context.Context, schoolID uint, id uint, req UpdateClassRequest) (*ClassResponse, error) {
	// Get existing class
	class, err := s.repo.FindClassByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, ErrNameRequired
		}
		class.Name = name
	}
	if req.Grade != nil {
		if *req.Grade <= 0 {
			return nil, ErrGradeRequired
		}
		class.Grade = *req.Grade
	}
	if req.Year != nil {
		year := strings.TrimSpace(*req.Year)
		if year == "" {
			return nil, ErrYearRequired
		}
		class.Year = year
	}

	// Check for duplicate if name, grade, or year changed
	existing, err := s.repo.FindClassByNameGradeYear(ctx, schoolID, class.Name, class.Grade, class.Year)
	if err != nil && !errors.Is(err, ErrClassNotFound) {
		return nil, err
	}
	if existing != nil && existing.ID != class.ID {
		return nil, ErrDuplicateClass
	}

	// Validate and update homeroom teacher if provided
	if req.HomeroomTeacherID != nil {
		if *req.HomeroomTeacherID == 0 {
			class.HomeroomTeacherID = nil
		} else {
			if err := s.repo.ValidateHomeroomTeacher(ctx, schoolID, *req.HomeroomTeacherID); err != nil {
				return nil, err
			}
			class.HomeroomTeacherID = req.HomeroomTeacherID
		}
	}

	if err := s.repo.UpdateClass(ctx, class); err != nil {
		return nil, err
	}

	// Reload with relations
	class, err = s.repo.FindClassByID(ctx, schoolID, class.ID)
	if err != nil {
		return nil, err
	}

	return s.toClassResponse(ctx, class), nil
}

// DeleteClass deletes a class
func (s *service) DeleteClass(ctx context.Context, schoolID uint, id uint) error {
	// Check if class has students
	count, err := s.repo.GetClassStudentCount(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrClassHasStudents
	}

	return s.repo.DeleteClass(ctx, schoolID, id)
}

// AssignHomeroomTeacher assigns a homeroom teacher to a class
// Requirements: 4.3 - WHEN an Admin_Sekolah assigns Wali_Kelas role, THE System SHALL require class assignment
func (s *service) AssignHomeroomTeacher(ctx context.Context, schoolID uint, classID uint, teacherID uint) (*ClassResponse, error) {
	// Validate teacher
	if err := s.repo.ValidateHomeroomTeacher(ctx, schoolID, teacherID); err != nil {
		return nil, err
	}

	// Get class
	class, err := s.repo.FindClassByID(ctx, schoolID, classID)
	if err != nil {
		return nil, err
	}

	// Assign teacher
	class.HomeroomTeacherID = &teacherID

	if err := s.repo.UpdateClass(ctx, class); err != nil {
		return nil, err
	}

	// Reload with relations
	class, err = s.repo.FindClassByID(ctx, schoolID, class.ID)
	if err != nil {
		return nil, err
	}

	return s.toClassResponse(ctx, class), nil
}

// toClassResponse converts a Class model to ClassResponse DTO
func (s *service) toClassResponse(ctx context.Context, class *models.Class) *ClassResponse {
	response := &ClassResponse{
		ID:                class.ID,
		SchoolID:          class.SchoolID,
		Name:              class.Name,
		Grade:             class.Grade,
		Year:              class.Year,
		HomeroomTeacherID: class.HomeroomTeacherID,
		CreatedAt:         class.CreatedAt,
		UpdatedAt:         class.UpdatedAt,
	}

	// Get student count
	count, _ := s.repo.GetClassStudentCount(ctx, class.ID)
	response.StudentCount = count

	// Add homeroom teacher info if available
	if class.HomeroomTeacher != nil {
		response.HomeroomTeacher = &TeacherResponse{
			ID:       class.HomeroomTeacher.ID,
			Name:     class.HomeroomTeacher.Name,
			Username: class.HomeroomTeacher.Username,
		}
	}

	return response
}


// ==================== Student Service Methods ====================

// CreateStudent creates a new student
// Requirements: 3.2 - WHEN an Admin_Sekolah registers a student, THE System SHALL require NIS, NISN, name, and class assignment
// Requirements: 3.5 - IF duplicate NISN is detected within the system, THEN THE System SHALL reject the registration with an error message
// If CreateAccount is true, also creates a user account with NIS as username
func (s *service) CreateStudent(ctx context.Context, schoolID uint, req CreateStudentRequest) (*StudentResponse, error) {
	// Validate required fields
	nis := strings.TrimSpace(req.NIS)
	if nis == "" {
		return nil, ErrNISRequired
	}
	nisn := strings.TrimSpace(req.NISN)
	if nisn == "" {
		return nil, ErrNISNRequired
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrNameRequired
	}
	if req.ClassID == 0 {
		return nil, ErrClassIDRequired
	}

	// Validate class exists
	_, err := s.repo.FindClassByID(ctx, schoolID, req.ClassID)
	if err != nil {
		return nil, err
	}

	// Check for duplicate NISN (globally unique)
	existing, err := s.repo.FindStudentByNISN(ctx, nisn)
	if err != nil && !errors.Is(err, ErrStudentNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrDuplicateNISN
	}

	// Check for duplicate NIS within school
	existing, err = s.repo.FindStudentByNIS(ctx, schoolID, nis)
	if err != nil && !errors.Is(err, ErrStudentNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrDuplicateNIS
	}

	// Create user account if requested
	var userID *uint
	var tempPassword string
	if req.CreateAccount {
		// Check if username already exists
		_, err := s.userRepo.FindByUsername(ctx, nis)
		if err == nil {
			return nil, errors.New("NIS sudah digunakan sebagai username siswa lain")
		}

		// Use default password
		tempPassword = DefaultPassword

		// Hash password
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		// Create user account
		user := &models.User{
			SchoolID:     &schoolID,
			Role:         models.RoleStudent,
			Username:     nis,
			PasswordHash: string(passwordHash),
			Name:         name,
			IsActive:     true,
			MustResetPwd: true,
		}

		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
		userID = &user.ID
	}

	// Create student with ClassID as pointer
	classID := req.ClassID
	student := &models.Student{
		SchoolID: schoolID,
		ClassID:  &classID,
		NIS:      nis,
		NISN:     nisn,
		Name:     name,
		RFIDCode: strings.TrimSpace(req.RFIDCode),
		IsActive: true, // IsActive is true when ClassID is set
		UserID:   userID,
	}

	if err := s.repo.CreateStudent(ctx, student); err != nil {
		return nil, err
	}

	// Reload with relations
	student, err = s.repo.FindStudentByID(ctx, schoolID, student.ID)
	if err != nil {
		return nil, err
	}

	response := s.toStudentResponse(student)
	if req.CreateAccount {
		response.Username = nis
		response.TemporaryPassword = tempPassword
	}

	return response, nil
}

// GetAllStudents retrieves all students for a school
func (s *service) GetAllStudents(ctx context.Context, schoolID uint, filter StudentFilter) (*StudentListResponse, error) {
	// Set defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	students, total, err := s.repo.FindAllStudents(ctx, schoolID, filter)
	if err != nil {
		return nil, err
	}

	// Convert to response
	studentResponses := make([]StudentResponse, len(students))
	for i, student := range students {
		studentResponses[i] = *s.toStudentResponse(&student)
	}

	// Calculate total pages
	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &StudentListResponse{
		Students: studentResponses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetStudentByID retrieves a student by ID
func (s *service) GetStudentByID(ctx context.Context, schoolID uint, id uint) (*StudentResponse, error) {
	student, err := s.repo.FindStudentByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}
	return s.toStudentResponse(student), nil
}

// GetStudentByNISN retrieves a student by NISN
func (s *service) GetStudentByNISN(ctx context.Context, nisn string) (*StudentResponse, error) {
	student, err := s.repo.FindStudentByNISN(ctx, nisn)
	if err != nil {
		return nil, err
	}
	return s.toStudentResponse(student), nil
}

// GetStudentsByClass retrieves all students in a class
func (s *service) GetStudentsByClass(ctx context.Context, schoolID uint, classID uint) ([]StudentResponse, error) {
	// Validate class exists
	_, err := s.repo.FindClassByID(ctx, schoolID, classID)
	if err != nil {
		return nil, err
	}

	students, err := s.repo.FindStudentsByClass(ctx, classID)
	if err != nil {
		return nil, err
	}

	// Convert to response
	responses := make([]StudentResponse, len(students))
	for i, student := range students {
		responses[i] = *s.toStudentResponse(&student)
	}

	return responses, nil
}

// UpdateStudent updates a student
// Requirements: 3.9, 3.10, 8.2, 8.3 - Handle ClassID changes and IsActive based on ClassID presence
func (s *service) UpdateStudent(ctx context.Context, schoolID uint, id uint, req UpdateStudentRequest) (*StudentResponse, error) {
	// Get existing student
	student, err := s.repo.FindStudentByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.NIS != nil {
		nis := strings.TrimSpace(*req.NIS)
		if nis == "" {
			return nil, ErrNISRequired
		}
		// Check for duplicate NIS if changed
		if nis != student.NIS {
			existing, err := s.repo.FindStudentByNIS(ctx, schoolID, nis)
			if err != nil && !errors.Is(err, ErrStudentNotFound) {
				return nil, err
			}
			if existing != nil && existing.ID != student.ID {
				return nil, ErrDuplicateNIS
			}
		}
		student.NIS = nis
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, ErrNameRequired
		}
		student.Name = name
	}
	if req.ClassID != nil {
		if *req.ClassID <= 0 {
			return nil, ErrClassIDRequired
		}
		classID := uint(*req.ClassID)
		// Validate class exists
		_, err := s.repo.FindClassByID(ctx, schoolID, classID)
		if err != nil {
			return nil, err
		}
		student.ClassID = &classID
		// When ClassID is set, student can be active (Requirements: 8.3)
		student.IsActive = true
	}
	if req.RFIDCode != nil {
		student.RFIDCode = strings.TrimSpace(*req.RFIDCode)
	}
	if req.IsActive != nil {
		// Only allow setting IsActive to true if ClassID is set (Requirements: 8.2)
		if *req.IsActive && !student.CanBeActive() {
			// Cannot set IsActive to true without ClassID
			student.IsActive = false
		} else {
			student.IsActive = *req.IsActive
		}
	}

	if err := s.repo.UpdateStudent(ctx, student); err != nil {
		return nil, err
	}

	// Reload with relations
	student, err = s.repo.FindStudentByID(ctx, schoolID, student.ID)
	if err != nil {
		return nil, err
	}

	return s.toStudentResponse(student), nil
}

// DeleteStudent deletes a student
func (s *service) DeleteStudent(ctx context.Context, schoolID uint, id uint) error {
	return s.repo.DeleteStudent(ctx, schoolID, id)
}

// toStudentResponse converts a Student model to StudentResponse DTO
func (s *service) toStudentResponse(student *models.Student) *StudentResponse {
	response := &StudentResponse{
		ID:         student.ID,
		SchoolID:   student.SchoolID,
		ClassID:    student.ClassID,
		NIS:        student.NIS,
		NISN:       student.NISN,
		Name:       student.Name,
		RFIDCode:   student.RFIDCode,
		IsActive:   student.IsActive,
		HasAccount: student.UserID != nil,
		CreatedAt:  student.CreatedAt,
		UpdatedAt:  student.UpdatedAt,
	}

	// Add username if student has account
	if student.UserID != nil {
		response.Username = student.NIS // Username is NIS
	}

	// Add class info if available (Class is now a pointer)
	if student.Class != nil && student.Class.ID != 0 {
		response.ClassName = student.Class.Name
		response.Class = &ClassResponse{
			ID:       student.Class.ID,
			SchoolID: student.Class.SchoolID,
			Name:     student.Class.Name,
			Grade:    student.Class.Grade,
			Year:     student.Class.Year,
		}
	}

	return response
}


// ==================== Parent Service Methods ====================

// CreateParent creates a new parent with a user account
// Requirements: 3.3 - WHEN an Admin_Sekolah registers a parent, THE System SHALL link the parent to one or more students
// Username: Phone number (primary), Password: Auto-generated or manual input
func (s *service) CreateParent(ctx context.Context, schoolID uint, req CreateParentRequest) (*ParentResponse, error) {
	// Validate required fields
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrNameRequired
	}
	
	// Phone is required for parent (used as primary username)
	phone := strings.TrimSpace(req.Phone)
	if phone == "" {
		return nil, ErrPhoneRequired
	}
	
	email := strings.TrimSpace(req.Email)

	// Username is phone number
	username := phone

	// Check if username (phone) already exists
	_, err := s.userRepo.FindByUsername(ctx, username)
	if err == nil {
		return nil, errors.New("nomor HP sudah digunakan oleh orang tua lain")
	}

	// Use default password if not provided
	password := strings.TrimSpace(req.Password)
	if password == "" {
		password = DefaultPassword
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user account for parent
	user := &models.User{
		SchoolID:     &schoolID,
		Role:         models.RoleParent,
		Username:     username,
		PasswordHash: string(passwordHash),
		Email:        email,
		Name:         name,
		IsActive:     true,
		MustResetPwd: true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Create parent
	parent := &models.Parent{
		SchoolID: schoolID,
		UserID:   user.ID,
		Name:     name,
		Phone:    phone,
	}

	if err := s.repo.CreateParent(ctx, parent); err != nil {
		return nil, err
	}

	// Link to students if provided
	if len(req.StudentIDs) > 0 {
		if err := s.repo.LinkParentToStudents(ctx, parent.ID, req.StudentIDs); err != nil {
			return nil, err
		}
	}

	// Reload with relations
	parent, err = s.repo.FindParentByID(ctx, schoolID, parent.ID)
	if err != nil {
		return nil, err
	}

	response := s.toParentResponse(parent)
	response.Username = username
	response.TemporaryPassword = password
	
	return response, nil
}

// GetAllParents retrieves all parents for a school
func (s *service) GetAllParents(ctx context.Context, schoolID uint, filter ParentFilter) (*ParentListResponse, error) {
	// Set defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	parents, total, err := s.repo.FindAllParents(ctx, schoolID, filter)
	if err != nil {
		return nil, err
	}

	// Convert to response
	parentResponses := make([]ParentResponse, len(parents))
	for i, parent := range parents {
		parentResponses[i] = *s.toParentResponse(&parent)
	}

	// Calculate total pages
	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &ParentListResponse{
		Parents: parentResponses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetParentByID retrieves a parent by ID
func (s *service) GetParentByID(ctx context.Context, schoolID uint, id uint) (*ParentResponse, error) {
	parent, err := s.repo.FindParentByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}
	return s.toParentResponse(parent), nil
}

// UpdateParent updates a parent
func (s *service) UpdateParent(ctx context.Context, schoolID uint, id uint, req UpdateParentRequest) (*ParentResponse, error) {
	// Get existing parent
	parent, err := s.repo.FindParentByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, ErrNameRequired
		}
		parent.Name = name
	}
	if req.Phone != nil {
		parent.Phone = strings.TrimSpace(*req.Phone)
	}

	if err := s.repo.UpdateParent(ctx, parent); err != nil {
		return nil, err
	}

	// Update email in user table if provided
	if req.Email != nil {
		email := strings.TrimSpace(*req.Email)
		if err := s.repo.UpdateParentUserEmail(ctx, parent.UserID, email); err != nil {
			return nil, err
		}
	}

	// Always update student links (frontend always sends student_ids)
	// Validate all students exist and belong to the same school
	for _, studentID := range req.StudentIDs {
		_, err := s.repo.FindStudentByID(ctx, schoolID, studentID)
		if err != nil {
			return nil, err
		}
	}
	if err := s.repo.LinkParentToStudents(ctx, parent.ID, req.StudentIDs); err != nil {
		return nil, err
	}

	// Reload with relations
	parent, err = s.repo.FindParentByID(ctx, schoolID, parent.ID)
	if err != nil {
		return nil, err
	}

	return s.toParentResponse(parent), nil
}

// DeleteParent deletes a parent
func (s *service) DeleteParent(ctx context.Context, schoolID uint, id uint) error {
	return s.repo.DeleteParent(ctx, schoolID, id)
}

// ResetParentPassword resets a parent's password to default password
func (s *service) ResetParentPassword(ctx context.Context, schoolID uint, parentID uint) (*ResetPasswordResponse, error) {
	// Get parent
	parent, err := s.repo.FindParentByID(ctx, schoolID, parentID)
	if err != nil {
		return nil, err
	}

	// Use default password
	newPassword := DefaultPassword

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Update user password
	if err := s.repo.ResetUserPassword(ctx, parent.UserID, string(passwordHash)); err != nil {
		return nil, err
	}

	return &ResetPasswordResponse{
		Username:          parent.Phone,
		TemporaryPassword: newPassword,
		Message:           "Password berhasil direset ke default. User wajib mengganti password saat login pertama.",
	}, nil
}

// CreateStudentAccount creates a user account for a student
func (s *service) CreateStudentAccount(ctx context.Context, schoolID uint, studentID uint) (*StudentResponse, error) {
	// Get student
	student, err := s.repo.FindStudentByID(ctx, schoolID, studentID)
	if err != nil {
		return nil, err
	}

	// Check if student already has an account
	if student.UserID != nil {
		return nil, ErrStudentHasAccount
	}

	// Username is NIS
	username := student.NIS

	// Check if username already exists
	_, err = s.userRepo.FindByUsername(ctx, username)
	if err == nil {
		return nil, errors.New("NIS sudah digunakan sebagai username siswa lain")
	}

	// Use default password
	password := DefaultPassword

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user account for student
	user := &models.User{
		SchoolID:     &schoolID,
		Role:         models.RoleStudent,
		Username:     username,
		PasswordHash: string(passwordHash),
		Name:         student.Name,
		IsActive:     true,
		MustResetPwd: true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Link user to student
	if err := s.repo.UpdateStudentUserID(ctx, studentID, user.ID); err != nil {
		return nil, err
	}

	// Reload student
	student, err = s.repo.FindStudentByID(ctx, schoolID, studentID)
	if err != nil {
		return nil, err
	}

	response := s.toStudentResponse(student)
	response.Username = username
	response.TemporaryPassword = password

	return response, nil
}

// ResetStudentPassword resets a student's password to a new generated password
func (s *service) ResetStudentPassword(ctx context.Context, schoolID uint, studentID uint) (*ResetPasswordResponse, error) {
	// Get student
	student, err := s.repo.FindStudentByID(ctx, schoolID, studentID)
	if err != nil {
		return nil, err
	}

	// Check if student has an account
	if student.UserID == nil {
		return nil, ErrStudentNoAccount
	}

	// Use default password
	newPassword := DefaultPassword

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Update user password
	if err := s.repo.ResetUserPassword(ctx, *student.UserID, string(passwordHash)); err != nil {
		return nil, err
	}

	return &ResetPasswordResponse{
		Username:          student.NIS,
		TemporaryPassword: newPassword,
		Message:           "Password berhasil direset. User wajib mengganti password saat login pertama.",
	}, nil
}

// LinkParentToStudents links a parent to students
// Requirements: 3.3 - WHEN an Admin_Sekolah registers a parent, THE System SHALL link the parent to one or more students
func (s *service) LinkParentToStudents(ctx context.Context, schoolID uint, parentID uint, studentIDs []uint) (*ParentResponse, error) {
	if len(studentIDs) == 0 {
		return nil, ErrStudentIDsRequired
	}

	// Validate parent exists
	parent, err := s.repo.FindParentByID(ctx, schoolID, parentID)
	if err != nil {
		return nil, err
	}

	// Validate all students exist and belong to the same school
	for _, studentID := range studentIDs {
		_, err := s.repo.FindStudentByID(ctx, schoolID, studentID)
		if err != nil {
			return nil, err
		}
	}

	// Link parent to students
	if err := s.repo.LinkParentToStudents(ctx, parentID, studentIDs); err != nil {
		return nil, err
	}

	// Reload with relations
	parent, err = s.repo.FindParentByID(ctx, schoolID, parent.ID)
	if err != nil {
		return nil, err
	}

	return s.toParentResponse(parent), nil
}

// toParentResponse converts a Parent model to ParentResponse DTO
func (s *service) toParentResponse(parent *models.Parent) *ParentResponse {
	response := &ParentResponse{
		ID:        parent.ID,
		SchoolID:  parent.SchoolID,
		UserID:    parent.UserID,
		Name:      parent.Name,
		Phone:     parent.Phone,
		CreatedAt: parent.CreatedAt,
		UpdatedAt: parent.UpdatedAt,
	}

	// Add email from user if available
	if parent.User.ID != 0 {
		response.Email = parent.User.Email
	}

	// Add students if available
	if len(parent.Students) > 0 {
		response.Students = make([]StudentResponse, len(parent.Students))
		for i, student := range parent.Students {
			response.Students[i] = *s.toStudentResponse(&student)
		}
	}

	return response
}

// ==================== Stats Service Methods ====================

// GetStats retrieves statistics for a school (for admin sekolah dashboard)
func (s *service) GetStats(ctx context.Context, schoolID uint) (*SchoolStatsResponse, error) {
	return s.repo.GetSchoolStats(ctx, schoolID)
}


// ==================== User Service Methods ====================

// GetAllUsers retrieves all users for a school
func (s *service) GetAllUsers(ctx context.Context, schoolID uint, filter UserFilter) (*UserListResponse, error) {
	// Set defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	users, total, err := s.repo.FindAllUsers(ctx, schoolID, filter)
	if err != nil {
		return nil, err
	}

	// Convert to response with assigned class info
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		response := s.toUserResponse(&user)
		
		// Get assigned class for wali_kelas
		if user.Role == models.RoleWaliKelas {
			class, err := s.repo.FindClassByHomeroomTeacher(ctx, schoolID, user.ID)
			if err == nil && class != nil {
				response.AssignedClassID = &class.ID
				response.AssignedClassName = class.Name
			}
		}

		// Get assigned classes for guru_bk
		if user.Role == models.RoleGuruBK {
			classes, err := s.repo.FindClassesByCounselor(ctx, schoolID, user.ID)
			if err == nil && len(classes) > 0 {
				response.AssignedClasses = make([]AssignedClassInfo, len(classes))
				for j, class := range classes {
					response.AssignedClasses[j] = AssignedClassInfo{
						ID:   class.ID,
						Name: class.Name,
					}
				}
			}
		}
		
		userResponses[i] = *response
	}

	// Calculate total pages
	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &UserListResponse{
		Users: userResponses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *service) GetUserByID(ctx context.Context, schoolID uint, id uint) (*UserResponse, error) {
	user, err := s.repo.FindUserByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}
	
	response := s.toUserResponse(user)
	
	// Get assigned class for wali_kelas
	if user.Role == models.RoleWaliKelas {
		class, err := s.repo.FindClassByHomeroomTeacher(ctx, *user.SchoolID, user.ID)
		if err == nil && class != nil {
			response.AssignedClassID = &class.ID
			response.AssignedClassName = class.Name
		}
	}

	// Get assigned classes for guru_bk
	if user.Role == models.RoleGuruBK {
		classes, err := s.repo.FindClassesByCounselor(ctx, *user.SchoolID, user.ID)
		if err == nil && len(classes) > 0 {
			response.AssignedClasses = make([]AssignedClassInfo, len(classes))
			for i, class := range classes {
				response.AssignedClasses[i] = AssignedClassInfo{
					ID:   class.ID,
					Name: class.Name,
				}
			}
		}
	}
	
	return response, nil
}

// CreateUser creates a new user
func (s *service) CreateUser(ctx context.Context, schoolID uint, req CreateUserRequest) (*UserResponse, error) {
	// Validate required fields
	username := strings.TrimSpace(req.Username)
	if username == "" {
		return nil, errors.New("Username wajib diisi")
	}
	password := strings.TrimSpace(req.Password)
	if password == "" {
		return nil, ErrPasswordRequired
	}
	if len(password) < 8 {
		return nil, ErrPasswordTooShort
	}

	// Check if username already exists
	existingUser, err := s.userRepo.FindByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return nil, errors.New("Username sudah digunakan")
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		SchoolID:     &schoolID,
		Role:         models.UserRole(req.Role),
		Username:     username,
		PasswordHash: string(passwordHash),
		Email:        strings.TrimSpace(req.Email),
		Name:         strings.TrimSpace(req.Name),
		IsActive:     true,
		MustResetPwd: true,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// If wali_kelas and assigned_class_id is provided, update the class
	if req.Role == "wali_kelas" && req.AssignedClassID != nil {
		class, err := s.repo.FindClassByID(ctx, schoolID, *req.AssignedClassID)
		if err == nil && class != nil {
			class.HomeroomTeacherID = &user.ID
			if err := s.repo.UpdateClass(ctx, class); err != nil {
				// Log error but don't fail the user creation
				// The class assignment can be done later
			}
		}
	}

	// If guru_bk and assigned_class_ids is provided, assign to classes
	if req.Role == "guru_bk" && len(req.AssignedClassIDs) > 0 {
		for _, classID := range req.AssignedClassIDs {
			s.repo.AssignCounselorToClass(ctx, schoolID, classID, user.ID)
		}
	}

	// Reload user
	user, err = s.repo.FindUserByID(ctx, schoolID, user.ID)
	if err != nil {
		return nil, err
	}

	response := s.toUserResponse(user)
	
	// Get assigned class for wali_kelas
	if user.Role == models.RoleWaliKelas {
		class, err := s.repo.FindClassByHomeroomTeacher(ctx, schoolID, user.ID)
		if err == nil && class != nil {
			response.AssignedClassID = &class.ID
			response.AssignedClassName = class.Name
		}
	}

	// Get assigned classes for guru_bk
	if user.Role == models.RoleGuruBK {
		classes, err := s.repo.FindClassesByCounselor(ctx, schoolID, user.ID)
		if err == nil && len(classes) > 0 {
			response.AssignedClasses = make([]AssignedClassInfo, len(classes))
			for i, class := range classes {
				response.AssignedClasses[i] = AssignedClassInfo{
					ID:   class.ID,
					Name: class.Name,
				}
			}
		}
	}

	return response, nil
}

// UpdateUser updates a user
func (s *service) UpdateUser(ctx context.Context, schoolID uint, id uint, req UpdateUserRequest) (*UserResponse, error) {
	// Get existing user
	user, err := s.repo.FindUserByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Email != nil {
		user.Email = strings.TrimSpace(*req.Email)
	}
	if req.Name != nil {
		user.Name = strings.TrimSpace(*req.Name)
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	// Handle assigned_class_id for wali_kelas
	if user.Role == models.RoleWaliKelas && req.AssignedClassID != nil {
		// First, remove this teacher from any existing class
		existingClass, _ := s.repo.FindClassByHomeroomTeacher(ctx, schoolID, user.ID)
		if existingClass != nil && existingClass.ID != *req.AssignedClassID {
			existingClass.HomeroomTeacherID = nil
			s.repo.UpdateClass(ctx, existingClass)
		}

		// Then assign to the new class
		if *req.AssignedClassID > 0 {
			newClass, err := s.repo.FindClassByID(ctx, schoolID, *req.AssignedClassID)
			if err == nil && newClass != nil {
				newClass.HomeroomTeacherID = &user.ID
				s.repo.UpdateClass(ctx, newClass)
			}
		}
	}

	// Handle assigned_class_ids for guru_bk
	if user.Role == models.RoleGuruBK && req.AssignedClassIDs != nil {
		// Remove from all existing classes first
		s.repo.RemoveCounselorFromAllClasses(ctx, schoolID, user.ID)

		// Then assign to new classes
		for _, classID := range req.AssignedClassIDs {
			s.repo.AssignCounselorToClass(ctx, schoolID, classID, user.ID)
		}
	}

	// Reload user
	user, err = s.repo.FindUserByID(ctx, schoolID, user.ID)
	if err != nil {
		return nil, err
	}

	response := s.toUserResponse(user)
	
	// Get assigned class for wali_kelas
	if user.Role == models.RoleWaliKelas {
		class, err := s.repo.FindClassByHomeroomTeacher(ctx, schoolID, user.ID)
		if err == nil && class != nil {
			response.AssignedClassID = &class.ID
			response.AssignedClassName = class.Name
		}
	}

	// Get assigned classes for guru_bk
	if user.Role == models.RoleGuruBK {
		classes, err := s.repo.FindClassesByCounselor(ctx, schoolID, user.ID)
		if err == nil && len(classes) > 0 {
			response.AssignedClasses = make([]AssignedClassInfo, len(classes))
			for i, class := range classes {
				response.AssignedClasses[i] = AssignedClassInfo{
					ID:   class.ID,
					Name: class.Name,
				}
			}
		}
	}

	return response, nil
}

// DeleteUser deletes a user
func (s *service) DeleteUser(ctx context.Context, schoolID uint, id uint) error {
	// Get user first to check role
	user, err := s.repo.FindUserByID(ctx, schoolID, id)
	if err != nil {
		return err
	}

	// If wali_kelas, remove from assigned class first
	if user.Role == models.RoleWaliKelas {
		class, _ := s.repo.FindClassByHomeroomTeacher(ctx, schoolID, id)
		if class != nil {
			class.HomeroomTeacherID = nil
			s.repo.UpdateClass(ctx, class)
		}
	}

	// If guru_bk, remove from all assigned classes
	if user.Role == models.RoleGuruBK {
		s.repo.RemoveCounselorFromAllClasses(ctx, schoolID, id)
	}

	return s.repo.DeleteUser(ctx, schoolID, id)
}

// ResetUserPassword resets a user's password
func (s *service) ResetUserPassword(ctx context.Context, schoolID uint, id uint) (*ResetPasswordResponse, error) {
	// Get existing user
	user, err := s.repo.FindUserByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	// Generate new password
	newPassword := DefaultPassword

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Update password
	user.PasswordHash = string(passwordHash)
	user.MustResetPwd = true

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return &ResetPasswordResponse{
		Username:          user.Username,
		TemporaryPassword: newPassword,
		Message:           "Password berhasil direset",
	}, nil
}

// toUserResponse converts a User model to UserResponse DTO
func (s *service) toUserResponse(user *models.User) *UserResponse {
	response := &UserResponse{
		ID:           user.ID,
		SchoolID:     *user.SchoolID,
		Role:         string(user.Role),
		Username:     user.Username,
		Email:        user.Email,
		Name:         user.Name,
		IsActive:     user.IsActive,
		MustResetPwd: user.MustResetPwd,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	if user.LastLoginAt != nil {
		lastLogin := user.LastLoginAt.Format("2006-01-02T15:04:05Z07:00")
		response.LastLoginAt = &lastLogin
	}

	return response
}

// ==================== Device Service Methods ====================

// GetSchoolDevices retrieves all devices for a school
func (s *service) GetSchoolDevices(ctx context.Context, schoolID uint) ([]DeviceResponse, error) {
	devices, err := s.repo.FindDevicesBySchool(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	responses := make([]DeviceResponse, len(devices))
	for i, device := range devices {
		responses[i] = DeviceResponse{
			ID:          device.ID,
			SchoolID:    device.SchoolID,
			DeviceCode:  device.DeviceCode,
			Description: device.Description,
			IsActive:    device.IsActive,
			LastSeenAt:  device.LastSeenAt,
			CreatedAt:   device.CreatedAt,
			UpdatedAt:   device.UpdatedAt,
		}
	}

	return responses, nil
}

// ClearStudentRFID clears the RFID code from a student
func (s *service) ClearStudentRFID(ctx context.Context, schoolID uint, studentID uint) error {
	// Verify student exists and belongs to school
	_, err := s.repo.FindStudentByID(ctx, schoolID, studentID)
	if err != nil {
		return err
	}

	return s.repo.ClearStudentRFID(ctx, studentID)
}


// ==================== Bulk Operations Service Methods ====================

// GetStudentsWithoutClass retrieves all students without class assignment
// Requirements: 6.1 - Filter for students without class assignment
func (s *service) GetStudentsWithoutClass(ctx context.Context, schoolID uint) ([]StudentResponse, error) {
	students, err := s.repo.FindStudentsWithoutClass(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	responses := make([]StudentResponse, len(students))
	for i, student := range students {
		responses[i] = *s.toStudentResponse(&student)
	}

	return responses, nil
}

// BulkAssignClass assigns a class to multiple students
// Requirements: 6.2, 6.3, 6.4, 6.5 - Bulk class assignment with validation
func (s *service) BulkAssignClass(ctx context.Context, schoolID uint, req BulkAssignClassRequest) (*BulkAssignClassResponse, error) {
	// Validate student IDs are provided
	if len(req.StudentIDs) == 0 {
		return nil, ErrStudentIDsRequired
	}

	// Validate class exists and belongs to the same school
	class, err := s.repo.FindClassByID(ctx, schoolID, req.ClassID)
	if err != nil {
		return nil, err
	}
	if class.SchoolID != schoolID {
		return nil, ErrClassNotFound
	}

	// Validate all students exist and belong to the same school
	for _, studentID := range req.StudentIDs {
		_, err := s.repo.FindStudentByID(ctx, schoolID, studentID)
		if err != nil {
			return nil, err
		}
	}

	// Perform bulk update
	if err := s.repo.BulkUpdateStudentClass(ctx, req.StudentIDs, req.ClassID); err != nil {
		return nil, err
	}

	// Reload updated students
	updatedStudents := make([]StudentResponse, 0, len(req.StudentIDs))
	for _, studentID := range req.StudentIDs {
		student, err := s.repo.FindStudentByID(ctx, schoolID, studentID)
		if err == nil {
			updatedStudents = append(updatedStudents, *s.toStudentResponse(student))
		}
	}

	return &BulkAssignClassResponse{
		UpdatedCount: len(updatedStudents),
		Students:     updatedStudents,
		Message:      "Kelas berhasil ditetapkan untuk siswa terpilih",
	}, nil
}


// ==================== Search Operations Service Methods ====================

// SearchStudents searches students by NISN or name for parent linking
// Requirements: 7.2 - Search students by NISN or name for parent linking
func (s *service) SearchStudents(ctx context.Context, schoolID uint, query string) ([]StudentSearchResponse, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []StudentSearchResponse{}, nil
	}

	students, err := s.repo.SearchStudents(ctx, schoolID, query, 20)
	if err != nil {
		return nil, err
	}

	responses := make([]StudentSearchResponse, len(students))
	for i, student := range students {
		responses[i] = StudentSearchResponse{
			ID:      student.ID,
			NIS:     student.NIS,
			NISN:    student.NISN,
			Name:    student.Name,
			ClassID: student.ClassID,
		}
		if student.Class != nil {
			responses[i].ClassName = student.Class.Name
		}
	}

	return responses, nil
}
