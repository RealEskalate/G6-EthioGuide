import 'package:ethioguide/features/AI%20chat/Domain/entities/conversation.dart';

class ProcedureModel extends Procedure {
  const ProcedureModel({required super.id, required super.name});

  factory ProcedureModel.fromJson(Map<String, dynamic> json) {
    return ProcedureModel(
      id: json['id'] as String,
      name: json['name'] as String,
    );
  }

  Map<String, dynamic> toJson() => {'id': id, 'name': name};

  Procedure toEntity() => Procedure(id: id, name: name);

  static ProcedureModel fromEntity(Procedure procedure) {
    return ProcedureModel(id: procedure.id, name: procedure.name);
  }
}

class ConversationModel extends Conversation {
  const ConversationModel({
    required super.id,
    required super.request,
    required super.response,
    required super.source,
    required super.procedures,
  });

  factory ConversationModel.fromJson(Map<String, dynamic> json) {
    return ConversationModel(
      id: json['id'] as String,
      request: json['request'] as String,
      response: json['response'] as String,
      source: json['source'] as String,
      procedures: (json['procedures'] as List)
          .map((e) => ProcedureModel.fromJson(e))
          .toList(),
    );
  }

  Map<String, dynamic> toJson() => {
    'id': id,
    'request': request,
    'response': response,
    'source': source,
    'procedures': procedures
        .map((e) => (e as ProcedureModel).toJson())
        .toList(),
  };

  Conversation toEntity() => Conversation(
    id: id,
    request: request,
    response: response,
    source: source,
    procedures: procedures
        .map((e) => ProcedureModel.fromEntity(e).toEntity())
        .toList(),
  );

  static ConversationModel fromEntity(Conversation conversation) {
    return ConversationModel(
      id: conversation.id,
      request: conversation.request,
      response: conversation.response,
      source: conversation.source,
      procedures: conversation.procedures
          .map((e) => ProcedureModel.fromEntity(e))
          .toList(),
    );
  }
}
