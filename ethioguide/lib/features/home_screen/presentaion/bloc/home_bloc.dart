import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:ethioguide/features/home_screen/domain/usecases/get_home_data.dart';
import 'home_event.dart';
import 'home_state.dart';

class HomeBloc extends Bloc<HomeEvent, HomeState> {
  final GetHomeData getHomeData;

  HomeBloc({required this.getHomeData}) : super(const HomeState()) {
    on<LoadHomeData>(_onLoadHomeData);
  }

  void _onLoadHomeData(LoadHomeData event, Emitter<HomeState> emit) {
    // Call the use case to get all the data lists
    final (quickActions, contentCards, popularServices) = getHomeData();
    // Emit a new state with the loaded data
    emit(HomeState(
      quickActions: quickActions,
      contentCards: contentCards,
      popularServices: popularServices,
    ));
  }
}