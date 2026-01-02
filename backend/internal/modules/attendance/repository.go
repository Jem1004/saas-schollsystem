package attendance

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrAttendanceNotFound = errors.New("data kehadiran tidak ditemukan")
	ErrStudentNotFound    = errors.New("siswa tidak ditemukan")
	ErrInvalidRFIDCode    = errors.New("kode RFID tidak valid")
	ErrDuplicateAttendance = errors.New("data kehadiran untuk tanggal ini sudah ada")
)

// Repository defines the interface for attendance data operations
type Repository interface {
	// Attendance CRUD
	Create(ctx context.Context, attendance *models.Attendance) error
	FindByID(ctx context.Context, id uint) (*models.Attendance, error)
	FindByStudentAndDate(ctx context.Context, studentID uint, date time.Time) (*models.Attendance, error)
	FindByStudentDateAndSchedule(ctx context.Context, studentID uint, date time.Time, scheduleID uint) (*models.Attendance, error)
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
	GetAttendanceSummaryByClass(ctx context.Context, schoolID uint, date time.Time) (*AttendanceSummaryResponse, []ClassSummaryItem, error)
	GetClassAttendanceSummary(ctx context.Context, classID uint, date time.Time) (*ClassAttendanceSummaryResponse, error)

	// Schedule operations
	// Requirements: 3.4, 3.5 - Find and associate active schedule with attendance
	FindActiveSchedule(ctx context.Context, schoolID uint, timestamp time.Time) (*models.AttendanceSchedule, error)
	FindDefaultSchedule(ctx context.Context, schoolID uint) (*models.AttendanceSchedule, error)

	// Export operations
	// Requirements: 1.2, 1.3 - Get attendance for export with filters
	GetAttendanceForExport(ctx context.Context, schoolID uint, filter ExportFilter) ([]ExportAttendanceRecord, error)
	// Requirements: 2.1 - Get monthly recap with aggregation
	GetMonthlyRecap(ctx context.Context, schoolID uint, filter MonthlyRecapFilter) (*MonthlyRecapResponse, error)
	// Get students by school for monthly recap calculation
	GetStudentsBySchool(ctx context.Context, schoolID uint, classID *uint) ([]models.Student, error)

	// Wali Kelas operations
	// Requirements: 2.7 - Find class assigned to wali_kelas
	FindClassByHomeroomTeacher(ctx context.Context, schoolID uint, teacherID uint) (*models.Class, error)
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
		Preload("Schedule").
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
		Preload("Schedule").
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

// FindByStudentDateAndSchedule retrieves attendance for a student on a specific date and schedule
// Used to check if student already has attendance for a specific schedule
func (r *repository) FindByStudentDateAndSchedule(ctx context.Context, studentID uint, date time.Time, scheduleID uint) (*models.Attendance, error) {
	var attendance models.Attendance
	
	// Normalize date to start of day
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Preload("Schedule").
		Where("student_id = ? AND date = ? AND schedule_id = ?", studentID, dateOnly, scheduleID).
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
		Preload("Schedule").
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
		Preload("Schedule").
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
		Preload("Schedule").
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
		Preload("Schedule").
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
		Where("rf_id_code = ? AND is_active = ?", rfidCode, true).
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

// GetAttendanceSummaryByClass retrieves attendance summary grouped by class
func (r *repository) GetAttendanceSummaryByClass(ctx context.Context, schoolID uint, date time.Time) (*AttendanceSummaryResponse, []ClassSummaryItem, error) {
	// Normalize date to start of day
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	// Get all classes in the school
	var classes []models.Class
	if err := r.db.WithContext(ctx).
		Where("school_id = ?", schoolID).
		Order("name ASC").
		Find(&classes).Error; err != nil {
		return nil, nil, err
	}

	// Build summary per class
	classSummaries := make([]ClassSummaryItem, 0, len(classes))
	totalSummary := &AttendanceSummaryResponse{
		Date: dateOnly.Format("2006-01-02"),
	}

	for _, class := range classes {
		// Get total students in class
		var totalStudents int64
		if err := r.db.WithContext(ctx).
			Model(&models.Student{}).
			Where("class_id = ? AND is_active = ?", class.ID, true).
			Count(&totalStudents).Error; err != nil {
			continue
		}

		if totalStudents == 0 {
			continue // Skip empty classes
		}

		// Get attendance counts by status for this class
		type StatusCount struct {
			Status string
			Count  int64
		}
		var statusCounts []StatusCount

		err := r.db.WithContext(ctx).
			Model(&models.Attendance{}).
			Select("status, COUNT(*) as count").
			Joins("JOIN students ON students.id = attendances.student_id").
			Where("students.class_id = ? AND attendances.date = ?", class.ID, dateOnly).
			Group("status").
			Scan(&statusCounts).Error

		if err != nil {
			continue
		}

		// Build class summary
		classSummary := ClassSummaryItem{
			ClassID:       class.ID,
			ClassName:     class.Name,
			TotalStudents: int(totalStudents),
		}

		var presentCount int64
		for _, sc := range statusCounts {
			switch models.AttendanceStatus(sc.Status) {
			case models.AttendanceStatusOnTime:
				classSummary.Present += int(sc.Count)
				presentCount += sc.Count
			case models.AttendanceStatusLate:
				classSummary.Late += int(sc.Count)
				presentCount += sc.Count
			case models.AttendanceStatusVeryLate:
				classSummary.Late += int(sc.Count) // Count very_late as late for simplicity
				presentCount += sc.Count
			}
		}

		classSummary.Absent = int(totalStudents) - int(presentCount)

		classSummaries = append(classSummaries, classSummary)

		// Add to total summary
		totalSummary.TotalCount += int(totalStudents)
		totalSummary.Present += classSummary.Present
		totalSummary.Late += classSummary.Late
		totalSummary.Absent += classSummary.Absent
	}

	return totalSummary, classSummaries, nil
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
		ID:         attendance.ID,
		StudentID:  attendance.StudentID,
		ScheduleID: attendance.ScheduleID,
		Date:       attendance.Date.Format("2006-01-02"),
		Status:     attendance.Status,
		Method:     attendance.Method,
		CreatedAt:  attendance.CreatedAt,
		UpdatedAt:  attendance.UpdatedAt,
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
		// Handle nullable Class pointer
		if attendance.Student.Class != nil && attendance.Student.Class.ID != 0 {
			response.ClassName = attendance.Student.Class.Name
		}
	}

	// Include schedule info if loaded
	// Requirements: 3.10 - Show which schedule the attendance belongs to
	if attendance.Schedule != nil && attendance.Schedule.ID != 0 {
		response.ScheduleName = attendance.Schedule.Name
	}

	return response
}

// FindActiveSchedule finds the active schedule for a given time and day
// Requirements: 3.4, 3.5 - Determine which schedule is currently active based on current time
func (r *repository) FindActiveSchedule(ctx context.Context, schoolID uint, timestamp time.Time) (*models.AttendanceSchedule, error) {
	var schedules []models.AttendanceSchedule
	
	// Get current time in HH:MM:SS format
	currentTime := timestamp.Format("15:04:05")
	
	// Find all active schedules for this school
	err := r.db.WithContext(ctx).
		Where("school_id = ? AND is_active = ?", schoolID, true).
		Find(&schedules).Error

	if err != nil {
		return nil, err
	}

	// Find the schedule that matches the current time and day
	for _, schedule := range schedules {
		// Check if schedule is active on this day
		if !schedule.IsActiveOnDay(timestamp.Weekday()) {
			continue
		}

		// Normalize times for comparison
		startTime := schedule.StartTime
		if len(startTime) == 5 {
			startTime += ":00"
		}
		endTime := schedule.EndTime
		if len(endTime) == 5 {
			endTime += ":00"
		}

		// Check if current time is within schedule range
		if currentTime >= startTime && currentTime <= endTime {
			return &schedule, nil
		}
	}

	// No active schedule found for current time, try to find default
	defaultSchedule, err := r.FindDefaultSchedule(ctx, schoolID)
	if err == nil && defaultSchedule != nil {
		// Check if default schedule is active on this day
		if defaultSchedule.IsActive && defaultSchedule.IsActiveOnDay(timestamp.Weekday()) {
			return defaultSchedule, nil
		}
	}

	// Return nil if no schedule matches
	return nil, nil
}

// FindDefaultSchedule finds the default schedule for a school
func (r *repository) FindDefaultSchedule(ctx context.Context, schoolID uint) (*models.AttendanceSchedule, error) {
	var schedule models.AttendanceSchedule
	err := r.db.WithContext(ctx).
		Where("school_id = ? AND is_default = ?", schoolID, true).
		First(&schedule).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No default schedule is not an error
		}
		return nil, err
	}

	return &schedule, nil
}

