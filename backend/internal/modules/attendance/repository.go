package attendance

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrAttendanceNotFound = errors.New("attendance record not found")
	ErrStudentNotFound    = errors.New("student not found")
	ErrInvalidRFIDCode    = errors.New("invalid RFID code")
	ErrDuplicateAttendance = errors.New("attendance record already exists for this date")
)

// Repository defines the interface for attendance data operations
type Repository interface {
	// Attendance CRUD
	Create(ctx context.Context, attendance *models.Attendance) error
	FindByID(ctx context.Context, id uint) (*models.Attendance, error)
	FindByStudentAndDate(ctx context.Context, studentID uint, date time.Time) (*models.Attendance, error)
	Update(ctx context.Context, attendance *models.Attendance) error
	Delete(ctx context.Context, id uint) error

	// Query operations
	FindByStudent(ctx context.Context, studentID uint, startDate, endDate time.Time) ([]models.Attendance, error)
	FindByClassAndDate(ctx context.Context, classID uint, date time.Time) ([]models.Attendance, error)
	FindBySchoolAndDate(ctx context.Context, schoolID uint, date time.Time) ([]models.Attendance, error)
	FindAll(ctx context.Context, schoolID uint, filter AttendanceFilter) ([]models.Attendance, int64, error)

	// Student lookup
	FindStudentByRFID(ctx context.Context, rfidCode string) (*models.Student, error)
	FindStudentByID(ctx context.Context, studentID uint) (*models.Student, error)
	FindStudentsByClass(ctx context.Context, classID uint) ([]models.Student, error)

	// Summary operations
	GetAttendanceSummary(ctx context.Context, schoolID uint, date time.Time) (*AttendanceSummaryResponse, error)
	GetClassAttendanceSummary(ctx context.Context, classID uint, date time.Time) (*ClassAttendanceSummaryResponse, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new attendance repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create creates a new attendance record
// Requirements: 5.2 - WHEN the backend receives attendance data, THE System SHALL record check-in or check-out
func (r *repository) Create(ctx context.Context, attendance *models.Attendance) error {
	return r.db.WithContext(ctx).Create(attendance).Error
}

// FindByID retrieves an attendance record by ID
func (r *repository) FindByID(ctx context.Context, id uint) (*models.Attendance, error) {
	var attendance models.Attendance
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Where("id = ?", id).
		First(&attendance).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAttendanceNotFound
		}
		return nil, err
	}

	return &attendance, nil
}

// FindByStudentAndDate retrieves attendance for a student on a specific date
// Requirements: 5.2 - Check-in or check-out based on existing records for that day
func (r *repository) FindByStudentAndDate(ctx context.Context, studentID uint, date time.Time) (*models.Attendance, error) {
	var attendance models.Attendance
	
	// Normalize date to start of day
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Where("student_id = ? AND date = ?", studentID, dateOnly).
		First(&attendance).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAttendanceNotFound
		}
		return nil, err
	}

	return &attendance, nil
}

