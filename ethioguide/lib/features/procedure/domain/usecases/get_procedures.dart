import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure.dart';
import 'package:ethioguide/features/procedure/domain/repositories/procedure_repository.dart';

class GetProcedures {
  final ProcedureRepository repository;

  GetProcedures(this.repository);

  Future<Either<Failure, List<Procedure>>> call(String? name ) {
    return repository.getProcedures(name );
  }
}


