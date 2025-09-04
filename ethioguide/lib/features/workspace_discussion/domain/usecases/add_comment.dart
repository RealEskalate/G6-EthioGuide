import 'package:dartz/dartz.dart';
import '../entities/comment.dart';
import '../repositories/workspace_discussion_repository.dart';

/// Use case for adding a comment to a discussion
class AddComment {
  final WorkspaceDiscussionRepository repository;

  const AddComment(this.repository);

  Future<Either<String, Comment>> call({
    required String discussionId,
    required String content,
  }) async {
    return await repository.addComment(
      discussionId: discussionId,
      content: content,
    );
  }
}
