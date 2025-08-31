import 'package:dartz/dartz.dart';
import '../entities/community_stats.dart';
import '../repositories/workspace_discussion_repository.dart';

/// Use case for fetching community statistics
class GetCommunityStats {
  final WorkspaceDiscussionRepository repository;

  const GetCommunityStats(this.repository);

  Future<Either<String, CommunityStats>> call() async {
    return await repository.getCommunityStats();
  }
}
