import 'package:bloc_test/bloc_test.dart';
import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_procedures.dart';
import 'package:ethioguide/features/procedure/presentation/bloc/procedure_bloc.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

import 'procedure_bloc_test.mocks.dart';

@GenerateMocks([GetProcedures])
void main() {
  late ProcedureBloc bloc;
  late MockGetProcedures mockUsecase;

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

  setUp(() {
    mockUsecase = MockGetProcedures();
    bloc = ProcedureBloc(getProcedures: mockUsecase);
  });

  blocTest<ProcedureBloc, ProcedureState>('emits [loading, success] on success',
      build: () {
        when(mockUsecase()).thenAnswer((_) async => Right(tProcedures));
        return bloc;
      },
      act: (b) => b.add(const LoadProceduresEvent()),
      expect: () => [
            const ProcedureState.initial().copyWith(status: ProcedureStatus.loading),
            ProcedureState.initial().copyWith(status: ProcedureStatus.success, procedures: tProcedures),
          ]);

  blocTest<ProcedureBloc, ProcedureState>('emits [loading, failure] on error',
      build: () {
        when(mockUsecase()).thenAnswer((_) async => const Left(ServerFailure()));
        return bloc;
      },
      act: (b) => b.add(const LoadProceduresEvent()),
      expect: () => [
            const ProcedureState.initial().copyWith(status: ProcedureStatus.loading),
            const ProcedureState.initial().copyWith(status: ProcedureStatus.failure, errorMessage: 'Server Failure'),
          ]);
}


