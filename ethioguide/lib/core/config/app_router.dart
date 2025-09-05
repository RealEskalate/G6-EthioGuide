import 'package:ethioguide/features/authentication/presentation/screens/auth_screen.dart';
import 'package:ethioguide/core/config/route_names.dart';
import 'package:ethioguide/features/AI%20chat/Presentation/screens/ai_chat_screen.dart';
import 'package:ethioguide/features/onboarding/presentation/screens/onboarding_screen.dart';
import 'package:ethioguide/features/profile/presentation/screen/profile_screen.dart';
import 'package:go_router/go_router.dart';
import 'package:ethioguide/features/onboarding/presentation/screens/onboarding_screen.dart';

import 'package:ethioguide/features/splashscreen/presentation/screens/splash_screen.dart';
import 'package:ethioguide/features/splashscreen/presentation/screens/placeholder_screen.dart'; 
import 'package:ethioguide/core/config/route_names.dart';
import 'package:ethioguide/features/home_screen/presentaion/screen/home_screen.dart';


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
    
    // This is the route for AI chat page
    GoRoute(
      path: '/aiChat',
      name: RouteNames.aiChat,
      builder: (context, state) => const ChatPage(),
    ),
     GoRoute(
      path: '/home',
      name: RouteNames.home,
      builder: (context, state) => const HomeScreen(),
    ),
    GoRoute(
      path: '/profile',
      name: 'profile', // You can add this to your RouteNames constants
      builder: (context, state) => const ProfileScreen(),
    ),
  
  ],
);
