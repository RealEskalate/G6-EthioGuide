import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../bloc/workspace_procedure_detail_bloc.dart';
import '../widgets/progress_overview_card.dart';
import '../widgets/step_instructions_list.dart';
import '../widgets/quick_tips_box.dart';
import '../widgets/required_documents_list.dart';

/// Page that displays detailed information about a workspace procedure
class WorkspaceProcedureDetailPage extends StatelessWidget {
  final String procedureId;

  const WorkspaceProcedureDetailPage({
    super.key,
    required this.procedureId,
  });

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (context) => context.read<WorkspaceProcedureDetailBloc>()
        ..add(FetchProcedureDetail(procedureId)),
      child: const _WorkspaceProcedureDetailView(),
    );
  }
}

class _WorkspaceProcedureDetailView extends StatelessWidget {
  const _WorkspaceProcedureDetailView();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Procedure Details'),
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
      ),
      body: BlocConsumer<WorkspaceProcedureDetailBloc, WorkspaceProcedureDetailState>(
        listener: (context, state) {
          if (state is ProcedureError) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.message),
                backgroundColor: Colors.red,
              ),
            );
          } else if (state is ProgressSaved) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(
                  state.success 
                    ? 'Progress saved successfully!' 
                    : 'Failed to save progress',
                ),
                backgroundColor: state.success ? Colors.green : Colors.red,
              ),
            );
          }
        },
        builder: (context, state) {
          if (state is ProcedureLoading) {
            return const Center(
              child: CircularProgressIndicator(),
            );
          } else if (state is ProcedureError) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(
                    Icons.error_outline,
                    size: 64,
                    color: Colors.grey[400],
                  ),
                  const SizedBox(height: 16),
                  Text(
                    'Error loading procedure details',
                    style: Theme.of(context).textTheme.titleMedium,
                  ),
                  const SizedBox(height: 8),
                  Text(
                    state.message,
                    style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                      color: Colors.grey[600],
                    ),
                    textAlign: TextAlign.center,
                  ),
                  const SizedBox(height: 16),
                  ElevatedButton(
                    onPressed: () {
                      context.read<WorkspaceProcedureDetailBloc>().add(
                        FetchProcedureDetail(
                          context.read<WorkspaceProcedureDetailBloc>().state is ProcedureLoaded
                            ? (context.read<WorkspaceProcedureDetailBloc>().state as ProcedureLoaded).procedureDetail.id
                            : '',
                        ),
                      );
                    },
                    child: const Text('Retry'),
                  ),
                ],
              ),
            );
          } else if (state is ProcedureLoaded || state is StepStatusUpdated) {
            final procedureDetail = state is ProcedureLoaded
                ? state.procedureDetail
                : (state as StepStatusUpdated).procedureDetail;

            return _ProcedureDetailContent(procedureDetail: procedureDetail);
          } else if (state is ProgressSaved) {
            // Return to the last loaded state
            final bloc = context.read<WorkspaceProcedureDetailBloc>();
            if (bloc.state is ProcedureLoaded) {
              return _ProcedureDetailContent(
                procedureDetail: (bloc.state as ProcedureLoaded).procedureDetail,
              );
            }
            return const Center(child: Text('No procedure data available'));
          }

          return const Center(
            child: Text('No procedure data available'),
          );
        },
      ),
    );
  }
}

class _ProcedureDetailContent extends StatelessWidget {
  final ProcedureDetail procedureDetail;

  const _ProcedureDetailContent({
    required this.procedureDetail,
  });

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Page title
          Text(
            procedureDetail.title,
            style: Theme.of(context).textTheme.headlineSmall?.copyWith(
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'Complete guide to ${procedureDetail.title.toLowerCase()} in Ethiopia.',
            style: Theme.of(context).textTheme.bodyMedium?.copyWith(
              color: Colors.grey[600],
            ),
          ),
          const SizedBox(height: 24),

          // Progress Overview Card
          ProgressOverviewCard(procedureDetail: procedureDetail),
          const SizedBox(height: 20),

          // Step-by-Step Instructions
          StepInstructionsList(procedureDetail: procedureDetail),
          const SizedBox(height: 20),

          // Quick Tips Box
          QuickTipsBox(procedureDetail: procedureDetail),
          const SizedBox(height: 20),

          // Required Documents List
          RequiredDocumentsList(procedureDetail: procedureDetail),
          const SizedBox(height: 32),

          // Save Progress Button
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: () {
                context.read<WorkspaceProcedureDetailBloc>().add(
                  SaveProgress(procedureDetail.id),
                );
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.blue,
                foregroundColor: Colors.white,
                padding: const EdgeInsets.symmetric(vertical: 16),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Icon(Icons.save),
                  const SizedBox(width: 8),
                  Text(
                    'Save My Progress',
                    style: Theme.of(context).textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 20),
        ],
      ),
    );
  }
}
