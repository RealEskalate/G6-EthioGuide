import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure.dart';

abstract class ProcedureRepository {
  Future<Either<Failure, List<Procedure>>> getProcedures(String? name );
  Future<Either<Failure,bool >> saveProcedure(String procedureId);
   Future<Either<Failure, Procedure>> getProceduresbyid(String procedureId);

   // feedback related
  Future<Either<Failure, List<FeedbackItem>>> getFeedbacks(String procedureId);

  Future<Either<Failure, bool>> saveFeedback(
    String procedureId,
    String feedback,
    List<String> tags,
    String type,
  );
}


