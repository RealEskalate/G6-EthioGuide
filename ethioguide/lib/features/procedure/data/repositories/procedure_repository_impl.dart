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
  Future<Either<Failure, List<Procedure>>> getProcedures() async {
    final isOnline = await networkInfo.isConnected;
    if (!isOnline) {
      return const Left(NetworkFailure());
    }
    try {
      final models = await remoteDataSource.getProcedures();
      return Right(models);
    } on ServerException catch (e) {
      return Left(ServerFailure(message: e.message));
    } catch (_) {
      return const Left(ServerFailure());
    }
  }
}


