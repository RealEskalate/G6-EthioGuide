import 'package:equatable/equatable.dart';

import '../../presentation/bloc/procedure_bloc.dart';

/// Domain entity representing a single step in a workspace procedure
class MyProcedureStep extends Equatable {
  final String id;
  final String title;
  final bool isChecked;


  const MyProcedureStep({
    required this.id,
    required this.title,
    required this.isChecked,
  });

  @override
  List<Object?> get props => [
        id,
        title,
        isChecked,
      ];
}

/// Domain entity representing procedure details with steps
enum ProcedureStatus {
  notStarted('Not Started'),
  inProgress('In Progress'),
  completed('Completed');

  const ProcedureStatus(this.displayName);
  
  final String displayName;
  
  String get displayValue => displayName;
}
