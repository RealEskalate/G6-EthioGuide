import '../../domain/entities/procedure_detail.dart';
import '../../domain/entities/procedure_step.dart';
import '../../domain/entities/workspace_procedure.dart';
import 'procedure_step_model.dart';

/// Data model for workspace procedures used by remote datasource and repository
class WorkspaceProcedureModel extends ProcedureDetail {
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
    super.steps = const [],
    super.estimatedTime = '',
    super.difficulty = '',
    super.officeType = '',
    super.quickTips = const [],
    super.requiredDocuments = const [],
  });

  factory WorkspaceProcedureModel.fromJson(Map<String, dynamic> json) {
    final stepsJson = json['steps'] as List<dynamic>? ?? [];
    final steps = stepsJson
        .map((s) => ProcedureStepModel.fromJson(s as Map<String, dynamic>))
        .toList();

    return WorkspaceProcedureModel(
      id: json['id'] as String,
      title: json['title'] as String,
      organization: json['organization'] as String,
      status: ProcedureStatus.values.firstWhere(
        (e) => e.name == json['status'],
        orElse: () => ProcedureStatus.notStarted,
      ),
      progressPercentage: json['progressPercentage'] as int? ??
          ((steps.isEmpty) ? 0 : ((steps.where((s) => s.isCompleted).length / steps.length) * 100).toInt()),
      documentsUploaded: json['documentsUploaded'] as int? ?? 0,
      totalDocuments: json['totalDocuments'] as int? ?? steps.length,
      startDate: DateTime.parse(json['startDate'] as String),
      estimatedCompletion: json['estimatedCompletion'] != null
          ? DateTime.parse(json['estimatedCompletion'] as String)
          : null,
      completedDate: json['completedDate'] != null
          ? DateTime.parse(json['completedDate'] as String)
          : null,
      notes: json['notes'] as String?,
      steps: steps,
      estimatedTime: json['estimatedTime'] as String? ?? '',
      difficulty: json['difficulty'] as String? ?? '',
      officeType: json['officeType'] as String? ?? '',
      quickTips: List<String>.from(json['quickTips'] ?? []),
      requiredDocuments: List<String>.from(json['requiredDocuments'] ?? []),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'organization': organization,
      'status': status.name,
      'progressPercentage': progressPercentage,
      'documentsUploaded': documentsUploaded,
      'totalDocuments': totalDocuments,
      'startDate': startDate.toIso8601String(),
      'estimatedCompletion': estimatedCompletion?.toIso8601String(),
      'completedDate': completedDate?.toIso8601String(),
      'notes': notes,
      'steps': steps.map((e) => (e as ProcedureStepModel).toJson()).toList(),
      'estimatedTime': estimatedTime,
      'difficulty': difficulty,
      'officeType': officeType,
      'quickTips': quickTips,
      'requiredDocuments': requiredDocuments,
    };
  }
}
