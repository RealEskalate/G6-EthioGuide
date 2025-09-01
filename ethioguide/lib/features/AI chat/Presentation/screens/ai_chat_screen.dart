import 'package:flutter/material.dart';

class Message {
  final bool isUser;
  final String text;
  final List<Widget>? cards; // For AI structured responses

  Message({required this.isUser, required this.text, this.cards});
}

class ChatPage extends StatefulWidget {
  const ChatPage({super.key});

  @override
  State<ChatPage> createState() => _ChatPageState();
}

class _ChatPageState extends State<ChatPage> {
  final List<Message> _messages = [];
  final TextEditingController _controller = TextEditingController();
  final FocusNode _focusNode = FocusNode();
  final ScrollController _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();

    // Add initial AI greeting
    _addMessage(
      isUser: false,
      text: '''Hello! I\'m your AI Assistant.
       I can help you navigate Ethiopian legal procedures,
       business registration, and more. What would you like to Know?''',
    );

    _focusNode.addListener(() {
      setState(() {}); // Rebuild to update border color on focus
    });
  }

  void _addMessage({
    required bool isUser,
    required String text,
    List<Widget>? cards,
  }) {
    setState(() {
      _messages.add(Message(isUser: isUser, text: text, cards: cards));
    });

    // Scroll to bottom
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _scrollController.animateTo(
        _scrollController.position.maxScrollExtent,
        duration: Duration(milliseconds: 300),
        curve: Curves.easeOut,
      );
    });
  }

//TODO:
      // Response card
      // Procedure card (should include two buttons to see procedure or start the procedure)
      // add query loading
      // add query failed or notice message if falure returned. (should be saved in the array)
  void _sendMessage() {
    if (_controller.text.isEmpty) return;
    String userQuery = _controller.text;
    _addMessage(isUser: true, text: userQuery);

    // TODO: Simulate AI response (replace with real API call in production)
    // Based on example: "How to register my business in Ethiopia?"
    if (userQuery.toLowerCase().contains('register my business')) {
      _addMessage(
        isUser: false,
        text:
            'I\'ll guide you through the complete business registration process in Ethiopia. Here\'s a step-by-step breakdown:',
        cards: [
          _buildStepCard(
            icon: Icons.description,
            title: 'Step 1: Prepare Documents',
            content: [
              'Valid ID or passport copy',
              'Business name reservation certificate',
              'Memorandum of Association',
              'Articles of Association',
            ],
          ),
          _buildStepCard(
            icon: Icons.payment,
            title: 'Step 2: Pay Fees',
            content: [
              'Registration Fee: 200 ETB',
              'Stamp Duty: 30 ETB',
              'Certificate Fee: 100 ETB',
            ],
          ),
          _buildStepCard(
            icon: Icons.location_on,
            title: 'Step 3: Visit Office',
            content: [
              'Submit documents at Regional Trade Office',
              'Wait for verification (5-7 business days)',
              'Collect business license certificate',
            ],
          ),
          _buildChecklistButton(), // "Save Checklist" with options
          // Quick Information cards
          _buildInfoCard(
            icon: Icons.info,
            title: 'Required Documents',
            content:
                'You\'ll need your national ID, proof of income, and residence certificate.',
          ),
          _buildInfoCard(
            icon: Icons.money,
            title: 'Processing Fee',
            content:
                'The application fee is 150 ETB, payable at the time of submission.',
          ),
          _buildInfoCard(
            icon: Icons.location_city,
            title: 'Office Location',
            content:
                'Visit your local district business office during business hours: 8:00 AM - 5:00 PM.',
          ),
        ],
      );
    } else {
      // Default response for other queries
      _addMessage(
        isUser: false,
        text: 'Thanks for your query! Here\'s some info...',
      );
    }

    _controller.clear();
    _focusNode.unfocus();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('AI Legal Assistant'),
        actions: [IconButton(icon: Icon(Icons.more_vert), onPressed: () {})],
      ),
      body: Column(
        children: [
          Expanded(
            child: ListView.builder(
              controller: _scrollController,
              padding: EdgeInsets.all(16),
              itemCount: _messages.length,
              itemBuilder: (context, index) {
                final msg = _messages[index];
                return Align(
                  alignment: msg.isUser
                      ? Alignment.centerRight
                      : Alignment.centerLeft,
                  child: Container(
                    margin: EdgeInsets.symmetric(vertical: 8),
                    padding: EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      color: msg.isUser ? Colors.teal : Colors.grey[200],
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Column(
                      crossAxisAlignment: msg.isUser
                          ? CrossAxisAlignment.end
                          : CrossAxisAlignment.start,
                      children: [
                        if (!msg.isUser)
                          Row(
                            children: [
                              Icon(
                                Icons.verified,
                                color: Colors.green,
                                size: 16,
                              ),
                              SizedBox(width: 4),
                              Text(
                                'Verified',
                                style: TextStyle(
                                  fontSize: 12,
                                  color: Colors.green,
                                ),
                              ),
                            ],
                          ),
                        Text(
                          msg.text,
                          style: TextStyle(
                            color: msg.isUser ? Colors.white : Colors.black,
                          ),
                        ),
                        if (msg.cards != null) ...msg.cards!,
                      ],
                    ),
                  ),
                );
              },
            ),
          ),
          Padding(
            padding: EdgeInsets.all(16),
            child: Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _controller,
                    focusNode: _focusNode,
                    decoration: InputDecoration(
                      hintText: 'Type your question here...',
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(24),
                        borderSide: BorderSide(
                          color: _focusNode.hasFocus
                              ? Colors.blue
                              : Colors.grey,
                        ),
                      ),
                      focusedBorder: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(24),
                        borderSide: BorderSide(color: Colors.blue, width: 2),
                      ),
                    ),
                    onSubmitted: (_) => _sendMessage(),
                  ),
                ),
                IconButton(
                  icon: Icon(Icons.send, color: Colors.teal),
                  onPressed: _sendMessage,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  //* custom widgtes
  Widget _buildStepCard({
    required IconData icon,
    required String title,
    required List<String> content,
  }) {
    return Card(
      color: Colors.teal[50],
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      margin: EdgeInsets.symmetric(vertical: 8),
      child: Padding(
        padding: EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(icon, color: Colors.teal),
                SizedBox(width: 8),
                Text(
                  title,
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                    color: Colors.teal,
                  ),
                ),
              ],
            ),
            SizedBox(height: 8),
            ...content.map(
              (item) => Text('• $item', style: TextStyle(fontSize: 14)),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoCard({
    required IconData icon,
    required String title,
    required String content,
  }) {
    return Card(
      color: Colors.yellow[50],
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      margin: EdgeInsets.symmetric(vertical: 4),
      child: Padding(
        padding: EdgeInsets.all(12),
        child: Row(
          children: [
            Icon(icon, color: Colors.yellow[700]),
            SizedBox(width: 8),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                  ),
                  SizedBox(height: 4),
                  Text(content, style: TextStyle(fontSize: 14)),
                ],
              ),
            ),
          ],
        ),
      ),
    );
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
