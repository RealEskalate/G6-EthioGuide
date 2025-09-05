import 'package:dio/dio.dart';
import 'package:ethioguide/core/config/end_points.dart';
import 'package:ethioguide/core/error/exception.dart';
import 'package:ethioguide/core/network/network_info.dart';
import 'package:ethioguide/features/AI%20chat/data/models/conversation_model.dart';
import 'package:flutter/widgets.dart';

abstract class AiRemoteDatasource {
  Future<ConversationModel> sendQuery(String query);
  Future<List<ConversationModel>> getHistory();
  Future<String> translateContent(String content, String lang);
}

class AiRemoteDataSourceImpl implements AiRemoteDatasource {
  final Dio dio;
  final NetworkInfo networkInfo;

  AiRemoteDataSourceImpl({required this.dio, required this.networkInfo});

  Exception throwsException(int statusCode) {
    if (statusCode == 400) {
      return ServerException(
        message: 'Bad request. Please check your input.',
        statusCode: statusCode,
      );
    } else if (statusCode == 401) {
      return ServerException(
        message: 'Couldn\'t authenticate user. Please log in again.',
        statusCode: statusCode,
      );
    } else if (statusCode == 403) {
      return ServerException(
        message: 'Forbidden content. User doesn\'t have permission.',
        statusCode: statusCode,
      );
    } else if (statusCode == 404) {
      return ServerException(
        message: 'Requested resource not found.',
        statusCode: statusCode,
      );
    } else if (statusCode == 409) {
      return ServerException(
        message: 'Conflict detected. The request could not be completed.',
        statusCode: statusCode,
      );
    } else if (statusCode == 422) {
      return ServerException(
        message: 'Unprocessable entity. Validation failed.',
        statusCode: statusCode,
      );
    } else if (statusCode == 429) {
      return ServerException(
        message: 'Too many requests. Please try again later.',
        statusCode: statusCode,
      );
    } else if (statusCode == 500) {
      return ServerException(
        message: 'Internal server error. Please try again later.',
        statusCode: statusCode,
      );
    } else if (statusCode == 502) {
      return ServerException(
        message: 'Bad gateway. Received invalid response from upstream server.',
        statusCode: statusCode,
      );
    } else if (statusCode == 503) {
      return ServerException(
        message: 'Service unavailable. Please try again later.',
        statusCode: statusCode,
      );
    } else if (statusCode == 504) {
      return ServerException(
        message: 'Gateway timeout. The server took too long to respond.',
        statusCode: statusCode,
      );
    } else {
      return ServerException(
        message: 'Unexpected error occurred. Status code: $statusCode',
      );
    }
  }

  @override
  Future<ConversationModel> sendQuery(String query) async {
    /// Check network connecitvity first
    if (!(await networkInfo.isConnected)) {
      throw ServerException(
        message: 'No internet connection',
        statusCode: null,
      );
    }

    /// if device is online
    try {
      final response = await dio.post(
        EndPoints.sendQueryEndPoint,
        data: {'query': query},
      );
      final statusCode = response.statusCode;

      if (statusCode! >= 300) {
        // TODO: remove debug print
        debugPrint('Status code: $statusCode');

        debugPrint('##################################################');
        debugPrint('ServerException at AiRemoteDataSourceImpl sendQuery');
        debugPrint('##################################################');
        throw throwsException(statusCode);
      }

      return ConversationModel.fromJson(response.data);
    } on DioException catch (e) {
      // TODO: remove debug print
      debugPrint('##################################################');
      debugPrint('DioException at AiRemoteDataSourceImpl sendQuery');
      debugPrint('Exception: ${e.message}');
      debugPrint('Response: ${e.response?.data}');
      debugPrint('##################################################');

      throw ServerException(
        message: e.response?.data['message'] ?? 'Failed to send query',
        statusCode: e.response?.statusCode,
      );
    } on ServerException {
      rethrow;
    } catch (e) {
      // TODO: remove debug print
      debugPrint('##################################################');
      debugPrint('Unexpected Exception at AiRemoteDataSourceImpl sendQuery');
      debugPrint('Exception: $e');
      debugPrint('##################################################');

      throw ServerException(message: 'Unexpected error: $e', statusCode: null);
    }
  }

  @override
  Future<List<ConversationModel>> getHistory() async {
    /// Check network connectivity
    if (!(await networkInfo.isConnected)) {
      throw ServerException(
        message: 'No internet connection',
        statusCode: null,
      );
    }

    /// if device is online
    try {
      final response = await dio.get(EndPoints.getHistoryEndPoint);
      final statusCode = response.statusCode;

      if (statusCode! >= 300) {
        // TODO: remove debug print
        debugPrint('##################################################');
        debugPrint('ServerException at AiRemoteDataSourceImpl getHisotry');
        debugPrint('##################################################');
        throw throwsException(statusCode);
      }

      final List<dynamic> jsonList = response.data['history'];
      return jsonList.map((json) => ConversationModel.fromJson(json)).toList();
    } on DioException catch (e) {
      // TODO: remove debug print
      debugPrint('##################################################');
      debugPrint('DioException at AiRemoteDataSourceImpl getHistory');
      debugPrint('Exception: ${e.message}');
      debugPrint('Response: ${e.response?.data}');
      debugPrint('##################################################');

      throw ServerException(
        message: e.response?.data['message'] ?? 'Failed to fetch history',
        statusCode: e.response?.statusCode,
      );
    } on ServerException {
      rethrow;
    } catch (e) {
      // TODO: remove debug print
      debugPrint('##################################################');
      debugPrint('Unexpected Exception at AiRemoteDataSourceImpl getHistory');
      debugPrint('Exception: $e');
      debugPrint('##################################################');

      throw ServerException(message: 'Unexpected error: $e', statusCode: null);
    }
  }

  @override
  Future<String> translateContent(String content, String lang) async {
    /// Check network connectivity
    if (!(await networkInfo.isConnected)) {
      throw ServerException(
        message: 'No internet connection',
        statusCode: null,
      );
    }

    /// if device is online
    try {
      final response = await dio.post(
        EndPoints.translateContentEndPoint,
        data: {'content': content, 'lang': lang},
      );

      final statusCode = response.statusCode;
      if (statusCode! >= 300) {
        // TODO: remove debug print
        debugPrint('##################################################');
        debugPrint('ServerException at AiRemoteDataSourceImpl translate');
        debugPrint('##################################################');
        throw throwsException(statusCode);
      }
      return response.data['translated'] as String;
    } on DioException catch (e) {
      // TODO: remove debug print
      debugPrint('##################################################');
      debugPrint('DioException at AiRemoteDataSourceImpl translateContent');
      debugPrint('Exception: ${e.message}');
      debugPrint('Response: ${e.response?.data}');
      debugPrint('##################################################');

      throw ServerException(
        message: e.response?.data['message'] ?? 'Failed to translate',
        statusCode: e.response?.statusCode,
      );
    } on ServerException {
      rethrow;
    } catch (e) {
      // TODO: remove debug print
      debugPrint('##################################################');
      debugPrint(
        'Unexpected Exception at AiRemoteDataSourceImpl translateContent',
      );
      debugPrint('Exception: $e');
      debugPrint('##################################################');

      throw ServerException(message: 'Unexpected error: $e', statusCode: null);
    }
  }
}
