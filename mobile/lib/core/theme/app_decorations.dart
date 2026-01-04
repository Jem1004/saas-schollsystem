import 'package:flutter/material.dart';

import '../constants/app_colors.dart';
import '../constants/app_spacing.dart';
import '../constants/app_typography.dart';

/// Application decoration styles for the School Management mobile app.
///
/// Provides reusable BoxDecoration, InputDecoration, and other decoration
/// styles following the design system.
///
/// Usage:
/// ```dart
/// Container(
///   decoration: AppDecorations.card,
///   child: // ...
/// )
/// ```
class AppDecorations {
  AppDecorations._();

  // ============================================
  // Card Decorations
  // ============================================

  /// Standard card decoration with subtle shadow.
  static BoxDecoration get card {
    return BoxDecoration(
      color: AppColors.surface,
      borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
      boxShadow: [
        BoxShadow(
          color: AppColors.black.withOpacity(0.04),
          blurRadius: 8,
          offset: const Offset(0, 2),
        ),
      ],
    );
  }

  /// Elevated card decoration with stronger shadow.
  static BoxDecoration get cardElevated {
    return BoxDecoration(
      color: AppColors.surface,
      borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
      boxShadow: [
        BoxShadow(
          color: AppColors.black.withOpacity(0.08),
          blurRadius: 16,
          offset: const Offset(0, 4),
        ),
      ],
    );
  }

  /// Flat card decoration without shadow.
  static BoxDecoration get cardFlat {
    return BoxDecoration(
      color: AppColors.surface,
      borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
      border: Border.all(color: AppColors.border),
    );
  }

  /// Card decoration with primary color border.
  static BoxDecoration get cardPrimary {
    return BoxDecoration(
      color: AppColors.primarySurface,
      borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
      border: Border.all(color: AppColors.primary.withOpacity(0.3)),
    );
  }

  /// Card decoration for selected/active state.
  static BoxDecoration get cardSelected {
    return BoxDecoration(
      color: AppColors.primarySurface,
      borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
      border: Border.all(color: AppColors.primary, width: AppSpacing.borderThick),
    );
  }

  // ============================================
  // Button Decorations
  // ============================================

  /// Primary button decoration.
  static BoxDecoration get primaryButton {
    return BoxDecoration(
      color: AppColors.primary,
      borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
    );
  }

  /// Primary button decoration with gradient.
  static BoxDecoration get primaryButtonGradient {
    return BoxDecoration(
      gradient: const LinearGradient(
        colors: [AppColors.primary, AppColors.primaryDark],
        begin: Alignment.topLeft,
        end: Alignment.bottomRight,
      ),
      borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
    );
  }

  /// Outlined button decoration.
  static BoxDecoration get outlinedButton {
    return BoxDecoration(
      color: AppColors.white,
      borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
      border: Border.all(color: AppColors.primary, width: AppSpacing.borderThick),
    );
  }

  /// Secondary button decoration.
  static BoxDecoration get secondaryButton {
    return BoxDecoration(
      color: AppColors.primarySurface,
      borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
    );
  }

  /// Disabled button decoration.
  static BoxDecoration get disabledButton {
    return BoxDecoration(
      color: AppColors.primary.withOpacity(0.5),
      borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
    );
  }

  // ============================================
  // Status Badge Decorations
  // ============================================

  /// Creates a status badge decoration with the given color.
  static BoxDecoration statusBadge(Color color, {double opacity = 0.15}) {
    return BoxDecoration(
      color: color.withOpacity(opacity),
      borderRadius: BorderRadius.circular(AppSpacing.radiusSm),
    );
  }

  /// Success badge decoration.
  static BoxDecoration get successBadge {
    return statusBadge(AppColors.success);
  }

  /// Warning badge decoration.
  static BoxDecoration get warningBadge {
    return statusBadge(AppColors.warning);
  }

  /// Error badge decoration.
  static BoxDecoration get errorBadge {
    return statusBadge(AppColors.error);
  }

  /// Info badge decoration.
  static BoxDecoration get infoBadge {
    return statusBadge(AppColors.info);
  }

  /// Primary badge decoration.
  static BoxDecoration get primaryBadge {
    return statusBadge(AppColors.primary);
  }

  // ============================================
  // Attendance Status Decorations
  // ============================================

  /// On time attendance badge.
  static BoxDecoration get onTimeBadge {
    return statusBadge(AppColors.onTime);
  }

  /// Late attendance badge.
  static BoxDecoration get lateBadge {
    return statusBadge(AppColors.late);
  }

  /// Very late attendance badge.
  static BoxDecoration get veryLateBadge {
    return statusBadge(AppColors.veryLate);
  }

  /// Absent attendance badge.
  static BoxDecoration get absentBadge {
    return statusBadge(AppColors.absent);
  }

  /// Sick attendance badge.
  static BoxDecoration get sickBadge {
    return statusBadge(AppColors.sick);
  }

  /// Excused attendance badge.
  static BoxDecoration get excusedBadge {
    return statusBadge(AppColors.excused);
  }

  /// Returns the appropriate badge decoration for an attendance status.
  static BoxDecoration getAttendanceBadge(String status) {
    switch (status.toLowerCase()) {
      case 'on_time':
      case 'tepat_waktu':
        return onTimeBadge;
      case 'late':
      case 'terlambat':
        return lateBadge;
      case 'very_late':
      case 'sangat_terlambat':
        return veryLateBadge;
      case 'absent':
      case 'tidak_hadir':
      case 'alpha':
        return absentBadge;
      case 'sick':
      case 'sakit':
        return sickBadge;
      case 'excused':
      case 'izin':
        return excusedBadge;
      default:
        return statusBadge(AppColors.textSecondary);
    }
  }

