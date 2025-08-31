import 'package:dio/dio.dart';
import 'package:ethioguide/core/error/exception.dart';
import 'package:ethioguide/features/procedure/data/datasources/procedure_remote_data_source.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

import 'procedure_remote_data_source_test.mocks.dart';

@GenerateMocks([Dio])
void main() {
  late MockDio mockDio;
  late ProcedureRemoteDataSourceImpl dataSource;

  setUp(() {
    mockDio = MockDio();
    dataSource = ProcedureRemoteDataSourceImpl(client: mockDio, baseUrl: 'https://api.test');
  });

  group('getProcedures', () {
    test('returns list when statusCode is 200 and maps json correctly', () async {
      final responsePayload = [
        {
          'id': 1,
          'title': 'Passport Renewal',
          'category': 'Travel',
          'duration': '2-3 weeks',
          'cost': '1200 ETB',
          'icon': 'badge',
          'isQuickAccess': true,
          'requiredDocuments': ['Passport Photo'],
          'steps': [
            {'number': 1, 'title': 'Fill', 'description': 'Form'}
          ],
          'resources': [
            {'name': 'Form', 'url': 'https://example.com/form.pdf'}
          ],
          'feedback': [
            {'user': 'Hanna', 'comment': 'Great', 'date': '2025-01-01', 'verified': true}
          ],
        }
      ];

      when(mockDio.get('https://api.test/procedures')).thenAnswer((_) async => Response(
            requestOptions: RequestOptions(path: ''),
            statusCode: 200,
            data: responsePayload,
          ));

      final result = await dataSource.getProcedures();

      expect(result.length, 1);
      final first = result.first;
      expect(first.id, '1');
      expect(first.title, 'Passport Renewal');
      expect(first.category, 'Travel');
      expect(first.duration, '2-3 weeks');
      expect(first.cost, '1200 ETB');
      expect(first.icon, 'badge');
      expect(first.isQuickAccess, true);
      expect(first.requiredDocuments, ['Passport Photo']);
      expect(first.steps.first.number, 1);
      expect(first.resources.first.url, 'https://example.com/form.pdf');
      expect(first.feedback.first.verified, true);
    });

    test('throws ServerException on non-200 status code', () async {
      when(mockDio.get('https://api.test/procedures')).thenAnswer((_) async => Response(
            requestOptions: RequestOptions(path: ''),
            statusCode: 500,
            data: {'error': 'server error'},
          ));

      expect(
        () => dataSource.getProcedures(),
        throwsA(isA<ServerException>()),
      );
    });

    test('maps DioException to ServerException', () async {
      when(mockDio.get('https://api.test/procedures')).thenThrow(DioException(
        requestOptions: RequestOptions(path: ''),
        type: DioExceptionType.connectionError,
        error: 'No internet',
      ));

      expect(
        () => dataSource.getProcedures(),
        throwsA(isA<ServerException>()),
      );
    });
  });
}


