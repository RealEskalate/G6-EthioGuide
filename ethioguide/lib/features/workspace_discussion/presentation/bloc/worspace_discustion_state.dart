import 'package:equatable/equatable.dart';
import '../../domain/entities/discussion.dart';
import '../../domain/entities/comment.dart';
import '../../domain/entities/community_stats.dart';

abstract class WorkspaceDiscussionState extends Equatable {
  const WorkspaceDiscussionState();

  @override
  List<Object?> get props => [];
}

class DiscussionInitial extends WorkspaceDiscussionState {}

class DiscussionLoading extends WorkspaceDiscussionState {}

class DiscussionLoaded extends WorkspaceDiscussionState {
  final List<Discussion> discussions;
  final CommunityStats? communityStats;

  const DiscussionLoaded({
    required this.discussions,
    this.communityStats,
  });

  @override
  List<Object?> get props => [discussions, communityStats];
}

class DiscussionError extends WorkspaceDiscussionState {
  final String message;

  const DiscussionError(this.message);

  @override
  List<Object?> get props => [message];
}

class CommentLoading extends WorkspaceDiscussionState {}

class CommentLoaded extends WorkspaceDiscussionState {
  final List<Comment> comments;
  final String discussionId;

  const CommentLoaded({
    required this.comments,
    required this.discussionId,
  });

  @override
  List<Object?> get props => [comments, discussionId];
}

class CommentError extends WorkspaceDiscussionState {
  final String message;

  const CommentError(this.message);

  @override
  List<Object?> get props => [message];
}

class ActionSuccess extends WorkspaceDiscussionState {
  final String message;

  const ActionSuccess(this.message);

  @override
  List<Object?> get props => [message];
}

class ActionFailure extends WorkspaceDiscussionState {
  final String message;

  const ActionFailure(this.message);

  @override
  List<Object?> get props => [message];
}
