package models

import (
	"errors"
	"strings"
	"time"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleSuperAdmin   UserRole = "super_admin"
	RoleAdminSekolah UserRole = "admin_sekolah"
	RoleGuruBK       UserRole = "guru_bk"
	RoleWaliKelas    UserRole = "wali_kelas"
	RoleGuru         UserRole = "guru"
	RoleParent       UserRole = "parent"
	RoleStudent      UserRole = "student"
)

// IsValid checks if the user role is valid
func (r UserRole) IsValid() bool {
	switch r {
	case RoleSuperAdmin, RoleAdminSekolah, RoleGuruBK, RoleWaliKelas, RoleGuru, RoleParent, RoleStudent:
		return true
	}
	return false
}

// User represents all system users
type User struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	SchoolID     *uint      `gorm:"index" json:"school_id"` // nil for super_admin
	Role         UserRole   `gorm:"type:varchar(20);not null" json:"role"`
	Username     string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	PasswordHash string     `gorm:"type:varchar(255);not null" json:"-"`
	Email        string     `gorm:"type:varchar(255)" json:"email"`
	Name         string     `gorm:"type:varchar(255)" json:"name"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	MustResetPwd bool       `gorm:"default:true" json:"must_reset_pwd"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// Relations
	School *School `gorm:"foreignKey:SchoolID" json:"school,omitempty"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// Validate validates the user data
// Requirements: 4.1 - User account creation SHALL assign a role and generate initial credentials
func (u *User) Validate() error {
	if strings.TrimSpace(u.Username) == "" {
		return errors.New("username wajib diisi")
	}
	if !u.Role.IsValid() {
		return errors.New("role tidak valid")
	}
	// Super admin doesn't need school_id
	if u.Role != RoleSuperAdmin && u.SchoolID == nil {
		return errors.New("ID sekolah wajib diisi untuk user non-super_admin")
	}
	return nil
}

// IsSuperAdmin checks if the user is a super admin
func (u *User) IsSuperAdmin() bool {
	return u.Role == RoleSuperAdmin
}

// IsAdminSekolah checks if the user is a school admin
func (u *User) IsAdminSekolah() bool {
	return u.Role == RoleAdminSekolah
}

// IsGuruBK checks if the user is a counseling teacher
func (u *User) IsGuruBK() bool {
	return u.Role == RoleGuruBK
}

// IsWaliKelas checks if the user is a homeroom teacher
func (u *User) IsWaliKelas() bool {
	return u.Role == RoleWaliKelas
}

// IsGuru checks if the user is a teacher
func (u *User) IsGuru() bool {
	return u.Role == RoleGuru
}

// IsParent checks if the user is a parent
func (u *User) IsParent() bool {
	return u.Role == RoleParent
}

// IsStudent checks if the user is a student
func (u *User) IsStudent() bool {
	return u.Role == RoleStudent
}

// Deactivate deactivates the user
// Requirements: 4.4 - Deactivating an account SHALL revoke all active sessions
func (u *User) Deactivate() {
	u.IsActive = false
}

// Activate activates the user
func (u *User) Activate() {
	u.IsActive = true
}

// UpdateLastLogin updates the last login timestamp
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
}

// MarkPasswordReset marks that the user must reset their password
func (u *User) MarkPasswordReset() {
	u.MustResetPwd = true
}

// ClearPasswordReset clears the password reset flag
func (u *User) ClearPasswordReset() {
	u.MustResetPwd = false
}
