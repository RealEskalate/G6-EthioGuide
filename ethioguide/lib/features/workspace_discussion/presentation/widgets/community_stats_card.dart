import 'package:flutter/material.dart';
import '../../domain/entities/community_stats.dart';

/// Widget that displays community statistics overview
class CommunityStatsCard extends StatelessWidget {
  final CommunityStats communityStats;

  const CommunityStatsCard({
    super.key,
    required this.communityStats,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 2,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
      ),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Row(
          children: [
            Expanded(
              child: _StatItem(
                icon: Icons.people,
                value: '${(communityStats.totalMembers / 1000).toStringAsFixed(1)}k',
                label: 'Members',
                color: Colors.blue,
              ),
            ),
            Expanded(
              child: _StatItem(
                icon: Icons.chat_bubble,
                value: '${communityStats.totalDiscussions}',
                label: 'Discussions',
                color: Colors.green,
              ),
            ),
            Expanded(
              child: _StatItem(
                icon: Icons.trending_up,
                value: '${communityStats.activeToday}',
                label: 'Active Today',
                color: Colors.orange,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _StatItem extends StatelessWidget {
  final IconData icon;
  final String value;
  final String label;
  final Color color;

  const _StatItem({
    required this.icon,
    required this.value,
    required this.label,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Icon(
          icon,
          color: color,
          size: 24,
        ),
        const SizedBox(height: 4),
        Text(
          value,
          style: Theme.of(context).textTheme.titleMedium?.copyWith(
            fontWeight: FontWeight.w600,
            color: color,
          ),
        ),
        const SizedBox(height: 2),
        Text(
          label,
          style: Theme.of(context).textTheme.bodySmall?.copyWith(
            color: Colors.grey[600],
          ),
          textAlign: TextAlign.center,
        ),
      ],
    );
  }
}
