package school

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrClassNotFound    = errors.New("kelas tidak ditemukan")
	ErrStudentNotFound  = errors.New("siswa tidak ditemukan")
	ErrParentNotFound   = errors.New("orang tua tidak ditemukan")
	ErrDuplicateNISN    = errors.New("siswa dengan NISN ini sudah terdaftar")
	ErrDuplicateNIS     = errors.New("siswa dengan NIS ini sudah terdaftar di sekolah ini")
	ErrDuplicateClass   = errors.New("kelas dengan nama ini sudah ada untuk tingkat dan tahun yang sama")
	ErrTeacherNotFound  = errors.New("guru tidak ditemukan")
	ErrInvalidTeacher   = errors.New("user bukan guru yang valid untuk ditugaskan sebagai wali kelas")
)

// Repository defines the interface for school data operations
type Repository interface {
	// Class operations
	CreateClass(ctx context.Context, class *models.Class) error
	FindAllClasses(ctx context.Context, schoolID uint, filter ClassFilter) ([]models.Class, int64, error)
	FindClassByID(ctx context.Context, schoolID uint, id uint) (*models.Class, error)
	FindClassByNameGradeYear(ctx context.Context, schoolID uint, name string, grade int, year string) (*models.Class, error)
	UpdateClass(ctx context.Context, class *models.Class) error
	DeleteClass(ctx context.Context, schoolID uint, id uint) error
	GetClassStudentCount(ctx context.Context, classID uint) (int64, error)

	// Student operations
	CreateStudent(ctx context.Context, student *models.Student) error
	FindAllStudents(ctx context.Context, schoolID uint, filter StudentFilter) ([]models.Student, int64, error)
	FindStudentByID(ctx context.Context, schoolID uint, id uint) (*models.Student, error)
	FindStudentByNISN(ctx context.Context, nisn string) (*models.Student, error)
	FindStudentByNIS(ctx context.Context, schoolID uint, nis string) (*models.Student, error)
	FindStudentsByClass(ctx context.Context, classID uint) ([]models.Student, error)
	UpdateStudent(ctx context.Context, student *models.Student) error
	UpdateStudentUserID(ctx context.Context, studentID uint, userID uint) error
	DeleteStudent(ctx context.Context, schoolID uint, id uint) error

	// Parent operations
	CreateParent(ctx context.Context, parent *models.Parent) error
	FindAllParents(ctx context.Context, schoolID uint, filter ParentFilter) ([]models.Parent, int64, error)
	FindParentByID(ctx context.Context, schoolID uint, id uint) (*models.Parent, error)
	FindParentByUserID(ctx context.Context, userID uint) (*models.Parent, error)
	UpdateParent(ctx context.Context, parent *models.Parent) error
	UpdateParentUserEmail(ctx context.Context, userID uint, email string) error
	ResetUserPassword(ctx context.Context, userID uint, passwordHash string) error
	DeleteParent(ctx context.Context, schoolID uint, id uint) error
	LinkParentToStudents(ctx context.Context, parentID uint, studentIDs []uint) error
	UnlinkParentFromStudent(ctx context.Context, parentID uint, studentID uint) error
	GetParentStudents(ctx context.Context, parentID uint) ([]models.Student, error)

	// Teacher operations
	FindTeacherByID(ctx context.Context, schoolID uint, teacherID uint) (*models.User, error)
	ValidateHomeroomTeacher(ctx context.Context, schoolID uint, teacherID uint) error

	// Stats operations
	GetSchoolStats(ctx context.Context, schoolID uint) (*SchoolStatsResponse, error)

	// User operations (for school staff management)
	FindAllUsers(ctx context.Context, schoolID uint, filter UserFilter) ([]models.User, int64, error)
	FindUserByID(ctx context.Context, schoolID uint, id uint) (*models.User, error)
	FindClassByHomeroomTeacher(ctx context.Context, schoolID uint, teacherID uint) (*models.Class, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, schoolID uint, id uint) error

	// Class Counselor operations
	FindClassCounselorsByClass(ctx context.Context, schoolID uint, classID uint) ([]models.ClassCounselor, error)
	FindClassesByCounselor(ctx context.Context, schoolID uint, counselorID uint) ([]models.Class, error)
	AssignCounselorToClass(ctx context.Context, schoolID uint, classID uint, counselorID uint) error
	RemoveCounselorFromClass(ctx context.Context, schoolID uint, classID uint, counselorID uint) error
	RemoveAllCounselorsFromClass(ctx context.Context, schoolID uint, classID uint) error
	RemoveCounselorFromAllClasses(ctx context.Context, schoolID uint, counselorID uint) error

	// Device operations
	FindDevicesBySchool(ctx context.Context, schoolID uint) ([]models.Device, error)

	// RFID operations
	ClearStudentRFID(ctx context.Context, studentID uint) error

	// Bulk operations for import
	FindStudentsWithoutClass(ctx context.Context, schoolID uint) ([]models.Student, error)
	BulkUpdateStudentClass(ctx context.Context, studentIDs []uint, classID uint) error

	// Search operations for parent linking
	SearchStudents(ctx context.Context, schoolID uint, query string, limit int) ([]models.Student, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new school repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}


// ==================== Class Repository Methods ====================

// CreateClass creates a new class
// Requirements: 3.1 - WHEN an Admin_Sekolah creates a class, THE System SHALL associate it with the school tenant
func (r *repository) CreateClass(ctx context.Context, class *models.Class) error {
	return r.db.WithContext(ctx).Create(class).Error
}

// FindAllClasses retrieves all classes for a school with pagination and filtering
func (r *repository) FindAllClasses(ctx context.Context, schoolID uint, filter ClassFilter) ([]models.Class, int64, error) {
	var classes []models.Class
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Class{}).Where("school_id = ?", schoolID)

	// Apply filters
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.Grade != nil {
		query = query.Where("grade = ?", *filter.Grade)
	}
	if filter.Year != "" {
		query = query.Where("year = ?", filter.Year)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (filter.Page - 1) * filter.PageSize
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	// Fetch records with homeroom teacher
	err := query.
		Preload("HomeroomTeacher").
		Order("grade ASC, name ASC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&classes).Error

	if err != nil {
		return nil, 0, err
	}

	return classes, total, nil
}

// FindClassByID retrieves a class by ID within a school
func (r *repository) FindClassByID(ctx context.Context, schoolID uint, id uint) (*models.Class, error) {
	var class models.Class
	err := r.db.WithContext(ctx).
		Preload("HomeroomTeacher").
		Where("id = ? AND school_id = ?", id, schoolID).
		First(&class).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrClassNotFound
		}
		return nil, err
	}

	return &class, nil
}

