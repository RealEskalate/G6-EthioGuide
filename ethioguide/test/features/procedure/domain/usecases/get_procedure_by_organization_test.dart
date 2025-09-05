import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_detail.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_step.dart';
import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';
import 'package:ethioguide/features/procedure/domain/repositories/workspace_procedure_repository.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_procedure_by_organization.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

import 'get_procedure_by_organization_test.mocks.dart';



@GenerateMocks([ProcedureDetailRepository])
void main() {
  late MockProcedureDetailRepository mockRepository;
  late GetProceduresByOrganization usecase;

  setUp(() {
    mockRepository = MockProcedureDetailRepository();
    usecase = GetProceduresByOrganization(mockRepository);
  });

  test('returns list on success', () async {
    when(mockRepository.getProceduresByOrganization('ETA')).thenAnswer((_) async => Right(const []));
    final result = await usecase('ETA');
    expect(result, Right(const []));
    verify(mockRepository.getProceduresByOrganization('ETA'));
    verifyNoMoreInteractions(mockRepository);
  });

  test('returns Failure on error', () async {
    when(mockRepository.getProceduresByOrganization('ETA')).thenAnswer((_) async => Left(ServerFailure()));
    final result = await usecase('ETA');
    expect(result, Left(ServerFailure()));
  });
}


