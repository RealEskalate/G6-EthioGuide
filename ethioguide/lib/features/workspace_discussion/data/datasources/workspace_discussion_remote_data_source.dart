import 'package:dio/dio.dart';
import '../models/discussion_model.dart';
import '../models/comment_model.dart';
import '../models/community_stats_model.dart';

/// Remote data source for workspace discussion operations
abstract class WorkspaceDiscussionRemoteDataSource {
  Future<CommunityStatsModel> getCommunityStats();
  Future<List<DiscussionModel>> getDiscussions({
    String? tag,
    String? category,
    String? filterType,
  });
  Future<DiscussionModel> createDiscussion({
    required String title,
    required String content,
    required List<String> tags,
    required String category,
  });
  Future<bool> likeDiscussion(String discussionId);
  Future<bool> reportDiscussion(String discussionId);
  Future<List<CommentModel>> getComments(String discussionId);
  Future<CommentModel> addComment({
    required String discussionId,
    required String content,
  });
  Future<bool> likeComment(String commentId);
  Future<bool> reportComment(String commentId);
}

class WorkspaceDiscussionRemoteDataSourceImpl implements WorkspaceDiscussionRemoteDataSource {
  final Dio dio;
  final String baseUrl;

  WorkspaceDiscussionRemoteDataSourceImpl({
    required this.dio,
    this.baseUrl = 'https://api.ethioguide.com',
  });

  @override
  Future<CommunityStatsModel> getCommunityStats() async {
    final response = await dio.get('$baseUrl/workspace/community/stats');
    if (response.statusCode == 200) {
      return CommunityStatsModel.fromJson(response.data as Map<String, dynamic>);
    }
    throw DioException(
      requestOptions: RequestOptions(path: '$baseUrl/workspace/community/stats'),
      response: response,
      error: 'Failed to fetch community stats',
    );
  }

  @override
  Future<List<DiscussionModel>> getDiscussions({
    String? tag,
    String? category,
    String? filterType,
  }) async {
    final response = await dio.get(
      '$baseUrl/workspace/discussions',
      queryParameters: {
        if (tag != null) 'tag': tag,
        if (category != null) 'category': category,
        if (filterType != null) 'filterType': filterType,
      },
    );
    if (response.statusCode == 200) {
      final data = response.data as List<dynamic>;
      return data.map((e) => DiscussionModel.fromJson(e as Map<String, dynamic>)).toList();
    }
    throw DioException(
      requestOptions: RequestOptions(path: '$baseUrl/workspace/discussions'),
      response: response,
      error: 'Failed to fetch discussions',
    );
  }

  @override
  Future<DiscussionModel> createDiscussion({
    required String title,
    required String content,
    required List<String> tags,
    required String category,
  }) async {
    final response = await dio.post(
      '$baseUrl/workspace/discussions',
      data: {
        'title': title,
        'content': content,
        'tags': tags,
        'category': category,
      },
    );
    if (response.statusCode == 200 || response.statusCode == 201) {
      return DiscussionModel.fromJson(response.data as Map<String, dynamic>);
    }
    throw DioException(
      requestOptions: RequestOptions(path: '$baseUrl/workspace/discussions'),
      response: response,
      error: 'Failed to create discussion',
    );
  }

  @override
  Future<bool> likeDiscussion(String discussionId) async {
    final response = await dio.post('$baseUrl/workspace/discussions/$discussionId/like');
    return (response.statusCode == 200 || response.statusCode == 204);
  }

  @override
  Future<bool> reportDiscussion(String discussionId) async {
    final response = await dio.post('$baseUrl/workspace/discussions/$discussionId/report');
    return (response.statusCode == 200 || response.statusCode == 204);
  }

  @override
  Future<List<CommentModel>> getComments(String discussionId) async {
    final response = await dio.get('$baseUrl/workspace/discussions/$discussionId/comments');
    if (response.statusCode == 200) {
      final data = response.data as List<dynamic>;
      return data.map((e) => CommentModel.fromJson(e as Map<String, dynamic>)).toList();
    }
    throw DioException(
      requestOptions: RequestOptions(path: '$baseUrl/workspace/discussions/$discussionId/comments'),
      response: response,
      error: 'Failed to fetch comments',
    );
  }

  @override
  Future<CommentModel> addComment({
    required String discussionId,
    required String content,
  }) async {
    final response = await dio.post(
      '$baseUrl/workspace/discussions/$discussionId/comments',
      data: {'content': content},
    );
    if (response.statusCode == 200 || response.statusCode == 201) {
      return CommentModel.fromJson(response.data as Map<String, dynamic>);
    }
    throw DioException(
      requestOptions: RequestOptions(path: '$baseUrl/workspace/discussions/$discussionId/comments'),
      response: response,
      error: 'Failed to add comment',
    );
  }

  @override
  Future<bool> likeComment(String commentId) async {
    final response = await dio.post('$baseUrl/workspace/comments/$commentId/like');
    return (response.statusCode == 200 || response.statusCode == 204);
  }

  @override
  Future<bool> reportComment(String commentId) async {
    final response = await dio.post('$baseUrl/workspace/comments/$commentId/report');
    return (response.statusCode == 200 || response.statusCode == 204);
  }
}
