import 'package:dio/dio.dart';
import 'package:ethioguide/core/config/end_points.dart';
import 'package:ethioguide/core/data/repositories/auth_repository_impl.dart';
import 'package:ethioguide/core/domain/repositories/auth_repository.dart';
import 'package:ethioguide/core/network/interceptors/auth_interceptor.dart';
import 'package:ethioguide/features/AI%20chat/Domain/repository/ai_repository.dart';
import 'package:ethioguide/features/AI%20chat/Domain/usecases/get_history.dart';
import 'package:ethioguide/features/AI%20chat/Domain/usecases/send_query.dart';
import 'package:ethioguide/features/AI%20chat/Domain/usecases/translate_content.dart';
import 'package:ethioguide/features/AI%20chat/Presentation/bloc/ai_bloc.dart';
import 'package:ethioguide/features/AI%20chat/data/datasources/ai_local_datasource.dart';
import 'package:ethioguide/features/AI%20chat/data/datasources/ai_remote_datasource.dart';
import 'package:ethioguide/features/AI%20chat/data/repository/ai_repository_impl.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:get_it/get_it.dart';

final sl = GetIt.instance;

Future<void> init() async {
  //* Features -
  // Bloc

  // Usecase

  // Repositories

  // Datasources

  //* Features - Ai chat
  // Bloc
  sl.registerFactory<AiBloc>(
    () => AiBloc(
      sendQueryUseCase: sl(),
      getHistoryUseCase: sl(),
      translateContentUseCase: sl(),
    ),
  );

  // Usecase
  sl.registerLazySingleton<SendQuery>(() => SendQuery(repository: sl()));
  sl.registerLazySingleton<GetHistory>(() => GetHistory(repository: sl()));
  sl.registerLazySingleton<TranslateContent>(
    () => TranslateContent(repository: sl()),
  );

  // Repositories
  sl.registerLazySingleton<AiRepository>(
    () => AiRepositoryImpl(
      remoteDatasource: sl(),
      localDatasource: sl(),
      networkInfo: sl(),
    ),
  );

  // Datasources
  sl.registerLazySingleton<AiRemoteDatasource>(
    () => AiRemoteDataSourceImpl(dio: sl(), networkInfo: sl()),
  );
  sl.registerLazySingleton<AiLocalDatasource>(
    () => AiLocalDataSourceImpl(secureStorage: sl()),
  );

  //! Core
  //* Data Layer (Repositories)
  sl.registerLazySingleton<AuthRepository>(
    () => AuthRepositoryImpl(secureStorage: sl()),
  );

  //* Network (Dio & Interceptors)
  sl.registerLazySingleton(() => AuthInterceptor());

  //! External
  sl.registerLazySingleton(() => const FlutterSecureStorage());
  sl.registerLazySingleton<Dio>(() {
    final dio = Dio(
      BaseOptions(
        baseUrl: EndPoints.baseUrl,
        headers: {'X-Platform': 'mobile'},
      ),
    );
    dio.interceptors.add(sl<AuthInterceptor>());
    // You might add other interceptors here too
    return dio;
  });
}
