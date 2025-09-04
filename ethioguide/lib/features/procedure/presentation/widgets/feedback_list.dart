import 'package:flutter/material.dart';
import '../../domain/entities/procedure.dart';

class FeedbackList extends StatelessWidget {
  final List<FeedbackItem> feedback;

  const FeedbackList({super.key, required this.feedback});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text('Feedback', style: Theme.of(context).textTheme.titleMedium),
        const SizedBox(height: 10),
        ...feedback.map((f) => Card(
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
              child: ListTile(
                title: Text(f.user),
                subtitle: Text(f.comment),
                trailing: f.verified ? const Chip(label: Text('Verified')) : null,
              ),
            )),
      ],
    );
  }
}


