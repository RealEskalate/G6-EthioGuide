import 'package:flutter/material.dart';

class ProcedureDetailHeader extends StatelessWidget {
  final String duration;
  final String cost;

  const ProcedureDetailHeader({super.key, required this.duration, required this.cost});

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          InfoCard(icon: Icons.timer, title: 'Processing Time', value: duration),
          InfoCard(icon: Icons.attach_money, title: 'Total Fees', value: cost),
          InfoCard(icon: Icons.chat, title: 'AI Assistant', value: 'Chat here'),
        ],
      ),
    );
  }
}

class InfoCard extends StatelessWidget {
  final IconData icon;
  final String title;
  final String value;

  const InfoCard({required this.icon, required this.title, required this.value});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(15.0),
      child: Expanded(
        child: Container(
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
              children: [
                Icon(icon, size: 28, color: Theme.of(context).colorScheme.primary),
                const SizedBox(height: 6),
                Text(title, style: Theme.of(context).textTheme.bodySmall),
                Text(value, style: Theme.of(context).textTheme.titleMedium),
              ],
            ),
          ),
        ),
      ),
    );
  }
}


