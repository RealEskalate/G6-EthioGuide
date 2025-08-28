import 'package:bloc_test/bloc_test.dart';
import 'package:mockito/mockito.dart';

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:ethioguide/features/onboarding/domain/entities/onboarding_page.dart';
import 'package:ethioguide/features/onboarding/domain/usecases/get_onboarding_pages.dart';
import 'package:ethioguide/features/onboarding/presentation/bloc/onboarding_bloc.dart';
import 'package:ethioguide/features/onboarding/presentation/bloc/onboarding_event.dart';
import 'package:ethioguide/features/onboarding/presentation/bloc/onboarding_state.dart';

@GenerateMocks([GetOnboardingPages])
import 'onboarding_bloc_test.mocks.dart';

void main() {
  late OnboardingBloc onboardingBloc;
  late MockGetOnboardingPages mockGetOnboardingPages;

  // Dummy data for our tests
  final tPages = [
    const OnboardingPage(icon: Icons.abc, title: 't1', subtitle: 's1', description: 'd1'),
    const OnboardingPage(icon: Icons.abc, title: 't2', subtitle: 's2', description: 'd2'),
    const OnboardingPage(icon: Icons.abc, title: 't3', subtitle: 's3', description: 'd3'),
  ];

  setUp(() {
    mockGetOnboardingPages = MockGetOnboardingPages();
    // Arrange: Mock the use case to return our test pages
    when(mockGetOnboardingPages()).thenReturn(tPages);
    
    onboardingBloc = OnboardingBloc(getOnboardingPages: mockGetOnboardingPages);
  });

  test('initial state should be an empty OnboardingState', () {
    expect(onboardingBloc.state, const OnboardingState());
  });

  group('LoadPages', () {
    blocTest<OnboardingBloc, OnboardingState>(
      'should get data from the use case and emit a new state',
      build: () => onboardingBloc,
      act: (bloc) => bloc.add(LoadPages()),
      expect: () => [
        OnboardingState(pages: tPages, pageIndex: 0, isLastPage: false),
      ],
      verify: (_) {
        verify(mockGetOnboardingPages()).called(1);
      },
    );
  });

  group('PageSwiped', () {
    blocTest<OnboardingBloc, OnboardingState>(
      'should update the pageIndex and isLastPage to false when not the last page',
      // Seed the BLoC with an initial loaded state
      seed: () => OnboardingState(pages: tPages, pageIndex: 0, isLastPage: false),
      build: () => onboardingBloc,
      act: (bloc) => bloc.add(const PageSwiped(1)), // Swipe to the middle page
      expect: () => [
        OnboardingState(pages: tPages, pageIndex: 1, isLastPage: false),
      ],
    );

    blocTest<OnboardingBloc, OnboardingState>(
      'should update the pageIndex and isLastPage to true when it is the last page',
      seed: () => OnboardingState(pages: tPages, pageIndex: 1, isLastPage: false),
      build: () => onboardingBloc,
      act: (bloc) => bloc.add(const PageSwiped(2)), // Swipe to the last page (index 2)
      expect: () => [
        OnboardingState(pages: tPages, pageIndex: 2, isLastPage: true),
      ],
    );
  });
}