import 'package:dartz/dartz.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/comment.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/user.dart';
import 'package:ethioguide/features/workspace_discussion/domain/repositories/workspace_discussion_repository.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/add_comment.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';




@GenerateMocks([WorkspaceDiscussionRepository])
void main() {
  late MockRepo repo;
  late AddComment usecase;

  setUp(() {
    repo = MockRepo();
    usecase = AddComment(repo);
  });

  test('success', () async {
    final comment =  Comment(id:'1', discussionId: '1', content: 'c', createdAt: DateTime.now(), createdBy: User(id: '1', name: 'name'), likeCount: 0, reportCount: 0, );
    when(repo.addComment(discussionId: '1', content: 'c')).thenAnswer((_) async => Right(comment));
    final result = await usecase(discussionId: '1', content: 'c');
    expect(result, Right(comment));
    verify(repo.addComment(discussionId: '1', content: 'c'));
    verifyNoMoreInteractions(repo);
  });

  test('failure', () async {
    when(repo.addComment(discussionId: '1', content: 'c')).thenAnswer((_) async => const Left('error'));
    final result = await usecase(discussionId: '1', content: 'c');
    expect(result, const Left('error'));
  });
}


