# Requirements Document

## Introduction

Sistem SaaS Manajemen Sekolah adalah platform multi-tenant berbasis web admin dan mobile app yang dirancang untuk mendukung pendataan sekolah, mengelola absensi siswa berbasis RFID, mencatat perkembangan akademik (nilai), serta menjalankan sistem konseling (BK). Sistem ini menempatkan transparansi kepada orang tua sebagai nilai utama melalui notifikasi real-time dan monitoring terstruktur.

## Glossary

- **Tenant**: Sekolah yang terdaftar dalam sistem SaaS
- **Super_Admin**: Administrator platform yang mengelola seluruh tenant
- **Admin_Sekolah**: Administrator yang mengelola operasional satu sekolah
- **Guru_BK**: Guru Bimbingan Konseling yang mengelola data perilaku siswa
- **Wali_Kelas**: Guru yang bertanggung jawab atas satu kelas
- **NISN**: Nomor Induk Siswa Nasional (identifier unik siswa)
- **NIS**: Nomor Induk Siswa (identifier internal sekolah)
- **RFID**: Radio Frequency Identification untuk absensi
- **ESP32**: Microcontroller untuk membaca kartu RFID
- **FCM**: Firebase Cloud Messaging untuk push notification
- **JWT**: JSON Web Token untuk autentikasi
- **BK**: Bimbingan Konseling

## Requirements

### Requirement 1: Multi-Tenant Management

**User Story:** As a Super_Admin, I want to create and manage schools (tenants), so that one platform can be used by many schools with isolated and secure data.

#### Acceptance Criteria

1. WHEN a Super_Admin creates a new tenant, THE System SHALL generate a unique school_id and isolate all data under that tenant
2. WHEN a Super_Admin views the tenant list, THE System SHALL display all registered schools with their status and basic information
3. WHEN a Super_Admin deactivates a tenant, THE System SHALL prevent all users under that tenant from accessing the system
4. WHEN data is queried, THE System SHALL filter results by school_id to ensure tenant isolation
5. IF a request attempts to access data from another tenant, THEN THE System SHALL reject the request with an authorization error

### Requirement 2: Device Management (RFID/ESP32)

**User Story:** As a Super_Admin, I want to manage RFID devices (ESP32) centrally, so that each device has a valid API key and can be controlled securely.

#### Acceptance Criteria

1. WHEN a Super_Admin registers a new device, THE System SHALL generate a unique API key for that device
2. WHEN a device sends attendance data, THE System SHALL validate the API key before processing
3. IF an invalid API key is used, THEN THE System SHALL reject the request and log the attempt
4. WHEN a Super_Admin revokes a device API key, THE System SHALL immediately invalidate that key
5. WHEN a Super_Admin views devices, THE System SHALL display device status, school assignment, and last activity

### Requirement 3: School Data Management

**User Story:** As an Admin_Sekolah, I want to map school data such as classes, students, and parents, so that the entire system can run with neat and consistent data structure.

#### Acceptance Criteria

1. WHEN an Admin_Sekolah creates a class, THE System SHALL associate it with the school tenant and allow student assignment
2. WHEN an Admin_Sekolah registers a student, THE System SHALL require NIS, NISN, name, and class assignment
3. WHEN an Admin_Sekolah registers a parent, THE System SHALL link the parent to one or more students
4. WHEN a student is assigned to a class, THE System SHALL update the class roster and notify the Wali_Kelas
5. IF duplicate NISN is detected within the system, THEN THE System SHALL reject the registration with an error message

### Requirement 4: User Account Management

**User Story:** As an Admin_Sekolah, I want to manage accounts for teachers, homeroom teachers, and counseling teachers, so that each role has access appropriate to their responsibilities.

#### Acceptance Criteria

1. WHEN an Admin_Sekolah creates a user account, THE System SHALL assign a role (guru, wali_kelas, guru_bk) and generate initial credentials
2. WHEN a user logs in for the first time, THE System SHALL require password change
3. WHEN an Admin_Sekolah assigns Wali_Kelas role, THE System SHALL require class assignment
4. WHEN an Admin_Sekolah deactivates an account, THE System SHALL revoke all active sessions immediately
5. THE System SHALL enforce role-based access control for all protected resources

### Requirement 5: RFID Attendance System

**User Story:** As an Admin_Sekolah, I want to monitor student attendance in summary, so that I can know the school condition without entering sensitive counseling details.

#### Acceptance Criteria

