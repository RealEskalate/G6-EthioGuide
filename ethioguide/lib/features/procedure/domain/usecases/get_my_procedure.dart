import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/procedure/domain/repositories/workspace_procedure_repository.dart';
import '../entities/procedure_detail.dart';


class GetProcedureDetails {
  final ProcedureDetailRepository repository;

  GetProcedureDetails(this.repository);

  Future<Either<Failure, List<ProcedureDetail>>> call() {
    return repository.getProcedure();
  }
}
