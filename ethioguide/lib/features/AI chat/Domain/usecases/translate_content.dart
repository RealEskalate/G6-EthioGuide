import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/AI%20chat/Domain/entities/conversation.dart';
import 'package:ethioguide/features/AI%20chat/Domain/entities/translated_conversation.dart';
import 'package:ethioguide/features/AI%20chat/Domain/repository/ai_repository.dart';
import 'package:dartz/dartz.dart';
import 'package:ethioguide/features/AI%20chat/data/models/translated_conversation_model.dart';

class TranslateContent {
  final AiRepository repository;

  TranslateContent({required this.repository});

  Future<Either<Failure, TranslatedConversation>> call({required String response, required List<Procedure?> procedures}) async {
    return await repository.translateContent(TranslatedConversationModel(response: response, procedures: procedures));
  }
}
