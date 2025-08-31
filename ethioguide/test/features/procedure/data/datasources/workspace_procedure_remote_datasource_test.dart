import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';
import 'package:ethioguide/features/procedure/data/datasources/workspace_procedure_remote_datasource.dart';
import 'package:ethioguide/features/procedure/data/models/workspace_procedure_model.dart';
import 'package:ethioguide/features/procedure/data/models/workspace_summary_model.dart';
import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';

import '../procedure_remote_data_source_test.mocks.dart';
import 'workspace_procedure_remote_datasource_test.mocks.dart';

@GenerateMocks([Dio])
void main() {
  late WorkspaceProcedureRemoteDataSourceImpl dataSource;
  late MockDio mockDio;

  setUp(() {
    mockDio = MockDio();
    dataSource = WorkspaceProcedureRemoteDataSourceImpl(dio: mockDio);
  });

  group('getWorkspaceProcedures', () {
    final tProcedures = [
       WorkspaceProcedureModel(
        id: '1',
        title: 'Test Procedure',
        organization: 'Test Org',
        status: ProcedureStatus.inProgress,
        progressPercentage: 50,
        documentsUploaded: 2,
        totalDocuments: 4,
        startDate: DateTime(2024, 1, 1),
      ),
    ];

    test('should return procedures when API call is successful', () async {
      // arrange
      when(mockDio.get(any, queryParameters: anyNamed('queryParameters')))
          .thenAnswer((_) async => Response(
                data: {'data': [
                  {
                    'id': '1',
                    'title': 'Test Procedure',
                    'organization': 'Test Org',
                    'status': 'In Progress',
                    'progressPercentage': 50,
                    'documentsUploaded': 2,
                    'totalDocuments': 4,
                    'startDate': '2024-01-01T00:00:00.000Z',
                  }
                ]},
                statusCode: 200,
                requestOptions: RequestOptions(path: ''),
              ));

      // act
      final result = await dataSource.getWorkspaceProcedures();

      // assert
      expect(result, equals(tProcedures));
      verify(mockDio.get('https://api.ethioguide.com/api/v1/workspace/procedures'));
    });

    test('should throw exception when API call fails', () async {
      // arrange
      when(mockDio.get(any, queryParameters: anyNamed('queryParameters')))
          .thenAnswer((_) async => Response(
                data: {'error': 'Server error'},
                statusCode: 500,
                requestOptions: RequestOptions(path: ''),
              ));

      // act & assert
      expect(
        () => dataSource.getWorkspaceProcedures(),
        throwsA(isA<Exception>()),
      );
    });

    test('should throw exception when network error occurs', () async {
      // arrange
      when(mockDio.get(any, queryParameters: anyNamed('queryParameters')))
          .thenThrow(DioException(
            requestOptions: RequestOptions(path: ''),
            error: 'Network error',
          ));

      // act & assert
      expect(
        () => dataSource.getWorkspaceProcedures(),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('getWorkspaceSummary', () {
    final tSummary =  WorkspaceSummaryModel(
      totalProcedures: 5,
      inProgress: 2,
      completed: 3,
      totalDocuments: 15,
    );

    test('should return summary when API call is successful', () async {
      // arrange
      when(mockDio.get(any, queryParameters: anyNamed('queryParameters')))
          .thenAnswer((_) async => Response(
                data: {'data': {
                  'totalProcedures': 5,
                  'inProgress': 2,
                  'completed': 3,
                  'totalDocuments': 15,
                }},
                statusCode: 200,
                requestOptions: RequestOptions(path: ''),
              ));

      // act
      final result = await dataSource.getWorkspaceSummary();

      // assert
      expect(result, equals(tSummary));
      verify(mockDio.get('https://api.ethioguide.com/api/v1/workspace/summary'));
    });
  });

  group('getProceduresByStatus', () {
    test('should return procedures filtered by status', () async {
      // arrange
      when(mockDio.get(any, queryParameters: anyNamed('queryParameters')))
          .thenAnswer((_) async => Response(
                data: {'data': []},
                statusCode: 200,
                requestOptions: RequestOptions(path: ''),
              ));

      // act
      await dataSource.getProceduresByStatus('In Progress');

      // assert
      verify(mockDio.get(
        'https://api.ethioguide.com/api/v1/workspace/procedures',
        queryParameters: {'status': 'In Progress'},
      ));
    });
  });

  group('getProceduresByOrganization', () {
    test('should return procedures filtered by organization', () async {
      // arrange
      when(mockDio.get(any, queryParameters: anyNamed('queryParameters')))
          .thenAnswer((_) async => Response(
                data: {'data': []},
                statusCode: 200,
                requestOptions: RequestOptions(path: ''),
              ));

      // act
      await dataSource.getProceduresByOrganization('Test Org');

      // assert
      verify(mockDio.get(
        'https://api.ethioguide.com/api/v1/workspace/procedures',
        queryParameters: {'organization': 'Test Org'},
      ));
    });
  });

  group('createWorkspaceProcedure', () {
    final tProcedure =  WorkspaceProcedureModel(
      id: '1',
      title: 'Test Procedure',
      organization: 'Test Org',
      status: ProcedureStatus.inProgress,
      progressPercentage: 50,
      documentsUploaded: 2,
      totalDocuments: 4,
      startDate: DateTime(2024, 1, 1),
    );

    test('should return created procedure when API call is successful', () async {
      // arrange
      when(mockDio.post(any, data: anyNamed('data')))
          .thenAnswer((_) async => Response(
                data: {'data': {
                  'id': '1',
                  'title': 'Test Procedure',
                  'organization': 'Test Org',
                  'status': 'In Progress',
                  'progressPercentage': 50,
                  'documentsUploaded': 2,
                  'totalDocuments': 4,
                  'startDate': '2024-01-01T00:00:00.000Z',
                }},
                statusCode: 201,
                requestOptions: RequestOptions(path: ''),
              ));

      // act
      final result = await dataSource.createWorkspaceProcedure(tProcedure);

      // assert
      expect(result, equals(tProcedure));
      verify(mockDio.post(
        'https://api.ethioguide.com/api/v1/workspace/procedures',
        data: tProcedure.toJson(),
      ));
    });
  });

  group('updateWorkspaceProcedure', () {
    final tProcedure =  WorkspaceProcedureModel(
      id: '1',
      title: 'Updated Procedure',
      organization: 'Test Org',
      status: ProcedureStatus.completed,
      progressPercentage: 100,
      documentsUploaded: 4,
      totalDocuments: 4,
      startDate: DateTime(2024, 1, 1),
    );

    test('should return updated procedure when API call is successful', () async {
      // arrange
      when(mockDio.put(any, data: anyNamed('data')))
          .thenAnswer((_) async => Response(
                data: {'data': tProcedure.toJson()},
                statusCode: 200,
                requestOptions: RequestOptions(path: ''),
              ));

      // act
      final result = await dataSource.updateWorkspaceProcedure(tProcedure);

      // assert
      expect(result, equals(tProcedure));
      verify(mockDio.put(
        'https://api.ethioguide.com/api/v1/workspace/procedures/1',
        data: tProcedure.toJson(),
      ));
    });
  });

  group('deleteWorkspaceProcedure', () {
    test('should return true when deletion is successful', () async {
      // arrange
      when(mockDio.delete(any))
          .thenAnswer((_) async => Response(
                data: {'success': true},
                statusCode: 200,
                requestOptions: RequestOptions(path: ''),
              ));

      // act
      final result = await dataSource.deleteWorkspaceProcedure('1');

      // assert
      expect(result, isTrue);
      verify(mockDio.delete('https://api.ethioguide.com/api/v1/workspace/procedures/1'));
    });
  });

  group('updateProgress', () {
    test('should return updated procedure when progress update is successful', () async {
      // arrange
      when(mockDio.patch(any, data: anyNamed('data')))
          .thenAnswer((_) async => Response(
                data: {'data': {
                  'id': '1',
                  'title': 'Test Procedure',
                  'organization': 'Test Org',
                  'status': 'In Progress',
                  'progressPercentage': 75,
                  'documentsUploaded': 3,
                  'totalDocuments': 4,
                  'startDate': '2024-01-01T00:00:00.000Z',
                }},
                statusCode: 200,
                requestOptions: RequestOptions(path: ''),
              ));

      // act
      final result = await dataSource.updateProgress('1', 75);

      // assert
      expect(result.progressPercentage, equals(75));
      verify(mockDio.patch(
        'https://api.ethioguide.com/api/v1/workspace/procedures/1/progress',
        data: {'progressPercentage': 75},
      ));
    });
  });
}


