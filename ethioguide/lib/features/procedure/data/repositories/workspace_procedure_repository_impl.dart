import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/core/network/network_info.dart';
import 'package:ethioguide/features/procedure/domain/repositories/workspace_procedure_repository.dart';
import '../../domain/entities/procedure_detail.dart';
import '../../domain/entities/workspace_procedure.dart';
import '../../domain/entities/procedure_step.dart';

import '../datasources/workspace_procedure_remote_data_source.dart';

/// Implementation of ProcedureDetailRepository using only remote data source
class WorkspaceProcedureRepositoryImpl implements ProcedureDetailRepository {
  final WorkspaceProcedureRemoteDataSource remoteDataSource;
  final NetworkInfo networkInfo;

  const WorkspaceProcedureRepositoryImpl({
    required this.remoteDataSource,
    required this.networkInfo,
  });

  @override
  Future<Either<Failure, List<ProcedureDetail>>> getProcedure() async {
    if (await networkInfo.isConnected) {
      try {
        final details = await remoteDataSource.getMyProcedures();
        return Right(details);
      } catch (e) {
        return Left(ServerFailure());
      }
    }
    return Left(ServerFailure());
  }

  @override
  Future<Either<Failure, WorkspaceSummary>> getWorkspaceSummary() async {
    if (await networkInfo.isConnected) {
      try {
        final summary = await remoteDataSource.getWorkspaceSummary();
        return Right(summary);
      } catch (e) {
        return Left(ServerFailure());
      }
    }
    return Left(ServerFailure());
  }

  @override
  Future<Either<Failure, List<ProcedureDetail>>> getProceduresByStatus(
    ProcedureStatus status,
  ) async {
    if (await networkInfo.isConnected) {
      try {
        final procedures = await remoteDataSource.getProceduresByStatus(
          status.displayName,
        );
        return Right(procedures);
      } catch (e) {
        return Left(ServerFailure());
      }
    }
    return Left(ServerFailure());
  }

  @override
  Future<Either<Failure, List<ProcedureDetail>>> getProceduresByOrganization(
    String organization,
  ) async {
    if (await networkInfo.isConnected) {
      try {
        final procedures = await remoteDataSource.getProceduresByOrganization(
          organization,
        );
        return Right(procedures);
      } catch (e) {
        return Left(ServerFailure());
      }
    }
    return Left(ServerFailure());
  }

  @override
  Future<Either<String,List<MyProcedureStep>>> getProcedureDetail(String id) async {
    if (await networkInfo.isConnected) {
      try {
        final detail = await remoteDataSource.getProcedureDetail(id);
        return Right(detail);
      } catch (e) {
        return Left(e.toString());
      }
    }
    return Left('No internet connection');
  }

  @override
  Future<Either<String, bool>> updateStepStatus(
    String procedureId,
    String stepId,
    bool isCompleted,
  ) async {
    if (await networkInfo.isConnected) {
      try {
        final result = await remoteDataSource.updateStepStatus(
          procedureId,
          stepId,
          isCompleted,
        );
        return Right(result);
      } catch (e) {
        return Left(e.toString());
      }
    }
    return Left('No internet connection');
  }

  @override
  Future<Either<String, bool>> saveProgress(String procedureId) async {
    if (await networkInfo.isConnected) {
      try {
        final result = await remoteDataSource.saveProgress(procedureId);
        return Right(result);
      } catch (e) {
        return Left(e.toString());
      }
    }
    return Left('No internet connection');
  }
}
