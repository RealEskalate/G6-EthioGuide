import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:dartz/dartz.dart';
import 'package:ethioguide/features/procedure/data/models/procedure_detail_model.dart';
import 'package:ethioguide/features/procedure/data/repositories/workspace_procedure_repository_impl.dart';
import 'package:ethioguide/features/procedure/data/datasources/workspace_procedure_remote_data_source.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_detail.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_step.dart';
import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';

import 'workspace_procedure_repository_impl_test.mocks.dart';

@GenerateMocks([WorkspaceProcedureRemoteDataSource])
void main() {
  late WorkspaceProcedureRepositoryImpl repository;
  late MockWorkspaceProcedureRemoteDataSource mockRemoteDataSource;

  setUp(() {
    mockRemoteDataSource = MockWorkspaceProcedureRemoteDataSource();
    repository = WorkspaceProcedureRepositoryImpl(mockRemoteDataSource);
  });

  group('getProcedureDetail', () {
    final testProcedureDetail = ProcedureDetailModel(
      id: '1',
      title: 'Test Procedure',
      organization: 'Test Org',
      status: ProcedureStatus.inProgress,
      progressPercentage: 40,
      documentsUploaded: 2,
      totalDocuments: 5,
      startDate: DateTime(2024, 1, 1),
      estimatedCompletion: DateTime(2024, 1, 5),
      completedDate: null,
      notes: 'Test notes',
      steps: [
        ProcedureStepModel(
          id: 'step1',
          title: 'Step 1',
          description: 'Test step 1',
          isCompleted: true,
          completionStatus: 'Completed',
          order: 1,
        ),
        ProcedureStepModel(
          id: 'step2',
          title: 'Step 2',
          description: 'Test step 2',
          isCompleted: false,
          completionStatus: null,
          order: 2,
        ),
      ],
      estimatedTime: '2-3 days',
      difficulty: 'Easy',
      officeType: 'Authority',
      quickTips: ['Tip 1', 'Tip 2'],
      requiredDocuments: ['Doc 1', 'Doc 2'],
    );

    test('should return ProcedureDetail when remote data source is successful', () async {
      // arrange
      when(mockRemoteDataSource.getProcedureDetail('1'))
          .thenAnswer((_) async => testProcedureDetail);

      // act
      final result = await repository.getProcedureDetail('1');

      // assert
      expect(result, equals(Right(testProcedureDetail)));
      verify(mockRemoteDataSource.getProcedureDetail('1'));
      verifyNoMoreInteractions(mockRemoteDataSource);
    });

    test('should return error message when remote data source throws exception', () async {
      // arrange
      when(mockRemoteDataSource.getProcedureDetail('1'))
          .thenThrow(Exception('Network error'));

      // act
      final result = await repository.getProcedureDetail('1');

      // assert
      expect(result, equals(Left('Exception: Network error')));
      verify(mockRemoteDataSource.getProcedureDetail('1'));
      verifyNoMoreInteractions(mockRemoteDataSource);
    });
  });

  group('updateStepStatus', () {
    test('should return true when remote data source is successful', () async {
      // arrange
      when(mockRemoteDataSource.updateStepStatus('1', 'step1', true))
          .thenAnswer((_) async => true);

      // act
      final result = await repository.updateStepStatus('1', 'step1', true);

      // assert
      expect(result, equals(Right(true)));
      verify(mockRemoteDataSource.updateStepStatus('1', 'step1', true));
      verifyNoMoreInteractions(mockRemoteDataSource);
    });

    test('should return error message when remote data source throws exception', () async {
      // arrange
      when(mockRemoteDataSource.updateStepStatus('1', 'step1', true))
          .thenThrow(Exception('Update failed'));

      // act
      final result = await repository.updateStepStatus('1', 'step1', true);

      // assert
      expect(result, equals(Left('Exception: Update failed')));
      verify(mockRemoteDataSource.updateStepStatus('1', 'step1', true));
      verifyNoMoreInteractions(mockRemoteDataSource);
    });
  });

  group('saveProgress', () {
    test('should return true when remote data source is successful', () async {
      // arrange
      when(mockRemoteDataSource.saveProgress('1'))
          .thenAnswer((_) async => true);

      // act
      final result = await repository.saveProgress('1');

      // assert
      expect(result, equals(Right(true)));
      verify(mockRemoteDataSource.saveProgress('1'));
      verifyNoMoreInteractions(mockRemoteDataSource);
    });

    test('should return error message when remote data source throws exception', () async {
      // arrange
      when(mockRemoteDataSource.saveProgress('1'))
          .thenThrow(Exception('Save failed'));

      // act
      final result = await repository.saveProgress('1');

      // assert
      expect(result, equals(Left('Exception: Save failed')));
      verify(mockRemoteDataSource.saveProgress('1'));
      verifyNoMoreInteractions(mockRemoteDataSource);
    });
  });
}


