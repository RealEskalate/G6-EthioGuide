part of 'workspace_procedure_bloc.dart';

/// Base class for workspace procedure states
abstract class WorkspaceProcedureState extends Equatable {
  const WorkspaceProcedureState();

  @override
  List<Object?> get props => [];
}

/// Initial state
class WorkspaceProcedureInitial extends WorkspaceProcedureState {}

/// Loading state
class WorkspaceProcedureLoading extends WorkspaceProcedureState {}

/// State when procedures are loaded successfully
class WorkspaceProceduresLoaded extends WorkspaceProcedureState {
  final List<WorkspaceProcedure> procedures;

  const WorkspaceProceduresLoaded(this.procedures);

  @override
  List<Object?> get props => [procedures];
}

/// State when procedures are filtered
class WorkspaceProceduresFiltered extends WorkspaceProcedureState {
  final List<WorkspaceProcedure> procedures;
  final ProcedureStatus? statusFilter;
  final String? organizationFilter;

  const WorkspaceProceduresFiltered(
    this.procedures, [
    this.statusFilter,
    this.organizationFilter,
  ]);

  @override
  List<Object?> get props => [procedures, statusFilter, organizationFilter];
}

/// State when summary is loaded successfully
class WorkspaceSummaryLoaded extends WorkspaceProcedureState {
  final WorkspaceSummary summary;

  const WorkspaceSummaryLoaded(this.summary);

  @override
  List<Object?> get props => [summary];
}

/// Error state for procedures
class WorkspaceProcedureError extends WorkspaceProcedureState {
  final String message;

  const WorkspaceProcedureError(this.message);

  @override
  List<Object?> get props => [message];
}

/// Error state for summary
class WorkspaceSummaryError extends WorkspaceProcedureState {
  final String message;

  const WorkspaceSummaryError(this.message);

  @override
  List<Object?> get props => [message];
}
