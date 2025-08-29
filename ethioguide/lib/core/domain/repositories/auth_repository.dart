abstract class AuthRepository {
  Future<void> saveTokens({
    required String accessToken,
    required String refreshToken,
  });

  Future<String?> getAccessToken();
  Future<void> updateAccessToken(String token); // update with new access token
  Future<String?> getRefreshToken();
  Future<bool> isAuthenticated();
  Future<void> clearTokens();
}
