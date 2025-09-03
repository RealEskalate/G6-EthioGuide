import 'package:bloc_test/bloc_test.dart';
import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/AI%20chat/Domain/usecases/get_history.dart';
import 'package:ethioguide/features/AI%20chat/Domain/usecases/send_query.dart';
import 'package:ethioguide/features/AI%20chat/Domain/usecases/translate_content.dart';
import 'package:ethioguide/features/AI%20chat/Presentation/bloc/ai_bloc.dart';
import 'package:ethioguide/features/AI%20chat/data/models/conversation_model.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

@GenerateMocks([SendQuery, GetHistory, TranslateContent])
import 'ai_bloc_test.mocks.dart';

void main() {
  late AiBloc bloc;
  late MockSendQuery mockSendQueryUseCase;
  late MockGetHistory mockGetHistoryUseCase;
  late MockTranslateContent mockTranslateContentUseCase;

  setUp(() {
    mockSendQueryUseCase = MockSendQuery();
    mockGetHistoryUseCase = MockGetHistory();
    mockTranslateContentUseCase = MockTranslateContent();
    bloc = AiBloc(
      sendQueryUseCase: mockSendQueryUseCase,
      getHistoryUseCase: mockGetHistoryUseCase,
      translateContentUseCase: mockTranslateContentUseCase,
    );
  });

  tearDown(() {
    bloc.close();
  });

  group('SendQueryEvent', () {
    const tQuery = 'How to get a passport?';
    final tConversation = ConversationModel(
      id: 'id',
      request: tQuery,
      response: 'Steps to get a passport...',
      source: 'official',
      procedures: [ProcedureModel(id: '1', name: 'Passport Application')],
    );

    blocTest<AiBloc, AiState>(
      'emits [AiLoading, AiQuerySuccess] when query is successful',
      build: () {
        when(
          mockSendQueryUseCase(tQuery),
        ).thenAnswer((_) async => Right(tConversation));
        return bloc;
      },
      act: (bloc) => bloc.add(const SendQueryEvent(query: tQuery)),
      expect: () => [AiLoading(), AiQuerySuccess(conversation: tConversation)],
      verify: (_) {
        verify(mockSendQueryUseCase(tQuery)).called(1);
      },
    );

    blocTest<AiBloc, AiState>(
      'emits [AiLoading, AiError] when network fails',
      build: () {
        when(mockSendQueryUseCase(tQuery)).thenAnswer(
          (_) async => Left(NetworkFailure(message: 'No internet connection')),
        );
        return bloc;
      },
      act: (bloc) => bloc.add(const SendQueryEvent(query: tQuery)),
      expect: () => [
        AiLoading(),
        const AiError('No internet connection. Please check your network.'),
      ],
      verify: (_) {
        verify(mockSendQueryUseCase(tQuery)).called(1);
      },
    );

    blocTest<AiBloc, AiState>(
      'emits [AiLoading, AiError] when server fails with 401',
      build: () {
        when(mockSendQueryUseCase(tQuery)).thenAnswer(
          (_) async => Left(
            ServerFailure(
              message: 'Couldn\'t authenticate user. Please log in again.',
            ),
          ),
        );
        return bloc;
      },
      act: (bloc) => bloc.add(const SendQueryEvent(query: tQuery)),
      expect: () => [
        AiLoading(),
        const AiError('Couldn\'t authenticate user. Please log in again.'),
      ],
      verify: (_) {
        verify(mockSendQueryUseCase(tQuery)).called(1);
      },
    );
  }); // send query event grouop

  group('GetHistoryEvent', () {
    final tHistory = [
      ConversationModel(
        id: 'id',
        request: 'How to get a passport?',
        response: 'Steps to get a passport...',
        source: 'official',
        procedures: [ProcedureModel(id: '1', name: 'Passport Application')],
      ),
    ];

    blocTest<AiBloc, AiState>(
      'emits [AiLoading, AiHistorySuccess] when history is successful',
      build: () {
        when(mockGetHistoryUseCase()).thenAnswer((_) async => Right(tHistory));
        return bloc;
      },
      act: (bloc) => bloc.add(GetHistoryEvent()),
      expect: () => [AiLoading(), AiHistorySuccess(history: tHistory)],
      verify: (_) {
        verify(mockGetHistoryUseCase()).called(1);
      },
    );

    blocTest<AiBloc, AiState>(
      'emits [AiLoading, AiError] when cache fails',
      build: () {
        when(
          mockGetHistoryUseCase(),
        ).thenAnswer((_) async => Left(CachedFailure(message: 'Cache error')));
        return bloc;
      },
      act: (bloc) => bloc.add(GetHistoryEvent()),
      expect: () => [
        AiLoading(),
        const AiError('Unable to load cached data. Please try again.'),
      ],
      verify: (_) {
        verify(mockGetHistoryUseCase()).called(1);
      },
    );
  }); // get history event group

  group('TranslateContentEvent', () {
    const tContent = 'Hello';
    const tLang = 'am';
    const tTranslated = 'ሰላም';

    blocTest<AiBloc, AiState>(
      'emits [AiLoading, AiTranslateSuccess] when translation is successful',
      build: () {
        when(
          mockTranslateContentUseCase(content: tContent, lang: tLang),
        ).thenAnswer((_) async => Right(tTranslated));
        return bloc;
      },
      act: (bloc) =>
          bloc.add(const TranslateContentEvent(content: tContent, lang: tLang)),
      expect: () => [
        AiLoading(),
        const AiTranslateSuccess(translated: tTranslated),
      ],
      verify: (_) {
        verify(
          mockTranslateContentUseCase(content: tContent, lang: tLang),
        ).called(1);
      },
    );

    blocTest<AiBloc, AiState>(
      'emits [AiLoading, AiError] when server fails with 429',
      build: () {
        when(
          mockTranslateContentUseCase(content: tContent, lang: tLang),
        ).thenAnswer(
          (_) async => Left(
            ServerFailure(
              message: 'Too many requests. Please try again later.',
            ),
          ),
        );
        return bloc;
      },
      act: (bloc) =>
          bloc.add(const TranslateContentEvent(content: tContent, lang: tLang)),
      expect: () => [
        AiLoading(),
        const AiError('Too many requests. Please try again later.'),
      ],
      verify: (_) {
        verify(
          mockTranslateContentUseCase(content: tContent, lang: tLang),
        ).called(1);
      },
    );
  }); // translate content event group
}
