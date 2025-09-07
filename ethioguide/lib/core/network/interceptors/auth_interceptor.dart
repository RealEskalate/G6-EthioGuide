import 'package:dio/dio.dart';
import 'package:ethioguide/core/config/end_points.dart';
import 'package:ethioguide/core/domain/repositories/auth_repository.dart';
// import 'package:get_it/get_it.dart';

// TODO: check if we can remove this instance
// final getIt =
//     GetIt.instance; // Re-declare getIt here if not globally accessible

class AuthInterceptor extends Interceptor {
  final CoreAuthRepository _authRepository;
  final Dio _dio;

  AuthInterceptor(this._authRepository, this._dio);

  @override
  void onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    if (await _authRepository.isAuthenticated()) {

      final accessToken =
          // 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjhiN2YyZmE1Mzc2MjI3YTc4ZmYxZTJiIiwicm9sZSI6InVzZXIiLCJzdWJzY3JpcHRpb24iOiIiLCJpc3MiOiJldGhpby1ndWlkZS1hcGkiLCJleHAiOjE3NTczMTk2MzcsImlhdCI6MTc1NzEwMzYzNywianRpIjoiZDgyNTFkYjYtNDdkZC00YTNjLWFjYmItMTg0YjE2M2ZmNGFjIn0.S0GJap96N1jBiHIrof-_64HYcfWvXFf5LrgmqCYh7y8';
          await _authRepository.getAccessToken();
      options.headers['Authorization'] = 'Bearer $accessToken';

      print(
        "Request[${options.method}] => PATH: ${options.path}, token: ${options.headers['Authorization']}",
      );

      //  'Bearer $accessToken';

    }
     
    return handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    final isRetry = err.requestOptions.extra['retry'] == true;

    if (err.response?.statusCode == 401 && !isRetry) {
      final refreshToken = await _authRepository.getRefreshToken();
      if (refreshToken != null) {
        try {
          final refreshResponse = await _dio.post(
            EndPoints.refreshTokenEndPoint,
            data: {'refreshToken': refreshToken},
          );

          final newAccessToken = refreshResponse.data['accessToken'];
          if (newAccessToken != null) {
            await _authRepository.updateAccessToken(newAccessToken);

            final retryOptions = err.requestOptions;
            retryOptions.headers['Authorization'] = 'Bearer $newAccessToken';
            retryOptions.extra['retry'] = true;

            final retryResponse = await _dio.fetch(retryOptions);
            return handler.resolve(retryResponse);
          }
        } catch (refreshError) {
          await _authRepository.clearTokens();
          return handler.next(err);
        }
      }
    }
    return handler.next(err);
  }
}
