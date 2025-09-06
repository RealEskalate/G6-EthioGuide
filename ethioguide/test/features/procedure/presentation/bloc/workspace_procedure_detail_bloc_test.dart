import 'package:bloc_test/bloc_test.dart';
import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_detail.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_step.dart';
import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_my_procedure.dart' as usecase_my;
import 'package:ethioguide/features/procedure/domain/usecases/get_procedure_bystattus.dart' as usecase_by_status;
import 'package:ethioguide/features/procedure/domain/usecases/get_procedure_by_organization.dart' as usecase_by_org;
import 'package:ethioguide/features/procedure/domain/usecases/get_procedure_detail.dart' as usecase_detail;
import 'package:ethioguide/features/procedure/domain/usecases/get_workspace_summary.dart' as usecase_summary;
import 'package:ethioguide/features/procedure/domain/usecases/save_progress.dart' as usecase_save;
import 'package:ethioguide/features/procedure/domain/usecases/update_step_status.dart' as usecase_update;
import 'package:ethioguide/features/procedure/presentation/bloc/workspace_procedure_detail_bloc.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';

ProcedureDetail _detail() => ProcedureDetail(
      id: 'p1',
      title: 't',
      organization: 'o',
      status: ProcedureStatus.inProgress,
      progressPercentage: 0,
      documentsUploaded: 0,
      totalDocuments: 0,
      startDate: DateTime(2024, 1, 1),
      estimatedCompletion: null,
      completedDate: null,
      notes: null,
      steps: const [],
      estimatedTime: '1d',
      difficulty: 'Easy',
      officeType: 'Authority',
      quickTips: const [],
      requiredDocuments: const [],
    );

class _MockGetDetail extends Mock implements usecase_detail.GetProcedureDetail {}
class _MockUpdateStep extends Mock implements usecase_update.UpdateStepStatus {}
class _MockSaveProgress extends Mock implements usecase_save.SaveProgress {}
class _MockGetMy extends Mock implements usecase_my.GetProcedureDetails {}
class _MockGetByStatus extends Mock implements usecase_by_status.GetProceduresByStatus {}
class _MockGetByOrg extends Mock implements usecase_by_org.GetProceduresByOrganization {}
class _MockGetSummary extends Mock implements usecase_summary.GetWorkspaceSummary {}

void main() {
  group('WorkspaceProcedureDetailBloc', () {
    late _MockGetDetail getDetail;
    late _MockUpdateStep updateStep;
    late _MockSaveProgress saveProgress;
    late _MockGetMy getMy;
    late _MockGetByStatus getByStatus;
    late _MockGetByOrg getByOrg;
    late _MockGetSummary getSummary;
    late WorkspaceProcedureDetailBloc bloc;

    setUp(() {
      getDetail = _MockGetDetail();
      updateStep = _MockUpdateStep();
      saveProgress = _MockSaveProgress();
      getMy = _MockGetMy();
      getByStatus = _MockGetByStatus();
      getByOrg = _MockGetByOrg();
      getSummary = _MockGetSummary();
      bloc = WorkspaceProcedureDetailBloc(
        getProcedureDetail: getDetail,
        updateStepStatusUseCase: updateStep,
        saveProgressUseCase: saveProgress,
        getMyProcedureDetails: getMy,
        getProceduresByStatus: getByStatus,
        getProceduresByOrganization: getByOrg,
        getWorkspaceSummary: getSummary,
      );
    });

    test('initial state is ProcedureInitial', () {
      expect(bloc.state, isA<ProcedureInitial>());
    });

    blocTest<WorkspaceProcedureDetailBloc, WorkspaceProcedureDetailState>(
      'FetchProcedureDetail -> [Loading, Loaded] on success',
      build: () {
        when(getDetail('p1')).thenAnswer((_) async => Right(_detail()));
        return bloc;
      },
      act: (b) => b.add(const FetchProcedureDetail('p1')),
      expect: () => [isA<ProcedureLoading>(), isA<ProcedureLoaded>()],
      verify: (_) {
        verify(getDetail('p1'));
      }
    );

    blocTest<WorkspaceProcedureDetailBloc, WorkspaceProcedureDetailState>(
      'UpdateStepStatus -> StepStatusUpdated when refresh succeeds',
      build: () {
        when(updateStep('p1', 's1', true)).thenAnswer((_) async => const Right(true));
        when(getDetail('p1')).thenAnswer((_) async => Right(_detail()));
        return bloc;
      },
      act: (b) => b.add(const UpdateStepStatus('p1', 's1', true)),
      expect: () => [isA<StepStatusUpdated>()],
      verify: (_) {
        verify(updateStep('p1', 's1', true));
        verify(getDetail('p1'));
      }
    );

    blocTest<WorkspaceProcedureDetailBloc, WorkspaceProcedureDetailState>(
      'SaveProgress -> ProgressSaved on success',
      build: () {
        when(saveProgress('p1')).thenAnswer((_) async => const Right(true));
        return bloc;
      },
      act: (b) => b.add(const SaveProgress('p1')),
      expect: () => [isA<ProgressSaved>()],
    );

    blocTest<WorkspaceProcedureDetailBloc, WorkspaceProcedureDetailState>(
      'FetchMyProcedures -> [Loading, ProceduresListLoaded] on success',
      build: () {
        when(getMy()).thenAnswer((_) async => Right(const []));
        return bloc;
      },
      act: (b) => b.add(const FetchMyProcedures()),
      expect: () => [isA<ProcedureLoading>(), isA<ProceduresListLoaded>()],
    );

    blocTest<WorkspaceProcedureDetailBloc, WorkspaceProcedureDetailState>(
      'FetchProceduresByStatus -> [Loading, ProceduresListLoaded] on success',
      build: () {
        when(getByStatus(ProcedureStatus.inProgress)).thenAnswer((_) async => Right(const []));
        return bloc;
      },
      act: (b) => b.add(const FetchProceduresByStatus(ProcedureStatus.inProgress)),
      expect: () => [isA<ProcedureLoading>(), isA<ProceduresListLoaded>()],
    );

    blocTest<WorkspaceProcedureDetailBloc, WorkspaceProcedureDetailState>(
      'FetchProceduresByOrganization -> [Loading, ProceduresListLoaded] on success',
      build: () {
        when(getByOrg('ETA')).thenAnswer((_) async => Right(const []));
        return bloc;
      },
      act: (b) => b.add(const FetchProceduresByOrganization('ETA')),
      expect: () => [isA<ProcedureLoading>(), isA<ProceduresListLoaded>()],
    );

    blocTest<WorkspaceProcedureDetailBloc, WorkspaceProcedureDetailState>(
      'FetchWorkspaceSummary -> [Loading, WorkspaceSummaryLoaded] on success',
      build: () {
        when(getSummary()).thenAnswer((_) async => Right(const WorkspaceSummary(totalProcedures: 0, inProgress: 0, completed: 0 , totalDocuments: 5)));
        return bloc;
      },
      act: (b) => b.add(const FetchWorkspaceSummary()),
      expect: () => [isA<ProcedureLoading>(), isA<WorkspaceSummaryLoaded>()],
    );
  });
}


