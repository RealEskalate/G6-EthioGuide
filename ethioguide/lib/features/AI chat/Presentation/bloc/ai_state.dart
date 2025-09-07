part of 'ai_bloc.dart';

sealed class AiState extends Equatable {
  const AiState();

  @override
  List<Object> get props => [];
}

final class AiInitial extends AiState {}

final class AiLoading extends AiState {}

class AiQuerySuccess extends AiState {
  final Conversation conversation;

  const AiQuerySuccess({required this.conversation});

  @override
  List<Object> get props => [conversation];
}

class AiHistorySuccess extends AiState {
  final List<Conversation> history;

  const AiHistorySuccess({required this.history});

  @override
  List<Object> get props => [history];
}

class AiTranslateSuccess extends AiState {
  final TranslatedConversation translated;
  final id;

  const AiTranslateSuccess({required this.translated, required this.id});

  @override
  List<Object> get props => [translated];
}

class AiError extends AiState {
  final String message;

  const AiError(this.message);

  @override
  List<Object> get props => [message];
}
