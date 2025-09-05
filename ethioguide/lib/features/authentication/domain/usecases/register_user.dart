import 'package:dartz/dartz.dart';
import 'package:equatable/equatable.dart';
import 'package:ethioguide/features/authentication/domain/repositories/auth_repositoryy.dart';
import 'package:ethioguide/core/error/failures.dart';

class RegisterUser {
  final AuthRepository repository;

  RegisterUser(this.repository);

  Future<Either<Failure, void>> call(RegisterParams params) async {
    return await repository.register(
      username: params.username,
      email: params.email,
      password: params.password,
      name: params.name,
      phone: params.phone,
    );
  }
}

class RegisterParams extends Equatable {
  final String username;
  final String email;
  final String password;
  final String name;
  final String? phone;

  const RegisterParams({
    required this.username,
    required this.email,
    required this.password,
    required this.name,
    this.phone,
  });
  
  @override
  List<Object?> get props => [username, email, password, name, phone];
}