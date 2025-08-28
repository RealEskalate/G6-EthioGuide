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
  });

  group('sendQuery', () {
    const tQuery = 'How to get a passport';
    final tConversationJson = {
      'title': 'Passport Guide',
      'content': 'Steps to get a Passport...',
      'source': 'official',
      'procedures': [
        {'id': 1, 'name': 'Passport Application'},
      ],
    };
    final tConversationModel = ConversationModel.fromJson(tConversationJson);

    test(
      'should return ConversationModel when the call is successful',
      () async {
        // Arrange
        when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
        when(mockDio.post('/ai/guide', data: anyNamed('data'))).thenAnswer(
          (_) async => Response(
            requestOptions: RequestOptions(path: '/ai/guide'),
            data: tConversationJson,
            statusCode: 200,
          ),
        );

        // Act
        final result = await dataSource.sendQuery(tQuery);

        // Assert
        expect(result, equals(tConversationModel));
        verify(mockNetworkInfo.isConnected);
        verify(mockDio.post('/ai/guide', data: {'query': tQuery}));
      },
    );

    test('should throw ServerException when offline', () async {
      // arrange
      when(mockNetworkInfo.isConnected).thenAnswer((_) async => false);

      // act
      final call = dataSource.sendQuery;

      // assert
      expect(() => call(tQuery), throwsA(isA<ServerException>()));
      verify(mockNetworkInfo.isConnected);
      verifyZeroInteractions(mockDio);
    });
  }); // send query group
}
