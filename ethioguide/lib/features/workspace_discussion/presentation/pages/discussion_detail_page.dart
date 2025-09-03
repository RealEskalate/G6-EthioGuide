import 'package:ethioguide/core/config/app_color.dart';
import 'package:ethioguide/features/workspace_discussion/presentation/widgets/replay_card.dart';
import 'package:flutter/material.dart';
import '../../domain/entities/discussion.dart';
import '../../domain/entities/comment.dart';

import 'package:flutter/material.dart';
import '../../domain/entities/discussion.dart';
import '../../domain/entities/comment.dart';

class DiscussionDetailPage extends StatefulWidget {
  final Discussion discussion;

  const DiscussionDetailPage({super.key, required this.discussion});

  @override
  State<DiscussionDetailPage> createState() => _DiscussionDetailPageState();
}

class _DiscussionDetailPageState extends State<DiscussionDetailPage> {
  final _commentController = TextEditingController();
  late int _likeCount;
  late List<Comment> _comments;

  @override
  void initState() {
    super.initState();
    _likeCount = widget.discussion.likeCount;
    _comments = List.generate(
      1,
      (i) => Comment(
        id: 'c$i',
        discussionId: widget.discussion.id,
        content:
            "You'll need your old passport, birth certificate, and two passport photos. The process usually takes 2–3 weeks.",
        createdAt: DateTime.now().subtract(const Duration(hours: 1)),
        createdBy: widget.discussion.createdBy,
        likeCount: 8,
        reportCount: 0,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return 
      
         AlertDialog(
          insetPadding: const EdgeInsets.all(12.0),
          contentPadding: EdgeInsets.zero,
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
          content: SizedBox(
             width: MediaQuery.of(context).size.width * 0.85,
              height: MediaQuery.of(context).size.height * 0.75,
            child: Column(
              children: [
                Expanded(
                  child: SingleChildScrollView(
                    padding: const EdgeInsets.all(16),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        _buildDiscussionHeader(),
                        const SizedBox(height: 12),
                        Text(
                          widget.discussion.title,
                          style: Theme.of(context).textTheme.titleMedium?.copyWith(
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          widget.discussion.content,
                          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                            height: 1.4,
                            color: Colors.grey[800],
                          ),
                        ),
                        const SizedBox(height: 12),
                        _buildTags(),
                        const SizedBox(height: 16),
                        _buildActionButtons(),
                        const Divider(height: 32),
            
                        Padding(
                          padding: const EdgeInsets.all(8.0),
                          child: Column(
                            children: [
                              Row(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  CircleAvatar(
                radius: 16,
                backgroundColor: Colors.grey[300],
                child: Text(
                  'You',
                  style: const TextStyle(fontSize: 12),
                ),
              ),
              SizedBox(width: 8),
                                  Expanded(
                                    child: TextField(
                                      maxLines: 3,
                                      decoration: InputDecoration(
                                        hintText: "Add a helpful reply...",
                                        
                                        contentPadding: const EdgeInsets.symmetric(
                                          horizontal: 12,
                                          vertical: 10,
                                        ),
                                        filled: true,
                                  
                                        fillColor: AppColors.darkGreenColor.withOpacity(0.05),
                                        border: OutlineInputBorder(
                                          // no border
                                  
                                          borderSide: BorderSide.none,
                                          borderRadius: BorderRadius.circular(12),
                                        ),
                                      ),
                                    ),
                                  ),
                                ],
                              ),
                              const SizedBox(height: 10),
                              Row(
                                crossAxisAlignment: CrossAxisAlignment.end,
                                mainAxisAlignment: MainAxisAlignment.end,
                                children: [
                                  ElevatedButton.icon(
                                   
                                    onPressed: () {},
                                    icon: const Icon(Icons.send, size: 18),
                                    label: const Text("Reply"),
                                    style: ElevatedButton.styleFrom(
                                      
                                      shape: RoundedRectangleBorder(
                                        borderRadius: BorderRadius.circular(12),
                                      ),
                                    ),
                                  ),
                                ],
                              ),
                            ],
                          ),
                        ),
            
                        Text(
                          "Replies (${_comments.length})",
                          style: Theme.of(context).textTheme.titleMedium?.copyWith(
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const SizedBox(height: 12),
                        ..._comments.map(_buildCommentItem).toList(),
                      ],
                    ),
                  ),
                ),
                // _buildCommentInput(),
              ],
            ),
          ),
        );
    
  }

  Widget _buildDiscussionHeader() {
    return Row(
      children: [
        CircleAvatar(
          radius: 20,
          backgroundColor: Colors.grey[300],
          child: Text(
            widget.discussion.createdBy.name[0].toUpperCase(),
            style: const TextStyle(fontWeight: FontWeight.bold),
          ),
        ),
        const SizedBox(width: 10),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                widget.discussion.createdBy.name,
                style: const TextStyle(fontWeight: FontWeight.w600),
              ),
              Row(
                children: [
                  Text(
                    _formatTimeAgo(widget.discussion.createdAt),
                    style: TextStyle(color: Colors.grey[600], fontSize: 12),
                  ),
                  const SizedBox(width: 8),
                  const Icon(
                    Icons.remove_red_eye,
                    size: 14,
                    color: Colors.grey,
                  ),
                  const SizedBox(width: 4),
                  Text(
                    "156",
                    style: TextStyle(color: Colors.grey[600], fontSize: 12),
                  ),
                ],
              ),
            ],
          ),
        ),
        IconButton(onPressed: () {}, icon: const Icon(Icons.close)),
      ],
    );
  }

  Widget _buildTags() {
    return Wrap(
      spacing: 8,
      runSpacing: 4,
      children: widget.discussion.tags.map((tag) {
        return Container(
          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
          decoration: BoxDecoration(
            color: Colors.grey[100],
            borderRadius: BorderRadius.circular(16),
            border: Border.all(color: Colors.grey[300]!),
          ),
          child: Text(
            '#$tag',
            style: Theme.of(context).textTheme.bodyMedium?.copyWith(
              color: Colors.grey[700],
              fontWeight: FontWeight.w500,
            ),
          ),
        );
      }).toList(),
    );
  }

  Widget _buildActionButtons() {
    return Row(
      children: [
        _ActionButton(
          icon: Icons.thumb_up_outlined,
          label: '$_likeCount',
          onTap: () {
            setState(() {
              _likeCount += 1;
            });
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(content: Text('You liked this (local UI only)')),
            );
          },
          color: Colors.blue,
        ),
        const SizedBox(width: 16),
        _ActionButton(
          icon: Icons.chat_bubble_outline,
          label: '${_comments.length}',
          onTap: () {},
          color: Colors.green,
        ),
        const Spacer(),
        _ActionButton(
          icon: Icons.flag_outlined,
          label: 'Report',
          onTap: () {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(content: Text('Reported (local UI only)')),
            );
          },
          color: Colors.red,
        ),
      ],
    );
  }

  /* Widget _buildTags() {
    return Wrap(
      spacing: 6,
      children: widget.discussion.tags
          .map(
            (tag) => Chip(
              label: Text("#$tag"),
              backgroundColor: Colors.grey[200],
              labelStyle: const TextStyle(fontSize: 12),
              padding: const EdgeInsets.symmetric(horizontal: 6),
            ),
          )
          .toList(),
    );
  }

  Widget _buildStats() {
    return Row(
      children: [
        Icon(Icons.thumb_up, size: 18, color: Colors.grey[700]),
        const SizedBox(width: 4),
        Text("$_likeCount"),
        const SizedBox(width: 16),
        Icon(Icons.chat_bubble_outline, size: 18, color: Colors.grey[700]),
        const SizedBox(width: 4),
        Text("${_comments.length} replies"),
        const Spacer(),
        IconButton(
          icon: const Icon(Icons.bookmark_border),
          onPressed: () {},
        ),
        IconButton(
          icon: const Icon(Icons.share_outlined),
          onPressed: () {},
        ),
      ],
    );
  }
 */
  Widget _buildCommentItem(Comment comment) {
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          CircleAvatar(
            radius: 16,
            backgroundColor: Colors.grey[300],
            child: Text(
              comment.createdBy.name[0].toUpperCase(),
              style: const TextStyle(fontSize: 12),
            ),
          ),
          const SizedBox(width: 8),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  comment.createdBy.name,
                  style: const TextStyle(fontWeight: FontWeight.w600),
                ),
                Text(
                  _formatTimeAgo(comment.createdAt),
                  style: TextStyle(color: Colors.grey[600], fontSize: 12),
                ),
                const SizedBox(height: 4),
                Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 12,
                    vertical: 10,
                  ),
                  decoration: BoxDecoration(
                    color: Colors.grey[100],
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(comment.content),
                ),
                const SizedBox(height: 6),
                Row(
                  children: [
                    Icon(
                      Icons.thumb_up_alt_outlined,
                      size: 16,
                      color: Colors.grey[600],
                    ),
                    const SizedBox(width: 4),
                    Text(
                      "${comment.likeCount}",
                      style: TextStyle(color: Colors.grey[600], fontSize: 12),
                    ),
                    const SizedBox(width: 16),
                    Text(
                      "Reply",
                      style: TextStyle(color: Colors.grey[600], fontSize: 12),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildCommentInput() {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        border: Border(top: BorderSide(color: Colors.grey[300]!)),
        color: Colors.white,
      ),
      child: Row(
        children: [
          Expanded(
            child: TextField(
              controller: _commentController,
              decoration: InputDecoration(
                hintText: "Add a helpful reply...",
                contentPadding: const EdgeInsets.symmetric(
                  horizontal: 16,
                  vertical: 10,
                ),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(25),
                  borderSide: BorderSide(color: Colors.grey[300]!),
                ),
              ),
            ),
          ),
          const SizedBox(width: 8),
          ElevatedButton.icon(
            onPressed: () {},
            icon: const Icon(Icons.send, size: 18),
            label: const Text("Reply"),
            style: ElevatedButton.styleFrom(
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(25),
              ),
              backgroundColor: Colors.blue,
            ),
          ),
        ],
      ),
    );
  }

  String _formatTimeAgo(DateTime dateTime) {
    final diff = DateTime.now().difference(dateTime);
    if (diff.inDays > 0) return "${diff.inDays} day(s) ago";
    if (diff.inHours > 0) return "${diff.inHours} hour(s) ago";
    if (diff.inMinutes > 0) return "${diff.inMinutes} min ago";
    return "Just now";
  }
}

