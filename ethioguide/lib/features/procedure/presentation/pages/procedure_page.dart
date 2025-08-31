import 'package:ethioguide/core/config/app_color.dart';
import 'package:ethioguide/core/config/app_theme.dart';
import 'package:ethioguide/features/procedure/presentation/widgets/procedure_card.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/config/route_names.dart';
import '../../domain/entities/procedure.dart';


/// Page that displays the list of procedures, styled similarly to the mock.
class ProcedurePage extends StatelessWidget {
  const ProcedurePage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar:
     AppBar(
  backgroundColor: AppTheme.lightTheme.appBarTheme.backgroundColor,
  elevation: 0,
   toolbarHeight: 90, // 🟢 increase height
  leading: IconButton(
    icon: const Icon(Icons.arrow_back),
    onPressed: () => Navigator.maybePop(context),
  ),
  title: Expanded(
    child: Padding(
      padding: const EdgeInsets.symmetric(horizontal: 4.0 , vertical: 20 ), // padding added
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          
          Text(
            'Welcome to EthioGuide!',
            style: Theme.of(context).textTheme.headlineSmall,
            softWrap: true,
            overflow: TextOverflow.visible, // prevent ellipsis
          ),
          const SizedBox(height: 6),
          Text(
            'Your trusted partner for navigating Ethiopian government procedures with ease.',
            style: Theme.of(context)
                .textTheme
                .bodyMedium
                ?.copyWith(color: AppColors.graycolor),
            softWrap: true,
            overflow: TextOverflow.visible, // allow wrapping
          ),
        ],
      ),
    ),
  ),
),

      body: SingleChildScrollView(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const SizedBox(height: 12),
            _SearchBar(),
            const SizedBox(height: 16),

            // Quick Access Section
            Text('Quick Access Procedures', style: Theme.of(context).textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w600)),
            const SizedBox(height: 12),
            _QuickAccessGrid(),
            const SizedBox(height: 16),
            // All Procedures Section
            Text('All Procedures', style: Theme.of(context).textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w600)),
            const SizedBox(height: 12),
            _DummyAllProceduresList(),
          ],
        ),
      ),
    );
  }
}

class _SearchBar extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return TextField(
      decoration: InputDecoration(
        
        hintText: 'Search government services...',
        prefixIcon: const Icon(Icons.search),
        filled: true,
        fillColor: AppColors.graycolor.withOpacity(0.2),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(30), borderSide: BorderSide.none),
        contentPadding: const EdgeInsets.symmetric(vertical: 0, horizontal: 16),
      ),
    );
  }
}

class _QuickAccessGrid extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    // Static 4 items to match the mock while backend loads the complete list

    final List<Procedure> items = const [
      Procedure(
        id: '1',
        title: 'Passport Renewal',
        category: 'Travel',
        duration: '2-3 weeks',
        cost: '1,200 ETB',
        icon: 'badge',
        isQuickAccess: false,
      ),
      Procedure(
        id: '2',
        title: 'Vehicle Registration',
        category: 'Transportation',
        duration: '1-2 days',
        cost: '800 ETB',
        icon: 'badge',
        isQuickAccess: false,
      ),
       Procedure(
        id: '3',
        title:  'Driving License',
        category: 'Transportation',
        duration: '1-2 days',
        cost: '800 ETB',
        icon: 'badge',
        isQuickAccess: false,
      ),
       Procedure(
        id: '4',
        title: "Driver's License",
        category: 'Transportation',
        duration: '1-2 days',
        cost: '800 ETB',
        icon: 'badge',
        isQuickAccess: false,
      ),
    ];

    return GridView.builder(
      itemCount: items.length,
      physics: const NeverScrollableScrollPhysics(),
      shrinkWrap: true,
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        crossAxisSpacing: 12,
        mainAxisSpacing: 12,
        childAspectRatio: 1.2,
      ),
      itemBuilder: (context, index) {
        final item = items[index];
        return ProcedureCard(
          gridVariant: true,
          procedure: item,
        );
      },
    );
  }

  // Helper to construct a Procedure quickl
}

class _DummyAllProceduresList extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    final List<Procedure> procedures = const [
      Procedure(
        id: '1',
        title: 'Passport Renewal',
        category: 'Travel',
        duration: '2-3 weeks',
        cost: '1,200 ETB',
        icon: 'badge',
        isQuickAccess: false,
      ),
      Procedure(
        id: '2',
        title: 'Vehicle Registration',
        category: 'Transportation',
        duration: '1-2 days',
        cost: '800 ETB',
        icon: 'badge',
        isQuickAccess: false,
      ),
    ];

    return ListView.separated(
      physics: const NeverScrollableScrollPhysics(),
      shrinkWrap: true,
      itemCount: procedures.length,
      separatorBuilder: (_, __) => const SizedBox(height: 12),
      itemBuilder: (context, index) {
        final p = procedures[index];
        return ProcedureCard(procedure: p, onTap: () {
          context.push(RouteNames.procedure_detail);
        });
      },
    );
  }
}

// Lightweight helper value-object only used internally for quick access grid
// class ProcedureCardData {
//   final String id;
//   final String title;
//   final String category;
//   final String duration;
//   final String cost;
//   const ProcedureCardData({required this.id, required this.title, required this.category, required this.duration, required this.cost});
//   // Convert to entity with default values for fields not displayed in grid
//   toEntity() => Procedure(
//         id: id,
//         title: title,
//         category: category,
//         duration: duration,
//         cost: cost,
//         icon: 'badge',
//         isQuickAccess: true,
//       );
// }


