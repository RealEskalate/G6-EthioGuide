import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/procedure/domain/repositories/procedure_repository.dart';
import '../repositories/workspace_procedure_repository.dart';

/// Use case for saving progress
class SaveProcedure {
  final ProcedureRepository repository;

  const SaveProcedure(this.repository);

  Future<Either<Failure, void>> call(String procedureId) async {
    return await repository.saveProcedure(procedureId);
  }
}
