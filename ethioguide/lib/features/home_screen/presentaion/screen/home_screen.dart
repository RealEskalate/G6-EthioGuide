import 'package:ethioguide/core/components/bottom_nav_bar.dart';
import 'package:ethioguide/features/home_screen/domain/entities/home_data.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:get_it/get_it.dart';
import 'package:go_router/go_router.dart';
import 'package:ethioguide/core/components/dashboard_card.dart';
import 'package:ethioguide/core/config/app_color.dart';
import 'package:ethioguide/features/home_screen/domain/entities/home_data.dart';
import '../bloc/home_bloc.dart';
import '../bloc/home_event.dart';
import '../bloc/home_state.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (context) => GetIt.instance<HomeBloc>()..add(LoadHomeData()),
      child: const HomeView(),
    );
  }
}

class HomeView extends StatelessWidget {
  const HomeView({super.key});
  final pageIndex = 0;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.grey[100],
      body: BlocBuilder<HomeBloc, HomeState>(
        builder: (context, state) {
          return CustomScrollView(
            slivers: [
              // --- APP BAR ---
              SliverAppBar(
                backgroundColor: Colors.grey[100],
                pinned: true,
                elevation: 0,
                leading: Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: CircleAvatar(
                    backgroundColor: Colors.white,
                    child: Icon(
                      Icons.home_outlined,
                      color: AppColors.darkGreenColor,
                    ),
                  ),
                ),
                title: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Image.asset('assets/images/light_logo.jpg', height: 20),
                    const Text(
                      "Welcome back!",
                      style: TextStyle(fontSize: 14, color: Colors.grey),
                    ),
                  ],
                ),
                actions: [
                  Stack(
                    alignment: Alignment.topRight,
                    children: [
                      IconButton(
                        onPressed: () {},
                        icon: const Icon(
                          Icons.notifications_outlined,
                          color: Colors.black54,
                        ),
                      ),
                      Container(
                        margin: const EdgeInsets.only(top: 8, right: 8),
                        padding: const EdgeInsets.all(4),
                        decoration: const BoxDecoration(
                          color: Colors.red,
                          shape: BoxShape.circle,
                        ),
                        child: const Text(
                          '5',
                          style: TextStyle(color: Colors.white, fontSize: 10),
                        ),
                      ),
                    ],
                  ),
                  IconButton(
                    onPressed: () {
                      // Use go_router to navigate to the profile page
                      context.pushNamed('profile');
                    },
                    icon: const Icon(
                      Icons.person_outline,
                      color: Colors.black54,
                    ),
                  ),
                ],
                bottom: PreferredSize(
                  preferredSize: const Size.fromHeight(60.0),
                  child: Padding(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 16.0,
                      vertical: 8.0,
                    ),
                    child: TextField(
                      decoration: InputDecoration(
                        hintText: 'Search for government services...',
                        prefixIcon: const Icon(Icons.search),
                        filled: true,
                        fillColor: Colors.white,
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(12),
                          borderSide: BorderSide.none,
                        ),
                      ),
                    ),
                  ),
                ),
              ),

              // --- MAIN CONTENT LIST ---
              SliverList(
                delegate: SliverChildListDelegate([
                  _SectionHeader(title: "Quick Actions"),
                  _QuickActionsGrid(actions: state.quickActions),
                  ...state.contentCards
                      .map((card) => _ContentCardWidget(card: card))
                      .toList(),
                  _SectionHeader(title: "Popular Services"),
                  ...state.popularServices
                      .map((service) => _PopularServiceCard(service: service))
                      .toList(),
                  const SizedBox(height: 40),
                ]),
              ),
            ],
          );
        },
      ),

      bottomNavigationBar: bottomNav(context: context, selectedIndex: pageIndex),
    );
  }
}

// --- HELPER WIDGETS ---
class _SectionHeader extends StatelessWidget {
  final String title;
  const _SectionHeader({required this.title});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 24, 16, 8),
      child: Text(
        title,
        style: const TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
      ),
    );
  }
}

class _QuickActionsGrid extends StatelessWidget {
  final List<QuickAction> actions;
  const _QuickActionsGrid({required this.actions});

  @override
  Widget build(BuildContext context) {
    return GridView.builder(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      itemCount: actions.length,
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        crossAxisSpacing: 16,
        mainAxisSpacing: 16,
        childAspectRatio: 0.9,
      ),
      itemBuilder: (context, index) {
        final action = actions[index];
        return DashboardCard(
          icon: action.icon,
          title: action.title,
          subtitle: action.subtitle,
          onTap: () => context.pushNamed(action.routeName),
        );
      },
    );
  }
}

class _ContentCardWidget extends StatelessWidget {
  final ContentCard card;
  const _ContentCardWidget({required this.card});

  @override
  Widget build(BuildContext context) {
    // This is a simplified card. You can create a more detailed custom widget for this later.
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        _SectionHeader(title: card.sectionTitle),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16.0),
          child: GestureDetector(
            onTap: () => context.pushNamed(card.routeName),
            child: Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(12),
              ),
              child: Row(
                children: [
                  Icon(card.icon, size: 32, color: AppColors.darkGreenColor),
                  const SizedBox(width: 16),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          card.title,
                          style: const TextStyle(
                            fontWeight: FontWeight.bold,
                            fontSize: 16,
                          ),
                        ),
                        Text(
                          card.subtitle,
                          style: const TextStyle(color: Colors.grey),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ],
    );
  }
}

class _PopularServiceCard extends StatelessWidget {
  final PopularService service;
  const _PopularServiceCard({required this.service});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 8.0),
      child: GestureDetector(
        onTap: () => context.pushNamed(service.routeName),
        child: Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.circular(12),
          ),
          child: Row(
            children: [
              CircleAvatar(child: Icon(service.icon)),
              const SizedBox(width: 16),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      service.title,
                      style: const TextStyle(
                        fontWeight: FontWeight.bold,
                        fontSize: 16,
                      ),
                    ),
                    Text(
                      service.category,
                      style: const TextStyle(color: Colors.grey),
                    ),
                    const SizedBox(height: 8),
                    Row(
                      children: [
                        const Icon(
                          Icons.timer_outlined,
                          size: 16,
                          color: Colors.grey,
                        ),
                        const SizedBox(width: 4),
                        Text(
                          service.timeEstimate,
                          style: const TextStyle(color: Colors.grey),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
              const Icon(Icons.bookmark_border_outlined),
            ],
          ),
        ),
      ),
    );
  }
}
