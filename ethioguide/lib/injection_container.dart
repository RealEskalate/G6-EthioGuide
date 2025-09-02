import 'package:dio/dio.dart';
import 'package:ethioguide/core/network/network_info.dart';
import 'package:ethioguide/features/authentication/data/datasources/auth_local_data_source.dart';
import 'package:ethioguide/features/authentication/data/datasources/auth_remote_data_source.dart';
import 'package:ethioguide/features/authentication/data/repositories/auth_repositoryy_impl.dart';
import 'package:ethioguide/features/authentication/domain/repositories/auth_repositoryy.dart';
import 'package:ethioguide/features/authentication/domain/usecases/forgot_password.dart';
import 'package:ethioguide/features/authentication/domain/usecases/login_user.dart';
import 'package:ethioguide/features/authentication/domain/usecases/register_user.dart';
import 'package:ethioguide/features/authentication/domain/usecases/reset_password.dart';
import 'package:ethioguide/features/authentication/domain/usecases/sign_in_with_google.dart';
import 'package:ethioguide/features/authentication/presentation/bloc/auth_bloc.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:get_it/get_it.dart';
// REMOVED: No longer need to import google_sign_in here.
import 'package:internet_connection_checker/internet_connection_checker.dart';

import 'core/network/interceptors/auth_interceptor.dart';

final sl = GetIt.instance;

Future<void> init() async {
  //* Features - Authentication
  sl.registerFactory(() => AuthBloc(loginUser: sl(), registerUser: sl(), forgotPassword: sl(), resetPassword: sl(), signInWithGoogle: sl()));
  sl.registerLazySingleton(() => LoginUser(sl()));
  sl.registerLazySingleton(() => RegisterUser(sl()));
  sl.registerLazySingleton(() => ForgotPassword(sl()));
  sl.registerLazySingleton(() => ResetPassword(sl()));
  sl.registerLazySingleton(() => SignInWithGoogle(sl()));
  sl.registerLazySingleton<AuthRepository>(() => AuthRepositoryImpl(remoteDataSource: sl(), localDataSource: sl(), secureStorage: sl(), networkInfo: sl()));
  sl.registerLazySingleton<AuthRemoteDataSource>(() => AuthRemoteDataSourceImpl(dio: sl()));
  
  // THE FIX: The AuthLocalDataSourceImpl now has an empty constructor.
  sl.registerLazySingleton<AuthLocalDataSource>(() => AuthLocalDataSourceImpl());

  //! Core
  sl.registerLazySingleton<NetworkInfo>(() => NetworkInfoImpl(sl()));
  sl.registerLazySingleton(() => AuthInterceptor());

  //! External
  sl.registerLazySingleton(() => const FlutterSecureStorage());
  sl.registerLazySingleton(() => InternetConnectionChecker.createInstance());
  // THE FIX: We completely REMOVED the registration for GoogleSignIn.
  // It is now an internal detail of the AuthLocalDataSourceImpl.

  // Dio registration
  sl.registerLazySingleton<Dio>(() {
    final dio = Dio(BaseOptions(baseUrl: 'https://mock.api.com'));
    // Interceptor logic will be added back when API is ready.
    return dio;
  });
}