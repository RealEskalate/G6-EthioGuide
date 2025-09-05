import 'package:dio/dio.dart';
import 'package:ethioguide/core/error/exception.dart';
import 'package:ethioguide/core/network/network_info.dart';
import 'package:ethioguide/features/AI%20chat/data/datasources/ai_remote_datasource.dart';
import 'package:ethioguide/features/AI%20chat/data/models/conversation_model.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

@GenerateMocks([Dio, NetworkInfo])
import 'ai_remote_datasource_test.mocks.dart';

void main() {
  late AiRemoteDataSourceImpl dataSource;
  late MockDio mockDio;
  late MockNetworkInfo mockNetworkInfo;

  setUp(() {
    mockDio = MockDio();
    mockNetworkInfo = MockNetworkInfo();
    dataSource = AiRemoteDataSourceImpl(
      dio: mockDio,
      networkInfo: mockNetworkInfo,
    );

    when(mockDio.options).thenReturn(BaseOptions());
  });

  group('sendQuery', () {
    const tQuery = 'How to get a passport?';
    final tConversationModel = ConversationModel(
      id: 'id',
      request: tQuery,
      response: 'Steps to get a passport...',
      source: 'official',
      procedures: [ProcedureModel(id: '1', name: 'Passport Application')],
    );
    final tResponseData = {
      'id': 'id',
      'request': tQuery,
      'response': 'Steps to get a passport...',
      'source': 'official',
      'procedures': [
        {'id': '1', 'name': 'Passport Application'},
      ],
    };

    test(
      'should return ConversationModel when the request is successful',
      () async {
        // Arrange
        when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
        when(mockDio.post(any, data: anyNamed('data'))).thenAnswer(
          (_) async => Response(
            data: tResponseData,
            statusCode: 200,
            requestOptions: RequestOptions(path: '/ai/guide'),
          ),
        );

        // Act
        final result = await dataSource.sendQuery(tQuery);

        // Assert
        expect(result, tConversationModel);
        verify(mockDio.post('/ai/guide', data: {'query': tQuery})).called(1);
        verifyNoMoreInteractions(mockDio);
      },
    );

    test('should throw ServerException when network is offline', () async {
      // Arrange
      when(mockNetworkInfo.isConnected).thenAnswer((_) async => false);

      // Act & Assert
      expect(
        () => dataSource.sendQuery(tQuery),
        throwsA(
          isA<ServerException>()
              .having((e) => e.message, 'message', 'No internet connection')
              .having((e) => e.statusCode, 'statusCode', null),
        ),
      );
      verifyNever(mockDio.post(any, data: anyNamed('data')));
    });

    test('should throw ServerException for 401 status code', () async {
      // Arrange
      when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
      when(mockDio.post(any, data: anyNamed('data'))).thenAnswer(
        (_) async => Response(
          data: {'message': 'Unauthorized'},
          statusCode: 401,
          requestOptions: RequestOptions(path: '/ai/guide'),
        ),
      );

      // Act & Assert
      await expectLater(
        () => dataSource.sendQuery(tQuery),
        throwsA(
          isA<ServerException>()
              .having(
                (e) => e.message,
                'message',
                'Couldn\'t authenticate user. Please log in again.',
              )
              .having((e) => e.statusCode, 'statusCode', 401),
        ),
      );
      verify(mockDio.post('/ai/guide', data: {'query': tQuery})).called(1);
    });

    test('should throw ServerException for unexpected DioException', () async {
      // Arrange
      when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
      when(mockDio.post(any, data: anyNamed('data'))).thenThrow(
        DioException(
          requestOptions: RequestOptions(path: '/ai/guide'),
          error: 'Network error',
        ),
      );

      // Act & Assert
      await expectLater(
        () => dataSource.sendQuery(tQuery),
        throwsA(
          isA<ServerException>()
              .having((e) => e.message, 'message', 'Failed to send query')
              .having((e) => e.statusCode, 'statusCode', null),
        ),
      );
      verify(mockDio.post('/ai/guide', data: {'query': tQuery})).called(1);
    });
  }); // send query group

  group('getHistory', () {
    final tHistory = [
      ConversationModel(
        id: 'id',
        request: 'How to renew passport?',
        response: 'Here are the verified steps...',
        source: 'official',
        procedures: [
          ProcedureModel(id: 'id', name: 'name'),
          ProcedureModel(id: 'id2', name: 'name2'),
        ],
      ),
    ];
    final tResponseData = {
      "history": [
        {
          "id": "id",
          "procedures": [
            {"id": "id", "name": "name"},
            {"id": "id2", "name": "name2"},
          ],
          "request": "How to renew passport?",
          "response": "Here are the verified steps...",
          "source": "official",
        },
      ],
    };

    test(
      'should return List<ConversationModel> when the request is successful',
      () async {
        // Arrange
        when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
        when(mockDio.get(any)).thenAnswer(
          (_) async => Response(
            data: tResponseData,
            statusCode: 200,
            requestOptions: RequestOptions(path: '/ai/history'),
          ),
        );

        // Act
        final result = await dataSource.getHistory();

        // Assert
        expect(result, tHistory);
        verify(mockDio.get('/ai/history')).called(1);
        verifyNoMoreInteractions(mockDio);
      },
    );

    test('should throw ServerException when network is offline', () async {
      // Arrange
      when(mockNetworkInfo.isConnected).thenAnswer((_) async => false);

      // Act & Assert
      await expectLater(
        () => dataSource.getHistory(),
        throwsA(
          isA<ServerException>()
              .having((e) => e.message, 'message', 'No internet connection')
              .having((e) => e.statusCode, 'statusCode', null),
        ),
      );
      verifyNever(mockDio.get(any));
    });

    test('should throw ServerException for 500 status code', () async {
      // Arrange
      when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
      when(mockDio.get(any)).thenAnswer(
        (_) async => Response(
          data: {'message': 'Server error'},
          statusCode: 500,
          requestOptions: RequestOptions(path: '/ai/history'),
        ),
      );

      // Act & Assert
      await expectLater(
        () => dataSource.getHistory(),
        throwsA(
          isA<ServerException>()
              .having(
                (e) => e.message,
                'message',
                'Internal server error. Please try again later.',
              )
              .having((e) => e.statusCode, 'statusCode', 500),
        ),
      );
      verify(mockDio.get('/ai/history')).called(1);
    });

    test('should throw ServerException for unexpected DioException', () async {
      // Arrange
      when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
      when(mockDio.get(any)).thenThrow(
        DioException(
          requestOptions: RequestOptions(path: '/ai/history'),
          error: 'Network error',
        ),
      );

      // Act & Assert
      await expectLater(
        () => dataSource.getHistory(),
        throwsA(
          isA<ServerException>()
              .having((e) => e.message, 'message', 'Failed to fetch history')
              .having((e) => e.statusCode, 'statusCode', null),
        ),
      );
      verify(mockDio.get('/ai/history')).called(1);
    });
  }); // get history group

  group('translateContent', () {
    const tContent = 'Hello';
    const tLang = 'am';
    const tTranslated = 'ሰላም';

    test(
      'should return translated string when the request is successful',
      () async {
        // Arrange
        when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
        when(mockDio.post(any, data: anyNamed('data'))).thenAnswer(
          (_) async => Response(
            data: {'translated': tTranslated},
            statusCode: 200,
            requestOptions: RequestOptions(path: '/translate'),
          ),
        );

        // Act
        final result = await dataSource.translateContent(tContent, tLang);

        // Assert
        expect(result, tTranslated);
        verify(
          mockDio.post(
            '/translate',
            data: {'content': tContent, 'lang': tLang},
          ),
        ).called(1);
        verifyNoMoreInteractions(mockDio);
      },
    );

    test('should throw ServerException when network is offline', () async {
      // Arrange
      when(mockNetworkInfo.isConnected).thenAnswer((_) async => false);

      // Act & Assert
      await expectLater(
        () => dataSource.translateContent(tContent, tLang),
        throwsA(
          isA<ServerException>()
              .having((e) => e.message, 'message', 'No internet connection')
              .having((e) => e.statusCode, 'statusCode', null),
        ),
      );
      verifyNever(mockDio.post(any, data: anyNamed('data')));
    });

    test('should throw ServerException for 429 status code', () async {
      // Arrange
      when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
      when(mockDio.post(any, data: anyNamed('data'))).thenAnswer(
        (_) async => Response(
          data: {'message': 'Rate limit exceeded'},
          statusCode: 429,
          requestOptions: RequestOptions(path: '/translate'),
        ),
      );

      // Act & Assert
      await expectLater(
        () => dataSource.translateContent(tContent, tLang),
        throwsA(
          isA<ServerException>()
              .having(
                (e) => e.message,
                'message',
                'Too many requests. Please try again later.',
              )
              .having((e) => e.statusCode, 'statusCode', 429),
        ),
      );
      verify(
        mockDio.post('/translate', data: {'content': tContent, 'lang': tLang}),
      ).called(1);
    });

    test('should throw ServerException for unexpected DioException', () async {
      // Arrange
      when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
      when(mockDio.post(any, data: anyNamed('data'))).thenThrow(
        DioException(
          requestOptions: RequestOptions(path: '/translate'),
          error: 'Network error',
        ),
      );

      // Act & Assert
      await expectLater(
        () => dataSource.translateContent(tContent, tLang),
        throwsA(
          isA<ServerException>()
              .having((e) => e.message, 'message', 'Failed to translate')
              .having((e) => e.statusCode, 'statusCode', null),
        ),
      );
      verify(
        mockDio.post('/translate', data: {'content': tContent, 'lang': tLang}),
      ).called(1);
    });
  }); // translate content group
}
