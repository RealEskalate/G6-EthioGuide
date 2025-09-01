import 'package:ethioguide/features/AI%20chat/Domain/entities/conversation.dart';
import 'package:ethioguide/features/AI%20chat/Presentation/bloc/ai_bloc.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

class ChatPage extends StatefulWidget {
  const ChatPage({super.key});

  @override
  State<ChatPage> createState() => _ChatPageState();
}

class _ChatPageState extends State<ChatPage> {
  final TextEditingController _queryController = TextEditingController();
  final ScrollController _scrollController = ScrollController();
  final FocusNode _queryFocusNode = FocusNode();
  List<Conversation> _history = [];

  @override
  void initState() {
    super.initState();

    // Add initial AI greeting
    _history.add(
      const Conversation(
        id: 'initial',
        request: '',
        response:
            '''Hello! I'm your AI Assistant. I can help you navigate Ethiopian legal procedures, business registration, and more. What would you like to know?''',
        source: 'ai-generated',
        procedures: [],
      ),
    );

    // Fetch history on init
    context.read<AiBloc>().add(GetHistoryEvent());
    // Update border color on focus
    _queryFocusNode.addListener(() => setState(() {}));
  }

  void _addMessage({required Conversation conversation}) {
    setState(() {
      _history.add(conversation);
    });
    _scrollToBottom();
  }

  void _addLoadingMessage(String query) {
    setState(() {
      _history.add(
        Conversation(
          id: 'loading-${DateTime.now().millisecondsSinceEpoch}',
          request: query,
          response: '',
          source: 'loading',
          procedures: [],
        ),
      );
    });
    _scrollToBottom();
  }

  void _removeLoadingMessage() {
    setState(() {
      _history.removeWhere((conv) => conv.source == 'loading');
    });
  }

