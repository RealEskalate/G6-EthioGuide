import '../../domain/entities/comment.dart';
import 'user_model.dart';

/// Data model for comments
class CommentModel extends Comment {
  const CommentModel({
    required super.id,
    required super.discussionId,
    required super.content,
    required super.createdAt,
    required super.createdBy,
    required super.likeCount,
    required super.reportCount,
  });

  factory CommentModel.fromJson(Map<String, dynamic> json) {
    return CommentModel(
      id: json['id'] as String,
      discussionId: json['discussionId'] as String,
      content: json['content'] as String,
      createdAt: DateTime.parse(json['createdAt'] as String),
      createdBy: UserModel.fromJson(json['createdBy'] as Map<String, dynamic>),
      likeCount: json['likeCount'] as int,
      reportCount: json['reportCount'] as int,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'discussionId': discussionId,
      'content': content,
      'createdAt': createdAt.toIso8601String(),
      'createdBy': (createdBy as UserModel).toJson(),
      'likeCount': likeCount,
      'reportCount': reportCount,
    };
  }
}