1. WHEN a student taps RFID card, THE ESP32 SHALL send student identifier and timestamp to the backend
2. WHEN the backend receives attendance data, THE System SHALL record check-in or check-out based on existing records for that day
3. WHEN attendance is recorded, THE System SHALL trigger an asynchronous notification to the parent
4. WHEN an Admin_Sekolah views attendance dashboard, THE System SHALL display summary statistics (present, absent, late)
5. IF RFID system fails, THEN THE System SHALL allow manual attendance entry by authorized users
6. WHEN parsing attendance data from ESP32, THE System SHALL validate it against the expected JSON format
7. THE Attendance_Formatter SHALL format attendance records with timestamp and student information for display

### Requirement 6: Violation Recording (BK)

**User Story:** As a Guru_BK, I want to record student violations, so that student behavior is officially documented and can be followed up.

#### Acceptance Criteria

1. WHEN a Guru_BK records a violation, THE System SHALL require category, level, description, and student identifier
2. WHEN a violation is saved, THE System SHALL trigger notification to the parent with summary (not full details)
3. WHEN a Guru_BK views violation history, THE System SHALL display all violations for a student sorted by date
4. THE System SHALL categorize violations by severity level (ringan, sedang, berat)
5. WHEN a Wali_Kelas views BK data, THE System SHALL provide read-only access to violations

### Requirement 7: Achievement Recording (BK)

**User Story:** As a Guru_BK, I want to record student achievements with points, so that students are not only judged by mistakes but also by accomplishments.

#### Acceptance Criteria

1. WHEN a Guru_BK records an achievement, THE System SHALL require title, point value, and description
2. WHEN an achievement is saved, THE System SHALL add points to student's accumulated achievement score
3. WHEN viewing student profile, THE System SHALL display total achievement points
4. WHEN an achievement is recorded, THE System SHALL trigger notification to the parent
5. THE System SHALL maintain achievement history per student

### Requirement 8: School Exit Permit

**User Story:** As a Guru_BK, I want to record student exit permits during school hours and generate permit documents, so that permits have administrative evidence that can be seen by parents.

#### Acceptance Criteria

1. WHEN a Guru_BK creates an exit permit, THE System SHALL require reason, exit time, and responsible teacher
2. WHEN a permit is created, THE System SHALL generate a PDF/receipt document with student info, reason, and timestamp
3. WHEN a permit is issued, THE System SHALL trigger notification to the parent with permit details
4. WHEN a student returns, THE System SHALL allow recording of return time
5. THE Permit_Document SHALL contain: student name, NIS/NISN, class, reason, exit time, responsible teacher, and timestamp

### Requirement 9: Counseling Notes

**User Story:** As a Guru_BK, I want to create counseling notes that are internal and summaries for parents, so that the counseling process remains professional without exposing sensitive information.

#### Acceptance Criteria

1. WHEN a Guru_BK creates a counseling note, THE System SHALL require internal_note and optional parent_summary
2. WHEN parent_summary is provided, THE System SHALL make it visible to parents via mobile app
3. THE System SHALL keep internal_note private and accessible only to Guru_BK
4. WHEN a Wali_Kelas views counseling data, THE System SHALL show only parent_summary (not internal notes)
5. THE System SHALL maintain counseling history per student with timestamps

### Requirement 10: Grade Input by Homeroom Teacher

**User Story:** As a Wali_Kelas, I want to input student grades periodically, so that parents can monitor academic progress before the official report card.

#### Acceptance Criteria

1. WHEN a Wali_Kelas inputs a grade, THE System SHALL require title, score, and optional description
2. WHEN a grade is saved, THE System SHALL associate it with the student and the Wali_Kelas who created it
3. WHEN a grade is recorded, THE System SHALL optionally trigger notification to the parent
4. WHEN parents view grades, THE System SHALL display all grades for their child sorted by date
5. THE System SHALL validate that Wali_Kelas can only input grades for students in their assigned class

### Requirement 11: Homeroom Teacher Notes

**User Story:** As a Wali_Kelas, I want to write homeroom notes that can be read by parents, so that communication about student development becomes more structured compared to manual liaison books.

#### Acceptance Criteria

1. WHEN a Wali_Kelas creates a note, THE System SHALL require content and associate it with a student
2. WHEN a note is saved, THE System SHALL trigger notification to the parent
3. WHEN parents view notes, THE System SHALL display all homeroom notes for their child
4. THE System SHALL validate that Wali_Kelas can only create notes for students in their assigned class
5. THE System SHALL maintain note history with timestamps and author information

