import 'package:dartz/dartz.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/comment.dart';
import 'package:ethioguide/features/workspace_discussion/domain/repositories/workspace_discussion_repository.dart';
import 'package:ethioguide/features/workspace_discussion/domain/usecases/get_comments.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';




@GenerateMocks([WorkspaceDiscussionRepository])
void main() {
  late MockRepo repo;
  late GetComments usecase;

  setUp(() {
    repo = MockRepo();
    usecase = GetComments(repo);
  });

  test('success', () async {
    when(repo.getComments('1')).thenAnswer((_) async => const Right(<Comment>[]));
    final result = await usecase('1');
    expect(result, const Right(<Comment>[]));
    verify(repo.getComments('1'));
    verifyNoMoreInteractions(repo);
  });

  test('failure', () async {
    when(repo.getComments('1')).thenAnswer((_) async => const Left('error'));
    final result = await usecase('1');
    expect(result, const Left('error'));
  });
}


