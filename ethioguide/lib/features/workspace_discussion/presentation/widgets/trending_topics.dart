import 'package:ethioguide/core/config/app_color.dart';
import 'package:flutter/material.dart';
import '../../domain/entities/community_stats.dart';

/// Widget that displays trending topics
class TrendingTopics extends StatelessWidget {
  final CommunityStats communityStats;
  final Function(String)? onTagTap;

  const TrendingTopics({
    super.key,
    required this.communityStats,
    this.onTagTap,
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
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Text(
                  'Trending Topics',
                  style: Theme.of(context).textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const Spacer(),
                Icon(
                  Icons.trending_up,
                  color: Colors.orange[600],
                  size: 20,
                ),
              ],
            ),
            const SizedBox(height: 12),
            Wrap(
              spacing: 8,
              runSpacing: 8,
              children: communityStats.trendingTags.map((tag) {
                return GestureDetector(
                  onTap: () => onTagTap?.call(tag),
                  child: Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 12,
                      vertical: 6,
                    ),
                    decoration: BoxDecoration(
                      color: AppColors.darkGreenColor.withOpacity(0.10),
                      borderRadius: BorderRadius.circular(20),
                      border: Border.all(
                        color: AppColors.darkGreenColor.withOpacity(0.10)!,
                        width: 1,
                      ),
                    ),
                    child: Text(
                     '#${tag}',
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                  ),
                );
              }).toList(),
            ),
          ],
        ),
      ),
    );
  }
}
