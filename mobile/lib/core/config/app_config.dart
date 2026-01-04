/// Environment configuration for the Flutter mobile app.
///
/// This file provides environment-specific configuration values
/// for different deployment environments (development, staging, production).
library;

/// Enum representing the available application environments.
enum Environment {
  development,
  staging,
  production,
}

/// Application configuration class that holds all environment-specific values.
///
/// Usage:
/// ```dart
/// // Initialize in main.dart before runApp()
/// AppConfig.initialize(Environment.development);
///
/// // Access configuration values
/// final baseUrl = AppConfig.instance.apiBaseUrl;
/// ```
class AppConfig {
  /// Private constructor to prevent direct instantiation.
  AppConfig._({
    required this.environment,
    required this.apiBaseUrl,
    required this.apiTimeout,
    required this.enableLogging,
    required this.appName,
    required this.appVersion,
  });

  /// Singleton instance of AppConfig.
  static AppConfig? _instance;

  /// Gets the current AppConfig instance.
  ///
  /// Throws [StateError] if [initialize] has not been called.
  static AppConfig get instance {
    if (_instance == null) {
      throw StateError(
        'AppConfig has not been initialized. '
        'Call AppConfig.initialize() before accessing instance.',
      );
    }
    return _instance!;
  }

  /// Checks if AppConfig has been initialized.
  static bool get isInitialized => _instance != null;

  /// The current environment.
  final Environment environment;

  /// Base URL for API requests.
  final String apiBaseUrl;

  /// API request timeout in milliseconds.
  final int apiTimeout;

  /// Whether to enable debug logging.
  final bool enableLogging;

  /// Application name.
  final String appName;

  /// Application version.
  final String appVersion;

  /// Initializes the AppConfig with the specified environment.
  ///
  /// This should be called once at app startup, typically in main.dart
  /// before calling runApp().
  ///
  /// [environment] - The target environment to configure for.
  /// [overrideBaseUrl] - Optional override for the API base URL.
  static void initialize(
    Environment environment, {
    String? overrideBaseUrl,
  }) {
    _instance = _createConfig(environment, overrideBaseUrl: overrideBaseUrl);
  }

  /// Creates a configuration for the specified environment.
  static AppConfig _createConfig(
    Environment environment, {
    String? overrideBaseUrl,
  }) {
    switch (environment) {
      case Environment.development:
        return AppConfig._(
          environment: environment,
          apiBaseUrl: overrideBaseUrl ?? 'http://localhost:8080/api/v1',
          apiTimeout: 30000, // 30 seconds
          enableLogging: true,
          appName: 'School Management (Dev)',
          appVersion: '1.0.0',
        );

      case Environment.staging:
        return AppConfig._(
          environment: environment,
          apiBaseUrl: overrideBaseUrl ?? 'https://staging-api.schoolmanagement.com/api/v1',
          apiTimeout: 30000, // 30 seconds
          enableLogging: true,
          appName: 'School Management (Staging)',
          appVersion: '1.0.0',
        );

      case Environment.production:
        return AppConfig._(
          environment: environment,
          apiBaseUrl: overrideBaseUrl ?? 'https://api.schoolmanagement.com/api/v1',
          apiTimeout: 15000, // 15 seconds
          enableLogging: false,
          appName: 'School Management',
          appVersion: '1.0.0',
        );
    }
  }

  /// Resets the configuration (useful for testing).
  static void reset() {
    _instance = null;
  }

  /// Whether the app is running in development mode.
  bool get isDevelopment => environment == Environment.development;

  /// Whether the app is running in staging mode.
  bool get isStaging => environment == Environment.staging;

  /// Whether the app is running in production mode.
  bool get isProduction => environment == Environment.production;

  /// Whether the app is running in a debug environment (dev or staging).
  bool get isDebugEnvironment => isDevelopment || isStaging;

  @override
  String toString() {
    return 'AppConfig('
        'environment: $environment, '
        'apiBaseUrl: $apiBaseUrl, '
        'apiTimeout: $apiTimeout, '
        'enableLogging: $enableLogging, '
        'appName: $appName, '
        'appVersion: $appVersion'
        ')';
  }
}
