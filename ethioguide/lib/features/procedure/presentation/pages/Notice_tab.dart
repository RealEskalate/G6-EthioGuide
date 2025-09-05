import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';

class NoticesTab extends StatelessWidget {
  const NoticesTab({super.key});

@override
  Widget build(BuildContext context) {
    final feedbackList = [
      {
        "name": "Abebe M.",
        "verified": true,
        "message":
            "The process was smooth and staff were helpful. Completed in 2 hours as expected.",
        "time": "2 days ago",
      },
      {
        "name": "Sarah T.",
        "verified": true,
        "message":
            "Make sure to bring all original documents. They were strict about photocopies.",
        "time": "5 days ago",
      },
    ];

    return ListView.builder(
      padding: const EdgeInsets.all(12),
      itemCount: feedbackList.length,
      itemBuilder: (context, index) {
        final feedback = feedbackList[index];
        return Card(
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
                      '${feedback["name"]}',
                      style: const TextStyle(fontWeight: FontWeight.bold),
                    ),
                    const SizedBox(width: 6),
                    if (feedback["verified"] == true)
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
                          style: TextStyle(fontSize: 12, color: Colors.blue),
                        ),
                      ),
                  ],
                ),
                const SizedBox(height: 8),
                // Feedback message
                Text('${feedback["message"]}'),
                const SizedBox(height: 8),
                // Time
                Text(
                  '${feedback["time"]}',
                  style: const TextStyle(color: Colors.grey, fontSize: 12),
                ),
              ],
            ),
          ),
        );
      },
    );
  }
}