package homeroom

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrStudentIDRequired      = errors.New("ID siswa wajib diisi")
	ErrContentRequired        = errors.New("konten wajib diisi")
	ErrStudentNotInSchool     = errors.New("siswa bukan dari sekolah ini")
	ErrStudentNotInClass      = errors.New("siswa bukan dari kelas yang Anda ampu")
	ErrNotAuthorized          = errors.New("tidak memiliki izin untuk melakukan aksi ini")
	ErrNoClassAssigned        = errors.New("tidak ada kelas yang ditugaskan untuk guru ini")
	ErrGradeNotFound          = errors.New("nilai tidak ditemukan")
	ErrTitleRequired          = errors.New("judul wajib diisi")
	ErrScoreInvalid           = errors.New("nilai harus antara 0 dan 100")
	ErrAttendanceNotFound     = errors.New("data absensi tidak ditemukan")
	ErrAttendanceAlreadyExists = errors.New("data absensi sudah ada")
	ErrInvalidStatus          = errors.New("status absensi tidak valid")
)

// Service defines the interface for Homeroom Note business logic
type Service interface {
	// Note operations
	CreateNote(ctx context.Context, schoolID, teacherID uint, req CreateNoteRequest) (*NoteResponse, error)
	GetNoteByID(ctx context.Context, id uint) (*NoteResponse, error)
	GetStudentNotes(ctx context.Context, studentID uint) ([]NoteResponse, error)
	GetNotes(ctx context.Context, schoolID uint, filter NoteFilter) (*NoteListResponse, error)
	UpdateNote(ctx context.Context, noteID uint, req UpdateNoteRequest) (*NoteResponse, error)
	DeleteNote(ctx context.Context, noteID uint) error

	// Grade operations
	CreateGrade(ctx context.Context, schoolID, teacherID uint, req CreateGradeRequest) (*GradeResponse, error)
	CreateBatchGrades(ctx context.Context, schoolID, teacherID uint, req BatchGradeRequest) ([]GradeResponse, error)
	GetGradeByID(ctx context.Context, id uint) (*GradeResponse, error)
	GetClassGrades(ctx context.Context, schoolID, teacherID uint, page, pageSize int) (*GradeListResponse, error)
	GetStudentGrades(ctx context.Context, teacherID, studentID uint, page, pageSize int) (*GradeListResponse, error)
	UpdateGrade(ctx context.Context, gradeID uint, req UpdateGradeRequest) (*GradeResponse, error)
	DeleteGrade(ctx context.Context, gradeID uint) error

	// Teacher validation
	ValidateTeacherAccess(ctx context.Context, teacherID, studentID uint) error
	GetTeacherClassID(ctx context.Context, teacherID uint) (*uint, error)

	// Summary
	GetStudentNoteSummary(ctx context.Context, studentID uint) (*StudentNoteSummary, error)
	GetClassNotes(ctx context.Context, classID uint, filter NoteFilter) (*NoteListResponse, error)

	// Wali Kelas Dashboard
	GetHomeroomStats(ctx context.Context, schoolID, teacherID uint) (*HomeroomStatsResponse, error)
	GetMyClass(ctx context.Context, teacherID uint) (*ClassInfoResponse, error)
	GetClassStudents(ctx context.Context, teacherID uint, page, pageSize int) (*ClassStudentListResponse, error)
	GetClassAttendance(ctx context.Context, teacherID uint, date string) (*ClassAttendanceListResponse, error)

	// Manual Attendance
	RecordManualAttendance(ctx context.Context, schoolID, teacherID uint, req ManualAttendanceRequest) (*StudentAttendanceResponse, error)
	UpdateAttendance(ctx context.Context, teacherID, attendanceID uint, req UpdateAttendanceRequest) (*StudentAttendanceResponse, error)

	// Schedules
	GetActiveSchedules(ctx context.Context, schoolID uint, date string) ([]ScheduleResponse, error)
}

// service implements the Service interface
type service struct {
	repo Repository
	db   *gorm.DB
}

