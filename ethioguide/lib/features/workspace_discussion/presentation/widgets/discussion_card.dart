import 'package:flutter/material.dart';
import '../../domain/entities/discussion.dart';
import '../../domain/entities/user.dart';

/// Widget that displays a single discussion card
class DiscussionCard extends StatelessWidget {
  final Discussion discussions;
  final VoidCallback? onTap;
  final VoidCallback? onLike;
  final VoidCallback? onReport;
  final VoidCallback? onComment;

   DiscussionCard({
    super.key,
    required this.discussions,
    this.onTap,
    this.onLike,
    this.onReport,
    this.onComment,
  });


final discussion = 
   Discussion(
      id: 'd',
      title: 'How to complete step ${1 + 1}?',
      content: 'I need help with the requirements for step ${1 + 1}. Any tips?',
      tags: ['licensing', 'process'],
      category: 1 % 2 == 0 ? 'General' : 'Business',
      createdAt: DateTime.now().subtract(Duration(days: 1)),
      createdBy:  User(id: 'u1', name: 'Test User'),
      likeCount: 1 * 3,
      reportCount: 0,
      commentsCount: 1 + 1,
    );
  
  
  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 2,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
      ),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Header with user info and metadata
              Row(
                children: [
                  // User avatar
                  CircleAvatar(
                    backgroundColor: Colors.grey[300],
                    child: Text(
                      discussion.createdBy.name[0].toUpperCase(),
                      style: TextStyle(
                        color: Colors.grey[700],
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                  const SizedBox(width: 12),
                  // User info
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          children: [
                            Text(
                              discussion.createdBy.name,
                              style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                            if (discussion.createdBy.role != null) ...[
                              const SizedBox(width: 8),
                              Container(
                                padding: const EdgeInsets.symmetric(
                                  horizontal: 8,
                                  vertical: 2,
                                ),
                                decoration: BoxDecoration(
                                  color: Colors.green[100],
                                  borderRadius: BorderRadius.circular(12),
                                ),
                                child: Text(
                                  discussion.createdBy.role!,
                                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                                    color: Colors.green[700],
                                    fontWeight: FontWeight.w500,
                                  ),
                                ),
                              ),
                            ],
                          ],
                        ),
                        const SizedBox(height: 2),
                        Row(
                          children: [
                            Text(
                              _formatTimeAgo(discussions.createdAt),
                              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                                color: Colors.grey[600],
                              ),
                            ),
                            const SizedBox(width: 16),
                            Icon(
                              Icons.visibility,
                              size: 16,
                              color: Colors.grey[600],
                            ),
                            const SizedBox(width: 4),
                            Text(
                              '${discussion.likeCount + discussion.reportCount}',
                              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                                color: Colors.grey[600],
                              ),
                            ),
                          ],
                        ),
                      ],
                    ),
                  ),
                  // More options
                  IconButton(
                    onPressed: () {},
                    icon: const Icon(Icons.more_vert),
                    padding: EdgeInsets.zero,
                    constraints: const BoxConstraints(),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              
              // Title with pin indicator
              Row(
                children: [
                  if (discussion.isPinned) ...[
                    Icon(
                      Icons.push_pin,
                      color: Colors.red[600],
                      size: 20,
                    ),
                    const SizedBox(width: 8),
                  ],
                  Expanded(
                    child: Text(
                      discussions.title,
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 8),
              
              // Content snippet
              Text(
                discussions.content,
                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                  color: Colors.grey[700],
                ),
                maxLines: 3,
                overflow: TextOverflow.ellipsis,
              ),
              const SizedBox(height: 12),
              
              // Tags
              Wrap(
                spacing: 8,
                runSpacing: 4,
                children: discussions.tags.take(4).map((tag) {
                  return Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 8,
                      vertical: 4,
                    ),
                    decoration: BoxDecoration(
                      color: Colors.grey[100],
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(
                      '#$tag',
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        color: Colors.grey[700],
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                  );
                }).toList(),
              ),
              const SizedBox(height: 16),
              
              // Action buttons
              Row(
                children: [
                  // Like button
                  _ActionButton(
                    icon: Icons.thumb_up_outlined,
                    label: '${discussion.likeCount}',
                    onTap: onLike,
                    color: Colors.blue,
                  ),
                  const SizedBox(width: 16),
                  // Comment button
                  _ActionButton(
                    icon: Icons.chat_bubble_outline,
                    label: '${discussion.commentsCount}',
                    onTap: onComment,
                    color: Colors.green,
                  ),
                  const Spacer(),
                  // Report button
                  _ActionButton(
                    icon: Icons.flag_outlined,
                    label: 'Report',
                    onTap: onReport,
                    color: Colors.red,
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  String _formatTimeAgo(DateTime dateTime) {
    final now = DateTime.now();
    final difference = now.difference(dateTime);
    
    if (difference.inDays > 0) {
      return '${difference.inDays} day${difference.inDays > 1 ? 's' : ''} ago';
    } else if (difference.inHours > 0) {
      return '${difference.inHours} hour${difference.inHours > 1 ? 's' : ''} ago';
    } else if (difference.inMinutes > 0) {
      return '${difference.inMinutes} minute${difference.inMinutes > 1 ? 's' : ''} ago';
    } else {
      return 'Just now';
    }
  }
}

class _ActionButton extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback? onTap;
  final Color color;

  const _ActionButton({
    required this.icon,
    required this.label,
    this.onTap,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(8),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              icon,
              size: 18,
              color: color,
            ),
            const SizedBox(width: 4),
            Text(
              label,
              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                color: color,
                fontWeight: FontWeight.w500,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
