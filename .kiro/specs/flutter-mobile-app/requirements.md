# Requirements Document

## Introduction

Aplikasi mobile Flutter untuk Sistem SaaS Manajemen Sekolah yang ditujukan untuk dua jenis pengguna: **Orang Tua** dan **Siswa**. Aplikasi ini berfungsi sebagai platform monitoring real-time yang memungkinkan orang tua memantau kehadiran, nilai akademik, dan informasi BK (Bimbingan Konseling) anak mereka, serta siswa untuk melihat data perkembangan diri sendiri.

Aplikasi menggunakan desain modern minimalis dengan warna primary orange dan putih, terintegrasi dengan backend Golang yang sudah ada, dan mendukung push notification via Firebase Cloud Messaging (FCM).

## Glossary

- **Parent_App**: Aplikasi mobile untuk orang tua
- **Student_App**: Aplikasi mobile untuk siswa
- **NISN**: Nomor Induk Siswa Nasional (digunakan sebagai username login)
- **FCM**: Firebase Cloud Messaging untuk push notification
- **JWT**: JSON Web Token untuk autentikasi
- **Dashboard**: Halaman utama yang menampilkan ringkasan informasi
- **BK**: Bimbingan Konseling (counseling)
- **Child_Selector**: Komponen untuk memilih anak (jika orang tua punya lebih dari satu anak)
- **Bottom_Navigation**: Navigasi utama di bagian bawah layar
- **Pull_To_Refresh**: Gesture tarik ke bawah untuk memperbarui data
- **Skeleton_Loading**: Animasi loading berbentuk placeholder

## Requirements

### Requirement 1: Authentication

**User Story:** As a parent or student, I want to login using NISN and password, so that I can access the app without complicated registration process.

#### Acceptance Criteria

1. WHEN a user opens the app for the first time, THE App SHALL display a splash screen with logo followed by login page
2. WHEN a user enters valid NISN and password, THE App SHALL authenticate via API and store JWT tokens securely
3. WHEN authentication succeeds with must_reset_pwd flag true, THE App SHALL redirect to change password page
4. WHEN a user changes password successfully, THE App SHALL clear the flag and redirect to main screen
5. IF authentication fails, THEN THE App SHALL display appropriate error message in Indonesian
6. WHEN a user is already logged in, THE App SHALL skip login and show main screen directly
7. WHEN JWT token expires, THE App SHALL automatically refresh using refresh token
8. IF refresh token is invalid, THEN THE App SHALL redirect to login page
9. WHEN a user taps logout, THE App SHALL clear all stored tokens and redirect to login page

### Requirement 2: Push Notification (FCM)

**User Story:** As a parent, I want to receive push notifications when important events happen, so that I am immediately informed about my child's activities at school.

#### Acceptance Criteria

1. WHEN a user logs in successfully, THE App SHALL register FCM token to backend
2. WHEN a push notification arrives while app is in foreground, THE App SHALL display in-app notification banner
3. WHEN a push notification arrives while app is in background, THE App SHALL display system notification
4. WHEN a user taps on notification, THE App SHALL navigate to relevant screen based on notification type
5. WHEN a user logs out, THE App SHALL deactivate FCM token on backend
6. THE App SHALL handle notification types: attendance_in, attendance_out, violation, achievement, permit, grade, homeroom_note, counseling

### Requirement 3: Parent - Child Management

**User Story:** As a parent with multiple children, I want to easily switch between my children's profiles, so that I can monitor each child separately.

#### Acceptance Criteria

1. WHEN a parent logs in, THE App SHALL fetch list of linked children from API
2. IF parent has only one child, THEN THE App SHALL automatically select that child
3. IF parent has multiple children, THEN THE App SHALL display child selector on first login
4. WHEN parent is on dashboard, THE App SHALL display currently selected child with option to switch
5. WHEN parent switches child, THE App SHALL reload all data for the newly selected child
6. THE App SHALL persist selected child preference locally

### Requirement 4: Parent - Dashboard

**User Story:** As a parent, I want to see a summary of my child's status on the home screen, so that I can quickly understand their current situation at school.

#### Acceptance Criteria

1. WHEN parent opens dashboard, THE App SHALL display child's name, class, and school
2. WHEN parent opens dashboard, THE App SHALL display today's attendance status (check-in/check-out time)
3. WHEN parent opens dashboard, THE App SHALL display summary cards: attendance percentage, average grade, achievement points
4. WHEN parent opens dashboard, THE App SHALL display recent activities list (latest 5 items)
5. WHEN parent pulls down on dashboard, THE App SHALL refresh all data
6. WHEN data is loading, THE App SHALL display skeleton loading animation
7. IF no data available, THEN THE App SHALL display appropriate empty state with illustration

### Requirement 5: Parent - Attendance Monitoring

**User Story:** As a parent, I want to view my child's attendance history, so that I can monitor their discipline and punctuality.

#### Acceptance Criteria

1. WHEN parent opens attendance page, THE App SHALL display attendance summary (on_time, late, very_late, absent, sick, excused counts)
2. WHEN parent opens attendance page, THE App SHALL display attendance list sorted by date descending
3. WHEN parent views attendance item, THE App SHALL show date, check-in time, check-out time, and status with color indicator
4. WHEN parent filters by date range, THE App SHALL reload attendance data for selected period
5. WHEN parent scrolls to bottom, THE App SHALL load more attendance records (pagination)
6. THE App SHALL use color coding: green for on_time, yellow for late, orange for very_late, red for absent

### Requirement 6: Parent - Grade Monitoring

**User Story:** As a parent, I want to view my child's grades, so that I can monitor their academic progress.

#### Acceptance Criteria

