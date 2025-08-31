import 'package:dartz/dartz.dart';
import '../entities/procedure_detail.dart';
import '../repositories/workspace_procedure_repository.dart';

/// Use case for fetching procedure details
class GetProcedureDetail {
  final WorkspaceProcedureRepository repository;

  const GetProcedureDetail(this.repository);

  Future<Either<String, ProcedureDetail>> call(String procedureId) async {
    return await repository.getProcedureDetail(procedureId);
  }
}
