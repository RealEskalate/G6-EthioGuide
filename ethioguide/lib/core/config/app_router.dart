import 'package:ethioguide/core/config/route_names.dart';
import 'package:ethioguide/features/procedure/presentation/pages/workspace_page.dart';
import 'package:ethioguide/features/procedure/presentation/pages/workspace_procedure_detail_page.dart';
import 'package:ethioguide/features/workspace_discussion/presentation/pages/workspace_discussion_page.dart';
import 'package:go_router/go_router.dart';
import 'package:ethioguide/features/splashscreen/presentation/screens/placeholder_screen.dart';
import 'package:ethioguide/features/splashscreen/presentation/screens/splash_screen.dart';

import '../../features/procedure/presentation/pages/procedure_detail_page.dart';
import '../../features/procedure/presentation/pages/procedure_page.dart';

// This is the central router configuration for the entire application.
final GoRouter router = GoRouter(
  initialLocation: '/',

  // The 'routes' list contains all the possible pages the user can navigate to.
  routes: [
    //  This is the route for the splash screen.
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
            final procedureId = state.extra as String;
            return WorkspaceProcedureDetailPage(procedureId: procedureId);
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
    // This is the route for our temporary placeholder screen.
    GoRoute(
      path: '/placeholder',
      name: 'placeholder',
      builder: (context, state) => const PlaceholderScreen(),
    ),
  ],
);
