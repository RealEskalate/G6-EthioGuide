import 'package:equatable/equatable.dart';

/// Domain entity representing a workspace procedure with progress tracking
class WorkspaceProcedure extends Equatable {
  final String id;
  final String title;
  final String organization;
  final ProcedureStatus status;
  final int progressPercentage;
  final int documentsUploaded;
  final int totalDocuments;
  final DateTime startDate;
  final DateTime? estimatedCompletion;
  final DateTime? completedDate;
  final String? notes;

  const WorkspaceProcedure({
    required this.id,
    required this.title,
    required this.organization,
    required this.status,
    required this.progressPercentage,
    required this.documentsUploaded,
    required this.totalDocuments,
    required this.startDate,
    this.estimatedCompletion,
    this.completedDate,
    this.notes,
  });

  @override
  List<Object?> get props => [
        id,
        title,
        organization,
        status,
        progressPercentage,
        documentsUploaded,
        totalDocuments,
        startDate,
        estimatedCompletion,
        completedDate,
        notes,
      ];
}

/// Enum for procedure status
enum ProcedureStatus {
  notStarted('Not Started'),
  inProgress('In Progress'),
  completed('Completed');

  const ProcedureStatus(this.displayName);
  
  final String displayName;
  
  String get displayValue => displayName;
}

/// Workspace summary statistics
class WorkspaceSummary extends Equatable {
  final int totalProcedures;
  final int inProgress;
  final int completed;
  final int totalDocuments;

  const WorkspaceSummary({
    required this.totalProcedures,
    required this.inProgress,
    required this.completed,
    required this.totalDocuments,
  });

  @override
  List<Object?> get props => [
        totalProcedures,
        inProgress,
        completed,
        totalDocuments,
      ];
}
