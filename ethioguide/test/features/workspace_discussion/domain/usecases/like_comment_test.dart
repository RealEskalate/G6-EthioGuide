import 'package:dartz/dartz.dart';
import 'package:ethioguide/features/workspace_discussion/domain/repositories/workspace_discussion_repository.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/like_comment.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';



@GenerateMocks([WorkspaceDiscussionRepository])
void main() {
  late MockRepo repo;
  late LikeComment usecase;

  setUp(() {
    repo = MockRepo();
    usecase = LikeComment(repo);
  });

  test('success', () async {
    when(repo.likeComment('1')).thenAnswer((_) async => const Right(true));
    final result = await usecase('1');
    expect(result, const Right(true));
    verify(repo.likeComment('1'));
    verifyNoMoreInteractions(repo);
  });

  test('failure', () async {
    when(repo.likeComment('1')).thenAnswer((_) async => const Left('error'));
    final result = await usecase('1');
    expect(result, const Left('error'));
  });
}


