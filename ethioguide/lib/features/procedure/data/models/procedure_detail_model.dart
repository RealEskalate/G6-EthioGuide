import '../../domain/entities/procedure_detail.dart';
import '../../domain/entities/procedure_step.dart';
import '../../domain/entities/workspace_procedure.dart';
import 'procedure_step_model.dart';

/// Data model for procedure details
class ProcedureDetailModel extends ProcedureDetail {
  const ProcedureDetailModel({
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
    required super.steps,
    required super.estimatedTime,
    required super.difficulty,
    required super.officeType,
    required super.quickTips,
    required super.requiredDocuments,
  });

  factory ProcedureDetailModel.fromJson(Map<String, dynamic> json) {
    return ProcedureDetailModel(
      id: json['id'] as String,
      title: json['title'] as String,
      organization: json['organization'] as String,
      status: ProcedureStatus.values.firstWhere(
        (e) => e.name == json['status'],
        orElse: () => ProcedureStatus.notStarted,
      ),
      progressPercentage: json['progressPercentage'] as int,
      documentsUploaded: json['documentsUploaded'] as int,
      totalDocuments: json['totalDocuments'] as int,
      startDate: DateTime.parse(json['startDate'] as String),
      estimatedCompletion: json['estimatedCompletion'] != null
          ? DateTime.parse(json['estimatedCompletion'] as String)
          : null,
      completedDate: json['completedDate'] != null
          ? DateTime.parse(json['completedDate'] as String)
          : null,
      notes: json['notes'] as String?,
      steps: (json['steps'] as List<dynamic>)
          .map((step) => ProcedureStepModel.fromJson(step as Map<String, dynamic>))
          .toList(),
      estimatedTime: json['estimatedTime'] as String,
      difficulty: json['difficulty'] as String,
      officeType: json['officeType'] as String,
      quickTips: (json['quickTips'] as List<dynamic>).cast<String>(),
      requiredDocuments: (json['requiredDocuments'] as List<dynamic>).cast<String>(),
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
      'steps': steps.map((step) => (step as ProcedureStepModel).toJson()).toList(),
      'estimatedTime': estimatedTime,
      'difficulty': difficulty,
      'officeType': officeType,
      'quickTips': quickTips,
      'requiredDocuments': requiredDocuments,
    };
  }
}
