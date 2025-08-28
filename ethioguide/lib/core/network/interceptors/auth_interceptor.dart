import 'package:dio/dio.dart';
import 'package:ethioguide/core/domain/repositories/auth_repository.dart';
import 'package:get_it/get_it.dart';

// todo: check if we can remove this instance
final getIt =
    GetIt.instance; // Re-declare getIt here if not globally accessible

class AuthInterceptor extends Interceptor {
  final AuthRepository _authRepository = getIt<AuthRepository>();

  @override
  void onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    if (await _authRepository.isAuthenticated()) {
      options.headers['Authorization'] =
          'Bearer ${_authRepository.getAccessToken()}';
    }
    super.onRequest(options, handler);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    if (err.response?.statusCode == 401 || err.response?.statusCode == 403) {
      // todo: Handle token refresh or redirect to login
    }
    super.onError(err, handler);
  }
}
