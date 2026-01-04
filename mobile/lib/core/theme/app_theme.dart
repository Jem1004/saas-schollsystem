import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../constants/app_colors.dart';
import '../constants/app_spacing.dart';
import '../constants/app_typography.dart';

/// Application theme configuration for the School Management mobile app.
///
/// Provides a complete Material 3 theme with custom colors, typography,
/// and component styles following the design system.
///
/// Usage:
/// ```dart
/// MaterialApp(
///   theme: AppTheme.lightTheme,
///   // ...
/// )
/// ```
class AppTheme {
  AppTheme._();

  /// Light theme for the application.
  static ThemeData get lightTheme {
    return ThemeData(
      useMaterial3: true,
      brightness: Brightness.light,

      // Color Scheme
      colorScheme: _lightColorScheme,

      // Scaffold
      scaffoldBackgroundColor: AppColors.background,

      // App Bar
      appBarTheme: _appBarTheme,

      // Bottom Navigation
      bottomNavigationBarTheme: _bottomNavigationBarTheme,

      // Tab Bar
      tabBarTheme: _tabBarTheme,

      // Card
      cardTheme: _cardTheme,

      // Elevated Button
      elevatedButtonTheme: _elevatedButtonTheme,

      // Outlined Button
      outlinedButtonTheme: _outlinedButtonTheme,

      // Text Button
      textButtonTheme: _textButtonTheme,

      // Input Decoration
      inputDecorationTheme: _inputDecorationTheme,

      // Floating Action Button
      floatingActionButtonTheme: _floatingActionButtonTheme,

      // Chip
      chipTheme: _chipTheme,

      // Dialog
      dialogTheme: _dialogTheme,

      // Bottom Sheet
      bottomSheetTheme: _bottomSheetTheme,

      // Snack Bar
      snackBarTheme: _snackBarTheme,

      // Divider
      dividerTheme: _dividerTheme,

      // List Tile
      listTileTheme: _listTileTheme,

      // Icon
      iconTheme: _iconTheme,

      // Text Theme
      textTheme: _textTheme,

      // Font Family
      fontFamily: AppTypography.fontFamily,

      // Splash and Highlight
      splashColor: AppColors.primary.withOpacity(0.1),
      highlightColor: AppColors.primary.withOpacity(0.05),

      // Visual Density
      visualDensity: VisualDensity.adaptivePlatformDensity,
    );
  }

  // ============================================
  // Color Scheme
  // ============================================

  static ColorScheme get _lightColorScheme {
    return ColorScheme.light(
      primary: AppColors.primary,
      onPrimary: AppColors.textOnPrimary,
      primaryContainer: AppColors.primarySurface,
      onPrimaryContainer: AppColors.primaryDark,
      secondary: AppColors.primary,
      onSecondary: AppColors.textOnPrimary,
      secondaryContainer: AppColors.primarySurface,
      onSecondaryContainer: AppColors.primaryDark,
      tertiary: AppColors.info,
      onTertiary: AppColors.white,
      tertiaryContainer: AppColors.infoLight,
      onTertiaryContainer: AppColors.info,
      error: AppColors.error,
      onError: AppColors.white,
      errorContainer: AppColors.errorLight,
      onErrorContainer: AppColors.error,
      surface: AppColors.surface,
      onSurface: AppColors.textPrimary,
      surfaceContainerHighest: AppColors.background,
      onSurfaceVariant: AppColors.textSecondary,
      outline: AppColors.border,
      outlineVariant: AppColors.divider,
      shadow: AppColors.black.withOpacity(0.1),
      scrim: AppColors.black.withOpacity(0.5),
    );
  }

  // ============================================
  // App Bar Theme
  // ============================================

  static AppBarTheme get _appBarTheme {
    return AppBarTheme(
      elevation: 0,
      scrolledUnderElevation: 1,
      centerTitle: false,
      backgroundColor: AppColors.surface,
      foregroundColor: AppColors.textPrimary,
      surfaceTintColor: Colors.transparent,
      shadowColor: AppColors.black.withOpacity(0.1),
      titleTextStyle: AppTypography.appBarTitle,
      iconTheme: const IconThemeData(
        color: AppColors.textPrimary,
        size: AppSpacing.iconMd,
      ),
      actionsIconTheme: const IconThemeData(
        color: AppColors.textPrimary,
        size: AppSpacing.iconMd,
      ),
      systemOverlayStyle: SystemUiOverlayStyle.dark.copyWith(
        statusBarColor: Colors.transparent,
        statusBarIconBrightness: Brightness.dark,
      ),
    );
  }

  // ============================================
  // Bottom Navigation Bar Theme
  // ============================================

