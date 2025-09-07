import 'package:ethioguide/features/authentication/presentation/widgets/reset_password_view.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:ethioguide/core/config/route_names.dart';
import 'package:ethioguide/features/procedure/domain/entities/procedure_detail.dart';

// Import all your screens
import 'package:ethioguide/features/splashscreen/presentation/screens/splash_screen.dart';
import 'package:ethioguide/features/onboarding/presentation/screens/onboarding_screen.dart';
import 'package:ethioguide/features/authentication/presentation/screens/auth_screen.dart';
import 'package:ethioguide/features/home_screen/presentaion/screen/home_screen.dart';
import 'package:ethioguide/features/profile/presentation/screen/profile_screen.dart';
import 'package:ethioguide/features/AI%20chat/Presentation/screens/ai_chat_screen.dart';
import 'package:ethioguide/features/procedure/presentation/pages/procedure_page.dart';
import 'package:ethioguide/features/procedure/presentation/pages/procedure_detail_page.dart';
import 'package:ethioguide/features/procedure/presentation/pages/workspace_page.dart';
import 'package:ethioguide/features/procedure/presentation/pages/workspace_procedure_detail_page.dart';
import 'package:ethioguide/features/workspace_discussion/presentation/pages/workspace_discussion_page.dart';
import 'package:ethioguide/features/splashscreen/presentation/screens/placeholder_screen.dart';


import 'package:ethioguide/features/AI%20chat/Presentation/screens/ai_chat_screen.dart';
import 'package:ethioguide/features/onboarding/presentation/screens/onboarding_screen.dart';

import 'package:go_router/go_router.dart';
import 'package:ethioguide/features/onboarding/presentation/screens/onboarding_screen.dart';
import 'package:ethioguide/features/splashscreen/presentation/screens/splash_screen.dart';
import 'package:ethioguide/features/splashscreen/presentation/screens/placeholder_screen.dart';
import 'package:ethioguide/core/config/route_names.dart';

import '../../features/procedure/presentation/pages/procedure_detail_page.dart';
import '../../features/procedure/presentation/pages/procedure_page.dart';

// This is the central router configuration for the entire application.

final GoRouter router = GoRouter(
  initialLocation: '/',

  routes: [
    // --- STANDALONE ROUTES (No persistent UI like a back button to a shell) ---
    GoRoute(
      path: '/',
      name: RouteNames.splash,

      builder: (context, state) =>
          const SplashScreen(), // The function that builds the widget for this screen.
    ),

    GoRoute(
      path: '/Procedure',
      name: RouteNames.procedure,
      builder: (context, state) => const ProcedurePage(),
      routes: [
        GoRoute(
          path: 'detail',
          name: RouteNames.procedure_detail,
          builder: (context, state) => const ProcedureDetailPage(),
        ),
      ], // The function that builds the widget for this screen.
    ),

    GoRoute(
      path: '/workspace',
      name: RouteNames.workspace,
      builder: (context, state) => const WorkspacePage(),
      routes: [
        GoRoute(
          path: 'detail',
          name: RouteNames.workspace_detail,
          builder: (context, state) {
            final procedure = state.extra as ProcedureDetail;
            return WorkspaceProcedureDetailPage(procedureDetail: procedure);
          },
        ),
      ], // The function that builds the widget for this screen.
    ),

    GoRoute(
      path: '/discussion',
      name: RouteNames.workspacediscussion,
      builder: (context, state) =>
          const WorkspaceDiscussionPage(), // The function that builds the widget for this screen.

    ),
    GoRoute(
      path: '/onboarding',
      name: RouteNames.onboarding,
      builder: (context, state) => const OnboardingScreen(),
    ),
    GoRoute(
      path: '/auth',
      name: RouteNames.auth,

      builder: (context, state) =>
          AuthScreen(verificationToken: state.uri.queryParameters['token']),

    ),
    // This deep link for password reset is a special case and should be top-level.
    GoRoute(
      path: '/auth/reset-password',
      name: 'reset_password',
      builder: (context, state) {
        final token = state.uri.queryParameters['token'] ?? '';
        return ResetPasswordView(resetToken: token);
      },
    ),


    GoRoute(
      path: '/home',
      name: RouteNames.home,
      builder: (context, state) => const HomeScreen(),
      routes: [
        // These are now CHILD routes of '/home'.
        // To navigate to them, you use their path RELATIVE to the parent.
        // e.g., context.push('/home/profile') or context.pushNamed('profile')
        GoRoute(
          path: 'profile', // No leading '/'
          name: 'profile',
          builder: (context, state) => const ProfileScreen(),
        ),
        GoRoute(
          path: 'aiChat', // No leading '/'
          name: RouteNames.aiChat,
          builder: (context, state) => const ChatPage(),
        ),
        GoRoute(
          path: 'procedure', // No leading '/'
          name: RouteNames.procedure,
          builder: (context, state) => const ProcedurePage(),
          routes: [
            GoRoute(
              path: 'detail',
              name: RouteNames.procedure_detail,
              builder: (context, state) => const ProcedureDetailPage(),
            ),
          ],
        ),
        GoRoute(
          path: 'workspace', // No leading '/'
          name: RouteNames.workspace,
          builder: (context, state) => const WorkspacePage(),
          routes: [
            GoRoute(
              path: 'detail',
              name: RouteNames.workspace_detail,
              builder: (context, state) {
                final procedureId = state.extra as String;
                return WorkspaceProcedureDetailPage(procedureId: procedureId);
              },
            ),
          ],
        ),
        GoRoute(
          path: 'discussion', // No leading '/'
          name: RouteNames.workspacediscussion,
          builder: (context, state) => const WorkspaceDiscussionPage(),
        ),
        GoRoute(
          path: 'placeholder', // A generic placeholder child route
          name: 'placeholder',
          builder: (context, state) => const PlaceholderScreen(),
        ),
      ],
    ),
  ],
);
