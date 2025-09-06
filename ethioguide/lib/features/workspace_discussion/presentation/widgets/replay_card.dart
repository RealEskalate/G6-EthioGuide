import 'package:flutter/material.dart';

class ReplayCard extends StatelessWidget {
  const ReplayCard({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Column(
        children: [
          // 🔹 Input field for reply
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: Row(
              children: [
                Expanded(
                  child: TextField(
                    decoration: InputDecoration(
                      hintText: "Add a helpful reply...",
                      contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
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
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                  ),
                ),
              ],
            ),
          ),

          // 🔹 List of replies
          Expanded(
            child: ListView(
              children: const [
                ReplyCard(
                  initials: "S",
                  name: "Sarah Bekele",
                  timeAgo: "1 hour ago",
                  comment:
                      "You'll need your old passport, birth certificate, and two passport photos. The process usually takes 2-3 weeks.",
                  likes: 8,
                ),
                ReplyCard(
                  initials: "D",
                  name: "Daniel Mekonnen",
                  timeAgo: "45 minutes ago",
                  comment:
                      "Don’t forget to bring copies of everything! They always ask for copies.",
                  likes: 3,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class ReplyCard extends StatelessWidget {
  final String initials;
  final String name;
  final String timeAgo;
  final String comment;
  final int likes;

  const ReplyCard({
    super.key,
    required this.initials,
    required this.name,
    required this.timeAgo,
    required this.comment,
    required this.likes,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 🔹 Avatar with initials
          CircleAvatar(
            child: Text(initials),
          ),
          const SizedBox(width: 12),

          // 🔹 Comment box
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Name & time
                Row(
                  children: [
                    Text(
                      name,
                      style: const TextStyle(fontWeight: FontWeight.bold),
                    ),
                    const SizedBox(width: 8),
                    Text(
                      timeAgo,
                      style: TextStyle(color: Colors.grey[600], fontSize: 12),
                    ),
                  ],
                ),
                const SizedBox(height: 4),

                // Comment text
                Text(comment),

                const SizedBox(height: 8),

                // 🔹 Actions (Like + Reply)
                Row(
                  children: [
                    Icon(Icons.thumb_up_alt_outlined, size: 16, color: Colors.grey[600]),
                    const SizedBox(width: 4),
                    Text(likes.toString(), style: TextStyle(color: Colors.grey[700], fontSize: 12)),
                    const SizedBox(width: 16),
                    Icon(Icons.reply, size: 16, color: Colors.grey[600]),
                    const SizedBox(width: 4),
                    Text("Reply", style: TextStyle(color: Colors.grey[700], fontSize: 12)),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
