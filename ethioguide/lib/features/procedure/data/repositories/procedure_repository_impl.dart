import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/exception.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/core/network/network_info.dart';
import 'package:ethioguide/features/procedure/data/datasources/procedure_remote_data_source.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure.dart';
import 'package:ethioguide/features/procedure/domain/repositories/procedure_repository.dart';

class ProcedureRepositoryImpl implements ProcedureRepository {
  final ProcedureRemoteDataSource remoteDataSource;
  final NetworkInfo networkInfo;

  ProcedureRepositoryImpl({required this.remoteDataSource, required this.networkInfo});

  @override
  Future<Either<Failure, List<Procedure>>> getProcedures(String? name ) async {
    final isOnline = await networkInfo.isConnected;
    if (!isOnline) {
      return const Left(NetworkFailure());
    }
    try {
      final models = await remoteDataSource.getProcedures(
        name
      );
      return Right(models);
    } on ServerException catch (e) {
      return Left(ServerFailure(message: e.message));
    } catch (_) {
      return const Left(ServerFailure());
    }
  }

  @override

  Future<Either<Failure, Procedure>> getProceduresbyid(String procedureId) async {
    final isOnline = await networkInfo.isConnected;
    if (!isOnline) {
      return const Left(NetworkFailure());
    }
    try {
      final models = await remoteDataSource.getProceduresbyid(procedureId);
      return Right(models);
    } on ServerException catch (e) {
      return Left(ServerFailure(message: e.message));
    } catch (_) {
      return const Left(ServerFailure());
    }
    
  }
  

  @override
  Future<Either<Failure, bool>> saveProcedure(String procedureId) async {
    final isOnline = await networkInfo.isConnected;
    if (!isOnline) {
      return const Left(NetworkFailure());
    }
    try {
      await remoteDataSource.saveProcedure(procedureId);
      return const Right(true);
    } on ServerException catch (e) {
      return Left(ServerFailure(message: e.message));
    } catch (_) {
      return const Left(ServerFailure());
    }
  }

  @override
  Future<Either<Failure, List<FeedbackItem>>> getFeedbacks(String procedureId) async {
    final isOnline = await networkInfo.isConnected;
    if (!isOnline) {
      return const Left(NetworkFailure());
    }
    try {
      final models = await remoteDataSource.getfeadback(procedureId);
      return Right(models);
    } on ServerException catch (e) {
      return Left(ServerFailure(message: e.message));
    } catch (_) {
      return const Left(ServerFailure());
    }
  }

  @override
  Future<Either<Failure, bool>> saveFeedback(
    String procedureId,
    String feedback,
    List<String> tags,
    String type,
  ) async {
    if (await networkInfo.isConnected) {
      try {
        final result = await remoteDataSource.savefeedback(
          procedureId,
          feedback,
          tags,
          type,
        );
        return Right(result);
      } catch (e) {
        return Left(ServerFailure());
      }
    }
    return Left(ServerFailure());
  }


}


