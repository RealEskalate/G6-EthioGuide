import 'package:equatable/equatable.dart';
import 'user.dart';

/// Domain entity representing a comment on a discussion
class Comment extends Equatable {
  final String id;
  final String discussionId;
  final String content;
  final DateTime createdAt;
  final User createdBy;
  final int likeCount;
  final int reportCount;

  const Comment({
    required this.id,
    required this.discussionId,
    required this.content,
    required this.createdAt,
    required this.createdBy,
    required this.likeCount,
    required this.reportCount,
  });

  @override
  List<Object?> get props => [
        id,
        discussionId,
        content,
        createdAt,
        createdBy,
        likeCount,
        reportCount,
      ];
}
