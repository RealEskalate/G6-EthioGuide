import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../domain/entities/procedure_detail.dart';
import '../bloc/workspace_procedure_detail_bloc.dart';

/// Widget that displays the step-by-step instructions
class StepInstructionsList extends StatelessWidget {
  final ProcedureDetail procedureDetail;

  const StepInstructionsList({
    super.key,
    required this.procedureDetail,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 2,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
      ),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(
                  Icons.book,
                  color: Colors.blue[600],
                  size: 24,
                ),
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
            ...procedureDetail.steps.map((step) => _StepItem(
              step: step,
              procedureId: procedureDetail.id,
            )),
          ],
        ),
      ),
    );
  }
}

class _StepItem extends StatelessWidget {
  final ProcedureStep step;
  final String procedureId;

  const _StepItem({
    required this.step,
    required this.procedureId,
  });

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
              color: step.isCompleted ? Colors.green : Colors.grey[300],
            ),
            child: Center(
              child: step.isCompleted
                  ? const Icon(
                      Icons.check,
                      color: Colors.white,
                      size: 20,
                    )
                  : Text(
                      '${step.order}',
                      style: TextStyle(
                        color: Colors.grey[700],
                        fontWeight: FontWeight.w600,
                      ),
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
                  step.title,
                  style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  step.description,
                  style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                    color: Colors.grey[600],
                  ),
                ),
                const SizedBox(height: 8),
                
                // Completion status
                if (step.completionStatus != null)
                  Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 8,
                      vertical: 4,
                    ),
                    decoration: BoxDecoration(
                      color: step.isCompleted ? Colors.blue[100] : Colors.grey[200],
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Icon(
                          step.isCompleted ? Icons.check_circle : Icons.radio_button_unchecked,
                          size: 16,
                          color: step.isCompleted ? Colors.blue[600] : Colors.grey[600],
                        ),
                        const SizedBox(width: 4),
                        Text(
                          step.completionStatus!,
                          style: TextStyle(
                            fontSize: 12,
                            color: step.isCompleted ? Colors.blue[700] : Colors.grey[700],
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                      ],
                    ),
                  ),
              ],
            ),
          ),
          
          // Toggle button
          Checkbox(
            value: step.isCompleted,
            onChanged: (value) {
              if (value != null) {
                context.read<WorkspaceProcedureDetailBloc>().add(
                  UpdateStepStatus(
                    procedureId: procedureId,
                    stepId: step.id,
                    isCompleted: value,
                  ),
                );
              }
            },
            activeColor: Colors.green,
          ),
        ],
      ),
    );
  }
}
