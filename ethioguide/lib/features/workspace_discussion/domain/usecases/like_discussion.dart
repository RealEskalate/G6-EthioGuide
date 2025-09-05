import 'package:dartz/dartz.dart';
import '../repositories/workspace_discussion_repository.dart';

/// Use case for liking a discussion
class LikeDiscussion {
  final WorkspaceDiscussionRepository repository;

  const LikeDiscussion(this.repository);

  Future<Either<String, bool>> call(String discussionId) async {
    return await repository.likeDiscussion(discussionId);
  }
}
