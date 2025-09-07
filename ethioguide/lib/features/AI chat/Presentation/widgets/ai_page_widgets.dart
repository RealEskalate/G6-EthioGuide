import 'package:ethioguide/features/AI%20chat/Domain/entities/conversation.dart';
import 'package:ethioguide/features/AI%20chat/Presentation/bloc/ai_bloc.dart';
import 'package:ethioguide/features/AI%20chat/data/models/translated_conversation_model.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

//?############################################################################################
//?#                                                                                          #
//?#                                     App Bar                                              #
//?#                                                                                          #
//?############################################################################################

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

//?############################################################################################
//?#                                                                                          #
//?#                               Single Conversation Card                                   #
//?#                                                                                          #
//?############################################################################################

Widget buildMessage({
  required Conversation conv,
  required BuildContext context,
}) {
  final hasRequest = conv.request.isNotEmpty;
  final hasResponse = conv.response.isNotEmpty;
  final isLoading = conv.source == 'loading';
  final hasFailed = conv.source == 'failed';

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
                if (hasFailed) ...[
                  const SizedBox(width: 8),
                  const SizedBox(
                    child: Icon(
                      Icons.error,
                      color: Color.fromARGB(255, 245, 102, 91),
                    ),
                  ),
                ],
              ],
            ),
          ),
        ),
      // AI response, or initial greeting (left-aligned)
      if (hasResponse)
        Align(
          alignment: Alignment.centerLeft,
          child: Container(
            margin: const EdgeInsets.symmetric(vertical: 8),
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: Colors.grey[200],
              borderRadius: BorderRadius.circular(12),
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                if (!isLoading)
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
                  icon: Icons.assistant,
                  title: 'AI Response',
                  content: conv.response,
                ),
                if (!isLoading && conv.id != 'initial')
                  _buildChecklistButton(context: context, conv: conv),
                if (conv.procedures.isNotEmpty && !isLoading)
                  _buildRelatedProcedures(
                    procedures: conv.procedures,
                    context: context,
                  ),
              ],
            ),
          ),
        ),
    ],
  );
}

//?############################################################################################
//?#                                                                                          #
//?#                                     Error response                                       #
//?#                                                                                          #
//?############################################################################################

Widget errorCard(String errorMessage) {
  return Align(
    alignment: Alignment.centerLeft,
    child: Card(
      color: Colors.red[50],
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
        side: BorderSide(color: Colors.red.shade300, width: 1),
      ),
      margin: const EdgeInsets.symmetric(vertical: 6, horizontal: 4),
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Icon(Icons.error_outline, color: Colors.red, size: 22),
            const SizedBox(width: 10),
            Expanded(
              child: Text(
                errorMessage,
                style: TextStyle(
                  fontSize: 13,
                  fontWeight: FontWeight.w500,
                  color: Colors.red.shade900,
                ),
              ),
            ),
          ],
        ),
      ),
    ),
  );
}

//?############################################################################################
//?#                                                                                          #
//?#                                     AI Resonse Content                                   #
//?#                                                                                          #
//?############################################################################################

