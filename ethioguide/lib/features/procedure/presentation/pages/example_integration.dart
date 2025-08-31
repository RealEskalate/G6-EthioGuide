import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:get_it/get_it.dart';
import 'package:dio/dio.dart';
import 'package:internet_connection_checker/internet_connection_checker.dart';
import '../bloc/workspace_procedure_detail_bloc.dart';
import '../data/datasources/workspace_procedure_remote_data_source.dart';
import '../data/datasources/workspace_procedure_local_datasource.dart';
import '../data/repositories/workspace_procedure_repository_impl.dart';
import '../domain/usecases/get_procedure_detail.dart';
import '../domain/usecases/update_step_status.dart';
import '../domain/usecases/save_progress.dart';
import 'workspace_procedure_detail_page.dart';
import '../../../../core/network/network_info.dart';

/// Example of how to set up dependency injection and integrate the feature
class ExampleIntegration {
  static void setupDependencies() {
    final getIt = GetIt.instance;

    // Register Dio instance
    getIt.registerLazySingleton<Dio>(() => Dio());

    // Register network info
    getIt.registerLazySingleton<NetworkInfo>(() => NetworkInfoImpl(getIt<InternetConnectionChecker>()));

    // Register data sources
    getIt.registerLazySingleton<WorkspaceProcedureRemoteDataSource>(
      () => WorkspaceProcedureRemoteDataSourceImpl(dio: getIt<Dio>()),
    );
    
    getIt.registerLazySingleton<WorkspaceProcedureLocalDataSource>(
      () => WorkspaceProcedureLocalDataSourceImpl(),
    );

    // Register repository
    getIt.registerLazySingleton<WorkspaceProcedureRepositoryImpl>(
      () => WorkspaceProcedureRepositoryImpl(
        remoteDataSource: getIt<WorkspaceProcedureRemoteDataSource>(),
        localDataSource: getIt<WorkspaceProcedureLocalDataSource>(),
        networkInfo: getIt<NetworkInfo>(),
      ),
    );

    // Register use cases
    getIt.registerLazySingleton<GetProcedureDetail>(
      () => GetProcedureDetail(getIt<WorkspaceProcedureRepositoryImpl>()),
    );
    getIt.registerLazySingleton<UpdateStepStatus>(
      () => UpdateStepStatus(getIt<WorkspaceProcedureRepositoryImpl>()),
    );
    getIt.registerLazySingleton<SaveProgress>(
      () => SaveProgress(getIt<WorkspaceProcedureRepositoryImpl>()),
    );

    // Register bloc
    getIt.registerFactory<WorkspaceProcedureDetailBloc>(
      () => WorkspaceProcedureDetailBloc(
        getProcedureDetail: getIt<GetProcedureDetail>(),
        updateStepStatus: getIt<UpdateStepStatus>(),
        saveProgress: getIt<SaveProgress>(),
      ),
    );
  }

  /// Example of how to navigate to the workspace procedure detail page
  static void navigateToProcedureDetail(BuildContext context, String procedureId) {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => BlocProvider(
          create: (context) => GetIt.instance<WorkspaceProcedureDetailBloc>(),
          child: WorkspaceProcedureDetailPage(procedureId: procedureId),
        ),
      ),
    );
  }

  /// Example of how to use the feature in a list item
  static Widget buildProcedureListItem(BuildContext context, String procedureId, String title) {
    return ListTile(
      title: Text(title),
      subtitle: const Text('Tap to view details'),
      trailing: const Icon(Icons.arrow_forward_ios),
      onTap: () => navigateToProcedureDetail(context, procedureId),
    );
  }
}

/// Example usage in a main app or other page
class ExampleUsagePage extends StatelessWidget {
  const ExampleUsagePage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Example Usage'),
      ),
      body: ListView(
        children: [
          ExampleIntegration.buildProcedureListItem(
            context,
            '1',
            'Driver\'s License Renewal',
          ),
          ExampleIntegration.buildProcedureListItem(
            context,
            '2',
            'Passport Application',
          ),
          ExampleIntegration.buildProcedureListItem(
            context,
            '3',
            'Vehicle Registration',
          ),
        ],
      ),
    );
  }
}

/// Example of how to initialize the app with dependencies
class ExampleApp extends StatelessWidget {
  const ExampleApp({super.key});

  @override
  Widget build(BuildContext context) {
    // Setup dependencies when app starts
    ExampleIntegration.setupDependencies();

    return MaterialApp(
      title: 'EthioGuide Example',
      theme: ThemeData(
        primarySwatch: Colors.blue,
        useMaterial3: true,
      ),
      home: const ExampleUsagePage(),
    );
  }
}