/// Page for viewing discussion details and comments
// class DiscussionDetailPage extends StatefulWidget {
//   final Discussion discussion;

//   const DiscussionDetailPage({
//     super.key,
//     required this.discussion,
//   });

//   @override
//   State<DiscussionDetailPage> createState() => _DiscussionDetailPageState();
// }

// class _DiscussionDetailPageState extends State<DiscussionDetailPage> {
//   final _commentController = TextEditingController();
//   late int _likeCount;
//   late List<Comment> _comments;

//   @override
//   void initState() {
//     super.initState();
//     // Local dummy data for UI preview
//     _likeCount = widget.discussion.likeCount;
//     _comments = List.generate(
//       2,
//       (i) => Comment(
//         id: 'c$i',
//         discussionId: widget.discussion.id,
//         content: 'This is a sample reply #${i + 1}. Helpful information here.',
//         createdAt: DateTime.now().subtract(Duration(hours: i * 5 + 1)),
//         createdBy: widget.discussion.createdBy,
//         likeCount: i,
//         reportCount: 1, // small dummy
//       ),
//     );
//   }

//   @override
//   void dispose() {
//     _commentController.dispose();
//     super.dispose();
//   }

//   @override
//   Widget build(BuildContext context) {
//     return Scaffold(
//       appBar: AppBar(
//         title: const Text('Discussion'),
//         backgroundColor: Colors.white,
//         elevation: 0,
//       ),
//       body: Column(
//         children: [
//           // Discussion content
//           Expanded(
//             child: SingleChildScrollView(
//               padding: const EdgeInsets.all(16.0),
//               child: Column(
//                 crossAxisAlignment: CrossAxisAlignment.start,
//                 children: [
//                   // Discussion header
//                   _buildDiscussionHeader(),
//                   const SizedBox(height: 20),

