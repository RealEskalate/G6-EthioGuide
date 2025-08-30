import 'dart:convert';

import 'package:ethioguide/core/config/cache_key_names.dart';
import 'package:ethioguide/core/error/exception.dart';
import 'package:ethioguide/features/AI%20chat/data/datasources/ai_local_datasource.dart';
import 'package:ethioguide/features/AI%20chat/data/models/conversation_model.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

@GenerateMocks([FlutterSecureStorage])
import 'ai_local_datasource_test.mocks.dart';

void main() {
  late MockFlutterSecureStorage mockSecureStorage;
  late AiLocalDataSourceImpl dataSource;

  setUp(() {
    mockSecureStorage = MockFlutterSecureStorage();
    dataSource = AiLocalDataSourceImpl(secureStorage: mockSecureStorage);
  });

  group('getCachedHistory', () {
    final tHistory = [
      ConversationModel(
        request: 'How to get a passport?',
        response: 'Steps to get a passport...',
        source: 'official',
        procedures: [ProcedureModel(id: '1', name: 'Passport Application')],
      ),
    ];
    final tHistoryJson = jsonEncode(tHistory.map((e) => e.toJson()).toList());

    test(
      'should return List<ConversationModel> when cache contains valid data',
      () async {
        // Arrange
        when(
          mockSecureStorage.read(key: CacheKeyNames.aiHistoryKey),
        ).thenAnswer((_) async => tHistoryJson);

        // Act
        final result = await dataSource.getCachedHistory();

        // Assert
        expect(result, tHistory);
        verify(
          mockSecureStorage.read(key: CacheKeyNames.aiHistoryKey),
        ).called(1);
        verifyNoMoreInteractions(mockSecureStorage);
      },
    );

    test('should return empty list when cache is empty', () async {
      // Arrange
      when(
        mockSecureStorage.read(key: CacheKeyNames.aiHistoryKey),
      ).thenAnswer((_) async => null);

      // Act
      final result = await dataSource.getCachedHistory();

      // Assert
      expect(result, []);
      verify(mockSecureStorage.read(key: CacheKeyNames.aiHistoryKey)).called(1);
      verifyNoMoreInteractions(mockSecureStorage);
    });

    test(
      'should throw CacheException when secure storage read fails',
      () async {
        // Arrange
        when(
          mockSecureStorage.read(key: CacheKeyNames.aiHistoryKey),
        ).thenThrow(Exception('Storage error'));

        // Act & Assert
        await expectLater(
          () => dataSource.getCachedHistory(),
          throwsA(
            isA<CacheException>().having(
              (e) => e.message,
              'message',
              'Failed to read cached history: Exception: Storage error',
            ),
          ),
        );
        verify(
          mockSecureStorage.read(key: CacheKeyNames.aiHistoryKey),
        ).called(1);
        verifyNoMoreInteractions(mockSecureStorage);
      },
    );
  }); // get chached history group

  group('cacheHistory', () {
    final tHistory = [
      ConversationModel(
        request: 'How to get a passport?',
        response: 'Steps to get a passport...',
        source: 'official',
        procedures: [ProcedureModel(id: '1', name: 'Passport Application')],
      ),
    ];
    final tHistoryJson = jsonEncode(tHistory.map((e) => e.toJson()).toList());

    test(
      'should write history to secure storage when cache is successful',
      () async {
        // Arrange
        when(
          mockSecureStorage.write(
            key: anyNamed('key'),
            value: anyNamed('value'),
          ),
        ).thenAnswer((_) async => Future.value());

        // Act
        await dataSource.cacheHistory(tHistory);

        // Assert
        verify(
          mockSecureStorage.write(
            key: CacheKeyNames.aiHistoryKey,
            value: tHistoryJson,
          ),
        ).called(1);
        verifyNoMoreInteractions(mockSecureStorage);
      },
    );

    test(
      'should throw CacheException when secure storage write fails',
      () async {
        // Arrange
        when(
          mockSecureStorage.write(
            key: anyNamed('key'),
            value: anyNamed('value'),
          ),
        ).thenThrow(Exception('Storage error'));

        // Act & Assert
        await expectLater(
          () => dataSource.cacheHistory(tHistory),
          throwsA(
            isA<CacheException>().having(
              (e) => e.message,
              'message',
              'Failed to cache hisory: Exception: Storage error',
            ),
          ),
        );
        verify(
          mockSecureStorage.write(
            key: CacheKeyNames.aiHistoryKey,
            value: tHistoryJson,
          ),
        ).called(1);
        verifyNoMoreInteractions(mockSecureStorage);
      },
    );
  }); // cache history group
}
