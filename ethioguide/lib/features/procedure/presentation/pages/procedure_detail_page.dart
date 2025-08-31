import 'package:flutter/material.dart';
import '../../domain/entities/procedure.dart';
import '../widgets/procedure_detail_header.dart';
import '../widgets/document_list.dart';
import '../widgets/step_list.dart';
import '../widgets/feedback_list.dart';

class ProcedureDetailPage extends StatelessWidget {
  

  const ProcedureDetailPage({super.key,});

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
      requiredDocuments: const ['Passport Photo', 'Old Passport', 'National ID'],
      steps: const [
        ProcedureStep(number: 1, title: 'Fill application', description: 'Complete the online form with accurate details.'),
        ProcedureStep(number: 2, title: 'Upload documents', description: 'Attach all required documents.'),
        ProcedureStep(number: 3, title: 'Payment', description: 'Pay fees via the supported channels.'),
      ],
      resources: const [Resource(name: 'Application Form', url: 'https://example.com/form.pdf')],
      feedback: const [
        FeedbackItem(user: 'Hanna', comment: 'Quick and easy!', date: '2025-01-10', verified: true),
        FeedbackItem(user: 'Samuel', comment: 'Got it in two weeks.', date: '2025-02-05', verified: false),
      ],
    );

    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          ProcedureDetailHeader(duration: procedure.duration, cost: procedure.cost),
          const SizedBox(height: 20),
          DocumentList(documents: procedure.requiredDocuments),
          const SizedBox(height: 20),
          StepList(steps: procedure.steps),
          const SizedBox(height: 20),
          if (procedure.resources.isNotEmpty)
            ElevatedButton.icon(
              onPressed: () {},
              icon: const Icon(Icons.download),
              label: const Text('Download Application Form'),
            ),
          const SizedBox(height: 20),
          FeedbackList(feedback: procedure.feedback),
          const SizedBox(height: 20),
          ElevatedButton(onPressed: () {}, child: const Text('Give Feedback')),
        ],
      ),
    );
  }
}


