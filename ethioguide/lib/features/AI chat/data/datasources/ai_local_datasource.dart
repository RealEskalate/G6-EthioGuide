import 'dart:convert';

import 'package:ethioguide/core/config/cache_key_names.dart';
import 'package:ethioguide/core/error/exception.dart';
import 'package:ethioguide/features/AI%20chat/data/models/conversation_model.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

abstract class AiLocalDatasource {
  Future<List<ConversationModel>> getCachedHistory();
  Future<void> cacheHistory(List<ConversationModel> history);
}

class AiLocalDataSourceImpl implements AiLocalDatasource {
  final FlutterSecureStorage secureStorage;
  AiLocalDataSourceImpl({required this.secureStorage});

  @override
  Future<List<ConversationModel>> getCachedHistory() async {
    try {
      final jsonString = await secureStorage.read(
        key: CacheKeyNames.aiHistoryKey,
      );

      if (jsonString != null) {
        final List<dynamic> jsonList = jsonDecode(jsonString);
        return jsonList
            .map((json) => ConversationModel.fromJson(json))
            .toList();
      }

      return []; // returh Empty if no messages were saved
    } on CacheException {
      rethrow;
    } catch (e) {
      // TODO: remove debug print
      debugPrint('##################################################');
      debugPrint('CacheException at AiLocalDataSourceImpl getCachedHistory');
      debugPrint('Exception: $e');
      debugPrint('##################################################');
      throw CacheException(message: 'Failed to read cached history: $e');
    }
  }

  @override
  Future<void> cacheHistory(List<ConversationModel> history) async {
    try {
      final jsonString = jsonEncode(history.map((e) => e.toJson()).toList());
      await secureStorage.write(
        key: CacheKeyNames.aiHistoryKey,
        value: jsonString,
      );
    } catch (e) {
      // TODO: remove debug print
      debugPrint('##################################################');
      debugPrint('CacheException at AiLocalDataSourceImpl CacheHistory');
      debugPrint('Exception: $e');
      debugPrint('##################################################');
      throw CacheException(message: 'Failed to cache hisory: $e');
    }
  }
}