### Requirement 12: Parent Mobile App Authentication

**User Story:** As a Parent, I want to login using my child's NISN, so that I don't need a complicated account registration process.

#### Acceptance Criteria

1. WHEN a parent enters NISN and password, THE System SHALL authenticate and return JWT tokens
2. WHEN authentication succeeds, THE System SHALL return access to all linked children's data
3. IF NISN is not found, THEN THE System SHALL return appropriate error message
4. WHEN a parent has multiple children, THE System SHALL allow switching between children's profiles
5. THE System SHALL enforce password reset on first login

### Requirement 13: Real-time Attendance Notification

**User Story:** As a Parent, I want to receive notifications when my child arrives or leaves school, so that I feel at ease knowing my child's whereabouts.

#### Acceptance Criteria

1. WHEN a student checks in, THE System SHALL send push notification to parent's device via FCM
2. WHEN a student checks out, THE System SHALL send push notification to parent's device via FCM
3. THE Notification SHALL contain student name, time, and attendance type (masuk/pulang)
4. IF parent's device is offline, THEN THE System SHALL queue the notification for delivery
5. WHEN parent opens the app, THE System SHALL display notification history

### Requirement 14: Parent Information Access

**User Story:** As a Parent, I want to know if my child commits a violation, receives an achievement, or gets an exit permit, so that I can take action or follow up communication with the school.

#### Acceptance Criteria

1. WHEN a violation is recorded, THE System SHALL send notification to parent with summary
2. WHEN an achievement is recorded, THE System SHALL send notification to parent
3. WHEN an exit permit is issued, THE System SHALL send notification with permit details and document link
4. WHEN parent views BK information, THE System SHALL display violations, achievements, and permits for their child
5. THE System SHALL only show parent_summary for counseling notes (not internal notes)

### Requirement 15: Academic Monitoring for Parents

**User Story:** As a Parent, I want to read grades and homeroom notes, so that I know my child's academic progress and attitude at school.

#### Acceptance Criteria

1. WHEN parent opens grades section, THE System SHALL display all grades for their child
2. WHEN parent opens notes section, THE System SHALL display all homeroom notes for their child
3. WHEN new grade is added, THE System SHALL optionally notify the parent
4. WHEN new note is added, THE System SHALL notify the parent
5. THE System SHALL display grades and notes sorted by most recent first

### Requirement 16: Student Self-Monitoring (MVP+)

**User Story:** As a Student, I want to view my attendance history, grades, achievements, and violations, so that I can understand my responsibilities and personal development.

#### Acceptance Criteria

1. WHEN a student logs in, THE System SHALL authenticate using NISN and password
2. WHEN student views profile, THE System SHALL display personal information and summary statistics
3. WHEN student views attendance, THE System SHALL display attendance history
4. WHEN student views grades, THE System SHALL display all grades
5. WHEN student views BK info, THE System SHALL display achievements and violations (summary only)

### Requirement 17: Notification System

**User Story:** As a System, I want to send asynchronous notifications for all important events, so that parents and students are informed in real-time without blocking operations.

#### Acceptance Criteria

1. WHEN an event triggers notification, THE System SHALL queue the notification in Redis
2. WHEN a worker processes the queue, THE System SHALL send notification via FCM
3. THE System SHALL store notification history in database with read status
4. WHEN user views notifications, THE System SHALL display all notifications with read/unread status
5. IF FCM delivery fails, THEN THE System SHALL retry with exponential backoff
6. THE Notification_Serializer SHALL serialize notification objects to JSON for storage
7. THE Notification_Deserializer SHALL deserialize JSON to notification objects when retrieved

### Requirement 18: Data Persistence

**User Story:** As a System, I want to persist all data reliably, so that school information is never lost and can be retrieved accurately.

#### Acceptance Criteria

1. WHEN any entity is created, THE System SHALL store it in PostgreSQL with timestamps
2. WHEN any entity is updated, THE System SHALL update the record and maintain audit trail
3. THE System SHALL use GORM for all database operations
4. WHEN storing objects to database, THE System SHALL serialize complex fields to JSON format
5. WHEN retrieving objects from database, THE System SHALL deserialize JSON fields to proper objects
6. FOR ALL valid data objects, serializing then deserializing SHALL produce an equivalent object (round-trip property)
