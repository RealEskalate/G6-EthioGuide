import 'package:dio/dio.dart';
import 'package:ethioguide/core/config/end_points.dart';
import 'package:ethioguide/core/domain/repositories/auth_repository.dart';
import 'package:get_it/get_it.dart';

// todo: check if we can remove this instance
final getIt =
    GetIt.instance; // Re-declare getIt here if not globally accessible

class AuthInterceptor extends Interceptor {
  final AuthRepository _authRepository = getIt<AuthRepository>();
  final Dio _dio = getIt<Dio>();

  @override
  void onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    if (await _authRepository.isAuthenticated()) {
      final accessToken = await _authRepository.getAccessToken();
      options.headers['Authorization'] = 'Bearer $accessToken';
    }
    return handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.response?.statusCode == 401) {
      // if token expired (401), try to refresh
      final refreshToken = await _authRepository.getRefreshToken();
      if (refreshToken != null) {
        try {
          final refreshResponse = await _dio.post(
            EndPoints.refreshTokenEndPoint,
            data: {'refreshToken': refreshToken},
          );

          final newAccessToken = refreshResponse.data['accessToken'];
          if (newAccessToken != null) {
            // save new accesss token
            await _authRepository.updateAccessToken(newAccessToken);

            // Update headers and retry the original request
            err.requestOptions.headers['Authorization'] =
                'Bearer $newAccessToken';
            final retryResponse = await _dio.fetch(err.requestOptions);
            return handler.resolve(retryResponse);
          }
        } catch (refreshError) {
          // Refresh failed -> logout user
          await _authRepository.clearTokens();
          return handler.next(err);
        }
      }
      return handler.next(err);
    }
    return handler.next(err);
  }
}
