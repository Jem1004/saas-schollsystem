/// Custom exceptions for the School Management mobile app.
///
/// These exceptions are thrown by data sources and caught by repositories
/// to be converted into [Failure] objects.
library;

/// Base exception class for all app exceptions.
abstract class AppException implements Exception {
  final String message;
  final String? code;
  final dynamic originalError;

  const AppException({
    required this.message,
    this.code,
    this.originalError,
  });

  @override
  String toString() => 'AppException: $message (code: $code)';
}

/// Exception thrown when a server request fails.
///
/// Contains the HTTP status code and error code from the backend.
class ServerException extends AppException {
  final int? statusCode;

  const ServerException({
    required super.message,
    super.code,
    this.statusCode,
    super.originalError,
  });

  @override
  String toString() =>
      'ServerException: $message (code: $code, statusCode: $statusCode)';
}

/// Exception thrown when there is no network connection.
class NetworkException extends AppException {
  const NetworkException({
    super.message = 'Tidak ada koneksi internet',
    super.code = 'NETWORK_ERROR',
    super.originalError,
  });
}

/// Exception thrown when a request times out.
class TimeoutException extends AppException {
  const TimeoutException({
    super.message = 'Koneksi timeout, coba lagi',
    super.code = 'TIMEOUT',
    super.originalError,
  });
}

/// Exception thrown when authentication fails.
class AuthException extends AppException {
  const AuthException({
    required super.message,
    super.code,
    super.originalError,
  });
}

/// Exception thrown when the access token has expired.
class TokenExpiredException extends AuthException {
  const TokenExpiredException({
    super.message = 'Sesi telah berakhir, silakan login kembali',
    super.code = 'AUTH_TOKEN_EXPIRED',
    super.originalError,
  });
}

/// Exception thrown when the refresh token is invalid.
class RefreshTokenException extends AuthException {
  const RefreshTokenException({
    super.message = 'Sesi tidak valid, silakan login kembali',
    super.code = 'AUTH_TOKEN_INVALID',
    super.originalError,
  });
}

/// Exception thrown when cached data is not found or corrupted.
class CacheException extends AppException {
  const CacheException({
    super.message = 'Gagal mengambil data tersimpan',
    super.code = 'CACHE_ERROR',
    super.originalError,
  });
}

/// Exception thrown when input validation fails.
class ValidationException extends AppException {
  final Map<String, List<String>>? fieldErrors;

  const ValidationException({
    required super.message,
    super.code = 'VALIDATION_ERROR',
    this.fieldErrors,
    super.originalError,
  });
}

/// Exception thrown when a requested resource is not found.
class NotFoundException extends AppException {
  const NotFoundException({
    required super.message,
    super.code = 'NOT_FOUND',
    super.originalError,
  });
}

/// Exception thrown when the user is not authorized to access a resource.
class UnauthorizedException extends AppException {
  const UnauthorizedException({
    super.message = 'Anda tidak memiliki akses ke data ini',
    super.code = 'UNAUTHORIZED',
    super.originalError,
  });
}

/// Exception thrown when the user is forbidden from accessing a resource.
class ForbiddenException extends AppException {
  const ForbiddenException({
    super.message = 'Akses ditolak',
    super.code = 'FORBIDDEN',
    super.originalError,
  });
}
