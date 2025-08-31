import '../../domain/entities/procedure_step.dart';

/// Data model for procedure steps
class ProcedureStepModel extends ProcedureStep {
  const ProcedureStepModel({
    required super.id,
    required super.title,
    required super.description,
    required super.isCompleted,
    super.completionStatus,
    required super.order,
  });

  factory ProcedureStepModel.fromJson(Map<String, dynamic> json) {
    return ProcedureStepModel(
      id: json['id'] as String,
      title: json['title'] as String,
      description: json['description'] as String,
      isCompleted: json['isCompleted'] as bool,
      completionStatus: json['completionStatus'] as String?,
      order: json['order'] as int,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'description': description,
      'isCompleted': isCompleted,
      'completionStatus': completionStatus,
      'order': order,
    };
  }
}