  static BottomNavigationBarThemeData get _bottomNavigationBarTheme {
    return BottomNavigationBarThemeData(
      type: BottomNavigationBarType.fixed,
      backgroundColor: AppColors.surface,
      elevation: 8,
      selectedItemColor: AppColors.primary,
      unselectedItemColor: AppColors.textTertiary,
      selectedLabelStyle: AppTypography.bottomNavLabel.copyWith(
        color: AppColors.primary,
      ),
      unselectedLabelStyle: AppTypography.bottomNavLabel.copyWith(
        color: AppColors.textTertiary,
      ),
      showSelectedLabels: true,
      showUnselectedLabels: true,
      selectedIconTheme: const IconThemeData(
        size: AppSpacing.iconMd,
        color: AppColors.primary,
      ),
      unselectedIconTheme: const IconThemeData(
        size: AppSpacing.iconMd,
        color: AppColors.textTertiary,
      ),
    );
  }

  // ============================================
  // Tab Bar Theme
  // ============================================

  static TabBarTheme get _tabBarTheme {
    return TabBarTheme(
      labelColor: AppColors.primary,
      unselectedLabelColor: AppColors.textSecondary,
      labelStyle: AppTypography.tabLabel.copyWith(
        fontWeight: FontWeight.w600,
      ),
      unselectedLabelStyle: AppTypography.tabLabel,
      indicatorColor: AppColors.primary,
      indicatorSize: TabBarIndicatorSize.tab,
      dividerColor: AppColors.divider,
      overlayColor: WidgetStateProperty.all(
        AppColors.primary.withOpacity(0.1),
      ),
    );
  }

  // ============================================
  // Card Theme
  // ============================================

