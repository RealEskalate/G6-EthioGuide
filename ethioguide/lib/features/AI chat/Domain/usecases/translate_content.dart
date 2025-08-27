import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/AI%20chat/Domain/repository/ai_repository.dart';
import 'package:dartz/dartz.dart';

class TranslateContent {
  final AiRepository repository;

  TranslateContent({required this.repository});

  Future<Either<Failure, String>> call(String content, String lang) async {
    return await repository.translateContent(content, lang);
  }
}


