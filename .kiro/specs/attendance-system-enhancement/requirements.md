# Requirements Document

## Introduction

Dokumen ini mendefinisikan requirements untuk pengembangan fitur tambahan sistem absensi RFID pada aplikasi manajemen sekolah. Fitur-fitur ini bertujuan untuk meningkatkan fleksibilitas, monitoring real-time, dan kemampuan pelaporan sistem absensi yang sudah ada.

## Glossary

- **System**: Aplikasi manajemen sekolah (backend + frontend)
- **Admin_Sekolah**: Administrator sekolah yang mengelola data sekolah
- **Wali_Kelas**: Guru yang ditugaskan sebagai wali kelas
- **Guru_BK**: Guru Bimbingan Konseling
- **Attendance_Schedule**: Jadwal waktu absensi yang dapat dikustomisasi (masuk, pulang, sholat, dll)
- **Display_Token**: Token unik untuk mengakses halaman public display tanpa login
- **Public_Display**: Halaman khusus untuk ditampilkan di LCD/monitor publik sekolah
- **Live_Feed**: Daftar siswa yang baru melakukan absensi secara real-time
- **Leaderboard**: Daftar peringkat siswa yang datang paling awal
- **Export_Report**: Laporan absensi dalam format file yang dapat diunduh

## Requirements

### Requirement 1: Export Absensi ke Excel

**User Story:** As an Admin_Sekolah or Wali_Kelas, I want to export attendance data to Excel format, so that I can create reports and analyze attendance patterns offline.

#### Acceptance Criteria

1. WHEN an Admin_Sekolah or Wali_Kelas requests attendance export, THE System SHALL generate an Excel file (.xlsx) containing attendance records
2. WHEN exporting attendance data, THE System SHALL allow filtering by date range (start date and end date)
3. WHEN exporting attendance data, THE System SHALL allow filtering by class
4. WHEN exporting attendance data, THE System SHALL include student information (NIS, NISN, name, class)
5. WHEN exporting attendance data, THE System SHALL include attendance details (date, check-in time, check-out time, status)
6. IF no attendance records match the filter criteria, THEN THE System SHALL return an empty file with headers only
7. WHEN generating export file, THE System SHALL name the file with format `attendance_{school}_{date_range}.xlsx`

### Requirement 2: Rekap Bulanan Absensi

**User Story:** As an Admin_Sekolah or Wali_Kelas, I want to view and export monthly attendance summary, so that I can evaluate student attendance performance.

#### Acceptance Criteria

1. WHEN an Admin_Sekolah or Wali_Kelas views monthly recap, THE System SHALL display summary per student including total days present, late, very late, and absent
2. WHEN viewing monthly recap, THE System SHALL calculate and display attendance percentage for each student
3. WHEN viewing monthly recap, THE System SHALL allow filtering by month and year
4. WHEN viewing monthly recap, THE System SHALL allow filtering by class
5. WHEN exporting monthly recap, THE System SHALL generate an Excel file with summary statistics and daily breakdown
6. WHEN displaying monthly recap, THE System SHALL sort students by attendance percentage (highest to lowest) by default
7. WHEN a Wali_Kelas views monthly recap, THE System SHALL only show students from their assigned class

### Requirement 3: Pengaturan Waktu Absensi (Multi-Schedule)

**User Story:** As an Admin_Sekolah, I want to configure multiple attendance schedules, so that the school can track attendance for different activities (morning entry, afternoon dismissal, prayer times, etc.).

#### Acceptance Criteria

1. WHEN an Admin_Sekolah creates an attendance schedule, THE System SHALL require schedule name, start time, end time, and late threshold
2. WHEN an Admin_Sekolah creates an attendance schedule, THE System SHALL allow setting very late threshold (optional)
3. WHEN an Admin_Sekolah creates an attendance schedule, THE System SHALL allow selecting active days of week
4. WHEN a student taps RFID card, THE System SHALL determine which schedule is currently active based on current time
5. WHEN recording attendance, THE System SHALL associate the attendance record with the active schedule
6. IF no schedule is active at tap time, THEN THE System SHALL use default schedule or reject the tap with appropriate message
7. WHEN an Admin_Sekolah updates a schedule, THE System SHALL not affect existing attendance records
8. WHEN an Admin_Sekolah deactivates a schedule, THE System SHALL stop using it for new attendance records
9. THE System SHALL allow maximum 10 schedules per school
10. WHEN displaying attendance records, THE System SHALL show which schedule the attendance belongs to

