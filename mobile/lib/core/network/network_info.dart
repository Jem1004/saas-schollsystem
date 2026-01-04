/// Network connectivity information for the School Management mobile app.
///
/// Provides utilities to check network connectivity status.
library;

import 'dart:io';

/// Abstract interface for checking network connectivity.
abstract class NetworkInfo {
  /// Returns true if the device has an active network connection.
  Future<bool> get isConnected;
}

/// Implementation of [NetworkInfo] using dart:io.
///
/// Checks connectivity by attempting to lookup a known host.
class NetworkInfoImpl implements NetworkInfo {
  /// The host to use for connectivity checks.
  final String _lookupHost;

  /// Creates a [NetworkInfoImpl] with an optional custom lookup host.
  ///
  /// Defaults to 'google.com' for connectivity checks.
  NetworkInfoImpl({String lookupHost = 'google.com'}) : _lookupHost = lookupHost;

  @override
  Future<bool> get isConnected async {
    try {
      final result = await InternetAddress.lookup(_lookupHost);
      return result.isNotEmpty && result[0].rawAddress.isNotEmpty;
    } on SocketException catch (_) {
      return false;
    } catch (_) {
      return false;
    }
  }
}
