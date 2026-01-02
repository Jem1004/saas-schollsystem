package importmodule

import (
	"context"
	"strings"

	"github.com/school-management/backend/internal/domain/models"
	"gorm.io/gorm"
)

// ClassMatcher defines the interface for matching class names
// Requirements: 3.6, 3.7
type ClassMatcher interface {
	FindClassByName(ctx context.Context, schoolID uint, className string) (*models.Class, error)
}

// classMatcher implements ClassMatcher interface
type classMatcher struct {
	db *gorm.DB
}

// NewClassMatcher creates a new class matcher
func NewClassMatcher(db *gorm.DB) ClassMatcher {
	return &classMatcher{db: db}
}

// FindClassByName finds a class by name (case-insensitive)
// Requirements: 3.6 - Match class by name case-insensitively
// Returns nil if class not found (no error)
func (m *classMatcher) FindClassByName(ctx context.Context, schoolID uint, className string) (*models.Class, error) {
	if strings.TrimSpace(className) == "" {
		return nil, nil
	}

	var class models.Class
	err := m.db.WithContext(ctx).
		Where("school_id = ? AND LOWER(name) = LOWER(?)", schoolID, strings.TrimSpace(className)).
		First(&class).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Class not found is not an error
		}
		return nil, err
	}

	return &class, nil
}
