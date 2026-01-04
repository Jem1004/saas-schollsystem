import 'package:flutter/material.dart';

/// Application spacing constants for the School Management mobile app.
///
/// Uses an 8px grid system for consistent spacing throughout the app.
/// All spacing values follow the design system defined in the design document.
///
/// Usage:
/// ```dart
/// Padding(
///   padding: AppSpacing.pagePadding,
///   child: Column(
///     children: [
///       SizedBox(height: AppSpacing.md),
///       // content
///     ],
///   ),
/// )
/// ```
class AppSpacing {
  AppSpacing._();

  // ============================================
  // Base Spacing Values (8px Grid)
  // ============================================

  /// Extra small spacing - 4px
  static const double xs = 4.0;

  /// Small spacing - 8px
  static const double sm = 8.0;

  /// Medium spacing - 16px (base unit)
  static const double md = 16.0;

  /// Large spacing - 24px
  static const double lg = 24.0;

  /// Extra large spacing - 32px
  static const double xl = 32.0;

  /// Extra extra large spacing - 48px
  static const double xxl = 48.0;

  /// Extra extra extra large spacing - 64px
  static const double xxxl = 64.0;

  // ============================================
  // Common Padding Presets
  // ============================================

  /// Page padding - standard padding for screen content
  static const EdgeInsets pagePadding = EdgeInsets.all(md);

  /// Page horizontal padding only
  static const EdgeInsets pageHorizontalPadding = EdgeInsets.symmetric(horizontal: md);

  /// Card padding - internal padding for cards
  static const EdgeInsets cardPadding = EdgeInsets.all(md);

  /// Card padding small - for compact cards
  static const EdgeInsets cardPaddingSmall = EdgeInsets.all(sm);

  /// List item padding - for list tiles and items
  static const EdgeInsets listItemPadding = EdgeInsets.symmetric(
    horizontal: md,
    vertical: sm + xs, // 12px
  );

  /// List item padding compact
  static const EdgeInsets listItemPaddingCompact = EdgeInsets.symmetric(
    horizontal: md,
    vertical: sm,
  );

  /// Button padding - internal padding for buttons
  static const EdgeInsets buttonPadding = EdgeInsets.symmetric(
    horizontal: lg,
    vertical: sm + xs, // 12px
  );

  /// Button padding small
  static const EdgeInsets buttonPaddingSmall = EdgeInsets.symmetric(
    horizontal: md,
    vertical: sm,
  );

  /// Input field padding
  static const EdgeInsets inputPadding = EdgeInsets.symmetric(
    horizontal: md,
    vertical: sm + xs + xs, // 14px
  );

  /// Dialog padding
  static const EdgeInsets dialogPadding = EdgeInsets.all(lg);

  /// Bottom sheet padding
  static const EdgeInsets bottomSheetPadding = EdgeInsets.fromLTRB(md, lg, md, md);

  /// Section padding - for separating content sections
  static const EdgeInsets sectionPadding = EdgeInsets.symmetric(vertical: lg);

  // ============================================
  // Border Radius Values
  // ============================================

  /// Small border radius - 8px (for badges, small elements)
  static const double radiusSm = 8.0;

  /// Medium border radius - 12px (for cards, buttons)
  static const double radiusMd = 12.0;

  /// Large border radius - 16px (for modals, large cards)
  static const double radiusLg = 16.0;

  /// Extra large border radius - 24px (for bottom sheets)
  static const double radiusXl = 24.0;

  /// Full/circular border radius - 999px
  static const double radiusFull = 999.0;

  // ============================================
  // Border Radius Presets
  // ============================================

  /// Small border radius for all corners
  static const BorderRadius borderRadiusSm = BorderRadius.all(Radius.circular(radiusSm));

  /// Medium border radius for all corners
  static const BorderRadius borderRadiusMd = BorderRadius.all(Radius.circular(radiusMd));

  /// Large border radius for all corners
  static const BorderRadius borderRadiusLg = BorderRadius.all(Radius.circular(radiusLg));

  /// Extra large border radius for all corners
  static const BorderRadius borderRadiusXl = BorderRadius.all(Radius.circular(radiusXl));

