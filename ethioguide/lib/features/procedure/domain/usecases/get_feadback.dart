import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure.dart';
import '../repositories/procedure_repository.dart';

class GetFeedbacks {
  final ProcedureRepository repository;

  GetFeedbacks(this.repository);

  Future<Either<Failure, List<FeedbackItem>>> call(String procedureId) {
    return repository.getFeedbacks(procedureId);
  }
}