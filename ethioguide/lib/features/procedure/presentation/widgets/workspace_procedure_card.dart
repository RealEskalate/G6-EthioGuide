import 'package:ethioguide/features/procedure/domain/entities/procedure_detail.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_step.dart';
import 'package:flutter/material.dart';
import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';

/// Widget that displays a workspace procedure card
class WorkspaceProcedureCard extends StatelessWidget {
  final ProcedureDetail procedure;
  final VoidCallback onTap;

  const WorkspaceProcedureCard({
    super.key,
    required this.procedure,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 2,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
      ),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Header with icon, title, and status
            Row(
              children: [
                _getStatusIcon(),
                const SizedBox(width: 12),
                Expanded(
                  child: Text(
                    procedure.procedure.title,
                    style: Theme.of(context).textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
                // _getStatusChip(),
              ],
            ),
            
        
            
            // Organization
           
            
            const SizedBox(height: 16),
            
            // Progress section
            Row(
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Progress',
                        style: Theme.of(context).textTheme.bodySmall?.copyWith(
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                      const SizedBox(height: 4),
                      LinearProgressIndicator(
                        value: procedure.progressPercentage / 100,
                        backgroundColor: Colors.grey[300],
                        valueColor: AlwaysStoppedAnimation<Color>(
                          _getProgressColor(),
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(width: 16),
                Text(
                  '${procedure.progressPercentage}% Complete',
                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    color: _getProgressColor(),
                    fontWeight: FontWeight.w500,
                  ),
                ),
              ],
            ),
            
            const SizedBox(height: 16),
            
            // Details
           /*  _DetailRow(
              icon: Icons.calendar_today,
              text: 'Started: ${_formatDate(procedure.startDate)}',
            ),
             */
           /*  if (procedure.estimatedCompletion != null) ...[
              const SizedBox(height: 8),
              _DetailRow(
                icon: Icons.track_changes,
                text: 'Est. completion: ${_formatDate(procedure.estimatedCompletion!)}',
              ),
            ],
            
            if (procedure.status == ProcedureStatus.completed && procedure.completedDate != null) ...[
              const SizedBox(height: 8),
              _DetailRow(
                icon: Icons.check_circle,
                text: 'Completed: ${_formatDate(procedure.completedDate!)}',
              ),
            ], */
            
            // const SizedBox(height: 8),
          /*   _DetailRow(
              icon: Icons.description,
              text: '${procedure.documentsUploaded}/${procedure.totalDocuments} documents uploaded',
            ), */
            
            const SizedBox(height: 16),
            
            // Action button
            Align(
              alignment: Alignment.centerRight,
              child: ElevatedButton(
                onPressed: onTap,
                style: ElevatedButton.styleFrom(
                  backgroundColor: Colors.teal,
                  foregroundColor: Colors.white,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(8),
                  ),
                ),
                child: const Text('View Checklist'),
              ),
            ),
          ],
        ),
      ),
    );
  }

  /* Widget _getStatusIcon() {
    IconData iconData;
    Color iconColor;
    
    switch (procedure.status) {
      case ProcedureStatus.completed:
        iconData = Icons.check_circle;
        iconColor = Colors.green;
        break;
      case ProcedureStatus.inProgress:
        iconData = Icons.description;
        iconColor = Colors.blue;
        break;
      case ProcedureStatus.notStarted:
        iconData = Icons.people;
        iconColor = Colors.purple;
        break;
    }
    
    return Container(
      padding: const EdgeInsets.all(8),
      decoration: BoxDecoration(
        color: iconColor.withOpacity(0.1),
        shape: BoxShape.circle,
      ),
      child: Icon(iconData, color: iconColor, size: 24),
    );
  }

  Widget _getStatusChip() {
    Color chipColor;
    Color textColor;
    
    switch (procedure.status) {
      case ProcedureStatus.completed:
        chipColor = Colors.green.shade100;
        textColor = Colors.green.shade800;
        break;
      case ProcedureStatus.inProgress:
        chipColor = Colors.orange.shade100;
        textColor = Colors.orange.shade800;
        break;
      case ProcedureStatus.notStarted:
        chipColor = Colors.grey.shade100;
        textColor = Colors.grey.shade800;
        break;
    }
    
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: chipColor,
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: textColor.withOpacity(0.3)),
      ),
      child: Text(
        procedure.status.displayName,
        style: TextStyle(
          color: textColor,
          fontSize: 12,
          fontWeight: FontWeight.w500,
        ),
      ),
    );
  }



class _DetailRow extends StatelessWidget {
  final IconData icon;
  final String text;

  const _DetailRow({
    required this.icon,
    required this.text,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Icon(
          icon,
          size: 16,
          color: Colors.grey[600],
        ),
        const SizedBox(width: 8),
        Text(
          text,
          style: Theme.of(context).textTheme.bodySmall?.copyWith(
            color: Colors.grey[600],
          ),
        ),
      ],
    );
  }
}
 */

  Color _getProgressColor() {
    if (procedure.progressPercentage == 100) return Colors.green;
    if (procedure.progressPercentage > 50) return Colors.orange;
    return Colors.grey;
  }

  String _formatDate(DateTime date) {
    return '${_getMonthName(date.month)} ${date.day}, ${date.year}';
  }

  String _getMonthName(int month) {
    const months = [
      'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
      'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'
    ];
    return months[month - 1];
  }


   Widget _getStatusIcon() {
   late IconData iconData;
  late  Color iconColor;
    
    switch (procedure.status) {
      case ProcedureStatus.completed:
        iconData = Icons.check_circle;
        iconColor = Colors.green;
        break;
      case ProcedureStatus.inProgress:
        iconData = Icons.description;
        iconColor = Colors.blue;
        break;
      case ProcedureStatus.notStarted:
        iconData = Icons.people;
        iconColor = Colors.purple;
        break;
    }
    
    return Container(
      padding: const EdgeInsets.all(8),
      decoration: BoxDecoration(
        color: iconColor.withOpacity(0.1),
        shape: BoxShape.circle,
      ),
      child: Icon(iconData, color: iconColor, size: 24),
    );
  }

  /* 
  } */

}