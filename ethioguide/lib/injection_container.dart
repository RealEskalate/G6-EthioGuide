import 'package:dio/dio.dart';
import 'package:ethioguide/core/data/repositories/auth_repository_impl.dart';
import 'package:ethioguide/core/domain/repositories/auth_repository.dart';
import 'package:ethioguide/core/network/interceptors/auth_interceptor.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:get_it/get_it.dart';

final sl = GetIt.instance;

Future<void> init() async {
  //* Features -
  // Bloc

  // Usecase

  // Repositories

  // Datasources

  //! Core
  //* Data Layer (Repositories)
  sl.registerLazySingleton<AuthRepository>(
    () => AuthRepositoryImpl(secureStorage: sl()),
  );

  //* Network (Dio & Interceptors)
  sl.registerLazySingleton(() => AuthInterceptor());

  //! External
  sl.registerLazySingleton(() => const FlutterSecureStorage());
  // todo: replace with actural base url
  sl.registerLazySingleton<Dio>(() {
    final dio = Dio(BaseOptions(baseUrl: 'YOUR_API_BASE_URL'));
    dio.interceptors.add(sl<AuthInterceptor>());
    // You might add other interceptors here too
    return dio;
  });
}
