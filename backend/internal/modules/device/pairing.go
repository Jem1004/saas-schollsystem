package device

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrNoPairingSession    = errors.New("tidak ada sesi pairing aktif")
	ErrPairingSessionExpired = errors.New("sesi pairing sudah kadaluarsa")
	ErrStudentAlreadyPaired = errors.New("siswa sudah memiliki kartu RFID")
	ErrRFIDAlreadyUsed     = errors.New("kartu RFID sudah digunakan siswa lain")
)

// PairingSession represents an active RFID pairing session
type PairingSession struct {
	StudentID   uint      `json:"student_id"`
	StudentName string    `json:"student_name"`
	SchoolID    uint      `json:"school_id"`
	DeviceID    uint      `json:"device_id"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// PairingManager manages RFID pairing sessions
type PairingManager struct {
	sessions map[uint]*PairingSession // key: device_id
	mu       sync.RWMutex
}

// NewPairingManager creates a new pairing manager
func NewPairingManager() *PairingManager {
	pm := &PairingManager{
		sessions: make(map[uint]*PairingSession),
	}
	// Start cleanup goroutine
	go pm.cleanupExpiredSessions()
	return pm
}

// StartPairingSession starts a new pairing session for a device
func (pm *PairingManager) StartPairingSession(deviceID, studentID uint, studentName string, schoolID uint, duration time.Duration) *PairingSession {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	session := &PairingSession{
		StudentID:   studentID,
		StudentName: studentName,
		SchoolID:    schoolID,
		DeviceID:    deviceID,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(duration),
	}

	pm.sessions[deviceID] = session
	return session
}

// GetPairingSession gets the active pairing session for a device
func (pm *PairingManager) GetPairingSession(deviceID uint) (*PairingSession, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	session, exists := pm.sessions[deviceID]
	if !exists {
		return nil, ErrNoPairingSession
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, ErrPairingSessionExpired
	}

	return session, nil
}

// CompletePairingSession completes and removes a pairing session
func (pm *PairingManager) CompletePairingSession(deviceID uint) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.sessions, deviceID)
}

// CancelPairingSession cancels a pairing session
func (pm *PairingManager) CancelPairingSession(deviceID uint) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.sessions, deviceID)
}

// GetAllActiveSessions returns all active pairing sessions
func (pm *PairingManager) GetAllActiveSessions() []*PairingSession {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	sessions := make([]*PairingSession, 0)
	now := time.Now()
	for _, session := range pm.sessions {
		if now.Before(session.ExpiresAt) {
			sessions = append(sessions, session)
		}
	}
	return sessions
}

// cleanupExpiredSessions periodically removes expired sessions
func (pm *PairingManager) cleanupExpiredSessions() {
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		pm.mu.Lock()
		now := time.Now()
		for deviceID, session := range pm.sessions {
			if now.After(session.ExpiresAt) {
				delete(pm.sessions, deviceID)
			}
		}
		pm.mu.Unlock()
	}
}

// PairingService handles RFID card pairing operations
type PairingService interface {
	StartPairing(ctx context.Context, req StartPairingRequest) (*PairingSessionResponse, error)
	ProcessRFIDPairing(ctx context.Context, req RFIDPairingRequest) (*RFIDPairingResponse, error)
	CancelPairing(ctx context.Context, deviceID uint) error
	GetPairingStatus(ctx context.Context, deviceID uint) (*PairingSessionResponse, error)
}

// StartPairingRequest represents the request to start a pairing session
type StartPairingRequest struct {
	DeviceID  uint `json:"device_id" validate:"required"`
	StudentID uint `json:"student_id" validate:"required"`
}

// RFIDPairingRequest represents the request from ESP32 during pairing mode
type RFIDPairingRequest struct {
	APIKey   string `json:"api_key" validate:"required"`
	RFIDCode string `json:"rfid_code" validate:"required"`
}

// PairingSessionResponse represents the pairing session status
type PairingSessionResponse struct {
	Active      bool      `json:"active"`
	StudentID   uint      `json:"student_id,omitempty"`
	StudentName string    `json:"student_name,omitempty"`
	DeviceID    uint      `json:"device_id,omitempty"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
	Message     string    `json:"message"`
}

// RFIDPairingResponse represents the response after RFID pairing
type RFIDPairingResponse struct {
	Success     bool   `json:"success"`
	StudentID   uint   `json:"student_id,omitempty"`
	StudentName string `json:"student_name,omitempty"`
	RFIDCode    string `json:"rfid_code,omitempty"`
	Message     string `json:"message"`
}
