import 'package:equatable/equatable.dart';

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
