import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/domain/repositories/auth_repository.dart';
import 'package:ethioguide/core/error/exception.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/core/network/network_info.dart';
import 'package:ethioguide/features/authentication/domain/entities/user.dart';
import 'package:ethioguide/features/profile/domain/repositories/profile_repository.dart';
import '../datasources/profile_remote_data_source.dart';

class ProfileRepositoryImpl implements ProfileRepository {
  final ProfileRemoteDataSource remoteDataSource;
  final CoreAuthRepository coreAuthRepository; // For logging out
  final NetworkInfo networkInfo;

  ProfileRepositoryImpl({
    required this.remoteDataSource,
    required this.coreAuthRepository,
    required this.networkInfo,
  });

  @override
  Future<Either<Failure, User>> getUserProfile() async {
    if (await networkInfo.isConnected) {
      try {
        final userModel = await remoteDataSource.getUserProfile();
        return Right(userModel);
      } on ServerException catch (e) {
        return Left(ServerFailure(message: e.message));
      }
    } else {
      return Left(NetworkFailure());
    }
  }

  @override
  Future<Either<Failure, void>> logout() async {
    try {
      // Delegate the logout action to the central core repository.
      await coreAuthRepository.clearTokens();
      return const Right(null);
    } catch (e) {
      return Left(const CachedFailure(message: 'Failed to clear tokens.'));
    }
  }
}