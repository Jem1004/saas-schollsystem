# Design Document

## Overview

Aplikasi mobile Flutter untuk Sistem SaaS Manajemen Sekolah dengan target pengguna Orang Tua dan Siswa. Aplikasi ini menggunakan arsitektur Clean Architecture dengan state management menggunakan flutter_bloc, terintegrasi dengan backend Golang via REST API, dan mendukung push notification via Firebase Cloud Messaging.

Desain UI menggunakan konsep modern minimalis dengan warna primary orange (#FF6B00) dan putih (#FFFFFF), typography clean, dan komponen dengan rounded corners untuk kesan friendly dan approachable.

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           FLUTTER MOBILE APP                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                      PRESENTATION LAYER                              │    │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐               │    │
│  │  │    Pages     │  │    Widgets   │  │    BLoCs     │               │    │
│  │  │  (Screens)   │  │ (Components) │  │   (State)    │               │    │
│  │  └──────────────┘  └──────────────┘  └──────────────┘               │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                    │                                         │
│                                    ▼                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                        DOMAIN LAYER                                  │    │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐               │    │
│  │  │   Entities   │  │   UseCases   │  │ Repositories │               │    │
│  │  │   (Models)   │  │  (Business)  │  │ (Interfaces) │               │    │
│  │  └──────────────┘  └──────────────┘  └──────────────┘               │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                    │                                         │
│                                    ▼                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                         DATA LAYER                                   │    │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐               │    │
│  │  │   Remote     │  │    Local     │  │  Repository  │               │    │
│  │  │ DataSource   │  │  DataSource  │  │    Impl      │               │    │
│  │  └──────────────┘  └──────────────┘  └──────────────┘               │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                    │                                         │
└────────────────────────────────────┼────────────────────────────────────────┘
                                     │
                    ┌────────────────┼────────────────┐
                    │                │                │
                    ▼                ▼                ▼
             ┌───────────┐    ┌───────────┐    ┌───────────┐
             │  Backend  │    │  Secure   │    │   FCM     │
             │  REST API │    │  Storage  │    │  (Push)   │
             └───────────┘    └───────────┘    └───────────┘
```

### Project Structure

```
mobile/
├── lib/
│   ├── main.dart
│   ├── app.dart
│   │
│   ├── core/
│   │   ├── config/
│   │   │   ├── app_config.dart
│   │   │   ├── api_config.dart
│   │   │   └── routes.dart
│   │   ├── constants/
│   │   │   ├── api_endpoints.dart
│   │   │   ├── app_colors.dart
│   │   │   ├── app_typography.dart
│   │   │   ├── app_spacing.dart
│   │   │   └── app_strings.dart
│   │   ├── errors/
│   │   │   ├── exceptions.dart
│   │   │   └── failures.dart
│   │   ├── network/
│   │   │   ├── api_client.dart
│   │   │   ├── api_interceptor.dart
│   │   │   └── network_info.dart
│   │   ├── theme/
│   │   │   ├── app_theme.dart
│   │   │   └── app_decorations.dart
│   │   └── utils/
│   │       ├── date_formatter.dart
│   │       ├── validators.dart
│   │       └── extensions.dart
│   │
│   ├── data/
│   │   ├── models/
│   │   │   ├── user_model.dart
│   │   │   ├── student_model.dart
│   │   │   ├── attendance_model.dart
│   │   │   ├── grade_model.dart
│   │   │   ├── violation_model.dart
│   │   │   ├── achievement_model.dart
│   │   │   ├── permit_model.dart
│   │   │   ├── homeroom_note_model.dart
│   │   │   ├── counseling_note_model.dart
│   │   │   └── notification_model.dart
│   │   ├── datasources/
│   │   │   ├── remote/
│   │   │   │   ├── auth_remote_datasource.dart
│   │   │   │   ├── parent_remote_datasource.dart
│   │   │   │   ├── student_remote_datasource.dart
│   │   │   │   └── notification_remote_datasource.dart
│   │   │   └── local/
│   │   │       ├── auth_local_datasource.dart
│   │   │       └── cache_local_datasource.dart
│   │   └── repositories/
│   │       ├── auth_repository_impl.dart
│   │       ├── parent_repository_impl.dart
│   │       ├── student_repository_impl.dart
│   │       └── notification_repository_impl.dart
│   │
│   ├── domain/
│   │   ├── entities/
│   │   │   ├── user.dart
│   │   │   ├── child.dart
│   │   │   ├── attendance.dart
│   │   │   ├── grade.dart
│   │   │   ├── bk_info.dart
│   │   │   └── notification.dart
│   │   ├── repositories/
│   │   │   ├── auth_repository.dart
│   │   │   ├── parent_repository.dart
│   │   │   ├── student_repository.dart
│   │   │   └── notification_repository.dart
│   │   └── usecases/
│   │       ├── auth/
│   │       │   ├── login_usecase.dart
│   │       │   ├── logout_usecase.dart
│   │       │   └── change_password_usecase.dart
│   │       ├── parent/
│   │       │   ├── get_children_usecase.dart
│   │       │   ├── get_child_dashboard_usecase.dart
│   │       │   ├── get_child_attendance_usecase.dart
│   │       │   ├── get_child_grades_usecase.dart
│   │       │   ├── get_child_notes_usecase.dart
│   │       │   └── get_child_bk_usecase.dart
│   │       ├── student/
│   │       │   ├── get_profile_usecase.dart
│   │       │   ├── get_dashboard_usecase.dart
│   │       │   ├── get_attendance_usecase.dart
│   │       │   ├── get_grades_usecase.dart
│   │       │   └── get_bk_usecase.dart
│   │       └── notification/
│   │           ├── get_notifications_usecase.dart
│   │           └── mark_as_read_usecase.dart
│   │
│   ├── presentation/
│   │   ├── blocs/
│   │   │   ├── auth/
│   │   │   │   ├── auth_bloc.dart
│   │   │   │   ├── auth_event.dart
│   │   │   │   └── auth_state.dart
│   │   │   ├── parent/
│   │   │   │   ├── dashboard_bloc.dart
│   │   │   │   ├── attendance_bloc.dart
│   │   │   │   ├── grades_bloc.dart
│   │   │   │   ├── notes_bloc.dart
│   │   │   │   └── bk_bloc.dart
│   │   │   ├── student/
│   │   │   │   └── student_bloc.dart
│   │   │   ├── notification/
│   │   │   │   └── notification_bloc.dart
│   │   │   └── child_selector/
│   │   │       └── child_selector_bloc.dart
│   │   ├── pages/
│   │   │   ├── splash/
│   │   │   │   └── splash_page.dart
│   │   │   ├── auth/
│   │   │   │   ├── login_page.dart
│   │   │   │   └── change_password_page.dart
│   │   │   ├── parent/
│   │   │   │   ├── parent_main_page.dart
│   │   │   │   ├── child_selector_page.dart
│   │   │   │   ├── dashboard_page.dart
│   │   │   │   ├── attendance_page.dart
│   │   │   │   ├── grades_page.dart
│   │   │   │   ├── notes_page.dart
│   │   │   │   ├── bk_page.dart
│   │   │   │   └── permit_detail_page.dart
│   │   │   ├── student/
│   │   │   │   ├── student_main_page.dart
│   │   │   │   ├── student_dashboard_page.dart
│   │   │   │   ├── student_attendance_page.dart
│   │   │   │   ├── student_grades_page.dart
│   │   │   │   └── student_bk_page.dart
│   │   │   ├── notification/
│   │   │   │   └── notification_page.dart
│   │   │   └── settings/
│   │   │       └── settings_page.dart
│   │   └── widgets/
│   │       ├── common/
│   │       │   ├── app_scaffold.dart
│   │       │   ├── app_app_bar.dart
│   │       │   ├── app_bottom_nav.dart
│   │       │   ├── app_button.dart
│   │       │   ├── app_text_field.dart
│   │       │   ├── app_card.dart
│   │       │   ├── loading_widget.dart
│   │       │   ├── skeleton_widget.dart
│   │       │   ├── error_widget.dart
│   │       │   ├── empty_state_widget.dart
│   │       │   └── refresh_indicator.dart
│   │       ├── cards/
│   │       │   ├── summary_card.dart
│   │       │   ├── attendance_card.dart
│   │       │   ├── grade_card.dart
│   │       │   ├── note_card.dart
│   │       │   ├── achievement_card.dart
│   │       │   ├── violation_card.dart
│   │       │   ├── permit_card.dart
│   │       │   ├── activity_card.dart
│   │       │   └── notification_card.dart
│   │       ├── headers/
│   │       │   ├── child_header.dart
│   │       │   ├── profile_header.dart
│   │       │   └── attendance_status_header.dart
│   │       └── dialogs/
│   │           ├── child_selector_dialog.dart
│   │           ├── confirm_dialog.dart
│   │           └── error_dialog.dart
│   │
│   └── services/
│       ├── fcm_service.dart
│       ├── local_notification_service.dart
│       ├── deep_link_service.dart
│       └── injection_container.dart
│
├── assets/
│   ├── images/
│   │   ├── logo.png
│   │   ├── empty_attendance.svg
│   │   ├── empty_grades.svg
│   │   ├── empty_notes.svg
│   │   ├── empty_notification.svg
│   │   └── error_illustration.svg
│   ├── icons/
│   │   └── (custom icons if needed)
│   └── fonts/
│       └── (Poppins or Inter if used)
│
├── test/
│   ├── unit/
│   ├── widget/
│   └── integration/
│
├── pubspec.yaml
└── README.md
```

## Components and Interfaces

### 1. API Client

```dart
// Dio-based API client with interceptors
abstract class ApiClient {
  Future<Response<T>> get<T>(String path, {Map<String, dynamic>? queryParams});
  Future<Response<T>> post<T>(String path, {dynamic data});
  Future<Response<T>> put<T>(String path, {dynamic data});
  Future<Response<T>> delete<T>(String path);
}

// Auth interceptor for JWT token injection
class AuthInterceptor extends Interceptor {
  final AuthLocalDataSource authLocalDataSource;
  
  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler);
  
  @override
  void onError(DioException err, ErrorInterceptorHandler handler);
}
```

### 2. Repository Interfaces

```dart
// Auth Repository
abstract class AuthRepository {
  Future<Either<Failure, LoginResponse>> login(String username, String password);
  Future<Either<Failure, TokenPair>> refreshToken(String refreshToken);
  Future<Either<Failure, void>> logout();
  Future<Either<Failure, void>> changePassword(String oldPassword, String newPassword);
  Future<Either<Failure, User>> getCurrentUser();
  Future<bool> isLoggedIn();
}

