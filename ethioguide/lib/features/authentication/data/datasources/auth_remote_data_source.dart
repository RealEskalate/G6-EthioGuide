import 'package:dio/dio.dart';
import 'package:ethioguide/core/error/exception.dart';
import '../models/tokens_model.dart';
import '../models/user_model.dart';

abstract class AuthRemoteDataSource {
  Future<(UserModel, TokensModel)> login(String identifier, String password);
  Future<void> register({
    required String username,
    required String email,
    required String password,
    required String name,
    String? phone,
  });

  Future<void> forgotPassword(String email);
  Future<void> resetPassword({
    required String email,
    required String token,
    required String newPassword,
  });
  Future<(UserModel, TokensModel)> signInWithGoogle(String idToken);
}

class AuthRemoteDataSourceImpl implements AuthRemoteDataSource {
  final Dio dio;

  AuthRemoteDataSourceImpl({required this.dio});

  @override
  Future<(UserModel, TokensModel)> login(String identifier, String password) async {
    // --- MOCKED RESPONSE FOR DEVELOPMENT ---
    await Future.delayed(const Duration(seconds: 1)); 

    if (identifier == "test@test.com" && password == "password") {
      final mockJsonResponse = {
        "token": "mock_access_token_12345",
        "refreshToken": "mock_refresh_token_67890",
        "user": {
          "id": "123",
          "email": "test@test.com",
          "name": "Lidiya Test",
          "username": "lidiyatest"
        }
      };
      final user = UserModel.fromJson(mockJsonResponse['user'] as Map<String, dynamic>);
      final tokens = TokensModel.fromJson(mockJsonResponse);
      return (user, tokens);
    } else {
      // CORRECTED: Provide a message and status code when throwing the exception.
      throw ServerException(message: 'Invalid credentials', statusCode: 401);
    }
  }

  @override
  Future<void> register({
    required String username,
    required String email,
    required String password,
    required String name,
    String? phone,
  }) async {
    // --- MOCKED RESPONSE FOR DEVELOPMENT ---
    await Future.delayed(const Duration(seconds: 1));
    
    // You could add logic here to simulate errors, for example:
    if (email == "exists@test.com") {
      // CORRECTED: Provide a message for this specific error case.
      throw ServerException(message: 'Email already exists', statusCode: 409);
    }

    // Simulate a successful registration.
    return;
  }

  @override
  Future<void> forgotPassword(String email) async {
    // --- REAL API CALL (commented out) ---
    // final response = await dio.post('/auth/forgot', data: {'email': email});
    // if (response.statusCode != 204) { // As per your API doc
    //   throw ServerException(message: 'Failed to send reset link');
    // }

    // --- MOCKED RESPONSE ---
    await Future.delayed(const Duration(seconds: 1));
    if (email == 'notfound@test.com') {
      throw ServerException(message: 'Email not found', statusCode: 404);
    }
    // Simulate success
    return;
  }

  @override
  Future<void> resetPassword({
    required String email,
    required String token,
    required String newPassword,
  }) async {
    // --- REAL API CALL (commented out) ---
    // final response = await dio.post('/auth/reset', data: {
    //   'email': email,
    //   'reset token': token, // Key from your API doc
    //   'newPassword': newPassword,
    // });
    // if (response.statusCode != 204) {
    //   throw ServerException(message: 'Failed to reset password');
    // }

    // --- MOCKED RESPONSE ---
    await Future.delayed(const Duration(seconds: 1));
    if (token == 'invalid_token') {
      throw ServerException(message: 'Invalid or expired token', statusCode: 400);
    }
    // Simulate success
    return;
  }

   // ... inside class AuthRemoteDataSourceImpl ...

  @override
  Future<(UserModel, TokensModel)> signInWithGoogle(String authCode) async {
    // According to your backend team's new spec, we will call a new endpoint
    // and send the 'code' in the body. We will use Dio, not http.
    // final response = await dio.post(
    //   '/auth/google', // The endpoint your backend team provided
    //   data: {'code': authCode},
    // );
    //
    // if (response.statusCode == 200) {
    //   final user = UserModel.fromJson(response.data['user']);
    //   final tokens = TokensModel.fromJson(response.data); // Or however they return tokens
    //   return (user, tokens);
    // } else {
    //   throw ServerException(message: 'Backend Google authentication failed');
    // }

    // --- MOCKED RESPONSE FOR DEVELOPMENT ---
    await Future.delayed(const Duration(seconds: 1));
    final mockJsonResponse = {
        "token": "backend_access_token_google",
        "refreshToken": "backend_refresh_token_google",
        "user": {
          "id": "google_user_123",
          "email": "google@test.com",
          "name": "Google User",
          "username": "googleuser"
        }
      };
    final user = UserModel.fromJson(mockJsonResponse['user'] as Map<String, dynamic>);
    final tokens = TokensModel.fromJson(mockJsonResponse);
    return (user, tokens);
  }
}