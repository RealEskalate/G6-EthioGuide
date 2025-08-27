import 'package:ethioguide/features/splashscreen/domain/repositories/splash_repository.dart';

class SplashRepositoryImpl implements SplashRepository {
  @override
  Future<void> completeSplash() async {
    // This is the concrete implementation: wait for 3 seconds.
    await Future.delayed(const Duration(seconds: 3));
  }
}