// Parent Repository
abstract class ParentRepository {
  Future<Either<Failure, List<Child>>> getChildren();
  Future<Either<Failure, ChildDashboard>> getChildDashboard(int childId);
  Future<Either<Failure, PaginatedList<Attendance>>> getChildAttendance(int childId, AttendanceFilter filter);
  Future<Either<Failure, AttendanceSummary>> getChildAttendanceSummary(int childId, DateRange range);
  Future<Either<Failure, PaginatedList<Grade>>> getChildGrades(int childId, GradeFilter filter);
  Future<Either<Failure, GradeSummary>> getChildGradeSummary(int childId);
  Future<Either<Failure, PaginatedList<HomeroomNote>>> getChildNotes(int childId, NoteFilter filter);
  Future<Either<Failure, BKInfo>> getChildBKInfo(int childId);
}

// Student Repository
abstract class StudentRepository {
  Future<Either<Failure, StudentProfile>> getProfile();
  Future<Either<Failure, StudentDashboard>> getDashboard();
  Future<Either<Failure, PaginatedList<Attendance>>> getAttendance(AttendanceFilter filter);
  Future<Either<Failure, AttendanceSummary>> getAttendanceSummary(DateRange range);
  Future<Either<Failure, PaginatedList<Grade>>> getGrades(GradeFilter filter);
  Future<Either<Failure, GradeSummary>> getGradeSummary();
  Future<Either<Failure, StudentBKInfo>> getBKInfo();
}