### Requirement 4: Real-Time Attendance Dashboard

**User Story:** As an Admin_Sekolah, Wali_Kelas, or Guru_BK, I want to monitor attendance in real-time, so that I can immediately see which students have arrived and identify those who are late or absent.

#### Acceptance Criteria

1. WHEN a user opens the live attendance page, THE System SHALL display current day's attendance statistics (present, late, very late, absent count)
2. WHEN a student taps RFID card, THE System SHALL update the dashboard within 3 seconds without page refresh
3. WHEN displaying live feed, THE System SHALL show the 20 most recent attendance records with student name, class, time, and status
4. WHEN a new attendance is recorded, THE System SHALL add it to the top of the live feed with visual highlight
5. WHEN viewing live attendance, THE System SHALL allow filtering by class
6. WHEN a Wali_Kelas views live attendance, THE System SHALL only show students from their assigned class by default
7. WHEN a Guru_BK views live attendance, THE System SHALL only show students from their assigned classes by default
8. THE System SHALL use WebSocket or Server-Sent Events for real-time updates
9. IF connection is lost, THEN THE System SHALL display connection status and attempt to reconnect automatically
10. WHEN displaying statistics, THE System SHALL show percentage of attendance completion (attended/total students)

### Requirement 5: Public Display Mode (Kiosk)

**User Story:** As an Admin_Sekolah, I want to set up a public display screen showing live attendance, so that students and staff can see attendance information on LCD monitors in public areas.

#### Acceptance Criteria

1. WHEN an Admin_Sekolah generates a display token, THE System SHALL create a unique token associated with the school
2. WHEN generating display token, THE System SHALL allow setting a name/description for the display location
3. WHEN accessing public display URL with valid token, THE System SHALL show attendance data without requiring login
4. WHEN displaying public view, THE System SHALL show live feed of recent attendance (last 10 records)
5. WHEN displaying public view, THE System SHALL show real-time statistics (present, late, absent count)
6. WHEN displaying public view, THE System SHALL show leaderboard of top 10 earliest arrivals for the day
7. WHEN displaying public view, THE System SHALL show current date and time
8. WHEN displaying public view, THE System SHALL show school name
9. WHEN a new attendance is recorded, THE System SHALL update the public display within 3 seconds
10. WHEN an Admin_Sekolah revokes a display token, THE System SHALL immediately invalidate access to that display
11. THE System SHALL allow multiple display tokens per school (for different locations)
12. WHEN displaying public view, THE System SHALL use full-screen optimized layout with large fonts
13. IF display token is invalid or revoked, THEN THE System SHALL show an error message instead of attendance data
14. THE System SHALL NOT expose sensitive student information (only name, class, and attendance time)

### Requirement 6: Display Token Management

**User Story:** As an Admin_Sekolah, I want to manage display tokens, so that I can control which displays have access to attendance data.

#### Acceptance Criteria

1. WHEN an Admin_Sekolah views display tokens, THE System SHALL list all tokens for their school with name, status, and last accessed time
2. WHEN an Admin_Sekolah creates a display token, THE System SHALL generate a cryptographically secure random token
3. WHEN creating display token, THE System SHALL show the full token only once (similar to API key)
4. WHEN an Admin_Sekolah revokes a display token, THE System SHALL mark it as inactive immediately
5. WHEN an Admin_Sekolah regenerates a display token, THE System SHALL create a new token and invalidate the old one
6. WHEN an Admin_Sekolah deletes a display token, THE System SHALL permanently remove it from the system
7. THE System SHALL track last accessed timestamp for each display token
8. THE System SHALL allow setting optional expiration date for display tokens

