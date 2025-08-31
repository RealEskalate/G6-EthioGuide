import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:dartz/dartz.dart';
import 'package:ethioguide/features/workspace_discussion/data/repositories/workspace_discussion_repository_impl.dart';
import 'package:ethioguide/features/workspace_discussion/data/datasources/workspace_discussion_remote_data_source.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/discussion.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/comment.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/community_stats.dart';
import 'package:ethioguide/features/workspace_discussion/data/models/discussion_model.dart';
import 'package:ethioguide/features/workspace_discussion/data/models/comment_model.dart';
import 'package:ethioguide/features/workspace_discussion/data/models/community_stats_model.dart';
import 'package:ethioguide/features/workspace_discussion/data/models/user_model.dart';

import 'workspace_discussion_repository_impl_test.mocks.dart';

@GenerateMocks([WorkspaceDiscussionRemoteDataSource])
void main() {
  late WorkspaceDiscussionRepositoryImpl repository;
  late MockWorkspaceDiscussionRemoteDataSource mockRemoteDataSource;

  setUp(() {
    mockRemoteDataSource = MockWorkspaceDiscussionRemoteDataSource();
    repository = WorkspaceDiscussionRepositoryImpl(mockRemoteDataSource);
  });

  group('getCommunityStats', () {
    test('should return CommunityStats when remote data source is successful', () async {
      // arrange
      final mockCommunityStats = CommunityStatsModel(
        totalMembers: 1200,
        totalDiscussions: 3,
        activeToday: 42,
        trendingTags: ['#passport', '#renewal'],
      );
      when(mockRemoteDataSource.getCommunityStats())
          .thenAnswer((_) async => mockCommunityStats);

      // act
      final result = await repository.getCommunityStats();

      // assert
      expect(result, Right(mockCommunityStats));
      verify(mockRemoteDataSource.getCommunityStats());
      verifyNoMoreInteractions(mockRemoteDataSource);
    });

    test('should return error message when remote data source throws exception', () async {
      // arrange
      when(mockRemoteDataSource.getCommunityStats())
          .thenThrow(Exception('Network error'));

      // act
      final result = await repository.getCommunityStats();

      // assert
      expect(result, Left('Exception: Network error'));
      verify(mockRemoteDataSource.getCommunityStats());
      verifyNoMoreInteractions(mockRemoteDataSource);
    });
  });

  group('getDiscussions', () {
    test('should return list of discussions when remote data source is successful', () async {
      // arrange
      final mockDiscussions = [
        DiscussionModel(
          id: '1',
          title: 'Test Discussion',
          content: 'Test content',
          tags: ['test'],
          category: 'General',
          createdAt: DateTime.now(),
          createdBy: const UserModel(id: '1', name: 'Test User'),
          likeCount: 5,
          reportCount: 0,
          commentsCount: 2,
        ),
      ];
      when(mockRemoteDataSource.getDiscussions(
        tag: 'test',
        category: 'General',
        filterType: 'recent',
      )).thenAnswer((_) async => mockDiscussions);

      // act
      final result = await repository.getDiscussions(
        tag: 'test',
        category: 'General',
        filterType: 'recent',
      );

      // assert
      expect(result, Right(mockDiscussions));
      verify(mockRemoteDataSource.getDiscussions(
        tag: 'test',
        category: 'General',
        filterType: 'recent',
      ));
      verifyNoMoreInteractions(mockRemoteDataSource);
    });

    test('should return error message when remote data source throws exception', () async {
      // arrange
      when(mockRemoteDataSource.getDiscussions())
          .thenThrow(Exception('Network error'));

      // act
      final result = await repository.getDiscussions();

      // assert
      expect(result, Left('Exception: Network error'));
      verify(mockRemoteDataSource.getDiscussions());
      verifyNoMoreInteractions(mockRemoteDataSource);
    });
  });

  group('createDiscussion', () {
    test('should return discussion when remote data source is successful', () async {
      // arrange
      final mockDiscussion = DiscussionModel(
        id: '1',
        title: 'New Discussion',
        content: 'New content',
        tags: ['new'],
        category: 'General',
        createdAt: DateTime.now(),
        createdBy: const UserModel(id: '1', name: 'Test User'),
        likeCount: 0,
        reportCount: 0,
        commentsCount: 0,
      );
      when(mockRemoteDataSource.createDiscussion(
        title: 'New Discussion',
        content: 'New content',
        tags: ['new'],
        category: 'General',
      )).thenAnswer((_) async => mockDiscussion);

      // act
      final result = await repository.createDiscussion(
        title: 'New Discussion',
        content: 'New content',
        tags: ['new'],
        category: 'General',
      );

      // assert
      expect(result, Right(mockDiscussion));
      verify(mockRemoteDataSource.createDiscussion(
        title: 'New Discussion',
        content: 'New content',
        tags: ['new'],
        category: 'General',
      ));
      verifyNoMoreInteractions(mockRemoteDataSource);
    });

    test('should return error message when remote data source throws exception', () async {
      // arrange
      when(mockRemoteDataSource.createDiscussion(
        title: 'New Discussion',
        content: 'New content',
        tags: ['new'],
        category: 'General',
      )).thenThrow(Exception('Network error'));

      // act
      final result = await repository.createDiscussion(
        title: 'New Discussion',
        content: 'New content',
        tags: ['new'],
        category: 'General',
      );

      // assert
      expect(result, Left('Exception: Network error'));
      verify(mockRemoteDataSource.createDiscussion(
        title: 'New Discussion',
        content: 'New content',
        tags: ['new'],
        category: 'General',
      ));
      verifyNoMoreInteractions(mockRemoteDataSource);
    });
  });

  group('likeDiscussion', () {
    test('should return true when remote data source is successful', () async {
      // arrange
      when(mockRemoteDataSource.likeDiscussion('1'))
          .thenAnswer((_) async => true);

      // act
      final result = await repository.likeDiscussion('1');

      // assert
      expect(result, Right(true));
      verify(mockRemoteDataSource.likeDiscussion('1'));
      verifyNoMoreInteractions(mockRemoteDataSource);
    });

    test('should return error message when remote data source throws exception', () async {
      // arrange
      when(mockRemoteDataSource.likeDiscussion('1'))
          .thenThrow(Exception('Network error'));

      // act
      final result = await repository.likeDiscussion('1');

      // assert
      expect(result, Left('Exception: Network error'));
      verify(mockRemoteDataSource.likeDiscussion('1'));
      verifyNoMoreInteractions(mockRemoteDataSource);
    });
  });

  group('reportDiscussion', () {
    test('should return true when remote data source is successful', () async {
      // arrange
      when(mockRemoteDataSource.reportDiscussion('1'))
          .thenAnswer((_) async => true);

      // act
      final result = await repository.reportDiscussion('1');

      // assert
      expect(result, Right(true));
      verify(mockRemoteDataSource.reportDiscussion('1'));
      verifyNoMoreInteractions(mockRemoteDataSource);
    });

    test('should return error message when remote data source throws exception', () async {
      // arrange
      when(mockRemoteDataSource.reportDiscussion('1'))
          .thenThrow(Exception('Network error'));

      // act
      final result = await repository.reportDiscussion('1');

      // assert
      expect(result, Left('Exception: Network error'));
      verify(mockRemoteDataSource.reportDiscussion('1'));
      verifyNoMoreInteractions(mockRemoteDataSource);
    });
  });

  group('getComments', () {
    test('should return list of comments when remote data source is successful', () async {
      // arrange
      final mockComments = [
        CommentModel(
          id: '1',
          discussionId: '1',
          content: 'Test comment',
          createdAt: DateTime.now(),
          createdBy: const UserModel(id: '1', name: 'Test User'),
          likeCount: 2,
          reportCount: 0,
        ),
      ];
      when(mockRemoteDataSource.getComments('1'))
          .thenAnswer((_) async => mockComments);

      // act
      final result = await repository.getComments('1');

      // assert
      expect(result, Right(mockComments));
      verify(mockRemoteDataSource.getComments('1'));
      verifyNoMoreInteractions(mockRemoteDataSource);
    });

    test('should return error message when remote data source throws exception', () async {
      // arrange
      when(mockRemoteDataSource.getComments('1'))
          .thenThrow(Exception('Network error'));

      // act
      final result = await repository.getComments('1');

      // assert
      expect(result, Left('Exception: Network error'));
      verify(mockRemoteDataSource.getComments('1'));
      verifyNoMoreInteractions(mockRemoteDataSource);
    });
  });

  group('addComment', () {
    test('should return comment when remote data source is successful', () async {
      // arrange
      final mockComment = CommentModel(
        id: '1',
        discussionId: '1',
        content: 'New comment',
        createdAt: DateTime.now(),
        createdBy: const UserModel(id: '1', name: 'Test User'),
        likeCount: 0,
        reportCount: 0,
      );
      when(mockRemoteDataSource.addComment(
        discussionId: '1',
        content: 'New comment',
      )).thenAnswer((_) async => mockComment);

      // act
      final result = await repository.addComment(
        discussionId: '1',
        content: 'New comment',
      );

      // assert
      expect(result, Right(mockComment));
      verify(mockRemoteDataSource.addComment(
        discussionId: '1',
        content: 'New comment',
      ));
      verifyNoMoreInteractions(mockRemoteDataSource);
    });

    test('should return error message when remote data source throws exception', () async {
      // arrange
      when(mockRemoteDataSource.addComment(
        discussionId: '1',
        content: 'New comment',
      )).thenThrow(Exception('Network error'));

      // act
      final result = await repository.addComment(
        discussionId: '1',
        content: 'New comment',
      );

      // assert
      expect(result, Left('Exception: Network error'));
      verify(mockRemoteDataSource.addComment(
        discussionId: '1',
        content: 'New comment',
      ));
      verifyNoMoreInteractions(mockRemoteDataSource);
    });
  });

  group('likeComment', () {
    test('should return true when remote data source is successful', () async {
      // arrange
      when(mockRemoteDataSource.likeComment('1'))
          .thenAnswer((_) async => true);

      // act
      final result = await repository.likeComment('1');

      // assert
      expect(result, Right(true));
      verify(mockRemoteDataSource.likeComment('1'));
      verifyNoMoreInteractions(mockRemoteDataSource);
    });

    test('should return error message when remote data source throws exception', () async {
      // arrange
      when(mockRemoteDataSource.likeComment('1'))
          .thenThrow(Exception('Network error'));

      // act
      final result = await repository.likeComment('1');

      // assert
      expect(result, Left('Exception: Network error'));
      verify(mockRemoteDataSource.likeComment('1'));
      verifyNoMoreInteractions(mockRemoteDataSource);
    });
  });

  group('reportComment', () {
    test('should return true when remote data source is successful', () async {
      // arrange
      when(mockRemoteDataSource.reportComment('1'))
          .thenAnswer((_) async => true);

      // act
      final result = await repository.reportComment('1');

      // assert
      expect(result, Right(true));
      verify(mockRemoteDataSource.reportComment('1'));
      verifyNoMoreInteractions(mockRemoteDataSource);
    });

    test('should return error message when remote data source throws exception', () async {
      // arrange
      when(mockRemoteDataSource.reportComment('1'))
          .thenThrow(Exception('Network error'));

      // act
      final result = await repository.reportComment('1');

      // assert
      expect(result, Left('Exception: Network error'));
      verify(mockRemoteDataSource.reportComment('1'));
      verifyNoMoreInteractions(mockRemoteDataSource);
    });
  });
}
