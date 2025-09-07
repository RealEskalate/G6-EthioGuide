import 'package:dio/dio.dart';
import 'package:ethioguide/core/error/exception.dart';
import 'package:ethioguide/features/authentication/data/models/user_model.dart';
import 'package:ethioguide/features/procedure/data/models/procedure_model.dart';

abstract class ProcedureRemoteDataSource {
  /// Throws [ServerException] on non-200 responses or network errors.
  Future<List<ProcedureModel>> getProcedures(
    String? name 
  );
  Future<void> saveProcedure(String procedureid);
  Future<ProcedureModel> getProceduresbyid(String id);
  Future<List<FeedbackItemModel>> getfeadback(String procedureid);

  Future<bool> savefeedback(
    String procedureid,
    String feedback,
    List<String> tag,
    String type,
  );
}

class ProcedureRemoteDataSourceImpl implements ProcedureRemoteDataSource {
  final Dio client;

  ProcedureRemoteDataSourceImpl({required this.client});

  @override
  Future<List<ProcedureModel>> getProcedures(String? name ) async {
    try {
      final response = await client.get('procedures' ,  queryParameters: {
        if (name != null) 'name': name,
        
      },);

      print('Response data: ${response.statusCode}, ${response.data}');

      if (response.statusCode == 200) {
        final body = response.data;

        // print('Response body type: ${body.runtimeType}');
        print('Response body content: $body');

        // Ensure we extract `data` properly
        final List<dynamic> rawList;
        if (body is Map<String, dynamic>) {
          rawList = body['data'] as List<dynamic>? ?? [];
        } else if (body is List) {
          rawList = body;
        } else {
          throw ServerException(message: 'Unexpected response format');
        }

        print('Raw list content: $rawList');
        return rawList
            .map((e) => ProcedureModel.fromJson(e as Map<String, dynamic>))
            .toList();
      }

      throw ServerException(
        message: 'Unexpected status code',
        statusCode: response.statusCode,
      );
    } on DioException catch (e) {
      print('DioException: ${e.message}, Response: ${e.response}');
      final statusCode = e.response?.statusCode;
      final message = e.message ?? 'Network error';
      throw ServerException(message: message, statusCode: statusCode);
    }
  }

  @override
  Future<ProcedureModel> getProceduresbyid(String id) async {
    try {
      final response = await client.get('procedures/$id');

      print('Response data: ${response.statusCode}, ${response.data}');

      if (response.statusCode == 200) {
        final json = response.data as Map<String, dynamic>;

        final content = json['Content'] as Map<String, dynamic>? ?? {};
        final fees = json['Fees'] as Map<String, dynamic>? ?? {};
        final processingTime =
            json['ProcessingTime'] as Map<String, dynamic>? ?? {};

        return ProcedureModel(
          id: json['ID'] ?? '',
          title: json['Name'] ?? '',
          duration: ProcessTimeModel(
            maxday: processingTime['MaxDays'] ?? 0,
            minday: processingTime['MinDays'] ?? 0,
          ),
          cost: FeeModel(
            amount: fees['Amount'] ?? 0,
            currency: fees['Currency'] ?? '',
          ).toDisplayString(),
          requiredDocuments: (content['Prerequisites'] as List<dynamic>? ?? [])
              .map((doc) => doc.toString())
              .toList(),
          steps:
              (content['Steps'] as Map<String, dynamic>? ?? {}).entries
                  .map((entry) => ProcedureStepModel.fromJson(entry))
                  .toList()
                ..sort((a, b) => a.number.compareTo(b.number)),
        );
      }

      throw ServerException(
        message: 'Unexpected status code',
        statusCode: response.statusCode,
      );
    } on DioException catch (e) {
      print('DioException: ${e.message}, Response: ${e.response}');
      final statusCode = e.response?.statusCode;
      final message = e.message ?? 'Network error';
      throw ServerException(message: message, statusCode: statusCode);
    }
  }

  // save a procedure

  @override
  Future<void> saveProcedure(String procedureid) async {
    try {
      final response = await client.post(
        'checklists',
        data: {"procedure_id": procedureid},
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        return;
      }
      throw ServerException(
        message: 'Unexpected status code',
        statusCode: response.statusCode,
      );
    } on DioException catch (e) {
      final statusCode = e.response?.statusCode;
      final message = e.message ?? 'Network error while saving procedure';
      throw ServerException(message: message, statusCode: statusCode);
    } catch (e) {
      throw ServerException(message: e.toString());
    }
  }

  @override
  Future<List<FeedbackItemModel>> getfeadback(String procedureid) async {
    try {
      final response = await client.get('procedures/$procedureid/feedback');

      if (response.statusCode == 200) {
        final json = response.data as Map<String, dynamic>;

        final data = json['feedbacks'] as Map<String, dynamic>;

        print('Feedback data: $data');
        return FeedbackItemModel.fromJsonList(data);
      }
      throw ServerException(
        message: 'Unexpected status code',
        statusCode: response.statusCode,
      );
    } on DioException catch (e) {
      final statusCode = e.response?.statusCode;
      final message = e.message ?? 'Network error while fetching feedback';
      throw ServerException(message: message, statusCode: statusCode);
    } catch (e) {
      throw ServerException(message: e.toString());
    }
  }

  @override
  Future<bool> savefeedback(
    String procedureid,
    String feedback,
    List<String> tag,
    String type,
  ) async {
    try {
      final response = await client.post(
        'feedbacks/$procedureid',
        data: {"content": feedback, "tags": tag, "type": type},
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        return true;
      }
      throw ServerException(
        message: 'Unexpected status code',
        statusCode: response.statusCode,
      );
    } on DioException catch (e) {
      final statusCode = e.response?.statusCode;
      final message = e.message ?? 'Network error while saving feedback';
      throw ServerException(message: message, statusCode: statusCode);
    } catch (e) {
      throw ServerException(message: e.toString());
    }
  }
}
