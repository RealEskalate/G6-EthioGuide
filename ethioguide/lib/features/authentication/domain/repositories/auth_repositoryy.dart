import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart'; 
import 'package:ethioguide/features/authentication/domain/entities/user.dart';

abstract class AuthRepository {
  Future<Either<Failure, User>> login(String identifier, String password);
  Future<Either<Failure, void>> register({
    required String username,
    required String email,
    required String password,
    required String name,
    String? phone,
  });
  
  Future<Either<Failure, void>> forgotPassword(String email);
  Future<Either<Failure, void>> resetPassword({
    required String resetToken,
    required String newPassword,
  });
  Future<Either<Failure, User>> signInWithGoogle();
  Future<Either<Failure, User>> verifyAccount(String activationToken);

  // ... existing token methods ...
  Future<void> saveTokens({ required String accessToken, required String refreshToken });
  Future<String?> getAccessToken();
  Future<void> updateAccessToken(String token);
  Future<String?> getRefreshToken();
  Future<bool> isAuthenticated();
  Future<void> clearTokens();
}