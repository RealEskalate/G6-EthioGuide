import 'package:dartz/dartz.dart';
import 'package:equatable/equatable.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/authentication/domain/repositories/auth_repositoryy.dart';

class ResetPassword {
  final AuthRepository repository;
  ResetPassword(this.repository);

  Future<Either<Failure, void>> call(ResetPasswordParams params) async {
    // THE FIX: Pass the correct parameter name to the repository
    return await repository.resetPassword(
      resetToken: params.resetToken,
      newPassword: params.newPassword,
    );
  }
}

class ResetPasswordParams extends Equatable {
  // THE FIX: Change the property name from 'token' to 'resetToken'
  final String resetToken;
  final String newPassword;

  const ResetPasswordParams({
    required this.resetToken,
    required this.newPassword,
  });

  @override
  List<Object?> get props => [resetToken, newPassword];
}