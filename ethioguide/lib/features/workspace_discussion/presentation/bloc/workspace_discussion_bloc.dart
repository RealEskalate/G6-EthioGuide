import 'package:ethioguide/features/workspace_discussion/presentation/bloc/worspace_discustion_state.dart';
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



// -------- Events ----------
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
  final List<String> procedure;
  

  const CreateDiscussionEvent({
    required this.title,
    required this.content,
    required this.tags,
    required this.procedure
    
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

// -------- Bloc ----------
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
    // --- Community Stats ---
    on<FetchCommunityStats>((event, emit) async {
      emit(DiscussionLoading());
      final result = await getCommunityStats();
      result.fold(
        (failure) => emit(DiscussionError(failure)),
        (stats) => emit(DiscussionLoaded(discussions: [], communityStats: stats)),
      );
    });

    // --- Discussions ---
    on<FetchDiscussions>((event, emit) async {
      emit(DiscussionLoading());
      final result = await getDiscussions(
        tag: event.tag,
        category: event.category,
        filterType: event.filterType,
      );
      result.fold(
        (failure) => emit(DiscussionError(failure)),
        (discussions) => emit(DiscussionLoaded(discussions: discussions)),
      );
    });

    on<CreateDiscussionEvent>((event, emit) async {
      final result = await createDiscussion(
        title: event.title,
        content: event.content,
        tags: event.tags,
        
        procedure: event.procedure
      );

      print(result);


      result.fold(
        (failure) => emit(ActionFailure(failure)),
        (_) => emit(const ActionSuccess('Discussion created successfully!')),
      );
    });

    on<LikeDiscussionEvent>((event, emit) async {
      final result = await likeDiscussion(event.discussionId);
      result.fold(
        (failure) => emit(ActionFailure(failure)),
        (_) => emit(const ActionSuccess('Discussion liked!')),
      );
    });

    on<ReportDiscussionEvent>((event, emit) async {
      final result = await reportDiscussion(event.discussionId);
      result.fold(
        (failure) => emit(ActionFailure(failure)),
        (_) => emit(const ActionSuccess('Discussion reported!')),
      );
    });

    // --- Comments ---
    on<FetchCommentsEvent>((event, emit) async {
      emit(CommentLoading());
      final result = await getComments(event.discussionId);
      result.fold(
        (failure) => emit(CommentError(failure)),
        (comments) => emit(CommentLoaded(
          comments: comments,
          discussionId: event.discussionId,
        )),
      );
    });

    on<AddCommentEvent>((event, emit) async {
      final result = await addComment(
        
        discussionId: event.discussionId,
        content: event.content,
      );
      result.fold(
        (failure) => emit(ActionFailure(failure)),
        (_) => emit(const ActionSuccess('Comment added successfully!')),
      );
    });

    on<LikeCommentEvent>((event, emit) async {
      final result = await likeComment(event.commentId);
      result.fold(
        (failure) => emit(ActionFailure(failure)),
        (_) => emit(const ActionSuccess('Comment liked!')),
      );
    });

    on<ReportCommentEvent>((event, emit) async {
      final result = await reportComment(event.commentId);
      result.fold(
        (failure) => emit(ActionFailure(failure)),
        (_) => emit(const ActionSuccess('Comment reported!')),
      );
    });
  }
}