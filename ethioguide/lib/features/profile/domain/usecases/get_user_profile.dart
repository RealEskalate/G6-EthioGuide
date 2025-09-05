import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/authentication/domain/entities/user.dart';
import '../repositories/profile_repository.dart';

class GetUserProfile {
  final ProfileRepository repository;
  GetUserProfile(this.repository);

  Future<Either<Failure, User>> call() async {
    return await repository.getUserProfile();
  }
}