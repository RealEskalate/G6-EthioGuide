part of 'procedure_bloc.dart';

enum ProcedureStatus { initial, loading, success, failure }

class ProcedureState extends Equatable {
  final ProcedureStatus status;
  final List<Procedure> procedures;
  final Procedure? selectedProcedure;
  final String? errorMessage;
  final List<FeedbackItem>? feedbacks;

  const ProcedureState({
    required this.status,
    required this.procedures,
    this.errorMessage,
    this.selectedProcedure,
    this.feedbacks,
  });

  const ProcedureState.initial()
      : status = ProcedureStatus.initial,
        procedures = const [],
        selectedProcedure = null,
        feedbacks = const [],
        errorMessage = null;

  ProcedureState copyWith({
    ProcedureStatus? status,
    List<Procedure>? procedures,
    String? errorMessage,
    Procedure? selectedProcedure,
    List<FeedbackItem>? feedbacks,
  }) {
    return ProcedureState(
      status: status ?? this.status,
      procedures: procedures ?? this.procedures,
      errorMessage: errorMessage,
      selectedProcedure: selectedProcedure ?? this.selectedProcedure,
      feedbacks: feedbacks ?? this.feedbacks,
    );
  }

  @override
  List<Object?> get props => [status, procedures, errorMessage , selectedProcedure , feedbacks];
}
