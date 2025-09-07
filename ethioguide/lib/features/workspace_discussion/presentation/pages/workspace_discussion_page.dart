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
    trendingTags: ['license', 'process', 'All'],
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

  title:  Text(
              'Community Discussions',
              style: Theme.of(context)
                  .textTheme
                  .titleLarge
                  ?.copyWith(fontWeight: FontWeight.bold , fontSize: 15),
              overflow: TextOverflow.ellipsis, // truncate long text
            ),
  actions: [
    Padding(
      padding: const EdgeInsets.only(right: 8.0),
      child: ElevatedButton.icon(
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
          TextField(
            onSubmitted: (query) {
              print(query);
              context.read<WorkspaceDiscussionBloc>().add(
                FetchDiscussions(
                  filterType: query, // pass search term here
                ),
              );
            },
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
              if (tag == 'All') {
                // no filter
                context.read<WorkspaceDiscussionBloc>().add(
                  const FetchDiscussions(),
                );
                return;
              }
              // Implement tag filter logic here
              context.read<WorkspaceDiscussionBloc>().add(
                FetchDiscussions(tag: tag),
              );

              // no-op for UI preview
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

          BlocBuilder<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
            builder: (context, state) {
              if (state is DiscussionLoading) {
                return const Center(child: CircularProgressIndicator());
              } else if (state is DiscussionLoaded) {
                return ListView.builder(
                  shrinkWrap:
                      true, // ✅ Important when inside SingleChildScrollView
                  physics:
                      const NeverScrollableScrollPhysics(), // ✅ Disable nested scrolling
                  itemCount: state.discussions.length,
                  itemBuilder: (context, index) {
                    final discussion = state.discussions[index];
                    return Padding(
                      padding: const EdgeInsets.only(bottom: 16.0),
                      child: DiscussionCard(
                        discussions: discussion,
                        onTap: () {},
                           
                        onLike: () {},
                        onReport: () {},
                        onComment: () {}
                          
                      ),
                    );
                  },
                );
              } else if (state is DiscussionError) {
                return Center(child: Text("Error: ${state.message}"));
              }
              return const SizedBox.shrink();
            },
          ),

          /*  ...discussions.map(
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
          ), */
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

  /* void _navigateToDiscussionDetail(
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
  } */
}
