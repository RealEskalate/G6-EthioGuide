import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure.dart';
import 'package:ethioguide/features/procedure/domain/repositories/procedure_repository.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_procedures.dart';
import 'package:flutter_test/flutter_test.dart';

class _MockProcedureRepository implements ProcedureRepository {
  Either<Failure, List<Procedure>> result;
  _MockProcedureRepository(this.result);

  @override
  Future<Either<Failure, List<Procedure>>> getProcedures() async {
    return result;
  }
}

void main() {
  group('GetProcedures', () {
    test('returns list of procedures on success', () async {
      final procedures = [
        const Procedure(
          id: '1',
          title: 'Passport Renewal',
          category: 'Travel',
          duration: '2 weeks',
          cost: '1200 ETB',
          icon: 'passport',
          isQuickAccess: true,
        ),
      ];
      final repository = _MockProcedureRepository(Right(procedures));
      final usecase = GetProcedures(repository);

      final result = await usecase();

      expect(result.isRight(), true);
      result.fold((_) => fail('Expected Right'), (value) {
        expect(value, procedures);
      });
    });

    test('returns failure on error', () async {
      final repository = _MockProcedureRepository(const Left(ServerFailure()));
      final usecase = GetProcedures(repository);

      final result = await usecase();

      expect(result.isLeft(), true);
      result.fold((failure) => expect(failure, isA<ServerFailure>()), (_) => fail('Expected Left'));
    });
  });
}