// Notification Repository
abstract class NotificationRepository {
  Future<Either<Failure, PaginatedList<AppNotification>>> getNotifications(NotificationFilter filter);
  Future<Either<Failure, NotificationSummary>> getSummary();
  Future<Either<Failure, int>> getUnreadCount();
  Future<Either<Failure, void>> markAsRead(int notificationId);
  Future<Either<Failure, void>> markAllAsRead();
  Future<Either<Failure, void>> registerFCMToken(String token, String platform);
  Future<Either<Failure, void>> deactivateFCMToken(String token);
}
```

### 3. BLoC Interfaces

```dart
// Auth BLoC
class AuthBloc extends Bloc<AuthEvent, AuthState> {
  // Events: LoginRequested, LogoutRequested, CheckAuthStatus, ChangePasswordRequested
  // States: AuthInitial, AuthLoading, Authenticated, Unauthenticated, AuthError, MustChangePassword
}

// Dashboard BLoC (Parent)
class DashboardBloc extends Bloc<DashboardEvent, DashboardState> {
  // Events: LoadDashboard, RefreshDashboard, SelectChild
  // States: DashboardInitial, DashboardLoading, DashboardLoaded, DashboardError
}

// Child Selector BLoC
class ChildSelectorBloc extends Bloc<ChildSelectorEvent, ChildSelectorState> {
  // Events: LoadChildren, SelectChild
  // States: ChildSelectorInitial, ChildSelectorLoading, ChildSelectorLoaded, ChildSelectorError
}
```

## Data Models

### Core Models

```dart
// User Model
class UserModel {
  final int id;
  final int? schoolId;
  final String role; // 'parent' or 'student'
  final String username;
  final String? email;
  final String name;
  final bool isActive;
  final bool mustResetPwd;
  final DateTime? lastLoginAt;
  
  factory UserModel.fromJson(Map<String, dynamic> json);
  Map<String, dynamic> toJson();
}

// Login Response
class LoginResponse {
  final String accessToken;
  final String refreshToken;
  final UserModel user;
  final bool mustResetPassword;
  
  factory LoginResponse.fromJson(Map<String, dynamic> json);
}

// Child Model (for Parent)
class ChildModel {
  final int id;
  final String nis;
  final String nisn;
  final String name;
  final String? className;
  final String? schoolName;
  final bool isActive;
  
  factory ChildModel.fromJson(Map<String, dynamic> json);
}
```

### Attendance Models

```dart
// Attendance Model
class AttendanceModel {
  final int id;
  final int studentId;
  final DateTime date;
  final DateTime? checkInTime;
  final DateTime? checkOutTime;
  final String status; // on_time, late, very_late, absent, sick, excused
  final String method; // rfid, manual
  
  factory AttendanceModel.fromJson(Map<String, dynamic> json);
  
  // Helper getters
  bool get isPresent => checkInTime != null;
  bool get isOnTime => status == 'on_time';
  bool get isLate => status == 'late' || status == 'very_late';
  Color get statusColor; // Returns appropriate color based on status
}

// Attendance Summary
class AttendanceSummary {
  final int totalDays;
  final int onTime;
  final int late;
  final int veryLate;
  final int absent;
  final int sick;
  final int excused;
  
  double get attendancePercentage => (onTime + late + veryLate) / totalDays * 100;
  
  factory AttendanceSummary.fromJson(Map<String, dynamic> json);
}
```

### Academic Models

```dart
// Grade Model
class GradeModel {
  final int id;
  final int studentId;
  final String title;
  final double score;
  final String? description;
  final String? teacherName;
  final DateTime createdAt;
  
  factory GradeModel.fromJson(Map<String, dynamic> json);
  
  Color get scoreColor {
    if (score >= 80) return AppColors.success;
    if (score >= 60) return AppColors.warning;
    return AppColors.error;
  }
}

// Grade Summary
class GradeSummary {
  final int totalGrades;
  final double averageScore;
  final double highestScore;
  final double lowestScore;
  
  factory GradeSummary.fromJson(Map<String, dynamic> json);
}

// Homeroom Note Model
class HomeroomNoteModel {
  final int id;
  final int studentId;
  final String content;
  final String teacherName;
  final DateTime createdAt;
  
  factory HomeroomNoteModel.fromJson(Map<String, dynamic> json);
}
```

### BK Models

```dart
// Achievement Model
class AchievementModel {
  final int id;
  final int studentId;
  final String title;
  final int point;
  final String? description;
  final DateTime createdAt;
  
  factory AchievementModel.fromJson(Map<String, dynamic> json);
}

// Violation Model
class ViolationModel {
  final int id;
  final int studentId;
  final String category;
  final String level; // ringan, sedang, berat
  final int point; // negative value
  final String description;
  final DateTime createdAt;
  
  factory ViolationModel.fromJson(Map<String, dynamic> json);
  
  Color get levelColor {
    switch (level) {
      case 'ringan': return AppColors.warning;
      case 'sedang': return AppColors.orange;
      case 'berat': return AppColors.error;
      default: return AppColors.textSecondary;
    }
  }
}

// Permit Model
class PermitModel {
  final int id;
  final int studentId;
  final String reason;
  final DateTime exitTime;
  final DateTime? returnTime;
  final String teacherName;
  final String? documentUrl;
  final DateTime createdAt;
  
  factory PermitModel.fromJson(Map<String, dynamic> json);
  
  bool get hasReturned => returnTime != null;
}

// Counseling Note Model (Parent view - only summary)
class CounselingNoteModel {
  final int id;
  final int studentId;
  final String parentSummary; // internal_note is hidden
  final DateTime createdAt;
  
  factory CounselingNoteModel.fromJson(Map<String, dynamic> json);
}

// BK Info (Combined)
class BKInfo {
  final int totalAchievementPoints;
  final int totalViolationPoints;
  final List<AchievementModel> achievements;
  final List<ViolationModel> violations;
  final List<PermitModel> permits;
  final List<CounselingNoteModel> counselingNotes;
  
