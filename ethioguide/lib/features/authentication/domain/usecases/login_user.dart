import 'package:dartz/dartz.dart';
import 'package:equatable/equatable.dart';
import 'package:ethioguide/features/authentication/domain/repositories/auth_repositoryy.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/authentication/domain/entities/user.dart';

// Usecase for user login.
class LoginUser {
  final AuthRepository repository;

  LoginUser(this.repository);

  Future<Either<Failure, User>> call(LoginParams params) async {
    return await repository.login(params.identifier, params.password);
  }
}

// Parameters required for the LoginUser usecase.
class LoginParams extends Equatable {
  final String identifier; // Can be username, email, or phone
  final String password;

  const LoginParams({required this.identifier, required this.password});

  @override
  List<Object?> get props => [identifier, password];
}