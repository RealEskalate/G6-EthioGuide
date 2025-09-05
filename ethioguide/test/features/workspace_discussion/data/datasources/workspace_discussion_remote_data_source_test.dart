import 'package:dio/dio.dart';
import 'package:ethioguide/features/workspace_discussion/data/datasources/workspace_discussion_remote_data_source.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

import 'workspace_discussion_remote_data_source_test.mocks.dart';



@GenerateMocks([Dio])
void main() {
  late MockDio dio;
  late WorkspaceDiscussionRemoteDataSourceImpl ds;

  setUp(() {
    dio = MockDio();
    ds = WorkspaceDiscussionRemoteDataSourceImpl(dio: dio, baseUrl: 'https://api');
  });

  test('getCommunityStats success', () async {
    when(dio.get('')).thenAnswer((_) async => Response(requestOptions: RequestOptions(path: ''), statusCode: 200, data: {'totalMembers': 1, 'activeUsers': 1, 'totalDiscussions': 0}));
    final result = await ds.getCommunityStats();
    expect(result.totalMembers, 1);
    verify(dio.get('https://api/workspace/community/stats'));
  });
}