  void _scrollToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_scrollController.hasClients) {
        _scrollController.animateTo(
          _scrollController.position.maxScrollExtent,
          duration: const Duration(milliseconds: 300),
          curve: Curves.easeOut,
        );
      }
    });
  }

  void _showTranslateDialog() {
    final TextEditingController translateController = TextEditingController();
    String selectedLang = 'am'; // Default to Amharic
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Translate Text'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: translateController,
              decoration: const InputDecoration(
                hintText: 'Enter text to translate...',
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 16),
            DropdownButton<String>(
              value: selectedLang,
              isExpanded: true,
              items: const [
                DropdownMenuItem(value: 'am', child: Text('Amharic')),
                DropdownMenuItem(value: 'ti', child: Text('Tigrinya')),
                DropdownMenuItem(value: 'en', child: Text('English')),
              ],
              onChanged: (value) {
                setState(() {
                  selectedLang = value!;
                });
              },
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              final content = translateController.text.trim();
              if (content.isNotEmpty) {
                context.read<AiBloc>().add(
                  TranslateContentEvent(content: content, lang: selectedLang),
                );
                Navigator.pop(context);
              }
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: Colors.teal,
              foregroundColor: Colors.white,
            ),
            child: const Text('Translate'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return BlocListener<AiBloc, AiState>(
      listener: (context, state) {
        if (state is AiHistorySuccess) {
          setState(() {
            _history = [
              ...state.history,
              _history.first,
            ]; // Keep initial greeting
          });
          _scrollToBottom();
        } else if (state is AiQuerySuccess) {
          _removeLoadingMessage();
          _addMessage(conversation: state.conversation);
        } else if (state is AiError) {
          _removeLoadingMessage();
          _addMessage(
            conversation: Conversation(
              id: 'error-${DateTime.now().millisecondsSinceEpoch}',
              request: _queryController.text,
              response: state.message,
              source: 'error',
              procedures: [],
            ),
          );
        } else if (state is AiTranslateSuccess) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Translated: ${state.translated}')),
          );
        }
      },
      child: Scaffold(
        appBar: AppBar(
          title: const Text('AI Legal Assistant'),
          actions: [
            IconButton(
              icon: const Icon(Icons.history),
              onPressed: () {
                context.read<AiBloc>().add(GetHistoryEvent());
              },
            ),
            IconButton(
              icon: const Icon(Icons.more_vert),
              onPressed: () {
                // TODO: Add more options (e.g., clear history)
              },
            ),
          ],
        ),
        body: Column(
          children: [
            Expanded(
              child: BlocBuilder<AiBloc, AiState>(
                builder: (context, state) {
                  return ListView.builder(
                    controller: _scrollController,
                    padding: const EdgeInsets.all(16),
                    itemCount: _history.length,
                    itemBuilder: (context, index) {
                      final conv = _history[index];
                      return _buildMessage(conv);
                    },
                  );
                },
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(16),
              child: Row(
                children: [
                  Expanded(
                    child: TextField(
                      controller: _queryController,
                      focusNode: _queryFocusNode,
                      decoration: InputDecoration(
                        hintText: 'Type your question here...',
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(24),
                          borderSide: BorderSide(
                            color: _queryFocusNode.hasFocus
                                ? Colors.blue
                                : Colors.grey,
                          ),
                        ),
                        focusedBorder: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(24),
                          borderSide: const BorderSide(
                            color: Colors.blue,
                            width: 2,
                          ),
                        ),
                      ),
                      onSubmitted: (_) => _sendQuery(),
                    ),
                  ),
                  IconButton(
                    icon: const Icon(Icons.send, color: Colors.teal),
                    onPressed: _sendQuery,
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  //* custom widgtes

  Widget _buildMessage(Conversation conv) {
    final isUser = conv.request.isNotEmpty && conv.source != 'error';
    final isError = conv.source == 'error';
    final isLoading = conv.source == 'loading';

    return Align(
      alignment: isUser ? Alignment.centerRight : Alignment.centerLeft,
      child: Container(
        margin: const EdgeInsets.symmetric(vertical: 8),
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          color: isUser
              ? Colors.teal
              : isError
              ? Colors.red[100]
              : Colors.grey[200],
          borderRadius: BorderRadius.circular(12),
        ),
        child: Column(
          crossAxisAlignment: isUser
              ? CrossAxisAlignment.end
              : CrossAxisAlignment.start,
          children: [
            if (!isUser && !isError && !isLoading)
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
            if (isUser)
              Row(
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
            if (!isLoading && !isUser)
              _buildStepCard(
                icon: isError ? Icons.error : Icons.assistant,
                title: isError ? 'Error' : 'AI Response',
                content: conv.response,
              ),
            if (conv.procedures.isNotEmpty && !isUser && !isError && !isLoading)
              ...conv.procedures.map(
                (procedure) => _buildInfoCard(procedure: procedure!),
              ),
            if (!isUser && !isError && !isLoading) _buildChecklistButton(),
          ],
        ),
      ),
    );
  }

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
        padding: const EdgeInsets.all(16),
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
  Widget _buildInfoCard({required Procedure procedure}) {
    return Card(
      color: Colors.yellow[50],
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      margin: const EdgeInsets.symmetric(vertical: 4),
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Row(
          children: [
            const Icon(Icons.info, color: Colors.yellow),
            const SizedBox(width: 8),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    procedure.name,
                    style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                  ),
                  const SizedBox(height: 8),
                  Row(
                    children: [
                      ElevatedButton(
                        onPressed: () {
                          ScaffoldMessenger.of(context).showSnackBar(
                            SnackBar(content: Text('Viewing procedure: ${procedure.name}')),
                          );
                          // TODO: Navigate to procedure details page
                        },
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.teal,
                          foregroundColor: Colors.white,
                        ),
                        child: const Text('View'),
                      ),
                      const SizedBox(width: 8),
                      OutlinedButton(
                        onPressed: () {
                          ScaffoldMessenger.of(context).showSnackBar(
                            SnackBar(content: Text('Starting procedure: ${procedure.name}')),
                          );
                          // TODO: Navigate to procedure start page
                        },
                        style: OutlinedButton.styleFrom(
                          side: const BorderSide(color: Colors.teal),
                        ),
                        child: const Text('Start Procedure', style: TextStyle(color: Colors.teal)),
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

  void _sendQuery() {
    final query = _queryController.text.trim();
    if (query.isNotEmpty) {
      _addLoadingMessage(query);
      context.read<AiBloc>().add(SendQueryEvent(query: query));
      _queryController.clear();
      _queryFocusNode.unfocus();
    }
  }

  @override
  void dispose() {
    _queryController.dispose();
    _scrollController.dispose();
    _queryFocusNode.dispose();
    super.dispose();
  }

  Widget _buildChecklistButton() {
    return Card(
      color: Colors.teal[100],
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      margin: EdgeInsets.symmetric(vertical: 8),
      child: ExpansionTile(
        leading: Icon(Icons.checklist, color: Colors.teal),
        title: Text(
          'Save Checklist',
          style: TextStyle(color: Colors.teal, fontWeight: FontWeight.bold),
        ),
        children: [
          ListTile(
            leading: Icon(Icons.play_arrow),
            title: Text('Start Procedure'),
            onTap: () {
              // Handle action: e.g., navigate to procedure page
              ScaffoldMessenger.of(
                context,
              ).showSnackBar(SnackBar(content: Text('Starting procedure...')));
            },
          ),
          ListTile(
            leading: Icon(Icons.translate),
            title: Text('Translate'),
            onTap: () {
              // Handle translation (e.g., to Amharic)
              ScaffoldMessenger.of(
                context,
              ).showSnackBar(SnackBar(content: Text('Translating...')));
            },
          ),
        ],
      ),
    );
  }
}