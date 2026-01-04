/// Centralized API endpoint paths for the School Management mobile app.
///
/// All endpoint paths are relative to the base URL (e.g., '/api/v1').
/// Use with [ApiConfig.getUrl] to get the full URL.
///
/// Example usage:
/// ```dart
/// import 'package:mobile/core/constants/api_endpoints.dart';
/// import 'package:mobile/core/config/api_config.dart';
///
/// final loginUrl = ApiConfig.getUrl(ApiEndpoints.authLogin);
/// final dashboardUrl = ApiConfig.getUrlWithParams(
///   ApiEndpoints.parentChildDashboard,
///   {'id': '123'},
/// );
/// ```
library;

/// API endpoint paths organized by feature module.
class ApiEndpoints {
  ApiEndpoints._();

  // ============================================
  // Authentication Endpoints
  // ============================================

  /// POST - Login with NISN and password
  /// Request: { username: string, password: string }
  /// Response: { access_token: string, refresh_token: string, user: User, must_reset_password: bool }
  static const String authLogin = 'auth/login';

  /// POST - Refresh access token
  /// Request: { refresh_token: string }
  /// Response: { access_token: string, refresh_token: string }
  static const String authRefresh = 'auth/refresh';

  /// POST - Logout and invalidate tokens
  /// Request: (empty)
  /// Response: { success: bool }
  static const String authLogout = 'auth/logout';

  /// POST - Change password
  /// Request: { old_password: string, new_password: string }
  /// Response: { success: bool }
  static const String authChangePassword = 'auth/change-password';

  /// GET - Get current user profile
  /// Response: { user: User }
  static const String authMe = 'auth/me';

  // ============================================
  // Parent Endpoints
  // ============================================

  /// GET - Get list of linked children
  /// Response: { children: Child[] }
  static const String parentChildren = 'parent/children';

  /// GET - Get child dashboard data
  /// Path param: {id} - Child/Student ID
  /// Response: { dashboard: ChildDashboard }
  static const String parentChildDashboard = 'parent/children/{id}/dashboard';

  /// GET - Get child attendance list
  /// Path param: {id} - Child/Student ID
  /// Query params: page, limit, start_date, end_date
  /// Response: { data: Attendance[], pagination: Pagination }
  static const String parentChildAttendance = 'parent/children/{id}/attendance';

  /// GET - Get child attendance summary
  /// Path param: {id} - Child/Student ID
  /// Query params: start_date, end_date
  /// Response: { summary: AttendanceSummary }
  static const String parentChildAttendanceSummary = 'parent/children/{id}/attendance/summary';

  /// GET - Get child grades list
  /// Path param: {id} - Child/Student ID
  /// Query params: page, limit
  /// Response: { data: Grade[], pagination: Pagination }
  static const String parentChildGrades = 'parent/children/{id}/grades';

  /// GET - Get child grades summary
  /// Path param: {id} - Child/Student ID
  /// Response: { summary: GradeSummary }
  static const String parentChildGradesSummary = 'parent/children/{id}/grades/summary';

  /// GET - Get child homeroom notes
  /// Path param: {id} - Child/Student ID
  /// Query params: page, limit
  /// Response: { data: HomeroomNote[], pagination: Pagination }
  static const String parentChildHomeroomNotes = 'parent/children/{id}/homeroom-notes';

  /// GET - Get child achievements
  /// Path param: {id} - Child/Student ID
  /// Query params: page, limit
  /// Response: { data: Achievement[], pagination: Pagination }
  static const String parentChildAchievements = 'parent/children/{id}/achievements';

  /// GET - Get child violations
  /// Path param: {id} - Child/Student ID
  /// Query params: page, limit
  /// Response: { data: Violation[], pagination: Pagination }
  static const String parentChildViolations = 'parent/children/{id}/violations';

  /// GET - Get child permits
  /// Path param: {id} - Child/Student ID
  /// Query params: page, limit
  /// Response: { data: Permit[], pagination: Pagination }
  static const String parentChildPermits = 'parent/children/{id}/permits';

  /// GET - Get child counseling notes (parent summary only)
  /// Path param: {id} - Child/Student ID
  /// Query params: page, limit
  /// Response: { data: CounselingNote[], pagination: Pagination }
  static const String parentChildCounselingNotes = 'parent/children/{id}/counseling-notes';

  // ============================================
  // Student Endpoints
  // ============================================

  /// GET - Get student profile
  /// Response: { profile: StudentProfile }
  static const String studentProfile = 'student/profile';

  /// GET - Get student dashboard
  /// Response: { dashboard: StudentDashboard }
  static const String studentDashboard = 'student/dashboard';

