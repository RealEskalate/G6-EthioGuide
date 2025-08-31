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
  State<WorkspaceDiscussionPage> createState() => _WorkspaceDiscussionPageState();
}

class _WorkspaceDiscussionPageState extends State<WorkspaceDiscussionPage> {
  String? selectedCategory;
  String? selectedFilter = 'recent';
  final List<String> categories = ['Business', 'General', 'Travel', 'Transportation'];

  @override
  void initState() {
    super.initState();
    // Load initial data
    context.read<WorkspaceDiscussionBloc>().add(FetchCommunityStats());
    context.read<WorkspaceDiscussionBloc>().add(FetchDiscussions());
  }

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
            Icon(
              Icons.chat_bubble,
              color: Colors.teal[600],
              size: 28,
            ),
            const SizedBox(width: 12),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Community Discussions',
                  style: Theme.of(context).textTheme.titleLarge?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
                ),
                Text(
                  'Share knowledge and get help',
                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    color: Colors.grey[600],
                  ),
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
              backgroundColor: Colors.blue[600],
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(20),
              ),
            ),
          ),
        ],
      ),
      body: BlocConsumer<WorkspaceDiscussionBloc, WorkspaceDiscussionState>(
        listener: (context, state) {
          if (state is ActionSuccess) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.message),
                backgroundColor: Colors.green,
              ),
            );
          } else if (state is ActionFailure) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.message),
                backgroundColor: Colors.red,
              ),
            );
          }
        },
        builder: (context, state) {
          if (state is DiscussionLoading) {
            return const Center(child: CircularProgressIndicator());
          } else if (state is DiscussionError) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(
                    Icons.error_outline,
                    size: 64,
                    color: Colors.grey[400],
                  ),
                  const SizedBox(height: 16),
                  Text(
                    'Error loading discussions',
                    style: Theme.of(context).textTheme.titleMedium,
                  ),
                  const SizedBox(height: 8),
                  Text(
                    state.message,
                    style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                      color: Colors.grey[600],
                    ),
                    textAlign: TextAlign.center,
                  ),
                  const SizedBox(height: 16),
                  ElevatedButton(
                    onPressed: () {
                      context.read<WorkspaceDiscussionBloc>().add(FetchDiscussions());
                    },
                    child: const Text('Retry'),
                  ),
                ],
              ),
            );
          } else if (state is DiscussionLoaded) {
            return _buildContent(context, state);
          }
          
          return const Center(child: Text('No discussions available'));
        },
      ),
    );
  }

  Widget _buildContent(BuildContext context, DiscussionLoaded state) {
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

          // Community stats card
          if (state.communityStats != null) ...[
            CommunityStatsCard(communityStats: state.communityStats!),
            const SizedBox(height: 20),
          ],

          // Trending topics
          if (state.communityStats != null) ...[
            TrendingTopics(
              communityStats: state.communityStats!,
              onTagTap: (tag) {
                // Filter by tag
                context.read<WorkspaceDiscussionBloc>().add(
                  FetchDiscussions(tag: tag),
                );
              },
            ),
            const SizedBox(height: 20),
          ],

          // Filter controls
          FilterControls(
            selectedCategory: selectedCategory,
            selectedFilter: selectedFilter,
            categories: categories,
            onCategoryChanged: (category) {
              setState(() {
                selectedCategory = category;
              });
              context.read<WorkspaceDiscussionBloc>().add(
                FetchDiscussions(category: category),
              );
            },
            onFilterChanged: (filter) {
              setState(() {
                selectedFilter = filter;
              });
              context.read<WorkspaceDiscussionBloc>().add(
                FetchDiscussions(filterType: filter),
              );
            },
          ),
          const SizedBox(height: 20),

          // Discussions list
          Text(
            'Discussions',
            style: Theme.of(context).textTheme.titleMedium?.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 12),
          
          ...state.discussions.map((discussion) => Padding(
            padding: const EdgeInsets.only(bottom: 16.0),
            child: DiscussionCard(
              discussion: discussion,
              onTap: () => _navigateToDiscussionDetail(context, discussion),
              onLike: () {
                context.read<WorkspaceDiscussionBloc>().add(
                  LikeDiscussion(discussion.id),
                );
              },
              onReport: () {
                context.read<WorkspaceDiscussionBloc>().add(
                  ReportDiscussion(discussion.id),
                );
              },
              onComment: () => _navigateToDiscussionDetail(context, discussion),
            ),
          )),
        ],
      ),
    );
  }

  void _navigateToCreateDiscussion(BuildContext context) {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => const CreateDiscussionPage(),
      ),
    );
  }

  void _navigateToDiscussionDetail(BuildContext context, Discussion discussion) {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => DiscussionDetailPage(discussion: discussion),
      ),
    );
  }
}
