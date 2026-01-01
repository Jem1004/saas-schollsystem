package attendance

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/school-management/backend/internal/domain/models"
	"github.com/school-management/backend/internal/modules/device"
)

var (
	ErrInvalidDate        = errors.New("format tanggal tidak valid")
	ErrInvalidTime        = errors.New("format waktu tidak valid")
	ErrStudentIDRequired  = errors.New("ID siswa wajib diisi")
	ErrDateRequired       = errors.New("tanggal wajib diisi")
	ErrAPIKeyRequired     = errors.New("API key wajib diisi")
	ErrRFIDCodeRequired   = errors.New("kode RFID wajib diisi")
	ErrCheckOutBeforeIn   = errors.New("waktu check-out tidak boleh sebelum waktu check-in")
	ErrNoCheckIn          = errors.New("tidak dapat mencatat check-out tanpa check-in")
	ErrAlreadyCheckedOut  = errors.New("siswa sudah melakukan check-out")
)

// Service defines the interface for attendance business logic
type Service interface {
	// RFID attendance (from ESP32 devices)
	RecordRFIDAttendance(ctx context.Context, req RFIDAttendanceRequest) (*RFIDAttendanceResponse, error)
	
	// Manual attendance (fallback)
	RecordManualAttendance(ctx context.Context, schoolID uint, req ManualAttendanceRequest) (*AttendanceResponse, error)
	RecordBulkManualAttendance(ctx context.Context, schoolID uint, req BulkManualAttendanceRequest) ([]AttendanceResponse, error)
	
	// Query operations
	GetAttendanceByID(ctx context.Context, id uint) (*AttendanceResponse, error)
	GetStudentAttendance(ctx context.Context, studentID uint, startDate, endDate time.Time) ([]AttendanceResponse, error)
	GetClassAttendance(ctx context.Context, classID uint, date time.Time) (*ClassAttendanceSummaryResponse, error)
	GetSchoolAttendanceSummary(ctx context.Context, schoolID uint, date time.Time) (*SchoolAttendanceSummaryResponse, error)
	GetAllAttendance(ctx context.Context, schoolID uint, filter AttendanceFilter) (*AttendanceListResponse, error)
	
	// Delete
	DeleteAttendance(ctx context.Context, id uint) error
}

// service implements the Service interface
type service struct {
	repo          Repository
	deviceService device.Service
	policy        AttendancePolicy
}

// NewService creates a new attendance service
func NewService(repo Repository, deviceService device.Service, policy AttendancePolicy) Service {
	return &service{
		repo:          repo,
		deviceService: deviceService,
		policy:        policy,
	}
}

