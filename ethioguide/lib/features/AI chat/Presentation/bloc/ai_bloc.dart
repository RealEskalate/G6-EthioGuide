import 'package:bloc/bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/AI%20chat/Domain/entities/conversation.dart';
import 'package:ethioguide/features/AI%20chat/Domain/entities/translated_conversation.dart';
import 'package:ethioguide/features/AI%20chat/Domain/usecases/get_history.dart';
import 'package:ethioguide/features/AI%20chat/Domain/usecases/send_query.dart';
import 'package:ethioguide/features/AI%20chat/Domain/usecases/translate_content.dart';
import 'package:ethioguide/features/AI%20chat/data/models/translated_conversation_model.dart';

part 'ai_event.dart';
part 'ai_state.dart';

class AiBloc extends Bloc<AiEvent, AiState> {
  final SendQuery sendQueryUseCase;
  final GetHistory getHistoryUseCase;
  final TranslateContent translateContentUseCase;

  AiBloc({
    required this.sendQueryUseCase,
    required this.getHistoryUseCase,
    required this.translateContentUseCase,
  }) : super(AiInitial()) {
    on<SendQueryEvent>(_onSendQuery);
    on<GetHistoryEvent>(_onGetHistory);
    on<TranslateContentEvent>(_onTranslateContent);
    on<CancleQueryEvent>(_onCancleQuery);
  }

  Future<void> _onSendQuery(SendQueryEvent event, Emitter<AiState> emit) async {
    emit(AiLoading());
    final result = await sendQueryUseCase(event.query);
    emit(
      (result.fold(
        (failure) => AiError(_mapFailureToMessage(failure)),
        (conversation) => AiQuerySuccess(conversation: conversation),
      )),
    );
  }

  Future<void> _onCancleQuery(CancleQueryEvent event, Emitter<AiState> emit) async {
    emit(AiInitial());
  }

  Future<void> _onGetHistory(
    GetHistoryEvent event,
    Emitter<AiState> emit,
  ) async {
    emit(AiLoading());
    final result = await getHistoryUseCase();
    emit(
      result.fold(
        (failure) => AiError(_mapFailureToMessage(failure)),
        (history) => AiHistorySuccess(history: history),
      ),
    );
  }

  Future<void> _onTranslateContent(
    TranslateContentEvent event,
    Emitter<AiState> emit,
  ) async {
    emit(AiLoading());
    final result = await translateContentUseCase(
      response: event.conversation.response, procedures: event.conversation.procedures
    );
    emit(
      result.fold(
        (failure) => AiError(_mapFailureToMessage(failure)),
        (translated) => AiTranslateSuccess(translated: translated, id: event.id),
      ),
    );
  }

  String _mapFailureToMessage(Failure failure) {
    if (failure is NetworkFailure) {
      return 'No internet connection. Please check your network.';
    } else if (failure is ServerFailure) {
      return failure.message;
    } else if (failure is CachedFailure) {
      return 'Unable to load cached data. Please try again.';
    }
    return 'Unexpected error. Please try again.';
  }
}
