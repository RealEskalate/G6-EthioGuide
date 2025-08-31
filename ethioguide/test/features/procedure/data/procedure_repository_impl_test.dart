import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/exception.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/core/network/network_info.dart';
import 'package:ethioguide/features/procedure/data/datasources/procedure_remote_data_source.dart';
import 'package:ethioguide/features/procedure/data/models/procedure_model.dart';
import 'package:ethioguide/features/procedure/data/repositories/procedure_repository_impl.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

import 'procedure_repository_impl_test.mocks.dart';

@GenerateMocks([ProcedureRemoteDataSource, NetworkInfo])
void main() {
  late ProcedureRepositoryImpl repository;
  late MockProcedureRemoteDataSource mockRemote;
  late MockNetworkInfo mockNetworkInfo;

  setUp(() {
    mockRemote = MockProcedureRemoteDataSource();
    mockNetworkInfo = MockNetworkInfo();
    repository = ProcedureRepositoryImpl(remoteDataSource: mockRemote, networkInfo: mockNetworkInfo);
  });

  final tModels = [
    const ProcedureModel(
      id: '1',
      title: 'Passport Renewal',
      category: 'Travel',
      duration: '2-3 weeks',
      cost: '1,200 ETB',
      icon: 'passport',
      isQuickAccess: true,
    ),
  ];

  test('should return NetworkFailure when offline', () async {
    when(mockNetworkInfo.isConnected).thenAnswer((_) async => false);

    final result = await repository.getProcedures();

    expect(result, const Left(NetworkFailure()));
    verifyZeroInteractions(mockRemote);
  });

  test('should return list when remote call succeeds', () async {
    when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
    when(mockRemote.getProcedures()).thenAnswer((_) async => tModels);

    final result = await repository.getProcedures();

    expect(result, isA<Right<Failure, List<Procedure>>>());
  });

  test('should return ServerFailure on exception', () async {
    when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
    when(mockRemote.getProcedures()).thenThrow(ServerException(message: 'error'));

    final result = await repository.getProcedures();

    expect(result, isA<Left<Failure, List<Procedure>>>());
  });
}


