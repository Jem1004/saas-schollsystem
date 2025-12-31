# Design Document

## Overview

Sistem SaaS Manajemen Sekolah adalah aplikasi multi-tenant yang terdiri dari tiga komponen utama:
1. **Backend API** - Golang dengan Fiber framework, arsitektur Modular Monolith
2. **Web Admin** - Vue 3 + TypeScript + Ant Design Vue untuk admin sekolah
3. **Mobile App** - Flutter untuk orang tua dan siswa

Sistem ini menggunakan PostgreSQL sebagai database utama, Redis untuk queue dan cache, serta Firebase Cloud Messaging (FCM) untuk push notification.

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              CLIENTS                                         │
├─────────────────┬─────────────────┬─────────────────┬───────────────────────┤
│   Web Admin     │   Mobile App    │   Mobile App    │      ESP32/RFID       │
│   (Vue 3)       │   (Flutter)     │   (Flutter)     │      Devices          │
│   Super Admin   │   Orang Tua     │   Siswa         │                       │
│   Admin Sekolah │                 │                 │                       │
│   Guru BK       │                 │                 │                       │
│   Wali Kelas    │                 │                 │                       │
└────────┬────────┴────────┬────────┴────────┬────────┴───────────┬───────────┘
         │                 │                 │                     │
         │            HTTPS/REST API         │                     │
         └─────────────────┼─────────────────┴─────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         BACKEND (Golang + Fiber)                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                        API Gateway Layer                             │    │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐               │    │
│  │  │ Auth         │  │ Rate Limiter │  │ Tenant       │               │    │
│  │  │ Middleware   │  │              │  │ Middleware   │               │    │
│  │  └──────────────┘  └──────────────┘  └──────────────┘               │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                        Domain Modules                                │    │
│  │  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────┐        │    │
│  │  │ Auth       │ │ Tenant     │ │ School     │ │ Attendance │        │    │
│  │  │ Module     │ │ Module     │ │ Module     │ │ Module     │        │    │
│  │  └────────────┘ └────────────┘ └────────────┘ └────────────┘        │    │
│  │  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────┐        │    │
│  │  │ BK         │ │ Grade      │ │ Homeroom   │ │ Notification│       │    │
│  │  │ Module     │ │ Module     │ │ Module     │ │ Module     │        │    │
│  │  └────────────┘ └────────────┘ └────────────┘ └────────────┘        │    │
│  │  ┌────────────┐                                                      │    │
│  │  │ Device     │                                                      │    │
│  │  │ Module     │                                                      │    │
│  │  └────────────┘                                                      │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                     Background Workers                               │    │
│  │  ┌──────────────────┐  ┌──────────────────┐                         │    │
│  │  │ Notification     │  │ Document         │                         │    │
│  │  │ Worker           │  │ Generator        │                         │    │
│  │  └──────────────────┘  └──────────────────┘                         │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                 │
         ▼                 ▼                 ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│   PostgreSQL    │ │     Redis       │ │      FCM        │
│   (Primary DB)  │ │ (Queue/Cache)   │ │ (Push Notif)    │
└─────────────────┘ └─────────────────┘ └─────────────────┘
```

### Modular Monolith Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── tenant.go
│   │   └── rbac.go
│   ├── policy/                    # Policy Layer
│   │   ├── attendance_policy.go
│   │   ├── notification_policy.go
│   │   └── access_policy.go
│   ├── modules/
│   │   ├── auth/
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── dto.go
│   │   ├── tenant/
│   │   ├── school/
│   │   ├── attendance/
│   │   ├── bk/
│   │   ├── grade/
│   │   ├── homeroom/
│   │   ├── notification/
│   │   ├── device/
│   │   └── settings/              # School Settings Module
│   ├── outbox/                    # Event Outbox Pattern
│   │   ├── outbox.go
│   │   ├── publisher.go
│   │   └── worker.go
│   ├── shared/
│   │   ├── database/
│   │   ├── redis/
│   │   ├── fcm/
│   │   └── utils/
│   └── domain/
│       └── models/
├── pkg/
│   └── validator/
└── go.mod
```

### Policy Layer

Policy Layer memisahkan business rules dari service logic, memudahkan perubahan aturan tanpa mengubah core logic.

