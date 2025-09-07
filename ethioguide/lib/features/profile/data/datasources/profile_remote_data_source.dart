import 'package:dio/dio.dart';
import 'package:ethioguide/core/error/exception.dart';

import 'package:ethioguide/features/authentication/data/models/user_model.dart';

abstract class ProfileRemoteDataSource {
  Future<UserModel> getUserProfile();
}

class ProfileRemoteDataSourceImpl implements ProfileRemoteDataSource {
  final Dio dio;
  ProfileRemoteDataSourceImpl({required this.dio});

  @override
  Future<UserModel> getUserProfile() async {
    try {
      // Your AuthInterceptor will automatically add the required Bearer token to this request.
      final response = await dio.get('/auth/me');
      
      if (response.statusCode == 200 && response.data != null) {
        return UserModel.fromJson(response.data);
      } else {
        throw ServerException(message: 'Failed to get profile', statusCode: response.statusCode);
      }
    } on DioException catch (e) {
      final errorMessage = e.response?.data?['message'] ?? 'Failed to get profile.';
      throw ServerException(message: errorMessage, statusCode: e.response?.statusCode);
    }
  }
}