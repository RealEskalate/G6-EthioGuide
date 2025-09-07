import 'package:flutter/material.dart';
import '../../domain/entities/procedure.dart';

class StepList extends StatelessWidget {
  final List<ProcedureStep> steps;

  const StepList({super.key, required this.steps});

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
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Step-by-Step Instructions', style: Theme.of(context).textTheme.titleMedium),
            const SizedBox(height: 10),
            ...steps.map((s) => ListTile(
                  leading: CircleAvatar(child: Text(s.number.toString())),
                  title: Text(s.title, style: Theme.of(context).textTheme.bodyLarge?.copyWith(fontWeight: FontWeight.w600)),
                )),
          ],
        ),
      ),
    );
  }
}


