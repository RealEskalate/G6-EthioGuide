import 'package:ethioguide/core/components/button.dart';
import 'package:ethioguide/core/config/app_color.dart';
import 'package:ethioguide/core/config/route_names.dart';
import 'package:ethioguide/features/procedure/presentation/bloc/procedure_bloc.dart';
import 'package:ethioguide/features/procedure/presentation/pages/Discussin_tab.dart';
import 'package:ethioguide/features/procedure/presentation/pages/Feedback_tab.dart';
import 'package:ethioguide/features/procedure/presentation/pages/Notice_tab.dart';
import 'package:ethioguide/features/procedure/presentation/widgets/docment_form.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import '../../domain/entities/procedure.dart';
import '../widgets/procedure_detail_header.dart';
import '../widgets/document_list.dart';
import '../widgets/step_list.dart';

class ProcedureDetailPage extends StatelessWidget {
  const ProcedureDetailPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Procedure Detail')),
      body: _DummyDetailBody(),
    );
  }
}

class _DummyDetailBody extends StatelessWidget {
  const _DummyDetailBody();

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: BlocBuilder<ProcedureBloc, ProcedureState>(
        builder: (context, state) {
          if (state.status == ProcedureStatus.loading) {
            return const Center(child: CircularProgressIndicator());
          } else if (state.status == ProcedureStatus.failure) {
            return Center(
              child: Text(state.errorMessage ?? 'An error occurred'),
            );
          } else if (state.status == ProcedureStatus.success &&
              state.selectedProcedure == null) {
            return const Center(child: Text('No procedure selected'));
          } else if (state.status == ProcedureStatus.success &&
              state.selectedProcedure != null) {
            final procedure = state.selectedProcedure!;
            final feedback = state.feedbacks;
            return Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                ProcedureDetailHeader(
                  duration:
                      "${procedure.duration.minday} - ${procedure.duration.maxday} days",
                  cost: procedure.cost,
                ),
                const SizedBox(height: 20),
                DocumentList(documents: procedure.requiredDocuments),
                const SizedBox(height: 20),
                StepList(steps: procedure.steps),
                const SizedBox(height: 20),
                if (procedure.resources.isNotEmpty) DocmnetForm(),
                const SizedBox(height: 20),
                Center(
                  child: CustomButton(
                    text: 'Download Application Form',
                    icon: Icons.download,
                    onTap: () {
                      // Handle download action
                    },
                  ),
                ),

                const SizedBox(height: 20),

                BlocListener<ProcedureBloc, ProcedureState>(
  listener: (context, state) {
    if (state.errorMessage != null) {
      showDialog(
        context: context,
        builder: (context) => AlertDialog(
          title: const Text('Error'),
          content: Text(state.errorMessage!),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('OK'),
            ),
          ],
        ),
      );
      
    } else  {
      showDialog(
        context: context,
        builder: (context) => AlertDialog(
          title: const Text('Success'),
          content: const Text('Procedure saved successfully!'),
          actions: [
            TextButton(
              onPressed: () => context.push(RouteNames.workspace),
              child: const Text('Start Working'),
            ),
          ],
        ),
      );
      
    }
  },
  child: CustomButton(
    text: 'Save Procedure',
                    icon: Icons.download,
                    onTap: () {
                      context.read<ProcedureBloc>().add(
                        SaveProcedureEvent(procedure.id),
                      );
                    }
  ),
),


                

                DefaultTabController(
                  length: 3,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      const TabBar(
                        labelColor: AppColors.darkGreenColor,
                        unselectedLabelColor: Colors.grey,
                        indicatorColor: AppColors.darkGreenColor,
                        tabs: [
                          Tab(text: "Feedback"),
                          Tab(text: "Discussions"),
                        ],
                      ),
                      // ✅ Remove Expanded (can't use inside scroll)
                      SizedBox(
                        height:
                            300, // give a fixed height (or MediaQuery height fraction)
                        child: TabBarView(
                          children: [
                            FeedbackTab(feedbacklist: feedback ?? []),
                            const DiscussionTab(),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            );
          }
          return const Center(child: Text('Unknown state'));
        },
      ),
    );
  }
}
