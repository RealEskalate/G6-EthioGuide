import 'package:dio/dio.dart';
import 'package:ethioguide/features/procedure/data/models/workspace_procedure_model.dart';
import 'package:ethioguide/features/procedure/data/models/workspace_summary_model.dart';

/// Remote data source for workspace procedures
abstract class WorkspaceProcedureRemoteDataSource {
  /// Get all workspace procedures from remote API
  Future<List<WorkspaceProcedureModel>> getWorkspaceProcedures();
  
  /// Get workspace summary from remote API
  Future<WorkspaceSummaryModel> getWorkspaceSummary();
  
  /// Get procedures by status from remote API
  Future<List<WorkspaceProcedureModel>> getProceduresByStatus(String status);
  
  /// Get procedures by organization from remote API
  Future<List<WorkspaceProcedureModel>> getProceduresByOrganization(String organization);
  
  /// Create a new workspace procedure via remote API
  Future<WorkspaceProcedureModel> createWorkspaceProcedure(WorkspaceProcedureModel procedure);
  
  /// Update an existing workspace procedure via remote API
  Future<WorkspaceProcedureModel> updateWorkspaceProcedure(WorkspaceProcedureModel procedure);
  
  /// Delete a workspace procedure via remote API
  Future<bool> deleteWorkspaceProcedure(String id);
  
  /// Update procedure progress via remote API
  Future<WorkspaceProcedureModel> updateProgress(String id, int progressPercentage);
  
  /// Upload document for a procedure via remote API
  Future<bool> uploadDocument(String procedureId, String documentPath);
}

/// Implementation of remote data source
class WorkspaceProcedureRemoteDataSourceImpl implements WorkspaceProcedureRemoteDataSource {
  final Dio dio;
  final String baseUrl;

  const WorkspaceProcedureRemoteDataSourceImpl({
    required this.dio,
    this.baseUrl = 'https://api.ethioguide.com/api/v1',
  });

  @override
  Future<List<WorkspaceProcedureModel>> getWorkspaceProcedures() async {
    try {
      final response = await dio.get('$baseUrl/workspace/procedures');
      
      if (response.statusCode == 200) {
        final List<dynamic> data = response.data['data'] ?? [];
        return data.map((json) => WorkspaceProcedureModel.fromJson(json)).toList();
      } else {
        throw Exception('Failed to load workspace procedures');
      }
    } catch (e) {
      throw Exception('Network error: $e');
    }
  }

  @override
  Future<WorkspaceSummaryModel> getWorkspaceSummary() async {
    try {
      final response = await dio.get('$baseUrl/workspace/summary');
      
      if (response.statusCode == 200) {
        return WorkspaceSummaryModel.fromJson(response.data['data'] ?? {});
      } else {
        throw Exception('Failed to load workspace summary');
      }
    } catch (e) {
      throw Exception('Network error: $e');
    }
  }

  @override
  Future<List<WorkspaceProcedureModel>> getProceduresByStatus(String status) async {
    try {
      final response = await dio.get('$baseUrl/workspace/procedures', 
        queryParameters: {'status': status});
      
      if (response.statusCode == 200) {
        final List<dynamic> data = response.data['data'] ?? [];
        return data.map((json) => WorkspaceProcedureModel.fromJson(json)).toList();
      } else {
        throw Exception('Failed to load procedures by status');
      }
    } catch (e) {
      throw Exception('Network error: $e');
    }
  }

  @override
  Future<List<WorkspaceProcedureModel>> getProceduresByOrganization(String organization) async {
    try {
      final response = await dio.get('$baseUrl/workspace/procedures', 
        queryParameters: {'organization': organization});
      
      if (response.statusCode == 200) {
        final List<dynamic> data = response.data['data'] ?? [];
        return data.map((json) => WorkspaceProcedureModel.fromJson(json)).toList();
      } else {
        throw Exception('Failed to load procedures by organization');
      }
    } catch (e) {
      throw Exception('Network error: $e');
    }
  }

  @override
  Future<WorkspaceProcedureModel> createWorkspaceProcedure(WorkspaceProcedureModel procedure) async {
    try {
      final response = await dio.post(
        '$baseUrl/workspace/procedures',
        data: procedure.toJson(),
      );
      
      if (response.statusCode == 201) {
        return WorkspaceProcedureModel.fromJson(response.data['data'] ?? {});
      } else {
        throw Exception('Failed to create workspace procedure');
      }
    } catch (e) {
      throw Exception('Network error: $e');
    }
  }

  @override
  Future<WorkspaceProcedureModel> updateWorkspaceProcedure(WorkspaceProcedureModel procedure) async {
    try {
      final response = await dio.put(
        '$baseUrl/workspace/procedures/${procedure.id}',
        data: procedure.toJson(),
      );
      
      if (response.statusCode == 200) {
        return WorkspaceProcedureModel.fromJson(response.data['data'] ?? {});
      } else {
        throw Exception('Failed to update workspace procedure');
      }
    } catch (e) {
      throw Exception('Network error: $e');
    }
  }

  @override
  Future<bool> deleteWorkspaceProcedure(String id) async {
    try {
      final response = await dio.delete('$baseUrl/workspace/procedures/$id');
      
      if (response.statusCode == 200) {
        return true;
      } else {
        throw Exception('Failed to delete workspace procedure');
      }
    } catch (e) {
      throw Exception('Network error: $e');
    }
  }

  @override
  Future<WorkspaceProcedureModel> updateProgress(String id, int progressPercentage) async {
    try {
      final response = await dio.patch(
        '$baseUrl/workspace/procedures/$id/progress',
        data: {'progressPercentage': progressPercentage},
      );
      
      if (response.statusCode == 200) {
        return WorkspaceProcedureModel.fromJson(response.data['data'] ?? {});
      } else {
        throw Exception('Failed to update progress');
      }
    } catch (e) {
      throw Exception('Network error: $e');
    }
  }

  @override
  Future<bool> uploadDocument(String procedureId, String documentPath) async {
    try {
      final formData = FormData.fromMap({
        'document': await MultipartFile.fromFile(documentPath),
        'procedureId': procedureId,
      });
      
      final response = await dio.post(
        '$baseUrl/workspace/procedures/$procedureId/documents',
        data: formData,
      );
      
      if (response.statusCode == 200) {
        return true;
      } else {
        throw Exception('Failed to upload document');
      }
    } catch (e) {
      throw Exception('Network error: $e');
    }
  }
}
