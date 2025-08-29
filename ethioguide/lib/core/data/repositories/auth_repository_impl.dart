import 'package:ethioguide/core/config/cache_key_names.dart';
import 'package:ethioguide/core/domain/repositories/auth_repository.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class AuthRepositoryImpl implements AuthRepository {
  final FlutterSecureStorage secureStorage;

  AuthRepositoryImpl({required this.secureStorage});

  @override
  Future<void> clearTokens() async {
    // Delete both tokens (logout use case)
    await secureStorage.delete(key: CacheKeyNames.accessTokenKey);
    await secureStorage.delete(key: CacheKeyNames.refreshTokenKey);
  }

  @override
  Future<String?> getAccessToken() async {
    // Read the stored access token
    return await secureStorage.read(key: CacheKeyNames.accessTokenKey);
  }

  @override
  Future<String?> getRefreshToken() async {
    // Read the stored refresh token
    return await secureStorage.read(key: CacheKeyNames.refreshTokenKey);
  }

  @override
  Future<bool> isAuthenticated() async {
    // If access token exists -> consider user authenticated
    final accessToken = await getAccessToken();
    return accessToken != null && accessToken.isNotEmpty;
  }

  @override
  Future<void> saveTokens({
    required String accessToken,
    required String refreshToken,
  }) async {
    // Save both tokens securely
    await secureStorage.write(
      key: CacheKeyNames.accessTokenKey,
      value: accessToken,
    );
    await secureStorage.write(
      key: CacheKeyNames.refreshTokenKey,
      value: refreshToken,
    );
  }

  @override
  Future<void> updateAccessToken(String token) async {
    // Overwrite only the access token
    await secureStorage.write(key: CacheKeyNames.accessTokenKey, value: token);
  }
}