  int get netPoints => totalAchievementPoints + totalViolationPoints;
  
  factory BKInfo.fromJson(Map<String, dynamic> json);
}
```

### Notification Model

```dart
// Notification Model
class NotificationModel {
  final int id;
  final int userId;
  final String type; // attendance_in, attendance_out, violation, achievement, permit, grade, homeroom_note, counseling
  final String title;
  final String message;
  final Map<String, dynamic>? data;
  final bool isRead;
  final DateTime createdAt;
  
  factory NotificationModel.fromJson(Map<String, dynamic> json);
  
  IconData get typeIcon {
    switch (type) {
      case 'attendance_in':
      case 'attendance_out':
        return Icons.access_time;
      case 'violation':
        return Icons.warning;
      case 'achievement':
        return Icons.emoji_events;
      case 'permit':
        return Icons.exit_to_app;
      case 'grade':
        return Icons.school;
      case 'homeroom_note':
        return Icons.note;
      case 'counseling':
        return Icons.psychology;
      default:
        return Icons.notifications;
    }
  }
}

// Notification Summary
class NotificationSummary {
  final int total;
  final int unread;
  
  factory NotificationSummary.fromJson(Map<String, dynamic> json);
}
```


## UI/UX Design System

### Color Palette

```dart
class AppColors {
  // Primary Colors
  static const Color primary = Color(0xFFFF6B00);        // Orange
  static const Color primaryLight = Color(0xFFFF8A3D);   // Light Orange
  static const Color primaryDark = Color(0xFFE55A00);    // Dark Orange
  static const Color primarySurface = Color(0xFFFFF3E8); // Very Light Orange (backgrounds)
  
  // Neutral Colors
  static const Color white = Color(0xFFFFFFFF);
  static const Color background = Color(0xFFF8F9FA);     // Light Gray Background
  static const Color surface = Color(0xFFFFFFFF);        // Card Background
  static const Color border = Color(0xFFE8E8E8);         // Border Color
  static const Color divider = Color(0xFFF0F0F0);        // Divider Color
  
  // Text Colors
  static const Color textPrimary = Color(0xFF1A1A1A);    // Main Text
  static const Color textSecondary = Color(0xFF6B7280);  // Secondary Text
  static const Color textTertiary = Color(0xFF9CA3AF);   // Hint Text
  static const Color textOnPrimary = Color(0xFFFFFFFF);  // Text on Orange
  
  // Status Colors
  static const Color success = Color(0xFF22C55E);        // Green
  static const Color successLight = Color(0xFFDCFCE7);   // Light Green
  static const Color warning = Color(0xFFFACC15);        // Yellow
  static const Color warningLight = Color(0xFFFEF9C3);   // Light Yellow
  static const Color error = Color(0xFFEF4444);          // Red
  static const Color errorLight = Color(0xFFFEE2E2);     // Light Red
  static const Color info = Color(0xFF3B82F6);           // Blue
  static const Color infoLight = Color(0xFFDBEAFE);      // Light Blue
  
  // Attendance Status Colors
  static const Color onTime = Color(0xFF22C55E);         // Green
  static const Color late = Color(0xFFFACC15);           // Yellow
  static const Color veryLate = Color(0xFFF97316);       // Orange
  static const Color absent = Color(0xFFEF4444);         // Red
  
  // Skeleton Loading
  static const Color skeleton = Color(0xFFE5E7EB);
  static const Color skeletonHighlight = Color(0xFFF3F4F6);
}
```

### Typography

```dart
class AppTypography {
  // Font Family
  static const String fontFamily = 'Poppins'; // or system default
  
  // Headings
  static const TextStyle h1 = TextStyle(
    fontSize: 28,
    fontWeight: FontWeight.w700,
    color: AppColors.textPrimary,
    height: 1.3,
  );
  
  static const TextStyle h2 = TextStyle(
    fontSize: 24,
    fontWeight: FontWeight.w600,
    color: AppColors.textPrimary,
    height: 1.3,
  );
  
  static const TextStyle h3 = TextStyle(
    fontSize: 20,
    fontWeight: FontWeight.w600,
    color: AppColors.textPrimary,
    height: 1.4,
  );
  
  static const TextStyle h4 = TextStyle(
    fontSize: 18,
    fontWeight: FontWeight.w600,
    color: AppColors.textPrimary,
    height: 1.4,
  );
  
  // Body Text
  static const TextStyle bodyLarge = TextStyle(
    fontSize: 16,
    fontWeight: FontWeight.w400,
    color: AppColors.textPrimary,
    height: 1.5,
  );
  
  static const TextStyle bodyMedium = TextStyle(
    fontSize: 14,
    fontWeight: FontWeight.w400,
    color: AppColors.textPrimary,
    height: 1.5,
  );
  
  static const TextStyle bodySmall = TextStyle(
    fontSize: 12,
    fontWeight: FontWeight.w400,
    color: AppColors.textSecondary,
    height: 1.5,
  );
  
  // Labels
  static const TextStyle labelLarge = TextStyle(
    fontSize: 14,
    fontWeight: FontWeight.w500,
    color: AppColors.textPrimary,
    height: 1.4,
  );
  
  static const TextStyle labelMedium = TextStyle(
    fontSize: 12,
    fontWeight: FontWeight.w500,
    color: AppColors.textSecondary,
    height: 1.4,
  );
  
  static const TextStyle labelSmall = TextStyle(
    fontSize: 10,
    fontWeight: FontWeight.w500,
    color: AppColors.textTertiary,
    height: 1.4,
  );
  
  // Button Text
  static const TextStyle button = TextStyle(
    fontSize: 16,
    fontWeight: FontWeight.w600,
    color: AppColors.white,
    height: 1.2,
  );
}
```

### Spacing System (8px Grid)

```dart
class AppSpacing {
  static const double xs = 4.0;
  static const double sm = 8.0;
  static const double md = 16.0;
  static const double lg = 24.0;
  static const double xl = 32.0;
  static const double xxl = 48.0;
  
