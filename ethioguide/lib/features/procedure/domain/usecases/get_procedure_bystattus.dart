import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_step.dart';
import 'package:ethioguide/features/procedure/domain/repositories/workspace_procedure_repository.dart';
import '../entities/procedure_detail.dart';
import '../entities/workspace_procedure.dart';


class GetProceduresByStatus {
  final ProcedureDetailRepository repository;

  GetProceduresByStatus(this.repository);

  Future<Either<Failure, List<ProcedureDetail>>> call(ProcedureStatus status) {
    return repository.getProceduresByStatus(status);
  }
}
