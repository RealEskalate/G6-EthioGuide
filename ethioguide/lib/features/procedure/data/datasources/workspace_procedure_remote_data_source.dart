import 'package:dio/dio.dart';
import 'package:ethioguide/core/error/exception.dart';
import 'package:ethioguide/features/AI%20chat/data/models/conversation_model.dart';
import 'package:ethioguide/features/procedure/data/models/procedure_step_model.dart';
import '../../domain/entities/workspace_procedure.dart';
import '../models/workspace_procedure_model.dart';
import '../models/workspace_summary_model.dart';

/// Remote data source for workspace procedure operations
abstract class WorkspaceProcedureRemoteDataSource {
  // Original methods for workspace procedures
  Future<List<WorkspaceProcedureModel>> getMyProcedures();
  Future<WorkspaceSummaryModel> getWorkspaceSummary();
  Future<List<WorkspaceProcedureModel>> getProceduresByStatus(String status);
  Future<List<WorkspaceProcedureModel>> getProceduresByOrganization(String organization);
  Future<List<MyProcedureStepModel>> getProcedureDetail(String id);
  Future<bool> updateStepStatus(String procedureId, String stepId, bool isCompleted);
  Future<bool> saveProgress(String procedureId);
}


class WorkspaceProcedureRemoteDataSourceImpl implements WorkspaceProcedureRemoteDataSource {
  final Dio dio;


  WorkspaceProcedureRemoteDataSourceImpl({
    required this.dio
  });

  // Original methods implementation
  @override
  Future<List<WorkspaceProcedureModel>> getMyProcedures() async {
    try {
      final response = await dio.get('checklists/myProcedures');
      if (response.statusCode == 200) {
        final data = response.data as List<dynamic>;

              final futures = data.map((element) async {
        final procId = element['procedure_id'] as String;

        final procResponse = await dio.get('procedures/$procId');

        if (procResponse.statusCode == 200) {
          final procedureJson = procResponse.data as Map<String, dynamic>;

          // Merge into WorkspaceProcedureModel
          return WorkspaceProcedureModel.fromJson(element, procedureJson);
        /*   (
            id: element['id'] as String,
            procedure: ProcedureModel.fromJson(procedureJson),
            status: element['status'] as String,
            progressPercentage: element['percent'] as int? ?? 0,
          ); */

        } else {
          throw ServerException(
            message: 'Failed to fetch procedure $procId',
            statusCode: procResponse.statusCode,
          );
        }
      }).toList();


        return data.map((e) => WorkspaceProcedureModel.fromJson(e as Map<String, dynamic> , {})).toList();
        
      }
      throw ServerException(message: 'Unexpected status code', statusCode: response.statusCode);
    } on DioException catch (e) {
      final status = e.response?.statusCode;
      final msg = e.message ?? 'Network error while fetching workspace procedures';
      throw ServerException(message: msg, statusCode: status);
    } catch (e) {
      throw ServerException(message: e.toString());
    }
  }

  @override
  Future<WorkspaceSummaryModel> getWorkspaceSummary() async {
    try {
      final response = await dio.get('workspace/summary');
      if (response.statusCode == 200) {
        return WorkspaceSummaryModel.fromJson(response.data as Map<String, dynamic>);
      }
      throw ServerException(message: 'Unexpected status code', statusCode: response.statusCode);
    } on DioException catch (e) {
      final status = e.response?.statusCode;
      final msg = e.message ?? 'Network error while fetching workspace summary';
      throw ServerException(message: msg, statusCode: status);
    } catch (e) {
      throw ServerException(message: e.toString());
    }
  }

  @override
  Future<List<WorkspaceProcedureModel>> getProceduresByStatus(String status) async {
    try {
      final response = await dio.get('workspace/procedures', queryParameters: {'status': status});
      if (response.statusCode == 200) {
        final data = response.data as List<dynamic>;
        return data.map((e) => WorkspaceProcedureModel.fromJson(e as Map<String, dynamic>, {})).toList();
      }
      throw ServerException(message: 'Unexpected status code', statusCode: response.statusCode);
    } on DioException catch (e) {
      final statusCode = e.response?.statusCode;
      final message = e.message ?? 'Network error while fetching procedures by status';
      throw ServerException(message: message, statusCode: statusCode);
    } catch (e) {
      throw ServerException(message: e.toString());
    }
  }

