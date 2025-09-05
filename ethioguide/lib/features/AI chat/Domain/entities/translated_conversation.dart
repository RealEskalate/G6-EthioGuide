import 'package:equatable/equatable.dart';
import 'package:ethioguide/features/AI%20chat/Domain/entities/conversation.dart';

class TranslatedConversation extends Equatable{
  final String response;
  final List<Procedure?> procedures;

  const TranslatedConversation({required this.response, required this.procedures});
  
  @override
  List<Object?> get props => [response, procedures];
}

