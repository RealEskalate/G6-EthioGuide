import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/core/network/network_info.dart';
import '../../domain/entities/procedure_detail.dart';
import '../../domain/entities/workspace_procedure.dart';
import '../../domain/repositories/workspace_procedure_repository.dart';
import '../datasources/workspace_procedure_local_datasource.dart';
import '../datasources/workspace_procedure_remote_data_source.dart';
import '../models/workspace_procedure_model.dart';
import '../models/workspace_summary_model.dart';

/// Repository implementation for workspace procedure operations
class WorkspaceProcedureRepositoryImpl implements WorkspaceProcedureRepository {
  final WorkspaceProcedureRemoteDataSource remoteDataSource;
  final WorkspaceProcedureLocalDataSource localDataSource;
  final NetworkInfo networkInfo;

  const WorkspaceProcedureRepositoryImpl({
    required this.remoteDataSource,
    required this.localDataSource,
    required this.networkInfo,
  });

  @override
  Future<Either<Failure, List<WorkspaceProcedure>>> getWorkspaceProcedures() async {
    try {
      if (await networkInfo.isConnected) {
        // Try to get from remote
        final remoteProcedures = await remoteDataSource.getWorkspaceProcedures();
        // Cache locally
        await localDataSource.cacheWorkspaceProcedures(remoteProcedures);
        return Right(remoteProcedures);
      } else {
        // Get from local cache
        final localProcedures = await localDataSource.getWorkspaceProcedures();
        if (localProcedures.isNotEmpty) {
          return Right(localProcedures);
        } else {
          return Left(ServerFailure('No internet connection and no cached data'));
        }
      }
    } catch (e) {
      // Try to get from local cache as fallback
      try {
        final localProcedures = await localDataSource.getWorkspaceProcedures();
        if (localProcedures.isNotEmpty) {
          return Right(localProcedures);
        }
      } catch (_) {}
      
      return Left(ServerFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, WorkspaceSummary>> getWorkspaceSummary() async {
    try {
      if (await networkInfo.isConnected) {
        // Try to get from remote
        final remoteSummary = await remoteDataSource.getWorkspaceSummary();
        // Cache locally
        await localDataSource.cacheWorkspaceSummary(remoteSummary);
        return Right(remoteSummary);
      } else {
        // Get from local cache
        final localSummary = await localDataSource.getWorkspaceSummary();
        if (localSummary != null) {
          return Right(localSummary);
        } else {
          return Left(ServerFailure('No internet connection and no cached data'));
        }
      }
    } catch (e) {
      // Try to get from local cache as fallback
      try {
        final localSummary = await localDataSource.getWorkspaceSummary();
        if (localSummary != null) {
          return Right(localSummary);
        }
      } catch (_) {}
      
      return Left(ServerFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, List<WorkspaceProcedure>>> getProceduresByStatus(ProcedureStatus status) async {
    try {
      if (await networkInfo.isConnected) {
        final remoteProcedures = await remoteDataSource.getProceduresByStatus(status.displayName);
        return Right(remoteProcedures);
      } else {
        final localProcedures = await localDataSource.getProceduresByStatus(status.displayName);
        return Right(localProcedures);
      }
    } catch (e) {
      return Left(ServerFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, List<WorkspaceProcedure>>> getProceduresByOrganization(String organization) async {
    try {
      if (await networkInfo.isConnected) {
        final remoteProcedures = await remoteDataSource.getProceduresByOrganization(organization);
        return Right(remoteProcedures);
      } else {
        final localProcedures = await localDataSource.getProceduresByOrganization(organization);
        return Right(localProcedures);
      }
    } catch (e) {
      return Left(ServerFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, WorkspaceProcedure>> createWorkspaceProcedure(WorkspaceProcedure procedure) async {
    try {
      if (await networkInfo.isConnected) {
        final createdProcedure = await remoteDataSource.createWorkspaceProcedure(
          WorkspaceProcedureModel(
            id: procedure.id,
            title: procedure.title,
            organization: procedure.organization,
            status: procedure.status,
            progressPercentage: procedure.progressPercentage,
            documentsUploaded: procedure.documentsUploaded,
            totalDocuments: procedure.totalDocuments,
            startDate: procedure.startDate,
            estimatedCompletion: procedure.estimatedCompletion,
            completedDate: procedure.completedDate,
            notes: procedure.notes,
          ),
        );
        
        // Cache locally
        await localDataSource.saveWorkspaceProcedure(createdProcedure);
        return Right(createdProcedure);
      } else {
        return Left(ServerFailure('No internet connection'));
      }
    } catch (e) {
      return Left(ServerFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, WorkspaceProcedure>> updateWorkspaceProcedure(WorkspaceProcedure procedure) async {
    try {
      if (await networkInfo.isConnected) {
        final updatedProcedure = await remoteDataSource.updateWorkspaceProcedure(
          WorkspaceProcedureModel(
            id: procedure.id,
            title: procedure.title,
            organization: procedure.organization,
            status: procedure.status,
            progressPercentage: procedure.progressPercentage,
            documentsUploaded: procedure.documentsUploaded,
            totalDocuments: procedure.totalDocuments,
            startDate: procedure.startDate,
            estimatedCompletion: procedure.estimatedCompletion,
            completedDate: procedure.completedDate,
            notes: procedure.notes,
          ),
        );
        
        // Update local cache
        await localDataSource.updateWorkspaceProcedure(updatedProcedure);
        return Right(updatedProcedure);
      } else {
        return Left(ServerFailure('No internet connection'));
      }
    } catch (e) {
      return Left(ServerFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, bool>> deleteWorkspaceProcedure(String id) async {
    try {
      if (await networkInfo.isConnected) {
        final success = await remoteDataSource.deleteWorkspaceProcedure(id);
        if (success) {
          await localDataSource.deleteWorkspaceProcedure(id);
        }
        return Right(success);
      } else {
        return Left(ServerFailure('No internet connection'));
      }
    } catch (e) {
      return Left(ServerFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, WorkspaceProcedure>> updateProgress(String id, int progressPercentage) async {
    try {
      if (await networkInfo.isConnected) {
        final updatedProcedure = await remoteDataSource.updateProgress(id, progressPercentage);
        
        // Update local cache
        await localDataSource.updateWorkspaceProcedure(updatedProcedure);
        return Right(updatedProcedure);
      } else {
        return Left(ServerFailure('No internet connection'));
      }
    } catch (e) {
      return Left(ServerFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, bool>> uploadDocument(String procedureId, String documentPath) async {
    try {
      if (await networkInfo.isConnected) {
        final success = await remoteDataSource.uploadDocument(procedureId, documentPath);
        return Right(success);
      } else {
        return Left(ServerFailure('No internet connection'));
      }
    } catch (e) {
      return Left(ServerFailure(e.toString()));
    }
  }

  // New methods for workspace procedure detail feature
  @override
  Future<Either<String, ProcedureDetail>> getProcedureDetail(String id) async {
    try {
      final result = await remoteDataSource.getProcedureDetail(id);
      return Right(result);
    } catch (e) {
      return Left(e.toString());
    }
  }

  @override
  Future<Either<String, bool>> updateStepStatus(String procedureId, String stepId, bool isCompleted) async {
    try {
      final result = await remoteDataSource.updateStepStatus(procedureId, stepId, isCompleted);
      return Right(result);
    } catch (e) {
      return Left(e.toString());
    }
  }

  @override
  Future<Either<String, bool>> saveProgress(String procedureId) async {
    try {
      final result = await remoteDataSource.saveProgress(procedureId);
      return Right(result);
    } catch (e) {
      return Left(e.toString());
    }
  }
}
