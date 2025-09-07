import 'package:ethioguide/features/procedure/domain/entities/procedure.dart';

import 'package:equatable/equatable.dart';
import '../../domain/entities/procedure.dart';

class ProcedureModel extends Procedure {
  const ProcedureModel({
    required super.id,
    required super.title,
    required super.duration,
    required super.cost,
    super.requiredDocuments = const [],
    super.steps = const [],
  });

   factory ProcedureModel.fromJson(Map<String, dynamic> json) {
    final content = json['content'] as Map<String, dynamic>? ?? {};
    final fees = json['fees'] as Map<String, dynamic>? ?? {};
    final processingTime =
        json['processingTime'] as Map<String, dynamic>? ?? {};

    return ProcedureModel(
      id: json['id'] ?? '',
      title: json['name'] ?? '',
      duration: ProcessTimeModel.fromJson(processingTime),
      cost: FeeModel.fromJson(fees).toDisplayString(),
      requiredDocuments:
          (content['prerequisites'] as List<dynamic>?)
              ?.map((doc) => doc.toString())
              .toList() ??
          [],
      steps:
          (content['steps'] as Map<String, dynamic>? ?? {}).entries
              .map((entry) => ProcedureStepModel.fromJson(entry))
              .toList()
            ..sort((a, b) => a.number.compareTo(b.number)),
    );
  } 



  Map<String, dynamic> toJson() {
    return {
      "id": id,
      "name": title,
      "processingTime": (duration as ProcessTimeModel).toJson(),
      "fees": FeeModel.fromDisplayString(cost).toJson(),
      "content": {
        "prerequisites": requiredDocuments,
        "steps": {
          for (var step in steps)
            (step as ProcedureStepModel).number.toString(): step.title,
        },
      },
    };
  }

  // to entity
  Procedure toEntity() {
    return Procedure(
      id: id,
      title: title,
      duration: duration,
      cost: cost,

      requiredDocuments: requiredDocuments,
      steps: steps,
    );
  }

}

class ProcedureStepModel extends ProcedureStep {
  const ProcedureStepModel({required super.number, required super.title});

  factory ProcedureStepModel.fromJson(MapEntry<String, dynamic> entry) {
    return ProcedureStepModel(
      number: int.tryParse(entry.key) ?? 0,
      title: entry.value.toString(),
    );
  }

  Map<String, dynamic> toJson() {
    return {number.toString(): title};
  }
}

class ProcessTimeModel extends ProcessTime {
  const ProcessTimeModel({required super.minday, required super.maxday});

  factory ProcessTimeModel.fromJson(Map<String, dynamic> json) {
    return ProcessTimeModel(
      minday: json['minDays'] ?? 0,
      maxday: json['maxDays'] ?? 0,
    );
  }

  Map<String, dynamic> toJson() {
    return {"minDays": minday, "maxDays": maxday};
  }
}


class FeedbackItemModel extends FeedbackItem {
  const FeedbackItemModel({
    required super.user,
    required super.comment,
    required super.date,
  });

  factory FeedbackItemModel.fromJson(Map<String, dynamic> json ) {
    return FeedbackItemModel(
      user: json['user_id'] ?? 'Unknown User',
      comment: json['content'] ?? '',
      date: json['created_at'] ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      "user_id": user,
      "content": comment,
      "created_at": date,
    };
  }

  /// Convert list of feedbacks from JSON
  static List<FeedbackItemModel> fromJsonList(Map<String, dynamic> json) {
    final feedbacks = json['feedbacks'] as List<dynamic>? ?? [];
    return feedbacks
        .map((f) => FeedbackItemModel.fromJson(f as Map<String, dynamic>))
        .toList();
  }
}




class FeeModel extends Fee {
  const FeeModel({required super.amount, required super.currency});

  factory FeeModel.fromJson(Map<String, dynamic> json) {
    return FeeModel(
      amount: json['amount'] ?? 0,
      currency: json['currency'] ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    return {"amount": amount, "currency": currency};
  }

  /// Helper for display: "300 ETB"
  String toDisplayString() => "$amount $currency".trim();

  /// Convert from "300 ETB" string back to model
  factory FeeModel.fromDisplayString(String value) {
    final parts = value.split(" ");
    final amount = int.tryParse(parts.isNotEmpty ? parts.first : "0") ?? 0;
    final currency = parts.length > 1 ? parts.last : '';
    return FeeModel(amount: amount, currency: currency);
  }
}