  // Padding
  static const EdgeInsets pagePadding = EdgeInsets.all(16.0);
  static const EdgeInsets cardPadding = EdgeInsets.all(16.0);
  static const EdgeInsets listItemPadding = EdgeInsets.symmetric(horizontal: 16.0, vertical: 12.0);
  
  // Border Radius
  static const double radiusSm = 8.0;
  static const double radiusMd = 12.0;
  static const double radiusLg = 16.0;
  static const double radiusXl = 24.0;
  static const double radiusFull = 999.0;
}
```

### Component Styles

```dart
class AppDecorations {
  // Card Decoration
  static BoxDecoration card = BoxDecoration(
    color: AppColors.surface,
    borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
    boxShadow: [
      BoxShadow(
        color: Colors.black.withOpacity(0.04),
        blurRadius: 8,
        offset: const Offset(0, 2),
      ),
    ],
  );
  
  // Elevated Card
  static BoxDecoration cardElevated = BoxDecoration(
    color: AppColors.surface,
    borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
    boxShadow: [
      BoxShadow(
        color: Colors.black.withOpacity(0.08),
        blurRadius: 16,
        offset: const Offset(0, 4),
      ),
    ],
  );
  
  // Primary Button
  static BoxDecoration primaryButton = BoxDecoration(
    color: AppColors.primary,
    borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
  );
  
  // Outlined Button
  static BoxDecoration outlinedButton = BoxDecoration(
    color: AppColors.white,
    borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
    border: Border.all(color: AppColors.primary, width: 1.5),
  );
  
