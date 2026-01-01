# Design Document: Attendance System Enhancement

## Overview

Dokumen ini menjelaskan desain teknis untuk pengembangan 5 fitur tambahan sistem absensi RFID:
1. Export Absensi ke Excel
2. Rekap Bulanan Absensi
3. Pengaturan Waktu Absensi (Multi-Schedule)
4. Real-Time Attendance Dashboard
5. Public Display Mode (Kiosk)

Sistem ini dibangun menggunakan Go (Fiber) untuk backend dan Vue 3 (Ant Design Vue) untuk frontend. Database menggunakan PostgreSQL dengan GORM sebagai ORM.

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend (Vue 3)                        │
├─────────────────────────────────────────────────────────────────┤
│  Admin Dashboard  │  Live Attendance  │  Public Display (Kiosk) │
└────────┬──────────┴────────┬──────────┴────────┬────────────────┘
         │                   │                   │
         │ REST API          │ WebSocket/SSE     │ REST API (Token)
         ▼                   ▼                   ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Backend (Go Fiber)                         │
├─────────────────────────────────────────────────────────────────┤
│  Attendance  │  Schedule  │  Export  │  Display  │  Realtime   │
│   Module     │   Module   │  Module  │  Module   │   Hub       │
└──────┬───────┴─────┬──────┴────┬─────┴─────┬─────┴──────┬──────┘
       │             │           │           │            │
       ▼             ▼           ▼           ▼            ▼
┌─────────────────────────────────────────────────────────────────┐
│                      PostgreSQL Database                        │
├─────────────────────────────────────────────────────────────────┤
│ attendances │ attendance_schedules │ display_tokens │ students  │
└─────────────────────────────────────────────────────────────────┘
```

### Real-Time Communication Flow

```
┌──────────┐     RFID Tap      ┌──────────┐
│  ESP32   │ ─────────────────▶│  Backend │
│  Device  │                   │  Server  │
└──────────┘                   └────┬─────┘
                                    │
                    ┌───────────────┼───────────────┐
                    │               │               │
                    ▼               ▼               ▼
              ┌──────────┐   ┌──────────┐   ┌──────────┐
              │  Admin   │   │  Wali    │   │  Public  │
              │ Dashboard│   │  Kelas   │   │ Display  │
              └──────────┘   └──────────┘   └──────────┘
                    ▲               ▲               ▲
                    │               │               │
                    └───────────────┴───────────────┘
                              WebSocket/SSE
```

## Components and Interfaces

### 1. Attendance Schedule Module

#### Data Model

```go
// AttendanceSchedule represents a configurable attendance time slot
type AttendanceSchedule struct {
    ID               uint      `gorm:"primaryKey"`
    SchoolID         uint      `gorm:"index;not null"`
    Name             string    `gorm:"type:varchar(100);not null"` // "Masuk Pagi", "Pulang", "Sholat Dzuhur"
    StartTime        string    `gorm:"type:time;not null"`         // "07:00"
    EndTime          string    `gorm:"type:time;not null"`         // "08:00"
    LateThreshold    int       `gorm:"not null"`                   // minutes after start_time
    VeryLateThreshold *int     `gorm:""`                           // optional
    DaysOfWeek       string    `gorm:"type:varchar(20)"`           // "1,2,3,4,5" (Mon-Fri)
    IsActive         bool      `gorm:"default:true"`
    IsDefault        bool      `gorm:"default:false"`              // default schedule for the school
    CreatedAt        time.Time
    UpdatedAt        time.Time
    
    School           School    `gorm:"foreignKey:SchoolID"`
}
```

#### Service Interface

```go
type ScheduleService interface {
    CreateSchedule(ctx context.Context, schoolID uint, req CreateScheduleRequest) (*ScheduleResponse, error)
    GetAllSchedules(ctx context.Context, schoolID uint) ([]ScheduleResponse, error)
    GetScheduleByID(ctx context.Context, schoolID uint, id uint) (*ScheduleResponse, error)
    UpdateSchedule(ctx context.Context, schoolID uint, id uint, req UpdateScheduleRequest) (*ScheduleResponse, error)
    DeleteSchedule(ctx context.Context, schoolID uint, id uint) error
    GetActiveSchedule(ctx context.Context, schoolID uint, timestamp time.Time) (*ScheduleResponse, error)
    SetDefaultSchedule(ctx context.Context, schoolID uint, id uint) error
}
```

#### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/schedules` | List all schedules for school |
| POST | `/api/v1/schedules` | Create new schedule |
| GET | `/api/v1/schedules/:id` | Get schedule by ID |
| PUT | `/api/v1/schedules/:id` | Update schedule |
| DELETE | `/api/v1/schedules/:id` | Delete schedule |
| POST | `/api/v1/schedules/:id/default` | Set as default schedule |

