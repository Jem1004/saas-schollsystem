import 'package:flutter/material.dart';

/// Application color palette for the School Management mobile app.
///
/// Uses a modern minimalist design with primary orange (#FF6B00) and white.
/// All colors follow the design system defined in the design document.
///
/// Usage:
/// ```dart
/// Container(
///   color: AppColors.primary,
///   child: Text('Hello', style: TextStyle(color: AppColors.textOnPrimary)),
/// )
/// ```
class AppColors {
  AppColors._();

  // ============================================
  // Primary Colors
  // ============================================

  /// Primary orange color - main brand color
  static const Color primary = Color(0xFFFF6B00);

  /// Light orange - for hover states and accents
  static const Color primaryLight = Color(0xFFFF8A3D);

  /// Dark orange - for pressed states
  static const Color primaryDark = Color(0xFFE55A00);

  /// Very light orange - for backgrounds and surfaces
  static const Color primarySurface = Color(0xFFFFF3E8);

  // ============================================
  // Neutral Colors
  // ============================================

  /// Pure white
  static const Color white = Color(0xFFFFFFFF);

  /// Light gray background - main app background
  static const Color background = Color(0xFFF8F9FA);

  /// Card/surface background - white
  static const Color surface = Color(0xFFFFFFFF);

  /// Border color for cards and inputs
  static const Color border = Color(0xFFE8E8E8);

  /// Divider color - lighter than border
  static const Color divider = Color(0xFFF0F0F0);

  /// Black color for overlays
  static const Color black = Color(0xFF000000);

  // ============================================
  // Text Colors
  // ============================================

  /// Primary text color - main content
  static const Color textPrimary = Color(0xFF1A1A1A);

  /// Secondary text color - supporting content
  static const Color textSecondary = Color(0xFF6B7280);

  /// Tertiary text color - hints and placeholders
  static const Color textTertiary = Color(0xFF9CA3AF);

  /// Text color on primary (orange) backgrounds
  static const Color textOnPrimary = Color(0xFFFFFFFF);

  // ============================================
  // Status Colors
  // ============================================

  /// Success/positive - green
  static const Color success = Color(0xFF22C55E);

  /// Light success background
  static const Color successLight = Color(0xFFDCFCE7);

  /// Warning - yellow
  static const Color warning = Color(0xFFFACC15);

  /// Light warning background
  static const Color warningLight = Color(0xFFFEF9C3);

  /// Error/negative - red
  static const Color error = Color(0xFFEF4444);

  /// Light error background
  static const Color errorLight = Color(0xFFFEE2E2);

  /// Info - blue
  static const Color info = Color(0xFF3B82F6);

  /// Light info background
  static const Color infoLight = Color(0xFFDBEAFE);

  /// Orange for medium severity (violations, very late)
  static const Color orange = Color(0xFFF97316);

  /// Light orange background
  static const Color orangeLight = Color(0xFFFFF7ED);

  // ============================================
  // Attendance Status Colors
  // ============================================

  /// On time - green (same as success)
  static const Color onTime = Color(0xFF22C55E);

  /// Late - yellow (same as warning)
  static const Color late = Color(0xFFFACC15);

  /// Very late - orange
  static const Color veryLate = Color(0xFFF97316);

  /// Absent - red (same as error)
  static const Color absent = Color(0xFFEF4444);

  /// Sick - blue (same as info)
  static const Color sick = Color(0xFF3B82F6);

  /// Excused - purple
  static const Color excused = Color(0xFF8B5CF6);

  // ============================================
  // Skeleton Loading Colors
  // ============================================

  /// Base skeleton color
  static const Color skeleton = Color(0xFFE5E7EB);

  /// Skeleton highlight/shimmer color
  static const Color skeletonHighlight = Color(0xFFF3F4F6);

  // ============================================
  // Grade Colors
  // ============================================

  /// High grade (>=80) - green
  static const Color gradeHigh = Color(0xFF22C55E);

  /// Medium grade (60-79) - yellow
  static const Color gradeMedium = Color(0xFFFACC15);

  /// Low grade (<60) - red
  static const Color gradeLow = Color(0xFFEF4444);

  // ============================================
  // Violation Level Colors
  // ============================================

  /// Light violation (ringan) - yellow
  static const Color violationLight = Color(0xFFFACC15);

  /// Medium violation (sedang) - orange
  static const Color violationMedium = Color(0xFFF97316);

  /// Severe violation (berat) - red
  static const Color violationSevere = Color(0xFFEF4444);

  // ============================================
  // Helper Methods
  // ============================================

  /// Returns the appropriate color for an attendance status.
  static Color getAttendanceStatusColor(String status) {
    switch (status.toLowerCase()) {
      case 'on_time':
      case 'tepat_waktu':
        return onTime;
      case 'late':
      case 'terlambat':
        return late;
      case 'very_late':
      case 'sangat_terlambat':
        return veryLate;
      case 'absent':
      case 'tidak_hadir':
      case 'alpha':
        return absent;
      case 'sick':
      case 'sakit':
        return sick;
      case 'excused':
      case 'izin':
        return excused;
      default:
        return textSecondary;
    }
  }

  /// Returns the appropriate color for a grade score.
  static Color getGradeColor(double score) {
    if (score >= 80) return gradeHigh;
    if (score >= 60) return gradeMedium;
    return gradeLow;
  }

  /// Returns the appropriate color for a violation level.
  static Color getViolationLevelColor(String level) {
    switch (level.toLowerCase()) {
      case 'ringan':
      case 'light':
        return violationLight;
      case 'sedang':
      case 'medium':
        return violationMedium;
      case 'berat':
      case 'severe':
        return violationSevere;
      default:
        return textSecondary;
    }
  }

  /// Returns a color with reduced opacity for backgrounds.
  static Color withLightBackground(Color color, [double opacity = 0.15]) {
    return color.withOpacity(opacity);
  }
}
