import 'package:ethioguide/features/procedure/presentation/widgets/procedure_detail_header.dart';
import 'package:ethioguide/features/procedure/presentation/widgets/workspace_summary_cards.dart';
import 'package:flutter/material.dart';
import '../../domain/entities/community_stats.dart';

/// Widget that displays community statistics overview
class CommunityStatsCard extends StatelessWidget {
  final CommunityStats communityStats;

  const CommunityStatsCard({super.key, required this.communityStats});

  @override
  Widget build(BuildContext context) {
    return 
      SingleChildScrollView(
        scrollDirection: Axis.horizontal,
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            InfoCard(
              title: 'Members',
              value: communityStats.totalMembers.toString(),
              icon: Icons.people,
            ),
            InfoCard(
              title: 'Discussions',
              value: communityStats.totalDiscussions.toString(),
              icon: Icons.schedule,
            ),
            InfoCard(
              title: 'Active Today',
              value: communityStats.activeToday.toString(),
              icon: Icons.check_circle,
            ),
          ],
        ),
      );
    
  }
}
