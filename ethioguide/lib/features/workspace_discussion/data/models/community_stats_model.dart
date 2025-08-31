import '../../domain/entities/community_stats.dart';

/// Data model for community statistics
class CommunityStatsModel extends CommunityStats {
  const CommunityStatsModel({
    required super.totalMembers,
    required super.totalDiscussions,
    required super.activeToday,
    required super.trendingTags,
  });

  factory CommunityStatsModel.fromJson(Map<String, dynamic> json) {
    return CommunityStatsModel(
      totalMembers: json['totalMembers'] as int,
      totalDiscussions: json['totalDiscussions'] as int,
      activeToday: json['activeToday'] as int,
      trendingTags: (json['trendingTags'] as List<dynamic>).cast<String>(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'totalMembers': totalMembers,
      'totalDiscussions': totalDiscussions,
      'activeToday': activeToday,
      'trendingTags': trendingTags,
    };
  }
}
