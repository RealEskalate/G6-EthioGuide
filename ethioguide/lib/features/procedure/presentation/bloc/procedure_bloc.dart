import 'package:bloc/bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_feadback.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_procedures.dart';
import 'package:ethioguide/features/procedure/domain/usecases/getprocedurebyid.dart';
import 'package:ethioguide/features/procedure/domain/usecases/save_feedback.dart';
import 'package:ethioguide/features/procedure/domain/usecases/save_procedure.dart';

part 'procedure_event.dart';
part 'procedure_state.dart';

class ProcedureBloc extends Bloc<ProcedureEvent, ProcedureState> {
  final GetProcedures getProcedures;
  final SaveProcedure saveProcedure;
  final GetProceduresbyid getProceduresbyid;
  final GetFeedbacks getFeedbacks;
  final SaveFeedback saveFeedback;

  ProcedureBloc({
    required this.getFeedbacks,
    required this.saveFeedback,
    required this.getProcedures,
    required this.saveProcedure,
    required this.getProceduresbyid,
  }) : super(const ProcedureState.initial()) {

    on<SaveProcedureEvent>((event, emit) async {
      emit(state.copyWith(status: ProcedureStatus.loading));
      final result = await saveProcedure(event.procedureId);

      print('SaveProcedure result: $result');

      result.fold(
        (failure) => emit(
          state.copyWith(
            status: ProcedureStatus.failure,
            errorMessage: failure.message,
          ),
        ),
        (_) => emit(state.copyWith(status: ProcedureStatus.success)),
      );
    });

    on<LoadProcedureByIdEvent>((event, emit) async {
      emit(state.copyWith(status: ProcedureStatus.loading));
      final result = await getProceduresbyid(event.procedureId);
      final feedback = await getFeedbacks(event.procedureId);

      print('GetProceduresbyid result: $result');

      result.fold(
        (failure) => emit(
          state.copyWith(
            status: ProcedureStatus.failure,
            errorMessage: failure.message,
          ),
        ),
        (procedure) => emit(
          state.copyWith(
            status: ProcedureStatus.success,
            selectedProcedure: procedure,
            feedbacks: feedback.isRight() ? feedback.getOrElse(() => []) : [],
          ),
        ),
      );
    });

    on<LoadProceduresEvent>((event, emit) async {
      emit(state.copyWith(status: ProcedureStatus.loading));
      final result = await getProcedures(event.name);

      print('GetProcedures result: $result');

      result.fold(
        (failure) => emit(
          state.copyWith(
            status: ProcedureStatus.failure,
            errorMessage: failure.message,
          ),
        ),
        (procedures) => emit(
          state.copyWith(
            status: ProcedureStatus.success,
            procedures: procedures,
          ),
        ),
      );
    });

    

    on<SaveFeedbackEvent>((event, emit) async {
      emit(state.copyWith(status: ProcedureStatus.loading));
      final result = await saveFeedback(
        feedback: event.feadback,
        procedureId: event.procedureId,
        tags: event.tags,
        type: event.type,
      );

      print('SaveFeedback result: $result');

      result.fold(
        (failure) => emit(
          state.copyWith(
            status: ProcedureStatus.failure,
            errorMessage: failure.message,
          ),
        ),
        (_) => emit(state.copyWith(status: ProcedureStatus.success)),
      );
    });
  }
}
