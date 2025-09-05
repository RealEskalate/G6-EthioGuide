import 'package:dartz/dartz.dart';
import 'package:ethioguide/features/procedure/domain/repositories/workspace_procedure_repository.dart';
import 'package:ethioguide/features/procedure/domain/usecases/save_progress.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

import 'save_progress_test.mocks.dart';





@GenerateMocks([ProcedureDetailRepository])
void main() {
  late MockProcedureDetailRepository mockRepository;
  late SaveProgress usecase;

  setUp(() {
    mockRepository = MockProcedureDetailRepository();
    usecase = SaveProgress(mockRepository);
  });

  test('returns true on success', () async {
    when(mockRepository.saveProgress('p1')).thenAnswer((_) async => const Right(true));
    final result = await usecase('p1');
    expect(result, const Right(true));
    verify(mockRepository.saveProgress('p1'));
    verifyNoMoreInteractions(mockRepository);
  });

  test('returns error string on failure', () async {
    when(mockRepository.saveProgress('p1')).thenAnswer((_) async => const Left('error'));
    final result = await usecase('p1');
    expect(result, const Left('error'));
  });
}


