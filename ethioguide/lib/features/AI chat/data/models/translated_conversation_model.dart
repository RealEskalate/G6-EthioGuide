import 'package:ethioguide/features/AI%20chat/Domain/entities/translated_conversation.dart';
import 'package:ethioguide/features/AI%20chat/data/models/conversation_model.dart';

class TranslatedConversationModel extends TranslatedConversation {
  const TranslatedConversationModel({
    required super.response,
    required super.procedures,
  });

  factory TranslatedConversationModel.fromJson({
    required Map<String, dynamic> json,
  }) {
    return TranslatedConversationModel(
      response: json['response'],

      procedures: (json['procedures'] as List)
          .map((procedure) => ProcedureModel.fromJson(procedure).toEntity())
          .toList(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'response': response,
      'procedures': procedures
          .map((e) => (ProcedureModel.fromEntity(e!)).toJson())
          .toList(),
    };
  }

  TranslatedConversation toEntity() {
    return TranslatedConversation(response: response, procedures: procedures);
  }

  static TranslatedConversationModel fromEntity({
    required TranslatedConversation entity,
  }) {
    return TranslatedConversationModel(
      response: entity.response,
      procedures: entity.procedures,
    );
  }
}