// GetAttendanceForExport retrieves attendance records for export with filters
// Requirements: 1.2, 1.3 - Allow filtering by date range and class
func (r *repository) GetAttendanceForExport(ctx context.Context, schoolID uint, filter ExportFilter) ([]ExportAttendanceRecord, error) {
	var records []ExportAttendanceRecord

	// Parse dates
	startDate, err := time.Parse("2006-01-02", filter.StartDate)
	if err != nil {
		return nil, errors.New("invalid start_date format, expected YYYY-MM-DD")
	}
	endDate, err := time.Parse("2006-01-02", filter.EndDate)
	if err != nil {
		return nil, errors.New("invalid end_date format, expected YYYY-MM-DD")
	}

	// Build query
	query := r.db.WithContext(ctx).
		Table("attendances").
		Select(`
			students.nis as student_nis,
			students.nisn as student_nisn,
			students.name as student_name,
			classes.name as class_name,
			attendances.date,
			attendances.check_in_time,
			attendances.check_out_time,
			attendances.status,
			COALESCE(attendance_schedules.name, '') as schedule_name
		`).
		Joins("JOIN students ON students.id = attendances.student_id").
		Joins("JOIN classes ON classes.id = students.class_id").
		Joins("LEFT JOIN attendance_schedules ON attendance_schedules.id = attendances.schedule_id").
		Where("students.school_id = ?", schoolID).
		Where("attendances.date >= ? AND attendances.date <= ?", startDate, endDate)

	// Apply class filter if provided
	if filter.ClassID != nil {
		query = query.Where("students.class_id = ?", *filter.ClassID)
	}

	// Order by date and student name
	query = query.Order("attendances.date ASC, students.name ASC")

	// Execute query
	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var record ExportAttendanceRecord
		var date time.Time
		var checkInTime, checkOutTime *time.Time

		err := rows.Scan(
			&record.StudentNIS,
			&record.StudentNISN,
			&record.StudentName,
			&record.ClassName,
			&date,
			&checkInTime,
			&checkOutTime,
			&record.Status,
			&record.ScheduleName,
		)
		if err != nil {
			return nil, err
		}

		record.Date = date.Format("2006-01-02")
		if checkInTime != nil {
			record.CheckInTime = checkInTime.Format("15:04")
		}
		if checkOutTime != nil {
			record.CheckOutTime = checkOutTime.Format("15:04")
		}

		records = append(records, record)
	}

	return records, nil
}

