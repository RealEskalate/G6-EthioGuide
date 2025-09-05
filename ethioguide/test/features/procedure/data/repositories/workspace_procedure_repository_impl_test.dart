import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/core/network/network_info.dart';
import 'package:ethioguide/features/procedure/data/datasources/workspace_procedure_remote_data_source.dart';
import 'package:ethioguide/features/procedure/data/models/workspace_procedure_model.dart';
import 'package:ethioguide/features/procedure/data/models/workspace_summary_model.dart';
import 'package:ethioguide/features/procedure/data/repositories/workspace_procedure_repository_impl.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_detail.dart';
import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

import 'procedure_repository_impl_test.dart';
import 'workspace_procedure_repository_impl_test.mocks.dart';
// 👇 This tells Mockito to generate mocks for these classes
@GenerateMocks([NetworkInfo, WorkspaceProcedureRemoteDataSource])
void main() {
  late MockWorkspaceProcedureRemoteDataSource remote;
  late MockNetworkInfo net;
  late WorkspaceProcedureRepositoryImpl repo;

  setUp(() {
    remote = MockWorkspaceProcedureRemoteDataSource();
    net = MockNetworkInfo();
    repo = WorkspaceProcedureRepositoryImpl(
      remoteDataSource: remote,
      networkInfo: net,
    );
  });

  group('getProcedure', () {
    test('returns Right when online and success', () async {
      when(net.isConnected).thenAnswer((_) async => true);
      when(remote.getMyProcedures()).thenAnswer((_) async => <WorkspaceProcedureModel>[]);

      final r = await repo.getProcedure();

      expect(r, isA<Right<Failure, List<ProcedureDetail>>>());
    });

    test('returns Left when offline', () async {
      when(net.isConnected).thenAnswer((_) async => false);

      final r = await repo.getProcedure();

      expect(r, isA<Left<Failure, List<ProcedureDetail>>>());
    });
  });

  group('getWorkspaceSummary', () {
    test('returns Right on success', () async {
      when(net.isConnected).thenAnswer((_) async => true);
      when(remote.getWorkspaceSummary()).thenAnswer(
        (_) async => const WorkspaceSummaryModel(
          totalProcedures: 0,
          inProgress: 0,
          completed: 0,
          totalDocuments: 0,
        ),
      );

      final r = await repo.getWorkspaceSummary();

      expect(r, isA<Right<Failure, WorkspaceSummary>>());
    });
  });

  group('getProcedureDetail', () {
    test('returns Left when offline', () async {
      when(net.isConnected).thenAnswer((_) async => false);

      final r = await repo.getProcedureDetail('p1');

      expect(r, isA<Left<String, ProcedureDetail>>());
    });
  });

  group('updateStepStatus', () {
    test('returns Right on success', () async {
      when(net.isConnected).thenAnswer((_) async => true);
      when(remote.updateStepStatus(any, any, any)).thenAnswer((_) async => true);

      final r = await repo.updateStepStatus('p1', 's1', true);

      expect(r, isA<Right<String, bool>>());
    });
  });

  group('saveProgress', () {
    test('returns Right on success', () async {
      when(net.isConnected).thenAnswer((_) async => true);
      when(remote.saveProgress(any)).thenAnswer((_) async => true);

      final r = await repo.saveProgress('p1');

      expect(r, isA<Right<String, bool>>());
    });
  });
}