// FindClassByNameGradeYear finds a class by name, grade, and year within a school
func (r *repository) FindClassByNameGradeYear(ctx context.Context, schoolID uint, name string, grade int, year string) (*models.Class, error) {
	var class models.Class
	err := r.db.WithContext(ctx).
		Where("school_id = ? AND name = ? AND grade = ? AND year = ?", schoolID, name, grade, year).
		First(&class).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrClassNotFound
		}
		return nil, err
	}

	return &class, nil
}

// UpdateClass updates a class
func (r *repository) UpdateClass(ctx context.Context, class *models.Class) error {
	result := r.db.WithContext(ctx).
		Model(&models.Class{}).
		Where("id = ?", class.ID).
		Updates(map[string]interface{}{
			"name":                class.Name,
			"grade":               class.Grade,
			"year":                class.Year,
			"homeroom_teacher_id": class.HomeroomTeacherID,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrClassNotFound
	}
	return nil
}

// DeleteClass deletes a class
func (r *repository) DeleteClass(ctx context.Context, schoolID uint, id uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", id, schoolID).
		Delete(&models.Class{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrClassNotFound
	}
	return nil
}

// GetClassStudentCount returns the number of students in a class
func (r *repository) GetClassStudentCount(ctx context.Context, classID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("class_id = ?", classID).
		Count(&count).Error
	return count, err
}


// ==================== Student Repository Methods ====================

// CreateStudent creates a new student
// Requirements: 3.2 - WHEN an Admin_Sekolah registers a student, THE System SHALL require NIS, NISN, name, and class assignment
func (r *repository) CreateStudent(ctx context.Context, student *models.Student) error {
	return r.db.WithContext(ctx).Create(student).Error
}

// FindAllStudents retrieves all students for a school with pagination and filtering
func (r *repository) FindAllStudents(ctx context.Context, schoolID uint, filter StudentFilter) ([]models.Student, int64, error) {
	var students []models.Student
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Student{}).Where("school_id = ?", schoolID)

	// Apply filters
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.NIS != "" {
		query = query.Where("nis ILIKE ?", "%"+filter.NIS+"%")
	}
	if filter.NISN != "" {
		query = query.Where("nisn ILIKE ?", "%"+filter.NISN+"%")
	}
	if filter.ClassID != nil {
		query = query.Where("class_id = ?", *filter.ClassID)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (filter.Page - 1) * filter.PageSize
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	// Fetch records with class
	err := query.
		Preload("Class").
		Order("name ASC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&students).Error

	if err != nil {
		return nil, 0, err
	}

	return students, total, nil
}

// FindStudentByID retrieves a student by ID within a school
func (r *repository) FindStudentByID(ctx context.Context, schoolID uint, id uint) (*models.Student, error) {
	var student models.Student
	err := r.db.WithContext(ctx).
		Preload("Class").
		Preload("Parents").
		Where("id = ? AND school_id = ?", id, schoolID).
		First(&student).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStudentNotFound
		}
		return nil, err
	}

	return &student, nil
}

// FindStudentByNISN retrieves a student by NISN (globally unique)
// Requirements: 3.5 - IF duplicate NISN is detected within the system, THEN THE System SHALL reject the registration
func (r *repository) FindStudentByNISN(ctx context.Context, nisn string) (*models.Student, error) {
	var student models.Student
	err := r.db.WithContext(ctx).
		Where("nisn = ?", nisn).
		First(&student).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStudentNotFound
		}
		return nil, err
	}

	return &student, nil
}

