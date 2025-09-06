import 'package:equatable/equatable.dart';

/// Domain entity representing community statistics
class CommunityStats extends Equatable {
  final int totalMembers;
  final int totalDiscussions;
  final int activeToday;
  final List<String> trendingTags;

  const CommunityStats({
    required this.totalMembers,
    required this.totalDiscussions,
    required this.activeToday,
    required this.trendingTags,
  });

  @override
  List<Object?> get props => [
        totalMembers,
        totalDiscussions,
        activeToday,
        trendingTags,
      ];
}