  static CardTheme get _cardTheme {
    return CardTheme(
      elevation: 0,
      color: AppColors.surface,
      surfaceTintColor: Colors.transparent,
      shadowColor: AppColors.black.withOpacity(0.08),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
      ),
      margin: EdgeInsets.zero,
    );
  }

  // ============================================
  // Elevated Button Theme
  // ============================================

  static ElevatedButtonThemeData get _elevatedButtonTheme {
    return ElevatedButtonThemeData(
      style: ElevatedButton.styleFrom(
        elevation: 0,
        backgroundColor: AppColors.primary,
        foregroundColor: AppColors.textOnPrimary,
        disabledBackgroundColor: AppColors.primary.withOpacity(0.5),
        disabledForegroundColor: AppColors.textOnPrimary.withOpacity(0.7),
        shadowColor: Colors.transparent,
        padding: AppSpacing.buttonPadding,
        minimumSize: const Size(double.infinity, AppSpacing.buttonHeight),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        ),
        textStyle: AppTypography.button,
      ),
    );
  }

  // ============================================
  // Outlined Button Theme
  // ============================================

  static OutlinedButtonThemeData get _outlinedButtonTheme {
    return OutlinedButtonThemeData(
      style: OutlinedButton.styleFrom(
        foregroundColor: AppColors.primary,
        disabledForegroundColor: AppColors.primary.withOpacity(0.5),
        padding: AppSpacing.buttonPadding,
        minimumSize: const Size(double.infinity, AppSpacing.buttonHeight),
        side: const BorderSide(
          color: AppColors.primary,
          width: AppSpacing.borderThick,
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        ),
        textStyle: AppTypography.button.copyWith(
          color: AppColors.primary,
        ),
      ),
    );
  }

  // ============================================
  // Text Button Theme
  // ============================================

  static TextButtonThemeData get _textButtonTheme {
    return TextButtonThemeData(
      style: TextButton.styleFrom(
        foregroundColor: AppColors.primary,
        disabledForegroundColor: AppColors.primary.withOpacity(0.5),
        padding: AppSpacing.buttonPaddingSmall,
        minimumSize: const Size(0, AppSpacing.buttonHeightSmall),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(AppSpacing.radiusSm),
        ),
        textStyle: AppTypography.labelLarge.copyWith(
          color: AppColors.primary,
        ),
      ),
    );
  }

  // ============================================
  // Input Decoration Theme
  // ============================================

  static InputDecorationTheme get _inputDecorationTheme {
    return InputDecorationTheme(
      filled: true,
      fillColor: AppColors.background,
      contentPadding: AppSpacing.inputPadding,
      hintStyle: AppTypography.bodyMedium.copyWith(
        color: AppColors.textTertiary,
      ),
      labelStyle: AppTypography.labelLarge.copyWith(
        color: AppColors.textSecondary,
      ),
      floatingLabelStyle: AppTypography.labelMedium.copyWith(
        color: AppColors.primary,
      ),
      errorStyle: AppTypography.error,
      prefixIconColor: AppColors.textSecondary,
      suffixIconColor: AppColors.textSecondary,
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
        borderSide: const BorderSide(
          color: AppColors.primary,
          width: AppSpacing.borderThick,
        ),
      ),
      errorBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        borderSide: const BorderSide(
          color: AppColors.error,
          width: AppSpacing.borderThick,
        ),
      ),
      focusedErrorBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        borderSide: const BorderSide(
          color: AppColors.error,
          width: AppSpacing.borderThick,
        ),
      ),
      disabledBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        borderSide: BorderSide.none,
      ),
    );
  }

  // ============================================
  // Floating Action Button Theme
  // ============================================

  static FloatingActionButtonThemeData get _floatingActionButtonTheme {
    return FloatingActionButtonThemeData(
      elevation: AppSpacing.elevationMedium,
      highlightElevation: AppSpacing.elevationHigh,
      backgroundColor: AppColors.primary,
      foregroundColor: AppColors.textOnPrimary,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
      ),
    );
  }

  // ============================================
  // Chip Theme
  // ============================================

  static ChipThemeData get _chipTheme {
    return ChipThemeData(
      backgroundColor: AppColors.background,
      selectedColor: AppColors.primarySurface,
      disabledColor: AppColors.background,
      labelStyle: AppTypography.labelMedium,
      secondaryLabelStyle: AppTypography.labelMedium.copyWith(
        color: AppColors.primary,
      ),
      padding: const EdgeInsets.symmetric(
        horizontal: AppSpacing.sm,
        vertical: AppSpacing.xs,
      ),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusSm),
      ),
      side: const BorderSide(color: AppColors.border),
    );
  }

  // ============================================
  // Dialog Theme
  // ============================================

  static DialogTheme get _dialogTheme {
    return DialogTheme(
      elevation: AppSpacing.elevationHigh,
      backgroundColor: AppColors.surface,
      surfaceTintColor: Colors.transparent,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusLg),
      ),
      titleTextStyle: AppTypography.h4,
      contentTextStyle: AppTypography.bodyMedium,
    );
  }

  // ============================================
  // Bottom Sheet Theme
  // ============================================

  static BottomSheetThemeData get _bottomSheetTheme {
    return BottomSheetThemeData(
      elevation: AppSpacing.elevationHigh,
      backgroundColor: AppColors.surface,
      surfaceTintColor: Colors.transparent,
      modalElevation: AppSpacing.elevationHigh,
      modalBackgroundColor: AppColors.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: AppSpacing.borderRadiusTopXl,
      ),
      dragHandleColor: AppColors.border,
      dragHandleSize: const Size(40, 4),
      showDragHandle: true,
    );
  }

  // ============================================
  // Snack Bar Theme
  // ============================================

  static SnackBarThemeData get _snackBarTheme {
    return SnackBarThemeData(
      elevation: AppSpacing.elevationMedium,
      backgroundColor: AppColors.textPrimary,
      contentTextStyle: AppTypography.bodyMedium.copyWith(
        color: AppColors.white,
      ),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusSm),
      ),
      behavior: SnackBarBehavior.floating,
      insetPadding: AppSpacing.pagePadding,
      actionTextColor: AppColors.primaryLight,
    );
  }

  // ============================================
  // Divider Theme
  // ============================================

  static DividerThemeData get _dividerTheme {
    return const DividerThemeData(
      color: AppColors.divider,
      thickness: AppSpacing.borderDefault,
      space: 0,
    );
  }

  // ============================================
  // List Tile Theme
  // ============================================

  static ListTileThemeData get _listTileTheme {
    return ListTileThemeData(
      contentPadding: AppSpacing.listItemPadding,
      minVerticalPadding: AppSpacing.sm,
      horizontalTitleGap: AppSpacing.md,
      tileColor: Colors.transparent,
      selectedTileColor: AppColors.primarySurface,
      iconColor: AppColors.textSecondary,
      textColor: AppColors.textPrimary,
      titleTextStyle: AppTypography.bodyLarge,
      subtitleTextStyle: AppTypography.bodySmall,
      leadingAndTrailingTextStyle: AppTypography.labelMedium,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusSm),
      ),
    );
  }

  // ============================================
  // Icon Theme
  // ============================================

  static IconThemeData get _iconTheme {
    return const IconThemeData(
      color: AppColors.textSecondary,
      size: AppSpacing.iconMd,
    );
  }

  // ============================================
  // Text Theme
  // ============================================

  static TextTheme get _textTheme {
    return TextTheme(
      displayLarge: AppTypography.displayLarge,
      displayMedium: AppTypography.displayMedium,
      displaySmall: AppTypography.displaySmall,
      headlineLarge: AppTypography.h1,
      headlineMedium: AppTypography.h2,
      headlineSmall: AppTypography.h3,
      titleLarge: AppTypography.h4,
      titleMedium: AppTypography.h5,
      titleSmall: AppTypography.labelLarge,
      bodyLarge: AppTypography.bodyLarge,
      bodyMedium: AppTypography.bodyMedium,
      bodySmall: AppTypography.bodySmall,
      labelLarge: AppTypography.labelLarge,
      labelMedium: AppTypography.labelMedium,
      labelSmall: AppTypography.labelSmall,
    );
  }
}
