import 'package:ethioguide/core/components/button.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure.dart';
import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';

class FeedbackTab extends StatelessWidget {
  final List<FeedbackItem> feedbacklist;
  const FeedbackTab({super.key, required this.feedbacklist});

  String timeAgo(String isoDate) {
  try {
    final dateTime = DateTime.parse(isoDate).toLocal();
    final now = DateTime.now();
    final diff = now.difference(dateTime);

    if (diff.inSeconds < 60) {
      return "${diff.inSeconds}s ago";
    } else if (diff.inMinutes < 60) {
      return "${diff.inMinutes}m ago";
    } else if (diff.inHours < 24) {
      return "${diff.inHours}h ago";
    } else if (diff.inDays < 7) {
      return "${diff.inDays}d ago";
    } else {
      return "${(diff.inDays / 7).floor()}w ago";
    }
  } catch (e) {
    return isoDate; // fallback to raw date string
  }
}


  @override
  Widget build(BuildContext context) {
    return ListView(
      children: [
        ListView.builder(
          physics: const NeverScrollableScrollPhysics(),
          shrinkWrap: true,
          padding: const EdgeInsets.all(12),
          itemCount: feedbacklist.length > 2 ? 2 : feedbacklist.length,
          itemBuilder: (context, index) {
            final feedback = feedbacklist[index];
            return Column(
              children: [
                Card(
                  margin: const EdgeInsets.only(bottom: 12),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Padding(
                    padding: const EdgeInsets.all(12),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        // Name + Verified badge
                        Row(
                          children: [
                            Text(
                              'User',
                              style: const TextStyle(
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                            const SizedBox(width: 6),

                            Container(
                              padding: const EdgeInsets.symmetric(
                                horizontal: 6,
                                vertical: 2,
                              ),
                              decoration: BoxDecoration(
                                color: Colors.blue[50],
                                borderRadius: BorderRadius.circular(8),
                              ),
                              child: const Text(
                                "Verified",
                                style: TextStyle(
                                  fontSize: 12,
                                  color: Colors.blue,
                                ),
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 8),
                        // Feedback message
                        Text('${feedback.comment}'),
                        const SizedBox(height: 8),
                        // Time ago
                       
                        Text(
                          '${timeAgo(feedback.date) }',
                          style: const TextStyle(
                            color: Colors.grey,
                            fontSize: 12,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ],
            );
          },
        ),
        SizedBox(height: 8),
        CustomButton(
          text: 'Give Feedback',
          icon: Icons.feedback,
          onTap: () {
            // Handle download action
          },
        ),
      ],
    );
  }
}
