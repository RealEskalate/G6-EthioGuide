import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_detail.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_step.dart';
import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';
import 'package:ethioguide/features/procedure/domain/repositories/workspace_procedure_repository.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_procedure_bystattus.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';

class MockProcedureDetailRepository extends Mock implements ProcedureDetailRepository {}

void main() {
  late MockProcedureDetailRepository mockRepository;
  late GetProceduresByStatus usecase;

  setUp(() {
    mockRepository = MockProcedureDetailRepository();
    usecase = GetProceduresByStatus(mockRepository);
  });

  test('returns list on success', () async {
    when(mockRepository.getProceduresByStatus(ProcedureStatus.inProgress))
        .thenAnswer((_) async => Right(const []));
    final result = await usecase(ProcedureStatus.inProgress);
    expect(result, Right(const []));
    verify(mockRepository.getProceduresByStatus(ProcedureStatus.inProgress));
    verifyNoMoreInteractions(mockRepository);
  });

  test('returns Failure on error', () async {
    when(mockRepository.getProceduresByStatus(ProcedureStatus.inProgress))
        .thenAnswer((_) async => Left(ServerFailure()));
    final result = await usecase(ProcedureStatus.inProgress);
    expect(result, Left(ServerFailure()));
  });
}


