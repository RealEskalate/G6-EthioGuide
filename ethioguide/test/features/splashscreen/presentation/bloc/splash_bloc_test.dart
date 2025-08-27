import 'package:bloc_test/bloc_test.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';
import 'package:ethioguide/features/splashscreen/domain/usecases/complete_splash_task.dart';
import 'package:ethioguide/features/splashscreen/presentation/bloc/splash_bloc.dart';
import 'package:ethioguide/features/splashscreen/presentation/bloc/splash_event.dart';
import 'package:ethioguide/features/splashscreen/presentation/bloc/splash_state.dart';


@GenerateMocks([CompleteSplashTask])
import 'splash_bloc_test.mocks.dart';

void main() {
  late SplashBloc splashBloc;
  late MockCompleteSplashTask mockCompleteSplashTask;

  setUp(() {
    mockCompleteSplashTask = MockCompleteSplashTask();
    splashBloc = SplashBloc(completeSplashTask: mockCompleteSplashTask);
  });

  test('initial state is SplashInitial', () {
    expect(splashBloc.state, equals(SplashInitial()));
  });

  blocTest<SplashBloc, SplashState>(
    'should emit [SplashCompleted] when StartSplashTimer event is added',
    build: () {
      // Arrange
      when(mockCompleteSplashTask()).thenAnswer((_) async => Future.value());
      return splashBloc;
    },
    act: (bloc) => bloc.add(StartSplashTimer()),
    expect: () => [
      SplashCompleted(),
    ],
    verify: (_) {
      verify(mockCompleteSplashTask()).called(1);
    },
  );
}