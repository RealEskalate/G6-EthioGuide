import 'package:equatable/equatable.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure.dart';



/// Domain entity representing full procedure details used by the UI and repository.
class ProcedureDetail extends Equatable {
  final String id;
  final Procedure procedure; // ✅ domain entity, not ProcedureModel
  final String status;
  final int progressPercentage;

  const ProcedureDetail({
    required this.id,
    required this.procedure,
    required this.status,
    required this.progressPercentage,
  });

  @override
  List<Object?> get props => [id, procedure, status, progressPercentage];
}
