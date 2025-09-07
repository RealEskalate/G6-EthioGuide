import 'package:dartz/dartz.dart';
import 'package:ethioguide/features/procedure/domain/repositories/workspace_procedure_repository.dart';
import 'package:ethioguide/features/procedure/domain/usecases/update_step_status.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

import 'update_step_status_test.mocks.dart';






@GenerateMocks([ProcedureDetailRepository])
void main() {
  late MockProcedureDetailRepository mockRepository;
  late UpdateStepStatus usecase;

  setUp(() {
    mockRepository = MockProcedureDetailRepository();
    usecase = UpdateStepStatus(mockRepository);
  });

  test('returns true on success', () async {
    when(mockRepository.updateStepStatus('p1', 's1', true)).thenAnswer((_) async => const Right(true));
    final result = await usecase('p1', 's1', true);
    expect(result, const Right(true));
    verify(mockRepository.updateStepStatus('p1', 's1', true));
    verifyNoMoreInteractions(mockRepository);
  });

  test('returns error string on failure', () async {
    when(mockRepository.updateStepStatus('p1', 's1', true)).thenAnswer((_) async => const Left('fail'));
    final result = await usecase('p1', 's1', true);
    expect(result, const Left('fail'));
  });
}