// Update updates an attendance record
func (r *repository) Update(ctx context.Context, attendance *models.Attendance) error {
	result := r.db.WithContext(ctx).
		Model(&models.Attendance{}).
		Where("id = ?", attendance.ID).
		Updates(map[string]interface{}{
			"student_id":     attendance.StudentID,
			"date":           attendance.Date,
			"check_in_time":  attendance.CheckInTime,
			"check_out_time": attendance.CheckOutTime,
			"status":         attendance.Status,
			"method":         attendance.Method,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrAttendanceNotFound
	}
	return nil
}

// Delete deletes an attendance record
func (r *repository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Attendance{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrAttendanceNotFound
	}
	return nil
}


// FindByStudent retrieves attendance records for a student within a date range
func (r *repository) FindByStudent(ctx context.Context, studentID uint, startDate, endDate time.Time) ([]models.Attendance, error) {
	var attendances []models.Attendance
	
	query := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Where("student_id = ?", studentID)

	if !startDate.IsZero() {
		query = query.Where("date >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("date <= ?", endDate)
	}

	err := query.Order("date DESC").Find(&attendances).Error
	return attendances, err
}

// FindByClassAndDate retrieves attendance records for a class on a specific date
func (r *repository) FindByClassAndDate(ctx context.Context, classID uint, date time.Time) ([]models.Attendance, error) {
	var attendances []models.Attendance
	
	// Normalize date to start of day
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Joins("JOIN students ON students.id = attendances.student_id").
		Where("students.class_id = ? AND attendances.date = ?", classID, dateOnly).
		Order("students.name ASC").
		Find(&attendances).Error

	return attendances, err
}

// FindBySchoolAndDate retrieves attendance records for a school on a specific date
func (r *repository) FindBySchoolAndDate(ctx context.Context, schoolID uint, date time.Time) ([]models.Attendance, error) {
	var attendances []models.Attendance
	
	// Normalize date to start of day
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Joins("JOIN students ON students.id = attendances.student_id").
		Where("students.school_id = ? AND attendances.date = ?", schoolID, dateOnly).
		Order("students.name ASC").
		Find(&attendances).Error

	return attendances, err
}

// FindAll retrieves all attendance records with pagination and filtering
func (r *repository) FindAll(ctx context.Context, schoolID uint, filter AttendanceFilter) ([]models.Attendance, int64, error) {
	var attendances []models.Attendance
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.Attendance{}).
		Joins("JOIN students ON students.id = attendances.student_id").
		Where("students.school_id = ?", schoolID)

	// Apply filters
	if filter.StudentID != nil {
		query = query.Where("attendances.student_id = ?", *filter.StudentID)
	}
	if filter.ClassID != nil {
		query = query.Where("students.class_id = ?", *filter.ClassID)
	}
	if filter.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *filter.StartDate)
		if err == nil {
			query = query.Where("attendances.date >= ?", startDate)
		}
	}
	if filter.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *filter.EndDate)
		if err == nil {
			query = query.Where("attendances.date <= ?", endDate)
		}
	}
	if filter.Status != nil && *filter.Status != "" {
		query = query.Where("attendances.status = ?", *filter.Status)
	}
	if filter.Method != nil && *filter.Method != "" {
		query = query.Where("attendances.method = ?", *filter.Method)
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
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Joins("JOIN students ON students.id = attendances.student_id").
		Where("students.school_id = ?", schoolID).
		Order("attendances.date DESC, students.name ASC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&attendances).Error

	if err != nil {
		return nil, 0, err
	}

	return attendances, total, nil
}

// FindStudentByRFID retrieves a student by RFID code
// Requirements: 5.1 - WHEN a student taps RFID card, THE ESP32 SHALL send student identifier
func (r *repository) FindStudentByRFID(ctx context.Context, rfidCode string) (*models.Student, error) {
	var student models.Student
	err := r.db.WithContext(ctx).
		Preload("Class").
		Preload("School").
		Preload("Parents").
		Where("rfid_code = ? AND is_active = ?", rfidCode, true).
		First(&student).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidRFIDCode
		}
		return nil, err
	}

	return &student, nil
}

// FindStudentByID retrieves a student by ID
func (r *repository) FindStudentByID(ctx context.Context, studentID uint) (*models.Student, error) {
	var student models.Student
	err := r.db.WithContext(ctx).
		Preload("Class").
		Preload("School").
		Preload("Parents").
		Where("id = ?", studentID).
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
		Where("class_id = ? AND is_active = ?", classID, true).
		Order("name ASC").
		Find(&students).Error

	return students, err
}