```go
// AttendancePolicy determines attendance rules based on school settings
type AttendancePolicy interface {
    IsWithinAttendanceWindow(schoolID uint, timestamp time.Time) bool
    DetermineAttendanceStatus(checkInTime time.Time, settings *SchoolSettings) string // on_time, late, very_late
    ShouldSendNotification(settings *SchoolSettings, eventType string) bool
}

// NotificationPolicy determines when and how to send notifications
type NotificationPolicy interface {
    ShouldNotifyParent(settings *SchoolSettings, eventType string) bool
    GetNotificationPriority(eventType string) string
    FormatNotificationMessage(eventType string, data map[string]interface{}) string
}

// AccessPolicy determines access rules beyond basic RBAC
type AccessPolicy interface {
    CanAccessStudentData(user *User, studentID uint) bool
    CanModifyGrade(user *User, gradeID uint) bool
    GetVisibleFields(user *User, resourceType string) []string
}
```

### Event Outbox Pattern (Optional)

Event Outbox ensures reliable event publishing with transactional guarantees.

```go
// OutboxEvent represents an event to be published
type OutboxEvent struct {
    ID          uint      `gorm:"primaryKey"`
    AggregateID uint      `gorm:"index;not null"`
    EventType   string    `gorm:"type:varchar(100);not null"`
    Payload     string    `gorm:"type:jsonb;not null"`
    Status      string    `gorm:"type:varchar(20);default:'pending'"` // pending, published, failed
    RetryCount  int       `gorm:"default:0"`
    CreatedAt   time.Time
    PublishedAt *time.Time
}

// OutboxPublisher interface
type OutboxPublisher interface {
    Publish(event *OutboxEvent) error
    MarkAsPublished(eventID uint) error
    MarkAsFailed(eventID uint) error
    GetPendingEvents(limit int) ([]OutboxEvent, error)
}

// OutboxWorker processes pending events
type OutboxWorker interface {
    Start()
    Stop()
    ProcessPendingEvents()
}
```

## Components and Interfaces

### 1. Auth Module

```go
// Handler Interface
type AuthHandler interface {
    Login(c *fiber.Ctx) error
    RefreshToken(c *fiber.Ctx) error
    Logout(c *fiber.Ctx) error
    ChangePassword(c *fiber.Ctx) error
}

// Service Interface
type AuthService interface {
    Authenticate(username, password string) (*TokenPair, error)
    RefreshAccessToken(refreshToken string) (*TokenPair, error)
    RevokeToken(userID uint) error
    ChangePassword(userID uint, oldPassword, newPassword string) error
}

// Repository Interface
type UserRepository interface {
    FindByUsername(username string) (*User, error)
    FindByID(id uint) (*User, error)
    UpdatePassword(id uint, passwordHash string) error
    UpdateLastLogin(id uint) error
}
```

### 2. Tenant Module

```go
// Handler Interface
type TenantHandler interface {
    CreateSchool(c *fiber.Ctx) error
    GetSchools(c *fiber.Ctx) error
    GetSchool(c *fiber.Ctx) error
    UpdateSchool(c *fiber.Ctx) error
    DeactivateSchool(c *fiber.Ctx) error
}

// Service Interface
type TenantService interface {
    CreateSchool(req CreateSchoolRequest) (*School, error)
    GetAllSchools(pagination Pagination) ([]School, int64, error)
    GetSchoolByID(id uint) (*School, error)
    UpdateSchool(id uint, req UpdateSchoolRequest) (*School, error)
    DeactivateSchool(id uint) error
}

// Repository Interface
type SchoolRepository interface {
    Create(school *School) error
    FindAll(pagination Pagination) ([]School, int64, error)
    FindByID(id uint) (*School, error)
    Update(school *School) error
    Deactivate(id uint) error
}
```

### 3. School Module (Master Data)

```go
// Handler Interface
type SchoolHandler interface {
    // Class Management
    CreateClass(c *fiber.Ctx) error
    GetClasses(c *fiber.Ctx) error
    UpdateClass(c *fiber.Ctx) error
    DeleteClass(c *fiber.Ctx) error
    
    // Student Management
    CreateStudent(c *fiber.Ctx) error
    GetStudents(c *fiber.Ctx) error
    GetStudent(c *fiber.Ctx) error
    UpdateStudent(c *fiber.Ctx) error
    
    // Parent Management
    CreateParent(c *fiber.Ctx) error
    LinkParentToStudent(c *fiber.Ctx) error
}

// Service Interface
type SchoolService interface {
    // Class
    CreateClass(schoolID uint, req CreateClassRequest) (*Class, error)
    GetClassesBySchool(schoolID uint) ([]Class, error)
    
    // Student
    CreateStudent(schoolID uint, req CreateStudentRequest) (*Student, error)
    GetStudentsByClass(classID uint) ([]Student, error)
    GetStudentByNISN(nisn string) (*Student, error)
    
    // Parent
    CreateParent(schoolID uint, req CreateParentRequest) (*Parent, error)
    LinkParentToStudent(parentID, studentID uint) error
}
```

### 4. Attendance Module

