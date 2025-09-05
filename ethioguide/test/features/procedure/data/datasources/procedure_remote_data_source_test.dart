import 'package:dio/dio.dart';
import 'package:ethioguide/core/error/exception.dart';
import 'package:ethioguide/features/procedure/data/datasources/procedure_remote_data_source.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

import 'procedure_remote_data_source_test.mocks.dart';




@GenerateMocks([Dio])
void main() {
  late MockDio dio;
  late ProcedureRemoteDataSourceImpl ds;

  setUp(() {
    dio = MockDio();
    ds = ProcedureRemoteDataSourceImpl(client: dio, baseUrl: 'https://api');
  });

  test('getProcedures returns models on 200', () async {
    when(dio.get('procedure')).thenAnswer((_) async => Response(requestOptions: RequestOptions(path: 'procedure'), statusCode: 200, data: [
          {
            'id': '1',
            'title': 'Passport',
            'category': 'Travel',
            'duration': '2 weeks',
            'cost': '1200 ETB',
            'icon': 'passport',
            'isQuickAccess': true,
          }
        ]));

    final result = await ds.getProcedures();
    expect(result.length, 1);
    verify(dio.get('https://api/procedures'));
  });

  test('throws ServerException on non-200', () async {
    when(dio.get('')).thenAnswer((_) async => Response(requestOptions: RequestOptions(path: ''), statusCode: 500));
    expect(() => ds.getProcedures(), throwsA(isA<ServerException>()));
  });
}


