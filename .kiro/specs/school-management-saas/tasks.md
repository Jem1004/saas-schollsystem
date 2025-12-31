# Implementation Plan: School Management SaaS

## Overview

Implementasi Sistem SaaS Manajemen Sekolah dimulai dari Frontend (Vue 3 + Ant Design Vue) untuk Web Admin, kemudian dilanjutkan dengan Backend (Golang + Fiber). Pendekatan ini memungkinkan validasi UI/UX lebih awal dan parallel development.

## Tasks

### Phase 1: Frontend Web Admin (Vue 3 + TypeScript + Ant Design Vue)

- [x] 1. Setup Frontend Project Structure
  - [x] 1.1 Initialize Vue 3 project dengan Vite
    - Create `web-admin/` directory
    - Initialize dengan `npm create vite@latest` (Vue + TypeScript)
    - Install dependencies: ant-design-vue, vue-router, pinia, axios
    - _Requirements: Technical Context 2.1_

  - [x] 1.2 Setup project structure dan konfigurasi
    - Create folder structure: views/, components/, stores/, services/, types/
    - Setup Ant Design Vue dengan tema putih + orange
    - Setup Vue Router dengan route guards
    - Setup Pinia untuk state management
    - _Requirements: Technical Context 2.1_

  - [x] 1.3 Create base layout components
    - Create MainLayout dengan sidebar dan header
    - Create role-based sidebar menu
    - Create header dengan profil dan logout
    - _Requirements: UI/UX 5.1_

- [x] 2. Implement Authentication Pages
  - [x] 2.1 Create Login page
    - Form login dengan email/username dan password
    - Validasi form dengan Ant Design Form
    - Handle login response dan store token
    - _Requirements: 12.1_

  - [x] 2.2 Create Change Password page
    - Form untuk first-time password change
    - Validasi password strength
    - _Requirements: 12.5, 4.2_

  - [x] 2.3 Setup auth store (Pinia)
    - Store untuk user state, tokens
    - Actions untuk login, logout, refresh token
    - _Requirements: 12.1_

  - [x] 2.4 Setup route guards
    - Protected routes berdasarkan authentication
    - Role-based route access
    - _Requirements: 4.5_

- [x] 3. Implement Super Admin Pages
  - [x] 3.1 Create Super Admin Dashboard
    - Overview statistics (total schools, devices, users)
    - System status indicators
    - _Requirements: 1.3_

  - [x] 3.2 Create School (Tenant) Management page
    - Table list sekolah dengan pagination
    - Modal/drawer untuk create/edit school
    - Activate/deactivate school actions
    - _Requirements: 1.1, 1.2, 1.3_

  - [x] 3.3 Create Device Management page
    - Table list devices dengan status
    - Modal untuk register device
    - Actions: revoke API key, regenerate
    - _Requirements: 2.1, 2.4, 2.5_

- [x] 4. Checkpoint - Super Admin UI
  - Review UI dengan mock data, ask the user if questions arise.

- [x] 5. Implement Admin Sekolah Pages
  - [x] 5.1 Create Admin Sekolah Dashboard
    - Overview: total siswa, kelas, absensi hari ini
    - Quick stats dan charts
    - _Requirements: 3.4_

  - [x] 5.2 Create Class Management page
    - Table list kelas
    - CRUD operations untuk kelas
    - Assign wali kelas
    - _Requirements: 3.1_

  - [x] 5.3 Create Student Management page
    - Table list siswa dengan filter per kelas
    - Modal/drawer untuk create/edit student
    - Import bulk students (optional)
    - _Requirements: 3.2_

  - [x] 5.4 Create Parent Management page
    - Table list orang tua
    - Link parent ke student
    - _Requirements: 3.3_

  - [x] 5.5 Create User Management page
    - Table list users (guru, wali kelas, guru BK)
    - CRUD operations untuk users
    - Role assignment
    - _Requirements: 4.1, 4.3_

  - [x] 5.6 Create Attendance Monitoring page
    - Daily attendance summary per kelas
    - Filter by date, class
    - Export attendance report
    - _Requirements: 5.4_

  - [x] 5.7 Create School Settings page
    - Form untuk attendance time settings
    - Notification toggles
    - Academic year settings
    - _Requirements: School Settings_

- [x] 6. Checkpoint - Admin Sekolah UI
  - Review UI dengan mock data, ask the user if questions arise.

