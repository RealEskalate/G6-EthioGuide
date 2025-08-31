import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import '../../domain/entities/procedure_detail.dart';
import '../../domain/usecases/get_procedure_detail.dart';
import '../../domain/usecases/update_step_status.dart';
import '../../domain/usecases/save_progress.dart';

// Events
abstract class WorkspaceProcedureDetailEvent extends Equatable {
  const WorkspaceProcedureDetailEvent();

  @override
  List<Object?> get props => [];
}

class FetchProcedureDetail extends WorkspaceProcedureDetailEvent {
  final String procedureId;

  const FetchProcedureDetail(this.procedureId);

  @override
  List<Object?> get props => [procedureId];
}

class UpdateStepStatus extends WorkspaceProcedureDetailEvent {
  final String procedureId;
  final String stepId;
  final bool isCompleted;

  const UpdateStepStatus({
    required this.procedureId,
    required this.stepId,
    required this.isCompleted,
  });

  @override
  List<Object?> get props => [procedureId, stepId, isCompleted];
}

class SaveProgress extends WorkspaceProcedureDetailEvent {
  final String procedureId;

  const SaveProgress(this.procedureId);

  @override
  List<Object?> get props => [procedureId];
}

// States
abstract class WorkspaceProcedureDetailState extends Equatable {
  const WorkspaceProcedureDetailState();

  @override
  List<Object?> get props => [];
}

class ProcedureInitial extends WorkspaceProcedureDetailState {}

class ProcedureLoading extends WorkspaceProcedureDetailState {}

class ProcedureLoaded extends WorkspaceProcedureDetailState {
  final ProcedureDetail procedureDetail;

  const ProcedureLoaded(this.procedureDetail);

  @override
  List<Object?> get props => [procedureDetail];
}

class ProcedureError extends WorkspaceProcedureDetailState {
  final String message;

  const ProcedureError(this.message);

  @override
  List<Object?> get props => [message];
}

class StepStatusUpdated extends WorkspaceProcedureDetailState {
  final ProcedureDetail procedureDetail;
  final String stepId;
  final bool isCompleted;

  const StepStatusUpdated({
    required this.procedureDetail,
    required this.stepId,
    required this.isCompleted,
  });

  @override
  List<Object?> get props => [procedureDetail, stepId, isCompleted];
}

class ProgressSaved extends WorkspaceProcedureDetailState {
  final bool success;

  const ProgressSaved(this.success);

  @override
  List<Object?> get props => [success];
}

// Bloc
class WorkspaceProcedureDetailBloc
    extends Bloc<WorkspaceProcedureDetailEvent, WorkspaceProcedureDetailState> {
  final GetProcedureDetail getProcedureDetail;
  final UpdateStepStatus updateStepStatus;
  final SaveProgress saveProgress;

  WorkspaceProcedureDetailBloc({
    required this.getProcedureDetail,
    required this.updateStepStatus,
    required this.saveProgress,
  }) : super(ProcedureInitial()) {
    on<FetchProcedureDetail>(_onFetchProcedureDetail);
    on<UpdateStepStatus>(_onUpdateStepStatus);
    on<SaveProgress>(_onSaveProgress);
  }

  Future<void> _onFetchProcedureDetail(
    FetchProcedureDetail event,
    Emitter<WorkspaceProcedureDetailState> emit,
  ) async {
    emit(ProcedureLoading());

    final result = await getProcedureDetail(event.procedureId);

    result.fold(
      (failure) => emit(ProcedureError(failure)),
      (procedureDetail) => emit(ProcedureLoaded(procedureDetail)),
    );
  }

  Future<void> _onUpdateStepStatus(
    UpdateStepStatus event,
    Emitter<WorkspaceProcedureDetailState> emit,
  ) async {
    final result = await updateStepStatus(
      event.procedureId,
      event.stepId,
      event.isCompleted,
    );

    result.fold(
      (failure) => emit(ProcedureError(failure)),
      (success) {
        if (success && state is ProcedureLoaded) {
          final currentState = state as ProcedureLoaded;
          final updatedSteps = currentState.procedureDetail.steps.map((step) {
            if (step.id == event.stepId) {
              return step.copyWith(isCompleted: event.isCompleted);
            }
            return step;
          }).toList();

          final updatedProcedure = currentState.procedureDetail.copyWith(
            steps: updatedSteps,
            progressPercentage: _calculateProgress(updatedSteps),
          );

          emit(StepStatusUpdated(
            procedureDetail: updatedProcedure,
            stepId: event.stepId,
            isCompleted: event.isCompleted,
          ));
        }
      },
    );
  }

  Future<void> _onSaveProgress(
    SaveProgress event,
    Emitter<WorkspaceProcedureDetailState> emit,
  ) async {
    final result = await saveProgress(event.procedureId);

    result.fold(
      (failure) => emit(ProcedureError(failure)),
      (success) => emit(ProgressSaved(success)),
    );
  }

  int _calculateProgress(List<ProcedureStep> steps) {
    if (steps.isEmpty) return 0;
    final completedSteps = steps.where((step) => step.isCompleted).length;
    return ((completedSteps / steps.length) * 100).round();
  }
}
