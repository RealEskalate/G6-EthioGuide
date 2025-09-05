import 'package:dartz/dartz.dart';
import 'package:equatable/equatable.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/authentication/domain/repositories/auth_repositoryy.dart';

class ResetPassword {
  final AuthRepository repository;

  ResetPassword(this.repository);

  Future<Either<Failure, void>> call(ResetPasswordParams params) async {
    return await repository.resetPassword(
      email: params.email,
      token: params.token,
      newPassword: params.newPassword,
    );
  }
}

class ResetPasswordParams extends Equatable {
  final String email;
  final String token;
  final String newPassword;

  const ResetPasswordParams({
    required this.email,
    required this.token,
    required this.newPassword,
  });

  @override
  List<Object?> get props => [email, token, newPassword];
}