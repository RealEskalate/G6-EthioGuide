import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/config/cache_key_names.dart';
import 'package:ethioguide/features/authentication/data/datasources/auth_local_data_source.dart';
import 'package:ethioguide/features/authentication/domain/repositories/auth_repositoryy.dart';
import 'package:ethioguide/core/error/exception.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/core/network/network_info.dart';
import 'package:ethioguide/features/authentication/data/datasources/auth_remote_data_source.dart';
import 'package:ethioguide/features/authentication/domain/entities/user.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class AuthRepositoryImpl implements AuthRepository {
  final AuthRemoteDataSource remoteDataSource;
  final FlutterSecureStorage secureStorage;
  final NetworkInfo networkInfo;
  
  final AuthLocalDataSource localDataSource;

  AuthRepositoryImpl({
    required this.remoteDataSource,
    required this.secureStorage,
    required this.localDataSource, 
    required this.networkInfo,
  });

  @override
  Future<Either<Failure, User>> login(String identifier, String password) async {
    if (await networkInfo.isConnected) {
      try {
        final (user, tokens) = await remoteDataSource.login(identifier, password);
        await saveTokens(accessToken: tokens.accessToken, refreshToken: tokens.refreshToken);
        return Right(user);
      } on ServerException catch (e) {
        // CORRECTED: Pass the exception's message to the failure object.
        return Left(ServerFailure(message: e.message));
      }
    } else {
      return Left(NetworkFailure());
    }
  }

  @override
  Future<Either<Failure, void>> register({
    required String username,
    required String email,
    required String password,
    required String name,
    String? phone,
  }) async {
    if (await networkInfo.isConnected) {
      try {
        await remoteDataSource.register(
          username: username,
          email: email,
          password: password,
          name: name,
          phone: phone,
        );
        return const Right(null);
      } on ServerException catch (e) {
        // CORRECTED: Pass the exception's message to the failure object.
        return Left(ServerFailure(message: e.message));
      }
    } else {
      return Left(NetworkFailure());
    }
  }

  @override
  Future<Either<Failure, void>> forgotPassword(String email) async {
    if (await networkInfo.isConnected) {
      try {
        await remoteDataSource.forgotPassword(email);
        return const Right(null);
      } on ServerException catch (e) {
        return Left(ServerFailure(message: e.message));
      }
    } else {
      return Left(NetworkFailure());
    }
  }

  @override
  Future<Either<Failure, void>> resetPassword({
    required String email,
    required String token,
    required String newPassword,
  }) async {
    if (await networkInfo.isConnected) {
      try {
        await remoteDataSource.resetPassword(
          email: email,
          token: token,
          newPassword: newPassword,
        );
        return const Right(null);
      } on ServerException catch (e) {
        return Left(ServerFailure(message: e.message));
      }
    } else {
      return Left(NetworkFailure());
    }
  }

// ... inside class AuthRepositoryImpl ...

  @override
  Future<Either<Failure, User>> signInWithGoogle() async {
    if (await networkInfo.isConnected) {
      try {
        // 1. THE CHANGE: Call the new method to get the server auth code.
        final googleAuthCode = await localDataSource.getGoogleServerAuthCode();
        
        // 2. Send this code to our backend. The rest of the logic is the same.
        final (user, tokens) = await remoteDataSource.signInWithGoogle(googleAuthCode);
        
        // 3. Save our app's tokens.
        await saveTokens(accessToken: tokens.accessToken, refreshToken: tokens.refreshToken);
        
        // 4. Return the user.
        return Right(user);
      } on CacheException catch (e) {
        return Left(CachedFailure(message: e.message));
      } on ServerException catch (e) {
        return Left(ServerFailure(message: e.message));
      }
    } else {
      return Left(NetworkFailure());
    }
  }

  // --- ALL YOUR EXISTING TOKEN MANAGEMENT METHODS ---
  // These do not need to change.
  @override
  Future<void> saveTokens({required String accessToken, required String refreshToken}) async {
    await secureStorage.write(key: CacheKeyNames.accessTokenKey, value: accessToken);
    await secureStorage.write(key: CacheKeyNames.refreshTokenKey, value: refreshToken);
  }

  @override
  Future<void> clearTokens() async {
    await secureStorage.delete(key: CacheKeyNames.accessTokenKey);
    await secureStorage.delete(key: CacheKeyNames.refreshTokenKey);
  }

  @override
  Future<String?> getAccessToken() async {
    return await secureStorage.read(key: CacheKeyNames.accessTokenKey);
  }

  @override
  Future<String?> getRefreshToken() async {
    return await secureStorage.read(key: CacheKeyNames.refreshTokenKey);
  }

  @override
  Future<bool> isAuthenticated() async {
    final token = await getAccessToken();
    return token != null && token.isNotEmpty;
  }
  
  @override
  Future<void> updateAccessToken(String token) async {
    await secureStorage.write(key: CacheKeyNames.accessTokenKey, value: token);
  }
}