- [x] 7. Implement Guru BK Pages
  - [x] 7.1 Create Guru BK Dashboard
    - Overview: recent violations, achievements
    - Students requiring attention
    - _Requirements: 6.1, 7.1_

  - [x] 7.2 Create Student BK Profile page
    - Tabs: Pelanggaran, Prestasi, Izin, Konseling
    - Timeline view per student
    - _Requirements: 6.3, 7.5, 8.4, 9.5_

  - [x] 7.3 Create Violation Management
    - Form input pelanggaran
    - Category dan level selection
    - History table
    - _Requirements: 6.1, 6.3, 6.4_

  - [x] 7.4 Create Achievement Management
    - Form input prestasi dengan point
    - Total points display
    - History table
    - _Requirements: 7.1, 7.2, 7.3_

  - [x] 7.5 Create Exit Permit Management
    - Form input izin keluar
    - Generate permit document (preview)
    - Record return time
    - _Requirements: 8.1, 8.2, 8.4_

  - [x] 7.6 Create Counseling Notes Management
    - Form dengan internal note dan parent summary
    - Clear separation UI
    - History view
    - _Requirements: 9.1, 9.2_

- [x] 8. Checkpoint - Guru BK UI
  - Review UI dengan mock data, ask the user if questions arise.

- [x] 9. Implement Wali Kelas Pages
  - [x] 9.1 Create Wali Kelas Dashboard
    - Class overview: attendance today, recent grades
    - Students list
    - _Requirements: 10.1, 11.1_

  - [x] 9.2 Create Class Attendance page
    - Daily attendance view untuk kelas
    - Manual attendance input
    - _Requirements: 5.4, 5.5_

  - [x] 9.3 Create Grade Input page
    - Form input nilai per student
    - Batch input option
    - Grade history per student
    - _Requirements: 10.1, 10.2_

  - [x] 9.4 Create Homeroom Notes page
    - Form input catatan wali kelas
    - Notes history per student
    - _Requirements: 11.1, 11.5_

  - [x] 9.5 Create BK Monitoring page (Read-only)
    - View violations, achievements, permits
    - No edit capability
    - _Requirements: 6.5_

- [x] 10. Checkpoint - Wali Kelas UI
  - Review UI dengan mock data, ask the user if questions arise.

- [x] 11. Implement Shared Components
  - [x] 11.1 Create reusable table components
    - DataTable dengan pagination, sorting, filtering
    - Export functionality
    - _Requirements: UI/UX_

  - [x] 11.2 Create form components
    - StudentSelect, ClassSelect, TeacherSelect
    - DatePicker dengan locale Indonesia
    - _Requirements: UI/UX_

  - [x] 11.3 Create notification components
    - Toast notifications
    - Notification dropdown di header
    - _Requirements: 17.4_

  - [x] 11.4 Create permit document preview
    - PDF preview component
    - Print functionality
    - _Requirements: 8.2, 8.5_

- [x] 12. Setup API Services
  - [x] 12.1 Create axios instance dengan interceptors
    - Base URL configuration
    - Token injection
    - Error handling
    - Refresh token logic
    - _Requirements: Technical Context_

  - [x] 12.2 Create API service modules
    - authService, schoolService, studentService
    - attendanceService, bkService, gradeService
    - notificationService, settingsService
    - _Requirements: All API endpoints_

  - [x] 12.3 Create TypeScript types/interfaces
    - Define all data models as TypeScript interfaces
    - API request/response types
    - _Requirements: Technical Context_

- [x] 13. Implement UX Guards dan Security Features
  - [x] 13.1 Implement Sensitive Data UX Guards
    - Blur/mask sensitive data (internal counseling notes)
    - Confirmation dialogs untuk akses data sensitif
    - Visual indicators untuk data confidential
    - _Requirements: 9.3, 9.4_

  - [x] 13.2 Implement Strict Role-Based Navigation & Guards
    - Hide menu items berdasarkan role
    - Redirect unauthorized access ke appropriate page
    - Show "Access Denied" page untuk unauthorized routes
    - Disable actions yang tidak diizinkan
    - _Requirements: 4.5, 6.5_

  - [x] 13.3 Implement Read-Only Mode Indicators
    - Visual badge/indicator untuk read-only views
    - Disable form inputs untuk read-only mode
    - Tooltip explaining why user can only view
    - Wali Kelas BK view dengan clear read-only indicator
    - _Requirements: 6.5_

  - [x] 13.4 Implement Empty & First-Time States
    - Empty state illustrations dan messages
    - First-time user onboarding hints
    - "No data yet" states dengan helpful actions
    - Welcome screen untuk new users
    - _Requirements: UI/UX_

  - [x] 13.5 Implement Global Error & Permission UI
    - Global error boundary component
    - Permission denied modal/page
    - Network error handling dengan retry option
    - Session expired handling dengan re-login prompt
    - _Requirements: Error Handling_

  - [x] 13.6 Implement Form Safety & Validation
    - Unsaved changes warning (beforeunload)
    - Form validation dengan clear error messages
    - Required field indicators
    - Input sanitization untuk XSS prevention
    - Confirmation dialogs untuk destructive actions
    - _Requirements: 3.2, 6.1, 7.1, 8.1, 9.1, 10.1, 11.1_

  - [ ]* 13.7 Implement Audit Awareness (Optional)
    - Show "Last modified by" information
    - Activity log viewer untuk admins
    - Change history untuk critical data
    - _Requirements: 18.2_

  - [ ]* 13.8 Implement Manual Override Warning (Optional)
    - Warning dialog untuk manual attendance override
    - Reason input untuk manual overrides
    - Log manual actions dengan justification
    - _Requirements: 5.5_

  - [ ]* 13.9 Implement Export Restrictions (Optional)
    - Role-based export permissions
    - Watermark pada exported documents
    - Export audit logging
    - Limit export data scope berdasarkan role
    - _Requirements: Security_

