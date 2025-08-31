import 'package:ethioguide/features/workspace_discussion/domain/repositories/workspace_discussion_repository.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import '../../domain/entities/discussion.dart';
import '../../domain/entities/comment.dart';
import '../../domain/entities/community_stats.dart';
import '../../domain/usecases/get_community_stats.dart';
import '../../domain/usecases/get_discussions.dart';
import '../../domain/usecases/create_discussion.dart';
import '../../domain/usecases/like_discussion.dart';
import '../../domain/usecases/report_discussion.dart';
import '../../domain/usecases/get_comments.dart';
import '../../domain/usecases/add_comment.dart';
import '../../domain/usecases/like_comment.dart';
import '../../domain/usecases/report_comment.dart';

// Events
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

  const FetchDiscussions({
    this.tag,
    this.category,
    this.filterType,
  });

  @override
  List<Object?> get props => [tag, category, filterType];
}

class CreateDiscussion extends WorkspaceDiscussionEvent {
  final String title;
  final String content;
  final List<String> tags;
  final String category;

  const CreateDiscussion(WorkspaceDiscussionRepository workspaceDiscussionRepository, {
    required this.title,
    required this.content,
    required this.tags,
    required this.category,
  });

  @override
  List<Object?> get props => [title, content, tags, category];
}

class LikeDiscussion extends WorkspaceDiscussionEvent {
  final String discussionId;

  const LikeDiscussion(this.discussionId);

  @override
  List<Object?> get props => [discussionId];
}

class ReportDiscussion extends WorkspaceDiscussionEvent {
  final String discussionId;

  const ReportDiscussion(this.discussionId);

  @override
  List<Object?> get props => [discussionId];
}

class FetchComments extends WorkspaceDiscussionEvent {
  final String discussionId;

  const FetchComments(this.discussionId);

  @override
  List<Object?> get props => [discussionId];
}

class AddComment extends WorkspaceDiscussionEvent {
  final String discussionId;
  final String content;

  const AddComment(WorkspaceDiscussionRepository workspaceDiscussionRepository, {
    required this.discussionId,
    required this.content,
  });

  @override
  List<Object?> get props => [discussionId, content];
}

class LikeComment extends WorkspaceDiscussionEvent {
  final String commentId;

  const LikeComment(this.commentId);

  @override
  List<Object?> get props => [commentId];
}

class ReportComment extends WorkspaceDiscussionEvent {
  final String commentId;

  const ReportComment(this.commentId);

  @override
  List<Object?> get props => [commentId];
}

// States
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

// Bloc
class WorkspaceDiscussionBloc
    extends Bloc<WorkspaceDiscussionEvent, WorkspaceDiscussionState> {
  final GetCommunityStats getCommunityStats;
  final GetDiscussions getDiscussions;
  final CreateDiscussion createDiscussion;
  final LikeDiscussion likeDiscussion;
  final ReportDiscussion reportDiscussion;
  final GetComments getComments;
  final AddComment addComment;
  final LikeComment likeComment;
  final ReportComment reportComment;

  WorkspaceDiscussionBloc({
    required this.getCommunityStats,
    required this.getDiscussions,
    required this.createDiscussion,
    required this.likeDiscussion,
    required this.reportDiscussion,
    required this.getComments,
    required this.addComment,
    required this.likeComment,
    required this.reportComment,
  }) : super(DiscussionInitial()) {
    on<FetchCommunityStats>(_onFetchCommunityStats);
    on<FetchDiscussions>(_onFetchDiscussions);
    on<CreateDiscussion>(_onCreateDiscussion);
    on<LikeDiscussion>(_onLikeDiscussion);
    on<ReportDiscussion>(_onReportDiscussion);
    on<FetchComments>(_onFetchComments);
    on<AddComment>(_onAddComment);
    on<LikeComment>(_onLikeComment);
    on<ReportComment>(_onReportComment);
  }

  Future<void> _onFetchCommunityStats(
    FetchCommunityStats event,
    Emitter<WorkspaceDiscussionState> emit,
  ) async {
    emit(DiscussionLoading());
    final result = await getCommunityStats();
    result.fold(
      (failure) => emit(DiscussionError(failure)),
      (communityStats) => emit(DiscussionLoaded(
        discussions: [],
        communityStats: communityStats,
      )),
    );
  }

  Future<void> _onFetchDiscussions(
    FetchDiscussions event,
    Emitter<WorkspaceDiscussionState> emit,
  ) async {
    emit(DiscussionLoading());
    final result = await getDiscussions(
      tag: event.tag,
      category: event.category,
      filterType: event.filterType,
    );
    result.fold(
      (failure) => emit(DiscussionError(failure)),
      (discussions) => emit(DiscussionLoaded(
        discussions: discussions,
        communityStats: null,
      )),
    );
  }

  Future<void> _onCreateDiscussion(
    CreateDiscussion event,
    Emitter<WorkspaceDiscussionState> emit,
  ) async {
    final result = await createDiscussion(
      title: event.title,
      content: event.content,
      tags: event.tags,
      category: event.category,

    );
    
    /* createDiscussion(
      title: event.title,
      content: event.content,
      tags: event.tags,
      category: event.category,
    ); */
    result.fold(
      (failure) => emit(ActionFailure(failure)),
      (discussion) => emit(ActionSuccess('Discussion created successfully!')),
    );
  }

  Future<void> _onLikeDiscussion(
    LikeDiscussion event,
    Emitter<WorkspaceDiscussionState> emit,
  ) async {
    final result = await likeDiscussion(event.discussionId);
    result.fold(
      (failure) => emit(ActionFailure(failure)),
      (success) => emit(ActionSuccess('Discussion liked!')),
    );
  }

  Future<void> _onReportDiscussion(
    ReportDiscussion event,
    Emitter<WorkspaceDiscussionState> emit,
  ) async {
    final result = await reportDiscussion(event.discussionId);
    result.fold(
      (failure) => emit(ActionFailure(failure)),
      (success) => emit(ActionSuccess('Discussion reported!')),
    );
  }

  Future<void> _onFetchComments(
    FetchComments event,
    Emitter<WorkspaceDiscussionState> emit,
  ) async {
    emit(CommentLoading());
    final result = await getComments(event.discussionId);
    result.fold(
      (failure) => emit(CommentError(failure)),
      (comments) => emit(CommentLoaded(
        comments: comments,
        discussionId: event.discussionId,
      )),
    );
  }

  Future<void> _onAddComment(
    AddComment event,
    Emitter<WorkspaceDiscussionState> emit,
  ) async {
    final result = await addComment(
      discussionId: event.discussionId,
      content: event.content,
    );
    result.fold(
      (failure) => emit(ActionFailure(failure)),
      (comment) => emit(ActionSuccess('Comment added successfully!')),
    );
  }

  Future<void> _onLikeComment(
    LikeComment event,
    Emitter<WorkspaceDiscussionState> emit,
  ) async {
    final result = await likeComment(event.commentId);
    result.fold(
      (failure) => emit(ActionFailure(failure)),
      (success) => emit(ActionSuccess('Comment liked!')),
    );
  }

  Future<void> _onReportComment(
    ReportComment event,
    Emitter<WorkspaceDiscussionState> emit,
  ) async {
    final result = await reportComment(event.commentId);
    result.fold(
      (failure) => emit(ActionFailure(failure)),
      (success) => emit(ActionSuccess('Comment reported!')),
    );
  }
}
