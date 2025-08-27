import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/AI%20chat/Domain/entities/conversation.dart';
import 'package:ethioguide/features/AI%20chat/Domain/repository/ai_repository.dart';
import 'package:dartz/dartz.dart';

class SendQuery {
  final AiRepository repository;

  SendQuery({required this.repository});

  Future<Either<Failure, Conversation>> call(String query) async {
    return await repository.sendQuery(query);
  }
}