//                   // Discussion content
//                   _buildDiscussionContent(),
//                   const SizedBox(height: 20),

//                   // Tags
//                   _buildTags(),
//                   const SizedBox(height: 20),

//                   // Action buttons
//                   _buildActionButtons(),
//                   const SizedBox(height: 24),

//                   // Comments section (local dummy)
//                   _buildCommentsSection(),

//                 ],
//               ),
//             ),
//           ),

//           // Comment input
//           _buildCommentInput(),
//         ],
//       ),
//     );
//   }

//   Widget _buildDiscussionHeader() {
//     return Row(
//       children: [
//         CircleAvatar(
//           backgroundColor: Colors.grey[300],
//           child: Text(
//             widget.discussion.createdBy.name[0].toUpperCase(),
//             style: TextStyle(
//               color: Colors.grey[700],
//               fontWeight: FontWeight.w600,
//             ),
//           ),
//         ),
//         const SizedBox(width: 12),
//         Expanded(
//           child: Column(
//             crossAxisAlignment: CrossAxisAlignment.start,
//             children: [
//               Row(
//                 children: [
//                   Text(
//                     widget.discussion.createdBy.name,
//                     style: Theme.of(context).textTheme.titleMedium?.copyWith(
//                       fontWeight: FontWeight.w600,
//                     ),
//                   ),
//                   if (widget.discussion.createdBy.role != null) ...[
//                     const SizedBox(width: 8),
//                     Container(
//                       padding: const EdgeInsets.symmetric(
//                         horizontal: 8,
//                         vertical: 2,
//                       ),
//                       decoration: BoxDecoration(
//                         color: Colors.green[100],
//                         borderRadius: BorderRadius.circular(12),
//                       ),
//                       child: Text(
//                         widget.discussion.createdBy.role!,
//                         style: Theme.of(context).textTheme.bodySmall?.copyWith(
//                           color: Colors.green[700],
//                           fontWeight: FontWeight.w500,
//                         ),
//                       ),
//                     ),
//                   ],
//                 ],
//               ),
//               const SizedBox(height: 4),
//               Text(
//                 _formatTimeAgo(widget.discussion.createdAt),
//                 style: Theme.of(context).textTheme.bodySmall?.copyWith(
//                   color: Colors.grey[600],
//                 ),
//               ),
//             ],
//           ),
//         ),
//         IconButton(
//           onPressed: () {},
//           icon: const Icon(Icons.more_vert),
//         ),
//       ],
//     );
//   }