// GetStudentsBySchool retrieves all active students in a school, optionally filtered by class
func (r *repository) GetStudentsBySchool(ctx context.Context, schoolID uint, classID *uint) ([]models.Student, error) {
	var students []models.Student

	query := r.db.WithContext(ctx).
		Preload("Class").
		Where("school_id = ? AND is_active = ?", schoolID, true)

	if classID != nil {
		query = query.Where("class_id = ?", *classID)
	}

	err := query.Order("name ASC").Find(&students).Error
	return students, err
}

// GetMonthlyRecap retrieves monthly attendance recap with aggregation
// Requirements: 2.1 - Display summary per student including total days present, late, very late, and absent
func (r *repository) GetMonthlyRecap(ctx context.Context, schoolID uint, filter MonthlyRecapFilter) (*MonthlyRecapResponse, error) {
	// Calculate date range for the month
	startDate := time.Date(filter.Year, time.Month(filter.Month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1) // Last day of the month

	// Get all students in the school (optionally filtered by class)
	students, err := r.GetStudentsBySchool(ctx, schoolID, filter.ClassID)
	if err != nil {
		return nil, err
	}

	// Calculate total school days (weekdays only, Mon-Fri)
	totalDays := countWeekdays(startDate, endDate)

	// Get class name if filtered by class
	var className string
	if filter.ClassID != nil {
		var class models.Class
		if err := r.db.WithContext(ctx).First(&class, *filter.ClassID).Error; err == nil {
			className = class.Name
		}
	}

	// Build response
	response := &MonthlyRecapResponse{
		Month:         filter.Month,
		Year:          filter.Year,
		TotalDays:     totalDays,
		ClassID:       filter.ClassID,
		ClassName:     className,
		StudentRecaps: make([]StudentRecapSummary, 0, len(students)),
	}

	// Get attendance summary for each student
	for _, student := range students {
		// Query attendance counts by status for this student in the month
		type StatusCount struct {
			Status string
			Count  int
		}
		var statusCounts []StatusCount

		err := r.db.WithContext(ctx).
			Model(&models.Attendance{}).
			Select("status, COUNT(*) as count").
			Where("student_id = ? AND date >= ? AND date <= ?", student.ID, startDate, endDate).
			Group("status").
			Scan(&statusCounts).Error

		if err != nil {
			continue
		}

		// Build student summary
		summary := StudentRecapSummary{
			StudentID:   student.ID,
			StudentNIS:  student.NIS,
			StudentNISN: student.NISN,
			StudentName: student.Name,
			ClassName:   student.Class.Name,
		}

		for _, sc := range statusCounts {
			switch models.AttendanceStatus(sc.Status) {
			case models.AttendanceStatusOnTime:
				summary.TotalPresent += sc.Count
			case models.AttendanceStatusLate:
				summary.TotalLate += sc.Count
			case models.AttendanceStatusVeryLate:
				summary.TotalVeryLate += sc.Count
			}
		}

		// Calculate absent days
		totalAttended := summary.TotalPresent + summary.TotalLate + summary.TotalVeryLate
		summary.TotalAbsent = totalDays - totalAttended
		if summary.TotalAbsent < 0 {
			summary.TotalAbsent = 0
		}

		// Calculate attendance percentage (present / total_days * 100)
		// Requirements: 2.2 - Calculate and display attendance percentage
		if totalDays > 0 {
			summary.AttendancePercent = float64(summary.TotalPresent) / float64(totalDays) * 100
		}

		response.StudentRecaps = append(response.StudentRecaps, summary)
	}

	// Sort by attendance percentage descending
	// Requirements: 2.6 - Sort students by attendance percentage (highest to lowest)
	sortStudentRecapsByPercentage(response.StudentRecaps)

	return response, nil
}

// countWeekdays counts the number of weekdays (Mon-Fri) between two dates
func countWeekdays(start, end time.Time) int {
	count := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		weekday := d.Weekday()
		if weekday != time.Saturday && weekday != time.Sunday {
			count++
		}
	}
	return count
}

// sortStudentRecapsByPercentage sorts student recaps by attendance percentage descending
func sortStudentRecapsByPercentage(recaps []StudentRecapSummary) {
	for i := 0; i < len(recaps)-1; i++ {
		for j := i + 1; j < len(recaps); j++ {
			if recaps[j].AttendancePercent > recaps[i].AttendancePercent {
				recaps[i], recaps[j] = recaps[j], recaps[i]
			}
		}
	}
}

// FindClassByHomeroomTeacher finds the class assigned to a wali_kelas (homeroom teacher)
// Requirements: 2.7 - Wali_Kelas only sees students from their assigned class
func (r *repository) FindClassByHomeroomTeacher(ctx context.Context, schoolID uint, teacherID uint) (*models.Class, error) {
	var class models.Class
	err := r.db.WithContext(ctx).
		Where("school_id = ? AND homeroom_teacher_id = ?", schoolID, teacherID).
		First(&class).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No class assigned is not an error
		}
		return nil, err
	}

	return &class, nil
}
