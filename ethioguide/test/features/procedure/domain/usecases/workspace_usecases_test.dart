import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_detail.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_step.dart';
import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';
import 'package:ethioguide/features/procedure/domain/repositories/workspace_procedure_repository.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_my_procedure.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_procedure_bystattus.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_procedure_by_organization.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_procedure_detail.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_workspace_summary.dart';
import 'package:ethioguide/features/procedure/domain/usecases/save_progress.dart';
import 'package:ethioguide/features/procedure/domain/usecases/update_step_status.dart';
import 'package:flutter_test/flutter_test.dart';

class _RepoStub implements ProcedureDetailRepository {
  Either<Failure, List<ProcedureDetail>>? listResult;
  Either<Failure, WorkspaceSummary>? summaryResult;
  Either<Failure, List<ProcedureDetail>>? byStatusResult;
  Either<Failure, List<ProcedureDetail>>? byOrgResult;
  Either<String, ProcedureDetail>? detailResult;
  Either<String, bool>? updateResult;
  Either<String, bool>? saveResult;

  @override
  Future<Either<String, bool>> saveProgress(String procedureId) async {
    return saveResult!;
  }

  @override
  Future<Either<String, bool>> updateStepStatus(String procedureId, String stepId, bool isCompleted) async {
    return updateResult!;
  }

  @override
  Future<Either<Failure, List<ProcedureDetail>>> getProcedure() async {
    return listResult!;
  }

  @override
  Future<Either<String, ProcedureDetail>> getProcedureDetail(String id) async {
    return detailResult!;
  }

  @override
  Future<Either<Failure, List<ProcedureDetail>>> getProceduresByOrganization(String organization) async {
    return byOrgResult!;
  }

  @override
  Future<Either<Failure, List<ProcedureDetail>>> getProceduresByStatus(ProcedureStatus status) async {
    return byStatusResult!;
  }

  @override
  Future<Either<Failure, WorkspaceSummary>> getWorkspaceSummary() async {
    return summaryResult!;
  }
}

ProcedureDetail _detail({String id = 'p1'}) => ProcedureDetail(
      id: id,
      title: 'Driver License',
      organization: 'ETA',
      status: ProcedureStatus.inProgress,
      progressPercentage: 40,
      documentsUploaded: 1,
      totalDocuments: 3,
      startDate: DateTime(2024, 1, 1),
      estimatedCompletion: null,
      completedDate: null,
      notes: null,
      steps: const [
        MyProcedureStep(
          id: 's1',
          title: 'Fill form',
          description: 'desc',
          isCompleted: false,
          completionStatus: null,
          order: 1,
        )
      ],
      estimatedTime: '2d',
      difficulty: 'Easy',
      officeType: 'Authority',
      quickTips: const [],
      requiredDocuments: const [],
    );

void main() {
  test('GetProcedureDetails returns list', () async {
    final repo = _RepoStub()..listResult = Right([_detail()]);
    final usecase = GetProcedureDetails(repo);
    final result = await usecase();
    expect(result.isRight(), true);
  });

  test('GetWorkspaceSummary returns summary', () async {
    final repo = _RepoStub()..summaryResult = Right(const WorkspaceSummary(totalProcedures: 1, inProgress: 1, completed: 0, totalDocuments: 5));
    final usecase = GetWorkspaceSummary(repo);
    final result = await usecase();
    expect(result.isRight(), true);
  });

  test('GetProceduresByStatus returns filtered list', () async {
    final repo = _RepoStub()..byStatusResult = Right([_detail()]);
    final usecase = GetProceduresByStatus(repo);
    final result = await usecase(ProcedureStatus.inProgress);
    expect(result.isRight(), true);
  });

  test('GetProceduresByOrganization returns filtered list', () async {
    final repo = _RepoStub()..byOrgResult = Right([_detail()]);
    final usecase = GetProceduresByOrganization(repo);
    final result = await usecase('ETA');
    expect(result.isRight(), true);
  });

  test('GetProcedureDetail returns detail', () async {
    final repo = _RepoStub()..detailResult = Right(_detail());
    final usecase = GetProcedureDetail(repo);
    final result = await usecase('p1');
    expect(result.isRight(), true);
  });

  test('UpdateStepStatus returns success', () async {
    final repo = _RepoStub()..updateResult = const Right(true);
    final usecase = UpdateStepStatus(repo);
    final result = await usecase('p1', 's1', true);
    expect(result.isRight(), true);
  });

  test('SaveProgress returns success', () async {
    final repo = _RepoStub()..saveResult = const Right(true);
    final usecase = SaveProgress(repo);
    final result = await usecase('p1');
    expect(result.isRight(), true);
  });
}


