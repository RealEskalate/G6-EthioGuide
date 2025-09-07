import 'package:dartz/dartz.dart';
import '../entities/discussion.dart';
import '../repositories/workspace_discussion_repository.dart';

/// Use case for creating a new discussion
class CreateDiscussion {
  final WorkspaceDiscussionRepository repository;

  const CreateDiscussion(this.repository);

  Future<Either<String, Discussion>> call({
    required String title,
    required String content,
    required List<String> tags,
    required List<String> procedure,
  }) async {
    return await repository.createDiscussion(
      title: title,
      content: content,
      tags: tags,
      procedure: procedure,
      // category: category,
    );
  }
}
