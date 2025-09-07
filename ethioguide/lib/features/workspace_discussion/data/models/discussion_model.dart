import 'package:ethioguide/features/workspace_discussion/domain/entities/user.dart';

import '../../domain/entities/discussion.dart';
import 'user_model.dart';

/// Data model for discussions
class DiscussionModel extends Discussion {
  const DiscussionModel({
    required super.id,
    required super.title,
    required super.content,
    required super.tags,
    required super.procedure,
    required super.category,
    required super.createdAt,
    required super.createdBy,
    required super.likeCount,
    required super.reportCount,
    required super.commentsCount,
    super.isPinned = false,
  });

factory DiscussionModel.fromJson(Map<String, dynamic> json) {
  final data = json['post'] ?? json; // handle wrapped vs unwrapped

  return DiscussionModel(
    id: data["ID"] as String,
    title: data['Title'] as String,
    content: data['Content'] as String,
    tags: (data["Tags"] != null)
        ? (data["Tags"] as List<dynamic>).cast<String>()
        : [],
    procedure: (data["Procedures"] != null)
        ? (data["Procedures"] as List<dynamic>).cast<String>()
        : [],
    category: '', // backend doesn’t send yet
    createdAt: DateTime.parse(data['CreatedAt'] as String),
    createdBy: User(id: data["UserID"] ?? '', name: ''), // fix user fallback
    likeCount: 0,
    reportCount: 0,
    commentsCount: 0,
    isPinned: false,
  );
}


  static List<DiscussionModel> listFromJson(Map<String, dynamic> json) {
    final postsJson = json["Posts"]?["posts"] as List<dynamic>? ?? [];
    return postsJson.map((post) => DiscussionModel.fromJson(post)).toList();
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
