import 'package:ethioguide/features/procedure/data/models/workspace_procedure_model.dart';
import 'package:ethioguide/features/procedure/data/models/workspace_summary_model.dart';
import 'package:ethioguide/features/procedure/data/sample_data/workspace_sample_data.dart';

/// Local data source for workspace procedures (cache/local storage)
abstract class WorkspaceProcedureLocalDataSource {
  /// Get all workspace procedures from local storage
  Future<List<WorkspaceProcedureModel>> getWorkspaceProcedures();
  
  /// Get workspace summary from local storage
  Future<WorkspaceSummaryModel?> getWorkspaceSummary();
  
  /// Cache workspace procedures locally
  Future<void> cacheWorkspaceProcedures(List<WorkspaceProcedureModel> procedures);
  
  /// Cache workspace summary locally
  Future<void> cacheWorkspaceSummary(WorkspaceSummaryModel summary);
  
  /// Get procedures by status from local storage
  Future<List<WorkspaceProcedureModel>> getProceduresByStatus(String status);
  
  /// Get procedures by organization from local storage
  Future<List<WorkspaceProcedureModel>> getProceduresByOrganization(String organization);
  
  /// Save a workspace procedure locally
  Future<void> saveWorkspaceProcedure(WorkspaceProcedureModel procedure);
  
  /// Update a workspace procedure locally
  Future<void> updateWorkspaceProcedure(WorkspaceProcedureModel procedure);
  
  /// Delete a workspace procedure locally
  Future<void> deleteWorkspaceProcedure(String id);
  
  /// Clear all cached data
  Future<void> clearCache();
}

/// Implementation of local data source using in-memory storage
class WorkspaceProcedureLocalDataSourceImpl implements WorkspaceProcedureLocalDataSource {
  // In-memory storage for demo purposes
  // In a real app, this would use SharedPreferences, Hive, or SQLite
  static final Map<String, WorkspaceProcedureModel> _procedures = {};
  static WorkspaceSummaryModel? _summary;
  
  // Initialize with sample data
static Future<void> init() async{
    final sampleProcedures = WorkspaceSampleData.getSampleProcedures();
    for (final procedure in sampleProcedures) {
      _procedures[procedure.id] = procedure;
    }
    _summary = WorkspaceSampleData.getSampleSummary();
  }

  @override
  Future<List<WorkspaceProcedureModel>> getWorkspaceProcedures() async {
    return _procedures.values.toList();
  }

  @override
  Future<WorkspaceSummaryModel?> getWorkspaceSummary() async {
    return _summary;
  }

  @override
  Future<void> cacheWorkspaceProcedures(List<WorkspaceProcedureModel> procedures) async {
    _procedures.clear();
    for (final procedure in procedures) {
      _procedures[procedure.id] = procedure;
    }
  }

  @override
  Future<void> cacheWorkspaceSummary(WorkspaceSummaryModel summary) async {
    _summary = summary;
  }

  @override
  Future<List<WorkspaceProcedureModel>> getProceduresByStatus(String status) async {
    return _procedures.values
        .where((procedure) => procedure.status.displayName == status)
        .toList();
  }

  @override
  Future<List<WorkspaceProcedureModel>> getProceduresByOrganization(String organization) async {
    return _procedures.values
        .where((procedure) => procedure.organization == organization)
        .toList();
  }

  @override
  Future<void> saveWorkspaceProcedure(WorkspaceProcedureModel procedure) async {
    _procedures[procedure.id] = procedure;
  }

  @override
  Future<void> updateWorkspaceProcedure(WorkspaceProcedureModel procedure) async {
    if (_procedures.containsKey(procedure.id)) {
      _procedures[procedure.id] = procedure;
    }
  }

  @override
  Future<void> deleteWorkspaceProcedure(String id) async {
    _procedures.remove(id);
  }

  @override
  Future<void> clearCache() async {
    _procedures.clear();
    _summary = null;
  }
}
