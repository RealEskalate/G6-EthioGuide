import 'package:equatable/equatable.dart';

import '../../domain/entities/onboarding_page.dart';

// We only need one state for this feature since data loading is instant.
class OnboardingState extends Equatable {
  final List<OnboardingPage> pages;
  final int pageIndex;
  final bool isLastPage;

  const OnboardingState({
    this.pages = const [], // Default to an empty list
    this.pageIndex = 0,    // Default to the first page
    this.isLastPage = false,
  });

  // copyWith is a powerful tool to create a new state instance based on the old one.
  OnboardingState copyWith({
    List<OnboardingPage>? pages,
    int? pageIndex,
    bool? isLastPage,
  }) {
    return OnboardingState(
      pages: pages ?? this.pages,
      pageIndex: pageIndex ?? this.pageIndex,
      isLastPage: isLastPage ?? this.isLastPage,
    );
  }

  @override
  List<Object> get props => [pages, pageIndex, isLastPage];
}