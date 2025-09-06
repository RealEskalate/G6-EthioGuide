import 'package:dartz/dartz.dart';
import '../repositories/workspace_discussion_repository.dart';

/// Use case for reporting a comment
class ReportComment {
  final WorkspaceDiscussionRepository repository;

  const ReportComment(this.repository);

  Future<Either<String, bool>> call(String commentId) async {
    return await repository.reportComment(commentId);
  }
}
