import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
// THE FIX: Import the User entity from its location in the authentication feature.
import 'package:ethioguide/features/authentication/domain/entities/user.dart';

abstract class ProfileRepository {
  Future<Either<Failure, User>> getUserProfile();
  Future<Either<Failure, void>> logout();
}