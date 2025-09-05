import 'package:dartz/dartz.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/community_stats.dart';
import 'package:ethioguide/features/workspace_discussion/domain/repositories/workspace_discussion_repository.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/get_community_stats.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';



@GenerateMocks([WorkspaceDiscussionRepository])
void main() {
  late MockRepo repo;
  late GetCommunityStats usecase;

  setUp(() {
    repo = MockRepo();
    usecase = GetCommunityStats(repo);
  });

  test('success', () async {
    when(repo.getCommunityStats()).thenAnswer((_) async => Right(const CommunityStats(totalMembers: 1, activeToday: 4 ,trendingTags: [], totalDiscussions: 0)));
    final result = await usecase();
    expect(result.isRight(), true);
    verify(repo.getCommunityStats());
    verifyNoMoreInteractions(repo);
  });

  test('failure', () async {
    when(repo.getCommunityStats()).thenAnswer((_) async => const Left('error'));
    final result = await usecase();
    expect(result, const Left('error'));
  });
}


