package models

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeAttendanceIn  NotificationType = "attendance_in"
	NotificationTypeAttendanceOut NotificationType = "attendance_out"
	NotificationTypeViolation     NotificationType = "violation"
	NotificationTypeAchievement   NotificationType = "achievement"
	NotificationTypePermit        NotificationType = "permit"
	NotificationTypeCounseling    NotificationType = "counseling"
	NotificationTypeGrade         NotificationType = "grade"
	NotificationTypeHomeroomNote  NotificationType = "homeroom_note"
)

// IsValid checks if the notification type is valid
func (t NotificationType) IsValid() bool {
	switch t {
	case NotificationTypeAttendanceIn, NotificationTypeAttendanceOut,
		NotificationTypeViolation, NotificationTypeAchievement,
		NotificationTypePermit, NotificationTypeCounseling,
		NotificationTypeGrade, NotificationTypeHomeroomNote:
		return true
	}
	return false
}

// Notification represents user notification
type Notification struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	UserID    uint             `gorm:"index;not null" json:"user_id"`
	Type      NotificationType `gorm:"type:varchar(50);not null" json:"type"`
	Title     string           `gorm:"type:varchar(255);not null" json:"title"`
	Message   string           `gorm:"type:text;not null" json:"message"`
	Data      string           `gorm:"type:jsonb" json:"data"` // Additional JSON data
	IsRead    bool             `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time        `json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for Notification
func (Notification) TableName() string {
	return "notifications"
}

// Validate validates the notification data
// Requirements: 17.3 - Notification SHALL store notification history in database
func (n *Notification) Validate() error {
	if n.UserID == 0 {
		return errors.New("ID user wajib diisi")
	}
	if !n.Type.IsValid() {
		return errors.New("tipe notifikasi tidak valid")
	}
	if strings.TrimSpace(n.Title) == "" {
		return errors.New("judul wajib diisi")
	}
	if strings.TrimSpace(n.Message) == "" {
		return errors.New("pesan wajib diisi")
	}
	return nil
}

// MarkAsRead marks the notification as read
func (n *Notification) MarkAsRead() {
	n.IsRead = true
}

// SetData sets the additional JSON data
func (n *Notification) SetData(data map[string]interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	n.Data = string(jsonData)
	return nil
}

// GetData retrieves the additional JSON data
func (n *Notification) GetData() (map[string]interface{}, error) {
	if n.Data == "" {
		return nil, nil
	}
	var data map[string]interface{}
	err := json.Unmarshal([]byte(n.Data), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// FCMToken represents user's FCM device token
type FCMToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Token     string    `gorm:"type:varchar(500);not null" json:"token"`
	Platform  string    `gorm:"type:varchar(20);not null" json:"platform"` // android, ios
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for FCMToken
func (FCMToken) TableName() string {
	return "fcm_tokens"
}

// Validate validates the FCM token data
func (f *FCMToken) Validate() error {
	if f.UserID == 0 {
		return errors.New("ID user wajib diisi")
	}
	if strings.TrimSpace(f.Token) == "" {
		return errors.New("token wajib diisi")
	}
	if f.Platform != "android" && f.Platform != "ios" {
		return errors.New("platform harus android atau ios")
	}
	return nil
}

// Deactivate deactivates the FCM token
func (f *FCMToken) Deactivate() {
	f.IsActive = false
}
