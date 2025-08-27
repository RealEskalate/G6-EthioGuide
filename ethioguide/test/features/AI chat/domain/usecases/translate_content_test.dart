import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/AI%20chat/Domain/usecases/translate_content.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';

import 'send_query_test.mocks.dart'; // common mocks for all usecases (only repository)

void main() {
  late TranslateContent usecase;
  late MockAiRepository mockAiRepository;

  setUp(() {
    mockAiRepository = MockAiRepository();
    usecase = TranslateContent(repository: mockAiRepository);
  });

  final String tContent = 'test content';
  final String tLang = 'amh';
  final String tResponse = 'mukera';
  test(
    'should return a translation of the content from repository if successful',
    () async {
      // Arrange
      when(
        mockAiRepository.translateContent(tContent, tLang),
      ).thenAnswer((_) async => Right(tResponse));

      // Act
      final result = await usecase(content: tContent, lang: tLang);

      // Assert
      expect(result, Right(tResponse));
      verify(mockAiRepository.translateContent(tContent, tLang));
      verifyNoMoreInteractions(mockAiRepository);
    },
  );

  test('should get a failure from repository if uncessful', () async {
    // Arrange
    when(
      mockAiRepository.translateContent(tContent, tLang),
    ).thenAnswer((_) async => Left(Failure()));

    // Act
    final result = await usecase(content: tContent, lang: tLang);

    // Assert
    expect(result, Left(Failure()));
    verify(mockAiRepository.translateContent(tContent, tLang));
    verifyNoMoreInteractions(mockAiRepository);
  });
}
