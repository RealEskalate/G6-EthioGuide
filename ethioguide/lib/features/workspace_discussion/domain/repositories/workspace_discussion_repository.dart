import 'package:dartz/dartz.dart';
import '../entities/discussion.dart';
import '../entities/comment.dart';
import '../entities/community_stats.dart';

/// Repository interface for workspace discussion operations
abstract class WorkspaceDiscussionRepository {
  /// Get community statistics
  Future<Either<String, CommunityStats>> getCommunityStats();
  
  /// Get all discussions with optional filtering
  Future<Either<String, List<Discussion>>> getDiscussions({
    String? tag,
    String? category,
    String? filterType, // 'trending', 'recent'
  });
  
  /// Create a new discussion
  Future<Either<String, Discussion>> createDiscussion({
    required String title,
    required String content,
    required List<String> tags,
    required List<String> procedure,
    // required String category,
  });
  
  /// Like a discussion
  Future<Either<String, bool>> likeDiscussion(String discussionId);
  
  /// Report a discussion
  Future<Either<String, bool>> reportDiscussion(String discussionId);
  
  /// Get comments for a discussion
  Future<Either<String, List<Comment>>> getComments(String discussionId);
  
  /// Add a comment to a discussion
  Future<Either<String, Comment>> addComment({
    required String discussionId,
    required String content,
  });
  
  /// Like a comment
  Future<Either<String, bool>> likeComment(String commentId);
  
  /// Report a comment
  Future<Either<String, bool>> reportComment(String commentId);
}
