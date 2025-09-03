import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/AI%20chat/Domain/entities/conversation.dart';
import 'package:ethioguide/features/AI%20chat/Domain/repository/ai_repository.dart';
import 'package:dartz/dartz.dart';

class GetHistory {
  final AiRepository repository;

  GetHistory({required this.repository});

  Future<Either<Failure, List<Conversation>>> call() async {
    return await repository.getHistory();
  }
}


