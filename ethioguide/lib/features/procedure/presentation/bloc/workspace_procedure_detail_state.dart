part of 'workspace_procedure_detail_bloc.dart';

abstract class WorkspaceProcedureDetailState extends Equatable {
  const WorkspaceProcedureDetailState();

  @override
  List<Object?> get props => [];
}

class ProcedureInitial extends WorkspaceProcedureDetailState {}

class ProcedureLoading extends WorkspaceProcedureDetailState {}

class ProcedureLoaded extends WorkspaceProcedureDetailState {
  final List<MyProcedureStep> procedureDetail;
  const ProcedureLoaded(this.procedureDetail);

  @override
  List<Object?> get props => [procedureDetail];
}

class ProceduresListLoaded extends WorkspaceProcedureDetailState {
  final List<ProcedureDetail> procedures;
  const ProceduresListLoaded(this.procedures);

  @override
  List<Object?> get props => [procedures];
}

class StepStatusUpdated extends WorkspaceProcedureDetailState {
  final ProcedureDetail procedureDetail;
  const StepStatusUpdated(this.procedureDetail);

  @override
  List<Object?> get props => [procedureDetail];
}

class ProgressSaved extends WorkspaceProcedureDetailState {
  final bool success;
  const ProgressSaved(this.success);

  @override
  List<Object?> get props => [success];
}

class WorkspaceSummaryLoaded extends WorkspaceProcedureDetailState {
  final WorkspaceSummary summary;
  const WorkspaceSummaryLoaded(this.summary);

  @override
  List<Object?> get props => [summary];
}

class ProcedureError extends WorkspaceProcedureDetailState {
  final String message;
  const ProcedureError(this.message);

  @override
  List<Object?> get props => [message];
}

class Stepupdate extends WorkspaceProcedureDetailState {
  final bool check;
  const Stepupdate(this.check);

  @override
  List<Object?> get props => [check];
}
