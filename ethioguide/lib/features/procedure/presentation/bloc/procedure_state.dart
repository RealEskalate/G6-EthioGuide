part of 'procedure_bloc.dart';

enum ProcedureStatus { initial, loading, success, failure }

class ProcedureState extends Equatable {
  final ProcedureStatus status;
  final List<Procedure> procedures;
  final String? errorMessage;

  const ProcedureState({
    required this.status,
    required this.procedures,
    this.errorMessage,
  });

  const ProcedureState.initial()
      : status = ProcedureStatus.initial,
        procedures = const [],
        errorMessage = null;

  ProcedureState copyWith({
    ProcedureStatus? status,
    List<Procedure>? procedures,
    String? errorMessage,
  }) {
    return ProcedureState(
      status: status ?? this.status,
      procedures: procedures ?? this.procedures,
      errorMessage: errorMessage,
    );
  }

  @override
  List<Object?> get props => [status, procedures, errorMessage];
}
