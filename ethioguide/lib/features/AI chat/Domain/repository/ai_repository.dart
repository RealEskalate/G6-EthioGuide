import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/AI%20chat/Domain/entities/translated_conversation.dart';
import 'package:ethioguide/features/AI%20chat/data/models/translated_conversation_model.dart';
import '../entities/conversation.dart';

abstract class AiRepository {
  Future<Either<Failure, Conversation>> sendQuery(String query);
  Future<Either<Failure, List<Conversation>>> getHistory();
  Future<Either<Failure, TranslatedConversation>> translateContent(TranslatedConversationModel conversation);
}