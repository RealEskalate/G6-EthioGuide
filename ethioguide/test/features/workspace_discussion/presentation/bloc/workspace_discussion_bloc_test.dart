import 'package:flutter_test/flutter_test.dart';
import 'package:bloc_test/bloc_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:dartz/dartz.dart';
import 'package:ethioguide/features/workspace_discussion/presentation/bloc/workspace_discussion_bloc.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/discussion.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/comment.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/community_stats.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/get_community_stats.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/get_discussions.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/create_discussion.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/like_discussion.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/report_discussion.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/get_comments.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/add_comment.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/like_comment.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/report_comment.dart';

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
  late WorkspaceDiscussionBloc bloc;
  late MockGetCommunityStats mockGetCommunityStats;
  late MockGetDiscussions mockGetDiscussions;
  late MockCreateDiscussion mockCreateDiscussion;
  late MockLikeDiscussion mockLikeDiscussion;
  late MockReportDiscussion mockReportDiscussion;
  late MockGetComments mockGetComments;
  late MockAddComment mockAddComment;
  late MockLikeComment mockLikeComment;
  late MockReportComment mockReportComment;

  setUp(() {
    mockGetCommunityStats = MockGetCommunityStats();
    mockGetDiscussions = MockGetDiscussions();
    mockCreateDiscussion = MockCreateDiscussion();
    mockLikeDiscussion = MockLikeDiscussion();
    mockReportDiscussion = MockReportDiscussion();
    mockGetComments = MockGetComments();
    mockAddComment = MockAddComment();
    mockLikeComment = MockLikeComment();
    mockReportComment = MockReportComment();

    bloc = WorkspaceDiscussionBloc(
      getCommunityStats: mockGetCommunityStats,
      getDiscussions: mockGetDiscussions,
      createDiscussion: mockCreateDiscussion,
      likeDiscussion: mockLikeDiscussion,
      reportDiscussion: mockReportDiscussion,
      getComments: mockGetComments,
      addComment: mockAddComment,
      likeComment: mockLikeComment,
      reportComment: mockReportComment,
    );
  });

  tearDown(() {
    bloc.close();
  });

  test('initial state should be DiscussionInitial', () {
    expect(bloc.state, equals(DiscussionInitial()));
  });

  group('FetchCommunityStats', () {
    final mockCommunityStats = CommunityStats(
      totalMembers: 1200,
      totalDiscussions: 3,
      activeToday: 42,
      trendingTags: ['#passport', '#renewal'],
    );

    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'emits [DiscussionLoading, DiscussionLoaded] when successful',
      build: () {
        when(mockGetCommunityStats())
            .thenAnswer((_) async => Right(mockCommunityStats));
        return bloc;
      },
      act: (bloc) => bloc.add(FetchCommunityStats()),
      expect: () => [
        DiscussionLoading(),
        DiscussionLoaded(
          discussions: [],
          communityStats: mockCommunityStats,
        ),
      ],
      verify: (_) {
        verify(mockGetCommunityStats()).called(1);
      },
    );

    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'emits [DiscussionLoading, DiscussionError] when unsuccessful',
      build: () {
        when(mockGetCommunityStats())
            .thenAnswer((_) async => Left('Error message'));
        return bloc;
      },
      act: (bloc) => bloc.add(FetchCommunityStats()),
      expect: () => [
        DiscussionLoading(),
        DiscussionError('Error message'),
      ],
      verify: (_) {
        verify(mockGetCommunityStats()).called(1);
      },
    );
  });

  group('FetchDiscussions', () {
    final mockDiscussions = [
      Discussion(
        id: '1',
        title: 'Test Discussion',
        content: 'Test content',
        tags: ['test'],
        category: 'General',
        createdAt: DateTime.now(),
        createdBy: const User(id: '1', name: 'Test User'),
        likeCount: 5,
        reportCount: 0,
        commentsCount: 2,
      ),
    ];

    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'emits [DiscussionLoading, DiscussionLoaded] when successful',
      build: () {
        when(mockGetDiscussions(
          tag: 'test',
          category: 'General',
          filterType: 'recent',
        )).thenAnswer((_) async => Right(mockDiscussions));
        return bloc;
      },
      act: (bloc) => bloc.add(FetchDiscussions(
        tag: 'test',
        category: 'General',
        filterType: 'recent',
      )),
      expect: () => [
        DiscussionLoading(),
        DiscussionLoaded(
          discussions: mockDiscussions,
          communityStats: null,
        ),
      ],
      verify: (_) {
        verify(mockGetDiscussions(
          tag: 'test',
          category: 'General',
          filterType: 'recent',
        )).called(1);
      },
    );

    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'emits [DiscussionLoading, DiscussionError] when unsuccessful',
      build: () {
        when(mockGetDiscussions())
            .thenAnswer((_) async => Left('Error message'));
        return bloc;
      },
      act: (bloc) => bloc.add(FetchDiscussions()),
      expect: () => [
        DiscussionLoading(),
        DiscussionError('Error message'),
      ],
      verify: (_) {
        verify(mockGetDiscussions()).called(1);
      },
    );
  });

  group('CreateDiscussion', () {
    final mockDiscussion = Discussion(
      id: '1',
      title: 'New Discussion',
      content: 'New content',
      tags: ['new'],
      category: 'General',
      createdAt: DateTime.now(),
      createdBy: const User(id: '1', name: 'Test User'),
      likeCount: 0,
      reportCount: 0,
      commentsCount: 0,
    );

    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'emits [ActionSuccess] when successful',
      build: () {
        when(mockCreateDiscussion(
          title: 'New Discussion',
          content: 'New content',
          tags: ['new'],
          category: 'General',
        )).thenAnswer((_) async => Right(mockDiscussion));
        return bloc;
      },
      act: (bloc) => bloc.add(CreateDiscussion(
        title: 'New Discussion',
        content: 'New content',
        tags: ['new'],
        category: 'General',
      )),
      expect: () => [
        ActionSuccess('Discussion created successfully!'),
      ],
      verify: (_) {
        verify(mockCreateDiscussion(
          title: 'New Discussion',
          content: 'New content',
          tags: ['new'],
          category: 'General',
        )).called(1);
      },
    );

    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'emits [ActionFailure] when unsuccessful',
      build: () {
        when(mockCreateDiscussion(
          title: 'New Discussion',
          content: 'New content',
          tags: ['new'],
          category: 'General',
        )).thenAnswer((_) async => Left('Error message'));
        return bloc;
      },
      act: (bloc) => bloc.add(CreateDiscussion(
        title: 'New Discussion',
        content: 'New content',
        tags: ['new'],
        category: 'General',
      )),
      expect: () => [
        ActionFailure('Error message'),
      ],
      verify: (_) {
        verify(mockCreateDiscussion(
          title: 'New Discussion',
          content: 'New content',
          tags: ['new'],
          category: 'General',
        )).called(1);
      },
    );
  });

  group('LikeDiscussion', () {
    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'emits [ActionSuccess] when successful',
      build: () {
        when(mockLikeDiscussion('1'))
            .thenAnswer((_) async => Right(true));
        return bloc;
      },
      act: (bloc) => bloc.add(LikeDiscussion('1')),
      expect: () => [
        ActionSuccess('Discussion liked!'),
      ],
      verify: (_) {
        verify(mockLikeDiscussion('1')).called(1);
      },
    );

    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'emits [ActionFailure] when unsuccessful',
      build: () {
        when(mockLikeDiscussion('1'))
            .thenAnswer((_) async => Left('Error message'));
        return bloc;
      },
      act: (bloc) => bloc.add(LikeDiscussion('1')),
      expect: () => [
        ActionFailure('Error message'),
      ],
      verify: (_) {
        verify(mockLikeDiscussion('1')).called(1);
      },
    );
  });

  group('ReportDiscussion', () {
    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'emits [ActionSuccess] when successful',
      build: () {
        when(mockReportDiscussion('1'))
            .thenAnswer((_) async => Right(true));
        return bloc;
      },
      act: (bloc) => bloc.add(ReportDiscussion('1')),
      expect: () => [
        ActionSuccess('Discussion reported!'),
      ],
      verify: (_) {
        verify(mockReportDiscussion('1')).called(1);
      },
    );

    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'emits [ActionFailure] when unsuccessful',
      build: () {
        when(mockReportDiscussion('1'))
            .thenAnswer((_) async => Left('Error message'));
        return bloc;
      },
      act: (bloc) => bloc.add(ReportDiscussion('1')),
      expect: () => [
        ActionFailure('Error message'),
      ],
      verify: (_) {
        verify(mockReportDiscussion('1')).called(1);
      },
    );
  });

  group('FetchComments', () {
    final mockComments = [
      Comment(
        id: '1',
        discussionId: '1',
        content: 'Test comment',
        createdAt: DateTime.now(),
        createdBy: const User(id: '1', name: 'Test User'),
        likeCount: 2,
        reportCount: 0,
      ),
    ];

    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'emits [CommentLoading, CommentLoaded] when successful',
      build: () {
        when(mockGetComments('1'))
            .thenAnswer((_) async => Right(mockComments));
        return bloc;
      },
      act: (bloc) => bloc.add(FetchComments('1')),
      expect: () => [
        CommentLoading(),
        CommentLoaded(
          comments: mockComments,
          discussionId: '1',
        ),
      ],
      verify: (_) {
        verify(mockGetComments('1')).called(1);
      },
    );

    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'emits [CommentLoading, CommentError] when unsuccessful',
      build: () {
        when(mockGetComments('1'))
            .thenAnswer((_) async => Left('Error message'));
        return bloc;
      },
      act: (bloc) => bloc.add(FetchComments('1')),
      expect: () => [
        CommentLoading(),
        CommentError('Error message'),
      ],
      verify: (_) {
        verify(mockGetComments('1')).called(1);
      },
    );
  });

  group('AddComment', () {
    final mockComment = Comment(
      id: '1',
      discussionId: '1',
      content: 'New comment',
      createdAt: DateTime.now(),
      createdBy: const User(id: '1', name: 'Test User'),
      likeCount: 0,
      reportCount: 0,
    );

    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'emits [ActionSuccess] when successful',
      build: () {
        when(mockAddComment(
          discussionId: '1',
          content: 'New comment',
        )).thenAnswer((_) async => Right(mockComment));
        return bloc;
      },
      act: (bloc) => bloc.add(AddComment(
        discussionId: '1',
        content: 'New comment',
      )),
      expect: () => [
        ActionSuccess('Comment added successfully!'),
      ],
      verify: (_) {
        verify(mockAddComment(
          discussionId: '1',
          content: 'New comment',
        )).called(1);
      },
    );

    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'emits [ActionFailure] when unsuccessful',
      build: () {
        when(mockAddComment(
          discussionId: '1',
          content: 'New comment',
        )).thenAnswer((_) async => Left('Error message'));
        return bloc;
      },
      act: (bloc) => bloc.add(AddComment(
        discussionId: '1',
        content: 'New comment',
      )),
      expect: () => [
        ActionFailure('Error message'),
      ],
      verify: (_) {
        verify(mockAddComment(
          discussionId: '1',
          content: 'New comment',
        )).called(1);
      },
    );
  });

  group('LikeComment', () {
    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'emits [ActionSuccess] when successful',
      build: () {
        when(mockLikeComment('1'))
            .thenAnswer((_) async => Right(true));
        return bloc;
      },
      act: (bloc) => bloc.add(LikeComment('1')),
      expect: () => [
        ActionSuccess('Comment liked!'),
      ],
      verify: (_) {
        verify(mockLikeComment('1')).called(1);
      },
    );

    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'emits [ActionFailure] when unsuccessful',
      build: () {
        when(mockLikeComment('1'))
            .thenAnswer((_) async => Left('Error message'));
        return bloc;
      },
      act: (bloc) => bloc.add(LikeComment('1')),
      expect: () => [
        ActionFailure('Error message'),
      ],
      verify: (_) {
        verify(mockLikeComment('1')).called(1);
      },
    );
  });

  group('ReportComment', () {
    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'emits [ActionSuccess] when successful',
      build: () {
        when(mockReportComment('1'))
            .thenAnswer((_) async => Right(true));
        return bloc;
      },
      act: (bloc) => bloc.add(ReportComment('1')),
      expect: () => [
        ActionSuccess('Comment reported!'),
      ],
      verify: (_) {
        verify(mockReportComment('1')).called(1);
      },
    );

    blocTest<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
      'emits [ActionFailure] when unsuccessful',
      build: () {
        when(mockReportComment('1'))
            .thenAnswer((_) async => Left('Error message'));
        return bloc;
      },
      act: (bloc) => bloc.add(ReportComment('1')),
      expect: () => [
        ActionFailure('Error message'),
      ],
      verify: (_) {
        verify(mockReportComment('1')).called(1);
      },
    );
  });
}
