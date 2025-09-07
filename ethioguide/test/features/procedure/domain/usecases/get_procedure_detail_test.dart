import 'package:dartz/dartz.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_detail.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_step.dart';
import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';
import 'package:ethioguide/features/procedure/domain/repositories/workspace_procedure_repository.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_procedure_detail.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

import 'get_procedure_detail_test.mocks.dart';






@GenerateMocks([ProcedureDetailRepository])
void main() {
  late MockProcedureDetailRepository mockRepository;
  late GetProcedureDetail usecase;

  setUp(() {
    mockRepository = MockProcedureDetailRepository();
    usecase = GetProcedureDetail(mockRepository);
  });

  final tDetail = ProcedureDetail(
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

  test('returns detail on success', () async {
    when(mockRepository.getProcedureDetail('p1')).thenAnswer((_) async => Right(tDetail));
    final result = await usecase('p1');
    expect(result, Right(tDetail));
    verify(mockRepository.getProcedureDetail('p1'));
    verifyNoMoreInteractions(mockRepository);
  });

  test('returns error string on failure', () async {
    when(mockRepository.getProcedureDetail('p1')).thenAnswer((_) async => const Left('error'));
    final result = await usecase('p1');
    expect(result, const Left('error'));
  });
}


