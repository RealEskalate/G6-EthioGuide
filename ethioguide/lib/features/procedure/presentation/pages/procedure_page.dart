import 'package:ethioguide/core/components/dashboard_card.dart';
import 'package:ethioguide/core/config/app_color.dart';
import 'package:ethioguide/core/config/app_theme.dart';
import 'package:ethioguide/features/procedure/presentation/bloc/procedure_bloc.dart';
import 'package:ethioguide/features/procedure/presentation/widgets/procedure_card.dart';
import 'package:ethioguide/features/workspace_discussion/presentation/bloc/workspace_discussion_bloc.dart';
import 'package:ethioguide/features/workspace_discussion/presentation/bloc/worspace_discustion_state.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/config/route_names.dart';
import '../../domain/entities/procedure.dart';

/// Page that displays the list of procedures, styled similarly to the mock.
class ProcedurePage extends StatelessWidget {
  const ProcedurePage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: AppTheme.lightTheme.appBarTheme.backgroundColor,
        elevation: 0,
        toolbarHeight: 90, 
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        title: Expanded(
          child: Padding(
            padding: const EdgeInsets.symmetric(
              horizontal: 4.0,
              vertical: 20,
            ), // padding added
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
                  style: Theme.of(
                    context,
                  ).textTheme.bodyMedium?.copyWith(color: AppColors.graycolor),
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
            Text(
              'Quick Access Procedures',
              style: Theme.of(
                context,
              ).textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w600),
            ),

            BlocBuilder<ProcedureBloc, ProcedureState>(
              builder: (context, state) {
                if (state.status == ProcedureStatus.loading) {
                  return const Center(child: CircularProgressIndicator());
                } else if (state.status == ProcedureStatus.failure) {
                  print(state.errorMessage);
                  return Center(
                    child: Text(
                      "Error: ${state.errorMessage ?? 'Unknown error'}",
                    ),
                  );
                } else if (state.status == ProcedureStatus.success) {
                  if (state.procedures.isEmpty) {
                    return const Center(child: Text("No procedures found."));
                  }
                  return Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const SizedBox(height: 12),
                      _QuickAccessGrid(items: state.procedures),
                      const SizedBox(height: 10),
                      // All Procedures Section
                      Text(
                        'All Procedures',
                        style: Theme.of(context).textTheme.titleMedium
                            ?.copyWith(fontWeight: FontWeight.w600),
                      ),
                      const SizedBox(height: 12),
                      _buildProceduresList(state.procedures),
                    ],
                  );
                }

                // Default: initial state
                return const Center(child: Text("Please load procedures."));
              },
            ),
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
      onSubmitted: (query) {
              print(query);
              context.read<ProcedureBloc>().add(
               LoadProceduresEvent(query)
              );
            },
      decoration: InputDecoration(
        hintText: 'Search government services...',
        prefixIcon: const Icon(Icons.search),
        filled: true,
        fillColor: AppColors.graycolor.withOpacity(0.2),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(30),
          borderSide: BorderSide.none,
        ),
        contentPadding: const EdgeInsets.symmetric(vertical: 0, horizontal: 16),
      ),
    );
  }
}

class _QuickAccessGrid extends StatelessWidget {
  final List<Procedure> items;

  const _QuickAccessGrid({Key? key, required this.items}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    // Show only the first 4 (or fewer if less than 4)
    final quickItems = items.take(4).toList();

    return SizedBox(
     // height: (quickItems.length / 2).ceil() * 210, // dynamic height
      child: GridView.builder(
        itemCount: quickItems.length,
        physics: const NeverScrollableScrollPhysics(),
        shrinkWrap: true,
        gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
          crossAxisCount: 2,
          crossAxisSpacing: 12,
          mainAxisSpacing: 12,
          childAspectRatio: 1.2,
          mainAxisExtent: 210, // 👈 fixed card height
        ),
        itemBuilder: (context, index) {
          final item = quickItems[index];
          return ProcedureCard(
            gridVariant: true,
            procedure: item,
            onTap: () {
              context.read<ProcedureBloc>().add(
                LoadProcedureByIdEvent(item.id),
              );

              context.push(RouteNames.procedure_detail);
            },
          );
        },
      ),
    );
  }
}

Widget _buildProceduresList(List<Procedure> procedures) {
  return ListView.separated(
    physics: const NeverScrollableScrollPhysics(),
    shrinkWrap: true,
    itemCount: procedures.length,
    separatorBuilder: (_, __) => const SizedBox(height: 12),
    itemBuilder: (context, index) {
      final p = procedures[index];
      return ProcedureCard(
        procedure: p,
        onTap: () {
          context.read<ProcedureBloc>().add(LoadProcedureByIdEvent(p.id));

          context.push(RouteNames.procedure_detail);
        },
      );
    },
  );
}
