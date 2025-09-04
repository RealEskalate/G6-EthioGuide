import 'package:ethioguide/features/AI%20chat/Presentation/bloc/ai_bloc.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

AppBar appBar({required BuildContext context}) {
  return AppBar(
    leading: Icon(Icons.menu),
    title: const Text('EthioGuide AI Assistant'),
    actions: [
      IconButton(
        icon: const Icon(Icons.history),
        tooltip: 'View History',
        onPressed: () {
          context.read<AiBloc>().add(GetHistoryEvent());
        },
      ),
      PopupMenuButton<String>(
        icon: const Icon(Icons.more_vert),
        tooltip: 'More Options',
        onSelected: (value) {
          if (value == 'clear_history') {
            // TODO: Implement ClearHistoryEvent
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(content: Text('Clearing history...')),
            );
          } else if (value == 'change_language') {
            // TODO: Show language selection dialog
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(content: Text('Changing language...')),
            );
          } else if (value == 'settings') {
            // TODO: Navigate to settings page
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(content: Text('Opening settings...')),
            );
          } else if (value == 'log_out') {
            // TODO: Trigger logout with AuthRepository
            ScaffoldMessenger.of(
              context,
            ).showSnackBar(const SnackBar(content: Text('Logging out...')));
          }
        },
        itemBuilder: (context) => [
          const PopupMenuItem(
            value: 'clear_history',
            child: Text('Clear History'),
          ),
          const PopupMenuItem(
            value: 'change_language',
            child: Text('Change Language'),
          ),
          const PopupMenuItem(value: 'settings', child: Text('Settings')),
          const PopupMenuItem(value: 'log_out', child: Text('Log Out')),
        ],
      ),
    ],
  );
}
