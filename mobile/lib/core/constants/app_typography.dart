import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

import 'app_colors.dart';

/// Application typography styles for the School Management mobile app.
///
/// Uses Poppins font family with a clean, modern look.
/// Falls back to google_fonts package if local fonts are not available.
/// All text styles follow the design system defined in the design document.
///
/// Usage:
/// ```dart
/// Text('Hello', style: AppTypography.h1)
/// Text('Subtitle', style: AppTypography.bodyMedium)
/// ```
class AppTypography {
  AppTypography._();

  /// Primary font family - Poppins
  static const String fontFamily = 'Poppins';

  /// Gets the Poppins TextStyle using google_fonts as fallback.
  /// This ensures the font works even if local font files are not bundled.
  static TextStyle _poppins({
    double fontSize = 14,
    FontWeight fontWeight = FontWeight.w400,
    Color color = AppColors.textPrimary,
    double height = 1.5,
    double? letterSpacing,
    TextDecoration? decoration,
  }) {
    return GoogleFonts.poppins(
      fontSize: fontSize,
      fontWeight: fontWeight,
      color: color,
      height: height,
      letterSpacing: letterSpacing,
      decoration: decoration,
    );
  }

  // ============================================
  // Headings
  // ============================================

  /// Heading 1 - Largest heading (28px, bold)
  static TextStyle get h1 => _poppins(
        fontSize: 28,
        fontWeight: FontWeight.w700,
        color: AppColors.textPrimary,
        height: 1.3,
        letterSpacing: -0.5,
      );

  /// Heading 2 - Large heading (24px, semi-bold)
  static TextStyle get h2 => _poppins(
        fontSize: 24,
        fontWeight: FontWeight.w600,
        color: AppColors.textPrimary,
        height: 1.3,
        letterSpacing: -0.3,
      );

  /// Heading 3 - Medium heading (20px, semi-bold)
  static TextStyle get h3 => _poppins(
        fontSize: 20,
        fontWeight: FontWeight.w600,
        color: AppColors.textPrimary,
        height: 1.4,
        letterSpacing: -0.2,
      );

  /// Heading 4 - Small heading (18px, semi-bold)
  static TextStyle get h4 => _poppins(
        fontSize: 18,
        fontWeight: FontWeight.w600,
        color: AppColors.textPrimary,
        height: 1.4,
      );

  /// Heading 5 - Extra small heading (16px, semi-bold)
  static TextStyle get h5 => _poppins(
        fontSize: 16,
        fontWeight: FontWeight.w600,
        color: AppColors.textPrimary,
        height: 1.4,
      );

  // ============================================
  // Body Text
  // ============================================

  /// Body large - Main content text (16px, regular)
  static TextStyle get bodyLarge => _poppins(
        fontSize: 16,
        fontWeight: FontWeight.w400,
        color: AppColors.textPrimary,
        height: 1.5,
      );

  /// Body medium - Standard content text (14px, regular)
  static TextStyle get bodyMedium => _poppins(
        fontSize: 14,
        fontWeight: FontWeight.w400,
        color: AppColors.textPrimary,
        height: 1.5,
      );

  /// Body small - Secondary content text (12px, regular)
  static TextStyle get bodySmall => _poppins(
        fontSize: 12,
        fontWeight: FontWeight.w400,
        color: AppColors.textSecondary,
        height: 1.5,
      );

  // ============================================
  // Labels
  // ============================================

  /// Label large - Form labels, card titles (14px, medium)
  static TextStyle get labelLarge => _poppins(
        fontSize: 14,
        fontWeight: FontWeight.w500,
        color: AppColors.textPrimary,
        height: 1.4,
      );

  /// Label medium - Secondary labels (12px, medium)
  static TextStyle get labelMedium => _poppins(
        fontSize: 12,
        fontWeight: FontWeight.w500,
        color: AppColors.textSecondary,
        height: 1.4,
      );

  /// Label small - Captions, timestamps (10px, medium)
  static TextStyle get labelSmall => _poppins(
        fontSize: 10,
        fontWeight: FontWeight.w500,
        color: AppColors.textTertiary,
        height: 1.4,
      );

  // ============================================
  // Button Text
  // ============================================

  /// Button text - Primary buttons (16px, semi-bold)
  static TextStyle get button => _poppins(
        fontSize: 16,
        fontWeight: FontWeight.w600,
        color: AppColors.white,
        height: 1.2,
      );

  /// Button text small - Secondary buttons (14px, semi-bold)
  static TextStyle get buttonSmall => _poppins(
        fontSize: 14,
        fontWeight: FontWeight.w600,
        color: AppColors.white,
        height: 1.2,
      );

  // ============================================
  // Special Text Styles
  // ============================================

  /// App bar title style
  static TextStyle get appBarTitle => _poppins(
        fontSize: 18,
        fontWeight: FontWeight.w600,
        color: AppColors.textPrimary,
        height: 1.3,
      );

  /// Bottom navigation label
  static TextStyle get bottomNavLabel => _poppins(
        fontSize: 12,
        fontWeight: FontWeight.w500,
        height: 1.2,
      );

  /// Tab label
  static TextStyle get tabLabel => _poppins(
        fontSize: 14,
        fontWeight: FontWeight.w500,
        height: 1.2,
      );

  /// Badge text (notification count, etc.)
  static TextStyle get badge => _poppins(
        fontSize: 10,
        fontWeight: FontWeight.w600,
        color: AppColors.white,
        height: 1.0,
      );

  /// Large number display (statistics, scores)
  static TextStyle get displayLarge => _poppins(
        fontSize: 32,
        fontWeight: FontWeight.w700,
        color: AppColors.textPrimary,
        height: 1.2,
      );

  /// Medium number display
  static TextStyle get displayMedium => _poppins(
        fontSize: 24,
        fontWeight: FontWeight.w700,
        color: AppColors.textPrimary,
        height: 1.2,
      );

  /// Small number display
  static TextStyle get displaySmall => _poppins(
        fontSize: 20,
        fontWeight: FontWeight.w700,
        color: AppColors.textPrimary,
        height: 1.2,
      );

  /// Link text style
  static TextStyle get link => _poppins(
        fontSize: 14,
        fontWeight: FontWeight.w500,
        color: AppColors.primary,
        height: 1.5,
        decoration: TextDecoration.underline,
      );

  /// Error text style
  static TextStyle get error => _poppins(
        fontSize: 12,
        fontWeight: FontWeight.w400,
        color: AppColors.error,
        height: 1.4,
      );

  /// Hint text style
  static TextStyle get hint => _poppins(
        fontSize: 14,
        fontWeight: FontWeight.w400,
        color: AppColors.textTertiary,
        height: 1.5,
      );

  // ============================================
  // Helper Methods
  // ============================================

  /// Returns a text style with a custom color.
  static TextStyle withColor(TextStyle style, Color color) {
    return style.copyWith(color: color);
  }

  /// Returns a text style with a custom font weight.
  static TextStyle withWeight(TextStyle style, FontWeight weight) {
    return style.copyWith(fontWeight: weight);
  }

  /// Returns a text style for secondary text (gray color).
  static TextStyle secondary(TextStyle style) {
    return style.copyWith(color: AppColors.textSecondary);
  }

  /// Returns a text style for tertiary text (light gray color).
  static TextStyle tertiary(TextStyle style) {
    return style.copyWith(color: AppColors.textTertiary);
  }
}
