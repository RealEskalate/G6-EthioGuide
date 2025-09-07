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
  initialLocation: '/auth',
  routes: [
    GoRoute(
      path: '/',
      name: RouteNames.splash,
      builder: (context, state) => const SplashScreen(),
    ),

    GoRoute(
      path: '/onboarding',
      name: RouteNames.onboarding,
      builder: (context, state) => const OnboardingScreen(),
    ),

    GoRoute(
      path: '/auth',
      name: RouteNames.auth,
      builder: (context, state) => AuthScreen(
        verificationToken: state.uri.queryParameters['token'],
      ),
      routes: [
        GoRoute(
          path: 'reset-password',
          name: 'reset_password',
          builder: (context, state) {
            final token = state.uri.queryParameters['token'] ?? '';
            return ResetPasswordView(resetToken: token);
          },
        ),
      ],
    ),



    GoRoute(
      path: '/home',
      name: RouteNames.home,
      builder: (context, state) => const HomeScreen(),
      routes: [
        GoRoute(
          path: 'profile',
          name: RouteNames.profile,
          builder: (context, state) => const ProfileScreen(),
        ),
        GoRoute(
          path: 'aiChat',
          name: RouteNames.aiChat,
          builder: (context, state) => const ChatPage(),
        ),
        GoRoute(
          path: 'procedure',
          name: RouteNames.procedure,
          builder: (context, state) => const ProcedurePage(),
          routes: [
            GoRoute(
              path: 'proceduredetail',
              name: RouteNames.procedureDetail,
              builder: (context, state) => const ProcedureDetailPage(),
            ),
          ],
        ),
        GoRoute(
          path: 'workspace',
          name: RouteNames.workspace,
          builder: (context, state) => const WorkspacePage(),
          routes: [
            GoRoute(
              path: 'detail',
              name: RouteNames.workspaceDetail,
              builder: (context, state) {
                final procedure = state.extra as ProcedureDetail;
                return WorkspaceProcedureDetailPage(procedureDetail: procedure);
              },
            ),
          ],
        ),
        GoRoute(
          path: 'discussion',
          name: RouteNames.workspaceDiscussion,
          builder: (context, state) => const WorkspaceDiscussionPage(),
        ),
        GoRoute(
          path: 'placeholder',
          name: 'placeholder',
          builder: (context, state) => const PlaceholderScreen(),
        ),
      ],
    ),
  ],
);
