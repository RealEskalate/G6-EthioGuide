import 'package:flutter/material.dart';
import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';

/// Widget that displays summary statistics in card format
class WorkspaceSummaryCards extends StatelessWidget {
  final WorkspaceSummary summary;

  const WorkspaceSummaryCards({super.key, required this.summary});

  @override
  Widget build(BuildContext context) {
    return GridView.count(
      crossAxisCount: 2,
      crossAxisSpacing: 12,
      mainAxisSpacing: 12,
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      childAspectRatio: 1.5,
      children: [
        SummaryCard(
          title: 'Total Procedures',
          value: summary.totalProcedures.toString(),
          icon: Icons.description,
          color: Colors.blue,
        ),
        SummaryCard(
          title: 'In Progress',
          value: summary.inProgress.toString(),
          icon: Icons.schedule,
          color: Colors.orange,
        ),
        SummaryCard(
          title: 'Completed',
          value: summary.completed.toString(),
          icon: Icons.check_circle,
          color: Colors.green,
        ),
        SummaryCard(
          title: 'Documents',
          value: summary.totalDocuments.toString(),
          icon: Icons.folder,
          color: Colors.purple,
        ),
      ],
    );
  }
}

class SummaryCard extends StatelessWidget {
  final String title;
  final String value;
  final IconData icon;
  final Color color;

  const SummaryCard({
    required this.title,
    required this.value,
    required this.icon,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 2,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              title,
              style: Theme.of(
                context,
              ).textTheme.bodySmall?.copyWith(color: Colors.grey[600]),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 4),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  value,
                  style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: color,
                  ),
                ),
                const SizedBox(height: 8),
                Container(
                  padding: const EdgeInsets.all(8),
                  decoration: BoxDecoration(
                    color: color.withOpacity(0.1),
                    shape: BoxShape.circle,
                  ),
                  child: Icon(icon, color: color, size: 24),
                ),
                /*  Icon(
                  icon,
                  size: 32,
                  color: color,
                ), */
              ],
            ),
          ],
        ),
      ),
    );
  }
}
