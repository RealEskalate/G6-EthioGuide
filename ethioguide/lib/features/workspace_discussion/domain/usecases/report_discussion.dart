import 'package:dartz/dartz.dart';
import '../repositories/workspace_discussion_repository.dart';

/// Use case for reporting a discussion
class ReportDiscussion {
  final WorkspaceDiscussionRepository repository;

  const ReportDiscussion(this.repository);

  Future<Either<String, bool>> call(String discussionId) async {
    return await repository.reportDiscussion(discussionId);
  }
}
