part of 'procedure_bloc.dart';

abstract class ProcedureEvent extends Equatable {
  const ProcedureEvent();

  @override
  List<Object?> get props => [];
}

class LoadProceduresEvent extends ProcedureEvent {
  const LoadProceduresEvent();
}


