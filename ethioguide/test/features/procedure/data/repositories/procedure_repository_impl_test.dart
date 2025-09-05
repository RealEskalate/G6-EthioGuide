import 'package:dartz/dartz.dart';
import 'package:dio/dio.dart';
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

import '../../domain/usecases/get_my_procedure_test.mocks.dart';
import 'procedure_repository_impl_test.mocks.dart';
import 'workspace_procedure_repository_impl_test.mocks.dart' hide MockNetworkInfo;




@GenerateMocks([NetworkInfo, ProcedureRemoteDataSource])
void main() {
  group('ProcedureRepositoryImpl.getProcedures', () {
    late MockProcedureRemoteDataSource remote;
    late MockNetworkInfo network;
    late ProcedureRepositoryImpl repository;

    setUp(() {
      remote = MockProcedureRemoteDataSource();
      network = MockNetworkInfo();
      repository = ProcedureRepositoryImpl(remoteDataSource: remote, networkInfo: network);
      when(network.isConnected).thenAnswer((_) async => true);
    });

    test('returns NetworkFailure when offline', () async {
      when(network.isConnected).thenAnswer((_) async => false);
      final result = await repository.getProcedures();

      expect(result, const Left(NetworkFailure()));
    });

    test('returns Right(list) when remote succeeds', () async {
      when(remote.getProcedures()).thenAnswer((_) async => const [
            ProcedureModel(
              id: '1',
              title: 'Passport',
              category: 'Travel',
              duration: '2 weeks',
              cost: '1200 ETB',
              icon: 'passport',
              isQuickAccess: true,
            ),
          ]);

      final result = await repository.getProcedures();

      expect(result.isRight(), true);
      result.fold(
        (_) => fail('Expected Right'),
        (value) => expect(value, isA<List<Procedure>>()),
      );
      verify(remote.getProcedures());
      verify(network.isConnected);
    });

    test('maps ServerException to ServerFailure', () async {
      when(remote.getProcedures()).thenThrow(ServerException(message: 'boom', statusCode: 500));

      final result = await repository.getProcedures();

      expect(result.isLeft(), true);
      result.fold(
        (failure) {
          expect(failure, isA<ServerFailure>());
          expect((failure as ServerFailure).message, 'boom');
        },
        (_) => fail('Expected Left'),
      );
    });

    test('maps any other error to generic ServerFailure', () async {
      when(remote.getProcedures()).thenThrow(DioException(requestOptions: RequestOptions(path: '/')));

      final result = await repository.getProcedures();

      expect(result.isLeft(), true);
      result.fold(
        (failure) => expect(failure, const ServerFailure()),
        (_) => fail('Expected Left'),
      );
    });
  });
}


