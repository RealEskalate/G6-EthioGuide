import 'package:ethioguide/core/config/app_color.dart';
import 'package:flutter/material.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure.dart';

class ProcedureCard extends StatelessWidget {
  final Procedure procedure;
  final VoidCallback? onTap;
  final bool gridVariant; // true for 2x2 quick access style
  const ProcedureCard({super.key, required this.procedure, this.onTap, this.gridVariant = false});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return InkWell(
      onTap: onTap,
      child: gridVariant ? _buildGridCard(theme) : _buildListCard(theme),
    );
  }

  Widget _buildGridCard(ThemeData theme) {
    return Container(
      height: 10,
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(color: Colors.black.withOpacity(0.05), blurRadius: 12, offset: const Offset(0, 6)),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Gradient header
          Container(
            height: 64,
            decoration: const BoxDecoration(
              gradient: AppColors.gradient,
              borderRadius: BorderRadius.only(topLeft: Radius.circular(16), topRight: Radius.circular(16)),
            ),
            alignment: Alignment.centerLeft,
            padding: const EdgeInsets.all(16),
            child: 
          /*   Image.asset(
              procedure.icon,
              fit: BoxFit.fill,
            ), */
            const Icon(Icons.badge, color: Colors.white),
          ),
          Padding(
            padding: const EdgeInsets.all(12),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(procedure.title, style: theme.textTheme.bodyLarge?.copyWith(fontWeight: FontWeight.w600)),
                const SizedBox(height: 6),
                // Text(procedure.category, style: theme.textTheme.labelLarge),
                // const SizedBox(height: 8),
                Column(
                  children: [
                    Row(
                      children: [
                        const Icon(Icons.access_time, size: 14, color: Colors.black45),
                        const SizedBox(width: 4),
                    
                        Text("${procedure.duration.minday} - ${procedure.duration.maxday} days", style: theme.textTheme.labelLarge),
                       
                      ],
                    ),
                        const SizedBox(height: 4),
                    Row(
                      children: [
                
                        const Icon(Icons.payments_outlined, size: 14, color: Colors.black45),
                        const SizedBox(width: 4),
                        Text(procedure.cost, style: theme.textTheme.labelLarge),
                      ],
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildListCard(ThemeData theme) {
    return Container(
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(color: Colors.black.withOpacity(0.05), blurRadius: 12, offset: const Offset(0, 6)),
        ],
      ),
      padding: const EdgeInsets.all(16),
      child: Row(
        children: [
          CircleAvatar(
            radius: 24,
            backgroundColor: theme.colorScheme.primary.withOpacity(0.1),
            child: Icon(Icons.badge, color: theme.colorScheme.primary),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Expanded(
                      child: Text(
                        procedure.title,
                        style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w600),
                      ),
                    ),
                    const Icon(Icons.arrow_forward_ios, size: 16),
                  ],
                ),
                const SizedBox(height: 6),
                Wrap(
                  spacing: 8,
                  runSpacing: -8,
                  children: [
                    // _Chip(label: procedure.category),
                    _Chip(label: "${procedure.duration.minday} - ${procedure.duration.maxday}", icon: Icons.access_time),
                    _Chip(label: procedure.cost, icon: Icons.payments_outlined),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _Chip extends StatelessWidget {
  final String label;
  final IconData? icon;
  const _Chip({required this.label, this.icon});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
      decoration: BoxDecoration(
        color: theme.colorScheme.surfaceVariant,
        borderRadius: BorderRadius.circular(20),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          if (icon != null) ...[
            Icon(icon, size: 14, color: theme.colorScheme.onSurfaceVariant),
            const SizedBox(width: 4),
          ],
          Text(label, style: theme.textTheme.labelSmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
        ],
      ),
    );
  }
}

class _StatusPill extends StatelessWidget {
  final String text;
  const _StatusPill({required this.text});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    Color bg = AppColors.blueTagColor.withOpacity(0.1);
    Color fg = AppColors.blueTagColor;
    if (text.toLowerCase() == 'pending') {
      bg = AppColors.yellowTagColor.withOpacity(0.1);
      fg = AppColors.yellowTagColor;
    } else if (text.toLowerCase() == 'completed') {
      bg = AppColors.greenTagColor.withOpacity(0.1);
      fg = AppColors.greenTagColor;
    }
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
      decoration: BoxDecoration(color: bg, borderRadius: BorderRadius.circular(20)),
      child: Text(text, style: theme.textTheme.labelSmall?.copyWith(color: fg, fontWeight: FontWeight.w600)),
    );
  }
}


