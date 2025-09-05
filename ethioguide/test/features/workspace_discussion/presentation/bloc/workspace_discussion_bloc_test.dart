import 'package:bloc_test/bloc_test.dart';
import 'package:dartz/dartz.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/comment.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/community_stats.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/discussion.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/user.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/add_comment.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/create_discussion.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/get_comments.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/get_community_stats.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/get_discussions.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/like_comment.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/like_discussion.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/report_comment.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/report_discussion.dart';
import 'package:ethioguide/features/workspace_discussion/presentation/bloc/workspace_discussion_bloc.dart';
import 'package:ethioguide/features/workspace_discussion/presentation/bloc/worspace_discustion_state.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

import 'workspace_discussion_bloc_test.mocks.dart';




@GenerateMocks([
  GetCommunityStats,
  GetDiscussions,
  CreateDiscussion,
  LikeDiscussion,
  ReportDiscussion,
  GetComments,
  AddComment,
  LikeComment,
  ReportComment,
])
void main() {
  group('WorkspaceDiscussionBloc', () {
    late MockGetCommunityStats getStats;
    late MockGetDiscussions getDiscussions;
    late MockCreateDiscussion createDiscussion;
    late MockLikeDiscussion likeDiscussion;
    late MockReportDiscussion reportDiscussion;
    late MockGetComments getComments;
    late MockAddComment addComment;
    late MockLikeComment likeComment;
    late MockReportComment reportComment;
    late WorkspaceDiscussionBloc bloc;

    setUp(() {
      getStats = MockGetCommunityStats();
      getDiscussions = MockGetDiscussions();
      createDiscussion = MockCreateDiscussion();
      likeDiscussion = MockLikeDiscussion();
      reportDiscussion = MockReportDiscussion();
      getComments = MockGetComments();
      addComment = MockAddComment();
      likeComment = MockLikeComment();
      reportComment = MockReportComment();
      bloc = WorkspaceDiscussionBloc(
        getCommunityStats: getStats,
        getDiscussions: getDiscussions,
        createDiscussion: createDiscussion,
        likeDiscussion: likeDiscussion,
        reportDiscussion: reportDiscussion,
        getComments: getComments,
        addComment: addComment,
        likeComment: likeComment,
        reportComment: reportComment,
      );
    });

    test('initial state is DiscussionInitial', () {
      expect(bloc.state, isA<DiscussionInitial>());
    });

    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'FetchCommunityStats -> [Loading, Loaded] on success',
      build: () {
        when(getStats()).thenAnswer(
          (_) async => const Right(
            CommunityStats(
              totalMembers: 0,
              totalDiscussions: 0,
              activeToday: 0,
              trendingTags: [],
            ),
          ),
        );
        return bloc;
      },
      act: (b) => b.add(FetchCommunityStats()),
      expect: () => [isA<DiscussionLoading>(), isA<DiscussionLoaded>()],
      verify: (_) {
        verify(getStats());
      },
    );

    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'FetchDiscussions -> [Loading, Loaded] on success',
      build: () {
        when(
          getDiscussions(
            tag: anyNamed('tag'),
            category: anyNamed('category'),
            filterType: anyNamed('filterType'),
          ),
        ).thenAnswer((_) async => const Right(<Discussion>[]));
        return bloc;
      },
      act: (b) => b.add(const FetchDiscussions()),
      expect: () => [isA<DiscussionLoading>(), isA<DiscussionLoaded>()],
    );

    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'CreateDiscussionEvent -> ActionSuccess on success',
      build: () {
        when(
          createDiscussion(
            title: 'title',
            content: 'content',
            tags: [],
            category: 'category',
          ),
        ).thenAnswer(
          (_) async => Right(
            Discussion(
              id: '1',
              title: 'title',
              content: 'content',
              tags: [],
              category: 'category',
              commentsCount: 0,
              likeCount: 0,
              reportCount: 0,
              createdAt: DateTime.now(),
              createdBy: User(id: '1', name: 'name'),
            ),
          ),
        );
        return bloc;
      },
      act: (b) => b.add(
        const CreateDiscussionEvent(
          title: 't',
          content: 'c',
          tags: [],
          category: 'cat',
        ),
      ),
      expect: () => [isA<ActionSuccess>()],
    );

    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'FetchCommentsEvent -> [CommentLoading, CommentLoaded] on success',
      build: () {
        when(
          getComments('1'),
        ).thenAnswer((_) async => const Right(<Comment>[]));
        return bloc;
      },
      act: (b) => b.add(const FetchCommentsEvent('1')),
      expect: () => [isA<CommentLoading>(), isA<CommentLoaded>()],
    );
  });
}
