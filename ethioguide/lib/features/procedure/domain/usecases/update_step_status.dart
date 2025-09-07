import 'package:dartz/dartz.dart';
import '../repositories/workspace_procedure_repository.dart';

/// Use case for updating step statusUpdateStepSta
class  UpdateStepStatus{
  final ProcedureDetailRepository repository;

  const UpdateStepStatus(this.repository);

  Future<Either<String, bool>> call(String procedureId, String stepId, bool isCompleted) async {
    return await repository.updateStepStatus(procedureId, stepId, isCompleted);
  }
}
