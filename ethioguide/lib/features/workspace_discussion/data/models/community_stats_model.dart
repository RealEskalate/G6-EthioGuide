import '../../domain/entities/community_stats.dart';

/// Data model for community statistics
class CommunityStatsModel extends CommunityStats {
  const CommunityStatsModel({
    required int totalMembers,
    required int totalDiscussions,
    required int activeToday,
    required List<String> trendingTags,
    
  }) : super(
         
          totalMembers: totalMembers,
          totalDiscussions: totalDiscussions,
          activeToday: activeToday,
          trendingTags: trendingTags,
         
        );

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


