import 'package:ethioguide/core/config/app_color.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure.dart';
import 'package:flutter/material.dart';
import '../../domain/entities/procedure_detail.dart';

/// Widget that displays the progress overview of a procedure
class ProgressOverviewCard extends StatelessWidget {
  final ProcedureDetail procedureDetail;

  const ProgressOverviewCard({
    super.key,
    required this.procedureDetail,
  });

  @override
  Widget build(BuildContext context) {
    // final completedSteps = procedureDetail..steps.where((step) => step.isCompleted).length;
    // final totalSteps = procedureDetail.steps.length;

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
            Text(
              'Progress Overview',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.w600,
              ),
            ),
            const SizedBox(height: 16),
            
            // Progress status
            Row(
              children: [
                /* Text(
                  '$completedSteps of $totalSteps steps completed',
                  style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                    fontWeight: FontWeight.w500,
                  ),
                ), */
                const Spacer(),
                Text(
                  '${procedureDetail.progressPercentage}%',
                  style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                    fontWeight: FontWeight.w600,
                    color: Colors.blue,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),
            
            // Progress bar
            LinearProgressIndicator(
              value: procedureDetail.progressPercentage / 100,
              backgroundColor: Colors.grey[300],
              valueColor: const AlwaysStoppedAnimation<Color>(AppColors.darkGreenColor),
              minHeight: 8,
            ),
            const SizedBox(height: 16),
            
            // Info grid
           /*  Row(
              children: [
                Expanded(
                  child: _InfoItem(
                    icon: Icons.access_time,
                    label: 'Est. Time',
                    value: procedureDetail.estimatedTime,
                    color: Colors.orange,
                  ),
                ),
                Expanded(
                  child: _InfoItem(
                    icon: Icons.people,
                    label: 'Difficulty',
                    value: procedureDetail.difficulty,
                    color: Colors.green,
                  ),
                ),
                Expanded(
                  child: _InfoItem(
                    icon: Icons.location_on,
                    label: 'Office',
                    value: procedureDetail.officeType,
                    color: Colors.blue,
                  ),
                ),
              ],
            ), */
          ],
        ),
      ),
    );
  }
}

class _InfoItem extends StatelessWidget {
  final IconData icon;
  final String label;
  final String value;
  final Color color;

  const _InfoItem({
    required this.icon,
    required this.label,
    required this.value,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Icon(
          icon,
          color: color,
          size: 24,
        ),
        const SizedBox(height: 4),
        Text(
          label,
          style: Theme.of(context).textTheme.bodySmall?.copyWith(
            color: Colors.grey[600],
          ),
          textAlign: TextAlign.center,
        ),
        const SizedBox(height: 2),
        Text(
          value,
          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
            fontWeight: FontWeight.w500,
          ),
          textAlign: TextAlign.center,
        ),
      ],
    );
  }
}
