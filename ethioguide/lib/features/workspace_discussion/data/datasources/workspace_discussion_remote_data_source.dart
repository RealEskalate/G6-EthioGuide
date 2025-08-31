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
    try {
      await Future.delayed(const Duration(milliseconds: 500));
      
      // Mock response data
      final mockData = {
        'totalMembers': 1200,
        'totalDiscussions': 3,
        'activeToday': 42,
        'trendingTags': ['#passport', '#renewal', '#business-license', '#id-card', '#timeline', '#documents', '#fees'],
      };
      
      return CommunityStatsModel.fromJson(mockData);
    } catch (e) {
      throw Exception('Failed to fetch community stats: $e');
    }
  }

  @override
  Future<List<DiscussionModel>> getDiscussions({
    String? tag,
    String? category,
    String? filterType,
  }) async {
    try {
      await Future.delayed(const Duration(milliseconds: 800));
      
      // Mock response data
      final mockData = [
        {
          'id': '1',
          'title': 'Business license application - timeline and costs?',
          'content': 'I\'m planning to start a small business and need to get a business license. Can anyone share their recent experience with the process? How long did it take and what were the actual costs involved?',
          'tags': ['Business Licenses', 'business', 'license', 'timeline', 'costs'],
          'category': 'Business',
          'createdAt': '2024-01-15T10:00:00.000Z',
          'createdBy': {
            'id': 'user1',
            'name': 'Yohannes Girma',
            'avatar': null,
            'role': null,
          },
          'likeCount': 13,
          'reportCount': 0,
          'commentsCount': 1,
          'isPinned': false,
        },
        {
          'id': '2',
          'title': 'Common mistakes when applying for government services',
          'content': 'Based on community feedback, here are the most common mistakes people make when applying for government services. Please read this before starting any application process to save time and avoid frustration.',
          'tags': ['General Help', 'guidelines', 'tips', 'common-mistakes'],
          'category': 'General',
          'createdAt': '2024-01-14T15:30:00.000Z',
          'createdBy': {
            'id': 'user2',
            'name': 'Community Manager',
            'avatar': null,
            'role': 'Moderator',
          },
          'likeCount': 89,
          'reportCount': 0,
          'commentsCount': 1,
          'isPinned': true,
        },
      ];
      
      return mockData.map((data) => DiscussionModel.fromJson(data)).toList();
    } catch (e) {
      throw Exception('Failed to fetch discussions: $e');
    }
  }

  @override
  Future<DiscussionModel> createDiscussion({
    required String title,
    required String content,
    required List<String> tags,
    required String category,
  }) async {
    try {
      await Future.delayed(const Duration(milliseconds: 1000));
      
      // Mock response data
      final mockData = {
        'id': DateTime.now().millisecondsSinceEpoch.toString(),
        'title': title,
        'content': content,
        'tags': tags,
        'category': category,
        'createdAt': DateTime.now().toIso8601String(),
        'createdBy': {
          'id': 'currentUser',
          'name': 'Current User',
          'avatar': null,
          'role': null,
        },
        'likeCount': 0,
        'reportCount': 0,
        'commentsCount': 0,
        'isPinned': false,
      };
      
      return DiscussionModel.fromJson(mockData);
    } catch (e) {
      throw Exception('Failed to create discussion: $e');
    }
  }

  @override
  Future<bool> likeDiscussion(String discussionId) async {
    try {
      await Future.delayed(const Duration(milliseconds: 300));
      return true;
    } catch (e) {
      throw Exception('Failed to like discussion: $e');
    }
  }

  @override
  Future<bool> reportDiscussion(String discussionId) async {
    try {
      await Future.delayed(const Duration(milliseconds: 300));
      return true;
    } catch (e) {
      throw Exception('Failed to report discussion: $e');
    }
  }

  @override
  Future<List<CommentModel>> getComments(String discussionId) async {
    try {
      await Future.delayed(const Duration(milliseconds: 600));
      
      // Mock response data
      final mockData = [
        {
          'id': 'comment1',
          'discussionId': discussionId,
          'content': 'yes',
          'createdAt': DateTime.now().toIso8601String(),
          'createdBy': {
            'id': 'currentUser',
            'name': 'You',
            'avatar': null,
            'role': null,
          },
          'likeCount': 0,
          'reportCount': 0,
        },
      ];
      
      return mockData.map((data) => CommentModel.fromJson(data)).toList();
    } catch (e) {
      throw Exception('Failed to fetch comments: $e');
    }
  }

  @override
  Future<CommentModel> addComment({
    required String discussionId,
    required String content,
  }) async {
    try {
      await Future.delayed(const Duration(milliseconds: 500));
      
      // Mock response data
      final mockData = {
        'id': DateTime.now().millisecondsSinceEpoch.toString(),
        'discussionId': discussionId,
        'content': content,
        'createdAt': DateTime.now().toIso8601String(),
        'createdBy': {
          'id': 'currentUser',
          'name': 'Current User',
          'avatar': null,
          'role': null,
        },
        'likeCount': 0,
        'reportCount': 0,
      };
      
      return CommentModel.fromJson(mockData);
    } catch (e) {
      throw Exception('Failed to add comment: $e');
    }
  }

  @override
  Future<bool> likeComment(String commentId) async {
    try {
      await Future.delayed(const Duration(milliseconds: 300));
      return true;
    } catch (e) {
      throw Exception('Failed to like comment: $e');
    }
  }

  @override
  Future<bool> reportComment(String commentId) async {
    try {
      await Future.delayed(const Duration(milliseconds: 300));
      return true;
    } catch (e) {
      throw Exception('Failed to report comment: $e');
    }
  }
}
