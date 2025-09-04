import 'package:bloc/bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_procedures.dart';

part 'procedure_event.dart';
part 'procedure_state.dart';

class ProcedureBloc extends Bloc<ProcedureEvent, ProcedureState> {
  final GetProcedures getProcedures;

  ProcedureBloc({required this.getProcedures}) : super(const ProcedureState.initial()) {
    on<LoadProceduresEvent>((event, emit) async {
      emit(state.copyWith(status: ProcedureStatus.loading));
      final result = await getProcedures();
      result.fold(
        (failure) => emit(state.copyWith(status: ProcedureStatus.failure, errorMessage: failure.message)),
        (procedures) => emit(state.copyWith(status: ProcedureStatus.success, procedures: procedures)),
      );
    });
  }
}


