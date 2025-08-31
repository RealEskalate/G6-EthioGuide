import 'package:ethioguide/features/procedure/data/models/workspace_procedure_model.dart';
import 'package:ethioguide/features/procedure/data/models/workspace_summary_model.dart';
import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';

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
       WorkspaceProcedureModel(
        id: '5',
        title: 'Business License Application',
        organization: 'Trade Ministry',
        status: ProcedureStatus.inProgress,
        progressPercentage: 80,
        documentsUploaded: 8,
        totalDocuments: 10,
        startDate: DateTime(2024, 11, 1),
        estimatedCompletion: DateTime(2024, 12, 31),
      ),
       WorkspaceProcedureModel(
        id: '6',
        title: 'Property Tax Payment',
        organization: 'Revenue Authority',
        status: ProcedureStatus.completed,
        progressPercentage: 100,
        documentsUploaded: 2,
        totalDocuments: 2,
        startDate: DateTime(2024, 10, 15),
        completedDate: DateTime(2024, 11, 30),
      ),
       WorkspaceProcedureModel(
        id: '7',
        title: 'Student ID Card',
        organization: 'Education Ministry',
        status: ProcedureStatus.completed,
        progressPercentage: 100,
        documentsUploaded: 1,
        totalDocuments: 1,
        startDate: DateTime(2024, 9, 1),
        completedDate: DateTime(2024, 9, 15),
      ),
       WorkspaceProcedureModel(
        id: '8',
        title: 'Health Insurance Registration',
        organization: 'Health Ministry',
        status: ProcedureStatus.inProgress,
        progressPercentage: 45,
        documentsUploaded: 3,
        totalDocuments: 6,
        startDate: DateTime(2024, 12, 10),
        estimatedCompletion: DateTime(2025, 1, 20),
      ),
       WorkspaceProcedureModel(
        id: '9',
        title: 'Building Permit',
        organization: 'Urban Development',
        status: ProcedureStatus.notStarted,
        progressPercentage: 0,
        documentsUploaded: 0,
        totalDocuments: 12,
        startDate: DateTime(2024, 12, 25),
        estimatedCompletion: DateTime(2025, 3, 15),
      ),
       WorkspaceProcedureModel(
        id: '10',
        title: 'Import License',
        organization: 'Customs Authority',
        status: ProcedureStatus.completed,
        progressPercentage: 100,
        documentsUploaded: 5,
        totalDocuments: 5,
        startDate: DateTime(2024, 8, 1),
        completedDate: DateTime(2024, 10, 15),
      ),
       WorkspaceProcedureModel(
        id: '11',
        title: 'Social Security Registration',
        organization: 'Labor Ministry',
        status: ProcedureStatus.completed,
        progressPercentage: 100,
        documentsUploaded: 2,
        totalDocuments: 2,
        startDate: DateTime(2024, 7, 1),
        completedDate: DateTime(2024, 7, 20),
      ),
       WorkspaceProcedureModel(
        id: '12',
        title: 'Telecom Service Activation',
        organization: 'Communication Authority',
        status: ProcedureStatus.inProgress,
        progressPercentage: 70,
        documentsUploaded: 2,
        totalDocuments: 3,
        startDate: DateTime(2024, 12, 5),
        estimatedCompletion: DateTime(2024, 12, 28),
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
