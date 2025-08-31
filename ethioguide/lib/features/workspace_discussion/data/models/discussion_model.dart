import '../../domain/entities/discussion.dart';
import 'user_model.dart';

/// Data model for discussions
class DiscussionModel extends Discussion {
  const DiscussionModel({
    required super.id,
    required super.title,
    required super.content,
    required super.tags,
    required super.category,
    required super.createdAt,
    required super.createdBy,
    required super.likeCount,
    required super.reportCount,
    required super.commentsCount,
    super.isPinned = false,
  });

  factory DiscussionModel.fromJson(Map<String, dynamic> json) {
    return DiscussionModel(
      id: json['id'] as String,
      title: json['title'] as String,
      content: json['content'] as String,
      tags: (json['tags'] as List<dynamic>).cast<String>(),
      category: json['category'] as String,
      createdAt: DateTime.parse(json['createdAt'] as String),
      createdBy: UserModel.fromJson(json['createdBy'] as Map<String, dynamic>),
      likeCount: json['likeCount'] as int,
      reportCount: json['reportCount'] as int,
      commentsCount: json['commentsCount'] as int,
      isPinned: json['isPinned'] as bool? ?? false,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'content': content,
      'tags': tags,
      'category': category,
      'createdAt': createdAt.toIso8601String(),
      'createdBy': (createdBy as UserModel).toJson(),
      'likeCount': likeCount,
      'reportCount': reportCount,
      'commentsCount': commentsCount,
      'isPinned': isPinned,
    };
  }
}
