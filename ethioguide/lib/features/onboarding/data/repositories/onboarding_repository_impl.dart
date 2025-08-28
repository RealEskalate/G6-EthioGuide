import 'package:flutter/material.dart';
import 'package:ethioguide/features/onboarding/domain/entities/onboarding_page.dart';
import 'package:ethioguide/features/onboarding/domain/repositories/onboarding_repository.dart';

class OnboardingRepositoryImpl implements OnboardingRepository {
  @override
  List<OnboardingPage> getOnboardingPages() {
    // This is a hardcoded list of static data built into the app.
    return [
      const OnboardingPage(
        icon: Icons.apps_outlined,
        title: 'EthioGuide',
        subtitle: 'Government Services Navigator',
        description:
            'Your comprehensive guide to Ethiopian government procedures. Access expert knowledge and step-by-step assistance for all administrative processes.',
      ),

      const OnboardingPage(
        icon: Icons.book_outlined,
        title: 'Expert Information',
        subtitle: 'Professional Guidance',
        description:
            'Access detailed procedures, requirements, and official information for all Ethiopian government services in an organized, easy-to-understand format.',
      ),
      const OnboardingPage(
        icon: Icons.navigation_outlined,
        title: 'Navigate with Confidence',
        subtitle: 'Comprehensive Support',
        description:
            'Find government offices, understand procedures, and complete your administrative tasks efficiently with our comprehensive guidance system.',
      ),
    ];
  }
}