- [x] 14. Checkpoint - Frontend Complete
  - Full UI review dengan mock data dan UX guards, ask the user if questions arise.

### Phase 2: Backend API (Golang + Fiber)

- [x] 15. Setup Backend Project Structure
  - [x] 15.1 Initialize Go module dan setup project structure
    - Create `backend/` directory dengan struktur modular monolith
    - Initialize `go.mod` dengan dependencies (fiber, gorm, redis, jwt)
    - Setup `cmd/server/main.go` entry point
    - _Requirements: Technical Context 2.1_

  - [x] 15.2 Implement database connection dan migrations
    - Create `internal/shared/database/` package
    - Implement PostgreSQL connection dengan GORM
    - Create migration files untuk semua models
    - _Requirements: 18.1, 18.3_

  - [x] 15.3 Implement Redis connection
    - Create `internal/shared/redis/` package
    - Implement Redis client untuk queue dan cache
    - _Requirements: 17.1_

  - [x] 15.4 Implement configuration management
    - Create `internal/config/config.go`
    - Load config dari environment variables
    - _Requirements: Technical Context_

- [x] 16. Implement Domain Models
  - [x] 16.1 Create all domain models
    - User, School, Class, Student, Parent
    - Attendance, Violation, Achievement, Permit, CounselingNote
    - Grade, HomeroomNote, Device, Notification, SchoolSettings
    - _Requirements: All data models_

  - [ ]* 16.2 Write property test for data serialization round-trip
    - **Property 9: Data Serialization Round-Trip**
    - **Validates: Requirements 17.6, 17.7, 18.4, 18.5, 18.6**

- [x] 17. Implement Auth Module
  - [x] 17.1 Create auth repository, service, handler
    - JWT token generation dan validation
    - Password hashing dengan bcrypt
    - Login, refresh, logout, change password
    - _Requirements: 12.1, 12.3, 12.5_

  - [x] 17.2 Implement auth dan tenant middleware
    - JWT validation middleware
    - Tenant isolation middleware
    - _Requirements: 1.4, 1.5, 4.5_

  - [ ]* 17.3 Write property test for tenant isolation
    - **Property 1: Tenant Data Isolation**
    - **Validates: Requirements 1.4, 1.5**

- [x] 18. Implement Tenant Module
  - [x] 18.1 Create tenant repository, service, handler
    - School CRUD operations
    - Tenant activation/deactivation
    - _Requirements: 1.1, 1.2, 1.3_

  - [ ]* 18.2 Write property test for tenant creation uniqueness
    - **Property 2: Tenant Creation Uniqueness**
    - **Validates: Requirements 1.1**

- [x] 19. Implement School Module
  - [x] 19.1 Create school repository, service, handler
    - Class, Student, Parent management
    - NISN uniqueness validation
    - _Requirements: 3.1, 3.2, 3.3, 3.5_

  - [ ]* 19.2 Write property test for NISN uniqueness
    - **Property 5: NISN Uniqueness**
    - **Validates: Requirements 3.5**

  - [ ]* 19.3 Write property test for required field validation
    - **Property 12: Required Field Validation**
    - **Validates: Requirements 3.2, 6.1, 7.1, 8.1, 9.1, 10.1, 11.1**

- [x] 20. Implement RBAC
  - [x] 20.1 Create RBAC middleware dan access policy
    - Role-based access control
    - Resource-level permissions
    - _Requirements: 4.5, 6.5, 9.3, 9.4, 10.5, 11.4_

  - [ ]* 20.2 Write property test for RBAC enforcement
    - **Property 6: Role-Based Access Control Enforcement**
    - **Validates: Requirements 4.5, 6.5, 9.3, 9.4, 10.5, 11.4, 14.5**

- [x] 21. Checkpoint - Core Backend
  - Ensure all tests pass, ask the user if questions arise.

- [x] 22. Implement Device Module
  - [x] 22.1 Create device repository, service, handler
    - Device registration dengan API key
    - API key validation dan revocation
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5_

  - [ ]* 22.2 Write property test for API key security
    - **Property 4: Device API Key Security**
    - **Validates: Requirements 2.1, 2.2, 2.3, 2.4**

