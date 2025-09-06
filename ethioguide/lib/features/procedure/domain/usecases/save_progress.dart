import 'package:dartz/dartz.dart';
import '../repositories/workspace_procedure_repository.dart';

/// Use case for saving progress
class SaveProgress {
  final ProcedureDetailRepository repository;

  const SaveProgress(this.repository);

  Future<Either<String, bool>> call(String procedureId) async {
    return await repository.saveProgress(procedureId);
  }
}