// FindStudentByNIS retrieves a student by NIS within a school
func (r *repository) FindStudentByNIS(ctx context.Context, schoolID uint, nis string) (*models.Student, error) {
	var student models.Student
	err := r.db.WithContext(ctx).
		Where("school_id = ? AND nis = ?", schoolID, nis).
		First(&student).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStudentNotFound
		}
		return nil, err
	}

	return &student, nil
}

// FindStudentsByClass retrieves all students in a class
func (r *repository) FindStudentsByClass(ctx context.Context, classID uint) ([]models.Student, error) {
	var students []models.Student
	err := r.db.WithContext(ctx).
		Where("class_id = ?", classID).
		Order("name ASC").
		Find(&students).Error

	return students, err
}

// UpdateStudent updates a student
func (r *repository) UpdateStudent(ctx context.Context, student *models.Student) error {
	result := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("id = ?", student.ID).
		Updates(map[string]interface{}{
			"class_id":   student.ClassID,
			"nis":        student.NIS,
			"nisn":       student.NISN,
			"name":       student.Name,
			"rf_id_code": student.RFIDCode,
			"is_active":  student.IsActive,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrStudentNotFound
	}
	return nil
}

// UpdateStudentUserID links a user account to a student
func (r *repository) UpdateStudentUserID(ctx context.Context, studentID uint, userID uint) error {
	result := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("id = ?", studentID).
		Update("user_id", userID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteStudent deletes a student
func (r *repository) DeleteStudent(ctx context.Context, schoolID uint, id uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", id, schoolID).
		Delete(&models.Student{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrStudentNotFound
	}
	return nil
}


// ==================== Parent Repository Methods ====================

// CreateParent creates a new parent
// Requirements: 3.3 - WHEN an Admin_Sekolah registers a parent, THE System SHALL link the parent to one or more students
func (r *repository) CreateParent(ctx context.Context, parent *models.Parent) error {
	return r.db.WithContext(ctx).Create(parent).Error
}

// FindAllParents retrieves all parents for a school with pagination and filtering
func (r *repository) FindAllParents(ctx context.Context, schoolID uint, filter ParentFilter) ([]models.Parent, int64, error) {
	var parents []models.Parent
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Parent{}).Where("school_id = ?", schoolID)

	// Apply filters
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.Phone != "" {
		query = query.Where("phone ILIKE ?", "%"+filter.Phone+"%")
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (filter.Page - 1) * filter.PageSize
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	// Fetch records with students and user
	err := query.
		Preload("User").
		Preload("Students").
		Order("name ASC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&parents).Error

	if err != nil {
		return nil, 0, err
	}

	return parents, total, nil
}

// FindParentByID retrieves a parent by ID within a school
func (r *repository) FindParentByID(ctx context.Context, schoolID uint, id uint) (*models.Parent, error) {
	var parent models.Parent
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Students").
		Preload("Students.Class").
		Where("id = ? AND school_id = ?", id, schoolID).
		First(&parent).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrParentNotFound
		}
		return nil, err
	}

	return &parent, nil
}

// FindParentByUserID retrieves a parent by user ID
func (r *repository) FindParentByUserID(ctx context.Context, userID uint) (*models.Parent, error) {
	var parent models.Parent
	err := r.db.WithContext(ctx).
		Preload("Students").
		Where("user_id = ?", userID).
		First(&parent).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrParentNotFound
		}
		return nil, err
	}

	return &parent, nil
}

// UpdateParent updates a parent
func (r *repository) UpdateParent(ctx context.Context, parent *models.Parent) error {
	result := r.db.WithContext(ctx).
		Model(&models.Parent{}).
		Where("id = ?", parent.ID).
		Updates(map[string]interface{}{
			"name":  parent.Name,
			"phone": parent.Phone,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrParentNotFound
	}
	return nil
}

// UpdateParentUserEmail updates the email of the user associated with a parent
func (r *repository) UpdateParentUserEmail(ctx context.Context, userID uint, email string) error {
	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("email", email)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// ResetUserPassword resets a user's password and sets must_reset_pwd to true
func (r *repository) ResetUserPassword(ctx context.Context, userID uint, passwordHash string) error {
	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"password_hash":  passwordHash,
			"must_reset_pwd": true,
		})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteParent deletes a parent
func (r *repository) DeleteParent(ctx context.Context, schoolID uint, id uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", id, schoolID).
		Delete(&models.Parent{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrParentNotFound
	}
	return nil
}

// LinkParentToStudents links a parent to multiple students
// Requirements: 3.3 - WHEN an Admin_Sekolah registers a parent, THE System SHALL link the parent to one or more students
func (r *repository) LinkParentToStudents(ctx context.Context, parentID uint, studentIDs []uint) error {
	// Get the parent
	var parent models.Parent
	if err := r.db.WithContext(ctx).First(&parent, parentID).Error; err != nil {
		return ErrParentNotFound
	}

	// Get the students (handle empty array case)
	var students []models.Student
	if len(studentIDs) > 0 {
		if err := r.db.WithContext(ctx).Where("id IN ?", studentIDs).Find(&students).Error; err != nil {
			return err
		}
	}

	// Replace associations (empty students array will clear all associations)
	return r.db.WithContext(ctx).Model(&parent).Association("Students").Replace(students)
}

// UnlinkParentFromStudent unlinks a parent from a student
func (r *repository) UnlinkParentFromStudent(ctx context.Context, parentID uint, studentID uint) error {
	var parent models.Parent
	if err := r.db.WithContext(ctx).First(&parent, parentID).Error; err != nil {
		return ErrParentNotFound
	}

	var student models.Student
	if err := r.db.WithContext(ctx).First(&student, studentID).Error; err != nil {
		return ErrStudentNotFound
	}

	return r.db.WithContext(ctx).Model(&parent).Association("Students").Delete(&student)
}

// GetParentStudents retrieves all students linked to a parent
func (r *repository) GetParentStudents(ctx context.Context, parentID uint) ([]models.Student, error) {
	var parent models.Parent
	if err := r.db.WithContext(ctx).Preload("Students").Preload("Students.Class").First(&parent, parentID).Error; err != nil {
		return nil, ErrParentNotFound
	}
	return parent.Students, nil
}

// ==================== Teacher Repository Methods ====================

// FindTeacherByID retrieves a teacher by ID within a school
func (r *repository) FindTeacherByID(ctx context.Context, schoolID uint, teacherID uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", teacherID, schoolID).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTeacherNotFound
		}
		return nil, err
	}

	return &user, nil
}

// ValidateHomeroomTeacher validates that a user can be assigned as homeroom teacher
// Requirements: 4.3 - WHEN an Admin_Sekolah assigns Wali_Kelas role, THE System SHALL require class assignment
func (r *repository) ValidateHomeroomTeacher(ctx context.Context, schoolID uint, teacherID uint) error {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", teacherID, schoolID).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTeacherNotFound
		}
		return err
	}

	// Check if user has a valid role for homeroom teacher
	validRoles := []models.UserRole{
		models.RoleWaliKelas,
		models.RoleGuru,
		models.RoleAdminSekolah,
	}

	isValid := false
	for _, role := range validRoles {
		if user.Role == role {
			isValid = true
			break
		}
	}

	if !isValid {
		return ErrInvalidTeacher
	}

	return nil
}

