import 'package:dartz/dartz.dart';
import 'package:ethioguide/features/procedure/domain/repositories/workspace_procedure_repository.dart';
import '../entities/procedure_detail.dart';


class GetProcedureDetail {
  final ProcedureDetailRepository repository;

  GetProcedureDetail(this.repository);

  Future<Either<String, ProcedureDetail>> call(String id) {
    return repository.getProcedureDetail(id);
  }
}



