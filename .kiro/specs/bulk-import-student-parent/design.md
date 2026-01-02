# Design Document: Bulk Import Student & Parent

## Overview

Fitur Bulk Import memungkinkan Admin Sekolah mengimpor data siswa dan orang tua secara massal melalui file Excel (.xlsx). Sistem menggunakan pendekatan hybrid untuk relasi kelas (auto-match jika cocok, warning jika tidak) dan import terpisah untuk orang tua dengan linking manual via UI.

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Frontend (Vue.js)                           │
├─────────────────────────────────────────────────────────────────────┤
│  StudentManagement.vue    │    ParentManagement.vue                 │
│  - Download template      │    - Download template                  │
│  - Upload Excel           │    - Upload Excel                       │
│  - View import results    │    - View import results                │
│  - Bulk assign class      │    - Link to students                   │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         Backend (Go/Fiber)                          │
├─────────────────────────────────────────────────────────────────────┤
│  import/handler.go        │    import/service.go                    │
│  - POST /import/students  │    - ImportStudents()                   │
│  - POST /import/parents   │    - ImportParents()                    │
│  - GET /import/template/* │    - GenerateTemplate()                 │
├─────────────────────────────────────────────────────────────────────┤
│  import/parser.go         │    import/validator.go                  │
│  - ParseExcelFile()       │    - ValidateStudentRow()               │
│  - ReadRows()             │    - ValidateParentRow()                │
├─────────────────────────────────────────────────────────────────────┤
│  school/service.go (existing)                                       │
│  - BulkAssignClass()      │    - LinkParentToStudents()             │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         Database (PostgreSQL)                       │
├─────────────────────────────────────────────────────────────────────┤
│  students (modified)      │    parents                              │
│  - class_id nullable      │    - unchanged                          │
│  - is_active logic        │                                         │
├─────────────────────────────────────────────────────────────────────┤
│  student_parents          │    users                                │
│  - many-to-many link      │    - parent accounts                    │
└─────────────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### 1. Import Module (New)

```go
// backend/internal/modules/import/service.go

type ImportService interface {
    // Template generation
    GenerateStudentTemplate() ([]byte, error)
    GenerateParentTemplate() ([]byte, error)
    
    // Import operations
    ImportStudents(ctx context.Context, schoolID uint, file multipart.File) (*ImportResult, error)
    ImportParents(ctx context.Context, schoolID uint, file multipart.File) (*ImportResult, error)
}

type importService struct {
    schoolRepo   school.Repository
    userRepo     school.UserRepository
    excelParser  ExcelParser
}
```

### 2. Excel Parser Component

```go
// backend/internal/modules/import/parser.go

type ExcelParser interface {
    ParseStudentFile(file multipart.File) ([]StudentRow, error)
    ParseParentFile(file multipart.File) ([]ParentRow, error)
}

type StudentRow struct {
    RowNumber int
    NIS       string
    NISN      string
    Name      string
    ClassName string  // Optional
}

type ParentRow struct {
    RowNumber int
    Name      string
    Phone     string
    Email     string
}
```

### 3. Class Matcher Component

```go
// backend/internal/modules/import/matcher.go

type ClassMatcher interface {
    FindClassByName(ctx context.Context, schoolID uint, className string) (*models.Class, error)
}
```

### 4. Import Handler

```go
// backend/internal/modules/import/handler.go

type ImportHandler struct {
    service ImportService
}

// Routes:
// GET  /api/v1/school/import/template/students  - Download student template
// GET  /api/v1/school/import/template/parents   - Download parent template
// POST /api/v1/school/import/students           - Import students from Excel
// POST /api/v1/school/import/parents            - Import parents from Excel
```

### 5. Extended School Service

```go
// backend/internal/modules/school/service.go (additions)

type Service interface {
    // ... existing methods ...
    
    // Bulk operations for import
    BulkAssignClass(ctx context.Context, schoolID uint, studentIDs []uint, classID uint) error
    GetStudentsWithoutClass(ctx context.Context, schoolID uint) ([]StudentResponse, error)
}
```

## Data Models

### Modified Student Model

```go
// backend/internal/domain/models/student.go

type Student struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    SchoolID  uint      `gorm:"index;not null" json:"school_id"`
    ClassID   *uint     `gorm:"index" json:"class_id"`  // Changed: nullable
    UserID    *uint     `gorm:"uniqueIndex" json:"user_id"`
    NIS       string    `gorm:"type:varchar(20);not null" json:"nis"`
    NISN      string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"nisn"`
    Name      string    `gorm:"type:varchar(255);not null" json:"name"`
    RFIDCode  string    `gorm:"type:varchar(50);index" json:"rfid_code"`
    IsActive  bool      `gorm:"default:false" json:"is_active"`  // Changed: default false
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    // Relations
    School  School   `gorm:"foreignKey:SchoolID" json:"school,omitempty"`
    Class   *Class   `gorm:"foreignKey:ClassID" json:"class,omitempty"`  // Changed: pointer
    User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Parents []Parent `gorm:"many2many:student_parents" json:"parents,omitempty"`
}

// Validate validates the student data (modified for import)
func (s *Student) Validate() error {
    if strings.TrimSpace(s.NIS) == "" {
        return errors.New("NIS is required")
    }
    if strings.TrimSpace(s.NISN) == "" {
        return errors.New("NISN is required")
    }
    if strings.TrimSpace(s.Name) == "" {
        return errors.New("name is required")
    }
    if s.SchoolID == 0 {
        return errors.New("school_id is required")
    }
    // ClassID validation removed - now optional
    return nil
}

// CanBeActive checks if student can be set to active
func (s *Student) CanBeActive() bool {
    return s.ClassID != nil && *s.ClassID > 0
}
```

### Import DTOs

```go
// backend/internal/modules/import/dto.go

// ImportResult represents the result of an import operation
type ImportResult struct {
    TotalRows           int              `json:"total_rows"`
    SuccessCount        int              `json:"success_count"`
    FailedCount         int              `json:"failed_count"`
    WarningCount        int              `json:"warning_count"`
    Errors              []ImportError    `json:"errors,omitempty"`
    Warnings            []ImportWarning  `json:"warnings,omitempty"`
    StudentsWithoutClass int             `json:"students_without_class,omitempty"`
}

// ImportError represents an error during import
type ImportError struct {
    Row     int    `json:"row"`
    Field   string `json:"field"`
    Message string `json:"message"`
}

// ImportWarning represents a warning during import
type ImportWarning struct {
    Row     int    `json:"row"`
    Field   string `json:"field"`
    Message string `json:"message"`
}

// BulkAssignClassRequest represents request to assign class to multiple students
type BulkAssignClassRequest struct {
    StudentIDs []uint `json:"student_ids" validate:"required,min=1"`
    ClassID    uint   `json:"class_id" validate:"required"`
}
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Required Field Validation

*For any* student import row, if any required field (NIS, NISN, or Name) is empty or whitespace-only, the row should be rejected and added to the error list with appropriate error message indicating which field is missing.

**Validates: Requirements 3.1, 3.2, 3.3**

### Property 2: Duplicate Detection

*For any* student import row where NISN or NIS already exists in the school's database, the import should skip that row, add it to the warning list with appropriate message, and leave the existing student record unchanged.

**Validates: Requirements 3.4, 3.5**

### Property 3: Class Matching Correctness

*For any* student import row with non-empty class_name, if the class_name matches an existing class in the school (case-insensitive), the imported student should have ClassID set to that class's ID.

**Validates: Requirements 3.6**

### Property 4: ClassID-IsActive Invariant

*For any* student record, IsActive should be true if and only if ClassID is not null. This invariant must hold after import, after bulk class assignment, and after any update operation.

**Validates: Requirements 3.9, 3.10, 6.4, 8.2, 8.3**

### Property 5: Tenant Isolation

*For any* import or linking operation performed by an admin, all created or modified records (students, parents, user accounts, student-parent links) should have SchoolID equal to the authenticated admin's school ID.

**Validates: Requirements 3.11, 4.6, 6.5, 7.6**

### Property 6: Parent Account Creation Consistency

*For any* valid parent import row, the system should create exactly one Parent record and exactly one User record, where the User has role "orang_tua" and username equals the parent's phone number.

**Validates: Requirements 4.4, 4.5**

### Property 7: Import Result Accuracy

*For any* import operation with N total rows, the returned ImportResult should satisfy: TotalRows == N, and SuccessCount + FailedCount + WarningCount (skipped duplicates) == N.

**Validates: Requirements 5.1, 5.2, 5.3, 5.4**

### Property 8: Parent-Student Many-to-Many Relationship

*For any* parent-student link operation, the system should correctly maintain many-to-many relationships: one parent can be linked to multiple students, and one student can be linked to multiple parents, with each link creating exactly one entry in student_parents table.

**Validates: Requirements 7.3, 7.4, 7.5**

### Property 9: Bulk Class Assignment Completeness

*For any* bulk class assignment operation with a list of student IDs and a valid class ID, all students in the list should have their ClassID updated to the specified value, and their IsActive set to true.

**Validates: Requirements 6.3, 6.4**

### Property 10: Attendance Restriction for Unassigned Students

*For any* student with null ClassID, the system should prevent attendance recording for that student.

**Validates: Requirements 8.4**

## Error Handling

### File Validation Errors

| Error Code | Condition | Response |
|------------|-----------|----------|
| `INVALID_FILE_FORMAT` | File is not .xlsx | 400 Bad Request |
| `FILE_TOO_LARGE` | File > 5MB | 400 Bad Request |
| `EMPTY_FILE` | No data rows | 400 Bad Request |
| `INVALID_HEADER` | Missing required columns | 400 Bad Request |

### Row-Level Errors

| Error Code | Condition | Handling |
|------------|-----------|----------|
| `MISSING_NIS` | NIS column empty | Skip row, add to errors |
| `MISSING_NISN` | NISN column empty | Skip row, add to errors |
| `MISSING_NAME` | Name column empty | Skip row, add to errors |
| `DUPLICATE_NISN` | NISN exists in school | Skip row, add to warnings |
| `DUPLICATE_NIS` | NIS exists in school | Skip row, add to warnings |
| `INVALID_EMAIL` | Email format invalid | Skip row, add to errors |
| `DUPLICATE_PHONE` | Phone exists in school | Skip row, add to warnings |

### Row-Level Warnings

| Warning Code | Condition | Handling |
|--------------|-----------|----------|
| `CLASS_NOT_FOUND` | class_name doesn't match | Import with null ClassID |
| `EMPTY_CLASS` | class_name column empty | Import with null ClassID |

## Testing Strategy

### Unit Tests

Unit tests will verify specific examples and edge cases:

1. **Template Generation Tests**
   - Student template has correct columns (NIS, NISN, Nama, Kelas)
   - Parent template has correct columns (Nama, No_HP, Email)
   - Templates include header row and example data row

2. **Parser Tests**
   - Parse valid Excel file with all columns
   - Parse Excel file with missing optional columns (class_name)
   - Handle empty rows gracefully
   - Handle special characters in names
   - Reject non-.xlsx files
   - Reject files larger than 5MB
   - Reject empty files

3. **Validator Tests**
   - Validate email format for parents
   - Handle edge cases (empty strings, whitespace-only)

4. **Class Matcher Tests**
   - Match exact class name
   - Handle case-insensitive matching
   - Return nil for non-existent class

### Property-Based Tests

Property-based tests will verify universal properties across many generated inputs using the `gopter` library for Go:

1. **Required Field Validation Property** (Property 1)
   - Generate student rows with random empty required fields
   - Verify all rows with empty required fields are rejected
   - **Feature: bulk-import-student-parent, Property 1: Required Field Validation**

2. **Duplicate Detection Property** (Property 2)
   - Generate data with known duplicate NISN/NIS
   - Import and verify duplicates are skipped with warnings
   - **Feature: bulk-import-student-parent, Property 2: Duplicate Detection**

3. **Class Matching Property** (Property 3)
   - Generate student rows with various class names
   - Verify correct ClassID assignment for matching names
   - **Feature: bulk-import-student-parent, Property 3: Class Matching Correctness**

4. **ClassID-IsActive Invariant Property** (Property 4)
   - Generate students with various ClassID states
   - Verify IsActive is always consistent with ClassID presence
   - **Feature: bulk-import-student-parent, Property 4: ClassID-IsActive Invariant**

5. **Tenant Isolation Property** (Property 5)
   - Generate import data for multiple schools
   - Verify all records have correct SchoolID
   - **Feature: bulk-import-student-parent, Property 5: Tenant Isolation**

6. **Parent Account Creation Property** (Property 6)
   - Generate valid parent rows
   - Verify Parent and User records are created correctly
   - **Feature: bulk-import-student-parent, Property 6: Parent Account Creation Consistency**

7. **Import Result Accuracy Property** (Property 7)
   - Generate import data with mix of valid, invalid, and duplicate rows
   - Verify result counts are accurate
   - **Feature: bulk-import-student-parent, Property 7: Import Result Accuracy**

8. **Parent-Student Linking Property** (Property 8)
   - Generate random parent-student link operations
   - Verify many-to-many relationships are maintained correctly
   - **Feature: bulk-import-student-parent, Property 8: Parent-Student Many-to-Many Relationship**

9. **Bulk Class Assignment Property** (Property 9)
   - Generate list of students without class
   - Perform bulk assignment
   - Verify all students have correct ClassID and IsActive=true
   - **Feature: bulk-import-student-parent, Property 9: Bulk Class Assignment Completeness**

### Integration Tests

1. **Full Student Import Flow**
   - Upload file → Parse → Validate → Import → Verify DB state
   - Test with mix of valid rows, invalid rows, and class matching scenarios

2. **Full Parent Import Flow**
   - Upload file → Parse → Validate → Import → Verify DB state and User accounts

3. **Bulk Class Assignment Flow**
   - Import students without class → Bulk assign → Verify activation

4. **Parent Linking Flow**
   - Import parents → Link to students → Verify relationships

### Test Configuration

- Property tests: minimum 100 iterations per property
- Use `gopter` library for property-based testing in Go
- Tag format: **Feature: bulk-import-student-parent, Property {number}: {property_text}**
- Each property test must reference its design document property number