// ==================== Stats Repository Methods ====================

// GetSchoolStats retrieves statistics for a school (for admin sekolah dashboard)
func (r *repository) GetSchoolStats(ctx context.Context, schoolID uint) (*SchoolStatsResponse, error) {
	stats := &SchoolStatsResponse{}

	// Count total students
	if err := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("school_id = ?", schoolID).
		Count(&stats.TotalStudents).Error; err != nil {
		return nil, err
	}

	// Count total classes
	if err := r.db.WithContext(ctx).
		Model(&models.Class{}).
		Where("school_id = ?", schoolID).
		Count(&stats.TotalClasses).Error; err != nil {
		return nil, err
	}

	// Count total teachers (users with role guru, wali_kelas, or guru_bk)
	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("school_id = ? AND role IN ?", schoolID, []models.UserRole{
			models.RoleGuru,
			models.RoleWaliKelas,
			models.RoleGuruBK,
		}).
		Count(&stats.TotalTeachers).Error; err != nil {
		return nil, err
	}

	// Count total parents
	if err := r.db.WithContext(ctx).
		Model(&models.Parent{}).
		Where("school_id = ?", schoolID).
		Count(&stats.TotalParents).Error; err != nil {
		return nil, err
	}

	// Get today's attendance stats
	today := r.db.NowFunc().Format("2006-01-02")

	// Total students for attendance
	stats.TodayAttendance.Total = stats.TotalStudents

	// Count present (status = 'present' or 'late')
	var presentCount int64
	if err := r.db.WithContext(ctx).
		Model(&models.Attendance{}).
		Joins("JOIN students ON students.id = attendances.student_id").
		Where("students.school_id = ? AND DATE(attendances.check_in_time) = ?", schoolID, today).
		Where("attendances.status IN ?", []string{"present", "late"}).
		Count(&presentCount).Error; err != nil {
		return nil, err
	}
	stats.TodayAttendance.Present = presentCount

	// Count late
	var lateCount int64
	if err := r.db.WithContext(ctx).
		Model(&models.Attendance{}).
		Joins("JOIN students ON students.id = attendances.student_id").
		Where("students.school_id = ? AND DATE(attendances.check_in_time) = ?", schoolID, today).
		Where("attendances.status = ?", "late").
		Count(&lateCount).Error; err != nil {
		return nil, err
	}
	stats.TodayAttendance.Late = lateCount

	// Absent = Total - Present (those who haven't checked in)
	stats.TodayAttendance.Absent = stats.TotalStudents - presentCount

	return stats, nil
}

