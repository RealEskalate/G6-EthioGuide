import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/exception.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/core/network/network_info.dart';
import 'package:ethioguide/features/AI%20chat/data/datasources/ai_local_datasource.dart';
import 'package:ethioguide/features/AI%20chat/data/datasources/ai_remote_datasource.dart';
import 'package:ethioguide/features/AI%20chat/data/models/conversation_model.dart';
import 'package:ethioguide/features/AI%20chat/data/repository/ai_repository_impl.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

@GenerateMocks([AiRemoteDatasource, AiLocalDatasource, NetworkInfo])
import 'ai_repository_impl_test.mocks.dart';

void main() {
  late AiRepositoryImpl repository;
  late MockAiRemoteDatasource mockRemoteDataSource;
  late MockAiLocalDatasource mockLocalDataSource;
  late MockNetworkInfo mockNetworkInfo;

  setUp(() {
    mockRemoteDataSource = MockAiRemoteDatasource();
    mockLocalDataSource = MockAiLocalDatasource();
    mockNetworkInfo = MockNetworkInfo();
    repository = AiRepositoryImpl(
      remoteDatasource: mockRemoteDataSource,
      localDatasource: mockLocalDataSource,
      networkInfo: mockNetworkInfo,
    );
  });

  group('sendQuery', () {
    const tQuery = 'How to get a passport?';
    final tConversation = ConversationModel(
      id: 'id',
      request: tQuery,
      response: 'Steps to get a passport...',
      source: 'official',
      procedures: [ProcedureModel(id: '1', name: 'Passport Application')],
    );
    final tCachedHistory = [
      ConversationModel(
        id: 'id',
        request: 'Previous query',
        response: 'Previous response',
        source: 'official',
        procedures: [],
      ),
    ];

    test('should return Conversation when online and successful', () async {
      // Arrange
      when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
      when(
        mockRemoteDataSource.sendQuery(tQuery),
      ).thenAnswer((_) async => tConversation);
      when(
        mockLocalDataSource.getCachedHistory(),
      ).thenAnswer((_) async => tCachedHistory);
      when(
        mockLocalDataSource.cacheHistory(any),
      ).thenAnswer((_) async => Future.value());

      // Act
      final result = await repository.sendQuery(tQuery);

      // Assert
      expect(result, Right(tConversation));
      verify(mockRemoteDataSource.sendQuery(tQuery)).called(1);
      verify(mockLocalDataSource.getCachedHistory()).called(1);
      verify(
        mockLocalDataSource.cacheHistory([tConversation, ...tCachedHistory]),
      ).called(1);
      verifyNoMoreInteractions(mockRemoteDataSource);
      verifyNoMoreInteractions(mockLocalDataSource);
    });

    test('should return Conversation when cache fails', () async {
      // Arrange
      when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
      when(
        mockRemoteDataSource.sendQuery(tQuery),
      ).thenAnswer((_) async => tConversation);
      when(
        mockLocalDataSource.getCachedHistory(),
      ).thenThrow(CacheException(message: 'Cache error'));

      // Act
      final result = await repository.sendQuery(tQuery);

      // Assert
      expect(result, Right(tConversation));
      verify(mockRemoteDataSource.sendQuery(tQuery)).called(1);
      verify(mockLocalDataSource.getCachedHistory()).called(1);
      verifyNever(mockLocalDataSource.cacheHistory(any));
      verifyNoMoreInteractions(mockRemoteDataSource);
      verifyNoMoreInteractions(mockLocalDataSource);
    });

    test('should return NetworkFailure when offline', () async {
      // Arrange
      when(mockNetworkInfo.isConnected).thenAnswer((_) async => false);

      // Act
      final result = await repository.sendQuery(tQuery);

      // Assert
      expect(result, Left(NetworkFailure(message: 'No internet connection')));
      verifyNever(mockRemoteDataSource.sendQuery(any));
      verifyNever(mockLocalDataSource.getCachedHistory());
      verifyNever(mockLocalDataSource.cacheHistory(any));
    });

    test('should return ServerFailure on ServerException', () async {
      // Arrange
      when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
      when(mockRemoteDataSource.sendQuery(tQuery)).thenThrow(
        ServerException(
          message: 'Couldn\'t authenticate user. Please log in again.',
          statusCode: 401,
        ),
      );

      // Act
      final result = await repository.sendQuery(tQuery);

      // Assert
      expect(
        result,
        Left(
          ServerFailure(
            message: 'Couldn\'t authenticate user. Please log in again.',
          ),
        ),
      );
      verify(mockRemoteDataSource.sendQuery(tQuery)).called(1);
      verifyNever(mockLocalDataSource.getCachedHistory());
      verifyNever(mockLocalDataSource.cacheHistory(any));
    });
  }); // send query group

  group('getHistory', () {
    final tHistory = [
      ConversationModel(
        id: 'id',
        request: 'How to get a passport?',
        response: 'Steps to get a passport...',
        source: 'official',
        procedures: [ProcedureModel(id: '1', name: 'Passport Application')],
      ),
    ];

    test(
      'should return List<Conversation> when online and successful',
      () async {
        // Arrange
        when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
        when(
          mockRemoteDataSource.getHistory(),
        ).thenAnswer((_) async => tHistory);
        when(
          mockLocalDataSource.cacheHistory(tHistory),
        ).thenAnswer((_) async => Future.value());

        // Act
        final result = await repository.getHistory();

        // Assert
        expect(result, Right(tHistory));
        verify(mockRemoteDataSource.getHistory()).called(1);
        verify(mockLocalDataSource.cacheHistory(tHistory)).called(1);
        verifyNoMoreInteractions(mockRemoteDataSource);
        verifyNoMoreInteractions(mockLocalDataSource);
      },
    );

    test(
      'should return List<Conversation> when online but cache fails',
      () async {
        // Arrange
        when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
        when(
          mockRemoteDataSource.getHistory(),
        ).thenAnswer((_) async => tHistory);
        when(
          mockLocalDataSource.cacheHistory(tHistory),
        ).thenThrow(CacheException(message: 'Cache error'));

        // Act
        final result = await repository.getHistory();

        // Assert
        expect(result, Right(tHistory));
        verify(mockRemoteDataSource.getHistory()).called(1);
        verify(mockLocalDataSource.cacheHistory(tHistory)).called(1);
        verifyNoMoreInteractions(mockRemoteDataSource);
        verifyNoMoreInteractions(mockLocalDataSource);
      },
    );

    test('should return cached List<Conversation> when offline', () async {
      // Arrange
      when(mockNetworkInfo.isConnected).thenAnswer((_) async => false);
      when(
        mockLocalDataSource.getCachedHistory(),
      ).thenAnswer((_) async => tHistory);

      // Act
      final result = await repository.getHistory();

      // Assert
      expect(result, Right(tHistory));
      verify(mockLocalDataSource.getCachedHistory()).called(1);
      verifyNever(mockRemoteDataSource.getHistory());
      verifyNever(mockLocalDataSource.cacheHistory(any));
    });

    test('should return CachedFailure when offline and cache fails', () async {
      // Arrange
      when(mockNetworkInfo.isConnected).thenAnswer((_) async => false);
      when(
        mockLocalDataSource.getCachedHistory(),
      ).thenThrow(CacheException(message: 'Cache error'));

      // Act
      final result = await repository.getHistory();

      // Assert
      expect(result, Left(CachedFailure(message: 'Cache error')));
      verify(mockLocalDataSource.getCachedHistory()).called(1);
      verifyNever(mockRemoteDataSource.getHistory());
      verifyNever(mockLocalDataSource.cacheHistory(any));
    });

    test('should return cached history when online but remote fails', () async {
      // Arrange
      when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
      when(mockRemoteDataSource.getHistory()).thenThrow(
        ServerException(
          message: 'Internal server error. Please try again later.',
          statusCode: 500,
        ),
      );
      when(
        mockLocalDataSource.getCachedHistory(),
      ).thenAnswer((_) async => tHistory);

      // Act
      final result = await repository.getHistory();

      // Assert
      expect(result, Right(tHistory));
      verify(mockRemoteDataSource.getHistory()).called(1);
      verify(mockLocalDataSource.getCachedHistory()).called(1);
      verifyNever(mockLocalDataSource.cacheHistory(any));
    });

    test(
      'should return ServerFailure when online, remote fails and no cache',
      () async {
        // Arrange
        when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
        when(mockRemoteDataSource.getHistory()).thenThrow(
          ServerException(
            message: 'Internal server error. Please try again later.',
            statusCode: 500,
          ),
        );
        when(
          mockLocalDataSource.getCachedHistory(),
        ).thenThrow(CacheException(message: 'No cache found'));

        // Act
        final result = await repository.getHistory();

        // Assert
        expect(
          result,
          Left(
            ServerFailure(
              message:
                  'Remote failed: ServerException: Internal server error. Please try again later. (Status code: 500)',
            ),
          ),
        );
        verify(mockRemoteDataSource.getHistory()).called(1);
        verify(mockLocalDataSource.getCachedHistory()).called(1);
        verifyNever(mockLocalDataSource.cacheHistory(any));
      },
    );
  }); // get history group

  group('translateContent', () {
    const tContent = 'Hello';
    const tLang = 'am';
    const tTranslated = 'ሰላም';

    test('should return translated string when successful', () async {
      // Arrange
      when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
      when(
        mockRemoteDataSource.translateContent(tContent, tLang),
      ).thenAnswer((_) async => tTranslated);

      // Act
      final result = await repository.translateContent(tContent, tLang);

      // Assert
      expect(result, Right(tTranslated));
      verify(mockRemoteDataSource.translateContent(tContent, tLang)).called(1);
      verifyNoMoreInteractions(mockRemoteDataSource);
      verifyNoMoreInteractions(mockLocalDataSource);
    });

    test('should return NetworkFailure when offline', () async {
      // Arrange
      when(mockNetworkInfo.isConnected).thenAnswer((_) async => false);

      // Act
      final result = await repository.translateContent(tContent, tLang);

      // Assert
      expect(result, Left(NetworkFailure(message: 'No internet connection')));
      verifyNever(mockRemoteDataSource.translateContent(any, any));
      verifyNoMoreInteractions(mockLocalDataSource);
    });

    test('should return ServerFailure on ServerException', () async {
      // Arrange
      when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
      when(mockRemoteDataSource.translateContent(tContent, tLang)).thenThrow(
        ServerException(
          message: 'Too many requests. Please try again later.',
          statusCode: 429,
        ),
      );

      // Act
      final result = await repository.translateContent(tContent, tLang);

      // Assert
      expect(
        result,
        Left(
          ServerFailure(message: 'Too many requests. Please try again later.'),
        ),
      );
      verify(mockRemoteDataSource.translateContent(tContent, tLang)).called(1);
      verifyNoMoreInteractions(mockRemoteDataSource);
      verifyNoMoreInteractions(mockLocalDataSource);
    });
  }); // translate content group
}
