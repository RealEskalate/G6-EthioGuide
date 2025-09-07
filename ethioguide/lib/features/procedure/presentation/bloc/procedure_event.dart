part of 'procedure_bloc.dart';

abstract class ProcedureEvent extends Equatable {
  const ProcedureEvent();

  @override
  List<Object?> get props => [];
}

class LoadProceduresEvent extends ProcedureEvent {
  String? name ;
   LoadProceduresEvent(this.name);
}

class SaveProcedureEvent extends ProcedureEvent {
  final String procedureId;

  const SaveProcedureEvent(this.procedureId);

  @override
  List<Object?> get props => [procedureId];
}

class LoadProcedureByIdEvent extends ProcedureEvent {
  final String procedureId;

  const LoadProcedureByIdEvent(this.procedureId);

  @override
  List<Object?> get props => [procedureId];
}

class LoadFeedbackEvent extends ProcedureEvent {
  final String procedureId;

  const LoadFeedbackEvent(this.procedureId);

  @override
  List<Object?> get props => [procedureId];
}

class SaveFeedbackEvent extends ProcedureEvent {
  final String procedureId;
  final String feadback;
  final String type;
  final List<String> tags;

  const SaveFeedbackEvent({
    required this.procedureId,
    required this.feadback,
    required this.type,
    required this.tags,

  });

  @override
  List<Object?> get props => [procedureId, feadback, type, tags];
}


