import 'package:internet_connection_checker/internet_connection_checker.dart';

// The abstract contract.
abstract class NetworkInfo {
  Future<bool> get isConnected;
}

// The concrete implementation for mobile.
// It requires an InternetConnectionChecker to be passed into its constructor.
class NetworkInfoImpl implements NetworkInfo {
  final InternetConnectionChecker connectionChecker;
  NetworkInfoImpl(this.connectionChecker);

  @override
  Future<bool> get isConnected => connectionChecker.hasConnection;
}