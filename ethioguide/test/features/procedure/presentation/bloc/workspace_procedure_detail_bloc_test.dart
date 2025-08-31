import 'package:flutter_test/flutter_test.dart';
import 'package:bloc_test/bloc_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:dartz/dartz.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_detail.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_step.dart';
import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_procedure_detail.dart';
import 'package:ethioguide/features/procedure/domain/usecases/update_step_status.dart';
import 'package:ethioguide/features/procedure/domain/usecases/save_progress.dart';
import 'package:ethioguide/features/procedure/presentation/bloc/workspace_procedure_detail_bloc.dart';

import 'workspace_procedure_detail_bloc_test.mocks.dart';

@GenerateMocks([
  GetProcedureDetail,
  UpdateStepStatus,
  SaveProgress,
])
void main() {
  late WorkspaceProcedureDetailBloc bloc;
  late MockGetProcedureDetail mockGetProcedureDetail;
  late MockUpdateStepStatus mockUpdateStepStatus;
  late MockSaveProgress mockSaveProgress;

  setUp(() {
    mockGetProcedureDetail = MockGetProcedureDetail();
    mockUpdateStepStatus = MockUpdateStepStatus();
    mockSaveProgress = MockSaveProgress();
    bloc = WorkspaceProcedureDetailBloc(
      getProcedureDetail: mockGetProcedureDetail,
      updateStepStatus: mockUpdateStepStatus,
      saveProgress: mockSaveProgress,
    );
  });

  tearDown(() {
    bloc.close();
  });

  test('initial state should be ProcedureInitial', () {
    expect(bloc.state, equals(ProcedureInitial()));
  });

  group('FetchProcedureDetail', () {
    final testProcedureDetail = ProcedureDetail(
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
        ProcedureStep(
          id: 'step1',
          title: 'Step 1',
          description: 'Test step 1',
          isCompleted: true,
          completionStatus: 'Completed',
          order: 1,
        ),
        ProcedureStep(
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

    blocTest<WorkspaceProcedureDetailBloc, WorkspaceProcedureDetailState>(
      'emits [ProcedureLoading, ProcedureLoaded] when FetchProcedureDetail succeeds',
      build: () {
        when(mockGetProcedureDetail('1'))
            .thenAnswer((_) async => Right(testProcedureDetail));
        return bloc;
      },
      act: (bloc) => bloc.add(FetchProcedureDetail('1')),
      expect: () => [
        ProcedureLoading(),
        ProcedureLoaded(testProcedureDetail),
      ],
    );

    blocTest<WorkspaceProcedureDetailBloc, WorkspaceProcedureDetailState>(
      'emits [ProcedureLoading, ProcedureError] when FetchProcedureDetail fails',
      build: () {
        when(mockGetProcedureDetail('1'))
            .thenAnswer((_) async => Left('Error message'));
        return bloc;
      },
      act: (bloc) => bloc.add(FetchProcedureDetail('1')),
      expect: () => [
        ProcedureLoading(),
        ProcedureError('Error message'),
      ],
    );
  });

  group('UpdateStepStatus', () {
    final testProcedureDetail = ProcedureDetail(
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
        ProcedureStep(
          id: 'step1',
          title: 'Step 1',
          description: 'Test step 1',
          isCompleted: true,
          completionStatus: 'Completed',
          order: 1,
        ),
        ProcedureStep(
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

    blocTest<WorkspaceProcedureDetailBloc, WorkspaceProcedureDetailState>(
      'emits [StepStatusUpdated] when UpdateStepStatus succeeds',
      build: () {
        when(mockUpdateStepStatus('1', 'step2', true))
            .thenAnswer((_) async => Right(true));
        return bloc;
      },
      seed: () => ProcedureLoaded(testProcedureDetail),
      act: (bloc) => bloc.add(UpdateStepStatus(
        procedureId: '1',
        stepId: 'step2',
        isCompleted: true,
      )),
      expect: () => [
        StepStatusUpdated(
          procedureDetail: testProcedureDetail.copyWith(
            steps: [
              testProcedureDetail.steps[0],
              testProcedureDetail.steps[1].copyWith(isCompleted: true),
            ],
            progressPercentage: 100,
          ),
          stepId: 'step2',
          isCompleted: true,
        ),
      ],
    );

    blocTest<WorkspaceProcedureDetailBloc, WorkspaceProcedureDetailState>(
      'emits [ProcedureError] when UpdateStepStatus fails',
      build: () {
        when(mockUpdateStepStatus('1', 'step2', true))
            .thenAnswer((_) async => Left('Error message'));
        return bloc;
      },
      seed: () => ProcedureLoaded(testProcedureDetail),
      act: (bloc) => bloc.add(UpdateStepStatus(
        procedureId: '1',
        stepId: 'step2',
        isCompleted: true,
      )),
      expect: () => [
        ProcedureError('Error message'),
      ],
    );
  });

  group('SaveProgress', () {
    blocTest<WorkspaceProcedureDetailBloc, WorkspaceProcedureDetailState>(
      'emits [ProgressSaved] when SaveProgress succeeds',
      build: () {
        when(mockSaveProgress('1'))
            .thenAnswer((_) async => Right(true));
        return bloc;
      },
      act: (bloc) => bloc.add(SaveProgress('1')),
      expect: () => [
        ProgressSaved(true),
      ],
    );

    blocTest<WorkspaceProcedureDetailBloc, WorkspaceProcedureDetailState>(
      'emits [ProgressSaved] when SaveProgress fails',
      build: () {
        when(mockSaveProgress('1'))
            .thenAnswer((_) async => Left('Error message'));
        return bloc;
      },
      act: (bloc) => bloc.add(SaveProgress('1')),
      expect: () => [
        ProgressSaved(false),
      ],
    );
  });
}
