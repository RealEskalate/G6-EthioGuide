import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure.dart';

abstract class ProcedureRepository {
  Future<Either<Failure, List<Procedure>>> getProcedures();
}


