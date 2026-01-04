/// API configuration for the Flutter mobile app.
///
/// This file provides centralized API configuration and URL building utilities
/// for communicating with the School Management backend API.
///
/// For endpoint definitions, see [ApiEndpoints] in `api_endpoints.dart`.
library;

import 'app_config.dart';

// Re-export ApiEndpoints for convenience
export '../constants/api_endpoints.dart';

/// API configuration class that provides URL building utilities and API settings.
///
/// Usage:
/// ```dart
/// import 'package:mobile/core/config/api_config.dart';
///
/// // Get full URL for an endpoint
/// final loginUrl = ApiConfig.getUrl(ApiEndpoints.authLogin);
///
/// // Get URL with path parameters
/// final dashboardUrl = ApiConfig.getUrlWithParams(
///   ApiEndpoints.parentChildDashboard,
///   {'id': '123'},
/// );
///
/// // Access API settings
/// final timeout = ApiConfig.timeout;
/// ```
class ApiConfig {
  ApiConfig._();

  /// Gets the base URL from AppConfig.
  static String get baseUrl => AppConfig.instance.apiBaseUrl;

  /// Gets the API timeout in milliseconds from AppConfig.
  static int get timeout => AppConfig.instance.apiTimeout;

  /// Constructs a full URL for the given endpoint path.
  ///
  /// [endpoint] - The endpoint path (e.g., 'auth/login')
  /// Returns the full URL (e.g., 'http://localhost:8080/api/v1/auth/login')
  static String getUrl(String endpoint) {
    // Remove leading slash if present to avoid double slashes
    final cleanEndpoint = endpoint.startsWith('/') ? endpoint.substring(1) : endpoint;
    return '$baseUrl/$cleanEndpoint';
  }

  /// Constructs a URL with path parameters replaced.
  ///
  /// [endpoint] - The endpoint path with placeholders (e.g., 'parent/children/{id}/dashboard')
  /// [params] - Map of parameter names to values (e.g., {'id': '123'})
  /// Returns the full URL with parameters replaced
  ///
  /// Example:
  /// ```dart
  /// final url = ApiConfig.getUrlWithParams(
  ///   ApiEndpoints.parentChildDashboard,
  ///   {'id': '123'},
  /// );
  /// // Returns: 'http://localhost:8080/api/v1/parent/children/123/dashboard'
  /// ```
  static String getUrlWithParams(String endpoint, Map<String, String> params) {
    var url = endpoint;
    params.forEach((key, value) {
      url = url.replaceAll('{$key}', value);
    });
    return getUrl(url);
  }

  /// Constructs a URL with query parameters.
  ///
  /// [endpoint] - The endpoint path
  /// [queryParams] - Map of query parameter names to values
  /// Returns the full URL with query string
  ///
  /// Example:
  /// ```dart
  /// final url = ApiConfig.getUrlWithQuery(
  ///   ApiEndpoints.parentChildAttendance.replaceAll('{id}', '123'),
  ///   {'page': 1, 'limit': 20, 'start_date': '2024-01-01'},
  /// );
  /// ```
  static String getUrlWithQuery(String endpoint, Map<String, dynamic> queryParams) {
    final baseEndpointUrl = getUrl(endpoint);
    if (queryParams.isEmpty) {
      return baseEndpointUrl;
    }

    final queryString = queryParams.entries
        .where((e) => e.value != null)
        .map((e) => '${Uri.encodeComponent(e.key)}=${Uri.encodeComponent(e.value.toString())}')
        .join('&');

    return '$baseEndpointUrl?$queryString';
  }

  /// Constructs a URL with both path parameters and query parameters.
  ///
  /// [endpoint] - The endpoint path with placeholders
  /// [pathParams] - Map of path parameter names to values
  /// [queryParams] - Map of query parameter names to values
  /// Returns the full URL with path params replaced and query string appended
  ///
  /// Example:
  /// ```dart
  /// final url = ApiConfig.getFullUrl(
  ///   ApiEndpoints.parentChildAttendance,
  ///   pathParams: {'id': '123'},
  ///   queryParams: {'page': 1, 'limit': 20},
  /// );
  /// ```
  static String getFullUrl(
    String endpoint, {
    Map<String, String>? pathParams,
    Map<String, dynamic>? queryParams,
  }) {
    var url = endpoint;
    
    // Replace path parameters
    if (pathParams != null && pathParams.isNotEmpty) {
      pathParams.forEach((key, value) {
        url = url.replaceAll('{$key}', value);
      });
    }
    
    // Add query parameters
    if (queryParams != null && queryParams.isNotEmpty) {
      return getUrlWithQuery(url, queryParams);
    }
    
    return getUrl(url);
  }
}
