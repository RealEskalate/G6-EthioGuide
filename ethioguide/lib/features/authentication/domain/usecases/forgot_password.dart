import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/authentication/domain/repositories/auth_repositoryy.dart';

class ForgotPassword {
  final AuthRepository repository;

  ForgotPassword(this.repository);

  // This use case takes the email as a direct parameter.
  Future<Either<Failure, void>> call(String email) async {
    return await repository.forgotPassword(email);
  }
}