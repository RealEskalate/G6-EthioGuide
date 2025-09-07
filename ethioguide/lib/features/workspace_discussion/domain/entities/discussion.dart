import 'package:equatable/equatable.dart';
import 'user.dart';

/// Domain entity representing a discussion in the workspace community
class Discussion extends Equatable {
  final String id;
  final String title;
  final String content;
  final List<String> tags;
  final List<String> procedure;
  final String category;
  final DateTime createdAt;
  final User createdBy;
  final int likeCount;
  final int reportCount;
  final int commentsCount;
  final bool isPinned; // For moderator pinned discussions

  const Discussion({
    required this.id,
    required this.title,
    required this.content,
    required this.tags,
    required this.procedure,
    required this.category,
    required this.createdAt,
    required this.createdBy,
    required this.likeCount,
    required this.reportCount,
    required this.commentsCount,
    this.isPinned = false,
  });

  @override
  List<Object?> get props => [
        id,
        title,
        content,
        tags,
        procedure,
        category,
        createdAt,
        createdBy,
        likeCount,
        reportCount,
        commentsCount,
        isPinned,
      ];
}