Widget _buildStepCard({
  required IconData icon,
  required String title,
  required String content,
}) {
  return Card(
    color: Colors.teal[50],
    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
    margin: const EdgeInsets.symmetric(vertical: 8),
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

//?############################################################################################
//?#                                                                                          #
//?#                                related Procedure cards                                   #
//?#                                                                                          #
//?############################################################################################

Widget _buildRelatedProcedures({
  required List<Procedure?> procedures,
  required BuildContext context,
}) {
  if (procedures.isEmpty) return const SizedBox.shrink();

  return Card(
    color: Colors.teal[100],
    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
    margin: const EdgeInsets.symmetric(vertical: 8),
    child: ExpansionTile(
      leading: const Icon(Icons.folder_copy, color: Colors.teal),
      title: const Text(
        'Related Procedures',
        style: TextStyle(color: Colors.teal, fontWeight: FontWeight.bold),
      ),
      children: procedures.map((procedure) {
        return ListTile(
          leading: const Icon(Icons.info, color: Colors.teal),
          title: Text(
            procedure!.name,
            style: const TextStyle(fontSize: 14, fontWeight: FontWeight.w500),
            overflow: TextOverflow.ellipsis,
          ),
          trailing: Wrap(
            spacing: 6,
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
        );
      }).toList(),
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

//?############################################################################################
//?#                                                                                          #
//?#                                   CheckList button                                       #
//?#                                                                                          #
//?############################################################################################

Widget _buildChecklistButton({
  required BuildContext context,
  required Conversation conv,
}) {
  TranslatedConversationModel conversation = TranslatedConversationModel(
    response: conv.response,
    procedures: conv.procedures,
  );
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
            context.read<AiBloc>().add(
              TranslateContentEvent(conversation: conversation, id: conv.id),
            );
          },
        ),
      ],
    ),
  );
}

//?############################################################################################
//?#                                                                                          #
//?#                                     For FOQ's                                            #
//?#                                                                                          #
//?############################################################################################
Widget questionCard(String question) {
  return Container(
    padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
    decoration: BoxDecoration(
      color: const Color(0xFFF1FAF9), // soft pale teal background
      borderRadius: BorderRadius.circular(12),
      border: Border.all(
        color: Colors.teal.withOpacity(0.3), // subtle border
        width: 1,
      ),
      boxShadow: [
        BoxShadow(
          color: Colors.black.withOpacity(0.05),
          blurRadius: 6,
          offset: const Offset(0, 3),
        ),
      ],
    ),
    child: Text(
      question,
      style: const TextStyle(
        color: Colors.black87,
        fontSize: 12,
        fontWeight: FontWeight.w500,
        letterSpacing: 0.3,
      ),
      textAlign: TextAlign.center,
      maxLines: 2,
      overflow: TextOverflow.ellipsis,
    ),
  );
}

//?############################################################################################
//?#                                                                                          #
//?#                             Loading Indicator Widget                                     #
//?#                                                                                          #
//?############################################################################################

Widget loadingCard() {
  return Align(
    alignment: Alignment.centerLeft,
    child: Card(
      color: Colors.blueGrey.shade50,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
        side: BorderSide(color: Colors.blueGrey.shade100, width: 1),
      ),
      margin: const EdgeInsets.symmetric(vertical: 6, horizontal: 4),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            Expanded(
              child: Text(
                'AI is processing your request',
                style: const TextStyle(
                  fontSize: 13,
                  fontWeight: FontWeight.w500,
                  color: Colors.black87,
                ),
              ),
            ),
            const SizedBox(width: 8),
            const SizedBox(
              width: 40,
              child: _AnimatedDots(), // bigger dots animation
            ),
          ],
        ),
      ),
    ),
  );
}

//?############################################################################################
//?#                                                                                          #
//?#                                  Loading Animation                                       #
//?#                                                                                          #
//?############################################################################################
/// Small animated dots widget (bigger version)
class _AnimatedDots extends StatefulWidget {
  const _AnimatedDots({Key? key}) : super(key: key);

  @override
  State<_AnimatedDots> createState() => _AnimatedDotsState();
}

class _AnimatedDotsState extends State<_AnimatedDots>
    with SingleTickerProviderStateMixin {
  late AnimationController _controller;
  late Animation<double> _animation;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 1),
    )..repeat();
    _animation = Tween<double>(begin: 0, end: 3).animate(_controller);
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: _animation,
      builder: (context, child) {
        int dots = (_animation.value % 3).ceil();
        return Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: List.generate(
            3,
            (index) => Container(
              margin: const EdgeInsets.symmetric(horizontal: 2),
              width: 8, // bigger dot
              height: 8, // bigger dot
              decoration: BoxDecoration(
                color: index < dots ? Colors.teal : Colors.teal.shade200,
                shape: BoxShape.circle,
              ),
            ),
          ),
        );
      },
    );
  }
}
