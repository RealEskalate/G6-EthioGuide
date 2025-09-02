import 'package:ethioguide/core/config/route_names.dart';
import 'package:ethioguide/features/AI%20chat/Presentation/screens/ai_chat_screen.dart';
import 'package:ethioguide/features/onboarding/presentation/screens/onboarding_screen.dart';
import 'package:go_router/go_router.dart';
import 'package:ethioguide/features/splashscreen/presentation/screens/placeholder_screen.dart';
import 'package:ethioguide/features/splashscreen/presentation/screens/splash_screen.dart';

// This is the central router configuration for the entire application.
final GoRouter router = GoRouter(
  // TODO: return it to root path
  initialLocation: '/aiChat',

  // The 'routes' list contains all the possible pages the user can navigate to.
  routes: [
    //  This is the route for the splash screen.
    GoRoute(
      path: '/',
      name: RouteNames.splash,
      builder: (context, state) =>
          const SplashScreen(), // The function that builds the widget for this screen.
    ),

    // This is the route for our temporary placeholder screen.
    GoRoute(
      path: '/placeholder',
      name: 'placeholder',
      builder: (context, state) => const PlaceholderScreen(),
    ),

     GoRoute(
      path: '/onboarding',
      name: 'onboarding',
      builder: (context, state) => const OnboardingScreen(),
    ),
    
    // This is the route for AI chat page
    GoRoute(
      path: '/aiChat',
      name: RouteNames.aiChat,
      builder: (context, state) => const ChatPage(),
    ),
  
  ],
);
