import 'package:dio/dio.dart';
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
  // TODO: Replace with actual base URL
  final String baseUrl = 'https://api.ethioguide.com';

  AiRemoteDataSourceImpl({required this.dio, required this.networkInfo}) {
    dio.options
      ..baseUrl = baseUrl
      ..headers = {'Content-Type': 'application/json'};
  }

  @override
  Future<ConversationModel> sendQuery(String query) async {
    /// if device is offline
    if (!(await networkInfo.isConnected)) {
      throw ServerException(
        message: 'No internet connection',
        statusCode: null,
      );
    }

    /// if device is online
    try {
      final response = await dio.post('/ai/guide', data: {'query': query});
      return ConversationModel.fromJson(response.data);
    } on DioException catch (e) {
      // TODO: remove debug print
      debugPrint('##################################################');
      debugPrint('DioException at ai_remote_repository sendQuery function');
      debugPrint('Exception: ${e.message}');
      debugPrint('##################################################');

      throw ServerException(
        message: e.response?.data['message'] ?? 'Failed to send query',
        statusCode: e.response?.statusCode,
      );
    } catch (e) {
      // TODO: remove debug print
      debugPrint('##################################################');
      debugPrint('Exception at ai_remote_repository sendQuery function');
      debugPrint('Exception: $e');
      debugPrint('##################################################');

      throw ServerException(message: 'Unexpected error: $e', statusCode: null);
    }
  }

  @override
  Future<List<ConversationModel>> getHistory() async {
    /// if device is offline
    if (!(await networkInfo.isConnected)) {
      throw ServerException(
        message: 'No internet connection',
        statusCode: null,
      );
    }

    /// if device is online
    try {
      final response = await dio.get('/ai/history');
      final List<dynamic> jsonList = response.data;
      return jsonList.map((json) => ConversationModel.fromJson(json)).toList();
    } on DioException catch (e) {
      // TODO: remove debug print
      debugPrint('##################################################');
      debugPrint('DioException at ai_remote_repository getHistory function');
      debugPrint('Exception: ${e.message}');
      debugPrint('##################################################');

      throw ServerException(
        message: e.response?.data['message'] ?? 'Failed to fetch history',
        statusCode: e.response?.statusCode,
      );
    } catch (e) {
      // TODO: remove debug print
      debugPrint('##################################################');
      debugPrint('DioException at ai_remote_repository getHistory function');
      debugPrint('Exception: $e');
      debugPrint('##################################################');

      throw ServerException(message: 'Unexpected error: $e', statusCode: null);
    }
  }

  @override
  Future<String> translateContent(String content, String lang) async {
    /// if device is offline
    if (!(await networkInfo.isConnected)) {
      throw ServerException(
        message: 'No internet connection',
        statusCode: null,
      );
    }

    /// if device is online
    try {
      final response = await dio.post(
        '/translate',
        data: {'content': content, 'lang': lang},
      );
      return response.data['translated'] as String;
    } on DioException catch (e) {
      // TODO: remove debug print
      debugPrint('##################################################');
      debugPrint(
        'DioException at ai_remote_repository translateContent function',
      );
      debugPrint('Exception: ${e.message}');
      debugPrint('##################################################');

      throw ServerException(
        message: e.response?.data['message'] ?? 'Failed to translate',
        statusCode: e.response?.statusCode,
      );
    } catch (e) {
      // TODO: remove debug print
      debugPrint('##################################################');
      debugPrint('DioException at ai_remote_repository translateContent function');
      debugPrint('Exception: $e');
      debugPrint('##################################################');

      throw ServerException(message: 'Unexpected error: $e', statusCode: null);
    }
  }
}
