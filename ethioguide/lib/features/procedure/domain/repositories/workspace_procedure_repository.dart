import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import '../entities/procedure_detail.dart';
import '../entities/workspace_procedure.dart';

/// Repository interface for workspace procedure operations
abstract class WorkspaceProcedureRepository {
  /// Get all workspace procedures
  Future<Either<Failure, List<WorkspaceProcedure>>> getWorkspaceProcedures();
  
  /// Get workspace summary statistics
  Future<Either<Failure, WorkspaceSummary>> getWorkspaceSummary();
  
  /// Get procedures by status
  Future<Either<Failure, List<WorkspaceProcedure>>> getProceduresByStatus(ProcedureStatus status);
  
  /// Get procedures by organization
  Future<Either<Failure, List<WorkspaceProcedure>>> getProceduresByOrganization(String organization);
  
  /// Create a new workspace procedure
  Future<Either<Failure, WorkspaceProcedure>> createWorkspaceProcedure(WorkspaceProcedure procedure);
  
  /// Update an existing workspace procedure
  Future<Either<Failure, WorkspaceProcedure>> updateWorkspaceProcedure(WorkspaceProcedure procedure);
  
  /// Delete a workspace procedure
  Future<Either<Failure, bool>> deleteWorkspaceProcedure(String id);
  
  /// Update procedure progress
  Future<Either<Failure, WorkspaceProcedure>> updateProgress(String id, int progressPercentage);
  
  /// Upload document for a procedure
  Future<Either<Failure, bool>> uploadDocument(String procedureId, String documentPath);

  /// Fetch procedure details by ID
  Future<Either<String, ProcedureDetail>> getProcedureDetail(String id);
  
  /// Update step status
  Future<Either<String, bool>> updateStepStatus(String procedureId, String stepId, bool isCompleted);
  
  /// Save progress
  Future<Either<String, bool>> saveProgress(String procedureId);
}
