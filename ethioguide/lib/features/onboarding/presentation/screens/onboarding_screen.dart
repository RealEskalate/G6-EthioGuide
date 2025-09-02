import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:ethioguide/core/components/button.dart';
import 'package:ethioguide/core/components/icon_container.dart';
import 'package:ethioguide/core/config/app_color.dart';
import '../../data/repositories/onboarding_repository_impl.dart';
import '../../domain/usecases/get_onboarding_pages.dart';
import '../bloc/onboarding_bloc.dart';
import '../bloc/onboarding_event.dart';
import '../bloc/onboarding_state.dart';

class OnboardingScreen extends StatelessWidget {
  const OnboardingScreen({super.key});

  @override
  Widget build(BuildContext context) {
    // 1. Provide the OnboardingBloc to the widget tree.
    return BlocProvider(
      create: (context) => OnboardingBloc(
        getOnboardingPages: GetOnboardingPages(OnboardingRepositoryImpl()),
      )..add(LoadPages()), // 2. Add the initial event to load the static data.
      child: const OnboardingView(),
    );
  }
}

class OnboardingView extends StatefulWidget {
  const OnboardingView({super.key});

  @override
  State<OnboardingView> createState() => _OnboardingViewState();
}

class _OnboardingViewState extends State<OnboardingView> {
  // 3. Create a PageController to programmatically control the PageView.
  final PageController _pageController = PageController();

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white, // Assuming white background
      // 4. Use BlocConsumer to both rebuild the UI and listen for actions.
      body: BlocConsumer<OnboardingBloc, OnboardingState>(
        listener: (context, state) {
          // Listener is for one-time actions, not needed here but good practice.
        },
        builder: (context, state) {
          // If there are no pages yet, show a loading indicator.
          if (state.pages.isEmpty) {
            return const Center(child: CircularProgressIndicator());
          }

          // --- Show the main content once pages are loaded ---
          return SafeArea(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
              child: Column(
                children: [
                  // --- "Skip" Button ---
                  Align(
                    alignment: Alignment.centerRight,
                    child: TextButton(
                      onPressed: () {
                        // 5. Navigate to the next flow when "Skip" is pressed.
                        context.go('/placeholder');
                      },
                      child: const Text(
                        'Skip',
                        style: TextStyle(
                          color: AppColors.secondaryText,
                          fontSize: 16,
                        ),
                      ),
                    ),
                  ),

                  // --- The Swipeable Pages ---
                  Expanded(
                    child: PageView.builder(
                      controller: _pageController,
                      itemCount: state.pages.length,
                      // 6. Notify the BLoC when the user swipes to a new page.
                      onPageChanged: (index) {
                        context.read<OnboardingBloc>().add(PageSwiped(index));
                      },
                      itemBuilder: (context, index) {
                        final page = state.pages[index];
                        return _OnboardingPageContent(
                          icon: page.icon,
                          title: page.title,
                          subtitle: page.subtitle,
                          description: page.description,
                        );
                      },
                    ),
                  ),
                  const SizedBox(height: 20),

                  // --- Bottom Navigation Row (Indicator + Button) ---
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    crossAxisAlignment: CrossAxisAlignment.center,
                    children: [
                      // --- Page Indicator Dots ---
                      _PageIndicator(
                        itemCount: state.pages.length,
                        currentIndex: state.pageIndex,

                        onDotTapped: (index) {
                          // Use the page controller to smoothly animate to the tapped page.
                          _pageController.animateToPage(
                            index,
                            duration: const Duration(milliseconds: 400),
                            curve: Curves.easeInOut,
                          );
                        },
                      ),

                      // --- Next/Continue Button ---
                      CustomButton(
                        // 7. Change button text based on whether it's the last page.
                        text: state.isLastPage ? 'Continue' : 'Next',
                        icon: Icons.arrow_forward,
                        onTap: () {
                          if (state.isLastPage) {
                            // 8. On the last page, navigate to the next flow.
                            context.go('/placeholder');
                          } else {
                            // 9. Otherwise, animate to the next page.
                            _pageController.nextPage(
                              duration: const Duration(milliseconds: 300),
                              curve: Curves.easeIn,
                            );
                          }
                        },
                      ),
                    ],
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}

class _OnboardingPageContent extends StatelessWidget {
  final IconData icon;
  final String title;
  final String subtitle;
  final String description;
  const _OnboardingPageContent({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.description,
  });
  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        IconContainer(icon: icon),
        const SizedBox(height: 40),
        Text(
          title,
          textAlign: TextAlign.center,
          style: const TextStyle(
            fontSize: 24,
            fontWeight: FontWeight.bold,
            color: AppColors.primaryText,
          ),
        ),
        const SizedBox(height: 12),
        Text(
          subtitle,
          textAlign: TextAlign.center,
          style: const TextStyle(fontSize: 18, color: AppColors.secondaryText),
        ),
        const SizedBox(height: 20),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
          child: Text(
            description,
            textAlign: TextAlign.center,
            style: const TextStyle(
              fontSize: 16,
              color: AppColors.secondaryText,
              height: 1.5,
            ),
          ),
        ),
      ],
    );
  }
}

// --- Helper widget for the page indicator dots ---
// --- Helper widget for the page indicator dots ---
class _PageIndicator extends StatelessWidget {
  final int itemCount;
  final int currentIndex;
  final Function(int) onDotTapped; // ADDED: The callback function

  const _PageIndicator({
    required this.itemCount,
    required this.currentIndex,
    required this.onDotTapped, // ADDED: Make the callback required
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      children: List.generate(
        itemCount,
        (index) => GestureDetector(
          // ADDED: To make the dot tappable
          onTap: () => onDotTapped(index), // ADDED: Call the function on tap
          child: AnimatedContainer(
            duration: const Duration(milliseconds: 300),
            margin: const EdgeInsets.symmetric(horizontal: 4.0),
            height: 10.0,
            width: currentIndex == index ? 24.0 : 10.0,
            decoration: BoxDecoration(
              color: currentIndex == index
                  ? AppColors.darkGreenColor
                  : AppColors.graycolor,
              borderRadius: BorderRadius.circular(12),
            ),
          ),
        ),
      ),
    );
  }
}