// ==================== User Repository Methods ====================

// FindAllUsers retrieves all users for a school with pagination and filtering
func (r *repository) FindAllUsers(ctx context.Context, schoolID uint, filter UserFilter) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.WithContext(ctx).Model(&models.User{}).Where("school_id = ?", schoolID)

	// Exclude parent and student roles (they have their own management)
	query = query.Where("role NOT IN ?", []models.UserRole{models.RoleParent, models.RoleStudent, models.RoleSuperAdmin})

	// Apply filters
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (filter.Page - 1) * filter.PageSize
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	// Fetch records
	err := query.
		Order("name ASC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// FindUserByID retrieves a user by ID within a school
func (r *repository) FindUserByID(ctx context.Context, schoolID uint, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", id, schoolID).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// CreateUser creates a new user
func (r *repository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// UpdateUser updates a user
func (r *repository) UpdateUser(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"email":         user.Email,
			"name":          user.Name,
			"is_active":     user.IsActive,
			"password_hash": user.PasswordHash,
			"must_reset_pwd": user.MustResetPwd,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

// DeleteUser deletes a user
func (r *repository) DeleteUser(ctx context.Context, schoolID uint, id uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", id, schoolID).
		Delete(&models.User{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

// FindClassByHomeroomTeacher finds the class assigned to a homeroom teacher
func (r *repository) FindClassByHomeroomTeacher(ctx context.Context, schoolID uint, teacherID uint) (*models.Class, error) {
	var class models.Class
	err := r.db.WithContext(ctx).
		Where("school_id = ? AND homeroom_teacher_id = ?", schoolID, teacherID).
		First(&class).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No class assigned, not an error
		}
		return nil, err
	}

	return &class, nil
}

// ==================== Class Counselor Repository Methods ====================

// FindClassCounselorsByClass retrieves all counselors assigned to a class
func (r *repository) FindClassCounselorsByClass(ctx context.Context, schoolID uint, classID uint) ([]models.ClassCounselor, error) {
	var counselors []models.ClassCounselor
	err := r.db.WithContext(ctx).
		Preload("Counselor").
		Where("school_id = ? AND class_id = ?", schoolID, classID).
		Find(&counselors).Error

	if err != nil {
		return nil, err
	}

	return counselors, nil
}

// FindClassesByCounselor retrieves all classes assigned to a counselor
func (r *repository) FindClassesByCounselor(ctx context.Context, schoolID uint, counselorID uint) ([]models.Class, error) {
	var classes []models.Class
	err := r.db.WithContext(ctx).
		Joins("JOIN class_counselors ON class_counselors.class_id = classes.id").
		Where("class_counselors.school_id = ? AND class_counselors.counselor_id = ?", schoolID, counselorID).
		Find(&classes).Error

	if err != nil {
		return nil, err
	}

	return classes, nil
}

// AssignCounselorToClass assigns a counselor to a class
func (r *repository) AssignCounselorToClass(ctx context.Context, schoolID uint, classID uint, counselorID uint) error {
	// Check if already assigned
	var existing models.ClassCounselor
	err := r.db.WithContext(ctx).
		Where("school_id = ? AND class_id = ? AND counselor_id = ?", schoolID, classID, counselorID).
		First(&existing).Error

	if err == nil {
		// Already assigned
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Create new assignment
	assignment := &models.ClassCounselor{
		SchoolID:    schoolID,
		ClassID:     classID,
		CounselorID: counselorID,
	}

	return r.db.WithContext(ctx).Create(assignment).Error
}

// RemoveCounselorFromClass removes a counselor from a class
func (r *repository) RemoveCounselorFromClass(ctx context.Context, schoolID uint, classID uint, counselorID uint) error {
	return r.db.WithContext(ctx).
		Where("school_id = ? AND class_id = ? AND counselor_id = ?", schoolID, classID, counselorID).
		Delete(&models.ClassCounselor{}).Error
}

// RemoveAllCounselorsFromClass removes all counselors from a class
func (r *repository) RemoveAllCounselorsFromClass(ctx context.Context, schoolID uint, classID uint) error {
	return r.db.WithContext(ctx).
		Where("school_id = ? AND class_id = ?", schoolID, classID).
		Delete(&models.ClassCounselor{}).Error
}

// RemoveCounselorFromAllClasses removes a counselor from all classes
func (r *repository) RemoveCounselorFromAllClasses(ctx context.Context, schoolID uint, counselorID uint) error {
	return r.db.WithContext(ctx).
		Where("school_id = ? AND counselor_id = ?", schoolID, counselorID).
		Delete(&models.ClassCounselor{}).Error
}

// ==================== Device Repository Methods ====================

// FindDevicesBySchool retrieves all devices for a school
func (r *repository) FindDevicesBySchool(ctx context.Context, schoolID uint) ([]models.Device, error) {
	var devices []models.Device
	err := r.db.WithContext(ctx).
		Where("school_id = ? AND is_active = ?", schoolID, true).
		Order("device_code ASC").
		Find(&devices).Error

	if err != nil {
		return nil, err
	}

	return devices, nil
}

// ==================== RFID Repository Methods ====================

// ClearStudentRFID clears the RFID code from a student
func (r *repository) ClearStudentRFID(ctx context.Context, studentID uint) error {
	result := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("id = ?", studentID).
		Update("rf_id_code", "")

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrStudentNotFound
	}
	return nil
}


// ==================== Bulk Operations Repository Methods ====================

// FindStudentsWithoutClass retrieves all students without class assignment
// Requirements: 6.1 - Filter for students without class assignment
func (r *repository) FindStudentsWithoutClass(ctx context.Context, schoolID uint) ([]models.Student, error) {
	var students []models.Student
	err := r.db.WithContext(ctx).
		Where("school_id = ? AND class_id IS NULL", schoolID).
		Order("name ASC").
		Find(&students).Error

	if err != nil {
		return nil, err
	}

	return students, nil
}

// BulkUpdateStudentClass updates ClassID and IsActive for multiple students
// Requirements: 6.3, 6.4 - Bulk class assignment with IsActive update
func (r *repository) BulkUpdateStudentClass(ctx context.Context, studentIDs []uint, classID uint) error {
	if len(studentIDs) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("id IN ?", studentIDs).
		Updates(map[string]interface{}{
			"class_id":  classID,
			"is_active": true,
		})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// SearchStudents searches students by NISN or name within a school
// Requirements: 7.2 - Search students by NISN or name for parent linking
func (r *repository) SearchStudents(ctx context.Context, schoolID uint, query string, limit int) ([]models.Student, error) {
	var students []models.Student

	if limit <= 0 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}

	searchQuery := "%" + query + "%"

	err := r.db.WithContext(ctx).
		Preload("Class").
		Where("school_id = ?", schoolID).
		Where("nisn ILIKE ? OR name ILIKE ?", searchQuery, searchQuery).
		Order("name ASC").
		Limit(limit).
		Find(&students).Error

	if err != nil {
		return nil, err
	}

	return students, nil
}
