import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';

/// Data model for workspace procedures
class WorkspaceProcedureModel extends WorkspaceProcedure {
  const WorkspaceProcedureModel({
    required super.id,
    required super.title,
    required super.organization,
    required super.status,
    required super.progressPercentage,
    required super.documentsUploaded,
    required super.totalDocuments,
    required super.startDate,
    super.estimatedCompletion,
    super.completedDate,
    super.notes,
  });

  /// Create from JSON
  factory WorkspaceProcedureModel.fromJson(Map<String, dynamic> json) {
    return WorkspaceProcedureModel(
      id: json['id'] as String,
      title: json['title'] as String,
      organization: json['organization'] as String,
      status: ProcedureStatus.values.firstWhere(
        (e) => e.displayName == json['status'],
        orElse: () => ProcedureStatus.notStarted,
      ),
      progressPercentage: json['progressPercentage'] as int? ?? 0,
      documentsUploaded: json['documentsUploaded'] as int? ?? 0,
      totalDocuments: json['totalDocuments'] as int? ?? 0,
      startDate: DateTime.parse(json['startDate'] as String),
      estimatedCompletion: json['estimatedCompletion'] != null
          ? DateTime.parse(json['estimatedCompletion'] as String)
          : null,
      completedDate: json['completedDate'] != null
          ? DateTime.parse(json['completedDate'] as String)
          : null,
      notes: json['notes'] as String?,
    );
  }

  /// Convert to JSON
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'organization': organization,
      'status': status.displayName,
      'progressPercentage': progressPercentage,
      'documentsUploaded': documentsUploaded,
      'totalDocuments': totalDocuments,
      'startDate': startDate.toIso8601String(),
      'estimatedCompletion': estimatedCompletion?.toIso8601String(),
      'completedDate': completedDate?.toIso8601String(),
      'notes': notes,
    };
  }

  /// Create a copy with updated values
  WorkspaceProcedureModel copyWith({
    String? id,
    String? title,
    String? organization,
    ProcedureStatus? status,
    int? progressPercentage,
    int? documentsUploaded,
    int? totalDocuments,
    DateTime? startDate,
    DateTime? estimatedCompletion,
    DateTime? completedDate,
    String? notes,
  }) {
    return WorkspaceProcedureModel(
      id: id ?? this.id,
      title: title ?? this.title,
      organization: organization ?? this.organization,
      status: status ?? this.status,
      progressPercentage: progressPercentage ?? this.progressPercentage,
      documentsUploaded: documentsUploaded ?? this.documentsUploaded,
      totalDocuments: totalDocuments ?? this.totalDocuments,
      startDate: startDate ?? this.startDate,
      estimatedCompletion: estimatedCompletion ?? this.estimatedCompletion,
      completedDate: completedDate ?? this.completedDate,
      notes: notes ?? this.notes,
    );
  }
}