  // Text Field
  static InputDecoration textField({String? hint, Widget? prefix, Widget? suffix}) {
    return InputDecoration(
      hintText: hint,
      hintStyle: AppTypography.bodyMedium.copyWith(color: AppColors.textTertiary),
      prefixIcon: prefix,
      suffixIcon: suffix,
      filled: true,
      fillColor: AppColors.background,
      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        borderSide: BorderSide.none,
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        borderSide: const BorderSide(color: AppColors.primary, width: 1.5),
      ),
      errorBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        borderSide: const BorderSide(color: AppColors.error, width: 1.5),
      ),
    );
  }
  
  // Status Badge
  static BoxDecoration statusBadge(Color color) {
    return BoxDecoration(
      color: color.withOpacity(0.15),
      borderRadius: BorderRadius.circular(AppSpacing.radiusSm),
    );
  }
}
```

## Screen Wireframes

### 1. Splash Screen

```
┌─────────────────────────────────────┐
│                                     │
│                                     │
│                                     │
│                                     │
│           ┌─────────┐               │
│           │  LOGO   │               │
│           │ (Orange)│               │
│           └─────────┘               │
│                                     │
│         School Management           │
│              App                    │
│                                     │
│                                     │
│            ◐ Loading...             │
│                                     │
│                                     │
│                                     │
└─────────────────────────────────────┘
```

### 2. Login Page

```
┌─────────────────────────────────────┐
│                                     │
│           ┌─────────┐               │
│           │  LOGO   │               │
│           └─────────┘               │
│                                     │
│         Selamat Datang              │
│    Masuk untuk melanjutkan          │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ 👤  NISN                    │    │
│  └─────────────────────────────┘    │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ 🔒  Password            👁  │    │
│  └─────────────────────────────┘    │
│                                     │
│  ┌─────────────────────────────┐    │
│  │          MASUK              │    │
│  │    (Orange Button)          │    │
│  └─────────────────────────────┘    │
│                                     │
│       Lupa password? Hubungi        │
│       admin sekolah                 │
│                                     │
└─────────────────────────────────────┘
```

### 3. Child Selector (Parent with Multiple Children)

```
┌─────────────────────────────────────┐
│                                     │
│         Pilih Anak                  │
│                                     │
│  ┌─────────────────────────────┐    │
│  │  👤  Ahmad Fauzi            │    │
│  │      Kelas 7A               │    │
│  │      SMP Negeri 1           │    │
│  │                         ▶   │    │
│  └─────────────────────────────┘    │
│                                     │
│  ┌─────────────────────────────┐    │
│  │  👤  Siti Aminah            │    │
│  │      Kelas 5B               │    │
│  │      SD Negeri 2            │    │
│  │                         ▶   │    │
│  └─────────────────────────────┘    │
│                                     │
│                                     │
│                                     │
│                                     │
│                                     │
│                                     │
└─────────────────────────────────────┘
```

### 4. Parent Dashboard

```
┌─────────────────────────────────────┐
│ ☰  Selamat Pagi           🔔(3)    │
├─────────────────────────────────────┤
│                                     │
│  ┌─────────────────────────────┐    │
│  │  👤 Ahmad Fauzi      [▼]    │    │
│  │     Kelas 7A                │    │
│  │     SMP Negeri 1            │    │
│  └─────────────────────────────┘    │
│                                     │
│  ── Status Hari Ini ──              │
│  ┌────────────┐ ┌────────────┐      │
│  │  ✅ Masuk  │ │  🕐 Pulang │      │
│  │   06:45    │ │    --:--   │      │
│  │ Tepat Waktu│ │   Belum    │      │
│  └────────────┘ └────────────┘      │
│                                     │
│  ── Ringkasan ──                    │
│  ┌────────┐┌────────┐┌────────┐     │
│  │Kehadiran││ Nilai  ││Prestasi│     │
│  │  95%   ││  85.5  ││  120   │     │
│  │ bulan  ││rata-rata││ poin  │     │
│  └────────┘└────────┘└────────┘     │
│                                     │
│  ── Aktivitas Terbaru ──            │
│  ┌─────────────────────────────┐    │
│  │ ✅ Masuk sekolah            │    │
│  │    Hari ini, 06:45          │    │
│  ├─────────────────────────────┤    │
│  │ 📚 Nilai Matematika: 90     │    │
│  │    Kemarin, 10:30           │    │
│  ├─────────────────────────────┤    │
│  │ 📝 Catatan wali kelas       │    │
│  │    2 hari lalu              │    │
│  └─────────────────────────────┘    │
│                                     │
├─────────────────────────────────────┤
│  🏠    📅    📊    📋    👤        │
│ Home  Absen  Nilai  BK   Profil    │
└─────────────────────────────────────┘
```

### 5. Attendance Page

```
┌─────────────────────────────────────┐
│  ←  Kehadiran           📅 Filter   │
├─────────────────────────────────────┤
│                                     │
│  ── Ringkasan Bulan Ini ──          │
│  ┌─────────────────────────────┐    │
│  │ ████████████████░░░░  95%   │    │
│  │                             │    │
│  │  ✅ 18   ⚠️ 2   ❌ 1        │    │
│  │ Hadir  Terlambat  Absen     │    │
│  └─────────────────────────────┘    │
│                                     │
│  ── Riwayat ──                      │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ Senin, 4 Januari 2026       │    │
│  │ ┌─────┐                     │    │
│  │ │ ✅  │ Masuk: 06:45        │    │
│  │ │     │ Pulang: 14:30       │    │
│  │ └─────┘ Tepat Waktu         │    │
│  └─────────────────────────────┘    │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ Jumat, 3 Januari 2026       │    │
│  │ ┌─────┐                     │    │
│  │ │ ⚠️  │ Masuk: 07:15        │    │
│  │ │     │ Pulang: 14:30       │    │
│  │ └─────┘ Terlambat           │    │
│  └─────────────────────────────┘    │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ Kamis, 2 Januari 2026       │    │
│  │ ┌─────┐                     │    │
│  │ │ ✅  │ Masuk: 06:40        │    │
│  │ │     │ Pulang: 14:25       │    │
│  │ └─────┘ Tepat Waktu         │    │
│  └─────────────────────────────┘    │
│                                     │
├─────────────────────────────────────┤
│  🏠    📅    📊    📋    👤        │
└─────────────────────────────────────┘
```

### 6. Grades Page

```
┌─────────────────────────────────────┐
│  ←  Nilai Akademik                  │
├─────────────────────────────────────┤
│                                     │
│  ── Ringkasan ──                    │
│  ┌─────────────────────────────┐    │
│  │      Rata-rata Nilai        │    │
│  │    ┌───────────────┐        │    │
│  │    │     85.5      │        │    │
│  │    └───────────────┘        │    │
│  │                             │    │
│  │  Tertinggi    Terendah      │    │
│  │     95           70         │    │
│  └─────────────────────────────┘    │
│                                     │
│  ── Nilai Terbaru ──                │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ 📚 Ulangan Matematika       │    │
│  │    ┌────┐                   │    │
│  │    │ 90 │  Bu Siti          │    │
│  │    └────┘  4 Jan 2026       │    │
│  │    "Bagus, pertahankan!"    │    │
│  └─────────────────────────────┘    │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ 📚 Tugas Bahasa Indonesia   │    │
│  │    ┌────┐                   │    │
│  │    │ 85 │  Pak Ahmad        │    │
│  │    └────┘  3 Jan 2026       │    │
│  └─────────────────────────────┘    │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ 📚 Quiz IPA                 │    │
│  │    ┌────┐                   │    │
│  │    │ 78 │  Bu Rina          │    │
│  │    └────┘  2 Jan 2026       │    │
│  └─────────────────────────────┘    │
│                                     │
├─────────────────────────────────────┤
│  🏠    📅    📊    📋    👤        │
└─────────────────────────────────────┘
```

### 7. BK Page (Tabs)

```
┌─────────────────────────────────────┐
│  ←  Informasi BK                    │
├─────────────────────────────────────┤
│                                     │
│  ┌────────────┐ ┌────────────┐      │
│  │  🏆 +120   │ │  ⚠️ -15    │      │
│  │  Prestasi  │ │ Pelanggaran│      │
│  └────────────┘ └────────────┘      │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ Prestasi │Pelang│ Izin │Kons│    │
│  │    ●    │  ○   │  ○   │ ○  │    │
│  └─────────────────────────────┘    │
│                                     │
│  ── Tab: Prestasi ──                │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ 🏆 Juara 1 Lomba Matematika │    │
│  │    +50 poin                 │    │
│  │    20 Des 2025              │    │
│  │    "Tingkat Kabupaten"      │    │
│  └─────────────────────────────┘    │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ 🏆 Siswa Teladan November   │    │
│  │    +30 poin                 │    │
│  │    30 Nov 2025              │    │
│  └─────────────────────────────┘    │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ 🏆 Juara 2 Lomba Cerdas     │    │
│  │    +40 poin                 │    │
│  │    15 Okt 2025              │    │
│  └─────────────────────────────┘    │
│                                     │
├─────────────────────────────────────┤
│  🏠    📅    📊    📋    👤        │
└─────────────────────────────────────┘
```

### 8. Notification Page

```
┌─────────────────────────────────────┐
│  ←  Notifikasi      [Tandai Baca]   │
├─────────────────────────────────────┤
│                                     │
│  ── Hari Ini ──                     │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ 🔵 ✅ Ahmad masuk sekolah   │    │
│  │       Pukul 06:45           │    │
│  │       5 menit lalu          │    │
│  └─────────────────────────────┘    │
│                                     │
│  ── Kemarin ──                      │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ ○ 📚 Nilai baru: Matematika │    │
│  │      Ahmad mendapat nilai 90│    │
│  │      Kemarin, 10:30         │    │
│  └─────────────────────────────┘    │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ ○ ✅ Ahmad pulang sekolah   │    │
│  │      Pukul 14:30            │    │
│  │      Kemarin, 14:30         │    │
│  └─────────────────────────────┘    │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ ○ 📝 Catatan wali kelas     │    │
│  │      Dari Bu Siti           │    │
│  │      Kemarin, 14:00         │    │
│  └─────────────────────────────┘    │
│                                     │
│  ── Minggu Lalu ──                  │
│  ...                                │
│                                     │
├─────────────────────────────────────┤
│  🏠    📅    📊    📋    👤        │
└─────────────────────────────────────┘
```

### 9. Profile/Settings Page

```
┌─────────────────────────────────────┐
│  ←  Profil                          │
├─────────────────────────────────────┤
│                                     │
│         ┌─────────────┐             │
│         │     👤      │             │
│         │   Avatar    │             │
│         └─────────────┘             │
│                                     │
│         Bapak/Ibu Ahmad             │
│         Orang Tua                   │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ 👤  Informasi Akun          │    │
│  │     Username: 0012345678    │    │
│  │                         ▶   │    │
│  └─────────────────────────────┘    │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ 🔒  Ubah Password           │    │
│  │                         ▶   │    │
│  └─────────────────────────────┘    │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ 🔔  Pengaturan Notifikasi   │    │
│  │                         ▶   │    │
│  └─────────────────────────────┘    │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ ℹ️  Tentang Aplikasi        │    │
│  │     Versi 1.0.0         ▶   │    │
│  └─────────────────────────────┘    │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ 🚪  Keluar                  │    │
│  │     (Red Text)          ▶   │    │
│  └─────────────────────────────┘    │
│                                     │
├─────────────────────────────────────┤
│  🏠    📅    📊    📋    👤        │
└─────────────────────────────────────┘
```


## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

Based on the prework analysis, the following correctness properties have been identified:

### Property 1: Token Lifecycle Management

*For any* user session, the following token behaviors SHALL hold:
- If valid tokens exist in secure storage, the app SHALL skip login and show main screen
- If access token expires, the app SHALL automatically attempt refresh using refresh token
- If logout is triggered, the app SHALL clear all stored tokens from secure storage

**Validates: Requirements 1.6, 1.7, 1.9**

### Property 2: Authentication Error Handling

*For any* authentication attempt with invalid credentials, the app SHALL display an appropriate error message in Indonesian language and NOT store any tokens.

**Validates: Requirements 1.5**

### Property 3: Child Selection Logic

*For any* parent user:
- If the parent has exactly one child, the app SHALL automatically select that child without showing selector
- If the parent has more than one child, the app SHALL display child selector for user to choose
- Selected child preference SHALL be persisted locally and restored on app restart

**Validates: Requirements 3.2, 3.3, 3.6**

### Property 4: Child Data Refresh on Switch

*For any* child switch action by a parent, the app SHALL reload ALL data (dashboard, attendance, grades, notes, BK) for the newly selected child. No stale data from the previous child SHALL be displayed.

**Validates: Requirements 3.5**

### Property 5: List Sorting Consistency

*For any* list display (attendance, grades, notes, notifications), the items SHALL be sorted by date in descending order (most recent first). This sorting SHALL be maintained after pagination loads.

**Validates: Requirements 5.2, 6.2, 7.1, 9.1**

### Property 6: Data Item Completeness

*For any* data item displayed in a list:
- Attendance items SHALL contain: date, check-in time, check-out time, status
- Grade items SHALL contain: title, score, teacher name, date
- Note items SHALL contain: content, teacher name, date
- Notification items SHALL contain: type icon, title, message, timestamp

**Validates: Requirements 5.3, 6.3, 7.2, 9.2**

### Property 7: Notification Deep Linking

*For any* notification tap action, the app SHALL:
1. Mark the notification as read
2. Navigate to the relevant screen based on notification type
3. The navigation target SHALL match the notification type (e.g., attendance_in → attendance page)

**Validates: Requirements 2.4, 9.3**

### Property 8: Notification Type Handling

*For any* notification type in the set {attendance_in, attendance_out, violation, achievement, permit, grade, homeroom_note, counseling}, the app SHALL have a defined handler and navigation target.

**Validates: Requirements 2.6**

### Property 9: Student Feature Restrictions

*For any* student user, the app SHALL NOT display:
- Homeroom notes (only visible to parents)
- Counseling notes (only visible to parents)
- Permit details (only visible to parents)

The student BK page SHALL only show achievements and violations.

**Validates: Requirements 11.3, 11.4, 11.5**

### Property 10: Parent Counseling Privacy

*For any* counseling note displayed to a parent, the app SHALL only show the `parent_summary` field. The `internal_note` field SHALL never be visible to parents.

**Validates: Requirements 8.5**

### Property 11: Role-Based Navigation

*For any* authenticated user:
- If role is 'parent', the app SHALL display parent dashboard with child selector capability
- If role is 'student', the app SHALL display student dashboard without child selector

**Validates: Requirements 10.1**

### Property 12: Error Handling with Retry

*For any* network request failure, the app SHALL:
1. Display an error message in Indonesian
2. Provide a retry button/action
3. If cached data exists, display cached data with "offline" indicator
4. Handle all API error codes with appropriate messages

**Validates: Requirements 13.1, 13.2, 13.3, 13.5**

### Property 13: Session Expiry Handling

*For any* session expiry (refresh token invalid), the app SHALL:
1. Display session expired message
2. Clear all stored tokens
3. Redirect to login page

**Validates: Requirements 13.4**

### Property 14: Data Serialization Round-Trip

*For any* valid data object (User, Child, Attendance, Grade, Notification, etc.), serializing to JSON then deserializing back SHALL produce an object equivalent to the original.

`deserialize(serialize(obj)) == obj`

**Validates: Requirements 14.5**

### Property 15: Secure Token Storage

*For any* JWT token (access or refresh), the app SHALL store it in encrypted secure storage, NOT in plain text local storage.

**Validates: Requirements 14.1**

### Property 16: Logout Data Cleanup

*For any* logout action, the app SHALL clear:
- All JWT tokens from secure storage
- Selected child ID from local storage
- All cached dashboard data
- FCM token registration on backend

**Validates: Requirements 14.4, 2.5**

### Property 17: Pull-to-Refresh Functionality

*For any* list screen (dashboard, attendance, grades, notes, BK, notifications), pull-to-refresh gesture SHALL trigger a data reload from the API.

**Validates: Requirements 12.9**

### Property 18: Empty State Handling

*For any* list that returns zero items, the app SHALL display an appropriate empty state illustration with helpful message, NOT a blank screen.

**Validates: Requirements 4.7**

### Property 19: Unread Notification Badge

*For any* state where unread notifications exist, the notification icon in bottom navigation SHALL display a badge with the unread count. When all notifications are read, the badge SHALL be hidden.

**Validates: Requirements 9.5**

### Property 20: BK Points Calculation

*For any* BK info display, the total achievement points SHALL equal the sum of all individual achievement point values, and total violation points SHALL equal the sum of all individual violation point values (negative).

**Validates: Requirements 8.7**

## Error Handling

### API Error Response Handling

```dart
// Error codes and Indonesian messages
class ApiErrorHandler {
  static String getErrorMessage(String code) {
    switch (code) {
      case 'AUTH_INVALID_CREDENTIALS':
        return 'Username atau password salah';
      case 'AUTH_ACCOUNT_INACTIVE':
        return 'Akun tidak aktif';
      case 'AUTH_SCHOOL_INACTIVE':
        return 'Sekolah tidak aktif';
      case 'AUTH_TOKEN_EXPIRED':
        return 'Sesi telah berakhir, silakan login kembali';
      case 'AUTH_TOKEN_INVALID':
        return 'Token tidak valid';
      case 'AUTH_PASSWORD_MISMATCH':
        return 'Password lama salah';
      case 'AUTH_SAME_PASSWORD':
        return 'Password baru harus berbeda dari password lama';
      case 'NOT_FOUND_STUDENT':
        return 'Siswa tidak ditemukan';
      case 'NOT_FOUND_PARENT':
        return 'Data orang tua tidak ditemukan';
      case 'AUTHZ_NOT_LINKED':
        return 'Anda tidak memiliki akses ke data siswa ini';
      case 'NETWORK_ERROR':
        return 'Tidak dapat terhubung ke server';
      case 'TIMEOUT':
        return 'Koneksi timeout, coba lagi';
      default:
        return 'Terjadi kesalahan, coba lagi nanti';
    }
  }
}
```

### Failure Types

```dart
// Domain layer failures
abstract class Failure {
  final String message;
  const Failure(this.message);
}