// GetAttendanceSummary retrieves attendance summary for a school on a specific date
// Requirements: 5.4 - THE System SHALL display summary statistics (present, absent, late)
func (r *repository) GetAttendanceSummary(ctx context.Context, schoolID uint, date time.Time) (*AttendanceSummaryResponse, error) {
	// Normalize date to start of day
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	// Get total active students in school
	var totalStudents int64
	if err := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("school_id = ? AND is_active = ?", schoolID, true).
		Count(&totalStudents).Error; err != nil {
		return nil, err
	}

	// Get attendance counts by status
	type StatusCount struct {
		Status string
		Count  int64
	}
	var statusCounts []StatusCount

	err := r.db.WithContext(ctx).
		Model(&models.Attendance{}).
		Select("status, COUNT(*) as count").
		Joins("JOIN students ON students.id = attendances.student_id").
		Where("students.school_id = ? AND attendances.date = ?", schoolID, dateOnly).
		Group("status").
		Scan(&statusCounts).Error

	if err != nil {
		return nil, err
	}

	// Build summary
	summary := &AttendanceSummaryResponse{
		Date:       dateOnly.Format("2006-01-02"),
		TotalCount: int(totalStudents),
	}

	var presentCount int64
	for _, sc := range statusCounts {
		switch models.AttendanceStatus(sc.Status) {
		case models.AttendanceStatusOnTime:
			summary.Present += int(sc.Count)
			presentCount += sc.Count
		case models.AttendanceStatusLate:
			summary.Late = int(sc.Count)
			presentCount += sc.Count
		case models.AttendanceStatusVeryLate:
			summary.VeryLate = int(sc.Count)
			presentCount += sc.Count
		}
	}

	// Calculate absent (students without attendance record)
	summary.Absent = int(totalStudents - presentCount)

	return summary, nil
}

// GetClassAttendanceSummary retrieves attendance summary for a class on a specific date
func (r *repository) GetClassAttendanceSummary(ctx context.Context, classID uint, date time.Time) (*ClassAttendanceSummaryResponse, error) {
	// Normalize date to start of day
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	// Get class info
	var class models.Class
	if err := r.db.WithContext(ctx).First(&class, classID).Error; err != nil {
		return nil, err
	}

	// Get total active students in class
	var totalStudents int64
	if err := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("class_id = ? AND is_active = ?", classID, true).
		Count(&totalStudents).Error; err != nil {
		return nil, err
	}

	// Get attendance records for the class
	attendances, err := r.FindByClassAndDate(ctx, classID, date)
	if err != nil {
		return nil, err
	}

	// Build summary
	summary := AttendanceSummaryResponse{
		Date:       dateOnly.Format("2006-01-02"),
		TotalCount: int(totalStudents),
	}

	for _, a := range attendances {
		switch a.Status {
		case models.AttendanceStatusOnTime:
			summary.Present++
		case models.AttendanceStatusLate:
			summary.Late++
		case models.AttendanceStatusVeryLate:
			summary.VeryLate++
		}
	}

	// Calculate absent
	summary.Absent = int(totalStudents) - summary.Present - summary.Late - summary.VeryLate

	// Convert attendances to response format
	attendanceResponses := make([]AttendanceResponse, len(attendances))
	for i, a := range attendances {
		attendanceResponses[i] = *toAttendanceResponse(&a)
	}

	return &ClassAttendanceSummaryResponse{
		ClassID:     classID,
		ClassName:   class.Name,
		Date:        dateOnly.Format("2006-01-02"),
		Summary:     summary,
		Attendances: attendanceResponses,
	}, nil
}

// toAttendanceResponse converts an Attendance model to AttendanceResponse DTO
func toAttendanceResponse(attendance *models.Attendance) *AttendanceResponse {
	response := &AttendanceResponse{
		ID:        attendance.ID,
		StudentID: attendance.StudentID,
		Date:      attendance.Date.Format("2006-01-02"),
		Status:    attendance.Status,
		Method:    attendance.Method,
		CreatedAt: attendance.CreatedAt,
		UpdatedAt: attendance.UpdatedAt,
	}

	if attendance.CheckInTime != nil {
		checkIn := attendance.CheckInTime.Format("15:04")
		response.CheckInTime = &checkIn
	}
	if attendance.CheckOutTime != nil {
		checkOut := attendance.CheckOutTime.Format("15:04")
		response.CheckOutTime = &checkOut
	}

	// Include student info if loaded
	if attendance.Student.ID != 0 {
		response.StudentName = attendance.Student.Name
		response.StudentNIS = attendance.Student.NIS
		response.StudentNISN = attendance.Student.NISN
		if attendance.Student.Class.ID != 0 {
			response.ClassName = attendance.Student.Class.Name
		}
	}

	return response
}