- [x] 23. Implement Attendance Module
  - [x] 23.1 Create attendance repository, service, handler
    - RFID attendance recording
    - Check-in/check-out logic
    - Manual attendance fallback
    - _Requirements: 5.2, 5.4, 5.5_

  - [x] 23.2 Create attendance policy
    - Status determination based on school settings
    - _Requirements: 5.2, School Settings_

  - [ ]* 23.3 Write property test for attendance logic
    - **Property 7: Attendance Check-In/Check-Out Logic**
    - **Validates: Requirements 5.2**

- [x] 24. Implement BK Module
  - [x] 24.1 Create BK repositories dan services
    - Violations, Achievements, Permits, Counseling Notes
    - Point accumulation
    - Permit document generation
    - _Requirements: 6.1-6.5, 7.1-7.5, 8.1-8.5, 9.1-9.5_

  - [x] 24.2 Create BK handler
    - All BK endpoints
    - _Requirements: 6.1-6.5, 7.1-7.5, 8.1-8.5, 9.1-9.5_

  - [ ]* 24.3 Write property test for achievement points
    - **Property 8: Achievement Point Accumulation**
    - **Validates: Requirements 7.2, 7.3**

  - [ ]* 24.4 Write property test for permit document
    - **Property 13: Permit Document Content Completeness**
    - **Validates: Requirements 8.2, 8.5**

- [x] 25. Implement Grade dan Homeroom Module
  - [x] 25.1 Create grade repository, service, handler
    - Grade CRUD dengan teacher validation
    - _Requirements: 10.1, 10.2, 10.4, 10.5_

  - [x] 25.2 Create homeroom repository, service, handler
    - Homeroom notes CRUD
    - _Requirements: 11.1, 11.3, 11.4, 11.5_

  - [ ]* 25.3 Write property test for history maintenance
    - **Property 11: History Maintenance with Timestamps**
    - **Validates: Requirements 6.3, 7.5, 9.5, 11.5, 18.1**

  - [ ]* 25.4 Write property test for data sorting
    - **Property 16: Data Sorting by Date**
    - **Validates: Requirements 6.3, 10.4, 15.5**

- [x] 26. Implement Settings Module
  - [x] 26.1 Create settings repository, service, handler
    - School settings management
    - Default settings creation
    - _Requirements: School Settings_

  - [ ]* 26.2 Write property test for settings policy
    - **Property 17: School Settings Policy Enforcement**
    - **Validates: School Settings requirements**

- [x] 27. Checkpoint - Domain Modules
  - Ensure all tests pass, ask the user if questions arise.

- [x] 28. Implement Notification Module
  - [x] 28.1 Create notification repository, service, handler
    - Notification CRUD
    - Mark as read
    - _Requirements: 17.3, 17.4_

  - [x] 28.2 Create FCM client
    - Firebase Cloud Messaging integration
    - _Requirements: 13.1, 13.2_

  - [x] 28.3 Create notification worker
    - Background queue processing
    - Retry dengan exponential backoff
    - _Requirements: 17.1, 17.2, 17.5_

  - [ ]* 28.4 Write property test for notification queue
    - **Property 10: Notification Queue Processing**
    - **Validates: Requirements 17.1, 17.2, 17.5**

  - [ ]* 28.5 Write property test for notification content
    - **Property 14: Notification Content Completeness**
    - **Validates: Requirements 13.3**

- [x] 29. Implement Parent dan Student API
  - [x] 29.1 Create parent API handler
    - Data aggregation untuk linked children
    - _Requirements: 12.2, 14.4, 15.1, 15.2_

  - [x] 29.2 Create student API handler
    - Student self-monitoring endpoints
    - _Requirements: 16.1-16.5_

  - [ ]* 29.3 Write property test for parent data access
    - **Property 15: Parent Data Access Completeness**
    - **Validates: Requirements 12.2, 14.4, 15.1, 15.2**

- [x] 30. Integration dan Wiring
  - [x] 30.1 Wire all modules di main.go
    - Initialize dependencies
    - Register routes
    - Start workers
    - _Requirements: All_

  - [ ]* 30.2 Write integration tests
    - Full API flow tests
    - Multi-tenant scenarios
    - _Requirements: All_

- [x] 31. Final Checkpoint
  - Ensure all tests pass, frontend-backend integration complete, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Phase 1 (Frontend) dapat dikerjakan dengan mock data terlebih dahulu
- Phase 2 (Backend) akan menyediakan real API untuk frontend
- Checkpoints memastikan validasi incremental
- Property tests validate universal correctness properties
- Frontend menggunakan Vue 3 + TypeScript + Ant Design Vue
- Backend menggunakan Golang + Fiber + GORM + Redis + FCM
- UX Guards (Task 13) memastikan keamanan dan user experience yang baik
- Total: 31 task groups dengan ~90 sub-tasks