```go
// Handler Interface
type AttendanceHandler interface {
    RecordAttendance(c *fiber.Ctx) error      // From ESP32
    ManualAttendance(c *fiber.Ctx) error      // Manual fallback
    GetAttendanceByStudent(c *fiber.Ctx) error
    GetAttendanceByClass(c *fiber.Ctx) error
    GetAttendanceSummary(c *fiber.Ctx) error
}

// Service Interface
type AttendanceService interface {
    RecordRFIDAttendance(deviceCode string, studentRFID string, timestamp time.Time) (*Attendance, error)
    RecordManualAttendance(studentID uint, date time.Time, checkIn, checkOut *time.Time) (*Attendance, error)
    GetStudentAttendance(studentID uint, startDate, endDate time.Time) ([]Attendance, error)
    GetClassAttendance(classID uint, date time.Time) ([]AttendanceSummary, error)
    GetSchoolAttendanceSummary(schoolID uint, date time.Time) (*SchoolAttendanceSummary, error)
}

// Repository Interface
type AttendanceRepository interface {
    Create(attendance *Attendance) error
    FindByStudentAndDate(studentID uint, date time.Time) (*Attendance, error)
    UpdateCheckOut(id uint, checkOutTime time.Time) error
    FindByStudent(studentID uint, startDate, endDate time.Time) ([]Attendance, error)
    FindByClassAndDate(classID uint, date time.Time) ([]Attendance, error)
}
```

### 5. BK Module (Counseling)

```go
// Handler Interface
type BKHandler interface {
    // Violations
    CreateViolation(c *fiber.Ctx) error
    GetViolations(c *fiber.Ctx) error
    
    // Achievements
    CreateAchievement(c *fiber.Ctx) error
    GetAchievements(c *fiber.Ctx) error
    GetAchievementPoints(c *fiber.Ctx) error
    
    // Permits
    CreatePermit(c *fiber.Ctx) error
    GetPermits(c *fiber.Ctx) error
    RecordReturn(c *fiber.Ctx) error
    GetPermitDocument(c *fiber.Ctx) error
    
    // Counseling Notes
    CreateCounselingNote(c *fiber.Ctx) error
    GetCounselingNotes(c *fiber.Ctx) error
}

// Service Interface
type BKService interface {
    // Violations
    CreateViolation(req CreateViolationRequest) (*Violation, error)
    GetStudentViolations(studentID uint) ([]Violation, error)
    
    // Achievements
    CreateAchievement(req CreateAchievementRequest) (*Achievement, error)
    GetStudentAchievements(studentID uint) ([]Achievement, error)
    GetStudentAchievementPoints(studentID uint) (int, error)
    
    // Permits
    CreatePermit(req CreatePermitRequest) (*Permit, error)
    GetStudentPermits(studentID uint) ([]Permit, error)
    RecordReturn(permitID uint, returnTime time.Time) error
    GeneratePermitDocument(permitID uint) ([]byte, error)
    
    // Counseling
    CreateCounselingNote(req CreateCounselingNoteRequest) (*CounselingNote, error)
    GetStudentCounselingNotes(studentID uint, includeInternal bool) ([]CounselingNote, error)
}
```

### 6. Grade Module

```go
// Handler Interface
type GradeHandler interface {
    CreateGrade(c *fiber.Ctx) error
    GetGrades(c *fiber.Ctx) error
    UpdateGrade(c *fiber.Ctx) error
    DeleteGrade(c *fiber.Ctx) error
}

// Service Interface
type GradeService interface {
    CreateGrade(teacherID uint, req CreateGradeRequest) (*Grade, error)
    GetStudentGrades(studentID uint) ([]Grade, error)
    UpdateGrade(gradeID uint, req UpdateGradeRequest) (*Grade, error)
    DeleteGrade(gradeID uint) error
    ValidateTeacherAccess(teacherID, studentID uint) error
}

// Repository Interface
type GradeRepository interface {
    Create(grade *Grade) error
    FindByStudent(studentID uint) ([]Grade, error)
    FindByID(id uint) (*Grade, error)
    Update(grade *Grade) error
    Delete(id uint) error
}
```

### 7. Homeroom Module

```go
// Handler Interface
type HomeroomHandler interface {
    CreateNote(c *fiber.Ctx) error
    GetNotes(c *fiber.Ctx) error
    UpdateNote(c *fiber.Ctx) error
    DeleteNote(c *fiber.Ctx) error
}

// Service Interface
type HomeroomService interface {
    CreateNote(teacherID uint, req CreateNoteRequest) (*HomeroomNote, error)
    GetStudentNotes(studentID uint) ([]HomeroomNote, error)
    UpdateNote(noteID uint, req UpdateNoteRequest) (*HomeroomNote, error)
    DeleteNote(noteID uint) error
    ValidateTeacherAccess(teacherID, studentID uint) error
}
```

