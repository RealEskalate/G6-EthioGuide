import 'package:dio/dio.dart';

// A minimal interceptor that does nothing for now.
class AuthInterceptor extends Interceptor {
  @override
  void onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) {
    // We will add the token logic back here later when the API is ready.
    return handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    // We will add the refresh logic back here later.
    return handler.next(err);
  }
}