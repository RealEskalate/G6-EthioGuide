import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';

/// Data model for workspace summary
class WorkspaceSummaryModel extends WorkspaceSummary {
  const WorkspaceSummaryModel({
    required super.totalProcedures,
    required super.inProgress,
    required super.completed,
    required super.totalDocuments,
  });

  /// Create from JSON
  factory WorkspaceSummaryModel.fromJson(Map<String, dynamic> json) {
    return WorkspaceSummaryModel(
      totalProcedures: json['totalProcedures'] as int? ?? 0,
      inProgress: json['inProgress'] as int? ?? 0,
      completed: json['completed'] as int? ?? 0,
      totalDocuments: json['totalDocuments'] as int? ?? 0,
    );
  }

  /// Convert to JSON
  Map<String, dynamic> toJson() {
    return {
      'totalProcedures': totalProcedures,
      'inProgress': inProgress,
      'completed': completed,
      'totalDocuments': totalDocuments,
    };
  }
}


