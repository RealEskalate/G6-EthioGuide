import 'package:equatable/equatable.dart';



abstract class OnboardingEvent extends Equatable {
  const OnboardingEvent();

  @override
  List<Object> get props => [];
}

// Event to tell the BLoC to load the initial page data.
class LoadPages extends OnboardingEvent {}

// Event to tell the BLoC that the user has swiped to a new page.
class PageSwiped extends OnboardingEvent {
  final int newIndex;

  const PageSwiped(this.newIndex);

  @override
  List<Object> get props => [newIndex];
}