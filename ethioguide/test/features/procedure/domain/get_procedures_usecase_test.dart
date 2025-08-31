import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure.dart';
import 'package:ethioguide/features/procedure/domain/repositories/procedure_repository.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_procedures.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

import 'get_procedures_usecase_test.mocks.dart';

@GenerateMocks([ProcedureRepository])
void main() {
  late GetProcedures usecase;
  late MockProcedureRepository mockRepository;

  setUp(() {
    mockRepository = MockProcedureRepository();
    usecase = GetProcedures(mockRepository);
  });

  final tProcedures = [
    const Procedure(
      id: '1',
      title: 'Passport Renewal',
      category: 'Travel',
      duration: '2-3 weeks',
      cost: '1,200 ETB',
      icon: 'passport',
      isQuickAccess: true,
    ),
  ];

  test('should return procedures from repository', () async {
    when(mockRepository.getProcedures()).thenAnswer((_) async => Right(tProcedures));

    final result = await usecase();

    expect(result, Right<Failure, List<Procedure>>(tProcedures));
    verify(mockRepository.getProcedures());
    verifyNoMoreInteractions(mockRepository);
  });
}


