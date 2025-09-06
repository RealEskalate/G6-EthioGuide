import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import '../repositories/profile_repository.dart';

class LogoutUser {
  final ProfileRepository repository;
  LogoutUser(this.repository);

  Future<Either<Failure, void>> call() async {
    return await repository.logout();
  }
}