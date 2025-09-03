import 'package:ethioguide/features/AI%20chat/data/models/conversation_model.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:ethioguide/features/AI chat/Domain/entities/conversation.dart';

void main() {
  group('ProcedureModel', () {
    const procedureJson = {'id': '1', 'name': 'Passport Application'};

    test('fromJson should return valid model', () {
      final result = ProcedureModel.fromJson(procedureJson);

      expect(result, isA<ProcedureModel>());
      expect(result.id, '1');
      expect(result.name, 'Passport Application');
    });

    test('toJson should return valid json', () {
      const procedure = ProcedureModel(id: '1', name: 'Passport Application');

      final json = procedure.toJson();

      expect(json, procedureJson);
    });

    test('toEntity should return valid entity', () {
      const procedure = ProcedureModel(id: '1', name: 'Passport Application');

      final entity = procedure.toEntity();

      expect(entity, isA<Procedure>());
      expect(entity.id, '1');
      expect(entity.name, 'Passport Application');
    });
  });

  group('ConversationModel', () {
    final conversationJson = {
      'id': 'id',
      'request': 'How to get a passport',
      'response': 'Steps to get a Passport...',
      'source': 'official',
      'procedures': [
        {'id': '1', 'name': 'Passport Application'},
      ],
    };

    test('fromJson should return valid model', () {
      final result = ConversationModel.fromJson(conversationJson);

      expect(result, isA<ConversationModel>());
      expect(result.request, 'How to get a passport');
      expect(result.response, 'Steps to get a Passport...');
      expect(result.source, 'official');
      expect(result.procedures.first, isA<ProcedureModel>());
    });

    test('toJson should return valid json', () {
      final conversation = ConversationModel(
        id: 'id',
        request: 'How to get a passport',
        response: 'Steps to get a Passport...',
        source: 'official',
        procedures: [
          const ProcedureModel(id: '1', name: 'Passport Application'),
        ],
      );

      final json = conversation.toJson();

      expect(json, conversationJson);
    });

    test('toEntity should return valid entity', () {
      final conversation = ConversationModel.fromJson(conversationJson);

      final entity = conversation.toEntity();

      expect(entity, isA<Conversation>());
      expect(entity.request, 'How to get a passport');
      expect(entity.procedures.first, isA<Procedure>());
    });
  });
}