//   Widget _buildDiscussionContent() {
//     return Column(
//       crossAxisAlignment: CrossAxisAlignment.start,
//       children: [
//         if (widget.discussion.isPinned) ...[
//           Row(
//             children: [
//               Icon(
//                 Icons.push_pin,
//                 color: Colors.red[600],
//                 size: 20,
//               ),
//               const SizedBox(width: 8),
//               Text(
//                 'Pinned by moderator',
//                 style: Theme.of(context).textTheme.bodySmall?.copyWith(
//                   color: Colors.red[600],
//                   fontWeight: FontWeight.w500,
//                 ),
//               ),
//             ],
//           ),
//           const SizedBox(height: 12),
//         ],
//         Text(
//           widget.discussion.title,
//           style: Theme.of(context).textTheme.headlineSmall?.copyWith(
//             fontWeight: FontWeight.bold,
//           ),
//         ),
//         const SizedBox(height: 16),
//         Text(
//           widget.discussion.content,
//           style: Theme.of(context).textTheme.bodyLarge?.copyWith(
//             color: Colors.grey[700],
//             height: 1.5,
//           ),
//         ),
//       ],
//     );
//   }

//   Widget _buildCommentsSection() {
//     return Column(
//       crossAxisAlignment: CrossAxisAlignment.start,
//       children: [
//         Text(
//           'Replies (${_comments.length})',
//           style: Theme.of(context).textTheme.titleMedium?.copyWith(
//             fontWeight: FontWeight.w600,
//           ),
//         ),
//         const SizedBox(height: 16),
//   ..._comments.map((comment) => _buildCommentItem(comment)),
//       ],
//     );
//   }

