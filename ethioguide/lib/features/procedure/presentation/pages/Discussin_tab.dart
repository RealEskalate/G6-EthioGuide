


import 'package:ethioguide/core/config/app_color.dart';
import 'package:ethioguide/core/config/route_names.dart';
import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:go_router/go_router.dart';

class DiscussionTab extends StatelessWidget {
  const DiscussionTab({super.key});

  @override
  Widget build(BuildContext context) {
    return  Center(
      child: SizedBox(
        height: 200,
        child: Card(
          
          
          margin:  EdgeInsets.only(bottom: 12),
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(12),
            ),
            child: Padding(
              padding: const EdgeInsets.all(12),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  // Name + Verified badge
                  Row(
                    children: [
                      Icon(
                        Icons.chat_bubble_outline_rounded,
                        color: AppColors.darkGreenColor,
                        
                        ),
                      Text(
                        'Question about this process',
                        style: const TextStyle(fontWeight: FontWeight.bold),
                      ),
            
            
                    ],
                  ),
                  const SizedBox(height: 8),
                  // Feedback message
                  Text('Join the community discussion to ask questions and share experiences.'),
                  const SizedBox(height: 8),
                  // Time
                ElevatedButton(onPressed: () {
              context
                  .push(RouteNames.workspacediscussion);
                }, child: const Text('Join discussion')),
                ],
              ),
            ),
        ),
      )
    );
  }
}