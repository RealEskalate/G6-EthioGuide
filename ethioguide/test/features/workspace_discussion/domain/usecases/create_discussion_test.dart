import 'package:dartz/dartz.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/discussion.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/user.dart';
import 'package:ethioguide/features/workspace_discussion/domain/repositories/workspace_discussion_repository.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/create_discussion.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';



@GenerateMocks([WorkspaceDiscussionRepository])
void main() {
  late MockRepo repo;
  late CreateDiscussion usecase;

  setUp(() {
    repo = MockRepo();
    usecase = CreateDiscussion(repo);
  });

  test('success', () async {
    final discussion = Discussion(
      id: '1',
      title: 't',
      content: 'c',
      tags: [],
      category: 'cat',
      commentsCount: 10,
      likeCount: 0,
      reportCount: 0,
      createdAt: DateTime.now(),
      createdBy: User(id: '1', name: 'name'),
    );
    when(
      repo.createDiscussion(
        title: 't',
        content: 'c',
        tags: const [],
        category: 'cat',
      ),
    ).thenAnswer((_) async => Right(discussion));
    final result = await usecase(
      title: 't',
      content: 'c',
      tags: const [],
      category: 'cat',
    );
    expect(result, Right(discussion));
    verify(
      repo.createDiscussion(
        title: 't',
        content: 'c',
        tags: const [],
        category: 'cat',
      ),
    );
    verifyNoMoreInteractions(repo);
  });

  test('failure', () async {
    when(
      repo.createDiscussion(
        title: 't',
        content: 'c',
        tags: const [],
        category: 'cat',
      ),
    ).thenAnswer((_) async => const Left('error'));
    final result = await usecase(
      title: 't',
      content: 'c',
      tags: const [],
      category: 'cat',
    );
    expect(result, const Left('error'));
  });
}
