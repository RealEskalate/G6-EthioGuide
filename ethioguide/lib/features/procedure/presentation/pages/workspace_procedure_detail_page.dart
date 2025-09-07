import 'package:ethioguide/features/procedure/presentation/bloc/workspace_procedure_detail_bloc.dart';
import 'package:flutter/material.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_detail.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_step.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import '../widgets/progress_overview_card.dart';
import '../widgets/step_instructions_list.dart';
import '../widgets/quick_tips_box.dart';
import '../widgets/required_documents_list.dart';

/// Page that displays detailed information about a workspace procedure


class WorkspaceProcedureDetailPage extends StatelessWidget {
final ProcedureDetail procedureDetail;

  const WorkspaceProcedureDetailPage({
    super.key,
    required this.procedureDetail,
  });


  @override
  Widget build(BuildContext context) {


    return Scaffold(
     /*  appBar: AppBar(
        title: Column(children: [
          
        ],),
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
      ), */
      body:  _ProcedureDetailContent(procedureDetail: procedureDetail),
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
            procedureDetail.procedure.title,
            style: Theme.of(context).textTheme.headlineSmall?.copyWith(
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'Complete guide to ${procedureDetail.procedure.title.toLowerCase()} in Ethiopia.',
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

           BlocBuilder<
                WorkspaceProcedureDetailBloc,
                WorkspaceProcedureDetailState
              >(
                builder: (context, state) {
                  if (state is ProcedureLoading) {
                    return const Center(child: CircularProgressIndicator());
                  } else if (state is ProcedureError) {
                    return Text(state.message);
                  } else if (state is ProcedureLoaded) {
                    return StepInstructionsList(procedureDetail: state.procedureDetail);
                  }
                  return const Center(child: Text('No procedures found'));
                },
              ),

          // Step-by-Step Instructions
          //StepInstructionsList(procedureDetail: procedureDetail.procedure),
          // const SizedBox(height: 20),

          // Quick Tips Box
          // QuickTipsBox(procedureDetail: procedureDetail),
          // const SizedBox(height: 20),

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
