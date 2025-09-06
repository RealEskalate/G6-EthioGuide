import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/authentication/domain/entities/user.dart';
import 'package:ethioguide/features/authentication/domain/repositories/auth_repositoryy.dart';

class VerifyAccount {
  final AuthRepository repository;
  VerifyAccount(this.repository);

  Future<Either<Failure, User>> call(String activationToken) async {
    return await repository.verifyAccount(activationToken);
  }
}