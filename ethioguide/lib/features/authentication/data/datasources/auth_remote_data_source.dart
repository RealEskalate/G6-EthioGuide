import 'package:dio/dio.dart';
import 'package:ethioguide/core/error/exception.dart';
import '../models/tokens_model.dart';
import '../models/user_model.dart';
import 'package:ethioguide/core/config/end_points.dart';

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
     required String resetToken,
    required String newPassword,
  });
  Future<(UserModel, TokensModel)> signInWithGoogle(String idToken);
}

class AuthRemoteDataSourceImpl implements AuthRemoteDataSource {
  final Dio dio;

  AuthRemoteDataSourceImpl({required this.dio});

  @override
  Future<(UserModel, TokensModel)> login(
    String identifier,
    String password,
  ) async {
    try {
      // Make the real network request to the login endpoint.
      final response = await dio.post(
        EndPoints.loginEndPoint, // Using the constant from your EndPoints class
        data: {'identifier': identifier, 'password': password},
      );

      // Check for a successful response code (e.g., 200 OK).
      if (response.statusCode == 200 && response.data != null) {
        // Parse the JSON response into our data models.
        final user = UserModel.fromJson(response.data['user']);
        final tokens = TokensModel.fromJson(response.data);
        return (user, tokens);
      } else {
       
        throw ServerException(
          message:
              'Login failed with an unexpected status code: ${response.statusCode}',
          statusCode: response.statusCode,
        );
      }
    } on DioException catch (e) {
       print("----------- DIO REGISTER ERROR -----------");
  print("STATUS CODE: ${e.response?.statusCode}");
  print("RESPONSE DATA: ${e.response?.data}");
  print("--------------------------------------");
    
final errorMessage = e.response?.data?['message'] ?? e.response?.data?['error'] ?? 'An unknown server error occurred';
      throw ServerException(
        message: errorMessage,
        statusCode: e.response?.statusCode,
      );
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
    try {
      // Make the real network request to the register endpoint.
      final response = await dio.post(
        EndPoints.registerEndPoint, // Use the constant
        data: {
          'username': username,
          'email': email,
          'password': password,
          'name': name,
          'phone': phone, // Will be null if not provided
        },
      );

      // According to your API docs, a successful registration is 201 Created.
      if (response.statusCode != 201) {
        throw ServerException(
          message:
              'Registration failed with an unexpected status code: ${response.statusCode}',
          statusCode: response.statusCode,
        );
      }
      // If the status code is 201, the method completes successfully.
    } on DioException catch (e) {
      print("----------- DIO ERROR -----------");
      print("ERROR TYPE: ${e.type}");
      print("ERROR MESSAGE: ${e.message}");
      print("SERVER RESPONSE: ${e.response}");
      print("---------------------------------");
      final errorMessage =
          e.response?.data?['message'] ?? 'An unknown server error occurred';
      throw ServerException(
        message: errorMessage,
        statusCode: e.response?.statusCode,
      );
    }
  }



  @override
  Future<void> forgotPassword(String email) async {
    // --- REAL API CALL (commented out) ---
    // final response = await dio.post('/auth/forgot', data: {'email': email});
    // if (response.statusCode != 204) { // As per your API doc
    //   throw ServerException(message: 'Failed to send reset link');
    // }
    try {
      final response = await dio.post(
        EndPoints.forgotPassword,
        data: {'email': email},
      );
      // The API doc says 200 is success for this endpoint.
      if (response.statusCode != 200) {
        throw ServerException(message: 'Failed to send reset email.', statusCode: response.statusCode);
      }
    } on DioException catch (e) {
      final errorMessage = e.response?.data?['message'] ?? e.response?.data?['error'] ?? 'An unknown server error occurred';
      throw ServerException(message: errorMessage, statusCode: e.response?.statusCode);
    }
  }

  
  @override
  Future<void> resetPassword({
    required String resetToken,
    required String newPassword,
  }) async {
   
    try {
      // IMPORTANT: Based on your API docs, the /auth/reset endpoint only takes 'resetToken' and 'new_password'.
      // Please confirm with your backend team if 'email' is also needed.
      final response = await dio.post(
        EndPoints.resetPassword,
        data: {
          'resetToken': resetToken,
          'new_password': newPassword,
        },
      );
      
      if (response.statusCode != 200) {
        throw ServerException(message: 'Failed to reset password.', statusCode: response.statusCode);
      }
    } on DioException catch (e) {
      final errorMessage = e.response?.data?['message'] ?? e.response?.data?['error'] ?? 'An unknown server error occurred';
      throw ServerException(message: errorMessage, statusCode: e.response?.statusCode);
    }
  }
  

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
    "username": "googleuser",
  },
};
final user = UserModel.fromJson(
  mockJsonResponse['user'] as Map<String, dynamic>,
);
final tokens = TokensModel.fromJson(mockJsonResponse);
return (user, tokens);

    
  }

}