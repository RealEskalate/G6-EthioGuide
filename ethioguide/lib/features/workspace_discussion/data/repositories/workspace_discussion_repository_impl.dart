import 'package:dartz/dartz.dart';
import '../../domain/entities/discussion.dart';
import '../../domain/entities/comment.dart';
import '../../domain/entities/community_stats.dart';
import '../../domain/repositories/workspace_discussion_repository.dart';
import '../datasources/workspace_discussion_remote_data_source.dart';

/// Repository implementation for workspace discussion operations
class WorkspaceDiscussionRepositoryImpl implements WorkspaceDiscussionRepository {
  final WorkspaceDiscussionRemoteDataSource remoteDataSource;

  const WorkspaceDiscussionRepositoryImpl(this.remoteDataSource);

  @override
  Future<Either<String, CommunityStats>> getCommunityStats() async {
    try {
      final result = await remoteDataSource.getCommunityStats();
      return Right(result);
    } catch (e) {
      return Left(e.toString());
    }
  }

  @override
  Future<Either<String, List<Discussion>>> getDiscussions({
    String? tag,
    String? category,
    String? filterType,
  }) async {
    try {
      final result = await remoteDataSource.getDiscussions(
        tag: tag,
        category: category,
        filterType: filterType,
      );
      return Right(result);
    } catch (e) {
      return Left(e.toString());
    }
  }

  @override
  Future<Either<String, Discussion>> createDiscussion({
    required String title,
    required String content,
    required List<String> tags,
    required List<String> procedure,
    // required String category,
  }) async {
    try {
      final result = await remoteDataSource.createDiscussion(
        title: title,
        content: content,
        tags: tags,
        procedure: procedure,
        // category: category,
      );
      return Right(result);
    } catch (e) {
      return Left(e.toString());
    }
  }

  @override
  Future<Either<String, bool>> likeDiscussion(String discussionId) async {
    try {
      final result = await remoteDataSource.likeDiscussion(discussionId);
      return Right(result);
    } catch (e) {
      return Left(e.toString());
    }
  }

  @override
  Future<Either<String, bool>> reportDiscussion(String discussionId) async {
    try {
      final result = await remoteDataSource.reportDiscussion(discussionId);
      return Right(result);
    } catch (e) {
      return Left(e.toString());
    }
  }

  @override
  Future<Either<String, List<Comment>>> getComments(String discussionId) async {
    try {
      final result = await remoteDataSource.getComments(discussionId);
      return Right(result);
    } catch (e) {
      return Left(e.toString());
    }
  }

  @override
  Future<Either<String, Comment>> addComment({
    required String discussionId,
    required String content,
  }) async {
    try {
      final result = await remoteDataSource.addComment(
        discussionId: discussionId,
        content: content,
      );
      return Right(result);
    } catch (e) {
      return Left(e.toString());
    }
  }

  @override
  Future<Either<String, bool>> likeComment(String commentId) async {
    try {
      final result = await remoteDataSource.likeComment(commentId);
      return Right(result);
    } catch (e) {
      return Left(e.toString());
    }
  }

  @override
  Future<Either<String, bool>> reportComment(String commentId) async {
    try {
      final result = await remoteDataSource.reportComment(commentId);
      return Right(result);
    } catch (e) {
      return Left(e.toString());
    }
  }
}
