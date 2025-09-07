
import 'package:equatable/equatable.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_step.dart';

class MyProcedureStepModel extends MyProcedureStep {
  const MyProcedureStepModel({
    required super.id,
    required super.title,
    required super.isChecked,
  });

  // From JSON
  factory MyProcedureStepModel.fromJson(Map<String, dynamic> json) {
    return MyProcedureStepModel(
      id: json['id'] as String,
      title: json['content'] as String,
      isChecked: json['is_checked'] as bool,
    );
  }

  // To JSON
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'content': title,
      'is_checked': isChecked,
    };
  }

  @override
  List<Object?> get props => [id, title, isChecked];
}

// accrmony