### 2. Display Token Module

#### Data Model

```go
// DisplayToken represents a token for public display access
type DisplayToken struct {
    ID             uint       `gorm:"primaryKey"`
    SchoolID       uint       `gorm:"index;not null"`
    Token          string     `gorm:"type:varchar(64);uniqueIndex;not null"`
    Name           string     `gorm:"type:varchar(100)"` // "Display Pintu Utama"
    IsActive       bool       `gorm:"default:true"`
    LastAccessedAt *time.Time
    ExpiresAt      *time.Time // optional expiration
    CreatedAt      time.Time
    UpdatedAt      time.Time
    
    School         School     `gorm:"foreignKey:SchoolID"`
}
```

#### Service Interface

```go
type DisplayTokenService interface {
    CreateToken(ctx context.Context, schoolID uint, req CreateDisplayTokenRequest) (*DisplayTokenWithSecretResponse, error)
    GetAllTokens(ctx context.Context, schoolID uint) ([]DisplayTokenResponse, error)
    GetTokenByID(ctx context.Context, schoolID uint, id uint) (*DisplayTokenResponse, error)
    ValidateToken(ctx context.Context, token string) (*DisplayTokenValidation, error)
    RevokeToken(ctx context.Context, schoolID uint, id uint) error
    RegenerateToken(ctx context.Context, schoolID uint, id uint) (*DisplayTokenWithSecretResponse, error)
    DeleteToken(ctx context.Context, schoolID uint, id uint) error
    UpdateLastAccessed(ctx context.Context, token string) error
}
```

#### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/display-tokens` | List all display tokens |
| POST | `/api/v1/display-tokens` | Create new display token |
| GET | `/api/v1/display-tokens/:id` | Get token by ID |
| POST | `/api/v1/display-tokens/:id/revoke` | Revoke token |
| POST | `/api/v1/display-tokens/:id/regenerate` | Regenerate token |
| DELETE | `/api/v1/display-tokens/:id` | Delete token |

### 3. Export Module

#### Service Interface

```go
type ExportService interface {
    ExportAttendanceToExcel(ctx context.Context, schoolID uint, filter ExportFilter) ([]byte, string, error)
    ExportMonthlyRecap(ctx context.Context, schoolID uint, filter MonthlyRecapFilter) ([]byte, string, error)
}

type ExportFilter struct {
    StartDate string // YYYY-MM-DD
    EndDate   string // YYYY-MM-DD
    ClassID   *uint
}

type MonthlyRecapFilter struct {
    Month   int   // 1-12
    Year    int   // 2024
    ClassID *uint
}
```

#### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/attendance/export` | Export attendance to Excel |
| GET | `/api/v1/attendance/monthly-recap` | Get monthly recap data |
| GET | `/api/v1/attendance/monthly-recap/export` | Export monthly recap to Excel |

### 4. Real-Time Module

#### WebSocket Hub

```go
type RealtimeHub struct {
    clients    map[*Client]bool
    broadcast  chan *AttendanceEvent
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

type Client struct {
    hub      *RealtimeHub
    conn     *websocket.Conn
    send     chan []byte
    schoolID uint
    classID  *uint // optional filter
    isPublic bool  // for public display
    token    string // display token for public
}

type AttendanceEvent struct {
    Type        string              `json:"type"` // "new_attendance", "stats_update"
    SchoolID    uint                `json:"school_id"`
    Attendance  *AttendanceResponse `json:"attendance,omitempty"`
    Stats       *AttendanceStats    `json:"stats,omitempty"`
    Leaderboard []LeaderboardEntry  `json:"leaderboard,omitempty"`
}
```

#### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/ws/attendance` | WebSocket for authenticated users |
| GET | `/api/v1/public/display/:token` | Public display data (REST) |
| GET | `/api/v1/public/display/:token/ws` | WebSocket for public display |

### 5. Public Display Module

#### Response DTOs

