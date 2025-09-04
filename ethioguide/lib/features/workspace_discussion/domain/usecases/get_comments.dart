import 'package:dartz/dartz.dart';
import '../entities/comment.dart';
import '../repositories/workspace_discussion_repository.dart';

/// Use case for fetching comments for a discussion
class GetComments {
  final WorkspaceDiscussionRepository repository;

  const GetComments(this.repository);

  Future<Either<String, List<Comment>>> call(String discussionId) async {
    return await repository.getComments(discussionId);
  }
}
