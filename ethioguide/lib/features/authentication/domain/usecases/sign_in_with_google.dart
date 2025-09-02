import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/authentication/domain/entities/user.dart';
import 'package:ethioguide/features/authentication/domain/repositories/auth_repositoryy.dart';

// No parameters are needed, so we can create a simple UseCase class.
class SignInWithGoogle {
  final AuthRepository repository;

  SignInWithGoogle(this.repository);

  Future<Either<Failure, User>> call() async {
    return await repository.signInWithGoogle();
  }
}