### 8. Notification Module

```go
// Handler Interface
type NotificationHandler interface {
    GetNotifications(c *fiber.Ctx) error
    MarkAsRead(c *fiber.Ctx) error
    MarkAllAsRead(c *fiber.Ctx) error
}

// Service Interface
type NotificationService interface {
    SendNotification(userID uint, notifType string, title, message string) error
    QueueNotification(notification *Notification) error
    GetUserNotifications(userID uint, pagination Pagination) ([]Notification, int64, error)
    MarkAsRead(notificationID uint) error
    MarkAllAsRead(userID uint) error
}

// Worker Interface
type NotificationWorker interface {
    Start()
    Stop()
    ProcessQueue()
}

// FCM Client Interface
type FCMClient interface {
    SendPushNotification(deviceToken string, title, body string, data map[string]string) error
}
```

### 9. Device Module

```go
// Handler Interface
type DeviceHandler interface {
    RegisterDevice(c *fiber.Ctx) error
    GetDevices(c *fiber.Ctx) error
    UpdateDevice(c *fiber.Ctx) error
    RevokeAPIKey(c *fiber.Ctx) error
    RegenerateAPIKey(c *fiber.Ctx) error
}

// Service Interface
type DeviceService interface {
    RegisterDevice(schoolID uint, req RegisterDeviceRequest) (*Device, error)
    GetDevicesBySchool(schoolID uint) ([]Device, error)
    ValidateAPIKey(apiKey string) (*Device, error)
    RevokeAPIKey(deviceID uint) error
    RegenerateAPIKey(deviceID uint) (*Device, error)
}
```

## Data Models

### Core Models

```go
// User represents all system users
type User struct {
    ID           uint      `gorm:"primaryKey"`
    SchoolID     *uint     `gorm:"index"` // nil for super_admin
    Role         string    `gorm:"type:varchar(20);not null"`
    Username     string    `gorm:"type:varchar(100);uniqueIndex;not null"`
    PasswordHash string    `gorm:"type:varchar(255);not null"`
    Email        string    `gorm:"type:varchar(255)"`
    IsActive     bool      `gorm:"default:true"`
    MustResetPwd bool      `gorm:"default:true"`
    LastLoginAt  *time.Time
    CreatedAt    time.Time
    UpdatedAt    time.Time
    
    School       *School   `gorm:"foreignKey:SchoolID"`
}

// School represents a tenant
type School struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"type:varchar(255);not null"`
    Address   string    `gorm:"type:text"`
    Phone     string    `gorm:"type:varchar(20)"`
    Email     string    `gorm:"type:varchar(255)"`
    IsActive  bool      `gorm:"default:true"`
    CreatedAt time.Time
    UpdatedAt time.Time
    
    Classes   []Class   `gorm:"foreignKey:SchoolID"`
    Students  []Student `gorm:"foreignKey:SchoolID"`
    Users     []User    `gorm:"foreignKey:SchoolID"`
    Devices   []Device  `gorm:"foreignKey:SchoolID"`
}

// Class represents a school class
type Class struct {
    ID        uint      `gorm:"primaryKey"`
    SchoolID  uint      `gorm:"index;not null"`
    Name      string    `gorm:"type:varchar(50);not null"`
    Grade     int       `gorm:"not null"` // e.g., 7, 8, 9 for SMP
    Year      string    `gorm:"type:varchar(10);not null"` // e.g., "2024/2025"
    CreatedAt time.Time
    UpdatedAt time.Time
    
    School    School    `gorm:"foreignKey:SchoolID"`
    Students  []Student `gorm:"foreignKey:ClassID"`
    Teacher   *User     `gorm:"foreignKey:HomeroomTeacherID"`
    HomeroomTeacherID *uint
}

// Student represents a student
type Student struct {
    ID        uint      `gorm:"primaryKey"`
    SchoolID  uint      `gorm:"index;not null"`
    ClassID   uint      `gorm:"index;not null"`
    NIS       string    `gorm:"type:varchar(20);not null"`
    NISN      string    `gorm:"type:varchar(20);uniqueIndex;not null"`
    Name      string    `gorm:"type:varchar(255);not null"`
    RFIDCode  string    `gorm:"type:varchar(50);index"`
    IsActive  bool      `gorm:"default:true"`
    CreatedAt time.Time
    UpdatedAt time.Time
    
    School    School    `gorm:"foreignKey:SchoolID"`
    Class     Class     `gorm:"foreignKey:ClassID"`
    Parents   []Parent  `gorm:"many2many:student_parents"`
}

