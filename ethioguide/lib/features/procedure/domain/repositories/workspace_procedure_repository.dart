import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import '../entities/procedure_detail.dart';
import '../entities/procedure_step.dart';
import '../entities/workspace_procedure.dart';

/// Repository interface for workspace procedure operations
abstract class ProcedureDetailRepository {
  /// Get all workspace procedures
  Future<Either<Failure, List<ProcedureDetail>>> getProcedure();
  
  /// Get workspace summary statistics
  Future<Either<Failure, WorkspaceSummary>> getWorkspaceSummary();
  
  /// Get procedures by status
  Future<Either<Failure, List<ProcedureDetail>>> getProceduresByStatus(ProcedureStatus status);
  
  /// Get procedures by organization
  Future<Either<Failure, List<ProcedureDetail>>> getProceduresByOrganization(String organization);

  /// Fetch procedure details by ID
  Future<Either<String,List<MyProcedureStep>>> getProcedureDetail(String id);
  
  /// Update step status
  Future<Either<String, bool>> updateStepStatus(String procedureId, String stepId, bool isCompleted);
  
  /// Save progress
  Future<Either<String, bool>> saveProgress(String procedureId);
}
