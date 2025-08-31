import 'package:bloc_test/bloc_test.dart';
import 'package:dartz/dartz.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_workspace_procedures.dart';
import 'package:ethioguide/features/procedure/domain/usecases/get_workspace_summary.dart';
import 'package:ethioguide/features/procedure/presentation/bloc/workspace_procedure_bloc.dart';

import 'workspace_procedure_bloc_test.mocks.dart';

@GenerateMocks([GetWorkspaceProcedures, GetWorkspaceSummary])
void main() {
  late WorkspaceProcedureBloc bloc;
  late MockGetWorkspaceProcedures mockGetWorkspaceProcedures;
  late MockGetWorkspaceSummary mockGetWorkspaceSummary;

  setUp(() {
    mockGetWorkspaceProcedures = MockGetWorkspaceProcedures();
    mockGetWorkspaceSummary = MockGetWorkspaceSummary();
    bloc = WorkspaceProcedureBloc(
      getWorkspaceProcedures: mockGetWorkspaceProcedures,
      getWorkspaceSummary: mockGetWorkspaceSummary,
    );
  });

  tearDown(() {
    bloc.close();
  });

  test('initial state should be WorkspaceProcedureInitial', () {
    expect(bloc.state, equals(WorkspaceProcedureInitial()));
  });

  group('LoadWorkspaceProcedures', () {
    final tProcedures = [
      const WorkspaceProcedure(
        id: '1',
        title: 'Test Procedure',
        organization: 'Test Org',
        status: ProcedureStatus.inProgress,
        progressPercentage: 50,
        documentsUploaded: 2,
        totalDocuments: 4,
        startDate: DateTime(2024, 1, 1),
      ),
    ];

    blocTest<WorkspaceProcedureBloc, WorkspaceProcedureState>(
      'emits [WorkspaceProcedureLoading, WorkspaceProceduresLoaded] when successful',
      build: () {
        when(mockGetWorkspaceProcedures(any))
            .thenAnswer((_) async => Right(tProcedures));
        return bloc;
      },
      act: (bloc) => bloc.add(const LoadWorkspaceProcedures()),
      expect: () => [
        WorkspaceProcedureLoading(),
        WorkspaceProceduresLoaded(tProcedures),
      ],
      verify: (_) {
        verify(mockGetWorkspaceProcedures(NoParams()));
      },
    );

    blocTest<WorkspaceProcedureBloc, WorkspaceProcedureState>(
      'emits [WorkspaceProcedureLoading, WorkspaceProcedureError] when failure',
      build: () {
        when(mockGetWorkspaceProcedures(any))
            .thenAnswer((_) async => Left(ServerFailure('Server error')));
        return bloc;
      },
      act: (bloc) => bloc.add(const LoadWorkspaceProcedures()),
      expect: () => [
        WorkspaceProcedureLoading(),
        const WorkspaceProcedureError('Server error'),
      ],
      verify: (_) {
        verify(mockGetWorkspaceProcedures(NoParams()));
      },
    );
  });

  group('LoadWorkspaceSummary', () {
    final tSummary = const WorkspaceSummary(
      totalProcedures: 5,
      inProgress: 2,
      completed: 3,
      totalDocuments: 15,
    );

    blocTest<WorkspaceProcedureBloc, WorkspaceProcedureState>(
      'emits [WorkspaceSummaryLoaded] when successful',
      build: () {
        when(mockGetWorkspaceSummary(any))
            .thenAnswer((_) async => Right(tSummary));
        return bloc;
      },
      act: (bloc) => bloc.add(const LoadWorkspaceSummary()),
      expect: () => [
        WorkspaceSummaryLoaded(tSummary),
      ],
      verify: (_) {
        verify(mockGetWorkspaceSummary(NoParams()));
      },
    );

    blocTest<WorkspaceProcedureBloc, WorkspaceProcedureState>(
      'emits [WorkspaceSummaryError] when failure',
      build: () {
        when(mockGetWorkspaceSummary(any))
            .thenAnswer((_) async => Left(ServerFailure('Server error')));
        return bloc;
      },
      act: (bloc) => bloc.add(const LoadWorkspaceSummary()),
      expect: () => [
        const WorkspaceSummaryError('Server error'),
      ],
      verify: (_) {
        verify(mockGetWorkspaceSummary(NoParams()));
      },
    );
  });

  group('FilterProceduresByStatus', () {
    final tProcedures = [
      const WorkspaceProcedure(
        id: '1',
        title: 'Test Procedure',
        organization: 'Test Org',
        status: ProcedureStatus.inProgress,
        progressPercentage: 50,
        documentsUploaded: 2,
        totalDocuments: 4,
        startDate: DateTime(2024, 1, 1),
      ),
      const WorkspaceProcedure(
        id: '2',
        title: 'Another Procedure',
        organization: 'Test Org',
        status: ProcedureStatus.completed,
        progressPercentage: 100,
        documentsUploaded: 4,
        totalDocuments: 4,
        startDate: DateTime(2024, 1, 1),
      ),
    ];

    blocTest<WorkspaceProcedureBloc, WorkspaceProcedureState>(
      'emits [WorkspaceProceduresFiltered] when filtering by status',
      build: () {
        return bloc;
      },
      seed: () => WorkspaceProceduresLoaded(tProcedures),
      act: (bloc) => bloc.add(const FilterProceduresByStatus(ProcedureStatus.inProgress)),
      expect: () => [
        WorkspaceProceduresFiltered([tProcedures[0]], ProcedureStatus.inProgress),
      ],
    );

    blocTest<WorkspaceProcedureBloc, WorkspaceProcedureState>(
      'does not emit when state is not WorkspaceProceduresLoaded',
      build: () {
        return bloc;
      },
      seed: () => WorkspaceProcedureInitial(),
      act: (bloc) => bloc.add(const FilterProceduresByStatus(ProcedureStatus.inProgress)),
      expect: () => [],
    );
  });

  group('FilterProceduresByOrganization', () {
    final tProcedures = [
      const WorkspaceProcedure(
        id: '1',
        title: 'Test Procedure',
        organization: 'Org A',
        status: ProcedureStatus.inProgress,
        progressPercentage: 50,
        documentsUploaded: 2,
        totalDocuments: 4,
        startDate: DateTime(2024, 1, 1),
      ),
      const WorkspaceProcedure(
        id: '2',
        title: 'Another Procedure',
        organization: 'Org B',
        status: ProcedureStatus.completed,
        progressPercentage: 100,
        documentsUploaded: 4,
        totalDocuments: 4,
        startDate: DateTime(2024, 1, 1),
      ),
    ];

    blocTest<WorkspaceProcedureBloc, WorkspaceProcedureState>(
      'emits [WorkspaceProceduresFiltered] when filtering by organization',
      build: () {
        return bloc;
      },
      seed: () => WorkspaceProceduresLoaded(tProcedures),
      act: (bloc) => bloc.add(const FilterProceduresByOrganization('Org A')),
      expect: () => [
        WorkspaceProceduresFiltered([tProcedures[0]], null, 'Org A'),
      ],
    );
  });

  group('RefreshWorkspaceData', () {
    blocTest<WorkspaceProcedureBloc, WorkspaceProcedureState>(
      'emits [WorkspaceProcedureLoading, WorkspaceProceduresLoaded] when refreshing',
      build: () {
        when(mockGetWorkspaceProcedures(any))
            .thenAnswer((_) async => const Right([]));
        when(mockGetWorkspaceSummary(any))
            .thenAnswer((_) async => const Right(WorkspaceSummary(
                  totalProcedures: 0,
                  inProgress: 0,
                  completed: 0,
                  totalDocuments: 0,
                )));
        return bloc;
      },
      act: (bloc) => bloc.add(const RefreshWorkspaceData()),
      expect: () => [
        WorkspaceProcedureLoading(),
        const WorkspaceProceduresLoaded([]),
        const WorkspaceSummaryLoaded(WorkspaceSummary(
          totalProcedures: 0,
          inProgress: 0,
          completed: 0,
          totalDocuments: 0,
        )),
      ],
      verify: (_) {
        verify(mockGetWorkspaceProcedures(NoParams()));
        verify(mockGetWorkspaceSummary(NoParams()));
      },
    );
  });
}


