import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_detail.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_step.dart';
import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';
import 'package:ethioguide/features/procedure/domain/repositories/workspace_procedure_repository.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_my_procedure.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

import 'get_my_procedure_test.mocks.dart';

@GenerateMocks([ProcedureDetailRepository])
void main() {
  late MockProcedureDetailRepository mockRepository;
  late GetProcedureDetails usecase;

  setUp(() {
    mockRepository = MockProcedureDetailRepository();
    usecase = GetProcedureDetails(mockRepository);
  });

  final tDetails = [
    ProcedureDetail(
      id: 'p1',
      title: 't',
      organization: 'o',
      status: ProcedureStatus.inProgress,
      progressPercentage: 10,
      documentsUploaded: 1,
      totalDocuments: 3,
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
    ),
  ];

  test('returns list on success', () async {
    // arrange
    when(mockRepository.getProcedure()).thenAnswer((_) async => Right(tDetails));
    // act
    final result = await usecase();
    // assert
    expect(result, Right(tDetails));
    verify(mockRepository.getProcedure());
    verifyNoMoreInteractions(mockRepository);
  });

  test('returns Failure on error', () async {
    // arrange
    when(mockRepository.getProcedure()).thenAnswer((_) async => Left(ServerFailure()));
    // act
    final result = await usecase();
    // assert
    expect(result, Left(ServerFailure()));
    verify(mockRepository.getProcedure());
    verifyNoMoreInteractions(mockRepository);
  });
}


