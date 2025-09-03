import 'package:flutter_bloc/flutter_bloc.dart';
import '../../domain/usecases/get_onboarding_pages.dart';
import 'onboarding_event.dart';
import 'onboarding_state.dart';

class OnboardingBloc extends Bloc<OnboardingEvent, OnboardingState> {
  final GetOnboardingPages getOnboardingPages;

  OnboardingBloc({required this.getOnboardingPages}) : super(const OnboardingState()) {
    on<LoadPages>(_onLoadPages);
    on<PageSwiped>(_onPageSwiped);
  }

  void _onLoadPages(LoadPages event, Emitter<OnboardingState> emit) {
    final pages = getOnboardingPages();
    emit(state.copyWith(
      pages: pages,
      pageIndex: 0,
      isLastPage: pages.length == 1, // Handle edge case of only one page
    ));
  }

  void _onPageSwiped(PageSwiped event, Emitter<OnboardingState> emit) {
    emit(state.copyWith(
      pageIndex: event.newIndex,
      isLastPage: event.newIndex == state.pages.length - 1,
    ));
  }
}