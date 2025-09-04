import 'package:bloc_test/bloc_test.dart';
import 'package:ethioguide/features/home_screen/domain/entities/home_data.dart';
import 'package:ethioguide/features/home_screen/domain/usecases/get_home_data.dart';

import 'package:ethioguide/features/home_screen/presentaion/bloc/home_bloc.dart';
import 'package:ethioguide/features/home_screen/presentaion/bloc/home_event.dart';
import 'package:ethioguide/features/home_screen/presentaion/bloc/home_state.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

@GenerateMocks([GetHomeData])
import 'home_bloc_test.mocks.dart';

void main() {
  late HomeBloc homeBloc;
  late MockGetHomeData mockGetHomeData;

  setUp(() {
    mockGetHomeData = MockGetHomeData();
    homeBloc = HomeBloc(getHomeData: mockGetHomeData);
  });

  // Create some dummy data to be returned by the mock use case
  final tQuickActions = [
    const QuickAction(icon: Icons.add, title: 't1', subtitle: 's1', routeName: 'r1'),
  ];
  final tContentCards = [
    const ContentCard(sectionTitle: 'st1', icon: Icons.pages, title: 't1', subtitle: 's1', details: [], routeName: 'r1'),
  ];
  final tPopularServices = [
    const PopularService(icon: Icons.abc, title: 't1', category: 'c1', timeEstimate: '1hr', routeName: 'r1'),
  ];

  test('initial state should be an empty HomeState', () {
    expect(homeBloc.state, const HomeState());
  });

  blocTest<HomeBloc, HomeState>(
    'should call GetHomeData use case and emit HomeState with data when LoadHomeData is added',
    build: () {
      // Arrange: Program the mock use case to return our dummy data.
      when(mockGetHomeData())
          .thenReturn((tQuickActions, tContentCards, tPopularServices));
      return homeBloc;
    },
    act: (bloc) => bloc.add(LoadHomeData()),
    expect: () => [
      // Assert: Expect the BLoC to emit a single state containing all the data.
      HomeState(
        quickActions: tQuickActions,
        contentCards: tContentCards,
        popularServices: tPopularServices,
      ),
    ],
    verify: (_) {
      // Verify that the use case's call() method was executed exactly once.
      verify(mockGetHomeData()).called(1);
    },
  );
}