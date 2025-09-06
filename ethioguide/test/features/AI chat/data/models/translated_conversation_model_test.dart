import 'package:ethioguide/features/AI%20chat/Domain/entities/conversation.dart';
import 'package:ethioguide/features/AI%20chat/Domain/entities/translated_conversation.dart';
import 'package:ethioguide/features/AI%20chat/data/models/translated_conversation_model.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('TranslatedConversaionModel', () {
    final translationJson = {
      "response": "Ai response",
      "procedures": [
        {"id": "id1", "name": "related procedure 1"},
        {"id": "id2", "name": "related procedure 2"},
      ],
    };

    final translationModel = TranslatedConversationModel(
      response: 'Ai response',
      procedures: [
        Procedure(id: 'id1', name: 'related procedure 1'),
        Procedure(id: 'id2', name: 'related procedure 2'),
      ],
    );

    final translationEntity = TranslatedConversation(
      response: 'Ai response',
      procedures: [
        Procedure(id: 'id1', name: 'related procedure 1'),
        Procedure(id: 'id2', name: 'related procedure 2'),
      ],
    );

    test('should return valid model from Json', () {
      final result = TranslatedConversationModel.fromJson(
        json: translationJson,
      );

      expect(result, translationModel);
    });

    test('toJson should return a valid json', () {
      final json = translationModel.toJson();
      expect(json, translationJson);
    });

    test('should return a valid entity when toEntity is called', () {
      final translation = TranslatedConversationModel.fromJson(
        json: translationJson,
      );
      final entity = translation.toEntity();
      expect(entity, translationEntity);
    });

    test('should return a valid model from entity', () {
      final result = TranslatedConversationModel.fromEntity(
        entity: translationEntity,
      );

      expect(result, translationModel);
    });
  }); // translated conversation model
}
