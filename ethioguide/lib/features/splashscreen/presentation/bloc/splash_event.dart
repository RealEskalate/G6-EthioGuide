import 'package:equatable/equatable.dart';

abstract class SplashEvent extends Equatable {
  const SplashEvent();

  @override
  List<Object> get props => [];
}

/// Event to notify the BLoC to start the splash screen timer process.
class StartSplashTimer extends SplashEvent {}