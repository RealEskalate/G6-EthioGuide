part of 'ai_bloc.dart';

sealed class AiEvent extends Equatable {
  const AiEvent();

  @override
  List<Object> get props => [];
}

class SendQueryEvent extends AiEvent {
  final String query;

  const SendQueryEvent({required this.query});

  @override
  List<Object> get props => [query];
}

class CancleQueryEvent extends AiEvent {}

class GetHistoryEvent extends AiEvent {}

class TranslateContentEvent extends AiEvent {
  final TranslatedConversationModel conversation;
  final id;

  const TranslateContentEvent({required this.conversation, required this.id});
}
