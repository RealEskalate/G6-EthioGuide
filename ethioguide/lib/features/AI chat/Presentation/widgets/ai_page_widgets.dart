import 'package:ethioguide/features/AI%20chat/Domain/entities/conversation.dart';
import 'package:ethioguide/features/AI%20chat/Presentation/bloc/ai_bloc.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

//############################################################################################
//#                                                                                          #
//#                                     App Bar                                              #
//#                                                                                          #
//############################################################################################

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

//############################################################################################
//#                                                                                          #
//#                                     Single Conversation Card                             #
//#                                                                                          #
//############################################################################################

Widget buildMessage({
  required Conversation conv,
  required BuildContext context,
}) {
  final hasRequest = conv.request.isNotEmpty;
  final hasResponse = conv.response.isNotEmpty || conv.source == 'loading';
  final isError = conv.source == 'error';
  final isLoading = conv.source == 'loading';

  return Column(
    children: [
      // User query (right-aligned)
      if (hasRequest)
        Align(
          alignment: Alignment.centerRight,
          child: Container(
            margin: const EdgeInsets.symmetric(vertical: 8),
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: Colors.teal,
              borderRadius: BorderRadius.circular(12),
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                Text(
                  'You: ${conv.request}',
                  style: const TextStyle(color: Colors.white),
                ),
                if (isLoading) ...[
                  const SizedBox(width: 8),
                  const SizedBox(
                    width: 16,
                    height: 16,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  ),
                ],
              ],
            ),
          ),
        ),
      // AI response, error, or initial greeting (left-aligned)
      if (hasResponse)
        Align(
          alignment: Alignment.centerLeft,
          child: Container(
            margin: const EdgeInsets.symmetric(vertical: 8),
            padding: isError
                ? const EdgeInsets.all(6)
                : const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: isError ? Colors.red[100] : Colors.grey[200],
              borderRadius: BorderRadius.circular(12),
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                if (!isError && !isLoading)
                  Row(
                    children: const [
                      Icon(Icons.verified, color: Colors.green, size: 16),
                      SizedBox(width: 4),
                      Text(
                        'Verified',
                        style: TextStyle(fontSize: 12, color: Colors.green),
                      ),
                    ],
                  ),
                _buildStepCard(
                  icon: isError ? Icons.error : Icons.assistant,
                  title: isError ? 'Error' : 'AI Response',
                  content: conv.response,
                  isError: isError,
                ),
                if (!isError && !isLoading && conv.id != 'initial')
                  _buildChecklistButton(context: context),
                if (conv.procedures.isNotEmpty && !isError && !isLoading)
                  ...conv.procedures.map(
                    (procedure) =>
                        _buildInfoCard(procedure: procedure!, context: context),
                  ),
              ],
            ),
          ),
        ),
    ],
  );
}

//############################################################################################
//#                                                                                          #
//#                                     AI Resonse Content                                   #
//#                                                                                          #
//############################################################################################

Widget _buildStepCard({
  required IconData icon,
  required String title,
  required String content,
  required bool isError,
}) {
  return Card(
    color: Colors.teal[50],
    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
    margin: isError
        ? const EdgeInsets.symmetric(vertical: 3)
        : const EdgeInsets.symmetric(vertical: 8),
    child: Padding(
      padding: const EdgeInsets.all(10),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Icon(icon, color: Colors.teal),
              const SizedBox(width: 8),
              Text(
                title,
                style: const TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                  color: Colors.teal,
                ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Text(content, style: const TextStyle(fontSize: 14)),
        ],
      ),
    ),
  );
}

//############################################################################################
//#                                                                                          #
//#                                     related Procedure cards                              #
//#                                                                                          #
//############################################################################################

Widget _buildInfoCard({
  required Procedure procedure,
  required BuildContext context,
}) {
  return Card(
    elevation: 2,
    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
    margin: const EdgeInsets.symmetric(vertical: 6, horizontal: 4),
    color: const Color(0xFFF1FAF9), // pale teal background
    shadowColor: Colors.black.withOpacity(0.05),
    child: Padding(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          const Icon(
            Icons.info,
            color: Colors.teal,
            size: 22,
          ),
          const SizedBox(width: 10),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  procedure.name,
                  style: const TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w600,
                    color: Colors.black87,
                  ),
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 4),
                Row(
                  children: [
                    _buildCompactButton(
                      context: context,
                      label: 'View',
                      onPressed: () {
                        ScaffoldMessenger.of(context).showSnackBar(
                          SnackBar(
                            content: Text('Viewing procedure: ${procedure.name}'),
                            duration: const Duration(seconds: 2),
                          ),
                        );
                      },
                    ),
                    const SizedBox(width: 6),
                    _buildCompactButton(
                      context: context,
                      label: 'Start',
                      onPressed: () {
                        ScaffoldMessenger.of(context).showSnackBar(
                          SnackBar(
                            content: Text('Starting procedure: ${procedure.name}'),
                            duration: const Duration(seconds: 2),
                          ),
                        );
                      },
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    ),
  );
}

Widget _buildCompactButton({
  required BuildContext context,
  required String label,
  required VoidCallback onPressed,
}) {
  return ElevatedButton(
    onPressed: onPressed,
    style: ElevatedButton.styleFrom(
      backgroundColor: const Color.fromARGB(255, 18, 159, 145),
      foregroundColor: Colors.white,
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      minimumSize: const Size(64, 30),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
      elevation: 0,
      textStyle: const TextStyle(fontSize: 12, fontWeight: FontWeight.w500),
    ),
    child: Text(label),
  );
}


//############################################################################################
//#                                                                                          #
//#                                     CheckList button                                     #
//#                                                                                          #
//############################################################################################

Widget _buildChecklistButton({required BuildContext context}) {
  return Card(
    color: Colors.teal[100],
    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
    margin: const EdgeInsets.symmetric(vertical: 8),
    child: ExpansionTile(
      leading: const Icon(Icons.checklist, color: Colors.teal),
      title: const Text(
        'Save Checklist',
        style: TextStyle(color: Colors.teal, fontWeight: FontWeight.bold),
      ),
      children: [
        ListTile(
          leading: const Icon(Icons.play_arrow),
          title: const Text('Start Procedure'),
          onTap: () {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(content: Text('Starting procedure...')),
            );
            // TODO: Navigate to procedure start page
          },
        ),
        ListTile(
          leading: const Icon(Icons.translate),
          title: const Text('Translate'),
          onTap: () {
            SnackBar(content: Text('Translating response'));
          },
        ),
      ],
    ),
  );
}


//############################################################################################
//#                                                                                          #
//#                                     For FOQ's                                            #
//#                                                                                          #
//############################################################################################
Widget questionCard(String question) {
  return Container(
    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
    decoration: BoxDecoration(
      gradient: const LinearGradient(
        colors: [
          Color(0xFF00695C), // Teal
          Color.fromARGB(255, 23, 89, 124), // Blue
        ],
        begin: Alignment.topLeft,
        end: Alignment.bottomRight,
      ),
      borderRadius: BorderRadius.circular(16),
      boxShadow: [
        BoxShadow(
          color: Colors.black.withOpacity(0.2),
          blurRadius: 8,
          offset: const Offset(0, 4),
        ),
      ],
    ),
    child: Text(
      question,
      style: const TextStyle(
        color: Colors.white,
        fontSize: 10,
        fontWeight: FontWeight.w600,
        letterSpacing: 0.5,
      ),
      textAlign: TextAlign.center,
      maxLines: 2,
      overflow: TextOverflow.ellipsis,
    ),
  );
}
