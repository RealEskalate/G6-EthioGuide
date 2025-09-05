import 'package:equatable/equatable.dart';
import 'package:ethioguide/features/authentication/domain/entities/user.dart';

enum ProfileStatus { initial, loading, success, failure, loggedOut }

class ProfileState extends Equatable {
  final ProfileStatus status;
  final User? user;
  final String errorMessage;

  const ProfileState({
    this.status = ProfileStatus.initial,
    this.user,
    this.errorMessage = '',
  });

  ProfileState copyWith({
    ProfileStatus? status,
    User? user,
    String? errorMessage,
  }) {
    return ProfileState(
      status: status ?? this.status,
      user: user ?? this.user,
      errorMessage: errorMessage ?? this.errorMessage,
    );
  }

  @override
  List<Object?> get props => [status, user, errorMessage];
}