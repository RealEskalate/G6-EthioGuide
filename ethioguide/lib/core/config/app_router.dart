import 'package:ethioguide/features/authentication/presentation/screens/auth_screen.dart';

import 'package:go_router/go_router.dart';
import 'package:ethioguide/features/onboarding/presentation/screens/onboarding_screen.dart';
import 'package:ethioguide/features/splashscreen/presentation/screens/splash_screen.dart';
import 'package:ethioguide/features/splashscreen/presentation/screens/placeholder_screen.dart'; 
import 'package:ethioguide/core/config/route_names.dart';

// This is the central router configuration for the entire application.
final GoRouter router = GoRouter(
  // The initial route that the app will open on.
  initialLocation: '/',

  // The list of all available routes in the app.
  routes: [
    GoRoute(
      path: '/',
      name: RouteNames.splash, 
      builder: (context, state) => const SplashScreen(),
    ),
    GoRoute(
      path: '/onboarding',
      name: RouteNames.onboarding, // Assuming you have this
      builder: (context, state) => const OnboardingScreen(),
    ),
    // A single, clean route for the entire authentication flow (Login, Sign Up, Forgot Password).
    GoRoute(
      path: '/auth',
      name: RouteNames.auth, 
      builder: (context, state) => const AuthScreen(),
    ),
    GoRoute(
      path: '/placeholder', 
      name: 'placeholder',
      builder: (context, state) => const PlaceholderScreen(),
    ),
    
  ],
);
