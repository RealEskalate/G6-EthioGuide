import 'package:ethioguide/core/components/bottom_nav_bar.dart';
import 'package:ethioguide/features/AI%20chat/Domain/entities/conversation.dart';
import 'package:ethioguide/features/AI%20chat/Presentation/bloc/ai_bloc.dart';
import 'package:ethioguide/features/AI%20chat/Presentation/widgets/ai_page_widgets.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:speech_to_text/speech_to_text.dart' as stt;

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
  // For speech to text
  late stt.SpeechToText _speech;
  bool _isListening = false;
  // For bottom naviagation bar
  final pageIndex = 2;

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
    // Initialize SpeechToText object
    _speech = stt.SpeechToText();
  }

  void _listen() async {
    if (!_isListening) {
      bool available = await _speech.initialize(
        onStatus: (val) => print('onStatus: $val'),
        onError: (val) => print('onError: $val'),
      );

      if (available) {
        setState(() => _isListening = true);
        _speech.listen(
          onResult: (val) => setState(() {
            _queryController.text = val.recognizedWords;
          }),
        );
      }
    } else {
      setState(() => _isListening = false);
      _speech.stop();
    }
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

  void _removeLoadingMessage({required bool queryFailed}) {
    setState(() {
      final conv = _history.removeLast();
      // Add Just the query to display as failed
      if (queryFailed) {
        _history.add(
          Conversation(
            id: conv.id,
            request: conv.request,
            response: '',
            source: 'failed',
            procedures: [],
          ),
        );
      }
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

  @override
  Widget build(BuildContext context) {
    return BlocListener<AiBloc, AiState>(
      listener: (context, state) {
        if (state is AiHistorySuccess) {
          setState(() {
            _history = [
              _history.first,
              ...state.history,
            ]; // Keep initial greeting
          });
          _scrollToBottom();
        } else if (state is AiQuerySuccess) {
          _removeLoadingMessage(queryFailed: false);
          _addMessage(conversation: state.conversation);
        } else if (state is AiError) {
          _removeLoadingMessage(queryFailed: true);
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
          final newConv = state.translated;
          for (int i = 0; i < _history.length; i++) {
            if (_history[i].id == state.id) {
              setState(() {
                final curConv = _history[i];
                _history[i] = Conversation(
                  id: curConv.id,
                  request: curConv.request,
                  response: newConv.response,
                  source: curConv.source,
                  procedures: newConv.procedures,
                );
              });
              break;
            }
          }
        }
      },
      child: Scaffold(
        appBar: appBar(context: context),
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
                      return conv.source != 'error'
                          ? conv.source != 'loading'
                                ? buildMessage(conv: conv, context: context)
                                : Column(
                                    children: [
                                      buildMessage(
                                        conv: conv,
                                        context: context,
                                      ),
                                      loadingCard(),
                                    ],
                                  )
                          : errorCard(
                              conv.response,
                            ); // Determine whether the response is error or response
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
                    child: Container(
                      decoration: BoxDecoration(
                        color: Colors.grey[50],
                        borderRadius: BorderRadius.circular(24),
                        boxShadow: [
                          BoxShadow(
                            color: Colors.grey.withOpacity(0.2),
                            spreadRadius: 1,
                            blurRadius: 4,
                            offset: const Offset(0, 2),
                          ),
                        ],
                      ),
                      child: TextField(
                        controller: _queryController,
                        focusNode: _queryFocusNode,
                        decoration: InputDecoration(
                          hintText: 'Type your question here...',
                          prefixIcon: IconButton(
                            onPressed: _listen,
                            icon: Icon(
                              _isListening ? Icons.mic : Icons.mic_none,
                            ),
                          ),
                          border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(18),
                            borderSide: BorderSide(
                              color: _queryFocusNode.hasFocus
                                  ? Colors.teal
                                  : Colors.grey,
                            ),
                          ),
                          focusedBorder: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(18),
                            borderSide: const BorderSide(
                              color: Colors.teal,
                              width: 1,
                            ),
                          ),
                          contentPadding: const EdgeInsets.symmetric(
                            horizontal: 16,
                            vertical: 12,
                          ),
                        ),
                        onSubmitted: (_) => _sendQuery(),
                      ),
                    ),
                  ),
                  const SizedBox(width: 8),
                  BlocBuilder<AiBloc, AiState>(
                    builder: (context, state) {
                      final isLoading = state is AiLoading;
                      return GestureDetector(
                        onTap: isLoading ? _cancelQuery : _sendQuery,
                        child: Container(
                          decoration: const BoxDecoration(
                            shape: BoxShape.circle,
                            color: Colors.teal,
                          ),
                          padding: const EdgeInsets.all(10),
                          child: isLoading
                              ? const SizedBox(
                                  width: 24,
                                  height: 24,
                                  child: Icon(
                                    Icons.square,
                                    color: Colors.white,
                                  ),
                                )
                              : const Icon(
                                  Icons.send,
                                  color: Colors.white,
                                  size: 24,
                                ),
                        ),
                      );
                    },
                  ),
                ],
              ),
            ),
          ],
        ),

        bottomNavigationBar: bottomNav(
          context: context,
          selectedIndex: pageIndex,
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

  void _cancelQuery() {
    context.read<AiBloc>().add(CancleQueryEvent());
  }

  @override
  void dispose() {
    _queryController.dispose();
    _scrollController.dispose();
    _queryFocusNode.dispose();
    super.dispose();
  }
}