  // ============================================
  // Container Decorations
  // ============================================

  /// Rounded container decoration.
  static BoxDecoration get roundedContainer {
    return BoxDecoration(
      color: AppColors.background,
      borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
    );
  }

  /// Circular container decoration.
  static BoxDecoration get circularContainer {
    return const BoxDecoration(
      color: AppColors.background,
      shape: BoxShape.circle,
    );
  }

  /// Primary circular container decoration.
  static BoxDecoration get primaryCircle {
    return const BoxDecoration(
      color: AppColors.primary,
      shape: BoxShape.circle,
    );
  }

  /// Light primary circular container decoration.
  static BoxDecoration get primaryCircleLight {
    return const BoxDecoration(
      color: AppColors.primarySurface,
      shape: BoxShape.circle,
    );
  }

  // ============================================
  // Bottom Sheet Decorations
  // ============================================

  /// Bottom sheet decoration with top rounded corners.
  static BoxDecoration get bottomSheet {
    return const BoxDecoration(
      color: AppColors.surface,
      borderRadius: AppSpacing.borderRadiusTopXl,
    );
  }

  // ============================================
  // Input Decorations
  // ============================================

  /// Standard text field input decoration.
  static InputDecoration textField({
    String? hint,
    String? label,
    Widget? prefix,
    Widget? suffix,
    Widget? prefixIcon,
    Widget? suffixIcon,
    String? errorText,
    bool enabled = true,
  }) {
    return InputDecoration(
      hintText: hint,
      labelText: label,
      hintStyle: AppTypography.bodyMedium.copyWith(color: AppColors.textTertiary),
      labelStyle: AppTypography.labelLarge.copyWith(color: AppColors.textSecondary),
      prefix: prefix,
      suffix: suffix,
      prefixIcon: prefixIcon,
      suffixIcon: suffixIcon,
      errorText: errorText,
      filled: true,
      fillColor: enabled ? AppColors.background : AppColors.background.withOpacity(0.5),
      contentPadding: AppSpacing.inputPadding,
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        borderSide: BorderSide.none,
      ),
      enabledBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        borderSide: BorderSide.none,
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        borderSide: const BorderSide(color: AppColors.primary, width: AppSpacing.borderThick),
      ),
      errorBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        borderSide: const BorderSide(color: AppColors.error, width: AppSpacing.borderThick),
      ),
      focusedErrorBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        borderSide: const BorderSide(color: AppColors.error, width: AppSpacing.borderThick),
      ),
      disabledBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        borderSide: BorderSide.none,
      ),
    );
  }

  /// Search field input decoration.
  static InputDecoration searchField({
    String hint = 'Cari...',
    VoidCallback? onClear,
  }) {
    return InputDecoration(
      hintText: hint,
      hintStyle: AppTypography.bodyMedium.copyWith(color: AppColors.textTertiary),
      prefixIcon: const Icon(Icons.search, color: AppColors.textSecondary),
      suffixIcon: onClear != null
          ? IconButton(
              icon: const Icon(Icons.clear, color: AppColors.textSecondary),
              onPressed: onClear,
            )
          : null,
      filled: true,
      fillColor: AppColors.background,
      contentPadding: const EdgeInsets.symmetric(
        horizontal: AppSpacing.md,
        vertical: AppSpacing.sm,
      ),
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusFull),
        borderSide: BorderSide.none,
      ),
      enabledBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusFull),
        borderSide: BorderSide.none,
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusFull),
        borderSide: const BorderSide(color: AppColors.primary, width: AppSpacing.borderThick),
      ),
    );
  }

  // ============================================
  // Shadow Decorations
  // ============================================

  /// Light shadow for subtle elevation.
  static List<BoxShadow> get shadowLight {
    return [
      BoxShadow(
        color: AppColors.black.withOpacity(0.04),
        blurRadius: 8,
        offset: const Offset(0, 2),
      ),
    ];
  }

  /// Medium shadow for cards.
  static List<BoxShadow> get shadowMedium {
    return [
      BoxShadow(
        color: AppColors.black.withOpacity(0.08),
        blurRadius: 16,
        offset: const Offset(0, 4),
      ),
    ];
  }

  /// Heavy shadow for modals and dialogs.
  static List<BoxShadow> get shadowHeavy {
    return [
      BoxShadow(
        color: AppColors.black.withOpacity(0.12),
        blurRadius: 24,
        offset: const Offset(0, 8),
      ),
    ];
  }

  // ============================================
  // Gradient Decorations
  // ============================================

  /// Primary gradient (orange).
  static LinearGradient get primaryGradient {
    return const LinearGradient(
      colors: [AppColors.primary, AppColors.primaryDark],
      begin: Alignment.topLeft,
      end: Alignment.bottomRight,
    );
  }

  /// Success gradient (green).
  static LinearGradient get successGradient {
    return LinearGradient(
      colors: [AppColors.success, AppColors.success.withGreen(180)],
      begin: Alignment.topLeft,
      end: Alignment.bottomRight,
    );
  }

  /// Background gradient for headers.
  static LinearGradient get headerGradient {
    return const LinearGradient(
      colors: [AppColors.primary, AppColors.primaryLight],
      begin: Alignment.topCenter,
      end: Alignment.bottomCenter,
    );
  }

  // ============================================
  // Skeleton Loading Decorations
  // ============================================

  /// Skeleton loading decoration.
  static BoxDecoration get skeleton {
    return BoxDecoration(
      color: AppColors.skeleton,
      borderRadius: BorderRadius.circular(AppSpacing.radiusSm),
    );
  }

  /// Skeleton loading decoration for circular elements.
  static BoxDecoration get skeletonCircle {
    return const BoxDecoration(
      color: AppColors.skeleton,
      shape: BoxShape.circle,
    );
  }
}
