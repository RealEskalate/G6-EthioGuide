import 'package:equatable/equatable.dart';

// super class of all the failures
class Failure extends Equatable {
  final String message;
  const Failure({this.message = 'Failure'}); // default value for failure message

  @override
  List<Object?> get props => [message];
}

class ServerFailure extends Failure {
  const ServerFailure({super.message = 'Server Failure'}); // default value for server failure message
}

class CachedFailure extends Failure {
  const CachedFailure({super.message = 'Cache Failure'}); // default value for server failure message
}

class NetworkFailure extends Failure {
  const NetworkFailure({super.message = 'Network Failure'}); // default value for server failure message
}
