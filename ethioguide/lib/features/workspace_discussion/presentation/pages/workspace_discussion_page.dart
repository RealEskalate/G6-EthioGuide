import 'package:ethioguide/core/config/app_color.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/community_stats.dart';
import 'package:ethioguide/features/workspace_discussion/domain/entities/user.dart';
import 'package:ethioguide/features/workspace_discussion/presentation/bloc/worspace_discustion_state.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../domain/entities/discussion.dart';
import '../bloc/workspace_discussion_bloc.dart';
import '../widgets/community_stats_card.dart';
import '../widgets/trending_topics.dart';
import '../widgets/filter_controls.dart';
import '../widgets/discussion_card.dart';
import 'create_discussion_page.dart';
import 'discussion_detail_page.dart';

/// Main page for workspace community discussions
class WorkspaceDiscussionPage extends StatefulWidget {
  const WorkspaceDiscussionPage({super.key});

  @override
  State<WorkspaceDiscussionPage> createState() =>
      _WorkspaceDiscussionPageState();
}

class _WorkspaceDiscussionPageState extends State<WorkspaceDiscussionPage> {
  String? selectedCategory;
  String? selectedFilter = 'recent';
  final List<String> categories = [
    'Business',
    'General',
    'Travel',
    'Transportation',
  ];

  // Dummy data for UI-only preview
  final CommunityStats communityStats = CommunityStats(
    totalMembers: 100,
    totalDiscussions: 150,
    activeToday: 10,
    trendingTags: ['licensing', 'process'],
  );

  final discussions = List.generate(
    4,
    (i) => Discussion(
      id: 'd$i',
      title: 'How to complete step ${i + 1}?',
      content: 'I need help with the requirements for step ${i + 1}. Any tips?',
      tags: ['licensing', 'process'],
      category: i % 2 == 0 ? 'General' : 'Business',
      createdAt: DateTime.now().subtract(Duration(days: i)),
      createdBy: const User(id: 'u1', name: 'Test User'),
      likeCount: i * 3,
      reportCount: 0,
      commentsCount: i + 1,
    ),
  );

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.maybePop(context),
        ),
        title: Row(
          children: [
            Icon(Icons.chat_bubble_outline, color: Colors.teal[600], size: 28),
            const SizedBox(width: 12),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Community Discussions',
                  style: Theme.of(
                    context,
                  ).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold),
                ),
                Text(
                  'Share knowledge and get help',
                  style: Theme.of(
                    context,
                  ).textTheme.bodySmall?.copyWith(color: Colors.grey[600]),
                ),
              ],
            ),
          ],
        ),
        actions: [
          ElevatedButton.icon(
            onPressed: () => _navigateToCreateDiscussion(context),
            icon: const Icon(Icons.add, color: Colors.white),
            label: const Text(
              'New Post',
              style: TextStyle(color: Colors.white),
            ),
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.darkGreenColor,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(20),
              ),
            ),
          ),
        ],
      ),
      body: _buildContent(context),
    );
  }

  Widget _buildContent(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Search bar
          TextField(
            decoration: InputDecoration(
              hintText: 'Search discussions and procedures...',
              prefixIcon: const Icon(Icons.search),
              filled: true,
              fillColor: Colors.grey[100],
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(30),
                borderSide: BorderSide.none,
              ),
              contentPadding: const EdgeInsets.symmetric(
                vertical: 0,
                horizontal: 16,
              ),
            ),
          ),
          const SizedBox(height: 20),

          // Community stats card (dummy)
          CommunityStatsCard(communityStats: communityStats),
          const SizedBox(height: 20),

          // Trending topics (dummy uses communityStats)
          TrendingTopics(
            communityStats: communityStats,
            onTagTap: (tag) {
              // no-op for UI preview
            },
          ),
          const SizedBox(height: 20),

          // Filter controls
          FilterControls(
            selectedCategory: selectedCategory,
            selectedFilter: selectedFilter,
            categories: categories,
            onCategoryChanged: (category) {
              setState(() {
                selectedCategory = category;
              });
            },
            onFilterChanged: (filter) {
              setState(() {
                selectedFilter = filter;
              });
            },
          ),
          const SizedBox(height: 20),

          // Discussions list
          Text(
            'Discussions',
            style: Theme.of(
              context,
            ).textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w600),
          ),
          const SizedBox(height: 12),

          ...discussions.map(
            (discussion) => Padding(
              padding: const EdgeInsets.only(bottom: 16.0),
              child: DiscussionCard(
                discussion: discussion,
                onTap: () => _navigateToDiscussionDetail(context, discussion),
                onLike: () {},
                onReport: () {},
                onComment: () =>
                    _navigateToDiscussionDetail(context, discussion),
              ),
            ),
          ),
        ],
      ),
    );
  }

  void _navigateToCreateDiscussion(BuildContext context) {
    showDialog(
      context: context,
      barrierDismissible: true, // tap outside to close
      builder: (context) {
        return CreateDiscussionPage();
      },
    );
  }

  void _navigateToDiscussionDetail(
    BuildContext context,
    Discussion discussion,
  ) {
    showDialog(
      context: context,
      barrierDismissible: true, // tap outside to close
      builder: (context) {
        return DiscussionDetailPage(discussion: discussion);
      },
    );

    /*  Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => DiscussionDetailPage(discussion: discussion),
      ),
    ); */
  }
}
