package models

import (
	"errors"
	"strings"
	"time"
)

// ViolationCategory represents a customizable violation category per school
type ViolationCategory struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	SchoolID     uint           `gorm:"index;not null" json:"school_id"`
	Name         string         `gorm:"type:varchar(100);not null" json:"name"`
	DefaultPoint int            `gorm:"not null;default:-5" json:"default_point"`
	DefaultLevel ViolationLevel `gorm:"type:varchar(20);not null;default:'ringan'" json:"default_level"`
	Description  string         `gorm:"type:text" json:"description"`
	IsActive     bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`

	// Relations
	School School `gorm:"foreignKey:SchoolID" json:"school,omitempty"`
}

// TableName specifies the table name for ViolationCategory
func (ViolationCategory) TableName() string {
	return "violation_categories"
}

// Validate validates the violation category data
func (vc *ViolationCategory) Validate() error {
	if vc.SchoolID == 0 {
		return errors.New("school_id is required")
	}
	if strings.TrimSpace(vc.Name) == "" {
		return errors.New("name is required")
	}
	if vc.DefaultPoint > 0 {
		return errors.New("default_point must be zero or negative for violations")
	}
	if !vc.DefaultLevel.IsValid() {
		return errors.New("default_level must be one of: ringan, sedang, berat")
	}
	return nil
}

// DefaultViolationCategories returns the default violation categories for a new school
func DefaultViolationCategories(schoolID uint) []ViolationCategory {
	return []ViolationCategory{
		{SchoolID: schoolID, Name: "Keterlambatan", DefaultPoint: -5, DefaultLevel: ViolationLevelRingan, Description: "Terlambat masuk sekolah", IsActive: true},
		{SchoolID: schoolID, Name: "Bolos", DefaultPoint: -15, DefaultLevel: ViolationLevelSedang, Description: "Tidak masuk tanpa keterangan", IsActive: true},
		{SchoolID: schoolID, Name: "Seragam", DefaultPoint: -5, DefaultLevel: ViolationLevelRingan, Description: "Tidak memakai seragam sesuai aturan", IsActive: true},
		{SchoolID: schoolID, Name: "Perilaku", DefaultPoint: -10, DefaultLevel: ViolationLevelSedang, Description: "Perilaku tidak sopan", IsActive: true},
		{SchoolID: schoolID, Name: "Kekerasan", DefaultPoint: -30, DefaultLevel: ViolationLevelBerat, Description: "Melakukan kekerasan fisik", IsActive: true},
		{SchoolID: schoolID, Name: "Bullying", DefaultPoint: -25, DefaultLevel: ViolationLevelBerat, Description: "Melakukan perundungan", IsActive: true},
		{SchoolID: schoolID, Name: "Merokok", DefaultPoint: -20, DefaultLevel: ViolationLevelBerat, Description: "Merokok di lingkungan sekolah", IsActive: true},
		{SchoolID: schoolID, Name: "Narkoba", DefaultPoint: -50, DefaultLevel: ViolationLevelBerat, Description: "Terlibat narkoba", IsActive: true},
		{SchoolID: schoolID, Name: "Pencurian", DefaultPoint: -30, DefaultLevel: ViolationLevelBerat, Description: "Melakukan pencurian", IsActive: true},
		{SchoolID: schoolID, Name: "Vandalisme", DefaultPoint: -20, DefaultLevel: ViolationLevelSedang, Description: "Merusak fasilitas sekolah", IsActive: true},
		{SchoolID: schoolID, Name: "Lainnya", DefaultPoint: -5, DefaultLevel: ViolationLevelRingan, Description: "Pelanggaran lainnya", IsActive: true},
	}
}
