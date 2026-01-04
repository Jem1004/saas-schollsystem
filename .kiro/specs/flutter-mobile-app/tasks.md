# Implementation Tasks

## Task 1: Project Setup & Core Infrastructure

- [x] Create Flutter project with `flutter create --org com.schoolmanagement mobile`
- [x] Setup `pubspec.yaml` with dependencies:
  - flutter_bloc (state management)
  - dio (HTTP client)
  - get_it (dependency injection)
  - flutter_secure_storage (secure token storage)
  - shared_preferences (local storage)
  - firebase_core, firebase_messaging (FCM)
  - flutter_local_notifications
  - dartz (functional programming - Either)
  - equatable (value equality)
  - intl (date formatting)
  - cached_network_image
  - shimmer (skeleton loading)
  - pull_to_refresh
- [x] Create folder structure as defined in design.md
- [x] Setup `lib/core/config/app_config.dart` with environment configuration
- [x] Setup `lib/core/config/api_config.dart` with base URL and endpoints
- [x] Create `lib/core/constants/api_endpoints.dart`

**Requirements:** Foundation for all requirements

---

## Task 2: Theme & Design System

- [x] Create `lib/core/constants/app_colors.dart` with color palette (primary orange #FF6B00)
- [x] Create `lib/core/constants/app_typography.dart` with text styles
- [x] Create `lib/core/constants/app_spacing.dart` with spacing constants
- [x] Create `lib/core/constants/app_strings.dart` with Indonesian strings
- [x] Create `lib/core/theme/app_theme.dart` with ThemeData
- [x] Create `lib/core/theme/app_decorations.dart` with BoxDecoration styles
- [x] Add Poppins font to `assets/fonts/` and configure in pubspec.yaml

**Requirements:** REQ-12 (UI/UX Modern Minimalist Design)

---

## Task 3: Network Layer & API Client

- [-] Create `lib/core/network/api_client.dart` with Dio setup
- [ ] Create `lib/core/network/api_interceptor.dart` for JWT token injection
- [ ] Create `lib/core/network/network_info.dart` for connectivity check
- [ ] Create `lib/core/errors/exceptions.dart` with custom exceptions
- [ ] Create `lib/core/errors/failures.dart` with Failure classes
- [ ] Implement token refresh logic in interceptor
- [ ] Add request/response logging for debug mode

**Requirements:** REQ-1 (Authentication), REQ-13 (Error Handling)

---

## Task 4: Data Models

- [ ] Create `lib/data/models/user_model.dart`
- [ ] Create `lib/data/models/student_model.dart` (Child model)
- [ ] Create `lib/data/models/attendance_model.dart` with AttendanceSummary
- [ ] Create `lib/data/models/grade_model.dart` with GradeSummary
- [ ] Create `lib/data/models/homeroom_note_model.dart`
- [ ] Create `lib/data/models/achievement_model.dart`
- [ ] Create `lib/data/models/violation_model.dart`
- [ ] Create `lib/data/models/permit_model.dart`
- [ ] Create `lib/data/models/counseling_note_model.dart`
- [ ] Create `lib/data/models/notification_model.dart`
- [ ] Create `lib/data/models/bk_info_model.dart` (combined BK data)
- [ ] Add JSON serialization (fromJson/toJson) for all models

**Requirements:** REQ-4 to REQ-11 (Data display features)

---

## Task 5: Domain Layer - Entities & Repository Interfaces

- [ ] Create `lib/domain/entities/user.dart`
- [ ] Create `lib/domain/entities/child.dart`
- [ ] Create `lib/domain/entities/attendance.dart`
- [ ] Create `lib/domain/entities/grade.dart`
- [ ] Create `lib/domain/entities/bk_info.dart`
- [ ] Create `lib/domain/entities/notification.dart`
- [ ] Create `lib/domain/repositories/auth_repository.dart` interface
- [ ] Create `lib/domain/repositories/parent_repository.dart` interface
- [ ] Create `lib/domain/repositories/student_repository.dart` interface
- [ ] Create `lib/domain/repositories/notification_repository.dart` interface

**Requirements:** Foundation for all data requirements

---

## Task 6: Local Data Sources

- [ ] Create `lib/data/datasources/local/auth_local_datasource.dart`
  - Store/retrieve JWT tokens in secure storage
  - Store/retrieve selected child ID
- [ ] Create `lib/data/datasources/local/cache_local_datasource.dart`
  - Cache dashboard data
  - Cache last fetched data for offline mode
- [ ] Implement secure storage encryption for tokens

**Requirements:** REQ-14 (Data Persistence)

---

## Task 7: Remote Data Sources - Auth

- [ ] Create `lib/data/datasources/remote/auth_remote_datasource.dart`
  - POST /api/v1/auth/login
  - POST /api/v1/auth/refresh
  - POST /api/v1/auth/logout
  - POST /api/v1/auth/change-password
  - GET /api/v1/auth/me

**Requirements:** REQ-1 (Authentication)

---

## Task 8: Remote Data Sources - Parent

- [ ] Create `lib/data/datasources/remote/parent_remote_datasource.dart`
  - GET /api/v1/parent/children
  - GET /api/v1/parent/children/{id}/dashboard
  - GET /api/v1/parent/children/{id}/attendance
  - GET /api/v1/parent/children/{id}/attendance/summary
  - GET /api/v1/parent/children/{id}/grades
  - GET /api/v1/parent/children/{id}/grades/summary
  - GET /api/v1/parent/children/{id}/homeroom-notes
  - GET /api/v1/parent/children/{id}/achievements
  - GET /api/v1/parent/children/{id}/violations
  - GET /api/v1/parent/children/{id}/permits
  - GET /api/v1/parent/children/{id}/counseling-notes

**Requirements:** REQ-3 to REQ-8 (Parent features)

---

## Task 9: Remote Data Sources - Student

- [ ] Create `lib/data/datasources/remote/student_remote_datasource.dart`
  - GET /api/v1/student/profile
  - GET /api/v1/student/dashboard
  - GET /api/v1/student/attendance
  - GET /api/v1/student/attendance/summary
  - GET /api/v1/student/grades
  - GET /api/v1/student/grades/summary
  - GET /api/v1/student/achievements
  - GET /api/v1/student/violations

**Requirements:** REQ-10, REQ-11 (Student features)

---

## Task 10: Remote Data Sources - Notification

- [ ] Create `lib/data/datasources/remote/notification_remote_datasource.dart`
  - GET /api/v1/notifications
  - GET /api/v1/notifications/summary
  - GET /api/v1/notifications/unread-count
  - PUT /api/v1/notifications/{id}/read
  - PUT /api/v1/notifications/mark-all-read
  - POST /api/v1/fcm/register
  - POST /api/v1/fcm/deactivate

**Requirements:** REQ-2 (Push Notification), REQ-9 (Notification Center)

---

## Task 11: Repository Implementations

- [ ] Create `lib/data/repositories/auth_repository_impl.dart`
- [ ] Create `lib/data/repositories/parent_repository_impl.dart`
- [ ] Create `lib/data/repositories/student_repository_impl.dart`
- [ ] Create `lib/data/repositories/notification_repository_impl.dart`
- [ ] Implement error handling with Either<Failure, T>
- [ ] Implement offline fallback with cached data

**Requirements:** REQ-13 (Error Handling), REQ-14 (Data Persistence)

---

## Task 12: Use Cases

- [ ] Create auth use cases:
  - `lib/domain/usecases/auth/login_usecase.dart`
  - `lib/domain/usecases/auth/logout_usecase.dart`
  - `lib/domain/usecases/auth/change_password_usecase.dart`
- [ ] Create parent use cases:
  - `lib/domain/usecases/parent/get_children_usecase.dart`
  - `lib/domain/usecases/parent/get_child_dashboard_usecase.dart`
  - `lib/domain/usecases/parent/get_child_attendance_usecase.dart`
  - `lib/domain/usecases/parent/get_child_grades_usecase.dart`
  - `lib/domain/usecases/parent/get_child_notes_usecase.dart`
  - `lib/domain/usecases/parent/get_child_bk_usecase.dart`
- [ ] Create student use cases:
  - `lib/domain/usecases/student/get_profile_usecase.dart`
  - `lib/domain/usecases/student/get_dashboard_usecase.dart`
  - `lib/domain/usecases/student/get_attendance_usecase.dart`
  - `lib/domain/usecases/student/get_grades_usecase.dart`
  - `lib/domain/usecases/student/get_bk_usecase.dart`
- [ ] Create notification use cases:
  - `lib/domain/usecases/notification/get_notifications_usecase.dart`
  - `lib/domain/usecases/notification/mark_as_read_usecase.dart`

**Requirements:** All feature requirements

---

## Task 13: BLoCs - Auth & Child Selector

- [ ] Create `lib/presentation/blocs/auth/auth_bloc.dart`
- [ ] Create `lib/presentation/blocs/auth/auth_event.dart`
- [ ] Create `lib/presentation/blocs/auth/auth_state.dart`
- [ ] Create `lib/presentation/blocs/child_selector/child_selector_bloc.dart`
- [ ] Create `lib/presentation/blocs/child_selector/child_selector_event.dart`
- [ ] Create `lib/presentation/blocs/child_selector/child_selector_state.dart`
- [ ] Implement auto-select for single child
- [ ] Implement child persistence

**Requirements:** REQ-1 (Authentication), REQ-3 (Child Management)

---

## Task 14: BLoCs - Parent Features

- [ ] Create `lib/presentation/blocs/parent/dashboard_bloc.dart` (with events/states)
- [ ] Create `lib/presentation/blocs/parent/attendance_bloc.dart` (with events/states)
- [ ] Create `lib/presentation/blocs/parent/grades_bloc.dart` (with events/states)
- [ ] Create `lib/presentation/blocs/parent/notes_bloc.dart` (with events/states)
- [ ] Create `lib/presentation/blocs/parent/bk_bloc.dart` (with events/states)
- [ ] Implement pagination for list blocs
- [ ] Implement refresh functionality

**Requirements:** REQ-4 to REQ-8 (Parent features)

---

## Task 15: BLoCs - Student & Notification

- [ ] Create `lib/presentation/blocs/student/student_bloc.dart` (with events/states)
- [ ] Create `lib/presentation/blocs/notification/notification_bloc.dart` (with events/states)
- [ ] Implement unread count tracking
- [ ] Implement mark as read functionality

**Requirements:** REQ-10, REQ-11 (Student), REQ-9 (Notifications)

---

## Task 16: Common Widgets

- [ ] Create `lib/presentation/widgets/common/app_scaffold.dart`
- [ ] Create `lib/presentation/widgets/common/app_app_bar.dart`
- [ ] Create `lib/presentation/widgets/common/app_bottom_nav.dart`
- [ ] Create `lib/presentation/widgets/common/app_button.dart`
- [ ] Create `lib/presentation/widgets/common/app_text_field.dart`
- [ ] Create `lib/presentation/widgets/common/app_card.dart`
- [ ] Create `lib/presentation/widgets/common/loading_widget.dart`
- [ ] Create `lib/presentation/widgets/common/skeleton_widget.dart`
- [ ] Create `lib/presentation/widgets/common/error_widget.dart`
- [ ] Create `lib/presentation/widgets/common/empty_state_widget.dart`

**Requirements:** REQ-12 (UI/UX Design)

---

## Task 17: Card Widgets

- [ ] Create `lib/presentation/widgets/cards/summary_card.dart`
- [ ] Create `lib/presentation/widgets/cards/attendance_card.dart`
- [ ] Create `lib/presentation/widgets/cards/grade_card.dart`
- [ ] Create `lib/presentation/widgets/cards/note_card.dart`
- [ ] Create `lib/presentation/widgets/cards/achievement_card.dart`
- [ ] Create `lib/presentation/widgets/cards/violation_card.dart`
- [ ] Create `lib/presentation/widgets/cards/permit_card.dart`
- [ ] Create `lib/presentation/widgets/cards/activity_card.dart`
- [ ] Create `lib/presentation/widgets/cards/notification_card.dart`
- [ ] Implement color coding for status indicators

**Requirements:** REQ-5, REQ-6, REQ-8, REQ-9 (Data display with colors)

---

## Task 18: Header & Dialog Widgets

- [ ] Create `lib/presentation/widgets/headers/child_header.dart`
- [ ] Create `lib/presentation/widgets/headers/profile_header.dart`
- [ ] Create `lib/presentation/widgets/headers/attendance_status_header.dart`
- [ ] Create `lib/presentation/widgets/dialogs/child_selector_dialog.dart`
- [ ] Create `lib/presentation/widgets/dialogs/confirm_dialog.dart`
- [ ] Create `lib/presentation/widgets/dialogs/error_dialog.dart`

**Requirements:** REQ-3 (Child selector), REQ-4 (Dashboard headers)

---

## Task 19: Auth Pages

- [ ] Create `lib/presentation/pages/splash/splash_page.dart`
- [ ] Create `lib/presentation/pages/auth/login_page.dart`
- [ ] Create `lib/presentation/pages/auth/change_password_page.dart`
- [ ] Implement form validation
- [ ] Implement loading states
- [ ] Implement error display

**Requirements:** REQ-1 (Authentication)

---

## Task 20: Parent Pages

- [ ] Create `lib/presentation/pages/parent/parent_main_page.dart` (with bottom nav)
- [ ] Create `lib/presentation/pages/parent/child_selector_page.dart`
- [ ] Create `lib/presentation/pages/parent/dashboard_page.dart`
- [ ] Create `lib/presentation/pages/parent/attendance_page.dart`
- [ ] Create `lib/presentation/pages/parent/grades_page.dart`
- [ ] Create `lib/presentation/pages/parent/notes_page.dart`
- [ ] Create `lib/presentation/pages/parent/bk_page.dart` (with tabs)
- [ ] Create `lib/presentation/pages/parent/permit_detail_page.dart`
- [ ] Implement pull-to-refresh on all list pages
- [ ] Implement pagination scroll loading

**Requirements:** REQ-3 to REQ-8 (Parent features)

---

## Task 21: Student Pages

- [ ] Create `lib/presentation/pages/student/student_main_page.dart` (with bottom nav)
- [ ] Create `lib/presentation/pages/student/student_dashboard_page.dart`
- [ ] Create `lib/presentation/pages/student/student_attendance_page.dart`
- [ ] Create `lib/presentation/pages/student/student_grades_page.dart`
- [ ] Create `lib/presentation/pages/student/student_bk_page.dart` (achievements & violations only)
- [ ] Implement pull-to-refresh
- [ ] Implement pagination

**Requirements:** REQ-10, REQ-11 (Student features)

---

## Task 22: Notification & Settings Pages

- [ ] Create `lib/presentation/pages/notification/notification_page.dart`
- [ ] Create `lib/presentation/pages/settings/settings_page.dart`
- [ ] Implement notification grouping by date
- [ ] Implement mark all as read
- [ ] Implement logout confirmation

**Requirements:** REQ-9 (Notifications)

---

## Task 23: Navigation & Routing

- [ ] Create `lib/core/config/routes.dart` with named routes
- [ ] Setup `lib/app.dart` with MaterialApp and route configuration
- [ ] Implement role-based navigation (parent vs student)
- [ ] Implement deep linking for notifications
- [ ] Implement auth guard for protected routes

**Requirements:** REQ-1 (Auth flow), REQ-2 (Notification navigation)

---

## Task 24: FCM Service

- [ ] Create `lib/services/fcm_service.dart`
- [ ] Setup Firebase initialization in `main.dart`
- [ ] Implement foreground notification handling
- [ ] Implement background notification handling
- [ ] Implement notification tap handling with navigation
- [ ] Create `lib/services/local_notification_service.dart`

**Requirements:** REQ-2 (Push Notification)

---

## Task 25: Dependency Injection

- [ ] Create `lib/services/injection_container.dart`
- [ ] Register all data sources
- [ ] Register all repositories
- [ ] Register all use cases
- [ ] Register all BLoCs
- [ ] Setup GetIt initialization in main.dart

**Requirements:** Foundation for all features

---

## Task 26: Assets & Empty States

- [ ] Add app logo to `assets/images/logo.png`
- [ ] Add empty state illustrations:
  - `assets/images/empty_attendance.svg`
  - `assets/images/empty_grades.svg`
  - `assets/images/empty_notes.svg`
  - `assets/images/empty_notification.svg`
  - `assets/images/error_illustration.svg`
- [ ] Configure assets in pubspec.yaml

**Requirements:** REQ-12 (UI/UX - empty states)

---

## Task 27: Utilities & Helpers

- [ ] Create `lib/core/utils/date_formatter.dart` (Indonesian date format)
- [ ] Create `lib/core/utils/validators.dart` (form validation)
- [ ] Create `lib/core/utils/extensions.dart` (String, DateTime extensions)
- [ ] Create error message handler with Indonesian translations

**Requirements:** REQ-13 (Error Handling)

---

## Task 28: Unit Tests - Models

- [ ] Write tests for UserModel serialization
- [ ] Write tests for AttendanceModel serialization
- [ ] Write tests for GradeModel serialization
- [ ] Write tests for all BK models serialization
- [ ] Write tests for NotificationModel serialization
- [ ] Verify round-trip property (Property 14)

**Requirements:** REQ-14.5 (Round-trip property)

---

## Task 29: Unit Tests - BLoCs

- [ ] Write tests for AuthBloc state transitions
- [ ] Write tests for ChildSelectorBloc (auto-select single child)
- [ ] Write tests for DashboardBloc
- [ ] Write tests for AttendanceBloc (pagination, refresh)
- [ ] Write tests for NotificationBloc (unread count, mark read)

**Requirements:** All feature requirements

---

## Task 30: Widget Tests

- [ ] Write tests for LoginPage (form validation, error display)
- [ ] Write tests for DashboardPage (data display, child switch)
- [ ] Write tests for AttendancePage (list rendering, color coding)
- [ ] Write tests for NotificationPage (read/unread states)
- [ ] Write tests for common widgets (buttons, cards, empty states)

**Requirements:** REQ-12 (UI/UX)

---

## Task 31: Integration Tests

- [ ] Write integration test for login flow
- [ ] Write integration test for parent dashboard flow
- [ ] Write integration test for child switching
- [ ] Write integration test for student flow
- [ ] Write integration test for notification handling

**Requirements:** All feature requirements

---

## Task 32: Final Polish & Documentation

- [ ] Review all screens for design consistency
- [ ] Verify all Indonesian strings are correct
- [ ] Test on multiple screen sizes
- [ ] Create README.md with setup instructions
- [ ] Document API endpoint requirements
- [ ] Final QA testing

**Requirements:** All requirements