//   Widget _buildCommentItem(Comment comment) {
//     return Container(
//       margin: const EdgeInsets.only(bottom: 16),
//       padding: const EdgeInsets.all(16),
//       decoration: BoxDecoration(
//         color: Colors.grey[50],
//         borderRadius: BorderRadius.circular(12),
//         border: Border.all(color: Colors.grey[200]!),
//       ),
//       child: Column(
//         crossAxisAlignment: CrossAxisAlignment.start,
//         children: [
//           Row(
//             children: [
//               CircleAvatar(
//                 backgroundColor: Colors.grey[300],
//                 radius: 16,
//                 child: Text(
//                   comment.createdBy.name[0].toUpperCase(),
//                   style: TextStyle(
//                     color: Colors.grey[700],
//                     fontWeight: FontWeight.w600,
//                     fontSize: 12,
//                   ),
//                 ),
//               ),
//               const SizedBox(width: 12),
//               Expanded(
//                 child: Column(
//                   crossAxisAlignment: CrossAxisAlignment.start,
//                   children: [
//                     Text(
//                       comment.createdBy.name,
//                       style: Theme.of(context).textTheme.bodyMedium?.copyWith(
//                         fontWeight: FontWeight.w600,
//                       ),
//                     ),
//                     Text(
//                       _formatTimeAgo(comment.createdAt),
//                       style: Theme.of(context).textTheme.bodySmall?.copyWith(
//                         color: Colors.grey[600],
//                       ),
//                     ),
//                   ],
//                 ),
//               ),
//               IconButton(
//                 onPressed: () {},
//                 icon: const Icon(Icons.more_vert, size: 20),
//                 padding: EdgeInsets.zero,
//                 constraints: const BoxConstraints(),
//               ),
//             ],
//           ),
//           const SizedBox(height: 12),
//           Text(
//             comment.content,
//             style: Theme.of(context).textTheme.bodyMedium?.copyWith(
//               color: Colors.grey[700],
//             ),
//           ),
//           const SizedBox(height: 12),
//           Row(
//             children: [
//               _ActionButton(
//                 icon: Icons.thumb_up_outlined,
//                 label: '${comment.likeCount}',
//                 onTap: () {
//                   // local increment for UI preview
//                   final idx = _comments.indexWhere((c) => c.id == comment.id);
//                   if (idx != -1) {
//                     final old = _comments[idx];
//                     setState(() {
//                       _comments[idx] = Comment(
//                         id: old.id,
//                         discussionId: old.discussionId,
//                         content: old.content,
//                         createdAt: old.createdAt,
//                         createdBy: old.createdBy,
//                         likeCount: old.likeCount + 1,
//                         reportCount: old.reportCount,
//                       );
//                     });
//                     ScaffoldMessenger.of(context).showSnackBar(
//                       const SnackBar(content: Text('Liked comment (local UI only)')),
//                     );
//                   }
//                 },
//                 color: Colors.blue,
//               ),
//               const SizedBox(width: 16),
//               _ActionButton(
//                 icon: Icons.flag_outlined,
//                 label: 'Report',
//                 onTap: () {
//                   ScaffoldMessenger.of(context).showSnackBar(
//                     const SnackBar(content: Text('Reported comment (local UI only)')),
//                   );
//                 },
//                 color: Colors.red,
//               ),
//             ],
//           ),
//         ],
//       ),
//     );
//   }

//   Widget _buildCommentInput() {
//     return Container(
//       padding: const EdgeInsets.all(16),
//       decoration: BoxDecoration(
//         color: Colors.white,
//         border: Border(
//           top: BorderSide(color: Colors.grey[300]!),
//         ),
//       ),
//       child: Row(
//         children: [
//           Expanded(
//             child: TextField(
//               controller: _commentController,
//               decoration: InputDecoration(
//                 hintText: 'Add a helpful reply...',
//                 border: OutlineInputBorder(
//                   borderRadius: BorderRadius.circular(25),
//                   borderSide: BorderSide(color: Colors.grey[300]!),
//                 ),
//                 contentPadding: const EdgeInsets.symmetric(
//                   horizontal: 16,
//                   vertical: 12,
//                 ),
//               ),
//               maxLines: null,
//             ),
//           ),
//           const SizedBox(width: 12),
//           ElevatedButton(
//             onPressed: _commentController.text.trim().isEmpty
//                 ? null
//                 : _submitComment,
//             style: ElevatedButton.styleFrom(
//               backgroundColor: Colors.blue[600],
//               foregroundColor: Colors.white,
//               shape: RoundedRectangleBorder(
//                 borderRadius: BorderRadius.circular(25),
//               ),
//               padding: const EdgeInsets.symmetric(
//                 horizontal: 20,
//                 vertical: 12,
//               ),
//             ),
//             child: const Text('Reply'),
//           ),
//         ],
//       ),
//     );
//   }

//   void _submitComment() {
//     final text = _commentController.text.trim();
//     if (text.isEmpty) return;
//     final newComment = Comment(
//       id: 'c${_comments.length + 1}',
//       discussionId: widget.discussion.id,
//       content: text,
//       createdAt: DateTime.now(),
//       createdBy: widget.discussion.createdBy,
//       likeCount: 0,
//       reportCount: 0,
//     );
//     setState(() {
//       _comments.insert(0, newComment);
//       _commentController.clear();
//     });
//     ScaffoldMessenger.of(context).showSnackBar(
//       const SnackBar(content: Text('Comment added (local UI only)')),
//     );
//   }

//   String _formatTimeAgo(DateTime dateTime) {
//     final now = DateTime.now();
//     final difference = now.difference(dateTime);

//     if (difference.inDays > 0) {
//       return '${difference.inDays} day${difference.inDays > 1 ? 's' : ''} ago';
//     } else if (difference.inHours > 0) {
//       return '${difference.inHours} hour${difference.inHours > 1 ? 's' : ''} ago';
//     } else if (difference.inMinutes > 0) {
//       return '${difference.inMinutes} minute${difference.inMinutes > 1 ? 's' : ''} ago';
//     } else {
//       return 'Just now';
//     }
//   }
// }

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
            Icon(icon, size: 18, color: color),
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
