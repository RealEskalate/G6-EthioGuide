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

class GetHistoryEvent extends AiEvent {}

class TranslateContentEvent extends AiEvent {
  final String content;
  final String lang;

  const TranslateContentEvent(this.content, this.lang);

  @override
  List<Object> get props => [content, lang];
}
