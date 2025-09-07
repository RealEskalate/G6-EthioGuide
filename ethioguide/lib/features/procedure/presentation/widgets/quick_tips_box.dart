import 'package:flutter/material.dart';
import '../../domain/entities/procedure_detail.dart';

/// Widget that displays quick tips for the procedure
// class QuickTipsBox extends StatelessWidget {
//   final ProcedureDetail procedureDetail;

//   const QuickTipsBox({
//     super.key,
//     required this.procedureDetail,
//   });

//   @override
//   Widget build(BuildContext context) {
//     return Card(
//       elevation: 2,
//       shape: RoundedRectangleBorder(
//         borderRadius: BorderRadius.circular(12),
//       ),
//       color: Colors.blue[50],
//       child: Padding(
//         padding: const EdgeInsets.all(16.0),
//         child: Column(
//           crossAxisAlignment: CrossAxisAlignment.start,
//           children: [
//             Row(
//               children: [
//                 Icon(
//                   Icons.lightbulb,
//                   color: Colors.amber[600],
//                   size: 24,
//                 ),
//                 const SizedBox(width: 8),
//                 Text(
//                   'Quick Tips',
//                   style: Theme.of(context).textTheme.titleMedium?.copyWith(
//                     fontWeight: FontWeight.w600,
//                     color: Colors.blue[800],
//                   ),
//                 ),
//               ],
//             ),
//             const SizedBox(height: 12),
            
//             // Tips list
//             ...procedureDetail.quickTips.map((tip) => Padding(
//               padding: const EdgeInsets.only(bottom: 8.0),
//               child: Row(
//                 crossAxisAlignment: CrossAxisAlignment.start,
//                 children: [
//                   Container(
//                     width: 6,
//                     height: 6,
//                     margin: const EdgeInsets.only(top: 8),
//                     decoration: BoxDecoration(
//                       shape: BoxShape.circle,
//                       color: Colors.blue[600],
//                     ),
//                   ),
//                   const SizedBox(width: 12),
//                   Expanded(
//                     child: Text(
//                       tip,
//                       style: Theme.of(context).textTheme.bodyMedium?.copyWith(
//                         color: Colors.blue[700],
//                       ),
//                     ),
//                   ),
//                 ],
//               ),
//             )),
//           ],
//         ),
//       ),
//     );
//   }
// }