  /// GET - Get student attendance list
  /// Query params: page, limit, start_date, end_date
  /// Response: { data: Attendance[], pagination: Pagination }
  static const String studentAttendance = 'student/attendance';

  /// GET - Get student attendance summary
  /// Query params: start_date, end_date
  /// Response: { summary: AttendanceSummary }
  static const String studentAttendanceSummary = 'student/attendance/summary';

  /// GET - Get student grades list
  /// Query params: page, limit
  /// Response: { data: Grade[], pagination: Pagination }
  static const String studentGrades = 'student/grades';

  /// GET - Get student grades summary
  /// Response: { summary: GradeSummary }
  static const String studentGradesSummary = 'student/grades/summary';

  /// GET - Get student achievements
  /// Query params: page, limit
  /// Response: { data: Achievement[], pagination: Pagination }
  static const String studentAchievements = 'student/achievements';

  /// GET - Get student violations
  /// Query params: page, limit
  /// Response: { data: Violation[], pagination: Pagination }
  static const String studentViolations = 'student/violations';

  // ============================================
  // Notification Endpoints
  // ============================================

  /// GET - Get notifications list
  /// Query params: page, limit, is_read
  /// Response: { data: Notification[], pagination: Pagination }
  static const String notifications = 'notifications';

  /// GET - Get notification summary
  /// Response: { summary: NotificationSummary }
  static const String notificationsSummary = 'notifications/summary';

  /// GET - Get unread notification count
  /// Response: { count: int }
  static const String notificationsUnreadCount = 'notifications/unread-count';

  /// PUT - Mark notification as read
  /// Path param: {id} - Notification ID
  /// Response: { success: bool }
  static const String notificationMarkRead = 'notifications/{id}/read';

  /// PUT - Mark all notifications as read
  /// Response: { success: bool }
  static const String notificationsMarkAllRead = 'notifications/mark-all-read';

  // ============================================
  // FCM (Firebase Cloud Messaging) Endpoints
  // ============================================

  /// POST - Register FCM token
  /// Request: { token: string, platform: string }
  /// Response: { success: bool }
  static const String fcmRegister = 'fcm/register';

  /// POST - Deactivate FCM token
  /// Request: { token: string }
  /// Response: { success: bool }
  static const String fcmDeactivate = 'fcm/deactivate';
}

/// HTTP methods used in API requests.
enum HttpMethod {
  get,
  post,
  put,
  delete,
  patch,
}

/// API response status codes and their meanings.
class ApiStatusCodes {
  ApiStatusCodes._();

  // Success codes
  static const int ok = 200;
  static const int created = 201;
  static const int noContent = 204;

  // Client error codes
  static const int badRequest = 400;
  static const int unauthorized = 401;
  static const int forbidden = 403;
  static const int notFound = 404;
  static const int conflict = 409;
  static const int unprocessableEntity = 422;
  static const int tooManyRequests = 429;

  // Server error codes
  static const int internalServerError = 500;
  static const int badGateway = 502;
  static const int serviceUnavailable = 503;
  static const int gatewayTimeout = 504;
}

/// API error codes returned by the backend.
///
/// These codes are used to provide specific error messages in Indonesian.
class ApiErrorCodes {
  ApiErrorCodes._();

  // Authentication errors
  static const String invalidCredentials = 'AUTH_INVALID_CREDENTIALS';
  static const String accountInactive = 'AUTH_ACCOUNT_INACTIVE';
  static const String schoolInactive = 'AUTH_SCHOOL_INACTIVE';
  static const String tokenExpired = 'AUTH_TOKEN_EXPIRED';
  static const String tokenInvalid = 'AUTH_TOKEN_INVALID';
  static const String passwordMismatch = 'AUTH_PASSWORD_MISMATCH';
  static const String samePassword = 'AUTH_SAME_PASSWORD';

  // Not found errors
  static const String studentNotFound = 'NOT_FOUND_STUDENT';
  static const String parentNotFound = 'NOT_FOUND_PARENT';
  static const String notificationNotFound = 'NOT_FOUND_NOTIFICATION';

  // Authorization errors
  static const String notLinked = 'AUTHZ_NOT_LINKED';
  static const String forbidden = 'AUTHZ_FORBIDDEN';

  // Network errors (client-side)
  static const String networkError = 'NETWORK_ERROR';
  static const String timeout = 'TIMEOUT';

  // Validation errors
  static const String validationError = 'VALIDATION_ERROR';
  static const String invalidInput = 'INVALID_INPUT';
}
