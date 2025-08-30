import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/exception.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/core/network/network_info.dart';
import 'package:ethioguide/features/AI%20chat/Domain/entities/conversation.dart';
import 'package:ethioguide/features/AI%20chat/Domain/repository/ai_repository.dart';
import 'package:ethioguide/features/AI%20chat/data/datasources/ai_local_datasource.dart';
import 'package:ethioguide/features/AI%20chat/data/datasources/ai_remote_datasource.dart';

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
    if (await networkInfo.isConnected) {
      try {
        final remoteHistory = await remoteDatasource.getHistory();
        try {
          await localDatasource.cacheHistory(remoteHistory);
        } on CacheException {
          // If caching fails, return the sucessful remote result
          //TODO: log this
        }
        return Right(remoteHistory);
      } on ServerException catch (e) {
        return Left(ServerFailure(message: e.message));
      } catch (e) {
        return Left(ServerFailure(message: 'Unexpected error: $e'));
      }
    } else {
      // Get cached chat if offline
      try {
        final localHistory = await localDatasource.getCachedHistory();
        return Right(localHistory);
      } on CacheException catch (e) {
        return Left(CachedFailure(message: e.message));
      } catch (e) {
        return Left(CachedFailure(message: 'Unexpected cache error: $e'));
      }
    }
  }

  @override
  Future<Either<Failure, String>> translateContent(
    String content,
    String lang,
  ) async {
    try {
      if (!(await networkInfo.isConnected)) {
        return Left(NetworkFailure(message: 'No internet connection'));
      }
      final translated = await remoteDatasource.translateContent(content, lang);
      return Right(translated);
    } on ServerException catch (e) {
      return Left(ServerFailure(message: e.message));
    } catch (e) {
      return Left(ServerFailure(message: 'Unexpected error: $e'));
    }
  }
}
