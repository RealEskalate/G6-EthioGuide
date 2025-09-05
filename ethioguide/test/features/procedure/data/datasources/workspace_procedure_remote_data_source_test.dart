import 'package:dio/dio.dart';
import 'package:ethioguide/features/procedure/data/datasources/workspace_procedure_remote_data_source.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

import 'workspace_procedure_remote_data_source_test.mocks.dart';




@GenerateMocks([Dio])
void main() {
  late MockDio dio;
  late WorkspaceProcedureRemoteDataSourceImpl ds;

  setUp(() {
    dio = MockDio();
    ds = WorkspaceProcedureRemoteDataSourceImpl(dio: dio, baseUrl: 'https://api');
  });

  test('getWorkspaceSummary success', () async {
    when(dio.get('workspace/summary')).thenAnswer((_) async => Response(requestOptions: RequestOptions(path: ''), statusCode: 200, data: {
          'totalProcedures': 1,
          'inProgress': 1,
          'completed': 0,
          'pending': 0,
        }));
    final result = await ds.getWorkspaceSummary();
    expect(result.totalProcedures, 1);
    verify(dio.get('https://api/workspace/summary'));
  });

  test('getMyProcedures success', () async {
    when(dio.get('')).thenAnswer((_) async => Response(requestOptions: RequestOptions(path: ''), statusCode: 200, data: [
          {
            'id': 'p1',
            'title': 't',
            'organization': 'o',
            'status': 'inProgress',
            'startDate': DateTime(2024, 1, 1).toIso8601String(),
            'steps': const [],
          }
        ]));
    final result = await ds.getMyProcedures();
    expect(result.length, 1);
    verify(dio.get('https://api/workspace/procedures'));
  });

  test('getProceduresByStatus success', () async {
    when(dio.get('', queryParameters: anyNamed('queryParameters'))).thenAnswer((_) async => Response(requestOptions: RequestOptions(path: ''), statusCode: 200, data: const []));
    final result = await ds.getProceduresByStatus('inProgress');
    expect(result, isA<List>());
    verify(dio.get('https://api/workspace/procedures', queryParameters: {'status': 'inProgress'}));
  });

  test('getProceduresByOrganization success', () async {
    when(dio.get('', queryParameters: anyNamed('queryParameters'))).thenAnswer((_) async => Response(requestOptions: RequestOptions(path: ''), statusCode: 200, data: const []));
    final result = await ds.getProceduresByOrganization('ETA');
    expect(result, isA<List>());
    verify(dio.get('https://api/workspace/procedures', queryParameters: {'organization': 'ETA'}));
  });

  test('getProcedureDetail success', () async {
    when(dio.get('')).thenAnswer((_) async => Response(requestOptions: RequestOptions(path: ''), statusCode: 200, data: {
          'id': 'p1',
          'title': 't',
          'organization': 'o',
          'status': 'inProgress',
          'startDate': DateTime(2024, 1, 1).toIso8601String(),
          'steps': const [],
        }));
    final result = await ds.getProcedureDetail('p1');
    expect(result.id, 'p1');
    verify(dio.get('https://api/workspace/procedures/p1'));
  });

  test('updateStepStatus success', () async {
    when(dio.patch('', data: anyNamed('data'))).thenAnswer((_) async => Response(requestOptions: RequestOptions(path: ''), statusCode: 200));
    final result = await ds.updateStepStatus('p1', 's1', true);
    expect(result, true);
    verify(dio.patch('https://api/workspace/procedures/p1/steps/s1', data: {'isCompleted': true}));
  });

  test('saveProgress success', () async {
    when(dio.post('')).thenAnswer((_) async => Response(requestOptions: RequestOptions(path: ''), statusCode: 201));
    final result = await ds.saveProgress('p1');
    expect(result, true);
    verify(dio.post('https://api/workspace/procedures/p1/progress'));
  });
}


