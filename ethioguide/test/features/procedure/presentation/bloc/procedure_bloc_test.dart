import 'package:bloc_test/bloc_test.dart';
import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_procedures.dart';
import 'package:ethioguide/features/procedure/presentation/bloc/procedure_bloc.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';

class _MockGetProcedures extends Mock implements GetProcedures {}

void main() {
  group('ProcedureBloc', () {
    late _MockGetProcedures getProcedures;
    late ProcedureBloc bloc;

    setUp(() {
      getProcedures = _MockGetProcedures();
      bloc = ProcedureBloc(getProcedures: getProcedures);
    });

    test('initial state is initial()', () {
      expect(bloc.state.status, ProcedureStatus.initial);
    });

    blocTest<ProcedureBloc, ProcedureState>(
      'emits [loading, success] with data',
      build: () {
        when(getProcedures()).thenAnswer((_) async => Right(const [
              Procedure(
                id: '1',
                title: 't',
                category: 'c',
                duration: 'd',
                cost: '0',
                icon: 'i',
                isQuickAccess: false,
              )
            ]));
        return bloc;
      },
      act: (b) => b.add(LoadProceduresEvent()),
      expect: () => [
        isA<ProcedureState>().having((s) => s.status, 'status', ProcedureStatus.loading),
        isA<ProcedureState>()
            .having((s) => s.status, 'status', ProcedureStatus.success)
            .having((s) => s.procedures.length, 'procedures', 1),
      ],
      verify: (_) {
        verify(getProcedures());
        verifyNoMoreInteractions(getProcedures);
      }
    );

    blocTest<ProcedureBloc, ProcedureState>(
      'emits [loading, failure] on error',
      build: () {
        when(getProcedures()).thenAnswer((_) async => const Left(ServerFailure()));
        return bloc;
      },
      act: (b) => b.add(LoadProceduresEvent()),
      expect: () => [
        isA<ProcedureState>().having((s) => s.status, 'status', ProcedureStatus.loading),
        isA<ProcedureState>().having((s) => s.status, 'status', ProcedureStatus.failure),
      ],
    );
  });
}


