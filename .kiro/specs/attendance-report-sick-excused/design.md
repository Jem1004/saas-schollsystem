# Design Document: Attendance Report Sick & Excused Status

## Overview

Dokumen ini menjelaskan desain teknis untuk menambahkan tampilan status "Sakit" (sick) dan "Izin" (excused) pada halaman Laporan Absensi Admin Sekolah. Perubahan mencakup update pada backend API response, frontend types, dan komponen UI.

## Architecture

Sistem menggunakan arsitektur existing:
- **Backend**: Go dengan Fiber framework dan GORM ORM
- **Frontend**: Vue 3 dengan TypeScript dan Ant Design Vue
- **Database**: PostgreSQL dengan tabel `attendances` yang sudah memiliki kolom `status` dengan nilai `sick` dan `excused`

```
┌─────────────────────────────────────────────────────────────────┐
│                        Frontend (Vue 3)                         │
├─────────────────────────────────────────────────────────────────┤
│  AttendanceReport.vue                                           │
│  ├── Tab Harian (Daily)                                         │
│  │   ├── Summary Cards (+ Sakit, Izin)                         │
│  │   └── Class Table (+ kolom Sakit, Izin)                     │
│  ├── Tab Bulanan (Monthly)                                      │
│  │   ├── Summary Cards (+ Rata-rata Sakit, Izin)               │
│  │   └── Student Table (+ kolom Sakit, Izin)                   │
│  └── Tab Detail                                                 │
│      └── Status Filter (+ Sakit, Izin options)                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Backend API (Go/Fiber)                      │
├─────────────────────────────────────────────────────────────────┤
│  GET /attendance/summary                                        │
│  └── Response: + sick, excused counts                          │
│  GET /export/monthly-recap                                      │
│  └── Response: + total_sick, total_excused per student         │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Database (PostgreSQL)                        │
├─────────────────────────────────────────────────────────────────┤
│  attendances.status: on_time|late|very_late|absent|sick|excused│
└─────────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### Backend Changes

#### 1. Update DTO - `backend/internal/modules/attendance/dto.go`

```go
// AttendanceSummaryResponse - tambah field Sick dan Excused
type AttendanceSummaryResponse struct {
    Date       string `json:"date"`
    TotalCount int    `json:"total_count"`
    Present    int    `json:"present"`
    Late       int    `json:"late"`
    VeryLate   int    `json:"very_late"`
    Absent     int    `json:"absent"`
    Sick       int    `json:"sick"`      // NEW
    Excused    int    `json:"excused"`   // NEW
}

// ClassSummaryItem - tambah field Sick dan Excused
type ClassSummaryItem struct {
    ClassID       uint   `json:"class_id"`
    ClassName     string `json:"class_name"`
    TotalStudents int    `json:"total_students"`
    Present       int    `json:"present"`
    Late          int    `json:"late"`
    Absent        int    `json:"absent"`
    Sick          int    `json:"sick"`     // NEW
    Excused       int    `json:"excused"`  // NEW
}

