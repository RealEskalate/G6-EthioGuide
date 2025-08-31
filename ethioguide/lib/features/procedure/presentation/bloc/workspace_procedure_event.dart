part of 'workspace_procedure_bloc.dart';

/// Base class for workspace procedure events
abstract class WorkspaceProcedureEvent extends Equatable {
  const WorkspaceProcedureEvent();

  @override
  List<Object?> get props => [];
}

/// Event to load workspace procedures
class LoadWorkspaceProcedures extends WorkspaceProcedureEvent {
  const LoadWorkspaceProcedures();
}

/// Event to load workspace summary
class LoadWorkspaceSummary extends WorkspaceProcedureEvent {
  const LoadWorkspaceSummary();
}

/// Event to filter procedures by status
class FilterProceduresByStatus extends WorkspaceProcedureEvent {
  final ProcedureStatus status;

  const FilterProceduresByStatus(this.status);

  @override
  List<Object?> get props => [status];
}

/// Event to filter procedures by organization
class FilterProceduresByOrganization extends WorkspaceProcedureEvent {
  final String organization;

  const FilterProceduresByOrganization(this.organization);

  @override
  List<Object?> get props => [organization];
}

/// Event to refresh workspace data
class RefreshWorkspaceData extends WorkspaceProcedureEvent {
  const RefreshWorkspaceData();
}
