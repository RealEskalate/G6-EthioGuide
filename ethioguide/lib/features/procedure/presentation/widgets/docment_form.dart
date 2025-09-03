import 'package:ethioguide/core/config/app_color.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_step.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../domain/entities/procedure_detail.dart';
import '../bloc/workspace_procedure_detail_bloc.dart';

/// Widget that displays the step-by-step instructions
class DocmnetForm extends StatelessWidget {
  const DocmnetForm({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: Theme.of(context).cardColor,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          // central soft shadow
          BoxShadow(
            color: Colors.black.withOpacity(0.12),
            offset: const Offset(0, 4),
            blurRadius: 8,
            spreadRadius: 1,
          ),
          // left
          BoxShadow(
            color: Colors.black.withOpacity(0.06),
            offset: const Offset(-4, 0),
            blurRadius: 6,
            spreadRadius: 0,
          ),
          // right
          BoxShadow(
            color: Colors.black.withOpacity(0.06),
            offset: const Offset(4, 0),
            blurRadius: 6,
            spreadRadius: 0,
          ),
          // top
          BoxShadow(
            color: Colors.black.withOpacity(0.06),
            offset: const Offset(0, -3),
            blurRadius: 6,
            spreadRadius: 0,
          ),
        ],
      ),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Resources & Forms',
              style: Theme.of(
                context,
              ).textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w600),
            ),

            const SizedBox(height: 16),

            Center(
              child: ElevatedButton.icon(
                onPressed: () {},
                icon: const Icon(Icons.download),
                label: const Text('Download Application Form'),
              ),
            ),

            const SizedBox(height: 16),

            Text(
              'Official application forms and additional resources will be available here.',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.w200,
                fontSize: 12,
              ),
            ),

            // Steps list
            // ...procedureDetail.steps.map((step) => _StepItem(
            //   step: step,
            //   procedureId: procedureDetail.id,
            // )),
          ],
        ),
      ),
    );
  }
}
