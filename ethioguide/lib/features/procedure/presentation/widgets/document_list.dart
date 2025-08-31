import 'package:ethioguide/core/config/app_color.dart';
import 'package:flutter/material.dart';

class DocumentList extends StatelessWidget {
  final List<String> documents;

  const DocumentList({super.key, required this.documents});

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
            Text(
              'Required Documents',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 10),
            ...documents.map(
              (doc) => ListTile(
                leading: const Icon(Icons.description),
                title: Text(doc, style: Theme.of(context).textTheme.bodyMedium),
                trailing: Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 18,
                    vertical: 6,
                  ),
                  decoration: BoxDecoration(
                    color: AppColors.darkGreenColor,
                    borderRadius: BorderRadius.circular(15),
                  ),
                  child: const Text(
                    'required',
                    style: TextStyle(color: Colors.white, fontSize: 12),
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