// Parent represents a parent/guardian
type Parent struct {
    ID        uint      `gorm:"primaryKey"`
    SchoolID  uint      `gorm:"index;not null"`
    UserID    uint      `gorm:"uniqueIndex;not null"`
    Name      string    `gorm:"type:varchar(255);not null"`
    Phone     string    `gorm:"type:varchar(20)"`
    CreatedAt time.Time
    UpdatedAt time.Time
    
    User      User      `gorm:"foreignKey:UserID"`
    Students  []Student `gorm:"many2many:student_parents"`
}
```

### Attendance Model

```go
// Attendance represents daily attendance record
type Attendance struct {
    ID           uint       `gorm:"primaryKey"`
    StudentID    uint       `gorm:"index;not null"`
    Date         time.Time  `gorm:"type:date;index;not null"`
    CheckInTime  *time.Time
    CheckOutTime *time.Time
    Method       string     `gorm:"type:varchar(10);not null"` // rfid, manual
    CreatedAt    time.Time
    UpdatedAt    time.Time
    
    Student      Student    `gorm:"foreignKey:StudentID"`
}
```

### BK Models

```go
// Violation represents student violation record
type Violation struct {
    ID          uint      `gorm:"primaryKey"`
    StudentID   uint      `gorm:"index;not null"`
    Category    string    `gorm:"type:varchar(100);not null"`
    Level       string    `gorm:"type:varchar(20);not null"` // ringan, sedang, berat
    Description string    `gorm:"type:text;not null"`
    CreatedBy   uint      `gorm:"not null"`
    CreatedAt   time.Time
    
    Student     Student   `gorm:"foreignKey:StudentID"`
    Creator     User      `gorm:"foreignKey:CreatedBy"`
}

// Achievement represents student achievement record
type Achievement struct {
    ID          uint      `gorm:"primaryKey"`
    StudentID   uint      `gorm:"index;not null"`
    Title       string    `gorm:"type:varchar(255);not null"`
    Point       int       `gorm:"not null"`
    Description string    `gorm:"type:text"`
    CreatedBy   uint      `gorm:"not null"`
    CreatedAt   time.Time
    
    Student     Student   `gorm:"foreignKey:StudentID"`
    Creator     User      `gorm:"foreignKey:CreatedBy"`
}

// Permit represents school exit permit
type Permit struct {
    ID              uint       `gorm:"primaryKey"`
    StudentID       uint       `gorm:"index;not null"`
    Reason          string     `gorm:"type:text;not null"`
    ExitTime        time.Time  `gorm:"not null"`
    ReturnTime      *time.Time
    ResponsibleTeacher uint    `gorm:"not null"`
    DocumentURL     string     `gorm:"type:varchar(500)"`
    CreatedBy       uint       `gorm:"not null"`
    CreatedAt       time.Time
    
    Student         Student    `gorm:"foreignKey:StudentID"`
    Teacher         User       `gorm:"foreignKey:ResponsibleTeacher"`
    Creator         User       `gorm:"foreignKey:CreatedBy"`
}

// CounselingNote represents counseling session notes
type CounselingNote struct {
    ID            uint      `gorm:"primaryKey"`
    StudentID     uint      `gorm:"index;not null"`
    InternalNote  string    `gorm:"type:text;not null"`
    ParentSummary string    `gorm:"type:text"`
    CreatedBy     uint      `gorm:"not null"`
    CreatedAt     time.Time
    
    Student       Student   `gorm:"foreignKey:StudentID"`
    Creator       User      `gorm:"foreignKey:CreatedBy"`
}
```

### Academic Models

```go
// Grade represents student grade entry
type Grade struct {
    ID          uint      `gorm:"primaryKey"`
    StudentID   uint      `gorm:"index;not null"`
    Title       string    `gorm:"type:varchar(255);not null"`
    Score       float64   `gorm:"not null"`
    Description string    `gorm:"type:text"`
    CreatedBy   uint      `gorm:"not null"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    
    Student     Student   `gorm:"foreignKey:StudentID"`
    Creator     User      `gorm:"foreignKey:CreatedBy"`
}

// HomeroomNote represents homeroom teacher notes
type HomeroomNote struct {
    ID        uint      `gorm:"primaryKey"`
    StudentID uint      `gorm:"index;not null"`
    TeacherID uint      `gorm:"not null"`
    Content   string    `gorm:"type:text;not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
    
    Student   Student   `gorm:"foreignKey:StudentID"`
    Teacher   User      `gorm:"foreignKey:TeacherID"`
}
```

### Device & Notification Models

```go
// Device represents RFID device (ESP32)
type Device struct {
    ID          uint      `gorm:"primaryKey"`
    SchoolID    uint      `gorm:"index;not null"`
    DeviceCode  string    `gorm:"type:varchar(50);uniqueIndex;not null"`
    APIKey      string    `gorm:"type:varchar(255);uniqueIndex;not null"`
    Description string    `gorm:"type:varchar(255)"`
    IsActive    bool      `gorm:"default:true"`
    LastSeenAt  *time.Time
    CreatedAt   time.Time
    UpdatedAt   time.Time
    
    School      School    `gorm:"foreignKey:SchoolID"`
}

