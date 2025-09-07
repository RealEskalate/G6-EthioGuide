import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure.dart';
import 'package:ethioguide/features/procedure/domain/repositories/procedure_repository.dart';

class GetProceduresbyid {
  final ProcedureRepository repository;

  GetProceduresbyid(this.repository);

  Future<Either<Failure, Procedure>> call(String procedureId) {
    return repository.getProceduresbyid(procedureId);
  }
}