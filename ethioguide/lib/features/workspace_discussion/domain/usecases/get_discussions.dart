import 'package:dartz/dartz.dart';
import '../entities/discussion.dart';
import '../repositories/workspace_discussion_repository.dart';

/// Use case for fetching discussions with optional filtering
class GetDiscussions {
  final WorkspaceDiscussionRepository repository;

  const GetDiscussions(this.repository);

  Future<Either<String, List<Discussion>>> call({
    String? tag,
    String? category,
    String? filterType,
  }) async {
    return await repository.getDiscussions(
      tag: tag,
      category: category,
      filterType: filterType,
    );
  }
}
