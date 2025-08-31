import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:ethioguide/core/config/app_color.dart';
import 'package:ethioguide/core/config/app_theme.dart';
import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';
import 'package:ethioguide/features/procedure/presentation/bloc/workspace_procedure_bloc.dart';
import 'package:ethioguide/features/procedure/presentation/widgets/workspace_summary_cards.dart';
import 'package:ethioguide/features/procedure/presentation/widgets/workspace_procedure_card.dart';
import 'package:ethioguide/features/procedure/presentation/widgets/workspace_filters.dart';

/// Page that displays the workspace with procedures tracking
class WorkspacePage extends StatefulWidget {
  const WorkspacePage({super.key});

  @override
  State<WorkspacePage> createState() => _WorkspacePageState();
}

class _WorkspacePageState extends State<WorkspacePage> {
  ProcedureStatus? selectedStatus;
  String? selectedOrganization;

  @override
  void initState() {
    super.initState();
    context.read<WorkspaceProcedureBloc>().add(const LoadWorkspaceProcedures());
    context.read<WorkspaceProcedureBloc>().add(const LoadWorkspaceSummary());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: AppTheme.lightTheme.appBarTheme.backgroundColor,
        elevation: 0,
        toolbarHeight: 90,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.maybePop(context),
        ),
        title: Expanded(
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 4.0, vertical: 20),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'My Workspace',
                  style: Theme.of(context).textTheme.headlineSmall,
                  softWrap: true,
                  overflow: TextOverflow.visible,
                ),
                const SizedBox(height: 6),
                Text(
                  'Track and manage your ongoing procedures.',
                  style: Theme.of(context)
                      .textTheme
                      .bodyMedium
                      ?.copyWith(color: AppColors.graycolor),
                  softWrap: true,
                  overflow: TextOverflow.visible,
                ),
              ],
            ),
          ),
        ),
        actions: [
          Padding(
            padding: const EdgeInsets.only(right: 16.0, top: 20, bottom: 20),
            child: ElevatedButton.icon(
              onPressed: () {},
              icon: const Icon(Icons.add, color: Colors.white),
              label: const Text('Add New', style: TextStyle(color: Colors.white)),
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.teal,
                foregroundColor: Colors.white,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(8),
                ),
              ),
            ),
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          context.read<WorkspaceProcedureBloc>().add(const RefreshWorkspaceData());
        },
        child: SingleChildScrollView(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const SizedBox(height: 12),
              BlocBuilder<WorkspaceProcedureBloc, WorkspaceProcedureState>(
                builder: (context, state) {
                  if (state is WorkspaceSummaryLoaded) {
                    return WorkspaceSummaryCards(summary: state.summary);
                  }
                  return const WorkspaceSummaryCards(
                    summary: WorkspaceSummary(
                      totalProcedures: 0,
                      inProgress: 0,
                      completed: 0,
                      totalDocuments: 0,
                    ),
                  );
                },
              ),
              const SizedBox(height: 24),
              WorkspaceFilters(
                selectedStatus: selectedStatus,
                selectedOrganization: selectedOrganization,
                onStatusChanged: (status) {
                  setState(() => selectedStatus = status);
                  if (status != null) {
                    context.read<WorkspaceProcedureBloc>().add(
                      FilterProceduresByStatus(status),
                    );
                  }
                },
                onOrganizationChanged: (organization) {
                  setState(() => selectedOrganization = organization);
                  if (organization != null) {
                    context.read<WorkspaceProcedureBloc>().add(
                      FilterProceduresByOrganization(organization),
                    );
                  }
                },
              ),
              const SizedBox(height: 24),
              BlocBuilder<WorkspaceProcedureBloc, WorkspaceProcedureState>(
                builder: (context, state) {
                  if (state is WorkspaceProcedureLoading) {
                    return const Center(child: CircularProgressIndicator());
                  } else if (state is WorkspaceProceduresLoaded) {
                    return _buildProceduresList(state.procedures);
                  } else if (state is WorkspaceProceduresFiltered) {
                    return _buildProceduresList(state.procedures);
                  }
                  return const Center(child: Text('No procedures found'));
                },
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildProceduresList(List<WorkspaceProcedure> procedures) {
    if (procedures.isEmpty) {
      return const Center(
        child: Column(
          children: [
            Icon(Icons.inbox_outlined, size: 64, color: Colors.grey),
            SizedBox(height: 16),
            Text('No procedures found', style: TextStyle(fontSize: 18, color: Colors.grey)),
          ],
        ),
      );
    }

    return ListView.separated(
      physics: const NeverScrollableScrollPhysics(),
      shrinkWrap: true,
      itemCount: procedures.length,
      separatorBuilder: (_, __) => const SizedBox(height: 16),
      itemBuilder: (context, index) {
        final procedure = procedures[index];
        return WorkspaceProcedureCard(
          procedure: procedure,
          onTap: () {},
        );
      },
    );
  }
}
