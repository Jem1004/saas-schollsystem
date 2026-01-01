package device

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/school-management/backend/internal/domain/models"
)

const (
	DefaultPairingDuration = 60 * time.Second // 60 seconds to tap the card
)

// pairingService implements the PairingService interface
type pairingService struct {
	deviceRepo     Repository
	studentRepo    StudentRepository
	pairingManager *PairingManager
}

// StudentRepository defines the interface for student operations needed by pairing service
type StudentRepository interface {
	FindByID(ctx context.Context, id uint) (*models.Student, error)
	FindByRFIDCode(ctx context.Context, rfidCode string) (*models.Student, error)
	UpdateRFIDCode(ctx context.Context, studentID uint, rfidCode string) error
}

// NewPairingService creates a new pairing service
func NewPairingService(deviceRepo Repository, studentRepo StudentRepository) PairingService {
	return &pairingService{
		deviceRepo:     deviceRepo,
		studentRepo:    studentRepo,
		pairingManager: NewPairingManager(),
	}
}

// StartPairing starts a new pairing session
func (s *pairingService) StartPairing(ctx context.Context, req StartPairingRequest) (*PairingSessionResponse, error) {
	// Validate device exists and is active
	device, err := s.deviceRepo.FindByID(ctx, req.DeviceID)
	if err != nil {
		return nil, err
	}
	if !device.IsActive {
		return nil, ErrDeviceInactive
	}

	// Validate student exists
	student, err := s.studentRepo.FindByID(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}

	// Check if student already has RFID card
	if student.RFIDCode != "" {
		return nil, ErrStudentAlreadyPaired
	}

	// Verify student belongs to the same school as the device
	if student.SchoolID != device.SchoolID {
		return nil, errors.New("siswa tidak terdaftar di sekolah yang sama dengan perangkat")
	}

	// Start pairing session
	session := s.pairingManager.StartPairingSession(
		device.ID,
		student.ID,
		student.Name,
		device.SchoolID,
		DefaultPairingDuration,
	)

	log.Printf("Pairing session started: device=%d, student=%d (%s), expires=%s",
		device.ID, student.ID, student.Name, session.ExpiresAt.Format("15:04:05"))

	return &PairingSessionResponse{
		Active:      true,
		StudentID:   session.StudentID,
		StudentName: session.StudentName,
		DeviceID:    session.DeviceID,
		ExpiresAt:   session.ExpiresAt,
		Message:     "Sesi pairing dimulai. Silakan tap kartu RFID pada perangkat dalam 60 detik.",
	}, nil
}

// ProcessRFIDPairing processes an RFID tap during pairing mode
func (s *pairingService) ProcessRFIDPairing(ctx context.Context, req RFIDPairingRequest) (*RFIDPairingResponse, error) {
	// Validate API key
	device, err := s.deviceRepo.FindByAPIKey(ctx, req.APIKey)
	if err != nil {
		return nil, ErrInvalidAPIKey
	}

	// Check for active pairing session
	session, err := s.pairingManager.GetPairingSession(device.ID)
	if err != nil {
		// No pairing session - this is a normal attendance tap, not pairing
		return &RFIDPairingResponse{
			Success: false,
			Message: "Tidak ada sesi pairing aktif untuk perangkat ini",
		}, nil
	}

	// Check if RFID code is already used by another student
	existingStudent, err := s.studentRepo.FindByRFIDCode(ctx, req.RFIDCode)
	if err == nil && existingStudent != nil {
		return &RFIDPairingResponse{
			Success: false,
			Message: "Kartu RFID sudah digunakan oleh siswa lain: " + existingStudent.Name,
		}, ErrRFIDAlreadyUsed
	}

	// Update student's RFID code
	if err := s.studentRepo.UpdateRFIDCode(ctx, session.StudentID, req.RFIDCode); err != nil {
		return nil, err
	}

	// Complete pairing session
	s.pairingManager.CompletePairingSession(device.ID)

	log.Printf("RFID pairing completed: student=%d (%s), rfid=%s",
		session.StudentID, session.StudentName, req.RFIDCode)

	return &RFIDPairingResponse{
		Success:     true,
		StudentID:   session.StudentID,
		StudentName: session.StudentName,
		RFIDCode:    req.RFIDCode,
		Message:     "Kartu RFID berhasil dipasangkan dengan siswa " + session.StudentName,
	}, nil
}

// CancelPairing cancels an active pairing session
func (s *pairingService) CancelPairing(ctx context.Context, deviceID uint) error {
	s.pairingManager.CancelPairingSession(deviceID)
	log.Printf("Pairing session cancelled: device=%d", deviceID)
	return nil
}

// GetPairingStatus gets the status of a pairing session
func (s *pairingService) GetPairingStatus(ctx context.Context, deviceID uint) (*PairingSessionResponse, error) {
	session, err := s.pairingManager.GetPairingSession(deviceID)
	if err != nil {
		return &PairingSessionResponse{
			Active:  false,
			Message: "Tidak ada sesi pairing aktif",
		}, nil
	}

	return &PairingSessionResponse{
		Active:      true,
		StudentID:   session.StudentID,
		StudentName: session.StudentName,
		DeviceID:    session.DeviceID,
		ExpiresAt:   session.ExpiresAt,
		Message:     "Sesi pairing aktif. Menunggu tap kartu RFID.",
	}, nil
}