  @override
  Future<List<WorkspaceProcedureModel>> getProceduresByOrganization(String organization) async {
    try {
      final response = await dio.get('workspace/procedures', queryParameters: {'organization': organization});
      if (response.statusCode == 200) {
        final data = response.data as List<dynamic>;
        return data.map((e) => WorkspaceProcedureModel.fromJson(e as Map<String, dynamic> , {})).toList();
      }
      throw ServerException(message: 'Unexpected status code', statusCode: response.statusCode);
    } on DioException catch (e) {
      final statusCode = e.response?.statusCode;
      final message = e.message ?? 'Network error while fetching procedures by organization';
      throw ServerException(message: message, statusCode: statusCode);
    } catch (e) {
      throw ServerException(message: e.toString());
    }
  }

  // New methods for workspace procedure detail feature
  @override
  Future<List<MyProcedureStepModel>> getProcedureDetail(String id) async {
    try {
      final response = await dio.get('checklists/$id');
      if (response.statusCode == 200) {

        final data = response.data as List<dynamic>;
        return  data.map((e) => MyProcedureStepModel.fromJson(e as Map<String, dynamic>)).toList();
        
        
      }
      throw ServerException(message: 'Unexpected status code', statusCode: response.statusCode);
    } on DioException catch (e) {
      final status = e.response?.statusCode;
      final msg = e.message ?? 'Network error while fetching procedure detail';
      throw ServerException(message: msg, statusCode: status);
    } catch (e) {
      throw ServerException(message: e.toString());
    }
  }

  @override
  Future<bool> updateStepStatus(String procedureId, String stepId, bool isCompleted) async {
    try {
      final response = await dio.patch('workspace/procedures/$procedureId/steps/$stepId', data: {'isCompleted': isCompleted});
      if (response.statusCode == 200 || response.statusCode == 204) {
        return true;
      }
      throw ServerException(message: 'Unexpected status code', statusCode: response.statusCode);
    } on DioException catch (e) {
      final status = e.response?.statusCode;
      final msg = e.message ?? 'Network error while updating step status';
      throw ServerException(message: msg, statusCode: status);
    } catch (e) {
      throw ServerException(message: e.toString());
    }
  }

  @override
  Future<bool> saveProgress(String procedureId) async {
    try {
      final response = await dio.post('workspace/procedures/$procedureId/progress');
      if (response.statusCode == 200 || response.statusCode == 201) {
        return true;
      }
      throw ServerException(message: 'Unexpected status code', statusCode: response.statusCode);
    } on DioException catch (e) {
      final status = e.response?.statusCode;
      final msg = e.message ?? 'Network error while saving progress';
      throw ServerException(message: msg, statusCode: status);
    } catch (e) {
      throw ServerException(message: e.toString());
    }
  }

  Map<String, dynamic> _getMockProcedureDetail(String id) {
    return {
      'id': id,
      'title': 'Driver\'s License Renewal',
      'organization': 'Ethiopian Transport Authority',
      'status': 'inProgress',
      'progressPercentage': 40,
      'documentsUploaded': 2,
      'totalDocuments': 5,
      'startDate': '2024-01-15T00:00:00.000Z',
      'estimatedCompletion': '2024-01-20T00:00:00.000Z',
      'completedDate': null,
      'notes': 'Renewal process in progress',
      'steps': [
        {
          'id': 'step1',
          'title': 'Fill Application Form',
          'description': 'Download and complete the official renewal form with accurate information.',
          'isCompleted': true,
          'completionStatus': 'Completed',
          'order': 1,
        },
        {
          'id': 'step2',
          'title': 'Prepare Required Documents',
          'description': 'Gather all necessary documents and make copies as needed.',
          'isCompleted': true,
          'completionStatus': 'Documents Ready',
          'order': 2,
        },
        {
          'id': 'step3',
          'title': 'Visit License Office',
          'description': 'Submit your application and documents at the designated office.',
          'isCompleted': false,
          'completionStatus': 'Submitted',
          'order': 3,
        },
        {
          'id': 'step4',
          'title': 'Pay Renewal Fees',
          'description': 'Complete payment at the cashier or through online payment if available.',
          'isCompleted': false,
          'completionStatus': null,
          'order': 4,
        },
        {
          'id': 'step5',
          'title': 'Collect New License',
          'description': 'Return after processing period to collect your renewed license.',
          'isCompleted': false,
          'completionStatus': null,
          'order': 5,
        },
      ],
      'estimatedTime': '2-3 days',
      'difficulty': 'Easy',
      'officeType': 'Authority',
      'quickTips': [
        'Bring original documents and photocopies',
        'Visit early in the morning for shorter queues',
        'Keep your receipt until you collect the license',
      ],
      'requiredDocuments': [
        'Original expired driving license',
        'National ID card (original + copy)',
        'Passport-size photos (2 copies)',
        'Medical certificate',
        'Renewal application form',
      ],
    };
  }
}