1. WHEN parent opens grades page, THE App SHALL display grade summary (average, highest, lowest, total count)
2. WHEN parent opens grades page, THE App SHALL display grades list sorted by date descending
3. WHEN parent views grade item, THE App SHALL show title, score, description, teacher name, and date
4. WHEN parent scrolls to bottom, THE App SHALL load more grade records (pagination)
5. THE App SHALL use color coding for scores: green (>=80), yellow (60-79), red (<60)

### Requirement 7: Parent - Homeroom Notes

**User Story:** As a parent, I want to read notes from the homeroom teacher, so that I can understand my child's development and behavior at school.

#### Acceptance Criteria

1. WHEN parent opens notes page, THE App SHALL display homeroom notes list sorted by date descending
2. WHEN parent views note item, THE App SHALL show content, teacher name, and date
3. WHEN parent taps on note, THE App SHALL display full note content in detail view
4. WHEN parent scrolls to bottom, THE App SHALL load more notes (pagination)

### Requirement 8: Parent - BK Information

**User Story:** As a parent, I want to view BK information including violations, achievements, permits, and counseling summaries, so that I can be aware of my child's behavior and accomplishments.

#### Acceptance Criteria

1. WHEN parent opens BK page, THE App SHALL display tabs: Prestasi, Pelanggaran, Izin, Konseling
2. WHEN parent views Prestasi tab, THE App SHALL display achievements list with title, points, and date
3. WHEN parent views Pelanggaran tab, THE App SHALL display violations list with category, level, points, and date
4. WHEN parent views Izin tab, THE App SHALL display permits list with reason, exit time, return time, and document link
5. WHEN parent views Konseling tab, THE App SHALL display counseling notes with parent_summary only (not internal notes)
6. WHEN parent taps on permit, THE App SHALL display permit detail with option to download PDF document
7. THE App SHALL display total achievement points and violation points at top of BK page

### Requirement 9: Parent - Notification Center

**User Story:** As a parent, I want to view all notifications in one place, so that I can review past notifications and mark them as read.

#### Acceptance Criteria

1. WHEN parent opens notification page, THE App SHALL display notifications list sorted by date descending
2. WHEN parent views notification, THE App SHALL show type icon, title, message, and timestamp
3. WHEN parent taps on notification, THE App SHALL mark it as read and navigate to relevant screen
4. WHEN parent taps "Mark All Read", THE App SHALL mark all notifications as read
5. THE App SHALL display unread count badge on notification icon in bottom navigation
6. THE App SHALL visually distinguish read and unread notifications

### Requirement 10: Student - Profile & Dashboard

**User Story:** As a student, I want to view my own profile and summary statistics, so that I can monitor my own progress and development.

#### Acceptance Criteria

1. WHEN student logs in, THE App SHALL display student dashboard (not child selector)
2. WHEN student opens dashboard, THE App SHALL display profile: name, NIS, NISN, class, school
3. WHEN student opens dashboard, THE App SHALL display today's attendance status
4. WHEN student opens dashboard, THE App SHALL display summary: attendance percentage, average grade, achievement points
5. WHEN student opens dashboard, THE App SHALL display recent activities

### Requirement 11: Student - Self Monitoring

**User Story:** As a student, I want to view my attendance history, grades, and BK summary, so that I can be responsible for my own development.

#### Acceptance Criteria

1. WHEN student opens attendance page, THE App SHALL display own attendance history with summary
2. WHEN student opens grades page, THE App SHALL display own grades with summary
3. WHEN student opens BK page, THE App SHALL display own achievements and violations only (no permits, no counseling)
4. THE Student_App SHALL NOT display homeroom notes (only for parents)
5. THE Student_App SHALL NOT display counseling notes (only for parents)

### Requirement 12: UI/UX - Modern Minimalist Design

**User Story:** As a user, I want a clean and modern interface, so that I can easily navigate and understand the information presented.

#### Acceptance Criteria

1. THE App SHALL use primary color orange (#FF6B00) and white (#FFFFFF) as main colors
2. THE App SHALL use clean sans-serif typography (system font or Poppins/Inter)
3. THE App SHALL use rounded corners (border-radius: 12-16px) for cards and buttons
4. THE App SHALL use subtle shadows for card elevation
5. THE App SHALL use consistent spacing (8px grid system)
6. THE App SHALL use meaningful icons from Material Icons or custom SVG
7. THE App SHALL display skeleton loading animation while fetching data
8. THE App SHALL display empty state illustrations when no data available
9. THE App SHALL support pull-to-refresh on all list screens
10. THE App SHALL use bottom navigation with 5 tabs maximum

### Requirement 13: Offline & Error Handling

**User Story:** As a user, I want the app to handle errors gracefully, so that I understand what went wrong and can take appropriate action.

#### Acceptance Criteria

1. WHEN network request fails, THE App SHALL display error message with retry button
2. WHEN user has no internet connection, THE App SHALL display offline indicator
3. IF cached data is available, THEN THE App SHALL display cached data with "offline" badge
4. WHEN session expires, THE App SHALL display session expired message and redirect to login
5. THE App SHALL handle all API error codes and display appropriate Indonesian messages

### Requirement 14: Data Persistence

**User Story:** As a user, I want the app to remember my preferences and cache recent data, so that I have a faster experience.

#### Acceptance Criteria

1. THE App SHALL store JWT tokens in secure storage (encrypted)
2. THE App SHALL store selected child ID in local storage
3. THE App SHALL cache dashboard data for offline viewing
4. THE App SHALL clear all cached data on logout
5. FOR ALL cached data, storing then retrieving SHALL produce equivalent data (round-trip property)
