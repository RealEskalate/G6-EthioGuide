import '../../domain/entities/procedure_detail.dart';
import 'procedure_model.dart';

class WorkspaceProcedureModel extends ProcedureDetail {
  const WorkspaceProcedureModel({
    required super.id,
    required super.procedure,
    required super.status,
    required super.progressPercentage,
  });

  factory WorkspaceProcedureModel.fromJson(
    Map<String, dynamic> json,
    Map<String, dynamic> procedureJson,
  ) {
    final content = json['Content'] as Map<String, dynamic>? ?? {};
    final fees = json['Fees'] as Map<String, dynamic>? ?? {};
    final processingTime =
        json['ProcessingTime'] as Map<String, dynamic>? ?? {};

    return WorkspaceProcedureModel(
      id: procedureJson['id'] as String,
      procedure: ProcedureModel(
        id: json['ID'] ?? '',
        title: json['Name'] ?? '',
        duration: ProcessTimeModel(
          maxday: processingTime['MaxDays'] ?? 0,
          minday: processingTime['MinDays'] ?? 0,
        ),
        cost: FeeModel(
          amount: fees['Amount'] ?? 0,
          currency: fees['Currency'] ?? '',
        ).toDisplayString(),
        requiredDocuments: (content['Prerequisites'] as List<dynamic>? ?? [])
            .map((doc) => doc.toString())
            .toList(),
        steps:
            (content['Steps'] as Map<String, dynamic>? ?? {}).entries
                .map((entry) => ProcedureStepModel.fromJson(entry))
                .toList()
              ..sort((a, b) => a.number.compareTo(b.number)),
      ),

      status: procedureJson['status'] as String,
      progressPercentage: (procedureJson['percent'] ?? 0) as int,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'procedure': (procedure as ProcedureModel).toJson(), // cast back to model
      'status': status,
      'progressPercentage': progressPercentage,
    };
  }

  /// Convert back to pure domain entity if needed
  ProcedureDetail toEntity() {
    return ProcedureDetail(
      id: id,
      procedure:
          procedure, // already a Procedure because ProcedureModel extends Procedure
      status: status,
      progressPercentage: progressPercentage,
    );
  }

  @override
  List<Object?> get props => [id, procedure, status, progressPercentage];
}


// class WorkspaceProcedureModel extends ProcedureDetail {


//   const WorkspaceProcedureModel({
//     required super.id,
//     required super.procedure,
//     required super.status,
//     required super.progressPercentage,
//   });

//   factory WorkspaceProcedureModel.fromJson(Map<String, dynamic> json) {
//     return WorkspaceProcedureModel(
//       id: json['id'] as String,
//       procedure: ,
//       status: json['status'] as String,
//       progressPercentage: (json['progressPercentage'] ?? 0) as int,
//     );
//   }

//   Map<String, dynamic> toJson() {
//     return {
//       'id': id,
//       'procedure': procedure.toJson(),
//       'status': status,
//       'progressPercentage': progressPercentage,
//     };
//   }

 

//   @override
//   List<Object?> get props => [id, procedure, status, progressPercentage];
// }