// Notification represents user notification
type Notification struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"index;not null"`
    Type      string    `gorm:"type:varchar(50);not null"`
    Title     string    `gorm:"type:varchar(255);not null"`
    Message   string    `gorm:"type:text;not null"`
    Data      string    `gorm:"type:jsonb"` // Additional JSON data
    IsRead    bool      `gorm:"default:false"`
    CreatedAt time.Time
    
    User      User      `gorm:"foreignKey:UserID"`
}

// FCMToken represents user's FCM device token
type FCMToken struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"index;not null"`
    Token     string    `gorm:"type:varchar(500);not null"`
    Platform  string    `gorm:"type:varchar(20);not null"` // android, ios
    IsActive  bool      `gorm:"default:true"`
    CreatedAt time.Time
    UpdatedAt time.Time
    
    User      User      `gorm:"foreignKey:UserID"`
}
```

### School Settings Model

```go
// SchoolSettings represents configurable settings per school
type SchoolSettings struct {
    ID                      uint      `gorm:"primaryKey"`
    SchoolID                uint      `gorm:"uniqueIndex;not null"`
    
    // Attendance Settings
    AttendanceStartTime     string    `gorm:"type:varchar(5);default:'07:00'"` // HH:MM format
    AttendanceEndTime       string    `gorm:"type:varchar(5);default:'07:30'"` // Late after this
    AttendanceLateThreshold int       `gorm:"default:30"` // Minutes after start to be considered late
    AttendanceVeryLateThreshold int   `gorm:"default:60"` // Minutes after start to be considered very late
    
    // Notification Settings
    EnableAttendanceNotification bool `gorm:"default:true"`
    EnableGradeNotification      bool `gorm:"default:true"`
    EnableBKNotification         bool `gorm:"default:true"`
    EnableHomeroomNotification   bool `gorm:"default:true"`
    
    // General Settings
    AcademicYear            string    `gorm:"type:varchar(10)"` // e.g., "2024/2025"
    Semester                int       `gorm:"default:1"` // 1 or 2
    
    CreatedAt               time.Time
    UpdatedAt               time.Time
    
    School                  School    `gorm:"foreignKey:SchoolID"`
}

// AttendanceTimeWindow represents the valid time window for attendance
type AttendanceTimeWindow struct {
    StartTime time.Time
    EndTime   time.Time
    LateTime  time.Time
    VeryLateTime time.Time
}
```

### Settings Module Interface

```go
// Handler Interface
type SettingsHandler interface {
    GetSettings(c *fiber.Ctx) error
    UpdateSettings(c *fiber.Ctx) error
    ResetToDefaults(c *fiber.Ctx) error
}

// Service Interface
type SettingsService interface {
    GetSchoolSettings(schoolID uint) (*SchoolSettings, error)
    UpdateSchoolSettings(schoolID uint, req UpdateSettingsRequest) (*SchoolSettings, error)
    GetAttendanceTimeWindow(schoolID uint, date time.Time) (*AttendanceTimeWindow, error)
    ShouldSendNotification(schoolID uint, notificationType string) (bool, error)
    ResetToDefaults(schoolID uint) (*SchoolSettings, error)
}

