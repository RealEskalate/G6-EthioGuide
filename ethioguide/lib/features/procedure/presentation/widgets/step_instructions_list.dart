import 'package:ethioguide/core/config/app_color.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_step.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../domain/entities/procedure_detail.dart';
import '../bloc/workspace_procedure_detail_bloc.dart';

/// Widget that displays the step-by-step instructions
///
class StepInstructionsList extends StatefulWidget {
  final List<MyProcedureStep> procedureDetail;
  const StepInstructionsList({super.key, required this.procedureDetail});

  @override
  State<StepInstructionsList> createState() => _StepInstructionsList();
}

class _StepInstructionsList extends State<StepInstructionsList> {
  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 2,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(Icons.book, color: Colors.blue[600], size: 24),
                const SizedBox(width: 8),
                Text(
                  'Step-by-Step Instructions',
                  style: Theme.of(context).textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),

            // Steps list
            ...widget.procedureDetail.map(
              (step) => _StepItem(step: step, procedureId: step.id),
            ),
          ],
        ),
      ),
    );
  }
}

class _StepItem extends StatefulWidget {
  final MyProcedureStep step;
  final String procedureId;

  const _StepItem({Key? key, required this.step, required this.procedureId})
    : super(key: key);

  @override
  State<_StepItem> createState() => _StepItemState();
}

class _StepItemState extends State<_StepItem> {
  late bool isChecked;

  @override
  void initState() {
    super.initState();
    isChecked = widget.step.isChecked;
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16.0),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Step number or completion icon
          Container(
            width: 32,
            height: 32,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              color: isChecked ? Colors.green : AppColors.darkGreenColor,
            ),
            child: Center(
              child: isChecked
                  ? const Icon(Icons.check, color: Colors.white, size: 20)
                  : const Icon(
                      Icons.check_box_outline_blank,
                      color: Colors.white,
                      size: 20,
                    ),
            ),
          ),
          const SizedBox(width: 12),

          // Step content
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  widget.step.title,
                  style: Theme.of(
                    context,
                  ).textTheme.bodyLarge?.copyWith(fontWeight: FontWeight.w600),
                ),

                const SizedBox(height: 8),
              ],
            ),
          ),

          // Toggle button
          Checkbox(
            value: isChecked,
            onChanged: (value) {
              if (value != null) {
                setState(() {
                  isChecked = value;
                });
              }
            },
            activeColor: Colors.green,
          ),
        ],
      ),
    );
  }
}