// NewService creates a new Homeroom service
func NewService(repo Repository, db *gorm.DB) Service {
	return &service{repo: repo, db: db}
}

// CreateNote creates a new homeroom note
// Requirements: 11.1 - WHEN a Wali_Kelas creates a note, THE System SHALL require content and associate it with a student
// Requirements: 11.4 - THE System SHALL validate that Wali_Kelas can only create notes for students in their assigned class
func (s *service) CreateNote(ctx context.Context, schoolID, teacherID uint, req CreateNoteRequest) (*NoteResponse, error) {
	// Validate required fields
	if req.StudentID == 0 {
		return nil, ErrStudentIDRequired
	}
	if req.Content == "" {
		return nil, ErrContentRequired
	}

	// Verify student belongs to the school
	student, err := s.repo.FindStudentByID(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if student.SchoolID != schoolID {
		return nil, ErrStudentNotInSchool
	}

	// Validate teacher has access to this student (wali kelas validation)
	if err := s.ValidateTeacherAccess(ctx, teacherID, req.StudentID); err != nil {
		return nil, err
	}

	note := &models.HomeroomNote{
		StudentID: req.StudentID,
		TeacherID: teacherID,
		Content:   req.Content,
	}

	if err := s.repo.Create(ctx, note); err != nil {
		return nil, err
	}

	// Reload with relations
	note, err = s.repo.FindByID(ctx, note.ID)
	if err != nil {
		return nil, err
	}

	// TODO: Trigger notification to parent (async)
	// Requirements: 11.2 - WHEN a note is saved, THE System SHALL trigger notification to the parent

	return toNoteResponse(note), nil
}


// GetNoteByID retrieves a homeroom note by ID
func (s *service) GetNoteByID(ctx context.Context, id uint) (*NoteResponse, error) {
	note, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toNoteResponse(note), nil
}

// GetStudentNotes retrieves all homeroom notes for a student
// Requirements: 11.3 - WHEN parents view notes, THE System SHALL display all homeroom notes for their child
// Requirements: 11.5 - THE System SHALL maintain note history with timestamps and author information
func (s *service) GetStudentNotes(ctx context.Context, studentID uint) ([]NoteResponse, error) {
	notes, err := s.repo.FindByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	responses := make([]NoteResponse, len(notes))
	for i, n := range notes {
		responses[i] = *toNoteResponse(&n)
	}
	return responses, nil
}

// GetNotes retrieves homeroom notes with pagination and filtering
func (s *service) GetNotes(ctx context.Context, schoolID uint, filter NoteFilter) (*NoteListResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	notes, total, err := s.repo.FindAll(ctx, schoolID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]NoteResponse, len(notes))
	for i, n := range notes {
		responses[i] = *toNoteResponse(&n)
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &NoteListResponse{
		Notes: responses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateNote updates a homeroom note
func (s *service) UpdateNote(ctx context.Context, noteID uint, req UpdateNoteRequest) (*NoteResponse, error) {
	note, err := s.repo.FindByID(ctx, noteID)
	if err != nil {
		return nil, err
	}

	if req.Content != "" {
		note.Content = req.Content
	}

	if err := s.repo.Update(ctx, note); err != nil {
		return nil, err
	}

	return toNoteResponse(note), nil
}

// DeleteNote deletes a homeroom note
func (s *service) DeleteNote(ctx context.Context, noteID uint) error {
	return s.repo.Delete(ctx, noteID)
}

// ValidateTeacherAccess validates that a teacher can access a student's notes
// Requirements: 11.4 - THE System SHALL validate that Wali_Kelas can only create notes for students in their assigned class
func (s *service) ValidateTeacherAccess(ctx context.Context, teacherID, studentID uint) error {
	// Get teacher's assigned class
	classID, err := s.GetTeacherClassID(ctx, teacherID)
	if err != nil {
		return err
	}
	if classID == nil {
		return ErrNoClassAssigned
	}

	// Get student's class
	student, err := s.repo.FindStudentByID(ctx, studentID)
	if err != nil {
		return err
	}

	// Check if student is in teacher's class (handle nullable ClassID)
	if student.ClassID == nil || *student.ClassID != *classID {
		return ErrStudentNotInClass
	}

	return nil
}

// GetTeacherClassID returns the class ID assigned to a wali kelas
func (s *service) GetTeacherClassID(ctx context.Context, teacherID uint) (*uint, error) {
	var class models.Class
	err := s.db.WithContext(ctx).
		Where("homeroom_teacher_id = ?", teacherID).
		First(&class).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &class.ID, nil
}

// GetStudentNoteSummary retrieves note summary for a student
func (s *service) GetStudentNoteSummary(ctx context.Context, studentID uint) (*StudentNoteSummary, error) {
	return s.repo.GetStudentNoteSummary(ctx, studentID)
}

// GetClassNotes retrieves homeroom notes for all students in a class
func (s *service) GetClassNotes(ctx context.Context, classID uint, filter NoteFilter) (*NoteListResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	notes, total, err := s.repo.GetClassNotes(ctx, classID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]NoteResponse, len(notes))
	for i, n := range notes {
		responses[i] = *toNoteResponse(&n)
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &NoteListResponse{
		Notes: responses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// ==================== Response Converters ====================

func toNoteResponse(n *models.HomeroomNote) *NoteResponse {
	response := &NoteResponse{
		ID:        n.ID,
		StudentID: n.StudentID,
		TeacherID: n.TeacherID,
		Content:   n.Content,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}

	if n.Student.ID != 0 {
		response.StudentName = n.Student.Name
		response.StudentNIS = n.Student.NIS
		response.StudentNISN = n.Student.NISN
		// Handle nullable Class pointer
		if n.Student.Class != nil && n.Student.Class.ID != 0 {
			response.ClassName = n.Student.Class.Name
		}
	}

	if n.Teacher.ID != 0 {
		response.TeacherName = n.Teacher.Username
	}

	return response
}

// ==================== Wali Kelas Dashboard Methods ====================

// GetHomeroomStats retrieves dashboard statistics for wali kelas
func (s *service) GetHomeroomStats(ctx context.Context, schoolID, teacherID uint) (*HomeroomStatsResponse, error) {
	// Get teacher's assigned class
	classID, err := s.GetTeacherClassID(ctx, teacherID)
	if err != nil {
		return nil, err
	}
	if classID == nil {
		return nil, ErrNoClassAssigned
	}

	// Get class info
	var class models.Class
	if err := s.db.WithContext(ctx).
		Preload("HomeroomTeacher").
		First(&class, *classID).Error; err != nil {
		return nil, err
	}

	// Count total students in class
	var totalStudents int64
	if err := s.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("class_id = ? AND is_active = ?", *classID, true).
		Count(&totalStudents).Error; err != nil {
		return nil, err
	}

	// Get today's attendance stats
	today := s.db.NowFunc().Format("2006-01-02")
	todayAttendance := TodayAttendanceStats{}

	// Count on_time
	var onTimeCount int64
	s.db.WithContext(ctx).
		Model(&models.Attendance{}).
		Joins("JOIN students ON students.id = attendances.student_id").
		Where("students.class_id = ? AND attendances.date = ? AND attendances.status = ?", *classID, today, "on_time").
		Count(&onTimeCount)
	todayAttendance.Present = onTimeCount

	// Count late + very_late
	var lateCount int64
	s.db.WithContext(ctx).
		Model(&models.Attendance{}).
		Joins("JOIN students ON students.id = attendances.student_id").
		Where("students.class_id = ? AND attendances.date = ? AND attendances.status IN ?", *classID, today, []string{"late", "very_late"}).
		Count(&lateCount)
	todayAttendance.Late = lateCount

	// Absent = total - (on_time + late)
	todayAttendance.Absent = totalStudents - onTimeCount - lateCount

	// Get recent grades (last 5)
	var grades []models.Grade
	s.db.WithContext(ctx).
		Preload("Student").
		Joins("JOIN students ON students.id = grades.student_id").
		Where("students.class_id = ?", *classID).
		Order("grades.created_at DESC").
		Limit(5).
		Find(&grades)

	recentGrades := make([]GradeResponse, len(grades))
	for i, g := range grades {
		recentGrades[i] = GradeResponse{
			ID:          g.ID,
			StudentID:   g.StudentID,
			StudentName: g.Student.Name,
			StudentNIS:  g.Student.NIS,
			Title:       g.Title,
			Score:       g.Score,
			CreatedBy:   g.CreatedBy,
			CreatedAt:   g.CreatedAt,
			UpdatedAt:   g.UpdatedAt,
		}
	}

	// Get recent notes (last 5)
	var notes []models.HomeroomNote
	s.db.WithContext(ctx).
		Preload("Student").
		Preload("Teacher").
		Joins("JOIN students ON students.id = homeroom_notes.student_id").
		Where("students.class_id = ?", *classID).
		Order("homeroom_notes.created_at DESC").
		Limit(5).
		Find(&notes)

	recentNotes := make([]NoteResponse, len(notes))
	for i, n := range notes {
		recentNotes[i] = *toNoteResponse(&n)
	}

	return &HomeroomStatsResponse{
		ClassID:         *classID,
		ClassName:       class.Name,
		TotalStudents:   totalStudents,
		TodayAttendance: todayAttendance,
		RecentGrades:    recentGrades,
		RecentNotes:     recentNotes,
	}, nil
}

// GetMyClass retrieves class information for wali kelas
func (s *service) GetMyClass(ctx context.Context, teacherID uint) (*ClassInfoResponse, error) {
	// Get teacher's assigned class
	var class models.Class
	if err := s.db.WithContext(ctx).
		Preload("HomeroomTeacher").
		Where("homeroom_teacher_id = ?", teacherID).
		First(&class).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoClassAssigned
		}
		return nil, err
	}

	// Count students
	var totalStudents int64
	s.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("class_id = ? AND is_active = ?", class.ID, true).
		Count(&totalStudents)

	response := &ClassInfoResponse{
		ID:                class.ID,
		Name:              class.Name,
		Grade:             class.Grade,
		Year:              class.Year,
		TotalStudents:     totalStudents,
		HomeroomTeacherID: teacherID,
	}

	if class.HomeroomTeacher != nil {
		response.HomeroomTeacherName = class.HomeroomTeacher.Username
	}

	return response, nil
}

// GetClassStudents retrieves students in wali kelas's class
func (s *service) GetClassStudents(ctx context.Context, teacherID uint, page, pageSize int) (*ClassStudentListResponse, error) {
	// Get teacher's assigned class
	classID, err := s.GetTeacherClassID(ctx, teacherID)
	if err != nil {
		return nil, err
	}
	if classID == nil {
		return nil, ErrNoClassAssigned
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	// Count total
	var total int64
	s.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("class_id = ? AND is_active = ?", *classID, true).
		Count(&total)

	// Get students
	var students []models.Student
	offset := (page - 1) * pageSize
	s.db.WithContext(ctx).
		Where("class_id = ? AND is_active = ?", *classID, true).
		Order("name ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&students)

	data := make([]ClassStudentResponse, len(students))
	for i, st := range students {
		data[i] = ClassStudentResponse{
			ID:       st.ID,
			NIS:      st.NIS,
			NISN:     st.NISN,
			Name:     st.Name,
			IsActive: st.IsActive,
		}
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &ClassStudentListResponse{
		Data: data,
		Pagination: PaginationMeta{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetClassAttendance retrieves attendance for wali kelas's class on a specific date
func (s *service) GetClassAttendance(ctx context.Context, teacherID uint, date string) (*ClassAttendanceListResponse, error) {
	// Get teacher's assigned class
	classID, err := s.GetTeacherClassID(ctx, teacherID)
	if err != nil {
		return nil, err
	}
	if classID == nil {
		return nil, ErrNoClassAssigned
	}

	// Get all students in class
	var students []models.Student
	s.db.WithContext(ctx).
		Where("class_id = ? AND is_active = ?", *classID, true).
		Order("name ASC").
		Find(&students)

	// Get attendance records for the date
	var attendances []models.Attendance
	s.db.WithContext(ctx).
		Preload("Student").
		Joins("JOIN students ON students.id = attendances.student_id").
		Where("students.class_id = ? AND attendances.date = ?", *classID, date).
		Find(&attendances)

	// Create a map of student_id -> attendance
	attendanceMap := make(map[uint]*models.Attendance)
	for i := range attendances {
		attendanceMap[attendances[i].StudentID] = &attendances[i]
	}

	// Build response for all students
	data := make([]StudentAttendanceResponse, len(students))
	for i, st := range students {
		resp := StudentAttendanceResponse{
			StudentID:   st.ID,
			StudentName: st.Name,
			StudentNIS:  st.NIS,
			Date:        date,
			Status:      "absent", // Default to absent
			Method:      "",
		}

		if att, ok := attendanceMap[st.ID]; ok {
			resp.ID = att.ID
			resp.Status = string(att.Status)
			resp.Method = string(att.Method)
			resp.CreatedAt = att.CreatedAt.Format("2006-01-02T15:04:05Z")
			resp.UpdatedAt = att.UpdatedAt.Format("2006-01-02T15:04:05Z")

			if att.CheckInTime != nil {
				checkIn := att.CheckInTime.Format("15:04")
				resp.CheckInTime = &checkIn
			}
			if att.CheckOutTime != nil {
				checkOut := att.CheckOutTime.Format("15:04")
				resp.CheckOutTime = &checkOut
			}
		}

		data[i] = resp
	}

	return &ClassAttendanceListResponse{
		Data: data,
		Pagination: PaginationMeta{
			Page:       1,
			PageSize:   len(data),
			Total:      int64(len(data)),
			TotalPages: 1,
		},
	}, nil
}

// ==================== Grade Methods ====================

// CreateGrade creates a new grade for a student
func (s *service) CreateGrade(ctx context.Context, schoolID, teacherID uint, req CreateGradeRequest) (*GradeResponse, error) {
	// Validate required fields
	if req.StudentID == 0 {
		return nil, ErrStudentIDRequired
	}
	if req.Title == "" {
		return nil, ErrTitleRequired
	}
	if req.Score < 0 || req.Score > 100 {
		return nil, ErrScoreInvalid
	}

	// Verify student belongs to the school
	student, err := s.repo.FindStudentByID(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if student.SchoolID != schoolID {
		return nil, ErrStudentNotInSchool
	}

	// Validate teacher has access to this student
	if err := s.ValidateTeacherAccess(ctx, teacherID, req.StudentID); err != nil {
		return nil, err
	}

	grade := &models.Grade{
		StudentID:   req.StudentID,
		Title:       req.Title,
		Score:       req.Score,
		Description: req.Description,
		CreatedBy:   teacherID,
	}

	if err := s.db.WithContext(ctx).Create(grade).Error; err != nil {
		return nil, err
	}

	// Reload with relations
	if err := s.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Preload("Creator").
		First(grade, grade.ID).Error; err != nil {
		return nil, err
	}

	return toGradeResponse(grade), nil
}

// CreateBatchGrades creates multiple grades at once
func (s *service) CreateBatchGrades(ctx context.Context, schoolID, teacherID uint, req BatchGradeRequest) ([]GradeResponse, error) {
	if req.Title == "" {
		return nil, ErrTitleRequired
	}
	if len(req.Grades) == 0 {
		return nil, ErrStudentIDRequired
	}

	// Get teacher's class
	classID, err := s.GetTeacherClassID(ctx, teacherID)
	if err != nil {
		return nil, err
	}
	if classID == nil {
		return nil, ErrNoClassAssigned
	}

	var responses []GradeResponse

	for _, entry := range req.Grades {
		if entry.Score < 0 || entry.Score > 100 {
			continue // Skip invalid scores
		}

		// Verify student belongs to the school and class
		student, err := s.repo.FindStudentByID(ctx, entry.StudentID)
		if err != nil {
			continue
		}
		if student.SchoolID != schoolID {
			continue
		}
		if student.ClassID == nil || *student.ClassID != *classID {
			continue
		}

		grade := &models.Grade{
			StudentID:   entry.StudentID,
			Title:       req.Title,
			Score:       entry.Score,
			Description: req.Description,
			CreatedBy:   teacherID,
		}

		if err := s.db.WithContext(ctx).Create(grade).Error; err != nil {
			continue
		}

		// Reload with relations
		if err := s.db.WithContext(ctx).
			Preload("Student").
			Preload("Student.Class").
			Preload("Creator").
			First(grade, grade.ID).Error; err != nil {
			continue
		}

		responses = append(responses, *toGradeResponse(grade))
	}

	return responses, nil
}

// GetGradeByID retrieves a grade by ID
func (s *service) GetGradeByID(ctx context.Context, id uint) (*GradeResponse, error) {
	var grade models.Grade
	if err := s.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Preload("Creator").
		First(&grade, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGradeNotFound
		}
		return nil, err
	}
	return toGradeResponse(&grade), nil
}

// GetClassGrades retrieves grades for wali kelas's class
func (s *service) GetClassGrades(ctx context.Context, schoolID, teacherID uint, page, pageSize int) (*GradeListResponse, error) {
	// Get teacher's assigned class
	classID, err := s.GetTeacherClassID(ctx, teacherID)
	if err != nil {
		return nil, err
	}
	if classID == nil {
		return nil, ErrNoClassAssigned
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	// Count total
	var total int64
	s.db.WithContext(ctx).
		Model(&models.Grade{}).
		Joins("JOIN students ON students.id = grades.student_id").
		Where("students.class_id = ?", *classID).
		Count(&total)

	// Get grades
	var grades []models.Grade
	offset := (page - 1) * pageSize
	s.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Preload("Creator").
		Joins("JOIN students ON students.id = grades.student_id").
		Where("students.class_id = ?", *classID).
		Order("grades.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&grades)

	data := make([]GradeResponse, len(grades))
	for i, g := range grades {
		data[i] = *toGradeResponse(&g)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &GradeListResponse{
		Data: data,
		Pagination: PaginationMeta{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetStudentGrades retrieves grades for a specific student
func (s *service) GetStudentGrades(ctx context.Context, teacherID, studentID uint, page, pageSize int) (*GradeListResponse, error) {
	// Validate teacher has access to this student
	if err := s.ValidateTeacherAccess(ctx, teacherID, studentID); err != nil {
		return nil, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}

	// Count total
	var total int64
	s.db.WithContext(ctx).
		Model(&models.Grade{}).
		Where("student_id = ?", studentID).
		Count(&total)

	// Get grades
	var grades []models.Grade
	offset := (page - 1) * pageSize
	s.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Preload("Creator").
		Where("student_id = ?", studentID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&grades)

	data := make([]GradeResponse, len(grades))
	for i, g := range grades {
		data[i] = *toGradeResponse(&g)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &GradeListResponse{
		Data: data,
		Pagination: PaginationMeta{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateGrade updates a grade record
func (s *service) UpdateGrade(ctx context.Context, gradeID uint, req UpdateGradeRequest) (*GradeResponse, error) {
	var grade models.Grade
	if err := s.db.WithContext(ctx).First(&grade, gradeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGradeNotFound
		}
		return nil, err
	}

	if req.Title != "" {
		grade.Title = req.Title
	}
	if req.Score >= 0 && req.Score <= 100 {
		grade.Score = req.Score
	}
	if req.Description != "" {
		grade.Description = req.Description
	}

	if err := s.db.WithContext(ctx).Save(&grade).Error; err != nil {
		return nil, err
	}

	// Reload with relations
	if err := s.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Preload("Creator").
		First(&grade, gradeID).Error; err != nil {
		return nil, err
	}

	return toGradeResponse(&grade), nil
}

// DeleteGrade deletes a grade record
func (s *service) DeleteGrade(ctx context.Context, gradeID uint) error {
	result := s.db.WithContext(ctx).Delete(&models.Grade{}, gradeID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrGradeNotFound
	}
	return nil
}

// toGradeResponse converts a Grade model to GradeResponse
func toGradeResponse(g *models.Grade) *GradeResponse {
	response := &GradeResponse{
		ID:          g.ID,
		StudentID:   g.StudentID,
		Title:       g.Title,
		Score:       g.Score,
		Description: g.Description,
		CreatedBy:   g.CreatedBy,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}

	if g.Student.ID != 0 {
		response.StudentName = g.Student.Name
		response.StudentNIS = g.Student.NIS
	}

	return response
}

// ==================== Manual Attendance Methods ====================

// RecordManualAttendance records manual attendance for a student
func (s *service) RecordManualAttendance(ctx context.Context, schoolID, teacherID uint, req ManualAttendanceRequest) (*StudentAttendanceResponse, error) {
	// Validate required fields
	if req.StudentID == 0 {
		return nil, ErrStudentIDRequired
	}
	if req.ScheduleID == 0 {
		return nil, errors.New("schedule_id wajib diisi")
	}

	// Validate status
	status := models.AttendanceStatus(req.Status)
	if !status.IsValid() {
		return nil, ErrInvalidStatus
	}

	// Verify student belongs to the school
	student, err := s.repo.FindStudentByID(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if student.SchoolID != schoolID {
		return nil, ErrStudentNotInSchool
	}

	// Validate teacher has access to this student
	if err := s.ValidateTeacherAccess(ctx, teacherID, req.StudentID); err != nil {
		return nil, err
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("format tanggal tidak valid")
	}

	// Validate schedule exists and belongs to the school
	var schedule models.AttendanceSchedule
	if err := s.db.WithContext(ctx).
		Where("id = ? AND school_id = ? AND is_active = ?", req.ScheduleID, schoolID, true).
		First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("jadwal tidak ditemukan atau tidak aktif")
		}
		return nil, err
	}

	// Check if schedule is active on the given day
	if !schedule.IsActiveOnDay(date.Weekday()) {
		return nil, errors.New("jadwal tidak aktif pada hari tersebut")
	}

	// Check if attendance already exists for this student on this date and schedule
	var existingAttendance models.Attendance
	err = s.db.WithContext(ctx).
		Where("student_id = ? AND date = ? AND schedule_id = ?", req.StudentID, date, req.ScheduleID).
		First(&existingAttendance).Error
	if err == nil {
		return nil, ErrAttendanceAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create attendance record
	attendance := &models.Attendance{
		StudentID:  req.StudentID,
		ScheduleID: &req.ScheduleID,
		Date:       date,
		Status:     status,
		Method:     models.AttendanceMethodManual,
	}

	// Parse check-in time if provided
	if req.CheckInTime != "" {
		checkInTime, err := time.Parse("15:04", req.CheckInTime)
		if err == nil {
			// Combine date with time
			fullCheckIn := time.Date(date.Year(), date.Month(), date.Day(),
				checkInTime.Hour(), checkInTime.Minute(), 0, 0, date.Location())
			attendance.CheckInTime = &fullCheckIn
		}
	}

	// Parse check-out time if provided
	if req.CheckOutTime != "" {
		checkOutTime, err := time.Parse("15:04", req.CheckOutTime)
		if err == nil {
			// Combine date with time
			fullCheckOut := time.Date(date.Year(), date.Month(), date.Day(),
				checkOutTime.Hour(), checkOutTime.Minute(), 0, 0, date.Location())
			attendance.CheckOutTime = &fullCheckOut
		}
	}

	if err := s.db.WithContext(ctx).Create(attendance).Error; err != nil {
		return nil, err
	}

	// Build response
	resp := &StudentAttendanceResponse{
		ID:          attendance.ID,
		StudentID:   attendance.StudentID,
		StudentName: student.Name,
		StudentNIS:  student.NIS,
		Date:        req.Date,
		Status:      string(attendance.Status),
		Method:      string(attendance.Method),
		CreatedAt:   attendance.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   attendance.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if attendance.CheckInTime != nil {
		checkIn := attendance.CheckInTime.Format("15:04")
		resp.CheckInTime = &checkIn
	}
	if attendance.CheckOutTime != nil {
		checkOut := attendance.CheckOutTime.Format("15:04")
		resp.CheckOutTime = &checkOut
	}

	return resp, nil
}

// UpdateAttendance updates an existing attendance record
func (s *service) UpdateAttendance(ctx context.Context, teacherID, attendanceID uint, req UpdateAttendanceRequest) (*StudentAttendanceResponse, error) {
	// Get existing attendance
	var attendance models.Attendance
	if err := s.db.WithContext(ctx).
		Preload("Student").
		First(&attendance, attendanceID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAttendanceNotFound
		}
		return nil, err
	}

	// Validate teacher has access to this student
	if err := s.ValidateTeacherAccess(ctx, teacherID, attendance.StudentID); err != nil {
		return nil, err
	}

	// Update status if provided
	if req.Status != "" {
		status := models.AttendanceStatus(req.Status)
		if !status.IsValid() {
			return nil, ErrInvalidStatus
		}
		attendance.Status = status
	}

	// Update check-in time if provided
	if req.CheckInTime != "" {
		checkInTime, err := time.Parse("15:04", req.CheckInTime)
		if err == nil {
			fullCheckIn := time.Date(attendance.Date.Year(), attendance.Date.Month(), attendance.Date.Day(),
				checkInTime.Hour(), checkInTime.Minute(), 0, 0, attendance.Date.Location())
			attendance.CheckInTime = &fullCheckIn
		}
	}

	// Update check-out time if provided
	if req.CheckOutTime != "" {
		checkOutTime, err := time.Parse("15:04", req.CheckOutTime)
		if err == nil {
			fullCheckOut := time.Date(attendance.Date.Year(), attendance.Date.Month(), attendance.Date.Day(),
				checkOutTime.Hour(), checkOutTime.Minute(), 0, 0, attendance.Date.Location())
			attendance.CheckOutTime = &fullCheckOut
		}
	}

	// Mark as manual if updating
	attendance.Method = models.AttendanceMethodManual

	if err := s.db.WithContext(ctx).Save(&attendance).Error; err != nil {
		return nil, err
	}

	// Build response
	resp := &StudentAttendanceResponse{
		ID:          attendance.ID,
		StudentID:   attendance.StudentID,
		StudentName: attendance.Student.Name,
		StudentNIS:  attendance.Student.NIS,
		Date:        attendance.Date.Format("2006-01-02"),
		Status:      string(attendance.Status),
		Method:      string(attendance.Method),
		CreatedAt:   attendance.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   attendance.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if attendance.CheckInTime != nil {
		checkIn := attendance.CheckInTime.Format("15:04")
		resp.CheckInTime = &checkIn
	}
	if attendance.CheckOutTime != nil {
		checkOut := attendance.CheckOutTime.Format("15:04")
		resp.CheckOutTime = &checkOut
	}

	return resp, nil
}

// GetActiveSchedules retrieves active attendance schedules for a school on a specific date
func (s *service) GetActiveSchedules(ctx context.Context, schoolID uint, date string) ([]ScheduleResponse, error) {
	// Parse the date to get the day of week
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("format tanggal tidak valid")
	}

	// Get all active schedules for the school
	var schedules []models.AttendanceSchedule
	if err := s.db.WithContext(ctx).
		Where("school_id = ? AND is_active = ?", schoolID, true).
		Order("start_time ASC").
		Find(&schedules).Error; err != nil {
		return nil, err
	}

	// Filter schedules that are active on the given day
	var result []ScheduleResponse
	for _, schedule := range schedules {
		if schedule.IsActiveOnDay(parsedDate.Weekday()) {
			result = append(result, ScheduleResponse{
				ID:                schedule.ID,
				Name:              schedule.Name,
				StartTime:         schedule.StartTime,
				EndTime:           schedule.EndTime,
				LateThreshold:     schedule.LateThreshold,
				VeryLateThreshold: schedule.VeryLateThreshold,
				IsDefault:         schedule.IsDefault,
			})
		}
	}

	// If no schedules found for the day, return empty array
	if result == nil {
		result = []ScheduleResponse{}
	}

	return result, nil
}
