import 'package:ethioguide/core/config/app_theme.dart';
import 'package:ethioguide/features/authentication/presentation/bloc/auth_bloc.dart';
import 'package:ethioguide/features/AI%20chat/Presentation/bloc/ai_bloc.dart';
import 'package:ethioguide/features/procedure/presentation/bloc/procedure_bloc.dart';
import 'package:ethioguide/features/procedure/presentation/bloc/workspace_procedure_detail_bloc.dart';
import 'package:ethioguide/features/workspace_discussion/presentation/bloc/workspace_discussion_bloc.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'core/config/app_router.dart'; 
import 'injection_container.dart' as di;

void main() async {
  WidgetsFlutterBinding.ensureInitialized;
  await di.init();

  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  MyApp({super.key});
  final ThemeMode _themeMode = ThemeMode.light;

  @override
  Widget build(BuildContext context) {
    // 2. Use the MaterialApp.router constructor.
    return MultiBlocProvider(
      providers: [
        BlocProvider(create: (context) => di.sl<AiBloc>()),

        BlocProvider(
          create: (context) =>
              di.sl<WorkspaceDiscussionBloc>()..add(const FetchDiscussions()),
        ),
        BlocProvider(
          create: (context) =>
              di.sl<ProcedureBloc>()..add(LoadProceduresEvent(null)),
        ),
        BlocProvider(
          create: (context) =>
              di.sl<WorkspaceProcedureDetailBloc>()..add(FetchMyProcedures()),
        ),

        BlocProvider(create: (context) => di.sl<AuthBloc>()),
      ],
      child: MaterialApp.router(
        themeMode: _themeMode,
        theme: AppTheme.lightTheme,
        darkTheme: AppTheme.darkTheme,
        routerConfig: router,
        title: 'EthioGuide',
        debugShowCheckedModeBanner: false,
      ),
    );
  }
}