// StudentRecapSummary - tambah field TotalSick dan TotalExcused
type StudentRecapSummary struct {
    StudentID         uint    `json:"student_id"`
    StudentNIS        string  `json:"student_nis"`
    StudentNISN       string  `json:"student_nisn"`
    StudentName       string  `json:"student_name"`
    ClassName         string  `json:"class_name"`
    TotalPresent      int     `json:"total_present"`
    TotalLate         int     `json:"total_late"`
    TotalVeryLate     int     `json:"total_very_late"`
    TotalAbsent       int     `json:"total_absent"`
    TotalSick         int     `json:"total_sick"`     // NEW
    TotalExcused      int     `json:"total_excused"`  // NEW
    AttendancePercent float64 `json:"attendance_percent"`
}
```

#### 2. Update Repository - `backend/internal/modules/attendance/repository.go`

Update fungsi `GetAttendanceSummary`, `GetAttendanceSummaryByClass`, dan `GetMonthlyRecap` untuk menghitung status `sick` dan `excused`.

### Frontend Changes

#### 1. Update Types - `web-admin/src/types/school.ts`

```typescript
// Update AttendanceSummary interface
export interface AttendanceSummary {
  classId: number
  className: string
  totalStudents: number
  present: number
  late: number
  absent: number
  sick: number      // NEW
  excused: number   // NEW
}
```

#### 2. Update Types - `web-admin/src/types/export.ts`

```typescript
// Update StudentRecapSummary interface
export interface StudentRecapSummary {
  student_id: number
  student_nis: string
  student_nisn: string
  student_name: string
  class_name: string
  total_present: number
  total_late: number
  total_very_late: number
  total_absent: number
  total_sick: number      // NEW
  total_excused: number   // NEW
  attendance_percent: number
}
```

#### 3. Update AttendanceReport.vue

- Tambah 2 Summary Cards untuk Sakit dan Izin di tab Harian
- Tambah 2 kolom di tabel Harian
- Tambah 2 Summary Cards untuk rata-rata Sakit dan Izin di tab Bulanan
- Tambah 2 kolom di tabel Bulanan
- Tambah 2 opsi filter status di tab Detail

## Data Models

### Existing Database Schema (No Changes Required)

```sql
-- attendances table sudah support status sick dan excused
CREATE TABLE attendances (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL,
    schedule_id INTEGER,
    date DATE NOT NULL,
    check_in_time TIME,
    check_out_time TIME,
    status VARCHAR(20), -- on_time, late, very_late, absent, sick, excused
    method VARCHAR(10) NOT NULL, -- rfid, manual
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

### API Response Changes

#### Daily Summary Response (Updated)

```json
{
  "success": true,
  "data": [
    {
      "class_id": 1,
      "class_name": "X-A",
      "total_students": 30,
      "present": 25,
      "late": 2,
      "absent": 1,
      "sick": 1,
      "excused": 1
    }
  ]
}
```

#### Monthly Recap Response (Updated)

```json
{
  "month": 1,
  "year": 2026,
  "total_days": 22,
  "student_recaps": [
    {
      "student_id": 1,
      "student_nis": "12345",
      "student_name": "John Doe",
      "class_name": "X-A",
      "total_present": 18,
      "total_late": 2,
      "total_very_late": 0,
      "total_absent": 0,
      "total_sick": 1,
      "total_excused": 1,
      "attendance_percent": 81.82
    }
  ]
}
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Daily Summary Counts Consistency

*For any* daily attendance summary response, the sum of present + late + absent + sick + excused SHALL equal total_students for each class.

**Validates: Requirements 1.3, 1.4, 2.1, 2.2**

### Property 2: Monthly Recap Counts Consistency

*For any* student monthly recap, the sum of total_present + total_late + total_very_late + total_absent + total_sick + total_excused SHALL be less than or equal to total_days.

**Validates: Requirements 3.3, 3.4, 4.1, 4.2**

### Property 3: Status Aggregation Accuracy

*For any* attendance status aggregation query, the count returned SHALL match the actual count of records with that status in the database for the given date range and filters.

**Validates: Requirements 5.1, 5.2, 5.3, 5.4, 6.1, 6.2, 6.3, 6.4**

### Property 4: Filter Results Correctness

*For any* status filter applied in the detail tab, all returned attendance records SHALL have the selected status value.

**Validates: Requirements 7.3, 7.4**

### Property 5: Export Data Completeness

*For any* exported attendance report, the sick and excused columns SHALL contain the same values as displayed in the UI for the same data set.

**Validates: Requirements 8.1, 8.2, 8.3, 8.4**

## Error Handling

### Backend Error Handling

1. **Invalid Status Value**: Jika status tidak valid, backend akan mengembalikan error 400 Bad Request
2. **Database Query Error**: Jika query gagal, backend akan log error dan mengembalikan 500 Internal Server Error
3. **Empty Data**: Jika tidak ada data, backend akan mengembalikan array kosong dengan status 200

### Frontend Error Handling

1. **API Error**: Menampilkan toast message error dan fallback ke nilai 0 untuk counts
2. **Missing Fields**: Menggunakan optional chaining dan default values untuk backward compatibility
3. **Loading State**: Menampilkan skeleton/spinner saat data sedang dimuat

## Testing Strategy

### Unit Tests

1. **Backend Repository Tests**
   - Test `GetAttendanceSummary` returns correct sick and excused counts
   - Test `GetAttendanceSummaryByClass` aggregates sick/excused per class correctly
   - Test `GetMonthlyRecap` calculates total_sick and total_excused per student

2. **Frontend Component Tests**
   - Test summary cards render with correct values
   - Test table columns display sick and excused data
   - Test filter options include sick and excused

### Property-Based Tests

Property-based testing akan digunakan untuk memvalidasi correctness properties di atas dengan minimum 100 iterasi per test.

1. **PBT: Daily Summary Consistency**
   - Generate random attendance data
   - Verify sum of all status counts equals total students

2. **PBT: Monthly Recap Consistency**
   - Generate random monthly attendance records
   - Verify total counts don't exceed total days

3. **PBT: Status Filter Accuracy**
   - Generate random attendance records with various statuses
   - Apply filter and verify all results match filter criteria

### Integration Tests

1. **API Integration Tests**
   - Test full flow from API call to database query
   - Verify response structure matches expected schema

2. **E2E Tests**
   - Test user flow viewing daily summary with sick/excused
   - Test user flow viewing monthly recap with sick/excused
   - Test filter functionality in detail tab