```go
type PublicDisplayData struct {
    SchoolName   string             `json:"school_name"`
    CurrentTime  time.Time          `json:"current_time"`
    Date         string             `json:"date"`
    Stats        AttendanceStats    `json:"stats"`
    LiveFeed     []LiveFeedEntry    `json:"live_feed"`
    Leaderboard  []LeaderboardEntry `json:"leaderboard"`
}

type LiveFeedEntry struct {
    StudentName string    `json:"student_name"`
    ClassName   string    `json:"class_name"`
    Time        time.Time `json:"time"`
    Status      string    `json:"status"` // "on_time", "late", "very_late"
    Type        string    `json:"type"`   // "check_in", "check_out"
}

type LeaderboardEntry struct {
    Rank        int       `json:"rank"`
    StudentName string    `json:"student_name"`
    ClassName   string    `json:"class_name"`
    ArrivalTime time.Time `json:"arrival_time"`
}

type AttendanceStats struct {
    TotalStudents int     `json:"total_students"`
    Present       int     `json:"present"`
    Late          int     `json:"late"`
    VeryLate      int     `json:"very_late"`
    Absent        int     `json:"absent"`
    Percentage    float64 `json:"percentage"`
}
```

## Data Models

### Updated Attendance Model

```go
type Attendance struct {
    ID           uint              `gorm:"primaryKey"`
    StudentID    uint              `gorm:"index;not null"`
    ScheduleID   *uint             `gorm:"index"` // NEW: link to schedule
    Date         time.Time         `gorm:"type:date;index;not null"`
    CheckInTime  *time.Time
    CheckOutTime *time.Time
    Status       AttendanceStatus  `gorm:"type:varchar(20)"`
    Method       AttendanceMethod  `gorm:"type:varchar(20)"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
    
    Student      Student           `gorm:"foreignKey:StudentID"`
    Schedule     *AttendanceSchedule `gorm:"foreignKey:ScheduleID"` // NEW
}
```

### Database Schema Changes

```sql
-- New table: attendance_schedules
CREATE TABLE attendance_schedules (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    late_threshold INT NOT NULL DEFAULT 15,
    very_late_threshold INT,
    days_of_week VARCHAR(20) DEFAULT '1,2,3,4,5',
    is_active BOOLEAN DEFAULT true,
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_attendance_schedules_school ON attendance_schedules(school_id);

-- New table: display_tokens
CREATE TABLE display_tokens (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    token VARCHAR(64) UNIQUE NOT NULL,
    name VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    last_accessed_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_display_tokens_school ON display_tokens(school_id);
CREATE INDEX idx_display_tokens_token ON display_tokens(token);

-- Add schedule_id to attendances table
ALTER TABLE attendances ADD COLUMN schedule_id INT REFERENCES attendance_schedules(id);
CREATE INDEX idx_attendances_schedule ON attendances(schedule_id);
```



## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Export Filter Correctness

*For any* attendance export request with date range filter, all records in the exported file SHALL have dates within the specified range (inclusive).

**Validates: Requirements 1.2, 1.3**

### Property 2: Export Data Completeness

*For any* attendance export, each record SHALL contain all required fields: student NIS, NISN, name, class, date, check-in time, check-out time, and status.

**Validates: Requirements 1.4, 1.5**

### Property 3: Export Filename Format

*For any* generated export file, the filename SHALL match the pattern `attendance_{school_name}_{start_date}_to_{end_date}.xlsx`.

**Validates: Requirements 1.7**

### Property 4: Monthly Recap Calculation Correctness

*For any* monthly recap, the sum of (present + late + very_late + absent) for each student SHALL equal the total school days in that month, and the percentage SHALL equal (present / total_days) * 100.

**Validates: Requirements 2.1, 2.2**

### Property 5: Monthly Recap Sorting

*For any* monthly recap response, students SHALL be sorted by attendance percentage in descending order.

**Validates: Requirements 2.6**

### Property 6: Wali Kelas Data Isolation

*For any* request from a Wali_Kelas user, the response SHALL only contain students from their assigned class.

**Validates: Requirements 2.7, 4.6**

### Property 7: Schedule Required Fields Validation

*For any* schedule creation request missing name, start_time, end_time, or late_threshold, the system SHALL reject the request with a validation error.

**Validates: Requirements 3.1**

### Property 8: Active Schedule Selection

*For any* RFID tap at time T, the system SHALL select the schedule where start_time <= T <= end_time and the current day is in days_of_week and is_active is true.

**Validates: Requirements 3.4, 3.5**

### Property 9: Schedule Update Data Integrity

*For any* schedule update, existing attendance records linked to that schedule SHALL remain unchanged (schedule_id and recorded times preserved).

**Validates: Requirements 3.7**

### Property 10: Live Feed Limit and Order

*For any* live attendance feed request, the response SHALL contain at most 20 records ordered by check-in time descending (most recent first).

**Validates: Requirements 4.3**

### Property 11: Attendance Stats Calculation

*For any* attendance statistics response, the percentage SHALL equal (present / total_students) * 100, and present + late + very_late + absent SHALL equal total_students.

**Validates: Requirements 4.1, 4.10**

### Property 12: Display Token Uniqueness

*For any* two display tokens in the system, their token values SHALL be different.

**Validates: Requirements 5.1, 6.2**

### Property 13: Public Display Data Privacy

*For any* public display response, student records SHALL only contain name, class_name, time, and status. No NIS, NISN, or other sensitive data SHALL be exposed.

**Validates: Requirements 5.14**

### Property 14: Display Token Access Control

*For any* revoked or expired display token, accessing the public display endpoint SHALL return an error response, not attendance data.

**Validates: Requirements 5.10, 5.13**

### Property 15: Leaderboard Correctness

*For any* public display leaderboard, it SHALL contain at most 10 students sorted by arrival time ascending (earliest first), and all entries SHALL have check-in time on the current date.

**Validates: Requirements 5.6**

### Property 16: Token Lifecycle Round-Trip

*For any* display token, after regeneration the old token SHALL be invalid and the new token SHALL be valid for accessing public display.

**Validates: Requirements 6.5**

### Property 17: Token Last Accessed Tracking

*For any* successful public display access, the display token's last_accessed_at timestamp SHALL be updated to a time >= the access time.

**Validates: Requirements 6.7**

## Error Handling

### Export Module Errors

| Error Code | HTTP Status | Description |
|------------|-------------|-------------|
| `EXPORT_INVALID_DATE_RANGE` | 400 | Start date is after end date |
| `EXPORT_DATE_RANGE_TOO_LARGE` | 400 | Date range exceeds 1 year |
| `EXPORT_NO_DATA` | 200 | No records found (returns empty file) |
| `EXPORT_GENERATION_FAILED` | 500 | Failed to generate Excel file |

### Schedule Module Errors

| Error Code | HTTP Status | Description |
|------------|-------------|-------------|
| `SCHEDULE_NOT_FOUND` | 404 | Schedule ID not found |
| `SCHEDULE_LIMIT_EXCEEDED` | 400 | School already has 10 schedules |
| `SCHEDULE_TIME_OVERLAP` | 400 | Schedule times overlap with existing |
| `SCHEDULE_INVALID_TIME` | 400 | End time is before start time |
| `SCHEDULE_IN_USE` | 400 | Cannot delete schedule with attendance records |

### Display Token Errors

| Error Code | HTTP Status | Description |
|------------|-------------|-------------|
| `TOKEN_NOT_FOUND` | 404 | Token ID not found |
| `TOKEN_INVALID` | 401 | Token is invalid or not found |
| `TOKEN_REVOKED` | 401 | Token has been revoked |
| `TOKEN_EXPIRED` | 401 | Token has expired |

### Real-Time Errors

| Error Code | Description |
|------------|-------------|
| `WS_AUTH_FAILED` | WebSocket authentication failed |
| `WS_CONNECTION_LOST` | Connection to server lost |
| `WS_RECONNECTING` | Attempting to reconnect |

## Testing Strategy

### Unit Tests

Unit tests will cover:
- Schedule time overlap detection
- Active schedule selection logic
- Attendance percentage calculation
- Token generation and validation
- Export file generation
- Date range validation

### Property-Based Tests

Property-based tests will use Go's `testing/quick` package or a library like `gopter` to verify:
- Filter correctness across random date ranges
- Calculation accuracy with random attendance data
- Token uniqueness across many generations
- Data privacy in public display responses

### Integration Tests

Integration tests will verify:
- End-to-end export flow
- WebSocket connection and message delivery
- Public display access with valid/invalid tokens
- Schedule-attendance association

### Test Configuration

- Minimum 100 iterations per property test
- Each property test tagged with: **Feature: attendance-system-enhancement, Property {number}: {property_text}**

