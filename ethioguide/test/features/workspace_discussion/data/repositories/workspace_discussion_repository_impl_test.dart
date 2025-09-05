import 'package:dartz/dartz.dart';
import 'package:ethioguide/features/workspace_discussion/data/datasources/workspace_discussion_remote_data_source.dart';
import 'package:ethioguide/features/workspace_discussion/data/models/comment_model.dart';
import 'package:ethioguide/features/workspace_discussion/data/models/community_stats_model.dart';
import 'package:ethioguide/features/workspace_discussion/data/models/discussion_model.dart';
import 'package:ethioguide/features/workspace_discussion/data/repositories/workspace_discussion_repository_impl.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/comment.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/community_stats.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/discussion.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/user.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';




@GenerateMocks([WorkspaceDiscussionRemoteDataSource])
void main() {
  group('WorkspaceDiscussionRepositoryImpl', () {
    late MockWorkspaceDiscussionRemoteDataSource remote;
    late WorkspaceDiscussionRepositoryImpl repository;

    setUp(() {
      remote = MockWorkspaceDiscussionRemoteDataSource();
      repository = WorkspaceDiscussionRepositoryImpl(remote);
    });

    test('getCommunityStats returns Right on success', () async {
      when(remote.getCommunityStats()).thenAnswer(
        (_) async => const CommunityStatsModel(
          totalMembers: 0,
          totalDiscussions: 0,
          activeToday: 0,
          trendingTags: [],
        ),
      );

      final result = await repository.getCommunityStats();

      expect(result, isA<Right>());
      verify(remote.getCommunityStats());
    });

    test('getDiscussions returns Right on success', () async {
      when(remote.getDiscussions(
        tag: anyNamed('tag'),
        category: anyNamed('category'),
        filterType: anyNamed('filterType'),
      )).thenAnswer((_) async => const <DiscussionModel>[]);

      final result = await repository.getDiscussions();

      expect(result, isA<Right>());
    });

    test('createDiscussion returns Right on success', () async {
      when(remote.createDiscussion(
        title: 'title',
        content: 'content',
        tags: [],
        category: 'category',
      )).thenAnswer(
        (_) async => DiscussionModel(
          id: '1',
          title: 'title',
          content: 'content',
          tags: const [],
          category: 'category',
          commentsCount: 0,
          likeCount: 0,
          reportCount: 0,
          createdAt: DateTime.now(),
          createdBy: const User(id: '1', name: 'name'),
        ),
      );

      final result = await repository.createDiscussion(
        title: 'title',
        content: 'content',
        tags: const [],
        category: 'category',
      );

      expect(result, isA<Right>());
    });

    test('likeDiscussion returns Right on success', () async {
      when(remote.likeDiscussion('1')).thenAnswer((_) async => true);

      final result = await repository.likeDiscussion('1');

      expect(result, isA<Right>());
    });

    test('reportDiscussion returns Right on success', () async {
      when(remote.reportDiscussion('1')).thenAnswer((_) async => true);

      final result = await repository.reportDiscussion('1');

      expect(result, isA<Right>());
    });

    test('getComments returns Right on success', () async {
      when(remote.getComments('1')).thenAnswer((_) async => const <CommentModel>[]);

      final result = await repository.getComments('1');

      expect(result, isA<Right>());
    });

    test('addComment returns Right on success', () async {
      when(remote.addComment(
        discussionId: 'discussionId',
        content: 'content',
      )).thenAnswer(
        (_) async => CommentModel(
          id: '1',
          discussionId: '1',
          content: 'c',
          createdAt: DateTime.now(),
          createdBy: const User(id: '1', name: 'name'),
          likeCount: 0,
          reportCount: 0,
        ),
      );

      final result = await repository.addComment(
        discussionId: '1',
        content: 'c',
      );

      expect(result, isA<Right>());
    });

    test('likeComment returns Right on success', () async {
      when(remote.likeComment('1')).thenAnswer((_) async => true);

      final result = await repository.likeComment('1');

      expect(result, isA<Right>());
    });

    test('reportComment returns Right on success', () async {
      when(remote.reportComment('1')).thenAnswer((_) async => true);

      final result = await repository.reportComment('1');

      expect(result, isA<Right>());
    });

    test('propagates Left on error', () async {
      when(remote.getCommunityStats()).thenThrow(Exception('boom'));

      final r1 = await repository.getCommunityStats();

      expect(r1, isA<Left>());
    });
  });
}