// RecordRFIDAttendance records attendance from RFID device
// Requirements: 5.1, 5.2 - WHEN a student taps RFID card, record check-in or check-out
func (s *service) RecordRFIDAttendance(ctx context.Context, req RFIDAttendanceRequest) (*RFIDAttendanceResponse, error) {
	// Validate required fields
	if req.APIKey == "" {
		return nil, ErrAPIKeyRequired
	}
	if req.RFIDCode == "" {
		return nil, ErrRFIDCodeRequired
	}

	// Validate API key and get device info
	validation, err := s.deviceService.ValidateAPIKey(ctx, req.APIKey)
	if err != nil {
		return nil, err
	}
	if !validation.Valid {
		return nil, device.ErrInvalidAPIKey
	}

	// Find student by RFID code
	student, err := s.repo.FindStudentByRFID(ctx, req.RFIDCode)
	if err != nil {
		log.Printf("RFID attendance failed: student not found for RFID %s", req.RFIDCode)
		return nil, err
	}

	// Verify student belongs to the same school as the device
	if student.SchoolID != validation.SchoolID {
		log.Printf("RFID attendance failed: student school %d doesn't match device school %d", student.SchoolID, validation.SchoolID)
		return nil, ErrInvalidRFIDCode
	}

	// Use current time if timestamp is zero
	timestamp := req.Timestamp
	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	// Get date from timestamp
	date := time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, timestamp.Location())

	// Check for existing attendance record for today
	existing, err := s.repo.FindByStudentAndDate(ctx, student.ID, date)
	
	if err != nil && !errors.Is(err, ErrAttendanceNotFound) {
		return nil, err
	}

	var response *RFIDAttendanceResponse

	if existing == nil {
		// First tap of the day - record check-in
		// Requirements: 5.2 - First attendance record SHALL be recorded as check-in
		status := s.policy.DetermineAttendanceStatus(timestamp, student.SchoolID)
		
		attendance := &models.Attendance{
			StudentID: student.ID,
			Date:      date,
			Method:    models.AttendanceMethodRFID,
			Status:    status,
		}
		attendance.SetCheckIn(timestamp)

		if err := s.repo.Create(ctx, attendance); err != nil {
			return nil, err
		}

		response = &RFIDAttendanceResponse{
			Success:     true,
			StudentID:   student.ID,
			StudentName: student.Name,
			Type:        "check_in",
			Status:      status,
			Time:        timestamp,
			Message:     "Check-in recorded successfully",
		}

		log.Printf("RFID check-in recorded: student %s (%d) at %s, status: %s", 
			student.Name, student.ID, timestamp.Format("15:04"), status)

	} else if existing.CheckOutTime == nil {
		// Second tap - record check-out
		// Requirements: 5.2 - Second attendance record SHALL update with check-out time
		if err := existing.SetCheckOut(timestamp); err != nil {
			return nil, err
		}

		if err := s.repo.Update(ctx, existing); err != nil {
			return nil, err
		}

		response = &RFIDAttendanceResponse{
			Success:     true,
			StudentID:   student.ID,
			StudentName: student.Name,
			Type:        "check_out",
			Time:        timestamp,
			Message:     "Check-out recorded successfully",
		}

		log.Printf("RFID check-out recorded: student %s (%d) at %s", 
			student.Name, student.ID, timestamp.Format("15:04"))

	} else {
		// Already checked out
		response = &RFIDAttendanceResponse{
			Success:     false,
			StudentID:   student.ID,
			StudentName: student.Name,
			Type:        "already_checked_out",
			Time:        timestamp,
			Message:     "Student has already checked out today",
		}

		log.Printf("RFID tap ignored: student %s (%d) already checked out", student.Name, student.ID)
	}

	// TODO: Trigger notification to parent (async)
	// Requirements: 5.3 - WHEN attendance is recorded, THE System SHALL trigger notification to parent

	return response, nil
}


// RecordManualAttendance records manual attendance entry
// Requirements: 5.5 - IF RFID system fails, THEN THE System SHALL allow manual attendance entry
func (s *service) RecordManualAttendance(ctx context.Context, schoolID uint, req ManualAttendanceRequest) (*AttendanceResponse, error) {
	// Validate required fields
	if req.StudentID == 0 {
		return nil, ErrStudentIDRequired
	}
	if req.Date == "" {
		return nil, ErrDateRequired
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, ErrInvalidDate
	}

	// Verify student exists and belongs to the school
	student, err := s.repo.FindStudentByID(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if student.SchoolID != schoolID {
		return nil, ErrStudentNotFound
	}

	// Parse check-in time if provided
	var checkInTime *time.Time
	if req.CheckInTime != nil && *req.CheckInTime != "" {
		t, err := time.Parse("15:04", *req.CheckInTime)
		if err != nil {
			return nil, ErrInvalidTime
		}
		// Combine date and time
		combined := time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), 0, 0, date.Location())
		checkInTime = &combined
	}

	// Parse check-out time if provided
	var checkOutTime *time.Time
	if req.CheckOutTime != nil && *req.CheckOutTime != "" {
		t, err := time.Parse("15:04", *req.CheckOutTime)
		if err != nil {
			return nil, ErrInvalidTime
		}
		// Combine date and time
		combined := time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), 0, 0, date.Location())
		checkOutTime = &combined
	}

	// Validate check-out is after check-in
	if checkInTime != nil && checkOutTime != nil && checkOutTime.Before(*checkInTime) {
		return nil, ErrCheckOutBeforeIn
	}

	// Check for existing attendance record
	existing, err := s.repo.FindByStudentAndDate(ctx, req.StudentID, date)
	
	var attendance *models.Attendance

	if err != nil && errors.Is(err, ErrAttendanceNotFound) {
		// Create new attendance record
		status := models.AttendanceStatusAbsent
		if checkInTime != nil {
			status = s.policy.DetermineAttendanceStatus(*checkInTime, schoolID)
		}

		attendance = &models.Attendance{
			StudentID:    req.StudentID,
			Date:         date,
			CheckInTime:  checkInTime,
			CheckOutTime: checkOutTime,
			Method:       models.AttendanceMethodManual,
			Status:       status,
		}

		if err := s.repo.Create(ctx, attendance); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		// Update existing record
		if checkInTime != nil {
			existing.CheckInTime = checkInTime
			existing.Status = s.policy.DetermineAttendanceStatus(*checkInTime, schoolID)
		}
		if checkOutTime != nil {
			existing.CheckOutTime = checkOutTime
		}
		existing.Method = models.AttendanceMethodManual

		if err := s.repo.Update(ctx, existing); err != nil {
			return nil, err
		}
		attendance = existing
	}

	// Reload with relations
	attendance, err = s.repo.FindByID(ctx, attendance.ID)
	if err != nil {
		return nil, err
	}

	return toAttendanceResponse(attendance), nil
}

