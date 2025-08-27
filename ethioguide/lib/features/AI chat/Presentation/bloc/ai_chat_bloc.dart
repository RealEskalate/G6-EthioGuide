import 'package:bloc/bloc.dart';
import 'package:equatable/equatable.dart';

part 'ai_chat_event.dart';
part 'ai_chat_state.dart';

class AiChatBloc extends Bloc<AiChatEvent, AiChatState> {
  AiChatBloc() : super(AiChatInitial()) {
    on<AiChatEvent>((event, emit) {
      // TODO: implement event handler
    });
  }
}
