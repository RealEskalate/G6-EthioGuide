import 'package:dio/dio.dart';
import '../../domain/entities/workspace_procedure.dart';
import '../models/procedure_detail_model.dart';
import '../models/workspace_procedure_model.dart';
import '../models/workspace_summary_model.dart';

/// Remote data source for workspace procedure operations
abstract class WorkspaceProcedureRemoteDataSource {
  // Original methods for workspace procedures
  Future<List<WorkspaceProcedureModel>> getWorkspaceProcedures();
  Future<WorkspaceSummaryModel> getWorkspaceSummary();
  Future<List<WorkspaceProcedureModel>> getProceduresByStatus(String status);
  Future<List<WorkspaceProcedureModel>> getProceduresByOrganization(String organization);
  Future<WorkspaceProcedureModel> createWorkspaceProcedure(WorkspaceProcedureModel procedure);
  Future<WorkspaceProcedureModel> updateWorkspaceProcedure(WorkspaceProcedureModel procedure);
  Future<bool> deleteWorkspaceProcedure(String id);
  Future<WorkspaceProcedureModel> updateProgress(String id, int progressPercentage);
  Future<bool> uploadDocument(String procedureId, String documentPath);

  // New methods for workspace procedure detail feature
  Future<ProcedureDetailModel> getProcedureDetail(String id);
  Future<bool> updateStepStatus(String procedureId, String stepId, bool isCompleted);
  Future<bool> saveProgress(String procedureId);
}

class WorkspaceProcedureRemoteDataSourceImpl implements WorkspaceProcedureRemoteDataSource {
  final Dio dio;
  final String baseUrl;

  WorkspaceProcedureRemoteDataSourceImpl({
    required this.dio,
    this.baseUrl = 'https://api.ethioguide.com', // Mock API base URL
  });

  // Original methods implementation
  @override
  Future<List<WorkspaceProcedureModel>> getWorkspaceProcedures() async {
    try {
      await Future.delayed(const Duration(milliseconds: 500));
      // Mock implementation - return empty list for now
      return [];
    } catch (e) {
      throw Exception('Failed to fetch workspace procedures: $e');
    }
  }

  @override
  Future<WorkspaceSummaryModel> getWorkspaceSummary() async {
    try {
      await Future.delayed(const Duration(milliseconds: 300));
      // Mock implementation - return default summary
      return const WorkspaceSummaryModel(
        totalProcedures: 0,
        inProgress: 0,
        completed: 0,
        totalDocuments: 0,
      );
    } catch (e) {
      throw Exception('Failed to fetch workspace summary: $e');
    }
  }

  @override
  Future<List<WorkspaceProcedureModel>> getProceduresByStatus(String status) async {
    try {
      await Future.delayed(const Duration(milliseconds: 400));
      // Mock implementation - return empty list for now
      return [];
    } catch (e) {
      throw Exception('Failed to fetch procedures by status: $e');
    }
  }

  @override
  Future<List<WorkspaceProcedureModel>> getProceduresByOrganization(String organization) async {
    try {
      await Future.delayed(const Duration(milliseconds: 400));
      // Mock implementation - return empty list for now
      return [];
    } catch (e) {
      throw Exception('Failed to fetch procedures by organization: $e');
    }
  }

  @override
  Future<WorkspaceProcedureModel> createWorkspaceProcedure(WorkspaceProcedureModel procedure) async {
    try {
      await Future.delayed(const Duration(milliseconds: 600));
      // Mock implementation - return the same procedure
      return procedure;
    } catch (e) {
      throw Exception('Failed to create workspace procedure: $e');
    }
  }

  @override
  Future<WorkspaceProcedureModel> updateWorkspaceProcedure(WorkspaceProcedureModel procedure) async {
    try {
      await Future.delayed(const Duration(milliseconds: 600));
      // Mock implementation - return the same procedure
      return procedure;
    } catch (e) {
      throw Exception('Failed to update workspace procedure: $e');
    }
  }

  @override
  Future<bool> deleteWorkspaceProcedure(String id) async {
    try {
      await Future.delayed(const Duration(milliseconds: 400));
      // Mock implementation - return true
      return true;
    } catch (e) {
      throw Exception('Failed to delete workspace procedure: $e');
    }
  }

  @override
  Future<WorkspaceProcedureModel> updateProgress(String id, int progressPercentage) async {
    try {
      await Future.delayed(const Duration(milliseconds: 500));
      // Mock implementation - return a default procedure
      return  WorkspaceProcedureModel(
        id: '1',
        title: 'Test Procedure',
        organization: 'Test Org',
        status: ProcedureStatus.inProgress,
        progressPercentage: 50,
        documentsUploaded: 2,
        totalDocuments: 4,
        startDate: DateTime.now(),
      );
    } catch (e) {
      throw Exception('Failed to update progress: $e');
    }
  }

  @override
  Future<bool> uploadDocument(String procedureId, String documentPath) async {
    try {
      await Future.delayed(const Duration(milliseconds: 800));
      // Mock implementation - return true
      return true;
    } catch (e) {
      throw Exception('Failed to upload document: $e');
    }
  }

  // New methods for workspace procedure detail feature
  @override
  Future<ProcedureDetailModel> getProcedureDetail(String id) async {
    try {
      // Simulate API delay
      await Future.delayed(const Duration(milliseconds: 800));
      
      // Mock response data
      final mockData = _getMockProcedureDetail(id);
      return ProcedureDetailModel.fromJson(mockData);
    } catch (e) {
      throw Exception('Failed to fetch procedure detail: $e');
    }
  }

  @override
  Future<bool> updateStepStatus(String procedureId, String stepId, bool isCompleted) async {
    try {
      // Simulate API delay
      await Future.delayed(const Duration(milliseconds: 500));
      
      // Mock successful update
      return true;
    } catch (e) {
      throw Exception('Failed to update step status: $e');
    }
  }

  @override
  Future<bool> saveProgress(String procedureId) async {
    try {
      // Simulate API delay
      await Future.delayed(const Duration(milliseconds: 600));
      
      // Mock successful save
      return true;
    } catch (e) {
      throw Exception('Failed to save progress: $e');
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