class ServerFailure extends Failure {
  final String code;
  const ServerFailure(this.code, String message) : super(message);
}

class NetworkFailure extends Failure {
  const NetworkFailure() : super('Tidak ada koneksi internet');
}

class CacheFailure extends Failure {
  const CacheFailure() : super('Gagal mengambil data tersimpan');
}

class AuthFailure extends Failure {
  const AuthFailure(String message) : super(message);
}

class ValidationFailure extends Failure {
  const ValidationFailure(String message) : super(message);
}
```

## Testing Strategy

### Dual Testing Approach

This application requires both unit tests and widget tests for comprehensive coverage:

1. **Unit Tests**: Verify business logic, data transformations, and repository implementations
2. **Widget Tests**: Verify UI components render correctly and respond to user interactions
3. **Integration Tests**: Verify complete user flows work end-to-end

### Property-Based Testing Configuration

- **Library**: `fast_check` for Dart (or custom generators)
- **Minimum Iterations**: 100 per property test
- **Tag Format**: `Feature: flutter-mobile-app, Property {number}: {property_text}`

### Test Categories

#### Unit Tests
- Model serialization/deserialization
- BLoC state transitions
- UseCase business logic
- Repository implementations (with mocked data sources)
- Error handling and mapping

#### Widget Tests
- Component rendering with various states
- User interaction handling
- Navigation flows
- Form validation
- Loading and error states

#### Integration Tests
- Login flow (valid/invalid credentials)
- Child selection flow
- Data refresh flows
- Notification handling
- Offline mode behavior

### Test File Structure

```
test/
├── unit/
│   ├── data/
│   │   ├── models/
│   │   │   ├── user_model_test.dart
│   │   │   ├── attendance_model_test.dart
│   │   │   └── ...
│   │   └── repositories/
│   │       ├── auth_repository_test.dart
│   │       └── ...
│   ├── domain/
│   │   └── usecases/
│   │       ├── login_usecase_test.dart
│   │       └── ...
│   └── presentation/
│       └── blocs/
│           ├── auth_bloc_test.dart
│           └── ...
├── widget/
│   ├── pages/
│   │   ├── login_page_test.dart
│   │   ├── dashboard_page_test.dart
│   │   └── ...
│   └── widgets/
│       ├── attendance_card_test.dart
│       └── ...
├── integration/
│   ├── login_flow_test.dart
│   ├── parent_flow_test.dart
│   └── student_flow_test.dart
└── property/
    ├── serialization_property_test.dart
    ├── sorting_property_test.dart
    └── token_lifecycle_property_test.dart
```

### Property Test Example

```dart
// Feature: flutter-mobile-app, Property 14: Data Serialization Round-Trip
void main() {
  group('Data Serialization Round-Trip', () {
    test('UserModel serialization round-trip', () {
      // Generate random valid UserModel instances
      for (var i = 0; i < 100; i++) {
        final original = generateRandomUserModel();
        final json = original.toJson();
        final deserialized = UserModel.fromJson(json);
        
        expect(deserialized.id, equals(original.id));
        expect(deserialized.username, equals(original.username));
        expect(deserialized.role, equals(original.role));
        expect(deserialized.name, equals(original.name));
      }
    });
    
    test('AttendanceModel serialization round-trip', () {
      for (var i = 0; i < 100; i++) {
        final original = generateRandomAttendanceModel();
        final json = original.toJson();
        final deserialized = AttendanceModel.fromJson(json);
        
        expect(deserialized.id, equals(original.id));
        expect(deserialized.date, equals(original.date));
        expect(deserialized.status, equals(original.status));
      }
    });
  });
}
```

### Coverage Requirements

- Unit test coverage: minimum 80%
- All correctness properties must have corresponding tests
- All BLoC state transitions must be tested
- Critical paths (auth, data display) must have 100% coverage
