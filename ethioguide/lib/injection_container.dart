import 'package:dio/dio.dart';

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
import 'package:ethioguide/features/workspace_discussion/data/datasources/workspace_discussion_remote_data_source.dart';
import 'package:ethioguide/features/workspace_discussion/data/repositories/workspace_discussion_repository_impl.dart';
import 'package:ethioguide/features/workspace_discussion/domain/repositories/workspace_discussion_repository.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/add_comment.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/create_discussion.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/get_comments.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/get_community_stats.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/get_discussions.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/like_comment.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/like_discussion.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/report_comment.dart';
import 'package:ethioguide/features/workspace_discussion/presentation/bloc/workspace_discussion_bloc.dart';
// REMOVED: No longer need to import google_sign_in here.

import 'core/network/interceptors/auth_interceptor.dart';

import 'package:ethioguide/core/config/end_points.dart';
import 'package:ethioguide/core/data/repositories/auth_repository_impl.dart';
import 'package:ethioguide/core/domain/repositories/auth_repository.dart';
import 'package:ethioguide/core/network/interceptors/auth_interceptor.dart';
import 'package:ethioguide/core/network/network_info.dart';
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
import 'package:internet_connection_checker/internet_connection_checker.dart';

final sl = GetIt.instance;

Future<void> init() async {
  //* Features - Authentication
  sl.registerFactory(
    () => AuthBloc(
      loginUser: sl(),
      registerUser: sl(),
      forgotPassword: sl(),
      resetPassword: sl(),
      signInWithGoogle: sl(),
    ),
  );
  sl.registerLazySingleton(() => LoginUser(sl()));
  sl.registerLazySingleton(() => RegisterUser(sl()));
  sl.registerLazySingleton(() => ForgotPassword(sl()));
  sl.registerLazySingleton(() => ResetPassword(sl()));
  sl.registerLazySingleton(() => SignInWithGoogle(sl()));
  sl.registerLazySingleton<AuthRepository>(
    () => AuthRepositoryImpl(
      remoteDataSource: sl(),
      localDataSource: sl(),
      secureStorage: sl(),
      networkInfo: sl(),
    ),
  );
  sl.registerLazySingleton<AuthRemoteDataSource>(
    () => AuthRemoteDataSourceImpl(dio: sl()),
  );

  // THE FIX: The AuthLocalDataSourceImpl now has an empty constructor.
  sl.registerLazySingleton<AuthLocalDataSource>(
    () => AuthLocalDataSourceImpl(),
  );

  //* Features - Ai chat
  // Bloc
  sl.registerFactory<AiBloc>(
    () => AiBloc(
      sendQueryUseCase: sl(),
      getHistoryUseCase: sl(),
      translateContentUseCase: sl(),
    ),
  );
  sl.registerFactory<WorkspaceDiscussionBloc>(
    () => WorkspaceDiscussionBloc(
      getCommunityStats: sl(),
      getDiscussions: sl(),
      createDiscussion: sl(),
      addComment: sl(),
      getComments: sl(),
      likeDiscussion: sl(),
      likeComment: sl(),
      reportComment: sl(),
      reportDiscussion: sl(),
    ),
  );

  // Usecase
  sl.registerLazySingleton<SendQuery>(() => SendQuery(repository: sl()));
  sl.registerLazySingleton<GetHistory>(() => GetHistory(repository: sl()));
  sl.registerLazySingleton<TranslateContent>(
    () => TranslateContent(repository: sl()),
  );
  sl.registerLazySingleton<AddComment>(() => AddComment(sl()));
  sl.registerLazySingleton<GetComments>(() => GetComments(sl()));
  sl.registerLazySingleton<CreateDiscussion>(() => CreateDiscussion(sl()));
  sl.registerLazySingleton<GetCommunityStats>(() => GetCommunityStats(sl()));
  sl.registerLazySingleton<LikeDiscussion>(() => LikeDiscussion(sl()));
  sl.registerLazySingleton<GetCommunityStats>(() => GetCommunityStats(sl()));
  sl.registerLazySingleton<ReportComment>(() => ReportComment(sl()));

  sl.registerLazySingleton<GetDiscussions>(() => GetDiscussions(sl()));
  sl.registerLazySingleton<LikeComment>(() => LikeComment(sl()));

  // Repositories
  sl.registerLazySingleton<AiRepository>(
    () => AiRepositoryImpl(
      remoteDatasource: sl(),
      localDatasource: sl(),
      networkInfo: sl(),
    ),
  );
  sl.registerLazySingleton<WorkspaceDiscussionRepository>(
    () => WorkspaceDiscussionRepositoryImpl(sl()),
  );

  // Datasources
  sl.registerLazySingleton<AiRemoteDatasource>(
    () => AiRemoteDataSourceImpl(dio: sl(), networkInfo: sl()),
  );
  sl.registerLazySingleton<AiLocalDatasource>(
    () => AiLocalDataSourceImpl(secureStorage: sl()),
  );
  sl.registerLazySingleton<WorkspaceDiscussionRemoteDataSource>(
    () => WorkspaceDiscussionRemoteDataSourceImpl(dio: sl()),
  );

  //! Core
  sl.registerLazySingleton<NetworkInfo>(() => NetworkInfoImpl(sl()));

  //* Data Layer (Repositories)
  sl.registerLazySingleton<CoreAuthRepository>(
    () => CoreAuthRepositoryImpl(secureStorage: sl()),
  );

  //* Network (Dio & Interceptors)
  // sl.registerLazySingleton(() => AuthInterceptor());

  //! External
  sl.registerLazySingleton(() => const FlutterSecureStorage());
  sl.registerLazySingleton<Dio>(() {
    final dio = Dio(
      BaseOptions(
        baseUrl: EndPoints.baseUrl,
        headers: {'X-Client-Type': 'mobile'},
      ),
    );
    dio.interceptors.add(AuthInterceptor(sl<CoreAuthRepository>(), dio));
    return dio;
  });
  sl.registerLazySingleton(() => InternetConnectionChecker.createInstance());
}
