import 'package:flutter/material.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_detail.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_step.dart';
import 'package:go_router/go_router.dart';
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
    return const _WorkspaceProcedureDetailView();
  }
}

class _WorkspaceProcedureDetailView extends StatelessWidget {
  const _WorkspaceProcedureDetailView();

  @override
  Widget build(BuildContext context) {
    final mockDetail = ProcedureDetail(
      id: 'mock-1',
      title: 'Driver\'s License Renewal',
      organization: 'Transport Authority',
      status: ProcedureStatus.inProgress,
      progressPercentage: 40,
      documentsUploaded: 1,
      totalDocuments: 3,
      startDate: DateTime.now().subtract(const Duration(days: 10)),
      estimatedCompletion: DateTime.now().add(const Duration(days: 5)),
      completedDate: null,
      notes: 'Bring original ID.',
      steps: const [
        MyProcedureStep(id: 's1', title: 'Fill application form', description: 'Complete the online form.', isCompleted: true, order: 2),
        MyProcedureStep(id: 's2', title: 'Upload ID', description: 'Upload scanned national ID.', isCompleted: false, order: 1),
        MyProcedureStep(id: 's3', title: 'Pay fee', description: 'Pay via bank or mobile.', isCompleted: false , order: 3),
      ],
      estimatedTime: '2-3 hours',
      difficulty: 'Medium',
      officeType: 'Government Office',
      quickTips: const [
        'Go early to avoid queues',
        'Carry extra photocopies',
      ],
      requiredDocuments: const [
        'National ID',
        'Old License',
        'Passport Photo',
      ],
    );

    return Scaffold(
     /*  appBar: AppBar(
        title: Column(children: [
          
        ],),
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
      ), */
      body: _ProcedureDetailContent(procedureDetail: mockDetail),
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

          Row(children: [

            IconButton(onPressed: (){
              context.pop();
            },
             icon: Icon(Icons.arrow_back)
            ),
            const SizedBox(width: 8),

            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [

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

            ],)

          ],),
          // Page title
          
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
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('Progress action (mock)')),
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
