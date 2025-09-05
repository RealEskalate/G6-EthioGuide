import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';
import 'package:ethioguide/features/procedure/domain/repositories/workspace_procedure_repository.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_workspace_summary.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

import 'get_workspace_summary_test.mocks.dart';






@GenerateMocks([ProcedureDetailRepository])
void main() {
  late MockProcedureDetailRepository mockRepository;
  late GetWorkspaceSummary usecase;

  setUp(() {
    mockRepository = MockProcedureDetailRepository();
    usecase = GetWorkspaceSummary(mockRepository);
  });

  test('returns summary on success', () async {
    when(mockRepository.getWorkspaceSummary()).thenAnswer((_) async => Right(const WorkspaceSummary(
      totalProcedures: 1, inProgress: 1, completed: 0, totalDocuments: 5)));
    final result = await usecase();
    expect(result.isRight(), true);
    verify(mockRepository.getWorkspaceSummary());
    verifyNoMoreInteractions(mockRepository);
  });

  test('returns Failure on error', () async {
    when(mockRepository.getWorkspaceSummary()).thenAnswer((_) async => Left(ServerFailure()));
    final result = await usecase();
    expect(result, Left(ServerFailure()));
  });
}


