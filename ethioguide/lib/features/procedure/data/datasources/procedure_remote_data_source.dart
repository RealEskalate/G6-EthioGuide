import 'package:dio/dio.dart';
import 'package:ethioguide/core/error/exception.dart';
import 'package:ethioguide/features/procedure/data/models/procedure_model.dart';

abstract class ProcedureRemoteDataSource {
  /// Throws [ServerException] on non-200 responses or network errors.
  Future<List<ProcedureModel>> getProcedures();
}

class ProcedureRemoteDataSourceImpl implements ProcedureRemoteDataSource {
  final Dio client;
  final String baseUrl;

  ProcedureRemoteDataSourceImpl({required this.client, this.baseUrl = 'https://example.com/api'});

  @override
  Future<List<ProcedureModel>> getProcedures() async {
    try {
      final response = await client.get('$baseUrl/procedures');
      if (response.statusCode == 200) {
        final data = response.data as List<dynamic>;
        return data.map((e) => ProcedureModel.fromJson(e as Map<String, dynamic>)).toList();
      }
      throw ServerException(message: 'Unexpected status code', statusCode: response.statusCode);
    } on DioException catch (e) {
      final statusCode = e.response?.statusCode;
      final message = e.message ?? 'Network error';
      throw ServerException(message: message, statusCode: statusCode);
    }
  }
}


