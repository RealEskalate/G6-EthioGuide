import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_workspace_procedures.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_workspace_summary.dart';

part 'workspace_procedure_event.dart';
part 'workspace_procedure_state.dart';

/// BLoC for managing workspace procedures
class WorkspaceProcedureBloc extends Bloc<WorkspaceProcedureEvent, WorkspaceProcedureState> {
  final GetWorkspaceProcedures getWorkspaceProcedures;
  final GetWorkspaceSummary getWorkspaceSummary;

  WorkspaceProcedureBloc({
    required this.getWorkspaceProcedures,
    required this.getWorkspaceSummary,
  }) : super(WorkspaceProcedureInitial()) {
    on<LoadWorkspaceProcedures>(_onLoadWorkspaceProcedures);
    on<LoadWorkspaceSummary>(_onLoadWorkspaceSummary);
    on<FilterProceduresByStatus>(_onFilterProceduresByStatus);
    on<FilterProceduresByOrganization>(_onFilterProceduresByOrganization);
    on<RefreshWorkspaceData>(_onRefreshWorkspaceData);
  }

  Future<void> _onLoadWorkspaceProcedures(
    LoadWorkspaceProcedures event,
    Emitter<WorkspaceProcedureState> emit,
  ) async {
    emit(WorkspaceProcedureLoading());
    
    final result = await getWorkspaceProcedures();
    
    result.fold(
      (failure) => emit(WorkspaceProcedureError(failure.message)),
      (procedures) => emit(WorkspaceProceduresLoaded(procedures)),
    );
  }

  Future<void> _onLoadWorkspaceSummary(
    LoadWorkspaceSummary event,
    Emitter<WorkspaceProcedureState> emit,
  ) async {
    final result = await getWorkspaceSummary();
    
    result.fold(
      (failure) => emit(WorkspaceSummaryError(failure.message)),
      (summary) => emit(WorkspaceSummaryLoaded(summary)),
    );
  }

  Future<void> _onFilterProceduresByStatus(
    FilterProceduresByStatus event,
    Emitter<WorkspaceProcedureState> emit,
  ) async {
    if (state is WorkspaceProceduresLoaded) {
      final currentState = state as WorkspaceProceduresLoaded;
      final filteredProcedures = currentState.procedures
          .where((procedure) => procedure.status == event.status)
          .toList();
      
      emit(WorkspaceProceduresFiltered(filteredProcedures, event.status));
    }
  }

  Future<void> _onFilterProceduresByOrganization(
    FilterProceduresByOrganization event,
    Emitter<WorkspaceProcedureState> emit,
  ) async {
    if (state is WorkspaceProceduresLoaded) {
      final currentState = state as WorkspaceProceduresLoaded;
      final filteredProcedures = currentState.procedures
          .where((procedure) => procedure.organization == event.organization)
          .toList();
      
      emit(WorkspaceProceduresFiltered(filteredProcedures, null, event.organization));
    }
  }

  Future<void> _onRefreshWorkspaceData(
    RefreshWorkspaceData event,
    Emitter<WorkspaceProcedureState> emit,
  ) async {
    add(LoadWorkspaceProcedures());
    add(LoadWorkspaceSummary());
  }
}
