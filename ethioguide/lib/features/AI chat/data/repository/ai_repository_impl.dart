import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/exception.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/core/network/network_info.dart';
import 'package:ethioguide/features/AI%20chat/Domain/entities/conversation.dart';
import 'package:ethioguide/features/AI%20chat/Domain/entities/translated_conversation.dart';
import 'package:ethioguide/features/AI%20chat/Domain/repository/ai_repository.dart';
import 'package:ethioguide/features/AI%20chat/data/datasources/ai_local_datasource.dart';
import 'package:ethioguide/features/AI%20chat/data/datasources/ai_remote_datasource.dart';
import 'package:ethioguide/features/AI%20chat/data/models/translated_conversation_model.dart';

class AiRepositoryImpl implements AiRepository {
  final AiRemoteDatasource remoteDatasource;
  final AiLocalDatasource localDatasource;
  final NetworkInfo networkInfo;

  AiRepositoryImpl({
    required this.remoteDatasource,
    required this.localDatasource,
    required this.networkInfo,
  });

  @override
  Future<Either<Failure, Conversation>> sendQuery(String query) async {
    try {
      // check network connectivity
      if (!(await networkInfo.isConnected)) {
        return Left(NetworkFailure(message: 'No internet connection'));
      }

      final conversation = await remoteDatasource.sendQuery(query);
      // Cache the new conversation
      try {
        final currentHistory = await localDatasource.getCachedHistory();
        final updatedHistory = [conversation, ...currentHistory];
        await localDatasource.cacheHistory(updatedHistory);
      } on CacheException {
        //TODO: log this
      }
      return Right(conversation);
    } on ServerException catch (e) {
      return Left(ServerFailure(message: e.message));
    } catch (e) {
      return Left(ServerFailure(message: 'Unexpected error: $e'));
    }
  }

  @override
  Future<Either<Failure, List<Conversation>>> getHistory() async {
    try {
      final localHistory = await localDatasource.getCachedHistory();
      if (localHistory.isNotEmpty) {
        return Right(localHistory);
      }
    } on CacheException {
      // Ignore → will fallback to remote
    }

    // If cache empty or missing → go to remote
    if (await networkInfo.isConnected) {
      try {
        final remoteHistory = await remoteDatasource.getHistory();
        await localDatasource.cacheHistory(remoteHistory);
        return Right(remoteHistory);
      } catch (e) {
        return Right([]);
      }
    } else {
      return Right([]);
    }
  }

  @override
  Future<Either<Failure, TranslatedConversation>> translateContent(
    TranslatedConversationModel conversation,
  ) async {
    try {
      if (!(await networkInfo.isConnected)) {
        return Left(NetworkFailure(message: 'No internet connection'));
      }
      final translated = await remoteDatasource.translateContent(conversation);
      return Right(translated);
    } on ServerException catch (e) {
      return Left(ServerFailure(message: e.message));
    } catch (e) {
      return Left(ServerFailure(message: 'Unexpected error: $e'));
    }
  }

  // @override
  // Future<Either<Failure, String>> translateContent(
  //   String content,
  //   String lang,
  // ) async {
  //   try {
  //     if (!(await networkInfo.isConnected)) {
  //       return Left(NetworkFailure(message: 'No internet connection'));
  //     }
  //     final translated = await remoteDatasource.translateContent(content, lang);
  //     return Right(translated);
  //   } on ServerException catch (e) {
  //     return Left(ServerFailure(message: e.message));
  //   } catch (e) {
  //     return Left(ServerFailure(message: 'Unexpected error: $e'));
  //   }
  // }
}
