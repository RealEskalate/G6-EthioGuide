import 'package:ethioguide/features/procedure/data/models/workspace_procedure_model.dart';
import 'package:ethioguide/features/procedure/data/models/workspace_summary_model.dart';
import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';

import '../../domain/entities/procedure_step.dart';

/// Sample data for workspace procedures
class WorkspaceSampleData {
  static List<WorkspaceProcedureModel> getSampleProcedures() {
    return [
      WorkspaceProcedureModel(
        id: '1',
        title: 'New Passport Application',
        organization: 'Immigration Department',
        status: ProcedureStatus.inProgress,
        progressPercentage: 60,
        documentsUploaded: 4,
        totalDocuments: 6,
        startDate: DateTime(2024, 12, 15),
        estimatedCompletion: DateTime(2025, 1, 30),
      ),
       WorkspaceProcedureModel(
        id: '2',
        title: "Driver's License Renewal",
        organization: 'Road Authority',
        status: ProcedureStatus.completed,
        progressPercentage: 100,
        documentsUploaded: 3,
        totalDocuments: 3,
        startDate: DateTime(2024, 11, 20),
        completedDate: DateTime(2024, 12, 10),
      ),
       WorkspaceProcedureModel(
        id: '3',
        title: 'Bank Account Opening',
        organization: 'National Bank',
        status: ProcedureStatus.notStarted,
        progressPercentage: 0,
        documentsUploaded: 0,
        totalDocuments: 5,
        startDate: DateTime(2024, 12, 20),
        estimatedCompletion: DateTime(2025, 1, 15),
      ),
       WorkspaceProcedureModel(
        id: '4',
        title: 'Vehicle Registration',
        organization: 'Road Authority',
        status: ProcedureStatus.inProgress,
        progressPercentage: 30,
        documentsUploaded: 2,
        totalDocuments: 7,
        startDate: DateTime(2024, 12, 1),
        estimatedCompletion: DateTime(2025, 2, 15),
      ),

    ];
  }

  static WorkspaceSummaryModel getSampleSummary() {
    final procedures = getSampleProcedures();
    final inProgress = procedures.where((p) => p.status == ProcedureStatus.inProgress).length;
    final completed = procedures.where((p) => p.status == ProcedureStatus.completed).length;
    final totalDocuments = procedures.fold<int>(0, (sum, p) => sum + p.totalDocuments);

    return WorkspaceSummaryModel(
      totalProcedures: procedures.length,
      inProgress: inProgress,
      completed: completed,
      totalDocuments: totalDocuments,
    );
  }

  static List<String> getSampleOrganizations() {
    return [
      'Immigration Department',
      'Road Authority',
      'National Bank',
      'Trade Ministry',
      'Revenue Authority',
      'Education Ministry',
      'Health Ministry',
      'Urban Development',
      'Customs Authority',
      'Labor Ministry',
      'Communication Authority',
    ];
  }
}
