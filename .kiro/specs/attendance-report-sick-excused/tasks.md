# Implementation Plan: Attendance Report Sick & Excused Status

## Overview

Implementasi penambahan status "Sakit" dan "Izin" pada halaman Laporan Absensi Admin Sekolah. Perubahan dilakukan secara incremental dari backend ke frontend.

## Tasks

- [x] 1. Update Backend DTOs
  - [x] 1.1 Update AttendanceSummaryResponse dengan field Sick dan Excused
    - Tambah field `Sick int` dan `Excused int` pada struct
    - _Requirements: 5.1, 5.2_
  - [x] 1.2 Update ClassSummaryItem dengan field Sick dan Excused
    - Tambah field `Sick int` dan `Excused int` pada struct
    - _Requirements: 5.3, 5.4_
  - [x] 1.3 Update StudentRecapSummary dengan field TotalSick dan TotalExcused
    - Tambah field `TotalSick int` dan `TotalExcused int` pada struct
    - _Requirements: 6.1, 6.2_

- [x] 2. Update Backend Repository - Daily Summary
  - [x] 2.1 Update GetAttendanceSummary untuk menghitung sick dan excused
    - Tambah case untuk AttendanceStatusSick dan AttendanceStatusExcused
    - _Requirements: 5.1, 5.2_
  - [x] 2.2 Update GetAttendanceSummaryByClass untuk menghitung sick dan excused per kelas
    - Tambah aggregation untuk status sick dan excused
    - _Requirements: 5.3, 5.4_
  - [ ]* 2.3 Write property test untuk daily summary counts consistency
    - **Property 1: Daily Summary Counts Consistency**
    - **Validates: Requirements 1.3, 1.4, 2.1, 2.2**

- [x] 3. Update Backend Repository - Monthly Recap
  - [x] 3.1 Update GetMonthlyRecap untuk menghitung total_sick dan total_excused per siswa
    - Tambah case untuk AttendanceStatusSick dan AttendanceStatusExcused dalam loop status counts
    - _Requirements: 6.3, 6.4_
  - [ ]* 3.2 Write property test untuk monthly recap counts consistency
    - **Property 2: Monthly Recap Counts Consistency**
    - **Validates: Requirements 3.3, 3.4, 4.1, 4.2**

- [x] 4. Checkpoint - Backend Tests
  - Ensure all backend tests pass, ask the user if questions arise.

- [x] 5. Update Frontend Types
  - [x] 5.1 Update AttendanceSummary interface di types/school.ts
    - Tambah field `sick: number` dan `excused: number`
    - _Requirements: 1.3, 1.4, 2.1, 2.2_
  - [x] 5.2 Update StudentRecapSummary interface di types/export.ts
    - Tambah field `total_sick: number` dan `total_excused: number`
    - _Requirements: 3.3, 3.4, 4.1, 4.2_
  - [x] 5.3 Update ATTENDANCE_STATUS_CONFIG di types/attendance.ts
    - Tambah config untuk sick dan excused dengan label dan color
    - _Requirements: 7.1, 7.2_

- [x] 6. Update AttendanceReport.vue - Tab Harian
  - [x] 6.1 Tambah Summary Cards untuk Sakit dan Izin
    - Tambah 2 Card components dengan Statistic untuk sick dan excused
    - Gunakan warna biru untuk Sakit dan cyan untuk Izin
    - _Requirements: 1.1, 1.2_
  - [x] 6.2 Update dailyStats computed untuk menghitung totalSick dan totalExcused
    - Tambah reduce untuk sick dan excused dari filteredDailyAttendance
    - _Requirements: 1.3, 1.4_
  - [x] 6.3 Tambah kolom Sakit dan Izin pada dailyColumns
    - Tambah column definition untuk sick dan excused
    - _Requirements: 2.1, 2.2_
  - [x] 6.4 Update table bodyCell template untuk render sick dan excused
    - Tambah template untuk column.key === 'sick' dan 'excused'
    - Gunakan Tag component dengan warna yang sesuai
    - _Requirements: 2.3, 2.4_
  - [x] 6.5 Update table summary footer untuk menampilkan total sick dan excused
    - Tambah Table.Summary.Cell untuk totalSick dan totalExcused
    - _Requirements: 2.5_

- [x] 7. Update AttendanceReport.vue - Tab Bulanan
  - [x] 7.1 Tambah Summary Cards untuk Rata-rata Sakit dan Izin
    - Tambah 2 Card components dengan Statistic untuk avgSick dan avgExcused
    - _Requirements: 3.1, 3.2_
  - [x] 7.2 Update monthlyStats computed untuk menghitung avgSick dan avgExcused
    - Tambah reduce dan average calculation untuk sick dan excused
    - _Requirements: 3.3, 3.4_
  - [x] 7.3 Tambah kolom Sakit dan Izin pada monthlyColumns
    - Tambah column definition untuk total_sick dan total_excused
    - _Requirements: 4.1, 4.2_
  - [x] 7.4 Update table bodyCell template untuk render total_sick dan total_excused
    - Tambah template untuk column.key === 'total_sick' dan 'total_excused'
    - _Requirements: 4.3, 4.4_

- [x] 8. Update AttendanceReport.vue - Tab Detail Filter
  - [x] 8.1 Tambah opsi Sakit dan Izin pada statusOptions
    - Tambah { value: 'sick', label: 'Sakit' } dan { value: 'excused', label: 'Izin' }
    - _Requirements: 7.1, 7.2_
  - [x] 8.2 Update getStatusTag function untuk handle sick dan excused
    - Tambah case untuk sick dan excused dengan color dan label yang sesuai
    - _Requirements: 7.3, 7.4_

- [x] 9. Checkpoint - Frontend Integration
  - Ensure frontend compiles without errors, ask the user if questions arise.

- [ ]* 10. Write Integration Tests
  - [ ]* 10.1 Write test untuk filter status sick mengembalikan data yang benar
    - **Property 4: Filter Results Correctness**
    - **Validates: Requirements 7.3, 7.4**
  - [ ]* 10.2 Write test untuk export data include sick dan excused
    - **Property 5: Export Data Completeness**
    - **Validates: Requirements 8.1, 8.2, 8.3, 8.4**

- [x] 11. Final Checkpoint
  - Ensure all tests pass and application works correctly, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Backend changes harus dilakukan terlebih dahulu karena frontend bergantung pada API response
- Database schema tidak perlu diubah karena status sick dan excused sudah ada
- Property tests validate universal correctness properties
- Unit tests validate specific examples and edge cases