// RecordBulkManualAttendance records multiple manual attendance entries
func (s *service) RecordBulkManualAttendance(ctx context.Context, schoolID uint, req BulkManualAttendanceRequest) ([]AttendanceResponse, error) {
	if req.Date == "" {
		return nil, ErrDateRequired
	}

	responses := make([]AttendanceResponse, 0, len(req.Attendances))

	for _, item := range req.Attendances {
		manualReq := ManualAttendanceRequest{
			StudentID:    item.StudentID,
			Date:         req.Date,
			CheckInTime:  item.CheckInTime,
			CheckOutTime: item.CheckOutTime,
		}

		response, err := s.RecordManualAttendance(ctx, schoolID, manualReq)
		if err != nil {
			log.Printf("Failed to record manual attendance for student %d: %v", item.StudentID, err)
			continue
		}
		responses = append(responses, *response)
	}

	return responses, nil
}

// GetAttendanceByID retrieves an attendance record by ID
func (s *service) GetAttendanceByID(ctx context.Context, id uint) (*AttendanceResponse, error) {
	attendance, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toAttendanceResponse(attendance), nil
}

// GetStudentAttendance retrieves attendance records for a student
func (s *service) GetStudentAttendance(ctx context.Context, studentID uint, startDate, endDate time.Time) ([]AttendanceResponse, error) {
	attendances, err := s.repo.FindByStudent(ctx, studentID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	responses := make([]AttendanceResponse, len(attendances))
	for i, a := range attendances {
		responses[i] = *toAttendanceResponse(&a)
	}

	return responses, nil
}

// GetClassAttendance retrieves attendance for a class on a specific date
func (s *service) GetClassAttendance(ctx context.Context, classID uint, date time.Time) (*ClassAttendanceSummaryResponse, error) {
	return s.repo.GetClassAttendanceSummary(ctx, classID, date)
}

// GetSchoolAttendanceSummary retrieves school-wide attendance summary
// Requirements: 5.4 - WHEN an Admin_Sekolah views attendance dashboard, THE System SHALL display summary statistics
func (s *service) GetSchoolAttendanceSummary(ctx context.Context, schoolID uint, date time.Time) (*SchoolAttendanceSummaryResponse, error) {
	summary, byClass, err := s.repo.GetAttendanceSummaryByClass(ctx, schoolID, date)
	if err != nil {
		return nil, err
	}

	return &SchoolAttendanceSummaryResponse{
		SchoolID: schoolID,
		Date:     date.Format("2006-01-02"),
		Summary:  *summary,
		ByClass:  byClass,
	}, nil
}

// GetAllAttendance retrieves all attendance records with pagination
func (s *service) GetAllAttendance(ctx context.Context, schoolID uint, filter AttendanceFilter) (*AttendanceListResponse, error) {
	// Set defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	attendances, total, err := s.repo.FindAll(ctx, schoolID, filter)
	if err != nil {
		return nil, err
	}

	// Convert to response
	responses := make([]AttendanceResponse, len(attendances))
	for i, a := range attendances {
		responses[i] = *toAttendanceResponse(&a)
	}

	// Calculate total pages
	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &AttendanceListResponse{
		Attendances: responses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// DeleteAttendance deletes an attendance record
func (s *service) DeleteAttendance(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
