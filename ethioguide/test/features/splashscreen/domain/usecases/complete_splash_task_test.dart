import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';
import 'package:ethioguide/features/splashscreen/domain/repositories/splash_repository.dart';
import 'package:ethioguide/features/splashscreen/domain/usecases/complete_splash_task.dart';

@GenerateMocks([SplashRepository])
import 'complete_splash_task_test.mocks.dart';

void main() {
  late CompleteSplashTask usecase;
  late MockSplashRepository mockSplashRepository;

  setUp(() {
    mockSplashRepository = MockSplashRepository();
    usecase = CompleteSplashTask(mockSplashRepository);
  });

  test('should call completeSplash on the repository', () async {
    // Arrange
    when(mockSplashRepository.completeSplash()).thenAnswer((_) async => Future.value());

    // Act
    await usecase();

    // Assert
    verify(mockSplashRepository.completeSplash());
    verifyNoMoreInteractions(mockSplashRepository);
  });
}