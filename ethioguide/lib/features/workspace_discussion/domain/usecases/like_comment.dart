import 'package:dartz/dartz.dart';
import '../repositories/workspace_discussion_repository.dart';

/// Use case for liking a comment
class LikeComment {
  final WorkspaceDiscussionRepository repository;

  const LikeComment(this.repository);

  Future<Either<String, bool>> call(String commentId) async {
    return await repository.likeComment(commentId);
  }
}
