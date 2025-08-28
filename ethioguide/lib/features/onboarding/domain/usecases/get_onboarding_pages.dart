import '../entities/onboarding_page.dart';
import '../repositories/onboarding_repository.dart';

class GetOnboardingPages {
  final OnboardingRepository repository;

  GetOnboardingPages(this.repository);

  List<OnboardingPage> call() {
    return repository.getOnboardingPages();
  }
}