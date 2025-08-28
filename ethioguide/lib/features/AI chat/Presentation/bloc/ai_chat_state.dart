part of 'ai_chat_bloc.dart';

sealed class AiChatState extends Equatable {
  const AiChatState();
  
  @override
  List<Object> get props => [];
}

final class AiChatInitial extends AiChatState {}
