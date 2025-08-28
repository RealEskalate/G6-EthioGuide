import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/AI%20chat/Domain/entities/conversation.dart';
import 'package:ethioguide/features/AI%20chat/Domain/usecases/get_history.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';

import 'send_query_test.mocks.dart'; // common mocks for all usecases (only repository)

void main() {
  late GetHistory usecase;
  late MockAiRepository mockAiRepository;

  setUp(() {
    mockAiRepository = MockAiRepository();
    usecase = GetHistory(repository: mockAiRepository);
  });

  final String tQuery = 'How to get a passport';
  final Conversation tConversation = Conversation(
    request: tQuery,
    response: 'Steps to get a Passport...',
    source: 'official',
    procedures: [Procedure(id: '1', name: 'Passport Application')],
  );

  final List<Conversation> tconversations = [tConversation];

  test(
    'should return list of Conversation from repository if successful',
    () async {
      // Arrange
      when(
        mockAiRepository.getHistory(),
      ).thenAnswer((_) async => Right(tconversations));

      // Act
      final result = await usecase();

      // Assert
      expect(result, Right(tconversations));
      verify(mockAiRepository.getHistory());
      verifyNoMoreInteractions(mockAiRepository);
    },
  );

  test(
    'should get a failure from repository if uncessful',
    () async {
      // Arrange
      when(
        mockAiRepository.getHistory(),
      ).thenAnswer((_) async => Left(Failure()));

      // Act
      final result = await usecase();

      // Assert
      expect(result, Left(Failure()));
      verify(mockAiRepository.getHistory());
      verifyNoMoreInteractions(mockAiRepository);
    },
  );
}
