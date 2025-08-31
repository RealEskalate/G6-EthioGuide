import 'package:equatable/equatable.dart';

import '../../presentation/bloc/procedure_bloc.dart';

/// Domain entity representing a single step in a workspace procedure
class ProcedureStep extends Equatable {
  final String id;
  final String title;
  final String description;
  final bool isCompleted;
  final String? completionStatus;
  final int order;

  const ProcedureStep({
    required this.id,
    required this.title,
    required this.description,
    required this.isCompleted,
    this.completionStatus,
    required this.order,
  });

  @override
  List<Object?> get props => [
        id,
        title,
        description,
        isCompleted,
        completionStatus,
        order,
      ];
}

/// Domain entity representing procedure details with steps
class ProcedureDetail extends Equatable {
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
  final List<ProcedureStep> steps;
  final String estimatedTime;
  final String difficulty;
  final String officeType;
  final List<String> quickTips;
  final List<String> requiredDocuments;

  const ProcedureDetail({
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
    required this.steps,
    required this.estimatedTime,
    required this.difficulty,
    required this.officeType,
    required this.quickTips,
    required this.requiredDocuments,
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
        steps,
        estimatedTime,
        difficulty,
        officeType,
        quickTips,
        requiredDocuments,
      ];
}
