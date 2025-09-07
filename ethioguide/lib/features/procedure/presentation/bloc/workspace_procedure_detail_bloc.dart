import 'package:bloc/bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_detail.dart';
import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_step.dart';

// Import use cases (alias only for those that clash with event names)
import 'package:ethioguide/features/procedure/domain/usecases/get_procedure_detail.dart' as usecase_detail;
import 'package:ethioguide/features/procedure/domain/usecases/update_step_status.dart' as usecase_update;
import 'package:ethioguide/features/procedure/domain/usecases/save_progress.dart' as usecase_save;
import 'package:ethioguide/features/procedure/domain/usecases/get_my_procedure.dart' as usecase_my;
import 'package:ethioguide/features/procedure/domain/usecases/get_procedure_bystattus.dart' as usecase_by_status;
import 'package:ethioguide/features/procedure/domain/usecases/get_procedure_by_organization.dart' as usecase_by_org;
import 'package:ethioguide/features/procedure/domain/usecases/get_workspace_summary.dart' as usecase_summary;

part 'workspace_procedure_detail_event.dart';
part 'workspace_procedure_detail_state.dart';

class WorkspaceProcedureDetailBloc extends Bloc<WorkspaceProcedureDetailEvent, WorkspaceProcedureDetailState> {
  final usecase_detail.GetProcedureDetail getProcedureDetail;
  // final usecase_update.UpdateStepStatus updateStepStatusUseCase;
  // final usecase_save.SaveProgress saveProgressUseCase;
  final usecase_my.GetProcedureDetails getMyProcedureDetails;
  // final usecase_by_status.GetProceduresByStatus getProceduresByStatus;
  // final usecase_by_org.GetProceduresByOrganization getProceduresByOrganization;
  // final usecase_summary.GetWorkspaceSummary getWorkspaceSummary;

  WorkspaceProcedureDetailBloc({
    required this.getProcedureDetail,
    // required this.updateStepStatusUseCase,
    // required this.saveProgressUseCase,
    required this.getMyProcedureDetails,
    // required this.getProceduresByStatus,
    // required this.getProceduresByOrganization,
    // required this.getWorkspaceSummary,
  }) : super(ProcedureInitial()) {


    on<FetchProcedureDetail>((event, emit) async {
      emit(ProcedureLoading());
      final result = await getProcedureDetail(event.id);
      result.fold(
        (error) => emit(ProcedureError(error)),
        (detail) => emit(ProcedureLoaded(detail)),
      );
    });
    
       on<FetchMyProcedures>((event, emit) async {
      emit(ProcedureLoading());
      final result = await getMyProcedureDetails();
      result.fold(
        (failure) => emit(ProcedureError(failure.message)),
        (procedures) => emit(ProceduresListLoaded(procedures)),
      );
    });

    

   /*  on<UpdateStepStatus>((event, emit) async {
      // Keep last loaded detail, show loading only for update if needed
      final current = state;
      final result = await updateStepStatusUseCase(event.procedureId, event.stepId, event.isCompleted);
      result.fold(
        (error) => emit(ProcedureError(error)),
        (success) async {
          if (success) {
            // Refresh detail to reflect latest progress
            final refreshed = await getProcedureDetail(event.procedureId);
            refreshed.fold(
              (error) => emit(ProcedureError(error)),
              (detail) => emit(StepStatusUpdated(detail)),
            );
          } else {
            emit(current);
          }
        },
      );
    });

    on<SaveProgress>((event, emit) async {
      final result = await saveProgressUseCase(event.procedureId);
      result.fold(
        (error) => emit(ProcedureError(error)),
        (success) => emit(ProgressSaved(success)),
      );
    });

 

    on<FetchProceduresByStatus>((event, emit) async {
      emit(ProcedureLoading());
      final result = await getProceduresByStatus(event.status);
      result.fold(
        (failure) => emit(ProcedureError(failure.message)),
        (procedures) => emit(ProceduresListLoaded(procedures)),
      );
    });



    on<FetchWorkspaceSummary>((event, emit) async {
      emit(ProcedureLoading());
      final result = await getWorkspaceSummary();
      result.fold(
        (failure) => emit(ProcedureError(failure.message)),
        (summary) => emit(WorkspaceSummaryLoaded(summary)),
      );
    });
  } */
}


}