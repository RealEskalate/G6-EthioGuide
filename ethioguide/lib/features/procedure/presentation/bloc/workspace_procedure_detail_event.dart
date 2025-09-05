part of 'workspace_procedure_detail_bloc.dart';

abstract class WorkspaceProcedureDetailEvent extends Equatable {
  const WorkspaceProcedureDetailEvent();

  @override
  List<Object?> get props => [];
}

class FetchProcedureDetail extends WorkspaceProcedureDetailEvent {
  final String id;
  const FetchProcedureDetail(this.id);

  @override
  List<Object?> get props => [id];
}

class UpdateStepStatus extends WorkspaceProcedureDetailEvent {
  final String procedureId;
  final String stepId;
  final bool isCompleted;

  const UpdateStepStatus(this.procedureId, this.stepId, this.isCompleted);

  @override
  List<Object?> get props => [procedureId, stepId, isCompleted];
}

class SaveProgress extends WorkspaceProcedureDetailEvent {
  final String procedureId;
  const SaveProgress(this.procedureId);

  @override
  List<Object?> get props => [procedureId];
}

class FetchMyProcedures extends WorkspaceProcedureDetailEvent {
  const FetchMyProcedures();
}

class FetchProceduresByStatus extends WorkspaceProcedureDetailEvent {
  final ProcedureStatus status;
  const FetchProceduresByStatus(this.status);

  @override
  List<Object?> get props => [status];
}

class FetchProceduresByOrganization extends WorkspaceProcedureDetailEvent {
  final String organization;
  const FetchProceduresByOrganization(this.organization);

  @override
  List<Object?> get props => [organization];
}

class FetchWorkspaceSummary extends WorkspaceProcedureDetailEvent {
  const FetchWorkspaceSummary();
}


