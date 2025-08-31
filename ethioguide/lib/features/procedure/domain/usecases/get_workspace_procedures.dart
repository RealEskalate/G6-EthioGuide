import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';
import 'package:ethioguide/features/procedure/domain/repositories/workspace_procedure_repository.dart';

/// Use case for getting all workspace procedures
class GetWorkspaceProcedures  {
  final WorkspaceProcedureRepository repository;

  const GetWorkspaceProcedures(this.repository);

  @override
  Future<Either<Failure, List<WorkspaceProcedure>>> call() async {
    return await repository.getWorkspaceProcedures();
  }
}
