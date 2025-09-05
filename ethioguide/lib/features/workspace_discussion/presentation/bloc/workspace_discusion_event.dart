import 'package:equatable/equatable.dart';

abstract class WorkspaceDiscussionEvent extends Equatable {
  const WorkspaceDiscussionEvent();

  @override
  List<Object?> get props => [];
}

class FetchCommunityStats extends WorkspaceDiscussionEvent {}

class FetchDiscussions extends WorkspaceDiscussionEvent {
  final String? tag;
  final String? category;
  final String? filterType;

  const FetchDiscussions({this.tag, this.category, this.filterType});
}

class CreateDiscussionEvent extends WorkspaceDiscussionEvent {
  final String title;
  final String content;
  final List<String> tags;
  final String category;

  const CreateDiscussionEvent({
    required this.title,
    required this.content,
    required this.tags,
    required this.category,
  });
}

class LikeDiscussionEvent extends WorkspaceDiscussionEvent {
  final String discussionId;
  const LikeDiscussionEvent(this.discussionId);
}

class ReportDiscussionEvent extends WorkspaceDiscussionEvent {
  final String discussionId;
  const ReportDiscussionEvent(this.discussionId);
}

class FetchCommentsEvent extends WorkspaceDiscussionEvent {
  final String discussionId;
  const FetchCommentsEvent(this.discussionId);
}

class AddCommentEvent extends WorkspaceDiscussionEvent {
  final String discussionId;
  final String content;

  const AddCommentEvent({
    required this.discussionId,
    required this.content,
  });
}

class LikeCommentEvent extends WorkspaceDiscussionEvent {
  final String commentId;
  const LikeCommentEvent(this.commentId);
}

class ReportCommentEvent extends WorkspaceDiscussionEvent {
  final String commentId;
  const ReportCommentEvent(this.commentId);
}
