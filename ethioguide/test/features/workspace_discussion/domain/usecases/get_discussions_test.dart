import 'package:dartz/dartz.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/discussion.dart';
import 'package:ethioguide/features/workspace_discussion/domain/repositories/workspace_discussion_repository.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/get_discussions.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';



@GenerateMocks([WorkspaceDiscussionRepository])
void main() {
  late MockRepo repo;
  late GetDiscussions usecase;

  setUp(() {
    repo = MockRepo();
    usecase = GetDiscussions(repo);
  });

  test('success with filters', () async {
    when(repo.getDiscussions(tag: 't', category: 'c', filterType: 'recent')).thenAnswer((_) async => const Right(<Discussion>[]));
    final result = await usecase(tag: 't', category: 'c', filterType: 'recent');
    expect(result, const Right(<Discussion>[]));
    verify(repo.getDiscussions(tag: 't', category: 'c', filterType: 'recent'));
    verifyNoMoreInteractions(repo);
  });

  test('failure', () async {
    when(repo.getDiscussions()).thenAnswer((_) async => const Left('error'));
    final result = await usecase();
    expect(result, const Left('error'));
  });
}


