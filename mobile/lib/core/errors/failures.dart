/// Failure classes for the School Management mobile app.
///
/// Failures are used in the domain layer to represent errors in a type-safe way.
/// They are returned as the Left side of an Either<Failure, T> from repositories.
library;

import 'package:equatable/equatable.dart';

import '../constants/api_endpoints.dart';

/// Base failure class for all app failures.
///
/// All failures extend this class and provide a user-friendly message.
abstract class Failure extends Equatable {
  final String message;
  final String? code;

  const Failure({
    required this.message,
    this.code,
  });

  @override
  List<Object?> get props => [message, code];
}

/// Failure representing a server error.
class ServerFailure extends Failure {
  final int? statusCode;

  const ServerFailure({
    required super.message,
    super.code,
    this.statusCode,
  });

  @override
  List<Object?> get props => [message, code, statusCode];
}

/// Failure representing a network connectivity issue.
class NetworkFailure extends Failure {
  const NetworkFailure({
    super.message = 'Tidak ada koneksi internet',
    super.code = 'NETWORK_ERROR',
  });
}

/// Failure representing a request timeout.
class TimeoutFailure extends Failure {
  const TimeoutFailure({
    super.message = 'Koneksi timeout, coba lagi',
    super.code = 'TIMEOUT',
  });
}

/// Failure representing an authentication error.
class AuthFailure extends Failure {
  const AuthFailure({
    required super.message,
    super.code,
  });
}

/// Failure representing a cache/local storage error.
class CacheFailure extends Failure {
  const CacheFailure({
    super.message = 'Gagal mengambil data tersimpan',
    super.code = 'CACHE_ERROR',
  });
}

/// Failure representing a validation error.
class ValidationFailure extends Failure {
  final Map<String, List<String>>? fieldErrors;

  const ValidationFailure({
    required super.message,
    super.code = 'VALIDATION_ERROR',
    this.fieldErrors,
  });

  @override
  List<Object?> get props => [message, code, fieldErrors];
}

/// Failure representing a not found error.
class NotFoundFailure extends Failure {
  const NotFoundFailure({
    required super.message,
    super.code = 'NOT_FOUND',
  });
}

/// Failure representing an unauthorized access error.
class UnauthorizedFailure extends Failure {
  const UnauthorizedFailure({
    super.message = 'Anda tidak memiliki akses ke data ini',
    super.code = 'UNAUTHORIZED',
  });
}

/// Failure representing a forbidden access error.
class ForbiddenFailure extends Failure {
  const ForbiddenFailure({
    super.message = 'Akses ditolak',
    super.code = 'FORBIDDEN',
  });
}

/// Failure representing a session expiry.
class SessionExpiredFailure extends Failure {
  const SessionExpiredFailure({
    super.message = 'Sesi telah berakhir, silakan login kembali',
    super.code = 'SESSION_EXPIRED',
  });
}

/// Utility class for mapping API error codes to user-friendly messages.
class FailureMessageMapper {
  FailureMessageMapper._();

  /// Maps an API error code to a user-friendly Indonesian message.
  static String getMessage(String? code) {
    switch (code) {
      // Authentication errors
      case ApiErrorCodes.invalidCredentials:
        return 'Username atau password salah';
      case ApiErrorCodes.accountInactive:
        return 'Akun tidak aktif';
      case ApiErrorCodes.schoolInactive:
        return 'Sekolah tidak aktif';
      case ApiErrorCodes.tokenExpired:
        return 'Sesi telah berakhir, silakan login kembali';
      case ApiErrorCodes.tokenInvalid:
        return 'Token tidak valid';
      case ApiErrorCodes.passwordMismatch:
        return 'Password lama salah';
      case ApiErrorCodes.samePassword:
        return 'Password baru harus berbeda dari password lama';

      // Not found errors
      case ApiErrorCodes.studentNotFound:
        return 'Siswa tidak ditemukan';
      case ApiErrorCodes.parentNotFound:
        return 'Data orang tua tidak ditemukan';
      case ApiErrorCodes.notificationNotFound:
        return 'Notifikasi tidak ditemukan';

      // Authorization errors
      case ApiErrorCodes.notLinked:
        return 'Anda tidak memiliki akses ke data siswa ini';
      case ApiErrorCodes.forbidden:
        return 'Akses ditolak';

      // Network errors
      case ApiErrorCodes.networkError:
        return 'Tidak dapat terhubung ke server';
      case ApiErrorCodes.timeout:
        return 'Koneksi timeout, coba lagi';

      // Validation errors
      case ApiErrorCodes.validationError:
      case ApiErrorCodes.invalidInput:
        return 'Data yang dimasukkan tidak valid';

      default:
        return 'Terjadi kesalahan, coba lagi nanti';
    }
  }
}
