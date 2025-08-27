import 'package:flutter_bloc/flutter_bloc.dart';
import '../../domain/usecases/complete_splash_task.dart';
import 'splash_event.dart';
import 'splash_state.dart';

class SplashBloc extends Bloc<SplashEvent, SplashState> {
  final CompleteSplashTask completeSplashTask;

  // 1. The BLoC starts with the SplashInitial state.
  SplashBloc({required this.completeSplashTask}) : super(SplashInitial()) {
    // 2. Register the event handler.
    on<StartSplashTimer>(_onStartSplashTimer);
  }

  // 3. This function runs when a StartSplashTimer event is added.
  Future<void> _onStartSplashTimer(
    StartSplashTimer event,
    Emitter<SplashState> emit,
  ) async {
    // 4. It calls our business logic (the use case).
    await completeSplashTask();
    
    // 5. After the logic is complete, it emits the new state to the UI.
    emit(SplashCompleted());
  }
}