// Repository Interface
type SettingsRepository interface {
    FindBySchoolID(schoolID uint) (*SchoolSettings, error)
    Create(settings *SchoolSettings) error
    Update(settings *SchoolSettings) error
    Delete(schoolID uint) error
}
```



## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

Based on the prework analysis, the following correctness properties have been identified:

### Property 1: Tenant Data Isolation

*For any* data query in the system, the results SHALL only contain records belonging to the requesting user's tenant (school_id). Cross-tenant data access SHALL always be rejected with an authorization error.

**Validates: Requirements 1.4, 1.5**

### Property 2: Tenant Creation Uniqueness

*For any* newly created tenant, the system SHALL generate a unique school_id that does not exist in the database. No two tenants SHALL share the same school_id.

**Validates: Requirements 1.1**

### Property 3: Tenant Deactivation Access Control

*For any* deactivated tenant, all users belonging to that tenant SHALL be denied access to the system. Authentication attempts by users of deactivated tenants SHALL fail.

**Validates: Requirements 1.3**

### Property 4: Device API Key Security

*For any* registered device, the system SHALL generate a unique API key. *For any* attendance request, valid API keys SHALL be accepted and invalid/revoked API keys SHALL be rejected. Revoking an API key SHALL immediately invalidate it for all subsequent requests.

**Validates: Requirements 2.1, 2.2, 2.3, 2.4**

### Property 5: NISN Uniqueness

*For any* student registration, if the NISN already exists in the system, the registration SHALL be rejected with an error. No two students SHALL share the same NISN across the entire system.

**Validates: Requirements 3.5**

### Property 6: Role-Based Access Control Enforcement

*For any* protected resource access:
- Guru_BK SHALL have read/write access to BK data (violations, achievements, permits, counseling notes)
- Wali_Kelas SHALL have read-only access to BK data and read/write access only to students in their assigned class
- Parents SHALL only see parent_summary of counseling notes, never internal_note
- Users SHALL only access resources permitted by their role

**Validates: Requirements 4.5, 6.5, 9.3, 9.4, 10.5, 11.4, 14.5**

### Property 7: Attendance Check-In/Check-Out Logic

*For any* student on any given day:
- The first attendance record SHALL be recorded as check-in
- The second attendance record SHALL update the existing record with check-out time
- Check-out time SHALL always be after check-in time

**Validates: Requirements 5.2**

### Property 8: Achievement Point Accumulation

*For any* student, the total achievement points SHALL equal the sum of all individual achievement point values. Adding a new achievement SHALL increase the total by exactly the new achievement's point value.

**Validates: Requirements 7.2, 7.3**

### Property 9: Data Serialization Round-Trip

*For any* valid data object (Notification, complex JSON fields), serializing to JSON then deserializing back SHALL produce an object equivalent to the original. `deserialize(serialize(obj)) == obj`

**Validates: Requirements 17.6, 17.7, 18.4, 18.5, 18.6**

### Property 10: Notification Queue Processing

*For any* event that triggers a notification:
- The notification SHALL be queued in Redis
- The worker SHALL process and send via FCM
- Failed deliveries SHALL be retried with exponential backoff
- The notification SHALL eventually be stored in the database

**Validates: Requirements 17.1, 17.2, 17.5**

### Property 11: History Maintenance with Timestamps

*For any* entity that maintains history (violations, achievements, counseling notes, homeroom notes, grades), all records SHALL have a created_at timestamp. Records SHALL be retrievable in chronological order.

**Validates: Requirements 6.3, 7.5, 9.5, 11.5, 18.1**

### Property 12: Required Field Validation

*For any* entity creation:
- Student registration SHALL require NIS, NISN, name, and class_id
- Violation SHALL require category, level, description, and student_id
- Achievement SHALL require title, point, and student_id
- Permit SHALL require reason, exit_time, and responsible_teacher
- Counseling note SHALL require internal_note and student_id
- Grade SHALL require title, score, and student_id
- Homeroom note SHALL require content and student_id

Missing required fields SHALL result in validation error.

**Validates: Requirements 3.2, 6.1, 7.1, 8.1, 9.1, 10.1, 11.1**

### Property 13: Permit Document Content Completeness

*For any* generated permit document, the document SHALL contain: student name, NIS, NISN, class name, reason, exit time, responsible teacher name, and generation timestamp.

**Validates: Requirements 8.2, 8.5**

### Property 14: Notification Content Completeness

*For any* attendance notification, the message SHALL contain: student name, timestamp, and attendance type (masuk/pulang).

**Validates: Requirements 13.3**

### Property 15: Parent Data Access Completeness

*For any* parent with linked children, the parent SHALL be able to access all data (attendance, grades, notes, BK info) for all their linked children. Switching between children SHALL show the correct child's data.

**Validates: Requirements 12.2, 14.4, 15.1, 15.2**

### Property 16: Data Sorting by Date

*For any* list of historical records (violations, grades, notes, attendance), the records SHALL be sorted by date in descending order (most recent first) when displayed to users.

**Validates: Requirements 6.3, 10.4, 15.5**

### Property 17: School Settings Policy Enforcement

*For any* school with configured settings:
- Attendance status (on_time, late, very_late) SHALL be determined based on the school's configured time thresholds
- Notifications SHALL only be sent if the corresponding notification setting is enabled
- Default settings SHALL be applied when a new school is created

**Validates: School-specific configuration requirements from PRD**

### Property 18: Event Outbox Reliability

*For any* event published through the outbox:
- The event SHALL be persisted in the database within the same transaction as the business operation
- Pending events SHALL eventually be published (at-least-once delivery)
- Failed events SHALL be retried with exponential backoff

**Validates: Reliable event delivery for notifications**

## Error Handling

### API Error Response Format

```go
type ErrorResponse struct {
    Success bool   `json:"success"`
    Error   struct {
        Code    string `json:"code"`
        Message string `json:"message"`
        Details any    `json:"details,omitempty"`
    } `json:"error"`
}
```

### Error Categories

| Category | HTTP Status | Error Codes |
|----------|-------------|-------------|
| Authentication | 401 | `AUTH_INVALID_CREDENTIALS`, `AUTH_TOKEN_EXPIRED`, `AUTH_TOKEN_INVALID` |
| Authorization | 403 | `AUTHZ_FORBIDDEN`, `AUTHZ_TENANT_MISMATCH`, `AUTHZ_ROLE_DENIED` |
| Validation | 400 | `VAL_REQUIRED_FIELD`, `VAL_INVALID_FORMAT`, `VAL_DUPLICATE_ENTRY` |
| Not Found | 404 | `NOT_FOUND_STUDENT`, `NOT_FOUND_CLASS`, `NOT_FOUND_USER` |
| Server Error | 500 | `INTERNAL_ERROR`, `DB_ERROR`, `QUEUE_ERROR` |

### Error Handling Strategy

1. **Validation Errors**: Return immediately with specific field errors
2. **Authentication Errors**: Log attempt, return generic message (security)
3. **Authorization Errors**: Log attempt with user context, return forbidden
4. **Database Errors**: Log full error, return generic message to client
5. **External Service Errors (FCM, Redis)**: Retry with backoff, queue for later if persistent

### Retry Strategy for Notifications

```go
type RetryConfig struct {
    MaxRetries     int           // 5
    InitialDelay   time.Duration // 1 second
    MaxDelay       time.Duration // 5 minutes
    BackoffFactor  float64       // 2.0
}
```

## Testing Strategy

### Dual Testing Approach

This system requires both unit tests and property-based tests for comprehensive coverage:

1. **Unit Tests**: Verify specific examples, edge cases, and error conditions
2. **Property-Based Tests**: Verify universal properties across all valid inputs

### Property-Based Testing Configuration

- **Library**: `github.com/leanovate/gopter` for Go
- **Minimum Iterations**: 100 per property test
- **Tag Format**: `Feature: school-management-saas, Property {number}: {property_text}`

### Test Categories

#### Unit Tests
- Handler tests with mocked services
- Service tests with mocked repositories
- Repository tests with test database
- Middleware tests (auth, tenant, RBAC)
- Utility function tests

#### Property-Based Tests
- Tenant isolation property
- API key uniqueness and validation
- NISN uniqueness
- RBAC enforcement
- Attendance logic
- Point accumulation
- Serialization round-trip
- Required field validation

#### Integration Tests
- Full API flow tests
- Notification queue processing
- Multi-tenant scenarios
- Authentication flows

### Test File Structure

```
backend/
├── internal/
│   ├── modules/
│   │   ├── auth/
│   │   │   ├── handler_test.go
│   │   │   ├── service_test.go
│   │   │   └── service_property_test.go
│   │   ├── tenant/
│   │   │   ├── handler_test.go
│   │   │   ├── service_test.go
│   │   │   └── service_property_test.go
│   │   └── ...
│   └── middleware/
│       ├── auth_test.go
│       ├── tenant_test.go
│       └── rbac_test.go
└── tests/
    └── integration/
        ├── auth_flow_test.go
        ├── attendance_flow_test.go
        └── notification_flow_test.go
```

### Property Test Example

```go
// Feature: school-management-saas, Property 9: Data Serialization Round-Trip
func TestNotificationSerializationRoundTrip(t *testing.T) {
    parameters := gopter.DefaultTestParameters()
    parameters.MinSuccessfulTests = 100
    
    properties := gopter.NewProperties(parameters)
    
    properties.Property("serialize then deserialize produces equivalent notification", 
        prop.ForAll(
            func(n Notification) bool {
                serialized, err := json.Marshal(n)
                if err != nil {
                    return false
                }
                var deserialized Notification
                err = json.Unmarshal(serialized, &deserialized)
                if err != nil {
                    return false
                }
                return reflect.DeepEqual(n, deserialized)
            },
            genNotification(),
        ))
    
    properties.TestingRun(t)
}
```

### Coverage Requirements

- Unit test coverage: minimum 80%
- All correctness properties must have corresponding property tests
- All API endpoints must have integration tests
- Critical paths (auth, tenant isolation) must have 100% coverage
