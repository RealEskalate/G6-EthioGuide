import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';
import 'package:ethioguide/features/onboarding/domain/entities/onboarding_page.dart';
import 'package:ethioguide/features/onboarding/domain/repositories/onboarding_repository.dart';
import 'package:ethioguide/features/onboarding/domain/usecases/get_onboarding_pages.dart';

@GenerateMocks([OnboardingRepository])
import 'get_onboarding_pages_test.mocks.dart';

void main() {
  late GetOnboardingPages usecase;
  late MockOnboardingRepository mockOnboardingRepository;

  setUp(() {
    mockOnboardingRepository = MockOnboardingRepository();
    usecase = GetOnboardingPages(mockOnboardingRepository);
  });

  final testPages = [
    const OnboardingPage(
      icon: Icons.navigation,
      title: 'Test Title 1',
      subtitle: 'Test Subtitle 1',
      description: 'Test Description 1',
    ),
  ];

  
  test('should get list of onboarding pages from the repository', () {
    // Arrange: When `getOnboardingPages` is called, return our test data directly.
    when(mockOnboardingRepository.getOnboardingPages()).thenReturn(testPages);

    // Act: Execute the use case.
    final result = usecase();

    // Assert
    expect(result, testPages);
    verify(mockOnboardingRepository.getOnboardingPages());
    verifyNoMoreInteractions(mockOnboardingRepository);
  });
}