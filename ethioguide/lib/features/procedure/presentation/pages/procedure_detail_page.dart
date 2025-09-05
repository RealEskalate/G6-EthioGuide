import 'package:ethioguide/core/components/button.dart';
import 'package:ethioguide/core/config/app_color.dart';
import 'package:ethioguide/features/procedure/presentation/pages/Discussin_tab.dart';
import 'package:ethioguide/features/procedure/presentation/pages/Feedback_tab.dart';
import 'package:ethioguide/features/procedure/presentation/pages/Notice_tab.dart';
import 'package:ethioguide/features/procedure/presentation/widgets/docment_form.dart';
import 'package:flutter/material.dart';
import '../../domain/entities/procedure.dart';
import '../widgets/procedure_detail_header.dart';
import '../widgets/document_list.dart';
import '../widgets/step_list.dart';
import '../widgets/feedback_list.dart';

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
    // Dummy entity filled with demonstration data
    final procedure = Procedure(
      id: '1',
      title: 'Passport Renewal',
      category: 'Travel',
      duration: '2-3 weeks',
      cost: '1,200 ETB',
      icon: 'badge',
      isQuickAccess: false,
      requiredDocuments: const [
        'Passport Photo',
        'Old Passport',
        'National ID',
      ],
      steps: const [
        ProcedureStep(
          number: 1,
          title: 'Fill application',
          description: 'Complete the online form with accurate details.',
        ),
        ProcedureStep(
          number: 2,
          title: 'Upload documents',
          description: 'Attach all required documents.',
        ),
        ProcedureStep(
          number: 3,
          title: 'Payment',
          description: 'Pay fees via the supported channels.',
        ),
      ],
      resources: const [
        Resource(name: 'Application Form', url: 'https://example.com/form.pdf'),
      ],
      feedback: const [
        FeedbackItem(
          user: 'Hanna',
          comment: 'Quick and easy!',
          date: '2025-01-10',
          verified: true,
        ),
        FeedbackItem(
          user: 'Samuel',
          comment: 'Got it in two weeks.',
          date: '2025-02-05',
          verified: false,
        ),
      ],
    );

    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          ProcedureDetailHeader(
            duration: procedure.duration,
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
                Tab(text: "Notices"),
                Tab(text: "Feedback"),
                Tab(text: "Discussions"),
              ],
            ),
            // ✅ Remove Expanded (can't use inside scroll)
            SizedBox(
              height: 300, // give a fixed height (or MediaQuery height fraction)
              child: const TabBarView(
                children: [
                  NoticesTab(),
                  FeedbackTab(),
                  DiscussionTab(),
                ],
              ),
            ),
          ],
        ),
      ),
         

         
        ],
      ),
    );
  }
}
