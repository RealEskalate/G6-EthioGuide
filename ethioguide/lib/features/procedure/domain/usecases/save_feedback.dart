import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import '../repositories/procedure_repository.dart';

class SaveFeedback {
  final ProcedureRepository repository;

  SaveFeedback(this.repository);

  Future<Either<Failure, bool>> call({
    required String procedureId,
    required String feedback,
    required List<String> tags,
    required String type,
  }) {
    return repository.saveFeedback(procedureId, feedback, tags, type);
  }
}