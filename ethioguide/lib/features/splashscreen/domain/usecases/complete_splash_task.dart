import '../repositories/splash_repository.dart';

class CompleteSplashTask {
  final SplashRepository repository;

  CompleteSplashTask(this.repository);

  Future<void> call() async {
    return await repository.completeSplash();
  }
}