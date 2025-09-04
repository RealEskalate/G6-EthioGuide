import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/procedure/domain/repositories/workspace_procedure_repository.dart';
import '../entities/procedure_detail.dart';

class GetProceduresByOrganization {
  final ProcedureDetailRepository repository;

  GetProceduresByOrganization(this.repository);

  Future<Either<Failure, List<ProcedureDetail>>> call(String organization) {
    return repository.getProceduresByOrganization(organization);
  }
}
