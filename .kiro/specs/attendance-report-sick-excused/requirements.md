# Requirements Document

## Introduction

Fitur ini menambahkan tampilan status "Sakit" dan "Izin" pada halaman Laporan Absensi untuk role Admin Sekolah. Saat ini, status tersebut sudah tersedia di backend (`sick` dan `excused`) dan ditampilkan di tab Detail Absensi, namun belum ditampilkan pada tab Harian dan Rekap Bulanan. Enhancement ini akan memberikan visibilitas lengkap terhadap semua jenis status kehadiran siswa.

## Glossary

- **Admin_Sekolah**: Administrator sekolah yang memiliki akses ke laporan absensi
- **Attendance_Report**: Halaman laporan absensi yang terdiri dari tab Harian, Rekap Bulanan, dan Detail Absensi
- **Daily_Summary**: Ringkasan kehadiran harian per kelas
- **Monthly_Recap**: Rekap kehadiran bulanan per siswa
- **Sick_Status**: Status kehadiran "Sakit" (sick) - siswa tidak hadir karena sakit
- **Excused_Status**: Status kehadiran "Izin" (excused) - siswa tidak hadir dengan izin
- **Summary_Card**: Komponen kartu statistik yang menampilkan ringkasan data
- **Class_Summary_Table**: Tabel yang menampilkan ringkasan kehadiran per kelas

## Requirements

### Requirement 1: Tampilan Status Sakit dan Izin pada Summary Cards Harian

**User Story:** As an Admin_Sekolah, I want to see sick and excused counts in daily summary cards, so that I can quickly understand the complete attendance breakdown for the day.

#### Acceptance Criteria

1. WHEN the Admin_Sekolah views the daily attendance tab, THE Daily_Summary SHALL display a summary card for "Sakit" count with blue color styling
2. WHEN the Admin_Sekolah views the daily attendance tab, THE Daily_Summary SHALL display a summary card for "Izin" count with cyan color styling
3. WHEN attendance data is loaded, THE System SHALL calculate total sick count from all classes
4. WHEN attendance data is loaded, THE System SHALL calculate total excused count from all classes

### Requirement 2: Tampilan Status Sakit dan Izin pada Tabel Harian

**User Story:** As an Admin_Sekolah, I want to see sick and excused columns in the daily attendance table, so that I can see the breakdown per class.

#### Acceptance Criteria

1. WHEN the daily attendance table is displayed, THE Class_Summary_Table SHALL include a "Sakit" column showing sick count per class
2. WHEN the daily attendance table is displayed, THE Class_Summary_Table SHALL include an "Izin" column showing excused count per class
3. WHEN a class has sick students, THE System SHALL display the count with appropriate tag styling
4. WHEN a class has excused students, THE System SHALL display the count with appropriate tag styling
5. WHEN the table footer summary is displayed, THE System SHALL show total sick and excused counts

### Requirement 3: Tampilan Status Sakit dan Izin pada Summary Cards Bulanan

**User Story:** As an Admin_Sekolah, I want to see average sick and excused days in monthly summary cards, so that I can understand monthly attendance patterns.

#### Acceptance Criteria

1. WHEN the Admin_Sekolah views the monthly recap tab, THE Monthly_Recap SHALL display a summary card for "Rata-rata Sakit" with blue color styling
2. WHEN the Admin_Sekolah views the monthly recap tab, THE Monthly_Recap SHALL display a summary card for "Rata-rata Izin" with cyan color styling
3. WHEN monthly recap data is loaded, THE System SHALL calculate average sick days across all students
4. WHEN monthly recap data is loaded, THE System SHALL calculate average excused days across all students

### Requirement 4: Tampilan Status Sakit dan Izin pada Tabel Rekap Bulanan

**User Story:** As an Admin_Sekolah, I want to see sick and excused columns in the monthly recap table, so that I can see each student's sick and excused days.

#### Acceptance Criteria

1. WHEN the monthly recap table is displayed, THE System SHALL include a "Sakit" column showing total sick days per student
2. WHEN the monthly recap table is displayed, THE System SHALL include an "Izin" column showing total excused days per student
3. WHEN a student has sick days, THE System SHALL display the count with blue tag styling
4. WHEN a student has excused days, THE System SHALL display the count with cyan tag styling

### Requirement 5: Backend API Update untuk Summary Harian

**User Story:** As a System, I want to include sick and excused counts in the daily summary API response, so that the frontend can display complete attendance data.

#### Acceptance Criteria

1. WHEN the daily attendance summary API is called, THE Backend SHALL return sick count in the response
2. WHEN the daily attendance summary API is called, THE Backend SHALL return excused count in the response
3. WHEN calculating class summary, THE Backend SHALL aggregate sick status counts per class
4. WHEN calculating class summary, THE Backend SHALL aggregate excused status counts per class

### Requirement 6: Backend API Update untuk Rekap Bulanan

**User Story:** As a System, I want to include sick and excused totals in the monthly recap API response, so that the frontend can display complete monthly data.

#### Acceptance Criteria

1. WHEN the monthly recap API is called, THE Backend SHALL return total_sick in each student's recap summary
2. WHEN the monthly recap API is called, THE Backend SHALL return total_excused in each student's recap summary
3. WHEN calculating student recap, THE Backend SHALL count all attendance records with sick status
4. WHEN calculating student recap, THE Backend SHALL count all attendance records with excused status

### Requirement 7: Filter Status pada Tab Detail

**User Story:** As an Admin_Sekolah, I want to filter by sick and excused status in the detail tab, so that I can find specific attendance records.

#### Acceptance Criteria

1. WHEN the status filter dropdown is displayed, THE System SHALL include "Sakit" option
2. WHEN the status filter dropdown is displayed, THE System SHALL include "Izin" option
3. WHEN Admin_Sekolah selects "Sakit" filter, THE System SHALL show only attendance records with sick status
4. WHEN Admin_Sekolah selects "Izin" filter, THE System SHALL show only attendance records with excused status

### Requirement 8: Export Data dengan Status Lengkap

**User Story:** As an Admin_Sekolah, I want exported reports to include sick and excused data, so that I have complete attendance records in exported files.

#### Acceptance Criteria

1. WHEN exporting daily attendance, THE System SHALL include sick and excused columns in the export
2. WHEN exporting monthly recap, THE System SHALL include total_sick and total_excused columns per student
3. WHEN the export file is generated, THE System SHALL properly label sick status as "Sakit"
4. WHEN the export file is generated, THE System SHALL properly label excused status as "Izin"