  /// Top corners only - for bottom sheets
  static const BorderRadius borderRadiusTopLg = BorderRadius.only(
    topLeft: Radius.circular(radiusLg),
    topRight: Radius.circular(radiusLg),
  );

  /// Top corners only - extra large for modals
  static const BorderRadius borderRadiusTopXl = BorderRadius.only(
    topLeft: Radius.circular(radiusXl),
    topRight: Radius.circular(radiusXl),
  );

  // ============================================
  // Icon Sizes
  // ============================================

  /// Extra small icon - 16px
  static const double iconXs = 16.0;

  /// Small icon - 20px
  static const double iconSm = 20.0;

  /// Medium icon - 24px (default)
  static const double iconMd = 24.0;

  /// Large icon - 32px
  static const double iconLg = 32.0;

  /// Extra large icon - 48px
  static const double iconXl = 48.0;

  /// Extra extra large icon - 64px (for empty states)
  static const double iconXxl = 64.0;

  // ============================================
  // Avatar/Image Sizes
  // ============================================

  /// Small avatar - 32px
  static const double avatarSm = 32.0;

  /// Medium avatar - 48px
  static const double avatarMd = 48.0;

  /// Large avatar - 64px
  static const double avatarLg = 64.0;

  /// Extra large avatar - 96px
  static const double avatarXl = 96.0;

  // ============================================
  // Component Heights
  // ============================================

  /// Button height - 48px
  static const double buttonHeight = 48.0;

  /// Button height small - 40px
  static const double buttonHeightSmall = 40.0;

  /// Input field height - 52px
  static const double inputHeight = 52.0;

  /// App bar height - 56px
  static const double appBarHeight = 56.0;

  /// Bottom navigation height - 64px
  static const double bottomNavHeight = 64.0;

  /// Tab bar height - 48px
  static const double tabBarHeight = 48.0;

  /// List item height - 72px
  static const double listItemHeight = 72.0;

  /// List item height compact - 56px
  static const double listItemHeightCompact = 56.0;

  // ============================================
  // Divider/Border Widths
  // ============================================

  /// Thin border - 0.5px
  static const double borderThin = 0.5;

  /// Default border - 1px
  static const double borderDefault = 1.0;

  /// Thick border - 1.5px
  static const double borderThick = 1.5;

  /// Extra thick border - 2px
  static const double borderExtraThick = 2.0;

  // ============================================
  // Shadow Elevations
  // ============================================

  /// No elevation
  static const double elevationNone = 0.0;

  /// Low elevation - subtle shadow
  static const double elevationLow = 2.0;

  /// Medium elevation - card shadow
  static const double elevationMedium = 4.0;

  /// High elevation - modal shadow
  static const double elevationHigh = 8.0;

  /// Extra high elevation - floating action button
  static const double elevationExtraHigh = 12.0;

  // ============================================
  // Helper Methods
  // ============================================

  /// Returns a SizedBox with the specified height for vertical spacing.
  static SizedBox verticalSpace(double height) => SizedBox(height: height);

  /// Returns a SizedBox with the specified width for horizontal spacing.
  static SizedBox horizontalSpace(double width) => SizedBox(width: width);

  /// Vertical spacing presets
  static const SizedBox verticalSpaceXs = SizedBox(height: xs);
  static const SizedBox verticalSpaceSm = SizedBox(height: sm);
  static const SizedBox verticalSpaceMd = SizedBox(height: md);
  static const SizedBox verticalSpaceLg = SizedBox(height: lg);
  static const SizedBox verticalSpaceXl = SizedBox(height: xl);
  static const SizedBox verticalSpaceXxl = SizedBox(height: xxl);

  /// Horizontal spacing presets
  static const SizedBox horizontalSpaceXs = SizedBox(width: xs);
  static const SizedBox horizontalSpaceSm = SizedBox(width: sm);
  static const SizedBox horizontalSpaceMd = SizedBox(width: md);
  static const SizedBox horizontalSpaceLg = SizedBox(width: lg);
  static const SizedBox horizontalSpaceXl = SizedBox(width: xl);
}
